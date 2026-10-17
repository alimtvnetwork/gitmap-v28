# Specification: GitMap PAE Post-Pull Infinite Hang & Heartbeat Teardown RCA Fix

## Status
- **Task ID**: `285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix`
- **Specification Version**: 1.0.0
- **Status**: APPROVED
- **Date**: 2026-10-18

---

## 1. User Request (Verbatim)

```text
gitmap pae bug , find the root cause and fix and rleae with a minor bump please
```
*Accompanying Screenshot:* Terminal showing `gitmap pae` stuck at `100% (46/46 repos)` with `0 active workers • 0 failures so far...`, repeatedly emitting heartbeat lines every 5s up to `580s elapsed` (`~10 minutes`) without completing.

---

## 2. Root Cause Analysis (4-Part RCA)

### 2.1 Problem Description
Running `gitmap pae` successfully pulls all active repositories (e.g. 46/46 repos completed). However, instead of exiting cleanly and rendering the final summary, the process hangs indefinitely. The terminal repeatedly prints:
`✓ [████████████████████] 100% (46/46 repos) | ✔ wp-onboarding-v17: up-to-date (580.0s) [! 580s elapsed] 46/46 completed • 0 active workers • 0 failures so far...`

### 2.2 Root Cause Identification
1. **Deferred Lifecycle Teardown in `executeActiveEfficientBatch` (`cli/cmdpull/pull_efficient.go:244-259`):**
   `bar.Stop()` and `stopHeartbeat()` were registered using `defer` to the exit of `executeActiveEfficientBatch`. When `executePull` returned with 100% completion and 0 active workers, neither the progress bar ticker (80ms) nor the heartbeat ticker (5s) was stopped.
2. **Unbounded Blocking Channel in Background Ignore Scan (`cli/cmdpull/pull_ignore_scan.go:48-53`):**
   Directly after `executePull`, line 256 calls `collectAndRemediateEfficientIgnore(ignoreHandle, ...)`, which calls `ignoreHandle.Collect()`. `Collect()` executed an unbounded blocking channel receive `return <-h.done` with zero timeout.
3. **Sequential Background Git Subprocess Stalling (`cli/cmdpull/pull_concurrency.go` & `pull_ignore_scan.go`):**
   The background scan (`scanColdRecordsSequentially`) iterates through cold repositories spawning 5 sequential `git ls-files` subprocesses per repo (up to 230 git process spawns on Windows). While pulling, git index locks (`.git/index.lock`) cause these checks to stall.
4. **Interactive Prompt Terminal Erasure:**
   If `Collect()` did return issues, `handleIgnoreRemediation` called `promptIgnoreRemediationChoice()` reading from `os.Stdin`. Because `bar.Stop()` was deferred, the progress bar ticker continued clearing the line every 80ms (`\r\033[K`), rendering the stdin prompt completely invisible to the user while blocking execution.
5. **Lack of Heartbeat Self-Termination (`cli/cmdpull/pull_heartbeat.go:43-52`):**
   `runHeartbeatLoop` lacked a self-termination guard when `bar.Completed() >= bar.Total() && len(bar.ActiveWorkers()) == 0`.

### 2.3 Proposed Changes & Invariants
1. **Immediate Teardown in `executeActiveEfficientBatch`:** Stop the heartbeat and progress bar immediately when `executePull` finishes, before any post-pull ignore checks or summary rendering.
2. **Bounded Timeout on `Collect()`:** Guard `<-h.done` with a 3-second `select` timeout in `cli/cmdpull/pull_ignore_scan.go`. Background ignore checks must never block a user pull.
3. **Heartbeat Self-Termination:** In `cli/cmdpull/pull_heartbeat.go`, self-terminate the loop when progress reaches 100% and active workers is 0.
4. **Non-Blocking Ignore Remediation in PAE:** Do not prompt interactive stdin during autonomous `gitmap pae` execution; print non-interactive guidance if issues are detected.
5. **Minor Version Bump:** Execute minor version bump to `v6.568.0` and release.

---

## 3. File Ownership & Boundaries

- **Worker 01:**
  - `cli/cmdpull/pull_efficient.go` (immediate teardown of heartbeat and progress bar)
  - `cli/cmdpull/pull_efficient_types.go` (non-interactive ignore handling in PAE)
  - `cli/cmdpull/pull_ignore_scan.go` (add 3s timeout to `Collect()` and panic recovery)
- **Worker 02:**
  - `cli/cmdpull/pull_heartbeat.go` (self-termination condition when 100% completed & 0 active workers)
  - `cli/cmdpull/pull_concurrency.go` (deduplicate ignore check paths, add panic recovery)

---

## 4. Acceptance Criteria
1. `gitmap pae` terminates immediately upon reaching 100% with 0 active workers.
2. Heartbeat loop automatically stops and never prints when work is done.
3. `IgnoreScanHandle.Collect()` never blocks longer than 3 seconds under any circumstances.
4. All modified Go files strictly <= 300 lines with nesting depth <= 1.
5. FastGate passes 100%.
6. Released under minor version bump (`v6.568.0`).
