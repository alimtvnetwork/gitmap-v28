# Issue 21 / RCA-55: Test Destination Auto-Provisioning Polluting VS Code Project Manager and Antigravity Projects

**Status:** Resolved
**Date:** 2026-10-01
**Affected Subsystem:** `cli/cmd/` (`create_ops.go`), `cli/workspacesync/` (`workspacesync.go`, `path_guard.go`), `cli/vscodepm/` (`sync.go`, `path_filter.go`)
**Associated Spec:** [02-spec/22-app-issues/55-test-auto-dest-repo-pollution-of-vscode-project-manager-and-agy-rca.md](../../02-spec/22-app-issues/55-test-auto-dest-repo-pollution-of-vscode-project-manager-and-agy-rca.md)

---

## 1. Symptom

During automated unit test runs (`go test ./cmd/...`), test temporary directories (`auto-dest-repo`, `new-project-dir`) created under `%LOCALAPPDATA%\Temp\` were unconditionally synchronized into VS Code's `projects.json` file and Antigravity IDE configuration databases, causing persistent project clutter and mock repository pollution on developer workstations.

---

## 2. Root Cause

1. `provisionMissingDestination` in `cli/cmd/create_ops.go` called `workspacesync.SyncAll` without checking if the repository path was a temporary test directory.
2. `workspacesync.SyncAll` lacked guards for OS temporary directories (`os.TempDir()`, `$TEMP`, `$TMP`, `AppData\Local\Temp`).
3. `vscodepm.SyncMode` lacked pair filtering against temporary project roots.

---

## 3. Resolution

1. Added `cli/workspacesync/path_guard.go` with `IsTempOrTestPath` and `IsRestrictedPath` to abort workspace synchronization for any path located inside system temp directories.
2. Added `cli/vscodepm/path_filter.go` with `IsDisallowedProjectPath` and test bypass flags to prevent test mock fixtures from entering production project files.
3. Added unit tests in `cli/workspacesync/path_guard_test.go` confirming zero sync execution for temporary repository targets.
