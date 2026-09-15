# gitmap servers-clients

Broadcast shell commands, git operations, machine joins, or installations across all server and client nodes in the cluster.

## Usage

```bash
gitmap servers-clients <subcommand> [args] [flags]
```

Aliases: `sc`, `servers-client`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `bash <cmd>` | Execute a Bash command across all cluster nodes |
| `sh <cmd>` / `shell <cmd>` | Execute a POSIX shell command across all cluster nodes |
| `ps <cmd>` | Execute a PowerShell command on all nodes |
| `cmd <cmd>` | Execute a Windows Command Prompt command on all nodes |
| `join <target> [alias]` | Join and enroll a machine into the cluster registry (`add`, `enroll`) |
| `nodes` / `ls` / `list` | List all registered machines joined to the cluster (`joined`, `machines`) |
| `remove <target>` / `rm` | Remove a machine from the cluster registry (`delete`) |
| `ping` | Check reachability and ping latency for all cluster nodes (`health`) |
| `install <pkgs>` | Install packages (comma-separated list) on all nodes |
| `pull --all` | Run git pull --all on all nodes |
| `push --all` | Run git push --all on all nodes |
| `commit --all` | Run git commit --all on all nodes |
| `status --all` | Show combined dirty/clean status across all nodes |
| `proj <name> run` | Run project-level automation on all nodes |

## Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet |
| `--ip <list>` | Target specific node IP addresses |
| `--id <list>` | Target specific node Display IDs |
| `--json` | Output machine list or results in JSON format |
| `--yes`, `-Y` | Bypass preflight confirmation prompt |
| `--dry-run` | Preview execution plan without running commands |
| `--verbose` | Print detailed real-time execution logs |

## Examples

### Bash & POSIX Shell Execution

Execute bash or shell commands across Linux and macOS cluster nodes:
```bash
# Check system kernel, release, and uptime across all nodes
gitmap servers-clients bash "uname -a && uptime"

# Check disk space across all machines, excluding node 2
gitmap sc shell "df -h" --except 2

# Inspect running Docker containers on all cluster workers
gitmap sc bash "docker ps --format 'table {{.Names}}\t{{.Status}}'"

# Execute quick POSIX shell commands
gitmap sc sh "cat /etc/os-release | grep PRETTY_NAME"
```

### Machine Joining & Enrollment (Join)

Join and enroll target machines into the cluster:
```bash
# Join a machine by IP address with custom alias 'u1'
gitmap servers-clients join 192.168.1.14 u1

# Join using SSH username, IP, and key authentication
gitmap sc join alim@192.168.1.14 u1 --key ~/.ssh/id_rsa

# Join with root credentials and password
gitmap sc join root@192.168.1.50 worker-1 --password secret

# Enroll an additional node using the add alias
gitmap sc add 10.0.0.5 dev-box
```

### Listing Joined Machines (List / Nodes / LS)

Inspect all registered machines currently joined to the cluster:
```bash
# View table of joined nodes (ID, ALIAS, IP, USER, ROLE, CREATED_AT)
gitmap servers-clients nodes
gitmap sc ls
gitmap sc list
gitmap sc joined

# Output joined machines in JSON format for scripts and automation
gitmap sc ls --json
gitmap servers-clients list --json
```

### Machine Removal (Remove / RM)

Remove a decommissioned machine from the cluster registry:
```bash
# Remove machine by alias
gitmap sc rm u1

# Remove machine by IP address
gitmap servers-clients remove 192.168.1.14
```

### Node Reachability & Health Status (Ping)

Check reachability and round-trip ping latency across all nodes:
```bash
gitmap servers-clients ping
gitmap sc ping
```

### PowerShell & Windows Command Prompt

Execute PowerShell and cmd commands on Windows nodes:
```bash
# Query active services via PowerShell
gitmap servers-clients ps "Get-Service | Where Status -eq Running"

# Check IP configuration via Windows Command Prompt, excluding specific nodes
gitmap sc cmd "ipconfig /all" --except 24,151

# Query top 5 CPU-intensive processes
gitmap sc ps "Get-Process | Sort-Object CPU -Descending | Select-Object -First 5"
```

### Package Installation

Install software packages across all machines simultaneously:
```bash
# Install common developer tools across cluster
gitmap servers-clients install "git,nodejs,curl,htop" --except 2
```

### Distributed Git Operations

Delegate git synchronization across all joined machines:
```bash
# Pull latest changes on all repositories
gitmap servers-clients pull --all

# Push branches across all nodes, excluding node 3
gitmap sc push --all --except 3

# Check working tree cleanliness across cluster
gitmap sc status --all
```

### Project Automation

Run local build scripts and CI pipelines across nodes:
```bash
gitmap servers-clients proj "api-backend" run --except 2
```

See also: `gitmap sc`, `gitmap clients`, `gitmap cluster`, `gitmap cluster exec`, `gitmap cluster nodes`
