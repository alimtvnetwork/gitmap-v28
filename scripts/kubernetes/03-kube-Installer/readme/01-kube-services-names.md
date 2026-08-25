# Kube Services

1. kube-apiserver
2. kube-controller-manager
3. kube-scheduler
4. kubelet
5. kube-proxy
6. etcd
7. crio (or containerd)
8. kubeadm

## Start Stop

```bash
# kube-apiserver
sudo systemctl enable kube-apiserver && \
sudo systemctl start kube-apiserver && \
sudo systemctl restart kube-apiserver && \
sudo systemctl status kube-apiserver && \
sudo systemctl stop kube-apiserver

# kube-controller-manager
sudo systemctl enable kube-controller-manager && \
sudo systemctl start kube-controller-manager && \
sudo systemctl restart kube-controller-manager && \
sudo systemctl status kube-controller-manager && \
sudo systemctl stop kube-controller-manager

# kube-scheduler
sudo systemctl enable kube-scheduler && \
sudo systemctl start kube-scheduler && \
sudo systemctl restart kube-scheduler && \
sudo systemctl status kube-scheduler && \
sudo systemctl stop kube-scheduler

# kubelet
sudo systemctl enable kubelet && \
sudo systemctl start kubelet && \
sudo systemctl restart kubelet && \
sudo systemctl status kubelet && \
sudo systemctl stop kubelet

# kube-proxy
sudo systemctl enable kube-proxy && \
sudo systemctl start kube-proxy && \
sudo systemctl restart kube-proxy && \
sudo systemctl status kube-proxy && \
sudo systemctl stop kube-proxy

# etcd
sudo systemctl enable etcd && \
sudo systemctl start etcd && \
sudo systemctl restart etcd && \
sudo systemctl status etcd && \
sudo systemctl stop etcd

# crio (or containerd)
sudo systemctl enable crio && \
sudo systemctl start crio && \
sudo systemctl restart crio && \
sudo systemctl status crio && \
sudo systemctl stop crio

# kubeadm
sudo systemctl enable kubeadm && \
sudo systemctl start kubeadm && \
sudo systemctl restart kubeadm && \
sudo systemctl status kubeadm && \
sudo systemctl stop kubeadm
```

Now, here are the `systemctl` commands for each service in sequence:

```bash
# Enable services
sudo systemctl enable kube-apiserver && \
sudo systemctl enable kube-controller-manager && \
sudo systemctl enable kube-scheduler && \
sudo systemctl enable kubelet && \
sudo systemctl enable kube-proxy && \
sudo systemctl enable etcd && \
sudo systemctl enable crio && \
sudo systemctl enable kubeadm

# Start services
sudo systemctl start kube-apiserver && \
sudo systemctl start kube-controller-manager && \
sudo systemctl start kube-scheduler && \
sudo systemctl start kubelet && \
sudo systemctl start kube-proxy && \
sudo systemctl start etcd && \
sudo systemctl start crio && \
sudo systemctl start kubeadm

# Restart services
sudo systemctl restart kube-apiserver && \
sudo systemctl restart kube-controller-manager && \
sudo systemctl restart kube-scheduler && \
sudo systemctl restart kubelet && \
sudo systemctl restart kube-proxy && \
sudo systemctl restart etcd && \
sudo systemctl restart crio && \
sudo systemctl restart kubeadm

# Check status of services
sudo systemctl status kube-apiserver && \
sudo systemctl status kube-controller-manager && \
sudo systemctl status kube-scheduler && \
sudo systemctl status kubelet && \
sudo systemctl status kube-proxy && \
sudo systemctl status etcd && \
sudo systemctl status crio && \
sudo systemctl status kubeadm

# Stop services
sudo systemctl stop kube-apiserver && \
sudo systemctl stop kube-controller-manager && \
sudo systemctl stop kube-scheduler && \
sudo systemctl stop kubelet && \
sudo systemctl stop kube-proxy && \
sudo systemctl stop etcd && \
sudo systemctl stop crio && \
sudo systemctl stop kubeadm
```
