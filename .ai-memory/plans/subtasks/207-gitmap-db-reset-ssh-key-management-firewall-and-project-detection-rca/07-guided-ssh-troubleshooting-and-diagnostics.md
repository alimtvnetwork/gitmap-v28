# Subtask 77.7: Guided SSH Troubleshooting & Diagnostics

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdssh/`, `cli/apperror/`, `cli/termpad/`

## Objective
Implement intelligent SSH connection failure diagnostics and an automated 7-step remediation HUD displayed in the CLI footer whenever an SSH connection fails, along with a dedicated diagnostic command (`gitmap ssh troubleshoot <target>`, aliases: `doctor`, `diagnose`) featuring multi-probe network analysis and recovery advice.

## Scope & Implementation Details

### 1. Connection Failure Interception & Footer HUD Hook (`cli/cmdssh/ssh_troubleshoot.go`)
- Intercept connection failures across all SSH operations:
  - `gitmap ssh <target>`
  - `gitmap ssh join <target>`
  - `gitmap ssh exec <target>`
  - `gitmap ssh copy-id <target>`
  - `gitmap ssh probe <target>`
- Render structured 7-step remediation card in the CLI footer via `PrintSevenStepTroubleshootGuide(target, err)`:
  ```text
  ╔══════════════════════════════════════════════════════════════════╗
  ║           SSH CONNECTION TROUBLESHOOTING & REMEDIATION           ║
  ╚══════════════════════════════════════════════════════════════════╝
    Failure Cause: dial tcp node-u1:22: connectex: A connection attempt failed...

    7-Step Remediation Checklist for 'node-u1':
      1. Check Host Power & Ping: Confirm the machine is online.
         → ping node-u1
      2. Enable OpenSSH Server: Install GitMap or enable sshd on the target.
         → Run 'gitmap ssh enable --port 22' (or 'sudo systemctl enable --now ssh') on target
      3. Check Target SSH Port: Verify what port sshd is listening on.
         → Run 'gitmap ssh port ls' on target machine (default: 22)
      4. Verify Firewall Rules: Ensure inbound TCP port is allowed through the firewall.
         → Run 'gitmap ssh port add <port>' or 'gitmap ssh enable-public <port>'
      5. Deploy Authentication Key: Push this machine's public key to the target.
         → gitmap ssh copy-id node-u1 (or gitmap ssh view to inspect key)
      6. Run Deep Diagnostics: Execute automated reachability and port probe.
         → gitmap ssh troubleshoot node-u1
      7. Manual Verbose Connection: Test with OpenSSH client directly.
         → ssh -vvv node-u1
  ```

### 2. Deep Diagnostic CLI Command (`cli/cmdssh/ssh_troubleshoot.go`)
- Command: `gitmap ssh troubleshoot <target>` (aliases: `doctor`, `diagnose`).
- Target parsing: handles bare IP (`node-u1`), host:port (`node-u1:2222`), user@host (`a@node-u1`), or enrolled node alias (`u1`).
- Multi-probe diagnostic pipeline:
  - **Ping Probe:** ICMP echo test (1-second timeout) to determine if host hardware and OS network stack are alive.
  - **SSH TCP Probe:** Direct dial to target port (default 22 or parsed port).
  - **Auxiliary Service Probes:**
    - Probe WinRM (port 5985 HTTP, 5986 HTTPS) to detect Windows hosts with remote management active.
    - Probe RDP (port 3389) to detect active desktop access.
- Failure Mode Classification Engine:
  - `FIREWALL INBOUND DROP`: Host is pingable, but SSH TCP port times out. Inbound firewall is dropping packets.
  - `SSH DAEMON NOT RUNNING`: Host is pingable, but SSH TCP connection is actively refused (`RST`). sshd service is inactive or bound to localhost.
  - `HOST UNREACHABLE / OFFLINE`: Ping fails and TCP connection fails. Machine is powered off, suspended, or routing is severed.
  - `SSH PORT OPEN & ACCESSIBLE`: TCP handshake succeeds. Connection issue is authentication-level (public key rejected, bad password, or shell misconfiguration).
- Troubleshoot HUD Output:
  - Reachability matrix: Ping status, SSH TCP status, WinRM/RDP status.
  - Plain-English diagnosis explanation.
  - Target machine commands to resolve the issue (e.g. `gitmap ssh enable`, `gitmap ssh port add`).
  - Alternative remote connection options (e.g. `Enter-PSSession -ComputerName <ip>`, `mstsc /v:<ip>`).
  - Next steps checklist.

### 3. Command Routing & Error Integration (`cli/cmdssh/ssh.go`)
- Register `troubleshoot`, `doctor`, `diagnose` in `dispatchDaemonSSH`.
- Wrap SSH dialers to catch errors and invoke `PrintSevenStepTroubleshootGuide` before bubbling the structured `AppError`.

## Acceptance Criteria
- [ ] Any failed SSH connection dial automatically triggers the 7-step remediation guide in the CLI footer.
- [ ] `gitmap ssh troubleshoot <target>` executes multi-probe checks without crashing or hanging.
- [ ] Aliases `gitmap ssh doctor <target>` and `gitmap ssh diagnose <target>` execute identical diagnostics.
- [ ] Firewall drop vs daemon inactive states are correctly classified based on timeout vs refused TCP responses.
- [ ] Alternative connections (WinRM, RDP) are accurately detected and reported when available.
- [ ] All advice commands printed in the HUD are copy-pasteable and accurate.
