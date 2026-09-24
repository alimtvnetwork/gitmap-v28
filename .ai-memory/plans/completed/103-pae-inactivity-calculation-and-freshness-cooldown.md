# Plan 103: Inactivity Calculation, Freshness Cooldown, and Global Column Width Stability for `gitmap pae`

## Status: Completed
**Version:** `v6.332.0`  
**Spec Reference:** [Spec 154](../../../02-spec/21-app/154-pae-inactivity-calculation-and-freshness-cooldown.md)  
**Parent Task:** GitMap Pull All-Efficient Optimization

---

## 1. Objectives & Steps
- [x] **Step 1: Store Enhancement in `cli/store/split_db_pull_activity.go`**
  - Added `GetRecentActualRepoPullHistory` filtering out synthetic `skipped-inactive` rows.
  - Implemented `EvaluateRepoActivityStatus(repoPath string, windowHours int, cooldownMinutes int) (RepoInactivityStatus, error)`.
  - Cooldown: 5 minutes (skips redundant rapid successive network pulls).
  - 24h Window: single clean pull marks repo inactive for remainder of day (no 3-run trap).
- [x] **Step 2: Command Integration in `cli/cmdpull/pull_efficient.go`**
  - Updated `classifyRepoActivity` to use `EvaluateRepoActivityStatus`.
  - Passed `records []model.ScanRecord` through lifecycle to enable workspace-wide stable column width calculation.
  - Updated skip summary messages to accurately reflect 0 changes or recent checks within 24h window.
- [x] **Step 3: Render Stability in `cli/cmdpull/pull_efficient_render.go`**
  - Updated `ResolveConciseRepoColWidth(states []*PullRepoState, allRecords ...[]model.ScanRecord) int` to use tracked records when present.
- [x] **Step 4: Unit & E2E Tests**
  - Unit tests in `cli/store/split_db_pull_test.go` verifying cooldown and single-pull 24h window inactivity (100% PASS).
  - Unit tests in `cli/cmdpull/pull_efficient_render_test.go` verifying layout stability with `allRecords` (100% PASS).
  - Isolated temporary E2E test in `cli/tests/e2e/pae_activity_calculation_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`, 100% PASS).
- [ ] **Step 5: Release & Deployment**
  - Bump version to `v6.332.0`.
  - Compile and deploy to `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.
  - Single atomic commit, tag, and push.
