# Subtask 03: Multi-IDE Registration, Pinned Projects & Conversation Sync Hook

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/229-nodes-deploy-repos-and-multi-ide-fleet-sync.md`  
> **Owned Files:**  
> - `cli/cmdnodes/nodes_deploy_ide.go` (NEW)  

---

## 1. Objectives

1. Implement `RegisterRemoteIDEs(client *ssh.Client, conn db.SSHConnection, repoName, remoteRepoPath string, withPinned bool) ([]string, error)` in `cli/cmdnodes/nodes_deploy_ide.go`:
   - Translates local repo path to remote OS path (`D:\work\<repoName>` on Windows, `~/work/<repoName>` on Linux).
   - Generates and executes remote script:
     - VS Code: Adds/updates entry in `projects.json` (`alefragnani.project-manager`).
     - Cursor: Adds/updates entry in `projects.json` (`alefragnani.project-manager`).
     - Antigravity IDE: Writes `~/.gemini/config/projects/<uuid>.json` with URL-encoded file URI.
     - GitHub Desktop: Executes `github <remoteRepoPath>` CLI if available.
   - If `withPinned` is true, updates remote `~/.gemini/config/pinned_projects.json` and adds `"pinned"` tag.
2. Implement conversation sync helper:
   - Packages local Antigravity conversation (`<conv-id>.db` and `brain/<conv-id>/`) associated with the repository.
   - Streams and unpacks onto target machine's `~/.gemini/antigravity/conversations/` and `~/.gemini/antigravity/brain/`.
