# Subtask 02: Heartbeat Self-Termination & Ignore Scan Hardening

## Status: DONE
- **Parent Plan**: `.ai-memory/plans/285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix.md`
- **Owner**: Worker 02
- **Owned Files**:
  - `cli/cmdpull/pull_heartbeat.go`
  - `cli/cmdpull/pull_concurrency.go`

## Objectives
1. In `cli/cmdpull/pull_heartbeat.go`:
   - In `runHeartbeatLoop`:
     On every ticker event, inspect the progress message or state. If it indicates `0 active workers` and all completed (`isProgressCompleted` verifying `0 active workers` and `completed >= total` or completed substring), terminate the loop immediately and do not emit trailing heartbeats. (DONE)
2. In `cli/cmdpull/pull_concurrency.go`:
   - In `runThrottledIgnoreScanWorkflow`:
     Added `defer func() { if r := recover(); r != nil { handle.done <- nil } }()` so `handle.done` is always guaranteed to send even if an unhandled panic occurs during scan. (DONE)
   - Deduplicated repo records by canonical path in `filterColdIgnoreRecords` and sanitized ignore check path lists via `sanitizeIgnoreIssue` and `deduplicateStringList`. (DONE)
3. FastGate: Ensure all modified Go files stay strictly <= 300 lines with nesting depth <= 1:
   - `cli/cmdpull/pull_heartbeat.go`: 98 lines (depth <= 1)
   - `cli/cmdpull/pull_concurrency.go`: 136 lines (depth <= 1)
   - FastGate 100% PASS in 0.011s. (DONE)
