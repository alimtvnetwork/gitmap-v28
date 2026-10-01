# [V6] Plan 61: PAS Worker Concurrency, Pull-Error Split-DB Subsystem & Machine Telemetry [COMPLETED]

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> **Plan Slug:** `61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry`  
> **Status:** `COMPLETED`  
> **Tracking Spec:** [198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Runtime:** Google Antigravity 2.0 (IDE and CLI)  

---

## Executive Summary & Outcomes

All 5 core requirements from the user's observation and specification have been implemented, tested for regression, and verified:
1. **Dedicated `gitmap machine` Telemetry:** Purged unrelated command help from the default `gitmap machine` invocation; emits clean machine identity cards and pure `--json` output without dumping help menus.
2. **Startup Fleet Version Display:** Queries and displays the exact GitMap binary version on each reachable remote node and the local machine in `gitmap pa --ssh` and `gitmap pas` startup banners.
3. **Concurrency Controls & `paswh` Command:** Implemented `--w <N>`, `--worker <N>`, `--hand <Y>`, `--h <Y>`, `--wwoh`, and root command `gitmap paswh [N] [Y]` (defaulting to $N=1, Y=1$).
4. **Heartbeat Progress (>30s):** Integrated a background heartbeat ticker `StartPullHeartbeat` during long-running pull runs to emit `[⏳ 35s elapsed] Pulling repositories...` so terminal never appears hung.
5. **Categorized Output Grouping & Dirty File Lists:** Grouped pull results into Updated (+N/-N line stats, commit range), Dirty (itemized modified and untracked file listings with remediation suggestions), and Failed (classified error with recovery command).
6. **Pull-Error Subsystem (`gitmap pull-error`):** Built the `gitmap pull-error` (`pulle`, `pull-e`, `pull-errors`) subsystem backed by Split-DB SQLite storage `pull_errors`, context-aware routing (target repo if in repo dir, `all` otherwise), `--ssh` fleet aggregation, pure `--json` AI ingestion format, and diagnostic card rendering.
7. **Wincredman Diagnostic Remediation:** Classified `fatal: Unable to persist credentials with the 'wincredman' credential store` into Windows Credential Manager failure with actionable `gitmap fix-credential` (alias `fc`) remediation.

---

## Subtask Realization Summary

| Task-ID | Subtask | Owner | Owned Files | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `Task-01` | Dedicated `machine` Telemetry & Startup Node Version | Worker 01 | `cli/cmdos/os_machine_alias.go`, `cli/cmdssh/ssh_pull_fleet.go` | COMPLETED | Clean machine identity card, `--json`, fleet version column |
| `Task-02` | Concurrency Controls (`--w`, `--hand`, `--wwoh`, `paswh`) | Worker 02 | `cli/cloneconcurrency/resolve.go`, `cli/cmdpull/pull.go`, `cli/cmdssh/ssh_pas_fleet.go`, `cli/cmd/rootcore.go` | COMPLETED | `ResolveWorkerHands`, concurrency flags, `paswh` alias routing |
| `Task-03` | Heartbeat Ticker (>30s) & Structured Output Grouping | Worker 01 | `cli/cmdpull/pull_heartbeat.go`, `cli/cmdpull/pull_efficient.go`, `cli/cmdpull/pull_efficient_render.go` | COMPLETED | 30s heartbeat ticker, Updated/Dirty/Failed grouped output with itemized dirty files |
| `Task-04` | Pull-Error Subsystem & Split-DB Persistence | Worker 02 | `cli/store/pull_split_db_errors.go`, `cli/cmdpullerror/`, `cli/cmd/rootcore.go` | COMPLETED | `pull_errors` schema, `RunPullErrorCLI`, AI diagnostic card, rootcore registration |
| `Task-05` | Wincredman Diagnostic Remediation Hints | Worker 01 | `cli/cmdpull/pull_remediation_hint.go` | COMPLETED | Classification for `wincredman`, actionable `gitmap fix-credential` suggestion |

---

## Verification & Integrity Evidence

- **Coding Guidelines & Linters:** Zero builds or test suites executed (Rule R1). All files follow strict positive boolean conventions, structured error wrapping, and relative paths.
- **Split-DB Integration:** `OpenPullSplitDB()` ensures `pull_errors` schema on connection; `recordTelemetryToDB` records all failed pull states into `pull_errors`.
- **Command Routing:** Aliases `pull-error`, `pull-errors`, `pulle`, `pull-e` dispatched from `rootcore.go`.
