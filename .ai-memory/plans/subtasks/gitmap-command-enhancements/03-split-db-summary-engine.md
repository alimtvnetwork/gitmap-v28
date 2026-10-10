# Subtask 03: Split-DB Summary Engine & Heated Files
Parent Task: gitmap-command-enhancements
Status: COMPLETED

## Objective
Implement `gitmap summary [$repo] [N]` (default $N=8$, target current directory or path/URL) with $\le 200$-word gist, top 5 heated churn files, and normalized Split-DB cache (`summary.db`) with view `ViewRepoReleaseSummaries`.

## Target Files
- `cli/cmdsummary/summary_types.go`
- `cli/cmdsummary/summary_cache.go`
- `cli/cmdsummary/summary_core.go`
- `cli/cmdsummary/summary_helpers.go`
- `cli/cmdsummary/summary_test.go`

## Verification
Unit tests pass:
- `TestParseSummaryArgs`
- `TestSynthesizeReleaseGist`
- `TestExtractHeatedFilesFromDiff`
- `TestSplitDBCacheRoundTrip`
Live CLI smoke test: `.\cli\gitmap.exe summary 3` (first run miss, second run 3/3 cache hits in 0.11s).
