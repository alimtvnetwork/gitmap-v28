# Subtask 05: Unit Tests, Parity Verification & CI Local Runner

## Objective
Author comprehensive regression tests and verify against all repository quality gates:
1. `gitmap/cmd/corrupted_dirs_test.go`:
   - Test detecting corrupted directories with embedded newlines, ANSI escape sequences, `Default:`, and literal `~`.
   - Test that protected paths (`/`, `$HOME`, `/usr/local/bin`, `/tmp`) are strictly rejected from being deleted.
   - Test CWD escaping: when process is inside corrupted directory, verify `CleanCorruptedDirs` changes directory out before unlinking.
   - Test asset recovery: verify binary is copied to target folder with `0755` permissions before the corrupted source directory is deleted.
   - Test `--dry-run` flag leaves filesystem untouched.
2. Verify shell script cleanup:
   - Test `cleanup_corrupted_install_dirs` in a mock folder with actual corrupted directory names containing newlines and ANSI sequences.
3. Quality gates:
   - Run `go test -C gitmap ./cmd/... -run Corrupted -v`.
   - Run `go vet -C gitmap ./...`.
   - Run `python linter-scripts/check-nested-ifs.py`.
   - Run `python linter-scripts/check-enum-and-boolean.py`.

## Files Affected
- `gitmap/cmd/corrupted_dirs_test.go`
