# App Issue 65: Commit Non-Git Repository Execution Crash, Child Batch Delegation & Terminal UI Output RCA

**Issue ID:** 65  
**Date:** 2026-10-03  
**Status:** Resolved  
**Affected Subsystem:** `cli/cmd` (`commit_cmd.go`, `commit_push.go`, `commit_batch.go`, `commit_ui.go`)  
**Reference Commit:** `02-spec/22-app-issues/65-commit-non-git-repo-crash-and-commit-all-delegation-rca.md` (PR #65 / Task 65)

---

## 1. Reproduction

When executing `gitmap commit all` (or `gitmap commit -a`, `gitmap commit -m "..."`, `gitmap cpf "..."`) in a directory that is not itself a Git repository (e.g. a parent workspace root containing multiple child repositories, or an arbitrary non-git folder):

```text
PS D:\work> gitmap commit all
  INFO Staging all changes...
fatal: not a git repository (or any of the parent directories): .git
gitmap commit: execute failed: [E9000:EXECUTION] git add failed:: exit status 128 (at=cmd/commit_cmd.go:52)
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.executeCommit (cmd/commit_cmd.go:52)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runCommit (cmd/commit_cmd.go:24)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.coreBasicOpEntries.func21 (cmd/rootcore.go:113)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatchTable (cmd/rootdispatch.go:24)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatchCore (cmd/rootcore.go:22)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:481)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:145)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:112)
    at main.main (cli/main.go:7)
```

Additionally, when running standard feature commits (`gitmap cpf "message"` or `gitmap commit-push`):
```text
warning: in the working copy of 'file.go', LF will be replaced by CRLF the next time Git touches it
```
Raw Git warning noise leaked into standard console output, and success responses ended abruptly without clear structured summary cards showing what branch and commit were updated.

### Symptoms:
1. `gitmap commit` and `gitmap commit-push` failed with an unhandled fatal application error (`[E9000:EXECUTION]`) and full Go stack trace when launched outside a Git repository root.
2. Running `gitmap commit all` from a workspace root containing multiple child Git repositories crashed instead of recognizing the intent to commit across child repositories.
3. Git's internal OS CRLF translation warnings (`LF will be replaced by CRLF`) polluted the command line output during staging and committing.
4. Output lacked indentation and visual hierarchy, and successful operations ended without a formatted summary card.

---

## 2. Root Cause Analysis (4-Part RCA)

### Symptom:
Executing `gitmap commit all` or `gitmap commit` in a non-git directory crashes with exit status 128 and prints a full internal stack trace `[E9000:EXECUTION] git add failed:: exit status 128`. In addition, Git CRLF noise floods output during standard commit operations.

### Direct Cause:
1. **Missing Pre-Flight Repository Validation:** `executeCommit` and `executeCommitPush` directly executed `git add -A` via `execGitInheritCP` without verifying that the current working directory (`CWD`) was inside an initialized Git worktree (`.git`).
2. **Missing Multi-Repo Delegation for `all`:** When users ran `gitmap commit all` from a multi-repo root, `gitmap` treated `all` as an argument or commit message rather than detecting a request for recursive/batch child repository commits.
3. **Unfiltered Git Subprocess Output:** Output streams from Git subprocesses were inherited directly (`os.Stdout`, `os.Stderr`), allowing platform-specific line ending warnings to be dumped directly into the user terminal.

### Compound Failure:
1. **Developer Friction on High-Frequency Commands:** `gitmap commit` and `gitmap cpf` are among the most frequently executed developer commands. Crashing with raw stack traces when outside a repo degrades user trust.
2. **Workspace Root Inconsistency:** Commands like `gitmap status all`, `gitmap pull all`, and `gitmap fetch all` successfully discover and process child repositories, but `gitmap commit all` lacked equivalent workspace-level delegation.
3. **Noisy Terminal Transcripts:** CRLF warnings obscure meaningful git status and commit messages in automated CI logs and agent terminals.

### Impact:
Users and automated agents running `gitmap commit all` at workspace level encountered hard execution failures. Terminal transcripts were noisy and unformatted.

---

## 3. Fix & Remediation

1. **Preflight Git Repository Check (`isGitRepoCWD`)**:
   - Added pre-flight git repository validation in `cli/cmd/commit_cmd.go` and `cli/cmd/commit_push.go`.
   - When the current working directory is not a Git repository, commands avoid invoking raw `git add` and gracefully divert to non-git handlers.

2. **Workspace Child Repository Batch Delegation (`cli/cmd/commit_batch.go`)**:
   - Implemented `isAllCommitRequested(args []string) bool` to identify when the user specified `all` or `--all`.
   - Implemented `handleNonGitRepoCommit(cleanArgs []string, hasPush, isDryRun bool) error`:
     - Discovers all child Git repositories within the workspace root.
     - For each child repo with staged or unstaged changes, runs batch commit (and push if `--push` was specified).
     - In `--dry-run` mode, displays dirty child repos that would be committed.
     - If no child Git repositories are found or if `all` was not requested, prints a clean guidance message and returns an abort AppError (`reported: true`) without dumping internal stack traces.

3. **Terminal Output Sanitization & Indentation (`cli/cmd/commit_ui.go`)**:
   - Implemented `isGitNoiseLine(line string) bool` to filter out verbose `LF will be replaced by CRLF` line ending warnings.
   - Implemented `execGitPaddedFiltered(gitArgs ...string) error` to capture Git output, filter noise, and indent remaining status lines with 4 spaces for clear visual hierarchy.

4. **Completion Summary Card (`renderCommitPushSummaryCard`)**:
   - Formatted success output showing branch, short commit hash, commit message, and remote push status.

---

## 4. Prevention & Learnings

1. **Pre-flight Guards for Tool Subprocesses:** Any CLI command wrapping an external executable (`git`, `ssh`, `docker`) must verify the executable's prerequisites (such as repository existence or network configuration) before launching subprocesses.
2. **Consistent Workspace-Level Semantics:** All `* all` commands in GitMap (`status all`, `pull all`, `commit all`) must consistently implement child repository auto-discovery.
3. **Clean Terminal Discipline:** Never dump raw subprocess stderr or unformatted noise directly to the user when wrapping external CLI tools; filter known benign warnings and maintain indented visual styling.
