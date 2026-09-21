# RCA-13: Pipeline Cache TTL Bypass and Table Chronological Order Inversion

**Date:** 2026-09-21
**Component:** `cli/cmdpipeline`, `cli/pipelinedb`
**Status:** Resolved
**Target Release:** `v6.288.0`

---

## 1. Symptom & User Context

When running `gitmap pe` from outside a git repository (such as `C:\Users\Administrator`), users observed two distinct anomalies:
1. **Infinite Stale Cache (TTL ignored):** Despite a 5-second TTL specification (`GITMAP_PIPELINE_CACHE_TTL_SEC`), running `gitmap pe` repeatedly across several minutes continued serving old cached telemetry: `⚡ Served from local SQLITE DB cache (commit 874f85c)`.
2. **Table Chronological Inversion:** In the `● Recent Commits Pipeline Summary` table, the previous commit (`874f85c`) was labeled as `latest`, while the newer commit (`2764db6`) was labeled as `-1`, showing inverted chronological ordering.

---

## 2. Grounded Root Cause Analysis

### Part A: Indefinite Completed Success Cache Hit
In `cli/cmdpipeline/pipeline_cache_eval.go`:
```go
func evaluateDecisionFromRuns(...) {
    ...
    if checkLatestCompletedSuccessCacheHit(latest) {
        return buildCacheHitDecision(dbRuns, latest.Sha, "latest_completed_success")
    }
}
```
`checkLatestCompletedSuccessCacheHit` checked whether `latest.Status == "completed" && latest.Conclusion == "success"`. Because the last run for commit `874f85c` had succeeded, this condition was permanently true on all subsequent invocations, completely bypassing the 5-second TTL check (`checkTtlCacheHit`).

### Part B: Primary Key Autoincrement Order Inversion
In `cli/pipelinedb/pipeline_split_ops.go`:
```sql
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha, ...
FROM PipelineRun ORDER BY PipelineRunId DESC LIMIT ?;
```
Runs were ordered by `PipelineRunId DESC` (SQLite autoincrement integer). When multiple runs were recorded or updated across turns, older runs had higher `PipelineRunId` values than newer runs. Consequently, `dbRuns[0]` was the older commit (`874f85c`), not the newer commit (`2764db6`).
Furthermore, `resolveCachedRunsOrFetch` re-queried the database and passed these inverted rows to `RenderHistorySummaryTable`, which grouped and rendered `874f85c` as offset `0` (`latest`) and `2764db6` as `-1`.

### Part C: Unformatted Go Test File Causing GitHub CI Failure
In commit `2764db6`, `cli/cmd/llm/llm_train_test.go` was not formatted with `gofmt`, triggering a failure in GitHub Actions step `test_gofmt_check_clean_repo`.

---

## 3. Corrective Implementation

1. **Enforce Strict 5s TTL in Cache Evaluation:**
   - In `cli/cmdpipeline/pipeline_cache_eval.go`: Completely eliminated `checkLatestCompletedSuccessCacheHit`. If no target index is requested, cache hits are strictly restricted to `checkTtlCacheHit(db.Path)`. Once 5 seconds elapse, cache expires and fresh runs are queried from GitHub.
   - In `cli/cmdpipeline/pipeline_status.go`: Removed `checkCommitMatchCacheHit` bypass; status cache strictly honors 5s TTL.
2. **True Chronological Sorting in SQLite DB:**
   - In `cli/pipelinedb/pipeline_split_ops.go` and `cli/pipelinedb/pipeline_prune.go`: Updated all queries to `ORDER BY CreatedAt DESC, RunId DESC` instead of `ORDER BY PipelineRunId DESC`.
3. **Defensive Group Timestamp Tracking & In-Memory Payload Propagation:**
   - In `cli/cmdpipeline/pipeline_commit_groups.go`: `addRun` now tracks the newest `CreatedAt` across workflows in each group, and `build()` defensively sorts groups descending by `CreatedAt`.
   - In `cli/cmdpipeline/pipeline.go` and `cli/cmdpipeline/pipeline_logs.go`: Added `Runs []ghRunItem` to `PipelineErrorLogsPayload`. Terminal renderers directly consume the in-memory runs evaluated during the request, avoiding inconsistent second-pass DB reads.
   - In `cli/cmdpipeline/pipeline_query.go`: `resolveCachedRunsOrFetch` verifies TTL freshness before serving cached runs.
4. **Formatting Fix:**
   - Formatted `cli/cmd/llm/llm_train_test.go` with `gofmt`.

---

## 4. Verification

- `golangci-lint run ./...`: 0 errors.
- `python linter-scripts/check-nested-ifs.py`: 0 violations across all candidate files.
- `python linter-scripts/check-boolean-guidelines.py`: 0 violations.
- `python linter-scripts/check-newline-styling.py`: 100% Unix LF.
- `go test ./cmdpipeline`: PASS (100% green).
- `go test ./pipelinedb`: PASS (100% green).
- Ran `gitmap pe` from home directory:
  - Output accurately identifies `2764db6` as latest commit and displays `latest 2764db6` in the summary table.
  - Subsequent execution within 5 seconds cleanly hits cache; execution after 5 seconds refreshes from GitHub.
