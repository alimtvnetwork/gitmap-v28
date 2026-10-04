# Plan: 01-ports-and-ssh-enablement

## User Request (Verbatim)
```
PS C:\Users\Alim> gitmap ssh join administrator@node-t1
✓ Machine 'host-node-t1' (administrator@node-t1) joined successfully.
  Recall anytime: gitmap ssh host-node-t1
  Or connect directly: gitmap ssh administrator@node-t1
  Undo anytime: gitmap ssh undo
PS C:\Users\Alim> gitmap ssh join administrator@node-t1 t1
✓ Machine 't1' (administrator@node-t1) joined successfully.
  Recall anytime: gitmap ssh t1
  Or connect directly: gitmap ssh administrator@node-t1
  Undo anytime: gitmap ssh undo
PS C:\Users\Alim> gitmap nodes

  ALIAS            ROLE           HOST (IP:PORT)         USER           STATUS                ENROLLED
  --------------------------------------------------------------------------------------------------------------
  w1               worker         node-w1:22         Administrator  ○ offline (timeout)   2026-09-30 06:04:29
  w2               worker         node-w2:22         Administrator  ○ offline (timeout)   2026-09-30 06:04:29
  w4               worker         node-w4:22        Administrator  ○ offline (timeout)   2026-09-30 06:04:29
  u1               worker         node-u1:22        a              ○ offline (timeout)   2026-09-30 06:04:29
  t1               worker         node-t1:22        administrator  ○ offline (timeout)   2026-10-03 07:00:23
  ip               worker         exec:22                Alim           ○ offline (timeout)   2026-10-01 04:26:08
  main             worker         node-main:22        administrator  ○ offline (timeout)   2026-09-30 06:06:48
  w3               worker         node-w3:22        Administrator  ● ready               2026-09-30 06:04:29
```

Target Machine:
```
PS C:\Users\Administrator> gitmap ports
gitmap: [E1001:VALIDATION] cmd.dispatch: Unknown command: ports
It is not there, but here is a suggestion you can try: stats, nodes, paet
```

## 1. Overview & Architectural Goals
1. Root Cause Analysis: Diagnose why `gitmap ssh t1` failed with connection timeout (OpenSSH Client present, but OpenSSH Server `sshd` missing/disabled/blocked by firewall).
2. `gitmap ports` command: Inspect open/listening ports and firewall rule status cross-platform (Windows 10/11/Server, Ubuntu/Debian, CentOS/RHEL).
3. `gitmap ssh enable`: Cross-platform command to install OpenSSH Server capability/package, start service, set Automatic startup, and configure inbound firewall rule.
4. `gitmap ssh port <port>`: Change SSH listening port in `sshd_config`, update firewall rules, and restart daemon.
5. Smart Troubleshooting HUD: Intercept SSH exit status 255 / timeout and display clear actionable troubleshooting steps, alternatives (WinRM, PowerShell Remoting), and next steps.
6. Minor version bump and release ceremony.

## 2. Work Breakdown & Subtasks
- **Subtask 01 (`Task-01`)**: Author Canonical Specification `02-spec/21-app/203-ports-inspection-and-ssh-daemon-enablement.md`.
- **Subtask 02 (`Task-02`)**: Author 4-Part RCA in `02-spec/22-app-issues/64-ssh-connection-timeout-and-target-firewall-sshd-rca.md` and `.ai-memory/issues/29-ssh-connection-timeout-and-target-firewall-sshd-rca.md`.
- **Subtask 03 (`Task-03`)**: Implement `gitmap ports` in `cli/cmdports/` and wire into `cli/cmd/roottooling.go`, `cli/cmd/rootsuggest.go`.
- **Subtask 04 (`Task-04`)**: Implement `gitmap ssh enable` in `cli/cmdssh/ssh_daemon_enable.go`.
- **Subtask 05 (`Task-05`)**: Implement `gitmap ssh port <port>` in `cli/cmdssh/ssh_port_config.go`.
- **Subtask 06 (`Task-06`)**: Implement `gitmap ssh troubleshoot` & enhance `SpawnSSHWithPassword` in `cli/cmdssh/ssh_client.go` / `cli/cmdssh/ssh_troubleshoot.go`.
- **Subtask 07 (`Task-07`)**: Unit tests in `cli/cmdports/` and `cli/cmdssh/`, pass linters.
- **Subtask 08 (`Task-08`)**: Minor version bump and atomic commit/push.

## 3. Disjoint Worker Waves
- **Wave 1 (Specs & Core Commands)**:
  - Worker 01: RCA 64, RCA 29, `cli/cmdports/` package.
  - Worker 02: `cli/cmdssh/ssh_daemon_enable.go`, `cli/cmdssh/ssh_port_config.go`.
- **Wave 2 (Troubleshooting HUD & Integration)**:
  - Worker 01: `cli/cmd/roottooling.go`, `cli/cmd/rootsuggest.go`, `cli/cmdports/ports_test.go`.
  - Worker 02: `cli/cmdssh/ssh_troubleshoot.go`, `cli/cmdssh/ssh_client.go`, `cli/cmdssh/ssh_daemon_test.go`.
- **Wave 3 (Lead)**:
  - Lead: Index updates, version bump, targeted linters, atomic commit.
