package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate_TLSConfig(t *testing.T) {
	tmpDir := t.TempDir()

	caCert := filepath.Join(tmpDir, "ca.pem")
	_ = os.WriteFile(caCert, []byte("test"), 0644)

	clientCert := filepath.Join(tmpDir, "client.pem")
	_ = os.WriteFile(clientCert, []byte("test"), 0644)

	clientKey := filepath.Join(tmpDir, "client-key.pem")
	_ = os.WriteFile(clientKey, []byte("test"), 0600)

	caPath := filepath.Join(tmpDir, "certs")
	_ = os.Mkdir(caPath, 0755)

	tests := []struct {
		name    string
		store   SecretStore
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid CA cert",
			store: SecretStore{
				Address:    "https://vault.example.com",
				AuthMethod: "token",
				Token:      "test",
				TLSCACert:  caCert,
			},
			wantErr: false,
		},
		{
			name: "invalid CA cert path",
			store: SecretStore{
				Address:    "https://vault.example.com",
				AuthMethod: "token",
				Token:      "test",
				TLSCACert:  "/nonexistent/ca.pem",
			},
			wantErr: true,
			errMsg:  "tlsCACert file does not exist",
		},
		{
			name: "valid CA path",
			store: SecretStore{
				Address:    "https://vault.example.com",
				AuthMethod: "token",
				Token:      "test",
				TLSCAPath:  caPath,
			},
			wantErr: false,
		},
		{
			name: "invalid CA path",
			store: SecretStore{
				Address:    "https://vault.example.com",
				AuthMethod: "token",
				Token:      "test",
				TLSCAPath:  "/nonexistent/certs",
			},
			wantErr: true,
			errMsg:  "tlsCAPath directory does not exist",
		},
		{
			name: "CA path is file not directory",
			store: SecretStore{
				Address:    "https://vault.example.com",
				AuthMethod: "token",
				Token:      "test",
				TLSCAPath:  caCert,
			},
			wantErr: true,
			errMsg:  "tlsCAPath is not a directory",
		},
		{
			name: "valid client cert and key",
			store: SecretStore{
				Address:       "https://vault.example.com",
				AuthMethod:    "token",
				Token:         "test",
				TLSClientCert: clientCert,
				TLSClientKey:  clientKey,
			},
			wantErr: false,
		},
		{
			name: "client cert without key",
			store: SecretStore{
				Address:       "https://vault.example.com",
				AuthMethod:    "token",
				Token:         "test",
				TLSClientCert: clientCert,
			},
			wantErr: true,
			errMsg:  "tlsClientKey is required when tlsClientCert is set",
		},
		{
			name: "client key without cert",
			store: SecretStore{
				Address:      "https://vault.example.com",
				AuthMethod:   "token",
				Token:        "test",
				TLSClientKey: clientKey,
			},
			wantErr: true,
			errMsg:  "tlsClientCert is required when tlsClientKey is set",
		},
		{
			name: "invalid client cert path",
			store: SecretStore{
				Address:       "https://vault.example.com",
				AuthMethod:    "token",
				Token:         "test",
				TLSClientCert: "/nonexistent/client.pem",
				TLSClientKey:  clientKey,
			},
			wantErr: true,
			errMsg:  "tlsClientCert file does not exist",
		},
		{
			name: "invalid client key path",
			store: SecretStore{
				Address:       "https://vault.example.com",
				AuthMethod:    "token",
				Token:         "test",
				TLSClientCert: clientCert,
				TLSClientKey:  "/nonexistent/client-key.pem",
			},
			wantErr: true,
			errMsg:  "tlsClientKey file does not exist",
		},
		{
			name: "skip verify",
			store: SecretStore{
				Address:       "https://vault.example.com",
				AuthMethod:    "token",
				Token:         "test",
				TLSSkipVerify: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				SecretStore: tt.store,
				Secrets: []Secret{
					{
						Name:            "test",
						Key:             "test/path",
						MountPath:       "secret",
						KVVersion:       "v2",
						RefreshInterval: 60000000000,
						Template: Template{
							Data: map[string]string{"key": "{{ .value }}"},
						},
						Files: []File{
							{Path: "/test", Mode: "0600"},
						},
					},
				},
			}

			err := Validate(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
