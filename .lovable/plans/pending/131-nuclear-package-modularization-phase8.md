# Plan 131: Nuclear Package Modularization (Phase 8), Heavy Test Isolation & Test Inventory Duration Estimation

> **Version:** 1.0.0
> **Scope:** Monolith Package Decomposition (`cli/cmd`), Heavy Test Segregation (`cli/tests/heavy_test`), and Test Inventory Profiling
> **Status:** pending

## Objectives & High-Level Architecture

1. **Nuclear Monolith Package Decomposition**:
   - Decompose oversized command clusters out of `cli/cmd` (754 files) into clean, acyclic subpackages:
     - `cli/cmdclone/` (~74 files: `clone*.go` and tests)
     - `cli/cmdupdate/` (~23 files: `update*.go` and tests)
     - `cli/cmdpull/` (~21 files: `pull*.go`, `push*.go` and tests)
   - Total files extracted: ~118 files (~15,000+ lines of code).
   - `cli/cmd` size reduced from 754 to ~636 files.

2. **Strict DAG (Zero Circular Dependencies)**:
   - Root `cli/cmd` package depends on domain packages (`cmdclone`, `cmdupdate`, `cmdpull`).
   - Domain packages only depend downward on leaf packages (`constants`, `model`, `store`, `apperror`, `fsutil`, `cliexit`).
   - Bridge forwarders and constructor hooks (`exports.go`) ensure full backward compatibility and zero cycles.

3. **Heavy Test Isolation**:
   - Extract remaining tests with duration >= 1.0s or real subprocess/git calls into `cli/tests/heavy_test/` (`package heavy_test`).
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

- [01-heavy-test-isolation-phase8.md](subtasks/131-nuclear-package-modularization-phase8/01-heavy-test-isolation-phase8.md)
- [02-modularize-cmdclone.md](subtasks/131-nuclear-package-modularization-phase8/02-modularize-cmdclone.md)
- [03-modularize-cmdupdate.md](subtasks/131-nuclear-package-modularization-phase8/03-modularize-cmdupdate.md)
- [04-modularize-cmdpull.md](subtasks/131-nuclear-package-modularization-phase8/04-modularize-cmdpull.md)
- [05-test-inventory-sync-and-verification.md](subtasks/131-nuclear-package-modularization-phase8/05-test-inventory-sync-and-verification.md)
