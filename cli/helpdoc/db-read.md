# gitmap db-read

Read and inspect SQLite databases across Split-DB tiers with safe read-only queries and schema reflection.

## Aliases

`dbr`

## Usage

```bash
gitmap db-read <path-or-alias> "sql" [flags]
gitmap db-read <path-or-alias> ls [flags]
gitmap db-read <path-or-alias> table <t1> [t2...] [flags]
gitmap db-read <path> -as <aliasName> [flags]
```

## Description

`gitmap db-read` provides safe, read-only inspection of SQLite databases across GitMap's Split-DB tiers (`master`, `pipeline`, `tasks`, `heatmap`, `history`) and custom SQLite database paths.

The command enforces strict read-only execution: any mutation attempts (`INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`) are rejected with `ERR_READONLY_VIOLATION`. For write operations, use `gitmap db-add` or `gitmap db-edit`.

## Flags

| Flag | Default | Description |
|---|---|---|
| `-as, --alias <name>` | master | Register or target database alias name |
| `-l, --limit, -n <N>` | 100 | Maximum rows to return for query |
| `-c, --columns` | false | List columns, types, notnull, default, and primary key |
| `-j, --json` | false | Output results in structured JSON envelope |
| `-h, --help` | false | Show command help text |

## Operations

### SQL Query Execution
Pass any valid SQLite `SELECT`, `PRAGMA`, or `EXPLAIN` statement in quotes. Queries are automatically wrapped with limit constraints if `--limit` is specified.

### Table Listing (`ls`)
Passing `ls` lists all user tables and row counts within the target database.

### Table Schema Inspection (`table`)
Passing `table <t1> [t2...]` outputs column definitions and schema details. Use `--columns` / `-c` for detailed column properties.

### Alias Registration (`-as`)
Passing a database file path with `-as <aliasName>` registers the database path under the specified alias for subsequent commands.

## Examples

### Run a read-only SQL query

```bash
gitmap db-read master "SELECT Id, Name, Path FROM Repository LIMIT 5"
```

### List all tables in a split database

```bash
gitmap dbr tasks ls
```

### Inspect table schema and columns

```bash
gitmap db-read tasks table Subtask --columns
```

### Register a database file alias

```bash
gitmap db-read ./custom.db -as custom
```

### Query with limit and JSON output

```bash
gitmap db-read tasks "SELECT * FROM Subtask" --limit 10 --json
```

## See Also

- `gitmap db-add` — Insert records into SQLite tables
- `gitmap db-edit` — Update existing records by primary key ID
- `gitmap db-find` — Search table columns by wildcard or regex
- `gitmap db-listing` — Discover and inspect registered SQLite databases
- `gitmap db-view` — Query multi-table joins and filtered views

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-read
```
