# Subtask 02: Installation Split DB Telemetry & Command Runner

**Parent Plan:** [77-scripts-fixer-installation-split-db-and-tooling-engine.md](../../pending/77-scripts-fixer-installation-split-db-and-tooling-engine.md)
**Status:** Completed
**Target Files:**
- `gitmap/store/installation_split_log.go`
- `gitmap/store/installation_split_db_test.go`
- `gitmap/cmd/installtools.go`
- `gitmap/cmd/install_packages.go`
- `gitmap/cmd/install_audit_runner.go`
- `gitmap/cmd/install_linux_apps.go`
- `gitmap/cmd/install_log_file.go`
- `gitmap/cmd/install_audit_runner_test.go`

## Objectives Completed
1. Added `RecordExecution` and `GetFailedLogs` methods to `InstallationSplitDB` in `gitmap/store/installation_split_log.go`.
2. Added output bounding `SanitizeLogOutput` (64 KB cap per stream) to prevent SQLite database bloat.
3. Implemented `executeCommandWithAudit(args []string, verbose bool) commandExecutionResult` with dual-stream capturing in `gitmap/cmd/install_audit_runner.go`.
4. Updated `runInstallCommand` and `handleInstallError` to record execution telemetry (`DurationMs`, `ExitCode`, `Stdout`, `Stderr`, `CommandLine`, `IsSuccess`) for both passing and failing runs into `installation.db`.
5. Fixed missing install telemetry in `runInstallVSCodeLinux` and `runInstallGitHubDesktopLinux` with step-by-step phase execution and audit logging in `gitmap/cmd/install_linux_apps.go`.
6. Refactored `installtools.go` into modular sub-200-line files complying strictly with coding guidelines (zero nested ifs, functions <= 15 lines, blank lines before returns, no variable reassignment, zero linter violations).
