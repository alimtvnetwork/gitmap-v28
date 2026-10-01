# Subtask 01: Dedicated Machine Telemetry & Startup Node Version

> **Parent Plan:** [61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry](../../completed/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Tracking Spec:** [198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Primary File Targets:** `cli/cmdos/os_machine_alias.go`, `cli/cmdssh/ssh_pull_fleet.go`  

---

## 1. Objective

Implement dedicated `gitmap machine` telemetry and remote node version discovery on fleet pull startup:
1. `gitmap machine` defaults to rendering the local machine's identity card (Alias, Hostname, IP, OS, Architecture, CPU count, User, GitMap Version) without leaking generic command help.
2. Supports `--json` (`-j`) for pure structured JSON output.
3. In `gitmap pa --ssh` and `gitmap pas`, probe and display the exact GitMap binary version for each reachable remote fleet node and local node in the startup enqueue banner.

---

## 2. Implementation Scope

- **`cli/cmdos/os_machine_alias.go`:**
  - Route bare `gitmap machine` and `gitmap machine --json` to dedicated telemetry renderer `renderMachineIdentity`.
  - Format structured table showing node identity, IP, platform, core count, and version.
  - Retain flag parsing for `-h`/`--help` to display help usage when explicitly requested.
- **`cli/cmdssh/ssh_pull_fleet.go`:**
  - Enhance node liveness probe to query `gitmap --version` asynchronously.
  - Format startup banner lines: `• Remote Node [w1] (192.168.1.3) [GitMap v6.442.0]: Online → Enqueued (async)`.
  - Fall back to `[GitMap unknown]` on timeout or connection error.

---

## 3. Verification

- Run `gitmap machine` locally and verify clean table output.
- Run `gitmap machine --json` and verify valid JSON output with `gitmap_version`.
- Verify no compiler regressions with `go build ./...` in `cli/`.
