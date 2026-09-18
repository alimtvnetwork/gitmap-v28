# Subtask 01: Kubernetes Recipe Generators

## Objective
Implement modular bash recipe generator functions in Go ported from `kubernetes-training/03-kube-Installer/` for prerequisites, CRI-O runtime, kubelet/kubeadm/kubectl installation, cluster initialization, CNI network overlays, NFS storage, and Helm package management.

## Disjoint Files Assigned
- `cli/cmdssh/cluster_k8s_recipes.go`
- `cli/cmdssh/cluster_k8s_recipes_test.go`

## Implementation Details
1. `cli/cmdssh/cluster_k8s_recipes.go`:
   - `GenerateK8sPrereqScript() string`:
     Configures `/etc/modules-load.d/k8s.conf` (`overlay`, `br_netfilter`), runs `modprobe`, configures `/etc/sysctl.d/k8s.conf` (`net.bridge.bridge-nf-call-iptables = 1`, `net.bridge.bridge-nf-call-ip6tables = 1`, `net.ipv4.ip_forward = 1`), and runs `sysctl --system`.
   - `GenerateK8sInstallScript(version string, newHostname string) string`:
     Optionally changes hostname, configures `pkgs.k8s.io` apt keyrings and repo for Kubernetes `v1.31` (or passed version), configures CRI-O repo and keyring, installs `cri-o kubelet kubeadm kubectl socat apt-transport-https ca-certificates curl gpg`, enables and starts `crio.service`, and executes `kubeadm config images pull`.
   - `GenerateK8sInitScript(podCIDR string) string`:
     Ensures CRI-O is active, runs `kubeadm init`, copies `/etc/kubernetes/admin.conf` to `$HOME/.kube/config`, and adjusts permissions.
   - `GenerateCNIScript(plugin string) string`:
     If plugin is `"calico"` applies Calico manifest, otherwise applies Weave Net manifest (`https://reweave.azurewebsites.net/k8s/v1.31/net.yaml`).
   - `GenerateK8sJoinScript(joinCmd string) string`:
     Executes the join command with sudo.
   - `GenerateNFSServerScript(exportDir string) string`:
     Installs `nfs-server`, creates export directory (defaults to `/nfsexport`), appends to `/etc/exports` with `(rw,no_root_squash)`, and restarts `nfs-server`.
   - `GenerateHelmInstallScript(version string) string`:
     Downloads official Helm tarball (`https://get.helm.sh/helm-v<version>-linux-amd64.tar.gz`), unpacks, moves binary to `/usr/local/bin/helm`, and verifies installation.
   - `GenerateHelmNFSScript(nfsServerIP string, exportDir string) string`:
     Adds `nfs-subdir-external-provisioner` helm repo, uninstalls any existing release, and installs `nfs-provisioner` with specified server IP and path.
   - `GenerateK8sResetScript() string`:
     Runs `kubeadm reset --force && rm -rf $HOME/.kube`.

2. `cli/cmdssh/cluster_k8s_recipes_test.go`:
   - Add unit tests verifying script generation, parameter interpolation, and default values.

## Constraints
- Max function lines <= 15 (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal AppError wrapping.
- TOTAL BAN on running `go test`, `go build`, or local runner during execution.
- Strict Unix LF line endings.
