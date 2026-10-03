# RCA-28: Commit-Push Clean Working Tree Exit Status 1 Failure & Missing Staged Changes Guard

Spec Reference: [02-spec/22-app-issues/63-commit-push-clean-working-tree-exit-failure-rca.md](../../02-spec/22-app-issues/63-commit-push-clean-working-tree-exit-failure-rca.md)

---

## 1. Symptom

Running commit commands (`gitmap cpf "message"`, `gitmap cp`, `gitmap commit -p`, `gitmap cpb`, `gitmap cpr`) in a repository with a clean working tree resulted in an unhandled exit code 1 crash:

```text
nothing to commit, working tree clean
[E9000:EXECUTION] git commit failed:: exit status 1 (at=cmd/commit_push.go:334)
```

Furthermore, any previously committed but unpushed commits were blocked from being pushed because the command terminated immediately at the failed `git commit` step.

---

## 2. Root Cause

1. `git commit` without `--allow-empty` returns exit code 1 when there are no staged changes.
2. `executeCommitPush` (`cli/cmd/commit_push.go`) and `executeCommit` (`cli/cmd/commit_cmd.go`) invoked `git commit` unconditionally after `git add -A` without checking if any changes were staged.
3. The lack of staged-change detection broke composite commands (`commit-push`), preventing subsequent push operations for already existing local commits.

---

## 3. Resolution

1. **Staged Changes Detection Helper (`hasStagedChangesCP`)**:
   - Added `hasStagedChangesCP() (bool, error)` in `cli/cmd/commit_push.go` utilizing `git status --porcelain`.
2. **Unpushed Commits Counter (`countUnpushedCommitsCP`)**:
   - Added `countUnpushedCommitsCP() int` in `cli/cmd/commit_push.go` utilizing `git rev-list --count @{u}..HEAD`.
3. **Graceful Handling in `commit_push.go`**:
   - In `executeCommitPush`, guarded commit execution: if working tree is clean, prints `Working tree clean, nothing to commit.`
   - Pushes unpushed commits if `unpushed > 0`, or prints `Everything is up to date.` and exits cleanly with 0.
4. **Graceful Handling in `commit_cmd.go`**:
   - In `executeCommit`, checks `hasStagedChangesCP()`. If clean, delegates to `handleCleanWorkingTree(hasPush)`.
   - In `handleCleanWorkingTree`, reports clean status, pushes unpushed commits when `hasPush == true`, and exits cleanly with 0.

---

## 4. Prevention & Learnings

- Always verify prerequisite states before executing non-idempotent or non-zero exiting subprocesses like `git commit`.
- Decouple composite multi-stage actions (e.g. commit then push) so that a benign no-op in stage one does not cancel stage two.
- Standardize zero-exit responses with informative logging for idempotent developer commands.
