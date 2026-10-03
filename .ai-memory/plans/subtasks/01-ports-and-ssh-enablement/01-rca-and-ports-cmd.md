# Subtask 01: RCA 64 and gitmap ports CLI Implementation

**Parent Plan:** [01-ports-and-ssh-enablement.md](../../pending/01-ports-and-ssh-enablement.md)  
**Assigned Role:** Worker 01  
**Status:** PENDING  

## Scope:
1. Document 4-Part Root Cause Analysis in `02-spec/22-app-issues/64-ssh-connection-timeout-and-target-firewall-sshd-rca.md` and `.ai-memory/issues/29-ssh-connection-timeout-and-target-firewall-sshd-rca.md`.
2. Implement `cli/cmdports/` package:
   - `ports.go`: CLI parsing, flags (`--port`, `--common`, `--json`), TermTable rendering.
   - `ports_windows.go`: PowerShell `Get-NetTCPConnection` & `Get-NetFirewallRule` or `netstat -ano` fallback.
   - `ports_unix.go`: `ss -tulpn` & `ufw` / `firewall-cmd` fallback.
   - `ports_types.go`: structs and models.
3. Register `ports` in `cli/cmd/roottooling.go` and `cli/cmd/rootsuggest.go`.
