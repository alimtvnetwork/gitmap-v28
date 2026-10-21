# Subtask-03: Remote Node Clone IDE Registration

- **Subtask Path:** `.ai-memory/plans/subtasks/gitmap-enhancement-and-issue-resolution/03-subtask-nodes-clone-ide-registration.md`
- **Parent Plan:** [.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md](../../gitmap-enhancement-and-issue-resolution.md)
- **Primary Owner:** Worker 03
- **Component Scope:** `cli/cmdnodes/`
- **Exclusive File Ownership:**
  - `cli/cmdnodes/nodes_clone_remote.go`
  - `cli/cmdnodes/nodes_clone_ide_hook.go`
  - `cli/cmdnodes/nodes_clone_ide_hook_test.go`
- **Status:** Pending Execution
- **Target Release:** GitMap v6.590.0

---

## 1. Context & Architectural Root Cause

When running `gitmap nodes clone <repo>` across remote SSH nodes, the repository is cloned successfully to the remote filesystem (e.g. `D:\work\<repo>` on Windows or `~/work/<repo>` on Linux). However, unlike local clones which automatically register in local developer tools, remote clones were failing to register on the remote node's:
1. **GitHub Desktop** (`github open` command or repository registry database).
2. **VS Code Project Manager** (`~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json` or Windows equivalent `%APPDATA%\Code\User\globalStorage\alefragnani.project-manager\projects.json`).
3. **Cursor IDE** recent projects registry.
4. **Antigravity IDE** instance workspace associations.

### Root Cause Analysis (RCA)
In `cli/cmdnodes/nodes_deploy_repo.go`, `RegisterRemoteIDEs` (defined in `cli/cmdnodes/nodes_deploy_ide.go:27`) is already implemented and proven for repo deployment (`runRemoteIDEsHook`). However, `cli/cmdnodes/nodes_clone_remote.go` merely runs `gitmap clone` over SSH and populates `RemoteCloneNodeResult`, completely omitting the invocation of `RegisterRemoteIDEs`.

---

## 2. Implementation Blueprint

### Step 1: Create Dedicated Hook Helper (`cli/cmdnodes/nodes_clone_ide_hook.go`)
Create a modular, single-responsibility helper to extract repository metadata and dispatch remote IDE registration without cluttering `nodes_clone_remote.go`:

```go
package cmdnodes

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// resolveClonedRepoName extracts the target folder name from clone options.
func resolveClonedRepoName(opts NodesCloneOptions, fileName string) string {
	if len(opts.PassArgs) > 0 {
		first := opts.PassArgs[0]
		base := filepath.Base(strings.TrimSuffix(first, ".git"))
		if base != "" && base != "." && base != "/" {
			return base
		}
	}
	if fileName != "" {
		return strings.TrimSuffix(filepath.Base(fileName), ".git")
	}
	return "repository"
}

// resolveRemoteRepoFullPath constructs the remote destination path for the repository.
func resolveRemoteRepoFullPath(conn db.SSHConnection, opts NodesCloneOptions, repoName string) string {
	targetDir := resolveTargetDirForOS(isWindowsNode(conn), opts)
	if isWindowsNode(conn) {
		return targetDir + `\` + repoName
	}
	return strings.TrimSuffix(targetDir, "/") + "/" + repoName
}

// executeRemoteCloneIDERegistration invokes RegisterRemoteIDEs on successful remote clones.
func executeRemoteCloneIDERegistration(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileName string) ([]string, error) {
	repoName := resolveClonedRepoName(opts, fileName)
	remotePath := resolveRemoteRepoFullPath(conn, opts, repoName)
	return RegisterRemoteIDEs(client, conn, repoName, remotePath, false)
}
```

### Step 2: Integrate into `nodes_clone_remote.go`
In `cli/cmdnodes/nodes_clone_remote.go`, update `executeRemoteCloneSession` to trigger the registration hook whenever clone execution succeeds:

1. Check if `res.Status == "success"`.
2. Dispatch `executeRemoteCloneIDERegistration(client, conn, opts, fileName)`.
3. If successful, append registered IDE names (e.g. `Registered in: GitHub Desktop, VS Code, Cursor, Antigravity`) to `res.Details`.
4. If an error occurs during IDE registration, record the warning in `res.Details` without marking the git clone as failed.

### Step 3: Coding Guideline Invariants
- **Function Cap:** Every newly added function must be `<= 15` lines.
- **Positive Booleans:** Variable names must use `isWin`, `hasIDESuccess`, `hasRegistrationError`. Negative booleans are strictly forbidden.
- **Error Sovereignty:** Internal helpers must return `error` or `([]string, error)` and let the session runner handle reporting.

---

## 3. Detailed Step-by-Step Action Items

- [ ] **Step 3.1:** Author `cli/cmdnodes/nodes_clone_ide_hook.go` with `resolveClonedRepoName`, `resolveRemoteRepoFullPath`, and `executeRemoteCloneIDERegistration`.
- [ ] **Step 3.2:** Modify `cli/cmdnodes/nodes_clone_remote.go` to invoke `executeRemoteCloneIDERegistration` upon successful clone completion.
- [ ] **Step 3.3:** Ensure `res.Details` reflects IDE registration results cleanly for CLI terminal output and JSON responses.
- [ ] **Step 3.4:** Author unit tests in `cli/cmdnodes/nodes_clone_ide_hook_test.go` testing repo name resolution, remote path formatting for Windows/Linux, and mock registration handling.
- [ ] **Step 3.5:** Validate that all functions in `cli/cmdnodes/nodes_clone_ide_hook.go` and modified regions of `nodes_clone_remote.go` remain `<= 15` lines with zero nested `if` statements.

---

## 4. Acceptance Criteria & Verification

1. **AC-3.1:** Given a Windows SSH node `w2`, when `gitmap nodes clone <repo>` succeeds, PowerShell registration commands for GitHub Desktop (`github open`), VS Code Project Manager, Cursor, and Antigravity execute cleanly over SSH.
2. **AC-3.2:** Given a Linux SSH node `u1`, when `gitmap nodes clone <repo>` succeeds, shell registration scripts for VS Code and Cursor update `projects.json` cleanly.
3. **AC-3.3:** If the remote node does not have GitHub Desktop or VS Code installed, `RegisterRemoteIDEs` returns a graceful partial registration without failing the clone task.
4. **Verification Command:**
   ```bash
   go test -v -run TestRemoteCloneIDEHook ./cli/cmdnodes
   ```
