# gitmap start-fresh

Permanently wipe all databases and caches, reinitializing clean schemas from scratch.

## Usage

```bash
gitmap start-fresh [flags]
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--yes`, `-y` | false | Bypass interactive irreversible confirmation prompt |
| `--force`, `-f` | false | Alias for `--yes` |

## Description

`gitmap start-fresh` performs a complete purge of all SQLite database files, journal files (`-wal`, `-shm`), and per-repository split databases in `repo_search/`. It then re-runs all schema migrations to create clean, empty tables ready for a new repository scan.

## Examples

### Interactive Start Fresh

```bash
gitmap start-fresh
```

**Output:**

```text
  ⚠ WARNING: Irreversible Database Transaction!
  This will permanently delete all tracked repositories, scan histories,
  search caches, profiles, and split databases across your entire system.

  Are you sure you want to start fresh? [y/N]: y

  ✓ Removed 4 previous database and cache file(s).
  ✓ Fresh master database initialized: C:\Users\Alim\bin\data\gitmap.db
  ✓ Rebuilt clean database schemas and migrations.
```

### Non-Interactive Start Fresh

```bash
gitmap start-fresh --yes
```

**Output:**

```text
  ✓ Removed 4 previous database and cache file(s).
  ✓ Fresh master database initialized: C:\Users\Alim\bin\data\gitmap.db
  ✓ Rebuilt clean database schemas and migrations.
```

## See Also

- `gitmap db`: Manage, inspect, and reset Gitmap databases
- `gitmap scan`: Scan directories to rebuild repository index
