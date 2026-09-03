# Reconcile Projects (`gitmap agy reconcile`)

Reconciles missing Antigravity projects by cross-referencing their names and directory slugs against GitMap's repository index and database. When a project folder has moved or was relocated to another directory, `reconcile` automatically discovers the new location and updates the project's workspace JSON.

## Synopses & Aliases

```bash
gitmap agy reconcile [flags]
gitmap agy recon [flags]
gitmap agy reconcile-projects [flags]
```

## Flags

| Flag | Shorthand | Type | Default | Description |
|---|---|---|---|---|
| `--dry-run` | `-d` | bool | `false` | Preview discovered path re-links without saving |
| `--yes` | `-y` | bool | `false` | Apply re-linked paths without interactive prompt |

## Examples

### 1. Reconcile missing projects (standard execution)
```bash
gitmap agy reconcile
```
**Alias:**
```bash
gitmap agy recon
```

### 2. Preview path reconciliations safely (dry run)
```bash
gitmap agy recon --dry-run
```

### 3. Non-interactive automatic reconciliation
```bash
gitmap agy recon -y
```

## Related Recovery Commands

When projects are unresolved (permanently deleted):
- **Purge missing projects:** `gitmap agy remove-missing-projects` (alias: `rm-missing`)
- **Scan directory tree:** `gitmap agy scan [path]`
