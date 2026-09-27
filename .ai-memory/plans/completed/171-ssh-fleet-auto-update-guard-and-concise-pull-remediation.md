# Plan: 171 — SSH Fleet GitMap Auto-Update Guard & Concise Pull Remediation

## Spec Reference
[02-spec/21-app/171-ssh-fleet-auto-update-guard-and-concise-pull-remediation.md](../../../02-spec/21-app/171-ssh-fleet-auto-update-guard-and-concise-pull-remediation.md)

## Status
`completed`

## Traceability & Deliverables
- **Task-01:** Remote Node Gitmap Version Verification & Auto-Update Guard
- **Task-02:** Local & Fleet Pull-All Active-Only Concise Summary & Remediation Gate
- **Task-03:** Minor SemVer Release (v6.350.0)
- **Task-04:** Remote CI/CD Pipeline Monitoring via `gitmap pe`

## Execution Summary
1. `cli/cmdssh/ssh_pull_fleet.go`: Added `ensureRemoteNodeGitmap` and `runRemotePullJSON` to probe remote node OS, verify GitMap version, trigger remote update when outdated (< 6.349.0) or uninstalled, and automatically retry on `-json` unrecognized flag errors.
2. `cli/cmdpull/pull.go`: Added guard in `handlePullRemediation` for `opts.all && !opts.showStatus` to skip interactive remediation prompt during fast batch pull operations.
3. `cli/cmdpull/pull_efficient.go`: Removed transport rewrite `maybeApplyTransportToRecords` during active batch pull, preventing origin URL logging spam.
4. `cli/cmdssh/ssh_pull_fleet_test.go`: Added `TestResolveTargetNodeOS` unit test.
5. All 18 CI unit tests and linters passed clean.
