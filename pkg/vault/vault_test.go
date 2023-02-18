package vault

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

var ctx = context.TODO()
var mountPath = "unit-tests-mount-path"
var secretPath = "unit-tests-secret-path"
var inputData = map[string]interface{}{
	"bacon": "delicious",
	"steak": "true",
}
var vaultUrl = os.Getenv("VAULT_URL")
var vaultToken = os.Getenv("VAULT_TOKEN")
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func TestWriteDelete(t *testing.T) {
	vClient, err := getVaultClient()
	if err != nil {
		t.Error(err.Error())
	}
	writtenSecret, err := WriteSecret(mountPath, secretPath, inputData, vClient, ctx)
	if err != nil {
		t.Error(err.Error())
	}
	if err := DeleteSecretVersioned(mountPath, secretPath, writtenSecret.VersionMetadata.Version, vClient, ctx); err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetDeleteLastVersion(t *testing.T) {
	vClient, err := getVaultClient()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = WriteSecret(mountPath, secretPath, inputData, vClient, ctx)
	if err != nil {
		t.Error(err.Error())
	}
	secret, err := GetLastVersionSecret(mountPath, secretPath, vClient, ctx)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(secret.Data, inputData) {
		t.Errorf("%v != %v", secret.Data, inputData)
	}

	if err := DeleteLastVersionSecret(mountPath, secretPath, vClient, ctx); err != nil {
		t.Error(err.Error())
	}
}

func TestGetNonExistingSecret(t *testing.T) {
	vClient, err := getVaultClient()
	if err != nil {
		t.Error(err.Error())
	}
	if _, err := GetSecretVersioned(mountPath, secretPath, 999, vClient, ctx); err == nil {
		t.Error("Succeeded (!) in getting non existing secret.")
	}
	if _, err := GetSecretVersioned("non-existing", secretPath, 1, vClient, ctx); err == nil {
		t.Error("Succeeded (!) in getting non existing secret.")
	}
	if _, err := GetSecretVersioned(mountPath, "non-existing", 1, vClient, ctx); err == nil {
		t.Error("Succeeded (!) in getting non existing secret.")
	}
}

// ==================
// Helper functions
// ==================
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
