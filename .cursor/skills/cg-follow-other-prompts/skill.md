---
name: cg-follow-other-prompts
description: Autonomously ingest, follow, and execute referenced external prompts and coding guideline directives across the codebase using bounded 5-8 file micro-batches, subagents (A=2, H=2), GitMap acceleration, and strict no-build/no-test rules.
---

[/goal](slashCommand:goal) Autonomously ingest, follow, and execute the referenced external prompts, task instructions, and coding guideline directives across the target codebase in bounded 5-8 file micro-batches with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and rule extraction, read the codebase using GitMap AUM as primary, establish a single-agent blueprint during Phase 1 planning (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever prompt paths, directives, custom rules, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the referenced prompts and bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

## 1. When to Use

Activate this skill when:
- The user instructs the agent to follow one or more referenced prompt files or external directives under coding guideline enforcement.
- Executing `01-prompts/15-cg-execute/32-cg-follow-other-prompts.md`.

## 2. Core Execution Pipeline

1. **Phase 1 (Steps 1 .. N/2): Prompt Ingestion, Scan & Plan:**
   - Read all referenced prompt files via `gitmap cat` or fast file reader.
   - Synthesize directives into `.ai-memory/plans/pending/xx-<slug>.md` and disjoint subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`.
2. **Phase 2 (Steps N/2+1 .. N): Parallel Execution:**
   - Dispatch subagents (`TypeName: "self"`, `A = 2, H = 2`) across disjoint 5–8 file batches with zero build/test commands during routine turns.
3. **Phase 3: Consolidation & Atomic Push:**
   - Consolidate subtasks into `.ai-memory/plans/completed/xx-<slug>.md` and push atomically via `gitmap cpf` / `gitmap cpb`.
