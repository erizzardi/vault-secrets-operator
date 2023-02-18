package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type VaultSecretSpecData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type VaultSecretSpec struct {
	MountPath  string                `json:"secretEngine"`
	SecretPath string                `json:"secretPath"`
	Data       []VaultSecretSpecData `json:"data"`
}

type VaultSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VaultSecretSpec `json:"spec"`
}

type VaultSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []VaultSecret `json:"items"`
}

func dataEqual(a, b []VaultSecretSpecData) bool {
	for key, da := range a {
		if equal := da.Name == b[key].Name && da.Value == b[key].Value; !equal {
			return equal
		}
	}
	return true
}

func SpecEqual(a, b VaultSecretSpec) bool {
	return a.MountPath == b.MountPath &&
		a.SecretPath == b.SecretPath &&
		dataEqual(a.Data, b.Data)
}
