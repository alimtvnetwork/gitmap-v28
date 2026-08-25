# 02-kube-commands: Kubernetes Shell Script Integration

## 1. Overview
This specification details the integration of existing bare-metal Kubernetes shell scripts (located in `scripts/kubernetes/`) into the `gitmap` CLI via Go bindings. It bridges the local Go CLI to remote execution over SSH using the `gitmap ssh` layer.

## 2. Command Definitions

### 2.1 Kubernetes Install (`gitmap kube install <target>`)
**Behavior:** Connects to the remote target and executes the `02-ubuntu-prereq` and `03-kube-install` scripts sequentially.
- Pre-flight checks: ensures target OS is Ubuntu/Debian.
- Pipes the scripts over SSH.

### 2.2 Kubernetes Init Master (`gitmap kube init-master <target>`)
**Behavior:** Executes the `04-kube-init/init-master.sh` script on the specified target.
- Captures the generated `kubeadm join` command from stdout.
- Stores the join token and hash in the SQLite database (`kube_cluster` table).

### 2.3 Kubernetes Join Worker (`gitmap kube join-worker <target>`)
**Behavior:** Fetches the stored token from SQLite and executes `kubeadm join` on the remote target.
- Retrieves `<master-ip>`, `<token>`, and `<hash>` securely.
- Connects to the worker node via SSH and runs the join command.

### 2.4 Kubernetes Rollout (`gitmap kube rollout config.json`)
**Behavior:** Fully automates the cluster provisioning from the configuration template.
- Parses `config.json` (master IPs, worker IPs, user credentials).
- Executes `install` on all nodes concurrently.
- Executes `init-master` on the control plane.
- Executes `join-worker` on all worker nodes.
- Orchestrates via goroutines.

## 3. Database Schema Extensions (SQLite)

### `kube_cluster` Table
- `id` (TEXT PRIMARY KEY)
- `master_ip` (TEXT)
- `join_token` (TEXT)
- `ca_cert_hash` (TEXT)
- `created_at` (DATETIME)

## 4. Execution Model
- Reuses the `SpawnSSH` and `executeRemoteInstall` infrastructure built during the SSH integration phase.
- Uses `go:embed` or reads from the local `scripts/kubernetes/` directory to stream payload contents into the SSH stdin.
