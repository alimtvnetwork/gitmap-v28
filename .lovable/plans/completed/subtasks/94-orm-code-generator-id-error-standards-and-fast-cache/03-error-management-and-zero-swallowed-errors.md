# Subtask 94.03: Zero-Swallowed Errors & Universal *apperror.AppError Wrapping

## Goal
Eliminate all swallowed errors (`_ = `) and missing `*apperror.AppError` return types across `gitmap/pipelinedb/pipeline_split_ops.go`, `gitmap/dbengine/wrapper.go`, and `gitmap/repodb/repo_db.go`.

## Files Impacted
- `gitmap/pipelinedb/pipeline_split_ops.go`
- `gitmap/dbengine/wrapper.go`
- `gitmap/repodb/repo_db.go`

## Acceptance Criteria
1. `GetStats` in `pipelinedb/pipeline_split_ops.go`: properly handle and propagate errors from all 6 count queries instead of discarding them.
2. `QueryRecentRuns` and `QueryRecentErrorLogs`: row scan errors returned as `*apperror.AppError` rather than silently ignored.
3. `dbengine/wrapper.go`: `WithTransaction` inspects and logs/returns rollback error if rollback fails.
4. `repodb/repo_db.go`: `OptimizeRepoDB`, `OpenRepoDB`, and `InitRepoSchema` wrap errors in `*apperror.AppError` and handle checkpoint/optimize errors cleanly.
5. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
