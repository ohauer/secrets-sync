package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewClientWithTLS(t *testing.T) {
	tmpDir := t.TempDir()

	caCert := filepath.Join(tmpDir, "ca.pem")
	_ = os.WriteFile(caCert, []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"), 0644)

	clientCert := filepath.Join(tmpDir, "client.pem")
	_ = os.WriteFile(clientCert, []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"), 0644)

	clientKey := filepath.Join(tmpDir, "client-key.pem")
	_ = os.WriteFile(clientKey, []byte("-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"), 0600)

	caPath := filepath.Join(tmpDir, "certs")
	_ = os.Mkdir(caPath, 0755)

	tests := []struct {
		name      string
		address   string
		tlsConfig *TLSConfig
		wantErr   bool
	}{
		{
			name:      "nil TLS config",
			address:   "https://vault.example.com",
			tlsConfig: nil,
			wantErr:   false,
		},
		{
			name:    "skip verify",
			address: "https://vault.example.com",
			tlsConfig: &TLSConfig{
				SkipVerify: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithTLS(tt.address, tt.tlsConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithTLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClientWithTLS() returned nil client")
			}
		})
	}
}

func TestTLSConfig(t *testing.T) {
	tlsConfig := &TLSConfig{
		CACert:     "/certs/ca.pem",
		CAPath:     "/etc/ssl/certs",
		ClientCert: "/certs/client.pem",
		ClientKey:  "/certs/client-key.pem",
		SkipVerify: false,
	}

	if tlsConfig.CACert != "/certs/ca.pem" {
		t.Errorf("CACert = %v, want /certs/ca.pem", tlsConfig.CACert)
	}
	if tlsConfig.CAPath != "/etc/ssl/certs" {
		t.Errorf("CAPath = %v, want /etc/ssl/certs", tlsConfig.CAPath)
	}
	if tlsConfig.ClientCert != "/certs/client.pem" {
		t.Errorf("ClientCert = %v, want /certs/client.pem", tlsConfig.ClientCert)
	}
	if tlsConfig.ClientKey != "/certs/client-key.pem" {
		t.Errorf("ClientKey = %v, want /certs/client-key.pem", tlsConfig.ClientKey)
	}
	if tlsConfig.SkipVerify != false {
		t.Errorf("SkipVerify = %v, want false", tlsConfig.SkipVerify)
	}
}
