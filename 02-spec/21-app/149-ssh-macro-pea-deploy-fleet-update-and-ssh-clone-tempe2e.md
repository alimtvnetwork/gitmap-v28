# Specification 149: SSH Macro/PEA/PEAT Fleet Deployment, Multi-Node Update Telemetry, Remote SSH Clone & Isolated Temporary E2E Validation

> **Spec Status:** Active  
> **Traceability IDs:** Task-01, Task-02, Task-03, Task-04, Task-05  
> **Canonical Path:** `02-spec/21-app/149-ssh-macro-pea-deploy-fleet-update-and-ssh-clone-tempe2e.md`  
> **Parent Spec:** `02-spec/21-app/01-index.md`

---

## 1. Domain Architecture & System Context

This specification formalizes the SSH fleet deployment commands, remote package update JSON telemetry and inventory tables, remote SSH repository cloning with default workdir resolution, and the `//go:build tempe2e` isolated validation harness:

1. **Macro, PEAT & PEA SSH Fleet Deployment (`gitmap macro|peat|pea deploy ssh`)**:
   - Supports `gitmap macro deploy ssh --except id,ip,alias`, `gitmap peat deploy ssh --except id,ip,alias`, and `gitmap pea deploy ssh --except id,ip,alias` (including `--excep` and `--exclude` shorthands).
   - Distributes macro payloads and pipeline-error-ai / tasks configurations concurrently across all active cluster nodes while excluding matching IDs, IPs, or aliases.
   - Renders an aligned terminal summary table detailing per-node status, macro counts, duration, and excluded counts.

2. **Fleet-Wide Multi-App Update (`gitmap update --all` / `gitmap update all` / `gitmap ua`)**:
   - Supports `gitmap update --all`, `gitmap update all`, and `gitmap ua` with `--except`/`--excep` node exclusion.
   - Invokes remote `gitmap` update workflows across SSH nodes in parallel, parses structured `FleetUpdateTelemetry` JSON responses from each machine, and renders a unified summary table.

3. **Targeted App Update (`gitmap update <name> --except`) & Fleet Inventory (`gitmap update ls`)**:
   - Supports `gitmap update <name> --except id,alias,ip` (and `--excep`) to update a specific installed component across the fleet.
   - Supports `gitmap update ls` to query `FleetNodeInventory` JSON payloads across all SSH nodes in parallel and display all installed items in a formatted terminal table.

4. **Remote SSH Repository Clone (`gitmap ssh clone` / `ssh-clone` / `ssh-c`)**:
   - Supports `gitmap ssh clone <repo-name|url> [git|path]`, `gitmap ssh-clone`, and `gitmap ssh-c`.
   - When the repo argument is omitted or `.` / `git` is passed inside a tracked repository, automatically resolves the current repository's `origin` remote URL.
   - When the destination path is omitted or `git` is specified, clones into the remote node's default work directory (`~\git\<repo>` on Windows or `~/git/<repo>` on Linux).

5. **Temporary End-to-End Isolation Standard (`//go:build tempe2e`)**:
   - Implements `cli/tests/e2e/ssh_fleet_deploy_update_clone_tempe2e_test.go` guarded by both `//go:build tempe2e` and runtime `if os.Getenv("RUN_TEMP_E2E") != "1" { t.Skip(...) }`.
   - Guarantees zero CI/CD execution and zero impact on routine local `go test ./...` runs.

---

## 2. Verification & Acceptance Criteria

### AC-APP-149-001: Macro, PEAT, and PEA Exclusion Filtering
**Given** Joined SSH nodes `w1 (192.168.1.3)`, `w2 (192.168.1.7)`, and `w3 (192.168.1.12)`.  
**When** `gitmap macro deploy ssh --except w3,192.168.1.12` or `gitmap peat deploy ssh --excep w2` is executed.  
**Then** Excluded nodes are skipped, active nodes receive parallel deployments, and the summary table reports exact succeeded and excluded counts.

### AC-APP-149-002: Fleet Update & Inventory JSON-to-Table Rendering
**Given** Joined SSH nodes returning JSON `FleetUpdateTelemetry` and `FleetNodeInventory`.  
**When** `gitmap ua --except w3`, `gitmap update gitmap --excep w2`, or `gitmap update ls` is executed.  
**Then** JSON payloads are parsed and rendered into an aligned terminal summary table.

### AC-APP-149-003: Skip-by-Default Temporary E2E Quarantine
**Given** `cli/tests/e2e/ssh_fleet_deploy_update_clone_tempe2e_test.go`.  
**When** Executed with `RUN_TEMP_E2E=1 go test -tags=tempe2e -v ./tests/e2e/... -run TempE2E`.  
**Then** All 3 E2E suites pass, and when `RUN_TEMP_E2E` is unset, all 3 tests skip automatically.

---

_Specification last updated: 2026-09-24_
