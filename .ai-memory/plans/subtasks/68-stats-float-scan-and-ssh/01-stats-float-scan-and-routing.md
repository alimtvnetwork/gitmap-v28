# Subtask 68-01: Stats AvgDuration Float64 Scan Fix & SS Routing

**Subtask ID:** 68-01  
**Status:** Pending  
**Spec Reference:** `cli/store/stats.go`, `cli/cmd/stats.go`  
**Assigned Worker:** Worker 01  
**Target Files:**  
- `cli/store/stats.go` [MODIFY]  
- `cli/cmd/stats.go` [MODIFY]  
- `cli/cmd/stats_test.go` [MODIFY]  

---

## 1. Technical Context & Scope
In SQLite, `COALESCE(AVG(DurationMs), 0)` returns a `float64` value (e.g. `26040.96023391813`). Scanning this value directly into an `int64` target causes:
`sql: Scan error on column index 5, name "AvgDuration": converting driver.Value type float64 ("26040.96023391813") to a int64: invalid syntax`
This causes `cmd/stats.go` to enter `handleStatsError` and misdiagnose the issue as `constants.MsgLegacyProjectData`.

Additionally, when a user types `gitmap ss fix-auth t1`, the `ss` alias currently dispatches to `stats`, which ignores the arguments and executes `loadStats`, resulting in the crash above.

## 2. Implementation Directives
1. In `cli/store/stats.go`:
   - In `QueryOverallStats()`:
     Scan column 5 into a `var avgDuration float64`.
     Assign `s.AvgDuration = int64(math.Round(avgDuration))`.
   - In `scanStatsRows()`:
     Scan column 5 into a `var avgDuration float64`.
     Assign `s.AvgDuration = int64(math.Round(avgDuration))`.
2. In `cli/cmd/stats.go`:
   - Implement `isSSHForwardCandidate(args []string) bool`:
     Returns true if `len(args) > 0` and `args[0]` matches `fix-auth`, `login`, `host`, `nodes`, `enable`, `port`, `troubleshoot`, `key`, `copy-id`, `auth-key`.
   - In `runStats(args []string)`:
     If `isSSHForwardCandidate(args)`, return `runSSH(args)`.
3. In `cli/cmd/stats_test.go`:
   - Add unit tests verifying `isSSHForwardCandidate` and float64 scan resilience.
