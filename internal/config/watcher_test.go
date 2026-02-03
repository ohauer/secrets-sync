package config

import (
	"os"
	"testing"
	"time"
)

func TestWatcher_DetectsChanges(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	initialConfig := `secretStore:
  address: "https://vault.example.com"
  authMethod: "token"
  token: "test-token"

secrets:
  - name: "test-secret"
    key: "test/path"
    mountPath: "secret"
    kvVersion: "v2"
    refreshInterval: "5m"
    template:
      data:
        key: '{{ .value }}'
    files:
      - path: "/test/key"
        mode: "0600"
`
	if _, err := tmpFile.WriteString(initialConfig); err != nil {
		t.Fatalf("failed to write initial config: %v", err)
	}
	_ = tmpFile.Close()

	changeDetected := make(chan *Config, 1)
	onChange := func(cfg *Config) error {
		changeDetected <- cfg
		return nil
	}
	onError := func(err error) {}

	watcher, err := NewWatcher(tmpFile.Name(), onChange, onError)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	watcher.Start()

	time.Sleep(100 * time.Millisecond)

	updatedConfig := `secretStore:
  address: "https://vault.example.com"
  authMethod: "token"
  token: "updated-token"

secrets:
  - name: "updated-secret"
    key: "test/path"
    mountPath: "secret"
    kvVersion: "v2"
    refreshInterval: "10m"
    template:
      data:
        key: '{{ .value }}'
    files:
      - path: "/test/key"
        mode: "0600"
`
	if err := os.WriteFile(tmpFile.Name(), []byte(updatedConfig), 0644); err != nil {
		t.Fatalf("failed to update config: %v", err)
	}

	select {
	case cfg := <-changeDetected:
		if cfg.Secrets[0].Name != "updated-secret" {
			t.Errorf("expected secret name 'updated-secret', got: %s", cfg.Secrets[0].Name)
		}
		if cfg.Secrets[0].RefreshInterval != 10*time.Minute {
			t.Errorf("expected refresh interval 10m, got: %v", cfg.Secrets[0].RefreshInterval)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for config change detection")
	}
}

func TestWatcher_InvalidConfigIgnored(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	validConfig := `secretStore:
  address: "https://vault.example.com"
  authMethod: "token"
  token: "test-token"

secrets:
  - name: "test-secret"
    key: "test/path"
    mountPath: "secret"
    kvVersion: "v2"
    refreshInterval: "5m"
    template:
      data:
        key: '{{ .value }}'
    files:
      - path: "/test/key"
        mode: "0600"
`
	if _, err := tmpFile.WriteString(validConfig); err != nil {
		t.Fatalf("failed to write initial config: %v", err)
	}
	_ = tmpFile.Close()

	changeDetected := make(chan *Config, 1)
	onChange := func(cfg *Config) error {
		changeDetected <- cfg
		return nil
	}
	onError := func(err error) {}

	watcher, err := NewWatcher(tmpFile.Name(), onChange, onError)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Stop()

	watcher.Start()

	time.Sleep(100 * time.Millisecond)

	invalidConfig := `secretStore:
  address: ""
secrets: []
`
	if err := os.WriteFile(tmpFile.Name(), []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("failed to update config: %v", err)
	}

	select {
	case <-changeDetected:
		t.Fatal("invalid config should not trigger onChange callback")
	case <-time.After(500 * time.Millisecond):
	}
}
