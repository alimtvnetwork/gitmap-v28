# CI/CD RCA 117: Unused `gitmap pe all` Cache and Error Symbols in `cli/cmdpipeline/`

## 1. Issue Summary & Failure Telemetry
- **Failed Workflow Runs:**
  - CI: Run `#38154388787` (`Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)` and `Lint Baseline Diff`)
- **Extracted Error Trace:**
  ```text
  X [unused] func `collectAllPipelineSummaryWithWorkers` is unused (NEW vs baseline)
  Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive): cmdpipeline/pipeline_all_errors.go#124

  X [unused] func `recordCommitSha` is unused (NEW vs baseline)
  Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive): cmdpipeline/pipeline_all_cache.go#214

  X [unused] func `getCachedCommitSha` is unused (NEW vs baseline)
  Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive): cmdpipeline/pipeline_all_cache.go#204
  ```

---

## 2. Root Cause Analysis (4-Part RCA)

### A. Immediate Cause
The CI workflow job `Lint Baseline Guard` runs `golangci-lint run --disable-all -E unused` and diffs findings against the rolling baseline cache (`golangci-unused-baseline-main-v4`). Three newly added private functions in `cli/cmdpipeline/` were not referenced by any callers or tests in the codebase:
1. `collectAllPipelineSummaryWithWorkers` in `cli/cmdpipeline/pipeline_all_errors.go`
2. `getCachedCommitSha` in `cli/cmdpipeline/pipeline_all_cache.go`
3. `recordCommitSha` in `cli/cmdpipeline/pipeline_all_cache.go`

### B. Underlying Mechanism & Propagation
During the refactoring and optimization of `gitmap pe all` cache in commit `d6c9884a2`, `getCachedInspection` and `recordCachedInspection` were implemented to handle full structured cache records (including SHA, clean state, failed items, and step jobs). The simpler helper functions `getCachedCommitSha` and `recordCommitSha` were made redundant and left unused in `pipeline_all_cache.go`. Similarly, `collectAllPipelineSummaryWithWorkers` was declared as a worker-configurable wrapper around `collectAllPipelineSummaryWithOptions`, but only `collectAllPipelineSummary` (default workers) and direct `collectAllPipelineSummaryWithOptions` were actually called.

### C. Detection Latency & Tooling Defect
`golangci-lint run --disable-all -E unused ./cmdpipeline/...` was not executed locally in `cli/` prior to pushing commit `d6c9884a2`.

### D. Prevention & Mitigation
- Eliminate dead and redundant helpers when refactoring caching logic.
- Run `golangci-lint run --disable-all -E unused ./cmdpipeline/...` locally before committing changes to `cli/cmdpipeline/`.

---

## 3. Remediation & Code Fix
1. Removed unused `collectAllPipelineSummaryWithWorkers` from `cli/cmdpipeline/pipeline_all_errors.go`.
2. Removed unused `getCachedCommitSha` and `recordCommitSha` from `cli/cmdpipeline/pipeline_all_cache.go`.

---

## 4. Verification & Prevention
- Verified locally with `golangci-lint run --disable-all -E unused ./cmdpipeline/...`: exits 0 with 0 findings.
- Verified with `go vet ./...` across the entire `cli` module: exits 0 with zero warnings.
- Verified guideline autofixer, relative paths, newline styling, and boolean/enum linters: all pass.
- Executed minor release bump per continuous self-healing instructions.
