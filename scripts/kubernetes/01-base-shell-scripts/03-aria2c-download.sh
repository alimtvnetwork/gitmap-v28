#!/bin/bash

# source "01-logger.sh"
# source "02-install-apt.sh"

# Constants for aria2c download settings
SPLIT=40
CONNECTIONS=16

# Function to install aria2c if not installed
install_aria2c() {
    # install_apt 
    install_apt_no_msg
}

# Function to download a file using aria2c
# Usage:
#   aria2c_download <url> [<output_path>] [<split>] [<connections>]
#   - <url>: The URL to download the file from.
#   - <output_path> (optional): Directory where the file will be downloaded. Default is the current directory.
#   - <split> (optional): Number of split connections. Default is 40.
#   - <connections> (optional): Number of parallel connections. Default is 16.
#
# Example:
#   aria2c_download "https://example.com/file.tar.gz" "/downloads" 50 20
aria2c_download() {
    local url="$1"
    local output_path="${2:-.}"
    local split="${3:-$SPLIT}"
    local connections="${4:-$CONNECTIONS}"
    
    install_aria2c

    echo
    log_message "Downloading ${url} to ${output_path}"
    aria2c -x "$connections" -s "$split" -d "$output_path" "$url"

    log_message "Download of ${url} completed!"
    echo
}
