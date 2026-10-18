# SSH Fleet Pull UX & Local Node Detection Specification

## 1. Overview & Problem Definition

When executing `gitmap pa --ssh` (or `gitmap pull-all --ssh`), two issues were identified:
1. **Terminal Ticker Spam**:
   During long-running fleet pull operations, the heartbeat ticker printed on a new line every 5 seconds after 30 seconds (`[⏳ 30s elapsed] Pulling across SSH fleet (2 online nodes, 1 local)...`), resulting in up to 80+ redundant lines flooding the terminal over a 400s operation.
2. **False Online Reporting & Loopback Pull of Local Host ("main machine was off")**:
   When an SSH connection is registered in the database pointing to the local machine's IP address (e.g. `192.168.1.20`), GitMap previously probed and connected to that IP over SSH without checking if it matched any of the current machine's network interface IPs. Consequently:
   - GitMap connected to itself over SSH, ran a redundant secondary `pull-all` process against the same local repositories, and reported duplicate results.
   - If the user had powered off a physical machine previously associated with that IP, GitMap falsely reported the node as `Online → Enqueued (async)` and `66 pulled` because the current host was responding to port 22 on its own interface.

## 2. Technical Architecture & Invariants

### 2.1 Terminal Interface In-Place Rendering
- **TTY Mode (`termout.IsStdoutTTY() == true`)**:
  - Heartbeat status updates MUST render in-place using carriage return and line erase: `\r\033[K    [⏳ %ds elapsed] Pulling across SSH fleet (%d online nodes, 1 local)...`.
  - Upon task completion or cancellation, the heartbeat line MUST be cleared cleanly via `\r\033[K` so final output headers do not collide with leftover ticker text.
- **Non-TTY Mode (`termout.IsStdoutTTY() == false`)**:
  - The heartbeat ticker MUST throttle output to 30-second intervals (e.g. 30s, 60s, 90s) with newline formatting to avoid log explosion in CI or automated script logs.

### 2.2 Local Interface IP Detection
- Function: `IsLocalHostAddress(ip string) bool`
  - Evaluates `127.0.0.1`, `::1`, `localhost`, and the current machine's hostname.
  - Enumerates `net.Interfaces()` and inspects all unicast IPv4 and IPv6 interface addresses (`net.IPNet` / `net.IPAddr`).
- Fleet Liveness & Enqueue Behavior:
  - When probing fleet connections, if `IsLocalHostAddress(targetConn.IPAddress)` evaluates to `true`:
    - The node is classified as the local machine.
    - It MUST NOT be enqueued into the async remote SSH worker pool.
    - The enqueue banner displays:
      `• Remote Node [%s] (%s): Local machine (covered by local execution, skipped)`
    - This eliminates loopback SSH execution, prevents lock contention on local repos, and prevents false "Online" claims when an external machine is offline.

### 2.3 Modular File Size Hygiene (<300 Lines)
To comply strictly with the 300-line limit enforced by `52-file-size-check.py`, `cli/cmdssh/ssh_pull_fleet.go` (693 lines) is decomposed into 5 focused modules:
1. `cli/cmdssh/ssh_pull_fleet.go` (~120 lines): Core entrypoint `RunSSHPullAllFleet`, options, error tallying.
2. `cli/cmdssh/ssh_pull_fleet_heartbeat.go` (~90 lines): Dynamic heartbeat loop, TTY in-place carriage return rendering, non-TTY 30s throttling.
3. `cli/cmdssh/ssh_pull_fleet_liveness.go` (~180 lines): Network interface resolution, `IsLocalHostAddress`, node version queries, probe orchestration.
4. `cli/cmdssh/ssh_pull_fleet_dispatch.go` (~170 lines): Local execution and bounded remote worker dispatch.
5. `cli/cmdssh/ssh_pull_fleet_outcome.go` (~190 lines): Summary extraction, active state filtering, terminal summary output formatting.
