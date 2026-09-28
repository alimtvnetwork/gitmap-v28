# Plan 188: AGY, AGM, Copilot, and Edge Complete Uninstallation Suite & Enhanced DevTool Cache Cleaner (Completed)

Spec Reference: [02-spec/21-app/180-agy-agm-copilot-edge-uninstall-and-devtool-clean/01-overview.md](../../../02-spec/21-app/180-agy-agm-copilot-edge-uninstall-and-devtool-clean/01-overview.md)
Status: Completed
Completion Date: 2026-09-28

## Objective
Implement a robust, safe uninstallation engine for Antigravity (`gitmap uninstall agy`, `gitmap uninstall agy-all`), Antigravity Manager (`gitmap uninstall agm`, `gitmap uninstall agm-all`), Windows Copilot (`gitmap uninstall copilot`), and Microsoft Edge via Chris Titus WinUtil parity (`gitmap uninstall edge`). Ensure full state snapshotting (project workspaces, conversation IDs/names) prior to AGY full purge, strict protection of `d:\work` and active work directories, enhance the multi-category DevTool cache cleaner (`gitmap devtool clear`, `gitmap dt clear`), and provide hermetic local E2E test suites with an out-of-IDE PowerShell runner.

---

## Subtask Execution Summary

### Subtask 01: AGY State Snapshot and Dual-Tier Uninstaller
- **Traceability ID:** Task-01 (AC-SPEC180-01, AC-SPEC180-02, AC-SPEC180-06)
- **Target Files:** `cli/cmdagy/agy_snapshot.go`, `cli/cmdagy/agy_uninstall.go`, `cli/cmdagy/agy_types.go`, `cli/cmdagy/agy_cmd.go`, `cli/cmd/uninstall.go`
- **Delivered Capabilities:**
  - `ExportAGYRestoreSnapshot(outPath string)` scans Antigravity workspaces and conversation summaries from SQLite store and brain transcript directory, exporting metadata to `~/.gitmap/agy-snapshot-<timestamp>.json` or custom `--backup` location.
  - Dual-tier uninstaller `RunAGYUninstall(isFullPurge, isForce, isDryRun bool, backupPath string)`:
    - Standard tier (`gitmap uninstall agy`): Removes application binaries and launcher scripts while preserving user data and conversation history.
    - Full purge tier (`gitmap uninstall agy-all`, `gitmap agy uninstall-all`, `gitmap agy uninstall --all`): Mandates snapshot generation first, then purges `.gemini`, brain transcripts, and cache trees.
  - Strict work directory protection invariant: `IsPathSafeToDelete` and `IsWorkDirectoryOverlap` reject any deletion candidate matching or containing `d:\work`, root filesystem volumes, or active git repositories.
- **Verification:** Unit tests and boolean/nested-if linters passed with 0 violations.

### Subtask 02: Antigravity Manager (AGM) Uninstaller
- **Traceability ID:** Task-02 (AC-SPEC180-03)
- **Target Files:** `cli/cmdinstall/agm_uninstall.go`, `cli/cmdinstall/agm_install.go`, `cli/cmd/uninstall.go`
- **Delivered Capabilities:**
  - `RunAGMUninstall(isFullPurge, isForce, isDryRun bool)` discovers running AGM processes (`antigravity-manager`, `agm.exe`, `agm`) and terminates them cleanly.
  - Purges binary wrappers and desktop shortcuts.
  - Purges configuration folders (`~/.agm`, `%APPDATA%\agm`, `%LOCALAPPDATA%\agm`) when full purge is flagged.
  - Integrated into root dispatch (`gitmap uninstall agm`, `gitmap uninstall agm-all`, `gitmap agm uninstall`).
- **Verification:** Linters verified, function lengths <= 15 lines.

### Subtask 03: Windows Copilot Removal & Registry Blocker
- **Traceability ID:** Task-03 (AC-SPEC180-04)
- **Target Files:** `cli/cmdwinutil/winutil_copilot.go`, `cli/cmdwinutil/winutil_types.go`, `cli/cmdwinutil/winutil_cmd.go`, `cli/cmd/uninstall.go`
- **Delivered Capabilities:**
  - `RunCopilotUninstall(isDryRun bool)` executes PowerShell Appx removal for `*Microsoft.Windows.Copilot*`.
  - Sets Group Policy `TurnOffWindowsCopilot = 1` in both HKCU and HKLM.
  - Disables taskbar Copilot button (`ShowCopilotButton = 0`).
  - Supports `--dry-run` simulation mode and provides safe fallback for non-Windows platforms.
- **Verification:** Dry-run unit tests pass cleanly without system alteration.

### Subtask 04: Microsoft Edge Uninstallation (Chris Titus WinUtil Parity)
- **Traceability ID:** Task-04 (AC-SPEC180-05)
- **Target Files:** `cli/cmdwinutil/winutil_edge.go`, `cli/cmdwinutil/winutil_cmd.go`, `cli/cmd/uninstall.go`
- **Delivered Capabilities:**
  - `RunEdgeUninstall(hasKeepWebView2, isDryRun bool)` implements full Chris Titus WinUtil removal flow:
    - Force terminates `msedge.exe` (preserving `msedgewebview2.exe` by default).
    - Detects Edge versioned `setup.exe` installer across Program Files directories.
    - Executes silent force uninstall (`setup.exe --uninstall --system-level --verbose-logging --force-uninstall`).
    - Removes Edge Appx packages (`Microsoft.MicrosoftEdge`, `MicrosoftEdgeDevToolsClient`).
    - Disables background update services (`edgeupdate`, `edgeupdatem`).
    - Sets registry blocker `HKLM\SOFTWARE\Microsoft\EdgeUpdate` -> `DoNotUpdateToEdgeWithChromium = 1`.
  - Integrated into CLI via `gitmap uninstall edge`, `gitmap winutil edge uninstall`.
- **Verification:** Unit tests confirm WebView2 preservation and dry-run safety.

### Subtask 05: Enhanced DevTool Cache Cleaner
- **Traceability ID:** Task-05 (AC-SPEC180-07)
- **Target Files:** `cli/osclean/cleaner_enhanced.go`, `cli/osclean/cleaner_enhanced_categories.go`, `cli/osclean/cleaner_types.go`, `cli/cmdos/devclean_cmd.go`, `cli/cmd/clean_dev_entry.go`
- **Delivered Capabilities:**
  - Concurrent 10-category cache cleaner covering:
    1. Go build cache & module cache
    2. Node/npm/pnpm/yarn/bun caches
    3. Python/pip/uv/poetry caches & `__pycache__`
    4. Webpack/Vite/Turbopack caches
    5. Antigravity temporary logs & transcripts
    6. VS Code workspace storage caches
    7. Chrome/Chromium dev profile caches
    8. Git dangling objects & index caches
    9. System temp developer artifacts
    10. Stale test binaries (`*.test`, `*.test.exe`)
  - Parallel size scanning worker pool, dry-run previews, human-readable size formatting (`FormatCleanSize`), and aligned terminal summary tables.
  - CLI routing via `gitmap devtool clear`, `gitmap dt clear`, and `gitmap clean-dev`.
- **Verification:** Dry-run test verifies all 10 categories scanned concurrently.

### Subtask 06: Hermetic Local E2E Verification Suites & Out-of-IDE Runner
- **Traceability ID:** Task-06 (AC-SPEC180-08, AC-SPEC180-09)
- **Target Files:** `cli/tests/uninstall_e2e_test.go`, `scripts/test-uninstall-e2e.ps1`
- **Delivered Capabilities:**
  - `cli/tests/uninstall_e2e_test.go` guarded by `//go:build e2e` verifying:
    - AGY snapshot generation creates valid JSON with project counts and timestamps.
    - Safety invariant verifies `d:\work` and root directories are strictly protected.
    - Dry-run uninstallations execute cleanly without errors.
    - Copilot, Edge, and 10-category DevTool cleaner dry-runs produce correct status results.
  - Out-of-IDE PowerShell runner script `scripts/test-uninstall-e2e.ps1` executes non-destructive validations and CLI help checks.
- **Verification:** Hermetic tests run without touching real system data; linters confirm 0 violations.

---

## Acceptance Criteria Traceability Matrix

| AC ID | Description | Verified Status |
|---|---|---|
| AC-SPEC180-01 | AGY state snapshot serializes workspace paths and conversation IDs to JSON | Verified |
| AC-SPEC180-02 | AGY dual-tier uninstaller supports standard bin removal & full purge | Verified |
| AC-SPEC180-03 | AGM process termination, binary wrapper removal, and config purging | Verified |
| AC-SPEC180-04 | Windows Copilot Appx removal, registry Group Policy blocker, taskbar toggle | Verified |
| AC-SPEC180-05 | Microsoft Edge Chris Titus WinUtil force uninstall, Appx removal, registry blocker | Verified |
| AC-SPEC180-06 | Work directory protection invariant strictly prevents deleting `d:\work` | Verified |
| AC-SPEC180-07 | 10-category DevTool cache cleaner with parallel scan and formatted summary table | Verified |
| AC-SPEC180-08 | Hermetic E2E test suite isolated behind `//go:build e2e` | Verified |
| AC-SPEC180-09 | PowerShell out-of-IDE runner `scripts/test-uninstall-e2e.ps1` | Verified |
