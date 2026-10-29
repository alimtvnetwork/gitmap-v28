---
name: next-100-code-improvements
description: Autonomously scan, audit, discover, and document 100 grounded code improvements across 10 coding guideline families using a 10-field specification schema, GitMap streaming search, micro-batching, and pending plan enqueueing.
---

# Next 100 Code Improvements

> **Invocation:** `@[next-100-code-improvements](file;.cursor/skills/next-100-code-improvements)`

Source Prompt: `01-prompts/30-improvements/02-next-100-code-improvements.md`

This skill governs the systematic, autonomous discovery, technical specification, and backlog enqueueing of exactly 100 high-impact code improvements across the repository. It enforces structured 10-field specification cards, balances discovery across 10 canonical rule families, and writes actionable execution plans into `.ai-memory/plans/pending/`.

## When to Use

- When tasked with running a deep, grounded code audit across a codebase or module.
- When generating a structured backlog of 100 refactoring opportunities grounded in repository standards.
- When invoked via slash command `[/next-100-code-improvements](file;.cursor/skills/next-100-code-improvements)` or direct user request.
- When preparing an architectural roadmap before major refactoring sprints.
- When auditing codebases against `02-spec/02-coding-guidelines/` and `02-spec/03-error-manage/`.

## Core Architecture

The audit engine operates under the canonical `execute-parent-task-with-n-steps-v6` architecture:
- **Concurrency Budget:** `N = 400` steps, `A = 3` concurrent subagents (`invoke_subagent`), `H = 2` operational hands per agent (Capacity = 6 concurrent operations).
- **Tree Swarm Hierarchy:** Root Lead Orchestrator -> Branch Coordinators (Rule Family Owners) -> Leaf Workers & Dynamic Joining Assistants.
- **Micro-Batching:** All 100 improvements are discovered and authored in discrete micro-batches of 5–8 files to prevent context exhaustion and token limits.
- **Sequential Wave Generation:** Each execution discovers and enqueues a distinct wave of 100 improvements (`wave-01`, `wave-02`, etc.) without duplicating previously addressed patterns.

## 10 Guideline Families & 10-Field Schema

Exactly 10 improvements are identified under each of the 10 rule families below:

1. **Family 01: Error Management & Structured `*appfault.AppError` Returns (IMP-001 to IMP-010):**
   - Grounded in `02-spec/03-error-manage/02-error-architecture/06-apperror-package/02-apperror-reference.md`.
   - Functions returning error metadata must return `*appfault.AppError`, zero swallowed errors, package stutter elimination (`appfault.AppError`).
2. **Family 02: Monadic `Result[T]` Envelopes & Types Extraction (IMP-011 to IMP-020):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/25-generic-return-types.md`.
   - Replace `(T, error)` with `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`; centralize structs and generic aliases into `types.go`.
3. **Family 03: Positive Booleans (Rules P1–P4) (IMP-021 to IMP-030):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-naming-prefixes.md`.
   - Affirmative `is*` and `has*` prefixes only. Banned single-letter flags (`v bool`) and negative flags (`isNotValid`). Replace `!isEmpty` with `isDefined`.
4. **Family 04: Condition Simplification (Rules P5–P9) (IMP-031 to IMP-040):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md`.
   - Total ban on mixed polarity (`if isA && !isB`). Maximum nesting depth 1 with guard clauses. No inline assignment in `if`.
5. **Family 05: Function Sizing (IMP-041 to IMP-050):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`.
   - Hard ceiling <=15 lines (target 8–10 lines). Map dispatch tables replacing switch blocks.
6. **Family 06: Parameter Structs (The 3-Parameter Cap) (IMP-051 to IMP-060):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`.
   - Refactor functions with >2-3 parameters into dedicated `*Params` / `*Options` structs.
7. **Family 07: File Modularity & Sizing (IMP-061 to IMP-070):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`.
   - Hard ceiling <=100 lines for components, <=300 lines for packages. Split monolithic files into cohesive siblings.
8. **Family 08: Vertical Line Gaps & Whitespace Spacing (IMP-071 to IMP-080):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md`.
   - Blank line before `if`, blank line after `}`, blank line before `return`.
9. **Family 09: String Efficiency & Enums (IMP-081 to IMP-090):**
   - Grounded in `02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/03-boolean-and-enum.md`.
   - Zero-allocation `strutil.EqualFoldAny`, PascalCase enums with mandatory `Type` suffix.
10. **Family 10: Testing Hygiene & Pathutil Abstraction (IMP-091 to IMP-100):**
    - Grounded in `02-spec/02-coding-guidelines/01-cross-language/14-test-naming-and-structure.md`.
    - Unit tests must never trigger destructive OS calls (shutdown, reboot, kill). Mock duration verification, pathutil abstraction.

### The 10-Field Specification Card Schema

Every improvement card must provide all 10 fields verbatim:
1. **Target File Path:** `[strictly/relative/path/to/file.ext]`
2. **File Name:** `file.ext`
3. **Target Function / Symbol / Scope:** `Package.Function`
4. **Line Numbers:** `L10–L25`
5. **Rule Violated & Canonical Citation:** `[Rule Code: Name]` — Citation: `02-spec/...`
6. **Violation Severity:** `[Critical | Major | Minor]`
7. **Current Code Context:** Verbatim offending code block
8. **Violation Analysis & Architectural Rationale:** In-depth explanation of risk and benefit
9. **Exact Fix Proposal:** Copy-paste ready refactored code satisfying all guidelines
10. **Blast Radius, Downstream Callers & Verification Contract:** Callers via `gitmap aum search`, blast radius assessment, and targeted exit 0 verification command

## Workflow & GitMap Execution

### Step 1: Preflight & Working Tree Check
- Confirm `git status --porcelain` is clean.
- Check Heatmap telemetry: `gitmap heatmap status` and `gitmap heatmap tables`.
- Initialize SQLite Task DB: `gitmap agent task init --name "100-improvements" --budget 400`.

### Step 2: Phase 1 Codebase Scanning (A = 3 `research` Subagents)
- Dispatch 3 parallel discovery subagents using GitMap streaming search:
  * Directory file inventory: `gitmap lf <dir>` or `gitmap find "*.*" <dir>`
  * Live symbol and pattern search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
  * Substring file finder: `gitmap find "<pattern>" [dir]`
  * Zero-disk file streaming: `gitmap cat <filepath>`
- Total Ban on shell search tools: never run `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`, `Get-ChildItem`, `dir`, or `ls`.

### Step 3: Phase 2 Card Authoring & Wave Assembly (A = 3 `self` Subagents)
- Partition authoring across subagents with disjoint file ownership (Rule R6).
- Compile the Master Catalog Summary Table (100 rows).
- Author modular subtasks for each rule family.

### Step 4: Phase 3 Verification & Atomic Push
- Validate relative path hygiene: zero absolute paths or `file:///` URIs.
- Verify zero secrets via `python linter-scripts/check-forbidden-strings.py`.
- Atomic GitMap commit and push: `gitmap cpf "<module> - <summary>"` using hyphen format.

## Enqueueing Plans to `.ai-memory/plans/pending/`

All 100 improvements must be enqueued as pending plans in the repository:
1. **Master Plan:** Authored at `.ai-memory/plans/pending/XX-next-100-code-improvements-wave-YY.md`.
2. **Subtasks:** Authored at `.ai-memory/plans/subtasks/XX-next-100-code-improvements-wave-YY/` (10 family files).
3. **Register:** Listed in `.ai-memory/plans/readme.md` under `## Pending Plans`.
4. **Context Register:** Enqueued in `.ai-memory/what-to-read.md`.
5. **Next 100 Continuity:** Subsequent executions read existing waves and generate wave `YY + 1`.

## Verification Criteria

- [ ] All 100 improvements are documented using the strict 10-field specification schema.
- [ ] Exactly 10 items are authored per rule family across all 10 families.
- [ ] Master Catalog Summary Table contains 100 valid rows matching the cards.
- [ ] All code snippets are verbatim extractions from the repository (via `gitmap cat`).
- [ ] Caller blast radius is mapped via `gitmap aum search` for every symbol.
- [ ] Master plan and 10 subtask files are committed to Git and pushed upstream via `gitmap cpf`.
- [ ] All markdown links strictly follow relative Git path syntax.

## Related Skills

- [execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)
- [coding-guidelines](file;.cursor/skills/coding-guidelines)
- [cg-error-management](file;.cursor/skills/cg-error-management)
- [cg-boolean-and-naming](file;.cursor/skills/cg-boolean-and-naming)
- [gitmap](file;.cursor/skills/gitmap)
- [ai-verification](file;.cursor/skills/ai-verification)
