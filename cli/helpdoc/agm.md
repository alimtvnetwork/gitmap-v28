# agm — Antigravity Manager GUI & Tools Management

Manage and install the Antigravity Manager GUI desktop application and tools locally or across remote SSH fleet nodes concurrently.

## Synopsis

```bash
gitmap agm <subcommand> [flags]
gitmap ag-manager <subcommand> [flags]
gitmap antigravity-manager <subcommand> [flags]
```

## Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `install` | `in`, `i` | Install Antigravity Manager GUI / Tools |
| `update` | `up`, `u` | Update Antigravity Manager locally or on remote node (`--remote`, `--ssh`) |
| `update-all` | `updateall`, `up-all` | Update Antigravity Manager across all nodes concurrently via SSH |

## Flags

### Update Flags (`gitmap agm update` / `gitmap agm update-all`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--ssh` | `-s` | `false` | Update remote fleet machines via SSH concurrently |
| `--remote` | `-r` | `""` | Target remote node alias or IP to update |
| `--except` | `-e` | `""` | Exclude machines by alias or IP (comma separated) |
| `--dry-run` | `-n` | `false` | Simulate update without downloading |
| `--yes` | `-y` | `false` | Automatic yes to prompts |
| `--verbose` | `-v` | `false` | Enable verbose output |

## Examples

```bash
# Update Antigravity Manager locally
gitmap agm update

# Update Antigravity Manager across all SSH fleet nodes concurrently
gitmap agm update-all ssh
gitmap agm update-all --ssh

# Update fleet nodes excluding specific workers
gitmap agm update-all ssh -e worker-2,192.168.1.15

# Update a single remote node
gitmap agm update --remote worker-1
```

## See Also

- `gitmap agy` — Antigravity CLI and prompt management
- `gitmap ssh` — SSH fleet execution and management
