# gitmap db-view

Query tables with automatic join inference, optional views, and condition filtering.

## Aliases

`dbv`

## Usage

```bash
gitmap db-view <alias-or-path> -t "<tables>" [flags]
```

## Description

`gitmap db-view` simplifies complex relational queries across SQLite tables in GitMap Split-DB tiers (`master`, `pipeline`, `tasks`, `heatmap`, `history`) or custom SQLite databases.

Key capabilities include:
- **Automatic Join Inference**: Given a comma-separated list of tables and aliases in `-t`, foreign keys and entity identifiers (e.g., `RepoId`, `PlanId`) are analyzed to generate proper `JOIN` clauses automatically.
- **Natural Filter Conditions**: Specify expressions via `--find` with support for equality, comparison operators, and wildcards.
- **View Materialization**: Persist the resolved SQL query as a SQLite `VIEW` using `--viewname`.
- **JSON Envelope**: Export structured query results with column metadata and execution timing in milliseconds.

## Flags

| Flag | Default | Description |
|---|---|---|
| `-t, --table, -table <spec>` | (required) | Table declarations with aliases (e.g. `"Repository as r, ScanHistory as s"`) |
| `--viewname, -v, --name <name>` | "" | Create a persistent named VIEW in the database |
| `--find, -f <condition>` | "" | Filter condition expression (e.g. `"r.IsActive = 1"`) |
| `--limit, -l <N>` | 100 | Maximum rows to return |
| `--json, -j` | false | Output result as structured JSON with execution timing |
| `-h, --help` | false | Show command help text |

## Multi-Table Join Inference

When multiple tables are declared in the `-t` parameter:

```bash
gitmap db-view master -t "Repository as r, ScanHistory as s"
```

The join engine:
1. Resolves foreign key references declared in the SQLite schema.
2. Matches conventional identifier column pairs (such as `r.Id = s.RepoId` or `r.Id = s.RepositoryId`).
3. Assembles the final `SELECT` query with joined projections.

## Filter Conditions

The `--find` (`-f`) flag evaluates conditions with support for:
- Comparisons: `=`, `!=`, `<`, `<=`, `>`, `>=`
- Pattern matching: `LIKE 'prefix*'`, `LIKE '*suffix'`
- Boolean conjunctions: `AND`, `OR`

## Examples

### Multi-table join query with filter condition

```bash
gitmap db-view master -t "Repository as r, ScanHistory as s" --find "r.IsActive = 1"
```

### Create a named view with filtering

```bash
gitmap db-view tasks -t "Subtask as s, Plan as p" --viewname v_tasks --find "s.Status = 'DONE'"
```

### Complex filter condition with wildcard and JSON output

```bash
gitmap db-view tasks -t "Subtask as s, Plan as p" --find "s.Status = DONE and p.Title like 'v6*'" --json
```

### Query using short alias and limit

```bash
gitmap dbv master -t "Repository as r" --find "r.Name like 'gitmap*'" --limit 20
```

## See Also

- `gitmap db-read` — Read and inspect SQLite databases
- `gitmap db-find` — Search table columns by wildcard or regex
- `gitmap db-listing` — Discover and inspect registered SQLite databases

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-view
```
