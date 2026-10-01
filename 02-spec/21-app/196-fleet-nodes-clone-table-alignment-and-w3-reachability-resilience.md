# 196 — Fleet Nodes Clone Table Spacing Alignment & W3 Reachability Resilience

## Status: Approved
- **Spec ID:** 196
- **Subsystems:** `cli/cmdnodes`, `cli/cmdssh`
- **Date:** 2026-10-01

---

## 1. Architectural Overview & Context

During fleet clone operations (`gitmap nodes clone <repo>`), GitMap orchestrates concurrent clone actions across the local machine and registered remote SSH worker nodes.

Two primary deficiencies were identified:
1. **Column Boundary & Spacing Distortion in Clone Results Table**:
   - The `DURATION` column formatting width was set to 10 characters (`%-10s`). When local execution rendered `"in-process"`, it filled all 10 characters, leaving zero padding spaces before the separator space. As a result, only a single space separated `"in-process"` from `"already exists on disk"`, distorting the gutter and making column alignment jagged relative to remote node rows.
   - Remote node rows formatted duration as `3687ms` (6 characters), creating a 5-character gap to the `DETAILS` column.
   - The local row displayed `"in-process"` under `DURATION` rather than actual measured elapsed execution duration.
2. **False Node Inaccessibility & Misleading Offline Diagnostics**:
   - When a registered remote node (e.g. `w3`) fails Layer-2 ARP resolution or network ping, GitMap emitted `(machine is off)` during execution.
   - In virtualized or multi-homed environments (e.g. VMware Workstation), machines may be powered on and operational locally while having virtual network adapters isolated, switched to NAT, or reassigned alternate DHCP addresses. Displaying `(machine is off)` is factually misleading to users who can access the VM console.

---

## 2. Technical Specification

### 2.1 Table Column Width Standards
- Table format strings across `renderFleetResultsTable`, `renderLocalRow`, and `renderRemoteRow` MUST use a minimum width of 12 characters for the `DURATION` column (`%-12s`).
- `STATUS` column visual padding MUST use `padVisual(statusTag, 14)` to prevent ANSI escape sequence byte lengths from offsetting visual alignment.
- The `DETAILS` column across all rows (both local and remote) MUST start at the exact same column index (column 81).
- Table divider horizontal lines MUST be extended to match the table width (114 characters).

### 2.2 Local Execution Elapsed Duration Measurement
- `executeLocalClone` MUST measure elapsed execution time using `time.Now()` and `time.Since(start).Milliseconds()`.
- Measured duration MUST be passed to `renderLocalRow` and rendered as `"%dms"` (or `"<1ms"` if elapsed time is zero).
- In the event local execution is skipped (`--skip-local` / `except-self`), the duration MUST display `"-"`.

### 2.3 Diagnostic Message Accuracy
- In `cli/cmdssh/ssh_exec_ui.go`, unreachable nodes MUST be reported as `(unreachable or port 22 closed)` rather than `(machine is off)`.
- In `cli/cmdnodes/nodes_clone_remote.go`, connection failures resulting from timeouts or unresolvable networks MUST be classified accurately with diagnostic context.

---

## 3. Data Contracts & Visual Representation

```text
  NODE (ALIAS)     HOST                   ROLE       STATUS         DURATION     DETAILS
  ------------------------------------------------------------------------------------------------------------------
  local (current)  127.0.0.1              master     ● success      14ms         already exists on disk (D:\work\awansoft-v10)
  w4               192.168.1.13           worker     ● success      3687ms       already exists on disk (D:\work\awansoft-v10)
  w1               192.168.1.3            worker     ● success      3688ms       already exists on disk (D:\work\awansoft-v10)
  w3               192.168.1.12           worker     ○ offline      4516ms       node unreachable (host offline or IP changed)
  main             192.168.1.20           worker     ○ offline      4517ms       node unreachable (host offline or IP changed)
  u1               192.168.1.22           worker     ○ offline      4525ms       node unreachable (host offline or IP changed)
  ------------------------------------------------------------------------------------------------------------------
```

---

## 4. Acceptance Criteria
- **AC-01**: `gitmap nodes clone` renders uniform `DETAILS` column alignment starting at column 81 across local and remote rows.
- **AC-02**: Local clone row displays actual elapsed duration (`%dms` or `<1ms`) under `DURATION`.
- **AC-03**: No row displays `"in-process"` with single-space gutter bleed.
- **AC-04**: Machine offline diagnostics report `(unreachable or port 22 closed)` without claiming the machine is powered off.
- **AC-05**: SemVer minor release bump (`v6.437.0`) executed cleanly per Rule 1.
