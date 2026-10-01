# 25 — Fleet Nodes Clone Table Spacing Alignment and W3 Reachability Diagnostics RCA

## Status: Resolved
- **Issue ID:** RCA-59 / Issue 25
- **Spec Reference:** [02-spec/21-app/196-fleet-nodes-clone-table-alignment-and-w3-reachability-resilience.md](../../02-spec/21-app/196-fleet-nodes-clone-table-alignment-and-w3-reachability-resilience.md)
- **Target Subsystems:** `cli/cmdnodes`, `cli/cmdssh`
- **Date:** 2026-10-01

---

## 1. Symptoms & Resolution
- **Symptom 1**: In `gitmap nodes clone`, local host displayed `"in-process already exists on disk"` with collapsed 1-space gutter, while remote rows displayed 5-space gutters.
  - **Resolution**: Measured local duration in `executeLocalClone`, rendered `"%dms"` or `"<1ms"`, and widened DURATION column to `%-12s` with `DETAILS` aligning on column 81.
- **Symptom 2**: When `w3` machine failed network connection, GitMap reported `(machine is off)` even though user had direct access to the active VM console.
  - **Resolution**: Replaced `(machine is off)` with `(unreachable or port 22 closed)` and clarified ARP/network unreachable errors in `nodes_clone_remote.go`.
