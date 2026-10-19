# Specification: Antigravity Verification, Deadlock Regression & CI Testing

- **Spec Path:** `02-spec/21-app/command-and-database-integration-instructions/05-antigravity-verification-and-testing.md`
- **Module:** `cli/cmddb`, `cli/cmddb/tests`, `test/e2e`
- **Status:** APPROVED
- **Version:** 1.0.0
- **Line Count Cap:** <= 300 lines

---

## 1. User Request (Verbatim)

```text
Author the verification specification and parent execution plan:
1. 02-spec/21-app/command-and-database-integration-instructions/05-antigravity-verification-and-testing.md: End-to-end verification protocol testing against live Antigravity databases (C:\Users\Administrator\.gemini\antigravity\conversation_summaries.db and conversation trajectories), mock database tests, deadlock regression test for listTablesAndViews (MaxOpenConns=1), and CI integration.
2. .ai-memory/plans/command-and-database-integration-instructions.md and subtask files in .ai-memory/plans/subtasks/command-and-database-integration-instructions/ (01-spec-and-registry-engine.md, 02-db-read-listing-and-schema.md, 03-db-add-edit-and-mutation.md, 04-db-find-and-view-joins.md, 05-skills-prompts-and-minor-bump.md).
Include verbatim user prompt, 5 pillars, and subtask milestones.
Use strictly relative Git paths and lowercase filenames. All files <= 300 lines.
```

---

## 2. Testing Philosophy & Multi-Tier Test Pyramid

GitMap database commands (`db-read`, `db-add`, `db-edit`, `db-find`, `db-view`, `db-listing`) operate across sensitive developer databases, multi-agent workspaces, and CI pipelines. Verification adheres to a three-tier testing pyramid:

1. **Unit & Isolation Tier:** Pure Go tests using hermetic in-memory SQLite (`:memory:`) or temporary directory file databases (`t.TempDir()`). Fast execution (<5ms per test), zero external environment dependencies, and executed on every CI commit.
2. **Deadlock & Concurrency Tier:** Stress and pool constraint tests enforcing `MaxOpenConns = 1` to guarantee nested query routines never lock SQLite database connections.
3. **Live Antigravity Integration Tier:** Environment-aware verification interrogating live Google Antigravity workspaces (`conversation_summaries.db` and conversation trajectory stores). If live databases are absent (such as in headless GitHub Actions runners), tests gracefully degrade or synthesize trajectory fixtures.

---

## 3. Live Antigravity Database Verification Protocol

### 3.1 Antigravity Database Discovery & Ingestion

GitMap discovers the local Antigravity workspace via environment variables or default user paths:
- Windows: `%USERPROFILE%\.gemini\antigravity\conversation_summaries.db`
- Unix / macOS: `~/.gemini/antigravity/conversation_summaries.db`
- Override: Environment variable `ANTIGRAVITY_DB_PATH` or `GEMINI_HOME`.

### 3.2 Target Tables & Schemas

The Antigravity `conversation_summaries.db` database contains primary and virtual FTS tables:
- `conversation_summaries`: Primary metadata table containing:
  - `conversation_id` (TEXT PRIMARY KEY)
  - `title` (TEXT)
  - `preview` (TEXT)
  - `step_count` (INTEGER)
  - `last_modified_time` (INTEGER / DATETIME)
  - `project_id` (TEXT)
  - `agent_name` (TEXT)
  - `status` (TEXT)
- `conversations_fts` (FTS5 virtual table for full-text search)

### 3.3 Read-Only Safety Invariant

When reading live Antigravity or third-party IDE databases, GitMap MUST enforce strict query isolation:
```sql
PRAGMA query_only = ON;
PRAGMA busy_timeout = 5000;
```
Any command invoking `db-read`, `db-find`, or `db-view` against registered system databases is forbidden from executing write operations (`INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`). Attempted writes return `ErrReadOnlyDatabase` with exit code 1.

### 3.4 Live Verification CLI Invocations

```bash
# 1. Register live Antigravity database alias
gitmap db-read "$USERPROFILE/.gemini/antigravity/conversation_summaries.db" -as "antigravity-live"

# 2. Verify registration in listing table
gitmap db-listing --json

# 3. List tables and record counts from live database
gitmap db-read antigravity-live ls

# 4. Read schema and column types for conversation_summaries
gitmap db-read antigravity-live table conversation_summaries --columns --json

# 5. Query recent conversation summaries with limit
gitmap db-read antigravity-live "SELECT conversation_id, title, step_count FROM conversation_summaries ORDER BY last_modified_time DESC LIMIT 5" --json

# 6. Wildcard and regex search across conversation titles and previews
gitmap db-find antigravity-live -table "conversation_summaries" --find "heatmap*" --limit 10
```

---

## 4. Deadlock Regression Testing (`MaxOpenConns = 1`)

### 4.1 Root Cause Analysis of `listTablesAndViews` Deadlock

In `cli/cmddb/cmddb_tables.go`, `listTablesAndViews(conn *sql.DB)` executed:
1. `rows, err := conn.Query("SELECT name, type FROM sqlite_master ...")`
2. Inside `for rows.Next()`, it invoked `fetchTableStats(conn, name, tableType)`.
3. `fetchTableStats` invoked `countTableColumns` and `countTableRows`, executing `conn.QueryRow(...)`.

**Flaw:** While `rows.Next()` is actively streaming, the database connection is checked out of the pool. When `conn.SetMaxOpenConns(1)` is set (standard for SQLite to avoid lock contention or single-writer serialization), `conn.QueryRow` blocks waiting for a connection to return to the pool. Because `rows.Close()` is deferred until function exit, the caller deadlocks permanently.

### 4.2 Two-Phase Remediation Architecture

The query is partitioned into two distinct phases:
- **Phase 1 (Introspection):** Query `sqlite_master`, exhaustively scan all table/view names and types into a local struct slice, and explicitly call `rows.Close()`.
- **Phase 2 (Aggregation):** Iterate over the slice and execute column and row count queries with an idle pool connection.

```go
func listTablesAndViewsSafe(conn *sql.DB) ([]TableInfo, error) {
    query := `SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY type, name`
    rows, err := conn.Query(query)
    if err != nil {
        return nil, fmt.Errorf("query sqlite_master: %w", err)
    }
    
    type rawItem struct {
        name      string
        tableType string
    }
    var rawItems []rawItem
    for rows.Next() {
        var item rawItem
        if err := rows.Scan(&item.name, &item.tableType); err != nil {
            rows.Close()
            return nil, err
        }
        rawItems = append(rawItems, item)
    }
    rows.Close() // Explicit close releases connection back to pool

    var tables []TableInfo
    for _, item := range rawItems {
        colCount := countTableColumns(conn, item.name)
        rowCount := countTableRows(conn, item.name)
        tables = append(tables, TableInfo{
            Name:        item.name,
            Type:        item.tableType,
            ColumnCount: colCount,
            RowCount:    rowCount,
        })
    }
    return tables, nil
}
```

### 4.3 Automated Deadlock Regression Test

```go
func TestListTablesAndViews_DeadlockRegression_MaxOpenConns1(t *testing.T) {
    dbPath := filepath.Join(t.TempDir(), "deadlock_test.db")
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        t.Fatalf("failed to open sqlite test db: %v", err)
    }
    defer db.Close()

    // Strictly enforce single connection pool constraint
    db.SetMaxOpenConns(1)

    // Create 15 tables with multiple columns and records
    for i := 1; i <= 15; i++ {
        tableName := fmt.Sprintf("tbl_test_%02d", i)
        _, err = db.Exec(fmt.Sprintf("CREATE TABLE %s (id INTEGER PRIMARY KEY, title TEXT, score REAL);", tableName))
        if err != nil {
            t.Fatalf("failed to create table %s: %v", tableName, err)
        }
        _, err = db.Exec(fmt.Sprintf("INSERT INTO %s (title, score) VALUES ('alpha', 9.5), ('beta', 8.2);", tableName))
        if err != nil {
            t.Fatalf("failed to insert records into %s: %v", tableName, err)
        }
    }

    done := make(chan struct{})
    var tables []TableInfo
    var listErr error

    go func() {
        tables, listErr = listTablesAndViews(db)
        close(done)
    }()

    select {
    case <-done:
        if listErr != nil {
            t.Fatalf("listTablesAndViews returned unexpected error: %v", listErr)
        }
        if len(tables) != 15 {
            t.Fatalf("expected 15 tables, got %d", len(tables))
        }
    case <-time.After(3 * time.Second):
        t.Fatal("DEADLOCK DETECTED: listTablesAndViews blocked on MaxOpenConns=1 connection pool")
    }
}
```

---

## 5. Hermetic Mock Database Test Fixtures

### 5.1 Synthetic Test Dataset Generation

Unit test suites construct reproducible SQLite databases in temporary directories:
1. `Schema`:
   - `Agents` (`AgentId`, `AgentName`, `Role`, `Capacity`, `IsActive`)
   - `Tasks` (`TaskId`, `AgentId`, `Slug`, `Status`, `StepCount`, `CreatedAt`)
   - `Logs` (`LogId`, `TaskId`, `LogLevel`, `Message`, `Timestamp`)
2. `Seed Data`:
   - 5 Agents (`Worker 01` through `Worker 05`).
   - 25 Tasks across `PENDING`, `IN_PROGRESS`, `DONE`, `BLOCKED`.
   - 100 Log records containing searchable keywords (`error`, `timeout`, `success`, `heatmap`).

### 5.2 Mock Command Test Matrix

| Command Target | Test Function | Input Assertions | Expected Outcome |
| :--- | :--- | :--- | :--- |
| `db-read` (SQL) | `TestDbRead_RawSQL_QueryOnly` | `SELECT * FROM Agents WHERE IsActive=1` | Returns 5 rows in JSON; attempts to `DROP TABLE` fail with error |
| `db-read` (ls) | `TestDbRead_Ls_TableCounts` | `ls` command argument | Outputs table names with correct column & row counts |
| `db-read` (table) | `TestDbRead_TableSchema_Types` | `table Tasks --columns --json` | Validates column types (`INTEGER`, `TEXT`, `DATETIME`) |
| `db-add` (CSV) | `TestDbAdd_CsvValues_Insert` | `Tasks csv 1, "task-99", "DONE", 50` | Inserts new record; primary key incremented; row count increases |
| `db-edit` (id) | `TestDbEdit_UpdateRecord_ById` | `-table "Tasks" -id "1" csv ...` | Modifies targeted record; verifies updated column values |
| `db-find` (wildcard)| `TestDbFind_WildcardPatterns` | `--find "task*"` and `"*99"` | Finds matching records across multiple text columns |
| `db-find` (regex) | `TestDbFind_RegexSearch` | `--grep "^(Worker 0[1-3])$"` | Regex filter returns only workers 1, 2, and 3 |
| `db-view` (join) | `TestDbView_AutoJoin_Inference` | `-t "Tasks as t, Agents as a"` | Generates `ON t.AgentId = a.AgentId` join and queries projected columns |
| `db-listing` | `TestDbListing_AliasRegistry` | Register 3 test DBs; run `db-listing` | Validates JSON envelope with path, alias name, file size, query count |

---

## 6. Continuous Integration (CI) Workflow Integration

### 6.1 Test Execution Invariant

All database unit and deadlock regression tests must execute hermetically within standard CI pipelines:
```bash
go test ./cli/cmddb/... -v -race -timeout=60s
```

### 6.2 Performance Gates

1. **In-Memory Query Latency:** Sub-millisecond (< 2ms) response for `db-read` queries on schemas with < 10,000 records.
2. **Registry Introspection:** `db-listing` completes in < 15ms across up to 50 registered database aliases.
3. **Deadlock Test Duration:** `TestListTablesAndViews_DeadlockRegression_MaxOpenConns1` completes in < 100ms.
4. **Zero Flakiness:** Tests run deterministically without time-dependent races or unclosed file descriptors.
