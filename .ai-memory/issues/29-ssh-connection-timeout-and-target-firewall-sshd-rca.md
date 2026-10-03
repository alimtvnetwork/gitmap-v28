# RCA-29: SSH Connection Timeout, Target OpenSSH Server Absence & Firewall Block

Spec Reference: [02-spec/22-app-issues/64-ssh-connection-timeout-and-target-firewall-sshd-rca.md](../../02-spec/22-app-issues/64-ssh-connection-timeout-and-target-firewall-sshd-rca.md)

---

## 1. Symptom

Attempting to join or connect to a Windows target node (`t1` at `192.168.1.17`) failed with connection timeouts:

```text
$ gitmap ssh join administrator@192.168.1.17 t1
  WARN Node t1 responded with connection timeout. Registered as offline.

$ gitmap nodes
| 3  | t1   | 192.168.1.17          | admin   | 22     | offline | timeout (3000ms) |

$ gitmap ssh t1
ssh: connect to host 192.168.1.17 port 22: Connection timed out
[exit code 255]
```

On node `t1`, attempting to diagnose ports with `gitmap ports` failed with unknown command, despite the `ssh.exe` client binary being present on the system.

---

## 2. Root Cause

1. **Uninstalled OpenSSH Server Capability:** Modern Windows 10/11 includes the OpenSSH Client out-of-the-box, but the OpenSSH Server (`sshd`) is an optional feature that is uninstalled by default.
2. **Missing Inbound Firewall Rule:** Windows Defender Firewall silently drops incoming TCP packets on port 22 without an explicit inbound allow rule, resulting in client connection timeouts instead of quick refusal.
3. **Diagnostic Blindspot in GitMap:** GitMap lacked a native port and firewall inspection command (`gitmap ports`), forcing opaque manual troubleshooting when cluster nodes were unreachable.

---

## 3. Resolution

1. **Cross-Platform Port & Firewall Diagnostic Engine (`gitmap ports` / `cli/cmdports`)**:
   - Implemented `gitmap ports` in package `cmdports` supporting `--port / -p <port>`, `--common`, `--json`, and interactive table rendering via `termtable`.
   - Inspects TCP listening sockets and processes, checks inbound firewall allowance (Windows `netsh advfirewall` and Unix `ufw`/`firewall-cmd`), and provides contextual recommendations.
2. **Automated SSH Server Enablement (`gitmap ssh enable`)**:
   - Automated installation of the `OpenSSH.Server~~~~0.0.1.0` Windows capability, configured service autostart, and opened port 22 inbound in Windows Firewall.
3. **Port Check & Troubleshooting Tools**:
   - Enhanced remote verification workflows with `gitmap ssh port` and `gitmap ssh troubleshoot`.

---

## 4. Prevention & Learnings

- Built-in diagnostic tools (`gitmap ports`) eliminate guesswork when debugging node reachability and cluster membership.
- Node onboarding workflows should surface clear, actionable prerequisite steps (such as enabling `sshd` and firewall rules) when initial connection attempts fail.
