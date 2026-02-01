package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/sony/gobreaker"
)

// TLSConfig holds TLS configuration for Vault client
type TLSConfig struct {
	CACert     string
	CAPath     string
	ClientCert string
	ClientKey  string
	SkipVerify bool
}

// Client wraps the Vault API client
type Client struct {
	client  *api.Client
	breaker *gobreaker.CircuitBreaker
}

// NewClient creates a new Vault client
func NewClient(address string) (*Client, error) {
	return NewClientWithTLS(address, nil)
}

// NewClientWithTLS creates a new Vault client with TLS configuration
func NewClientWithTLS(address string, tlsConfig *TLSConfig) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = address

	if tlsConfig != nil {
		if err := configureTLS(config, tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	return &Client{client: client}, nil
}

func configureTLS(config *api.Config, tlsConfig *TLSConfig) error {
	tlsClientConfig := &api.TLSConfig{
		Insecure: tlsConfig.SkipVerify,
	}

	if tlsConfig.CACert != "" {
		tlsClientConfig.CACert = tlsConfig.CACert
	}

	if tlsConfig.CAPath != "" {
		tlsClientConfig.CAPath = tlsConfig.CAPath
	}

	if tlsConfig.ClientCert != "" {
		tlsClientConfig.ClientCert = tlsConfig.ClientCert
	}

	if tlsConfig.ClientKey != "" {
		tlsClientConfig.ClientKey = tlsConfig.ClientKey
	}

	return config.ConfigureTLS(tlsClientConfig)
}

// GetAPIClient returns the underlying Vault API client
func (c *Client) GetAPIClient() *api.Client {
	return c.client
}

// Ping checks if the Vault server is reachable
func (c *Client) Ping() error {
	_, err := c.executeWithBreaker(func() (interface{}, error) {
		return c.client.Sys().Health()
	})
	if err != nil {
		return fmt.Errorf("vault health check failed: %w", err)
	}
	return nil
}
