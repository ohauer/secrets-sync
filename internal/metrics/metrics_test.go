package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRecordFetchSuccess(t *testing.T) {
	RecordFetchSuccess("test-secret", "secret/test")

	count := testutil.ToFloat64(SecretFetchTotal.WithLabelValues("test-secret", "secret/test", "success"))
	if count < 1 {
		t.Errorf("expected count >= 1, got %f", count)
	}
}

func TestRecordFetchError(t *testing.T) {
	RecordFetchError("test-secret", "secret/test", "timeout")

	errorCount := testutil.ToFloat64(SecretFetchErrors.WithLabelValues("test-secret", "secret/test", "timeout"))
	if errorCount < 1 {
		t.Errorf("expected error count >= 1, got %f", errorCount)
	}
}

func TestRecordSyncDuration(t *testing.T) {
	RecordSyncDuration("test-secret", 1.5)
	RecordSyncDuration("test-secret", 2.5)

	// Just verify it doesn't panic
	t.Log("sync duration recorded successfully")
}

func TestSetCircuitBreakerState(t *testing.T) {
	tests := []struct {
		state    string
		expected float64
	}{
		{"closed", 0},
		{"half-open", 1},
		{"open", 2},
	}

	for _, tt := range tests {
		SetCircuitBreakerState("vault-client", tt.state)

		value := testutil.ToFloat64(CircuitBreakerState.WithLabelValues("vault-client"))
		if value != tt.expected {
			t.Errorf("state %s: expected %f, got %f", tt.state, tt.expected, value)
		}
	}
}

func TestSetSecretsConfigured(t *testing.T) {
	SetSecretsConfigured(5)

	value := testutil.ToFloat64(SecretsConfigured)
	if value != 5 {
		t.Errorf("expected 5, got %f", value)
	}
}

func TestSetSecretsSynced(t *testing.T) {
	SetSecretsSynced(3)

	value := testutil.ToFloat64(SecretsSynced)
	if value != 3 {
		t.Errorf("expected 3, got %f", value)
	}
}
