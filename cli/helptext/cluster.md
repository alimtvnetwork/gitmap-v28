# gitmap cluster

Comprehensive multi-node cluster management, topology enrollment, parallel command execution, remote script deployment, and Ubuntu node provisioning suite.

## What is GitMap Cluster?

`gitmap cluster` is GitMap's orchestration and infrastructure management engine designed for managing multi-machine environments, control-plane and worker topologies, automated Ubuntu provisioning, and complete Kubernetes lifecycles directly over SSH without heavyweight external daemons.

### Comparison: Local Git `exec` vs Macro `execute` vs Remote `cluster exec` / `sc`

To eliminate command ambiguity:

| Scope | Command | Purpose & Execution Context |
|-------|---------|-----------------------------|
| **Local Repos** | `gitmap exec <git-cmd>` (`x`) | Runs a `git` command locally across all indexed repositories on the current machine (e.g. `gitmap exec fetch --prune`). |
| **Local Automation** | `gitmap execute <macro>` | Replays recorded multistep workflow automation scripts sequentially on the local workstation. |
| **Remote Cluster** | `gitmap cluster exec <target> <cmd>` | Dispatches remote shell commands across role-targeted Kubernetes / Ubuntu cluster nodes (`all`, `control`, `workers`) over SSH. |
| **Fleet Fan-Out** | `gitmap sc bash <cmd>` (`sc`) | High-speed broadcast fan-out across all enrolled servers and client nodes with concurrency pools. |

---

## Architecture Overview

Gitmap Cluster provides an enterprise-grade orchestration layer directly over standard SSH transport without requiring agent daemons on target nodes:

- **Centralized Topology Storage**: Enrolled nodes are persisted in the local SQLite database (`ssh_hosts` and `ssh_history` tables) with host aliases, IP addresses, SSH ports, cluster roles (`control` vs `worker`), and creation timestamps.
- **Encrypted Credential Envelope**: Node passwords are encrypted at rest using RSA/AES-GCM encryption. Credentials are only decrypted transiently in memory during execution to feed `sudo -S` elevation prompts.
- **Concurrent Execution Engine**: Commands and scripts are dispatched across nodes using bounded worker pools (`--parallel <n>`, default 4) with real-time, interleaved stdout/stderr streaming prefixed by host alias (`[<alias>] <output>`).
- **Transactional Summary Telemetry**: Execution results across all targeted nodes are aggregated and rendered into an ASCII summary table detailing node alias, role, IP address, status (`SUCCESS` / `FAILED`), process exit code, and execution duration.

---

## GitMap Cluster Triad Architecture

GitMap orchestrates distributed multi-node infrastructure through a cohesive triad of complementary layers:

```
+-----------------------------------------------------------------------------+
|                          GITMAP CLUSTER TRIAD                               |
+-----------------------------------------------------------------------------+
|                                                                             |
|  1. ADMISSION & IDENTITY LAYER: `ssh-join` (`sj`)                           |
|     * Machine discovery, credential enrollment, and SSH key authorization   |
|     * Persistent host aliases (`devbox`, `prod-db`) and encrypted passwords |
|     * Connectivity probing, latency monitoring, and zero-collision recall   |
|                                                                             |
|  2. TOPOLOGY & ORCHESTRATION LAYER: `cluster`                               |
|     * Role-based cluster topology (`control` vs `worker`) & JSON imports   |
|     * Ubuntu provisioning recipes (Netplan IP, users, sudoers, apt cache)   |
|     * Complete Kubernetes lifecycle (CRI-O, kubeadm, Weave/Calico, Helm)    |
|                                                                             |
|  3. DISTRIBUTED FAN-OUT & BROADCAST LAYER: `servers-clients` (`sc`)         |
|     * Parallel command fan-out across all nodes or targeted subsets         |
|     * Multi-shell dispatch (Bash, POSIX sh, PowerShell, Windows cmd)        |
|     * Distributed Git operations (pull, push, status) & remote automation   |
|                                                                             |
+-----------------------------------------------------------------------------+
```

### Triad Roles & Workflow

- **`ssh-join` (`sj`) — Admission & Identity Layer**: First-class entrypoint for server admission. Enrolls remote machines into SQLite (`ssh_hosts`), securely manages RSA credentials, authorizes public keys, and prevents OpenSSH hostname collisions.
- **`cluster` — Topology & Orchestration Layer**: Manages multi-node architecture, control-plane vs worker roles, automated Ubuntu provisioning recipes, and full Kubernetes cluster setup.
- **`servers-clients` (`sc`) — Distributed Fan-Out Layer**: High-speed broadcast execution engine. Fans out shell commands, package installations, and git sync operations across cluster nodes in parallel.

---

## Usage

```bash
gitmap cluster <subcommand> [args...] [flags]
```

---

## Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `nodes` | `ls` | List registered nodes in the cluster |
| `add` | | Enroll new node into cluster topology |
| `join` | | Join node to cluster using SSH enrollment |
| `node` | | Manage cluster nodes and run Ubuntu provisioning recipes |
| `exec` | `run` | Execute remote shell command across targeted nodes |
| `run-script` | `script` | Deploy and execute local script file across targeted nodes |
| `bootstrap` | `bs` | Bootstrap nodes with RSA keys, passwordless sudo, and enrollment |
| `install` | | Install GitMap on remote target nodes via SSH (curl/PowerShell) |
| `init` | `template`, `init-config` | Generate starter cluster JSON schema (01-config.json) automatically |
| `import` | `import-config`, `import-cluster` | Import cluster topology from JSON and enroll hosts |
| `k8s` | `kube`, `kubernetes` | Manage Kubernetes cluster lifecycle, runtime, CNI, Helm, and NFS |
| `status` | `ping`, `health` | Display cluster health, active nodes, and ping latency |
| `remove` | `rm` | Remove node from cluster registry |
| `history` | `hi` | Inspect audit trail of past cluster executions |
| `export` | | Export node registry to JSON or CSV |
| `set-password` | | Set or update password credentials for nodes |
| `reset-password` | | Reset or clear stored credentials |
| `audit-clean` | | Clean legacy execution logs and expired node records |
| `stats` | | Show cluster statistics and execution metrics |

---

## Node Enrollment & Management (`add`, `join`, `nodes`, `status`, `remove`)

Gitmap provides first-class subcommands to enroll, list, probe, and remove cluster nodes:

- **Enroll Node (`add`)**: Register and enroll a new node into the cluster topology (`gitmap cluster add <user@ip|ip> [alias] [flags]`).
- **Join Node (`join`)**: Join a node to the cluster via SSH enrollment (`gitmap cluster join <user@ip|ip> [alias] [flags]`).
- **List Nodes (`nodes`, `ls`)**: List all registered cluster nodes with roles, OS, IPs, and status (`gitmap cluster nodes [--json]`).
- **Health & Ping (`status`, `ping`)**: Display heartbeat records and ping round-trip latency (`gitmap cluster status`, `gitmap cluster ping [target]`).
- **Remove Node (`remove`, `rm`)**: Unregister a node from the cluster topology (`gitmap cluster rm <alias|ip>`, `gitmap cluster node rm <alias|ip>`).

---

## Automated Cluster JSON Initialization (`init`, `template`)

Generate a starter cluster JSON configuration file automatically right from the terminal:

```bash
# Generate default 01-config.json with control plane and 2 workers
gitmap cluster init

# Custom control plane and worker IPs with custom username
gitmap cluster init --control 10.0.0.1 --workers 10.0.0.2,10.0.0.3 --user devadmin --out my-cluster.json

# Force overwrite existing file
gitmap cluster init --force
```

Terminal Output:
```
✓ Generated cluster topology configuration: 01-config.json
  Control Plane:
    • k8s-m1: 192.168.1.10
  Worker Nodes:
    • worker-1: 192.168.1.11
    • worker-2: 192.168.1.12
  SSH User: ubuntu

  Next: Edit credentials, then import into cluster:
    gitmap cluster import 01-config.json
```

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

## Remote Node Bootstrap & Sudoers NOPASSWD (`bootstrap`)

Gitmap provides an automated, idempotent remote node bootstrap command that configures SSH key authentication, grants passwordless sudo rules, verifies key access, and enrolls the node into the local SQLite topology database in a single step:

```bash
gitmap cluster bootstrap <target> [password] [flags]
gitmap cluster bs <target> [password] [flags]
gitmap sj bootstrap <target> [password] [flags]
```

### Bootstrap Lifecycle Contract

1. **RSA Cluster Keypair Verification**: Discovers existing 4096-bit RSA cluster keypair at `~/.ssh/id_rsa` or automatically generates a new keypair if missing.
2. **Target Resolution**: Supports semantic cluster groups (`all`, `control`, `workers`), host aliases, or direct IP addresses (`192.168.0.101` or `ubuntu@192.168.0.101`).
3. **Password Resolution**: Resolves credentials in priority order: CLI argument -> encrypted database password -> secure interactive prompt (`PromptSSHPassword`).
4. **Public Key Injection**: Pushes public key to remote `~/.ssh/authorized_keys` via `SSH_ASKPASS` (`attachAskPass`) with proper directory permissions (`0700` and `0600`).
5. **Passwordless Sudo Configuration**: When `--sudo` is enabled (default true), creates `/etc/sudoers.d/<user>` with `<user> ALL=(ALL) NOPASSWD:ALL` and sets permissions to `0440`.
6. **Key Authentication Verification**: Runs `ssh -i <key> -o BatchMode=yes <user>@<ip> "echo ssh_ok"` to verify passwordless SSH authentication without interactive prompts.
7. **Memory Hygiene**: Immediately clears plaintext password bytes and strings from memory upon completion.
8. **SQLite Host Enrollment**: Registers node credentials and aliases into `ssh_hosts` and records connection history into `ssh_history`.
9. **Summary Telemetry**: Formats and prints an execution summary table detailing `NODE`, `IP`, `SUDO`, `KEY_AUTH`, `STATUS`, and `DURATION`.

### Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--sudo` | `-s` | `true` | Configure `/etc/sudoers.d/<user>` with `NOPASSWD:ALL` (use `--sudo=false` or `--no-sudo` to skip) |
| `--user` | `-u` | | Remote SSH username override |
| `--port` | `-p` | `22` | Remote SSH port |
| `--password` | `-P` | | Remote user password (optional CLI override) |
| `--help` | `-h` | | Show bootstrap help message |

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

## Remote GitMap Installation & Auto-Bootstrapping (`install`)

Install or update GitMap on remote target nodes using official one-liners over SSH:

```bash
# Install GitMap across all cluster nodes in parallel
gitmap cluster install gitmap all

# Target specific node or alias
gitmap cluster install gitmap devbox
gitmap sj install gitmap k8s-w1

# Pin a specific GitMap version with custom concurrency
gitmap cluster install gitmap workers --version v6.256.0 --parallel 8
```

### Auto-Bootstrapping Engine
When executing commands via `gitmap cluster exec` or `gitmap sc bash`:
If a command requires `gitmap` and the binary is not present on the remote host (`command -v gitmap` fails), GitMap automatically detects this preflight condition, bootstraps and installs GitMap using curl (`install.sh`) or PowerShell (`install.ps1`), and then transparently runs the command.

---

## Examples

```bash
# 1. Enroll and join nodes into cluster topology
gitmap cluster add ubuntu@192.168.0.101 k8s-w1
gitmap cluster add 192.168.0.102
gitmap cluster join kube@192.168.0.103 k8s-w3 --auth
gitmap cluster node add 192.168.0.104 worker-4

# 2. List all registered cluster nodes
gitmap cluster nodes
gitmap cluster ls
gitmap cluster nodes --json
gitmap cluster node ls

# 3. Check cluster health, connectivity, and ping latency
gitmap cluster status
gitmap cluster ping
gitmap cluster ping k8s-w1
gitmap cluster health control

# 4. Import cluster topology from JSON configuration
gitmap cluster import 01-config.json

# 5. Bootstrap remote nodes with RSA keys, passwordless sudo, and enrollment
gitmap cluster bootstrap ubuntu@192.168.0.101 SecretPass123
gitmap cluster bootstrap workers
gitmap cluster bs control --sudo=false
gitmap sj bootstrap k8s-w1

# 6. Execute commands across targeted nodes in parallel with sudo elevation
gitmap cluster exec all "uptime"
gitmap cluster exec control "kubectl get nodes"
gitmap cluster exec workers "apt-get update -y" --sudo --parallel 6

# 7. Deploy and execute local script across targeted nodes
gitmap cluster run-script control ./scripts/install-k8s-control.sh --sudo
gitmap cluster run-script workers ./scripts/install-k8s-worker.sh --sudo

# 8. Ubuntu node provisioning recipes
gitmap cluster node set-ip k8s-w1 192.168.0.101 --route-ip 192.168.0.1
gitmap cluster node install-base all
gitmap cluster node create-user workers kube SecretPass123 --theme agnoster
gitmap cluster node set-theme all fletcherm
gitmap cluster node purge all

# 9. Kubernetes cluster lifecycle management
gitmap cluster k8s prereq all
gitmap cluster k8s install all --version 1.31
gitmap cluster k8s init control
gitmap cluster k8s cni control --plugin weave
gitmap cluster k8s join workers
gitmap cluster k8s nfs control --export-dir /nfsexport
gitmap cluster k8s helm-install control --version 3.16.2
gitmap cluster k8s helm-nfs control --export-dir /nfsexport
gitmap cluster k8s status control
gitmap cluster k8s reset k8s-w1

# 10. Remove nodes from cluster registry
gitmap cluster rm k8s-w1
gitmap cluster node rm worker-4
gitmap cluster remove --id node-123 --confirm

# 11. Inspect audit trail of past cluster executions
gitmap cluster history
gitmap cluster history RUN-20260817-001
gitmap cluster history --limit 10

# 12. Export node registry to JSON or CSV
gitmap cluster export --format json --output cluster-nodes.json
gitmap cluster export --format csv --output cluster-nodes.csv

# 13. Set or reset node password credentials
gitmap cluster set-password --id node-123
gitmap cluster reset-password --id node-123 --confirm

# 14. Clean legacy audit logs and inspect cluster statistics
gitmap cluster audit-clean --before 2026-01-01T00:00:00Z --confirm
gitmap cluster stats

# 15. GitMap Cluster Triad Operations
# Remote bash/shell execution across cluster nodes
gitmap sc bash "uname -a"
gitmap cluster exec all "apt-get update" --sudo
gitmap sc shell "df -h" --except control-plane

# Remote GitMap installation & auto-bootstrapping
gitmap cluster install gitmap all
gitmap sj install gitmap devbox

# Scheduled remote shutdown & restart
gitmap sc restart --except control-plane
gitmap schedule shutdown 1:45hr
gitmap schedule restart 2h
```

See also: `gitmap ssh`, `gitmap ssh-join`, `gitmap servers-clients`, `gitmap sc`
To view the full subsystems architecture matrix, run: `gitmap ssh compare`
