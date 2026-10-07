# Architecture Specification: GitMap PE Cache Invalidation, Relative Paths & Error Extraction Remediation

> **Specification ID:** 238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation  
> **Status:** Ratified  
> **Version:** 1.0.0  
> **Target Subsystem:** `cli/cmdpipeline/`  

---

## 1. Executive Summary

This architecture specification addresses three systemic defects identified in the `gitmap pe` (pipeline error inspection) command:
1. **Stale In-Progress Cache Evaluation**: In `cli/cmdpipeline/pipeline_cache_eval.go`, `checkHeadShaCacheHit`, `checkCommitTargetCacheHit`, and `checkTtlCacheHit` returned a cache hit if any single run was completed or if TTL was active, even when target workflow runs were actively `in_progress`. Consequently, `gitmap pe` served stale cached snapshots, locking the terminal into false `Active Pipeline is RUNNING (ETA: overtime +Xm)` banners without ever querying live GitHub Actions status.
2. **Explicit Target Bypassing in Cache**: When a specific commit SHA or offset was targeted (e.g. `gitmap pe 159351e` or `gitmap pe -6`), if the target was absent from the local cache, the evaluator fell through to `checkCompletedCommitCacheHit` or `checkTtlCacheHit` for the repository's latest HEAD SHA (`latest.Sha`), returning cached data for the wrong commit instead of executing a live GitHub API fetch for the requested target.
3. **Absolute Path Leakage in Error Output**: In `cli/cmdpipeline/pipeline_logs.go:1205`, `cli/cmdpipeline/pipeline_logs.go:1319`, `cli/cmdpipeline/pipeline_details.go:532`, and `cli/cmdpipeline/pipeline_error_extract.go:896`, log paths under user app data directories (e.g. `C:/Users/.../AppData/Local/gitmap-cli/data/pipeline/...`) were printed directly without transformation, violating relative path guidelines and exposing internal machine user profiles in terminal and clipboard buffers.
4. **Redundant Tab Prefix Artifacts in Error Summaries**: When GitHub Actions logs with `job\tstep\t` prefixes were formatted in section summaries, the raw tab-delimited string was emitted verbatim next to formatted `Job: ... | Step: ...` labels, resulting in redundant display noise.

---

## 2. Architectural Invariants

1. **In-Progress Run Invalidation Invariant**:
   - A cached commit evaluation MUST NEVER be considered a cache hit if ANY run for that commit is active (`isRunActive(r)`: `status == "in_progress"` or `status == "queued"`).
   - If an active run exists for the target commit in local SQLite DB records, `EvaluatePipelineErrorsCache` MUST return `IsFromCache: false, Reason: "active_run_in_progress"` to force live GitHub telemetry polling.
2. **Commit Target Cache Affinity Invariant**:
   - When an explicit commit SHA or negative offset index is provided in `flags.CommitTarget` or `flags.Index`, the cache decision MUST strictly match that target.
   - If the requested target is not present in cached records, the evaluator MUST return `IsFromCache: false, Reason: "target_not_in_cache"`. It must NEVER fall back to serving cached records for the repository HEAD commit.
3. **Complete Commit Runs Invariant**:
   - `checkHeadShaCacheHit` and `checkCommitTargetCacheHit` MUST enforce `isCommitRunsCompleted(dbRuns, sha)`, verifying that ALL runs for the specified SHA are completed before approving a cache hit.
4. **Relative Path Sanitization Invariant**:
   - All log paths emitted to stdout, stderr, or copied to the clipboard MUST pass through `FormatRelativeDbPath` / `toRelativeGitPath` to convert `%LOCALAPPDATA%` and home directory paths into clean relative format (`data/pipeline/<slug>/<run-id>.log`).
