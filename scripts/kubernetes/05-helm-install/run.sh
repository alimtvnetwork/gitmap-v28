#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 5 -- Install Helm
#  Downloads and installs Helm package manager.
#  Run as root:  sudo ./run.sh [version]
#  Example:      sudo ./run.sh 3.16.2
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

assert_root

HELM_VERSION="${1:-3.16.2}"

log_message "=== Install Helm v${HELM_VERSION} ===" "info"

# -- Check if already installed --------------------------------------------
if command -v helm &>/dev/null; then
    CURRENT_VER=$(helm version --short 2>/dev/null || echo "unknown")
    log_message "Helm is already installed: $CURRENT_VER" "success"
    exit 0
fi

# -- Download and install ---------------------------------------------------
HELM_FILE="helm-v${HELM_VERSION}-linux-amd64.tar.gz"
HELM_URL="https://get.helm.sh/${HELM_FILE}"
EXTRACT_DIR="/tmp/helm-install"

mkdir -p "$EXTRACT_DIR"

log_message "Downloading Helm v${HELM_VERSION}..." "info"
curl -fsSL "$HELM_URL" -o "$EXTRACT_DIR/$HELM_FILE"
log_message "Download complete." "success"

log_message "Extracting..." "info"
tar xzf "$EXTRACT_DIR/$HELM_FILE" -C "$EXTRACT_DIR"

log_message "Installing to /usr/local/bin..." "info"
mv "$EXTRACT_DIR/linux-amd64/helm" /usr/local/bin/helm
chmod +x /usr/local/bin/helm

# -- Cleanup ---------------------------------------------------------------
rm -rf "$EXTRACT_DIR"

# -- Verify ----------------------------------------------------------------
INSTALLED_VER=$(helm version --short 2>/dev/null)
log_message "Helm installed: $INSTALLED_VER" "success"
