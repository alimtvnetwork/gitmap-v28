# Subtask 77.5: Cross-OS SSH Port Management & Firewall Automation

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdssh/`, `cli/firewall/`, `cli/cmd/`

## Objective
Implement cross-platform SSH port lifecycle commands (`ls`, `add`, `rm`, `set`), daemon enable/disable controls (`gitmap ssh enable [--port <p>]`, `gitmap ssh disable`), public network exposure (`gitmap ssh enable-public <port>`), and unified host firewall rule synchronization across Windows (Netsh/PowerShell), Ubuntu/Debian (`ufw`), CentOS/RHEL (`firewalld`), and macOS (`pfctl`).

## Scope & Implementation Details

### 1. SSH Port Management CLI Suite (`cli/cmdssh/ssh_port_config.go`)
- `gitmap ssh port ls` (aliases: `list`, `show`):
  - Parse configured listening ports from `sshd_config`.
  - Check live firewall status for each port via `firewall.IsPortAllowed(port)`.
  - Render styled table displaying `PORT`, `PROTOCOL` (TCP), and `FIREWALL STATUS` (`✔ ALLOWED` in green, `✘ BLOCKED` in red).
  - Print quick-action command tips for `add`, `rm`, `set`, and `enable-public`.
- `gitmap ssh port add <port>`:
  - Validate port number range (1–65535).
  - Append new `Port <port>` directive to `sshd_config` if not already present.
  - Create automatic timestamped backup of `sshd_config` before modification.
  - Invoke `firewall.AllowPort(port, name)` to open inbound TCP rule.
  - Verify syntax via `sshd -t -f <configPath>`; revert to backup on error.
  - Restart SSH daemon service (`systemctl restart ssh/sshd` or `Restart-Service sshd`).
- `gitmap ssh port rm <port>` (aliases: `remove`, `delete`, `del`):
  - Remove target `Port <port>` directive from `sshd_config`.
  - Invoke `firewall.RemovePortRule(port, name)` to delete inbound firewall rule.
  - Validate syntax and restart SSH daemon service.
- `gitmap ssh port set <port>`:
  - Replace primary listening port in `sshd_config`.
  - Synchronize firewall rules and restart daemon service.

### 2. SSH Daemon Lifecycle & Exposure Controls (`cli/cmdssh/ssh_port_config.go`, `cli/cmdssh/ssh.go`)
- `gitmap ssh enable [--port <port>]` (aliases: `enable-server`, `sshd`, `enable-sshd`):
  - Ensure OpenSSH Server is installed:
    - Windows: `Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0`.
    - Linux (Debian/Ubuntu): verify `openssh-server` package; install if absent.
    - Linux (RHEL/CentOS): verify `openssh-server` package; install if absent.
  - Configure startup type to automatic:
    - Windows: `Set-Service -Name sshd -StartupType Automatic; Start-Service sshd`.
    - Linux: `systemctl enable --now ssh` (or `sshd`).
  - Set target port if `--port` flag is provided.
  - Ensure default port 22 or requested port is allowed in firewall.
- `gitmap ssh disable`:
  - Stop OpenSSH server:
    - Windows: `Stop-Service -Name sshd -Force; Set-Service -Name sshd -StartupType Disabled`.
    - Linux: `systemctl disable --now ssh` (or `sshd`).
  - Remove opened firewall rules for configured SSH ports.
  - Output clear confirmation banner.
- `gitmap ssh enable-public <port>` (alias: `public`):
  - Configure `sshd_config` to bind all interfaces: `GatewayPorts yes` and `ListenAddress 0.0.0.0`.
  - Update host firewall rule with public profile permissions:
    - Windows: `New-NetFirewallRule` / `Set-NetFirewallRule` with `-Profile Any` and `-EdgeTraversalPolicy Allow`.
    - Linux (ufw): `ufw allow <port>/tcp`.
    - Linux (firewalld): `firewall-cmd --permanent --zone=public --add-port=<port>/tcp && firewall-cmd --reload`.
    - macOS: update `pfctl` ruleset.
  - Restart SSH daemon service.

### 3. Cross-OS Firewall Synchronization Engine (`cli/firewall/`)
- Unified interface in `cli/firewall/firewall.go`:
  - `AllowPort(port int, name string) error`
  - `BlockPort(port int, name string) error`
  - `RemovePortRule(port int, name string) error`
  - `IsPortAllowed(port int) (bool, error)`
  - `EnablePublicSSH(port int) error`
  - `ListRules() ([]Rule, error)`
- OS Adapters:
  - Windows (`cli/firewall/windows.go`): PowerShell cmdlets `New-NetFirewallRule`, `Remove-NetFirewallRule`, `Get-NetFirewallPortFilter`, `Set-NetFirewallRule`.
  - Linux (`cli/firewall/linux.go`): Detection and execution of `ufw` vs `firewall-cmd` vs `iptables`.
  - macOS (`cli/firewall/darwin.go`): Anchor definitions in `/etc/pf.anchors/gitmap` and `pfctl -f /etc/pf.conf`.

### 4. Command Dispatch & Routing (`cli/cmdssh/ssh.go`)
- Route `port`, `ports`, `set-port` to `RunSSHPortCLI(args)`.
- Route `enable`, `enable-server`, `sshd`, `enable-sshd` to `RunSSHEnableCLI(args)`.
- Route `disable`, `disable-server`, `stop-sshd` to `RunSSHDisableCLI(args)`.
- Route `enable-public`, `public` to `RunSSHPublicCLI(args)`.

## Acceptance Criteria
- [ ] `gitmap ssh port ls` lists all active ports and correct firewall status.
- [ ] `gitmap ssh port add <port>` successfully updates `sshd_config`, allows port in firewall, and restarts sshd.
- [ ] `gitmap ssh port rm <port>` removes directive, cleans firewall rule, and restarts sshd.
- [ ] `gitmap ssh enable [--port <p>]` verifies installation, enables service startup, and ensures port open.
- [ ] `gitmap ssh disable` stops sshd, sets service startup to disabled, and cleans firewall rules.
- [ ] `gitmap ssh enable-public <port>` sets `ListenAddress 0.0.0.0` and public profile firewall access.
- [ ] Syntax check rollback (`sshd -t`) prevents broken `sshd_config` from breaking active daemon.
