# Development Environment

This directory contains docker-compose configurations and scripts for local development and testing.

## Vault Development Environment

Start Vault in development mode with test data:

```bash
docker-compose -f docker-compose.vault.yml up -d
```

This will:
- Start Vault in dev mode on http://localhost:8200
- Create test secrets (TLS certs, database credentials, app config)
- Configure AppRole authentication
- Display credentials in the logs

View initialization output:
```bash
docker logs vault-init
```

Access Vault:
- **URL**: http://localhost:8200
- **Root Token**: `dev-root-token`

Stop Vault:
```bash
docker-compose -f docker-compose.vault.yml down
```

## OpenBao Development Environment

Start OpenBao in development mode with test data:

```bash
docker-compose -f docker-compose.openbao.yml up -d
```

This will:
- Start OpenBao in dev mode on http://localhost:8200
- Create test secrets (TLS certs, database credentials, app config)
- Configure AppRole authentication
- Display credentials in the logs

View initialization output:
```bash
docker logs openbao-init
```

Access OpenBao:
- **URL**: http://localhost:8200
- **Root Token**: `dev-root-token`

Stop OpenBao:
```bash
docker-compose -f docker-compose.openbao.yml down
```

## Test Secrets

Both environments create the following test secrets:

### TLS Certificate
- **Path**: `secret/common/tls/example-cert`
- **Fields**: `tlsCrt`, `tlsKey`

### Database Credentials
- **Path**: `secret/database/prod/credentials`
- **Fields**: `username`, `password`

### Application Config
- **Path**: `secret/app/config`
- **Fields**: `apiKey`, `apiSecret`

## Authentication Methods

### Token Authentication
Use the root token for testing:
```bash
export VAULT_TOKEN=dev-root-token
# or
export BAO_TOKEN=dev-root-token
```

### AppRole Authentication
Get credentials from the init container logs:
```bash
docker logs vault-init | grep -A 2 "AppRole Credentials"
# or
docker logs openbao-init | grep -A 2 "AppRole Credentials"
```

## Testing the Secrets Sync Tool

Example configuration for testing with Vault:

```yaml
# testdata/dev-config.yaml
secretStore:
  address: "http://localhost:8200"
  authMethod: "token"
  token: "dev-root-token"
  kvVersion: "v2"
  mountPath: "secret"

secrets:
  - name: "tls-cert"
    path: "common/tls/example-cert"
    refreshInterval: "30s"
    template:
      data:
        tls.crt: '{{ .tlsCrt }}'
        tls.key: '{{ .tlsKey }}'
    files:
      - path: "/tmp/secrets/tls.crt"
        mode: "0644"
      - path: "/tmp/secrets/tls.key"
        mode: "0600"
```

Run the tool:
```bash
CONFIG_FILE=testdata/dev-config.yaml go run ./cmd/secrets-sync
```

## Notes

⚠️ **WARNING**: These are development environments only. Never use dev mode or these tokens in production!

- Dev mode stores data in memory (lost on restart)
- Root tokens have unlimited access
- TLS is disabled
- Audit logging is disabled
