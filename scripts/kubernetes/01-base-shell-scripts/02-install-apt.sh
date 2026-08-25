#!/bin/bash

# Source the log function from the external file
# source "01-logger.sh"

# Function to install packages using apt or apt-get
# Usage:
#   install_apt <package1> <package2> ...
# 
#   This function takes one or more package names as arguments and installs 
#   them using apt or apt-get if they are not already installed on the system.
#
# Parameters:
#   - <package1>, <package2>, ...: Names of the packages to install.
# 
# Example:
#   install_apt curl git
#
#   The above command will install `curl` and `git` packages if they are not 
#   already installed on the system.
#
# Note:
#   - This function will log messages for each package installation step 
#     using the `log_message` function from the sourced file.
#   - It performs a check to see if the package is already installed before 
#     attempting to install it.
install_apt() {
    local packages=("$@")  # Array of package names passed as arguments
    local ip=$(hostname -I | awk '{print $1}')
    for package in "${packages[@]}"; do
        if ! dpkg -l | grep -q "$package"; then  # Check if package is not installed
            log_message "$package could not be found. Installing $package on $ip..."
            sudo apt update -y  # Update package list
            sudo apt install -y "$package"  # Install the package
        else
            log_message "$package is already installed on $ip."
        fi
    done
}

# Function to install packages using apt or apt-get without logging messages
# 
# Md Alim Ul Karim -> https://alimkarim.com | https://riseup-pro.com
#
# Usage:
#   install_apt_no_msg <package1> <package2> ...
# 
#   This function takes one or more package names as arguments and installs 
#   them using apt or apt-get if they are not already installed on the system.
#
# Parameters:
#   - <package1>, <package2>, ...: Names of the packages to install.
# 
# Example:
#   install_apt_no_msg curl git
#
#   The above command will install `curl` and `git` packages if they are not 
#   already installed on the system.
#
# Note:
#   - This function will not log any messages for each package installation step.
#   - It performs a check to see if the package is already installed before 
#     attempting to install it.
install_apt_no_msg() {
    local packages=("$@")  # Array of package names passed as arguments
    for package in "${packages[@]}"; do
        if ! dpkg -l | grep -q "$package"; then  # Check if package is not installed
            sudo apt update -y  # Update package list
            sudo apt install -y "$package"  # Install the package
        fi
    done
}
