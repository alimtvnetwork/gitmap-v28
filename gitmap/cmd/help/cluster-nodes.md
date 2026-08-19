# gitmap cluster nodes

List registered nodes in the cluster.

## Usage

```
gitmap cluster nodes [--json]
```

## Description

Displays the list of nodes currently registered in the server's database.

## Table Columns

| Column | Description |
|--------|-------------|
| ID | Sequential 1-based integer Display ID |
| ALIAS | Human-readable alias (usually hostname) |
| IP | IP address of the node |
| ROLE | `server` or `client` |
| OS | `windows`, `linux`, or `darwin` |
| STATUS | `online`, `offline`, or `unreachable` |

## Flags

- `--json`: Emit the listing as a machine-readable JSON array of node objects.

## Examples

```
gitmap cluster nodes
gitmap cluster nodes --json
```
