# Subtask 63.2: CLI Routing & Conflict Resolution

- **Parent Plan:** [63-semantic-flat-commit-suite.md](../../completed/63-semantic-flat-commit-suite.md)
- **Spec Reference:** [02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md](../../../../02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md)
- **Status:** Completed
- **Files:** `cli/cmd/rootcore.go`, `cli/cmd/roottooling.go`, `cli/constants/constants_cli.go`

## Objective
Eliminate shadow routing where legacy passthrough in `rootcore.go` intercepted `"commit"` and `"cm"` before tooling dispatchers could execute.

## Implementation Details
1. Defined constants: `CmdCommit = "commit"`, `CmdCommitAlias = "cm"`, `CmdCommitAlias2 = "commit-all"`, `CmdCommitAlias3 = "ca"`.
2. Updated `rootcore.go` to route all four aliases directly to `runCommit(argsTail())`.
3. Removed duplicate routing entry from `roottooling.go`.
4. Removed deprecated `cli/cmd/commit_cli.go`.
