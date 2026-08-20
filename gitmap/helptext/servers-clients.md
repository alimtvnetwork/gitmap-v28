# gitmap servers-clients

Broadcast shell commands, git operations, or installations across all server and client nodes in the cluster.

## Usage

```bash
gitmap servers-clients <subcommand> [args] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ps <cmd>` | Execute a PowerShell command on all nodes |
| `cmd <cmd>` | Execute a Windows Command Prompt command on all nodes |
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
| `--yes`, `-Y` | Bypass preflight confirmation prompt |
| `--dry-run` | Preview execution plan without running commands |
| `--verbose` | Print detailed real-time execution logs |

## Examples

### PowerShell & Command Prompt
Execute PowerShell and cmd commands across the cluster:
```bash
gitmap servers-clients ps "Get-Service | Where Status -eq Running"
gitmap servers-clients cmd "ipconfig /all" --except 24,151
```

### Installation
Install packages across all nodes simultaneously:
```bash
gitmap servers-clients install "git,nodejs,dotnet" --except 2
```

### Git Operations
Delegate git operations to all joined machines:
```bash
gitmap servers-clients pull --all
gitmap servers-clients push --all --except 3
gitmap servers-clients status --all
```

### Project Automation
Run local build scripts across nodes:
```bash
gitmap servers-clients proj "api-backend" run --except 2
```

See also: `gitmap sc`, `gitmap clients`, `gitmap cluster`
