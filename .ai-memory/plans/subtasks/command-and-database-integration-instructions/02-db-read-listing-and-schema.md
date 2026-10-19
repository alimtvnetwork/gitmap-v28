# Subtask 02: `db-read` Engine, Schema Introspection & Deadlock Regression Fix

- **Subtask ID**: `subtask-02-db-read-listing-and-schema`
- **Parent Task**: `command-and-database-integration-instructions`
- **Specification Reference**:
  - `02-spec/21-app/command-and-database-integration-instructions/02-cli-command-contracts.md`
  - `02-spec/21-app/command-and-database-integration-instructions/05-antigravity-verification-and-testing.md`
- **Owned Files**:
  - `cli/cmddb/cmddb_read.go`
  - `cli/cmddb/cmddb_tables.go`
  - `cli/cmddb/cmddb_schema.go`
  - `cli/cmddb/cmddb_read_test.go`
  - `cli/cmddb/cmddb_deadlock_test.go`
- **Assigned Worker**: Worker 02
- **Status**: PENDING

---

## 1. Objectives

1. Implement the unified `gitmap db-read` (alias `gitmap dbr`) command in `cli/cmddb/cmddb_read.go`:
   - Resolve target database from first argument (alias or path).
   - Auto-register alias if `-as <alias>` flag is specified alongside a path.
   - Enforce read-only execution safety via `PRAGMA query_only = ON;` and `PRAGMA busy_timeout = 5000;`.
2. Support raw SQL query execution with optional `--limit/-l` flag (defaults to 100 rows).
3. Support schema introspection subcommands:
   - `gitmap db-read <db> ls`: Enumerate all user tables and views with column and row counts.
   - `gitmap db-read <db> table <table1> [table2...] [--columns/-c] [--json]`: Introspect column names, types, default values, and primary keys.
4. **Deadlock Regression Fix:** Refactor `listTablesAndViews` in `cli/cmddb/cmddb_tables.go` into a two-phase architecture (collect table names into a slice, explicitly close `rows`, then query individual column/row counts) to eliminate connection pool starvation when `MaxOpenConns = 1`.
5. Author `TestListTablesAndViews_DeadlockRegression_MaxOpenConns1` verifying that concurrent or sequential table listings never block or deadlock.

---

## 2. CLI Command Syntax

```bash
# Direct SQL query using alias
gitmap db-read master "SELECT * FROM Repositories" --limit 50 --json

# Direct SQL query using physical path with auto-alias assignment
gitmap db-read "./data/local.db" -as "local" "SELECT count(*) FROM users"

# Enumerate tables and row statistics
gitmap db-read master ls

# Inspect column definitions and types across tables
gitmap db-read master table Repositories PullErrors --columns --json
```

---

## 3. Verification Criteria

1. `go test ./cli/cmddb -run TestDbRead -v` passes cleanly.
2. `go test ./cli/cmddb -run TestListTablesAndViews_DeadlockRegression_MaxOpenConns1 -v` completes in < 1 second.
3. Executing write queries (`UPDATE`, `DELETE`, `INSERT`) via `db-read` fails immediately with `ErrReadOnlyDatabase` violation.
4. `--json` output returns valid JSON payload containing headers, records, row count, and execution duration.
