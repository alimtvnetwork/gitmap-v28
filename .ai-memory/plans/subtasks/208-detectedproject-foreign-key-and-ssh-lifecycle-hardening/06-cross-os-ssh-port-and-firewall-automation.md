# Subtask 78.6: Cross-OS SSH Port Management & Firewall Automation

## 1. Context & Objective
The user requires commands to manage SSH ports and firewall rules across all supported operating systems (Windows, Ubuntu, CentOS, macOS):
- `gitmap ssh port ls`: List active listening ports and inbound firewall rule status.
- `gitmap ssh port add <port>`: Add a listening port to `sshd_config` and open inbound firewall.
- `gitmap ssh port rm <port>`: Remove a listening port from `sshd_config` and delete firewall rule.
- `gitmap ssh port set <port>`: Switch default listening port.
- `gitmap ssh enable [--port <p>]` / `gitmap ssh disable`.
- `gitmap ssh enable-public <port>`: Open firewall, enable edge traversal, bind `0.0.0.0`, and allow public/external SSH access.

---

## 2. Technical Scope
- `cli/firewall/firewall.go` (unified interface)
- `cli/firewall/windows.go` (Netsh advfirewall)
- `cli/firewall/linux.go` (UFW / firewalld / iptables)
- `cli/firewall/darwin.go` (pfctl)
- `cli/cmdssh/ssh_port_config.go`

---

## 3. Remediation Checklist
- [ ] Ensure cross-OS firewall driver provides `AllowPort`, `BlockPort`, `RemovePortRule`, `IsPortAllowed`, `EnablePublicSSH`.
- [ ] Connect `sshd_config` port parsing and modification with automatic service reload.
- [ ] Wire CLI commands under `gitmap ssh port` and `gitmap ssh enable-public`.
