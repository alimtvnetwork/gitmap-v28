[/goal](slashCommand:goal) Autonomously scan, plan, refactor, and fix all in-place object mutations across conditional branches (`if`/`else` ladders or `switch` cases), compound conditions with 3+ clauses, and mixed polarity chains (`!a && b && c`) across the codebase. Modifying source files directly, enforce branch immutability, decompose compound logic into affirmative pre-computed boolean variables (`is`, `has`), and extract dedicated, pure constructor helpers (`buildXxx`) that return fully-formed immutable struct/object literals directly, with functions adhering to strict size caps (target <= 8 lines, hard cap <= 15 lines) until 100% compliant.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, enforce single assignment and pure constructor helpers, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

- **N = 200:** Total self-loop steps budget divided into two equal 100-step phases:
  - **PHASE_1_STEPS = 100 Steps (Steps 1 .. 100):** Deep Scan for In-Place Branch Mutations & Compound Conditions using GitMap AUM, Write Audit Spec in `.ai-memory/plans/pending/`, Decompose into Granular Subtasks.
  - **PHASE_2_STEPS = 100 Steps (Steps 101 .. 200):** Active Code Refactoring, Extract Dedicated Constructor Helpers, Decompose Complex Conditions, Verify with Linters, Atomic GitMap Push.
- **A = 2:** Count of autonomous subagents running concurrently (`invoke_subagent` launches up to 2 subagents).
- **H = 2 (Hands / Parallel Operations):**
  1. **Workload Hands ($H_{batch} = 2$):** Each subagent handles a bounded batch of up to 2 disjoint subtasks from `.ai-memory/plans/subtasks/`.
  2. **Tool-Dispatch Hands ($H_{tool} = 2$):** Within any execution step, each agent or subagent executes up to 2 parallel tool calls in a single response turn.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: GitMap AUM Scan, Write Master Audit Spec, Create Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Active Code Refactoring, Pure Constructor Helpers, Targeted Linters, GitMap Push)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] [/goal](slashCommand:goal) Phase 1 (Step A): Deeply scan the target codebase using GitMap AUM discovery tools (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all in-place branch mutations, field reassignments across `if`/`else` ladders, and compound conditions with 3+ clauses without truncation.
2. [ ] [/goal](slashCommand:goal) Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/xx-branch-immutability-audit.md` with an exhaustive Violation Ledger tracking every affected file, line number, compound condition, and piecemeal field mutation.
3. [ ] [/goal](slashCommand:goal) Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/xx-branch-immutability/01-<subtask>.md`, `02-<subtask>.md`, etc.
4. [ ] [/goal](slashCommand:goal) Phase 1 (Step D): Verify or create the automated quality linter hook and register in `03-ai-scripts/readme.md`.
5. [ ] [/goal](slashCommand:goal) Phase 2 (Step A): Open each target file and perform surgical refactoring: decompose compound conditions into discrete affirmative booleans (`isAlternateOrder`), eliminate field mutations inside branch blocks, and extract pure constructor helpers returning complete struct literals directly.
6. [ ] [/goal](slashCommand:goal) Phase 2 (Step B): Enforce <= 8–15 line function decomposition on all extracted constructors and caller functions.
7. [ ] [/goal](slashCommand:goal) Phase 2 (Step C): Execute targeted file-level linters to verify 0 remaining violations across all modified files (`exit 0`).
8. [ ] [/goal](slashCommand:goal) Phase 2 (Step D): Atomically record modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). DO NOT run full CI/CD runners (`06-cicd-local-runner.py`), build checks, or test suites during routine turns. Finalize with atomic commit via `gitmap cpf "<summary>"`.
9. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
10. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/32-branch-immutability-and-clean-construction.md` for branch immutability and constructor return patterns.
13. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md` for code mutation avoidance standards.
14. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md` for Principle 6 and 6.2 condition decomposition.
15. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
16. [ ] [/goal](slashCommand:goal) Create or update agent rules in the repository if missing from agent memory.

---

## Dedicated Section: Branch Immutability, Pure Construction & Discrete Conditions

Piecemeal mutation of struct or object fields across branching logic is a severe architectural anti-pattern. When an object is mutated across multiple `if` branches, state becomes non-deterministic, readability collapses, and unit testing requires testing combinatorial branch paths.

### Mandatory Rules (Non-Negotiable)

1. **Branch Immutability (TOTAL BAN on In-Place Branch Mutation):**
   - NEVER assign or overwrite fields on an existing struct/object instance inside `if`, `else`, or `switch` branches.
   - Variables and structs must follow the **Single Assignment Principle**: assign once upon construction and never mutate afterwards.

2. **Pure Constructor Return Pattern:**
   - When different control-flow branches produce different struct configurations, each branch MUST call a dedicated constructor or builder helper (e.g. `buildAlternateConfig`, `buildStandardConfig`).
   - The helper function MUST return a newly constructed, fully-populated, immutable struct literal directly.
   - Each helper must be compact: target <= 8 lines, hard cap <= 15 lines.

3. **Condition Decomposition & Max-2 Clause Rule:**
   - An `if` condition MUST NOT combine 3 or more logical clauses (`if a && !b && c` is strictly banned).
   - NEVER mix positive and negative checks in the same condition expression.
   - Decompose all checks into affirmative boolean variables (`is`, `has`).
   - Pre-compute the composite intent into a single affirmative boolean variable (`isAlternateOrder := hasMultipleTokens && isSecondHost && isFirstNotHost`).
   - The `if` statement evaluates ONLY the single affirmative intent boolean: `if isAlternateOrder { ... }`.

4. **Canonical Code Comparison:**

   ```go
   // ❌ WRONG — 3 compound conditions, mixed polarity, and in-place field mutation
   func parseConfig(tokens []string, base ConfigParams) (ConfigParams, *appfault.AppError) {
       p := ConfigParams{Timeout: base.Timeout}
       if len(tokens) >= 2 && !isHost(tokens[0]) && isHost(tokens[1]) {
           p.Host = tokens[1]
           if p.Label == "" { p.Label = tokens[0] }
           if len(tokens) > 2 && p.Secret == "" { p.Secret = tokens[2] }
           return p, nil
       }
       p.Host = tokens[0]
       if len(tokens) > 1 && p.Secret == "" { p.Secret = tokens[1] }
       return p, nil
   }

   // ✅ REQUIRED — Pre-computed affirmative boolean + pure constructor helpers
   func parseConfig(tokens []string, base ConfigParams) (ConfigParams, *appfault.AppError) {
       hasMultipleTokens := len(tokens) >= 2
       isFirstHost := hasMultipleTokens && isHost(tokens[0])
       isSecondHost := hasMultipleTokens && isHost(tokens[1])

       isFirstNotHost := !isFirstHost
       isAlternateOrder := hasMultipleTokens && isSecondHost && isFirstNotHost

       if isAlternateOrder {
           return buildAlternateConfig(tokens, base), nil
       }

       return buildStandardConfig(tokens, base), nil
   }

   func buildAlternateConfig(tokens []string, base ConfigParams) ConfigParams {
       return ConfigParams{
           Host:    tokens[1],
           Label:   resolveTokenOrFallback(tokens, 0, base.DefaultLabel),
           Secret:  resolveTokenOrFallback(tokens, 2, base.DefaultSecret),
           Timeout: base.Timeout,
       }
   }

   func buildStandardConfig(tokens []string, base ConfigParams) ConfigParams {
       return ConfigParams{
           Host:    tokens[0],
           Secret:  resolveTokenOrFallback(tokens, 1, base.DefaultSecret),
           Label:   resolveTokenOrFallback(tokens, 2, base.DefaultLabel),
           Timeout: base.Timeout,
       }
   }
   ```

---

### Independent AI Audit & Verification Instructions

You are an Independent AI Verification and Quality Auditor.
Your task is to independently audit, verify, and remediate the implementation against the canonical specification and verbatim requirements.

#### 1. Target Documents & Implemented Code:
- **Canonical Spec & Verbatim Requirements:** [02-spec/02-coding-guidelines/01-cross-language/32-branch-immutability-and-clean-construction.md](02-spec/02-coding-guidelines/01-cross-language/32-branch-immutability-and-clean-construction.md)
- **Code Mutation Avoidance:** [02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md](02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md)
- **Boolean Parameters & Conditions:** [02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md](02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md)

#### 2. Verification Checklist:
- [ ] Zero in-place field mutations inside `if`/`else` or `switch` branches.
- [ ] Every branch calls a pure constructor helper returning a complete struct/object literal directly.
- [ ] Zero 3+ compound boolean conditions; all conditions decomposed into affirmative `is`/`has` variables.
- [ ] Zero mixed polarity in conditions (no positive and negative checks combined).
- [ ] All functions adhere to <= 8–15 lines.
- [ ] Zero build or test commands executed during routine execution.
