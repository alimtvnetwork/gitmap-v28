# Naming Conventions, Boolean Prefixes & Anti-Ok Variables — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-naming`, `cg-execute naming`, `audit naming`, `fix boolean naming`, `fix naming conventions`, `fix ok boolean`, `affirmative naming`, `positive boolean naming`, `naming conventions audit`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all variable and boolean naming violations across the codebase, directly modifying source files to replace bare `ok` identifiers, eliminate negative boolean variables (`hasNo*`, `isNot*`), enforce affirmative prefixes (is and has only (can, should, was, etc. are banned)), apply positive framing with inverted `if` guard clauses, and normalize acronym casing until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase to inventory all architectural violations and anti-patterns.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md` for implicit positive booleans and anti-negative rules.
15. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/10-function-naming.md` for semantic verb and predicate prefix standards.
16. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md` for domain-specific architectural specifications.
17. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md` for domain-specific architectural specifications.
18. [ ] /learn Ingest `spec/02-coding-guidelines/` for domain-specific architectural specifications.
19. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
20. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase for Naming Violations, Write .lovable/plans/pending/ Spec, Create .lovable/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Rename Bare `ok`, Invert Negative Booleans, Enforce Positive Guards, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Variable & Boolean Naming Architecture (Zero Tolerance)

Naming clarity is the backbone of robust code. Vague identifiers, bare `ok` variables, and negative boolean flags cause cognitive fatigue and obscure critical edge-case bugs.

### 1. Mandatory Boolean Prefix Rule (Strict is/has ONLY)

Every boolean variable, parameter, struct field, or property MUST begin with `is` or `has` ONLY (PascalCase `Is` / `Has` for exported symbols; e.g. `isValid`, `hasPermission`, `isReady`, `hasData`); all other prefixes (`can`, `should`, `was`, `will`, `did`, `must`, etc. like `canExecute`, `shouldRetry`) and negative names are **strictly BANNED**.

- Go: `isValid`, `hasPermission`, `isExecutable`, `isRetryRequired` (or `IsValid`, `HasPermission`).
- TypeScript/JavaScript: `isLoaded`, `hasColors`, `hasPayload`, `isSubmittable`.
- Python: `is_valid`, `has_permission`, `is_ready_to_proceed`.
- C#: `IsValid`, `HasAccess`, `IsExecutable`.

---

### 2. TOTAL BAN on Bare `ok` (The `ok` Anti-Pattern)

In Go type assertions, map lookups, channel receives, and comma-ok idioms, the bare identifier `ok` is **strictly forbidden**. It carries zero domain semantics and violates the mandatory boolean prefix rule.

You MUST replace bare `ok` with a domain-specific boolean starting with `is` or `has`:

| Context | ❌ FORBIDDEN (Bare `ok`) | ✅ REQUIRED (Affirmative Semantic Boolean) |
|---|---|---|
| **Type Assertion** | `appErr, ok := err.(*apperror.AppError)` | `appErr, isAppErr := err.(*apperror.AppError)` |
| **Map Lookup** | `val, ok := userMap[id]` | `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]` |
| **Map Key Check** | `_, ok := headers["Authorization"]` | `_, hasAuthHeader := headers["Authorization"]` |
| **Channel Receive** | `msg, ok := <-msgChan` | `msg, hasMessage := <-msgChan` or `msg, isChannelOpen := <-msgChan` |
| **Type Switch / Cast** | `str, ok := val.(string)` | `str, isString := val.(string)` |
| **Status Tuples** | `data, ok := fetch()` | `data, isSuccess := fetch()` |

---

### 3. TOTAL BAN on Negative Boolean Identifiers (Anti-`hasNo*`, Anti-`isNot*`)

Never name a boolean variable or property with negative prefixes or inverted words:

- ❌ **FORBIDDEN:** `hasNoColors`, `hasNoPayload`, `isNotReady`, `isNotDisabled`, `hasNoAccess`, `isNoOp`, `disallowGuest`, `unauthorized`.
- ✅ **REQUIRED:** `hasColors`, `hasPayload`, `isReady`, `isEnabled`, `hasAccess`, `isOp`, `allowGuest`, `isAuthorized`.

---

### 4. The Clean Solution: Positive Framing with Inverted `if` Guard Clauses

When you need to handle the absence, empty state, or failure condition of a resource, **ALWAYS declare the variable positively** and perform the negative check inside the `if` guard clause:

#### Go Example: Type Assertion & Guard Inversion

```go
// ❌ FORBIDDEN: Nested if with bare ok and else branch
if appErr, ok := err.(*apperror.AppError); ok {
    if appErr.Code != "E_INTERNAL_ERROR" {
        t.Errorf("expected E_INTERNAL_ERROR, got %s", appErr.Code)
    }
} else {
    t.Errorf("expected AppError, got %T", err)
}

// ✅ REQUIRED: Semantic isAppErr boolean + inverted guard clause
appErr, isAppErr := err.(*apperror.AppError)
if !isAppErr {
    t.Fatalf("expected AppError, got %T", err)
}

if appErr.Code != "E_INTERNAL_ERROR" {
    t.Errorf("expected E_INTERNAL_ERROR, got %s", appErr.Code)
}
```

#### TypeScript / React Example: Positive Framing with Inverted Guard

```tsx
// ❌ FORBIDDEN: Negative boolean variables (hasNoColors, hasNoPayload)
const hasNoColors = !colorConfig.length;
if (hasNoColors) {
    return null;
}

const hasNoPayload = !payload?.length;
if (hideLabel || hasNoPayload) {
    return null;
}

// ✅ REQUIRED: Positive boolean variables + inverted guard conditions
const hasColors = colorConfig.length > 0;
if (!hasColors) {
    return null;
}

const hasPayload = Boolean(payload?.length);
if (hideLabel || !hasPayload) {
    return null;
}
```

#### Python Example: Dictionary Lookup & Guard Inversion

```python
# ❌ FORBIDDEN: Negative boolean flag
is_not_authorized = user.role != "admin"
if is_not_authorized:
    raise PermissionDenied()

# ✅ REQUIRED: Positive boolean + inverted condition
is_authorized = user.role == "admin"
if not is_authorized:
    raise PermissionDenied()
```

---

### 5. Pros vs Cons: Why Positive Naming + Inverted Guards Is Superior

| Dimension | ❌ Negative Naming (`hasNoColors = !len`) | ✅ Positive Framing (`hasColors = len > 0; if (!hasColors)`) |
|---|---|---|
| **Cognitive Load** | **High:** Requires mental inversion on every read. | **Low:** Matches natural human language and domain models. |
| **Double Negative Risk** | **Severe:** Leads to monstrosities like `if (!hasNoColors)`. | **Zero:** Negation is always single and explicit: `if (!hasColors)`. |
| **Boolean Algebra** | **Confusing:** Combining `hasNoColors && hasNoPayload` obscures truth tables. | **Intuitive:** De Morgan's laws and logical OR/AND remain obvious. |
| **Single Source of Truth** | **Fragmented:** Some files use `hasColors`, others use `hasNoColors`. | **Standardized:** All components evaluate the presence of state uniformly. |
| **Guard Clause Flow** | **Awkward:** Hides happy-path invariants inside inverted branches. | **Clean:** Early returns eliminate nesting depth to level 0. |

---

### 6. General Variable & File Naming Conventions

1. **Strict Lowercase Filenames:** All files, scripts, documentation, and system files MUST use strictly lowercase naming (e.g., `readme.md`, `01-file-manipulator.py`, `agents.md`, `skill.md`).
2. **Anti-Garbage Variable Naming:** Absolutely NO generic garbage variable names (`comp_100.go`, `temp`, `data`, `obj`, `val1`, `item_01`, `TestHandleComp100`). All names must be semantic and domain-specific.
3. **PascalCase Acronyms:** Acronyms are formatted as regular words: first letter capitalized, remaining letters lowercase (`UserId`, `ApiUrl`, `HttpServer`, `IpAddress`, `JsonData` — NOT `UserID`, `APIURL`, `HTTPServer`).

---

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.lovable/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.lovable/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`, `.lovable/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.lovable/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.lovable/plans/subtasks/` to `.lovable/plans/completed/` and update `.lovable/plans/01-index.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without cancelling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native Naming Autofixer:** If you need to scan and fix boolean naming or bare `ok` identifiers, use `python 03-ai-scripts/05-guideline-autofixer.py <file>` or create and index `03-ai-scripts/08-naming-autofixer.py`.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `05-naming-autofixer.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/01-index.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.lovable/plans/01-index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] **Zero Bare `ok` Identifiers:** Absolutely zero bare `ok` variables in type assertions, map lookups, or status returns. All renamed to `isAppErr`, `isFound`, `hasValue`, etc.
- [ ] **Positive Boolean Framing:** All booleans named positively (`hasColors`, `hasPayload`, `isReady`). Zero `hasNo*` or `isNot*` variables.
- [ ] **Inverted Guard Clauses:** Negative checks handled via inverted guard returns (`if (!hasColors) return null;`).
- [ ] **Boolean Prefixes:** All booleans begin with is or has ONLY (all other prefixes like can, should, was, will, did, must are banned).
- [ ] **PascalCase Acronyms:** All acronyms formatted as `UserId`, `ApiUrl`, `JsonData`.
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.lovable/coding-guidelines.md`.

1. [ ] /learn and apply as a /goal `.lovable/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.

- [ ] `python linter-scripts/check-enum-and-boolean.mjs` or `python linter-scripts/check-boolean-guidelines.py` exited with code 0.
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`, `spec/02-coding-guidelines/01-cross-language/10-function-naming.md`, `spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md`, `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`, and `.lovable/coding-guidelines.md`.
- [ ] Zero Bare `ok`: All type assertions and map lookups use affirmative boolean names (`isAppErr`, `isFound`).
- [ ] Positive Booleans & Inverted Guards: All booleans use affirmative names (`hasColors`, `hasPayload`); guard clauses invert condition (`if (!hasColors)`).

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Script:** `linter-scripts/check-enum-and-boolean.mjs` (or `linter-scripts/validate-guidelines.py`)
2. **Local Run Command:** `python linter-scripts/validate-guidelines.py`
3. **Autofixer Command:** `python 03-ai-scripts/05-guideline-autofixer.py <file>`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate Naming & Boolean Conventions
     run: python linter-scripts/validate-guidelines.py
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "Naming & Boolean Check": [sys.executable, "linter-scripts/validate-guidelines.py"],
   }
   ```

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.
