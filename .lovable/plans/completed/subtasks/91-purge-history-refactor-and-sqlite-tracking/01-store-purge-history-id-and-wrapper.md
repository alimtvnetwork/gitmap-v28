# Subtask 01: Store Purge History ID Normalization & QueryWrapper Pattern Reuse

## Objective
Refactor `gitmap/store/purge_history.go` to adhere to repository database guidelines:
- Rename all-caps `ID` to `PurgeHistoryLogId` (canonical PK) and provide `Id` (PascalCase) and `ID` (deprecated compatibility alias).
- Rename `Restored` boolean to `IsRestored` (canonical boolean) and provide `Restored` (deprecated compatibility alias).
- Change SQLite schema to use `PurgeHistoryLogId INTEGER PRIMARY KEY AUTOINCREMENT`, `IsRestored INTEGER NOT NULL DEFAULT 0`, `Notes TEXT NULL`, and `Comments TEXT NULL`.
- Add backward-compatible `migratePurgeHistoryColumns()` to handle existing databases with `ID` or `Restored`.
- Adopt `ExecWrapper` and `QueryRowWrapper` from `gitmap/store/wrapper.go` for all SQL executions with explicit error logging and structured `QueryResult`.
- Eliminate swallowed error from `res.LastInsertId()` by verifying `err != nil`.
- Wire `EnsurePurgeHistoryTable()` into `gitmap/store/store.go:Migrate()`.

## Files Affected
- `gitmap/store/purge_history.go`
- `gitmap/store/store.go`
