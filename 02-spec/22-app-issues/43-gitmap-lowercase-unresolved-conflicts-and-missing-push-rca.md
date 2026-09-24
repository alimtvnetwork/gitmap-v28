# Issue 43: GitMap Lowercase Unresolved Merge Conflicts, Working Tree Dirt Contamination & Missing Push RCA

> **Issue ID:** `ISSUE-43`  
> **Spec Reference:** [02-spec/21-app/156-gitmap-lowercase-preflight-hygiene-and-conflict-resolution.md](../21-app/156-gitmap-lowercase-preflight-hygiene-and-conflict-resolution.md)  
> **Status:** Resolved  
> **Date:** 2026-09-25  

---

## 1. Problem Reproduction & Symptom Analysis

### User Invocation & Observed Behavior:
```text
PS D:\work\presentations-repos\hiltrax> gitmap lowercase "*.md"
...
  ● Found: 18 uppercase file(s) to rename:
    • SKILL.md -> skill.md
    ... and 10 more file(s)

Type 'confirm' or 'yes' (or 'y') to proceed: y
  [1/18] SKILL.md -> skill.md (OK)
  ...
✔ Committed 18 lowercase file rename(s): "chore: rename 18 files to lowercase across repository" (9582c388)
```

### Problems Reported:
1. **Unresolved Git Conflicts:** The renamer was executed while the working tree had pending uncommitted changes or in-flight merge conflicts.
2. **Indiscriminate Git Staging:** `git add -A` staged pre-existing dirty files and conflicted markers, committing unmerged state into the rename commit.
3. **Hidden Files:** Only 8 of 18 files were displayed, truncating the rest as `... and 10 more file(s)`, preventing developers from knowing exactly which files were being renamed.
4. **Missing Push:** The commit was created locally but never pushed to the remote tracking branch, leading to divergent branches and merge conflicts on subsequent pulls.

---

## 2. Root Cause Analysis (4-Part RCA)

### 1. Why Did Merge Conflicts and Dirty Commits Occur?
`lowercasefix_ops.go` executed without any pre-flight inspection of `git status --porcelain`. When `git add -A` was called in `lowercasefix_commit.go`, Git staged **all** modified, conflicted, and untracked files across the entire repository—not just the renamed files.

### 2. Why Were Files Hidden?
In `lowercasefix_confirm.go`, `renderPreflightFound` imposed a hardcoded ceiling of `limit := 8`, truncating the list with `... and %d more file(s)`. On large batches (like 18 files), over half the files were hidden from the user.

### 3. Why Were Local Commits Not Pushed?
`maybeCommitRenames` in `lowercasefix_ops.go` called `commitRenames()`, which concluded after running `git commit -m ...`. No `git push` command was ever executed or scheduled, leaving the local branch ahead of origin.

### 4. Why Was Working Tree Dirt Not Discarded or Resolved?
There was no prompt or mechanism for the user to inspect uncommitted pending files or choose to discard them prior to renaming.

---

## 3. Resolution & Code Fix Strategy

1. **Pre-Flight Working Tree & Conflict Detection:**
   - Execute `git status --porcelain`.
   - If unmerged files exist (`UU`, `AA`, etc.), immediately abort with actionable error.
   - If uncommitted changes exist, list them explicitly and prompt: `Discard pending changes to proceed with clean renames? [y/N]` (or use `--discard-pending`).
2. **Complete File Disclosure:**
   - Remove the `limit := 8` cutoff in `renderPreflightFound`; display all candidate files with relative repository paths so there are zero hidden items.
3. **Automated Git Push:**
   - Add `git push` execution in `commitRenames` (with `--no-push` flag to opt out), ensuring the commit is pushed directly to the current tracking branch.

---

## 4. Prevention Architecture

- Unit tests in `lowercasefix_test.go` and `lowercasefix_confirm_test.go` verifying that dirty working trees trigger pre-flight warnings and that all matched files are reported without truncation.
