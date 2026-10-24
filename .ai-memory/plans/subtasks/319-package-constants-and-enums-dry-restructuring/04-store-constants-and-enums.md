# Subtask 04: `cli/store` Constants & Enums Refactoring

> **Parent Plan:** [.ai-memory/plans/pending/319-package-constants-and-enums-dry-restructuring.md](../../pending/319-package-constants-and-enums-dry-restructuring.md)  
> **Target Package:** `cli/store`  
> **Status:** Ready for Execution (Next Turn)  

---

## 1. Scope & Technical Deliverables
1. **Enum: `StoreTableType`** (`cli/store/enums/store_table_enum.go`):
   - Values: `StoreTableInvalid` (0), `StoreTableRepositories`, `StoreTablePipelines`, `StoreTableClusters`, `StoreTableWorkgroups`, `StoreTableAuditLogs`, `StoreTableCredentials`.
   - Methods: `IsValid() bool`, `String() string`, `TableName() string`.
2. **Enum: `DbOpType`** (`cli/store/enums/db_op_enum.go`):
   - Values: `DbOpInvalid` (0), `DbOpQuery`, `DbOpExec`, `DbOpBegin`, `DbOpCommit`, `DbOpRollback`.
   - Methods: `IsValid() bool`, `String() string`, `IsWrite() bool`, `IsTransaction() bool`.
3. **Constants: `cli/store/constants/constants.go`**:
   - SQLite Pragmas: `PragmaJournalWAL`, `PragmaSyncNormal`, `PragmaBusyTimeout`.
   - Table columns: `ColRepoID`, `ColRepoPath`, `ColRepoName`, `ColUpdatedAt`, `ColStatus`.
   - Limits: `DefaultMaxOpenConns`, `DefaultMaxIdleConns`, `DefaultConnMaxLifetimeSeconds`.
4. **Call-Site Refactoring:**
   - Update `cli/store/store.go`, `cli/store/db_sqlite.go`, `cli/store/repo_query.go` to eliminate inline SQL table strings.

---

## 2. Coding Guidelines Invariants
- Files $\le 300$ lines.
- Functions $\le 15$ lines.
- Positive booleans only (`is*`, `has*`).
- Zero nested `if` statements.

---

## 3. Verification Evidence Contract
- Run `gofmt -w cli/store/`.
- AST pass via FastGate.
