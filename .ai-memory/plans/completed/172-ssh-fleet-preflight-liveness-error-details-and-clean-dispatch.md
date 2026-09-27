# Plan: 172 — SSH Fleet Pre-Flight Liveness, Error Details, and Clean Dispatch

## Status
`completed`

## Spec Reference
`02-spec/21-app/172-ssh-fleet-preflight-liveness-error-details-and-clean-dispatch.md`

## Tasks & Verification Checklist
- ✅ 1. Silence stack trace dumps on offline/unreachable SSH nodes in `cli/cmdssh/ssh_target_nodes.go`.
- ✅ 2. Preserve `lastFail` git failure output, branch, exit code, and diagnosis in `cli/cloner/safe_pull.go`.
- ✅ 3. Enrich `PullRepoSummaryItem` and `fleetPullRepoItem` with `errorDetails` and `remediationHint` JSON fields.
- ✅ 4. Create `cli/cmdpull/pull_remediation_hint.go` to classify failure causes and generate actionable CLI commands.
- ✅ 5. Populate `ErrorDetails` and `RemediationHint` in `buildPullRepoSummaryItems` (`cli/cmdpull/pull_batch_json.go`).
- ✅ 6. Format `↳ Reason: ...` and `↳ Next Step: ...` in `RenderConciseActiveResultsTo` (`cli/cmdpull/pull_efficient_render.go`).
- ✅ 7. Implement concurrent pre-flight machine probing (`probeFleetLiveness`) in `cli/cmdssh/ssh_pull_fleet.go`.
- ✅ 8. Format pre-flight banner declaring local machine execution, online remote nodes enqueued, and offline nodes skipped.
- ✅ 9. Dispatch pull tasks ONLY to online remote connections and skip offline nodes without connecting.
- ✅ 10. Auto-recover remote nodes with stale pending tasks via `recoverRemotePendingTask` in `cli/cmdssh/ssh_pull_fleet.go`.
- ✅ 11. Format failure reasons and `↳ Next Step: ...` in `renderActiveStatesList` in `cli/cmdssh/ssh_pull_fleet.go`.
- ✅ 12. Run targeted linters (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`) with 0 errors.
- ✅ 13. Author Spec 172.
- ✅ 14. Execute minor release (`v6.351.0`) and verify CI/CD pipelines via `gitmap pe`.
