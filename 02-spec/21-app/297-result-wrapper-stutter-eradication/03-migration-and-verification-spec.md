# Migration & Verification Specification 297: FastGate Validation & Disjoint Swarm Schedule

## 1. Disjoint Wave Schedule

Execution is structured into three disjoint waves to prevent concurrent file modification collisions:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Core Subsystems Stutter Eradication (Subtasks 01–03)                 │
│ Files: cmdagy, prdb, osuser, osclean                                         │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Command & Automation Stutter Elimination (Subtasks 04–06)            │
│ Files: cmdautomation, committransfer, utility types                          │
├──────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: Static AST Verification, Documentation & Release (Subtasks 07–09)    │
│ Files: FastGate gate, plan registries, atomic gitmap cpf commit              │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Disjoint File Allocation Matrix

| Subtask ID | Wave | Exclusive File Box Boundaries |
| :--- | :--- | :--- |
| `subtask-01` | Wave 1 | `cli/cmdagy/types.go`<br/>`cli/cmdagy/agy_agentapi_ops.go`<br/>`cli/cmdagy/agy_open_paths.go`<br/>`cli/cmdagy/agy_cmd.go`<br/>`cli/cmdagy/agy_ping.go` |
| `subtask-02` | Wave 1 | `cli/prdb/pr_split_conn.go`<br/>`cli/prdb/pr_split_ops.go` |
| `subtask-03` | Wave 1 | `cli/osuser/group_portable.go`<br/>`cli/osuser/portable.go`<br/>`cli/osclean/dev_clean.go`<br/>`cli/osclean/terminal_clean.go` |
| `subtask-04` | Wave 2 | `cli/cmdautomation/phase1_types.go`<br/>`cli/cmdautomation/phase2_types.go`<br/>`cli/cmdautomation/phase3_types.go`<br/>`cli/cmdautomation/phase4_types.go`<br/>`cli/cmdautomation/phase5_types.go`<br/>`cli/cmdautomation/phase6_types.go`<br/>`cli/cmdautomation/types.go`<br/>`cli/cmdautomation/runtime_db.go`<br/>`cli/cmdautomation/worker_pool.go` |
| `subtask-05` | Wave 2 | `cli/committransfer/pr_engine.go` |
| `subtask-06` | Wave 2 | `cli/cmddownload/types.go`<br/>`cli/cmdupdate/types.go`<br/>`cli/cmdzsh/types.go`<br/>`cli/lazyregex/types.go`<br/>`cli/macro/types.go`<br/>`cli/netip/types.go`<br/>`cli/osfix/types.go`<br/>`cli/movemerge/types.go`<br/>`cli/cmdschedule/types.go` |
| `subtask-07` | Wave 3 | `03-ai-scripts/50-fastgate.py` (runner execution only) |
| `subtask-08` | Wave 3 | `.ai-memory/plans/readme.md`<br/>`.ai-memory/plans/pending/297-result-wrapper-stutter-eradication.md` |
| `subtask-09` | Wave 3 | Final Lead Orchestrator Commit & Push |

---

## 3. FastGate Invariants
- Functions $\le 15$ lines of executable body.
- Files $\le 300$ lines.
- Zero nested `if` statements.
- Positive boolean variables (`is*`, `has*` only, zero negative tokens).
- Zero compiler builds or unit tests.
