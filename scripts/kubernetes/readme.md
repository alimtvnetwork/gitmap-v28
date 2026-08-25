# Kubernetes Cluster Setup (Unix/Ubuntu)

Reference scripts for setting up a bare-metal Kubernetes cluster on Ubuntu
using kubeadm. Inspired by
[kubernetes-training-v1](https://github.com/aukgit/kubernetes-training-v1)
by Md Alim Ul Karim.

## Folder Structure

| Folder / File | Purpose |
|---------------|---------|
| `01-base-helpers/` | Reusable shell helpers (logger, apt installer, package checker) |
| `02-ubuntu-prereq/` | Ubuntu prerequisites and server bootstrap |
| `03-kube-install/` | kubeadm, kubelet, kubectl installation |
| `04-kube-init/` | Cluster initialization (master + worker join) |
| `05-helm-install/` | Helm package manager installation |
| `06-nfs-setup/` | NFS server + Helm NFS provisioner |
| `07-remote-commands/` | Multi-node SSH command executor |
| `config-sample.json` | Node IP configuration template |
| `cheat-sheet.md` | Quick-reference kubectl/kubeadm commands |

## Quick Start

```bash
# 1. Copy config
cp config-sample.json config.json
# Edit config.json with your node IPs and credentials

# 2. Run on EACH node (master + workers)
chmod +x 02-ubuntu-prereq/run.sh && sudo ./02-ubuntu-prereq/run.sh
chmod +x 03-kube-install/run.sh  && sudo ./03-kube-install/run.sh

# 3. Initialize master (run on master node only)
chmod +x 04-kube-init/init-master.sh && sudo ./04-kube-init/init-master.sh

# 4. Join workers (run on each worker node)
# Use the join command printed by step 3
sudo kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>

# 5. Apply network plugin (run on master)
kubectl apply -f https://reweave.azurewebsites.net/k8s/v1.31/net.yaml

# 6. Install Helm (run on master)
chmod +x 05-helm-install/run.sh && sudo ./05-helm-install/run.sh
```

## Prerequisites

- Ubuntu 20.04+ (or Debian-based)
- Root/sudo access
- Internet connectivity
- Minimum 2 GB RAM per node
- Unique hostname per node

## Network Topology

```
+-------------------+
|   Master Node     |  192.168.0.20  (control plane)
+-------------------+
        |
   +---------+---------+
   |         |         |
+------+ +------+ +------+
|  W1  | |  W2  | |  W3  |
+------+ +------+ +------+
.0.21    .0.22    .0.23
```

## Related Scripts

- **Script 45** (`install-docker`) -- Docker Engine installation
- **Script 46** (`install-kubernetes`) -- Windows-based kubectl/minikube/Helm via Chocolatey

## Credits

Shell script patterns and cluster setup flow adapted from
[aukgit/kubernetes-training-v1](https://github.com/aukgit/kubernetes-training-v1)
(MIT License) by Md Alim Ul Karim.
