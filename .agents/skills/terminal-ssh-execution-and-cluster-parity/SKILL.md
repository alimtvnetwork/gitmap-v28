---
name: terminal-ssh-execution-and-cluster-parity
description: >-
  Autonomously implement and verify terminal SSH execution, machine liveness scanning and caching,
  smart default key discovery, remote gitmap install/update, AGY/VS Code delegation, AppError stack traces,
  and SSH vs Cluster vs SC comparative help frameworks across GitMap.
---

# Terminal SSH Execution, Liveness Scan, Remote Install & Cluster Parity

Autonomously implement and verify terminal SSH execution, machine liveness scanning and caching,
smart default key discovery, remote gitmap install/update, AGY/VS Code delegation, AppError stack traces,
and SSH vs Cluster vs SC comparative help frameworks across GitMap.

## Core Directives

1. **Liveness Probe & Availability Cache**:
   - Fast ping/TCP dial with short TTL cache so offline machines are skipped gracefully with clear status.
   - Machines that are online run the command immediately.
2. **Smart Default Key Discovery**:
   - When no password or key is explicitly saved in DB, probe `~/.ssh/id_ed25519`, `~/.ssh/id_rsa`, `~/.ssh/id_ecdsa`, and SSH Agent.
   - Provide structured `AppError` envelopes with stack traces and actionable hints if auth fails.
3. **Remote GitMap Install & Update**:
   - `gitmap ssh install [gitmap]` installs if missing, updates if present.
   - `gitmap ssh update [gitmap]` updates fleet machines.
4. **Remote AGY & VS Code Commands**:
   - `gitmap ssh agy <args>` and `gitmap ssh code <args>` to open folders and run tools remotely.
5. **Comparative Architecture & UI Help**:
   - Comparing table for `ssh` vs `cluster` vs `sc` explaining when to use which, how to join, and how to monitor.
