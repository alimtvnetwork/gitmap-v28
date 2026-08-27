# Subtask 02: Git Executor & Business Logic

## Goal
Implement the Git execution logic that creates branches and tags.

## Requirements
1. **File**: `gitmap/release/scan_executor.go`
   - Define a struct `ScanCommitAction` containing `CommitHash`, `Version`, `IsBranchCreated`, `IsBranchSkipped`, `IsTagCreated`, `IsTagSkipped`.
   - Implement `ExecuteCommitActions(repoDir string, commits []ParsedCommit) ([]ScanCommitAction, error)`
   - For each commit, verify if `release/<version>` branch exists. If not, create it at `CommitHash`.
   - Verify if `<version>` tag exists. If not, create it at `CommitHash`.
   - Track creations vs skips in the returned `ScanCommitAction`.
2. **Constraints**:
   - `ParsedCommit` struct will have `Hash` and `Message` and `Version`.
   - Use `exec.Command` carefully. Wrap errors.
   - 15 line max per function.
