# 203: Ports Inspection, Firewall Audit & Cross-Platform OpenSSH Daemon Management

**Spec ID:** 203  
**Status:** Approved  
**Version:** 1.0.0  
**Updated:** 2026-10-03  
**Subsystem:** `cli/cmdports`, `cli/cmdssh`, `cli/cmdnodes`  
**Dependencies:** `cli/store`, `cli/apperror`, `cli/constants`

---

## 1. Overview & Problem Statement

In distributed fleet and cluster operations (`gitmap ssh`, `gitmap nodes`, `gitmap cluster exec`), operators frequently encounter connection timeouts:
```text
ssh: connect to host node-t1 port 22: Connection timed out
gitmap ssh: execute failed: [E_INTERNAL_ERROR:EXECUTION] SpawnSSH: exit status 255
```
Investigating the target host often reveals:
1. `ssh` client is present (e.g. built into Windows 10/11), leading operators to believe SSH is active.
2. The OpenSSH Server service (`sshd`) is either **not installed** (optional capability in Windows), stopped, or set to manual startup.
3. Inbound port 22 is blocked by host firewalls (Windows Defender Firewall, UFW, Firewalld, iptables).
4. Running `gitmap ports` on the target machine returns `Unknown command: ports` because no built-in port inspection tool existed.

This specification defines:
1. **`gitmap ports`**: A cross-platform CLI command to audit listening TCP/UDP ports, bound processes, and inbound firewall rules.
2. **`gitmap ssh enable`**: A cross-platform automation command to install, configure, start, and allow OpenSSH Server in the host firewall.
3. **`gitmap ssh port <port>`**: A command to reconfigure the SSH daemon port, update firewall rules, and restart the daemon safely.
4. **`gitmap ssh troubleshoot <target>`**: An interactive diagnostics tool that tests reachability, identifies firewall vs daemon blocks, and suggests alternative remote management protocols (WinRM, PowerShell Remoting).

---

## 2. CLI Command Specifications

### 2.1 `gitmap ports` [flags]

Displays listening ports, process IDs, protocol states, and firewall status on the local system.

#### Syntax & Options:
```text
gitmap ports [--port <n>] [--protocol <tcp|udp>] [--common] [--json]
Aliases: gitmap port, gitmap listening-ports, gitmap open-ports
```

- `--port <n>` (`-p <n>`): Filter for a specific port number (e.g., `gitmap ports -p 22`).
- `--common`: Filter for common administrative ports (22 SSH, 80 HTTP, 443 HTTPS, 3389 RDP, 5985/5986 WinRM, 8080/8443 Web).
- `--json`: Output structured JSON data for programmatic ingestion.

#### Output Format (TermTable):
```text
  PORT    PROTO  PROCESS / SERVICE   LISTENING STATE  FIREWALL STATUS  RECOMMENDATION
  ---------------------------------------------------------------------------------------
  22      TCP    sshd.exe (PID 1420)  ● LISTENING      ● ALLOWED        Ready for SSH
  80      TCP    -                    ○ CLOSED         -                Port inactive
  443     TCP    -                    ○ CLOSED         -                Port inactive
  3389    TCP    TermService          ● LISTENING      ● ALLOWED        RDP enabled
  5985    TCP    WinRM (PID 912)      ● LISTENING      ▲ BLOCKED        Open port in firewall
  ---------------------------------------------------------------------------------------
```

---

### 2.2 `gitmap ssh enable` [flags]

Installs, enables, configures, and starts OpenSSH Server on the local machine with automated firewall configuration.

#### Syntax & Options:
```text
gitmap ssh enable [--port <n>] [--start-now] [--force]
Aliases: gitmap ssh-enable, gitmap ssh enable-server, gitmap ssh sshd
```

- `--port <n>`: Specify listening port during enablement (default: 22).
- `--start-now`: Ensure service is immediately started and verified (default: true).
- `--force`: Reinstall or force firewall rule recreation if existing rule is broken.

#### OS Execution Flow:
1. **Windows (Client 10/11 & Server 2016-2025)**:
   - Queries `Get-WindowsCapability -Online -Name OpenSSH.Server*`.
   - If missing: executes `Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0`.
   - Sets service startup: `Set-Service -Name sshd -StartupType 'Automatic'`.
   - Configures firewall: `New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort <port>`.
   - Starts service: `Start-Service sshd`.
2. **Ubuntu / Debian**:
   - `apt-get update -y && apt-get install -y openssh-server`.
   - `systemctl enable --now ssh`.
   - `ufw allow <port>/tcp`.
3. **CentOS / RHEL / Fedora / Alma / Rocky**:
   - `dnf install -y openssh-server` (or `yum install -y openssh-server`).
   - `systemctl enable --now sshd`.
   - `firewall-cmd --permanent --add-port=<port>/tcp && firewall-cmd --reload`.
   - (If non-default port) `semanage port -a -t ssh_port_t -p tcp <port>`.

---

### 2.3 `gitmap ssh port <port>`

Changes the SSH daemon listening port in `sshd_config`, validates configuration syntax, updates firewall rules, and safely restarts the service.

#### Syntax:
```text
gitmap ssh port <port>
Aliases: gitmap ssh set-port <port>, gitmap ssh config port <port>
```

#### Safety Gates:
1. Validates that `<port>` is within valid range `1..65535`.
2. Checks whether `<port>` is already bound by another process.
3. Creates backup `sshd_config.bak.<timestamp>`.
4. Replaces or appends `Port <port>` directive.
5. Adds new firewall inbound rule for `<port>`.
6. Executes `sshd -t` configuration syntax validation before restarting.
7. If syntax check passes, restarts `sshd`; if it fails, reverts backup and preserves uptime.

---

### 2.4 `gitmap ssh troubleshoot <target>`

Interactive diagnostic tool and connection doctor for remote nodes and SSH targets.

#### Diagnostics Performed:
1. **DNS & IP Resolution**: Resolves target host to IPv4/IPv6.
2. **TCP Port 22 Liveness Dial**: Dials target with bounded timeout (2.0s).
   - If timeout: checks whether ICMP ping succeeds. If ping responds but port 22 times out -> flags **FIREWALL INBOUND DROP**.
   - If connection refused -> flags **SSH DAEMON NOT RUNNING**.
3. **Banner Handshake**: Reads SSH identification string (e.g. `SSH-2.0-OpenSSH_9.2p1`).
4. **Alternative Protocols Check**: Probes alternative remote management ports:
   - WinRM HTTP `5985` / HTTPS `5986` (Windows PowerShell Remoting).
   - RDP `3389` (Remote Desktop).
5. **Actionable Remediation HUD**:
   - Displays exact command to run on target machine (`gitmap ssh enable` or `gitmap ports`).
   - Provides PowerShell and Bash one-liners for remote administration.

---

## 3. OS Command Implementation Matrix

| Action | Windows (PowerShell/CMD) | Ubuntu/Debian (Bash) | CentOS/RHEL (Bash) |
| :--- | :--- | :--- | :--- |
| **Inspect Listening Ports** | `Get-NetTCPConnection -State Listen` | `ss -tulpn` | `ss -tulpn` |
| **Inspect Firewall** | `Get-NetFirewallRule -Direction Inbound` | `ufw status` | `firewall-cmd --list-all` |
| **Install SSH Server** | `Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0` | `apt install -y openssh-server` | `dnf install -y openssh-server` |
| **Enable & Start Service** | `Set-Service sshd -StartupType Automatic; Start-Service sshd` | `systemctl enable --now ssh` | `systemctl enable --now sshd` |
| **Open Port in Firewall** | `New-NetFirewallRule -Name sshd -LocalPort 22 -Protocol TCP -Action Allow` | `ufw allow 22/tcp` | `firewall-cmd --permanent --add-port=22/tcp && firewall-cmd --reload` |
| **Validate Config Syntax** | `& "$env:SystemRoot\System32\OpenSSH\sshd.exe" -t` | `sshd -t` | `sshd -t` |
| **Config File Path** | `C:\ProgramData\ssh\sshd_config` | `/etc/ssh/sshd_config` | `/etc/ssh/sshd_config` |

---

## 4. Verification & Quality Gates

- `cli/cmdports`: Unit tests covering port filtering, common port catalog, and JSON serialization.
- `cli/cmdssh`: Unit tests covering `ssh enable`, `ssh port`, and `ssh troubleshoot` routing and arguments.
- Zero nested `if` statements across all implementations.
- All code formatted and validated with repository linters.
