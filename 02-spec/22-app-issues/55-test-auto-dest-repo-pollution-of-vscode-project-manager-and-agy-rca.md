# RCA-55: Test Destination Auto-Provisioning Polluting VS Code Project Manager and Antigravity Projects

**Status:** Resolved
**Date:** 2026-10-01
**Affected Subsystem:** `cli/cmd/` (`create_ops.go`), `cli/workspacesync/` (`workspacesync.go`), `cli/vscodepm/` (`sync.go`)
**Associated Spec:** [02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md](../21-app/191-nodes-clone-except-self-windows-runner-and-path.md)
**Visual Proof:** VS Code Project Manager sidebar `Favorites` showing 10 duplicate `auto-dest-repo` entries pointing to `C:\Users\ADMINI~1\AppData\Local\Temp\TestProvisionDestTarget_MissingLocal<pid>\001\auto-dest-repo`

---

## 1. Reproduction Steps

1. Run standard unit tests for `cli/cmd`:
   ```powershell
   go test -v -run TestProvisionDestTarget_MissingLocal ./cmd
   go test -v -run TestEnsureOrProvisionDestinationRepo_MissingLocal ./cmd
   ```
2. Open VS Code with the `alefragnani.project-manager` extension installed.
3. Observe the Project Manager sidebar under `Favorites`:
   - Duplicate entries named `auto-dest-repo` and `new-project-dir` appear.
   - Hovering over `auto-dest-repo` reveals:
     `C:\Users\ADMINI~1\AppData\Local\Temp\TestProvisionDestTarget_MissingLocal...\001\auto-dest-repo`
   - Each successive test run appends another entry to `%APPDATA%\Code\User\globalStorage\alefragnani.project-manager\projects.json`.
   - Simultaneously, orphaned project config files are created in `~/.gemini/config/projects/<uuid>.json`.

---

## 2. Root Cause Analysis

### 2.1 Unchecked Workspace Sync in Auto-Provisioning
- In `cli/cmd/create_ops.go`, `provisionMissingDestination` handles auto-provisioning missing destination repositories during commit transfer and replay:
  ```go
  if err := initLocalRepo(params); err != nil {
      return "", err
  }
  tryPushRemoteProvisioned(params)
  workspacesync.SyncAll(absTarget, name)
  ```
- Both `TestProvisionDestTarget_MissingLocal` and `TestEnsureOrProvisionDestinationRepo_MissingLocal` invoke this function inside a temporary directory created by `t.TempDir()`.
- Because `provisionMissingDestination` had no check for temporary or test execution paths, it unconditionally invoked `workspacesync.SyncAll` for every mock test repository.

### 2.2 Inadequate Path Guarding in Workspace Sync
- In `cli/workspacesync/workspacesync.go`, `SyncAll`, `SyncWithoutDesktop`, and `SyncAgyOnly` did not check `isRestrictedPath` before calling `vscodepm.SyncMode`.
- Furthermore, `isRestrictedPath` verified only drive roots (`/`, `C:\`), user home dirs, and system folders (`WINDIR`, `ProgramFiles`), completely omitting OS temporary directories (`os.TempDir()`, `TEMP`, `TMP`, `AppData/Local/Temp`, `/tmp`).
- Consequently, test temporary paths passed straight through to `vscodepm.SyncMode` and `SyncAntigravity`.

### 2.3 Unfiltered Production Pair Ingestion in `vscodepm.SyncMode`
- `vscodepm.SyncMode` targets the host machine's live `projects.json` file.
- It passed all incoming `Pair` entries directly to `SyncAtMode` without verifying whether the `RootPath` was located in a temporary scratch or test directory.

---

## 3. Code Fix

1. **Dedicated Workspace Sync Path Guard:**
   - Created `cli/workspacesync/path_guard.go` with `IsTempOrTestPath(path string) bool` and `IsRestrictedPath(path string) bool`.
   - Detects `os.TempDir()`, environment temp directories (`TEMP`, `TMP`), standard temp locations (`AppData/Local/Temp`, `/tmp`, `/var/tmp`), and test paths (`TestProvisionDestTarget`, `TestEnsureOrProvision`).
   - `SyncAll`, `SyncWithoutDesktop`, and `SyncAgyOnly` now immediately abort if `IsRestrictedPath(repoPath)` is true.

2. **Creation Ops Isolation:**
   - In `cli/cmd/create_ops.go`, guarded `workspacesync.SyncAll(absTarget, name)` with `if !workspacesync.IsTempOrTestPath(absTarget)`.

3. **Production Pair Filter in `vscodepm`:**
   - Added `cli/vscodepm/path_filter.go` with `IsDisallowedProjectPath(path string) bool`.
   - Updated `vscodepm.SyncMode` to filter out any pairs located in temporary or disallowed directories before mutating `projects.json`.

4. **Fleet Clone Enhancements:**
   - In `cli/cmdclone/clonevscode.go`, added `if IsFleetCloneActive() { return }` in `openInVSCode` to prevent spawning GUI VS Code instances during fleet or batch operations.
   - In `cli/cmdclone/clonefixrepo_escape.go`, upgraded `escapeNestedGitRepo` to discover enclosing git roots upwards, correctly escaping out of git subdirectories to parent workspace roots.
   - In `cli/cmdnodes/nodes_clone_table.go`, aligned banner box width formatting (`%-94s`) and updated `sanitizeStdout` to ignore separator lines, stack traces, and pending task notices.

---

## 4. Prevention & Quality Guardrails

1. **Unit Test Verification:**
   - Added `cli/workspacesync/path_guard_test.go` verifying that temporary directories and test repository targets are recognized and rejected by workspace sync.
   - Verified that running `TestProvisionDestTarget_MissingLocal` leaves `projects.json` and Antigravity project configs completely untouched.
2. **Coding Guideline Conformance:**
   - Zero nested `if` statements.
   - All helper functions strictly $\le 15$ lines.
   - Positive boolean naming conventions maintained across all modified files.
