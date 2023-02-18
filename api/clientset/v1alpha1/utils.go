package v1alpha1

import (
	"encoding/json"

	"github.com/erizzardi/vault-secrets-operator/api/types/v1alpha1"
)

func FromDataToSecret(vs *v1alpha1.VaultSecret) (map[string]interface{}, error) {
	secret := make(map[string]interface{})
	for _, m := range vs.Spec.Data {
		if json.Valid([]byte(m.Value)) {
			if err := json.Unmarshal([]byte(m.Value), &secret); err != nil {
				return nil, err
			}
		} else {
			secret[m.Name] = m.Value
		}
	}
	return secret, nil
}
