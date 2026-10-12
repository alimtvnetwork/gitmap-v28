# Issue 31: Remediation Unborn Branch Stash Failure and Batch Resiliency RCA

> **Issue ID:** `ISSUE-31`
> **Target Subsystem:** `cli/cmdfix`, `cli/gitutil`, `cli/cmdremediation`, `cli/cmdreconcile`
> **Affected Repositories:** `repo-cache`, `repo-secrets`, and multi-repo remediation batches
> **Status:** Resolved
> **Date:** 2026-10-12

---

## 1. Problem Reproduction & Symptom Analysis

### User Invocation & Observed Error
When running `gitmap fix --all` across registered repositories:
```text
@a:~$ gitmap fix --all
ℹ Remediating 5 repository(ies) with action: stash

ℹ Applying Fix: Option 1 (Stash & Re-apply) on gitmap-v28
  Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

  [1/3] ➜ git stash -u ... ✔ ok
  [2/3] ➜ git pull ... ✔ ok
  [3/3] ➜ git stash pop ... • recovered (dropped stashed untracked collision)

✓ Fix applied successfully on gitmap-v28

ℹ Applying Fix: Option 1 (Stash & Re-apply) on gstack
  Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

  [1/3] ➜ git stash -u ... ✔ ok
  [2/3] ➜ git pull ... ✔ ok
  [3/3] ➜ git stash pop ... • clean (nothing to commit/apply)

✓ Fix applied successfully on gstack

ℹ Applying Fix: Option 1 (Stash & Re-apply) on repo-cache
  Plan:    Temporarily save local changes (including untracked), pull latest remote commits, then re-apply

  [1/3] ➜ git stash -u ... ✖ failed

✖ Remediation Step Failed on repo-cache: exit status 1
  Command:   git -C /home/a/git-work/repo-cache stash -u
  Output:
    You do not have the initial commit yet
  RCA (Root Cause): Git process exited with failure: exit status 1
  Known Solutions:
    • Inspect working tree status: git -C <repo> status
    • Check git log: git -C <repo> log -n 3 --oneline

gitmap fix: execute failed: exit status 1
```

---

## 2. Root Cause Analysis (4-Part RCA)

### 1. Symptom
`gitmap fix --all` crashed on `repo-cache` with `exit status 1` and `You do not have the initial commit yet` when attempting `git stash -u`. The failure aborted the entire batch loop immediately, leaving subsequent repositories un-remediated.

### 2. Root Cause
1. **Unborn HEAD Git Invariant:** Git fundamentally rejects `git stash` on repositories without commits (`HEAD` is unborn). `repo-cache` was initialized by `ensureSpecialGitRepo` via `git init` without commits or remote tracking.
2. **Internal Metadata Directory Leak:** GitMap stores SQLite tracking files under `<repo>/.gitmap/` (such as `pull_errors.db`). Without `.gitmap/` in `.git/info/exclude` or filtering in `dirty_inspect.go`, `git status --porcelain` reported `?? .gitmap/`, falsely tagging clean repositories as dirty (`+1 untracked`).
3. **Static Stash Recipe Generation:** `GenerateRemediationRecipes` unconditionally generated stash/pull recipes for all dirty repositories, irrespective of whether the repository possessed an initial commit or configured remotes.
4. **Queue Abort on Single Error:** `runFixAll` in `cli/cmdfix/fix_cmd.go` and `runReconcileAll` in `cli/cmdreconcile/reconcile_cmd.go` returned errors immediately inside the iteration loop instead of isolating failures and completing the remaining repositories in the batch.
5. **Generic RCA Diagnostic:** `analyzeGitErrorOutput` lacked a parser branch for unborn commit messages (`You do not have the initial commit yet`).

### 3. Resolution
1. **GitMap Internal Directory Isolation:**
   - Updated `classifyDirtyFile` in `cli/gitutil/dirty_inspect.go` to ignore `.gitmap` and `.gitmap/*` paths unconditionally.
   - Added `EnsureGitMapExcluded(repoRoot string)` in `cli/gitutil/exclude.go` and called it in `ensureSpecialGitRepo` (`cli/cmd/special_repos_ops.go`), `saveRepoPullErrorDetails` (`cli/store/pull_split_db_errors.go`), and `ensureGitExcluded` (`cli/cmd/commitin/workspace/paths.go`).
2. **Context-Aware Recipe Generation:**
   - Added `HasInitialCommit(repoPath string) bool` and `HasRemote(repoPath string) bool` in `cli/gitutil/gitutil.go`.
   - Updated `GenerateDefaultRecipes` in `cli/gitutil/remediation_generator.go`:
     - Repositories without initial commits receive `GenerateNoInitialCommitRecipes` (Option 1: Initial Commit via `git add -A && git commit -m "feat: initial commit"`, Option 2: `git clean -fd`, Option 3: `git status` or `git fetch`).
     - Repositories without remotes receive `GenerateLocalOnlyRecipes` (Option 1: Commit WIP, Option 2: Stash, Option 3: Discard).
3. **Execution Self-Healing & Recipe Sanitization:**
   - Implemented `isRepoAlreadyClean`, `sanitizeRecipeForRepo`, `recoverNoInitialCommit`, and `isBenignPullNoTracking` in `cli/cmdfix/fix_recover.go` and `cli/cmdremediation/helpers_shim.go`.
   - When stash fails on an unborn branch with uncommitted changes, GitMap automatically stages and creates the initial commit. If clean, it reports clean and continues.
   - When pulling without configured remotes or tracking information, GitMap marks the step as benign clean.
4. **Batch Queue Resiliency:**
   - Updated `runFixAll` in `cli/cmdfix/fix_cmd.go` and `runReconcileAll` in `cli/cmdreconcile/reconcile_cmd.go` to accumulate failures into `failedRepos []string` and continue processing all items, returning a consolidated summary.
5. **Precise 4-Part Diagnostic:**
   - Added `noInitialCommitDiagnostic()` in `cli/cmdfix/fix_diagnostics.go` and `cli/cmdremediation/helpers_shim.go` detailing unborn branch constraints and actionable solutions.

### 4. Prevention & Learnings
- **Pre-flight Validation:** Git commands that depend on repository invariants (`HEAD` commit existence, remote tracking branch) must be gated by state inspections (`HasInitialCommit`, `HasRemote`) before recipes are generated or executed.
- **Batch Isolation Invariant:** Multi-repository orchestrators must never let a localized repository failure terminate the entire batch process. Every failure must be isolated, accumulated, and presented in the final summary.
- **Internal Storage Exclusion:** Self-generated metadata folders (`.gitmap/`) must be natively excluded in `.git/info/exclude` at the point of creation so standard Git tools never perceive them as workspace modifications.
