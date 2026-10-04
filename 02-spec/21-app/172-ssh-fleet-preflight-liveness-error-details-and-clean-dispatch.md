# Specification: 172 — SSH Fleet Pre-Flight Liveness, Error Details, and Clean Dispatch

## Status
`active`

## Domain
`cli / ssh fleet / preflight liveness / error details / actionable remediation / safe-pull diagnosis / self-healing pending tasks`

## User Request (Verbatim)
```text
PS [REPO_ROOT]> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [alpha-win] (10.20.0.11): Enqueued (async)
    • Remote Node [beta-linux] (10.20.0.12): Enqueued (async)
    • Remote Node [gamma-mac] (10.20.0.13): Enqueued (async)
    • Remote Node [w1] (node-w1): Enqueued (async)
    • Remote Node [w2] (node-w2): Enqueued (async)
    • Remote Node [w3] (node-w3): Enqueued (async)
    • Local VM (127.0.0.1 - localhost): Running locally

  [alpha-win|10.20.0.11] Offline: [E9000:EXECUTION] execution: node [alpha-win|10.20.0.11] is unreachable: connection timed out (at=cmdssh/ssh_target_nodes.go:28)
  Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.checkRemoteNodeOnline (cmdssh/ssh_target_nodes.go:28)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeRemoteNodePull (cmdssh/ssh_pull_fleet.go:115)
...
```

## Architectural Context & Blast Radius
1. **Pre-Flight Fleet Machine Liveness Probing:**
   - In `cli/cmdssh/ssh_pull_fleet.go`, `probeFleetLiveness` probes port 22 TCP connectivity across all configured fleet connections concurrently before task dispatch.
   - Enqueue banner partitions nodes deterministically:
     - Online remote nodes: `Remote Node [alias] (IP): Online → Enqueued (async)`
     - Offline remote nodes: `Remote Node [alias] (IP): Offline (skipped, no task enqueued)`
     - Local machine: `Current Machine [hostname (127.0.0.1)]: Running locally (direct execution, not enqueued)`
   - Offline nodes are marked `IsSkipped: true` and are NEVER connected to via SSH or dispatched to goroutines.

2. **Total Stack Trace Suppression on Known Network States:**
   - In `cli/cmdssh/ssh_target_nodes.go`, unreachable nodes during `checkRemoteNodeOnline` return `false` silently without instantiating `apperror.NewExecutionError` or dumping 15-line stack traces to stdout/stderr.

3. **Safe-Pull Diagnostic Error Preservation:**
   - In `cli/cloner/safe_pull.go`, `safePullRepoWithProgress` captures and preserves `lastFail` across retry attempts, retaining the git stderr, exit code, branch name, and diagnosis instead of discarding them with a generic error.

4. **Error Details & Actionable Remediation in Terminal & JSON:**
   - In `cli/cmdpull/pull_remediation_hint.go`, `ResolvePullErrorDetails` classifies failure strings into human-readable summaries and `ResolvePullRemediationHint` generates actionable CLI commands.
   - `PullRepoSummaryItem` and `fleetPullRepoItem` include `ErrorDetails` and `RemediationHint`.
   - `RenderConciseActiveResultsTo` and `renderActiveStatesList` output `↳ Reason: <errorDetails>` and `↳ Next Step: <remediationHint>` for failed/dirty repositories.

5. **Remote Stale Pending Task Auto-Recovery:**
   - In `cli/cmdssh/ssh_pull_fleet.go`, `recoverRemotePendingTask` extracts lingering pending task IDs (`Id <N>`) from remote command output, executes `gitmap task cancel <N>`, and retries `gitmap pa --json` automatically.

## Verification Gates & Invariants
- `check-nested-ifs.py`: 0 nested if violations.
- `check-boolean-guidelines.py`: 0 boolean guideline violations.
- `check-relative-paths.py`: 0 absolute paths.
- Function sizing: all functions <= 15 lines.
- Version bump: SemVer minor release `v6.351.0`.
- GitHub Actions CI/CD workflows run green without regressions.
