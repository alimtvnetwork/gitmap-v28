#!/bin/bash

# Function to check if a package is installed
# Usage:
#   is_package_installed <package_name>
#
# This function checks if the specified package is installed on the system.
# If the package is installed, it logs a message and returns 0.
# If the package is not installed, it logs a message and returns 1.
#
# Example:
#   if is_package_installed helm; then
#     log_message "Helm is installed."
#     helm version
#   else
#     log_message "Helm is not installed."
#   fi
#
# Note:
#   - This function will log messages for package status using the `log_message` function from the sourced file.
is_package_installed() {
    local ip=$(hostname -I | awk '{print $1}')
    local package="$1"
    if command -v "$package" &> /dev/null; then
        log_message "$package is already installed on $ip."
        return 0
    else
        log_message "$package is not installed on $ip."
        return 1
    fi
}
