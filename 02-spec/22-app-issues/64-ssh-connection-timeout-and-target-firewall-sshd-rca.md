# App Issue 64: SSH Connection Timeout, Target OpenSSH Server Absence & Firewall Block RCA

**Issue ID:** 64  
**Date:** 2026-10-03  
**Status:** Resolved  
**Affected Subsystem:** `cli/cmdports`, `cli/cmd`, `cli/cmdssh`

---

## 1. Reproduction

When attempting to join a remote Windows machine (node `t1` at `node-t1`) into the GitMap SSH fleet cluster:

```text
$ gitmap ssh join administrator@node-t1 t1
  INFO Registering node t1 (administrator@node-t1:22)...
  INFO Testing SSH reachability to t1...
  WARN Node t1 responded with connection timeout. Registered as offline.
  SUCCESS Node t1 registered in cluster vault.
```

Subsequent node status inspections and connection attempts failed consistently:

```text
$ gitmap nodes
+----+------+-----------------------+---------+--------+---------+------------------+
| #  | NODE | HOST                  | USER    | PORT   | STATUS  | LATENCY          |
+----+------+-----------------------+---------+--------+---------+------------------+
| 1  | w1   | node-w3          | admin   | 22     | online  | 1.4ms            |
| 2  | w2   | node-w4          | admin   | 22     | online  | 2.1ms            |
| 3  | t1   | node-t1          | admin   | 22     | offline | timeout (3000ms) |
+----+------+-----------------------+---------+--------+---------+------------------+

$ gitmap ssh t1
ssh: connect to host node-t1 port 22: Connection timed out
[exit code 255]
```

When logging into `t1` directly to diagnose the network listener:

```text
PS C:\> gitmap ports
gitmap: 'ports' is not a gitmap command. See 'gitmap --help'.

PS C:\> ssh -V
OpenSSH_for_Windows_8.6p1, LibreSSL 3.4.3

PS C:\> Get-Service sshd
Get-Service : Cannot find any service with service name 'sshd'.
```

### Symptoms:
1. `gitmap ssh join` registered the node, but reachability ping failed with a 3000ms connection timeout.
2. Direct SSH attempts (`gitmap ssh t1` or `ssh administrator@node-t1`) timed out with exit code 255.
3. The target machine possessed the OpenSSH Client binary (`ssh.exe`), giving the false impression that SSH was enabled, but the OpenSSH Server capability (`sshd`) was neither installed nor running.
4. GitMap lacked an integrated cross-platform port and firewall diagnostic utility (`gitmap ports`), forcing developers to fumble with opaque manual PowerShell and `netstat` commands.

---

## 2. Root Cause Analysis (4-Part RCA)

### Symptom:
SSH operations to remote target node `t1` fail with `ssh: connect to host node-t1 port 22: Connection timed out` (exit code 255), and the target host cannot be managed in the GitMap multi-node fleet.

### Direct Cause:
Modern Windows client editions (Windows 10/11) install the OpenSSH Client capability by default, but the OpenSSH Server capability (`OpenSSH.Server~~~~0.0.1.0`) is optional and **uninstalled by default**. Consequently:
1. No process was listening on TCP port 22 (`sshd.exe` was absent).
2. The `sshd` Windows Service was not registered or running.
3. The Windows Advanced Firewall inbound rule for port 22 was missing or disabled, causing incoming TCP SYN packets to be silently dropped rather than rejected, resulting in client connection timeouts.

### Compound Failure:
1. **Asymmetric OpenSSH Defaults:** The presence of `ssh.exe` in `System32` creates false confidence that SSH is functional on Windows, whereas inbound SSH requires the discrete OpenSSH Server feature.
2. **Silent Packet Drops by Firewall:** Unsolicited TCP packets arriving at Windows without an explicit inbound allow rule are dropped silently by Windows Defender Firewall, converting what could be an instantaneous `Connection refused` into a prolonged multi-second `Connection timed out`.
3. **Diagnostic Gap in GitMap Tooling:** When node connectivity fails, developers had no built-in diagnostic tool inside GitMap to inspect open listening ports, process owners, and firewall rule permissions on either the local or remote machine. Running `gitmap ports` resulted in an unknown command error.

### Impact:
Multi-machine fleet expansion, remote command delegation (`gitmap se`), and cross-node Git repository synchronization (`gitmap pull-all`, `gitmap deploy`) were completely blocked whenever newly introduced Windows nodes were unconfigured for inbound SSH.

---

## 3. Fix & Remediation

1. **Cross-Platform Port & Firewall Diagnostic Engine (`gitmap ports` / `cli/cmdports`)**:
   - Implemented `gitmap ports` in package `cmdports` supporting `--port / -p <port>`, `--common`, `--json`, and interactive table rendering via `termtable`.
   - Inspects active TCP listening sockets, maps PIDs to process executable names, and validates inbound firewall policy across Windows (`netsh advfirewall` / `Get-NetFirewallRule`) and Unix (`ufw`, `firewall-cmd`, `ss`).
   - Categorizes common administration and fleet ports (22, 80, 443, 3389, 5985, 5986, 8080) and outputs contextual remediation recommendations.

2. **Automated SSH Server Provisioning (`gitmap ssh enable`)**:
   - Added SSH server enablement orchestration to automate installing the Windows OpenSSH Server capability:
     `Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0`
   - Configures `sshd` service startup type to Automatic and starts the service.
   - Automatically registers and enables inbound Windows Firewall rule:
     `New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH SSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22`

3. **Dedicated Port Inspection Flag (`gitmap ssh port <node>`)**:
   - Allows verifying remote port reachability before attempting SSH key authentication or full cluster joins.

4. **Self-Healing Diagnostics (`gitmap ssh troubleshoot <node>`)**:
   - Executes multi-step diagnostic probes (ICMP ping, TCP port probe, SSH handshake check) and prints immediate actionable commands when reachability fails.

---

## 4. Prevention & Learnings

1. **Provide Built-In Diagnostic Observability:** Never leave developers to debug low-level OS networking with mismatched external commands. Adding `gitmap ports` provides instant, deterministic visibility into listeners and firewall blocks.
2. **Proactive Pre-Flight Verification:** Node registration (`gitmap ssh join`) should probe port 22 and report precise root causes (e.g., "Port 22 unreachable: target firewall dropping packets or SSHD not running") rather than generic timeout messages.
3. **Provide One-Line Remediation:** Always pair diagnostic detections with turn-key fixes (e.g., advising `gitmap ssh enable` on the target machine).
