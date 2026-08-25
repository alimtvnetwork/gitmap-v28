#!/bin/bash

echo "Md Alim Ul Karim"
echo "https://alimkarim.com | https://riseup-pro.com"
echo

# Function to configure Kubernetes and optionally change hostname
configure_k8s() {
    local new_hostname=$1

    if [ -n "$new_hostname" ]; then
        echo "Changing hostname to $new_hostname"
        sudo hostnamectl set-hostname "$new_hostname"
    fi

    cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
    sudo modprobe overlay
    sudo modprobe br_netfilter

    # sysctl params required by setup, params persist across reboots
    cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

    sudo sysctl --system && \
    lsmod | grep br_netfilter && \
    lsmod | grep overlay && \
    sudo sysctl net.bridge.bridge-nf-call-iptables net.bridge.bridge-nf-call-ip6tables net.ipv4.ip_forward && \
    sudo sysctl -w net.ipv4.ip_forward=1 && \
    sudo curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list && \
    sudo curl -fsSL https://pkgs.k8s.io/addons:/cri-o:/prerelease:/main/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/cri-o-apt-keyring.gpg
    echo "deb [signed-by=/etc/apt/keyrings/cri-o-apt-keyring.gpg] https://pkgs.k8s.io/addons:/cri-o:/prerelease:/main/deb/ /" | sudo tee /etc/apt/sources.list.d/cri-o.list && \
    sudo apt-get install -y apt-transport-https ca-certificates curl gpg sshpass && \
    sudo apt update -y && \
    sudo apt-get update -y && \
    sudo apt upgrade -y && \
    sudo apt-get install -y cri-o kubelet kubeadm kubectl socat && \
    sudo sysctl -w net.ipv4.ip_forward=1 && \
    sudo systemctl enable crio.service && \
    sudo systemctl start crio.service && \
    sudo kubeadm config images pull && \
    echo "replace -192.168.0.20 with your ip" && \
    echo "run the below command :" && \
    echo "  $ kubeadm init"
    echo "  ========== after that, before running join command ==========="
    echo "  $ mkdir -p \$HOME/.kube"
    echo "  $ sudo cp -i /etc/kubernetes/admin.conf \$HOME/.kube/config"
    echo "  $ chown \$(id -u):\$(id -g) \$HOME/.kube/config"
    echo "  ========== before running join command ==========="
    echo "  ========== change hostname if cloned from same host ==========="
    echo "  hostnamectl set-hostname <your hostname>"
    echo "  hostnamectl set-hostname masterK"
    echo "  hostnamectl set-hostname worker1"
    echo "  hostnamectl set-hostname worker2"
    echo "  ========== after joining, run the below in master node ==========="
    echo "  $ kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml"
    echo "  $ kubectl -n kube-system get pods"
}

# Call the function with the optional hostname parameter
configure_k8s "$1"
