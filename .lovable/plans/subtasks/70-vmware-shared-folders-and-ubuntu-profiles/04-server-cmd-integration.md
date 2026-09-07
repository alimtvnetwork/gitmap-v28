# Subtask 04: Server-Cmd Cluster & SSH Delegation Command Integration

## Status
Pending

## Context & Objectives
Integrate `gitmap server-cmd` (and `server-cmds`, `scmd`) based on Kubernetes and server remote execution patterns (`aukgit/kubernetes-training-v1` `05-server-cmds/02-run-cmd-v2.sh`):
1. **Target Selectors**:
   - `all`: Execute across all nodes (servers and workers). Maps to `cluster.ServersClients`.
   - `control`: Execute across control plane / server nodes. Maps to `cluster.ServersOnly`.
   - `workers`: Execute across worker / client nodes. Maps to `cluster.ClientsOnly`.
   - `<node-id-or-alias>`: Execute on a specific node by ID, display ID, or hostname.
2. **Command & Script Modes**:
   - Direct command execution: `gitmap server-cmd <target> "<command>"`.
   - `--script`: Deploy on-the-fly script (`/tmp/on-the-fly-cmd/script-<timestamp>.sh`) to remote nodes, execute, capture stdout/stderr, and clean up.
   - `--sudo`: Run command with `sudo` / `sudo -S` non-interactive privilege escalation.
   - `--config <path>`: Fallback JSON topology file (such as `01-config-sample.json`), falling back to SQLite `ClusterNode` / `SSHConnection` pools.
3. **CLI Dispatching & AST Parity**:
   - Declare `CmdServerCmd = "server-cmd"`, `CmdServerCmds = "server-cmds"`, `CmdServerCmdAlias = "scmd"` in `constants_cli.go`.
   - Synchronize `cmd_constants_test.go` `topLevelCmds()`.
   - Register dispatch entry in `gitmap/cmd/rootcore.go` under `coreClusterEntries()`.

## Files to Create / Modify
- [MODIFY] `gitmap/constants/constants_cli.go`: Declare `CmdServerCmd`, `CmdServerCmds`, `CmdServerCmdAlias` under `// gitmap:cmd top-level`.
- [MODIFY] `gitmap/constants/cmd_constants_test.go`: Add entries to `topLevelCmds()`.
- [MODIFY] `gitmap/cmd/rootcore.go`: Register dispatch in `coreClusterEntries()`.
- [MODIFY] `gitmap/cmd/rootutility.go`: Map aliases in `canonicalCommandName`.
- [NEW] `gitmap/cmd/servercmd.go` (<= 200 lines): Flag parsing, target routing, and entry point.
- [NEW] `gitmap/cmd/servercmd_run.go` (<= 200 lines): Remote execution engine, `--script` deployer, and `--sudo` escalation.
- [NEW] `gitmap/cmd/servercmd_test.go` (<= 200 lines): Unit tests for target selector parsing and flag validation.

## Verification Steps
- AST parity test: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1` exits 0.
- Command tests: `go test -v ./gitmap/cmd/... -run "TestServerCmd"` exits 0.
