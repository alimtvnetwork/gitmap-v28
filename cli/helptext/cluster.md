# gitmap cluster

Comprehensive multi-node cluster management, topology enrollment, parallel command execution, remote script deployment, and Ubuntu node provisioning suite.

## Architecture Overview

Gitmap Cluster provides an enterprise-grade orchestration layer directly over standard SSH transport without requiring agent daemons on target nodes:

- **Centralized Topology Storage**: Enrolled nodes are persisted in the local SQLite database (`ssh_hosts` and `ssh_history` tables) with host aliases, IP addresses, SSH ports, cluster roles (`control` vs `worker`), and creation timestamps.
- **Encrypted Credential Envelope**: Node passwords are encrypted at rest using RSA/AES-GCM encryption. Credentials are only decrypted transiently in memory during execution to feed `sudo -S` elevation prompts.
- **Concurrent Execution Engine**: Commands and scripts are dispatched across nodes using bounded worker pools (`--parallel <n>`, default 4) with real-time, interleaved stdout/stderr streaming prefixed by host alias (`[<alias>] <output>`).
- **Transactional Summary Telemetry**: Execution results across all targeted nodes are aggregated and rendered into an ASCII summary table detailing node alias, role, IP address, status (`SUCCESS` / `FAILED`), process exit code, and execution duration.

---

## Usage

```bash
gitmap cluster <subcommand> [args...] [flags]
```

---

## Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `import` | `import-config`, `import-cluster` | Import cluster topology from JSON and enroll hosts |
| `exec` | `run` | Execute remote shell command across targeted nodes |
| `run-script` | `script` | Deploy and execute local script file across targeted nodes |
| `node` | | Execute Ubuntu node provisioning recipes (netplan, packages, user, zsh, purge) |
| `k8s` | `kube`, `kubernetes` | Manage Kubernetes cluster lifecycle, runtime, CNI, Helm, and NFS |
| `status` | | Display cluster health, active nodes, and ping latency |
| `nodes` | `ls` | List registered nodes in the cluster |
| `history` | `hi` | Inspect audit trail of past cluster executions |
| `export` | | Export node registry to JSON or CSV |
| `set-password` | | Set or update password credentials for nodes |
| `reset-password`| | Reset or clear stored credentials |
| `remove` | `rm` | Remove node from cluster registry |
| `audit-clean` | | Clean legacy execution logs and expired node records |
| `stats` | | Show cluster statistics and execution metrics |

---

## Cluster Topology JSON Import

You can declare and import entire cluster topologies using standard JSON manifests:

```bash
gitmap cluster import 01-config.json
```

### JSON Schema

The configuration file defines control-plane nodes, worker nodes, and initial SSH user credentials:

```json
{
  "control": {
    "k8s-m1": "192.168.0.100"
  },
  "nodes": {
    "k8s-w1": "192.168.0.101",
    "k8s-w2": "192.168.0.102",
    "k8s-w3": "192.168.0.103"
  },
  "user": {
    "name": "kube",
    "password": "SuperSecretPassword123"
  }
}
```

- Each control node entry is registered with `ClusterRole = "control"`.
- Each worker node entry is registered with `ClusterRole = "worker"`.
- The user password is encrypted and stored in `EncryptedPassword`.
- Host join events are immutably logged into `ssh_history`.

---

## Role Targeting

All cluster operations (`exec`, `run-script`, and `node`) support semantic target selectors:

| Target | Description | Example |
|--------|-------------|---------|
| `all` | Targets every node in the cluster (control and workers) | `gitmap cluster exec all "uptime"` |
| `control` | Targets only control-plane master nodes | `gitmap cluster exec control "kubectl get nodes"` |
| `workers` | Targets only worker nodes | `gitmap cluster exec workers "free -h"` |
| `<alias>` | Targets an individual node by alias or IP address | `gitmap cluster exec k8s-w1 "df -h"` |

---

## Sudo Elevation Rules

Commands requiring administrative privileges can specify the `-s` or `--sudo` flag:

```bash
gitmap cluster exec all "systemctl restart containerd" --sudo
```

### Security Contract
1. When `--sudo` is active, Gitmap decrypts the node's stored password into memory.
2. The remote command is piped securely into `sudo -S`:
   ```bash
   echo '<password>' | sudo -S bash -c '<command>'
   ```
3. Passwords are never written to disk or logged to shell histories.
4. If no password is stored, the command runs directly without password pipe.

---

## Remote Script Upload & Execution

Run arbitrary local bash, python, or shell scripts on cluster nodes:

```bash
gitmap cluster run-script <target> <local-script-path> [--sudo] [--parallel <n>]
```

### Execution Flow
1. **Validation**: Gitmap reads and validates the local script file.
2. **Staging**: Encodes script content into base64 and uploads it to remote `/tmp/gm_script_<timestamp>.sh`.
3. **Execution**: Runs the script using `bash /tmp/gm_script_...` (or elevated with `sudo -S`).
4. **Cleanup**: Automatically deletes the staged temporary script upon completion, even if the script failed.
5. **Output**: Live streams stdout/stderr and prints the final execution summary table.

---

## Ubuntu Node Provisioning Recipes

Gitmap provides built-in provisioning recipes to configure fresh Ubuntu nodes into Kubernetes-ready servers:

```bash
gitmap cluster node <recipe> [args...]
```

### Recipes

#### 1. Configure Static Netplan IP (`set-ip`)
Creates `/etc/netplan/00-installer-config.yaml` with ens33 static IP, default route gateway, and Google/Cloudflare DNS (`[8.8.8.8, 1.1.1.1]`), and applies configuration via `netplan apply`:

```bash
gitmap cluster node set-ip <target> <new-ip> [--route-ip <gw>]
```

- `<new-ip>`: Static IP with or without CIDR (defaults to `/24` if no slash is provided).
- `--route-ip <gw>`: Gateway route IP (defaults to `192.168.0.1`).

#### 2. Install Base Packages (`install-base`)
Installs essential system tools, development utilities, and shell runtimes:

```bash
gitmap cluster node install-base <target>
```

Packages installed: `curl`, `wget`, `git`, `zsh`, `net-tools`, `htop`, `build-essential`.

#### 3. Create Sudo User with Oh-My-Zsh (`create-user`)
Provisions a dedicated user with Zsh as default login shell, configures passwordless sudoers, installs Oh-My-Zsh, and configures the default theme:

```bash
gitmap cluster node create-user <target> <username> [password] [--theme <theme>]
```

- Configures `/etc/sudoers.d/<username>` with `0440` permissions.
- Adds user to the `sudo` group.
- Installs Oh-My-Zsh in unattended mode.
- `--theme`: ZSH theme name (defaults to `fletcherm`).

#### 4. Change ZSH Theme (`set-theme`)
Modifies `ZSH_THEME` in `~/.zshrc` across targeted nodes:

```bash
gitmap cluster node set-theme <target> <theme>
```

#### 5. Purge and Clean Apt Cache (`purge`)
Removes unused dependencies, unneeded kernels, and cleans the apt download cache to minimize VM storage:

```bash
gitmap cluster node purge <target>
```

---

## Kubernetes Cluster Lifecycle (`gitmap cluster k8s`)

Gitmap provides an end-to-end Kubernetes orchestration and lifecycle management suite ported from battle-tested production bootstrap recipes. It provisions lightweight, OCI-compliant CRI-O container runtimes, sets up multi-node Kubernetes clusters (`v1.31`+), deploys high-speed CNI overlay networks, configures automated worker joining via SSH token discovery, sets up centralized NFS shared storage, and deploys Helm with dynamic StorageClass provisioners.

### ASCII Cluster Topology

```
+---------------------------------------------------------------------------------------------------------+
|                                    KUBERNETES CLUSTER TOPOLOGY                                          |
+---------------------------------------------------------------------------------------------------------+
|                                                                                                         |
|   +-------------------------------------------------------------------+                                 |
|   |                  Control Plane (Master Node)                      |                                 |
|   |                  Alias: k8s-m1  |  IP: 192.168.0.100              |                                 |
|   |  +-------------------------------------------------------------+  |                                 |
|   |  | Kubernetes Control Components:                              |  |                                 |
|   |  |  * kube-apiserver          (:6443/TCP)                      |  |                                 |
|   |  |  * etcd cluster datastore  (:2379-2380/TCP)                 |  |                                 |
|   |  |  * kube-controller-manager (:10257/TCP)                    |  |                                 |
|   |  |  * kube-scheduler          (:10259/TCP)                     |  |                                 |
|   |  |  * kubelet API             (:10250/TCP)                     |  |                                 |
|   |  |  * CRI-O Container Runtime (crio.service)                   |  |                                 |
|   |  +-------------------------------------------------------------+  |                                 |
|   |  | Addons & Storage:                                           |  |                                 |
|   |  |  * Helm Package Manager    (v3.16.2 /usr/local/bin/helm)    |  |                                 |
|   |  |  * NFS Kernel Server       (:2049/TCP /nfsexport)           |  |                                 |
|   |  |  * NFS Subdir Provisioner  (Dynamic StorageClass)           |  |                                 |
|   |  +-------------------------------------------------------------+  |                                 |
|   +---------------------------------+---------------------------------+                                 |
|                                     |                                                                   |
|          ===========================+===========================                                        |
|         ||   Pod Network Overlay (CNI): Weave Net / Calico    ||                                        |
|         ||   Subnet CIDR: 10.244.0.0/16 or 10.32.0.0/12      ||                                        |
|          ===========================+===========================                                        |
|                                     |                                                                   |
|             +-----------------------+-----------------------+                                           |
|             |                                               |                                           |
|             v                                               v                                           |
|   +-----------------------------------+   +-----------------------------------+                         |
|   |           Worker Node 1           |   |           Worker Node 2           |                         |
|   |   Alias: k8s-w1                   |   |   Alias: k8s-w2                   |                         |
|   |   IP: 192.168.0.101               |   |   IP: 192.168.0.102               |                         |
|   |  * kubelet API (:10250/TCP)       |   |  * kubelet API (:10250/TCP)       |                         |
|   |  * kube-proxy                     |   |  * kube-proxy                     |                         |
|   |  * CRI-O Runtime (crio.service)   |   |  * CRI-O Runtime (crio.service)   |                         |
|   |  * NodePort Services (30000-32767)|   |  * NodePort Services (30000-32767)|                         |
|   |  * NFS Client Mount (/nfsexport)  |   |  * NFS Client Mount (/nfsexport)  |                         |
|   +-----------------------------------+   +-----------------------------------+                         |
|                                                                                                         |
+---------------------------------------------------------------------------------------------------------+
```

### Network Ports & Kernel Requirements

#### 1. Network Port Allocation Contract

| Scope | Port / Range | Protocol | Direction | Purpose |
|-------|--------------|----------|-----------|---------|
| Control Plane | `6443` | TCP | Inbound | Kubernetes API Server endpoint |
| Control Plane | `2379-2380` | TCP | Inbound | etcd server client & peer communication |
| Control & Workers | `10250` | TCP | Inbound | Kubelet API |
| Control Plane | `10259` | TCP | Inbound | kube-scheduler health and metrics |
| Control Plane | `10257` | TCP | Inbound | kube-controller-manager health and metrics |
| Worker Nodes | `30000-32767` | TCP | Inbound | NodePort service high port range |
| NFS Storage | `2049` | TCP | Inbound | NFS daemon RPC traffic |
| NFS Storage | `111` | TCP/UDP | Inbound | Portmapper (`rpcbind`) daemon |
| Weave Net CNI | `6783` | TCP | Inbound/Outbound | Weave control communication |
| Weave Net CNI | `6783-6784` | UDP | Inbound/Outbound | Weave fastdp / pcap data path overlay |
| Calico CNI | `179` | TCP | Inbound/Outbound | BGP routing protocol peering |
| Calico CNI | `4789` | UDP | Inbound/Outbound | VXLAN overlay encapsulation |
| Calico CNI | `51820` | UDP | Inbound/Outbound | WireGuard encrypted data plane |

#### 2. Kernel Modules (`/etc/modules-load.d/k8s.conf`)

- **`overlay`**: Required for union filesystem drivers that power OCI containers without storage degradation.
- **`br_netfilter`**: Ensures Linux bridge packets are intercepted and processed by iptables/ip6tables filters.

#### 3. Sysctl Network Forwarding (`/etc/sysctl.d/k8s.conf`)

- **`net.bridge.bridge-nf-call-iptables = 1`**: Enables bridge packets to be processed by iptables rules.
- **`net.bridge.bridge-nf-call-ip6tables = 1`**: Enables bridge IPv6 packets to be processed by ip6tables rules.
- **`net.ipv4.ip_forward = 1`**: Enables kernel IPv4 packet routing across virtual and physical network interfaces.

#### 4. Container Runtime Interface (CRI-O)

- **Repository**: Official `pkgs.k8s.io` repositories for Kubernetes `v1.31` core and CRI-O addons.
- **Daemon Management**: Managed through systemd as `crio.service`.
- **Pre-pulling**: Automated image caching via `kubeadm config images pull` ensures fast and reliable initialization.

---

### Subcommands & Syntax Reference

```bash
gitmap cluster k8s <subcommand> <target> [flags]
```

| Subcommand | Syntax | Description |
|------------|--------|-------------|
| `prereq` | `gitmap cluster k8s prereq <target>` | Load kernel modules (`overlay`, `br_netfilter`), set sysctl parameters, and install core network utilities |
| `install` | `gitmap cluster k8s install <target> [--version <ver>] [--hostname <name>]` | Configure apt keys, install CRI-O runtime and Kubernetes binaries (`kubelet`, `kubeadm`, `kubectl`, `socat`), enable services, and pre-pull container images |
| `init` | `gitmap cluster k8s init <target> [--pod-cidr <cidr>]` | Initialize control-plane node using `kubeadm init`, configure `$HOME/.kube/config`, and set ownership |
| `cni` | `gitmap cluster k8s cni <target> [--plugin weave\|calico]` | Deploy container network interface overlay (default: `weave`; optional: `calico`) |
| `join-command` | `gitmap cluster k8s join-command <target>` | Generate a fresh token on control plane and print the copy-pasteable join command |
| `join` | `gitmap cluster k8s join <target> [--command "<cmd>"]` | Join worker nodes to the control plane (automatically retrieves join token from control plane if `--command` is omitted) |
| `nfs` | `gitmap cluster k8s nfs <target> [--export-dir <dir>]` | Install NFS kernel server, export shared directory (`*(rw,no_root_squash)`), and restart services |
| `helm-install` | `gitmap cluster k8s helm-install <target> [--version <ver>]` | Download official Helm tarball, install binary to `/usr/local/bin/helm`, and verify version |
| `helm-nfs` | `gitmap cluster k8s helm-nfs <target> [--nfs-server <ip>] [--export-dir <dir>]` | Deploy `nfs-subdir-external-provisioner` Helm chart to enable dynamic PersistentVolume storage |
| `status` | `gitmap cluster k8s status <target>` | Inspect node status (`kubectl get nodes -o wide`), pods across all namespaces (`-A`), and CRI-O runtime health |
| `reset` | `gitmap cluster k8s reset <target>` | Tear down Kubernetes cluster state (`kubeadm reset --force`), clean `$HOME/.kube`, and reset network interfaces |

---

### Complete Lifecycle Workflow

#### Step 1: Kernel & System Prerequisites (`prereq`)

Prepare target nodes with required kernel modules, iptables bridge integration, and packet forwarding:

```bash
# Apply kernel modules and sysctl configuration to all cluster nodes
gitmap cluster k8s prereq all
```

This ensures `/etc/modules-load.d/k8s.conf` and `/etc/sysctl.d/k8s.conf` are populated and reloaded with `sysctl --system`.

#### Step 2: CRI-O & Kubernetes Binaries Installation (`install`)

Install container runtime and orchestration tooling across nodes:

```bash
# Install default Kubernetes v1.31 and CRI-O on all nodes
gitmap cluster k8s install all

# Override hostname and specify exact Kubernetes version
gitmap cluster k8s install control --version 1.31 --hostname masterK
gitmap cluster k8s install k8s-w1 --version 1.31 --hostname worker1
gitmap cluster k8s install k8s-w2 --version 1.31 --hostname worker2
```

This step imports keyrings from `pkgs.k8s.io`, enables and starts `crio.service`, and executes `kubeadm config images pull`.

#### Step 3: Control-Plane Initialization (`init`)

Bootstrap the control plane on the master node:

```bash
# Standard control plane initialization
gitmap cluster k8s init control

# Explicit pod network CIDR
gitmap cluster k8s init control --pod-cidr 10.244.0.0/16
```

The master node initializes etcd, API server, controller manager, and scheduler, automatically copying `/etc/kubernetes/admin.conf` to `$HOME/.kube/config` with correct user ownership.

#### Step 4: Container Network Interface Overlay Deployment (`cni`)

Deploy an overlay network so pods can communicate across nodes:

```bash
# Option A: Weave Net (Default - lightweight, automatic mesh)
gitmap cluster k8s cni control --plugin weave

# Option B: Calico (Production - enterprise BGP and NetworkPolicies)
gitmap cluster k8s cni control --plugin calico
```

Manifest URLs deployed:
- **Weave Net**: `https://reweave.azurewebsites.net/k8s/v1.31/net.yaml`
- **Calico**: `https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml`

#### Step 5: Automated Worker Joining (`join`)

Join worker nodes into the Kubernetes cluster:

```bash
# Fully Automated: Automatically extracts token from control plane and joins workers
gitmap cluster k8s join workers

# Explicit Command: Pass a pre-generated join command
gitmap cluster k8s join workers --command "kubeadm join 192.168.0.100:6443 --token abcdef.0123456789abcdef --discovery-token-ca-cert-hash sha256:4b3..."
```

> [!TIP]
> When executing `gitmap cluster k8s join workers` without `--command`, Gitmap uses SSH to query the control plane with `sudo kubeadm token create --print-join-command`, dynamically extracts the token and hash parameters, and dispatches the join command to all target worker nodes concurrently!

To inspect the raw join command at any time:
```bash
gitmap cluster k8s join-command control
```

#### Step 6: Shared NFS Storage Setup (`nfs`)

Provision a central NFS server on the control node or dedicated storage host:

```bash
# Set up NFS export directory (defaults to /nfsexport)
gitmap cluster k8s nfs control

# Set up custom export directory
gitmap cluster k8s nfs control --export-dir /data/k8s-storage
```

Configures `/etc/exports` with `<export-dir> *(rw,no_root_squash)` permissions, starts `nfs-server`, and verifies availability via `showmount -e`.

#### Step 7: Helm Package Manager Installation (`helm-install`)

Install the Helm v3 binary on the control plane:

```bash
# Install default Helm version (v3.16.2)
gitmap cluster k8s helm-install control

# Install custom Helm version
gitmap cluster k8s helm-install control --version 3.16.2
```

Downloads official binaries from `https://get.helm.sh/helm-v<version>-linux-amd64.tar.gz`, verifies unpacking, and installs to `/usr/local/bin/helm`.

#### Step 8: Dynamic NFS Storage Provisioner (`helm-nfs`)

Deploy the `nfs-subdir-external-provisioner` Helm chart to provide dynamic PersistentVolumes:

```bash
# Automatic server IP resolution
gitmap cluster k8s helm-nfs control

# Explicit NFS server IP and export directory
gitmap cluster k8s helm-nfs control --nfs-server 192.168.0.100 --export-dir /nfsexport
```

This creates a default `nfs-client` StorageClass. Any PVC requesting this StorageClass will automatically create a directory on the NFS server.

#### Step 9: Health & Cluster Status Inspection (`status`)

Verify complete cluster health, node statuses, and running pods:

```bash
gitmap cluster k8s status control
```

Outputs:
- `kubectl get nodes -o wide`: Verify all nodes are in `Ready` state with correct internal IPs.
- `kubectl get pods -A`: Verify all system pods (`kube-system`, `crio`, `nfs-provisioner`) are in `Running` state.
- Service telemetry: Confirms `crio.service` and `kubelet.service` active status.

#### Step 10: Teardown & Reset (`reset`)

Reset nodes to a clean state if you need to reconfigure or restart:

```bash
# Reset specific worker nodes
gitmap cluster k8s reset k8s-w1

# Reset all worker nodes
gitmap cluster k8s reset workers

# Full teardown of all nodes in the cluster
gitmap cluster k8s reset all
```

Executes `kubeadm reset --force`, purges `/etc/cni/net.d`, clears `$HOME/.kube`, and cleans local IPVS/iptables rules.

---

## Examples

```bash
# 1. Import topology from JSON
gitmap cluster import 01-config.json

# 2. Check cluster connectivity
gitmap cluster status

# 3. Configure static IP on a worker node
gitmap cluster node set-ip k8s-w1 192.168.0.101 --route-ip 192.168.0.1

# 4. Install base developer tools on all nodes in parallel
gitmap cluster node install-base all

# 5. Create kube user with zsh across all worker nodes
gitmap cluster node create-user workers kube SecretPass123 --theme agnoster

# 6. Run arbitrary commands in parallel with sudo
gitmap cluster exec all "apt-get update -y" --sudo --parallel 6

# 7. Upload and execute local Kubernetes install script
gitmap cluster run-script control ./scripts/install-k8s-control.sh --sudo
gitmap cluster run-script workers ./scripts/install-k8s-worker.sh --sudo

# 8. Clean up unused packages and apt cache
gitmap cluster node purge all

# 9. Prepare kernel modules and sysctl parameters on all nodes
gitmap cluster k8s prereq all

# 10. Install CRI-O runtime and Kubernetes v1.31 across all nodes
gitmap cluster k8s install all --version 1.31

# 11. Initialize control-plane node and deploy Weave Net CNI
gitmap cluster k8s init control
gitmap cluster k8s cni control --plugin weave

# 12. Automatically join all worker nodes using token auto-discovery
gitmap cluster k8s join workers

# 13. Set up NFS shared storage and deploy Helm dynamic provisioner
gitmap cluster k8s nfs control --export-dir /nfsexport
gitmap cluster k8s helm-install control --version 3.16.2
gitmap cluster k8s helm-nfs control --export-dir /nfsexport

# 14. Inspect cluster nodes, pods, and runtime health
gitmap cluster k8s status control
```

See also: `gitmap ssh`, `gitmap ssh-join`
