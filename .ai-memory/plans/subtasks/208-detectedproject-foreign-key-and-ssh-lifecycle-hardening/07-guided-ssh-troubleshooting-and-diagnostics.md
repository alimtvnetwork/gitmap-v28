# Subtask 78.7: Interactive SSH Web UI & 7-Step Troubleshooting Engine

## 1. Context & Objective
- **SSH Web UI (`gitmap ssh ui` / `web`)**: Launch browser dashboard displaying SSH keys, configured ports, firewall status, and cluster fleet nodes.
- **7-Step Troubleshooting Engine**: Whenever an SSH connection to a remote machine fails (e.g. timeout, target not found, auth rejected), render an intelligent, structured 7-step remediation workflow with concrete actions:
  1. Check host power and network reachability (`ping <target>`).
  2. Install GitMap / enable OpenSSH Server on target (`gitmap ssh enable`).
  3. Verify listening SSH port on target (`gitmap ssh port ls`).
  4. Verify firewall allows inbound port (`gitmap ssh port add <port>` / `gitmap ssh enable-public`).
  5. Deploy authentication public key (`gitmap ssh copy-id <target>` / inspect with `gitmap ssh view`).
  6. Run automated deep diagnostics (`gitmap ssh troubleshoot <target>`).
  7. Test manual verbose OpenSSH connection (`ssh -vvv <target>`).

---

## 2. Technical Scope
- `cli/cmdssh/ssh_ui.go`
- `cli/cmdssh/ssh_troubleshoot.go`
- `cli/cmdssh/ssh_target_exec.go`

---

## 3. Remediation Checklist
- [ ] Connect `gitmap ssh ui` to `cmdui.RunUI("ssh", 8080)`.
- [ ] Implement `PrintSevenStepTroubleshootGuide` in `ssh_troubleshoot.go`.
- [ ] Hook troubleshooting guide on connection failure in `ssh_target_exec.go`.
