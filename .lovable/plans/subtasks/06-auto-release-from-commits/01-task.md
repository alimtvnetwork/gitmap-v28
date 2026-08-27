# Subtask 01: State Management & Regex Parsing

## Goal
Implement the state manager and regex parsers for the commit scanning feature.

## Requirements
1. **File**: `gitmap/release/scan_state.go`
   - Implement `ReadLastScannedCommit(repoDir string) (string, error)` that reads `commit_scan_state.json` inside `.gitmap/`.
   - Implement `WriteLastScannedCommit(repoDir, commitHash string) error` that writes to it.
2. **File**: `gitmap/release/scan_parser.go`
   - Implement `ParseVersionFromCommit(commitMessage string) (string, bool)`
   - Must use Regex to detect patterns like `bump version to 1.2.0`, `chore(release): bump version to v1.2.0`, `release v1.14.0`.
   - Must return the extracted version (e.g., `v1.2.0` or `1.2.0`) and a boolean `isFound`.
3. **Tests**: Create `scan_state_test.go` and `scan_parser_test.go` with semantic test names.
4. **Constraints**:
   - 15 line max per function.
   - Use `apperror`.
   - Booleans must start with `is`, `has`, `can`, `should`.
