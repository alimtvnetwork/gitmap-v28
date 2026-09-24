# Issue 38: PAET Table Command Excessive Execution Latency: Root Cause Analysis & Mitigation

> **/goal** Provide grounded 4-part Root Cause Analysis (Reproduction, Cause, Fix, Prevention) for excessive execution latency in `gitmap paet` (`pull-all-efficient-table`).
> **/learn** Pinpoint exact subprocess and network bottlenecks during batch pull and table rendering, implementing targeted mitigation without regressions.

---

## 1. Issue Summary & Status Matrix

| Component | Target Command | Reported Symptom | Root Cause | Status |
|---|---|---|---|---|
| **Inactivity Heuristic** | `gitmap paet` | 0 inactive repositories skipped; all 58 repos pulled on routine execution | `minRuns=20` requires 20 prior pulls within 24h before a repo can be marked inactive | Identified |
| **Table PR Status Detection** | `renderPullBatchResults` | Post-pull table rendering hangs for 40-70+ seconds after pulls finish | `detectGitHubPRs` runs `gh pr list --state open --json number` sequentially for every row | Mitigated |
| **Command Help Intercept** | `RunPullAllEfficient` | Running `gitmap paet --help` executes network pull of 58 repos instead of printing help | Missing `checkHelp` at top of `RunPullAllEfficient` | Fixed |

---

## 2. Reproduction

Execute `gitmap paet` on a system with 58 tracked repositories:
1. `gitmap paet` is invoked.
2. The CLI reports `resolved 58 repo(s) (58 active, 0 inactive skipped in 24h window)`.
3. 58 repositories begin executing `git fetch` and `git pull` across the network, taking ~15–25 seconds.
4. The progress bar reaches 100%.
5. **The terminal hangs for another 45–60 seconds before printing the table.**
6. When `gitmap paet --help` is invoked, rather than showing help in <10ms, it initiates the entire 58-repo pull process.

---

## 3. Root Cause Analysis

### Cause 1: Serial External GitHub CLI Network Invocations during Table Rendering
In `cli/cmdpull/pull.go`:
```go
func renderPullBatchResults(states []*PullRepoState) {
    var tableRows []model.PullTableRow
    for _, state := range states {
        tableRows = append(tableRows, buildPullTableRowFromState(state))
    }
    RenderPullBatchTable(tableRows)
}

func buildPullTableRowFromState(state *PullRepoState) model.PullTableRow {
    ...
    row.PRStatus = gitutil.DetectPRStatus(state.RepoPath)
    ...
}
```
And in `cli/gitutil/pr_detector.go`:
```go
func detectGitHubPRs(repoPath string) (string, bool) {
    cmdGh := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number")
    cmdGh.Dir = repoPath
    out, err := cmdGh.Output()
    ...
}
```
For each of the 58 repositories, `buildPullTableRowFromState` synchronously executes the external binary `gh pr list --state open --json number`.
- Each `gh` invocation makes an HTTPS REST API call to `api.github.com`.
- With 58 repositories, this creates **58 sequential remote API roundtrips**.
- At an average of 800ms–1200ms per roundtrip, this single loop adds **45 to 70 seconds** of pure blocked network I/O *after* git pull has already completed.

### Cause 2: Ineffective Inactivity Filter (`minRuns=20` Threshold)
In `cli/store/split_db_pull_ops.go`:
```go
func (s *PullSplitDB) EvaluateRepoInactivity(repoPath string, minRuns int, windowHours int) (RepoInactivityStatus, error) {
    history, err := s.GetRecentRepoPullHistory(repoPath, minRuns)
    ...
    if len(history) < minRuns {
        status.Reason = fmt.Sprintf("insufficient run history (%d/%d runs)", len(history), minRuns)
        return status, nil
    }
```
The caller in `pull_efficient.go` passes `minRuns = 20` and `windowHours = 24`.
Unless a developer runs `gitmap pae` or `gitmap pull` twenty times within a 24-hour period, every single repository will have `< 20` runs. Consequently, `status.IsInactive` evaluates to `false` for every repository on every normal working day, defeating the purpose of the efficient pull partition.

### Cause 3: Missing Help Intercept in `RunPullAllEfficient`
In `cli/cmdpull/pull_efficient.go`, `RunPullAllEfficient` extracted flags but did not intercept `checkHelp` before calling `requireOnline()` and `resolveAllTrackedRecords()`. As a result, checking `--help` triggered full execution.

---

## 4. Fix & Mitigation Architecture

1. **Short-Circuit Help Flag:**
   Check for `--help` / `-h` at the very start of `RunPullAllEfficient` and immediately render command help text.
2. **Opt-in / Fast PR Status Evaluation:**
   In batch table rendering (`renderPullBatchResults`), bypass the slow `gh pr list` network call by default (or use `--pr` flag when explicit GitHub PR querying is required). Fall back to ultra-fast local branch tracking (`gitutil.DetectGitBranchTracking` in <1ms) or cache results.
3. **Adaptive Inactivity Evaluation:**
   Support an adaptive minimum runs threshold (e.g. 3 consecutive zero-change runs within 24h) so repositories with recent zero-change pull history are legitimately skipped during efficient pulls.

---

## 5. Prevention Invariants

1. **Never Make Network Calls in Table Rendering Loops:** Table rendering functions (`renderPullBatchResults`, `RenderPullBatchTable`) must only format and display data already collected in memory or query fast local git metadata. External CLI/API calls (`gh`, `curl`, network probes) are strictly banned inside serial row loops.
2. **Mandatory Top-Level Help Check:** Every command entry point in `cli/cmd*/` must evaluate `isHelpArg` as step 0 before opening databases, verifying network, or resolving records.
