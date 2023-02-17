package v1alpha1

import (
	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// This interface contains all the types clients
type V1Alpha1Interface interface {
	VaultSecrets(namespace string) VaultSecretInterface
	// For the future...
}

type v1Alpha1Client struct {
	restClient rest.Interface
	ns         string
}

func (c *v1Alpha1Client) VaultSecrets(namespace string) VaultSecretInterface {
	return &vaultSecretClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func NewForConfig(c *rest.Config) (V1Alpha1Interface, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &v1Alpha1Client{restClient: client}, nil
}
