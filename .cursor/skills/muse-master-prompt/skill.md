---
name: muse-master-prompt
description: >-
  Executes the Muse Master Onboarding & Canonical V6 Autonomous Execution prompt. Use when the user asks to run muse-master-prompt, onboard a Muse AI agent, execute with Muse turbo mode, or apply the V6 task confirmation protocol.
---

# [V6] Muse Master Onboarding & Autonomous Execution Skill

Source prompt: [`01-prompts/27-muse-prompts/01-muse-master-prompt.md`](01-prompts/27-muse-prompts/01-muse-master-prompt.md)
Canonical Execution Engine: [`01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`](01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md)

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
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt/skill (in the user preamble or header blocks) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously boot the session through Phases 0–5 without permission-seeking, FIRST showcase the confirmed task breakdown in visible chat during Turn 1, capture requests verbatim, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution of multi-part work is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish each task with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format pushed immediately.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives are provided ABOVE this prompt outrank everything below. Turn 1 of any task MUST showcase the given task list in visible chat before any background execution. Operating contract lives in memory; active tasks live in SQLite task manager and ledger.

[/plan](slashCommand;plan) Execute thorough step-by-step planning before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) at every stage of multi-part work:
  1. **Planning Step (A = 2 `research` subagents):** Spawn 2 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified plan.
  2. **Spec Step (A = 2 `self` subagents or lead):** Spawn 2 subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks.
  3. **Execution Step (A = 2 `self` worker subagents):** Spawn 2 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes solo for multi-part tasks without calling the `invoke_subagent` tool.

---

## Operating Instructions

### 1. Onboarding Protocol (Phases 0–5)
- **Phase 0:** Ask for the repository table in one message (ask once, all at once).
- **Phase 1:** Terminal check (`echo ok`), `git --version`, `gh auth status`, `gitmap login --status`. Stop if blocked on auth.
- **Phase 2:** GitMap intake: clone, build, learn the 5-phase SOP (discover → modify → verify → commit+push → telemetry). Reuse GitMap commands first (`gitmap aum search`, `gitmap find`, `gitmap cat`).
- **Phase 3:** Memory operating contract checklist: TURBO-01 through TURBO-05, GUARD-01 through GUARD-08 (never delete repos, confirm file deletions, prove claims).
- **Phase 4:** Guideline intake: booleans (`is*`/`has*`, no `== true`), guard clauses, `*appfault.AppError`, lowercase filenames, strict relative git paths.
- **Phase 5:** Send one short ready message, then ask for the first task.

### 2. Task Confirmation & Execution Protocol (Every Task)
1. **Verbatim Capture:** Record incoming request losslessly.
2. **Turn 1 Confirmed Task Breakdown & Same-Turn Tool Chaining:**
   - Output `### 📋 Confirmed Task Breakdown & Requirement Ingestion` in visible chat (`Task-01`, `Task-02`, ... each with `State: [IN PROGRESS]` and `Understood: [YES]`).
   - In the **EXACT SAME TURN**, invoke your first tool call (SQLite task DB initialization, ledger creation, preflight checks). NEVER emit text alone.
3. **SQLite Task DB & Ledger Preflight:**
   - Run `python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "<task name>" --budget 300`.
   - If `RESUME_FOUND`, inspect crash forensics and resume uncompleted subtasks.
   - Maintain human-readable `.ai-memory/temp-agents/<slug>/ledger.md`.
4. **Mandatory Multi-Agent Execution:**
   - Dispatch `A = 2` subagents in disjoint file boxes (`TypeName: "research"` for discovery, `TypeName: "self"` for execution).
   - Worker Git Ban: Subagents NEVER run git commands (`git add`, `git commit`, `git status`) to avoid `.git/index.lock` collisions.
5. **Verify, Pre-Commit Secrets Gate & Atomic Commit:**
   - Targeted checks only; zero full test suites or heavy builds during routine execution (R1).
   - Pre-commit secrets gate: `python linter-scripts/check-forbidden-strings.py` and `gitmap aum search` regex. Offload via `gitmap rs`.
   - Commit and push atomically via GitMap using hyphen format (no colons in arguments):
     - Feature: `gitmap cpf "<module> - <summary>"`
     - Bug Fix: `gitmap cpb "<module> - <summary>"`
   - Report briefly and ask for the next task.
