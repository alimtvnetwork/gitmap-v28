# Component Specification: GitMap PE Cache Invalidation, Relative Paths & Error Extraction Remediation

> **Specification ID:** 238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation  
> **Status:** Ratified  
> **Version:** 1.0.0  
> **Target Subsystem:** `cli/cmdpipeline/`  

---

## 1. Component Changes

### 1.1 `cli/cmdpipeline/pipeline_cache_eval.go`
1. **Target Affinity Guard in `evaluateDecisionFromRuns`**:
   ```go
   func evaluateDecisionFromRuns(db *pipelinedb.PipelineSplitDb, repo string, dbRuns []pipelinedb.PipelineRunRecord, flags PipelineErrorFlags) PipelineCacheDecision {
       latest := dbRuns[0]
       if flags.CommitTarget != "" || flags.HasIndex {
           if hit, decision := evaluateTargetOrHeadCacheHit(repo, dbRuns, flags, latest.Sha); hit {
               return decision
           }
           return PipelineCacheDecision{IsFromCache: false, Reason: "target_not_in_cache"}
       }
       if hasAnyActiveRun(dbRuns, latest.Sha) {
           return PipelineCacheDecision{IsFromCache: false, Reason: "active_run_in_progress"}
       }
       if checkCompletedCommitCacheHit(repo, dbRuns) {
           return buildCacheHitDecision(dbRuns, latest.Sha, "latest_commit_completed")
       }
       if checkTtlCacheHit(db.Path) && !hasAnyActiveRun(dbRuns, latest.Sha) {
           return buildCacheHitDecision(dbRuns, latest.Sha, "within_ttl")
       }
       return PipelineCacheDecision{IsFromCache: false, Reason: "cache_expired"}
   }
   ```
2. **Strict Complete Runs Checks**:
   - Refactor `checkHeadShaCacheHit` and `checkCommitTargetCacheHit` to use `isCommitRunsCompleted`:
   ```go
   func checkCommitTargetCacheHit(dbRuns []pipelinedb.PipelineRunRecord, target string) bool {
       if len(target) == 0 {
           return false
       }
       matched := false
       for _, r := range dbRuns {
           if matchCommitSha(r.Sha, target) {
               matched = true
               break
           }
       }
       return matched && isCommitRunsCompleted(dbRuns, target)
   }
   ```

### 1.2 `cli/cmdpipeline/pipeline_logs.go`
1. **Relative Path Transformation**:
   - In `renderSingleFailureTerminal`:
     `fmt.Printf("      Log:     %s\n", filepath.ToSlash(FormatRelativeDbPath(sec.SavedLogFile)))`
   - In `renderFailedRunCard`:
     `fmt.Printf("  │ Saved Log: %s\n", filepath.ToSlash(FormatRelativeDbPath(fr.SavedLogFile)))`
2. **Negative Offset Target Resolution**:
   - In `queryWorkflowRunsForTarget`: If `target` is a negative index like `-5` or `-6` and not present in initial 30 runs, resolve the commit SHA via `git rev-parse HEAD~<N>` and query `queryWorkflowRunsByCommit(repo, resolvedSha)`.

### 1.3 `cli/cmdpipeline/pipeline_details.go`
1. **Relative Path Transformation**:
   - In `renderDetailedSectionFailure`:
     `fmt.Printf("      Log:     %s\n", filepath.ToSlash(FormatRelativeDbPath(sec.SavedLogFile)))`

### 1.4 `cli/cmdpipeline/pipeline_error_extract.go`
1. **Relative Git Path Sanitization**:
   - In `toRelativeGitPath(targetPath string) string`:
     Check `FormatRelativeDbPath(targetPath)` if outside repo root to convert `AppData/Local/.../data/pipeline/...` to `data/pipeline/...`.
2. **Clean Redundant Tab Prefixes**:
   - In `cleanLogText` or `formatSingleSectionFailure`: Strip redundant `job + "\t" + step + "\t"` if present at the beginning of error summaries.
