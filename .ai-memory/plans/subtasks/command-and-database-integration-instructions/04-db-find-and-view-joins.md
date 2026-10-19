# Subtask 04: `db-find` & `db-view` Virtual Auto-Join Query Engine

- **Subtask ID**: `subtask-04-db-find-and-view-joins`
- **Parent Task**: `command-and-database-integration-instructions`
- **Specification Reference**:
  - `02-spec/21-app/command-and-database-integration-instructions/04-view-and-join-engine.md`
- **Owned Files**:
  - `cli/cmddb/cmddb_find.go`
  - `cli/cmddb/cmddb_view.go`
  - `cli/cmddb/cmddb_join_inference.go`
  - `cli/cmddb/cmddb_find_test.go`
  - `cli/cmddb/cmddb_view_test.go`
- **Assigned Worker**: Worker 01
- **Status**: PENDING

---

## 1. Objectives

1. Implement `gitmap db-find` in `cli/cmddb/cmddb_find.go`:
   - Multi-table search support: `-table "table1, table2"`.
   - Wildcard pattern matching:
     - Exact match: `"text"`
     - Suffix match: `"text*"`
     - Prefix match: `"*text"`
     - Substring match: `"*text*"`
   - Regular expression scanning: `--grep/-r "<pattern>"`.
   - Column projection filter: `--columns "col1, col2"` (validates column existence across specified tables).
   - Configurable row limit (`--limit/-l`, default 100) and structured `--json` output.
2. Implement `gitmap db-view` in `cli/cmddb/cmddb_view.go`:
   - Table aliasing argument: `-t "table1 as t1, table2 as t2"`.
   - Named view persistence: `--viewname "<view_name>"` creating or querying existing virtual views.
   - **Automatic Join Inference:** Engine inspects schema relationships, foreign keys, and shared column names (e.g., `ProjectId`, `RepoId`, `UserId`) to generate optimal SQL `INNER JOIN` conditions without manual boilerplate.
   - Condition filter execution: `--find "t1.status = 'ACTIVE' and t2.count > 5"`.
3. Provide unit test suites in `cli/cmddb/cmddb_find_test.go` and `cli/cmddb/cmddb_view_test.go` verifying join resolution, pattern compilation, and pagination limits.

---

## 2. CLI Command Syntax

```bash
# Wildcard find across multiple tables
gitmap db-find master -table "Repositories, DetectedProjects" --find "gitmap*" --limit 50 --json

# Regex pattern search across specified columns
gitmap db-find master -table "Repositories" --columns "Name, Path" --grep "^(lib|cli)-.*"

# Declarative dynamic view creation with auto-join inference
gitmap db-view master -t "Repositories as r, DetectedProjects as p" --viewname "v_active_projects" --find "r.Status = 'ACTIVE'" --json
```

---

## 3. Verification Criteria

1. `go test ./cli/cmddb -run TestDbFind -v` passes cleanly.
2. `go test ./cli/cmddb -run TestDbView -v` passes cleanly.
3. Wildcard `text*` executes efficient SQLite `LIKE 'text%'` query utilizing indices where applicable.
4. Regex `--grep` executes without crashing on invalid expressions and filters records deterministically.
5. Inferred joins accurately resolve parent-child keys across multiple tables.
