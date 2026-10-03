# Subtask 02: gitmap ssh enable and gitmap ssh port Implementation

**Parent Plan:** [01-ports-and-ssh-enablement.md](../../pending/01-ports-and-ssh-enablement.md)  
**Assigned Role:** Worker 02  
**Status:** PENDING  

## Scope:
1. Implement `cli/cmdssh/ssh_daemon_enable.go`:
   - `RunSSHEnableCLI(args []string) error`
   - Cross-platform OpenSSH Server installation:
     - Windows: `Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0`, firewall rule creation (`New-NetFirewallRule`), startup automatic, start service.
     - Linux (Debian/Ubuntu): `apt-get install -y openssh-server`, `systemctl enable --now ssh`, `ufw allow <port>`.
     - Linux (RHEL/CentOS): `dnf install -y openssh-server`, `systemctl enable --now sshd`, `firewall-cmd`.
2. Implement `cli/cmdssh/ssh_port_config.go`:
   - `RunSSHPortCLI(args []string) error`
   - Validates port `1..65535`.
   - Modifies `sshd_config` (`Port <port>`), tests syntax (`sshd -t`), updates firewall rules, restarts daemon.
3. Wire into `cli/cmdssh/ssh.go` and `cli/cmdssh/exports.go`.
