# gitmap cluster nodes

List registered nodes in the cluster database.

## Usage

```bash
gitmap cluster nodes [flags]
gitmap cluster ls [flags]
gitmap cluster node ls
```

## Description

Displays all registered nodes in the cluster SQLite database (`ssh_hosts` and `cluster_nodes` tables), including display ID, host alias, network IP address, operating system, role (`control`, `worker`, `server`, `client`), and current reachability status. If any nodes are offline or unreachable, a warning banner is automatically displayed.

## Table Columns

| Column | Description |
|--------|-------------|
| DisplayId | Sequential integer display ID |
| Alias | Hostname or memorable node label |
| IP | Network IP address of the node |
| OS | Operating system (`linux`, `windows`, `darwin`) |
| Role | Node role (`control`, `worker`, `server`, `client`) |
| Status | Connectivity status (`online`, `offline`, `unreachable`) |
| LastHeartbeat | Timestamp of the last recorded heartbeat |

## Flags

- `--json`: Emit node registry as structured JSON array.
- `-h`, `--help`: Show command help.

## Examples

```bash
# List all registered nodes in formatted ASCII table
gitmap cluster nodes

# Use short alias
gitmap cluster ls

# Output node registry as JSON array for scripting
gitmap cluster nodes --json

# List nodes via node subcommand
gitmap cluster node ls
```

See also: `gitmap cluster add`, `gitmap cluster status`, `gitmap cluster`
