---
name: cg-golang-pointer-reduction-and-value-semantics
description: Autonomously audits, plans, and refactors Golang codebases to reduce pointer usage (*T) and enforce value semantics across small structs, parameters, and query receivers using a 400-step 3-phase self-loop.
---

# Skill: Golang Pointer Reduction & Value Semantics (`cg-golang-pointer-reduction-and-value-semantics`)

> **Invocation:** `@[cg-golang-pointer-reduction-and-value-semantics](file;.cursor/skills/cg-golang-pointer-reduction-and-value-semantics)`

This skill governs autonomous execution for migrating Golang codebases to value semantics, reducing heap allocations, and eliminating pointer indirection and nil-dereference panics across all Go packages.

## 1. Concurrency Invariants & Header Parameters

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Fleet Sync, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

---

## 2. The 400-Step 3-Phase Architecture

- **Phase 1: Deep Audit & Struct Sizing (Steps 1–100):**
  - Inventory all struct declarations, pointer fields, pointer parameters (`*T`), pointer receivers (`func (s *T)`), and pointer return types.
  - Calculate memory footprints: classify structs into Value Candidates (<= 64 bytes) vs. Pointer Retainers (> 64–128 bytes, or containing mutexes).
  - Write master audit spec in `.ai-memory/plans/pending/xx-pointer-reduction-audit.md` with an exhaustive Violation Ledger.
- **Phase 2: Architectural Specification & Type Selection (Steps 101–250):**
  - Author canonical architecture spec in `02-spec/21-app/xx-golang-value-semantics.md`.
  - Define clear migration rules: query methods to value receivers, small parameter structs by value, slice/map pointer elimination.
  - Decompose into granular subtasks in `.ai-memory/plans/subtasks/`.
- **Phase 3: Active Refactoring & Quality Verification (Steps 251–400):**
  - Execute surgical refactoring in micro-batches (5–8 files).
  - Migrate small structs and read-only receivers to value semantics.
  - Verify with targeted file-level linters (`exit 0`).
  - Consolidate plans and finalize with an atomic push via GitMap hyphen commit.

---

## 3. Pointer Justification Matrix

| Construct | Recommendation | Justification |
| :--- | :---: | :--- |
| **Small Structs (<= 64 bytes)** | `T` (Value) | Copying 1–4 words is faster than heap allocation and pointer dereference. |
| **Large Structs (> 64–128 bytes)** | `*T` (Pointer) | Copy overhead exceeds pointer cost; pass by reference. |
| **Read-Only / Query Methods** | `(s T)` Value Receiver | Pure functions with no state mutation; prevents accidental modification. |
| **Stateful Mutation Methods** | `(s *T)` Pointer Receiver | Method genuinely mutates fields on the caller's struct instance. |
| **Structs with Mutexes** | `*T` (Pointer) | `sync.Mutex`, `sync.RWMutex`, or atomic fields MUST NOT be copied. |
| **Slices & Maps** | `[]T` / `map[K]V` (Value) | Slices and maps are already pointer-backed header structs; never use `*[]T`. |
| **Channels** | `chan T` (Value) | Channels are already runtime pointers; never use `*chan T`. |
| **Tri-State Nilability** | `*T` (Pointer) | Distinguishing absent `nil` from empty zero-value `""` / `0`. |
| **Structured Errors** | `*appfault.AppError` | Mandatory exception to conform to Go error return type standards. |

---

## 4. Step 0 Preflight with Data Schema Verification

Before enqueuing subtasks or claiming work, initialize and verify the SQLite task database schema:

```bash
gitmap agent task init --name "<task name>" --budget 400
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 5. 3-Tier Database Location Transparency & Diagnostic Suite

Whenever executing or logging subtask telemetry, fulfill the Database Location Transparency Contract:

| Database Tier | Physical File Location | Access Mode & Concurrency Protocol | Data Retained & Ownership |
| :--- | :--- | :--- | :--- |
| **Tier 1: Machine-Local Root Database** | `~/.gitmap/agent-root.db` | Machine sessions, cross-repo leases | Machine Sessions & Cross-Repo Leases |
| **Tier 2: Repo-Cache Companion Storage** | `repo-cache/.heatmap/<ai-slug>/<repo-slug>/heatmap.db` | Isolated Split-DB with periodic sync (`gitmap heatmap rc sync`) | Cross-repository multi-agent actions, tool usage traces, commit checkpoints |
| **Tier 3: In-Flight Local Run Database** | `.ai-memory/temp-agents/<slug>/agent-task.db` | WAL Mode (`PRAGMA journal_mode=WAL; busy_timeout=5000`) | Task status, worker claims, subtasks, in-flight action logs (<1ms Go engine) |
| **Tier 3: Local Heatmap Cache** | `.ai-memory/heatmap.db` | WAL Mode + Read-Only Queries (`mode=ro`) | Cognitive action telemetry, file exploration reads/edits, contention leases |

### Executable Diagnostic Commands:

```bash
# Pre-edit contention and lease check
gitmap heatmap discuss [--window 15m] [--threshold 2] [--json]

# Cognitive file exploration and editing telemetry
gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <explanation>" -s "EXPLORATION"
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <value-semantics-refactor>" -s "EXECUTION"

# Real-time continuous activity stream & tree
gitmap heatmap tree [-t <task-slug>]
gitmap heatmap live

# Read-only companion query and schema inspection
gitmap heatmap rc query "SELECT ActionType, FilePath, ReasoningText FROM AgentActionLog ORDER BY AgentActionLogId DESC LIMIT 5" --format table
gitmap heatmap rc tables
gitmap heatmap rc view AgentActionLog --limit 10
gitmap heatmap rc sync
```

---

## 6. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine turns. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run test suites (`go test ./...`, `pytest`) during routine execution.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`) after individual file edits.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 7. Verification Checklist

- [ ] Small structs (<= 64 bytes) converted to value semantics for parameters and query receivers.
- [ ] Heap allocations minimized; unnecessary pointer indirection eliminated.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
