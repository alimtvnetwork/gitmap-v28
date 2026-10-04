---
name: gitmap-ssh-cluster-fleet
description: Autonomously develop, test, and modify GitMap's multi-node SSH delegation, cluster management, node joining, credential vault, and termpad/termtable rendering.
---

# GitMap SSH & Fleet Cluster Skill (`gitmap-ssh-cluster-fleet`)

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for developing, auditing, and enhancing GitMap's multi-node SSH orchestration, cluster management, node joining, and remote task execution.

---

## 1. Key Architectural Components & Code Map

| Component | Primary Location | Key Responsibilities |
|---|---|---|
| **SSH Command Dispatch** | `cli/cmdssh/` | Token resolution (positional vs quoted commands), space-delimited node lists, remote execution, and subcommand routing. |
| **Cluster Triad Engine** | `cli/cluster/` | Multi-node cluster orchestration, node heartbeat monitoring, background task daemon, and cluster state sync. |
| **Node Join & Vault** | `cli/cmdssh/sshjoin_cmd.go`, `cli/crypto/` | Interactive machine joining, SSH config parsing, encrypted password/key vault storage, and host recall. |
| **Port 22 Liveness Prober** | `cli/cmdssh/ssh_health.go`, `cli/netip/` | Fast concurrent network probing (TCP port 22) to avoid hanging commands on unreachable or offline nodes. |
| **Port & Firewall Diagnostics Engine** | `cli/cmdports/` | Listening port discovery, process PID attribution, and firewall allow/block inspection (`netsh` on Windows, `ufw`/`firewall-cmd` on Linux). |
| **OpenSSH Daemon Management** | `cli/cmdssh/ssh_daemon_enable.go`, `cli/cmdssh/ssh_port_config.go`, `cli/cmdssh/ssh_troubleshoot.go` | Automated `sshd` installation (`OpenSSH.Server` capability on Windows), service auto-start configuration, inbound firewall rule creation, and remote host troubleshooting. |
| **Terminal Formatting** | `cli/termpad/`, `cli/termtable/` | Clean Unicode table rendering, column auto-padding, status badge colors, and AGY CLI help parity. |

---

## 2. Essential Commands

```bash
# Execute remote command across multiple space-delimited nodes
gitmap ssh exec node1 node2 "uptime && uname -a"

# Positional multi-target resolution with quoted command
gitmap ssh "ls -la /var/log" node1 node2

# Interactive machine join and key registration
gitmap ssh join user@192.168.1.50 --alias worker-node-1

# Check fleet liveness and port 22 connectivity
gitmap ssh health --all

# Target Node: Enable OpenSSH server, service auto-start & firewall allowance
gitmap ssh enable [--port <n>]

# Target Node: Inspect listening ports, process PIDs, and firewall rules
gitmap ports ssh
gitmap ports --common

# Host/Client: Troubleshoot remote node reachability and port 22 timeout
gitmap ssh troubleshoot <alias|ip>

# Inspect or reconfigure remote/local SSH port
gitmap ssh port <alias|ip> [new-port]

# Update target node binary when commands like 'ports' are unrecognized
gitmap self-update

# Cluster status inspection and table view
gitmap cluster status
gitmap sc list
```

---

## 3. Core Invariants & Engineering Guardrails

1. **Graceful Offline Node Handling:** When a target machine is offline or port 22 is closed, the executor must timeout promptly (using pre-flight liveness checks), record the failure in the status table, and proceed with remaining nodes rather than hanging indefinitely.
2. **Command Token Delimitation:** Always support both space-delimited host lists (`node1 node2 "cmd"`) and quoted commands (`"cmd" node1 node2`) via `cli/cmdssh/ssh_exec_resolve.go`.
3. **Hermetic Test Mocking:** Never execute real outbound network requests or SSH connections during unit tests. Unit tests must use mock transports or injectable runners with `defer` restoration blocks.
4. **Zero Destructive OS Calls:** Scheduled power state alterations (`schedule shutdown`, `restart`) must use injectable executors (`DefaultOSActionExecutor`) and never execute real system calls during testing.
5. **No Uppercase ID Acronyms:** Machine and node identifiers must strictly follow `NodeId`, `HostId`, `targetId` naming conventions. All-caps `ID` is prohibited.
6. **OpenSSH Client vs. Server Asymmetry on Windows:** The availability of the `ssh.exe` client binary (which ships by default with Windows 10/11) does NOT imply that OpenSSH Server (`sshd`) is installed or running. When a target node times out on port 22, never assume SSH is operational because `ssh` runs locally. Run `gitmap ssh enable` on the target to ensure the Windows Capability `OpenSSH.Server~~~~0.0.1.0` is installed, the service is started with `Automatic` startup, and inbound port 22 is allowed in Windows Defender Firewall.
7. **Timeout vs. Refusal Diagnostic Heuristic:** `Connection timed out` (exit code 255) almost always indicates packet drop by host or network firewalls, whereas `Connection refused` indicates an unblocked network path with no listening daemon. Always check firewall rules via `gitmap ports firewall` or `gitmap ssh troubleshoot`.
8. **Worker Node Binary Version Parity:** If a worker node outputs `[E1001:VALIDATION] cmd.dispatch: Unknown command: ports`, the remote binary is outdated (pre-v6.467.0). Perform `gitmap self-update` or redeploy the binary before diagnostic commands.
