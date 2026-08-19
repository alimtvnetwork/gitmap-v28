# gitmap servers-clients

Broadcast shell commands, git operations, or installations across all server and client nodes in the cluster.

## Usage

```
gitmap servers-clients <subcommand> [args] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ps` | Execute a PowerShell command |
| `cmd` | Execute a Windows Command Prompt command |
| `install` | Install packages (comma-separated list) |
| `pull` | Run git pull --all on all nodes |
| `push` | Run git push --all on all nodes |
| `commit` | Run git commit --all on all nodes |
| `status` | Show combined dirty/clean status across nodes |
| `proj` | Run project-level automation (run, create-cicd) |

## Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet |
| `--yes`, `-Y` | Bypass preflight confirmation prompt |

## Examples

### PowerShell & Command Prompt
Execute PowerShell and cmd commands across the cluster:
```
gitmap servers-clients ps "Get-Service | Where Status -eq Running"
gitmap servers-clients cmd "ipconfig /all" --except 24,151
```

### Installation
Install packages simultaneously:
```
gitmap servers-clients install "git,nodejs,dotnet" --except 2
```

### Git Operations
Delegate git operations to all joined machines:
```
gitmap servers-clients pull --all
gitmap servers-clients push --all --except 3
gitmap servers-clients status --all
```

### Project Automation
Run the local `run.ps1` or `run.sh` for specific projects:
```
gitmap servers-clients proj "api-backend,web-frontend" run
gitmap servers-clients proj "api-backend" run --except 2
```

See also: `gitmap sc`, `gitmap clients`
