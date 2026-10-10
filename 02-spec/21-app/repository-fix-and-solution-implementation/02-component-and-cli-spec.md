# Component and CLI Specification: Non-Repo Folder Detection, Special Repo Resolution & `gitmap fix` Overhaul

> **/goal** Specify the architectural components, CLI interfaces, domain models, and execution workflows for detecting non-repo folders, resolving special infrastructure repositories (`repo-cache`, `repo-secrets`, `rc`, `rs`), prescribing ANSI yellow dual-option subtree solutions, and overhauling `gitmap fix <repo>` / `gitmap fix --all`.
> **/learn** Eliminate invalid suggestions (such as `gitmap status` or `gitmap pull` on directories lacking `.git`), establish upstream remote resolution from Split-DB and GitHub CLI, format dual-option subtree remediation in bright yellow, and equip `gitmap fix` to execute smart clone and initialization workflows.

- **Specification Slug:** `repository-fix-and-solution-implementation`
- **Spec Document:** `02-spec/21-app/repository-fix-and-solution-implementation/02-component-and-cli-spec.md`
- **Spec Status:** `APPROVED`
- **Architecture Reference:** `02-spec/21-app/repository-fix-and-solution-implementation/01-architecture-spec.md`
- **Parent Plan:** `.ai-memory/plans/repository-fix-and-solution-implementation.md`
- **Subtask Plan Reference:** `.ai-memory/plans/subtasks/repository-fix-and-solution-implementation/02-non-repo-smart-clone-and-fix-engine.md`
- **Target Version:** `v6.529.0`

---

## 1. Executive Summary & Root Cause Analysis

### 1.1 Problem Statement
When `gitmap pull-all` (alias `gitmap pa`) executes parallel fetches across workspace repositories, it occasionally encounters target directories on disk that are not valid Git repositories—most frequently companion or infrastructure repositories such as `repo-cache` (`rc`) or `repo-secrets` (`rs`). In these scenarios, the pull process fails with `fatal: not a git repository (or any of the parent directories): .git`.

Historically, GitMap's diagnostic engine exhibited three severe defects:
1. **Misclassification as Missing Directory or Fallback**: The error diagnostic failed to distinguish between a directory that does not exist versus an existing directory that simply lacks a `.git` root. Consequently, it fell back to generic remediation hints:
   ```text
   Next Step: gitmap status repo-cache or gitmap fix repo-cache
   ```
2. **Invalid Suggested Commands**: Running `gitmap status repo-cache` or `gitmap pull repo-cache` on a directory without a `.git` folder immediately crashes or repeats the fatal git error, providing zero forward progress for the user.
3. **Flawed `gitmap fix` Execution**: `gitmap fix repo-cache` inspected only standard git dirty states via `gitutil.InspectDirtyState`. Finding no `.git` directory, it failed with `repository "repo-cache" not found or not a git repository`. Furthermore, the code-hygiene router in `cmdautofix.RunFixCmd` intercepted repository fix requests and aborted with `unknown subcommand "repo-cache"`.

### 1.2 Architectural Objectives
This component and CLI specification formalizes the solution across four foundational pillars:
- **Pillar 1: Accurate Non-Repo Folder State Detection**: Rigorously verify whether a path exists, is a directory, and contains a valid `.git` root or worktree pointer.
- **Pillar 2: Special Infrastructure Repository Upstream Resolution**: Automatically identify special companion repositories (`repo-cache`, `repo-secrets`, `rc`, `rs`) and resolve upstream remote URLs via SQLite Split-DB (`SpecialRepository` table) or GitHub CLI (`gh repo view`).
- **Pillar 3: Dual-Option Subtree Yellow Remediation**: Replace flat colon-separated suggestion strings with structured, nested sub-items rendered in bright ANSI yellow (`constants.ColorYellow`), providing Option 1 (Clone from remote) and Option 2 (Initialize local git repo & link remote).
- **Pillar 4: Unified `gitmap fix` Overhaul**: Route repository fix queries seamlessly to `cmdfix.RunFix`, equip `gitmap fix <repo>` with smart clone/init execution for non-repo directories, and provide a single batch resolution command `gitmap fix --all` at the end of the pull failure summary.

---

## 2. Component Architecture & Responsibility Matrix

```
┌────────────────────────────────────────────────────────────────────────┐
│                        CLI Dispatcher Layer                            │
│           (cli/cmdautofix/fix.go <-> cli/cmdfix/fix_cmd.go)            │
└───────────────────┬────────────────────────────────┬───────────────────┘
                    │                                │
                    ▼                                ▼
       ┌────────────────────────┐       ┌────────────────────────┐
       │   Code-Hygiene Flow    │       │  Repository Fix Flow   │
       │ (encoding, gofmt, etc) │       │ (gitmap fix <repo|all>)│
       └────────────────────────┘       └────────────┬───────────┘
                                                     │
                                                     ▼
┌────────────────────────────────────────────────────────────────────────┐
│                   Non-Repo Folder Detection Engine                     │
│                  (cli/cmdpull/pull_remediation_hint.go)                │
│                                                                        │
│   - InspectFolderGitState(path) -> FolderGitState                      │
│   - isDirPresent(path) && !hasGitRoot(path)                            │
└───────────────────┬────────────────────────────────┬───────────────────┘
                    │                                │
                    ▼                                ▼
┌───────────────────────────────┐       ┌────────────────────────────────┐
│ Special Repo Resolver         │       │ Standard Repo Matcher          │
│ (cli/store/special_repos_...  │       │ (Split-DB / Workspace Scan)    │
│  & cli/cmdpull/pull_remed...) │       │                                │
│ - Match rc, rs, repo-cache... │       │ - Match workspace registry     │
│ - Resolve upstream remote URL │       │ - Resolve remote tracking URL  │
└───────────────┬───────────────┘       └────────────┬───────────────────┘
                │                                    │
                └─────────────────┬──────────────────┘
                                  │
                                  ▼
┌────────────────────────────────────────────────────────────────────────┐
│                   Subtree Remediation Generator                        │
│               (cli/cmdpull/pull_remediation_hint.go)                   │
│                                                                        │
│   ├── Reason: directory exists but is not a Git repository             │
│   ├── Solutions:                                                       │
│   │   ├── Option 1 (Clone from Remote): gitmap clone <repo>            │
│   │   └── Option 2 (Init & Link Remote): git -C "<path>" init && ...   │
│   └── Diagnostic: gitmap pe -t <repo>                                  │
└─────────────────────────────────┬──────────────────────────────────────┘
                                  │
                                  ▼
┌────────────────────────────────────────────────────────────────────────┐
│                  Smart Clone & Fix Execution Engine                    │
│                      (cli/cmdfix/fix_execute.go)                       │
│                                                                        │
│   - Empty folder: git clone directly into target                       │
│   - Populated folder: git init in-place, link remote origin, fetch     │
│   - Batch resolution: gitmap fix --all iterates all pending fixes      │
└────────────────────────────────────────────────────────────────────────┘
```

### 2.1 Component Responsibility Matrix

| Component Name | Source File(s) | Primary Responsibility | Error Handling & Wrappers |
| :--- | :--- | :--- | :--- |
| **Folder State Inspector** | `cli/cmdpull/pull_remediation_hint.go` | Detects directory presence and validates `.git` folder/worktree existence. | Returns typed `FolderGitState` enum; returns positive booleans. |
| **Special Repo Resolver** | `cli/store/special_repos_split_db.go`, `cli/cmdpull/pull_remediation_hint.go` | Identifies `repo-cache`, `repo-secrets`, `rc`, `rs` and resolves upstream remote URLs. | `apperror.WrapSimple(err, "resolve special repo remote")` |
| **Pull Diagnostic Classifier** | `cli/cloner/pulldiag.go` | Classifies pull failure strings and isolates `not a git repository` conditions. | Sanitizes raw git stderr; appends clean diagnostic hints. |
| **Subtree Solution Formatter** | `cli/cmdpull/pull_remediation_hint.go`, `cli/cmdpull/pull_efficient_render.go` | Formats dual-option solutions in bright ANSI yellow as nested sub-items. | Strict formatting without colons; preserves indent hierarchy. |
| **Fix Command Router** | `cli/cmdautofix/fix.go` | Distinguishes code-hygiene subcommands from repo names/flags; delegates to `cmdfix`. | Transparent delegation; handles `--help` and usage errors. |
| **Non-Repo Fix Executor** | `cli/cmdfix/fix_cmd.go`, `cli/cmdfix/fix_execute.go` | Executes automated clone or init/remote linkage for non-repo directories. | `apperror.New("fix", "E_CLONE_FAILED", ...)` |

---

## 3. Detailed Component Specifications

### 3.1 Non-Repo Folder Detection Engine

#### 3.1.1 State Classification Enum
Folder state must be represented as a discrete enumeration using positive status definitions:

```go
package cmdpull

// FolderGitState represents the filesystem and git validation state of a directory path.
type FolderGitState int

const (
	// FolderStateMissing indicates the target path does not exist on disk.
	FolderStateMissing FolderGitState = iota

	// FolderStateNonRepo indicates the path exists as a directory but lacks a valid .git root.
	FolderStateNonRepo

	// FolderStateGitRepo indicates the path exists as a directory and contains a valid .git root.
	FolderStateGitRepo

	// FolderStateFileConflict indicates the path exists on disk but is a regular file, not a directory.
	FolderStateFileConflict
)
```

#### 3.1.2 Folder Inspection Logic
The inspector enforces positive boolean naming conventions and returns discrete state:

```go
// InspectFolderGitState inspects a path and determines its folder and git readiness.
func InspectFolderGitState(targetPath string) FolderGitState {
	if len(strings.TrimSpace(targetPath)) == 0 {
		return FolderStateMissing
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return FolderStateMissing
		}
		return FolderStateMissing
	}

	if !info.IsDir() {
		return FolderStateFileConflict
	}

	gitPath := filepath.Join(targetPath, ".git")
	gitInfo, gitErr := os.Stat(gitPath)
	if gitErr == nil && (gitInfo.IsDir() || isGitWorktreeFile(gitInfo)) {
		return FolderStateGitRepo
	}

	return FolderStateNonRepo
}

func isGitWorktreeFile(info os.FileInfo) bool {
	return !info.IsDir() && info.Size() > 0
}
```

### 3.2 Special Infrastructure Repository Resolution

#### 3.2.1 Identification of Special Repositories
GitMap recognizes companion infrastructure repositories by canonical key or accepted alias:
- **Cache / Storage**: `repo-cache`, `rc`, `repo-storage`, `cache`, `storage`
- **Secrets / Vault**: `repo-secrets`, `rs`, `secrets`, `vault`

```go
// IsSpecialInfrastructureRepo reports whether a repo name or alias is a companion repo.
func IsSpecialInfrastructureRepo(nameOrAlias string) bool {
	cleaned := strings.ToLower(strings.TrimSpace(nameOrAlias))
	switch cleaned {
	case "rc", "repo-cache", "repo-storage", "cache", "storage",
		"rs", "repo-secrets", "secrets", "vault":
		return true
	default:
		return false
	}
}

// CanonicalizeSpecialRepoKey returns the primary key ("repo-cache" or "repo-secrets") and short key ("rc" or "rs").
func CanonicalizeSpecialRepoKey(nameOrAlias string) (string, string) {
	cleaned := strings.ToLower(strings.TrimSpace(nameOrAlias))
	switch cleaned {
	case "rc", "repo-cache", "repo-storage", "cache", "storage":
		return "repo-cache", "rc"
	default:
		return "repo-secrets", "rs"
	}
}
```

#### 3.2.2 3-Tier Upstream Remote Resolution
When prescribing a clone command for a special repository, the system resolves its remote URL using a deterministic 3-tier lookup:
1. **Tier 1 (Split-DB Lookup)**: Query `gitmap-special-repos.db` table `SpecialRepository` where `ShortKey = ? OR RepoKey = ?`. If `RemoteURL` is non-empty, use it.
2. **Tier 2 (GitHub CLI Probe)**: Execute `gh repo view <repoName> --json url -q .url`. If execution succeeds and returns a valid URL, cache it in Split-DB and use it.
3. **Tier 3 (Workspace Sibling Probe)**: Inspect siblings in the active workspace root or configuration for remote origin patterns (e.g., `git@github.com:<org>/<repoName>.git`).

```go
// ResolveSpecialRepoRemote resolves the upstream git URL for a special companion repository.
func ResolveSpecialRepoRemote(repoName string) (string, bool) {
	canonicalKey, shortKey := CanonicalizeSpecialRepoKey(repoName)

	// Tier 1: SQLite Split-DB
	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, err := store.GetSpecialRepoByKey(db.Conn(), canonicalKey); err == nil && len(rec.RemoteURL) > 0 {
			return rec.RemoteURL, true
		}
	}

	// Tier 2: GitHub CLI lookup
	cmd := exec.Command("gh", "repo", "view", canonicalKey, "--json", "url", "-q", ".url")
	if out, err := cmd.Output(); err == nil {
		url := strings.TrimSpace(string(out))
		if len(url) > 0 {
			// Persist resolved URL back into Split-DB for fast offline lookup
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

### 3.3 Dual-Option Subtree Remediation Engine

#### 3.3.1 Domain Models
The remediation models are defined with clean JSON serialization and positive naming:

```go
// RemediationOption represents an actionable fix choice.
type RemediationOption struct {
	OptionNumber int    `json:"option_number"`
	Title        string `json:"title"`
	Command      string `json:"command"`
}

// StructuredRemediation captures error reason, dual-option solutions, and diagnostics.
type StructuredRemediation struct {
	Reason       string              `json:"reason"`
	Solutions    []RemediationOption `json:"solutions"`
	Diagnostic   string              `json:"diagnostic,omitempty"`
}
```

#### 3.3.2 Non-Repo Dual Remediation Hints
When `InspectFolderGitState` returns `FolderStateNonRepo`, the engine produces two best-practice solutions:
- **Special Repositories (`repo-cache`, `repo-secrets`)**:
  - Option 1 (Clone from Remote): `gitmap clone <repoName>`
  - Option 2 (Initialize in-place): `git -C "<path>" init && gitmap settings set special_repos.<key>_name <repoName>`
- **Standard Repositories with Known Remote**:
  - Option 1 (Clone from Remote into Target): `gitmap clone <repoName>` (or `git clone <remoteURL> "<path>"`)
  - Option 2 (Initialize Local Git Repo & Link Remote): `git -C "<path>" init && git -C "<path>" remote add origin <remoteURL> && git -C "<path>" fetch`
- **Standard Repositories without Known Remote**:
  - Option 1 (Initialize Local Repo): `git -C "<path>" init`
  - Option 2 (Remove Stale Folder / Unregister): `gitmap rm --db-only <repoName>`

```go
// ResolveNonRepoFolderDualHints builds dual remediation options for non-repo directories.
func ResolveNonRepoFolderDualHints(repoDir, repoName string) (string, string, string, string) {
	name := resolveEffectiveRepoName(repoName, repoDir)
	cleanDir := filepath.ToSlash(filepath.Clean(repoDir))

	if IsSpecialInfrastructureRepo(name) {
		canonicalKey, _ := CanonicalizeSpecialRepoKey(name)
		cloneCmd := fmt.Sprintf("gitmap clone %s", canonicalKey)
		initCmd := fmt.Sprintf("git -C \"%s\" init", cleanDir)
		return "Clone from Remote", cloneCmd, "Initialize Local Repo", initCmd
	}

	remoteURL, hasRemote := resolveKnownRemoteURL(name, repoDir)
	if hasRemote {
		cloneCmd := fmt.Sprintf("gitmap clone %s", name)
		linkCmd := fmt.Sprintf("git -C \"%s\" init && git -C \"%s\" remote add origin %s && git -C \"%s\" fetch",
			cleanDir, cleanDir, remoteURL, cleanDir)
		return "Clone from Remote", cloneCmd, "Init & Link Remote", linkCmd
	}

	initCmd := fmt.Sprintf("git -C \"%s\" init", cleanDir)
	removeCmd := fmt.Sprintf("gitmap rm --db-only %s", name)
	return "Initialize Local Repo", initCmd, "Remove from Registry", removeCmd
}
```

#### 3.3.3 Visual Subtree Layout Specification
The terminal rendering must strictly adhere to the nested subtree hierarchy:
- Root level: Repository failure banner with status badge (`[FAIL]`).
- Sub-item 1: `├── Reason: <clear diagnosis>` (in dim or white text).
- Sub-item 2: `├── Solutions:` (group header).
- Sub-items 2.1 & 2.2: `│   ├── Option 1 (<Title>): <command>` and `│   └── Option 2 (<Title>): <command>` (in bright ANSI yellow: `\033[1;33m`).
- Sub-item 3: `└── Diagnostic: gitmap pe -t <repo>` (in cyan or dim text).

```text
  ● repo-cache [FAIL]
    ├── Reason: Directory exists on disk but is not a Git repository (missing .git)
    ├── Solutions:
    │   ├── Option 1 (Clone from Remote): gitmap clone repo-cache
    │   └── Option 2 (Initialize Local Repo): git -C "repo-cache" init
    └── Diagnostic: gitmap pe -t repo-cache
```

#### 3.3.4 Unified Batch Resolution Command
At the very end of the failed repositories summary in `gitmap pull-all`, the renderer MUST output a single unified batch command in bright ANSI yellow:

```text
💡 Unified Resolution: gitmap fix --all (or: gitmap pull-fix --all)
```

---

## 4. `gitmap fix` Command Overhaul & Execution Engine

### 4.1 CLI Routing Architecture (`cli/cmdautofix` vs `cli/cmdfix`)
Historically, `cli/cmd/rootutility.go` routed `gitmap fix` directly to `cmdautofix.RunFixCmd`. When a user ran `gitmap fix repo-cache` or `gitmap fix --all`, `cmdautofix` rejected it because `repo-cache` is not a code-hygiene category.

The overhauled routing specification establishes clean boundary detection:

```
                  ┌──────────────────────────────┐
                  │      gitmap fix [args...]    │
                  └──────────────┬───────────────┘
                                 │
                                 ▼
                     Is args empty or --help?
                    ┌────────────┴────────────┐
                   YES                        NO
                    │                         │
                    ▼                         ▼
         Render Fix Category      Is args[0] in code-hygiene
         & Repo Fix Help          categories (encoding, gofmt...)?
                                 ┌────────────┴────────────┐
                                YES                        NO
                                 │                         │
                                 ▼                         ▼
                        Execute cmdautofix        Delegate to cmdfix
                        Code-Hygiene Flow         Repository Fix Flow
```

#### 4.1.1 Routing Predicate
```go
// IsCodeHygieneCategory reports whether a token matches an autofix linter category.
func IsCodeHygieneCategory(token string) bool {
	cleaned := strings.ToLower(strings.TrimSpace(token))
	switch cleaned {
	case "encoding", "newlines", "naming", "paths", "gofmt",
		"misspell", "markdown", "guidelines":
		return true
	default:
		return false
	}
}

// IsRepositoryFixTarget reports whether arguments represent a repo fix request.
func IsRepositoryFixTarget(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(strings.TrimSpace(args[0]))
	if first == "--all" || first == "-a" || first == "all" {
		return true
	}
	if IsCodeHygieneCategory(first) {
		return false
	}
	return true
}
```

### 4.2 Single Repository Fix Workflow (`gitmap fix <repo>`)

When invoked for a specific repository:
1. **Target Inspection**:
   - Determine local path from Split-DB or workspace relative path.
   - Run `InspectFolderGitState(targetPath)`.
2. **Branch A: Valid Git Repo (`FolderStateGitRepo`)**:
   - Inspect dirty state (`gitutil.InspectDirtyState(targetPath)`).
   - If dirty, apply requested recipe (stash / wip / discard).
   - If clean, inform user repository is clean.
3. **Branch B: Missing Directory (`FolderStateMissing`)**:
   - Resolve remote URL.
   - Execute `git clone <remoteURL> <targetPath>`.
4. **Branch C: Non-Repo Directory (`FolderStateNonRepo`)**:
   - Inspect folder contents (`os.ReadDir(targetPath)`).
   - **Case C1 (Empty Directory)**:
     - Clone directly into the target directory (`git clone <remoteURL> .` inside directory or remove empty folder and clone).
   - **Case C2 (Non-Empty Directory / Scratch Files Present)**:
     - Initialize Git repository: `git -C "<targetPath>" init`.
     - Resolve remote URL from Split-DB or GitHub CLI.
     - If remote URL exists:
       - Add remote: `git -C "<targetPath>" remote add origin <remoteURL>`.
       - Fetch remote: `git -C "<targetPath>" fetch origin`.
       - Link tracking: set default upstream branch if remote branch exists.
     - If no remote URL exists (local scratch repository):
       - Stage and create initial commit: `git -C "<targetPath>" add -A && git -C "<targetPath>" commit -m "init: local companion repository"`.
   - Update Split-DB record marking repository as active/initialized.

```go
// ExecuteNonRepoFolderFix fixes a directory that is not a valid git repository.
func ExecuteNonRepoFolderFix(repoName, repoPath string, isAutoAccept bool) error {
	cleanPath := filepath.Clean(repoPath)
	state := cmdpull.InspectFolderGitState(cleanPath)

	if state == cmdpull.FolderStateGitRepo {
		return nil
	}

	remoteURL, hasRemote := cmdpull.ResolveSpecialRepoRemote(repoName)
	if !hasRemote {
		remoteURL, hasRemote = resolveKnownRemoteURL(repoName, cleanPath)
	}

	entries, readErr := os.ReadDir(cleanPath)
	isEmpty := readErr == nil && len(entries) == 0

	if isEmpty && hasRemote {
		return executeDirectCloneIntoEmptyDir(remoteURL, cleanPath)
	}

	return executeInPlaceGitInitAndLink(repoName, cleanPath, remoteURL)
}
```

### 4.3 Batch Repository Fix Workflow (`gitmap fix --all`)

When invoked with `--all`:
1. **Load Remediation State**:
   - Read pending failed pull records from `cmdremediation.LoadRemediationState()`.
   - Also scan tracked workspace repositories for non-repo directories.
2. **Iterate Remediation Items**:
   - For each item, inspect `InspectFolderGitState`:
     - If `FolderStateNonRepo`: execute non-repo fix (clone or init+link).
     - If `FolderStateGitRepo` with dirty changes: execute default recipe (`stash` or user-specified action).
     - If `FolderStateMissing`: execute clone.
3. **Progress Reporting & Visual Summary**:
   - Render Catppuccin progress line per repository:
     ```text
     ✔ repo-cache: Cloned from https://github.com/org/repo-cache.git
     ✔ repo-secrets: Initialized git repository & linked remote origin
     ✔ gitmap: Stashed 2 modified file(s)
     ```
   - Print final success summary in green:
     ```text
     ✓ All 3 repository(ies) remediated successfully.
     ```

---

## 5. CLI Interface & Flag Specification

### 5.1 Command Signatures

```text
gitmap fix <repo> [action] [flags]
gitmap fix --all [action] [flags]
gitmap fix <category> [path] [flags]
```

### 5.2 Flag Definitions

| Flag | Shorthand | Type | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `--all` | `-a` | boolean | `false` | Batch remediate all pending failed and dirty repositories. |
| `--yes` | `-y` | boolean | `false` | Apply recommended remediation non-interactively without prompting. |
| `--prompt` | `-p` | boolean | `false` | Interactively prompt for each failed repository's remediation recipe. |
| `--clone` | `-c` | boolean | `false` | Force clone action for non-repo folders instead of in-place init. |
| `--init` | `-i` | boolean | `false` | Force in-place `git init` for non-repo folders instead of clone. |
| `--json` | — | boolean | `false` | Output machine-readable JSON remediation report. |

### 5.3 Exit Code Contract

| Exit Code | Classification | Condition |
| :---: | :--- | :--- |
| `0` | **Success** | All targets clean, or all requested fixes applied cleanly without errors. |
| `1` | **Warning / Incomplete** | User declined prompt (`[y/N]`), or report-only findings remain. |
| `2` | **Execution Error** | Filesystem permission error, network clone failure, or invalid argument. |

---

## 6. Error Management & AppError Taxonomy

All failure paths MUST wrap errors in typed `*apperror.AppError` instances with structured context:

| Error Code | Classification | Condition | Structured Context Fields |
| :--- | :--- | :--- | :--- |
| `E_NON_REPO_FOLDER` | Validation | Directory exists on disk but is not a valid git repository. | `path`, `repoName`, `hasRemote` |
| `E_SPECIAL_REPO_UNRESOLVED` | Lookup | Special repository remote URL could not be resolved. | `repoKey`, `shortKey` |
| `E_CLONE_FAILED` | Execution | Git clone command failed during remediation. | `repoName`, `targetDir`, `remoteUrl`, `stderr` |
| `E_INIT_FAILED` | Execution | Git init or remote add command failed during remediation. | `path`, `command`, `stderr` |
| `E_AMBIGUOUS_FIX` | Usage | Multiple repositories need remediation without `--all` or target name. | `pendingCount` |

```go
func newNonRepoFolderError(repoName, repoPath string, hasRemote bool) *apperror.AppError {
	return apperror.New("fix", "E_NON_REPO_FOLDER", map[string]any{
		"repoName":  repoName,
		"path":      repoPath,
		"hasRemote": hasRemote,
		"msg":       fmt.Sprintf("Directory %q exists but is not a git repository", repoPath),
	})
}
```

---

## 7. Quality Gates & Verification Matrix

| Gate ID | Area | Verification Criterion | Validation Command / Proof |
| :---: | :--- | :--- | :--- |
| **VG-01** | Non-Repo Detection | `InspectFolderGitState` returns `FolderStateNonRepo` for folders missing `.git`. | Unit test in `cli/cmdpull/pull_remediation_hint_test.go` |
| **VG-02** | Special Repo Resolver | `repo-cache` & `repo-secrets` resolve upstream remote URLs without errors. | Unit test verifying Split-DB and GitHub fallback |
| **VG-03** | Yellow Subtree | Solution options render in ANSI yellow (`constants.ColorYellow`) without colons. | Output inspection in `cli/cmdpull/pull_efficient_render_test.go` |
| **VG-04** | Dual Options | Non-repo folders display Option 1 (Clone) and Option 2 (Init & Link). | Unit test validating `ResolveStructuredRemediation` |
| **VG-05** | Single Fix Command | `gitmap fix repo-cache` detects non-repo folder and executes clone/init. | Integration test in `cli/cmdfix/fix_cmd_test.go` |
| **VG-06** | Batch Fix Command | `gitmap fix --all` processes dirty repos and non-repo folders in one run. | Integration test verifying batch remediation |
| **VG-07** | Relative Git Paths | Zero absolute filesystem paths or `file:///` URIs across code and docs. | Verification via `gitmap aum search` |
| **VG-08** | Positive Booleans | All new functions and fields use positive booleans (`isValid`, `hasGitRoot`). | Code review & linter check |
