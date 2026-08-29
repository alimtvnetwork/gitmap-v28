# Instruction (must follow): Plan: Coding Guideline Audit & Enforcement (v4)

/goal Deeply audit the entire codebase for coding guideline violations, boolean anti-patterns, missing enums, cyclomatic complexity, and error-handling flaws. Structure all findings into actionable, fine-grained tasks in `.lovable/plans/pending/` and `.lovable/plans/subtasks/` before stopping.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.** 
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with `[N]` placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memories/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to.

Stop doing that. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0.

---

## Variables — Auto-Discovered at Runtime

```text
N = 150 (Default number of steps the planning AI should take to generate the audit plan. The user may override this when triggering the prompt)
```

/learn Ingest, analyze, and internalize all coding guidelines, boolean principles, function size limits, and error handling architectures across the codebase, `.lovable/coding-guidelines/coding-guidelines.md`, `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `spec/17-consolidated-guidelines/31-compiled-simple-coding-guidelines.md`.

Autonomously self-loop and read:
- The master cross-language coding guidelines in `spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/01-naming-and-database.md` through `06-advanced-patterns.md`.
- The code style, braces, spacing, and multi-line rules in `spec/02-coding-guidelines/01-cross-language/04-code-style/01-braces-and-nesting.md` through `06-comments-and-documentation.md`.
- The strict function and type size caps (8 lines preferred, 15 lines max) in `spec/02-coding-guidelines/01-cross-language/04-code-style/04-function-and-type-size.md`.
- The boolean principles, prefixing rules (`is`, `has`, `can`, `should`), and guard extraction in `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md` through `05-exemptions-and-api.md`.
- The absolute prohibition against negative booleans and inverted logic in `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`.
- The strict identifier and file naming conventions in `spec/02-coding-guidelines/08-file-folder-naming/01-cross-language.md`.
- The DRY principles and duplication extraction patterns in `spec/02-coding-guidelines/01-cross-language/08-dry-principles.md`.
- The error management architecture and logging diagnostics in `spec/03-error-manage/00-overview.md` and `spec/03-error-manage/02-error-architecture/01-error-handling-reference.md`.
- The language-specific standards in `spec/02-coding-guidelines/` (TypeScript, Go, PHP, Rust, C#, Python, PowerShell).
- The anti-hallucination rules and common AI mistakes in `spec/02-coding-guidelines/06-ai-optimization/01-anti-hallucination-rules.md` and `03-common-ai-mistakes.md`.
- Read `.lovable/plans/index.md` and `.lovable/memory/index.md`.

---

## 1. Master Coding Standards & Audit Checklist (Zero Tolerance)

You MUST audit the codebase against this exact mechanical checklist derived from `.lovable/coding-guidelines/coding-guidelines.md`:

### Naming & Syntax (R1, R2, R8, R20)
- [ ] **R1 - Acronym Casing (PascalCase)**: Acronyms (`Id`, `Json`, `Url`, `Ip`, `Http`, `Api`, `Tls`, `Sql`, `Uuid`) MUST be PascalCase, NEVER all-caps. Write `UserId`, `HttpClient`, `SwapIpWindows`, `ParseJsonUrl`. Do NOT write `UserID`, `HTTPClient`, `swapIPWindows`, `parseJSONURL`.
- [ ] **R2 - Serialization & JSON Keys**: Serialized keys use PascalCase matching the Go/struct field name (`{"Id": "123", "ApiUrl": "...", "IsActive": true}`). Never use camelCase or snake_case for internal JSON schemas.
- [ ] **R8 - No Magic Strings or Numbers**: Every comparison, status check, or option literal MUST be against a named enum symbol or typed constant. Never compare against raw string literals (e.g. `status === "ACTIVE"`).
- [ ] **R20 - Anti-Garbage Naming**: Never generate arbitrary, generic, or sequential names like `temp`, `data`, `obj`, `comp_100.go`, or `TestHandleComp100`. All identifiers must semantically describe domain behavior.

### Boolean Standards & Logic (R3, R10, R14, R16, R17)
- [ ] **R3 - Mandatory Prefixes**: Every boolean variable, parameter, struct field, JSON key, or boolean-returning function MUST begin with `is`, `has`, `can`, `should`, `was`, `will`, `did`, or `must` (e.g., `isEnabled`, `hasAdminRole`, `isReady`).
- [ ] **Positive Framing Only**: Never use negative booleans (e.g., `isNotReady`, `disableCache`, `hasNoAccess` are banned). Invert them to positive (`isReady`, `isCacheEnabled`).
- [ ] **Never Invert Success**: Never use `!response.isSuccess` or `!ok`. Instead, use explicit failure checks (`response.isFail`, `isMissingEntry := !ok; if isMissingEntry == true { ... }`).
- [ ] **R10 - No Boolean Positional Arguments**: Never do `save(true)` or `render(false)`. Use explicit configuration objects or named methods: `save(SaveOptions{ IsForce: true })` or `renderExpanded()`.
- [ ] **R14 - No Inverted Complex Conditions**: Do not use `!` on complex conditions containing AND/OR. Extract them into named boolean variables using positive framing.
- [ ] **R16 - Strict Conditional Joins**: Keep `if` conditions to a maximum of one join (two operands). Extract complex joins into named boolean variables.
- [ ] **R17 - No Mixed Polarity**: Never mix positive and negative conditions in a single conditional join (e.g., `if (a && !b)` is forbidden; extract positive variables).

### Function Structure, Signatures & Spacing (R4, R5, R6, R13, R14, R15, R16, R17, R18, R19)
- [ ] **Function Length Cap**: 8 lines preferred, 15 lines hard cap. Functions exceeding 15 lines MUST be decomposed into single-responsibility helpers using Shell-and-Wire or Table-Driven Dispatch.
- [ ] **No Nested Ifs**: Zero tolerance for nested `if` blocks. Flatten with early returns and guard clauses.
- [ ] **R4 - Signature Line Splitting**: If a function has > 3 parameters OR the signature exceeds 100 characters, split into one parameter per line with trailing commas where allowed.
- [ ] **R5 - Parameter Grouping**: If a function has > 4 parameters OR 2+ adjacent same-typed parameters, group them into a dedicated struct / options object (e.g. `SwapIpParams`).
- [ ] **R6 - Unused Parameters**: Remove dead parameters or explicitly discard with an interface-required trailing comment.
- [ ] **R13 - Blank Line Before Return/Throw**: Exactly one blank line before every `return` or `throw` (unless it is the only statement in the block).
- [ ] **R14 - Blank Line After Closing Brace**: Exactly one blank line after a closing `}`, unless the next line is `}`, `else`, `case`, or `catch`.
- [ ] **R15 - No Double Blank Lines**: Never use two blank lines in a row anywhere.
- [ ] **R16 - No Edge Padding Inside Braces**: No blank line immediately after `{` or immediately before `}`.
- [ ] **R17 - Top-Level Declaration Spacing**: Exactly one blank line between top-level functions, classes, and structs.
- [ ] **R18 - Import Grouping**: Group imports strictly: standard library, third-party, first-party absolute, first-party relative.

### Enums & Type Safety (R19, R21, R22, R24)
- [ ] **R21 - Enum Type Suffix**: Every enum alias MUST end with the suffix `Type` (e.g., `FormatType`, `UserRoleType`, `NodeStateType`).
- [ ] **Dedicated Files**: Types, enums, constants, and interfaces MUST live in their own dedicated files, never inline alongside business logic.
- [ ] **Exhaustive Pattern Matching**: Switches on enums must handle every case or include a safe error-handling default branch.

### Error Management (AppError Architecture)
- [ ] **Zero Swallowed Errors**: Never use `_ = err` or empty `catch {}`. Every error must be logged or returned.
- [ ] **Context Wrapping**: Wrap every error with an operation label and domain parameters (`apperror.Wrap(err, "operationName", context)`). Original stack traces must survive.
- [ ] **Return AppError**: Core functions and CLI command handlers MUST return `*apperror.AppError` and avoid calling raw `os.Exit(1)` inside business logic so that command auditing can log failures.

### Language-Specific Rules
- [ ] **Go**: Return `*apperror.AppError`. Use result types for multi-returns containing booleans (`(T, bool)` forbidden; return a struct). Define enums as typed bytes/ints (`type StatusType byte`) with `iota`.
- [ ] **TypeScript / React**: `useEffect` conditions must extract guards into positively named booleans. State is strictly immutable: no in-place array/object mutations (`push`, `splice`, direct assignments banned; use spread or `structuredClone`). No public tuples.
- [ ] **Python**: Strict type hints on every public signature; use `@dataclass` or `pydantic` for structured data.
- [ ] **Rust**: `Result<T, E>` with `thiserror` enums; default to immutable `let`.
- [ ] **C#**: `I`-prefixed interfaces; PascalCase methods and properties; custom `AppException`.
- [ ] **PHP**: Strict enum comparisons using `->isEqual()`, never `===`.

---

## 2. Planning Loop (Deep N-Step Analysis)

This is not a quick glance. You must deeply read the codebase, looping yourself as much as needed (taking exactly `N = 150` steps of internal planning and reading). Each step MUST be followed using full multi-agent cognitive logic and memory retention.

You must dedicate this immense processing power to uncover:
- Every inverted boolean (`!isSuccess`, `!ok`, `!found`, `!has...`).
- Every magic string, number, or bare status literal.
- Every swallowed error or unhandled return.
- Every missing Enum `Type` suffix.
- Every monolithic function exceeding 15 lines.
- Every nested `if` statement.

If there are NO discrepancies, explicitly state: *"There are no coding guideline issues or discrepancies."* However, assume the codebase is a mess until proven otherwise.

---

## 3. Root Cause & Fallout Analysis

For every issue found:
1. **Root Cause**: Why was it written this way? What historical context allowed the violation?
2. **Blast Radius**: How many places does it touch?
3. **Fallout Check**: If we change this, what else breaks? Will it break the CI/CD pipeline? Will it break tests? Map the entire blast radius.

---

## 4. Enqueueing Tasks for Sub-Agents

Your final output must be a massively detailed plan stored at `.lovable/plans/pending/01-coding-guideline-fixes.md` and granular subtask files written to `.lovable/plans/subtasks/01-coding-guideline-fixes/XX-<subslug>.md`.

The plan must break the work down so granularly (exactly `N = 150` steps) that 3 concurrent sub-agents can be spawned later to safely execute the fixes.

### Sub-Agent Orchestration Requirements:
1. **Specific Titling**: Each sub-agent must be spawned with a highly specific title reflecting its exact task (e.g. `Enum Suffix Refactorer`, `Boolean Logic Normalizer`, `Guard Clause Flattener`).
2. **Micro-Tasking**: Sub-agents must only be assigned simple, small micro-tasks rather than monolithic passes.
3. **Agent Delegation**: Each subtask file MUST explicitly state that it will be executed by a separate standalone agent.
4. **End-of-Loop Commit Fix**: Instruct executing agents to verify that no temporary or binary artifacts are staged, run the local CI runner (`.lovable/ai-fix-scripts/02-cicd-local-runner.py`), commit with clean messages, and push immediately.
5. **No Code Fixes in This Turn**: Your job is ONLY to plan, audit, and enqueue. Do NOT modify application source code during the planning turn.

---

## 5. The 4-Part RCA Requirement (Mandatory Memory File)

Before any execution turn starts, you MUST document the issue in `.lovable/memory/issues/XX-<slug>.md` (where XX is the next available sequential number). The file MUST contain these exact four sections:

1. **Why it happened**: The high-level business, logical, or architectural breakdown of the failure.
2. **How it happened**: The technical execution flow that triggered the bug or pattern accumulation.
3. **Root Cause**: The exact file, line, and dependency responsible for the failure.
4. **Code Fix**: The exact code snippets showing what needed to be changed to fix the root cause.

---

## 6. Verification Checklist Before Finishing

- [ ] Plan written to `.lovable/plans/pending/01-coding-guideline-fixes.md` with EXACTLY `N = 150` steps.
- [ ] Subtask files created in `.lovable/plans/subtasks/01-coding-guideline-fixes/XX-<subslug>.md` (named `XX-<subslug>.md`).
- [ ] 4-part RCA memory issue documented in `.lovable/memory/issues/XX-<slug>.md`.
- [ ] `.lovable/plans/index.md` updated with the pending plan.
- [ ] `.lovable/memory/index.md` updated with the new memory issue.
- [ ] Confirmed root `readme.md` is strictly lowercase and synced with `what-to-read.md`.
- [ ] Zero application source code modified during this turn.
