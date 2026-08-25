#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 4b -- Join a Worker Node to the Cluster
#  Run as root on each WORKER node:
#    sudo ./join-worker.sh <hostname> <join-command>
#  Example:
#    sudo ./join-worker.sh worker1 "kubeadm join 192.168.0.20:6443 --token abc --discovery-token-ca-cert-hash sha256:xyz"
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

assert_root

NEW_HOSTNAME="${1:-}"
JOIN_COMMAND="${2:-}"

log_message "=== Join Worker to Kubernetes Cluster ===" "info"

# -- Set hostname ----------------------------------------------------------
if [[ -n "$NEW_HOSTNAME" ]]; then
    log_message "Setting hostname to: $NEW_HOSTNAME" "info"
    hostnamectl set-hostname "$NEW_HOSTNAME"
fi

# -- Ensure CRI-O is running -----------------------------------------------
systemctl enable crio.service
systemctl start crio.service

# -- Join the cluster -------------------------------------------------------
if [[ -z "$JOIN_COMMAND" ]]; then
    log_message "No join command provided." "error"
    echo ""
    echo "  Usage: sudo ./join-worker.sh <hostname> \"<join-command>\""
    echo ""
    echo "  Get the join command from the master node:"
    echo "    sudo kubeadm token create --print-join-command"
    echo ""
    exit 1
fi

log_message "Joining cluster..." "info"
eval "$JOIN_COMMAND"
log_message "Worker node joined successfully!" "success"
echo ""
echo "  Verify on the master node:"
echo "    kubectl get nodes"
echo ""
