/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	clientset_v1alpha1 "github.com/erizzardi/vault-secrets-operator/api/clientset/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/pkg/config"
	"github.com/erizzardi/vault-secrets-operator/pkg/vault"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {

	// ==============================================================
	// Read and validate configuration from flags and env variables
	//===============================================================
	cfg, err := config.GetConfigOrDie()
	if err != nil {
		log.Panic(err)
	}

	// Define context
	ctx := context.Background()

	logger := log.New()
	logger.Level = cfg.LogLevel

	if err := cfg.ValidateConfiguration(); err != nil {
		logger.Panicf("Panic: invalid configuration: %s", err.Error())
	}

	// =====================
	// Vault configuration
	// =====================
	vaultClient, err := api.NewClient(&api.Config{
		Address:    cfg.VaultUrl,
		HttpClient: httpClient,
		// Backoff: func(min time.Duration, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// },
		// CheckRetry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// },

	})
	if err != nil {
		panic(err)
	}
	// Add other authentication methods
	vaultClient.SetToken(cfg.VaultToken)

	// ==========================
	// Kubernetes configuration
	// ==========================
	var k8sCfg *rest.Config
	if cfg.LocalTesting {
		k8sCfg, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
		if err != nil {
			logger.Panicf("Panic: cannot build Kubernetes configuration: %s", err.Error())
		}
	} else {
		k8sCfg, err = rest.InClusterConfig()
		if err != nil {
			log.Panic(err.Error())
		}
	}

	v1AlphaClientSet, err := clientset_v1alpha1.NewForConfig(k8sCfg)
	if err != nil {
		panic(err)
	}

	// Register type definition
	v1alpha1.AddToScheme(scheme.Scheme)

	logger.Info("Controller started")

	// ===========================
	// Define and run controller
	// ===========================
	stopCh := make(chan struct{})
	bc := make(chan error)
	_, controller := vaultSecretsController(v1AlphaClientSet, "", cfg.ResyncPeriod, bc, logger, vaultClient, ctx)
	// Start controller
	go controller.Run(stopCh)

	for {
		// This handles all the non-handed errors, i.e. the critical ones
		// ATM there are none, but improve in the future
		if <-bc != nil {
			close(stopCh)
			break
		}
		time.Sleep(time.Duration(cfg.LoopPeriod) * time.Second)
	}
}

// watchResources leverages the Informer type to poll the k8s api and get the status of the VaultSecrets resources
func vaultSecretsController(clientSet clientset_v1alpha1.V1Alpha1Interface, namespace string, resyncPeriod int, bc chan error, logger *log.Logger, client *api.Client, ctx context.Context) (cache.Store, cache.Controller) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				res, err := clientSet.VaultSecrets(namespace).List(lo, ctx)
				if err != nil {
					logger.Errorf("Error unmarshaling object: %s", err.Error())
					return &v1alpha1.VaultSecretList{}, err
				}
				return res, nil
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				res, err := clientSet.VaultSecrets(namespace).Watch(lo, ctx)
				if err != nil {
					logger.Errorf("Error creating Watch for resource: %s", err.Error())
					return res, err
				}
				return res, nil
			},
		},
		&v1alpha1.VaultSecret{},
		time.Duration(resyncPeriod)*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				vs := obj.(*v1alpha1.VaultSecret)
				if _, ok := vs.Annotations[config.ManagedAnnotation]; !ok {
					if err := addFunc(vs, clientSet, client, logger, ctx); err != nil {
						logger.Errorf("Error at object adding: %s", err.Error())
					}
				} else {
					logger.Debugf("Skipped adding object %s", vs.Name)
				}
				bc <- nil
			},
			DeleteFunc: func(obj interface{}) {
				if err := deleteFunc(obj.(*v1alpha1.VaultSecret), client, logger, ctx); err != nil {
					logger.Errorf("Error at object deletion: %s", err.Error())
				}
				bc <- nil
			},
			UpdateFunc: func(obj interface{}, newObj interface{}) {
				if err := updateFunc(obj.(*v1alpha1.VaultSecret), clientSet, client, logger, ctx); err != nil {
					logger.Errorf("Error at object update: %s", err.Error())
				}
				bc <- nil
			},
		},
	)
	return store, controller
}

// addFunc = writeFunc + patchFunc
func addFunc(vs *v1alpha1.VaultSecret, clientSet clientset_v1alpha1.V1Alpha1Interface, client *api.Client, logger *log.Logger, ctx context.Context) error {
	vs, version, err := writeFunc(vs, clientSet, client, logger, ctx)
	if err != nil {
		return err
	}
	return patchFunc(vs, clientSet, version, client, logger, ctx)
}

func writeFunc(vs *v1alpha1.VaultSecret, clientSet clientset_v1alpha1.V1Alpha1Interface, client *api.Client, logger *log.Logger, ctx context.Context) (*v1alpha1.VaultSecret, int, error) {
	logger.Debugf("Write secret %s into Vault", vs.Name)
	// Build JSON from Data
	secret, err := clientset_v1alpha1.FromDataToSecret(vs)
	if err != nil {
		return nil, -1, errors.New("cannot unmarshal secret data into structure: " + err.Error())
	}
	// There's no error from WriteSecret that can stop the controller
	out, err := vault.WriteSecret(vs.Spec.MountPath, vs.Spec.SecretPath, secret, client, ctx)
	if err != nil {
		return nil, -1, errors.New("cannot write secret to Vault: " + err.Error())
	}
	return vs, out.VersionMetadata.Version, nil
}

func patchFunc(vs *v1alpha1.VaultSecret, clientSet clientset_v1alpha1.V1Alpha1Interface, version int, client *api.Client, logger *log.Logger, ctx context.Context) error {
	// Patch object with last-applied-configuration annotation - like kubectl, and secret version
	// Remove kubectl last-applied-configuration annotation, if present
	delete(vs.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	delete(vs.Annotations, config.LACAnnotation)

	// Marshal last applied configuration, and encode it in base64
	lacByte, err := json.Marshal(vs)
	if err != nil {
		return err
	}
	lacb64 := base64.StdEncoding.EncodeToString(lacByte)
	patchString := fmt.Sprintf("{\"metadata\":{\"annotations\":{\"%s\":\"%s\",\"%s\":\"true\",\"%s\":\"%d\"}}}", config.LACAnnotation, lacb64, config.ManagedAnnotation, config.SecretVersionAnnotation, version)
	logger.Debugf("Patching object %s with patch %s", vs.Name, patchString)
	if err := clientSet.VaultSecrets(vs.Namespace).Patch(vs.Name, types.MergePatchType, []byte(patchString), metav1.PatchOptions{}, ctx); err != nil {
		return err
	}
	logger.Infof("Written secret %s/%s version %d into Vault", vs.Spec.MountPath, vs.Spec.SecretPath, version)
	return nil
}

func deleteFunc(vs *v1alpha1.VaultSecret, client *api.Client, logger *log.Logger, ctx context.Context) error {
	version, err := strconv.Atoi(vs.Annotations[config.SecretVersionAnnotation])
	if err != nil {
		return err
	}
	if err := vault.DeleteSecretVersioned(vs.Spec.MountPath, vs.Spec.SecretPath, version, client, ctx); err != nil {
		return err
	}
	logger.Infof("Deleted secret %s/%d/%s version %s from Vault", vs.Spec.MountPath, version, vs.Spec.SecretPath, vs.Name)
	return nil
}

func updateFunc(vs *v1alpha1.VaultSecret, clientSet clientset_v1alpha1.V1Alpha1Interface, client *api.Client, logger *log.Logger, ctx context.Context) error {
	lac := &v1alpha1.VaultSecret{}
	// Read and decode last applied configuration from object annotations
	lacByte, err := base64.StdEncoding.DecodeString(vs.Annotations[config.LACAnnotation])
	if err != nil {
		return err
	}
	if string(lacByte) == "" {
		return nil
	}
	err = json.Unmarshal(lacByte, lac)
	if err != nil {
		return err
	}
	// If the new version has a different secret payload, then write the new secret version
	if !v1alpha1.SpecEqual(lac.Spec, vs.Spec) {
		return addFunc(vs, clientSet, client, logger, ctx)
	} else {
		logger.Debugf("Skipping update of object %s", vs.Name)
	}
	return nil
}

// // REMEMBER THIS for testing
// source := fcache.NewFakeControllerSource()
