# Subtask 03: Generic Terminal & Worker Package, Golang DRY Spec & Error Manage Enhancement

- **Subtask ID**: `subtask-03-generic-terminal-worker-package-and-coding-guidelines`
- **Parent Task**: `gitmap-efficiency-and-coding-guidelines-enhancement`
- **Owned Files**:
  - `cli/termgeneric/header.go`
  - `cli/termgeneric/color.go`
  - `cli/termgeneric/progress.go`
  - `cli/termgeneric/worker_pool.go`
  - `cli/termgeneric/termgeneric_test.go`
  - `02-spec/02-coding-guidelines/14-dry-code-architecture/readme.md`
  - `02-spec/03-error-manage/02-improvements.md`
- **Assigned Worker**: Worker 03

---

## 1. Objectives
1. Implement generic, reusable terminal & concurrency package in `cli/termgeneric/`:
   - `header.go`: Declarative box header renderer (`╭─ ... ──╮`) with customizable border styling, automatic padding calculation, and titles.
   - `color.go`: Clean ANSI coloring helpers (Red, Green, Cyan, Yellow, Magenta, Dim, Bold, Reset) that any codebase can drop in and use without heavy third-party dependencies.
   - `progress.go`: Clean, flicker-free terminal progress bar and spinner with deterministic runewidth handling and line-clear mechanics (`\r\033[2K`).
   - `worker_pool.go`: Generic Go worker pool pattern (`WorkerPool[T, R]`) using channels, bounded capacity, cancellation context, and error collection.
2. In `cli/termgeneric/termgeneric_test.go`:
   - Add unit tests for header rendering, ANSI color formatting, progress bar steps, and concurrent worker execution.
3. In `02-spec/02-coding-guidelines/14-dry-code-architecture/readme.md`:
   - Document the DRY code architecture principles for Golang packages, detailing the generic terminal framework and worker design patterns.
4. In `02-spec/03-error-manage/02-improvements.md`:
   - Enhance error management guidance with structured AppError wrappers, context preservation, and exit code conventions without removing any existing solutions.

---

## 2. Verification
- `go test ./termgeneric -v`
- Verify documentation integrity in `02-spec/02-coding-guidelines/` and `02-spec/03-error-manage/`.
