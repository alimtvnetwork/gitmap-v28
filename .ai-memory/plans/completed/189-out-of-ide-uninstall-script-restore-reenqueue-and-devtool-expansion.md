# Plan 189: Out-of-IDE Standalone Uninstaller, Snapshot Restore Re-Enqueue, and DevTool Cleaner Expansion (Completed)

Spec Reference: [02-spec/21-app/180-agy-agm-copilot-edge-uninstall-and-devtool-clean/01-overview.md](../../../02-spec/21-app/180-agy-agm-copilot-edge-uninstall-and-devtool-clean/01-overview.md)
Status: Completed
Completion Date: 2026-09-28

## Objective
Provide an out-of-IDE standalone PowerShell uninstallation script (`scripts/uninstall-antigravity.ps1`) to safely execute Antigravity full purge without terminating active IDE sessions mid-command; implement the snapshot restore and re-enqueue engine (`gitmap agy restore` / `gitmap agy reenq`); expand the DevTool cache cleaner to 12 categories including Rust/Cargo and IDE cache directories; verify all components hermetically via E2E test suites; and perform release ceremony.

---

## Subtask Execution Summary

### Subtask 01: Standalone Out-of-IDE PowerShell Antigravity & AGM Full Purge Script
- **Traceability ID:** Task-01
- **Target Files:** `scripts/uninstall-antigravity.ps1`
- **Delivered Capabilities:**
  - Standalone runner supporting execution outside the IDE terminal.
  - Generates restore snapshot JSON prior to any deletion.
  - Gracefully terminates running `antigravity.exe`, `agm.exe` processes.
  - Enforces hard invariant: strictly ignores/skips `D:\work` and active git directories.
  - Supports `-NewWindow` to spawn an independent PowerShell session.
- **Verification:** Dry-run executed successfully with zero data deletion.

### Subtask 02: AGY Restore & Re-Enqueue Engine from Snapshot
- **Traceability ID:** Task-02
- **Target Files:** `cli/cmdagy/agy_restore.go`, `cli/cmdagy/agy_cmd.go`
- **Delivered Capabilities:**
  - `RunAGYRestore(snapPath string) (int, int, error)` deserializes snapshot JSON.
  - Re-registers projects into Antigravity project config files (`~/.gemini/config/projects/<id>.json`).
  - Re-registers conversation records into SQLite `conversation_summaries` database.
  - Wired into Cobra CLI: `gitmap agy restore [path]` and `gitmap agy reenq [path]`.
- **Verification:** Unit tests and boolean linters passed with 0 violations.

### Subtask 03: DevTool Cache Cleaner Category Expansion
- **Traceability ID:** Task-03
- **Target Files:** `cli/osclean/cleaner_enhanced_categories.go`, `cli/osclean/cleaner_enhanced.go`
- **Delivered Capabilities:**
  - Expanded from 10 to 12 categories:
    - Added Rust / Cargo cache (`~/.cargo/registry/cache`, `~/.cargo/git/db`).
    - Added IDE scratch & updater logs (`~/.antigravity/logs`, `%LOCALAPPDATA%\antigravity-updater`).
  - Multi-category parallel scanner calculates space and formats aligned summary tables.
- **Verification:** Dry-run verified 12 categories scanned concurrently.

### Subtask 04: Hermetic E2E Verification & Release Ceremony
- **Traceability ID:** Task-04
- **Target Files:** `cli/tests/uninstall_e2e_test.go`, `scripts/test-uninstall-e2e.ps1`
- **Delivered Capabilities:**
  - Extended hermetic E2E tests (`TestE2EAGYRestoreFromSnapshot`, `TestE2EDevToolCleanerDryRun` with 12 categories).
  - All linters passed with 0 violations.
  - Pre-commit git pull and atomic feature commit executed.
- **Verification:** `go test ./tests -tags=e2e -v` passed all 6 tests.
