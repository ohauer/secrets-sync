# Go Project Best Practices

## Version Management

### Version Variables
Always define version variables that can be injected at build time:

```go
var (
    Version   = "0.1.0"
    GitCommit = "dev"
    BuildDate = "unknown"
)
```

### Makefile Version Injection
Use ldflags to inject version info:

```makefile
VERSION=0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/app ./cmd/app
```

### Version Display Format
- Release builds: `0.1.0`
- Development builds: `dev-<git-hash>`
- Show build date, Go version, OS/Arch

## CLI Command Structure

### Subcommand Pattern
Always implement these standard subcommands:

1. `help` / `-h` / `--help` - Usage information
2. `version` / `-v` / `--version` - Version info
3. `init` - Generate sample configuration
4. `validate` - Validate configuration without running

### Main Function Pattern
```go
func main() {
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "help", "-h", "--help":
            printHelp()
            os.Exit(0)
        case "version", "-v", "--version":
            printVersion()
            os.Exit(0)
        case "init":
            printInitConfig()
            os.Exit(0)
        case "validate":
            os.Exit(runValidate())
        default:
            printUsage()
            os.Exit(1)
        }
    }
    
    if err := run(); err != nil {
        // Handle error
        os.Exit(1)
    }
}
```

## Configuration Management

### Environment Variable Precedence
Environment variables should always override config file values:

```go
// Load from config file
tlsConfig := &TLSConfig{
    CACert: cfg.SecretStore.TLSCACert,
}

// Override with environment variables
if envCfg.VaultCACert != "" {
    tlsConfig.CACert = envCfg.VaultCACert
}
```

### Configuration Validation
Always validate configuration before use:
- Check required fields
- Validate file paths exist
- Validate paired values (e.g., cert + key)
- Provide clear error messages

### Sample Configuration Generation
The `init` command should generate a fully commented example configuration with:
- All available options
- Explanations for each field
- Common use cases
- Security best practices

## TLS/Certificate Handling

### Certificate Generation Script
Provide a script to generate test certificates:
- CA certificate and key
- Server certificates with SANs
- Client certificates for mTLS
- Proper file permissions (644 for certs, 600 for keys)

### TLS Configuration Fields
Standard TLS configuration should include:
- `tlsCACert` - Custom CA certificate file
- `tlsCAPath` - CA certificate directory
- `tlsSkipVerify` - Skip verification (dev only, warn user)
- `tlsClientCert` - Client certificate (mTLS)
- `tlsClientKey` - Client key (mTLS)

### TLS Validation
- Verify CA cert/path files exist
- Ensure client cert and key are both provided or both empty
- Check file permissions
- Validate certificate files are readable

## Testing

### Test Organization
- Unit tests: `*_test.go` in same package
- Integration tests: Separate test files with build tags
- Use `t.TempDir()` for temporary files
- Clean up resources in defer statements

### Test Certificates
For TLS tests, generate minimal valid certificates:
```go
tmpDir := t.TempDir()
caCert := filepath.Join(tmpDir, "ca.pem")
os.WriteFile(caCert, []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"), 0644)
```

### Race Detection
Always run tests with race detector:
```makefile
test:
	go test -v -race -coverprofile=coverage.out ./...
```

Use `atomic` package for shared counters in tests:
```go
var count int32
atomic.AddInt32(&count, 1)
```

## Docker Compose for Development

### Multiple Service Variants
Provide both HTTP and TLS variants:
```yaml
services:
  vault:        # HTTP on port 8200
  vault-tls:    # HTTPS on port 8201
```

### Health Checks
Use appropriate health check commands:
```yaml
healthcheck:
  test: ["CMD", "wget", "--no-check-certificate", "-q", "-O-", "https://127.0.0.1:8200/v1/sys/health"]
```

### Certificate Mounting
Mount certificates read-only:
```yaml
volumes:
  - ./test-run/certs:/certs:ro
```

## Code Quality

### Trailing Whitespace
Always trim trailing whitespace:
```bash
find . -type f \( -name "*.go" -o -name "*.md" \) -exec sed -i 's/[[:space:]]*$//' {} +
```

### Linting
- Run linter before commits: `make lint`
- Fix all linter issues before merging
- Use `golangci-lint` with standard configuration

### Error Handling
- Always check error returns
- Use `_` to explicitly ignore errors when intentional
- Wrap errors with context: `fmt.Errorf("context: %w", err)`

## Documentation

### README Structure
1. Features (with emojis for visual appeal)
2. Quick Start (minimal steps to get running)
3. Installation
4. Usage (command line examples)
5. Configuration (link to detailed docs)
6. Examples (docker-compose, configs)
7. Troubleshooting (link to guide)

### Separate Documentation Files
- `docs/configuration.md` - Detailed configuration reference
- `docs/environment-variables.md` - All env vars with descriptions
- `docs/troubleshooting.md` - Common issues and solutions

### Configuration Documentation
Always include:
- Required vs optional fields
- Default values
- Examples for common use cases
- Security considerations
- Environment variable equivalents

## Security

### Secrets in Logs
- Never log secret values
- Use structured logging with metadata only
- Log secret names/paths, not contents

### File Permissions
- Default to restrictive permissions (0600)
- Make permissions configurable
- Document security implications

### Container Security
- Use `FROM scratch` for minimal attack surface
- Run as non-root user (UID 65534)
- No shell or package manager in final image
- Static binary compilation

## Project Organization

### Test Environment
Create a `test-run/` directory for testing:
- Add to `.gitignore`
- Include setup scripts
- Provide sample configurations
- Document usage in `test-run/README.md`

### Scripts
Provide helper scripts:
- `setup-*.sh` - Setup development environment
- `run-*.sh` - Run the tool with test config
- `generate-*.sh` - Generate test data/certificates
- `verify-*.sh` - Verify results

All scripts should:
- Be POSIX-compliant
- Include logging with timestamps
- Have clear success/failure messages
- Be executable (`chmod +x`)
