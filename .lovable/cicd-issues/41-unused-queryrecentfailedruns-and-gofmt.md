# CI/CD Issue: Unused Function `queryRecentFailedRuns` & `gofmt` Formatting in CI

- Job: Lint Baseline Guard, Lint Baseline Diff, Full Suite Guard, Lint
- Type: FAIL
- Detected: 2026-09-11T21:19:00+08:00
- Status: resolved

## Error
```text
gitmap\cmd\pipeline_query.go:175:6: func `queryRecentFailedRuns` is unused (unused)
func queryRecentFailedRuns(repo string, limit int) []ghRunItem {
     ^
```

## Root Cause
`queryRecentFailedRuns` in `gitmap/cmd/pipeline_query.go` had its single caller removed when historical failed runs dredging was disabled in `pipeline error-logs`. In addition, `cmd/pipeline_errorlogs_test.go` and `cmd/pipeline_persist.go` had minor formatting differences detected by `gofmt`.

## Fix Applied
1. Removed `queryRecentFailedRuns` from `gitmap/cmd/pipeline_query.go`.
2. Formatted all Go files using `gofmt -w .`.
3. Verified zero lint issues with `python 03-ai-scripts/06-cicd-local-runner.py --filter "golangci-lint (strict)"` and full local runner.
