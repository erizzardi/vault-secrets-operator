package v1alpha1

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
)

type VaultSecretInterface interface {
	List(opts metav1.ListOptions, ctx context.Context) (*v1alpha1.VaultSecretList, error)
	Get(name string, options metav1.GetOptions, ctx context.Context) (*v1alpha1.VaultSecret, error)
	Create(vs *v1alpha1.VaultSecret, ctx context.Context) (*v1alpha1.VaultSecret, error)
	Watch(opts metav1.ListOptions, ctx context.Context) (watch.Interface, error)
	Delete(name string, opts metav1.ListOptions, ctx context.Context) error
	Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, ctx context.Context, subresources ...string) error
}

type vaultSecretClient struct {
	restClient rest.Interface
	ns         string
}

func (c *vaultSecretClient) List(opts metav1.ListOptions, ctx context.Context) (*v1alpha1.VaultSecretList, error) {
	result := v1alpha1.VaultSecretList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("vaultsecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *vaultSecretClient) Get(name string, opts metav1.GetOptions, ctx context.Context) (*v1alpha1.VaultSecret, error) {
	result := v1alpha1.VaultSecret{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("vaultsecrets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *vaultSecretClient) Create(vs *v1alpha1.VaultSecret, ctx context.Context) (*v1alpha1.VaultSecret, error) {
	result := v1alpha1.VaultSecret{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("vaultsecrets").
		Body(vs).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *vaultSecretClient) Watch(opts metav1.ListOptions, ctx context.Context) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("vaultsecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *vaultSecretClient) Delete(name string, opts metav1.ListOptions, ctx context.Context) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *vaultSecretClient) Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, ctx context.Context, subresources ...string) error {
	result := &v1alpha1.VaultSecret{}
	err := c.restClient.
		Patch(pt).
		Namespace(c.ns).
		Resource("vaultsecrets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return err
}
