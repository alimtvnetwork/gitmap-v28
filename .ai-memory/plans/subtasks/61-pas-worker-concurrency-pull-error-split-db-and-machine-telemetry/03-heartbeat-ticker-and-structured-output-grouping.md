# Subtask 03: Heartbeat Ticker & Structured Output Grouping

> **Parent Plan:** [61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry](../../completed/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)
> **Tracking Spec:** [198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)
> **Primary File Targets:** `cli/cmdpull/pull_heartbeat.go`, `cli/cmdpull/pull_efficient.go`, `cli/cmdpull/pull_efficient_render.go`

---

## 1. Objective

Prevent the appearance of silent hangs during long pull runs (>30s) and organize terminal pull results into categorized groups:
1. Emits background heartbeat progress ticks every 5 seconds when execution exceeds 30 seconds.
2. Formats pull summary into 3 clear categories:
   - **Updated Repositories:** Shows commit delta stats (`+N/-M`).
   - **Dirty Repositories:** Displays itemized dirty files (`modified:`, `untracked:`) and remediation commands (`gitmap cpar`, `gitmap stash`).
   - **Failed Repositories:** Displays failure cause and remediation hints.

---

## 2. Implementation Scope

- **`cli/cmdpull/pull_heartbeat.go`:**
  - Implement `StartHeartbeatTicker(threshold time.Duration, interval time.Duration) (*HeartbeatController, func())`.
  - Non-destructive ticker printing `[⏳ 35s elapsed] 14/48 completed...` without mangling active lines.
- **`cli/cmdpull/pull_efficient_render.go`:**
  - Group results into Updated, Dirty, and Failed buckets.
  - Query uncommitted modified and untracked files for dirty repos and render indented file listings.
  - Suggest actionable commands: `gitmap cpar` or `gitmap fix <repo>`.

---

## 3. Verification

- Verify pull results render with clean section headers.
- Verify heartbeat ticker cleans up properly on completion.
- Verify compilation with `go build ./...` in `cli/`.
