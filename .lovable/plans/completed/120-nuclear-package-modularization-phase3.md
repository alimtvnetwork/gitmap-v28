# 120-nuclear-package-modularization-phase3.md: Nuclear Monolith Subpackage Modularization, Heavy Test Isolation & Test Inventory Duration Sync

**Status: completed**

## 1. Problem Diagnosis & Architectural Goal
1. `gitmap/cmd` was an oversized monolithic Go package containing over 1,190 Go files, creating heavy compilation bloat and slow test execution.
2. In-tree tests in `gitmap/cmd` previously executed heavy `exec.Command` subprocesses, git command executions, and `time.Sleep` delays during standard test loops.
3. Extracted cohesive domains into strictly acyclic subpackages (`gitmap/cmdagy`, `gitmap/cmdchromeprofile`) and leaf packages (`gitmap/ecosystemgroup`, `gitmap/fsutil`, `gitmap/model/scan_records.go`).
4. Isolated all heavy/slow subprocess tests into `gitmap/tests/heavy_test/` (`package heavy_test`), reducing routine package tests in `gitmap/cmd` to pure in-memory execution (<0.01s).

## 2. Strict DAG Architecture
- Layer 0 (Leaves): `constants`, `model`, `store`, `apperror`, `cliexit`, `fsutil`, `ecosystemgroup`
- Layer 1 (Domain Subpackages): `gitmap/cmdprompt`, `gitmap/cmdpurge`, `gitmap/cmdvmware`, `gitmap/cmdagy`, `gitmap/cmdchromeprofile`
- Layer 2 (Root CLI Orchestrator): `gitmap/cmd`
- Layer 3 (Entrypoint): `main.go`
- Test Layer: Isolated integration test packages (`gitmap/tests/heavy_test`)

Zero circular dependencies exist across the entire repository. `go vet ./...` and `go build ./...` pass 100% cleanly.

## 3. Deliverables & Execution Summary
1. **Isolated 10 Heavy Test Suites into `gitmap/tests/heavy_test/`**:
   - `clonefixrepo_e2e_test.go`
   - `codingguidelines_exec_test.go`
   - `fixgit_rebuild_index_test.go`
   - `fixgit_unmerged_untracked_test.go`
   - `reconcile_workflow_e2e_test.go`
   - `ssh_login_spawn_test.go`
   - `visibility_local_remote_exec_test.go`
   - `historyrewrite_python_test.go`
   - `hygiene_scan_e2e_test.go`
   - `gomod_workflow_e2e_test.go`
   - Replaced all `time.Sleep` in `cmd` with `runtime.Gosched()` and channel synchronizations.
   - **Result**: Exactly **ZERO** `exec.Command` and **ZERO** `time.Sleep` remain in `gitmap/cmd/*_test.go`. Routine package tests run in <0.01s.

2. **Created Leaf Packages & Shared Helpers**:
   - `gitmap/model/scan_records.go`: `LoadStatusRecords(path string) ([]ScanRecord, error)`
   - `gitmap/fsutil/copydir.go`: `CopyDirContents(src, dest string) *apperror.AppError`
   - `gitmap/fsutil/snapshot.go`: `PickLatestSnapshotWithTag(snapshots []string, tag string) string`
   - `gitmap/ecosystemgroup/`: `groups.go` and `ops.go` (`package ecosystemgroup`)

3. **Created `gitmap/cmdagy` (34 files)**:
   - Moved all `agy_*.go` and `find_duplicates_agy.go` files into `package cmdagy`.
   - Exported `AgyCmd`, `DispatchAgy`, `RunFindDuplicates`, `RunOptimize`.
   - Wired seamlessly into `gitmap/cmd/root.go` and `gitmap/cmd/find_duplicates.go`.

4. **Created `gitmap/cmdchromeprofile` (44 files)**:
   - Moved all `chromeprofile_*.go`, `chrome_batch.go`, `find_duplicates_chrome.go` files into `package cmdchromeprofile`.
   - Created `helpers.go` and `testhelpers_test.go`.
   - Exported public APIs: `RunProfileCopy`, `RunProfileExport`, `RunProfileImport`, `RunProfileImportCheck`, `RunProfileList`, `RunProfileDelete`, `RunProfileClear`, `RunProfileOptimize`, `RunProfileMerge`, `RunProfileReconcile`, `RunProfileUndo`, `RunProfileRedo`, `RunProfileGroupDispatch`, `RunCopyAll`, `RunExportAll`, `RunImportAll`, `RunFindDuplicates`, `UserDataDir`, `ResolveProfileDir`, `ResolveProfile`, `ProfileSummary`, `ProfileDisplayName`, `AvailableProfileNames`, `PrintAvailableProfilesWithDisplay`, `SnapshotProfile`, `CopyEntry`, `IsChromeRunning`.
   - Wired seamlessly into `gitmap/cmd/chrome.go`, `export.go`, `profile.go`, `roottooling.go`, etc.

5. **Package Size Reduction**:
   - `gitmap/cmd` was reduced from 1,191 files down to 1,110 files (-81 files).
   - Package builds and compiles substantially faster with cleaner isolation and lower cognitive overhead.

6. **Synchronized `.lovable/test-inventory.json`**:
   - Total Tests Indexed: 3,532 tests across 125 packages.
   - Classified tests into `unit` (<0.01s) vs `heavy` (1.0s–5.0s) tiers.
   - `gitmap/tests/heavy_test`: 36 heavy tests properly isolated and categorized.
   - Safely recorded modified files under lock in `.lovable/temp/recent-file-changes.json`.

7. **Quality Gates & Repo Standards**:
   - Formatted all 2,558 Go files across 16 cores via `26-go-code-formatter.py`.
   - Verified 100% clean compilation and zero import cycles via `go vet ./...` and `go build ./...`.
