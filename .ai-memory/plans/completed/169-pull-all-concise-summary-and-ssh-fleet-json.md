# Plan 169: Pull-All Concise Summary Filtering, Noise Elimination & SSH Fleet JSON Orchestration

- **Status:** `completed`
- **Spec Reference:** [02-spec/21-app/169-pull-all-concise-summary-and-ssh-fleet-json.md](../../../02-spec/21-app/169-pull-all-concise-summary-and-ssh-fleet-json.md)
- **Completed Steps:** 2 subtasks consolidated

---

## 1. Executive Summary & Problem Context

When users executed `gitmap pa` / `gitmap pa --ssh`, several critical usability and orchestration regressions occurred:
1. **Summary Table Clutter:** Local `gitmap pa` across 64 repositories printed 60 verbose `up-to-date` rows, making it tedious to spot actual repository updates or failures.
2. **Diagnostic Log Flooding:** `maybeApplyTransportToRecords` invoked `ApplyTransportFlag` across all 64 repositories during batch pulls, printing 64 lines of `→ remote.origin.url: ... → ...` to the terminal.
3. **SSH Fleet Delegation Missing on Pull-All:** `gitmap pa --ssh` treated `--ssh` as a local transport coercion flag instead of delegating across SSH fleet nodes.

---

## 2. Consolidated Implementation & Architecture

### 2.1 Local Concise Filtering & Noise Removal (`Task-01`)
- **Noise Elimination:** In `cli/cmdpull/pull.go`, `executePullBatchLifecycle` guards `maybeApplyTransportToRecords` with `if !opts.all`, ensuring batch pull-all operations never rewrite remote origin URLs.
- **Concise Active Filtering:** In `cli/cmdpull/pull_efficient_render.go`, `RenderConciseActiveResultsTo` skips repositories with status label `up-to-date`. If all repositories are up-to-date, it renders `(all repositories are up-to-date)`.
- **Informative Summary Counts:** In `cli/cmdpull/pull.go`, `printPullAllFastSummary` reports active vs up-to-date counts:
  `✓ Pull all complete: 64 pulled (2 active, 62 up-to-date) (1.4s)`.

### 2.2 SSH Fleet Delegation & JSON Protocol (`Task-02`)
- **Fleet Interception:** `runPull` inspects `isPullAll && hasSSHFleetFlag(args)` and dispatches immediately to `RunRemoteSSHPullAllFleetFn` via `handleSSHFleetPullAll(args)`.
- **Enqueue Banner:** `printFleetEnqueueBanner` in `cli/cmdssh/ssh_pull_fleet.go` displays:
  ```text
  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [alpha-win] (10.20.0.11): Enqueued (async)
    • Remote Node [beta-linux] (10.20.0.12): Enqueued (async)
    • Local VM (127.0.0.1 - localhost): Running locally
  ```
- **Concurrent Execution:** Dispatches `gitmap pa --json` over SSH to all remote nodes in parallel with local VM execution (`exec.Command(selfExe, "pa", "--json")`).
- **Structured Aggregation:** Deserializes JSON payloads into `fleetPullJSONSummary`, filters active states per node, and renders a unified terminal summary table or combined JSON (`--json`).

---

## 3. Verification & Quality Gates

- `python linter-scripts/check-nested-ifs.py --changed-only`: PASSED (0 violations).
- `python linter-scripts/check-enum-and-boolean.py`: PASSED (zero explicit booleans, inverted checks, or unhandled enums).
- `go test ./cmdpull`: PASSED.
- `go test ./cmdssh`: PASSED.
