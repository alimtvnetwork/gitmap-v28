# Subtask 07: SSH Fleet Pull UX & Local Node Detection

- **Status**: COMPLETED
- **Owner**: Lead Orchestrator
- **Spec**: `02-spec/21-app/goemailsender-connection-and-testing/04-ssh-fleet-pull-ux-and-local-node-detection-spec.md`

## Objectives
1. Implement `IsLocalHostAddress(ip string) bool` to discover all local interface IPs (`net.Interfaces()`), localhost, and hostname.
2. In `probeFleetLiveness` and `printFleetEnqueueBanner`, detect nodes matching the current host:
   - Mark as `isLocalNode == true`.
   - Prevent remote SSH dispatch to localhost.
   - Display clear notification: `• Remote Node [%s] (%s): Local machine (covered by local execution, skipped)`.
3. In `startFleetPullHeartbeat`:
   - Detect TTY via `termout.IsStdoutTTY()`.
   - If TTY: render status in-place using `\r\033[K` without spamming newlines.
   - If non-TTY: throttle output to 30-second intervals.
   - On exit: cleanly clear heartbeat ticker line with `\r\033[K`.
4. Decompose `cli/cmdssh/ssh_pull_fleet.go` (693 lines) into 5 focused `.go` files strictly under 300 lines each:
   - `cli/cmdssh/ssh_pull_fleet.go`
   - `cli/cmdssh/ssh_pull_fleet_heartbeat.go`
   - `cli/cmdssh/ssh_pull_fleet_liveness.go`
   - `cli/cmdssh/ssh_pull_fleet_dispatch.go`
   - `cli/cmdssh/ssh_pull_fleet_outcome.go`
5. In `.github/workflows/ci.yml`:
   - Bump `lint_baseline_cache_version` from `v5` to `v6` (and all fallback keys) to resolve stale baseline diff cache and pass CI.
6. Verify locally:
   - Run unit tests in `cli/cmdssh/`.
   - Run `52-file-size-check.py --diff HEAD~1 --limit 300`.
   - Verify `gitmap nodes ping` and `gitmap pa --ssh`.
