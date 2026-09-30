# Plan 201: Fleet Nodes Clone, CFR, CFRP, Remote Manifest Staging, and Fleet Suggestions

## Objective
Implement `gitmap nodes clone`, `gitmap nodes cfr`, and `gitmap nodes cfrp` to execute repository cloning and CFR post-clone setup locally on the current machine and asynchronously across all enrolled SSH fleet nodes in parallel. Automatically stage local files/manifests to remote nodes' default work directories (`D:\work` / `~/work`) before executing remote commands. Provide actionable fleet clone suggestions when users perform `gitmap clone` or `gitmap cfr`. Conduct full release ceremony for `v6.426.0`.

## Architecture & Subtasks

### Subtask 1: `cli/cmdnodes/` Implementation
- `cli/cmdnodes/nodes_clone_types.go`: Definitions for `NodesCloneKind` (`clone`, `cfr`, `cfrp`), options, execution results, and interfaces.
- `cli/cmdnodes/nodes_clone_file.go`: Detection of local manifest/target file (`.json` or regular file), reading file payload, resolving target work directories on remote nodes (`D:\work` for Windows, `~/work` for Unix), and streaming file over SSH via `cmdssh.StreamFileToRemote`.
- `cli/cmdnodes/nodes_clone_remote.go`: Asynchronous parallel worker pool executing `gitmap <kind> <args>` across remote fleet nodes, capturing stdout/stderr, measuring latency, and rendering summary table.
- `cli/cmdnodes/nodes_clone.go`: Master orchestrator coordinating local execution and asynchronous remote fleet execution with recursion suppression.
- `cli/cmdnodes/nodes_clone_help.go`: Rich, structured help output with explicit examples for single repo, multi repo, and local manifest staging.

### Subtask 2: CLI Routing & Help Integration
- In `cli/cmd/nodes_cmd.go`: Wire `clone`, `cfr`, `cfrp` subcommand detection into `runUnifiedNodesCLI(args)`.
- Update `printUnifiedNodesHelp()` to document `gitmap nodes clone`, `gitmap nodes cfr`, `gitmap nodes cfrp`.
- In `cli/cmd/rootcore.go`: Register aliases `nodes-clone`, `nodes-cfr`, `nodes-cfrp`, `node clone`, `node cfr`, `node cfrp`.
- In `cli/cmd/nodes_cmd.go`: Update Supported Commands & Clusters Matrix to feature fleet clone commands.

### Subtask 3: Fleet Clone Suggestions in `cmdclone`
- In `cli/cmdclone/clone.go`: Add fleet suggestion at conclusion of direct `gitmap clone`.
- In `cli/cmdclone/clonefixrepo.go`: Add fleet suggestion at conclusion of direct `gitmap cfr` and `gitmap cfrp`.
- Add `SetFleetCloneActive` / `IsFleetCloneActive` to suppress redundant nested suggestions when invoked via `nodes clone`.

### Subtask 4: Unit Testing & Linter Verification
- Author `cli/cmdnodes/nodes_clone_test.go` verifying command parsing, file detection, remote command generation, and help text.
- Run `python .github/scripts/go-format-check.py --check-only`.
- Run `python .github/scripts/file-size-check.py`.
- Run `python .github/scripts/check-constants-collisions.py`.

### Subtask 5: Minor Version Bump (`v6.426.0`) & Release Ceremony
- Bump version to `v6.426.0` in `package.json`, `version.json`, `cli/constants/constants.go`.
- Update `changelog.md` with detailed entry for `v6.426.0`.
- Commit, tag, push to `main` and `release/v6.426.0`, and verify all CI/CD pipelines succeed with 100% green status.
