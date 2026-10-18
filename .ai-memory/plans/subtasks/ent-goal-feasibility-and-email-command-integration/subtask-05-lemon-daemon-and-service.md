# Subtask 05: Lemon Daemon Polling Loop & Cross-Platform Service Hooks

**Parent Plan:** `.ai-memory/plans/ent-goal-feasibility-and-email-command-integration.md`  
**State:** `COMPLETED`  
**Owned Files:** `cli/cmdlemon/interval_guard.go`, `cli/cmdlemon/lemon_daemon.go`, `cli/cmdlemon/service_manager.go`, `cli/cmdlemon/service_linux.go`, `cli/cmdlemon/service_darwin.go`, `cli/cmdlemon/service_windows.go`, `cli/cmdlemon/service_other.go`, `cli/cmdlemon/lemon_cmd.go`  

## Objectives & Delivered Implementation
1. Implemented Lemon daemon background polling runner in `cli/cmdlemon/lemon_daemon.go`.
2. Enforced polling interval constraints via `cli/cmdlemon/interval_guard.go`:
   - Configurable: default 8 seconds (`DefaultPollingInterval`).
   - Non-negotiable minimum 5 seconds guard (`MinPollingInterval`).
   - If requested interval is $< 5$s (e.g. 1s or 3s), automatically clamped to 5s with terminal notice:
     `Notice: interval clamped to minimum 5s safety threshold`.
3. Cross-platform OS service integration in `cli/cmdlemon/`:
   - Ubuntu / Linux: Generated systemd service unit (`gitmap-lemon.service`).
   - macOS: Generated launchd daemon plist (`com.gitmap.lemon.plist`).
   - Windows: Generated Windows service / scheduled task wrapper (`GitMapLemon`).
4. Wired Cobra subcommands in `cli/cmdlemon/lemon_cmd.go` and top-level routing in `cli/cmd/roottooling.go`:
   - `gitmap lemon run [--interval 8s] [--once] [--verbose]`
   - `gitmap lemon status [--json]`
   - `gitmap lemon poll`
   - `gitmap lemon service <install|uninstall|start|stop|restart|status>`

## Evidence
- `TestIntervalGuard_ClampingNotice` and `TestIntervalGuard_ValidIntervals` pass.
