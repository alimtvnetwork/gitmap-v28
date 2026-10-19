# Subtask 01: Node Fleet Telemetry, Account Mapping & Web UI SSH Ping

**Parent Plan:** `.ai-memory/plans/gitmap-node-and-prompt-system-issues.md`  
**Subtask ID:** `subtask-01`  
**Worker:** Worker 01  
**Status:** In Progress  

## Scope of Ownership
- `cli/cmdnodes/nodes_summary.go`
- `cli/cmdnodes/nodes_render.go`
- `cli/cmdnodes/nodes_agy_prompt.go`
- `cli/cmdui/ui_types.go`
- `cli/cmdui/ui_server_ssh.go`
- `cli/cmdui/ui_assets_js_nodes.go`
- `cli/cmdui/ui_assets_markup_ops.go`

## Implementation Checklist
- [ ] In `cli/cmdnodes/nodes_summary.go` and `cli/cmdnodes/nodes_render.go`, add `Account` and `Instances` fields to `UnifiedFleetNode` and `NodeSummaryResult`.
- [ ] Display GitHub node name / alias, bound account, and active instance count in `gitmap nodes` and `gitmap ssh` table output.
- [ ] In `cli/cmdnodes/nodes_agy_prompt.go`, print exact target node, account, and instance being invoked during dispatch.
- [ ] In `cli/cmdui/ui_types.go`, extend `NodeSummary` with `Account`, `Instances`, and `InstanceList`.
- [ ] In `cli/cmdui/ui_server_ssh.go`, populate `Account` and `Instances` and implement `/api/ssh/test` ping endpoint.
- [ ] In `cli/cmdui/ui_assets_js_nodes.go`, implement missing `testNode(alias)` function to call `/api/ssh/test`, update the badge dynamically, and show toast.
- [ ] In `cli/cmdui/ui_assets_markup_ops.go`, update Nodes table headers and columns to include Account and Instances.
- [ ] Verify with targeted test: `go test ./cmdnodes -run TestUnifiedNodes` and `go test ./cmdui -run TestAPISSHNodes`.
