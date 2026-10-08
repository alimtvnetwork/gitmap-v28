# Subtask 08: Smart Safe Repository Creation & Remote Clone Fallback

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `08-smart-repo-create`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md` (Section 4)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Upgrade the repository creation engine (`gitmap repo create`, alias `gitmap cr`) to prevent catastrophic data loss by introducing remote pre-flight checks, interactive clone confirmations, headless safety gates, and removing legacy force-pushing.

### The Problem
When running `gitmap repo create <name>`, if the target repository already exists on GitHub, `cli/cmd/repo_create_remote.go` previously caught the `gh` error and executed `git push -u origin main --force`. In multi-developer or automated AI agent workflows, this caused inadvertent destruction of existing remote repositories and commit histories without confirmation.

### The Solution
Implement the **Smart Safe Repository Creation Protocol**:
1. **Pre-Flight Remote Existence Check:** Prior to modifying local repository state or remotes, probe GitHub using `git ls-remote` or `gh repo view`.
2. **Interactive Confirmation:** If the remote exists and the session is interactive, prompt:
   ```text
   Remote repository 'org/repo' already exists. Clone existing repo instead? [y/N]: 
   ```
   If accepted, hand off execution cleanly to `cmdclone.ExecuteClone()`.
3. **Headless Safety Gate (`--clone-on-exists`):** If the session is headless (`!isInteractiveStdin()`):
   - With `--clone-on-exists`: Automatically fallback to cloning.
   - Without `--clone-on-exists`: Abort with error code `E1089` (`ErrRemoteRepoExists`).
4. **Permanent Force-Push Removal:** Eliminate all `push --force` fallback code from `repo_create_remote.go`.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmd/repo_create_remote.go` | Modify / Refactor | Add `probeRemoteExistence()`, replace `handleExistingRemoteRepo()` with interactive prompt and fallback logic. |
| `cli/cmd/create_cmd.go` | Modify | Parse `--clone-on-exists` flag, propagate to creation params. |
| `cli/cmd/create_ops.go` | Modify | Update `executeCreateRepo()` lifecycle to invoke pre-flight checks before local Git initialization. |
| `cli/cmd/repo_create_remote_test.go` | New | Comprehensive unit tests for remote probing, interactive mocking, headless gate validation, and error aborts. |

---

## 3. Detailed Implementation Requirements

### 3.1 Remote Probing Implementation
```go
package cmd

import (
	"os/exec"
	"strings"
	"time"
)

func probeRemoteRepository(slug string) (bool, error) {
	// First check via git ls-remote (faster, works without gh CLI auth)
	remoteURL := "https://github.com/" + slug + ".git"
	cmd := exec.Command("git", "ls-remote", "--heads", remoteURL)
	out, err := cmd.CombinedOutput()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return true, nil
	}

	// Secondary check via gh repo view
	ghCmd := exec.Command("gh", "repo", "view", slug, "--json", "name")
	if ghErr := ghCmd.Run(); ghErr == nil {
		return true, nil
	}

	return false, nil
}
```

### 3.2 Safe Fallback Flow in `pushRemoteRepo()`
```go
func handleRemoteCollision(absDir, slug string, isCloneOnExists bool) (string, error) {
	if isInteractiveStdin() {
		fmt.Printf("\n  %sNotice: Remote repository '%s' already exists.%s\n", constants.ColorYellow, slug, constants.ColorReset)
		fmt.Print("  Clone existing repo instead? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		answer := strings.ToLower(strings.TrimSpace(input))

		if answer == "y" || answer == "yes" {
			return dispatchCloneFallback(absDir, slug)
		}
		return "", apperror.NewSimple("repository creation aborted: remote already exists", "E1089")
	}

	if isCloneOnExists {
		fmt.Printf("  %sHeadless gate: --clone-on-exists active, cloning existing repo...%s\n", constants.ColorCyan, constants.ColorReset)
		return dispatchCloneFallback(absDir, slug)
	}

	return "", apperror.NewSimple("remote repository '"+slug+"' already exists; use --clone-on-exists or gitmap clone", "E1089")
}
```

---

## 4. Acceptance Criteria

- [ ] Force-push command (`git push ... --force`) is completely removed from `cli/cmd/repo_create_remote.go`.
- [ ] Pre-flight remote probe correctly detects existing repositories.
- [ ] Interactive prompt presents `"Clone existing repo instead? [y/N]"` and triggers `cmdclone` upon `y`.
- [ ] Headless environments without `--clone-on-exists` abort safely with code `E1089`.
- [ ] Headless environments with `--clone-on-exists` trigger `cmdclone` fallback.
- [ ] Unit tests in `repo_create_remote_test.go` cover interactive yes, interactive no, headless without flag, and headless with flag.

---

## 5. Verification Commands

```powershell
# 1. Run unit tests
go test -v ./cli/cmd/... -run TestRepoCreateRemote

# 2. Test headless failure when repo exists
gitmap repo create gitmap --clone-on-exists=false

# 3. Test headless clone-on-exists
gitmap repo create gitmap --clone-on-exists --dry-run

# 4. Guideline compliance check
python 03-ai-scripts/05-guideline-autofixer.py cli/cmd --check-only
```
