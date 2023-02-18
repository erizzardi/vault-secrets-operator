package main

import (
	"context"
	"os"
	"reflect"
	"testing"

	clientset_v1alpha1 "github.com/erizzardi/vault-secrets-operator/api/clientset/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig = os.Getenv("KUBECONFIG")
var namespace = os.Getenv("NAMESPACE")
var vaultUrl = os.Getenv("VAULT_URL")
var vaultToken = os.Getenv("VAULT_TOKEN")
var logger = logrus.New()
var ctx = context.TODO()
var vs = v1alpha1.VaultSecret{
	TypeMeta: v1.TypeMeta{
		Kind:       "VaultSecret",
		APIVersion: "erizzardi.mine.io/v1alpha1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name:      "name",
		Namespace: namespace,
	},
	Spec: v1alpha1.VaultSecretSpec{
		MountPath:  "unit-tests-mount-path",
		SecretPath: "unit-tests-secret-path",
		Data: []v1alpha1.VaultSecretSpecData{
			{Name: "test-name1", Value: "test-value1"},
			{Name: "test-name2", Value: "test-value2"},
		},
	},
}

func TestWriteDelete(t *testing.T) {
	clientSet, err := getv1Alpha1ClientSet(kubeconfig)
	if err != nil {
		t.Errorf(err.Error())
	}
	vClient, err := getVaultClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	vs, version, err := writeFunc(&vs, clientSet, vClient, logger, ctx)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Check if secret has been created on vault
	vsSecret, err := clientset_v1alpha1.FromDataToSecret(vs)
	if err != nil {
		t.Errorf(err.Error())
	}
	secretFromVault, err := vault.GetSecretVersioned(vs.Spec.MountPath, vs.Spec.SecretPath, version, vClient, ctx)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(vsSecret, secretFromVault.Data) {
		t.Errorf("Secret not written on vault properly")
	}
	// Delete secret
	if err := deleteFunc(vs, vClient, logger, ctx); err == nil {
		t.Errorf(err.Error())
	}

}

// ==================
// Helper functions
// ==================
func getv1Alpha1ClientSet(kubeconfig string) (clientset_v1alpha1.V1Alpha1Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	v1AlphaClientSet, err := clientset_v1alpha1.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	v1alpha1.AddToScheme(scheme.Scheme)
	return v1AlphaClientSet, nil
}

func getVaultClient() (*api.Client, error) {
	vaultClient, err := api.NewClient(&api.Config{
		Address:    vaultUrl,
		HttpClient: httpClient,
	})
	if err != nil {
		return nil, err
	}
	vaultClient.SetToken(vaultToken)
	return vaultClient, nil
}
