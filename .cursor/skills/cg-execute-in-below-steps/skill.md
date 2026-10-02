---
name: cg-execute-in-below-steps
description: Autonomously orchestrate and apply surgical coding guideline refactoring fixes for instructions appended below using bounded 5-8 file micro-batches, subagents (A=2, H=2), GitMap high-speed commands, and strict no-build/no-test rules.
---

[/goal](slashCommand:goal) Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations requested in the below instructions across the target codebase in bounded 5-8 file micro-batches with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and AST inspection, use GitMap AUM as primary, establish a single-agent blueprint during Phase 1 planning (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

## 1. When to Use

Activate this skill when:
- Executing coding guideline fixes where specific user instructions or target rules are appended at the bottom below `--`.
- Executing `01-prompts/15-cg-execute/31-cg-execute-in-below-steps.md`.

## 2. Core Execution Pipeline

1. **Phase 1 (Steps 1 .. N/2): Scan, Plan & Subtasks:**
   - Scan codebase using GitMap (`gitmap f`, `ff`, `ffa`, `lf`, `cat`, `search`, `ft`).
   - Author unified violation plan in `.ai-memory/plans/pending/xx-<slug>.md` and disjoint subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`.
2. **Phase 2 (Steps N/2+1 .. N): Parallel Refactoring:**
   - Dispatch subagents (`TypeName: "self"`, `A = 2, H = 2`) across disjoint 5–8 file batches.
   - Enforce positive booleans (`is`/`has`), `*appfault.AppError`, concrete types in `types.go`, <=8–15 line functions, and zero build/test commands during routine turns.
3. **Phase 3: Consolidation & Atomic Push:**
   - Consolidate subtasks into `.ai-memory/plans/completed/xx-<slug>.md` and push atomically via `gitmap cpf` / `gitmap cpb`.
