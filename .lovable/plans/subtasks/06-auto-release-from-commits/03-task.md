# Subtask 03: CLI Command & UI Rendering

## Goal
Wire the logic into the CLI.

## Requirements
1. **File**: `gitmap/cmd/release_scan_commits.go`
   - Implement `runReleaseScanCommits(args []string)`
   - Parse `--all` flag.
   - Fetch `git log <last_commit>..HEAD --oneline` (or all if `--all` or no last commit).
   - Pass lines through `ParseVersionFromCommit`.
   - Execute via `ExecuteCommitActions`.
   - Write new `HEAD` to state.
   - Render output using a nice tree summary (e.g., `✓` for created, `~` for skipped).
2. **File**: `gitmap/cmd/release.go` or `gitmap/cmd/releaseargs.go`
   - Register the `scan-commits` (alias `rsc`) command.
3. **Constraints**:
   - 15 line max per function.
