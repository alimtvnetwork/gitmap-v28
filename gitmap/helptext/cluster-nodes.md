# gitmap cluster nodes

List all registered nodes in the cluster database.

## Usage

```bash
gitmap cluster nodes [--json]
```

## Table Columns

| Column | Description |
|--------|-------------|
| ID | Sequential integer Display ID |
| ALIAS | Hostname or node label |
| IP | Network IP address of the node |
| ROLE | `server` or `client` |
| OS | `windows`, `linux`, or `darwin` |
| STATUS | `online`, `offline`, or `unreachable` |

## Flags

- `--json`: Emit node registry as structured JSON array.

## Examples

```bash
gitmap cluster nodes
gitmap cluster nodes --json
```

See also: `gitmap cluster`, `gitmap servers-clients`
