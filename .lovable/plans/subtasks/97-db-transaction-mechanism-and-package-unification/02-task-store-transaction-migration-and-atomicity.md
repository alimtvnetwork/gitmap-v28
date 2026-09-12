# Subtask 02: STORE Package Transaction & Atomicity Unification

Parent Plan: `plans/pending/97-db-transaction-mechanism-and-package-unification.md`  
Status: In Progress  
Ownership: Execution Subagent 2

---

## 1. Objectives

1. In `gitmap/store/makeallvisibility.go`:
   - Refactor `InsertMakeAllVisibilityPendingResults`: replace `db.conn.Begin()` and manual `Commit`/`Rollback` with `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
   - Refactor `MarkMakeAllVisibilityResultsExcluded`: replace `db.conn.Begin()` and manual `Commit`/`Rollback` with `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
   - Update `insertPendingResultsInTx` and `markIdsExcludedInTx` to work with `*dbengine.TxWrapper` or `dbengine.SqlExecutor`.

2. In `gitmap/store/owner_repo_name_index.go`:
   - Refactor `UpsertOwnerRepoNameIndex`: replace `db.conn.Begin()` with `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
   - Pass `*dbengine.TxWrapper` to `populateOwnerRepoIndexTx`.

3. In `gitmap/store/pendingtask.go`:
   - Refactor `CompleteTask`: replace `db.conn.Begin()` with `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
   - Pass `*dbengine.TxWrapper` to `executeCompleteTaskTx`.

4. In `gitmap/store/chromeprofile_delete.go`:
   - In `DeleteChromeProfile`: enclose both `sqlDeleteChromeProfileExports` and `sqlDeleteChromeProfile` deletions in a single transaction via `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.

5. In `gitmap/store/transaction.go`:
   - In `deleteTransactionRows`: wrap the iteration loop of `SQLDeleteTransaction` inside a transaction block via `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.

6. In `gitmap/store/ssh_repo.go`:
   - Adapt `InsertSSHHost(ctx context.Context, host SSHHost, tx *sql.Tx) error` to support callers passing either `*sql.Tx` or `*dbengine.TxWrapper` (or delegate to `SqlExecutor`).

7. In `gitmap/store/txhelpers.go`:
   - Deprecate/cleanup legacy `commitOrWrap` now that store transactions use `dbengine.WithTransaction`.

---

## 2. Acceptance Criteria

- Zero raw `Begin()` or `BeginTx()` calls in `gitmap/store/`.
- All multi-table deletions and operations execute inside transactions with automatic rollback on failure.
- Functions <= 15 lines.
- Affirmative booleans only (`is*`, `has*`).
- Single blank line before `return` and after closing `}`.
