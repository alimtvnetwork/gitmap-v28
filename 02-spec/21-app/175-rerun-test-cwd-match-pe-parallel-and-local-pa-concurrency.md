# Specification 175: Rerun Test CWD Match, PE Parallel Log Fetching & Local PA Concurrency

## Status
- **Status:** active
- **Domain:** app / cli / pipeline / agy / concurrency
- **Author:** Lead AI Architect
- **Date:** 2026-09-27

---

## 1. User Request (Verbatim)

```text
FIX it

================================================================================
GITMAP PIPELINE ERROR REPORT
================================================================================
Repo:                 alimtvnetwork/gitmap-v28
Repo URL:             https://github.com/alimtvnetwork/gitmap-v28
Branch:               main
Last Commit:          bc56920
Last Release:         v6.355.0
Open PRs:             0
Status:               completed (conclusion: failure)
Pipeline Run URL:     https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36306818978
================================================================================

FAILED PIPELINE CHECKS TREE:
  Commit bc56920:
    ├── ✖ CI (#36306818978)
    │   └── ✖ Full Suite Guard [Step: Run full test suite]
    └── ✖ Cross-Platform Build (#36306818855)
        ├── ✖ ubuntu-latest / go build + test [Step: go test ./... (no cache)]
        ├── ✖ windows-latest / go build + test [Step: go test ./... (no cache)]
        └── ✖ macos-latest / go build + test [Step: go test ./... (no cache)]

● Combined Pipeline Section Failures [4 failed section(s)]:
  ...
  Summary:   agy_rerun_test.go:87: expected target '1' to match CWD project 'match-1', got "other-1" seq 1
  --- FAIL: TestResolveClosestActiveProject_CwdMatch (0.00s)
  agy_rerun_test.go:87: expected target '1' to match CWD project 'match-1', got "other-1" seq 1
      FAIL
  FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdagy	0.199s

Okay, serious issue with GitMap. When I do the GitMap PE, it's taking insanely long amount of time. Probably you are not doing the parallel execution to get these logs, and you should, actually. If you are not, you have changed something recently based on the error that we faced. Please don't do that. So some positions are meant to run on parallelly. For example, GitMap PA. If it's running on the current machine, it will run on full-fledged as previous. Okay? But only thing, if it is running on SSH, then it will run differently. Do you understand my point? So if you have any reduction in the performance, I think you revert that back, okay, and try to fix it. And try to fix the root cause of the error, right? The root cause, why it happened, how it happened, and try to make sure that the errors do not repeat back again. Okay? And write proper code following the coding guideline. Do you understand? Can you please help me with this?
```

---

## 2. Architectural Blueprint & Requirements

### 2.1 Rerun Sequence Resolution & CWD Test Alignment
- In `cli/cmdagy/agy_rerun_project_resolve.go`, `isCwdTarget` explicitly matches `""`, `"."`, `"current"`, and `"here"`. Target `"1"` represents sequence number 1 and resolves to `projects[0]`.
- In `cli/cmdagy/agy_rerun_test.go`:
  - `TestResolveClosestActiveProject_CwdMatch` must assert that target `""`, `"."`, and `"current"` match the CWD project (`match-1`, seq 2).
  - Target `"1"` must be verified as resolving to sequence 1 (`other-1`, seq 1), proving that explicit sequence requests are never hijacked by CWD.

### 2.2 PE Concurrency & Pipeline Latency Optimization
- **Parallel Metadata Queries:** In `enrichErrorLogsMetadata` (`cli/cmdpipeline/pipeline_logs.go`), fetch `queryPendingPRs` and `queryLatestTagRelease` concurrently via goroutines and `sync.WaitGroup`.
- **Parallel Run Job & Log Extraction:** In `fetchAndBuildFailedRunItem` (`cli/cmdpipeline/pipeline_logs.go`), fetch `queryFailedRunLogs` and `queryRunJobs` concurrently so both network requests execute in parallel.
- **Parallel Failure Tree Extraction:** In `BuildCommitGroupFailureTree` (`cli/cmdpipeline/pipeline_failure_tree.go`), parallelize workflow log and job extraction across failing workflows instead of executing them sequentially in a single-threaded loop.

### 2.3 Restore Full-Fledged Local `gitmap pa` Concurrency
- In `cli/cmdssh/ssh_pull_fleet.go`:
  - `executeLocalVMPull` must run `gitmap pa --json` WITHOUT the forced `--parallel 2` argument.
  - When running locally, `cloneconcurrency.Resolve(0)` uses `runtime.NumCPU()`, utilizing 100% of available CPU cores.
  - Remote node execution (`remotePullCmd`) and incoming SSH sessions (`cloneconcurrency.IsSSHSession()`) retain half-priority throttling (capped at 2 workers) to prevent freezing remote hosts.

---

## 3. Verification Gates & Acceptance Criteria

1. **AC-175-1: Unit Tests Green**: `go test ./cli/cmdagy/...` passes with zero failures. `TestResolveClosestActiveProject_CwdMatch` correctly distinguishes sequence #1 from CWD match.
2. **AC-175-2: PE Parallel Execution**: `gitmap pe` retrieves and correlates pipeline logs, PR counts, and release tags using concurrent goroutines.
3. **AC-175-3: Local PA Full Concurrency**: `executeLocalVMPull` runs without `--parallel 2`, allowing local pulls to utilize full hardware concurrency while SSH pulls remain throttled.
4. **AC-175-4: Coding Guidelines**: All functions <= 8–15 lines, positive booleans only, structured error returns, zero guideline violations.
