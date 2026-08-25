#!/usr/bin/env bash
# --------------------------------------------------------------------------
#  Step 2 -- Ubuntu Prerequisites
#  Installs base packages needed before Kubernetes setup.
#  Run as root:  sudo ./run.sh
# --------------------------------------------------------------------------
set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../01-base-helpers/import-all.sh"

assert_root

log_message "=== Ubuntu Prerequisites ===" "info"

# -- Disable swap (required by kubelet) ------------------------------------
log_message "Disabling swap..." "info"
swapoff -a
sed -i '/ swap / s/^/#/' /etc/fstab
log_message "Swap disabled." "success"

# -- Load kernel modules ---------------------------------------------------
log_message "Loading kernel modules (overlay, br_netfilter)..." "info"
cat <<EOF | tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
modprobe overlay
modprobe br_netfilter
log_message "Kernel modules loaded." "success"

# -- Sysctl params for Kubernetes networking --------------------------------
log_message "Configuring sysctl for Kubernetes networking..." "info"
cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
sysctl --system >/dev/null 2>&1
log_message "Sysctl configured." "success"

# -- Install base packages -------------------------------------------------
log_message "Installing base packages..." "info"
install_apt curl apt-transport-https ca-certificates gnupg lsb-release \
    wget nano vim git sshpass jq software-properties-common

log_message "Ubuntu prerequisites complete." "success"
