# CI/CD Issue 53: Clean Success Pipeline Fallback Overwrite RCA

## 1. Symptom
In GitHub Actions CI run `#35263724019` on commit `e9cb7384`, four jobs failed with the exact same error in `TestBuildErrorLogsPayloadCleanSuccess`:
- `Full Suite Guard (Go Test Suite)` | Step: `Run full test suite with test log capture`
- `Build (macOS)` | Step: `Run tests`
- `Build (Ubuntu)` | Step: `Run tests`
- `Build (Windows)` | Step: `Run tests`

```text
=== RUN   TestBuildErrorLogsPayloadCleanSuccess
    pipeline_errorlogs_test.go:288: expected clean success payload, got failure
--- FAIL: TestBuildErrorLogsPayloadCleanSuccess (0.00s)
FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline	0.038s
```

---

## 2. Root Cause
In `cli/cmdpipeline/pipeline_logs.go`, commit `547573e` introduced fallback retrieval for pipeline errors when `len(failedRuns) == 0`:
```go
func populateRunsIntoPayload(repo string, runs []ghRunItem, p *PipelineErrorLogsPayload) {
	initLatestRunMeta(p, runs[0])
	checkAndApplyRunningState(p, runs)
	failedRuns := resolveFailedRunsForPayload(repo, runs)
	if len(failedRuns) > 0 {
		populateFailedRunsPayload(repo, failedRuns, p)

		return
	}
	_ = ApplyPreviousRunFallbackToPayload(p, repo, runs)
}
```
When the latest pipeline run completed cleanly with `Conclusion == "success"`, `failedRuns` is empty (`len(failedRuns) == 0`). Consequently, execution proceeded to `ApplyPreviousRunFallbackToPayload(p, repo, runs)`. This function inspected older runs in history, found an older failed run (in `TestBuildErrorLogsPayloadCleanSuccess`, run #9), and mutated `p.Conclusion = "failure"`, corrupting what should have been a clean success payload.

---

## 3. Resolution
1. Added guard clause in `cli/cmdpipeline/pipeline_logs.go` inside `populateRunsIntoPayload`:
```go
if p.Conclusion == "success" || p.IsRunning {
	return
}
_ = ApplyPreviousRunFallbackToPayload(p, repo, runs)
```
2. Ran `gofmt -w cli/cmdpipeline/pipeline_logs.go` and verified 100% cleanliness across all 2,853 Go files.
3. Verified zero violations with `check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`, and `check-enum-guidelines.py`.

---

## 4. Prevention & Learnings
- **Success & In-Progress Status Invariant**: Fallback lookups must never overwrite successful runs or active executions. Historical fallbacks are strictly for missing or failing runs where latest run data is unavailable or non-conclusive.
- **Sentinel Regression Tests**: Retain `TestBuildErrorLogsPayloadCleanSuccess` to prevent regressions in payload status resolution.
