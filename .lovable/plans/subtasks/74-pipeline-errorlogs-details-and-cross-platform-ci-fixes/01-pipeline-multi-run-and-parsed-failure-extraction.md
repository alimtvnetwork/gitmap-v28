# Subtask 01: Pipeline Multi-Run Query and Parsed Failure Extraction

## Scope
- Update `gitmap/cmd/pipeline_logs.go` and `gitmap/cmd/pipeline_query.go` to query all failed workflow runs corresponding to the latest commit/push rather than stopping at the first failure.
- In `gitmap/cmd/pipeline_error_extract.go`, parse raw `Job\tStep\tTimestamp\tLogText` lines into a structured model:
  - `FailedJobItem`: `JobName`, `StepName`, `FailureSummary`, `ErrorLines`, `URL`.
  - `FailedRunItem`: `WorkflowName`, `RunId`, `Conclusion`, `URL`, `Jobs`.
- Extract root-cause error lines and assertion details (`--- FAIL:`, `FAIL\t`, `Error:`, `Expected ...`, `exit status ...`, `panic:`) with context.
- Keep all functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.

## Files Touched
- `gitmap/cmd/pipeline_logs.go`
- `gitmap/cmd/pipeline_query.go`
- `gitmap/cmd/pipeline_error_extract.go`
- `gitmap/cmd/pipeline_helpers.go`
