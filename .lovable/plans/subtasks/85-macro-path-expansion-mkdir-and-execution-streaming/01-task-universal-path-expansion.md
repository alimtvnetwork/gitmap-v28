# Subtask 01: Universal Path & Environment Variable Expansion

## 1. Objective
Implement `NormalizeTargetPath(target, currentDir string) string` in `gitmap/macro/expand.go`, support temp aliases (`/temp`, `//temp`, `\temp`, `\\temp`, `/tmp`, `//tmp`), quote unwrapping, tilde expansion, and integrate into `gitmap/cmd/macro_add_helpers.go` (`handleCdCmd`, `executeInteractiveCd`) and `gitmap/macro/record_dir.go` (`DirTracker.resolveTarget`).

## 2. Target Files
- `gitmap/macro/expand.go`
- `gitmap/macro/expand_test.go`
- `gitmap/cmd/macro_add_helpers.go`
- `gitmap/cmd/macro_add_helpers_test.go`
- `gitmap/macro/record_dir.go`

## 3. Requirements
- Handle `%VAR%` (Windows case-insensitive), `$VAR` (Unix), `~` (home directory), and quoted paths (`"%temp%"`).
- Handle cross-platform temporary directory aliases: `/temp`, `//temp`, `\temp`, `\\temp`, `/tmp`, `//tmp` mapping to `os.TempDir()`.
- Support subpaths of temp aliases (e.g. `//temp/test` -> `<os.TempDir()>\test`).
- In `handleCdCmd`, do not truncate paths on spaces; use `extractCommandArgument(line)`.
- When `cd <target>` is executed in `macro_add_helpers.go`, stage the step in `*steps` so the saved macro stays synchronized.
- Ensure all Go functions <= 15 lines with blank lines before every return and positive booleans.
- Comprehensive unit tests in `expand_test.go` and `macro_add_helpers_test.go`.
