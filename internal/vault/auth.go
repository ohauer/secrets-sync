package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// AuthMethod represents the authentication method
type AuthMethod string

const (
	AuthMethodToken   AuthMethod = "token"
	AuthMethodAppRole AuthMethod = "approle"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Method   AuthMethod
	Token    string
	RoleID   string
	SecretID string
}

// Authenticate authenticates the client with Vault
func (c *Client) Authenticate(config AuthConfig) error {
	switch config.Method {
	case AuthMethodToken:
		return c.authenticateToken(config.Token)
	case AuthMethodAppRole:
		return c.authenticateAppRole(config.RoleID, config.SecretID)
	default:
		return fmt.Errorf("unsupported auth method: %s", config.Method)
	}
}

func (c *Client) authenticateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	c.client.SetToken(token)

	_, err := c.executeWithBreaker(func() (interface{}, error) {
		return c.client.Auth().Token().LookupSelf()
	})
	if err != nil {
		return fmt.Errorf("token authentication failed: %w", err)
	}

	return nil
}

func (c *Client) authenticateAppRole(roleID, secretID string) error {
	if roleID == "" {
		return fmt.Errorf("roleId is required")
	}
	if secretID == "" {
		return fmt.Errorf("secretId is required")
	}

	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	result, err := c.executeWithBreaker(func() (interface{}, error) {
		return c.client.Logical().Write("auth/approle/login", data)
	})
	if err != nil {
		return fmt.Errorf("approle authentication failed: %w", err)
	}

	resp, ok := result.(*api.Secret)
	if !ok || resp == nil || resp.Auth == nil {
		return fmt.Errorf("approle authentication returned no token")
	}

	c.client.SetToken(resp.Auth.ClientToken)
	return nil
}
