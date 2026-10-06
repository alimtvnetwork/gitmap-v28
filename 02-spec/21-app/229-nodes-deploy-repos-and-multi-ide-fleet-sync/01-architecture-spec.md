# Architecture Specification: Fleet Repository Deployment & Multi-IDE Sync Engine

> **Document Version:** 1.0.0  
> **Status:** Active  
> **Scope:** GitMap Core CLI (`cli/cmdnodes`, `cli/cmdssh`, `cli/workspacesync`, `cli/vscodepm`, `cli/desktop`, `cli/cmdcursor`, `cli/cmdagy`)  
> **Traceability:** Task-229  

---

## 1. Executive Summary & Problem Statement

GitMap operates across heterogeneous multi-node development environments (Windows and Linux nodes such as `main`, `u1`, `w1`, `w2`, `w3`). While GitMap supports:
- Cloning repositories from remote Git hosting via `gitmap nodes clone` / `cfr`,
- Pushing VS Code / Cursor workspace metadata via `gitmap nodes send-projects`,
- Deploying Antigravity Manager accounts via `gitmap nodes deploy agm-accounts`,

there was previously **no single unified command** to:
1. Deploy a repository directly from one machine to another machine across the fleet (supporting both peer-to-peer archive streaming for local/uncommitted work and remote Git clone delegation),
2. Trigger remote repository discovery (`gitmap scan` / `gitmap rescan`) on the target nodes so their local `gitmap.db` immediately indexes the new repository,
3. Automatically register the newly deployed repository across **all installed IDEs** on the target node:
   - **VS Code**: Added to `projects.json` (`alefragnani.project-manager`).
   - **Cursor**: Added to `projects.json` (`alefragnani.project-manager`).
   - **Antigravity IDE**: Added to `~/.gemini/config/projects/<uuid>.json` with URL-encoded file URI.
   - **GitHub Desktop**: Registered into GitHub Desktop via the `github <path>` CLI shim.
4. Support options to synchronize **pinned projects** (`--with-pinned`) from `~/.gemini/config/pinned_projects.json` and **conversations** (`--with-conversations`) from `~/.gemini/antigravity/conversations/<id>.db` and `~/.gemini/antigravity/brain/<id>/`.
5. Support granular machine filtering flags: `--target <alias>`, `--except <comma-separated>` (default excluding `main`), `--include <list>`, `--exclude <list>`, `--open-only`, `--dest <path>`, `--dry-run`, and `--json`.

---

## 2. Command Architecture & Syntax

### 2.1 Primary CLI Commands

```bash
# Deploy a single repository to fleet nodes
gitmap nodes deploy repo <slug|path|url> [dest] [flags]
gitmap nodes deploy-repo <slug|path|url> [dest] [flags]

# Deploy multiple repositories to fleet nodes
gitmap nodes deploy repos [slug1,slug2,...|all] [dest] [flags]
gitmap nodes deploy-repos [slug1,slug2,...|all] [dest] [flags]

# Fleet-wide remote repository scan / rescan
gitmap nodes scan [target] [flags]
gitmap nodes rescan [target] [flags]
```

### 2.2 Flag Specification

| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--target` | `-t` | `string` | `""` | Target a specific node alias (e.g. `u1`, `w1`). If omitted, targets all enrolled nodes. |
| `--except` | `-e` | `string` | `"main"` | Comma-separated list of nodes to exclude. Automatically includes `main` by default. |
| `--exclude` | | `string` | | Alias for `--except`. |
| `--include` | | `string` | `""` | Comma-separated whitelist of nodes to target. |
| `--accept` | | `string` | `""` | Alias for `--include`. |
| `--include-main`| | `bool` | `false` | Explicitly include the orchestrator `main` node. |
| `--open-only` | | `bool` | `true` | Skip offline or unreachable nodes immediately via fast preflight probe. |
| `--dest` | `-d` | `string` | `""` | Target root work directory on remote nodes (`D:\work` on Windows, `~/work` on Linux). |
| `--with-ides` | | `bool` | `true` | Register repository across VS Code, Cursor, Antigravity, and GitHub Desktop. |
| `--with-pinned`| | `bool` | `false` | Synchronize pinned project status to remote `pinned_projects.json`. |
| `--with-conversations` | | `bool` | `false` | Package, remap, and transfer associated Antigravity conversations and brain logs. |
| `--from-local` | | `bool` | `false` | Force direct local-to-remote archive transfer over SSH even if remote URL exists. |
| `--clone` | | `bool` | `false` | Force remote Git clone from origin URL over SSH. |
| `--clean` | | `bool` | `true` | Exclude temporary directories (`node_modules`, `dist`, `target`, `.venv`) during local transfer. |
| `--dry-run` | `-n` | `bool` | `false` | Simulate deployment and remote actions without writing to disk. |
| `--json` | `-j` | `bool` | `false` | Output structured JSON telemetry and results. |

---

## 3. High-Level Subsystem Architecture

```
[ Local GitMap Master (e.g. w3) ]
       │
       ├── 1. Repository Discovery & Resolution (store.OpenDefault() / filepath.Abs())
       ├── 2. Optional Pinned Status Check (~/.gemini/config/pinned_projects.json)
       ├── 3. Optional Conversation Bundling (~/.gemini/antigravity/conversations/ & brain/)
       ├── 4. Fleet Node Filtering & Fast Liveness Probe (cmdssh.ProbeTCP 1500ms)
       │
       ▼ (Parallel SSH Streams)
┌─────────────────────────────────┐       ┌─────────────────────────────────┐
│     Remote Windows Node (w1)     │       │      Remote Linux Node (u1)     │
│  - Staging: C:\Windows\Temp\... │       │  - Staging: /tmp/...            │
│  - Destination: D:\work\<repo>  │       │  - Destination: ~/work/<repo>   │
├─────────────────────────────────┤       ├─────────────────────────────────┤
│  Post-Deploy Actions:           │       │  Post-Deploy Actions:           │
│  1. Unpack / Clone Repo         │       │  1. Unpack / Clone Repo         │
│  2. Remote `gitmap rescan`      │       │  2. Remote `gitmap rescan`      │
│  3. Multi-IDE Registration:     │       │  3. Multi-IDE Registration:     │
│     - VS Code (projects.json)   │       │     - VS Code (projects.json)   │
│     - Cursor (projects.json)    │       │     - Cursor (projects.json)    │
│     - Antigravity (<uuid>.json) │       │     - Antigravity (<uuid>.json) │
│     - GitHub Desktop (`github`) │       │     - GitHub Desktop (if avail) │
│  4. Pinned Projects Sync (opt)  │       │  4. Pinned Projects Sync (opt)  │
│  5. Conversations Sync (opt)    │       │  5. Conversations Sync (opt)    │
└─────────────────────────────────┘       └─────────────────────────────────┘
```

---

## 4. Multi-IDE Registration Details

### 4.1 Visual Studio Code
- **Configuration File**: `%APPDATA%\Code\User\globalStorage\alefragnani.project-manager\projects.json` (Windows) / `~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json` (Linux).
- **Entry Schema**:
  ```json
  {
    "name": "<repo-name>",
    "rootPath": "<remote-repo-path>",
    "paths": ["<remote-repo-path>"],
    "tags": ["gitmap", "fleet"],
    "enabled": true
  }
  ```
- **Execution**: Upserted atomically via PowerShell on Windows or Shell script on POSIX.

### 4.2 Cursor IDE
- **Configuration File**: `%APPDATA%\Cursor\User\globalStorage\alefragnani.project-manager\projects.json` (Windows) / `~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json` (Linux).
- **Execution**: Simultaneously updated alongside VS Code during remote post-deploy execution.

### 4.3 Google Antigravity IDE
- **Configuration Directory**: `~/.gemini/config/projects/`
- **File Name**: `<uuid>.json`
- **Schema**:
  ```json
  {
    "id": "<uuid>",
    "name": "<repo-name>",
    "projectResources": {
      "resources": [
        {
          "gitFolder": {
            "folderUri": "file:///<url-encoded-path>",
            "defaultBranch": "main"
          }
        }
      ]
    },
    "settings": {},
    "updatedAt": "<RFC3339-timestamp>",
    "isWorkspaceOnly": false
  }
  ```

### 4.4 GitHub Desktop
- **Invocation**: Executes `github <remote-repo-path>` CLI shim.
- **Probe**: Checks `%LOCALAPPDATA%\GitHubDesktop\bin\github.bat` or PATH on Windows, and PATH on POSIX. If present, runs detached command to add repository to tracked database.

---

## 5. Security & Permission Hardening

1. **Token & Secret Exclusions**: All local transfer archives strictly exclude secrets (`.env`, `credentials.json`, `*.pem`, `repo-secrets/`, `node_modules/`, `dist/`).
2. **Permission Setting**: On POSIX targets, permissions are hardened via `chmod 700` for user project and config directories and `chmod 600` for sensitive conversation and token files.
3. **Orchestrator Safeguard**: The `main` node is excluded by default from mass fleet deployments to prevent accidental overrides of central repositories.
