# gitmap clients

Broadcast shell commands, git operations, or lifecycle actions across all client nodes in the cluster (excludes the server/orchestrator).

## Usage

```bash
gitmap clients <subcommand> [args] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ps <cmd>` | Run PowerShell command on client machines |
| `cmd <cmd>` | Run CMD command on client machines |
| `install <pkgs>` | Install packages across client machines |
| `pull --all` | Run git pull --all on client machines |
| `restart` | Reboot client machines (requires password) |
| `shutdown` | Power off client machines (requires password) |
| `logoff` | Log off active users from client machines |

## Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing octet |
| `--ip <list>` | Run only on specific IP addresses |
| `--id <list>` | Run only on specific Display IDs |
| `--force-lifecycle` | Required for restart, shutdown, logoff |
| `--yes`, `-Y` | Bypass preflight confirmation prompt |
| `--dry-run` | Preview command dispatch without executing |

## Examples

```bash
gitmap clients ps "Get-DiskUsage C:\" --except 24,151
gitmap clients cmd "whoami", ps "Get-Date" --ip 192.168.1.10,192.168.1.11
gitmap clients restart --except 1 --force-lifecycle -Y
```

See also: `gitmap servers-clients`, `gitmap cluster`
