package vault

import (
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
)

// SecretData represents the data retrieved from Vault
type SecretData map[string]interface{}

// FetchSecret fetches a secret from Vault KV v2
func (c *Client) FetchSecret(mountPath, secretPath string) (SecretData, error) {
	fullPath := path.Join(mountPath, "data", secretPath)

	result, err := c.executeWithBreaker(func() (interface{}, error) {
		return c.client.Logical().Read(fullPath)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("secret not found at path: %s", secretPath)
	}

	secret, ok := result.(*api.Secret)
	if !ok || secret == nil {
		return nil, fmt.Errorf("invalid secret response")
	}

	if secret.Data == nil {
		return nil, fmt.Errorf("secret has no data")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret data format")
	}

	return SecretData(data), nil
}
