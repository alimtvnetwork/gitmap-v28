# gitmap db

Inspect, manage, query, and reset Gitmap's primary SQLite database, Antigravity databases, and split databases.

## Usage

```bash
gitmap db [subcommand] [arguments...] [flags]
```

## Subcommands

| Subcommand | Description |
|---|---|
| `query "<SQL>"` | Execute read-only SQL queries (`SELECT`, `PRAGMA`, `EXPLAIN`, `WITH`) against any database or alias |
| `tables [db]` | List tables, views, column counts, and record counts |
| `schema [table]` | Dump `CREATE TABLE` and index DDL |
| `columns <table>` | Inspect column definitions (`PRAGMA table_info`) |
| `list <table>` | Display table rows with limit paging (`--limit N`) |
| `ls`, `list` | Show all database paths, sizes, split DB list, and architectural purposes |
| `repo-db list` | Show per-repository split database metrics and tracking status |
| `sizes list` | Show disk usage breakdown for all database files |
| `reset`, `clear` | Reset master database records and clear split DBs |
| `help` | Display command help |

## Flags

| Flag | Default | Description |
|---|---|---|
| `-d`, `--db <alias\|path>` | master | Target database: `agy` (`conversation_summaries.db`), `master`, `heatmap`, `pipeline`, `tasks`, or path |
| `-j`, `--json` | false | Output query results in structured JSON format |
| `-n`, `--limit <N>` | 20 | Maximum rows to return for `list` |
| `--yes`, `-y` | false | Bypass interactive confirmation prompt |
| `--force`, `-f` | false | Alias for `--yes` |

## System Database Aliases

| Alias | Database Path | Description |
|---|---|---|
| `agy`, `antigravity`, `convs` | `~/.gemini/antigravity/conversation_summaries.db` | Antigravity conversation summaries & IDE state |
| `master`, `gitmap` | `%LOCALAPPDATA%/gitmap-cli/data/gitmap.db` | Primary master repository & settings DB |
| `heatmap` | `.../data/heatmap.db` | AI agent heatmap & coordination DB |
| `pipeline` | `.../data/pipeline/.../sql.db` | CI/CD pipeline runs & telemetry DB |
| `tasks`, `agents` | `.../data/tasks/sql.db` | Master agent swarm task queue DB |

## Examples

### Query Antigravity Conversation Summaries

```bash
gitmap db query "SELECT conversation_id, title, status FROM conversation_summaries ORDER BY last_modified_time DESC LIMIT 5" --db agy
```

### Inspect Database Tables & Schemas

```bash
gitmap db tables --db agy
gitmap db schema conversation_summaries --db agy
gitmap db columns conversation_summaries --db agy
```

### Query Primary Gitmap Database

```bash
gitmap db query "SELECT RepoId, AbsolutePath FROM TrackedRepository LIMIT 10" --db master
```

### View Database Architecture and Split DBs

```bash
gitmap db ls
```

### View Database Sizes Breakdown

```bash
gitmap db sizes list
```

## See Also

- `gitmap agy summaries [N]`: Dedicated Antigravity conversation summary inspector
- `gitmap start-fresh`: Clear all databases and rebuild clean schemas
- `gitmap scan`: Scan disk and populate the database
