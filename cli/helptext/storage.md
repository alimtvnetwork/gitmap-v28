# storage

Display disk drives, storage volumes, mount points, filesystem space utilization, and Gitmap SQLite databases.

## Aliases

`disk`, `df`

## Usage

```bash
gitmap storage [subcommand] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| status (st, info) | Inspect drive capacity, filesystem metrics, and SQLite database summary |
| ls (list, db, dbs) | List all repository-scoped and global Gitmap SQLite databases |
| clean (clear, prune) | Clean pipeline error logs, temporary files, and vacuum SQLite databases |

## Options

### General Flags

    -a, --all      Include hidden, system, and virtual drives
    -j, --json     Output storage statistics in JSON format

### Clean Flags

    -n, --dry-run  Simulate cleanup without deleting any files
    -v, --verbose  Print each deleted file path and reclaimed bytes
    -f, --force    Bypass interactive confirmations
    -a, --vacuum   Vacuum and reclaim unused pages from SQLite databases

## Examples

```bash
# View storage consumption across all physical volumes and SQLite DB count
gitmap storage

# List all SQLite databases tracked across Gitmap repositories
gitmap storage ls

# Inspect storage metrics for a specific volume or path
gitmap storage status /var/lib

# Clean old pipeline logs and temp files (dry-run simulation)
gitmap storage clean --dry-run

# Force clean pipeline logs and temp files with verbose output
gitmap storage clean --force --verbose

# Clean temporary logs and vacuum pipeline SQLite database
gitmap storage clean --force --vacuum
```

