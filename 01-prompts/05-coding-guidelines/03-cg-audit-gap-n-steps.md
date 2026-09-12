# Gap Analysis & N-Step Guideline Audit — Planning Spec (must follow)

Trigger Keywords & Aliases: `cg-audit`, `audit gap`, `cg audit gap`, `cg-audit-gap`, `audit coding guidelines`, `verify coding guidelines`, `cg-audit-gap-n-steps`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 400
```

N = total self-loop steps budget across the codebase. The user may override this number when triggering the prompt (e.g., N = 100 or N = 200).

- [ ] /goal First N/2 steps (Phase 1): Deeply scan the entire codebase file-by-file, dividing N steps across files with 30-50 nested atomic checks per file, scoring guideline compliance from 0 to 100, and writing the master audit report to `spec/01-spec-authoring-guide/01-index.md`.
- [ ] /goal Second N/2 steps (Phase 2): Enqueue all identified gaps into `.lovable/plans/pending/XX-coding-guidelines-audit.md`, break them down into microscopic atomic subtasks inside `.lovable/plans/subtasks/XX-coding-guidelines/`, and register them in `.lovable/plans/01-index.md`.
- [ ] /learn Ingest `.lovable/coding-guidelines.md` (HIGH PRIORITY FIRST), `spec/02-coding-guidelines/`, `spec/03-error-manage/`, `spec/17-consolidated-guidelines/`, and `.lovable/strictly-avoid.md` before taking action.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Deep File-by-File Gap Audit & 0-100 Scoring Report)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Plan Generation & Atomic Subtask Enqueuing)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before any execution, check if this prompt is installed as a native Antigravity Skill.

1. Check if `.agents/skills/coding-guidelines/skill.md` exists in the workspace.
2. If it does NOT exist, create it now. Write the core instructions to `.agents/skills/coding-guidelines/skill.md` with frontmatter:
   ```yaml
   ---
   name: coding-guidelines
   description: >-
     Use this skill to audit, review, and enforce coding guidelines across all languages (Go, TS, Python, PHP, C#).
   ---
   ```
3. Once installed, load it on-demand via progressive disclosure for all future runs.

---

## 1. Non-Negotiable Master Guidelines Checklist (Read & Enforce High Priority)

You MUST verify and audit every item on this checklist across every file and function in the repository:

### Tier 1: Master Consolidated Guidelines (Highest Priority)

- [ ] **Master Consolidated File (`.lovable/coding-guidelines.md`):** Read and enforce all 29 cross-language chapters and language-specific sections.
- [ ] **Consolidated Review Spec (`spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md`):** Cross-verify against root spec truth.
- [ ] **Anti-Hallucination & AI Optimization (`spec/02-coding-guidelines/01-cross-language/01-index.md`):** Rule AH-N1 (Abbreviation casing: `Id`, `Url`, `Api`), AH-O1 (Zero placeholder/truncation stubs), AH-E1 (Implicit booleans).
- [ ] **Citation Requirement (`spec/02-coding-guidelines/01-cross-language/01-index.md`):** Every audit finding MUST cite the exact rule code and spec file path.

### Tier 2: Sizing, Nesting & Code Hygiene

- [ ] **Canonical Size Tiers (`spec/02-coding-guidelines/02-canonical-size-tier.md`):**
  - **Functions:** Ideal 8 lines or fewer, hard cap 15 coding lines max.
  - **Source Files:** Recommended 80 lines or fewer, standard cap 100 coding lines, absolute limit 200–300 lines max.
  - **React Components (`.tsx`):** Hard cap 100 lines max per component file.
  - **Class / Struct:** Hard cap 120 lines max.
  - **Anti-Line Compression Cheating:** STRICTLY BAN collapsing whitespace, removing indentation, merging if/else onto single lines, or stripping formatting to cheat line limits.
- [ ] **Braces & Nesting (`spec/02-coding-guidelines/01-cross-language/01-index.md`):** Zero nested if statements. Invert conditions into guard clauses and early returns. Maximum cyclomatic complexity 5 or less.
- [ ] **Return New Line Standards (R13–R16):** Exactly one blank line before every `return`/`throw`/`raise` (unless sole statement in block). Exactly one blank line after closing `}`. Never two consecutive blank lines.
- [ ] **Function Signatures (R4, R5, R9):** Functions with more than 3 parameters or signatures over 100 characters MUST be formatted with exactly one parameter per line.

### Tier 3: Boolean Principles & Logic

- [ ] **Implicit Positive Booleans (`spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`):** NEVER evaluate `== true` or `== false`. Implicit checks only (`if isReady`).
- [ ] **No Mixed Polarity:** NEVER combine positive and negative checks in the same condition (e.g., `if isA && !isB` is FORBIDDEN; extract to named boolean `isAWithoutB`).
- [ ] **No Inverted Success Checks (`spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`):** Ban `!response.isSuccess` or `!isFound`. Always use positive/inverse naming: `response.isFail`, `isMissing`.
- [ ] **Boolean Prefixes:** All booleans MUST start with is and has only (can, should, was, etc. are banned), `was`, `will`, `did`, or `must`.

### Tier 4: Error Management & Architecture

- [ ] **Error Management Architecture (`spec/03-error-manage/`):**
  - Zero swallowed errors (`catch (err) {}` or `_ = err` is a critical violation).
  - Universal `AppError` wrapping with operation context and key inputs.
  - Universal response envelopes with typed status codes.
  - Zero dual-handling (NEVER panic/log AND return error in the same branch).
  - Typed exit enums for error categories.

### Tier 5: Constants, Enums & Schema

- [ ] **Centralized Enums & Constants (`spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`):**
  - Zero magic strings or raw numeric literals.
  - All Enum names MUST end with the suffix `Type` (e.g., `UserRoleType`, `PaymentStatusType`).
  - TypeScript string unions (`type Role = "admin" | "user"`) are banned.
- [ ] **Database & Naming Conventions (`spec/04-database-conventions/`):** PascalCase table names, camelCase column names, `{TableName}Id` primary/foreign keys.

---

## 2. 0 to 100 Scoring & Gap Deduction System

The audit computes a mathematical score from **0 to 100** for every audited file, subdirectory, and the overall repository.

```text
Score = 100 - (Critical_Count * 10) - (Major_Count * 5) - (Minor_Count * 2)
Score is bounded between 0 and 100: final_score = max(0, min(100, Score))
```

### Violation Severity Matrix

| Severity | Deduction | Violation Types |
| :--- | :--- | :--- |
| **Critical** | **-10 pts** | Swallowed errors (`catch {}`, `_ = err`), panic + return dual-handling, explicit `== true` checks, files > 300 lines, disabled CI/CD or CLI linter commands. |
| **Major** | **-5 pts** | Monolithic functions > 15 lines, nested `if` statements, inverted booleans (`!isSuccess`), magic strings/numbers without enums, missing `*Type` enum suffix, line compression cheating. |
| **Minor** | **-2 pts** | Missing return new lines (R13), uppercase acronyms (`ID`/`URL` instead of `Id`/`Url`), missing parameter line breaks (> 3 params / > 100 chars), missing docstrings or semantic naming gaps. |

### Score Tier Ratings

- **95 – 100:** 🟢 **Exemplary (Production Ready)** — Zero architectural flaws, strict compliance.
- **80 – 94:** 🟡 **Acceptable (Minor Gaps)** — Minor style or line count adjustments required.
- **60 – 79:** 🟠 **Needs Refactoring** — Function length, nested ifs, or boolean inversions present.
- **0 – 59:** 🔴 **Critical Failure** — Swallowed errors, huge monolithic files, or anti-patterns detected.

---

## 3. Multi-Agent Parallelization & N-Step Division

To audit massive codebases quickly without context fatigue or hallucination:

1. **Step Allocation:**
   - The total budget `N` is divided among all active source files in the repository.
   - For each file, the agent dedicates 1 full self-loop step.
   - Inside that single file step, the agent executes **30 to 50 nested atomic checks** covering all checklist rules.
2. **Parallel Agent Execution (2 to 3 Threads Max):**
   - Spawn **2 to 3 parallel sub-agents** (never more than 3) to process distinct directory trees concurrently.
   - Example Thread Allocation:
     - **Sub-Agent A:** Audits backend/Go/PHP/Python files (`internal/`, `pkg/`, `api/`, `src/backend/`).
     - **Sub-Agent B:** Audits frontend/React/TypeScript files (`src/components/`, `src/hooks/`, `src/pages/`).
     - **Sub-Agent C:** Audits shared libraries, scripts, database schemas, and CI/CD pipelines (`scripts/`, `spec/`, `.github/`).
3. **Context Diet Protocol:**
   - Give each sub-agent a strictly bounded file list. Do NOT paste entire file contents or master specs into sub-agent prompts.
   - The sub-agent reads target files locally, calculates score deductions, and returns a structured JSON/Markdown violation table.

---

## 4. Phase 1: Deep Gap Audit & Report Generation (Steps 1 to PHASE_1_STEPS)

> [!IMPORTANT]
> **This phase is purely for analysis and auditing. DO NOT modify or fix application code in this prompt.**
> Your sole output is the comprehensive, rigorous audit report and structured plan tasks.

### Step 1: Discover & Map All Source Files

1. Execute `git ls-files` or recursive directory listing, excluding `node_modules`, `vendor`, `.git`, and build output.
2. Group files by module and directory.

### Step 2: Nested Sub-Step Audit per File (30–50 Checks per File)

For every file, perform the exhaustive check sequence:

1. File line count check (flag if > 100 lines, critical if > 300 lines). Check for line-compression cheating.
2. Function line count check (flag every function > 15 lines).
3. Braces and nested `if` check (flag every nested `if` or cyclomatic complexity > 5).
4. Boolean syntax check (search for `== true`, `== false`, `!isSuccess`, `isNot*`, `!is* && *`).
5. Error handling check (search for empty catch blocks, unchecked error returns, dual handling).
6. Constant & Enum check (search for raw magic strings, inline union types, missing `*Type` suffix).
7. Naming & casing check (search for uppercase acronyms `ID`/`URL`, generic names `temp`/`data`/`obj`).
8. Return new line check (verify blank line before return/throw and after closing braces).

### Step 3: Write Master Audit Report

Save the final comprehensive report to:
`spec/01-spec-authoring-guide/01-index.md`

#### Report Structure Template:

```markdown
# Master Coding Guideline Gap Audit Report

- **Audit Date:** YYYY-MM-DD
- **Total Files Audited:** X
- **Overall Codebase Compliance Score:** YY / 100 (Tier Rating)
- **Total Violations Found:** Critical: A | Major: B | Minor: C

## Executive Summary & Score Breakdown

| Module / Directory | Files | Critical (-10) | Major (-5) | Minor (-2) | Module Score | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `src/backend/` | 12 | 1 | 4 | 8 | 54 / 100 | 🔴 Critical |
| `src/frontend/` | 18 | 0 | 3 | 5 | 75 / 100 | 🟠 Needs Work |
| `src/shared/` | 6 | 0 | 0 | 2 | 96 / 100 | 🟢 Exemplary |

---

## Detailed Violation Ledger (Drop by Drop)

| Id | File Path | Line | Function / Component | Rule Code | Exact Snippet | Severity | Planned Remediation |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| V-001 | `src/auth/service.go` | 45 | `ValidateToken` | `CG-BOOL-001` | `if isValid == true` | Critical | Replace with implicit `if isValid` |
| V-002 | `src/auth/service.go` | 78 | `ProcessLogin` | `CG-SIZE-001` | 32 lines body | Major | Split into `verifyCredentials` & `issueSessionToken` |
| V-003 | `src/api/handler.ts` | 112 | `HandleOrder` | `CG-ERR-002` | `catch (e) { return null; }` | Critical | Wrap error in `AppError` and log context |

---

## Remediation Roadmap & Priority Matrix

1. **Sprint 1 (Critical Zero-Tolerance):** Eliminate all swallowed errors, explicit true checks, and dual-handling.
2. **Sprint 2 (Major Architecture):** Refactor monolithic functions > 15 lines and flatten nested `if` statements.
3. **Sprint 3 (Cleanliness & Hygiene):** Enforce `*Type` enums, PascalCase abbreviations (`Id`, `Url`), and return new lines.
```

---

## 5. Phase 2: Plan Task Enqueuing & Atomic Subtasks (Steps PHASE_1_STEPS+1 to N)

1. **Master Plan Enqueuing:**
   - Create `.lovable/plans/pending/XX-coding-guidelines-audit.md` (next sequential number).
   - Document all findings, total score, and the structured execution plan.
   - Update `.lovable/plans/01-index.md` with the new entry.

2. **Atomic Subtask Creation (`.lovable/plans/subtasks/XX-coding-guidelines/`):**
   - For every offending file or closely coupled group of violations, create a dedicated subtask:
     `.lovable/plans/subtasks/XX-coding-guidelines/01-<module>-<slug>.md`
   - Each subtask file must be **atomic and microscopic** (bounded strictly to 1 file or 1 function).
   - Subtask Template:
     ```markdown
     # Subtask XX.01: Refactor <Filename> Coding Guidelines

     ## Target
     - File: `<filepath>`
     - Function(s): `<function-names>`
     - Current File Score: YY / 100

     ## Violations to Fix
     - [ ] [Rule Code] Line XX: `<snippet>` -> `<fix-description>`
     - [ ] [Rule Code] Line YY: `<snippet>` -> `<fix-description>`

     ## Acceptance Criteria
     - [ ] Function lengths 15 lines max.
     - [ ] Implicit booleans only (no `== true`).
     - [ ] Zero swallowed errors; wrapped with `AppError`.
     - [ ] Return new lines verified (R13–R16).
     - [ ] All unit tests pass.
     ```

---

## STRICT AVOIDANCE: Never Disable CLI Linting or CI/CD Checks

> [!CAUTION]
> **TOTAL BAN ON DISABLING, SKIPPING, OR BYPASSING CLI LINTERS AND CI/CD GATES:**
>
> - **NEVER** disable, comment out, delete, or skip any CLI linting command (`golangci-lint`, `eslint`, `markdownlint`, `tsc`, `pytest`, `phpstan`, `mypy`, `check-*.py`), build step, or test suite.
> - **NEVER** add `|| true`, `continue-on-error: true`, `# nolint`, `// eslint-disable`, or ignore flags to fake a pipeline pass.
> - Disabling or bypassing any CI/CD or CLI lint check is an automatic and immediate rejection.

---

## No Automatic Releases (Strict Policy)

> [!CAUTION]
> This is an audit and verification workflow. You MUST NOT bump versions, update changelogs, or cut a release. All commits use `docs(audit): generate coding guideline gap report`.

---

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently (up to 2-3 threads). Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only audit directory X"). Never give a sub-agent a generic multi-directory task.

---

## Pre-Reply / Loop Checklist (Must Verify Every Turn)

- [ ] **Echo Back the Spec:** I have verified the exact acceptance criteria and rules.
- [ ] **Master Guidelines Read:** I have consulted `.lovable/coding-guidelines.md` at high priority.
- [ ] **Exhaustive Violation Ledger:** Maintained the exact table `| Id | File Path | Line | Function / Component | Rule Code | Exact Snippet | Severity | Planned Remediation |`.
- [ ] **0-100 Score Calculated:** Mathematically computed the score for every file and module.
- [ ] **Audit Report Saved:** Report written to `spec/01-spec-authoring-guide/01-index.md`.
- [ ] **Plans & Subtasks Enqueued:** Master plan written to `.lovable/plans/pending/` and atomic subtasks created in `.lovable/plans/subtasks/`.
- [ ] **Strict Lowercase Filenames:** All generated files use strictly lowercase naming.
- [ ] **No Code Modification in Audit Phase:** Ensured application source code was not modified during the audit.
