# Subtask 03: Work Directory Boundary Guard & Safe Humanized Stashing

## Status: COMPLETED
- **Parent Plan**: `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md`
- **Specification Reference**: `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/02-component-and-cli-spec.md`
- **Owner**: Worker 02 (CLI Operations & Git Engine Engineer)
- **Target Subsystem**: `cli/cmdfix`, `cli/cmdpull`, `cli/gitutil`
- **Owned Files**:
  - `cli/cmdfix/fix_nonrepo.go`
  - `cli/cmdfix/fix_casing.go`
  - `cli/cmdfix/fix_stashing.go` (new file)
  - `cli/cmdfix/fix_execute.go`
  - `cli/cmdpull/pull_cwd.go`
  - `cli/cmdpull/pull_efficient_render.go`
  - `cli/gitutil/status.go` (new file)

---

## 1. Context & Objectives

1. **Work Directory Commit Guard (`cli/cmdfix/fix_nonrepo.go`, `cli/cmdfix/fix_casing.go`, `cli/cmdpull/pull_cwd.go`)**:
   - Prevent `gitmap fix` and `gitmap fix --all` from attempting to stage, commit, or stash inside workspace root folders (e.g., `work/`) that do not contain a `.git` repository directory.
   - Enforce an invariant check (`IsValidGitRepository` and `IsWorkspaceContainerDirectory`) across all path resolution, casing reconciliation, and CWD pull targets.
   - Ensure clear, non-cryptic error reporting if a user or script accidentally targets a workspace container.

2. **Safe & Humanized Auto-Stashing (`cli/cmdfix/fix_stashing.go`, `cli/gitutil/status.go`)**:
   - Parse Git porcelain status into distinct, strongly-typed indicators:
     - `[M]` Modified tracked file
     - `[A]` Added / staged file
     - `[D]` Deleted tracked file
     - `[?]` Untracked file
   - Display a boxed Lipgloss/termpad notification card detailing the files to be stashed, the explicit intention, and the safety guarantee that work will be cleanly popped.
   - Limit display to 5 files with an overflow count (e.g., `... and 12 other file(s)`).
   - Hook into `fix_execute.go` to provide transparent transition steps (`Stashing...`, `Pulling...`, `Restoring...`).

3. **Condense Unchanged Repositories & Suggest One-Liner (`cli/cmdpull/pull_efficient_render.go`)**:
   - In `cli/cmdpull/pull_efficient_render.go`, replace verbose multi-line listings of unchanged repositories with a single, high-signal line:
     `✓ Updated all other 67 repositories (no changes).`
   - Print a suggested one-liner inspection command below the pull summary:
     `💡 Run 'gitmap pull --ls-all' or 'gitmap status --all' to inspect all unchanged repositories.`

4. **Code Quality & Coding Guideline Discipline**:
   - All functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses and early returns only).
   - Positive boolean naming (`isValid`, `hasDirtyWork`, `isStaged`).
   - AppError wrapping on all error return paths.
   - Zero test / zero build runs during specification authoring.

---

## 2. Step-by-Step Implementation Plan

### Step 3.1: Implement Directory Boundary Guards in `cli/cmdfix/fix_nonrepo.go`
1. Define `IsValidGitRepository(targetPath string) bool`:
   - Guard 1: Verify `targetPath != ""` and `info, err := os.Stat(cleanPath); err == nil && info.IsDir()`.
   - Guard 2: Verify `filepath.Join(cleanPath, ".git")` exists.
   - Guard 3: Fallback check via `git -C <path> rev-parse --is-inside-work-tree` == `"true"`.
2. Define `IsWorkspaceContainerDirectory(targetPath string) bool`:
   - Inspect folder base name: match `work`, `git-work`, `workspace`, `projects`.
   - If folder name matches and `.git` is absent, return `true`.
3. Update `resolveRepoPathForFix(repoQuery string) string`:
   - If `repoQuery` resolves to a path that is a workspace container or not a git repository, return `""` to prevent execution on non-git roots.
4. Add `GuardRepoExecution(targetPath string, opName string) *apperror.AppError`:
   - Returns validation error `apperror.NewValidationError("cannot execute " + opName + ": directory is not a Git repository")` if validation fails.

### Step 3.2: Guard Casing Reconcile and CWD Pull Target
1. In `cli/cmdfix/fix_casing.go`:
   - In `reconcileAllCasingDiscrepancies()`: verify `IsValidGitRepository(r.AbsolutePath)` before calling `DetectCasingDiscrepancies`.
   - In `DetectCasingDiscrepancies(repoPath string)`: add guard clause at line 40: return empty slice if `!IsValidGitRepository(repoPath)`.
2. In `cli/cmdpull/pull_cwd.go`:
   - In `isPullCWDEnabled`: ensure `isGitRepoCWD()` returns `false` if `IsWorkspaceContainerDirectory(cwd)` is `true`.
   - In `runPullCWD`: if current directory is not a git repo, return descriptive error rather than attempting Git commands.

### Step 3.3: Implement File Status Inspection in `cli/gitutil/status.go`
1. Define data structures:
   ```go
   type FileStatusCode string
   const (
       StatusModified  FileStatusCode = "[M]"
       StatusAdded     FileStatusCode = "[A]"
       StatusDeleted   FileStatusCode = "[D]"
       StatusUntracked FileStatusCode = "[?]"
       StatusRenamed   FileStatusCode = "[R]"
   )
   type DirtyFileInfo struct {
       Path       string
       StatusCode FileStatusCode
       IsStaged   bool
   }
   type RepoDirtyManifest struct {
       RepoPath       string
       RepoName       string
       HasDirtyWork   bool
       Files          []DirtyFileInfo
       ModifiedCount  int
       UntrackedCount int
       DeletedCount   int
   }
   ```
2. Implement `InspectDirtyManifest(repoPath string) (RepoDirtyManifest, error)`:
   - Run `git -C <path> status --porcelain=v1`.
   - Parse each line's 2-character status code:
     - `??` -> `StatusUntracked`
     - ` M` or `M ` -> `StatusModified`
     - ` D` or `D ` -> `StatusDeleted`
     - `A ` or `AM` -> `StatusAdded`
   - Populate `RepoDirtyManifest` with counts and file slice.

### Step 3.4: Implement Humanized Stash Card in `cli/cmdfix/fix_stashing.go`
1. Define rendering function:
   ```go
   func RenderHumanizedStashCard(manifest gitutil.RepoDirtyManifest) string
   ```
2. Format card using Lipgloss / termpad styling:
   - Box title: `┌── Safe Auto-Stash: <RepoName> ──┐`
   - Description line: `GitMap is temporarily stashing N uncommitted file(s) before sync:`
   - Up to 5 dirty file lines with cyan/yellow flags (`[M]`, `[D]`, `[?]`).
   - If `len(manifest.Files) > 5`, print `... and N other file(s)`.
   - Safety guarantee line: `ℹ Safety Guarantee: All changes will be re-applied after pull.`
3. Implement `ExecuteSafeStashLifecycle(repoPath string, pullFn func() error) error`:
   - Step 1: Inspect dirty manifest. If clean, proceed directly with `pullFn()`.
   - Step 2: Print `RenderHumanizedStashCard(manifest)`.
   - Step 3: Run `git stash -u` and print progress ticker.
   - Step 4: Execute `pullFn()`.
   - Step 5: Run `git stash pop` and print confirmation banner.
   - Step 6: Handle conflicts gracefully if pop detects merge issues.

### Step 3.5: Implement Unchanged Repo Collapsing in `cli/cmdpull/pull_efficient_render.go`
1. Update `renderUpToDateGroup(w io.Writer, upToDate []*PullRepoState, collisions map[string]bool)`:
   - If `len(upToDate) > 5` and `len(cat.updated) > 0 || len(cat.dirty) > 0 || len(cat.failed) > 0`:
     - Render collapsed summary:
       `fmt.Fprintf(w, "\n  %s✓ Updated all other %d repositories (no changes).%s\n", constants.ColorGreen, len(upToDate), constants.ColorReset)`
     - Print inspection suggestion:
       `fmt.Fprintf(w, "\n  %s💡 Run 'gitmap pull --ls-all' to inspect all unchanged repositories.%s\n", constants.ColorDim, constants.ColorReset)`
   - If `len(cat.updated) == 0 && len(cat.dirty) == 0 && len(cat.failed) == 0`:
     - Render standard up-to-date group or compact summary.

---

## 3. Acceptance Criteria & Quality Gates

- [x] `gitmap fix` executed from a parent folder (e.g., `work/`) immediately detects non-repository boundaries and refuses to stage/commit.
- [x] Auto-stashing prints dirty files with `[M]`, `[D]`, `[?]` flags, formatted in a clear box before stashing.
- [x] Pull summaries with >3 unchanged repositories collapse into a single line: `✓ Updated all other N repositories (no changes)`.
- [x] The suggested one-liner inspection command is rendered clearly below pull results.
- [x] All Go functions adhere to the 15-line limit (<= 8 lines preferred).
- [x] Positive boolean names (`isValid`, `hasDirtyWork`, `isStaged`) used exclusively.
- [x] Zero build or test commands executed during this spec authoring phase.
