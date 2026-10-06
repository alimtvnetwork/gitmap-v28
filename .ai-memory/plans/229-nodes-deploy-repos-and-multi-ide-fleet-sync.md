# Master Execution Plan: Fleet Repository Deployment & Multi-IDE Sync Engine

> **Plan Version:** 1.0.0  
> **Status:** Active  
> **Parent Task:** Task-229  
> **Associated Spec:** `02-spec/21-app/229-nodes-deploy-repos-and-multi-ide-fleet-sync/01-architecture-spec.md`  

---

## 1. Goal & Objectives

Enable GitMap to:
1. Deploy any repository across the fleet (`gitmap nodes deploy repo <slug>` / `gitmap nodes deploy repos <slugs>`), supporting both direct local archive streaming and remote Git origin cloning.
2. Automatically trigger remote repository scanning (`gitmap scan` / `gitmap rescan`) on the target nodes so their internal SQLite split-database (`gitmap.db`) immediately tracks the new repo.
3. Automatically register the deployed repository across all installed IDEs on the target node:
   - VS Code (`projects.json` Project Manager)
   - Cursor (`projects.json` Project Manager)
   - Antigravity IDE (`~/.gemini/config/projects/<uuid>.json`)
   - GitHub Desktop (via `github` CLI invocation)
4. Support optional flags for **pinned projects** (`--with-pinned`) and **conversations** (`--with-conversations`).
5. Support granular machine filtering flags: `--target <alias>`, `--except <list>` (default `main`), `--include <list>`, `--exclude <list>`, `--open-only`, `--dry-run`, `--json`.
6. Enable fleet-wide scanning commands: `gitmap nodes scan` and `gitmap nodes rescan`, plus fixing `scan` in `isGitmapCoreCommand`.

---

## 2. Subtask Breakdown

The work is partitioned into 4 discrete, disjoint subtasks:

- **Subtask 01:** `01-fleet-filtering-and-ssh-exec-scan.md`
  - Implement unified fleet node filtering in `cli/cmdnodes/nodes_filter.go` and unit tests in `cli/cmdnodes/nodes_filter_test.go`.
  - Fix `isGitmapCoreCommand` in `cli/cmdssh/ssh_exec_command.go` to recognize `"scan"` and `"rescan"`.
  - Implement `RunNodesScan` and `RunNodesRescan` in `cli/cmdnodes/nodes_scan.go`.

- **Subtask 02:** `02-nodes-deploy-repo-engine.md`
  - Implement core repository packaging and SSH streaming deployment in `cli/cmdnodes/nodes_deploy_repo.go`.
  - Implement `RunNodesDeployRepo` and `RunNodesDeployRepos` supporting `--from-local`, `--clone`, `--clean`, `--dry-run`, and `--json`.
  - Implement unit tests in `cli/cmdnodes/nodes_deploy_repo_test.go`.

- **Subtask 03:** `03-multi-ide-and-pinned-conv-sync.md`
  - Implement `cli/cmdnodes/nodes_deploy_ide.go` to orchestrate multi-IDE registration across VS Code, Cursor, Antigravity, and GitHub Desktop.
  - Implement pinned projects synchronization (`--with-pinned`) targeting `~/.gemini/config/pinned_projects.json`.
  - Implement conversation and brain log synchronization (`--with-conversations`) targeting `~/.gemini/antigravity/conversations/` and `~/.gemini/antigravity/brain/`.

- **Subtask 04:** `04-cli-routing-helptext-and-tests.md`
  - Wire subcommands and aliases into `cli/cmd/nodes_cmd.go` and `cli/cmd/rootcore.go`.
  - Document all commands and flags in `cli/helptext/nodes.md`.
  - Verify unit tests, run linter checks, and validate dry-run execution across the live cluster.

---

## 3. Worker Allocation & Concurrency (A = 2, H = 2)

- **Worker 01:** Owns Subtask 01 (`nodes_filter.go`, `nodes_filter_test.go`, `ssh_exec_command.go`, `nodes_scan.go`).
- **Worker 02:** Owns Subtask 02 & Subtask 03 (`nodes_deploy_repo.go`, `nodes_deploy_repo_test.go`, `nodes_deploy_ide.go`).
- **Lead Orchestrator:** Wires CLI routing in `cli/cmd/nodes_cmd.go`, updates help text, runs test suite and linters, performs version bump and release ceremony.

---

## 4. Execution Tracking Ledger

| Subtask ID | File | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| Subtask-01 | `.ai-memory/plans/subtasks/229-nodes-deploy-repos-and-multi-ide-fleet-sync/01-fleet-filtering-and-ssh-exec-scan.md` | Worker 01 | DONE | `nodes_filter.go`, `nodes_filter_test.go`, `nodes_scan.go`, `ssh_exec_command.go` (100% tests pass) |
| Subtask-02 | `.ai-memory/plans/subtasks/229-nodes-deploy-repos-and-multi-ide-fleet-sync/02-nodes-deploy-repo-engine.md` | Worker 02 | DONE | `nodes_deploy_repo.go`, `nodes_deploy_repo_test.go` (100% tests pass) |
| Subtask-03 | `.ai-memory/plans/subtasks/229-nodes-deploy-repos-and-multi-ide-fleet-sync/03-multi-ide-and-pinned-conv-sync.md` | Worker 02 | DONE | `nodes_deploy_ide.go` (VS Code, Cursor, Antigravity, GitHub Desktop) |
| Subtask-04 | `.ai-memory/plans/subtasks/229-nodes-deploy-repos-and-multi-ide-fleet-sync/04-cli-routing-helptext-and-tests.md` | Lead | DONE | CLI routing in `nodes_cmd.go`, root aliases in `rootcore.go`, `nodes.md` docs, dry-run verified |
