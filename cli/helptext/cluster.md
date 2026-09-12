# gitmap cluster

Manage multi-machine clusters, nodes, command history, and credential configurations.

## Usage

```bash
gitmap cluster <subcommand> [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `status` | Show cluster health, active nodes, and connectivity |
| `nodes` | List registered nodes in the cluster |
| `history` | Inspect audit trail of past cluster runs |
| `export` | Export node registry to JSON or CSV |
| `import` | Import node registry from JSON or CSV |
| `set-password` | Configure bcrypt password for privileged lifecycle ops |
| `set-default-path` | Set default working directory for a node |
| `set-path-alias` | Map path aliases across cluster nodes |
| `cat` | Read file from remote cluster node |
| `write` | Write file content to remote cluster node |
| `update` | Update gitmap binary on specific cluster node |
| `update-all` | Update gitmap binary across all cluster nodes |

## Examples

```bash
gitmap cluster status
gitmap cluster nodes --json
gitmap cluster history RUN-20260817-001
gitmap cluster export --format json --output nodes.json
```

See also: `gitmap servers-clients`, `gitmap clients`
