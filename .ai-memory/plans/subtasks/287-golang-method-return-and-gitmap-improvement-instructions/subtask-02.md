# Subtask 287.2: Codebase-Wide Survey of Void Methods & Swallowed Errors

- **Parent Plan:** [.ai-memory/plans/pending/287-golang-method-return-and-gitmap-improvement-instructions.md](../../pending/287-golang-method-return-and-gitmap-improvement-instructions.md)
- **Status:** Completed
- **Completed At:** 2026-10-22
- **Spec Reference:** [02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/02-component-and-error-return-spec.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/02-component-and-error-return-spec.md)

## Objectives & Outcomes
1. Conducted an exhaustive codebase-wide survey of void functions, swallowed errors (`_ = ...`), and unhandled return values across `cli/`:
   - Identified primary defects in `cli/cmdpull/pull_remediation_ssh.go` (`updateRepoTransportInDB`, `clearResolvedPullError`, `applySSHTransportRemediation`, `remediateCredentialAndConvertToSSH`, `loadStoreRepos`).
   - Identified secondary defects across `cli/cmdworkdir/workdir_cmd.go` (`ensureWorkDirsTableExists`), `cli/cmdvhost/vhost_rm.go` (`deleteSiteFromDB`), `cli/cmdssh/sshjoin_enroll.go` (`savePostInstallProfile`), `cli/cmdpull/pull_db_sync.go` (`syncPullFailuresToDB`), `cli/cmdpull/pull_ignore_remediation.go` (`markIssuesSkipped`), `cli/cmdpull/pull_path_heal.go` (`AutoHealRepoPathInDB`), `cli/committransfer/pr_engine.go` (`insertPRRecords`), `cli/cmdos/os_machine_alias.go` (`applyOSHostnameChange`), and `cli/cmdrelease/releaseremovebetween_exec.go` (`optimizeRepositoryPostPrune`).
2. Applied concrete remediation to the target files:
   - `cli/cmdpull/pull_remediation_ssh.go`: Converted all void functions to return `error`. Zero error swallowing.
   - `cli/cmdpull/pull_remediation.go`: Handled all returned errors in `applyPreemptiveCredentialFix`, `handleSuccessfulRemediation`, `tryRetryAfterCredentialRemediation`, and `handleSuccessfulRemediationSSH`. Centralized `execPullAutoMerge` with `gitutil.RunGitCombinedInDir`.
   - `cli/cmdworkdir/workdir_cmd.go`: Converted `ensureWorkDirsTableExists()` to return `error` and handled in `runWorkDir`.
3. Verified zero regression across all test suites.
