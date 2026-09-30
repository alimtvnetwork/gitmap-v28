# Subtask 02: CLI Routing & Help Integration

## Scope
1. In `cli/cmd/nodes_cmd.go`:
   - Detect `clone`, `cfr`, `cfrp`, `clone-fix-repo`, `clone-fix-repo-pub` in `runUnifiedNodesCLI`.
   - Dispatch to `cmdnodes.RunNodesClone(args)`.
   - Update `printUnifiedNodesHelp()` to describe `nodes clone`, `nodes cfr`, and `nodes cfrp`.
   - Update the Supported Commands & Clusters Matrix (Table 2 and footer suggestions) in `nodes_cmd.go`.
2. In `cli/cmd/rootcore.go`:
   - Register root command mappings and aliases (`nodes-clone`, `nodes-cfr`, `nodes-cfrp`, `node clone`, `node cfr`, `node cfrp`).
