# Subtask 02: GitMap Fix Ignore All (Local & SSH Fleet)

## Objectives
- Implement and verify `gitmap fix-ignore-all` (`fia`) and `gitmap fix-ignores-all-ssh` (`fias`).
- Deduplicate identical ignore patterns preserving comments and order.
- Verify commit state (`git ls-files`): only flag tracked files if present in the current commit state, and offer clean untracking via `git rm --cached`.
- Support `-y` non-interactive mode.
- In interactive mode without `-y`, display clean summary of issues per repository before prompting for confirmation.
- Wire `fias` to follow GitMap PAS Formula (TaskQueue registration, throttled fleet execution).

## Target Files
- `cli/cmdignore/fix_ignore.go`
- `cli/cmdignore/fias.go`
- `cli/cmd/rootcore.go`
