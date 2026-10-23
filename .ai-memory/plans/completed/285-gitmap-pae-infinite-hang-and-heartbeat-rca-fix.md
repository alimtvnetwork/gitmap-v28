# Plan: 285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix

## Metadata
- **Task ID**: `285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix`
- **Specification**: `02-spec/21-app/285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix/spec.md`
- **Total Steps**: 2 Subtasks across 1 Wave (A=2, H=2)
- **Status**: IN_PROGRESS

---

## User Request (Verbatim)

```text
gitmap pae bug , find the root cause and fix and rleae with a minor bump please
```

---

## Subtask Breakdown

### Subtask-01: PAE Bar Teardown, Bounded Ignore Collect & Non-Interactive Remediation
- **Owner**: Worker 01
- **Files Owned**:
  - `cli/cmdpull/pull_efficient.go`
  - `cli/cmdpull/pull_efficient_types.go`
  - `cli/cmdpull/pull_ignore_scan.go`
- **Actions**:
  1. In `pull_efficient.go`, immediately call `stopHeartbeat()` and `bar.Stop()` as soon as `executePull` returns, prior to any ignore collection or summary rendering.
  2. In `pull_ignore_scan.go`, update `IgnoreScanHandle.Collect()` to use `select` with a 3-second timeout (`time.After(3 * time.Second)`), preventing indefinite channel blocking.
  3. In `pull_efficient_types.go`, ensure `collectAndRemediateEfficientIgnore` does not block on interactive stdin prompts during autonomous `gitmap pae` execution.
  4. Enforce line counts <= 300 across all modified Go files.

### Subtask-02: Heartbeat Self-Termination, Ignore Scan Hardening & Minor Bump Preparation
- **Owner**: Worker 02
- **Files Owned**:
  - `cli/cmdpull/pull_heartbeat.go`
  - `cli/cmdpull/pull_concurrency.go`
- **Actions**:
  1. In `pull_heartbeat.go`, implement self-termination in `runHeartbeatLoop`: if progress indicator signals 100% completion with 0 active workers, stop the ticker and exit immediately.
  2. In `pull_concurrency.go`, add panic recovery to `runThrottledIgnoreScanWorkflow` to guarantee `handle.done` always receives a result.
  3. Deduplicate `antigravity_resume_task.json` in default ignore checks to reduce subprocess calls.
  4. Enforce line counts <= 300 across all modified Go files.
