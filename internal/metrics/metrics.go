package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// SecretFetchTotal tracks total secret fetch attempts
	SecretFetchTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "secret_fetch_total",
			Help: "Total number of secret fetch attempts",
		},
		[]string{"secret_name", "vault_path", "status"},
	)

	// SecretFetchErrors tracks secret fetch errors
	SecretFetchErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "secret_fetch_errors_total",
			Help: "Total number of secret fetch errors",
		},
		[]string{"secret_name", "vault_path", "error_type"},
	)

	// SecretSyncDuration tracks secret sync duration
	SecretSyncDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "secret_sync_duration_seconds",
			Help:    "Duration of secret sync operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"secret_name"},
	)

	// CircuitBreakerState tracks circuit breaker state
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
		},
		[]string{"name"},
	)

	// SecretsConfigured tracks number of configured secrets
	SecretsConfigured = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "secrets_configured",
			Help: "Number of configured secrets",
		},
	)

	// SecretsSynced tracks number of successfully synced secrets
	SecretsSynced = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "secrets_synced",
			Help: "Number of successfully synced secrets",
		},
	)
)

// RecordFetchSuccess records a successful secret fetch
func RecordFetchSuccess(secretName, vaultPath string) {
	SecretFetchTotal.WithLabelValues(secretName, vaultPath, "success").Inc()
}

// RecordFetchError records a failed secret fetch
func RecordFetchError(secretName, vaultPath, errorType string) {
	SecretFetchTotal.WithLabelValues(secretName, vaultPath, "error").Inc()
	SecretFetchErrors.WithLabelValues(secretName, vaultPath, errorType).Inc()
}

// RecordSyncDuration records the duration of a sync operation
func RecordSyncDuration(secretName string, duration float64) {
	SecretSyncDuration.WithLabelValues(secretName).Observe(duration)
}

// SetCircuitBreakerState sets the circuit breaker state
func SetCircuitBreakerState(name, state string) {
	var value float64
	switch state {
	case "closed":
		value = 0
	case "half-open":
		value = 1
	case "open":
		value = 2
	}
	CircuitBreakerState.WithLabelValues(name).Set(value)
}

// SetSecretsConfigured sets the number of configured secrets
func SetSecretsConfigured(count int) {
	SecretsConfigured.Set(float64(count))
}

// SetSecretsSynced sets the number of successfully synced secrets
func SetSecretsSynced(count int) {
	SecretsSynced.Set(float64(count))
}
