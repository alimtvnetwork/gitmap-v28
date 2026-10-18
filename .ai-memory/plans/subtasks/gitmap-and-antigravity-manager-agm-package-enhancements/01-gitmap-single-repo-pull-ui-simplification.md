# Subtask 01: GitMap Single-Repo Pull UI Simplification

**Parent Task:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Subtask ID:** `subtask-01`  
**Status:** `IN_PROGRESS`  
**Owned Files:**  
- `cli/cmdpull/pull_cwd.go`
- `cli/cmdpull/pull_worker.go`
- `cli/cloner/safe_pull.go`
- `cli/cmdpull/pull_summary_single.go`

## Objective
When running `gitmap pull` in a single git repository (CWD), replace the noisy 0%, 20%, 30% progress bar lines and the 8-column batch table with:
1. A simple, clean in-progress message: `  ↓ pulling remote changes...` (cleared on completion in TTY mode).
2. A clean, compact completion summary:
   - Commit delta (`from..to` or `up-to-date`)
   - Changed files and line diff count if updated (`+X/-Y, N files`)
   - Clean / active / dirty state (`[clean]` vs `[dirty: N uncommitted files]`)
   - Auto-merge indicator if executed (`[auto-merge executed]`)

## Implementation Details
1. In `cli/cloner/safe_pull.go`:
   - When `attemptAutoMergePull` succeeds, set `result.IsAutoMerged = true` or `state.IsAutoMerged = true`.
2. In `cli/cmdpull/pull_worker.go`:
   - Update `PullRepoState` to carry `IsAutoMerged bool`.
   - Propagate auto-merge status from `cloner.CloneResult`.
3. In `cli/cmdpull/pull_cwd.go`:
   - In `executeCWDTrackedPull(cwd)`:
     - DO NOT instantiate `PullProgressBar(1)`. Pass `bar = nil` to `ExecuteTrackedPull(rec, nil)`.
     - Print `  ↓ pulling remote changes...`
     - Call `RenderSingleRepoPullSummary(state, diag)` instead of `RenderPullBatchTable([]model.PullTableRow{row})`.
4. In `cli/cmdpull/pull_summary_single.go`:
   - Implement `RenderSingleRepoPullSummary(state *PullRepoState, cwd string)` rendering the clean, modern status line.
