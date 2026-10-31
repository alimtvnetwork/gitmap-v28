# Subtask 03: Cross-Node Fleet Software Updates (Cursor, VS Code, AGM, AGY Tools)

> **Subtask ID:** `364-03`  
> **Parent Plan:** [364-gitmap-command-enhancements.md](../../completed/364-gitmap-command-enhancements.md)  
> **Specification Reference:** [03-fleet-software-update-spec.md](../../../../02-spec/21-app/364-gitmap-command-enhancements/03-fleet-software-update-spec.md)  
> **Status:** `completed`  
> **Creation Date:** 31-Oct-26  
> **Lead Architecture:** Swarm Subagent 3 (Fleet Software Update Spec Author)  

---

## 1. Context & Objective

Development teams operating GitMap across hybrid multi-machine fleets (Windows developer workstations and Ubuntu/Debian remote nodes) require a frictionless mechanism to keep core developer tools updated. Currently, updating tools like Cursor IDE requires manual interaction on each node, leading to version drift and workflow disruptions.

This subtask implements centralized fleet update capabilities for:
1. **Cursor IDE:** `gitmap update cursor [node|all]` and `gitmap node update cursor [node|all]`
2. **Visual Studio Code:** `gitmap update vscode [node|all]` and `gitmap node update vscode [node|all]`
3. **Antigravity Manager:** `gitmap update agm [node|all]` and `gitmap node update agm [node|all]`
4. **Antigravity Tools:** `gitmap update agy [node|all]` / `gitmap update agy-tools [node|all]` and `gitmap node update agy [node|all]`

---

## 2. Work Breakdown Structure (WBS)

### Phase 1: CLI Dispatch & Routing Alignment
- **Target Files:**
  - `cli/cmdupdate/fleet_dispatch.go`
  - `cli/cmdupdate/update.go`
  - `cli/cmdnodes/nodes_cmd.go`
- **Deliverables:**
  - Extend `cmdupdate.IsFleetUpdateCommand` to recognize application update tokens (`cursor`, `vscode`, `agm`, `agy`, `agy-tools`).
  - Wire `gitmap update <app> [node|all]` to the fleet software update handler.
  - Wire `gitmap node update <app> [node|all]` within the `cmdnodes` command tree to invoke the shared engine.
  - Implement flag parser for `--force` (`-f`), `--dry-run` (`-d`), `--skip-offline` (`-s`), `--parallel` (`-p`), `--timeout` (`-t`), and `--json` (`-j`).

### Phase 2: Cross-Platform Update Execution Adapters
- **Target Files:**
  - `cli/cmdupdate/fleet_app_adapters.go`
  - `cli/cmdcursor/cursor_install.go`
  - `install.ps1`
  - `install.sh`
- **Deliverables:**
  - **Windows Adapter:**
    - Cursor: `winget upgrade Anysphere.Cursor --silent --accept-source-agreements --accept-package-agreements` with fallback to `install.ps1 -UpdateCursor`.
    - VS Code: `winget upgrade Microsoft.VisualStudioCode --silent --accept-source-agreements --accept-package-agreements`.
    - AGM: `install.ps1 -UpdateAGM` with process termination and binary swap.
    - AGY Tools: `install.ps1 -UpdateAGY`.
  - **Linux / Ubuntu Adapter:**
    - Cursor: `python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py --update` with fallback to `bash install.sh --update-cursor`.
    - VS Code: `sudo apt-get update -qq && sudo apt-get install --only-upgrade -y code`.
    - AGM: `bash install.sh --update-agm` and user systemd daemon reload.
    - AGY Tools: `bash install.sh --update-agy` or `pip3 install --upgrade antigravity-tools`.

### Phase 3: Concurrency Engine & Worker Pool
- **Target Files:**
  - `cli/cmdupdate/fleet_app_worker.go`
  - `cli/cmdupdate/fleet_app_types.go`
- **Deliverables:**
  - Define `AppUpdateType`, `FleetUpdateOptions`, `NodeUpdateResult`, and `FleetUpdateReport` structs.
  - Implement non-blocking TCP socket probe (port 22, timeout 1.5s) to detect unreachable nodes upfront.
  - Launch bounded goroutine worker pool (default concurrency: 4 workers).
  - Collect results with thread-safe mutex and compile comprehensive execution telemetry.

### Phase 4: Telemetry Aggregation & Terminal UI
- **Target Files:**
  - `cli/cmdupdate/fleet_app_ui.go`
  - `cli/cmdupdate/fleet_telemetry.go`
- **Deliverables:**
  - Interactive Lipgloss box table renderer displaying Node Alias, Host, Target OS, Status, Previous Version, Updated Version, and Latency.
  - JSON output serializer when `--json` flag is supplied.
  - Suggestion docket generator for failed or unreachable nodes providing remediation guidance.

### Phase 5: Documentation & Help Text Parity
- **Target Files:**
  - `cli/helpdoc/commands.go`
  - `src/data/commands.ts`
  - `02-spec/21-app/363-gitmap-command-enhancements/03-fleet-software-update-spec.md`
- **Deliverables:**
  - Add help entries for `gitmap update cursor`, `vscode`, `agm`, `agy` and their `node update` counterparts.
  - Update web dashboard command catalog in `src/data/commands.ts`.
  - Generate shell completions for newly added command routes.

---

## 3. Invariants & Coding Guidelines Checklist

| Standard | Rule Requirement | Compliance Verification |
| :--- | :--- | :--- |
| **Error Return Sovereignty** | Leaf functions return `*apperror.AppError` directly. Never call `cliexit.HandleError` or `os.Exit`. | Verified in `fleet_app_worker.go` & `fleet_app_adapters.go` |
| **Positive Booleans Only** | All booleans use `is*` or `has*` prefix (`isForce`, `isDryRun`, `isSkipOffline`, `isJSON`, `isSuccess`, `hasError`). | Verified across all Go structs and functions |
| **15-Line Function Limit** | Functions <= 8 lines preferred, <= 15 lines max. Decompose large routines into discrete helpers. | Verified via static decomposition |
| **Zero Nested Conditionals** | Flatten all nested conditionals with guard clauses and early returns. | Verified across dispatchers and parsers |
| **100% Relative Paths** | Zero hardcoded absolute paths (`C:\...`, `/home/...`, `file:///...`) in code, specs, and plans. | Verified in markdown and Go code |
| **Zero-Build & Zero-Test** | Targeted file checks only during normal iterations; no continuous heavy builds or test loops. | Verified; terminal release step handles full E2E |

---

## 4. Acceptance Criteria & Verification Plan

- [ ] **TC-01: Direct Local Cursor Update**
  - Command: `gitmap update cursor`
  - Expectation: Detects local OS, executes Windows winget or Linux setup script, reports status.
- [ ] **TC-02: Fleet-Wide Cursor Update**
  - Command: `gitmap update cursor all`
  - Expectation: Concurrently dispatches update across all fleet nodes in parallel worker pool.
- [ ] **TC-03: Single Remote Node Cursor Update**
  - Command: `gitmap update cursor u1`
  - Expectation: Connects via SSH to `u1`, executes Linux updater, reports version before/after.
- [ ] **TC-04: VS Code Fleet Update**
  - Command: `gitmap update vscode all`
  - Expectation: Updates VS Code across all nodes using respective platform managers.
- [ ] **TC-05: AGM Fleet Update**
  - Command: `gitmap update agm all`
  - Expectation: Triggers binary swap and daemon restart for Antigravity Manager across nodes.
- [ ] **TC-06: AGY Tools Update**
  - Command: `gitmap update agy all` and `gitmap update agy-tools all`
  - Expectation: Updates Antigravity CLI and environment tools.
- [ ] **TC-07: Node Subcommand Parity**
  - Command: `gitmap node update cursor all`
  - Expectation: Yields identical behavior and telemetry to `gitmap update cursor all`.
- [ ] **TC-08: Offline Node Bypassing**
  - Command: `gitmap update cursor all --skip-offline`
  - Expectation: Offline machines are detected within 1.5s TCP probe and skipped without error.
- [ ] **TC-09: Dry-Run Simulation**
  - Command: `gitmap update cursor all --dry-run`
  - Expectation: Displays exact remote commands that would be executed without invoking them.
- [ ] **TC-10: JSON Telemetry Emission**
  - Command: `gitmap update cursor all --json`
  - Expectation: Outputs valid JSON conforming to `FleetUpdateReport` schema.

---

## 5. Execution Dependencies

- **Pre-requisite:** Subtask 01 (`01-alias-and-copy-move`) establishes stable CLI dispatch trees.
- **Pre-requisite:** Subtask 02 (`02-deploy-and-os-compat`) establishes fast reachability probe utilities (`ProbeNodeReachability`).
- **Dependent:** Subtask 04 (`04-release-ceremony`) requires Subtask 03 documentation and command definitions before cutting `v6.548.0`.
