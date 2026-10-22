# Subtask 287.1: Audit & Refactor Void Error Swallowing in `cli/cmdpull/`

- **Parent Plan:** [.ai-memory/plans/pending/287-golang-method-return-and-gitmap-improvement-instructions.md](../../pending/287-golang-method-return-and-gitmap-improvement-instructions.md)
- **Status:** Completed
- **Completed At:** 2026-10-22
- **File Boundaries:**
  - `cli/cmdpull/pull_remediation_ssh.go`
  - `cli/cmdpull/pull_remediation.go`
  - `cli/cmdworkdir/workdir_cmd.go`

## Objectives & Outcomes
1. Refactored `updateRepoTransportInDB(repoPath, sshURL string) error`:
   - Returns `error` directly from `db.UpdateRepoTransport(...)`.
   - Propagates errors if `store.OpenDefault()` fails.
2. Refactored `clearResolvedPullError(repoName, repoPath string) error`:
   - Returns `error` directly from `db.ClearPullErrorsForRepo(...)`.
   - Propagates errors if `store.OpenPullSplitDB()` fails.
3. Refactored `applySSHTransportRemediation(repoPath string) error`:
   - Returns `error`. Returns `updateRepoTransportInDB(repoPath, newURL)` on URL transition.
4. Refactored `remediateCredentialAndConvertToSSH(repoPath, repoName string) error`:
   - Decomposed into `configureGlobalCredentialStore(storeType constants.CredentialStoreType) error`.
   - Returns errors directly instead of swallowing with `_ = ...`.
5. In `cli/cmdpull/pull_remediation.go`:
   - Handled errors in `applyPreemptiveCredentialFix`, `handleSuccessfulRemediation`, `tryRetryAfterCredentialRemediation`, and `handleSuccessfulRemediationSSH`.
   - Refactored `execPullAutoMerge` to use centralized `gitutil.RunGitCombinedInDir`.
6. In `cli/cmdworkdir/workdir_cmd.go`:
   - Refactored `ensureWorkDirsTableExists() error` to return errors from `store.OpenDefault()` and `db.SQL().Exec(...)`.
   - Updated `runWorkDir` caller to return `errEnsure`.
7. Verified:
   - FastGate AST linter passed with zero violations.
   - `go test -C cli -v ./cmdpull` passed with 100% green.
