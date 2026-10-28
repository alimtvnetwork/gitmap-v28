# Subtask 01: Casing Persistence and Auto-Commit
Parent Plan: 338-gitmap-pull-errors-and-ui-enhancement

## Scope
- In `cli/cmdfix/fix_casing.go`, when `ExecuteCasingNormalization` finishes and stages changes via `stageGitAddAll(repoPath)`, commit the staged normalization atomically.
- Prevent infinite re-stashing and re-normalizing cycles during batch pull operations.
- Verify with unit tests.
