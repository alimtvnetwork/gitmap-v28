# gitmap db-find

Search across SQLite table columns using wildcard patterns or regular expressions with column scoping.

## Aliases

`dbf`

## Usage

```bash
gitmap db-find <alias-or-path> -table <t1,t2> --find/-f <text> [flags]
gitmap db-find <alias-or-path> -table <t1,t2> --grep/-r <regex> [flags]
```

## Description

`gitmap db-find` scans SQLite databases across Split-DB tiers (`master`, `pipeline`, `tasks`, `heatmap`, `history`) and custom SQLite databases for records matching search patterns.

The command supports:
- **Wildcard Matching** via `--find` (`-f`): supports prefix (`text*`), suffix (`*text`), full wildcard (`*text*`), and substring matching.
- **Regular Expressions** via `--grep` (`-r`): evaluates Go-compatible regex patterns against column values.
- **Table Scoping**: comma-separated list of tables or `*` to search all tables in the database.
- **Column Scoping**: restrict search to specific column names using `--columns` (`-c`).

## Flags

| Flag | Default | Description |
|---|---|---|
| `-table, --table, -t <t1,t2>` | (required) | Comma-separated table names to search (or `*` for all) |
| `--find, -f <text>` | "" | Substring or wildcard pattern (`text`, `text*`, `*text`) |
| `--grep, -r <regex>` | "" | Regular expression to match against column values |
| `--columns, -c <c1,c2>` | all | Restrict search to comma-separated column names |
| `--limit, -l, -n <N>` | 100 | Maximum matching rows per table |
| `--json, -j` | false | Output results in JSON format |
| `-h, --help` | false | Show command help text |

## Wildcard Patterns

The `--find` (`-f`) flag interprets pattern markers:
- `text*` — Matches values beginning with `text` (prefix query).
- `*text` — Matches values ending with `text` (suffix query).
- `*text*` — Matches values containing `text` anywhere.
- `text` — Standard substring match across target columns.

## Regex Search

The `--grep` (`-r`) flag compiles a standard regular expression and scans column values:
- `^E[0-9]+` — Match error codes starting with E followed by digits.
- `(?i)failed|fatal` — Case-insensitive alternation search.

## Examples

### Wildcard find across a table

```bash
gitmap db-find master -table Repository --find "gitmap*"
```

### Regex find across multiple tables with column scoping

```bash
gitmap db-find tasks -table Subtask,Plan --grep "^E[0-9]+" --columns Title,Code
```

### Suffix wildcard find with limit and JSON output

```bash
gitmap db-find pipeline -table Runs -f "*failed" --limit 25 --json
```

### Search all tables using short alias and flags

```bash
gitmap dbf tasks -t "*" -f "refactor"
```

## See Also

- `gitmap db-read` — Read and inspect SQLite databases
- `gitmap db-view` — Query multi-table joins and filtered views
- `gitmap db-listing` — Discover and inspect registered SQLite databases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-find
```
