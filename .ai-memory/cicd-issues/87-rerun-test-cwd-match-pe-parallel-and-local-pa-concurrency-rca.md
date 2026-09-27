# CI/CD RCA: Rerun Test CWD Match, PE Parallel Log Fetching & Local PA Concurrency Restoration

- **Issue ID:** CI-87
- **Date:** 2026-09-27
- **Branch:** `main`
- **Failed Commit:** `bc56920aed75304af35e8971ff1e752261c185e4`
- **Failed Runs:**
  - CI (#36306818978) — Full Suite Guard (`go test ./...`)
  - Cross-Platform Build (#36306818855) — Ubuntu, Windows, macOS (`go test ./...`)
- **Status:** Resolved / Fixed

---

## 1. Reproduction & Error Telemetry

### Exact Failure Log:
```text
● Job: ubuntu-latest / go build + test | Step: go test ./... (no cache)
  Summary: agy_rerun_test.go:87: expected target '1' to match CWD project 'match-1', got "other-1" seq 1
  --- FAIL: TestResolveClosestActiveProject_CwdMatch (0.00s)
      agy_rerun_test.go:87: expected target '1' to match CWD project 'match-1', got "other-1" seq 1
      FAIL
  FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdagy	0.199s
```

### User Feedback & Performance Telemetry:
> "Okay, serious issue with GitMap. When I do the GitMap PE, it's taking insanely long amount of time. Probably you are not doing the parallel execution to get these logs, and you should, actually. If you are not, you have changed something recently based on the error that we faced. Please don't do that. So some positions are meant to run on parallelly. For example, GitMap PA. If it's running on the current machine, it will run on full-fledged as previous. Okay? But only thing, if it is running on SSH, then it will run differently. Do you understand my point? So if you have any reduction in the performance, I think you revert that back, okay, and try to fix it."

---

## 2. Root Cause Analysis (4-Part)

### Part 1: Why Did `TestResolveClosestActiveProject_CwdMatch` Fail?
In commit `17d3d58`, a previous author had added `"1"` into `isCwdTarget(target)` in `cli/cmdagy/agy_rerun_project_resolve.go`:
```go
func isCwdTarget(target string) bool {
    clean := strings.TrimSpace(target)
    return clean == "" || clean == "." || clean == "cwd" || clean == "here" || clean == "1"
}
```
And added a test asserting that passing target `"1"` must match CWD project `match-1`:
```go
p1, seq1 := resolveClosestActiveProject(projects, "1")
if p1.ID != "match-1" || seq1 != 2 {
    t.Errorf("expected target '1' to match CWD project 'match-1', got %q seq %d", p1.ID, seq1)
}
```
When we subsequently removed `"1"` from `isCwdTarget` to ensure `gitmap rerun 1` resolves to sequence number 1 (`projects[0]`) rather than hijacking the idle CWD project, `resolveClosestActiveProject(projects, "1")` correctly evaluated sequence number 1 via `strconv.Atoi("1")`, returning `other-1` (seq 1). However, the existing unit test in `agy_rerun_test.go:87` was still asserting the stale behavior where `"1"` was conflated with CWD.

### Part 2: Why Was `gitmap pe` Taking an Insanely Long Amount of Time?
In `cli/cmdpipeline/pipeline_logs.go` and `cli/cmdpipeline/pipeline_failure_tree.go`:
1. Metadata queries (`queryPendingPRs` and `queryLatestTagRelease`) were being executed serially in `enrichErrorLogsMetadata`, blocking on multiple consecutive GitHub CLI API requests.
2. For failing runs, `queryFailedRunLogs` (`gh run view --log-failed`) and `CorrelateRunFailedJobs` (`gh run view --json jobs`) were executed serially within each worker.
3. In `appendCommitGroupFailureTree` and `renderFailingWorkflowsDiagnostics`, failed workflow logs and job lists were queried redundantly and sequentially for every failing workflow in a commit group.

### Part 3: Why Was Local `gitmap pa` Artificially Throttled?
In `cli/cmdssh/ssh_pull_fleet.go`, `executeLocalVMPull` was passing hardcoded `--parallel 2`:
```go
subArgs := []string{"pa", "--json", "--parallel", "2"}
```
This forced the local machine to pull all repositories using only 2 workers instead of full-fledged machine concurrency (`runtime.NumCPU()`).

### Part 4: Root Cause Summary in One Sentence
Stale unit test assertions in `agy_rerun_test.go` conflated target `"1"` with CWD, while `gitmap pe` suffered from unparallelized GitHub API invocations and `executeLocalVMPull` improperly restricted local pull concurrency.

---

## 3. Code Fix & Remediation Plan

1. **Test Update in `cli/cmdagy/agy_rerun_test.go`:**
   - Update `TestResolveClosestActiveProject_CwdMatch` to verify that target `"1"` resolves to sequence #1 (`other-1`), while target `""`, `"."`, or `"current"` resolves to the CWD project (`match-1`, seq 2).
2. **Parallelize `gitmap pe` Metadata & Log Fetching:**
   - In `enrichErrorLogsMetadata` (`cli/cmdpipeline/pipeline_logs.go`), fetch `queryPendingPRs` and `queryLatestTagRelease` concurrently using `sync.WaitGroup`.
   - In `fetchAndBuildFailedRunItem` (`cli/cmdpipeline/pipeline_logs.go`), fetch `queryFailedRunLogs` and `queryRunJobs` in parallel for each failing run.
   - In `BuildCommitGroupFailureTree` (`cli/cmdpipeline/pipeline_failure_tree.go`), parallelize workflow log and job extraction across failing workflows.
3. **Restore Full-Fledged Local Concurrency in `executeLocalVMPull`:**
   - Remove `--parallel 2` from `executeLocalVMPull` in `cli/cmdssh/ssh_pull_fleet.go`, allowing local VM execution to run with full CPU capacity (`runtime.NumCPU()`).
   - Retain `--parallel 2` exclusively on remote SSH delegated commands (`remotePullCmd`) and incoming SSH sessions (`IsSSHSession()`).

---

## 4. Prevention & Quality Verification

- Unit tests in `cmdagy` will be verified against both sequence resolution and CWD targeting.
- Guideline autofixers (`python 03-ai-scripts/05-guideline-autofixer.py`) must pass with 0 errors.
- Version bump to `v6.355.1` with atomic commit and push to trigger clean green CI/CD.
