# gitmap db-record-del

Delete a record from a SQLite database table by primary key.

## Aliases

`dbrd`, `db-record-delete`, `db record-del`

## Usage

```bash
gitmap db-record-del <db> -table <table> -id <id> [-y] [--json]
gitmap db-record-del <db> <table> <id> [-y] [--json]
```

## Description

`gitmap db-record-del` deletes a single row from the designated table using its resolved primary key column. The command dynamically determines the primary key column name and executes the deletion inside an atomic transaction.

## Flags

| Flag | Description |
|---|---|
| `-t, --table <table>` | Target table name |
| `--id <id>` | Primary key value of the row to delete |
| `-y, --yes, --force` | Bypass confirmation prompt |
| `-j, --json` | Output result in structured Universal Response Envelope JSON |
| `-h, --help` | Show command help text |

## Examples

### Delete record with confirmation prompt
```bash
gitmap db-record-del tasks -table Subtask -id 42
```

### Bypass confirmation prompt
```bash
gitmap db-record-del tasks -table Subtask -id 42 -y
```

### Positional arguments with JSON envelope output
```bash
gitmap dbrd master Repositories 101 -y --json
```
