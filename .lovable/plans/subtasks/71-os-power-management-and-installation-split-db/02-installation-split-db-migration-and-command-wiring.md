# Subtask 02: Installation Split DB Migration Bridge & Command Wiring

## Status
Completed

## Context & Objectives
1. **Migration Bridge**:
   - `gitmap/store/installation_split_bridge.go`: Reads any legacy rows from root DB `InstalledTool` and migrates them into `installation.db` (`INSERT OR IGNORE`).
   - Maintains delegate wrappers in `gitmap/store/installedtool.go` on `*store.DB` so any remaining references seamlessly route to `InstallationSplitDB`.
2. **Command Wiring**:
   - `gitmap/cmd/installtools.go`: Writes to `store.OpenInstallationSplitDB()` and logs telemetry (`InstallationLogRecord`).
   - `gitmap/cmd/install_buildessential.go`: Writes `build-essential` profile and constituent tools to `installation.db` with execution timing.
   - `gitmap/cmd/installlist.go`: Reads installed tool status from `installation.db`.
   - `gitmap/cmd/uninstall.go`: Uninstalls and records audit log in `installation.db`.

## Verification Steps
- Unit tests verify root DB migration preserves version strings and timestamps.
- CLI installation commands write to `installation.db`.
