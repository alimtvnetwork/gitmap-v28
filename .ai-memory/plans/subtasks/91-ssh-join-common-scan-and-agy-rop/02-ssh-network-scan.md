# Subtask 91-02: SSH Network Subnet Scanner (`gitmap ssh scan`)

Spec Reference: [02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md](../../../../02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md)
Parent Plan: [.ai-memory/plans/completed/91-ssh-join-common-scan-and-agy-rop.md](../../completed/91-ssh-join-common-scan-and-agy-rop.md)

## Objective
Update `gitmap ssh scan` so that running `gitmap ssh scan [cidr]` discovers how many active SSH machines are available in the local network / subnet, showing reachable IPs, latency, enrollment status, and connection suggestions.

## Functional Requirements
1. **Network Discovery Overhaul:**
   - Modify `runSSHScanCLI` in `cli/cmdssh/ssh_scan_cmd.go` to invoke `ExecuteSubnetScan` from `cli/cmdssh/ssh_scanner.go` rather than only checking registered hosts.
   - Automatically detects the machine's active IPv4 subnet (e.g. `192.168.1.0/24`) or accepts a specified CIDR (e.g. `192.168.1.0/24`, `10.0.0.0/24`).
   - Concurrently probes port 22 with bounded timeout (e.g. 500-800ms) across 254 addresses using worker pool.
   - Cross-checks discovered IPs against registered hosts in `ssh_hosts` table to show `[ENROLLED: <alias>]` vs `[NEW]`.
   - Displays clear aligned table with total found count.

## Files to Create/Modify
- `cli/cmdssh/ssh_scan_cmd.go` [MODIFY]
- `cli/cmdssh/ssh_scanner.go` [MODIFY]
- `cli/cmdssh/ssh_scanner_test.go` [MODIFY]
