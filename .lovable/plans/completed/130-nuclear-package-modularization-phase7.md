# Plan 130: Nuclear Package Modularization (Phase 7), Heavy Test Isolation & Test Inventory Duration Estimation

> **Version:** 1.0.0
> **Scope:** Monolith Package Decomposition (`cli/cmd`), Heavy Test Segregation (`cli/tests/heavy_test`), and Test Inventory Profiling
> **Status:** completed

## Objectives & High-Level Architecture

1. **Nuclear Monolith Package Decomposition**:
   - Decompose oversized command clusters out of `cli/cmd` into clean, acyclic subpackages:
     - `cli/cmdinstaller/` (~34 files: `installer*.go` and tests)
     - `cli/cmdchrome/` (~22 files: `chrome*.go` and tests)
     - `cli/cmdsetup/` (~21 files: `setup*.go` and tests)
     - `cli/cmdinstall/` (~81 files: `install*.go` and tests)
   - Total files extracted from `cli/cmd`: ~158 files (~21,753 lines of code).
   - `cli/cmd` size drastically reduced from 916 files to ~758 files.

2. **Strict DAG (Zero Circular Dependencies)**:
   - Root `cli/cmd` package depends on domain packages (`cmdinstaller`, `cmdchrome`, `cmdsetup`, `cmdinstall`).
   - Domain packages only depend downward on leaf packages (`constants`, `model`, `store`, `apperror`, `fsutil`, `cliexit`).
   - Bridge forwarders and constructor hooks (`exports.go`) ensure full backward compatibility.

3. **Heavy Test Isolation**:
   - Extract tests executing `exec.Command`, real sockets, or `time.Sleep` into `cli/tests/heavy_test/` (`package heavy_test`).
   - Targets: `cli/tests/release_test`, `cli/tests/fixrepo_test`, `installctx_e2e_test`, and `probe/background_test`.
   - Core packages contain 100% fast in-memory unit tests (<0.01s).

4. **Test Inventory Duration Estimation & Dirty Tracking**:
   - Update `.lovable/test-inventory.json` using `03-ai-scripts/33-test-inventory-generator.py`.
   - Record estimated durations for all tests. Mark heavy/slow tests (`duration_sec >= 4.0s`).
   - Record modified files in `.lovable/temp/recent-file-changes.json` under atomic lock.

5. **Quality Gates & Verification**:
   - Verify `go vet -C cli ./...` and `go build -C cli .` pass with exit code 0.
   - Run full 38-gate suite via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests`.
   - Verify zero nested ifs, affirmative booleans, and relative paths.

## Subtask Directory

- [01-heavy-test-isolation.md](subtasks/130-nuclear-package-modularization-phase7/01-heavy-test-isolation.md)
- [02-modularize-cmdinstaller.md](subtasks/130-nuclear-package-modularization-phase7/02-modularize-cmdinstaller.md)
- [03-modularize-cmdchrome.md](subtasks/130-nuclear-package-modularization-phase7/03-modularize-cmdchrome.md)
- [04-modularize-cmdsetup.md](subtasks/130-nuclear-package-modularization-phase7/04-modularize-cmdsetup.md)
- [05-modularize-cmdinstall.md](subtasks/130-nuclear-package-modularization-phase7/05-modularize-cmdinstall.md)
- [06-test-inventory-sync-and-verification.md](subtasks/130-nuclear-package-modularization-phase7/06-test-inventory-sync-and-verification.md)


## Completion Summary
- **Decomposition**: 160 files extracted from `cli/cmd` into 4 new acyclic packages:
  - `cli/cmdinstaller`: 34 files
  - `cli/cmdchrome`: 22 files
  - `cli/cmdsetup`: 21 files
  - `cli/cmdinstall`: 83 files (including `osdetect*` and `agy_install`)
- **Heavy Test Segregation**: Isolated heavy subprocess/git tests into `cli/tests/heavy_test/` (`package heavy_test`).
- **Test Inventory**: Updated `.lovable/test-inventory.json` with 3,534 tests indexed across 136 packages (75 slow, 3459 fast).
- **DAG Hygiene**: Zero circular dependencies (`go vet -C cli ./...` exit 0). Bridge forwarders and delegate hooks in `cli/cmd/clihelpers.go`.
- **Quality Gates**: All 38 CI/CD quality gates passed (`--no-tests`).
