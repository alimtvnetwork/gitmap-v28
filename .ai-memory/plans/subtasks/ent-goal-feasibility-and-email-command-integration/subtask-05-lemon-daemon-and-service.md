# Subtask 05: GitMap Daemon Polling Loop, Settings Disablement & Cross-Platform Service Hooks

**Parent Plan:** `.ai-memory/plans/ent-goal-feasibility-and-email-command-integration.md`  
**State:** `COMPLETED`  
**Owned Files:** `cli/cmdlemon/interval_guard.go`, `cli/cmdlemon/lemon_daemon.go`, `cli/cmdlemon/daemon_settings.go`, `cli/cmdlemon/service_manager.go`, `cli/cmdlemon/service_linux.go`, `cli/cmdlemon/service_darwin.go`, `cli/cmdlemon/service_windows.go`, `cli/cmdlemon/service_other.go`, `cli/cmdlemon/lemon_cmd.go`  

## Objectives & Delivered Implementation
1. Implemented background daemon runner in `cli/cmdlemon/lemon_daemon.go` with `GitmapDaemon` and `LemonDaemon` naming.
2. Implemented application settings switch in `cli/cmdlemon/daemon_settings.go`:
   - Setting key: `daemon.enabled` (default `true`, supports aliases `gitmap-daemon.enabled`, `daemon_enabled`).
   - When switched off (`daemon.enabled = false`):
     * `gitmap daemon run` checks `IsDaemonEnabled()`, emits clean notice, and exits without starting background loops.
     * In-flight polling daemon checks `IsDaemonEnabled()` on every cycle and skips execution passes.
   - CLI management: `gitmap daemon enable` (`daemon.enabled = true`) and `gitmap daemon disable` (`daemon.enabled = false`).
3. Enforced polling interval constraints via `cli/cmdlemon/interval_guard.go`:
   - Configurable: default 8 seconds (`DefaultPollingInterval`).
   - Non-negotiable minimum 5 seconds guard (`MinPollingInterval`).
   - If requested interval is $< 5$s (e.g. 1s or 3s), automatically clamped to 5s with terminal notice:
     `Notice: interval clamped to minimum 5s safety threshold`.
4. Cross-platform OS service integration in `cli/cmdlemon/`:
   - Ubuntu / Linux: Generated systemd service unit (`gitmap-daemon.service` with `gitmap-lemon.service` fallback).
   - macOS: Generated launchd daemon plist (`com.gitmap.daemon.plist` with `com.gitmap.lemon.plist` fallback).
   - Windows: Generated Windows service / scheduled task wrapper (`GitMapDaemon` with `GitMapLemon` fallback).
5. Wired Cobra subcommands in `cli/cmdlemon/lemon_cmd.go` and top-level routing in `cli/cmd/roottooling.go`:
   - `gitmap daemon run [--interval 8s] [--once] [--verbose] [--force]`
   - `gitmap daemon status [--json]`
   - `gitmap daemon poll`
   - `gitmap daemon enable`
   - `gitmap daemon disable`
   - `gitmap daemon service <install|uninstall|start|stop|restart|status>`
   - Backward-compatible aliases: `gitmap gitmap-daemon`, `gitmap lemon`, `gitmap lm`.

## Evidence
- `TestIntervalGuard_ClampingNotice` and `TestIntervalGuard_ValidIntervals` pass.
- `TestGitmapDaemon_SettingsToggle`, `TestGitmapDaemon_CLISettingsToggle`, `TestGitmapDaemon_ExecutionSuppressionWhenDisabled`, `TestGitmapDaemon_StatusOutputReflectsSettings`, `TestGitmapDaemon_CurrentPlatformServiceConfig`, `TestGitmapDaemon_SettingsKeyAliases` pass.
