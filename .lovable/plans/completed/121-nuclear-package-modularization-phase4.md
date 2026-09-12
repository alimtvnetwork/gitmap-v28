# 121-nuclear-package-modularization-phase4.md: Nuclear Monolith Subpackage Modularization (cmdmacro, cmdvscode, cmdvhost, cmdzip), Heavy Test Isolation & Test Inventory Duration Estimation

**Status: completed**

> **Task Initiation & Execution:**
> Started in response to the user's "nuclear change" directive to decompose monolithic Go packages into small acyclic packages, isolate heavy/slow tests into `gitmap/tests/heavy_test/`, re-estimate test durations in `.lovable/test-inventory.json`, and ensure zero circular dependencies.
> Completed autonomously across 7 execution steps/subtasks.

## 1. Problem Diagnosis & Architectural Goal
1. `gitmap/cmd` was an oversized monolithic Go package containing 1,110 Go files, creating heavy compilation bloat and slow test iterations as Go compiles all files in a package as a single unit.
2. In-tree tests in several packages (`gitmap/clonefrom`, `gitmap/committransfer`, `gitmap/release`) executed real external `git` subprocesses, `exec.Command`, and long 220-commit history loops (taking up to 26.8s per test), slowing down regular development runs.
3. Extracted 4 cohesive domains into strictly acyclic domain subpackages (`gitmap/cmdmacro`, `gitmap/cmdvscode`, `gitmap/cmdvhost`, `gitmap/cmdzip`), reducing `gitmap/cmd` from 1,110 files down to 1,042 files (-68 files).
4. Isolated all slow git subprocess and commit replay tests into `gitmap/tests/heavy_test/` (`package heavy_test`), reducing routine package tests in `gitmap/clonefrom`, `gitmap/committransfer`, and `gitmap/release` to pure in-memory execution (<0.5s total).
5. Synchronized `.lovable/test-inventory.json` with updated package listings, test counts, and duration tier classifications (`unit` vs `heavy`).

## 2. Strict DAG Architecture (Zero Cycles)
- **Layer 0 (Leaves)**: `constants`, `model`, `store`, `apperror`, `cliexit`, `fsutil`, `ecosystemgroup`, `vscodepm`, `macro`, `archive`
- **Layer 1 (Domain Subpackages)**: `cmdprompt`, `cmdpurge`, `cmdvmware`, `cmdagy`, `cmdchromeprofile`, `cmdmacro`, `cmdvscode`, `cmdvhost`, `cmdzip`
- **Layer 2 (Root CLI Orchestrator)**: `gitmap/cmd` (dispatches commands downwards to Layer 1 domain packages; Layer 1 packages NEVER import `cmd`)
- **Layer 3 (Entrypoint)**: `main.go`
- **Isolated Test Layer**: `gitmap/tests/heavy_test` (`package heavy_test` blackbox tests)

Zero circular dependencies exist across the entire repository. `go vet ./...` and `go build ./...` pass with exit code 0.

## 3. Subtask Consolidation & Deliverables

### Subtask 1: Heavy Test Isolation & Duration Classification (Completed)
- Isolated subprocess and long-running tests into `gitmap/tests/heavy_test/`:
  - `clonefrom_execute_e2e_test.go`: git bare repo cloning (0.55s).
  - `committransfer_e2e_test.go`: git commit transfer, count parity, and dirty source rejection (7.73s).
  - `committransfer_idempotence_e2e_test.go`: 220-commit history bury loop (26.79s).
  - `release_scan_executor_e2e_test.go`: commit, branch, and tag generation (0.53s).
- Retained fast assembly unit tests in `gitmap/clonefrom/`, `gitmap/committransfer/`, and `gitmap/release/`, all executing in <0.5s total.
- Verified all heavy tests pass in `gitmap/tests/heavy_test/` (exit 0).

### Subtask 2: Extract `gitmap/cmdmacro` Subpackage (28 files) (Completed)
- Extracted 25 macro files from `gitmap/cmd/` to new package `gitmap/cmdmacro/`:
  - `macro_add*.go`, `macro_cmd.go`, `macro_edit*.go`, `macro_export*.go`, `macro_import*.go`, `macro_record*.go`, `macro_retry*.go`, `macro_tree*.go`, `macro_tui.go`, `macro_types.go`, `macro_ext_cmd.go`.
- Created `helpers.go` and `exports.go`.
- Exported entrypoints: `MacroCmd`, `RunMacroCmd`, `RunExecuteCmd`, `HandleMacroAdd`, `HandleMacroEdit`, `HandleMacroList`, `HandleMacroRecord`, `HandleMacroShow`, `HandleMacroDelete`, `RunMacroExport`, `RunMacroImport`, `RunMacroUntilSuccess`, `ExecuteDynamicMacro`, `RunCatCmd`, `RunTouchCmd`, `RunMkfileCmd`.
- Wired bridge functions in `gitmap/cmd/clihelpers.go` and updated `gitmap/cmd/macro_root_dispatch.go` and `gitmap/cmd/rootdata.go`.
- Verified: `go test ./cmdmacro/...` passes 100%.

### Subtask 3: Extract `gitmap/cmdvscode` Subpackage (28 files) (Completed)
- Extracted 26 VS Code workspace, PM sync, clone sync, and duplicate detection files from `gitmap/cmd/` to new package `gitmap/cmdvscode/`:
  - `clonepmsync*.go`, `find_duplicates_vscode.go`, `vscode*.go`, `vscodepm*.go`.
- Created `helpers.go`, `exports.go`, and `testhelpers_test.go`.
- Exported entrypoints: `VSCodeCmd`, `RunVSCode`, `RunVSCodePMSync`, `RunVSCodePMPath`, `RunVSCodeWorkspace`, `RunFindDuplicates`, `IsVSCodeSyncDisabled`, `ReportVSCodePMSoftError`, `SyncClonedReposToVSCodePM`, `SyncSingleClonedRepoToVSCodePM`, `BuildClonePMPair`, `CanonicalizePMPath`, `SyncRecordsToVSCodePM`, `RenameVSCodePMByPath`, `PrintVSCodeOptimizeResult`, `RunGitHubDesktopGroup`, `StripVSCodeSyncDisabledFlag`, `StripVSCodeTagFlags`, `ApplyDebugPathsEnv`.
- Wired bridge functions in `gitmap/cmd/clihelpers.go`.
- Verified: `go test ./cmdvscode/...` passes 100%.

### Subtask 4: Extract `gitmap/cmdvhost` Subpackage (13 files) (Completed)
- Extracted 12 virtual host configuration, templates, ops, and site management files (`vhost*.go` and `nginx_rm.go` -> `vhost_rm.go`) into new package `gitmap/cmdvhost/`.
- Created `helpers.go` and `exports.go`.
- Exported `VHostConfig` type alias, constants, `RenderVHostConfig`, and vhost command runners.
- Wired bridge functions in `gitmap/cmd/clihelpers.go`.
- Verified: `go test ./cmdvhost/...` passes 100%.

### Subtask 5: Extract `gitmap/cmdzip` Subpackage (8 files) (Completed)
- Extracted 6 archive compression, decompression, and zip-group management files (`zip.go`, `unzipcompact.go`, `zipgroup*.go`) into new package `gitmap/cmdzip/`.
- Created `helpers.go`, `exports.go`, and unit tests `zip_unit_test.go`.
- Exported `RunZip`, `RunUnzipCompact`, and `RunZipGroup`.
- Wired bridge functions in `gitmap/cmd/clihelpers.go`.
- Verified: `go test ./cmdzip/...` passes 100%.

### Subtask 6: Test Inventory Synchronization & Duration Re-Estimation (Completed)
- Regenerated `.lovable/test-inventory.json` using `03-ai-scripts/33-test-inventory-generator.py`.
- Result: **3,534 tests indexed across 129 packages** (up from 125 packages; 4 new domain packages: `cmdmacro`, `cmdvscode`, `cmdvhost`, `cmdzip`).
- Tests in `gitmap/tests/heavy_test/` categorized as `heavy` tier (~3.0s).
- In-tree package tests categorized as `unit` tier (~0.005s).
- Safely recorded all modified files under lock in `.lovable/temp/recent-file-changes.json`.

### Subtask 7: Repository Verification & Quality Gates (Completed)
- Strict acyclic compilation confirmed: `go vet ./...` and `go build ./...` exit 0 cleanly.
- Code formatting applied: `python 03-ai-scripts/26-go-code-formatter.py` formatted 2,569 Go files in 0.89s.
- Boolean, enum, and conditional compliance verified: `python linter-scripts/check-nested-ifs.py` (0 violations across 2,801 files) and `python linter-scripts/check-enum-and-boolean.py` (0 violations across 2,080 files).
- Multi-package test pass verified: `go test ./cmdmacro/... ./cmdvscode/... ./cmdvhost/... ./cmdzip/... ./tests/heavy_test/...` (100% pass).

## 4. Package Metrics & Impact Summary
- `gitmap/cmd/` file count reduced from **1,110 files** to **1,042 files** (-68 files).
- Created 4 dedicated, highly cohesive subpackages:
  1. `gitmap/cmdmacro/` (28 files)
  2. `gitmap/cmdvscode/` (28 files)
  3. `gitmap/cmdvhost/` (13 files)
  4. `gitmap/cmdzip/` (8 files)
- Isolated 31+ seconds of heavy git subprocess tests out of core packages into `gitmap/tests/heavy_test/`.
- `go test ./committransfer/...` reduced from >30s down to **0.599s**.
- `go test ./release/...` reduced to **0.216s**.
- Zero circular dependencies repository-wide.
