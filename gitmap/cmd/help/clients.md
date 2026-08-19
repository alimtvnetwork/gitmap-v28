# gitmap clients

Broadcast shell commands, git operations, or lifecycle actions across all **client** nodes in the cluster (excludes the server/orchestrator).

## Usage

```
gitmap clients <subcommand> [args] [flags]
```

## Subcommands

Supports all standard execution commands (`ps`, `cmd`, `install`, `pull`, `push`, `status`, `proj`) plus lifecycle commands:
| Subcommand | Description |
|------------|-------------|
| `restart` | Gracefully restart the OS |
| `shutdown` | Initiate a full power-off |
| `logoff` | Sign out the current user session |

## Target Filtering Flags

| Flag | Description |
|------|-------------|
| `--except <list>` | Exclude nodes by ID, IP, or trailing IP octet |
| `--ip <list>` | Run *only* on the named IPs |
| `--id <list>` | Run *only* on the named Display IDs |

*Note: `--except`, `--ip`, and `--id` are mutually exclusive.*

## Examples

### Targeting Specific Clients
```
gitmap clients ps "Get-DiskUsage C:\" --except 24,151
gitmap clients cmd "dir C:\Projects" --id 1,2,3
```

### Chaining Commands
Execute multiple sub-commands sequentially on a target:
```
gitmap clients cmd "whoami", ps "Get-Date" --ip 192.168.1.10,192.168.1.11
gitmap clients ps "Set-ExecutionPolicy Bypass", install "git,nodejs" --id 1,3
```

### Machine Lifecycle
Restart or shutdown machines (requires `--force-lifecycle` and stored password):
```
gitmap clients restart --except 1
gitmap clients shutdown --id 5
```
