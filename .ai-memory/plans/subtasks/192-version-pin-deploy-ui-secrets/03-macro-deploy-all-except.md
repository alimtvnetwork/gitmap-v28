# Subtask 192.3: Macro Fleet Deployment Commands (`deploy macro`, `macro deploy`)

## Spec Reference
- `02-spec/21-app/182-version-pinning-macro-deploy-ui-settings-secret-flags.md` §2.2

## Deliverables
1. **Routing & Dispatch**:
   - In `cli/cmdssh/ssh_deploy_cmd.go`: When `args[0]` is `"macro"` or `"macros"`, route to `cmdmacro.ExecuteMacroDeploySSH(args[1:])`.
   - Support `gitmap deploy macro all` and `gitmap macro deploy all` to deploy all macros across all nodes.
   - Support `gitmap deploy macro <macro-name> <node>` and `gitmap macro deploy <macro-name> <node>` to deploy a specific macro to a specific node.
2. **Filtering & Exclusions**:
   - Enforce `--except <id,alias,ip>` to sequence and exclude specific targets.
   - Confirm macro export (`gitmap macro export`) and import (`gitmap macro import`) work reliably.
