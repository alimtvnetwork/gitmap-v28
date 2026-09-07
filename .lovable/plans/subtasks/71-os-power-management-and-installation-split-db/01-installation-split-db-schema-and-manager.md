# Subtask 01: Installation Split DB Schema & Connection Manager

## Status
Completed

## Context & Objectives
Implement the dedicated split SQLite database for installations (`installation.db`):
1. **Path & Connection**:
   - Location: `filepath.Join(BinaryDataDir(), "installation.db")`.
   - Driver: `modernc.org/sqlite`.
   - Concurrency: `conn.SetMaxOpenConns(1)`.
2. **Schema Definition**:
   - `InstalledTool`: `InstalledToolId INTEGER PRIMARY KEY AUTOINCREMENT`, `Tool TEXT NOT NULL UNIQUE`, `VersionMajor INTEGER NOT NULL DEFAULT 0`, `VersionMinor INTEGER NOT NULL DEFAULT 0`, `VersionPatch INTEGER NOT NULL DEFAULT 0`, `VersionBuild INTEGER NOT NULL DEFAULT 0`, `VersionString TEXT NOT NULL DEFAULT ''`, `PackageManager TEXT NOT NULL DEFAULT ''`, `InstallPath TEXT NOT NULL DEFAULT ''`, `Description TEXT NULL` (Rule 10), `InstalledAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP`, `UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP`.
   - `InstallationLog`: `InstallationLogId INTEGER PRIMARY KEY AUTOINCREMENT`, `Tool TEXT NOT NULL`, `Action TEXT NOT NULL`, `Version TEXT NOT NULL DEFAULT ''`, `PackageManager TEXT NOT NULL DEFAULT ''`, `DurationMs INTEGER NOT NULL DEFAULT 0`, `IsSuccess INTEGER NOT NULL DEFAULT 1`, `ExitCode INTEGER NOT NULL DEFAULT 0`, `Stdout TEXT NULL`, `Stderr TEXT NULL`, `CommandLine TEXT NULL`, `Notes TEXT NULL` (Rule 11), `Comments TEXT NULL` (Rule 11), `CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP`.
   - Indexes: `IdxInstalledTool_Tool`, `IdxInstallationLog_Tool`, `IdxInstallationLog_Action`, `IdxInstallationLog_CreatedAt`.
3. **Manager Implementation**:
   - `gitmap/store/installation_split_db.go` (<= 200 lines): `InstallationSplitDB` struct, open, close, and schema initialization.
   - `gitmap/store/installation_split_tools.go` (<= 200 lines): `SaveInstalledTool`, `GetInstalledTool`, `ListInstalledTools`, `RemoveInstalledTool`, `IsToolInstalled`.
   - `gitmap/store/installation_split_log.go` (<= 200 lines): `RecordLog`, `GetLogs`, `GetLogsByTool`.

## Verification Steps
- Unit tests in `gitmap/store/installation_split_db_test.go` verify table creation, CRUD operations, and logs.
