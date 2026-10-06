# Issue Post-Mortem: Stale Pull Error Persistence and Clear Subsystem

> **Date:** 2026-10-06
> **Affected Subsystem:** `cli/cmdpull`, `cli/cmdpullerror`, `cli/store` (Pull Split-DB)
> **Repository:** `gitlogger-new-v2` (`d:\work\gitlogger-new`)
> **Cross-Reference:** `02-spec/22-app-issues/70-stale-pull-error-persistence-and-clear-subsystem-rca.md`

## 1. Symptom
`gitmap pull-error gitlogger-new-v2` continued to output stale historical merge conflict diagnostics (`[E_MERGE_CONFLICT:EXECUTION] pull_merge: git merge conflict detected: auto-merge pull aborted (at=cloner/safe_pull.go:28)`) after the repository divergence had been resolved and subsequent pulls succeeded cleanly.

## 2. Root Cause
1. `D:\work\gitlogger-new` local branch `main` had diverged from `origin/main` across 16 files, correctly triggering GitMap's `safe_pull` auto-merge abort and writing failure records to `pull_errors`.
2. GitMap's pull engine recorded errors on failure but lacked eviction logic on success in `cli/cmdpull/pull_db_sync.go` and `cli/cmdpull/pull_efficient.go`, leaving historical errors indefinitely.
3. `cli/cmdpullerror/pull_error_cmd.go` lacked an administrative `clear` command / `--clear` flag to allow users or agents to purge resolved diagnostics.

## 3. Resolution
1. Confirmed `d:\work\gitlogger-new` branch divergence was resolved, committed (`7017153`), documented (`02-pull-merge-conflict-and-divergence-resolution-rca.md`), and pushed to `origin/main`.
2. Added `ClearPullErrors` and `ClearPullErrorsForRepo` methods to `*PullSplitDB` in `cli/store/pull_split_db_errors_clear.go`.
3. Integrated automatic eviction into `cli/cmdpull/pull_db_sync.go` and `cli/cmdpull/pull_efficient.go` so successful pulls (`!isStateFailure`) clear the repository's errors from `pull_errors`.
4. Added `gitmap pull-error clear [target]` and `--clear` / `-c` flag in `cli/cmdpullerror/pull_error_cmd.go` and `cli/cmdpullerror/pull_error_clear.go`.
5. Purged the stale rows from `gitmap-pull.db`, verifying `gitmap pull-error gitlogger-new-v2` returns clean zero-error status.

## 4. Learnings
Telemetry stores must always implement symmetrical lifecycles: recording on failure and eviction on success, paired with manual administrative cleanup commands.
