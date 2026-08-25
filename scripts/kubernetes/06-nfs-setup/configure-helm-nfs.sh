#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 6b -- Configure Helm NFS Provisioner
#  Sets up dynamic PV provisioning via NFS.
#  Run on MASTER after Helm + NFS server are installed.
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

NFS_PATH="${1:-/nfsexport}"
NFS_SERVER="${2:-$(hostname -I | awk '{print $1}')}"

log_message "=== Configure Helm NFS Provisioner ===" "info"
log_message "NFS Server: $NFS_SERVER, Path: $NFS_PATH" "info"

# -- Add Helm repo ---------------------------------------------------------
helm repo add nfs-subdir-external-provisioner \
    https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/ 2>/dev/null || true

helm repo update

# -- Install/upgrade provisioner -------------------------------------------
helm upgrade --install nfs-provisioner \
    nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
    --set nfs.server="$NFS_SERVER" \
    --set nfs.path="$NFS_PATH"

log_message "NFS provisioner deployed." "success"
echo ""
echo "  Verify:"
echo "    kubectl get pods"
echo "    kubectl get storageclass"
echo ""
