package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

func WriteSecret(mountPath string, secretPath string, inputData map[string]interface{}, client *api.Client, ctx context.Context) (*api.KVSecret, error) {

	// fmt.Println(inputData)

	res, err := client.KVv2(mountPath).Put(ctx, secretPath, inputData)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteSecret(mountPath string, secretPath string, client *api.Client, ctx context.Context) error {
	return client.KVv2(mountPath).Delete(ctx, secretPath)
}

func DeleteSecretVersioned(mountPath string, secretPath string, version int, client *api.Client, ctx context.Context) error {
	return client.KVv2(mountPath).DeleteVersions(ctx, secretPath, []int{version})
}
