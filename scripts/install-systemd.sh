#!/bin/sh
# Install secrets-sync as a systemd service

set -e

SCRIPT_NAME="install-systemd.sh"
BINARY_SRC="bin/secrets-sync"
BINARY_DEST="/usr/local/bin/secrets-sync"
UNIT_FILE_SRC="examples/systemd/secrets-sync.service"
UNIT_FILE_DEST="/etc/systemd/system/secrets-sync.service"
ENV_FILE_SRC="examples/systemd/secrets-sync.env.example"
ENV_FILE_DEST="/etc/default/secrets-sync"
CONFIG_DIR="/etc/secrets-sync"

log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') ${SCRIPT_NAME} - $1"
}

check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        log_message "ERROR: This script must be run as root"
        exit 1
    fi
}

check_binary() {
    if [ ! -f "${BINARY_SRC}" ]; then
        log_message "ERROR: Binary not found at ${BINARY_SRC}"
        log_message "Run 'make build' first"
        exit 1
    fi
}

install_binary() {
    log_message "Installing binary to ${BINARY_DEST}"
    cp "${BINARY_SRC}" "${BINARY_DEST}"
    chmod 755 "${BINARY_DEST}"
    log_message "Binary installed successfully"
}

install_unit_file() {
    log_message "Installing systemd unit file to ${UNIT_FILE_DEST}"
    cp "${UNIT_FILE_SRC}" "${UNIT_FILE_DEST}"
    chmod 644 "${UNIT_FILE_DEST}"
    log_message "Unit file installed successfully"
}

create_config_dir() {
    if [ ! -d "${CONFIG_DIR}" ]; then
        log_message "Creating config directory ${CONFIG_DIR}"
        mkdir -p "${CONFIG_DIR}"
        chmod 755 "${CONFIG_DIR}"
    fi
}

install_env_file() {
    if [ ! -f "${ENV_FILE_DEST}" ]; then
        log_message "Installing environment file to ${ENV_FILE_DEST}"
        cp "${ENV_FILE_SRC}" "${ENV_FILE_DEST}"
        chmod 600 "${ENV_FILE_DEST}"
        log_message "Environment file installed (edit ${ENV_FILE_DEST} to configure)"
    else
        log_message "Environment file already exists at ${ENV_FILE_DEST}, skipping"
    fi
}

generate_config() {
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        log_message "Generating sample config at ${CONFIG_DIR}/config.yaml"
        "${BINARY_DEST}" init > "${CONFIG_DIR}/config.yaml"
        chmod 600 "${CONFIG_DIR}/config.yaml"
        log_message "Sample config generated (edit ${CONFIG_DIR}/config.yaml to configure)"
    else
        log_message "Config file already exists at ${CONFIG_DIR}/config.yaml, skipping"
    fi
}

reload_systemd() {
    log_message "Reloading systemd daemon"
    systemctl daemon-reload
}

enable_service() {
    log_message "Enabling secrets-sync service"
    systemctl enable secrets-sync
    log_message "Service enabled (will start on boot)"
}

main() {
    log_message "Starting secrets-sync systemd installation"
    
    check_root
    check_binary
    install_binary
    create_config_dir
    install_unit_file
    install_env_file
    generate_config
    reload_systemd
    enable_service
    
    log_message "Installation complete!"
    log_message ""
    log_message "Next steps:"
    log_message "  1. Edit configuration: ${CONFIG_DIR}/config.yaml"
    log_message "  2. Edit environment:   ${ENV_FILE_DEST}"
    log_message "  3. Start service:      systemctl start secrets-sync"
    log_message "  4. Check status:       systemctl status secrets-sync"
    log_message "  5. View logs:          journalctl -u secrets-sync -f"
}

main "$@"
