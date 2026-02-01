package main

import (
	"fmt"
	"os"
)

func printHelp() {
	fmt.Print(`secrets-sync - Sync secrets from Vault/OpenBao to filesystem

USAGE:
    secrets-sync [command] [flags]

COMMANDS:
    (none)      Run the secrets sync service (default)
    init        Generate example configuration file
    validate    Validate configuration file
    version     Show version information
    isready     Check if service is ready (for healthchecks)
    help        Show this help message

FLAGS:
    -h, --help  Show this help message

ENVIRONMENT VARIABLES:
    CONFIG_FILE              Path to configuration file (default: /config.yaml)
    VAULT_ADDR              Vault/OpenBao server address
    VAULT_TOKEN             Vault token for authentication
    VAULT_ROLE_ID           AppRole role ID
    VAULT_SECRET_ID         AppRole secret ID
    VAULT_CACERT            Path to CA certificate file
    VAULT_CAPATH            Path to CA certificate directory
    VAULT_SKIP_VERIFY       Skip TLS verification (insecure)
    VAULT_CLIENT_CERT       Path to client certificate (mTLS)
    VAULT_CLIENT_KEY        Path to client key (mTLS)
    LOG_LEVEL               Log level (debug, info, warn, error)
    HTTP_PORT               Health server port (default: 8080)
    WATCH_CONFIG            Enable config hot reload (default: false)

EXAMPLES:
    # Run with config file
    CONFIG_FILE=config.yaml secrets-sync

    # Generate example config
    secrets-sync init > config.yaml

    # Validate config
    secrets-sync validate

    # Check version
    secrets-sync version

    # Healthcheck
    secrets-sync isready

For more information, see: https://github.com/ohauer/docker-secrets
`)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: secrets-sync [command]\n")
	fmt.Fprintf(os.Stderr, "Run 'secrets-sync help' for more information.\n")
}
