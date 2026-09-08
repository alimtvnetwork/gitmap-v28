# Subtask 05: Install Logs CLI and Unit Tests

**Parent Plan:** [77-scripts-fixer-installation-split-db-and-tooling-engine.md](../../pending/77-scripts-fixer-installation-split-db-and-tooling-engine.md)
**Status:** Pending
**Target Files:**
- `gitmap/cmd/install.go`
- `gitmap/cmd/install_logs.go`
- `gitmap/cmd/install_unit_test.go`

## Objectives
1. Implement `runInstallLogs(args []string)` in `gitmap/cmd/install_logs.go`:
   - Prints table of recent installation execution logs from `installation.db` (`Tool`, `Action`, `Version`, `Manager`, `Duration`, `Status`, `CreatedAt`).
   - Supports `--failed` flag to filter only failed install attempts.
   - Supports `--tool <name>` to filter by tool.
2. Route `gitmap install logs`, `gitmap in logs`, and `gitmap in --logs` in `gitmap/cmd/install.go`.
3. Add CLI unit tests in `gitmap/cmd/install_unit_test.go`.
