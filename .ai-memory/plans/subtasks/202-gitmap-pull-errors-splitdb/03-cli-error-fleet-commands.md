---
title: "Fleet Error CLI Commands Routing"
status: "pending"
---

# Objective
Implement the fleet-aware CLI routing for error querying and management in cli/cmd/nodes_cmd.go and ensure cli/cmdsee/see.go routes correctly to the new architecture.

# Implementation Steps
1. **Nodes CLI Dispatcher:**
   - In cli/cmd/nodes_cmd.go (and/or cli/cmd/rootcore.go), register routing handlers for rrors, pull errors, and rrors clear under the 
odes subcommand.
   - Example: gitmap nodes errors -> cmdssh.RunFleetPASCommand("errors", "gitmap errors", ...)
2. **Command Handlers:**
   - Create unNodesErrorsCLI(args) which wraps cmderrors.RunErrorsCLI with PAS dispatch logic.
   - Create unNodesPullErrorsCLI(args) which wraps cmdpullerror.RunPullErrorCLI with PAS dispatch logic.
   - Wire gitmap nodes errors clear to propagate the clear signal to all nodes.
3. **JSON Support:**
   - Ensure --json is propagated properly to the remote nodes.
   - The PAS framework should automatically aggregate JSON outputs if structured properly, or the local CLI should stitch the remote JSON responses together.
4. **See Alias Integration:**
   - Update cli/cmdsee/see.go to ensure gitmap see errors perfectly consumes the new local multi-db stitched output.

# Validation
- gitmap nodes errors --json returns a valid merged JSON array of errors from all remote nodes.
- gitmap nodes errors clear reports success across the fleet.
- gitmap see errors locally returns the correct aggregated output.
