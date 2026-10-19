# gitmap db-rm

Remove a database alias mapping and optionally purge physical SQLite files.

## Aliases

`dbrm`, `db-delete`, `db rm`

## Usage

```bash
gitmap db-rm <alias> [--purge] [-y]
```

## Description

`gitmap db-rm` unregisters an alias from `DatabaseAlias` in `gitmap.db`. System aliases (`master`, `pipeline`, `heatmap`, `tasks`) are safeguarded against accidental removal.

If `--purge` or `--purge-file` is passed, the physical SQLite database file and associated `-wal`, `-shm`, and `-journal` files are permanently deleted from disk after confirmation.

## Flags

| Flag | Description |
|---|---|
| `--purge, --purge-file` | Permanently delete underlying database file from disk |
| `-y, --yes, --force` | Bypass confirmation prompt |
| `-h, --help` | Show command help text |

## Examples

### Unbind alias without deleting file
```bash
gitmap db-rm customdb
```

### Unbind alias and purge physical database file
```bash
gitmap db-rm tempdb --purge -y
```
