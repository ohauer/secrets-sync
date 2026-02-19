# Troubleshooting

## Common Issues

### Secret Not Found

**Symptom**: Error message "secret not found at path"

**Causes**:
- Incorrect secret path in configuration
- Secret doesn't exist in Vault
- Wrong KV mount path

**Solutions**:
1. Verify secret exists in Vault:
   ```bash
   vault kv get secret/path/to/secret
   ```

2. Check mount path configuration:
   ```yaml
   secretStore:
     mountPath: "secret"  # Default KV v2 mount
   ```

3. Ensure path doesn't include mount path:
   ```yaml
   secrets:
     - path: "common/tls/cert"  # Correct
     # NOT: "secret/data/common/tls/cert"
   ```

### Authentication Failed

**Symptom**: Error message "authentication failed"

**Token Auth**:
1. Verify token is valid:
   ```bash
   vault token lookup
   ```

2. Check token has required permissions:
   ```bash
   vault token capabilities secret/data/path/to/secret
   ```

**AppRole Auth**:
1. Verify role exists:
   ```bash
   vault read auth/approle/role/secrets-sync
   ```

2. Check role_id and secret_id are correct

3. Verify policy allows reading secrets:
   ```hcl
   path "secret/data/*" {
     capabilities = ["read", "list"]
   }
   ```

### Circuit Breaker Open

**Symptom**: Error message "circuit breaker: circuit breaker is open"

**Causes**:
- Too many consecutive failures
- Vault is unreachable
- Network issues

**Solutions**:
1. Check Vault connectivity:
   ```bash
   curl -k $VAULT_ADDR/v1/sys/health
   ```

2. Review logs for underlying errors

3. Wait for circuit breaker to transition to half-open (default: 30s)

4. Adjust circuit breaker settings if needed:
   ```bash
   CIRCUIT_BREAKER_TIMEOUT=60s
   CIRCUIT_BREAKER_MAX_REQUESTS=5
   ```

### File Permission Denied

**Symptom**: Error message "failed to write file: permission denied"

**Causes**:
- Container user doesn't have write permissions
- Parent directory doesn't exist
- SELinux/AppArmor restrictions

**Solutions**:
1. Ensure volume has correct permissions:
   ```yaml
   volumes:
     - secrets:/secrets  # Not :ro
   ```

2. Check container user can write:
   ```bash
   docker exec container-name touch /secrets/test
   ```

3. Create parent directories in advance

4. Check SELinux context (if applicable):
   ```bash
   ls -Z /path/to/secrets
   ```

### Template Rendering Failed

**Symptom**: Error message "failed to render template"

**Causes**:
- Invalid template syntax
- Missing field in secret data
- Typo in field name

**Solutions**:
1. Verify template syntax:
   ```yaml
   template:
     data:
       key: '{{ .fieldName }}'  # Correct
       # NOT: '{{ fieldName }}'
   ```

2. Check secret contains expected fields:
   ```bash
   vault kv get -format=json secret/path | jq .data.data
   ```

3. Use exact field names (case-sensitive)

### Configuration Not Reloading

**Symptom**: Changes to config file not applied

**Causes**:
- `WATCH_CONFIG` not enabled
- File watcher not detecting changes
- Invalid configuration (validation failed)

**Solutions**:
1. Enable config watching:
   ```bash
   WATCH_CONFIG=true
   ```

2. Check logs for validation errors

3. Restart container if hot reload not working

### Health Check Failing

**Symptom**: Container marked as unhealthy

**Causes**:
- No secrets synced yet
- All secrets failing to sync
- Status file not writable

**Solutions**:
1. Check if any secrets synced successfully:
   ```bash
   curl http://localhost:8080/ready
   ```

2. Review logs for sync errors

3. Verify status file path is writable:
   ```bash
   STATUS_FILE=/tmp/.ready-state
   ```

4. Increase healthcheck retries/interval:
   ```yaml
   healthcheck:
     retries: 10
     interval: 15s
   ```

### High Memory Usage

**Symptom**: Container using excessive memory

**Causes**:
- Too many secrets configured
- Very large secret values
- Memory leak (report bug)

**Solutions**:
1. Reduce number of concurrent secrets

2. Increase refresh intervals:
   ```yaml
   refreshInterval: "1h"  # Instead of "1m"
   ```

3. Monitor metrics:
   ```bash
   curl http://localhost:8080/metrics | grep memory
   ```

### Secrets Not Updating

**Symptom**: Old secret values remain after Vault update

**Causes**:
- Refresh interval not reached
- Sync errors (check logs)
- Circuit breaker open

**Solutions**:
1. Check refresh interval:
   ```yaml
   refreshInterval: "5m"  # Adjust as needed
   ```

2. Force restart to sync immediately:
   ```bash
   docker restart secrets-sidecar
   ```

3. Review logs for errors

## Debugging

### Enable Debug Logging

```bash
LOG_LEVEL=debug
```

### Check Metrics

```bash
curl http://localhost:8080/metrics
```

Key metrics:
- `secret_fetch_total` - Total fetch attempts
- `secret_fetch_errors_total` - Fetch errors
- `circuit_breaker_state` - Circuit breaker state
- `secrets_synced` - Successfully synced secrets

### View Logs

Docker:
```bash
docker logs secrets-sidecar
```

Docker Compose:
```bash
docker-compose logs -f secrets-sidecar
```

### Test Configuration

```bash
# Validate config loads
CONFIG_FILE=config.yaml ./secrets-sync --help

# Test with debug logging
LOG_LEVEL=debug CONFIG_FILE=config.yaml ./secrets-sync
```

### Check File Permissions

```bash
# Inside container
ls -la /secrets/

# Check file mode
stat -c '%a %n' /secrets/*
```

## Getting Help

If you're still experiencing issues:

1. Check [GitHub Issues](https://github.com/ohauer/secrets-sync/issues)
2. Review [Configuration Documentation](configuration.md)
3. Enable debug logging and collect logs
4. Include configuration (redact secrets) when reporting issues
