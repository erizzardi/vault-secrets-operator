package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

func WriteSecret(secretEngine string, secretPath string, inputData map[string]interface{}, client *api.Client, ctx context.Context) (*api.KVSecret, error) {

	// fmt.Println(inputData)

	res, err := client.KVv2(secretEngine).Put(ctx, secretPath, inputData)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteSecret(secretEngine string, secretPath string, client *api.Client, ctx context.Context) error {

	err := client.KVv2(secretEngine).Delete(ctx, secretPath)
	if err != nil {
		return err
	}
	return nil
}
