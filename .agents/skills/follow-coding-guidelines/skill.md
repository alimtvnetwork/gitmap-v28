---
name: follow-coding-guidelines
description: Autonomously enforce, audit, and follow repository coding guidelines, zero-nesting rules, leaf error sovereignty, parameter structs, positive booleans, and GitMap commit discipline.
---

# Follow Coding Guidelines Autonomous Skill

> **Invocation:** `@[follow-coding-guidelines](file;.agents/skills/follow-coding-guidelines)`

## Overview

This skill establishes the non-negotiable coding standards, architectural guardrails, and audit checklists across all codebases managed by Prompt Architect. Every AI agent, subagent, and human engineer operating in this repository MUST strictly adhere to these 12 grounded rules.

Whenever executing a coding task, invoke `[/learn](slashCommand;learn)` to ingest and internalize these guidelines, and `[/plan](slashCommand;plan)` first before touching code to eliminate friction, prevent rework, and maintain codebase purity.

---

## The `/learn` Continuous Learning Loop

Autonomous agents MUST maintain high cognitive awareness and learn from every task execution:

1. **Pre-Task Learning Intake:**
   - Review active conventions in `02-spec/02-coding-guidelines/`.
   - Check `.ai-memory/memory/` and `.ai-memory/what-to-read.md` to load project identity, architectural boundaries, and recent patterns.
   - Run `[/learn](slashCommand;learn)` to reflect on resolved ambiguities, anti-patterns, and project-specific invariants.

2. **Mid-Flight Guardrail Enforcement:**
   - Actively check each code edit against the 12 Grounded Rules Checklist below before staging or saving.
   - Never introduce temporary workarounds or shortcut antipatterns.

3. **Post-Task Retrospective Synthesis:**
   - Record newly discovered architectural boundaries or domain lessons into `.ai-memory/memory/learned/`.
   - Update repository documentation and specifications to keep knowledge evergreen.

---

## The 12 Grounded Coding Guideline Rules Checklist

### Rule 1: No Generated Code, Artifacts, or Temporary Fixtures in Version Control
- **Zero Artifact Commits:** Never commit `.tmp`, build caches, log files, node_modules, `.DS_Store`, or ephemeral test harnesses to standard git repositories.
- **Offload to `repo-cache`:** Diagnostic scripts, scratch runners, and temporary PowerShell harnesses belong in `repo-cache` (`gitmap rc file ...` or `gitmap rc text ...`).
- **Offload to `repo-secrets`:** Environment variables, PAT tokens, passwords, and API keys belong strictly in `repo-secrets` (`gitmap rs file ...` or `gitmap rs text ...`).
- **Clean Git Hygiene:** Tracked folders must immediately have redundant `.gitkeep` files removed.

### Rule 2: Strict Function & File Sizing Caps
- **Function Limit:** Standard functions MUST be 8 to 15 lines of executable logic maximum. Functions exceeding 15 lines must be systematically decomposed into discrete helper functions.
- **File Limit:** Standard source files MUST be 150 to 300 lines maximum. Files exceeding 300 lines must be split into modular sibling files.
- **React Component Limit:** React UI component files MUST be 100 lines maximum. Decompose sub-components, custom hooks, and styles into separate files.

### Rule 3: Zero Nested If Statements (Early Return & Guard Clause Primacy)
- **Flatten All Nesting:** Never nest `if` blocks more than 1 level deep.
- **Guard Clauses First:** Validate preconditions and error states at the top of functions using inverted conditions and immediate returns.
- **Discrete Handlers:** Replace deeply nested branch bodies with discrete, named helper functions.

### Rule 4: Positive Boolean Standard & Implicit Evaluation (TOTAL BAN on `== true`)
- **Prefix Standard:** All boolean variables, fields, and predicate functions MUST use strictly positive prefixes: `is` or `has` (e.g. `isReady`, `hasPermission`, `isValid`). Banned prefixes: `can`, `should`, `was`, `did`, `allow`.
- **Total Ban on Explicit True/False Checks:** NEVER evaluate a boolean explicitly against `true` or `false` (e.g. ❌ `if isReady == true`). Positive booleans MUST ALWAYS be evaluated implicitly (✅ `if isReady { ... }`).
- **No Mixed Polarity:** NEVER combine a positive check and a negative check in the same condition (e.g. ❌ `if isA && !isB`). Split into discrete sequential guard clauses.

### Rule 5: Structured Error Management & Monadic Results
- **Standard Error Return Type:** All Go packages returning structured failure metadata MUST return `*appfault.AppError`.
- **Package Naming:** Use package `appfault` (`04-code/golang/pkg/appfault`) with struct `AppError` to eliminate redundant package stutter (`appfault.AppError`).
- **Monadic Result Containers:** Prefer single-return Result wrappers (`result.Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`) over awkward multi-value return tuples `(T, error)`.
- **Leaf Error Sovereignty:** Catch and handle errors at the point of origin; never swallow errors or convert them into ambiguous strings without context.

### Rule 6: Parameter Structs / DTOs (Total Ban on Loose >2-3 Arguments)
- **Max Arguments:** Functions must not accept more than 2 to 3 loose parameters.
- **Parameter Struct Standard:** Any function requiring 3 or more arguments MUST encapsulate them into a dedicated value-based `*Params` or `*Config` struct (e.g. `func createWorkflow(params WorkflowCreateParams) Result[Workflow]`).
- **Vertical Spacing:** Multiline struct instantiations must have vertical whitespace separation before and after the invocation.

### Rule 7: Discrete Condition Decomposition & Branch Immutability
- **Branch Immutability:** Variables assigned within conditional branches must remain immutable. Never mutate shared in-place variables across alternative branches.
- **Pure Constructor Helpers:** Return newly constructed values or copy-on-write structs rather than mutating input arguments in place.
- **Discrete Logic:** Decompose complex multi-clause boolean expressions (`(a || b) && (c || d)`) into named boolean intermediate variables before passing to conditional statements.

### Rule 8: Mandatory Vertical Line Gaps (Spacing Discipline)
- **Mandatory Blank Lines:**
  - Mandatory blank line BEFORE every `if` statement (unless it is the very first line of a block).
  - Mandatory blank line AFTER every closing brace `}`.
  - Mandatory blank line BEFORE every `return` statement.
  - Mandatory blank lines around multiline struct initializations and function calls.
- **Clean Visual Rhythm:** Code must breathe vertically. Crammed, wall-of-code syntax is an auto-reject guideline failure.

### Rule 9: String Efficiency, Normalization & Enums
- **Zero Raw Literals:** Banned repeated string literals scattered across code. Define domain constants or string enums with `Type` suffix.
- **Case-Insensitive Normalization:** Use `strutil.EqualFoldAny` / `strings.EqualFold` for zero-allocation comparisons rather than chained `toLower() == "..." || toLower() == "..."`.
- **Loop Hoisting:** Hoist regular expressions and string transformations outside of hot loops.

### Rule 10: Map Dispatch Over Switch / If-Else Ladders
- **Declarative Lookup:** Replace repetitive chained `if-else if` ladders and large `switch` statements matching string keys or enum tokens with declarative map dispatch tables (`map[string]HandlerFunc`).
- **Sub-Millisecond Routing:** Constant-time `O(1)` map lookup ensures maximum throughput and clean separation of concerns.

### Rule 11: SQLite Split-DB Architecture & Isolation
- **Three-Tier Database Segregation:**
  - Global tier: `gitmap.db`
  - Installation/Tooling tier: `installation.db`
  - Repository-specific tier: `repodb/<repo>.db`
- **Zero DB File Lock Collisions:** Multi-AI telemetry MUST use WAL mode and route shared traces to `repo-cache` (`rc-heatmap.db`) to avoid workspace file locks.
- **Repository Pattern:** Encapsulate all SQL queries behind repository interfaces. Never leak raw database queries into presentation or CLI layers.

### Rule 12: GitMap Semantic Commit Discipline & 'No Push = Not Done' Invariant
- **Hyphen-Separated Message Format:** Never use colons inside commit arguments:
  - Feature: `gitmap cpf "<Module> - <summary>"`
  - Bug Fix: `gitmap cbf "<Module> - <summary>"`
  - Chore: `gitmap cpc "<Module> - <summary>"`
  - Release: `gitmap cpr "vX.Y.Z - <summary>"`
- **Subsystem Module Standard:** `<Module>` must represent a concrete architectural component in Title Case (`Auth`, `Storage`, `Scanner`, `Ui`, `Guidelines`) or lowercase (`auth`, `storage`, `scanner`). TOTAL BAN on `commit` as a module name.
- **Local vs Remote Commits:** Use flat commits (`gitmap commit "<msg>"`) for quick intermediate iteration.
- **Mandatory Completion Invariant:** Every task MUST conclude with a semantic auto-pushing commit or `gitmap push`. If code is uncommitted or unpushed, the task is strictly INCOMPLETE.

---

## AI Pre-Commit Self-Audit Matrix

Before finalizing any task or declaring success, every agent MUST verify:

| Check # | Guideline Requirement | Verified? |
| :--- | :--- | :---: |
| 1 | All new and modified filenames use strictly lowercase (`readme.md`, `skill.md`). | [ ] |
| 2 | All paths in specs, plans, code, and comments use relative Git paths (NO `file:///`). | [ ] |
| 3 | Functions adhere to 8–15 lines; files adhere to 150–300 lines (React 100). | [ ] |
| 4 | No nested `if` statements; guard clauses placed at top of functions. | [ ] |
| 5 | No explicit `== true` / `== false` checks; booleans use `is` / `has` prefixes. | [ ] |
| 6 | No mixed polarity in conditional statements. | [ ] |
| 7 | Functions with >2-3 parameters use dedicated parameter structs. | [ ] |
| 8 | Mandatory vertical line gaps (before `if`, after `}`, before `return`). | [ ] |
| 9 | Go functions return `*appfault.AppError` or monadic `Result[T]`. | [ ] |
| 10 | GitMap commit messages format as `<Module> - <summary>` with Title Case or lowercase module. | [ ] |
| 11 | Changes committed and pushed upstream (`gitmap cpf` / `cpc` / `cbf` / `cpr`). | [ ] |
| 12 | CI/CD pipelines remain enabled and green. | [ ] |
