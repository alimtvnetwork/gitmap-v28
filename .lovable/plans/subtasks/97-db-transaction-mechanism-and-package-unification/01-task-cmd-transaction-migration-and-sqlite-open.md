# Subtask 01: CMD Package Transaction & SQLite Open Unification

Parent Plan: `plans/pending/97-db-transaction-mechanism-and-package-unification.md`  
Status: In Progress  
Ownership: Execution Subagent 1

---

## 1. Objectives

1. In `gitmap/cmd/sequence_cmd.go`:
   - Replace `repoDB.BeginTx(ctx, nil)` with `dbengine.WrapDb(repoDB, dbengine.DbSQLite).WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError { ... })`.
   - Update `executeSaveSequenceTx` and `insertSequenceFilesTx` to accept `tx *dbengine.TxWrapper` or use `tx.Exec` / `tx.Prepare`.
   - Eliminate manual rollback/commit boilerplate.

2. In `gitmap/cmd/sshjoin_cmd.go`:
   - Replace `db.BeginTx(ctx, nil)` in `runJoinTransaction` with `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError { ... })`.
   - Update `insertJoinRecords` and `logSSHJoinInTx` to work with `tx *dbengine.TxWrapper`.

3. In `gitmap/cmd/ssh_alias_cmd.go`:
   - In `saveAliasCommand`, replace `db.Conn().BeginTx(ctx, nil)` with `dbengine.WrapDb(db.Conn(), dbengine.DbSQLite).WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError { ... })`.

4. In `gitmap/cmd/commitin/runlog/inputs.go`:
   - In `InsertSourceCommit`, replace `db.Begin()` with `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError { ... })`.
   - Update `executeCommitInsertTx`, `insertSourceCommitTx`, `insertSourceFilesTx` to execute against `tx *dbengine.TxWrapper` or `dbengine.SqlExecutor`.

5. In `gitmap/cmd/chromeprofile_export_all.go`:
   - In `writeAllChromeProfilesSQLite`, replace `sql.Open("sqlite", outPath)` with `store.OpenSQLiteDB(outPath)`.
   - In `populateProfilesInSQLite`, wrap the multi-profile insert loop inside `dbengine.WrapDb(db, dbengine.DbSQLite).WithTransaction(...)`.

6. In `gitmap/cmd/agy_conv_scanner.go`:
   - In `readSingleConvDB`, replace `sql.Open("sqlite", dbPath)` with `store.OpenSQLiteDB(dbPath)`.

---

## 2. Acceptance Criteria

- Zero raw `Begin()` or `BeginTx()` calls in `gitmap/cmd/`.
- All function bodies <= 15 lines.
- Affirmative booleans only (`is*`, `has*`).
- Single blank line before `return` and after closing `}`.
