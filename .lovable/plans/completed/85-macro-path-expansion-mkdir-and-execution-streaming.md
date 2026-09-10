# Plan 85: Macro Path Expansion, Enhanced Mkdir Engine, Stray Binary RCA & Live Execution Streaming (Consolidated)

> **Task Origin & Problem Initiation:**  
> Initiated from user session report demonstrating four distinct operational failures in Gitmap CLI:  
> 1. `cd %temp%` in interactive macro builder failed with `chdir %temp%: The system cannot find the file specified`.  
> 2. `gitmap mkdir` lacked deep path synthesis, file creation (`-f`), slash normalization, and step-by-step progress output.  
> 3. An obsolete 29.35 MB binary `gitmap.exe` resided in the workspace parent directory (`../gitmap.exe`).  
> 4. Interactive macro builder and macro runner (`gitmap macro run`) suppressed stdout/stderr output unless `--verbose` was passed, leaving commands like `gitmap install vscode` without visible feedback.  
>  
> **Execution Lifecycle:** Completed across 5 subtasks in Phase 1 planning and Phase 2 execution with 100% test pass rate and 0 linter violations.

---

## 1. Overview & Context

This consolidated specification documents the resolution of all four developer ergonomics and operational integrity issues:
1. **Universal Path & Environment Variable Expansion (`%TEMP%`, `//temp`, `~`, `$VAR`):**
   - In interactive macro builder (`gitmap macro add`) and macro execution (`macro.Execute` / `DirTracker.ProcessCd`), path targets containing `%temp%`, `/temp`, `//temp`, `$TEMP`, and `~` failed because environment variables and platform temp aliases were not resolved before calling Win32/POSIX `os.Chdir`.
   - Resolution: Centralized `macro.NormalizeTargetPath(target, currentDir)` supporting `%VAR%`, `$VAR`, `~`, `/temp`, `//temp`, `\temp`, `\\temp` aliases, and quote unwrapping.
2. **Enhanced `gitmap mkdir` Command (Deep Paths, File Creation, Step Feedback):**
   - Refactored `gitmap mkdir` to support deep multi-level directories (`-p` / `--parents`), file creation (`-f` / `--file`), slash normalization (`//`, `\\`, `/`, `\`), multiple path targets, and progressive console logging detailing directory and file creation hierarchy (`  [DIR]  created: ...`, `  [FILE] created: ...`).
3. **Root Cause Analysis & Cleanup of Stray Binary (`../gitmap.exe`):**
   - Removed obsolete 29.35 MB binary in workspace parent directory (`../gitmap.exe`). Documented RCA in `.lovable/issues/08-stray-parent-binary-rca.md`, updated `Makefile` build target to emit strictly into `bin/gitmap.exe`, and updated `.lovable/strictly-avoid.md`.
4. **Live Execution & Output Streaming in Macro Engine & Interactive Builder:**
   - Streamed child process stdout and stderr in `macro.Execute` via `io.MultiWriter` by default (unless structured output `--json`/`--yaml` is active), added diagnostic stderr dumping on step failures, and provided live command execution with real-time feedback in `macro add` interactive mode (`exec on`/`exec off` toggle).

---

## 2. Invariants & Rules Enforced

1. **Strict Relative Paths:** All markdown links and file paths strictly relative to the repository root. Zero absolute paths or `file:///` URIs.
2. **Coding Guidelines Adherence:** All Go functions <= 15 lines, blank line before every return, affirmative booleans (`is*`, `has*`), zero nested `if` blocks.
3. **Zero Error Swallowing:** All runtime and validation errors wrapped in `*apperror.AppError`.
4. **Binary Deployment Integrity:** The compiled `gitmap.exe` binary must strictly reside in `bin/gitmap.exe` and canonical user location `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`. Never deposit binaries in repository root or parent directories.
5. **No Automatic Releases:** Standard development commits; no automatic version bump or release tagging.

---

## 3. Consolidated Subtask Resolutions

### Subtask 01: Universal Path & Environment Variable Expansion
- Implemented `macro.NormalizeTargetPath(target, currentDir string) string` in `gitmap/macro/expand.go`.
- Automatically strips single and double quotes (`"%temp%"` -> `%temp%`).
- Expands Windows `%VAR%` (case-insensitive) and Unix `$VAR`.
- Expands tilde (`~`) to user home directory.
- Maps cross-platform temporary directory aliases (`//temp`, `\\temp`, `/temp`, `\temp`, `/tmp`, `//tmp`, and subpaths) to `os.TempDir()`.
- Wired into `gitmap/cmd/macro_add_helpers.go` (`handleCdCmd`, `executeInteractiveCd`) and `gitmap/macro/record_dir.go` (`DirTracker.resolveTarget`).
- In `macro add`, executing `cd <target>` updates current session working directory and records the step into the macro.
- Verified with unit tests in `gitmap/macro/expand_test.go` and `gitmap/cmd/macro_add_helpers_test.go`.

### Subtask 02: Enhanced Mkdir Engine (Deep Paths, File Creation, Step Feedback)
- Overhauled `gitmap/cmd/mkdir.go`:
  - Flag position independence: `-p` / `--parents`, `-f` / `--file`, and `-v` / `--verbose` can appear anywhere in arguments.
  - Deep recursive directory synthesis: creates all missing ancestor directories with progressive step logs (`  [DIR]  created: ...` or `  [DIR]  exists: ...`).
  - File creation: when `-f` is specified, creates the parent directory hierarchy and touches the file (`  [FILE] created: ...`).
  - Idempotency: re-running on existing paths succeeds gracefully.
- Added comprehensive unit tests in `gitmap/cmd/mkdir_test.go`.

### Subtask 03: Stray Binary Root Cause Analysis, Elimination & Build Target Guard
- Authored Root Cause Analysis in `.lovable/issues/08-stray-parent-binary-rca.md`.
- Safely deleted `../gitmap.exe` (29.35 MB, dated Sept 5, 2026).
- Updated `Makefile` build target to write strictly to `bin/gitmap.exe`.
- Added Total Ban entry in `.lovable/strictly-avoid.md` forbidding binaries outside canonical `bin/` and user AppData locations.

### Subtask 04: Macro Live Execution & Real-Time Output Streaming
- Updated `gitmap/macro/execute.go`:
  - Child process stdout and stderr stream live to `os.Stdout` and `os.Stderr` via `io.MultiWriter` by default (guarding structured output `--json`/`--yaml`).
  - On step failure, detailed stderr diagnostic logs are printed.
- Updated `gitmap/cmd/macro_add_interactive.go`:
  - Added live command execution: non-helper commands typed at `Step N>` execute immediately in the session working directory and stream output live to the terminal.
  - Added `exec on` / `exec off` in-builder toggle (default `exec on`).
  - Synchronized command steps recorded into `*steps`.

### Subtask 05: Verification, Quality Gates & Consolidation
- Ran full test suites across `macro` and `cmd` packages (100% pass).
- Passed all 4 linters:
  - `python linter-scripts/check-nested-ifs.py` (0 violations)
  - `python linter-scripts/check-boolean-guidelines.py` (0 violations)
  - `python linter-scripts/check-relative-paths.py` (0 violations)
  - `python linter-scripts/check-error-management.py` (0 violations)
- Compiled canonical binary: `bin/gitmap.exe`
- Deployed canonical binary: `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`

---

## 4. Final Verification Evidence

```powershell
# Testing mkdir -p with //temp alias
& "$env:LOCALAPPDATA\gitmap-cli\gitmap.exe" mkdir -p //temp/test-cli-mkdir/nested/sub
# Output:
#   [DIR]  created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir
#   [DIR]  created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir\nested
#   [DIR]  created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir\nested\sub
# ✔ Created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir\nested\sub

# Testing mkdir -f to create file with deep ancestor hierarchy
& "$env:LOCALAPPDATA\gitmap-cli\gitmap.exe" mkdir -f //temp/test-cli-mkdir/nested/sub/deep/hello.txt
# Output:
#   [DIR]  created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir\nested\sub\deep
#   [FILE] created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir\nested\sub\deep\hello.txt
# ✔ Created: C:\Users\Alim\AppData\Local\Temp\test-cli-mkdir\nested\sub\deep\hello.txt
```
