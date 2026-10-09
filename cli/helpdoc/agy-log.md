# agy log / agy logs

Inspect Antigravity decision audit logs and operational traces recorded in SQLite Split-DB (`gitmap-agy-log.db`).

## Synopsis

```bash
gitmap agy log [flags]
gitmap agy logs [flags]
```

## Options

- `--json`: Output decision logs in structured JSON format.
- `--ssh`: Query and aggregate decision logs from remote cluster nodes (e.g. `w1`).
- `-n, --limit <N>`: Maximum number of log entries to retrieve (default: 20).
- `-P, --project <name>`: Filter logs by target project name or directory path.
- `-c, --command <cmd>`: Filter logs by command (`rerun`, `ls`, `inject`, etc.).

## Examples

```bash
# View the latest 20 AGY decision logs
gitmap agy log

# Filter decision logs for a specific project
gitmap agy log -P white-presentation-v1

# Query logs across cluster nodes including w1
gitmap agy log --ssh -n 50

# Output decision traces as JSON
gitmap agy log --json
```
