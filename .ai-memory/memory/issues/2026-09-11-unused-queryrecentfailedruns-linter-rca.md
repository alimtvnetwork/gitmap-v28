# RCA: Unused Function `queryRecentFailedRuns` & `gofmt` Formatting in CI

## 1. Why it happened
During the refactoring of `pipeline error-logs` to prevent dredging up old historical failures when the current pipeline is passing cleanly, the call to `queryRecentFailedRuns` inside `cli/cmd/pipeline_logs.go` was removed in favor of `collectFailedRuns`. However, the helper definition `queryRecentFailedRuns` remained in `cli/cmd/pipeline_query.go` without any active callers, triggering `golangci-lint`'s `unused` rule. In addition, recent test and persist additions had minor formatting discrepancies detected by `gofmt`.

## 2. How it happened
When `golangci-lint` ran in CI (`Lint Baseline Guard`, `Lint Baseline Diff`, `Full Suite Guard`, `Lint`), it identified `func queryRecentFailedRuns is unused (unused)` as a new finding vs the baseline. This failed all 5 lint guard jobs in GitHub Actions.

## 3. Root Cause
- `queryRecentFailedRuns` in `cli/cmd/pipeline_query.go:175` had its single caller removed in `cli/cmd/pipeline_logs.go`. In Go, an unexported package-level function without callers is flagged as dead/unused code.
- `cli/cmd/pipeline_errorlogs_test.go` and `cli/cmd/pipeline_persist.go` had slight formatting deviations detected by `gofmt -l`.

## 4. Code Fix
1. Remove `queryRecentFailedRuns` from `cli/cmd/pipeline_query.go`.
2. Run `gofmt -w .` across `cli/` to format all Go source files.
3. Validate locally with `python 03-ai-scripts/06-cicd-local-runner.py --filter "golangci-lint (strict)"` and `--all`.
