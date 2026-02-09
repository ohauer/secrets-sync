#!/bin/sh
# Uninstall secrets-sync systemd service

set -e

SCRIPT_NAME="uninstall-systemd.sh"
BINARY_DEST="/usr/local/bin/secrets-sync"
UNIT_FILE_DEST="/etc/systemd/system/secrets-sync.service"
ENV_FILE_DEST="/etc/default/secrets-sync"
CONFIG_DIR="/etc/secrets-sync"
MAN_PAGE_DEST="/usr/share/man/man1/secrets-sync.1.gz"
DOC_DIR="/usr/share/doc/secrets-sync"

log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') ${SCRIPT_NAME} - $1"
}

check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        log_message "ERROR: This script must be run as root"
        exit 1
    fi
}

stop_service() {
    if systemctl is-active --quiet secrets-sync; then
        log_message "Stopping secrets-sync service"
        systemctl stop secrets-sync
    fi
}

disable_service() {
    if systemctl is-enabled --quiet secrets-sync 2>/dev/null; then
        log_message "Disabling secrets-sync service"
        systemctl disable secrets-sync
    fi
}

remove_unit_file() {
    if [ -f "${UNIT_FILE_DEST}" ]; then
        log_message "Removing unit file ${UNIT_FILE_DEST}"
        rm -f "${UNIT_FILE_DEST}"
    fi
}

remove_binary() {
    if [ -f "${BINARY_DEST}" ]; then
        log_message "Removing binary ${BINARY_DEST}"
        rm -f "${BINARY_DEST}"
    fi
}

remove_env_file() {
    if [ -f "${ENV_FILE_DEST}" ]; then
        log_message "Removing environment file ${ENV_FILE_DEST}"
        rm -f "${ENV_FILE_DEST}"
    fi
}

remove_config() {
    if [ -d "${CONFIG_DIR}" ]; then
        printf "Remove config directory %s? [y/N] " "${CONFIG_DIR}"
        read -r response
        case "${response}" in
            [yY][eE][sS]|[yY])
                log_message "Removing config directory ${CONFIG_DIR}"
                rm -rf "${CONFIG_DIR}"
                ;;
            *)
                log_message "Keeping config directory ${CONFIG_DIR}"
                ;;
        esac
    fi
}

reload_systemd() {
    log_message "Reloading systemd daemon"
    systemctl daemon-reload
}

remove_man_page() {
    if [ -f "${MAN_PAGE_DEST}" ]; then
        log_message "Removing man page ${MAN_PAGE_DEST}"
        rm -f "${MAN_PAGE_DEST}"
    fi
}

remove_documentation() {
    if [ -d "${DOC_DIR}" ]; then
        log_message "Removing documentation ${DOC_DIR}"
        rm -rf "${DOC_DIR}"
    fi
}

main() {
    log_message "Starting secrets-sync systemd uninstallation"

    check_root
    stop_service
    disable_service
    remove_unit_file
    remove_binary
    remove_env_file
    remove_config
    remove_man_page
    remove_documentation
    reload_systemd

    log_message "Uninstallation complete!"
}

main "$@"
