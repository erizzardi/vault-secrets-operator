package v1alpha1

import (
	"context"
	"os"
	"testing"

	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
	"github.com/erizzardi/vault-secrets-operator/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

// These tests require access to a working k8s cluster, with the VaultSecret CRD installed, and a local Vault instance. Specify the kubeconfig with the KUBECONFIG env var
var kubeconfig = os.Getenv("KUBECONFIG")
var namespace = os.Getenv("NAMESPACE")

var ctx = context.TODO()
var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func TestSpecDifferent(t *testing.T) {
	a := v1alpha1.VaultSecretSpec{
		MountPath:  "mp",
		SecretPath: "sp",
		Data: []v1alpha1.VaultSecretSpecData{{
			Name:  "name",
			Value: "value",
		}},
	}

	b := a
	b.MountPath = "something-different"
	if v1alpha1.SpecEqual(a, b) {
		t.Errorf("%v != %v", a, b)
	}
}

func TestCreateDeleteVS(t *testing.T) {
	v1AlphaClientSet, err := getv1Alpha1ClientSet(kubeconfig)
	if err != nil {
		t.Errorf(err.Error())
	}
	// Create the vaultsecret
	vs := v1alpha1.VaultSecret{
		TypeMeta: v1.TypeMeta{
			Kind:       "VaultSecret",
			APIVersion: "erizzardi.mine.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      utils.RandSeq(10, letters),
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
	res, err := v1AlphaClientSet.VaultSecrets(namespace).Create(&vs, ctx)
	if err != nil {
		t.Error(err.Error())
	}
	// Check if vs and res Specs are the same
	if !v1alpha1.SpecEqual(vs.Spec, res.Spec) {
		t.Errorf("%v != %v", vs.Spec, res.Spec)
	}

	// Delete the vaultsecret
	if err = v1AlphaClientSet.VaultSecrets(namespace).Delete(vs.Name, v1.DeleteOptions{}, ctx); err != nil {
		t.Error(err.Error())
	}
}

func TestDeleNonExistingVS(t *testing.T) {
	v1AlphaClientSet, err := getv1Alpha1ClientSet(kubeconfig)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Delete non existing vaultsecret
	if err = v1AlphaClientSet.VaultSecrets(namespace).Delete("non-existing", v1.DeleteOptions{}, ctx); err == nil {
		t.Error("Deleted a non existing resource!")
	}
}

func TestWatcher(t *testing.T) {
	v1AlphaClientSet, err := getv1Alpha1ClientSet(kubeconfig)
	if err != nil {
		t.Errorf(err.Error())
	}
	// Create the vaultsecret
	vs := v1alpha1.VaultSecret{
		TypeMeta: v1.TypeMeta{
			Kind:       "VaultSecret",
			APIVersion: "erizzardi.mine.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      utils.RandSeq(10, letters),
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
	if _, err = v1AlphaClientSet.VaultSecrets(namespace).Create(&vs, ctx); err != nil {
		t.Error(err.Error())
	}

	if _, err = v1AlphaClientSet.VaultSecrets(namespace).Watch(v1.ListOptions{}, ctx); err != nil {
		t.Error(err.Error())
	}
}

// ==================
// Helper functions
// ==================
func getv1Alpha1ClientSet(kubeconfig string) (V1Alpha1Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	v1AlphaClientSet, err := NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	v1alpha1.AddToScheme(scheme.Scheme)
	return v1AlphaClientSet, nil
}
