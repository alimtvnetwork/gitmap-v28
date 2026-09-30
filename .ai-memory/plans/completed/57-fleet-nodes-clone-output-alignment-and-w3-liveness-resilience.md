# 57 — Fleet Nodes Clone Output Alignment & W3 Liveness Resilience

## Traceability
- **Spec Reference:** [02-spec/21-app/194-fleet-nodes-clone-output-alignment-and-w3-liveness-resilience.md](../../../02-spec/21-app/194-fleet-nodes-clone-output-alignment-and-w3-liveness-resilience.md)
- **RCA Reference:** [02-spec/22-app-issues/57-fleet-nodes-clone-misalignment-and-w3-liveness-rca.md](../../../02-spec/22-app-issues/57-fleet-nodes-clone-misalignment-and-w3-liveness-rca.md)
- **Date:** 2026-10-01
- **Status:** In Progress

---

## User Request (Verbatim)

```text
Output is still wrong for 

gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10


I can still acess to w3 machine fix it and release minor please and 
```

---

## Blast Radius & Owned Files

- `cli/cmdnodes/nodes_clone_table.go`: Visual width string padding and pixel-perfect column alignment.
- `cli/cmdnodes/nodes_clone.go`: Dual stdout/stderr capture during local clone execution.
- `cli/cmdclone/clonefixrepo_escape.go`: Suppress escape notice when fleet clone is active.
- `cli/cmdssh/ssh_liveness.go`: 3s timeout, 1-shot retry, 5s failure TTL.
- `cli/cmdssh/ssh_health.go`: 3s default health timeout with retry.
- `cli/cmdssh/ssh_dial_fallback.go`: Pre-flight TCP liveness check to avoid 13.5s timeout cascade.
- `cli/cmdssh/sshexec.go`: Fallback to `OpenGlobalDefault()` in `queryHostPasswordFromDB`.
- `cli/cmdssh/ssh_pull_fleet.go`: Use resilient 3s liveness probe in fleet pull.

---

## Task Breakdown

- **Task-01: Table Alignment & Stderr Leak Prevention**
  - Implement `visibleWidth` and `padRight` in `cli/cmdnodes/nodes_clone_table.go`.
  - Capture both stdout and stderr in `cli/cmdnodes/nodes_clone.go`.
  - Suppress `MsgCFREscapeNested` when fleet clone is active in `cli/cmdclone/clonefixrepo_escape.go`.

- **Task-02: W3 Liveness Resilience, Fast-Fail & Vault Fallback**
  - Increase health timeout to 3000ms and add retry in `cli/cmdssh/ssh_health.go` and `cli/cmdssh/ssh_liveness.go`.
  - Separate cache TTL: 5s for offline/failure, 45s for online/success.
  - Query global vault fallback in `cli/cmdssh/sshexec.go`.
  - Perform fast pre-flight probe in `cli/cmdssh/ssh_dial_fallback.go` to eliminate 13.5s cascade on offline nodes.

- **Task-03: Verification, Minor Bump Release & Audit**
  - Verify `.\gitmap.exe nodes clone https://github.com/alimtvnetwork/awansoft-v10`.
  - Run targeted linters (`check-enum-and-boolean.py`, `check-nested-ifs.py`, `check-relative-paths.py`).
  - Perform minor version bump to `v6.436.0` and commit atomically.
