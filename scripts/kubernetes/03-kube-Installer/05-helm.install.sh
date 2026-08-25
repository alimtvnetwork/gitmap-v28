#!/bin/bash

# Source the log function from the external file
# source "$PWD/../base-shell-scripts/01-logger.sh"
# source "$PWD/../base-shell-scripts/02-install-apt.sh"
# source "$PWD/../base-shell-scripts/03-aria2c-download.sh"
# source "$PWD/../base-shell-scripts/04-is-package-installed.sh"
source "$PWD/../01-base-shell-scripts/00-import-all.sh"

# Function to download Helm based on specified version
download_helm() {
    local ip=$(hostname -I | awk '{print $1}')
    local version="${1:-3.16.2}"  # Default to 3.16.2 if no version is specified
    local helm_file="helm-v${version}-linux-amd64.tar.gz"
    local helm_url="https://get.helm.sh/${helm_file}"
    local extractPath="/temp/helm-install"
    local helm_bin="/usr/local/bin"
    local helm_dir="linux-amd64/helm"

    log_msg_ip "Checking for existing ${helm_file} in ${extractPath}."
    if [ -f "${extractPath}/${helm_file}" ]; then
        log_msg_ip "${helm_file} already exists, skipping download."
    else
        log_msg_ip "Downloading ${helm_file} at $PWD."

        # Create a temporary directory for the Helm installation
        mkdir -p "$extractPath"

        # Change to the temporary directory
        cd "$extractPath"

        # Download the specified version of Helm using aria2c
        aria2c_download "${helm_url}" "${extractPath}"

        log_msg_ip "Download of ${helm_file} completed!"
    fi

    log_msg_ip "Extracting at $extractPath - xvf ${helm_file}"

    # Unpack the Helm file
    tar xvf "${helm_file}"

    log_msg_ip "Moving ${helm_dir} to ${helm_bin}"
    # Move the Helm binary to the /usr/local/bin directory
    sudo mv "${helm_dir}" "${helm_bin}"

    # Go back to the previous directory
    cd -

    log_msg_ip "Removing download directory $extractPath"
    # Remove the temporary directory and its contents
    rm -rf "$extractPath"

    log_msg_ip "Helm v${version} installation completed successfully!"
}

# Check if Helm is already installed
if is_package_installed helm; then
    log_msg_ip "Helm version:"
    helm version
else
    # Call the function to download and install the specified or default version of Helm
    download_helm "${1:-3.16.2}"
    
    # Verify the installation
    log_msg_ip "Helm version:"
    helm version
fi
