#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 6a -- Install NFS Server on the Master Node
#  Provides persistent volume storage for the cluster.
#  Run as root on MASTER:  sudo ./install-nfs-server.sh
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

assert_root

NFS_EXPORT="/nfsexport"

log_message "=== Install NFS Server ===" "info"

# -- Install NFS server package --------------------------------------------
install_apt nfs-kernel-server

# -- Create export directory -----------------------------------------------
mkdir -p "$NFS_EXPORT"

# -- Add export rule (idempotent) ------------------------------------------
if ! grep -q "$NFS_EXPORT" /etc/exports 2>/dev/null; then
    echo "$NFS_EXPORT *(rw,no_root_squash,no_subtree_check)" >> /etc/exports
    log_message "Added NFS export: $NFS_EXPORT" "success"
else
    log_message "NFS export already configured." "info"
fi

# -- Restart NFS -----------------------------------------------------------
systemctl restart nfs-kernel-server
log_message "NFS server started." "success"

# -- Show exports ----------------------------------------------------------
CONTROL_IP=$(hostname -I | awk '{print $1}')
showmount -e "$CONTROL_IP"
log_message "NFS server ready at $CONTROL_IP:$NFS_EXPORT" "success"
