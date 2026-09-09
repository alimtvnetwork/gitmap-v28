# 71-os-power-management-and-installation-split-db.md: OS Screen Timeout & Sleep Management Framework & Installation Split DB

## 1. Executive Summary

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

## 2. Mandatory Rules & Invariants

1. **Rule 1 (Strict Relative Git Paths)**: All paths in markdown files, artifacts, and documentation must be strictly relative to the repository root (`/`). Zero `file:///` URIs.
2. **Rule 2 (Strict File & Function Sizing)**: Every function must be $\le 15$ lines (preferred $\le 8$ lines), with a mandatory blank line before every return statement. All new and modified Go files must remain $\le 200$ lines.
3. **Rule 3 (Database Conventions Compliance)**: All tables singular PascalCase (`InstalledTool`, `InstallationLog`, `PowerSetting`, `PowerSettingHistory`), integer auto-increment PKs (`{TableName}Id`), affirmative booleans (`IsNeverSleep`, `IsSuccess`), and Rule 10/11 context columns.
4. **Rule 4 (SQLite Concurrency Standard)**: All split databases must call `conn.SetMaxOpenConns(1)` and anchor paths via `store.BinaryDataDir()` / `filepath.EvalSymlinks(os.Executable())`.
5. **Rule 5 (AST Parity Guard)**: Top-level CLI command constants (`CmdPower = "power"`, `CmdPowerAlias = "pw"`, `CmdPowerAlias2 = "pwr"`) must be registered in `constants_cli.go` under `// gitmap:cmd top-level` and synchronized in `cmd_constants_test.go:topLevelCmds()`.
6. **Rule 6 (Zero Swallowed Errors)**: No blank error discards; wrap all errors using `apperror.WrapSimple` or `apperror.NewWithDetails`.

## 3. Subsystems & Architecture

### A. Pluggable Power Management Framework (`gitmap/power/`)
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

### B. SQLite Power Settings Persistence (`gitmap/store/`)
- `PowerSetting` table: records current and previous profiles (`baseline`, `current`, `previous`).
- `PowerSettingHistory` table: transactional audit log of every power change.
- Store helper: `gitmap/store/power.go`.

### C. Installation Split DB (`gitmap/store/installation_split_*.go`)
- Location: `<BinaryDataDir>/installation.db`
- Connection wrapper: `InstallationSplitDB` with `conn.SetMaxOpenConns(1)`.
- Tables:
  - `InstalledTool`: tracked packages and semantic version components.
  - `InstallationLog`: detailed execution logs, durations, exit codes, and stdout/stderr.
- Migration Bridge: `MigrateInstalledToolsFromRoot(rootDB, splitDB)`.

### D. CLI Command Integration (`gitmap/cmd/`)
- Power: `gitmap power [status|never-sleep|set|reset|history]` in `gitmap/cmd/power*.go`.
- Installation: wires `install.go`, `installtools.go`, `install_buildessential.go`, `installlist.go`, and `uninstall.go` to use `store.OpenInstallationSplitDB()`.

## 4. Subtasks Breakdown

1. [01-installation-split-db-schema-and-manager.md](subtasks/71-os-power-management-and-installation-split-db/01-installation-split-db-schema-and-manager.md): Implement `InstallationSplitDB` lifecycle, `InstalledTool`, `InstallationLog` schema, and CRUD in `gitmap/store/installation_split_*.go`.
2. [02-installation-split-db-migration-and-command-wiring.md](subtasks/71-os-power-management-and-installation-split-db/02-installation-split-db-migration-and-command-wiring.md): Implement migration bridge from root DB, backward-compatible delegation on `*store.DB`, and wire `install.go`, `install_buildessential.go`, `installtools.go`, `installlist.go`, and `uninstall.go`.
3. [03-os-power-framework-and-drivers.md](subtasks/71-os-power-management-and-installation-split-db/03-os-power-framework-and-drivers.md): Implement pluggable `power.Manager` interface, `types.go`, Windows `powercfg` driver, Linux 3-tier driver (`gsettings`/`xset`), and Darwin/fallback drivers.
4. [04-os-power-sqlite-persistence-and-commands.md](subtasks/71-os-power-management-and-installation-split-db/04-os-power-sqlite-persistence-and-commands.md): Implement `PowerSetting` and `PowerSettingHistory` store operations, CLI dispatcher `gitmap power` (`status`, `never-sleep`, `set`, `reset`, `history`), and AST constants.
5. [05-helptext-docs-and-ci-verification.md](subtasks/71-os-power-management-and-installation-split-db/05-helptext-docs-and-ci-verification.md): Author helptext files (`power.md`), update changelogs, verify AST parity, golden tests, and execute full CI quality runner.

## 5. Acceptance Criteria

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
