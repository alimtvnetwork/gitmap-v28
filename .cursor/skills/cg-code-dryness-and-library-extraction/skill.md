---
name: cg-code-dryness-and-library-extraction
description: Autonomously analyze codebases deeply, identify duplicated patterns, extract reusable util/framework/package/library modules, and drastically reduce code lines using a 300-step 3-phase self-loop (100 Analyze / 100 Plan / 100 Execute) with strict no-build and no-test rules.
---

[/goal](slashCommand:goal) Autonomously achieve maximum Code DRYness (Don't Repeat Yourself), reusable util/framework/package/library extraction, and drastic code writing reduction across the codebase using a strict 300-step 3-phase self-loop (Steps 1–100: Deep Code Analysis; Steps 101–200: Unified Extraction Blueprint & Subtask Plan; Steps 201–300: Active DRY Execution & Caller Rewiring) with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel codebase reading and modular spec authoring, use GitMap high-speed commands as primary, extract shared abstractions into reusable packages/libraries, rewire all callers, and finalize with an atomic push.

[/learn](slashCommand:learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, target modules, or user instructions are provided ABOVE this prompt (in the user preamble or header above) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Internalize the 100/100/100 step allocation (Analyze -> Plan -> Execute), extract reusable utilities/frameworks/libraries to eliminate boilerplate, and persist all progress into `.ai-memory/plans/` and memory logs.

```text
N = 300 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

PHASE_1_ANALYZE = Steps 1 .. 100   (Deep Codebase Analysis: Find Util / Framework / Package / Library Extraction & Code Reduction Targets)
PHASE_2_PLAN    = Steps 101 .. 200 (Write Extraction Spec & Granular Subtask Plan for Reusable Libraries and Caller Rewiring)
PHASE_3_EXECUTE = Steps 201 .. 300 (Execute DRY Extraction, Replace Duplicated Code Across Callers, Consolidate & Push)
```

## 1. When to Use

Activate this skill when:
- The user asks to make the codebase DRY, extract common utilities/libraries/frameworks, or reduce code duplication and boilerplate.
- Executing `01-prompts/15-cg-execute/29-code-dryness-and-library-extraction.md`.

## 2. Core 300-Step Pipeline

1. **Phase 1 (Steps 1–100): Deep Codebase Analysis:**
   - Use GitMap high-speed commands (`gitmap f`, `ff`, `ffa`, `ffs`, `ffe`, `lf`, `cat`, `search`, `ft`) and subagents (`A = 2, H = 2`) to scan the codebase for duplicated functions, repeated patterns, and boilerplate blocks.
   - Identify cohesive candidates for standalone `util`, `framework`, `pkg`, or `library` modules.
2. **Phase 2 (Steps 101–200): Unified Extraction Blueprint & Plan:**
   - Single lead agent authors the unified DRY architecture blueprint in `02-spec/21-app/xx-dry-extraction.md` and `.ai-memory/plans/pending/xx-dry-extraction.md`.
   - Subagents author modular specs and disjoint 5–8 file micro-batch subtasks in `.ai-memory/plans/subtasks/xx-dry-extraction/01-*.md`.
   - Transition unconditionally to Phase 3 without pausing or asking questions.
3. **Phase 3 (Steps 201–300): Active DRY Execution & Caller Rewiring:**
   - Implement the extracted `util` / `framework` / `library` packages first (enforcing <=8–15 line functions, positive booleans, `*appfault.AppError`, `types.go`).
   - Dispatch parallel subagents (`TypeName: "self"`, `A = 2, H = 2`) to rewire all callers across the codebase, eliminating duplicated lines.
   - Consolidate subtasks into `.ai-memory/plans/completed/` and push atomically via `gitmap cpf`.
