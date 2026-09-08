# Architectural Audit & Design: Split Database Registry & Installation Audit Logging

**Date:** 2026-09-08
**Audited Components:** `gitmap/store/` and `gitmap/cmd/` (`install.go`, `installtools.go`, `install_chrome_deb.go`)
**Compliance Standards:** `spec/04-database-conventions/`, `spec/05-split-db-architecture/`, `db-sqlite-conventions` skill.

---

## 1. Executive Summary & Gaps

1. **Main Database Disconnect from Split DBs:**
   - Root database (`gitmap.db` / `store.go`) had zero registry tracking for split databases (`installation.db`, `schedules/*.db`, `pipeline_db/*.db`).
   - The main DB could not verify the presence, size, record count, or health of `installation.db` without ad-hoc filesystem probes.

2. **Incomplete Installation Audit Logging in `installation.db` (`InstallationLog`):**
   - `recordInstallation()` in `cmd/installtools.go` only wrote `Tool`, `Action="install"`, `Version`, `PackageManager`, and `IsSuccess=true`.
   - `DurationMs`, `ExitCode`, `Stdout`, `Stderr`, and `CommandLine` were defaulted or left empty.
   - Failures were never written to `InstallationLog` in `installation.db` (only to flat text files).

3. **Inconsistent and Broken Handler Integrations:**
   - `runInstallVSCodeLinux()` never called `recordInstallation()`.
   - `runInstallChromeLinux()` did not log multi-stage metrics or failure telemetry.
   - User aborts and validation errors were not audited.

---

## 2. SplitDatabaseRegistry in Root Database (`store.go`)

### 2.1 Table Schema
```sql
CREATE TABLE IF NOT EXISTS SplitDatabaseRegistry (
    SplitDatabaseRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseType            TEXT NOT NULL,
    DatabaseKey             TEXT NOT NULL,
    DatabasePath            TEXT NOT NULL,
    SizeBytes               INTEGER NOT NULL DEFAULT 0,
    TableCount              INTEGER NOT NULL DEFAULT 0,
    RecordCount             INTEGER NOT NULL DEFAULT 0,
    SchemaVersion           INTEGER NOT NULL DEFAULT 1,
    Status                  TEXT NOT NULL DEFAULT 'active',
    IsActive                INTEGER NOT NULL DEFAULT 1,
    IsAttached              INTEGER NOT NULL DEFAULT 0,
    Description             TEXT NULL,
    Notes                   TEXT NULL,
    Comments                TEXT NULL,
    LastAccessedAt          INTEGER NOT NULL DEFAULT 0,
    LastSyncedAt            INTEGER NOT NULL DEFAULT 0,
    CreatedAt               INTEGER NOT NULL DEFAULT (unixepoch()),
    UpdatedAt               INTEGER NOT NULL DEFAULT (unixepoch()),
    UNIQUE(DatabaseType, DatabaseKey)
);

CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_Type ON SplitDatabaseRegistry(DatabaseType);
CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_Status ON SplitDatabaseRegistry(Status);
CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_UpdatedAt ON SplitDatabaseRegistry(UpdatedAt);
```

### 2.2 Go API on `*store.DB`
- `RegisterSplitDB(entry SplitDatabaseEntry) error`
- `GetSplitDB(dbType, dbKey string) (*SplitDatabaseEntry, error)`
- `ListSplitDBs(dbType string) ([]SplitDatabaseEntry, error)`
- `SyncKnownSplitDatabases() error`

---

## 3. Comprehensive Installation Audit Logging in `installation.db`

### 3.1 Enhanced Store Model & Methods (`store/installation_split_log.go`)
- `RecordExecution(tool, action, version, manager string, durationMs int64, isSuccess bool, exitCode int, stdout, stderr, cmdLine, notes, comments string) error`
- `GetFailedLogs(limit int) ([]InstallationLogRecord, error)`
- `SanitizeLogOutput(s string) string` (truncates at 64 KB per stream to prevent SQLite bloat)

### 3.2 Dual-Stream Command Runner in `cmd/installtools.go`
- `executeCommandWithAudit(args []string, verbose bool) commandExecutionResult` captures `Stdout`, `Stderr`, `ExitCode`, `DurationMs`, and `CommandLine`.
- Records every execution attempt into `InstallationLog` in `installation.db` for both success and failure cases.
