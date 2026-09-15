# gitmap sc

Short alias for `servers-clients`. Broadcasts shell commands, git operations, machine joins, or installations across all server and client nodes in the cluster.

## Usage

```bash
gitmap sc <subcommand> [args] [flags]
```

Full command: `gitmap servers-clients <subcommand> [args] [flags]`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `bash <cmd>` | Execute a Bash command across all cluster nodes |
| `sh <cmd>` / `shell <cmd>` | Execute a POSIX shell command across all cluster nodes |
| `ps <cmd>` | Execute a PowerShell command on all nodes |
| `cmd <cmd>` | Execute a Windows Command Prompt command on all nodes |
| `join <target> [alias]` | Join and enroll a machine into the cluster registry (`add`) |
| `nodes` / `ls` / `list` | List all registered machines joined to the cluster (`joined`, `machines`) |
| `remove <target>` / `rm` | Remove a machine from the cluster registry (`delete`) |
| `ping` | Check reachability and ping latency for all cluster nodes (`health`) |
| `install <pkgs>` | Install packages (comma-separated list) on all nodes |
| `pull --all` | Run git pull --all on all nodes |
| `push --all` | Run git push --all on all nodes |
| `commit --all` | Run git commit --all on all nodes |
| `status --all` | Show combined dirty/clean status across all nodes |
| `proj <name> run` | Run project-level automation on all nodes |

## Examples

### Bash & Shell Commands
```bash
gitmap sc bash "uname -a && uptime"
gitmap sc shell "df -h" --except 2
gitmap sc sh "docker ps"
```

### Join & Manage Machines
```bash
gitmap sc join 192.168.1.14 u1
gitmap sc join root@192.168.1.50 worker-1 --password secret
gitmap sc nodes
gitmap sc ls --json
gitmap sc rm u1
gitmap sc ping
```

### PowerShell & Windows Cmd
```bash
gitmap sc ps "Get-Date"
gitmap sc cmd "whoami" --except 2
gitmap sc ps "Get-Service | Where Status -eq Running"
```

### Git & Package Automation
```bash
gitmap sc pull --all
gitmap sc push --all --except 3
gitmap sc status --all
gitmap sc install "git,nodejs,curl" --except 2
```

See also: `gitmap servers-clients`, `gitmap clients`, `gitmap cluster`
