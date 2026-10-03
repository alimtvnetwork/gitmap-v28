# Subtask 03: gitmap ssh troubleshoot and Connection Failure Interception

**Parent Plan:** [01-ports-and-ssh-enablement.md](../../pending/01-ports-and-ssh-enablement.md)  
**Assigned Role:** Worker 02  
**Status:** PENDING  

## Scope:
1. Implement `cli/cmdssh/ssh_troubleshoot.go`:
   - `RunSSHTroubleshootCLI(ctx context.Context, args []string) error`
   - Host reachability probe: IP resolution, TCP dial on target port, ICMP ping.
   - Diagnoses: Timeout vs Connection Refused vs Firewall Drop.
   - Tests alternative ports: WinRM (5985/5986), RDP (3389).
   - Renders Troubleshooting HUD with clear next steps and commands to run on target.
2. Enhance `cli/cmdssh/ssh_client.go`:
   - Intercept OpenSSH exit status 255 and connection timeout errors in `SpawnSSHWithPassword` and `executeClientCmd`.
   - Print smart diagnostic guidance and link to `gitmap ssh troubleshoot <target>`.
3. Unit tests in `cli/cmdssh/ssh_daemon_test.go` and `cli/cmdports/ports_test.go`.
