# Plan 85: Macro Path Expansion, Enhanced Mkdir Engine, Stray Binary RCA & Live Execution Streaming

## 1. Overview & Context

This plan addresses four critical developer ergonomics and operational integrity gaps in Gitmap:
1. **Universal Path & Environment Variable Expansion (`%TEMP%`, `//temp`, `~`, `$VAR`):**
   - In interactive macro builder (`gitmap macro add`) and macro execution (`macro.Execute` / `DirTracker.ProcessCd`), path targets containing `%temp%`, `/temp`, `//temp`, `$TEMP`, and `~` fail because environment variables and platform temp aliases were not resolved before calling Win32/POSIX `os.Chdir`.
   - Resolution: Centralized `macro.NormalizeTargetPath(target, currentDir)` supporting `%VAR%`, `$VAR`, `~`, `/temp`, `//temp`, `\temp`, `\\temp` aliases, and quote unwrapping.
2. **Enhanced `gitmap mkdir` Command (Deep Paths, File Creation, Step Feedback):**
   - Refactored `gitmap mkdir` to support deep multi-level directories (`-p` / `--parents`), file creation (`-f` / `--file`), slash normalization (`//`, `\\`, `/`, `\`), multiple path targets, and progressive console logging detailing directory and file creation hierarchy (`  [DIR]  created: ...`, `  [FILE] created: ...`).
3. **Root Cause Analysis & Cleanup of Stray Binary (`../gitmap.exe`):**
   - Removed obsolete 29.35 MB binary in workspace parent directory (`../gitmap.exe`). Documented RCA in `.lovable/issues/08-stray-parent-binary-rca.md`, updated `Makefile` build target to emit strictly into `bin/gitmap.exe`, and updated `.lovable/strictly-avoid.md`.
4. **Live Execution & Output Streaming in Macro Engine & Interactive Builder:**
   - Streamed child process stdout and stderr in `macro.Execute` via `io.MultiWriter` by default (unless structured output `--json`/`--yaml` is active), added diagnostic stderr dumping on step failures, and provided live command execution with real-time feedback in `macro add` interactive mode (`exec on`/`exec off` toggle).

---

## 2. Task-Specific Rules & Invariants

1. **Strict Relative Paths:** All markdown links and file paths in subtasks and code comments MUST be strictly relative to the repository root. Zero absolute paths or `file:///` URIs.
2. **Coding Guidelines Adherence:** All Go functions <= 15 lines, blank line before every return, affirmative booleans (`is*`, `has*`), zero nested `if` blocks.
3. **Zero Error Swallowing:** All runtime and validation errors wrapped in `*apperror.AppError`.
4. **Binary Deployment Integrity:** The compiled `gitmap.exe` binary must strictly reside in `bin/gitmap.exe` and canonical user location `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`. Never deposit binaries in repository root or parent directories.
5. **No Automatic Releases:** Standard development commits; no automatic version bump or release tagging.

---

## 3. Subtask Decomposition & Execution Status

- [01-task-universal-path-expansion.md](../subtasks/85-macro-path-expansion-mkdir-and-execution-streaming/01-task-universal-path-expansion.md): [COMPLETED] Implemented `NormalizeTargetPath` in `gitmap/macro/expand.go`, wired into `cmd/macro_add_helpers.go` and `macro/record_dir.go`, and verified unit tests.
- [02-task-mkdir-deep-file-and-step-feedback.md](../subtasks/85-macro-path-expansion-mkdir-and-execution-streaming/02-task-mkdir-deep-file-and-step-feedback.md): [COMPLETED] Overhauled `gitmap/cmd/mkdir.go` with deep path handling, file creation (`-f`), multi-OS slash normalization, step-by-step feedback, and full unit test suite in `gitmap/cmd/mkdir_test.go`.
- [03-task-stray-binary-rca-and-cleanup.md](../subtasks/85-macro-path-expansion-mkdir-and-execution-streaming/03-task-stray-binary-rca-and-cleanup.md): [COMPLETED] Documented RCA in `.lovable/issues/08-stray-parent-binary-rca.md`, removed stray parent binary, updated `Makefile` to target `bin/gitmap.exe`, and updated `.lovable/strictly-avoid.md`.
- [04-task-macro-live-execution-and-output-streaming.md](../subtasks/85-macro-path-expansion-mkdir-and-execution-streaming/04-task-macro-live-execution-and-output-streaming.md): [COMPLETED] Overhauled `gitmap/macro/execute.go` to stream stdout/stderr live to terminal by default, dump stderr diagnostics on step failure, and added live execution mode to `gitmap/cmd/macro_add_interactive.go`.
- [05-task-verification-and-ci-gates.md](../subtasks/85-macro-path-expansion-mkdir-and-execution-streaming/05-task-verification-and-ci-gates.md): [COMPLETED] Ran comprehensive unit tests (100% pass), enforced all 4 linters (0 violations), built canonical binary `bin/gitmap.exe`, and deployed to `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.

---

## 4. Verification Summary

1. `macro.ExpandPathAndEnv` and `NormalizeTargetPath`:
   - `cd %temp%` resolves to `os.TempDir()`.
   - `cd //temp` and `cd /temp` resolve to `os.TempDir()`.
   - `cd ~` and `cd ~/subdir` resolve to user home dir.
   - Quoted paths `cd "%temp%"` and paths with spaces resolve correctly.
2. `gitmap mkdir`:
   - `gitmap mkdir -p //temp/test-cli-mkdir/nested/sub` creates hierarchy and displays step logs.
   - `gitmap mkdir -f //temp/test-cli-mkdir/nested/sub/deep/hello.txt` creates parent directories and touches file.
   - Idempotency verified: re-running prints `[DIR] exists:` and `[FILE] exists:`.
3. Stray binary elimination:
   - Parent directory `../gitmap.exe` deleted and verified absent.
   - `Makefile` build target verified to write to `bin/gitmap.exe`.
4. Macro live output streaming:
   - `macro.Execute` streams stdout and stderr live during step execution; dumps stderr diagnostics on error.
   - `macro add` interactive mode executes commands live with real-time output and records steps.
5. Quality Gates:
   - All Go unit tests pass 100%.
   - All 4 linters (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`, `check-error-management.py`) pass with 0 violations.
