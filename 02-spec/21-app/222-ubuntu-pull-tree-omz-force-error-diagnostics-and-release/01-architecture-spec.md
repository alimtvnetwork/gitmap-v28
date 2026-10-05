# Architecture Specification: Scanner Exclusions, Force-Include Flag, and Failure Tree Subtree Rendering

## 1. User Request (Verbatim)

> "In the Ubuntu system, I found these issues, but also at the same time, there is no stack trace so that you can remedy and find out what went wrong, and there is no solution. One more point. When we are in the Ubuntu machine, make sure that omizssh, we do not pick out automatically. And make sure that it is by default, automatically, it should not be catched by scan or something. But if the user adds it by default, manually add it to the scan, there should be specific force command or probably add a force for that repository mentioning this, and there should be command example for this. Then and only then this should be added. So make sure that you have a remedy for this. And also, I want you to connect back to the Ubuntu machine and try to see the root cause of it, why it happened, how it happened. Create scripts inside the repo secrets folder. Make sure that you implement those as well. And then you make sure that the solutions are integrated into the `gitmap`, and we make a final release bump and make a release. Also, write the root cause analysis, why it happened, how it happened, and also the, let's say, in the tree view, when you say next step, the `gitmap`, that actually go into another tree view, actually. So here, the next step, the option one, option two, that should be a subtree. Option one, option two. And next step, to clone this repository. The third section, the solution that you have provided, it feels a little bit buggy. So try to find the root cause. At the end, let me know the root cause, why it happened, how it happened, what was the solution, and `gitmap` should be able to fix these things. And also, if not, then it should give an error log and also hint the system to see the error log using error log methods. And you should also always do the error logs from your system and code so that it can be used for tracing, auditing, and things like that."

---

## 2. Executive Summary & Architecture Objectives

This specification defines the canonical architecture for two foundational components in GitMap:
1. **Scanner Exclusions & Explicit Force Override**:
   - Universal exclusion of Oh-My-Zsh directories (`.oh-my-zsh`, `oh-my-zsh`, `ohmyzsh`, `.ohmyzsh`, `omz`, `.omz`, `omizssh`) alongside dependency and build caches (`node_modules`, `vendor`, `.cache`, `.venv`, etc.).
   - Multi-layer defense across `constants.DefaultScanExcludeDirs`, `scanner.buildExcludeSet`, `model.DefaultConfig`, and directory traversal hooks (`fsutil.DiscoverChildGitRepos` and `fsutil.DiscoverTopLevelGitRepos`).
   - Controlled bypass via `--force-include <dir>` flag on `gitmap scan` and `--force` on `gitmap add`, complete with documented CLI command examples.
2. **Failure Tree Subtree Hierarchy & Context-Aware Dual Remediation**:
   - Structured visual tree rendering using box-drawing characters (`├──`, `└──`, `│   `) for pull failures in `cli/cmdpull/pull_efficient_render.go`.
   - Transformation of remediation options into a first-class nested subtree below `Next Step:`.
   - Eradication of broken commands (`git status`, `git pull`) when a repository directory does not exist on disk, replacing them with actionable recovery dual options: Option 1 (`gitmap clone <repo>`) and Option 2 (`gitmap rm --db-only <repo>`).
   - Explicit terminal guidance directing users to stack trace inspection via `gitmap pull-error <repo>` (`gitmap pe`) and `gitmap pe -t`.

---

## 3. High-Level System Data Flow

```
                                 +--------------------------------+
                                 |    Filesystem / CLI Trigger    |
                                 +--------------------------------+
                                                  |
                     +----------------------------+----------------------------+
                     |                                                         |
         [ gitmap scan / gitmap add ]                              [ gitmap pull / pa ]
                     |                                                         |
         +-----------------------+                                 +-----------------------+
         | constants_scan.go     |                                 | pull_worker.go        |
         | DefaultScanExcludeDirs|                                 | SafePull Execution    |
         +-----------------------+                                 +-----------------------+
                     |                                                         |
         +-----------------------+                                 +-----------------------+
         | scanner.go            |                                 | pull_remediation_hint |
         | buildExcludeSet       |                                 | Classify Error &      |
         | Check --force-include |                                 | resolveMissingRepo... |
         +-----------------------+                                 +-----------------------+
                     |                                                         |
         +-----------------------+                                 +-----------------------+
         | fsutil traversal      |                                 | pull_efficient_render |
         | isDefaultScanExcluded |                                 | Format Nested Subtree |
         +-----------------------+                                 | Box Drawing Connectors|
                     |                                             +-----------------------+
                     v                                                         |
        +-------------------------+                                            v
        | SQLite Registry Sync    |                               +-------------------------+
        | OMZ & Caches Ignored    |                               | Terminal UI:            |
        | unless Explicitly Forced|                               | • <repo> failed         |
        +-------------------------+                               |   ├── Reason            |
                                                                  |   ├── Next Step         |
                                                                  |   ├── Options:          |
                                                                  |   │   ├── Option 1      |
                                                                  |   │   └── Option 2      |
                                                                  |   └── Diagnostic: pe    |
                                                                  +-------------------------+
```

---

## 4. Component Architecture: Oh-My-Zsh & Default Scanner Exclusions

### 4.1 Canonical Exclusion Manifest (`cli/constants/constants_scan.go`)

In environments such as Ubuntu and macOS, users frequently have `.oh-my-zsh` installed in their home directory (`~/.oh-my-zsh`), alongside custom plugins or scripts. Because Oh-My-Zsh is internally a Git repository, running `gitmap scan ~` inadvertently captures `.oh-my-zsh` into the GitMap registry. When users execute `gitmap pull-all`, `.oh-my-zsh` causes noise, permission issues, or unwanted dirty state warnings.

To prevent automatic capture across all discovery paths, `DefaultScanExcludeDirs` in `cli/constants/constants_scan.go` serves as the single source of truth.

```go
package constants

// DefaultScanExcludeDirs defines canonical directories automatically ignored by scanner and discovery.
var DefaultScanExcludeDirs = []string{
	".oh-my-zsh",
	"oh-my-zsh",
	"ohmyzsh",
	".ohmyzsh",
	"omz",
	".omz",
	"omizssh",
	"node_modules",
	"vendor",
	".cache",
	".venv",
	"venv",
	"__pycache__",
	".cargo",
	".rustup",
	".git",
	".terraform",
	".next",
	".turbo",
	"dist",
	"build",
	"bin",
	"obj",
	"target",
}
```

#### Naming Variant Rationale:
- `.oh-my-zsh`: Standard default install path (`$ZSH`).
- `oh-my-zsh`: Cloned without dot or unpacked directory.
- `ohmyzsh` / `.ohmyzsh`: Common alternative repository naming convention.
- `omz` / `.omz`: Shorthand directory name used in containerized environments and prompt setups.
- `omizssh`: Common user phonetical misspelling captured in real-world environments.

### 4.2 Multi-Layer Exclusion Enforcement

Exclusions are enforced across four distinct architectural checkpoints:

```
[Checkpoint 1: Record Config Defaults] -> cli/model/record.go: DefaultConfig()
                                                 |
[Checkpoint 2: Core Directory Scanner] -> cli/scanner/scanner.go: buildExcludeSet()
                                                 |
[Checkpoint 3: Child Repo Discovery]   -> cli/fsutil/child_repos.go: isDefaultScanExcluded()
                                                 |
[Checkpoint 4: Recursive Top-Level]    -> cli/fsutil/recursive_top_level.go: filepath.SkipDir
```

1. **Model Defaults (`cli/model/record.go`)**:
   `DefaultConfig()` initializes `ExcludeDirs: append([]string{}, constants.DefaultScanExcludeDirs...)`, ensuring default configuration files written to disk inherit full exclusion coverage.
2. **Scanner Build Set (`cli/scanner/scanner.go`)**:
   `buildExcludeSet(customExcludes []string, forceIncludes []string)` converts all default directory names to a fast lookup hash map (`map[string]bool`). Both verbatim and lowercase (`strings.ToLower`) variants are inserted into the set to ensure cross-platform case-insensitive matching on Windows, macOS, and Linux.
3. **Child Repository Discovery (`cli/fsutil/child_repos.go`)**:
   `isDefaultScanExcluded(name string) bool` performs `strings.EqualFold` checks against `constants.DefaultScanExcludeDirs` before inspecting immediate subdirectories for `.git`.
4. **Recursive Top-Level Traversal (`cli/fsutil/recursive_top_level.go`)**:
   During `filepath.WalkDir`, any directory matching `isDefaultScanExcluded(name)` immediately returns `filepath.SkipDir`. This prunes entire subtrees at the root, completely preventing filesystem recursion into `.oh-my-zsh` plugins, themes, or node modules.

---

## 5. Component Architecture: Explicit Force Override (`--force-include`)

### 5.1 Architecture & Flag Definitions

When a developer explicitly intends to track an excluded repository (such as developing custom Oh-My-Zsh themes or auditing vendor dependencies), GitMap provides a deterministic override mechanism:

```go
const (
	FlagScanForceInclude      = "force-include"
	FlagScanForceIncludeAlias = "fi"
	FlagDescScanForceInclude  = "Comma-separated directories to include despite default exclusions (or 'all' to include everything)"
	FlagAddForce              = "force"
	FlagAddForceAlias         = "f"
	FlagDescAddForce          = "Force tracking of directory even if matching default exclusions"
)
```

### 5.2 Override Evaluation Logic (`cli/scanner/scanner.go`)

In `buildExcludeSet`:
1. If `forceIncludes` contains `"all"` or `"*"`, the exclude set is emptied, allowing all directories to be crawled.
2. Otherwise, a `forceSet` hash map is built with lowercase normalization.
3. When iterating over `constants.DefaultScanExcludeDirs` and `customExcludes`, any directory present in `forceSet` is skipped (`continue`), removing it from the exclude set:

```go
func buildExcludeSet(customExcludes []string, forceIncludes []string) map[string]bool {
	for _, fi := range forceIncludes {
		if fi == "all" || fi == "*" {
			return make(map[string]bool)
		}
	}

	forceSet := make(map[string]bool, len(forceIncludes)*2)
	for _, fi := range forceIncludes {
		trimmed := strings.TrimSpace(fi)
		if trimmed != "" {
			forceSet[trimmed] = true
			forceSet[strings.ToLower(trimmed)] = true
		}
	}

	set := make(map[string]bool)
	for _, d := range constants.DefaultScanExcludeDirs {
		isForced := forceSet[d] || forceSet[strings.ToLower(d)]
		if isForced {
			continue
		}
		set[d] = true
		set[strings.ToLower(d)] = true
	}
	// ... apply customExcludes similarly
	return set
}
```

### 5.3 CLI Interface & Concrete Usage Examples

To ensure clarity for end users and scripts, the following command patterns are established:

```bash
# Example 1: Scan home directory while explicitly forcing inclusion of .oh-my-zsh
gitmap scan ~ --force-include .oh-my-zsh

# Example 2: Using the short alias -fi for multiple directories
gitmap scan ~ -fi .oh-my-zsh,node_modules

# Example 3: Forcing all excluded directories to be scanned
gitmap scan /var/repos --force-include all

# Example 4: Forcibly adding a specific existing directory that matches exclusions
gitmap add ~/.oh-my-zsh --force
```

CLI help documentation and `--help` renderers in `cli/cmdscan` and `cli/cmdadd` MUST surface these exact command examples.

---

## 6. Component Architecture: Failure Tree Subtree Hierarchy

### 6.1 Tree Topology & Box Connector Standards

During batch pull operations (`gitmap pull-all` / `gitmap pa`), failed repositories must provide crystal-clear diagnostic information and remediation actions without confusing visual clutter.

Previous renderers flattened options or emitted inconsistent bullet indentation. The enhanced hierarchy utilizes standardized Unicode box-drawing connectors:

```
    • <repository-name>  failed
    ├── Reason: <error-summary>
    ├── Next Step: <primary-action-command>
    ├── Options:
    │   ├── Option 1 (<Option-1-Title>): <Command-1>
    │   └── Option 2 (<Option-2-Title>): <Command-2>
    └── Diagnostic: To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)
```

#### Box Connector Tokens (`cli/cmdpull/pull_efficient_render.go`):
- `treeBranch` = `"├── "` (Mid-branch connector)
- `treeTerminal` = `"└── "` (Last child connector)
- `treeContinuation` = `"│   "` (Vertical pass-through for nested subtrees)
- `treeIndent` = `"    "` (Indent spacing)

### 6.2 Structural Rendering Implementation

In `cli/cmdpull/pull_efficient_render.go`:
- `renderSingleFailedItem`:
  1. Renders the repository bullet header with ANSI color formatting: `    • %-*s  failed`.
  2. Renders `Reason:` using `treeBranch` and `constants.ColorDim`.
  3. Renders `Next Step:` using `treeBranch` and `constants.ColorCyan`.
  4. Renders `Options:` subtree via `renderStructuredOptions`:
     - Emits `    ├── Options:` header.
     - Loops over `structured.Options`:
       - Non-terminal options use `    │   ├── Option N (Title): Command`.
       - Final option uses `    │   └── Option N (Title): Command`.
  5. Renders `Diagnostic:` using `treeTerminal` and `constants.ColorYellow`:
     `    └── Diagnostic: To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`

```go
func renderSingleFailedItem(w io.Writer, colWidth int, s *PullRepoState, collisions map[string]bool) {
	displayName := resolveRepoDisplayName(s, collisions)
	fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, displayName, "failed"))
	errDetails := ResolvePullErrorDetails(s)
	if errDetails == "" {
		errDetails = "pull execution failed"
	}
	fmt.Fprintf(w, "    %sReason: %s%s%s\n", treeBranch, constants.ColorDim, errDetails, constants.ColorReset)
	remHint := ResolvePullRemediationHint(s)
	if remHint == "" {
		remHint = fmt.Sprintf("gitmap status %s or gitmap fix %s", s.RepoName, s.RepoName)
	}
	fmt.Fprintf(w, "    %sNext Step: %s%s%s\n", treeBranch, constants.ColorCyan, remHint, constants.ColorReset)
	renderStructuredOptions(w, s)
	repoTarget := s.RepoName
	if repoTarget == "" {
		repoTarget = "all"
	}
	fmt.Fprintf(w, "    %sDiagnostic: To inspect stack trace: %sgitmap pull-error %s%s (or: %sgitmap pe%s)\n",
		treeTerminal, constants.ColorYellow, repoTarget, constants.ColorReset, constants.ColorDim, constants.ColorReset)
}
```

---

## 7. Component Architecture: Context-Aware Recovery & Dual Options

### 7.1 Root Cause of Legacy "Buggy Solution"

In earlier versions, when a repository registered in SQLite had its underlying directory deleted or moved from disk, `gitmap pa` failed with:
`missing repository directory` or `no such file or directory`.

The legacy remediation logic fell back to:
`gitmap status <repo>` or `gitmap pull <repo>`.

**Why this was flawed:**
- `gitmap status <repo>` requires navigating into the repository directory on disk to run `git status`. Since the directory does not exist, running `gitmap status` immediately crashed or returned another `directory does not exist` error!
- `gitmap pull <repo>` likewise required the directory to exist on disk to run `git pull`. It resulted in an infinite failure loop.

### 7.2 Context-Aware Recovery Design (`cli/cmdpull/pull_remediation_hint.go`)

The remediation engine must recognize disk state and provide two logical paths:
1. **Option 1 (Clone from Remote)**: The user wants the code back. Command: `gitmap clone <repo>`.
2. **Option 2 (Remove from Registry)**: The user intentionally deleted the repository and wants to clean up GitMap's SQLite registry without touching disk. Command: `gitmap rm --db-only <repo>`.

```go
// ResolveMissingRepoDualHints produces dual remediation options for missing repository errors.
func resolveMissingRepoDualHints(repoDir, repoName string) (string, string, string, string) {
	name := repoName
	if name == "" && repoDir != "" {
		name = filepath.Base(repoDir)
	}
	cloneCmd := "gitmap clone " + name
	removeCmd := "gitmap rm --db-only " + name
	return "Clone from Remote", cloneCmd, "Remove from Registry", removeCmd
}
```

### 7.3 Multi-Condition Remediation Matrix

| Error Type / Condition | Primary Next Step | Subtree Option 1 | Subtree Option 2 |
|:---|:---|:---|:---|
| **Missing Directory** (`isMissingRepoFailure`) | `gitmap clone <repo>` | Clone from Remote (`gitmap clone <repo>`) | Remove from Registry (`gitmap rm --db-only <repo>`) |
| **Diverged Branches** (`isDivergedFailure`) | `gitmap pull <repo>` | Preserve Local / Rebase (`git -C <dir> pull --rebase`) | Discard Local / Hard Reset (`git -C <dir> reset --hard @{u}`) |
| **Untracked Only** (`resolveUntrackedDualHints`) | `gitmap fix <repo>` | Track / Stage (`git -C <dir> add .`) | Clean untracked (`git -C <dir> clean -fd`) |
| **Dirty Working Tree** (`isDirtyTreeError`) | `gitmap fix <repo>` | Commit WIP (`gitmap cpar "wip: save changes"`) | Stash Changes (`git -C <dir> stash`) |
| **Authentication / SSH** (`isAuthFailure`) | `gitmap status <repo>` | Fix Windows Credentials (`gitmap fix-credential`) | Deploy SSH Keys (`gitmap ssh deploy-keys`) |
| **Generic / Unknown** (`resolveFallbackDualHints`) | `gitmap status <repo>` | Inspect Status (`gitmap status <repo>`) | Re-pull Repo (`gitmap pull <repo>`) |

---

## 8. Error Diagnostics & Stack Trace Integration

To guarantee total traceability requested by the user:
1. Every failed pull step records structured telemetry into `gitmap-pull.db` and writes an entry to `.gitmap/logs/pull-errors.log` (JSONL).
2. The terminal failure tree concludes with the exact diagnostic pointer:
   `└── Diagnostic: To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`
3. Running `gitmap pe <repo>` displays the categorized failure, exit code, and captured stderr.
4. Running `gitmap pe -t` (or `--trace`) outputs full Go stack traces and execution breadcrumbs for debugging.

---

## 9. Verification & Acceptance Criteria

### 9.1 Scanner Exclusions
- [ ] Running `gitmap scan` across a directory containing `.oh-my-zsh` or any variant (`oh-my-zsh`, `ohmyzsh`, `omz`, `omizssh`) does NOT capture it in results.
- [ ] Passing `--force-include .oh-my-zsh` or `-fi .oh-my-zsh` explicitly captures the directory.
- [ ] Passing `--force-include all` bypasses all default exclusions.
- [ ] Traversal functions (`DiscoverChildGitRepos`, `DiscoverTopLevelGitRepos`) skip `.oh-my-zsh` subtrees without descending into them.

### 9.2 Failure Tree Rendering
- [ ] Failed repository outputs conform strictly to the 4-tier tree hierarchy with box connectors:
  - `• <repo>  failed`
  - `├── Reason: ...`
  - `├── Next Step: ...`
  - `├── Options:`
  - `│   ├── Option 1: ...`
  - `│   └── Option 2: ...`
  - `└── Diagnostic: To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)`
- [ ] Missing repository error emits `Option 1: gitmap clone <repo>` and `Option 2: gitmap rm --db-only <repo>`.
- [ ] Terminal colors and column alignment render cleanly without wrapping or broken box glyphs.
