# storage

Display disk drives, storage volumes, mount points, filesystem space utilization, repository SQLite databases, and restore database from snapshots.

## Synopsis

```bash
gitmap storage [subcommand] [flags]
gitmap disk [subcommand] [flags]
gitmap df [subcommand] [flags]
```

## Aliases & Shorthands

- `gitmap storage`
- `gitmap disk`
- `gitmap df`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `status (st, info)` | Inspect drive capacity, filesystem metrics, and SQLite database summary |
| `ls (list, db, dbs)` | List all system, split, and repository SQLite databases (repo DBs) |
| `space ls (space list)` | Direct shortcut to list all repository SQLite databases and space usage |
| `restore-db (restore, restoredb)` | Restore or auto-heal Gitmap SQLite databases from backup snapshot or schema |
| `clean (clear, prune)` | Clean pipeline error logs, temporary files, and vacuum SQLite databases |
| `help` | Show this storage command suite documentation |

---

## Options & Flags

### General Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--all` | `-a` | boolean | `false` | Include hidden, system, and virtual drives |
| `--json` | `-j` | boolean | `false` | Output storage statistics in structured JSON format |
| `--help` | `-h` | boolean | `false` | Show help for storage command |

### Restore Flags (`restore-db`)

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--force` | `-f` | boolean | `false` | Bypass confirmation when restoring database from snapshot |
| `--snapshot` | `-s` | string | `""` | Snapshot ID or relative index to restore (defaults to latest) |

### Clean Flags (`clean`)

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--dry-run` | `-n` | boolean | `false` | Simulate cleanup without deleting any files |
| `--verbose` | `-v` | boolean | `false` | Print each deleted file path and reclaimed bytes |
| `--force` | `-f` | boolean | `false` | Bypass interactive confirmations |
| `--vacuum` | `-a` | boolean | `false` | Vacuum and reclaim unused pages from SQLite databases |

---

## Repository SQLite Databases (`repo DBs`)

GitMap employs an isolated split-database architecture. Running `gitmap storage ls` or `gitmap storage space ls` outputs a structured inventory table across all managed SQLite databases:

| Column | Description |
|--------|-------------|
| `NAME` | Database identifier, split slug, or repository name |
| `TYPE` | Database classification: `master`, `split`, `repo`, or `custom` |
| `SIZE` | File size on disk formatted in human-readable bytes (KB / MB) |
| `TABLES` | Total count of SQLite tables defined in schema |
| `RECORDS` | Total record count aggregated across all tables |
| `PATH` | Relative path from GitMap data root |

### Database Categories

- **Master Database (`master`)**: `gitmap.db` stores primary repository paths, tracking configurations, and global metadata.
- **Split Databases (`split`)**: Dedicated domain databases such as `installation.db` (installer registry), `startup.db` (autostart manager), and `sites.db` (virtual hosts).
- **Subsystem Split DBs**: Isolated databases for task schedules (`schedules/<id>.db`) and pipeline telemetry (`pipeline/<repo>.db`).
- **Repository-Scoped Databases (`repo`)**: Local `.gitmap/gitmap.db` databases stored inside individual tracked repositories for per-project telemetry.

---

## Database Restore Engine (`restore-db`)

The `restore-db` command recovers GitMap SQLite databases from automated snapshot archives:

- **Auto-Healing**: Verifies schema integrity and repairs broken SQLite files using transactional schemas.
- **Snapshot Selection**: Restores from the most recent backup snapshot by default, or an explicit snapshot specified via `--snapshot <ID>`.
- **Zero Data Loss**: Safely creates a pre-restore backup before applying snapshot data.

---

## Examples

```bash
# View storage consumption across all physical volumes and SQLite DB summary
gitmap storage

# List all system, split, and repository SQLite databases (repo DBs)
gitmap storage ls

# List repository SQLite databases via storage space ls
gitmap storage space ls

# Inspect detailed storage metrics for a specific path or volume
gitmap storage status /var/lib

# Restore database from the latest cloud/local backup snapshot
gitmap storage restore-db

# Restore database from a specific snapshot ID with force confirmation bypass
gitmap storage restore-db --snapshot 20260918-120000 --force

# Clean old pipeline logs and temp files (dry-run simulation)
gitmap storage clean --dry-run

# Force clean pipeline logs and temp files with verbose output
gitmap storage clean --force --verbose

# Clean temporary logs and vacuum pipeline SQLite database
gitmap storage clean --force --vacuum
```

---

## See Also

- [agy](agy.md) — Google Antigravity workspaces, prompt templates, and rerun suite
- [pipeline](pipeline.md) — CI/CD pipeline status, logs, and database management
- [prompts-template](prompts-template.md) — Reusable prompt prefix and verification templates
