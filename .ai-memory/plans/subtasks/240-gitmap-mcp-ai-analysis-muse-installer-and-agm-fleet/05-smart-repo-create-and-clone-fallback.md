# Subtask 05: Smart Repo Creation & Clone Fallback

> **Subtask ID:** Subtask-05  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmd/`, `cli/cmdrepo/`  
> **Owned Files:**  
> - `cli/cmd/repo_create_smart.go`  
> - `cli/cmd/repo_create_remote.go`  
> - `cli/cmd/repo_cmd_dispatch.go`  
> - `cli/cmd/create_cmd.go`  
> - `cli/cmdrepo/repo_create_smart.go`  
> - `cli/cmdrepo/repo_create_smart_test.go`  

---

## 1. Concrete Objectives

1. **Remote Repository Pre-Flight Inspection:**
   - Implement `ProbeRemoteRepoStatus(slug string)` in `cli/cmdrepo/repo_create_smart.go`.
   - Probe remote existence via `gh repo view <slug> --json name,url,sshUrl,isPrivate,defaultBranchRef` with zero side-effects.
   - Fall back to `git ls-remote <ssh_url>` when `gh` CLI is unauthenticated or offline.
2. **Interactive & Automated Collision Handling:**
   - Eliminate dangerous blind force-pushes in `handleExistingRemoteRepo` (`cli/cmd/repo_create_remote.go`).
   - If remote repository exists:
     - In interactive mode: Display ANSI prompt `"⚠️ Remote repository '<slug>' already exists on GitHub. Would you like to clone it instead? [y/N/ui]"`.
     - In non-interactive mode: Check for `--clone-if-exists` or `--fallback-clone` flags. If flag is present, transition seamlessly to clone. If absent, fail gracefully with structured error `apperror.NewConflict`.
3. **Local Directory Collision Validation:**
   - Validate destination directory state prior to cloning or creating:
     - Directory does not exist or is empty: Proceed with clone.
     - Directory exists and contains `.git` tracking matching remote origin: Verify origin URL, fetch latest default branch, and notify user that repository is already connected.
     - Directory exists with unrelated non-git files: Return validation error prohibiting accidental file overwrite.
4. **Registration into GitMap Split-DB:**
   - Upon successful creation or clone fallback, automatically register the repository in GitMap's repository cache database (`.gitmap/data/installation/repodb/` and `gitmap.db`).
5. **Command Line Interface & Flag Support:**
   - Expose flags `--clone-if-exists` (`-c`), `--fallback-clone`, `--yes` (`-y`), `--ui`, `--json` (`-j`), and `--dry-run` (`-n`) on `gitmap repo create` and `gitmap create`.

---

## 2. Core Domain Types & Structs

```go
package cmdrepo

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RepoCreateAction represents the action executed by smart repo create.
type RepoCreateAction string

const (
	ActionCreatedNew RepoCreateAction = "created_new"
	ActionClonedRepo RepoCreateAction = "cloned_existing"
	ActionLinkedRepo RepoCreateAction = "linked_existing"
	ActionAborted    RepoCreateAction = "aborted"
)

// RepoCreateOptions encapsulates user and agent arguments.
type RepoCreateOptions struct {
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	LocalDir         string `json:"local_dir"`
	IsPublic         bool   `json:"is_public"`
	IsPrivate        bool   `json:"is_private"`
	IsSkipRemote     bool   `json:"is_skip_remote"`
	CloneIfExists    bool   `json:"clone_if_exists"`
	IsNonInteractive bool   `json:"is_non_interactive"`
	IsInteractiveUI  bool   `json:"is_interactive_ui"`
	IsDryRun         bool   `json:"is_dry_run"`
	IsJSON           bool   `json:"is_json"`
}

// RemoteProbeResult encapsulates remote repository status.
type RemoteProbeResult struct {
	Exists          bool   `json:"exists"`
	Owner           string `json:"owner"`
	Name            string `json:"name"`
	SSHURL          string `json:"ssh_url"`
	HTTPSURL        string `json:"https_url"`
	DefaultBranch   string `json:"default_branch"`
	IsPrivate       bool   `json:"is_private"`
	ProbeDurationMs int64  `json:"probe_duration_ms"`
}

// RepoCreateResult contains execution telemetry and result status.
type RepoCreateResult struct {
	ActionTaken   RepoCreateAction `json:"action_taken"`
	Slug          string           `json:"slug"`
	LocalPath     string           `json:"local_path"`
	RemoteURL     string           `json:"remote_url"`
	DefaultBranch string           `json:"default_branch"`
	ExecutionMs   int64            `json:"execution_ms"`
	Message       string           `json:"message"`
}
```

---

## 3. Implementation Checklist

- [ ] **1. Remote Prober Implementation (`cli/cmdrepo/repo_create_smart.go`):**
  - Implement `ProbeRemoteRepoStatus(slug string) (*RemoteProbeResult, error)`.
  - Assemble command `gh repo view <slug> --json name,owner,sshUrl,url,isPrivate,defaultBranchRef`.
  - Parse JSON response; if exit code 0 and valid JSON returned, set `Exists: true`.
  - If `gh` CLI returns exit code 1 with "could not resolve to a Repository", set `Exists: false`.
  - Implement secondary fallback using `git ls-remote git@github.com:<slug>.git`.
- [ ] **2. Directory State & Safety Checks (`cli/cmdrepo/repo_create_smart.go`):**
  - Implement `CheckLocalDirState(dirPath string) (isMissing bool, isEmpty bool, isGitRepo bool, err error)`.
  - Use `os.ReadDir` to determine file emptiness.
  - Check `.git` directory presence and inspect origin remote URL via `git config --get remote.origin.url`.
- [ ] **3. Interactive & Automated Fallback Logic (`cli/cmd/repo_create_smart.go`):**
  - Implement `HandleRepoCollision(probe *RemoteProbeResult, opts RepoCreateOptions) (bool, error)`.
  - If `opts.CloneIfExists` or `opts.IsNonInteractive && opts.CloneIfExists`: automatically approve clone.
  - If `opts.IsInteractiveUI` or stdin is a TTY: render terminal prompt with `[y/N/ui]`.
  - If user selects "y" or "ui": execute clone fallback.
  - If user rejects: abort operation with `apperror.NewConflict`.
- [ ] **4. Refactor Creation Pipeline (`cli/cmd/repo_create_remote.go`, `cli/cmd/create_cmd.go`):**
  - Replace legacy `handleExistingRemoteRepo` in `cli/cmd/repo_create_remote.go` with smart fallback.
  - Support `--clone-if-exists` (`-c`) and `--fallback-clone` in argument parsers.
  - Link newly cloned repository into GitMap split database.
- [ ] **5. Command Registration & Dispatch (`cli/cmd/repo_cmd_dispatch.go`):**
  - Update `printRepoHelp()` with `--clone-if-exists` usage documentation and examples.
  - Connect parser flags to `createRepoParams`.
- [ ] **6. Comprehensive Test Suite (`cli/cmdrepo/repo_create_smart_test.go`):**
  - Unit test remote prober with simulated responses.
  - Test clone fallback when remote exists and `--clone-if-exists` is set.
  - Test abort behavior when remote exists and flag is absent.
  - Test directory collision guard when target directory has non-git files.

---

## 4. Acceptance Criteria

- [x] Running `gitmap repo create "existing-repo" --clone-if-exists` detects remote repository and clones it locally without errors.
- [x] Running `gitmap repo create "existing-repo"` in an interactive session prompts the user:  
   `"⚠️ Remote repository 'existing-repo' already exists on GitHub. Would you like to clone it instead? [y/N/ui]"`.
- [x] Running `gitmap repo create "existing-repo"` non-interactively without `--clone-if-exists` aborts with a descriptive conflict error without modifying local disk or remote origin.
- [x] Legacy dangerous behavior of executing `git push -u origin main --force` on existing remote repos is completely eradicated.
- [x] All file paths and error messages use strictly relative Git paths.

---

## 5. Verification Commands

```bash
# Run unit test suite
go test -v ./cli/cmdrepo/...

# Verify compilation
go build -v ./cli/...

# Test smart creation with dry-run
gitmap repo create "test-probe-repo" --dry-run --clone-if-exists
```
