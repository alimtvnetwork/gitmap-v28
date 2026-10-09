# gitmap db

Inspect, manage, and reset Gitmap's primary SQLite database and split databases.

## Usage

```bash
gitmap db [subcommand] [flags]
```

## Subcommands

| Subcommand | Description |
|---|---|
| `ls`, `list` | Show all database paths, sizes, split DB list, and architectural purposes |
| `repo-db list` | Show per-repository split database metrics and tracking status |
| `sizes list` | Show disk usage breakdown for all database files |
| `reset`, `clear` | Reset master database records and clear split DBs |
| `help` | Display command help |

## Flags

| Flag | Default | Description |
|---|---|---|
| `--yes`, `-y` | false | Bypass interactive confirmation prompt |
| `--force`, `-f` | false | Alias for `--yes` |

## Examples

### View Database Architecture and Split DBs

```bash
gitmap db ls
```

**Output:**

```text
  ── Gitmap Database Architecture ──

  ● Primary Master Database:
    Name:        gitmap.db
    Type:        Primary Master DB
    Location:    C:\Users\Alim\bin\data\gitmap.db
    Size:        128.0 KB
    Description: Central SQLite database storing all global tracked repositories.
```

### View Database Sizes Breakdown

```bash
gitmap db sizes list
```

**Output:**

```text
  ── Gitmap Database Disk Sizes ──

  DATABASE FILE             CATEGORY      SIZE         PATH
  ──────────────────────────────────────────────────────────────────────────────────
  gitmap.db                 Master DB     128.0 KB     C:\Users\Alim\bin\data\gitmap.db
```

## See Also

- `gitmap start-fresh`: Clear all databases and rebuild clean schemas
- `gitmap scan`: Scan disk and populate the database
