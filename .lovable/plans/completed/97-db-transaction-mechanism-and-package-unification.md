# 97-db-transaction-mechanism-and-package-unification

## Execution Summary
- **Started By:** User command asking to confirm and ensure everywhere in DB transactions that the designed DB package and mechanism are used properly, following `.agents/skills/execute-parent-task-with-n-steps` (v2.1.0, N=130).
- **Duration / Loops:** Completed in 1 self-loop cycle across 2 concurrent execution subagents (CMD migrator and Store migrator) operating on disjoint file boundaries.
- **Status:** Completed & Consolidated.

---

## Acceptance Criteria (Verbatim Echo)

```text
1. Zero occurrences of raw `Begin()` or `BeginTx()` in `gitmap/cmd/` and `gitmap/store/`.
2. All database mutations requiring atomicity route through `dbengine.WrapDb(...).WithTransaction(...)` or `WithImmediateTransaction(...)`.
3. All SQLite connections opened in utility commands use `store.OpenSQLiteDB(...)`.
4. All functions <= 15 lines with affirmative booleans and strict blank-line spacing.
5. All 29 gates in `06-cicd-local-runner.py` pass 100% green (`exit 0`).
```

---

## Accomplishments

### 1. CMD Package Transaction & Connection Migrations
- **`gitmap/cmd/sequence_cmd.go`:**
  - Migrated `repoDB.BeginTx(ctx, nil)` to `dbengine.WrapDb(repoDB, dbengine.DbSQLite).WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError { ... })`.
  - Refactored `executeSaveSequenceTx` and `insertSequenceFilesTx` to execute queries via `tx *dbengine.TxWrapper`, eliminating manual commit and rollback calls.
- **`gitmap/cmd/sshjoin_cmd.go`:**
  - Migrated `db.BeginTx(ctx, nil)` in `runJoinTransaction` to canonical `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(...)`.
  - Refactored `insertJoinRecords` and `logSSHJoinInTx` to receive `tx *dbengine.TxWrapper` and return typed `*apperror.AppError`.
- **`gitmap/cmd/ssh_alias_cmd.go`:**
  - Migrated `db.Conn().BeginTx(ctx, nil)` to `dbengine.WrapDb(db.Conn(), dbengine.DbSQLite).WithTransaction(...)`.
  - Decomposed execution into `runSaveAliasTx` and `executeSaveAliasTx` adhering to <= 15 lines per function.
- **`gitmap/cmd/commitin/runlog/inputs.go`:**
  - Migrated raw `db.Begin()` in `InsertSourceCommit` to `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(...)`.
  - Updated `executeCommitInsertTx`, `insertSourceCommitTx`, and `insertSourceFilesTx` to execute against `tx *dbengine.TxWrapper` and removed legacy `commitRunlogTx`.
- **`gitmap/cmd/chromeprofile_export_all.go`:**
  - Replaced raw `sql.Open("sqlite", outPath)` with centralized factory `store.OpenSQLiteDB(outPath)`.
  - Wrapped multi-profile insertion loops (`populateProfilesInSQLite`, `insertProfileToSQLite`) inside a single atomic transaction via `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(...)`.
- **`gitmap/cmd/agy_conv_scanner.go`:**
  - Replaced raw `sql.Open("sqlite", dbPath)` with centralized factory `store.OpenSQLiteDB(dbPath)`.
  - Added unit test `gitmap/cmd/agy_conv_scanner_test.go`.

### 2. STORE Package Transaction & Connection Migrations
- **`gitmap/store/makeallvisibility.go`:**
  - Migrated `InsertMakeAllVisibilityPendingResults` and `MarkMakeAllVisibilityResultsExcluded` from raw `db.conn.Begin()` and manual `Commit`/`Rollback` to `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
  - Refactored `insertPendingRow`, `insertPendingResultsInTx`, and `markIdsExcludedInTx` to execute via `tx *dbengine.TxWrapper` returning typed `*apperror.AppError`.
- **`gitmap/store/owner_repo_name_index.go`:**
  - Migrated `UpsertOwnerRepoNameIndex` from raw `db.conn.Begin()` to `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
  - Updated `populateOwnerRepoIndexTx` and `insertOwnerRepoNames` to receive `tx *dbengine.TxWrapper`.
- **`gitmap/store/pendingtask.go`:**
  - Migrated `CompleteTask` from raw `db.conn.Begin()` to `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
  - Updated `executeCompleteTaskTx`, `insertCompletedTaskInTx`, `deletePendingTaskInTx`, and `findPendingTaskInTx` to receive `tx *dbengine.TxWrapper`.
- **`gitmap/store/chromeprofile_delete.go`:**
  - Enclosed cascaded `sqlDeleteChromeProfileExports` and `sqlDeleteChromeProfile` deletions inside a single transaction via `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
- **`gitmap/store/transaction.go`:**
  - Enclosed excess transaction pruning loop in `deleteTransactionRows` inside an atomic transaction via `dbengine.WrapDb(db.conn, dbengine.DbSQLite).WithTransaction(...)`.
- **`gitmap/store/ssh_repo.go`:**
  - Updated `InsertSSHHost` to support `tx any` (`*sql.Tx`, `*dbengine.TxWrapper`, `dbengine.SqlExecutor`, or `*sql.DB`), adding typed helper `InsertSSHHostTx(ctx, host, tx)`.
- **`gitmap/store/txhelpers.go`:**
  - Removed unused legacy `txhelpers.go` (`txExecer` and `commitOrWrap`) cleanly.

### 3. Verification & Quality Gates
- All unit tests across `cmd`, `store`, and `dbengine` pass cleanly.
- Full parallel quality runner `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests` passed with 29/29 gates green (exit code 0).
- Zero raw `Begin()` or `BeginTx()` calls remain in production code across the entire repository.
