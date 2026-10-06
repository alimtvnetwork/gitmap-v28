# Subtask 02: Fleet Repository Deployment Engine

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/229-nodes-deploy-repos-and-multi-ide-fleet-sync.md`  
> **Owned Files:**  
> - `cli/cmdnodes/nodes_deploy_repo.go` (NEW)  
> - `cli/cmdnodes/nodes_deploy_repo_test.go` (NEW)  

---

## 1. Objectives

1. Implement `RunNodesDeployRepo(args []string) error` in `cli/cmdnodes/nodes_deploy_repo.go`:
   - Resolve local repository by slug using `store.OpenDefault()`.
   - Parse flags: `--target`, `--except`, `--exclude`, `--include`, `--accept`, `--include-main`, `--open-only`, `--dest`, `--with-ides`, `--with-pinned`, `--with-conversations`, `--from-local`, `--clone`, `--clean`, `--dry-run`, `--json`.
   - Filter target nodes via `FilterFleetNodes()`.
   - In-memory tar.gz streaming of repo files via `cmdssh.StreamFileToRemote` (excluding `.git/objects` bloat if requested, or excluding build directories like `node_modules/`, `target/`, `.venv/`, `dist/`).
   - Trigger remote unpack and remote post-deploy hooks (remote `gitmap rescan` and multi-IDE registration).
   - Render telemetry table with status badges (`● SUCCESS`, `○ OFFLINE`, `▲ FAILED`, `◌ DRY-RUN`).
2. Implement `RunNodesDeployRepos(args []string) error` to support multiple repos or `--all`.
3. Add unit tests in `cli/cmdnodes/nodes_deploy_repo_test.go`.
