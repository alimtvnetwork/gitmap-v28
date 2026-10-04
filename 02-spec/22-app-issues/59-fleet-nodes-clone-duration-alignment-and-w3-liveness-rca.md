# 59 — Fleet Nodes Clone Duration Alignment, Column Gutter Spacing, and W3 Reachability Diagnostics RCA

## Status: Resolved
- **Issue ID:** RCA-59 / Issue 24
- **Spec Reference:** [02-spec/21-app/196-fleet-nodes-clone-table-alignment-and-w3-reachability-resilience.md](../21-app/196-fleet-nodes-clone-table-alignment-and-w3-reachability-resilience.md)
- **Target Subsystems:** `cli/cmdnodes`, `cli/cmdssh`
- **Date:** 2026-10-01

---

## 1. Reproduction & Symptoms

### 1.1 In-Process Duration String and Gutter Collision
When executing `gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10`, the local host row rendered:
```text
  local (current)  127.0.0.1              master     ● success      in-process already exists on disk (./awansoft-v10)
```
Notice that `"in-process"` and `"already exists on disk"` had only a single space between them. Because `"in-process"` contains 10 characters and `%-10s` was used, zero padding spaces were appended. Adjacent remote rows displayed `3687ms` with 4 padding spaces plus 1 delimiter space, resulting in a visible zigzag in the `DETAILS` column.

### 1.2 "Machine is off" Diagnostic Misleading Output
When remote node `w3` (`node-w3`) failed to connect, GitMap printed:
```text
Notice: The following machine(s) are currently OFF or unreachable: • [w3 | node-w3] (machine is off)
```
The user confirmed the machine was powered on, logged in, and accessible in their VMware console. Labeling an unreachable network endpoint as `(machine is off)` caused confusion because network isolation (e.g. Host-Only/NAT configuration, DHCP IP change, or firewall rules) was misdiagnosed as physical/virtual power-off.

---

## 2. Root Cause Analysis (4-Part RCA)

### 2.1 Direct Cause
1. `renderLocalRow` hardcoded `"in-process"` into the `DURATION` column rather than measuring and displaying elapsed time.
2. The format string `%-10s` allocated insufficient width for 10-character tokens, producing zero padding and collapsing the column gutter to 1 space.
3. `ssh_exec_ui.go` unconditionally appended `(machine is off)` to uncontactable hosts.

### 2.2 Indirect Cause
- `executeLocalClone` did not track start time or return execution duration, forcing formatting code to invent a placeholder string.
- Liveness detection had no distinction between an unreachable network interface and an actual confirmed shutdown.

### 2.3 Environmental Context
- Multi-VM VMware Workstation environment where guest VMs (`w1`, `w2`, `w3`, `w4`) obtain dynamic DHCP leases from the gateway (`gateway-node`). Any interface state change or network adapter bridge disconnection renders the IP unreachable via ARP without shutting down the VM.

### 2.4 Systemic Vulnerability
- Visual table column widths lacked explicit safety margin over longest possible value, allowing string length edge cases to compromise terminal alignment.

---

## 3. Corrective & Preventive Actions
1. **Measured Duration**: Measure `localDurMs` in `executeLocalClone` and render `"%dms"` or `"<1ms"`.
2. **Column Gutter Expansion**: Expand `DURATION` column width from 10 to 12 across headers and rows, ensuring column 81 alignment for all details text.
3. **Accurate Offline Diagnostics**: Replace `(machine is off)` with `(unreachable or port 22 closed)` in `ssh_exec_ui.go`.
4. **Resilient Error Categorization**: Provide granular categorization in `nodes_clone_remote.go` for network unreachable vs auth failures.
