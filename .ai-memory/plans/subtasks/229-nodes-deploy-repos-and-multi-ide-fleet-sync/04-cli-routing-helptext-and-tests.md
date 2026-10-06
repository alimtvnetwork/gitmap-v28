# Subtask 04: CLI Routing, Root Dispatch, Help Text & Verification

> **Subtask ID:** Subtask-04  
> **Parent Plan:** `.ai-memory/plans/229-nodes-deploy-repos-and-multi-ide-fleet-sync.md`  
> **Owned Files:**  
> - `cli/cmd/nodes_cmd.go` (MODIFY)  
> - `cli/cmd/rootcore.go` (MODIFY)  
> - `cli/helptext/nodes.md` (MODIFY)  

---

## 1. Objectives

1. Wire commands into `cli/cmd/nodes_cmd.go`:
   - Route `gitmap nodes deploy repo` / `deploy repos` to `cmdnodes.RunNodesDeployRepo` / `cmdnodes.RunNodesDeployRepos`.
   - Route `gitmap nodes scan` / `nodes rescan` to `cmdnodes.RunNodesScan` / `cmdnodes.RunNodesRescan`.
2. Wire root aliases in `cli/cmd/rootcore.go`:
   - `nodes-deploy-repo`, `deploy-repo`, `nodes-scan`, `fleet-scan`.
3. Update `cli/helptext/nodes.md` with complete documentation, flags, and usage examples.
4. Run Go tests across `cli/cmdnodes/` and run linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-relative-paths.py`).
