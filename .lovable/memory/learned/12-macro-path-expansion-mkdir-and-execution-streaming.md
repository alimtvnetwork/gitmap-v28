# 12: Macro Path Expansion, Enhanced Mkdir Engine, Stray Binary RCA & Live Execution Streaming

## Context & Problem Statement

1. **Interactive Macro Path Expansion Failure (`cd %temp%`):**
   In the interactive macro builder (`gitmap macro add <name>`), typing `cd %temp%` threw `chdir %temp%: The system cannot find the file specified`. `os.Chdir` was called with the literal string `"%temp%"` because environment variables, tilde (`~`), and cross-platform temp aliases were not normalized prior to invoking system directory change APIs.

2. **`gitmap mkdir` Command Gaps:**
   The `mkdir` command only checked positional `args[0] == "-p"`, had zero support for creating files alongside parent hierarchies (`-f`), failed when flags appeared after targets, lacked multi-OS slash normalization, and provided no step-by-step progress logging.

3. **Stray Parent Binary (`../gitmap.exe`):**
   An obsolete 29.35 MB binary resided in the parent workspace (`../gitmap.exe`). Build guides historically prescribed `cd gitmap && go build -o ../gitmap .` which, when executed from the repository root, deposited the binary outside the repository.

4. **Output Suppression & Missing Live Execution:**
   - In `macro.Execute`, stdout and stderr were buffered in memory unless `--verbose` was passed, preventing real-time feedback and obscuring diagnostic error details on failures.
   - In `gitmap macro add`, commands entered were only appended to in-memory steps without live execution, causing users to see empty directory listings when inspecting results (`ls`).

---

## Architectural Decisions & Implementations

1. **Universal Path Normalization (`macro.NormalizeTargetPath`):**
   - Created `macro.NormalizeTargetPath(target, currentDir string) string` in `gitmap/macro/expand.go`.
   - Automatically strips single and double quotes (`"%temp%"` -> `%temp%`).
   - Resolves `%VAR%` (Windows case-insensitive), `$VAR` (Unix), `~` (home directory).
   - Resolves cross-platform temp aliases (`//temp`, `\\temp`, `/temp`, `\temp`, `/tmp`, `//tmp`, and subpaths) to `os.TempDir()`.
   - Normalizes slashes and path traversal via `filepath.Clean`.
   - Wired into `gitmap/cmd/macro_add_helpers.go` (`handleCdCmd`, `executeInteractiveCd`) and `gitmap/macro/record_dir.go` (`DirTracker.resolveTarget`).

2. **Enhanced `gitmap mkdir` Engine:**
   - Overhauled `gitmap/cmd/mkdir.go` to support `-p` / `--parents`, `-f` / `--file`, and `-v` / `--verbose` in any argument position.
   - Traverses and synthesizes missing directory hierarchies, printing progressive logs (`  [DIR]  created: ...` or `  [DIR]  exists: ...`).
   - When `-f` / `--file` is supplied, synthesizes parent directories and touches the file (`  [FILE] created: ...`).
   - Fully idempotent: re-running on existing paths succeeds without error.

3. **Stray Binary Elimination & Build Target Guard:**
   - Authored Root Cause Analysis in `.lovable/issues/08-stray-parent-binary-rca.md`.
   - Safely deleted `../gitmap.exe`.
   - Updated `Makefile` build target to write strictly to `bin/gitmap.exe`.
   - Updated `.lovable/strictly-avoid.md` with a Total Ban on depositing binaries outside canonical `bin/` and user AppData locations.

4. **Live Execution & Output Streaming:**
   - Updated `gitmap/macro/execute.go` to stream child process output to `os.Stdout` and `os.Stderr` via `io.MultiWriter` by default (guarding structured output `--json`/`--yaml`).
   - Added diagnostic stderr dumping on step failure so errors are immediately visible.
   - Added live execution in `gitmap/cmd/macro_add_interactive.go` with interactive `exec on` / `exec off` toggle (default `exec on`).

---

## Verification & Quality Gates

- **Unit Tests:** All unit tests in `gitmap/macro` and `gitmap/cmd` pass 100%.
- **Linters:**
  - `python linter-scripts/check-nested-ifs.py` (0 violations)
  - `python linter-scripts/check-boolean-guidelines.py` (0 violations)
  - `python linter-scripts/check-relative-paths.py` (0 violations)
  - `python linter-scripts/check-error-management.py` (0 violations)
- **Deployment:** Compiled canonical binary to `bin/gitmap.exe` and deployed to `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.
