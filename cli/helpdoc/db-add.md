# gitmap db-add

Insert records into a SQLite database table using interactive prompts or comma-separated values.

## Aliases

`dba`

## Usage

```bash
gitmap db-add <alias-or-path> <table>
gitmap db-add <alias-or-path> <table> csv <val1>, <val2>, ...
gitmap db-add <alias-or-path> <table> col1=val1, col2=val2, ...
```

## Description

`gitmap db-add` inserts new rows into tables across any registered database tier (`master`, `pipeline`, `tasks`, `heatmap`) or custom SQLite file path.

The command supports three input workflows:
1. **Interactive Prompt Mode**: Launched when invoked without value parameters in an interactive terminal. Prompts for each table column sequentially.
2. **CSV Insertion Mode**: Fast positional insertion by providing comma-separated values prefixed by `csv`.
3. **Key-Value Mode**: Explicit named assignments (`col1=val1, col2=val2`) for partial or unordered column insertions.

## Arguments

| Argument | Description |
|---|---|
| `<alias-or-path>` | Target database alias (`master`, `tasks`, `pipeline`, etc.) or SQLite file path |
| `<table>` | Target table name |
| `csv <val1>, ...` | Positional comma-separated column values to insert |
| `<col1>=<val1>, ...` | Explicit key-value pairs mapping column names to values |

## Flags

| Flag | Default | Description |
|---|---|---|
| `-h, --help` | false | Show command help text |

## Interactive Prompt Mode

When invoked with only the database alias and table name in an interactive terminal session:

```bash
gitmap db-add tasks Subtask
```

The CLI prompts for each non-generated column:
- Column names and SQL data types are displayed.
- Defaults and nullable states are respected (press Enter to use default or null).
- Auto-incrementing primary keys are skipped automatically.

## Examples

### Insert record using CSV values

```bash
gitmap db-add tasks Subtask csv "Fix bug", "IN_PROGRESS", 1
```

### Insert record using key-value assignments

```bash
gitmap db-add tasks Subtask title="Fix bug", status="IN_PROGRESS"
```

### Launch interactive prompt mode using alias

```bash
gitmap dba tasks Subtask
```

### Insert record into custom SQLite file path

```bash
gitmap db-add ./custom.db Users csv "Alice", "alice@example.com"
```

## See Also

- `gitmap db-read` — Read and inspect SQLite databases
- `gitmap db-edit` — Update existing records by primary key ID
- `gitmap db-find` — Search table columns by wildcard or regex
- `gitmap db-listing` — Discover and inspect registered SQLite databases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-add
```
