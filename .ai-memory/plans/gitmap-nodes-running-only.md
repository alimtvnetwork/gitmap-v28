# Plan: GitMap Nodes Running Only (`gitmap-nodes-running-only`)

> Slug: `gitmap-nodes-running-only`  
> Parent Task: Implement `gitmap nodes-running-only` (`nro`), `--json` support, 10s node liveness caching, and running-only delegation commands  
> Date: 2026-10-16  
> Author: Lead Orchestrator  
> Status: IN_PROGRESS  

---

## 1. User Request (Verbatim)

```
PS C:\Users\Administrator> gitmap nodes --json
[
  {
    "alias": "main",
    "role": "worker",
    "host": "192.168.1.20",
    "port": 22,
    "user": "administrator",
    "status": "● ready",
    "subsystems": [
      "SSH",
      "Cluster",
      "SC"
    ],
    "enrolled_at": "2026-09-29 14:38:44",
    "is_online": true
  },
  {
    "alias": "w3",
    "role": "worker",
    "host": "192.168.1.12",
    "port": 22,
    "user": "Administrator",
    "status": "○ offline (timeout)",
    "subsystems": [
      "SSH",
      "Cluster",
      "SC"
    ],
    "enrolled_at": "2026-09-29 14:38:44",
    "is_online": false
  },
  {
    "alias": "w1",
    "role": "worker",
    "host": "192.168.1.3",
    "port": 22,
    "user": "Administrator",
    "status": "● ready",
    "subsystems": [
      "SSH",
      "Cluster",
      "SC"
    ],
    "enrolled_at": "2026-09-29 14:38:44",
    "is_online": true
  },
  {
    "alias": "w2",
    "role": "worker",
    "host": "192.168.1.7",
    "port": 22,
    "user": "Administrator",
    "status": "○ offline (timeout)",
    "subsystems": [
      "SSH",
      "Cluster",
      "SC"
    ],
    "enrolled_at": "2026-09-29 14:38:44",
    "is_online": false
  },
  {
    "alias": "w4",
    "role": "worker",
    "host": "192.168.1.13",
    "port": 22,
    "user": "Administrator",
    "status": "○ offline (timeout)",
    "subsystems": [
      "SSH",
      "Cluster",
      "SC"
    ],
    "enrolled_at": "2026-09-29 14:38:44",
    "is_online": false
  },
  {
    "alias": "u1",
    "role": "worker",
    "host": "192.168.1.22",
    "port": 22,
    "user": "a",
    "status": "○ offline (timeout)",
    "subsystems": [
      "SSH",
      "Cluster",
      "SC"
    ],
    "enrolled_at": "2026-09-29 14:38:44",
    "is_online": false
  }
]
PS C:\Users\Administrator>

also we should have new command

gitmap nodes-running-only (nro) --json 

or delegeation commands needs to be aded with this, which will only run on opened machines with cacheing and optimiation, clear???
```

---

## 2. Executive Summary & Architecture

The user observed that running `gitmap nodes --json` lists all registered cluster nodes, including offline nodes (`w3`, `w2`, `w4`, `u1`) which timed out (`"○ offline (timeout)"`).

This task introduces:
1. **Dedicated Command**: `gitmap nodes-running-only` and short alias `gitmap nro` (with additional aliases `nodesrunningonly`, `nodes-running`, `nodes-online`, `nodes-open`, `nodes-active`).
2. **Strict Running-Only Output**: Supports `--json` / `-j` returning a JSON array strictly containing nodes with `is_online == true` (e.g. `main` and `w1`), and an ANSI table omitting offline noise.
3. **10-Second Cross-Process Liveness Cache**: Avoids recurring 1200ms TCP timeouts per offline machine during frequent CLI calls by persisting probe results to `os.TempDir()/gitmap_nodes_liveness_cache.json` with an in-memory `sync.RWMutex`.
4. **Running-Only Delegation Commands**: `nro exec "<cmd>"`, `nro run "<cmd>"`, `nro tasks`, `nro resources`, `nro ping`, `nro pull`, `nro cp all` execute strictly across running/open nodes, skipping offline ones without delay or warning pollution.
5. **Flag Parity Across Fleet Commands**: Supports `--running-only`, `--open-only`, `-ro`, `--online-only` across existing node commands.

---

## 3. Subtask Decomposition (6 Subtasks)

| Subtask ID | Title | Owned Files | Status |
|---|---|---|---|
| **Subtask 01** | Fast Fleet Node Liveness Cache & Probe Optimization | `cli/cmdnodes/nodes_liveness_cache.go`, `cli/cmdnodes/nodes_probe.go`, `cli/cmdnodes/nodes_collect.go` | PENDING |
| **Subtask 02** | `gitmap nodes-running-only` & `gitmap nro` Dispatcher & CLI | `cli/cmdnodes/nodes_running_only.go`, `cli/cmdnodes/nodes_cmd.go`, `cli/cmd/rootcore.go` | PENDING |
| **Subtask 03** | Running-Only Delegation & Fast SSH Execution | `cli/cmdssh/sshexec.go`, `cli/cmdnodes/nodes_filter.go`, `cli/cmdnodes/nodes_scan.go` | PENDING |
| **Subtask 04** | Node Tasks & Resources Open-Only Filtering | `cli/cmdnodes/nodes_resources.go`, `cli/cmdnodes/nodes_tasks.go`, `cli/cmdnodes/nodes_tasks_types.go` | PENDING |
| **Subtask 05** | Documentation, Help, Completion & UI Registry | `cli/helpdoc/nodes-running-only.md`, `cli/helpdoc/print.go`, `cli/helpdoc/catalog.go`, `src/data/commands.ts`, `cli/completion/completion.go`, `cli/constants/constants_commands.go` | PENDING |
| **Subtask 06** | Verification, Release Bump (`v6.555.0`) & CI/CD Telemetry | `version.json`, `package.json`, `cli/constants/constants.go`, `changelog.md` | PENDING |

---

## 4. Invariants & Guardrails

- **Zero Nested Ifs**: Nesting depth $\le 1$ strictly enforced across all files.
- **Positive Booleans Only**: `is*`, `has*`, `can*`, `should*` prefixes; no explicit `== true`, no `!isSuccess` (use `.IsFail()`).
- **File Size Cap**: Strict $\le 300$ line cap per file.
- **Relative Paths**: No absolute filesystem paths or `file:///` URIs.
- **Error Wrapping**: All Go errors wrapped with `apperror.WrapSimple` or `*appfault.AppError`.
- **Atomic Commits**: `gitmap cpf "<module> - <summary>"`.
