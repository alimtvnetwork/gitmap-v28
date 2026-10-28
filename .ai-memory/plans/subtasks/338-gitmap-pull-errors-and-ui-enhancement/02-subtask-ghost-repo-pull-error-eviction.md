# Subtask 02: Ghost Repo and Pull Error Eviction
Parent Plan: 338-gitmap-pull-errors-and-ui-enhancement

## Scope
- In `cli/cmdpull/pull_missing_remediate.go`:
  - Update `purgeInvalidRepoEntry` to accept `repoName` and invoke `store.OpenPullSplitDB().ClearPullErrorsForRepo(repoName, repoPath)`.
  - When clone fallback cannot proceed due to missing remote URL or invalid target, evict the error entry so it never reappears in subsequent batch runs.
- Verify that `loadRecentPullFailuresFromDB` returns zero ghost entries after remediation.
