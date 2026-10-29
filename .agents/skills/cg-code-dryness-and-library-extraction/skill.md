---
name: cg-code-dryness-and-library-extraction
description: Autonomously analyze codebases deeply, identify duplicated patterns, extract reusable util/framework/package/library modules, and drastically reduce code lines using a 400-step 3-phase self-loop (100 Analyze / 150 Plan / 150 Execute) with strict no-build and no-test rules.
---

# Skill: Code DRYness & Reusable Library Extraction (`cg-code-dryness-and-library-extraction`)

> **Invocation:** `@[cg-code-dryness-and-library-extraction](file;.agents/skills/cg-code-dryness-and-library-extraction)`

This skill governs autonomous deep codebase analysis, identification of duplicated patterns, extraction of reusable utilities, frameworks, packages, and libraries, and caller rewiring across polyglot repositories.

## 1. Concurrency Invariants & Header Parameters

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_ANALYZE = Steps 1 .. 100   (Deep Codebase Analysis: Find Util / Framework / Package / Library Targets, Up-Front Brain & Roadmap)
PHASE_2_PLAN    = Steps 101 .. 250 (Write Extraction Spec & Granular Subtask Plan for Reusable Libraries and Caller Rewiring)
PHASE_3_EXECUTE = Steps 251 .. 400 (Execute DRY Extraction, Replace Duplicated Code Across Callers, Dynamic Join Saturation, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

[/goal](slashCommand;goal) Autonomously achieve maximum Code DRYness (Don't Repeat Yourself), reusable util/framework/package/library extraction, and drastic code writing reduction across the codebase using a strict 400-step 3-phase self-loop (Steps 1–100: Deep Code Analysis; Steps 101–250: Unified Extraction Blueprint & Subtask Plan; Steps 251–400: Active DRY Execution & Caller Rewiring) with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 3, H = 2, Capacity = 6) for parallel codebase reading and modular spec authoring, use GitMap high-speed commands as primary, extract shared abstractions into reusable packages/libraries, rewire all callers, and finalize with an atomic push via GitMap hyphen commit.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, target modules, or user instructions are provided ABOVE this prompt (in the user preamble or header above) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Internalize the 100/150/150 step allocation (Analyze -> Plan -> Execute), extract reusable utilities/frameworks/libraries to eliminate boilerplate, and persist all progress into `.ai-memory/plans/` and memory logs.

---

## 2. Core 400-Step Pipeline

1. **Phase 1 (Steps 1–100): Deep Codebase Analysis & Roadmap:**
   - Use GitMap high-speed commands (`gitmap find`, `gitmap aum search`, `gitmap cat`, `gitmap list-files`) and subagents (`A = 3, H = 2`) to scan the codebase for duplicated functions, repeated patterns, and boilerplate blocks.
   - Lead agent authors the Brain roadmap with target inventory, caller blast radius, and disjoint subtask partitions.
2. **Phase 2 (Steps 101–250): Unified Extraction Blueprint & Plan:**
   - Single lead agent authors the unified DRY architecture blueprint in `02-spec/21-app/xx-dry-extraction.md` and `.ai-memory/plans/pending/xx-dry-extraction.md`.
   - Subagents author modular specs and disjoint 5–8 file micro-batch subtasks in `.ai-memory/plans/subtasks/xx-dry-extraction/01-*.md`.
   - Transition unconditionally to Phase 3 without pausing or asking questions.
3. **Phase 3 (Steps 251–400): Active DRY Execution & Caller Rewiring:**
   - Implement the extracted `util` / `framework` / `library` packages first (enforcing <=8–15 line functions, positive booleans, `*appfault.AppError`, `types.go`).
   - Dispatch parallel subagents (`TypeName: "self"`, `A = 3, H = 2`, Capacity = 6) to rewire all callers across the codebase, eliminating duplicated lines. Immediate dynamic respawn and dynamic join modes maintain 100% capacity.
   - Consolidate subtasks into `.ai-memory/plans/completed/` and push atomically via `gitmap cpf "<module> - <summary>"`.

---

## 3. Step 0 Preflight with Data Schema Verification

Before enqueuing subtasks or claiming work, initialize and verify the SQLite task database schema:

```bash
gitmap agent task init --name "<task name>" --budget 400
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 4. 3-Tier Database Location Transparency & Diagnostic Suite

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

## 5. Routine Execution & GitMap Hyphen Commit Policy

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

## 6. Verification Checklist

- [ ] Extracted utilities eliminate duplicate code and boilerplate.
- [ ] Functions adhere to <= 8 lines (hard cap: 15 lines).
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
