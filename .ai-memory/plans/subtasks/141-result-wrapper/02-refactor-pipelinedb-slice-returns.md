# Subtask 02: Refactor Pipeline DB Subsystem Slice Returns to ResultSlice[...]

Parent Plan: [141-result-wrapper-and-slice-returns.md](../../pending/141-result-wrapper-and-slice-returns.md)

## Goals
1. Refactor `cli/pipelinedb/pipeline_split_ops.go`:
   - `QueryRecentRuns(limit int) result.ResultSlice[PipelineRunRecord]`
   - `collectRecentRuns(rows *sql.Rows) result.ResultSlice[PipelineRunRecord]`
   - `QueryRecentErrorLogs(limit int) result.ResultSlice[PipelineErrorRecord]`
   - `collectRecentErrors(rows *sql.Rows) result.ResultSlice[PipelineErrorRecord]`
   - `QueryDetailedErrors(limit int) result.ResultSlice[PipelineErrorRecord]`
   - `QueryCompactErrors(limit int) result.ResultSlice[PipelineCompactErrorRecord]`
   - `collectCompactErrors(rows *sql.Rows) result.ResultSlice[PipelineCompactErrorRecord]`
   - `QueryCachedErrorRunIds() result.ResultSlice[uint64]`
   - `QueryLastFailedRuns(limit int) result.ResultSlice[PipelineRunRecord]`
   - `QueryLastNFailedRuns(limit int) result.ResultSlice[PipelineRunRecord]`
   - `QueryErrorLogsByRunId(runId uint64) result.ResultSlice[PipelineErrorRecord]`
   - `QueryDetailedErrorLogsByRunId(runId uint64) result.ResultSlice[PipelineErrorRecord]`
   - `QueryCompactErrorLogsByRunId(runId uint64) result.ResultSlice[PipelineCompactErrorRecord]`
2. Modernize callers in:
   - `cli/cmdpipeline/pipeline_db_ops.go`
   - `cli/cmdpipeline/pipeline_history.go`
   - `cli/cmdpipeline/pipeline_compact_test.go`
   - `cli/pipelinedb/pipeline_split_db_test.go`
   - `cli/pipelinedb/pipeline_split_ops.go` internal callers (`QueryCachedErrorRunIdMap`)

## Acceptance Criteria
- [x] All 13 pipeline DB slice return functions use `result.ResultSlice[...]`.
- [x] Callers modernized with outer inspection methods.
- [x] Subtask completed and logged.
