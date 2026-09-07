# Subtask 01: Pipeline Error Logs, GH Timeout & Local Persistence

## Scope
- Update `gitmap/cmd/pipeline_query.go`:
  - Add `runGHCommandWithCustomTimeout(timeout time.Duration, args ...string) ([]byte, error)`.
  - Use 60-second timeout for `run view ... --log-failed` and `run view ... --log`.
  - Return informative error details if `gh` fails, without swallowing errors.
  - Implement `resolvePipelineDir` checking SQLite `Setting` table key `pipeline.dir`, defaulting to `.gitmap/pipeline`.
  - Check local disk `<pipelineDir>/<runId>.log` before querying remote `gh`.
  - Save downloaded log to `<pipelineDir>/<runId>.log` and structured metadata to `<pipelineDir>/<runId>.json`.
- Update `gitmap/cmd/pipeline_logs.go`:
  - Remove unsolicited interactive prompt `maybeOfferAutoFix`.
  - Only execute auto-repair checks when `--fix` or `-f` is explicitly specified.
  - Fix JSON output when `--json` flag is provided.

## Files Touched
- `gitmap/cmd/pipeline_query.go`
- `gitmap/cmd/pipeline_logs.go`
- `gitmap/cmd/pipeline_test.go`
