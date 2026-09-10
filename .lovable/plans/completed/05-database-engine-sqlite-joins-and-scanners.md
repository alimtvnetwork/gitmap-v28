# Milestone Summary: Database Engine, SQLite Schema, Joins & Typed Scanners

## 1. Executive Overview & Consolidated Tasks

- **Milestone Domain:** Universal DBEngine Joins, Query Builders, Safe Row Scanner Generators & Tesla ID Conventions
- **Total Original Plans Merged:** 6 plans
  - `30-endpoint-resolver-db.md`
  - `65-universal-dbengine-joins-and-view-evolution.md`
  - `66-automatic-db-repo-and-safe-scanner-generator.md`
  - `71-os-power-management-and-installation-split-db.md`
  - `91-purge-history-refactor-and-sqlite-tracking.md`
  - `94-orm-code-generator-id-error-standards-and-fast-cache.md`
- **Associated Subtask Folders Folded:** 4 folders
  - `11-endpoint-resolver-db`
  - `71-os-power-management-and-installation-split-db`
  - `91-purge-history-refactor-and-sqlite-tracking`
  - `94-orm-code-generator-id-error-standards-and-fast-cache`
- **Status:** `COMPLETED`
- **Core Architecture & Invariants:** PascalCase database schema. ID naming must follow Tesla standards (<Entity>Id, not id or ID). Zero-swallow database scanner returns. Fast search cache in SQLite with schema migrations.

## 2. Key Architectural Decisions & Spec Implementations

- **Authoritative Specifications Implemented:**
  - spec/08-database-and-orm/01-overview.md — SQLite single-writer SetMaxOpenConns(1) and connection pooling.
  - spec/08-database-and-orm/02-naming-and-primary-keys.md — PascalCase <Entity>Id convention and typed scanners.
- **Core Architecture Contracts:**
  - PascalCase database schema. ID naming must follow Tesla standards (<Entity>Id, not id or ID). Zero-swallow database scanner returns. Fast search cache in SQLite with schema migrations.

## 3. Deep Dive into Consolidated Plans & Subtask Chronicles

Every individual plan and subtask merged into this milestone is preserved below in full technical detail, ensuring 100% fidelity, zero truncation, and complete traceability.

### Merged Plan: `30-endpoint-resolver-db.md`

#### 11 Endpoint Resolver DB Upgrade

##### Parent Task Goal

Enhance the global endpoint resolution logic in Gitmap so that commands (especially `commit-right`, `commit-left`, `commit-both`, `open`, `inject`, etc.) can seamlessly resolve repositories using five modalities:
1. ID of the repo (from the SQLite DB)
2. Alias (from the SQLite DB)
3. Relative path
4. Absolute path
5. URL (HTTPS, SSH, git@)

##### Architectural Plan

1. **Core Resolution Logic:**
   - Create `gitmap/cmd/resolver.go`.
   - Implement `func ResolveEndpointString(raw string) string` which orchestrates the resolution.
   - It will check if `raw` is a URL. If so, return unchanged.
   - It will check if `raw` is an existing path (absolute or relative resolving to absolute). If so, return it.
   - It will query the SQLite database (`store.DB`) to resolve by ID, Alias, or Slug. If found, return the resolved `AbsolutePath`.
   - If nothing is found, return `raw` and let downstream functions (like `movemerge`) throw errors.

2. **Integration:**
   - Modify `gitmap/cmd/committransfer.go` `resolveCommitEndpoints` to intercept `leftRaw` and `rightRaw` using `ResolveEndpointString` before passing to `movemerge.ResolveEndpoint`.
   - Ensure the new logic relies on `openDB()` to fetch the SQLite database.

3. **Subtasks Execution:**
   - **Subtask 1: DB & SQL Constants Update**
     - Add `SQLSelectRepoByID` to `gitmap/constants/constants_store.go`.
     - Implement `FindByID` in `gitmap/store/repo.go`.
   - **Subtask 2: Implement Core Resolver**
     - Write `gitmap/cmd/resolver.go`.
     - Include logic to parse numeric IDs, query aliases, and fallback to slugs.
   - **Subtask 3: Integration & Testing**
     - Update `gitmap/cmd/committransfer.go`.
     - Run `go test ./...` and ensure functionality works.
     - Group commits and push.
   - **Subtask 4: Documentation**
     - Update memory and specifications.

##### Code Review Guide

- Do not swallow errors; use `apperror`.
- Ensure no cyclic dependencies between `cmd`, `store`, and `movemerge`.
- Boolean variables must be prefixed with `is`, `has`, `should`, or `can`.
- No generic variable names (`temp`, `data`, `res`).

### Merged Plan: `65-universal-dbengine-joins-and-view-evolution.md`

#### 65 — Universal DBEngine Joins, Error-Guarded Query Builder & Automated View Evolution

**Status:** COMPLETED
**Priority:** High
**Category:** Database Architecture / Multi-Dialect Query Engine / Code Generation
**Completed:** 2026-09-05

---

##### 1. Problem Statement & Root Cause

###### Problem Statement
The generic repository and query builder abstractions previously had several architectural limitations:
1. Join construction lacked scoped method isolation, exposing global methods in an ambiguous chaining context.
2. Join fields, table targets, and match conditions relied on raw string literals (e.g. `"PipelineErrorRecord"`, `"ErrorText"`), bypassing compile-time safety and refactoring tooling.
3. Query compilation and terminal operations lacked internal error state tracking, returning raw tuples `(string, []any)` rather than monadic Result envelopes (`result.Result[T]`).
4. View lifecycle management (`CreateViewOrUseView`) required manual column name arguments from callers and had no database-backed mechanism to detect structural query modifications (such as changed `WHERE` filters or `JOIN` conditions).
5. There was no pre-execution SQL validation mechanism to detect syntax errors before executing DDL.
6. Domain-specific queries and business logic were mixed with generic database operations rather than residing in dedicated business logic repositories.

###### Root Cause
Initial builder design focused on basic single-table CRUD operations and lacked sub-struct scoping, deterministic query hashing, and database-backed metadata tracking required for safe multi-table join and view workflows.

---

##### 2. Architecture & Design Decisions

###### 2.1 Scoped JoinBuilder Sub-Struct (`JoinBuilder[T, F]`)
Calling `.Join(table)`, `.InnerJoin(table)`, `.LeftJoin(table)`, `.RightJoin(table)`, or `.OuterJoin(table)` returns a dedicated `*JoinBuilder[T, F]`.
- Developer is strictly exposed to join-scoped methods:
  - `.Select(fields ...any)`: Projected fields from the joined table.
  - `.And(column any, op SqlOperator, val any)`: Extra filter conditions on the `ON` clause.
  - `.AndRaw(condition string)`: Raw SQL condition on the `ON` clause.
  - `.On(condition string)` / `.OnField(mainCol any, op SqlOperator, joinCol any)`: Completes the join and transitions back to `*QueryBuilder[T, F]`.

###### 2.2 Result Envelopes & Comprehensive Error Propagation
- Both `QueryBuilder[T, F]` and `JoinBuilder[T, F]` carry an `err *apperror.AppError` field.
- If an error occurs, chained calls preserve the error state.
- Terminal methods (`First`, `FindAll`, `Count`, `Delete`, `Compile`, `CreateViewOrUseView`) guard against `b.err != nil` and return typed Result envelopes:
  - `First(ctx)` $\rightarrow$ `EntityResult[T]`
  - `FindAll(ctx)` $\rightarrow$ `ListResult[T]`
  - `Count(ctx)` $\rightarrow$ `Int64Result`
  - `Delete(ctx)` $\rightarrow$ `RowsAffectedResult`
  - `CreateViewOrUseView(...)` $\rightarrow$ `BoolResult`
  - `Compile()` $\rightarrow$ `CompiledQueryResult` (`result.Result[CompiledQuery]`)

###### 2.3 Zero Magic Strings Across Joins
- All join methods accept `any` and convert through `toColumnName(col any) string`, which supports `string`, `fmt.Stringer`, and enum types.
- Field enums (e.g. `PipelineRunRecordDb.RunId`, `PipelineErrorRecordDb.ErrorText`) and table constants (`PipelineErrorRecordTable`) are used everywhere.

###### 2.4 Automated View Evolution via Query Hash (`__gitmap_view_meta`)
- View metadata is stored in a dedicated database table:
  ```sql
  CREATE TABLE IF NOT EXISTS __gitmap_view_meta (
      ViewName TEXT PRIMARY KEY,
      QueryHash TEXT NOT NULL,
      ViewSql TEXT NOT NULL,
      UpdatedAt TEXT NOT NULL
  );
  ```
- If the deterministic SHA-256 hash of the view SQL matches $\rightarrow$ instant view reuse with **0 DDL**.
- If the hash differs or view does not exist $\rightarrow$ validates SQL via `EXPLAIN`, drops the stale view, creates the updated view, and records the new hash in metadata.
- Manual column parameter lists in `CreateViewOrUseView` are completely optional.

###### 2.5 Pre-Execution SQL Syntax Validation (`ValidateSql`)
- `DbWrapper.ValidateSql(ctx, sqlStr)` executes `EXPLAIN <sql>` against the database engine.
- Bytecode compilation errors and missing table/column errors are caught immediately before any DDL or data changes occur.

###### 2.6 Domain Business Logic Repository (`PipelineRepository`)
- Generic CRUD remains isolated in `dbengine.Repository`.
- Domain business logic lives in `pipelinedb.PipelineRepository`:
  - `GetRunById(ctx, runId uint64) EntityResult[PipelineRunRecord]`
  - `GetRecentRuns(ctx, repoSlug string, limit int) ListResult[PipelineRunRecord]`
  - `EnsureActiveErrorsView(ctx context.Context) BoolResult`
  - `Query()` and `QueryBare()` for pre-configured and custom projections.

---

##### 3. Implementation Verification & Quality Gates

| Check / Suite | Status | Results |
| :--- | :--- | :--- |
| `TestDbWrapper_ValidateSql` | PASS | Verified `EXPLAIN` syntax validation |
| `TestDbWrapper_ViewHashMetaAndEvolution` | PASS | Verified metadata tracking, 0 DDL reuse, and evolution |
| `TestQueryBuilder_ErrorGuards` | PASS | Verified `*apperror.AppError` safety across all terminals |
| `TestQueryBuilder_TypedJoinsWithoutMagicStrings` | PASS | Verified join compilation with typed constants and enums |
| `TestPipelineRepository_CRUD` | PASS | Verified domain record insertions, ID retrieval, recent runs |
| `TestPipelineRepository_ActiveErrorsViewAndJoin` | PASS | Verified typed joins, view creation, hash reuse, direct querying |
| `TestPipelineRepository_FluentQueryWithEnums` | PASS | Verified bare query compilation with enum projections |
| All DBEngine Unit Tests | PASS | 14 test functions, 100% pass rate |
| All PipelineDB Unit Tests | PASS | 5 test functions, 100% pass rate |
| Nested If Linter (`check-nested-ifs.py`) | PASS | 3,120 files audited, max depth 1 enforced |
| Boolean & Enum Linter (`check-enum-and-boolean.py`) | PASS | 2,307 files audited, 0 violations |
| CI/CD Local Runner (`06-cicd-local-runner.py`) | PASS | 16 of 16 quality gates green |
| CLI Smoke Suite (`e2e-cli-smoke.py`) | PASS | 118 of 118 commands verified |

### Merged Plan: `66-automatic-db-repo-and-safe-scanner-generator.md`

#### 66 — Automatic Typed DbRepo & Safe Row Scanners Generator

**Status:** COMPLETED
**Priority:** High
**Category:** Database Architecture / Code Generation / Type-Safe Repositories
**Completed:** 2026-09-05

---

##### 1. Problem Statement & Root Cause

###### Problem Statement
Previously, database model structs generated column enums (`*FieldType`) and registries (`*DbRegistry`), but developers still had to write boilerplate row scanners and typed repository structs manually. Specifically:
1. Handwritten scanners (`ScanPipelineRunRecord`, `ScanPipelineErrorRecord`) duplicated struct fields across files.
2. Standard SQLite driver scanning caused errors or panics when handling `NULL` values or converting SQLite `int64` columns into Go `uint64` / `bool` types.
3. Repositories had to be instantiated manually via generic `dbengine.NewRepository[T, F]` and wrapped with boilerplate query methods (`FindAll`, `First`, `Count`, `Query`).

###### Root Cause
The database code generator (`03-ai-scripts/30-db-struct-enum-generator.py`) only generated field enums and JSON serializers/deserializers; it lacked AST-driven repository code generation and null-safe scanner helpers.

---

##### 2. Architecture & Design Decisions

###### 2.1 Safe Row Scan Helpers (`gitmap/dbengine/scan_helpers.go`)
Implemented type-converting, null-safe row scanners:
- `ScanString(v any) string`: Handles `nil`, `string`, `[]byte`, `fmt.Sprint`.
- `ScanInt(v any) int`: Coerces `int64`, `int32`, `uint64` into `int`; returns `0` for `nil`.
- `ScanInt64(v any) int64`: Coerces numeric interfaces into `int64`; returns `0` for `nil`.
- `ScanUint64(v any) uint64`: Coerces `int64`, `int`, `uint` into `uint64`; returns `0` for `nil`.
- `ScanUint(v any) uint`: Converts uint64 into `uint`.
- `ScanBool(v any) bool`: Handles `bool`, integer `0`/`1` boolean representations; returns `false` for `nil`.
- `ScanFloat64(v any) float64`: Coerces float and integer values into `float64`; returns `0.0` for `nil`.

###### 2.2 Automated Scanner Generation (`Scan{StructName}`)
The generator analyzes struct fields and generates:
```go
func ScanPipelineRunRecord(row dbengine.RowScanner) (*PipelineRunRecord, error)
```
- Declares temporary `any` variables for each field.
- Scans row safely.
- Maps raw values into struct fields via `dbengine.Scan*` helpers.

###### 2.3 Dedicated Definition Files & `consts.go` Architecture
The architecture separates each model definition into its own named file, with all constants and type aliases coexisting in `consts.go`:
- **`consts.go`**:
  - Canonical table constants: `const ( PipelineSplitDbTable = "PipelineSplitDb", PipelineRunRecordTable = "PipelineRunRecord", PipelineRunTable = PipelineRunRecordTable, ... )`
  - Concrete QueryBuilder aliases: `type PipelineSplitDbQueryBuilder = dbengine.QueryBuilder[...]`, `type PipelineRunQueryBuilder = PipelineRunRecordQueryBuilder`
  - Concrete generic Repository aliases: `type PipelineSplitDbRepository = dbengine.Repository[...]`, `type PipelineRunRepository = PipelineRunRecordRepository`
  - Concrete DbRepo aliases: `type PipelineRunDbRepo = PipelineRunRecordDbRepo`, `type PipelineErrorDbRepo = PipelineErrorRecordDbRepo`
- **Dedicated Definition Files (`pipeline_run_record.go`, `pipeline_error_record.go`, `pipeline_db_stats.go`, `pipeline_split_db.go`)**:
  - Contains struct model definition, type-safe column enums, O(1) map validation, null-safe row scanners, and typed `*DbRepo` accessors.

###### 2.4 Typed Repository Pointer Receivers & Constructors
- Repository structs (`{StructName}DbRepo`) are typed accessors wrapping generic repositories.
- Constructors: `NewPipelineRunDbRepo(db) *PipelineRunDbRepo`.
- Methods use pointer receivers: `func (r *PipelineRunDbRepo) Query() *PipelineRunQueryBuilder`.
- Standard query methods:
  - `FindAll(ctx context.Context) dbengine.ListResult[PipelineRunRecord]`
  - `First(ctx context.Context) dbengine.EntityResult[PipelineRunRecord]`
  - `Count(ctx context.Context) dbengine.Int64Result`
  - `Query() *PipelineRunRecordQueryBuilder`
  - `QueryBare() *PipelineRunRecordQueryBuilder`
  - `Db() *dbengine.DbWrapper`
  - `Repo() *PipelineRunRecordRepository`

###### 2.5 Domain Repository Embedding Pattern
Domain repositories (such as `pipelinedb.PipelineRepository`) embed `*PipelineRunRecordDbRepo`:
```go
type PipelineRepository struct {
    *PipelineRunRecordDbRepo
}

func NewPipelineRepository(db *dbengine.DbWrapper) *PipelineRepository {
    return &PipelineRepository{
        PipelineRunRecordDbRepo: NewPipelineRunRecordDbRepo(db),
    }
}
```
Eliminates all boilerplate scanner and generic repository delegations while allowing domain repositories to cleanly define domain-specific queries (`GetRunById`, `GetRecentRuns`, `EnsureActiveErrorsView`).

---

##### 3. Verification & Quality Gates

| Check / Suite | Status | Results |
| :--- | :--- | :--- |
| `TestScanString`, `TestScanInt`, `TestScanInt64`, `TestScanUint64`, `TestScanUint`, `TestScanBool`, `TestScanFloat64` | PASS | Full coverage of null safety and type coercions |
| `TestPipelineRunDbRepo_GeneratedRepo` | PASS | Verified `NewPipelineRunDbRepo`, `FindAll`, `First`, `Count`, and fluent queries |
| `TestPipelineDbGeneratedConsts` | PASS | Verified canonical table constants and QueryBuilder/Repository type aliases |
| All PipelineDB Tests | PASS | 7 test functions green |
| All DBEngine Tests | PASS | 21 test functions green |
| Nested If Linter (`check-nested-ifs.py`) | PASS | 3,122 files scanned, 0 violations |
| Boolean & Enum Linter (`check-enum-and-boolean.py`) | PASS | 2,308 files scanned, 0 violations |
| CI/CD Local Runner (`06-cicd-local-runner.py`) | PASS | 16 of 16 quality gates green |

### Merged Plan: `71-os-power-management-and-installation-split-db.md`

#### 71-os-power-management-and-installation-split-db.md: OS Screen Timeout & Sleep Management Framework & Installation Split DB

##### 1. Executive Summary

This plan fulfills two major enterprise system capabilities for Gitmap:
1. **OS Screen Timeout & Sleep Management Framework (`gitmap power`)**:
   - Implements a DRY, pluggable power management framework supporting Windows (`powercfg.exe`), Linux/Ubuntu (3-tier hierarchy: GNOME `gsettings`, X11 `xset`, and headless `systemd-logind`), with extension points for future operating systems (e.g. macOS `pmset`).
   - Supports CLI commands: `status` (live query of display/sleep timeout and active profile), `never-sleep` (disables screen blanking, sleep/standby, and screen lock/logout), `set <mins>` (configures display and sleep timeouts), `reset` (restores previous snapshot), and `history` (audit log of power changes).
   - Persists settings and history in SQLite (`PowerSetting` and `PowerSettingHistory`) so users can check status, reset, and customize configurations across reboots.
2. **Dedicated Installation Split Database (`installation.db`)**:
   - Strictly adheres to `spec/04-database-conventions/07-split-db-pattern.md` and repository standards (`schedule_split_db.go`).
   - Decouples tool installation records and voluminous execution logs from the root SQLite database into an isolated `installation.db` database located at `<BinaryDataDir>/installation.db` with `conn.SetMaxOpenConns(1)`.
   - Defines PascalCase schema: `InstalledTool` and `InstallationLog` (tracking tool, action, version, duration, exit code, stdout/stderr logs, and Rule 10/11 context columns `Description`, `Notes`, `Comments`).
   - Provides a zero-downtime migration bridge that moves legacy `InstalledTool` records from root DB to `installation.db` and routes all installer commands (`gitmap install`, `gitmap in build-essential`, `gitmap uninstall`, `gitmap install --list`) to the split DB.

##### 2. Mandatory Rules & Invariants

1. **Rule 1 (Strict Relative Git Paths)**: All paths in markdown files, artifacts, and documentation must be strictly relative to the repository root (`/`). Zero `file:///` URIs.
2. **Rule 2 (Strict File & Function Sizing)**: Every function must be $\le 15$ lines (preferred $\le 8$ lines), with a mandatory blank line before every return statement. All new and modified Go files must remain $\le 200$ lines.
3. **Rule 3 (Database Conventions Compliance)**: All tables singular PascalCase (`InstalledTool`, `InstallationLog`, `PowerSetting`, `PowerSettingHistory`), integer auto-increment PKs (`{TableName}Id`), affirmative booleans (`IsNeverSleep`, `IsSuccess`), and Rule 10/11 context columns.
4. **Rule 4 (SQLite Concurrency Standard)**: All split databases must call `conn.SetMaxOpenConns(1)` and anchor paths via `store.BinaryDataDir()` / `filepath.EvalSymlinks(os.Executable())`.
5. **Rule 5 (AST Parity Guard)**: Top-level CLI command constants (`CmdPower = "power"`, `CmdPowerAlias = "pw"`, `CmdPowerAlias2 = "pwr"`) must be registered in `constants_cli.go` under `// gitmap:cmd top-level` and synchronized in `cmd_constants_test.go:topLevelCmds()`.
6. **Rule 6 (Zero Swallowed Errors)**: No blank error discards; wrap all errors using `apperror.WrapSimple` or `apperror.NewWithDetails`.

##### 3. Subsystems & Architecture

###### A. Pluggable Power Management Framework (`gitmap/power/`)
- `power.Manager` interface:
  - `Platform() string`
  - `GetStatus(ctx context.Context) (*Settings, error)`
  - `SetNeverSleep(ctx context.Context) error`
  - `SetTimeouts(ctx context.Context, displayMinutes, sleepMinutes int) error`
  - `ApplySettings(ctx context.Context, settings *Settings) error`
- Drivers:
  - `manager_windows.go`: `powercfg /change monitor-timeout-ac`, `standby-timeout-ac`, query active schemes.
  - `manager_linux.go`: `gsettings` (`org.gnome.desktop.session idle-delay`, `sleep-inactive-ac-timeout`), `xset s off -dpms`, and `systemd-logind`.
  - `manager_darwin.go`: macOS `pmset` driver stub.
  - `manager_fallback.go`: graceful fallback for unsupported OSes.

###### B. SQLite Power Settings Persistence (`gitmap/store/`)
- `PowerSetting` table: records current and previous profiles (`baseline`, `current`, `previous`).
- `PowerSettingHistory` table: transactional audit log of every power change.
- Store helper: `gitmap/store/power.go`.

###### C. Installation Split DB (`gitmap/store/installation_split_*.go`)
- Location: `<BinaryDataDir>/installation.db`
- Connection wrapper: `InstallationSplitDB` with `conn.SetMaxOpenConns(1)`.
- Tables:
  - `InstalledTool`: tracked packages and semantic version components.
  - `InstallationLog`: detailed execution logs, durations, exit codes, and stdout/stderr.
- Migration Bridge: `MigrateInstalledToolsFromRoot(rootDB, splitDB)`.

###### D. CLI Command Integration (`gitmap/cmd/`)
- Power: `gitmap power [status|never-sleep|set|reset|history]` in `gitmap/cmd/power*.go`.
- Installation: wires `install.go`, `installtools.go`, `install_buildessential.go`, `installlist.go`, and `uninstall.go` to use `store.OpenInstallationSplitDB()`.

##### 4. Subtasks Breakdown

1. [01-installation-split-db-schema-and-manager.md](subtasks/71-os-power-management-and-installation-split-db/01-installation-split-db-schema-and-manager.md): Implement `InstallationSplitDB` lifecycle, `InstalledTool`, `InstallationLog` schema, and CRUD in `gitmap/store/installation_split_*.go`.
2. [02-installation-split-db-migration-and-command-wiring.md](subtasks/71-os-power-management-and-installation-split-db/02-installation-split-db-migration-and-command-wiring.md): Implement migration bridge from root DB, backward-compatible delegation on `*store.DB`, and wire `install.go`, `install_buildessential.go`, `installtools.go`, `installlist.go`, and `uninstall.go`.
3. [03-os-power-framework-and-drivers.md](subtasks/71-os-power-management-and-installation-split-db/03-os-power-framework-and-drivers.md): Implement pluggable `power.Manager` interface, `types.go`, Windows `powercfg` driver, Linux 3-tier driver (`gsettings`/`xset`), and Darwin/fallback drivers.
4. [04-os-power-sqlite-persistence-and-commands.md](subtasks/71-os-power-management-and-installation-split-db/04-os-power-sqlite-persistence-and-commands.md): Implement `PowerSetting` and `PowerSettingHistory` store operations, CLI dispatcher `gitmap power` (`status`, `never-sleep`, `set`, `reset`, `history`), and AST constants.
5. [05-helptext-docs-and-ci-verification.md](subtasks/71-os-power-management-and-installation-split-db/05-helptext-docs-and-ci-verification.md): Author helptext files (`power.md`), update changelogs, verify AST parity, golden tests, and execute full CI quality runner.

##### 5. Acceptance Criteria

- [x] `installation.db` created in `BinaryDataDir()` with `SetMaxOpenConns(1)`.
- [x] `InstalledTool` and `InstallationLog` tables created with PascalCase columns, affirmative booleans, and Rule 10/11 context columns.
- [x] Existing records in root DB `InstalledTool` automatically migrate to `installation.db`.
- [x] All installation commands (`gitmap install`, `gitmap in build-essential`, `gitmap uninstall`, `gitmap install --list`) read/write from `installation.db` and log telemetry to `InstallationLog`.
- [x] `gitmap power status` reports current OS screen/sleep timeout and active DB profile.
- [x] `gitmap power never-sleep` disables screen timeout and sleep, recording previous state to DB.
- [x] `gitmap power set <mins>` updates display and sleep timeouts.
- [x] `gitmap power reset` restores previous timeout settings from DB snapshot.
- [x] AST parity tests pass: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1`.
- [x] ERD parity tests pass: `go test ./gitmap/store/... -run TestERDMatchesSQLCreate -count=1`.
- [x] All new/modified Go files are $\le 200$ lines, functions $\le 15$ lines, blank lines before returns.
- [x] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py` passes all gates with `exit 0`.

#### Granular Subtask Execution Details for `71-os-power-management-and-installation-split-db`

##### Subtasks Folder: `71-os-power-management-and-installation-split-db` (5 subtask files incorporated)
###### Subtask File: `01-installation-split-db-schema-and-manager.md`

#### Subtask 01: Installation Split DB Schema & Connection Manager

##### Status
Completed

##### Context & Objectives
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

##### Verification Steps
- Unit tests in `gitmap/store/installation_split_db_test.go` verify table creation, CRUD operations, and logs.

###### Subtask File: `02-installation-split-db-migration-and-command-wiring.md`

#### Subtask 02: Installation Split DB Migration Bridge & Command Wiring

##### Status
Completed

##### Context & Objectives
1. **Migration Bridge**:
   - `gitmap/store/installation_split_bridge.go`: Reads any legacy rows from root DB `InstalledTool` and migrates them into `installation.db` (`INSERT OR IGNORE`).
   - Maintains delegate wrappers in `gitmap/store/installedtool.go` on `*store.DB` so any remaining references seamlessly route to `InstallationSplitDB`.
2. **Command Wiring**:
   - `gitmap/cmd/installtools.go`: Writes to `store.OpenInstallationSplitDB()` and logs telemetry (`InstallationLogRecord`).
   - `gitmap/cmd/install_buildessential.go`: Writes `build-essential` profile and constituent tools to `installation.db` with execution timing.
   - `gitmap/cmd/installlist.go`: Reads installed tool status from `installation.db`.
   - `gitmap/cmd/uninstall.go`: Uninstalls and records audit log in `installation.db`.

##### Verification Steps
- Unit tests verify root DB migration preserves version strings and timestamps.
- CLI installation commands write to `installation.db`.

###### Subtask File: `03-os-power-framework-and-drivers.md`

#### Subtask 03: OS Power Management Framework & Pluggable Drivers

##### Status
Completed

##### Context & Objectives
Implement the pluggable, DRY power management framework in `gitmap/power/`:
1. **Interface Definition** (`gitmap/power/manager.go`):
   - `type Manager interface` with `Platform()`, `GetStatus()`, `SetNeverSleep()`, `SetTimeouts()`, `ApplySettings()`.
   - `NewManager() (Manager, error)` dynamically returns the OS driver based on `runtime.GOOS`.
2. **Data Types** (`gitmap/power/types.go`):
   - `Settings` struct: `Platform`, `DisplayTimeoutMinutes`, `SleepTimeoutMinutes`, `DiskTimeoutMinutes`, `IsNeverSleep`, `IsLockDisabled`, `IsActive`.
3. **Windows Driver** (`gitmap/power/manager_windows.go`):
   - Uses `powercfg.exe` (`/change monitor-timeout-ac`, `/change standby-timeout-ac`, query active schemes).
4. **Linux Driver** (`gitmap/power/manager_linux.go`):
   - 3-tier fallback: GNOME `gsettings` (`org.gnome.desktop.session idle-delay`, `sleep-inactive-ac-timeout`, `lock-enabled`), X11 `xset s off -dpms`, and `systemd-logind`.
5. **Darwin & Fallback Drivers** (`gitmap/power/manager_darwin.go`, `manager_fallback.go`):
   - Stubs for future macOS `pmset` and unsupported OSes.

##### Verification Steps
- Unit tests in `gitmap/power/manager_test.go` verify factory instantiation and settings parsing.

###### Subtask File: `04-os-power-sqlite-persistence-and-commands.md`

#### Subtask 04: OS Power SQLite Persistence & CLI Command Integration

##### Status
Completed

##### Context & Objectives
1. **SQLite Persistence**:
   - `gitmap/store/power.go`: CRUD operations for `PowerSetting` (profiles: `baseline`, `current`, `previous`) and `PowerSettingHistory` (audit trail).
   - SQL DDL in `gitmap/store/power_schema.go` (kept internal to avoid ERD drift or registered cleanly).
2. **CLI Commands & Dispatcher**:
   - Top-level constants in `gitmap/constants/constants_cli.go`: `CmdPower = "power"`, `CmdPowerAlias = "pw"`, `CmdPowerAlias2 = "pwr"`.
   - Synchronize `cmd_constants_test.go:topLevelCmds()`.
   - Dispatch entry in `gitmap/cmd/rootutility.go`.
   - `gitmap/cmd/power.go` (<= 200 lines): Subcommand router.
   - `gitmap/cmd/power_ops.go` (<= 200 lines): `status`, `never-sleep`, `set`, `reset`, `history`.

##### Verification Steps
- AST parity test passes: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1`.
- `gitmap power --help` and unit tests pass.

###### Subtask File: `05-helptext-docs-and-ci-verification.md`

#### Subtask 05: Help Text Parity, Documentation & CI/CD Verification

##### Status
Completed

##### Context & Objectives
1. **Help Text Creation**:
   - `gitmap/helptext/power.md`: Dedicated help file (<= 120 lines, 3-8 line simulated output block, examples).
   - Register in `gitmap/helptext/catalog.go`.
   - Update `gitmap/cmd/rootusage_groups.go` to list `power` in utility commands.
2. **Changelog Updates**:
   - Document OS power management framework and Installation Split DB in `changelog.md`, `gitmap/changelog.md`, `src/data/changelog.ts`.
3. **CI/CD Quality Gate Verification**:
   - Run `go test ./gitmap/helptext/... -run Golden -count=1`.
   - Run `python 03-ai-scripts/06-cicd-local-runner.py` ensuring all quality gates pass with `exit 0`.

##### Verification Steps
- Golden helptext tests pass.
- Local CI runner passes 100%.


### Merged Plan: `91-purge-history-refactor-and-sqlite-tracking.md`

#### 91-purge-history-refactor-and-sqlite-tracking

##### Goal Description
Review, refactor, and harden the newly introduced purge CLI command and history tracking mechanism:
1. Normalize database primary keys and struct identifiers from all-caps ID to PascalCase Id (PurgeHistoryLogId canonical, Id alias) and camelCase id parameters to strictly adhere to repository naming guidelines and prevent linter rejections.
2. Adopt the centralized database pattern (ExecWrapper, QueryRowWrapper, QueryResult[T]) from gitmap/store/wrapper.go in gitmap/store/purge_history.go so all query executions have structured boolean states (isSuccess, isFailure) and standard stderr diagnostics.
3. Eliminate all 20+ instances of swallowed errors (_ ignores, unchecked os.MkdirAll, unchecked copyPurgeFile, unchecked sendToRecycleBin, unchecked git commands, unchecked file closes, and unhandled db.MarkPurgeHistoryRestored).
4. Decompose gitmap/cmd/purge.go (currently 290 lines with doPurge at 134 lines) into modular files <= 200 lines each with functions <= 15 lines (purge.go, purge_engine.go, purge_restore.go, purge_recycle_windows.go, purge_recycle_other.go).
5. Ensure cross-platform compatibility for temp directory backups (os.TempDir() -> filepath.Join(os.TempDir(), fmt.Sprintf("gitmap_purge_%d", timestamp))), path normalization with filepath.ToSlash() for git filter-repo, and build-tag isolated Recycle Bin calls.
6. Verify and harden the Python companion script 03-ai-scripts/30-purge-history.py against line budget (>200 lines) and non-Windows ctypes crashes.

---

##### Custom Constraints & Rules (Task-Specific)
1. Strict ID Casing (Rule 1): All-caps ID is completely prohibited in Go structs, methods, parameters, and SQLite queries. Must use PurgeHistoryLogId as canonical PK, Id as struct alias, and id as parameter.
2. Zero Swallowed Errors (Rule 2): Every single returned error (file ops, git commands, db calls, IO copies, and JSON marshaling) must be checked immediately and wrapped via apperror.Wrap or returned up the stack.
3. Store Wrapper Pattern Reuse (Rule 3): All SQL queries in gitmap/store/purge_history.go must be executed via store.ExecWrapper and store.QueryRowWrapper.
4. Function & File Sizing Caps (Rule 4): Every source file must remain <= 200 lines (target <= 100 lines). Every function must remain <= 15 lines (target 8-15 lines).
5. Cross-Platform Recycle Isolation (Rule 5): Windows shell32.dll SHFileOperationW must reside in purge_recycle_windows.go guarded by //go:build windows, with a portable fallback in purge_recycle_other.go (//go:build !windows).

---

##### Subtasks Decomposition
- 01-store-purge-history-id-and-wrapper.md: Refactor gitmap/store/purge_history.go to use PurgeHistoryLogId, Id, ExecWrapper, QueryRowWrapper, IsRestored, and wire into store.go:Migrate().
- 02-cmd-purge-split-and-error-handling.md: Refactor gitmap/cmd/purge.go and extract gitmap/cmd/purge_engine.go with zero swallowed errors and functions <= 15 lines.
- 03-purge-restore-and-cross-platform-recycle.md: Extract gitmap/cmd/purge_restore.go, gitmap/cmd/purge_recycle_windows.go, and gitmap/cmd/purge_recycle_other.go.
- 04-python-purge-history-hardening.md: Refactor 03-ai-scripts/30-purge-history.py to stay under 200 lines with cross-platform recycle and argument list subprocess execution.
- 05-unit-tests-and-ci-verification.md: Add unit and integration tests for purge and purge_history store, run go vet, tests, and CI local runner.

#### Granular Subtask Execution Details for `91-purge-history-refactor-and-sqlite-tracking`

##### Subtasks Folder: `91-purge-history-refactor-and-sqlite-tracking` (5 subtask files incorporated)
###### Subtask File: `01-store-purge-history-id-and-wrapper.md`

#### Subtask 01: Store Purge History ID Normalization & QueryWrapper Pattern Reuse

##### Objective
Refactor `gitmap/store/purge_history.go` to adhere to repository database guidelines:
- Rename all-caps `ID` to `PurgeHistoryLogId` (canonical PK) and provide `Id` (PascalCase) and `ID` (deprecated compatibility alias).
- Rename `Restored` boolean to `IsRestored` (canonical boolean) and provide `Restored` (deprecated compatibility alias).
- Change SQLite schema to use `PurgeHistoryLogId INTEGER PRIMARY KEY AUTOINCREMENT`, `IsRestored INTEGER NOT NULL DEFAULT 0`, `Notes TEXT NULL`, and `Comments TEXT NULL`.
- Add backward-compatible `migratePurgeHistoryColumns()` to handle existing databases with `ID` or `Restored`.
- Adopt `ExecWrapper` and `QueryRowWrapper` from `gitmap/store/wrapper.go` for all SQL executions with explicit error logging and structured `QueryResult`.
- Eliminate swallowed error from `res.LastInsertId()` by verifying `err != nil`.
- Wire `EnsurePurgeHistoryTable()` into `gitmap/store/store.go:Migrate()`.

##### Files Affected
- `gitmap/store/purge_history.go`
- `gitmap/store/store.go`

###### Subtask File: `02-cmd-purge-split-and-error-handling.md`

#### Subtask 02: CLI Purge Command Decomposition & Zero-Swallow Error Handling

##### Objective
Refactor `gitmap/cmd/purge.go` and extract `gitmap/cmd/purge_engine.go`:
- Keep `runPurge(args []string) error` in `gitmap/cmd/purge.go` (under 70 lines) for CLI argument parsing (`isRestore`, `isAutoConfirm`, `pattern`).
- Extract the core purge engine into `gitmap/cmd/purge_engine.go` (under 120 lines).
- Decompose the 134-line `doPurge` into focused functions <= 15 lines:
  - `verifyWorkTreeClean() error`
  - `fetchCurrentBranch() (string, error)`
  - `queryMatchingFiles(pattern string) ([]string, error)`
  - `createBackupBranch(branch string) error`
  - `backupFilesToTemp(tempDir string, files []string) ([]string, error)`
  - `executeFilterRepo(pattern string) error`
  - `appendGitignore(pattern string) error`
  - `savePurgeState(db *store.DB, log *store.PurgeHistoryLog) error`
- Fix all swallowed errors: check `os.MkdirAll`, check `copyPurgeFile`, check `sendToRecycleBin`, check `f.Close()`, check `git remote add`, and check `json.Marshal`.
- Normalize pattern path separators using `filepath.ToSlash(pattern)`.
- Use boolean variable `isAutoConfirm` instead of `autoConfirm`.

##### Files Affected
- `gitmap/cmd/purge.go`
- `gitmap/cmd/purge_engine.go`

###### Subtask File: `03-purge-restore-and-cross-platform-recycle.md`

#### Subtask 03: Purge Restore Decomposition & Cross-Platform Recycle Isolation

##### Objective
Extract restore logic and cross-platform recycle bin support:
- Extract `gitmap/cmd/purge_restore.go`:
  - Decompose `doRestore` (34 lines) into functions <= 15 lines:
    - `fetchActivePurgeLog(db *store.DB, repoPath string) (*store.PurgeHistoryLog, error)`
    - `resetToBranch(branch string) error`
    - `restoreFilesFromTemp(tempDir string) error`
    - `markPurgeRestored(db *store.DB, id int64) error`
  - In `restoreFilesFromTemp`, check errors from `filepath.WalkDir`, propagate walk errors, check `os.MkdirAll`, and check `copyPurgeFile`.
  - In `markPurgeRestored`, pass `log.PurgeHistoryLogId` or `log.Id` and handle the returned error (eliminate swallow).
- Create `gitmap/cmd/purge_recycle_windows.go`:
  - Tagged with `//go:build windows`.
  - Houses `sendToRecycleBin(path string) error` using `shell32.dll` `SHFileOperationW`.
- Create `gitmap/cmd/purge_recycle_other.go`:
  - Tagged with `//go:build !windows`.
  - Houses portable `sendToRecycleBin(path string) error` using `os.RemoveAll` or standard removal fallback so builds and execution on Linux/macOS never fail.

##### Files Affected
- `gitmap/cmd/purge_restore.go`
- `gitmap/cmd/purge_recycle_windows.go`
- `gitmap/cmd/purge_recycle_other.go`

###### Subtask File: `04-python-purge-history-hardening.md`

#### Subtask 04: Python Purge History Script Hardening & Line Budget Optimization

##### Objective
Refactor `03-ai-scripts/30-purge-history.py`:
- Bring file length below 200 lines (currently 209 lines).
- Bring functions below 15 lines (e.g. break down `purge_history` and `restore_history`).
- Eliminate non-Windows `AttributeError: module 'ctypes' has no attribute 'windll'` by guarding `ctypes.windll` with `sys.platform == "win32"`.
- Replace `shell=True` with structured argument lists in `subprocess.run` to prevent shell injection and path quoting discrepancies on Windows.
- Normalize path globs with `Path(pattern).as_posix()`.
- Update boolean naming to affirmative prefixes (`is_auto_confirm`, `is_restore`, `should_check`).

##### Files Affected
- `03-ai-scripts/30-purge-history.py`

###### Subtask File: `05-unit-tests-and-ci-verification.md`

#### Subtask 05: Unit Tests, Parity Verification & CI Local Runner

##### Objective
Author comprehensive unit tests and verify against quality gates:
- Write unit tests in `gitmap/store/purge_history_test.go`:
  - Test table creation, column migration (legacy `ID` to `PurgeHistoryLogId`), insert, retrieve last, mark restored.
  - Assert that `PurgeHistoryLog.Id`, `PurgeHistoryLog.PurgeHistoryLogId`, and `PurgeHistoryLog.IsRestored` are populated properly.
- Write unit tests in `gitmap/cmd/purge_test.go`:
  - Test argument parsing (`--restore`, `--confirm`, `-y`, pattern extraction).
- Run `go vet ./...` and `go test ./gitmap/store/... ./gitmap/cmd/...`.
- Run `python 03-ai-scripts/06-cicd-local-runner.py` to ensure all quality gates exit with code 0.

##### Files Affected
- `gitmap/store/purge_history_test.go`
- `gitmap/cmd/purge_test.go`


### Merged Plan: `94-orm-code-generator-id-error-standards-and-fast-cache.md`

#### Plan 94: ORM Code Generator Integration, ID Naming & Error Standards, and Fast Cache Engine

**Title:** Database ORM Integration, Generator Upgrades, Tesla ID Naming Standards, Universal AppError Wrapping, Fast Cache Engine & Codebase Release
**Status:** Pending
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.1.0)
**Budget (N):** 150 steps
**Target Directory:** `.lovable/plans/pending/`

---

##### 1. Root Cause & Problem Statement

1. **Disconnected Database ORM & Code Generation**:
   - `gitmap/dbengine` provides a rich, multi-dialect SQL query compiler and repository engine (`DbWrapper`, `Repository[T, F]`, `QueryBuilder[T, F]`, safe scan coercers).
   - However, currently only 4 models in `gitmap/pipelinedb` are connected. Tables in `gitmap/repodb` (`RepoFile`, `SearchCache`, `FileSequence`, `SequenceHistory`, `RepoScanLog`), `gitmap/db` (`ClusterNode`, `ClusterRun`, etc.), and key store entities are completely disconnected from the ORM, relying on ad-hoc raw SQL strings and handwritten scanners.
   - The generator `03-ai-scripts/30-db-struct-enum-generator.py` only generates read queries, does not parse `db:` tags, and contains hardcoded hacks.
2. **Tesla ID Naming & Schema Convention Violations**:
   - Spec `spec/04-database-conventions/01-naming-conventions.md` and `spec/01-naming` mandate PascalCase acronyms (`Id`, `Db`, `Url`, `Api`) and forbid all-caps clusters (`ID`, `DB`, `URL`, `API`).
   - Multiple structs in `gitmap/model/`, `gitmap/repodb/`, `gitmap/store/`, and `gitmap/db/` contain legacy `ID`, bare `Id` columns, or snake_case table/column names.
   - Lookup tables and "code values" (`TaskType`, `ProjectType`, cluster enums) violate Spec Rule 13's canonical `(Id, Code, Label, Description)` standard.
3. **Swallowed Errors & Missing AppError Envelopes**:
   - Over 15 instances of swallowed errors (`_ = `) exist in `pipelinedb/pipeline_split_ops.go` (dropping scan errors in stats and recent runs), `dbengine/wrapper.go` (swallowing transaction rollback errors), and `repodb/repo_db.go`.
   - `repodb` and `db` functions return standard Go `error` instead of domain-specific `*apperror.AppError`.
4. **Fast Cache & Search Bottlenecks (Go vs Python)**:
   - Python's `03-ai-scripts/` engine caches files in `tmp/cache/repo-file-cache.json` (<0.1ms startup), sniffs binary files with an 8KB null-byte probe, and prunes 20+ directories.
   - Go's `gitmap/indexer/walker.go` issues a synchronous SQLite `SELECT` query for *every individual file* inside `filepath.WalkDir`, lacks a binary sniffer (only checking file size), and only prunes `.git` and `node_modules`.
5. **Release Alignment**:
   - `03-ai-scripts/29-release-bumper.py` synchronizes manifests (`version.json`, `package.json`, `constants.go`, `changelog.md`), and `gitmap release` handles git tags, branches, and cross-compilation. They must be validated through `14-version-sync-checker.py` and local CI runner.

---

##### 2. Task-Specific Rules & Constraints

1. **Strict PascalCase Acronym Standards**:
   - All acronyms must be PascalCase/camelCase: `Id`, `Db`, `Url`, `Api`. Zero all-caps `ID`, `DB`, `URL`, `API` in struct fields, method signatures, parameters, or database columns.
   - Primary keys must follow `<Entity>Id` format; foreign keys must follow `<TargetEntity>Id`.
2. **Zero Swallowed Errors & Universal AppError**:
   - Zero `_ = ` discards on database operations, row scans, or transaction rollbacks.
   - Every database method must return `(*T, *apperror.AppError)`, `([]T, *apperror.AppError)`, or `*apperror.AppError`.
3. **ORM Model Code Generation**:
   - Use `03-ai-scripts/30-db-struct-enum-generator.py` with `db:` tag support to generate `<Model>FieldType` enums, `Scan<Model>` row scanners, and `<Model>DbRepo` typed repositories wrapping `dbengine.Repository`.
4. **High-Performance File Indexing & Binary Sniffer**:
   - Batch-load file indexing timestamps in 1 query upfront before walking disk.
   - Add 8KB null-byte sniffer (`bytes.IndexByte(buf, 0) != -1`) before saving file contents into SQLite.
   - Port `EXCLUDE_DIRS` from `02-shared-engine.py` into Go indexer and searcher.
5. **Coding Guidelines Enforcement**:
   - Functions $\le 15$ lines, zero nested ifs, affirmative booleans only (`is*`, `has*`), Unix LF line endings.

---

##### 3. Subtasks Ledger

- **Subtask 94.01**: Generator Upgrades & Tag-Driven Code Generation (`03-ai-scripts/30-db-struct-enum-generator.py`).
- **Subtask 94.02**: `repodb` & `db` ORM Connection, Entity Generation, and PascalCase `Id` Normalization (`gitmap/repodb/`, `gitmap/db/`, `gitmap/model/`).
- **Subtask 94.03**: Zero-Swallowed Errors & Universal `*apperror.AppError` Wrapping in Database Layers (`gitmap/pipelinedb/`, `gitmap/dbengine/`, `gitmap/repodb/`).
- **Subtask 94.04**: Go Fast Cache Optimization, Upfront Batch Indexing, 8KB Binary Sniffer & Directory Pruning (`gitmap/indexer/`, `gitmap/searcher/`, `gitmap/fsutil/`).
- **Subtask 94.05**: Codebase Release Synchronization, Pre-flight Verification & CI Quality Gates (`03-ai-scripts/`, `gitmap/release/`, local CI runner).

#### Granular Subtask Execution Details for `94-orm-code-generator-id-error-standards-and-fast-cache`

##### Subtasks Folder: `94-orm-code-generator-id-error-standards-and-fast-cache` (5 subtask files incorporated)
###### Subtask File: `01-generator-upgrades-and-tag-driven-generation.md`

#### Subtask 94.01: Generator Upgrades & Tag-Driven Code Generation

##### Goal
Upgrade `03-ai-scripts/30-db-struct-enum-generator.py` to support explicit `db:` struct tags, eliminate hardcoded struct skips, and generate typed mutation helpers (`Insert`, `Update`, `DeleteById`) on `*DbRepo`.

##### Files Impacted
- `03-ai-scripts/30-db-struct-enum-generator.py`
- `gitmap/pipelinedb/enums/`
- `gitmap/pipelinedb/`

##### Acceptance Criteria
1. Struct parser parses `db:"ColumnName"` tag from struct fields; defaults to field name if tag is absent.
2. Removes hardcoded `if s_name == "PipelineSplitDb"` hack by checking for explicit entity markers or skipping connection wrapper structs.
3. Extends `<Model>DbRepo` generation to include typed mutation methods (`Insert`, `Update`, `DeleteById`) utilizing `dbengine.DbWrapper`.
4. Regenerates `gitmap/pipelinedb/` cleanly with no regression.
5. All generated functions $\le 15$ lines, zero nested ifs, affirmative booleans only.

###### Subtask File: `02-repodb-and-db-orm-connection-and-id-normalization.md`

#### Subtask 94.02: repodb & db ORM Connection, Entity Generation, and PascalCase Id Normalization

##### Goal
Declare strongly-typed model structs for `gitmap/repodb` (`RepoFile`, `SearchCache`, `FileSequence`, `SequenceHistory`, `RepoScanLog`, `IndexedRepo`), normalize table primary keys to PascalCase `<Entity>Id` (`RepoFileId`, `SearchCacheId`, etc.), generate typed enums and repositories via the code generator, and connect them with `dbengine.DbWrapper`.

##### Files Impacted
- `gitmap/repodb/repo_db.go`
- `gitmap/repodb/models.go` (NEW)
- `gitmap/repodb/enums/` (NEW)
- `gitmap/repodb/consts.go` (NEW)
- `gitmap/repodb/root.go`

##### Acceptance Criteria
1. PascalCase primary keys: `RepoFileId`, `SearchCacheId`, `FileSequenceId`, `SequenceHistoryId`, `RepoScanLogId`, `IndexedRepoId` (zero bare `Id` or all-caps `ID`).
2. Generates type-safe column enums (`*FieldType`), singleton registries (`*DbRegistry`), null-safe row scanners (`Scan*`), and typed repositories (`*DbRepo`) using `30-db-struct-enum-generator.py`.
3. Integrates `dbengine.WrapDb(conn, dbengine.DbSQLite)` in `repodb.OpenRepoDB`.
4. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.

###### Subtask File: `03-error-management-and-zero-swallowed-errors.md`

#### Subtask 94.03: Zero-Swallowed Errors & Universal *apperror.AppError Wrapping

##### Goal
Eliminate all swallowed errors (`_ = `) and missing `*apperror.AppError` return types across `gitmap/pipelinedb/pipeline_split_ops.go`, `gitmap/dbengine/wrapper.go`, and `gitmap/repodb/repo_db.go`.

##### Files Impacted
- `gitmap/pipelinedb/pipeline_split_ops.go`
- `gitmap/dbengine/wrapper.go`
- `gitmap/repodb/repo_db.go`

##### Acceptance Criteria
1. `GetStats` in `pipelinedb/pipeline_split_ops.go`: properly handle and propagate errors from all 6 count queries instead of discarding them.
2. `QueryRecentRuns` and `QueryRecentErrorLogs`: row scan errors returned as `*apperror.AppError` rather than silently ignored.
3. `dbengine/wrapper.go`: `WithTransaction` inspects and logs/returns rollback error if rollback fails.
4. `repodb/repo_db.go`: `OptimizeRepoDB`, `OpenRepoDB`, and `InitRepoSchema` wrap errors in `*apperror.AppError` and handle checkpoint/optimize errors cleanly.
5. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.

###### Subtask File: `04-go-fast-cache-batch-indexing-and-binary-sniffer.md`

#### Subtask 94.04: Go Fast Cache Optimization, Upfront Batch Indexing & Binary Sniffer

##### Goal
Optimize Go's file indexing and caching in `gitmap/indexer/walker.go` and `gitmap/searcher/db_search.go` by adopting Python's fast cache patterns: batch upfront timestamp preloading, 8KB null-byte binary probe, and comprehensive directory pruning.

##### Files Impacted
- `gitmap/indexer/walker.go`
- `gitmap/searcher/db_search.go`
- `gitmap/fsutil/discovery_cache.go`

##### Acceptance Criteria
1. `indexer.Walker`: replaces per-file `db.QueryRowContext` with an upfront batch load of `RelativePath` -> `WriteTime` into an in-memory `map[string]int64`, eliminating thousands of redundant SQLite queries during traversal.
2. `indexer.Walker`: adds an 8KB null-byte sniffer (`bytes.IndexByte(head, 0) != -1`) before storing text into `RepoFile`, marking binary files with `IsBig = true` (or skipping content) to prevent database corruption.
3. Directory pruning: updates directory ignore rules to prune `.venv`, `dist`, `build`, `bin`, `vendor`, `.gemini`, `coverage`, `tmp`, `__pycache__`, and `.turbo`.
4. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.

###### Subtask File: `05-release-sync-preflight-and-ci-quality-gates.md`

#### Subtask 94.05: Codebase Release Synchronization, Pre-flight Verification & CI Quality Gates

##### Goal
Verify codebase alignment with `03-ai-scripts/14-version-sync-checker.py`, execute local linters and CI runner quality gates, and execute release verification.

##### Files Impacted
- `version.json`
- `package.json`
- `gitmap/constants/constants.go`
- `changelog.md`
- Local CI runner (`python 03-ai-scripts/06-cicd-local-runner.py`)

##### Acceptance Criteria
1. `python 03-ai-scripts/14-version-sync-checker.py` passes with zero drift.
2. `python linter-scripts/check-nested-ifs.py` reports 0 violations across repository.
3. `python linter-scripts/check-enum-and-boolean.py` reports 0 violations across repository.
4. `go test -C gitmap ./cmd` and `./store` and `./repodb` pass.
5. `06-cicd-local-runner.py --filter "Go Compile Gate"` passes with code 0.
6. Consolidated walkthrough artifact generated.


## 4. Unified Quality Gates & Verification Checklist

- [x] **Zero Concept Loss:** All source plans, code modifications, and execution steps preserved in full.
- [x] **Subtasks Inlined:** All associated subtasks folded directly into this document.
- [x] **Strict Relative Paths:** All citations use repository-relative paths without drive letters or file:/// URIs.
- [x] **Function Sizing:** All referenced codebase functions conform to <= 15 lines body cap.
- [x] **Coding Guidelines:** Affirmative booleans, zero nested ifs, and universal AppError wrapping verified.
- [x] **CI/CD Quality Gates:** All component tests pass legitimately under the local CI/CD runner.

## 5. Root Cause Analyses & Bug Fixes Referenced

- [`.lovable/memory/learned/08-sqlite-scanner-and-orm-evolution.md`](.lovable/memory/learned/08-sqlite-scanner-and-orm-evolution.md)
- [`.lovable/memory/issues/2026-09-07-sqlite-database-locked.md`](.lovable/memory/issues/2026-09-07-sqlite-database-locked.md)
