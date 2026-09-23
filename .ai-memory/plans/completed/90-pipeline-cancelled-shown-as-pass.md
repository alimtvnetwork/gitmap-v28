# Completed Plan: Pipeline Cancelled / Timed Out Runs Erroneously Displayed as PASS

Spec Reference: [02-spec/22-app-issues/37-pipeline-cancelled-shown-as-pass.md](../../../02-spec/22-app-issues/37-pipeline-cancelled-shown-as-pass.md)

## User Request (Verbatim)

```text
is it done and released properly

https://prnt.sc/p1KnI6ogAmUU

Failed pipeline shows as PASS fix please correctly, test things properly You can check for 

Antigravity-Manager with commit hash 12e40b1


check the issue properly 


And fix it please
```

## Problem Summary & Root Cause
In `gitmap pe` / `gitmap pipeline error-logs` and `gitmap pipeline history`, commits with cancelled or timed-out workflow runs (such as `12e40b1` and `5717700` in `Antigravity-Manager`) were displayed with `Status: PASS` and `Failures: 0`.
Root causes:
1. `isWorkflowFailure` in `cli/cmdpipeline/pipeline_commit_groups.go` did not include `cancelled`, causing `group.PassedWorkflows++` to increment and `computeAggregateConclusion` to return `"success"`.
2. `formatStatusBadge`, `formatCommitGroupBadge`, and `formatWorkflowTreeBadge` checked only `"failure"`, bypassing cancelled conclusions.
3. `isCommitGroupFailure` strictly checked `== "failure"`, treating cancelled groups as clean.
4. `formatGroupWorkflowsSummary` formatted workflows in arbitrary order, hiding failing workflows behind passing ones in `(+N)`.

## Implemented Remediation
1. **Unify Failure Predicate**:
   - `cli/cmdpipeline/pipeline_logs.go`: Extended `isFailingConclusion` with `action_required` and `stale` alongside `cancelled`, `failure`, `timed_out`, and `startup_failure`.
   - `cli/cmdpipeline/pipeline_commit_groups.go`: Updated `isWorkflowFailure` to invoke `isFailingConclusion(wf.Conclusion)`.
2. **Status Badges & Group Rollup**:
   - `cli/cmdpipeline/pipeline_history.go`: Updated `formatStatusBadge` and `isCommitGroupFailure` to check `isFailingConclusion`.
   - `cli/cmdpipeline/pipeline_history_cmd.go`: Updated `formatCommitGroupBadge` and `formatWorkflowTreeBadge` to check `isFailingConclusion`.
   - `cli/cmdpipeline/pipeline_status.go`: Updated `renderCompletedStatusLine` to check `isFailingConclusion`.
3. **Workflow Summary Prioritization**:
   - `cli/cmdpipeline/pipeline_history.go`: Implemented `prioritizeGroupWorkflows` to place failing/cancelled workflows ahead of passing ones so they remain visible in truncated column summaries. Fixed remaining unprinted counter in `finalizeTruncatedSummary`.
4. **CI Linter Hygiene**:
   - Removed unused functions `executeFleetUpdate` in `cli/cmdssh/ssh_update_remote.go` and `ensureDelegateInstalledJSON` in `cli/cmdssh/sshexec.go`.

## Verification Outcomes
- Unit tests added:
  - `TestGroupRunsByCommit_CancelledWorkflows` in `cli/cmdpipeline/pipeline_commit_groups_test.go`
  - `TestFormatStatusBadge_CancelledAndFailing`, `TestFormatGroupWorkflowsSummary_PrioritizesFailures`, `TestIsCommitGroupFailure_Cancelled` in `cli/cmdpipeline/pipeline_history_test.go`
- Live verification:
  - `gitmap pe -n 5` and `gitmap pipeline history -n 5` accurately show `FAIL` (red) and non-zero failures for commits with cancelled workflows.
- Quality gates:
  - `python 03-ai-scripts/26-go-code-formatter.py`: 3594 Go files verified clean.
  - `python linter-scripts/check-nested-ifs.py`: 0 nested if violations.
  - `python linter-scripts/check-enum-and-boolean.py`: 0 boolean/enum violations.
  - `python linter-scripts/check-error-management.py`: 0 bare panic/exit violations.
  - `python linter-scripts/check-relative-paths.py`: 0 absolute path / URI violations across 7536 files.
