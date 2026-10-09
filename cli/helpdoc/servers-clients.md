# gitmap servers-clients

Broadcast shell commands, git operations, machine joins, or installations across all server and client nodes in the cluster.

## Usage

```bash
gitmap servers-clients <subcommand> [args] [flags]
```

Aliases: `sc`, `servers-client`

## Subcommands Overview

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `bash <cmd>` | | Execute a Bash command or script pipeline across all cluster nodes |
| `sh <cmd>` | `shell` | Execute a POSIX shell command across all cluster nodes |
| `ps <cmd>` | | Execute a PowerShell command on Windows cluster nodes |
| `cmd <cmd>` | | Execute a Windows Command Prompt command on Windows nodes |
| `join <target> [alias]` | `add`, `enroll` | Join and enroll a machine into the cluster registry with key auth |
| `nodes` | `ls`, `list`, `joined`, `machines` | List all registered machines joined to the cluster topology |
| `remove <target>` | `rm`, `delete` | Remove a machine from the cluster registry and key database |
| `ping` | `status`, `health` | Check reachability and round-trip ping latency for all cluster nodes |
| `install <pkgs>` | | Install packages (comma-separated list) or GitMap across all nodes |
| `restart [duration]` | | Trigger or schedule system restart across all cluster nodes |
| `shutdown [duration]` | | Trigger or schedule system shutdown across all cluster nodes |
| `pull --all` | | Run git pull across all repositories on all nodes |
| `push --all` | | Run git push across all repositories on all nodes |
| `commit --all` | | Run automated git commit across all nodes |
| `status --all` | | Show combined dirty/clean repository status across all nodes |
| `proj <name> run` | | Run project-level automation or build scripts on all nodes |

## Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet (e.g. `--except 2,151`) |
| `--ip <list>` | Target specific node IP addresses exclusively |
| `--id <list>` | Target specific node Display IDs exclusively |
| `--json` | Output machine list or results in JSON format for scripting |
| `--yes`, `-Y` | Bypass preflight confirmation prompts |
| `--dry-run` | Preview execution plan and target nodes without running commands |
| `--verbose` | Print detailed real-time execution logs and SSH handshakes |

---

## Detailed Command Descriptions

### `bash <cmd>` & `sh <cmd>` (Linux & macOS Execution)
Dispatches remote shell commands across Linux and macOS cluster nodes in parallel. The command string is passed to `/bin/bash -c "<cmd>"` or `/bin/sh -c "<cmd>"` on each target host. Stdout and stderr from all machines are captured in real-time, prefixed by the target node's alias (e.g. `[worker-1] <output>`), and summarized at the end with process exit codes.

### `ps <cmd>` & `cmd <cmd>` (Windows Node Execution)
Dispatches native Windows commands to Windows cluster members. `ps` executes inside `powershell.exe -Command "<cmd>"`, while `cmd` executes inside `cmd.exe /c "<cmd>"`. Outputs are cleanly decoded and formatted.

### `join <target> [alias]` (Admission & Identity Layer)
Admits a remote machine into the cluster topology. Under the hood, `join`:
1. Validates network connectivity and SSH port reachability via fast TCP handshake.
2. Auto-accepts and records host keys into `~/.ssh/known_hosts` (preventing prompt hangs).
3. Discovers credentials: tests default SSH keys, or prompts interactively with hidden input.
4. Encrypts and securely vaults passwords at rest using RSA-OAEP / AES-GCM encryption (`gitmap ssh pass show <alias>`).
5. Probes and records remote OS metadata (`linux`, `windows`, `darwin`), OS version, and first-run timestamp.
6. Auto-checks or bootstraps GitMap agent on the remote machine.
7. Registers the host record in GitMap's SQLite database (`ssh_hosts` and `SSHConnection` tables).
Once joined, you can address this machine in all future `sc`, `clients`, and `cluster` commands by its friendly alias or IP address.

### `nodes` / `ls` / `list` (Cluster Registry Inspection)
Queries the local SQLite database and renders an ASCII table of all enrolled cluster machines. Columns include Display ID, Host Alias, IP Address, SSH User, Cluster Role (`control` vs `worker`), and Admission Timestamp. Pass `--json` to stream structured data into pipelines.

### `remove <target>` / `rm` (Machine Decommissioning)
De-registers a node from the cluster topology and removes its cryptographic keys and credentials from the local database. Accepts host alias, IP, or display ID.

### `ping` / `status` (Health & Latency Telemetry)
Sends lightweight non-destructive SSH probes to every registered node simultaneously. Accurately measures round-trip latency (RTT in milliseconds), verifies SSH authentication, and highlights unreachable or offline nodes.

### `install <packages>` (Multi-Node Package & GitMap Provisioning)
Dispatches software package installation across all cluster nodes simultaneously. It automatically detects the remote package manager (`apt-get`, `dnf`, `yum`, or `winget`) and installs specified packages (e.g. `gitmap sc install "git,nodejs,curl"`). You can also run `gitmap sc install gitmap` to distribute and bootstrap GitMap itself across all cluster machines!

### `pull --all`, `push --all`, `status --all` (Distributed Git Operations)
Runs fleet-wide Git management across every joined machine. Allows synchronizing codebases, checking repository dirty states, and pushing commits across entire developer teams or build server clusters in a single command.

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

---

## Examples

### 1. Bash & POSIX Shell Fan-Out
```bash
# Check system kernel and uptime across all nodes
gitmap servers-clients bash "uname -a && uptime"

# Check disk space across all machines, excluding node 2
gitmap sc shell "df -h" --except 2

# Inspect running Docker containers on all cluster workers
gitmap sc bash "docker ps --format 'table {{.Names}}\t{{.Status}}'"

# Execute quick POSIX shell commands
gitmap sc sh "cat /etc/os-release | grep PRETTY_NAME"
```

### 2. Machine Joining & Enrollment
```bash
# Join a machine by IP with alias 'worker-1'
gitmap servers-clients join 192.168.1.14 worker-1

# Join using SSH username, IP, and key authentication
gitmap sc join alim@192.168.1.14 devbox --key ~/.ssh/id_rsa

# Join with root credentials and password
gitmap sc join root@192.168.1.50 control-plane --password secret

# Enroll an additional node using the add alias
gitmap sc add 10.0.0.5 worker-2
```

### 3. Inspecting & Probing Joined Machines
```bash
# View table of joined nodes (ID, ALIAS, IP, USER, ROLE, CREATED_AT)
gitmap sc nodes
gitmap sc ls

# Output joined machines in JSON format for automation
gitmap sc ls --json

# Measure ping round-trip latency and health across all nodes
gitmap sc ping
```

### 4. Package Installation & Remote GitMap Setup
```bash
# Install common developer tools across cluster simultaneously
gitmap servers-clients install "git,nodejs,curl,htop" --except 2

# Install backend container and web stack
gitmap sc install "docker.io,build-essential,nginx"

# Install GitMap remotely across all cluster nodes via SSH
gitmap sc install gitmap

# Run remote one-liner installation script directly
gitmap sc bash "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/scripts/install.sh | bash"
```

### 5. Running GitMap Commands Across Nodes
```bash
# Check GitMap status across all remote nodes
gitmap sc bash "gitmap status"

# Run GitMap pull across all remote repositories
gitmap sc pull --all

# Push branches across all nodes, excluding node 3
gitmap sc push --all --except 3

# Check working tree cleanliness across cluster
gitmap sc status --all
```

### 6. Windows Node PowerShell & Command Prompt
```bash
# Query active services via PowerShell
gitmap sc ps "Get-Service | Where Status -eq Running"

# Check IP configuration via Windows Command Prompt, excluding specific nodes
gitmap sc cmd "ipconfig /all" --except 24,151

# Query top 5 CPU-intensive processes
gitmap sc ps "Get-Process | Sort-Object CPU -Descending | Select-Object -First 5"
```

### 7. Scheduled System Power Management
```bash
# Schedule reboot in 15 minutes across workers
gitmap sc restart 15m --except control-plane

# Schedule shutdown across all nodes at end of day
gitmap sc shutdown 2h
```

See also: `gitmap sc`, `gitmap clients`, `gitmap cluster`, `gitmap cluster exec`, `gitmap cluster nodes`, `gitmap ssh-join`, `gitmap sj`
