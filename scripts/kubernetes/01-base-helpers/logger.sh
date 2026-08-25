#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Logger -- Colorful timestamped logging for Kubernetes setup scripts
#  Source this file:  source "$(dirname "$0")/../01-base-helpers/logger.sh"
# --------------------------------------------------------------------------

# -- Colors ----------------------------------------------------------------
_CLR_RESET="\033[0m"
_CLR_GREEN="\033[0;32m"
_CLR_YELLOW="\033[0;33m"
_CLR_RED="\033[0;31m"
_CLR_CYAN="\033[0;36m"
_CLR_GRAY="\033[0;90m"

# -- Core log function -----------------------------------------------------
log_message() {
    local message="$1"
    local level="${2:-info}"
    local timestamp
    timestamp="$(date '+%Y-%m-%d %H:%M:%S')"

    case "$level" in
        success) echo -e "${_CLR_GREEN}[OK]${_CLR_RESET}  ${_CLR_GRAY}${timestamp}${_CLR_RESET} - $message" ;;
        warn)    echo -e "${_CLR_YELLOW}[!!]${_CLR_RESET}  ${_CLR_GRAY}${timestamp}${_CLR_RESET} - $message" ;;
        error)   echo -e "${_CLR_RED}[ERR]${_CLR_RESET} ${_CLR_GRAY}${timestamp}${_CLR_RESET} - $message" >&2 ;;
        *)       echo -e "${_CLR_CYAN}[>>]${_CLR_RESET}  ${_CLR_GRAY}${timestamp}${_CLR_RESET} - $message" ;;
    esac
}

# -- Log with hostname + IP ------------------------------------------------
log_msg_ip() {
    local message="$1"
    local level="${2:-info}"
    local hostname_val
    local ip_val
    hostname_val="$(hostname 2>/dev/null || echo 'unknown')"
    ip_val="$(hostname -I 2>/dev/null | awk '{print $1}' || echo '?.?.?.?')"
    log_message "[$hostname_val @ $ip_val] $message" "$level"
}

# -- Assert root -----------------------------------------------------------
assert_root() {
    if [[ $EUID -ne 0 ]]; then
        log_message "This script must be run as root (sudo)." "error"
        exit 1
    fi
}
