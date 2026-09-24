# Plan 94: SSH Macro & App Fleet Deployment, Remote Clone, REST Endpoint Triad, and Web UI Management Engine

Spec Reference: [02-spec/21-app/145-ssh-fleet-deploy-remote-clone-api-ui/01-overview.md](../../../02-spec/21-app/145-ssh-fleet-deploy-remote-clone-api-ui/01-overview.md)
Status: Completed
Execution Date: 2026-09-24

## Objective
Implement multi-node macro fleet distribution (`gitmap macro deploy ssh`, `peat deploy ssh`, `pea deploy ssh`), cross-node application update engine (`gitmap update --all`, `gitmap update all`, `gitmap ua`, `gitmap update <name>`, `gitmap update ls`), self-healing remote git clone with automatic key injection (`gitmap ssh clone`, `ssh-clone`, `ssh-c`), remote VS Code connection & file transfer bridge (`gitmap vscode remote`, `gitmap ssh cp`), embedded web-based management UI (`gitmap ui`, `gitmap <module> ui`) with REST API endpoints, inter-node REST communication triad daemon, and isolated local E2E verification suites.

## Consolidated Outcomes

### Subtask 01: Macro Fleet Deployment Commands
- **Files**: `cli/cmdmacro/macro_deploy_ssh.go`, `cli/cmdmacro/macro_deploy_ssh_test.go`, `cli/cmdmacro/macro_cmd.go`, `cli/cmd/rootutility.go`
- **Implemented**:
  - `gitmap macro deploy ssh --except <ids,ips,aliases>`
  - Command aliases: `gitmap peat deploy ssh`, `gitmap pea deploy ssh`, `gitmap deploy ssh`
  - Parallel dispatch to remote online nodes via SSH execution engine, broadcasting macro registries while honoring exception filters.
- **Verification**: Zero nested-if violations, positive boolean standards satisfied.

### Subtask 02: Multi-Node Update and Software Inventory Engine
- **Files**: `cli/cmdupdate/update_fleet.go`, `cli/cmdupdate/update_fleet_ls.go`, `cli/cmdupdate/update_fleet_test.go`, `cli/cmdupdate/exports.go`, `cli/cmd/rootutility.go`
- **Implemented**:
  - `gitmap update --all`, `gitmap update all`, `gitmap ua`: updates installed tools and packages across all remote SSH nodes in parallel, receiving JSON summaries from remote nodes and rendering an aligned terminal table.
  - `gitmap update <name> --except <id,alias,ip>`: targets a specific package across the fleet while excluding specified nodes.
  - `gitmap update ls`: queries all cluster nodes in parallel for installed software inventory via JSON responses and presents a unified table.
- **Verification**: Unit tests in `update_fleet_test.go`, boolean linter verified.

### Subtask 03: SSH Remote Clone and Self-Healing Auth
- **Files**: `cli/cmdssh/ssh_clone.go`, `cli/cmdssh/ssh_clone_auth.go`, `cli/cmdssh/ssh_clone_path.go`, `cli/cmdssh/ssh_clone_types.go`, `cli/cmdssh/ssh_clone_test.go`, `cli/cmdssh/exports.go`, `cli/cmdssh/ssh.go`, `cli/cmd/roottooling.go`
- **Implemented**:
  - `gitmap ssh clone`, `ssh-clone`, `ssh-c <repo|url> [git|path]`
  - Auto-resolves current repository git remote origin when omitted or specified as `git` / `.`.
  - Self-healing authentication: if remote clone fails due to publickey or permission denied, automatically detects failure, deploys local SSH public key to the target node or Git host, logs each diagnostic step, and resumes clone seamlessly.
- **Verification**: `ssh_clone_test.go`, flattened nested conditionals in `ssh_clone_path.go`.

### Subtask 04: VS Code Remote & File Transfer Bridge
- **Files**: `cli/cmdvscode/vscode_remote.go`, `cli/cmdvscode/vscode_cmd.go`, `cli/cmdssh/ssh_cp.go`, `cli/cmdssh/exports.go`, `cli/cmdssh/ssh.go`
- **Implemented**:
  - `gitmap vscode remote <node> [path]`: opens VS Code connected via `code --remote ssh-remote+<node> <path>`.
  - `gitmap vscode remote fix <node>`: diagnoses and clears stale remote VS Code server locks.
  - `gitmap ssh cp <src> <dest>`: bidirectional file transfer supporting `<local-path>` to `<node>:<remote-path>`, `<node>:<remote-path>` to `<local-path>`, `<node1>:<path>` to `<node2>:<path>`, and broadcast to all nodes.
- **Verification**: Boolean linter passed, nested-if linter passed.

### Subtask 05: Embedded Web UI & REST Endpoints
- **Files**: `cli/cmdui/ui_server.go`, `cli/cmdui/ui_assets.go`, `cli/cmdui/ui_editor.go`, `cli/cmdui/ui_types.go`, `cli/cmdui/ui_test.go`, `cli/cmd/ui_cmd.go`, `cli/cmd/rootutility.go`
- **Implemented**:
  - Embedded HTTP server supporting dynamic port binding and automated browser launch.
  - Endpoints: `gitmap ui`, `gitmap settings ui`, `gitmap commitin ui`, `gitmap ssh ui`, `gitmap macro ui`, `gitmap installer ui` (multi-OS CRUD with target node dropdown execution), `gitmap prompts ui` (templating, add empty, AI instruction JSON formatting, import/export), `gitmap schedules ui`, `gitmap import-export ui`.
  - In-browser syntax editor (`gitmap editor ui <node> <file>`) with save-back capabilities to remote SSH node filesystems.
- **Verification**: Unit tests in `ui_test.go`.

### Subtask 06: REST Triad Inter-Node Communication Daemon
- **Files**: `cli/cmddaemon/daemon_server.go`, `cli/cmddaemon/daemon_client.go`, `cli/cmddaemon/daemon_types.go`, `cli/cmddaemon/daemon_test.go`
- **Implemented**:
  - Post-join mutual authenticated REST API daemon (`gitmap daemon start`) running alongside GitMap instances.
  - Secure token authentication with transparent fallback to SSH execution (`sshexec`) if daemon is unavailable.
- **Verification**: Unit tests in `daemon_test.go`.

### Subtask 07: Isolated Local E2E Test Suites
- **Files**: `cli/tests/fleet_deploy_clone_ui_e2e_test.go`
- **Implemented**:
  - Dedicated E2E tests guarded by `//go:build e2e` verifying macro deployment, fleet update, remote clone argument parsing, and UI server initialization.
- **Verification**: Build-tag isolated to ensure zero disruption to CI test runners.

## Traceability
- **Spec**: `02-spec/21-app/145-ssh-fleet-deploy-remote-clone-api-ui/01-overview.md`
- **Index**: Registered in `.ai-memory/plans/01-index.md`
