#!/bin/sh
set -e

echo "Waiting for OpenBao to be ready..."
sleep 2

echo "Enabling KV v2 secrets engine..."
bao secrets enable -version=2 -path=secret kv || echo "KV engine already enabled"

echo "Creating test secrets..."

# TLS certificate example
bao kv put secret/common/tls/example-cert \
    tlsCrt="-----BEGIN CERTIFICATE-----
MIICljCCAX4CCQCKz8Zr8vJKZDANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJV
UzAeFw0yNjAyMDEwMDAwMDBaFw0yNzAyMDEwMDAwMDBaMA0xCzAJBgNVBAYTAlVT
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z8Z8Z8Z8Z8Z8Z8Z8Z8Z
-----END CERTIFICATE-----" \
    tlsKey="-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRnxnxnxnxnxnx
-----END PRIVATE KEY-----"

# Database credentials example
bao kv put secret/database/prod/credentials \
    username="dbuser" \
    password="dbpass123"

# Application config example
bao kv put secret/app/config \
    apiKey="test-api-key-12345" \
    apiSecret="test-api-secret-67890"

echo "Creating AppRole auth..."
bao auth enable approle || echo "AppRole already enabled"

bao write auth/approle/role/secrets-sync \
    token_ttl=1h \
    token_max_ttl=4h \
    policies=secrets-sync-policy

# Create policy for secrets-sync
bao policy write secrets-sync-policy - <<EOF
path "secret/data/*" {
  capabilities = ["read", "list"]
}
path "secret/metadata/*" {
  capabilities = ["read", "list"]
}
EOF

# Get RoleID and SecretID
ROLE_ID=$(bao read -field=role_id auth/approle/role/secrets-sync/role-id)
SECRET_ID=$(bao write -field=secret_id -f auth/approle/role/secrets-sync/secret-id)

echo ""
echo "=========================================="
echo "OpenBao initialized successfully!"
echo "=========================================="
echo "OpenBao Address: http://localhost:8200"
echo "Root Token: dev-root-token"
echo ""
echo "AppRole Credentials:"
echo "  Role ID:   $ROLE_ID"
echo "  Secret ID: $SECRET_ID"
echo ""
echo "Test secrets created:"
echo "  - secret/common/tls/example-cert (tlsCrt, tlsKey)"
echo "  - secret/database/prod/credentials (username, password)"
echo "  - secret/app/config (apiKey, apiSecret)"
echo "=========================================="
