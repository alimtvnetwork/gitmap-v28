# Subtask 02: Pipeline Error Logs JSON and Terminal Multi-Card Rendering

## Scope
- Update `PipelineErrorLogsPayload` to include structured `FailedRuns []FailedRunItem` and `ConcatenatedErrors string`.
- Update `writeOrRenderErrorLogs` and `formatErrorLogContent`:
  - In JSON view (`--json`): Render clean structured JSON payload containing `failedRuns`, with each job's name, step, and clean failure summary.
  - In Terminal view: Render formatted failure cards for each failed workflow run:
    - Workflow Name & Run ID
    - Job Name & Step Name
    - Extracted Assertion Failure & Error Detail
    - Direct Run URL
- If local errors exist (`.gitmap/last_error.log`), render them clearly alongside.
- Keep all functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.

## Files Touched
- `gitmap/cmd/pipeline_logs.go`
- `gitmap/cmd/pipeline_helpers.go`
- `gitmap/cmd/pipeline_errorlogs_test.go`
