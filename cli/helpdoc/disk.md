# gitmap disk

Inspect filesystem drive partitions, total capacity, available space, and usage percentages.

## Synopsis

```bash
gitmap disk [flags]
gitmap storage [subcommand] [flags]
gitmap os disk [flags]
```

## Description

Surfaces drive letter / mount points, filesystem types (NTFS, ext4, APFS), total capacity, free space, and usage percentages across all mounted volumes. Implements cross-process caching to prevent disk I/O contention.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output drive partitions and usage metrics as structured JSON |
| `--no-cache` (`-b`) | `false` | Bypass telemetry cache and execute live hardware probe |
| `--all` (`-a`) | `false` | Include hidden, system, and virtual drives |

## Examples

```bash
# Display disk drive capacity and usage percentages
gitmap disk

# Output structured JSON disk telemetry
gitmap disk --json

# Force live query bypassing cache
gitmap disk --no-cache
```

## See Also

- [cpu](cpu.md) — Query CPU utilization, core count, and hardware model
- [ram](ram.md) — Query system memory capacity and active utilization
- [storage](storage.md) — Manage storage volumes, backups, and SQLite databases
