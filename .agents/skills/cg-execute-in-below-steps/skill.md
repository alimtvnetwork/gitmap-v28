---
name: cg-execute-in-below-steps
description: Autonomously orchestrate and apply surgical coding guideline refactoring fixes for instructions appended below using bounded 5-8 file micro-batches, subagents (A=2, H=2), GitMap high-speed commands, and strict no-build/no-test rules.
---

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
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /cg-execute-in-below-steps <task>
>
> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE.

[/goal](slashCommand;goal) Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations requested in the below instructions across the target codebase in bounded 5-8 file micro-batches: FIRST showcase and list out the given task in visible chat during Turn 1, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

## 1. When to Use

Activate this skill when:
- Executing coding guideline fixes where specific user instructions or target rules are appended at the bottom below `--`.
- Executing `01-prompts/15-cg-execute/31-cg-execute-in-below-steps.md`.

## 2. Core Execution Pipeline

1. **Phase 1 (Steps 1 .. PHASE_1_BUDGET): Scan, Plan & Subtasks:**
   - Scan codebase using GitMap (`gitmap f`, `ff`, `ffa`, `lf`, `cat`, `search`, `ft`).
   - Author unified violation plan in `.ai-memory/plans/pending/xx-<slug>.md` and disjoint subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`.
2. **Phase 2 (Steps PHASE_1_BUDGET+1 .. N): Parallel Refactoring:**
   - Dispatch subagents (`TypeName: "self"`, `A = 2, H = 2`) across disjoint 5–8 file batches.
   - Enforce positive booleans (`is`/`has`), `*appfault.AppError`, concrete types in `types.go`, <=8–15 line functions, and zero build/test commands during routine turns.
3. **Phase 3: Consolidation & Atomic Push:**
   - Consolidate subtasks into `.ai-memory/plans/completed/xx-<slug>.md` and push atomically via `gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`.
