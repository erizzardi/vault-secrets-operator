/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	clientset_v1alpha1 "github.com/erizzardi/vault-secrets-operator/api/clientset/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/pkg/config"
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
		panic(err)
	}

	// Define context
	ctx := context.Background()

	logger := log.New()
	logger.Level = cfg.LogLevel

	if err := cfg.ValidateConfiguration(); err != nil {
		logger.Panicf("Panic: invalid configuration: %s", err.Error())
	}

	// Vault configuration
	client, err := api.NewClient(&api.Config{Address: cfg.VaultUrl, HttpClient: httpClient})
	if err != nil {
		panic(err)
	}
	client.SetToken(cfg.VaultToken)

	// Init Kubernetes configuration
	// k8sCfg, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Panic(err.Error())
	// }
	k8sCfg, err := clientcmd.BuildConfigFromFlags("", "/Users/erizzardi/.kube/config")
	if err != nil {
		logger.Panicf("Panic: cannot build Kubernetes configuration: %s", err.Error())
	}

	v1AlphaClientSet, err := clientset_v1alpha1.NewForConfig(k8sCfg)
	if err != nil {
		panic(err)
	}

	// Register type definition
	v1alpha1.AddToScheme(scheme.Scheme)

	stopCh := make(chan struct{})
	bc := make(chan bool)
	_, controller := vaultSecretsController(v1AlphaClientSet, "", bc, logger, ctx)
	// Start controller
	go controller.Run(stopCh)

	for {
		if !<-bc {
			close(stopCh)
			break
		}
		time.Sleep(1 * time.Second)
	}

	// count := 0
	// for {
	// 	time.Sleep(1 * time.Second)
	// 	if count == 5 {
	// 		<-stopCh
	// 		break
	// 	}
	// 	count++
	// }

	// // Reconciliation loop
	// prevIterVsList := v1alpha1.VaultSecretList{}
	// for {
	// 	// List all VaultSecrets
	// 	vsList := store.List()

	// 	// If there are new resources

	// 	for k, vs := range vsList {
	// 		secret := make(map[string]interface{})
	// 		// mountPath := vs.(*v1alpha1.VaultSecret).Spec.MountPath
	// 		// secretPath := vs.(*v1alpha1.VaultSecret).Spec.SecretPath
	// 		data := vs.(*v1alpha1.VaultSecret).Spec.Data

	// 		// Build JSON from Data
	// 		for _, m := range data {
	// 			if json.Valid([]byte(m.Value)) {
	// 				if err := json.Unmarshal([]byte(m.Value), &secret); err != nil {
	// 					panic(err)
	// 				}
	// 			} else {
	// 				secret[m.Name] = m.Value
	// 			}
	// 		}
	// 		// Create new secret version in Vault
	// 		// Copy current list in older iteration
	// 		vs.(*v1alpha1.VaultSecret).DeepCopyInto(&prevIterVsList.Items[k])
	// 	}
	// }

}

// watchResources leverages the Informer type to poll the k8s api and get the status of the VaultSecrets resources
func vaultSecretsController(clientSet clientset_v1alpha1.V1Alpha1Interface, namespace string, bc chan bool, logger *log.Logger, ctx context.Context) (cache.Store, cache.Controller) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.VaultSecrets(namespace).List(lo, ctx)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.VaultSecrets(namespace).Watch(lo, ctx)
			},
		},
		&v1alpha1.VaultSecret{},
		60*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				logger.Debugf("Write secret %s into Vault", obj.(*v1alpha1.VaultSecret).Name)
				bc <- true
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("Object deleted: %s\n", obj.(*v1alpha1.VaultSecret).Name)
				bc <- false
			},
			UpdateFunc: func(obj interface{}, newObj interface{}) {
				fmt.Printf("Object updated: %s\n", obj.(*v1alpha1.VaultSecret).Name)
			},
		},
	)
	return store, controller
}

// // REMEMBER THIS for testing
// source := fcache.NewFakeControllerSource()
