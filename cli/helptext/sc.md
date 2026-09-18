# gitmap sc

Short alias for `servers-clients`. Broadcasts shell commands, git operations, machine joins, or installations across all server and client nodes in the cluster.

## Usage

```bash
gitmap sc <subcommand> [args] [flags]
```

Full command: `gitmap servers-clients <subcommand> [args] [flags]`

## Subcommands Overview

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `bash <cmd>` | | Execute a Bash command across all cluster nodes |
| `sh <cmd>` | `shell` | Execute a POSIX shell command across all cluster nodes |
| `ps <cmd>` | | Execute a PowerShell command on Windows cluster nodes |
| `cmd <cmd>` | | Execute a Windows Command Prompt command on Windows nodes |
| `join <target> [alias]` | `add`, `enroll` | Join and enroll a machine into the cluster registry with key auth |
| `nodes` | `ls`, `list`, `joined`, `machines` | List all registered machines joined to the cluster |
| `remove <target>` | `rm`, `delete` | Remove a machine from the cluster registry |
| `ping` | `status`, `health` | Check reachability and ping latency for all cluster nodes |
| `install <pkgs>` | | Install packages (comma-separated list) or GitMap across all nodes |
| `restart [duration]` | | Trigger or schedule system restart across all cluster nodes |
| `shutdown [duration]` | | Trigger or schedule system shutdown across all cluster nodes |
| `pull --all` | | Run git pull --all across all repositories on all nodes |
| `push --all` | | Run git push --all across all repositories on all nodes |
| `commit --all` | | Run automated git commit across all nodes |
| `status --all` | | Show combined dirty/clean status across all nodes |
| `proj <name> run` | | Run project-level automation on all nodes |

## Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet (e.g. `--except 2,151`) |
| `--ip <list>` | Target specific node IP addresses |
| `--id <list>` | Target specific node Display IDs |
| `--json` | Output machine list or results in JSON format |
| `--yes`, `-Y` | Bypass preflight confirmation prompt |
| `--dry-run` | Preview execution plan without running commands |
| `--verbose` | Print detailed real-time execution logs |

---

## Detailed Command Descriptions

- **`bash <cmd>` / `sh <cmd>`**: Dispatches remote shell commands to Linux/macOS nodes in parallel. Stdout/stderr are streamed with node alias prefixes (`[worker-1] ...`).
- **`ps <cmd>` / `cmd <cmd>`**: Dispatches PowerShell or CMD commands to Windows cluster nodes.
- **`join <target> [alias]`**: Admits a node into the cluster topology, checks SSH connectivity, authorizes keys, encrypts credentials, and registers it in SQLite.
- **`nodes` / `ls`**: Lists all enrolled machines with display ID, alias, IP, user, role, and registration date.
- **`ping` / `status`**: Non-destructive network health check measuring round-trip latency in ms across all nodes.
- **`install <packages>`**: Installs software packages simultaneously (e.g. `"git,nodejs,curl"`) or deploys GitMap itself via `gitmap sc install gitmap`.
- **`pull --all` / `push --all` / `status --all`**: Synchronizes Git repositories across all machines in the cluster fleet.

---

## Examples

### Bash & Shell Commands
```bash
gitmap sc bash "uname -a && uptime"
gitmap sc shell "df -h" --except 2
gitmap sc sh "docker ps"
```

### Join & Manage Machines
```bash
gitmap sc join 192.168.1.14 worker-1
gitmap sc join root@192.168.1.50 control-plane --password secret
gitmap sc nodes
gitmap sc ls --json
gitmap sc rm worker-1
gitmap sc ping
```

### Running GitMap & Package Installation
```bash
# Install developer tools across all cluster nodes
gitmap sc install "git,nodejs,curl" --except 2

# Install GitMap remotely across all cluster nodes
gitmap sc install gitmap

# Run GitMap commands across remote fleet
gitmap sc bash "gitmap status"
gitmap sc pull --all
gitmap sc push --all --except 3
gitmap sc status --all
```

### PowerShell & Windows Cmd
```bash
gitmap sc ps "Get-Date"
gitmap sc cmd "whoami" --except 2
gitmap sc ps "Get-Service | Where Status -eq Running"
```

## Subsystems Architecture Comparison (compare / matrix)

Display comparison table between `ssh`, `cluster`, and `sc` (`servers-clients`):

    $ gitmap sc compare

| Subsystem | Primary Focus | Join Command | Exec Command | Monitoring | Best Used When |
|---|---|---|---|---|---|
| `gitmap ssh` | Direct node management | `gitmap ssh join <u@ip>` | `gitmap ssh exec <cmd>` | `gitmap ssh scan` | Ad-hoc terminal commands, install/update, AGY/code open |
| `gitmap cluster` | Multi-node cluster orchestration | `gitmap cluster node add <ip>` | `gitmap cluster exec <target> <cmd>` | `gitmap cluster node ls` | K8s bootstrap, cluster recipes, distributed scripts |
| `gitmap sc` | Servers-clients fleet daemon | `gitmap sc join <server-url>` | `gitmap sc exec <cmd>` | `gitmap sc status` | Master-worker topology, continuous sync, live telemetry |

#### When to use which command:
- **`ssh`**: Fast, lightweight, direct command execution over standard SSH. Ideal for developer workstations, ad-hoc maintenance, and AGY/VS Code remote opening.
- **`cluster`**: Role-based infrastructure orchestration (`control` vs `workers`), provisioning recipes (Netplan IP, users, apt purge), and full Kubernetes lifecycle.
- **`sc` (`servers-clients`)**: High-speed fan-out broadcasts across entire fleet with concurrency pools, multi-shell execution, and continuous daemon synchronization.

See also: `gitmap servers-clients`, `gitmap clients`, `gitmap cluster`, `gitmap cluster exec`, `gitmap ssh`, `gitmap ssh-join`
