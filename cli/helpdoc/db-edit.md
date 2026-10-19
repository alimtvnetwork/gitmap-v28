# gitmap db-edit

Update a record in a SQLite database table by primary key identifier.

## Aliases

`dbe`

## Usage

```bash
gitmap db-edit <alias-or-path> -table <table> -id <id> csv <col1>, <col2>, ...
gitmap db-edit <alias-or-path> -table <table> -id <id> "col1=val1, col2=val2, ..."
```

## Description

`gitmap db-edit` updates records in SQLite database tables across Split-DB tiers (`master`, `pipeline`, `tasks`, `heatmap`, `history`) and custom SQLite databases.

The command automatically resolves the primary key column for the target table (such as `Id`, `RowID`, or the first primary key in table schema) and applies updates to the record matching the specified `-id`.

Updates can be specified as positional comma-separated values following the `csv` token, or as explicit quoted key-value pairs (`"col1=val1, col2=val2"`).

## Flags and Arguments

| Flag / Parameter | Default | Description |
|---|---|---|
| `<alias-or-path>` | (required) | Database alias (`master`, `tasks`, `pipeline`, etc.) or SQLite file path |
| `-table, --table, -t <table>` | (required) | Target table name to update |
| `-id, --id <id>` | (required) | Primary key value of the row to modify |
| `csv <col1>, ...` | (optional) | Comma-separated positional values to apply |
| `"col1=val1, ..."` | (optional) | Quoted comma-separated column assignments |
| `-h, --help` | false | Show command help text |

## Update Formats

### Positional CSV Values
Supply `csv` followed by comma-separated values in column declaration order (skipping the primary key column):

```bash
gitmap db-edit tasks -table Subtask -id 42 csv "Updated title", "DONE"
```

### Key-Value Assignments
Supply quoted comma-separated `column=value` pairs to update specific fields without affecting unmentioned columns:

```bash
gitmap db-edit tasks -table Subtask -id 42 "Status=DONE, Priority=HIGH"
```

## Examples

### Update row using positional CSV format

```bash
gitmap db-edit tasks -table Subtask -id 42 csv "Updated title", "DONE"
```

### Update row using key-value assignments

```bash
gitmap db-edit tasks -table Subtask -id 42 "Status=DONE, Priority=HIGH"
```

### Use short alias and flags

```bash
gitmap dbe tasks -t Subtask -id 10 csv "Refactored module"
```

### Update record in a custom SQLite file path

```bash
gitmap db-edit ./custom.db -t Users -id 5 "Email=alice@newdomain.com"
```

## See Also

- `gitmap db-read` — Read and inspect SQLite databases
- `gitmap db-add` — Insert records into SQLite tables
- `gitmap db-find` — Search table columns by wildcard or regex
- `gitmap db-listing` — Discover and inspect registered SQLite databases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-edit
```
