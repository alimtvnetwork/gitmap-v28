# Subtask 05: Fleet Nodes SSH Parity & Peer GitMap Delegation

- **Task-ID**: Task-05
- **Subtask**: 05-nodes-remote-and-peer-delegation
- **Owner**: Worker 02
- **Status**: PENDING

## Owned Files (Exclusive)
- `cli/cmdnodes/nodes_remote.go`
- `cli/cluster/exec_delegate.go`
- `cli/cmdssh/ssh_exec_command.go`
- `cli/cmdnodes/nodes_cmd.go`
- `cli/cmdnodes/nodes_remote_test.go`

## Read-Only Context Files
- `02-spec/21-app/goemailsender-connection-and-testing/02-component-spec.md`
- `cli/cmdnodes/nodes_tasks.go`
- `cli/cmdnodes/nodes_filter.go`
- `cli/cluster/agent.go`

## Objectives
1. Fix command classification in `cli/cmdssh/ssh_exec_command.go`: recognize `space`, `heatmap`, and `agent` as GitMap subcommands.
2. Implement `RunNodesRemoteCLI` in `cli/cmdnodes/nodes_remote.go` handling `gitmap nodes run [target] "<command>" [flags]` with multi-node concurrent SSH dispatch.
3. Implement `cli/cluster/exec_delegate.go` defining `PeerDelegationPayload` and `PeerDelegationResponse` for peer GitMap delegation.
4. Add unit tests in `cli/cmdnodes/nodes_remote_test.go` with mock executor verifying command dispatch without live SSH.
