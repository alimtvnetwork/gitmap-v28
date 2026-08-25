#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  APT installer helper -- idempotent package installation
#  Source this file:  source "$(dirname "$0")/../01-base-helpers/install-apt.sh"
# --------------------------------------------------------------------------

# Requires logger.sh to be sourced first

# -- Install packages via apt (with logging) --------------------------------
install_apt() {
    local packages=("$@")
    for package in "${packages[@]}"; do
        if dpkg -l "$package" 2>/dev/null | grep -q "^ii"; then
            log_message "$package is already installed." "info"
        else
            log_message "Installing $package ..." "info"
            sudo apt-get update -qq -y
            sudo apt-get install -qq -y "$package"
            if [[ $? -eq 0 ]]; then
                log_message "$package installed successfully." "success"
            else
                log_message "Failed to install $package." "error"
            fi
        fi
    done
}

# -- Install silently (no log per package) ----------------------------------
install_apt_quiet() {
    local packages=("$@")
    for package in "${packages[@]}"; do
        if ! dpkg -l "$package" 2>/dev/null | grep -q "^ii"; then
            sudo apt-get install -qq -y "$package" >/dev/null 2>&1
        fi
    done
}

# -- Check if a package is installed ----------------------------------------
is_package_installed() {
    dpkg -l "$1" 2>/dev/null | grep -q "^ii"
}
