# `gitmap start-fresh`

Permanently wipe all databases, cached indices, and split DBs, re-running all schema migrations from scratch.

## Usage

```bash
gitmap start-fresh [flags]
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--yes`, `-y` | `false` | Bypass interactive irreversible confirmation prompt |
| `--force`, `-f` | `false` | Alias for `--yes` |

## Safety Guarantee

`gitmap start-fresh` displays an explicit warning modal requiring user confirmation:

```text
⚠️ WARNING: IRREVERSIBLE TRANSACTION ⚠️
This operation will permanently delete all GitMap databases, cached indices,
and split search databases. Tracked repository paths will need to be re-scanned.
Are you sure you want to proceed? [y/N]:
```

Upon confirmation, it cleanly deletes `gitmap.db`, `-wal`, `-shm`, and recursively purges `repo_search/`, rebuilding empty tables ready for a fresh scan.
