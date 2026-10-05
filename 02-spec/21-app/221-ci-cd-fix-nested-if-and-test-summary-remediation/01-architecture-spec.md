# 221 - CI/CD Fix: Nested If Linter & Test Failure Summary Remediation

## 1. Executive Summary & Root Cause Analysis

### 1.1 Finding 1: Nested If Violation in `pipeline_persist.go`
- **Location:** `cli/cmdpipeline/pipeline_persist.go:231`
- **Violation:**
  ```go
  if pipeDb.HasDetailErrorLog(runId) && pipeDb.HasCompactErrorLog(runId) {
      if existing, ok := queryDetailLogFromDb(pipeDb, runId); ok && !isCorruptOrFallbackErrorLog(existing) {
          return
      }
  }
  ```
- **Linters Failed:** `Nested If Linter` (`linter-scripts/check-nested-ifs.py`) and `Boolean & Enum Linter` (`linter-scripts/check-enum-and-boolean.py`).
- **Remediation:** Extract predicate `isLogAlreadyPersisted(pipeDb *pipelinedb.PipelineSplitDb, runId uint64) bool` and use an early return guard:
  ```go
  if isLogAlreadyPersisted(pipeDb, runId) {
      return
  }
  ```

### 1.2 Finding 2: `isStrongerSummary` Prevents Specific Test Location Summaries
- **Location:** `cli/cmdpipeline/pipeline_error_extract.go:505`
- **Violation:**
  In `isStrongerSummary`:
  `if isTestFailureSummary(current) && !isTestFailureSummary(candidate) { return false }`
  Because `isTestFailureSummary` treats `--- FAIL: ` as a test failure summary, a generic test banner like `--- FAIL: TestSample (0.00s)` prevents more specific assertion failures (e.g. `sample_test.go:12: Expected true to be false`) from replacing it.
- **Test Failed:** `TestParseFailedLogLines` in `cli/cmdpipeline/pipeline_errorlogs_test.go:123`.
- **Suite Impact:** Caused `Full Suite Guard` (`go test ./...`) and cross-platform build jobs (`ubuntu-latest`, `macos-latest`, `windows-latest`) to fail with exit code 1.
- **Remediation:**
  Ensure that location summaries (`isLocationSummary(candidate)`) and assertion summaries (containing `Expected `, `AssertionError`, etc.) take precedence over generic test banners (`--- FAIL: `, `FAIL\t`).

---

## 2. Architecture & Design Specifications

### 2.1 Flattening `pipeline_persist.go`
```go
func isLogAlreadyPersisted(pipeDb *pipelinedb.PipelineSplitDb, runId uint64) bool {
	if !pipeDb.HasDetailErrorLog(runId) || !pipeDb.HasCompactErrorLog(runId) {
		return false
	}
	existing, ok := queryDetailLogFromDb(pipeDb, runId)
	return ok && !isCorruptOrFallbackErrorLog(existing)
}
```
In `persistLogToRepoSplitDb`:
```go
func persistLogToRepoSplitDb(repo string, runId uint64, logContent string) {
	pipeDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return
	}
	defer pipeDb.Close()

	if isLogAlreadyPersisted(pipeDb, runId) {
		return
	}

	clean := extractCleanErrorLines(logContent)
	workflow := pipeDb.GetRunWorkflowName(runId)
	detail, compact := buildPersistRecords(repo, runId, workflow, logContent, clean)
	_ = pipeDb.RecordDualErrorLog(detail, compact)
}
```

### 2.2 Re-architecting `isStrongerSummary`
A test location with line number (e.g. `sample_test.go:12: Expected true to be false`) or an assertion error provides crucial debugging information that must override a generic top-level runner banner `--- FAIL: TestSample (0.00s)`.
1. Distinguish between a specific assertion failure / location summary and a generic test status banner (`isTestStatusBanner(s)`).
2. If `candidate` is `isLocationSummary(candidate)` or has explicit assertion details (`Expected `, `AssertionError`), and `current` is a generic banner (`--- FAIL: `, `FAIL\t`), `candidate` MUST be considered stronger.
3. If both are test failure summaries, a location-aware failure summary outranks a banner.

---

## 3. Verification Criteria
1. `python linter-scripts/check-nested-ifs.py` exits 0.
2. `python linter-scripts/check-enum-and-boolean.py` exits 0.
3. `go test -v -run TestParseFailedLogLines ./cmdpipeline` exits 0.
4. `go test -run=^$ ./...` (compile gate) exits 0.
5. All 5 policy linters exit 0.
