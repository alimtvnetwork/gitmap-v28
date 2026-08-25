#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 4a -- Initialize Kubernetes Master (Control Plane)
#  Run as root on the MASTER node only:  sudo ./init-master.sh [hostname]
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

assert_root

NEW_HOSTNAME="${1:-}"
POD_NETWORK_CIDR="${2:-10.244.0.0/16}"

log_message "=== Initialize Kubernetes Master ===" "info"

# -- Optionally set hostname -----------------------------------------------
if [[ -n "$NEW_HOSTNAME" ]]; then
    log_message "Setting hostname to: $NEW_HOSTNAME" "info"
    hostnamectl set-hostname "$NEW_HOSTNAME"
fi

# -- Ensure CRI-O is running -----------------------------------------------
log_message "Ensuring CRI-O service is active..." "info"
systemctl enable crio.service
systemctl start crio.service
log_message "CRI-O is running." "success"

# -- Initialize the cluster ------------------------------------------------
log_message "Running kubeadm init (pod CIDR: $POD_NETWORK_CIDR)..." "info"
kubeadm init --pod-network-cidr="$POD_NETWORK_CIDR"

# -- Configure kubectl for the current user --------------------------------
log_message "Configuring kubectl..." "info"
mkdir -p "$HOME/.kube"
cp -i /etc/kubernetes/admin.conf "$HOME/.kube/config"
chown "$(id -u):$(id -g)" "$HOME/.kube/config"
log_message "kubectl configured for $(whoami)." "success"

# -- Post-init instructions ------------------------------------------------
echo ""
log_message "========================================" "success"
log_message "Master initialized successfully!" "success"
log_message "========================================" "success"
echo ""
echo "  Next steps:"
echo ""
echo "  1. Apply a network plugin (Weave Net):"
echo "     kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml"
echo ""
echo "  2. On each WORKER node, run the join command printed above."
echo "     If you lost it, regenerate with:"
echo "     sudo kubeadm token create --print-join-command"
echo ""
echo "  3. Verify nodes:"
echo "     kubectl get nodes"
echo "     kubectl -n kube-system get pods"
echo ""
echo "  Troubleshooting:"
echo "     kubeadm reset --force    # Start fresh"
echo "     journalctl -xeu kubelet  # Check kubelet logs"
echo ""
