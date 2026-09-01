# Specification 16 - Part 7: Standalone Node Join CLI, Auto-Startup Agent & Cluster Connectivity Diagnostics

## 1. Architecture of Standalone Node Agent

### 1.1 `gitmap-node-join` Dedicated CLI

A lightweight, minimal-dependency binary compiled separately from the main `gitmap` CLI:
- **Executable Name**: `gitmap-node-join.exe` (Windows) / `gitmap-node-join` (Linux/macOS).
- **Core Purpose**: Connect client machines to the central cluster orchestrator daemon without requiring full gitmap database or git dependencies on worker nodes.
- **CLI Syntax**:
  ```bash
  # Interactive or parameterized join
  gitmap-node-join --server 192.168.1.10:9999 --token <join-token> [--alias "workstation-02"] [--auto-start]
  ```

### 1.2 Windows Auto-Startup Registration

When `--auto-start` is enabled (or via `gitmap-node-join install-startup`):
- Registers a Windows Run entry under `HKCU\Software\Microsoft\Windows\CurrentVersion\Run\GitMapNodeJoin` or creates a shortcut in `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup`.
- On system boot, `gitmap-node-join` launches in the background, reads cached credentials (`%USERPROFILE%\.gitmap\node-config.json`), establishes connection with the orchestrator, and maintains a persistent WebSocket/TCP heartbeat.

## 2. Orchestrator Daemon (`gitmap serve`) Integration

### 2.1 Node Lifecycle States

The orchestrator tracks node states in real-time:
- `online`: Heartbeat received within last 30 seconds.
- `idle`: Online and ready to receive commands.
- `busy`: Currently executing a delegated sub-command.
- `unreachable`: Heartbeat missed for > 60 seconds.
- `unauthorized`: Token expired or rejected.

## 3. Explicit Unreachable Machine Diagnostics

### 3.1 Warning Surfacing in Cluster Commands

Whenever a cluster command (`gitmap cluster nodes`, `gitmap servers-clients`, `gitmap clients`) executes, if any registered machines fail to respond or disconnect:
- The orchestrator surfaces an explicit diagnostic banner BEFORE and AFTER execution:

```text
  ────────────────────────────────────────────────────────────
  ▲ WARNING: 2 cluster machines are offline or unreachable:
     • Node 3 [ws-build-03] (192.168.1.13) - Last seen 42m ago (timeout)
     • Node 5 [render-box]  (192.168.1.15) - Connection refused

  These are the machines that cannot connect or be found.
  Excluded automatically from current broadcast.
  ────────────────────────────────────────────────────────────
```
