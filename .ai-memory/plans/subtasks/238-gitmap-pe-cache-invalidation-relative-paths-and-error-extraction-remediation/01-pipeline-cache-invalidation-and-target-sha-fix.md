# Subtask 01: Pipeline Cache Invalidation & Target SHA Fix

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation.md`  
> **Owned Files:**  
> - `cli/cmdpipeline/pipeline_cache_eval.go`  
> - `cli/cmdpipeline/pipeline_cache_eval_test.go`  

---

- [x] 1. In `cli/cmdpipeline/pipeline_cache_eval.go`:
   - Enforce target affinity in `evaluateDecisionFromRuns`: If `flags.CommitTarget != ""` or `flags.HasIndex`, return `IsFromCache: false` if not matched.
   - Enforce `isCommitRunsCompleted`: A commit is only completed if ALL runs for that commit are completed.
   - Disallow cache hit if any run for the target commit has `Status == "in_progress"` or `Status == "queued"`.
- [x] 2. In `cli/cmdpipeline/pipeline_cache_eval_test.go`:
   - Add unit tests verifying cache rejection when runs are in-progress or when target commit is not in cache.
- [x] 3. Verify with `go test -v ./cli/cmdpipeline/...`.

---

## 2. Verification Evidence

- `TestHasAnyActiveRun`: PASS.
- `TestCheckCommitTargetCacheHit_PartialCompletedRejection`: PASS.
- `TestEvaluateDecisionFromRuns_TargetNotInCache`: PASS.
- `TestEvaluateDecisionFromRuns_ActiveRunInProgress`: PASS.
- All cache evaluation tests pass in 0.161s.
- Status: **DONE**

