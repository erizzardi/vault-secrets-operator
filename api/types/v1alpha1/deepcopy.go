package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (vs *VaultSecret) DeepCopyInto(out *VaultSecret) {
	out.TypeMeta = vs.TypeMeta
	out.ObjectMeta = vs.ObjectMeta
	out.Spec = VaultSecretSpec{
		MountPath:  vs.Spec.MountPath,
		SecretPath: vs.Spec.SecretPath,
		Data:       vs.Spec.Data,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (vs *VaultSecret) DeepCopyObject() runtime.Object {
	out := VaultSecret{}
	vs.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (vs *VaultSecretList) DeepCopyObject() runtime.Object {
	out := VaultSecretList{}
	out.TypeMeta = vs.TypeMeta
	out.ListMeta = vs.ListMeta

	if vs.Items != nil {
		out.Items = make([]VaultSecret, len(vs.Items))
		for i := range vs.Items {
			vs.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
