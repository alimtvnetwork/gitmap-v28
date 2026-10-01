# Subtask 03: GitMap Commit-Push All Repos Suite (CPAR)

## Objectives
- Implement and verify `gitmap commit-push-all-repos` (`cpar`).
- Discover all repositories with uncommitted modifications, additions, or deletions.
- Support `-y` non-interactive mode.
- Support `--review` (`-r`) review mode displaying file change details.
- Support `--commit-only` (`--co`) mode: stage and commit locally across repositories without pushing upstream.
- Enqueue commit and push actions into `TaskQueue`.

## Target Files
- `cli/cmdcpar/cpar.go`
- `cli/cmdcpar/types.go`
- `cli/cmd/rootcore.go`
