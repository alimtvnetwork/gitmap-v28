---
name: cg-execute-in-below-steps
description: Autonomously orchestrate and apply surgical coding guideline refactoring fixes for instructions appended below using bounded 5-8 file micro-batches, subagents (A=3, H=2, Capacity=6), GitMap high-speed commands, and strict no-build/no-test rules.
---

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

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /cg-execute-in-below-steps <task>
>
> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

[/goal](slashCommand;goal) Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations requested in the below instructions across the target codebase in bounded 5-8 file micro-batches: FIRST showcase and list out the given task in visible chat during Turn 1, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

---

## 1. When to Use

Activate this skill when:
- Executing coding guideline fixes where specific user instructions or target rules are appended at the bottom below `--`.
- Executing `01-prompts/15-cg-execute/31-cg-execute-in-below-steps.md`.

---

## 2. Step 0 Preflight with Data Schema Verification

Before enqueuing subtasks or claiming work, initialize and verify the SQLite task database schema:

```bash
gitmap agent task init --name "<task name>" --budget 400
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 3. 3-Tier Database Location Transparency & Diagnostic Suite

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
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <refactor-rationale>" -s "EXECUTION"

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

## 4. Core Execution Pipeline

1. **Phase 1 (Steps 1 .. 100): Scan, Plan & Subtasks:**
   - Scan codebase using GitMap (`gitmap find`, `gitmap aum search`, `gitmap cat`, `gitmap list-files`).
   - Author unified violation plan in `.ai-memory/plans/pending/xx-<slug>.md` and disjoint subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`.
2. **Phase 2 (Steps 101 .. 250): Intermediate Checkpoints & Sanitization:**
   - Pre-clean stale caches and set up disjoint subtask files.
3. **Phase 3 (Steps 251 .. 400): Parallel Refactoring & Consolidation:**
   - Dispatch subagents (`TypeName: "self"`, `A = 3, H = 2`, Capacity = 6) across disjoint 5–8 file batches.
   - Enforce positive booleans (`is`/`has`), `*appfault.AppError`, concrete types in `types.go`, <=8–15 line functions, and zero build/test commands during routine turns.
   - Consolidate subtasks into `.ai-memory/plans/completed/xx-<slug>.md` and push atomically via `gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`.

---

## 5. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts, Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 6. Verification Checklist

- [ ] V6 Concurrency Invariants active (`budget 400`, `A = 3`, `H = 2`, `Capacity = 6`).
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.

--
## 🚨 High Priority Instructions Below
