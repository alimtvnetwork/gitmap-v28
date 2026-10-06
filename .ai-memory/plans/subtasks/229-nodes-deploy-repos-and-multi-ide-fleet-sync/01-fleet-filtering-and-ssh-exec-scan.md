# Subtask 01: Fleet Filtering Engine, SSH Exec Scan & Nodes Scan Commands

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/229-nodes-deploy-repos-and-multi-ide-fleet-sync.md`  
> **Owned Files:**  
> - `cli/cmdnodes/nodes_filter.go` (NEW)  
> - `cli/cmdnodes/nodes_filter_test.go` (NEW)  
> - `cli/cmdssh/ssh_exec_command.go` (MODIFY)  
> - `cli/cmdnodes/nodes_scan.go` (NEW)  

---

## 1. Objectives

1. Create a unified, canonical `NodeFilterOptions` struct and `FilterFleetNodes(conns, opts)` function in `cli/cmdnodes/nodes_filter.go`:
   - Supports `--target <alias>`, `--except <comma-separated>`, `--exclude <comma-separated>`, `--include <comma-separated>`, `--accept <comma-separated>`, `--include-main`, and `--open-only`.
   - Default excludes `"main"` unless `--include-main` is explicitly passed.
   - Automatically filters out local machine (`isLocalMachineConnection(c)`).
   - Fast TCP probe (1.5s timeout) when `OpenOnly` is true.
2. Add unit tests in `cli/cmdnodes/nodes_filter_test.go`.
3. Modify `cli/cmdssh/ssh_exec_command.go` around line 105:
   - Add `"scan"` and `"rescan"` to `isGitmapCoreCommand()`.
4. Implement `RunNodesScan` and `RunNodesRescan` in `cli/cmdnodes/nodes_scan.go`:
   - Broadcasts remote `gitmap scan` or `gitmap rescan` via SSH.
   - Outputs structured ANSI status table or JSON.
