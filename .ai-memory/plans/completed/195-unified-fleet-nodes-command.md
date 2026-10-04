# Plan 195: Unified Fleet Nodes Command (`gitmap nodes`)

## 1. Intent & Context
- The user requested a top-level command:
  `gitmap nodes # should show from gitmap ssh nodes, cluster nodes, or also from server-clients, do you understand it?`
- Unifies disparate node registries across 3 infrastructure subsystems:
  1. SSH connection vault / hosts (`SSHConnection`, `ssh_hosts` table)
  2. Cluster fleet database (`ClusterNode` table)
  3. Server-Clients (`sc`) peer-to-peer broadcast topology
- Supports deduplicated aggregation, live status probing, fast mode (`--fast`), JSON emission (`--json`), subsystem filtering (`--ssh`, `--cluster`, `--sc`), and target filtering.

## 2. Tasks Executed
1. **Export Live Probe Helper**:
   - Exported `ProbeHostsLiveStatus` in `cli/cmdssh/sshjoin_probe.go` for multi-package reuse.
2. **Top-Level Command Dispatch**:
   - Registered `nodes`, `node`, `allnodes`, `all-nodes`, and `fleet-nodes` in `coreClusterEntries()` within `cli/cmd/rootcore.go`.
3. **Core Aggregation & CLI Engine**:
   - Implemented `cli/cmd/nodes_cmd.go` with `collectUnifiedNodes`, `populateSSHHosts`, `populateSSHConnections`, `populateClusterNodes`, `probeUnifiedNodesLive`, `applyNodesFilter`, and ANSI table rendering.
   - Fixed `fmt.Fprintln` redundant newline formatting.
4. **Unit Tests**:
   - Implemented `cli/cmd/nodes_cmd_test.go` covering help displays, table rendering, subsystem filtering, and JSON serialization.
5. **Live Verification**:
   - Rebuilt `%USERPROFILE%\AppData\Local\gitmap-cli\gitmap.exe`.
   - Verified `gitmap nodes`, `gitmap nodes --fast`, `gitmap nodes --json`, `gitmap nodes main`, `gitmap node`, `gitmap allnodes`.
6. **Documentation**:
   - Authored specification `02-spec/21-app/185-unified-fleet-nodes-command.md`.
   - Updated index in `02-spec/21-app/readme.md`.
