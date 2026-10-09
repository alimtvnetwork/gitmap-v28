# gitmap cluster remove

Remove and unregister a node from the cluster topology and SSH host registry.

## Usage

```bash
gitmap cluster remove <alias|ip> [flags]
gitmap cluster rm <alias|ip> [flags]
gitmap cluster node rm <alias|ip> [flags]
gitmap cluster node remove <alias|ip> [flags]
```

## Description

`gitmap cluster remove` (or `rm`) removes a node from the local cluster database and SSH host registry. Once removed, the node will no longer be targeted by multi-node execution commands (`exec`, `run-script`, `node`, `k8s`).

Deletion can be performed by host alias, network IP address, or internal node ID.

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--id <id>` | | Target database node ID for cluster DB deletion |
| `--confirm` | | Confirmation flag required when deleting node via `--id` |
| `--force` | `-f` | Remove node without confirmation prompts |
| `--help` | `-h` | Show command help |

## Examples

```bash
# 1. Remove node by alias
gitmap cluster rm k8s-w1

# 2. Remove node using cluster remove command
gitmap cluster remove k8s-w1

# 3. Remove node using cluster node rm syntax
gitmap cluster node rm k8s-w2

# 4. Remove node by IP address
gitmap cluster rm 192.168.0.103

# 5. Remove node by database ID with explicit confirmation
gitmap cluster remove --id node-123 --confirm

# 6. Force remove node without confirmation prompt
gitmap cluster rm old-worker -f
```

See also: `gitmap cluster add`, `gitmap cluster nodes`, `gitmap cluster`
