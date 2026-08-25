#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 3 -- Install Kubernetes (kubeadm, kubelet, kubectl)
#  Uses the official Kubernetes apt repository.
#  Run as root:  sudo ./run.sh [version]
#  Example:      sudo ./run.sh 1.31
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

assert_root

KUBE_VERSION="${1:-1.31}"
log_message "=== Install Kubernetes v${KUBE_VERSION} ===" "info"

# -- Install CRI-O container runtime ---------------------------------------
log_message "Installing CRI-O container runtime..." "info"

# Add CRI-O repository
CRIO_VERSION="v${KUBE_VERSION}"
curl -fsSL "https://pkgs.k8s.io/addons:/cri-o:/stable:/${CRIO_VERSION}/deb/Release.key" \
    | gpg --dearmor -o /etc/apt/keyrings/cri-o-apt-keyring.gpg 2>/dev/null

echo "deb [signed-by=/etc/apt/keyrings/cri-o-apt-keyring.gpg] https://pkgs.k8s.io/addons:/cri-o:/stable:/${CRIO_VERSION}/deb/ /" \
    | tee /etc/apt/sources.list.d/cri-o.list

install_apt cri-o
systemctl daemon-reload
systemctl enable --now crio
log_message "CRI-O installed and started." "success"

# -- Add Kubernetes apt repository ------------------------------------------
log_message "Adding Kubernetes apt repository..." "info"
mkdir -p /etc/apt/keyrings
curl -fsSL "https://pkgs.k8s.io/core:/stable:/v${KUBE_VERSION}/deb/Release.key" \
    | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg 2>/dev/null

echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v${KUBE_VERSION}/deb/ /" \
    | tee /etc/apt/sources.list.d/kubernetes.list

# -- Install kubeadm, kubelet, kubectl -------------------------------------
log_message "Installing kubeadm, kubelet, kubectl..." "info"
apt-get update -qq
install_apt kubelet kubeadm kubectl

# -- Hold versions to prevent accidental upgrades --------------------------
apt-mark hold kubelet kubeadm kubectl
log_message "Packages held at current version." "info"

# -- Enable kubelet --------------------------------------------------------
systemctl enable --now kubelet

# -- Print versions --------------------------------------------------------
KUBECTL_VER=$(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null | head -1)
KUBEADM_VER=$(kubeadm version -o short 2>/dev/null || echo "installed")
log_message "kubectl: $KUBECTL_VER" "success"
log_message "kubeadm: $KUBEADM_VER" "success"
log_message "Kubernetes installation complete." "success"
