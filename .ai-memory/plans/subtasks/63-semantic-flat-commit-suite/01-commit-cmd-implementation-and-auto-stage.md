# Subtask 63.1: Commit Command Implementation & Auto-Stage Engine

- **Parent Plan:** [63-semantic-flat-commit-suite.md](../../completed/63-semantic-flat-commit-suite.md)
- **Spec Reference:** [02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md](../../../../02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md)
- **Status:** Completed
- **File:** `cli/cmd/commit_cmd.go`

## Objective
Implement `runCommit` in `cli/cmd/commit_cmd.go` to provide 1-step auto-staging (`git add -A`), message concatenation, dry-run previews, and optional remote push (`--push` / `-p`) while keeping all functions strictly <= 15 lines.

## Implementation Details
1. `parseCommitFlags(args)` separates clean commit message tokens from `--push`, `-p`, `--dry-run`, `-n`.
2. `previewCommitChanges()` executes `git status --short` when `--dry-run` is supplied.
3. `executeCommit(args, hasPush)` executes `git add -A` and invokes `dispatchGitCommit(args)`.
4. `handleOptionalPush(hasPush)` pushes to remote if `--push` was specified.
5. `dispatchGitCommit(args)` runs `git commit -m "<msg>"` or launches the commit editor if no message is provided.
