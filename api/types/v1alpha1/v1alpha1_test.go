package v1alpha1

import (
	"reflect"
	"testing"

	"github.com/erizzardi/vault-secrets-operator/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

var vs = VaultSecret{
	TypeMeta: v1.TypeMeta{
		Kind:       "VaultSecret",
		APIVersion: "erizzardi.mine.io/v1alpha1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name:      utils.RandSeq(10, letters),
		Namespace: "namespace",
	},
	Spec: VaultSecretSpec{
		MountPath:  "unit-tests-mount-path",
		SecretPath: "unit-tests-secret-path",
		Data: []VaultSecretSpecData{
			{Name: "test-name1", Value: "test-value1"},
			{Name: "test-name2", Value: "test-value2"},
		},
	},
}

var vsl = VaultSecretList{
	TypeMeta: v1.TypeMeta{
		Kind:       "VaultSecret",
		APIVersion: "erizzardi.mine.io/v1alpha1",
	},
	ListMeta: v1.ListMeta{},
	Items: []VaultSecret{
		vs,
		vs,
	},
}

func TestDeepCopyInto(t *testing.T) {
	out := VaultSecret{}

	vs.DeepCopyInto(&out)

	if !reflect.DeepEqual(vs, out) {
		t.Errorf("The two copies are not identical")
	}
}

func TestDeepCopyObject(t *testing.T) {
	obj := vs.DeepCopyObject()

	if !reflect.DeepEqual(&vs, obj.(*VaultSecret)) {
		t.Error("The two copies are not identical")
	}
}

func TestDeepCopyObjectList(t *testing.T) {
	obj := vsl.DeepCopyObject()

	if !reflect.DeepEqual(&vsl, obj.(*VaultSecretList)) {
		t.Error("The two copies are not identical")
	}
}

func TestSpecEqual(t *testing.T) {
	a := VaultSecretSpec{
		MountPath:  "mp",
		SecretPath: "sp",
		Data: []VaultSecretSpecData{
			{
				Name:  "name",
				Value: "value",
			},
			{
				Name:  "name",
				Value: "value",
			},
		},
	}

	b := a
	if !SpecEqual(a, b) {
		t.Errorf("%v != %v", a, b)
	}
}
