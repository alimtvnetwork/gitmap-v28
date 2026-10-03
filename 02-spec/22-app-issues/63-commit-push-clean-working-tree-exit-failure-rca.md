# App Issue 63: Commit-Push Clean Working Tree Exit Status 1 Failure & Missing Staged Changes Guard RCA

**Issue ID:** 63  
**Date:** 2026-10-03  
**Status:** Resolved  
**Affected Subsystem:** `cli/cmd` (`commit_push.go`, `commit_cmd.go`)

---

## 1. Reproduction

When executing `gitmap cpf "updated letterly prompts"` (or `gitmap cp`, `gitmap commit -p`, `gitmap commit-push-feature`, `gitmap commit-push-bug`, `gitmap commit-push-release`) in a repository where all working tree files are already up-to-date and there are no staged or unstaged modifications:

```text
$ gitmap cpf "updated letterly prompts"
  INFO Staging all changes...
  INFO Committing: Feature: updated letterly prompts
On branch master
Your branch is up to date with 'origin/master'.

nothing to commit, working tree clean

[E9000:EXECUTION] git commit failed:: exit status 1 (at=cmd/commit_push.go:334)
Stack Trace:
  [0] github.com/alimtvnetwork/gitmap-v28/cli/cmd.executeCommitPush
      d:/work/gitmap/cli/cmd/commit_push.go:334
  [1] github.com/alimtvnetwork/gitmap-v28/cli/cmd.runCommitPushFeature
      d:/work/gitmap/cli/cmd/commit_push.go:123
  [2] github.com/alimtvnetwork/gitmap-v28/cli/cmd.Execute
      d:/work/gitmap/cli/cmd/root.go:88
  [3] main.main
      d:/work/gitmap/cli/main.go:24
```

### Symptoms:
1. `git commit -m "..."` terminates with OS exit status 1 when there are no staged changes.
2. `executeCommitPush` and `executeCommit` blindly invoked `execGitInheritCP("commit", "-m", commitMessage)` immediately following `execGitInheritCP("add", "-A")`.
3. An `apperror.WrapSimple(err, "git commit failed:")` is returned, printing an unnecessary fatal error stack trace `[E9000:EXECUTION]` and exiting with non-zero status to the shell.
4. If there were local commits already committed but unpushed, the command aborted before running `git push`, leaving those commits unpushed.

---

## 2. Root Cause Analysis (4-Part RCA)

### Symptom:
Running any commit command variant (`commit-push`, `cpf`, `cpb`, `cpr`, or `cm --push`) on a clean working tree fails with `[E9000:EXECUTION] git commit failed:: exit status 1` instead of gracefully completing or pushing outstanding local commits.

### Direct Cause:
Standard Git behavior specifies that `git commit` exits with status code 1 when invoked with no staged changes (unless `--allow-empty` is supplied). Neither `executeCommitPush` (in `commit_push.go`) nor `executeCommit` (in `commit_cmd.go`) checked git status or staged changes prior to dispatching `git commit`.

### Compound Failure:
1. **Blind Subprocess Execution:** The command assumed `git add -A` would always stage files. When the working tree had no modifications, nothing was staged, and the subsequent `git commit` inevitably failed.
2. **Push Interruption:** When a user had committed changes previously (e.g., via `git commit` or an automated script) and ran `gitmap cp` to push them, the failure of `git commit` halted execution, completely bypassing the push step.
3. **Misleading Error Reporting:** The user was presented with an application-level fatal error `[E9000]` instead of an informative status indicating that the working tree is clean.

### Impact:
Developer workflows, automated scripts, CI/CD runners, and subagent loops that run `gitmap cpf` or `gitmap commit -p` were interrupted by false-alarm non-zero exit codes.

---

## 3. Fix & Remediation

1. **Staged Changes Detection Helper (`hasStagedChangesCP`)**:
   - Added `hasStagedChangesCP() (bool, error)` in `cli/cmd/commit_push.go` which runs `git status --porcelain`.
   - Returns `true` if any staged or unstaged modifications, untracked files, or deletions exist, or `false` when porcelain output is empty.

2. **Unpushed Commits Counter (`countUnpushedCommitsCP`)**:
   - Added `countUnpushedCommitsCP() int` in `cli/cmd/commit_push.go` which inspects `git rev-list --count @{u}..HEAD`.
   - Returns the integer count of local commits ahead of upstream remote tracking, or `0` if upstream is not set or errors occur.

3. **Guarded Commit Dispatch in `commit_push.go`**:
   - In `executeCommitPush(commitMessage string)`:
     - After `git add -A`, calls `hasStagedChangesCP()`.
     - If error occurs, wraps with `apperror.WrapSimple(err, "check git status failed:")`.
     - If `!hasChanges`:
       - Prints info: `Working tree clean, nothing to commit.`
       - Checks `unpushed := countUnpushedCommitsCP()`.
       - If `unpushed > 0`, prints `Pushing %d unpushed commit(s) to remote...`, executes `git push`, prints `Pushed %d commit(s) to remote.`, and returns `nil`.
       - If `unpushed == 0`, prints `Everything is up to date.` and returns `nil`.

4. **Guarded Commit Dispatch in `commit_cmd.go`**:
   - In `executeCommit(args []string, hasPush bool)`:
     - After `git add -A`, checks `hasStagedChangesCP()`.
     - If `!hasChanges`, prints `Working tree clean, nothing to commit.` and delegates to `handleCleanWorkingTree(hasPush)`.
   - In `handleCleanWorkingTree(hasPush bool)`:
     - If `!hasPush`, prints `Working tree clean.` and returns `nil`.
     - If `hasPush`:
       - Checks `unpushed := countUnpushedCommitsCP()`.
       - If `unpushed > 0`, pushes unpushed commits with friendly logs and returns `nil`.
       - Otherwise, prints `Everything is up to date.` and returns `nil`.

---

## 4. Prevention & Learnings

1. **Verify State Before Destructive / Guarded Subprocesses:** Before executing external tool commands that fail on no-op states (such as `git commit`), inspect the state beforehand with low-overhead queries (such as `git status --porcelain`).
2. **Preserve Dual-Action Intent:** For multi-action composite commands like `commit-push`, a no-op in the first stage (commit) must not arbitrarily discard the second stage (push) if there is pending work to be completed (such as unpushed commits).
3. **Graceful Exit Codes for Idempotent Operations:** Clean working tree operations are standard, benign conditions; they must report user-friendly status information and exit with code 0.
