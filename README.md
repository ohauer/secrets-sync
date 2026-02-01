# Docker Secrets Sidecar

A lightweight sidecar container for managing secrets from HashiCorp Vault or OpenBao in Docker/Podman environments. Continuously syncs secrets to the filesystem with configurable refresh intervals.

## Features

- 🔄 **Continuous Sync** - Automatically refreshes secrets at configurable intervals
- 🔐 **Multiple Auth Methods** - Supports Token and AppRole authentication
- 🔒 **TLS Support** - Custom CA certificates, mTLS, self-signed certificates
- 📝 **Template Engine** - Map secret fields to multiple files (external-secrets-operator style)
- 🛡️ **Circuit Breaker** - Prevents cascading failures with exponential backoff
- 📊 **Observability** - JSON logging, Prometheus metrics, optional OpenTelemetry tracing
- 🔧 **Hot Reload** - Configuration changes without restart
- 🐳 **Minimal Image** - FROM scratch, <20MB, runs as non-root
- ✅ **Health Checks** - Built-in healthcheck for docker-compose and Kubernetes

## Quick Start

### 1. Start Vault (for testing)

```bash
docker-compose -f docker-compose.vault.yml up -d
docker logs vault-init  # Get credentials
```

### 2. Create Configuration

```yaml
# config.yaml
secretStore:
  address: "http://localhost:8200"
  authMethod: "token"
  token: "dev-root-token"
  kvVersion: "v2"
  mountPath: "secret"

secrets:
  - name: "tls-cert"
    path: "common/tls/example-cert"
    refreshInterval: "30m"
    template:
      data:
        tls.crt: '{{ .tlsCrt }}'
        tls.key: '{{ .tlsKey }}'
    files:
      - path: "/secrets/tls.crt"
        mode: "0644"
      - path: "/secrets/tls.key"
        mode: "0600"
```

### 3. Run with Docker Compose

```bash
docker-compose -f examples/docker-compose.sidecar.yml up
```

## Installation

### Build from Source

```bash
make build
```

### Build Docker Image

```bash
make docker-build
# or
./build.sh
```

## Usage

### Command Line

```bash
# Run with config file
CONFIG_FILE=config.yaml ./secrets-sync

# Check readiness
./secrets-sync isready
```

### Docker Compose Sidecar Pattern

```yaml
services:
  app:
    image: your-app:latest
    volumes:
      - secrets:/secrets:ro
    depends_on:
      secrets-sidecar:
        condition: service_healthy

  secrets-sidecar:
    image: docker-secrets:v1.0.0
    volumes:
      - secrets:/secrets
      - ./config.yaml:/config.yaml:ro
    environment:
      VAULT_ADDR: http://vault:8200
      VAULT_TOKEN: ${VAULT_TOKEN}
      CONFIG_FILE: /config.yaml
    healthcheck:
      test: ["/app/secrets-sync", "isready"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  secrets:
```

## Configuration

See [docs/configuration.md](docs/configuration.md) for detailed configuration options.

### Environment Variables

See [docs/environment-variables.md](docs/environment-variables.md) for all available environment variables.

## Authentication

### Token Authentication

```yaml
secretStore:
  authMethod: "token"
  token: "${VAULT_TOKEN}"
```

### AppRole Authentication

```yaml
secretStore:
  authMethod: "approle"
  roleId: "${VAULT_ROLE_ID}"
  secretId: "${VAULT_SECRET_ID}"
```

### TLS with Custom CA (Self-Signed Certificates)

```yaml
secretStore:
  address: "https://vault.example.com"
  authMethod: "token"
  token: "${VAULT_TOKEN}"
  tlsCACert: "/certs/ca-bundle.pem"
```

Or using environment variables:

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=your-token
export VAULT_CACERT=/certs/ca-bundle.pem
./secrets-sync
```

See [docs/configuration.md](docs/configuration.md) for more TLS options (mTLS, CA path, skip verify).

## Observability

### Health Endpoints

- `GET /health` - Always returns 200 (liveness)
- `GET /ready` - Returns 200 when secrets synced (readiness)
- `GET /metrics` - Prometheus metrics

### Metrics

- `secret_fetch_total` - Total fetch attempts
- `secret_fetch_errors_total` - Total fetch errors
- `secret_sync_duration_seconds` - Sync duration histogram
- `circuit_breaker_state` - Circuit breaker state (0=closed, 1=half-open, 2=open)
- `secrets_configured` - Number of configured secrets
- `secrets_synced` - Number of successfully synced secrets

### Tracing

Enable OpenTelemetry tracing:

```bash
ENABLE_TRACING=true
OTEL_EXPORTER_ENDPOINT=http://jaeger:4318
```

## Examples

- [Docker Compose Sidecar](examples/docker-compose.sidecar.yml)
- [Vault Development Environment](docker-compose.vault.yml)
- [OpenBao Development Environment](docker-compose.openbao.yml)
- [Example Configurations](examples/)

## Troubleshooting

See [docs/troubleshooting.md](docs/troubleshooting.md) for common issues and solutions.

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build binary
make build

# Build Docker image
make docker-build
```

## Security

- Runs as non-root user (UID 65534)
- Minimal attack surface (FROM scratch)
- Atomic file writes
- Configurable file permissions
- Secrets never logged
- Circuit breaker prevents overwhelming Vault

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions welcome! Please read the contributing guidelines before submitting PRs.

## Support

- GitHub Issues: Report bugs and request features
- GitLab Issues: For internal deployments

## Acknowledgments

- Inspired by [external-secrets-operator](https://external-secrets.io/)
- Uses [HashiCorp Vault Go SDK](https://github.com/hashicorp/vault/tree/main/api)
- Circuit breaker by [sony/gobreaker](https://github.com/sony/gobreaker)
