# Subtask 02: Non-Repo Smart Clone & Fix Engine Implementation Plan

> **Subtask ID:** `subtask-02`  
> **Parent Plan:** `.ai-memory/plans/repository-fix-and-solution-implementation.md`  
> **Target Specification:** `02-spec/21-app/repository-fix-and-solution-implementation/02-component-and-cli-spec.md`  
> **Target Path:** `.ai-memory/plans/subtasks/repository-fix-and-solution-implementation/02-non-repo-smart-clone-and-fix-engine.md`  
> **Assigned Files:**
> - `cli/cloner/pulldiag.go`
> - `cli/cmdpull/pull_remediation_hint.go`
> - `cli/cmdfix/fix_cmd.go`
> - `cli/cmdautofix/fix.go`
> - `cli/cmdremediation/remediation_local.go`  
> **Status:** `PENDING`  

---

## 1. Objectives & Executive Scope

This subtask implements the backend inspection, upstream resolution, dual-option prescription, and automated execution engine for non-repo directories and companion repositories across GitMap:

1. **Non-Repo Folder Detection**:
   - Accurately determine when a filesystem path exists on disk as a directory but lacks a valid `.git` root or worktree file.
   - Prevent misleading classifications such as `missing repository directory` or fallback status checks.
2. **Special Infrastructure Repository Resolution**:
   - Detect companion infrastructure repositories (`repo-cache`, `repo-secrets`, `rc`, `rs`, `repo-storage`).
   - Query SQLite Split-DB (`gitmap-special-repos.db`) and GitHub CLI (`gh repo view`) to resolve upstream remote URLs.
   - Prescribe `gitmap clone <repo>` (Option 1) and in-place `git init` (Option 2) instead of invalid `gitmap status`/`gitmap pull`.
3. **Standard Non-Repo Dual-Option Generation**:
   - For standard folders with known upstream remotes, generate two best-practice sub-item options:
     - **Option 1**: `gitmap clone <repo>` (Clone from remote into target).
     - **Option 2**: `git -C "<path>" init && git -C "<path>" remote add origin <url> && git -C "<path>" fetch` (Initialize & link remote).
   - Render solutions in bright ANSI yellow (`constants.ColorYellow`) formatted as nested tree sub-items.
4. **`gitmap fix` & `gitmap fix --all` Overhaul**:
   - Update `cli/cmdautofix/fix.go` to transparently route repository targets and `--all` to `cmdfix.RunFix`.
   - Update `cli/cmdfix/fix_cmd.go` and `cli/cmdremediation/remediation_local.go` to inspect folder state, recognize non-repo directories, and execute automated clone or initialization workflows.
5. **Coding Guidelines Adherence**:
   - 100% relative Git paths (no absolute paths or `file:///` URIs).
   - Positive boolean naming (`isDirPresent`, `hasGitRoot`, `isValidGitRepo`).
   - Typed `*apperror.AppError` wrappers for all error scenarios.

---

## 2. Technical Architecture & Component Interactions

```
┌────────────────────────────────────────────────────────────────────────┐
│                        cli/cmdautofix/fix.go                           │
│   - Check if args contain repo target or --all                         │
│   - Route to cmdfix.RunFix if not a code-hygiene category              │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                          cli/cmdfix/fix_cmd.go                         │
│   - Inspect local target: cmdremediation.InspectLocalRepoItem          │
│   - If non-repo folder: ExecuteNonRepoFolderFix                        │
│   - If --all: batch iterate dirty repos + non-repo folders             │
└───────────────────┬────────────────────────────────┬───────────────────┘
                    │                                │
                    ▼                                ▼
┌──────────────────────────────────────┐  ┌──────────────────────────────┐
│ cli/cmdremediation/remediation_local │  │ cli/cmdpull/pull_remediation_│
│ - InspectLocalRepoItem updated to    │  │   hint.go                    │
│   recognize FolderStateNonRepo       │  │ - InspectFolderGitState      │
│ - Generate recipes for clone / init  │  │ - ResolveSpecialRepoRemote   │
└──────────────────────────────────────┘  │ - ResolveNonRepoFolderDual...│
                                          └──────────────────────────────┘
```

---

## 3. Step-by-Step Implementation Plan

### Phase 1: Non-Repo Folder Detection Engine (`cli/cmdpull/pull_remediation_hint.go`)

- [ ] **Step 1.1**: Define typed `FolderGitState` enumeration:
  ```go
  type FolderGitState int

  const (
      FolderStateMissing FolderGitState = iota
      FolderStateNonRepo
      FolderStateGitRepo
      FolderStateFileConflict
  )
  ```
- [ ] **Step 1.2**: Implement `InspectFolderGitState(targetPath string) FolderGitState`:
  - Return `FolderStateMissing` if `targetPath` is empty or does not exist on disk.
  - Return `FolderStateFileConflict` if `targetPath` exists but is a regular file.
  - Return `FolderStateGitRepo` if `targetPath/.git` exists as a directory or worktree file.
  - Return `FolderStateNonRepo` if `targetPath` exists as a directory but lacks `.git`.
- [ ] **Step 1.3**: Implement helper `hasValidGitRoot(dir string) bool`:
  - Checks presence of `.git` inside directory with positive boolean naming.
- [ ] **Step 1.4**: Update `ResolvePullErrorDetails(s *PullRepoState) string`:
  - Check `InspectFolderGitState(s.RepoPath)`.
  - When `FolderStateNonRepo`: return `"Directory exists on disk but is not a Git repository (missing .git)"`.
  - When `FolderStateMissing`: return `"Repository directory does not exist on disk"`.

### Phase 2: Pull Diagnostic Classification & Hints (`cli/cloner/pulldiag.go`)

- [ ] **Step 2.1**: Update `appendFailureHints(hints []string, repoDir, output string) []string`:
  - Add check for `not a git repository` in git command stderr:
  ```go
  func hasNotGitRepoFailure(output string) bool {
      lower := strings.ToLower(output)
      return strings.Contains(lower, "not a git repository") ||
          strings.Contains(lower, "fatal: not a git repo")
  }
  ```
- [ ] **Step 2.2**: When `hasNotGitRepoFailure(output)` is true:
  - Check if `repoDir` exists on disk.
  - If directory exists: append hint `"directory is not a git repository; run 'gitmap clone' or 'gitmap fix' to initialize"`.
  - If directory does not exist: append hint `"repository directory missing on disk; run 'gitmap clone'"`.

### Phase 3: Special Repo Upstream Resolution & Dual Options (`cli/cmdpull/pull_remediation_hint.go`)

- [ ] **Step 3.1**: Implement `IsSpecialInfrastructureRepo(name string) bool`:
  - Matches `repo-cache`, `repo-secrets`, `rc`, `rs`, `repo-storage`, `cache`, `storage`.
- [ ] **Step 3.2**: Implement `CanonicalizeSpecialRepoKey(name string) (string, string)`:
  - Maps aliases to canonical key (`repo-cache` or `repo-secrets`) and short key (`rc` or `rs`).
- [ ] **Step 3.3**: Implement `ResolveSpecialRepoRemote(name string) (string, bool)`:
  - Query SQLite Split-DB (`store.OpenSpecialReposSplitDB`).
  - Fall back to GitHub CLI lookup (`gh repo view <canonicalKey> --json url -q .url`).
  - Cache discovered URL in Split-DB for offline operation.
- [ ] **Step 3.4**: Implement `ResolveNonRepoFolderDualHints(repoDir, repoName string) (string, string, string, string)`:
  - For special repos:
    - Option 1 (Clone from Remote): `gitmap clone <canonicalKey>`
    - Option 2 (Initialize Local Repo): `git -C "<path>" init`
  - For standard repos with known remote:
    - Option 1 (Clone from Remote): `gitmap clone <repoName>`
    - Option 2 (Init & Link Remote): `git -C "<path>" init && git -C "<path>" remote add origin <url> && git -C "<path>" fetch`
  - For standard repos without known remote:
    - Option 1 (Initialize Local Repo): `git -C "<path>" init`
    - Option 2 (Remove from Registry): `gitmap rm --db-only <repoName>`
- [ ] **Step 3.5**: Update `ResolveStructuredRemediation(s *PullRepoState) StructuredRemediation`:
  - Route `FolderStateNonRepo` directly to `ResolveNonRepoFolderDualHints`.
  - Set diagnostic property to `"gitmap pe -t " + repoName`.

### Phase 4: Local Remediation Item Expansion (`cli/cmdremediation/remediation_local.go`)

- [ ] **Step 4.1**: Update `InspectLocalRepoItem(query string) (*RemediationItem, string)`:
  - When `.git` does not exist, do NOT return `nil, ""` blindly.
  - Check `InspectFolderGitState(absPath)`.
  - When `FolderStateNonRepo`:
    - Generate non-repo remediation recipes:
      - Recipe 1: Clone from remote (if remote URL known).
      - Recipe 2: Init git repository in-place and link remote.
    - Return `&RemediationItem{...}` with `SummaryReason: "Not a git repository (missing .git)"`.

### Phase 5: CLI Routing & Fix Command Overhaul (`cli/cmdautofix/fix.go` & `cli/cmdfix/fix_cmd.go`)

- [ ] **Step 5.1**: Overhaul `cli/cmdautofix/fix.go`:
  - When `args[0]` is NOT a valid code-hygiene subcommand (`encoding`, `newlines`, etc.):
    - Check if `args[0]` is a repo query, `--all`, `-a`, `all`, or special repo name.
    - If yes: delegate directly to `cmdfix.RunFix(args, "")`.
    - If no (unknown flag / help): render help text.
- [ ] **Step 5.2**: Equip `cli/cmdfix/fix_cmd.go` with Non-Repo Handler:
  - In `runFixDirect`:
    - When `cmdremediation.InspectLocalRepoItem` returns a non-repo item:
      - Execute smart clone or in-place init.
  - In `runFixAll`:
    - Process all pending remediation items, supporting both dirty git repositories and non-repo directories.
- [ ] **Step 5.3**: Implement `executeNonRepoFixRecipe(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error`:
  - If recipe is clone:
    - If target directory is empty: execute `git clone <remoteURL> <repoPath>`.
    - If target directory has files: prompt or execute in-place init & link remote.
  - If recipe is init:
    - Execute `git -C "<repoPath>" init`.
    - If remote URL exists: execute `git -C "<repoPath>" remote add origin <remoteURL>`.

---

## 4. Detailed Data Models & Go Signatures

### 4.1 Folder Inspection Contract

```go
package cmdpull

import (
	"os"
	"path/filepath"
	"strings"
)

// FolderGitState classifies directory presence and git initialization.
type FolderGitState int

const (
	FolderStateMissing FolderGitState = iota
	FolderStateNonRepo
	FolderStateGitRepo
	FolderStateFileConflict
)

// InspectFolderGitState inspects target directory and verifies .git presence.
func InspectFolderGitState(targetPath string) FolderGitState {
	if len(strings.TrimSpace(targetPath)) == 0 {
		return FolderStateMissing
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		return FolderStateMissing
	}
	if !info.IsDir() {
		return FolderStateFileConflict
	}
	gitPath := filepath.Join(targetPath, ".git")
	gitInfo, gitErr := os.Stat(gitPath)
	if gitErr == nil && (gitInfo.IsDir() || (!gitInfo.IsDir() && gitInfo.Size() > 0)) {
		return FolderStateGitRepo
	}
	return FolderStateNonRepo
}
```

### 4.2 Special Repository Upstream Lookup

```go
package cmdpull

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IsSpecialInfrastructureRepo reports whether name is a companion repository.
func IsSpecialInfrastructureRepo(name string) bool {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	switch cleaned {
	case "rc", "repo-cache", "repo-storage", "cache", "storage",
		"rs", "repo-secrets", "secrets", "vault":
		return true
	default:
		return false
	}
}

// CanonicalizeSpecialRepoKey returns canonical name and short key.
func CanonicalizeSpecialRepoKey(name string) (string, string) {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	switch cleaned {
	case "rc", "repo-cache", "repo-storage", "cache", "storage":
		return "repo-cache", "rc"
	default:
		return "repo-secrets", "rs"
	}
}

// ResolveSpecialRepoRemote resolves remote upstream URL from Split-DB or GitHub CLI.
func ResolveSpecialRepoRemote(name string) (string, bool) {
	canonicalKey, shortKey := CanonicalizeSpecialRepoKey(name)
	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, err := store.GetSpecialRepoByKey(db.Conn(), canonicalKey); err == nil && len(rec.RemoteURL) > 0 {
			return rec.RemoteURL, true
		}
	}

	cmd := exec.Command("gh", "repo", "view", canonicalKey, "--json", "url", "-q", ".url")
	if out, err := cmd.Output(); err == nil {
		url := strings.TrimSpace(string(out))
		if len(url) > 0 {
			if db, err := store.OpenSpecialReposSplitDB(); err == nil {
				defer db.Close()
				_ = store.UpdateSpecialRepoPath(db.Conn(), shortKey, "", url)
			}
			return url, true
		}
	}
	return "", false
}
```

---

## 5. Verification & Quality Gates

| Gate ID | Target Component | Acceptance Criteria |
| :---: | :--- | :--- |
| **VG-SUB-01** | `InspectFolderGitState` | Accurately returns `FolderStateNonRepo` when directory exists without `.git`. |
| **VG-SUB-02** | `IsSpecialInfrastructureRepo` | Correctly recognizes `rc`, `rs`, `repo-cache`, `repo-secrets`. |
| **VG-SUB-03** | `ResolveSpecialRepoRemote` | Retrieves cached remote URL from Split-DB or probes GitHub CLI. |
| **VG-SUB-04** | Dual Options | Returns Option 1 (Clone from Remote) and Option 2 (Initialize Local Repo). |
| **VG-SUB-05** | `cmdautofix` Router | Passes `gitmap fix repo-cache` and `gitmap fix --all` to `cmdfix.RunFix`. |
| **VG-SUB-06** | `gitmap fix <repo>` | Executes smart clone or in-place git init on non-repo folders without crashing. |
| **VG-SUB-07** | `gitmap fix --all` | Successfully processes batch resolution across all failed repositories. |
| **VG-SUB-08** | Code Style | Zero absolute paths, positive boolean names, and typed `*apperror.AppError`. |

---

## 6. Execution Traceability

- **Prerequisite Plans**: None (Subtask 01 handles UI rendering; Subtask 02 handles detection and engine).
- **Dependent Files**:
  - `cli/cloner/pulldiag.go`
  - `cli/cmdpull/pull_remediation_hint.go`
  - `cli/cmdfix/fix_cmd.go`
  - `cli/cmdautofix/fix.go`
  - `cli/cmdremediation/remediation_local.go`
- **Output Artifacts**:
  - Passing unit tests in `cli/cmdpull/pull_remediation_hint_test.go` and `cli/cmdfix/fix_cmd_test.go`.
