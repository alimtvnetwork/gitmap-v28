```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /find-100-codebase-improvements <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously audit, discover, catalog, and refactor 100 concrete improvements in the source codebase strictly based on official Coding Guidelines (`02-spec/02-coding-guidelines/`) and Error Management Specifications (`02-spec/03-error-manage/`): FIRST showcase and list out the given task in visible chat during Turn 1 with the 5 Mandatory Pillars and 6-part verification schema, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all 100 codebase improvements are cataloged with relative file paths, line numbers, function names, defect context, coding guideline violations, and concrete before/after code diffs before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 `self` authoring subagents):** Spawn 3 authoring subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool.

---

## Non-Negotiable Coding Standards & Error Management Rules

Every improvement discovered and refactored by this prompt MUST adhere strictly to the following 6 core pillars:

### 1. Error Return Sovereignty (CRITICAL)
- **Zero Void Methods on Fallible Operations:** Fallible methods that execute database updates, file I/O, subprocess execution, or API calls must NEVER be void (`func ...()`). They MUST return an `error`, `*appfault.AppError`, or a typed `Result[T]` wrapper.
- **Zero Error Swallowing:** Functions must NEVER swallow errors via blank assignment (`_ = db.Update(...)` or `_ = exec.Command(...)`) or return `nil` upon failure.
- **Leaf Error Sovereignty:** Leaf functions MUST return errors directly (`return err` or `return apperror.WrapSimple(...)`). Leaf functions MUST NEVER call exit handlers (`cliexit.HandleError`) or `os.Exit`.
- **Outer Caller Handles Exits:** Only the top-level orchestrator or root command dispatcher receives errors and decides how to terminate or report.

### 2. Function Sizing & Decomposition Limits
- **15-Line Hard Limit:** Functions <= 8 lines are preferred; <= 15 lines is the maximum allowed.
- **Single Responsibility Principle:** Decompose multi-step procedures into small, focused, independently testable helper functions.

### 3. Conditional Flattening & Guard Clauses
- **Zero Nested If Statements:** All nested conditionals must be flattened using guard clauses, early returns, and discrete helper functions.
- **Inverted Condition Pattern:** Check failure conditions first and exit early, keeping the primary execution flow at the top indentation level.

### 4. Affirmative Boolean Naming
- **Strict `is*` and `has*` Prefixes Only:** All boolean variables, fields, and function names MUST begin with `is` or `has` (`can`, `should`, `was`, `will` are banned).
- **Zero Negative Booleans:** Never name booleans negatively (e.g. `isNotReady`, `hasNoErrors` are banned; use `!isReady`, `!hasErrors`).

### 5. Application-Wide DB Repository Pattern
- **Zero Direct SQL in Business Logic:** Business logic commands (`cli/cmd*`) MUST NEVER execute direct SQL queries (`db.Conn().Exec(...)`, `SELECT ...`, `UPDATE ...`).
- **Encapsulated Store Methods:** All database mutations and queries must be encapsulated into typed repository methods on `*store.DB` (e.g. `db.UpdateRepoTransport(...)`).

### 6. Strongly-Typed Enums & Constant Centralization
- **No Magic Strings or Integers:** All string and integer constants must be centralized in dedicated constant packages (`cli/constants/`) and use strongly-typed enums ending in `*Type` (e.g. `GitProviderType`, `ExitCodeType`).

---

## The 100 Codebase Improvements Architecture: 10 Focus Categories

The AI must partition its discovery across exactly 10 categories, uncovering 10 distinct, grounded improvements per category (10 × 10 = 100 items):

### Category 1: Error Return Sovereignty & Swallowed Error Elimination (Items 01–10)
- Elimination of void functions performing mutations (`updateRepoTransportInDB`, `deleteSiteFromDB`, `ensureWorkDirsTableExists`).
- Elimination of blank error discards (`_ = ...`) in business logic commands.
- Wrapping SQLite and system errors with structured `apperror.WrapSimple`.
- Ensuring leaf functions return errors instead of calling exit handlers.

### Category 2: Function Length Reduction & Decomposition (Items 11–20)
- Decomposing oversized functions (> 15 lines) into discrete <= 8-line helper functions.
- Extracting complex parameter construction into dedicated helper functions.
- Isolating loop bodies into standalone single-responsibility functions.

### Category 3: Guard Clauses & Nested Condition Flattening (Items 21–30)
- Refactoring nested `if` blocks into linear guard clauses with early returns.
- Eliminating chained `else if` ladders in favor of map dispatch tables or early exits.
- De-nesting nil-checks and error checks.

### Category 4: Boolean Logic Modernization & Prefix Enforcement (Items 31–40)
- Renaming non-conforming booleans (`canExecute`, `shouldSkip`, `done`) to `is*` or `has*`.
- Inverting negative booleans (`isNotValid` -> `!isValid`).
- Replacing multi-condition boolean expressions with named affirmative helper functions.

### Category 5: Split-DB Repository Pattern Encapsulation (Items 41–50)
- Eliminating raw SQL strings (`SELECT ...`, `UPDATE ...`, `INSERT ...`) from `cli/cmd*`.
- Implementing encapsulated methods on `*store.DB` in `cli/store/`.
- Fixing schema drift defects (e.g. references to non-existent columns in SQLite queries).
- Enforcing parameterization across all database queries.

### Category 6: Strongly-Typed Enums & Constant Centralization (Items 51–60)
- Converting raw string literals (`"vscode"`, `"github"`, `"dpapi"`) into typed enums ending in `Type`.
- Centralizing command names, flag names, and file paths into `cli/constants/`.
- Eliminating magic exit integers (`1`, `2`) in favor of typed `constants.ExitCodeType`.

### Category 7: Subprocess Centralization & Safe Git Execution (Items 61–70)
- Eliminating raw ad-hoc `exec.Command("git", ...)` invocations in favor of `gitutil.RunGitInDir`.
- Guaranteeing safe environment variable injection via `gitutil.BuildSafeGitEnv()`.
- Capturing and parsing standard error streams consistently across all subprocesses.

### Category 8: Golang Pointer Reduction & Value Semantics (Items 71–80)
- Replacing unnecessary pointer parameters (`*string`, `*int`) with value semantics for small structs.
- Eliminating heap escape allocations in high-frequency loops.
- Standardizing slice parameters to variadic `...T` where appropriate.

### Category 9: String Operations & Zero-Allocation Comparisons (Items 81–90)
- Replacing repetitive `strings.ToLower(strings.TrimSpace(a)) == strings.ToLower(...)` with `strutil.EqualFoldAnyTrim`.
- Hoisting string conversions and allocations out of inner loops.
- Using `strings.Builder` and pre-allocated byte buffers for string concatenation.

### Category 10: Database Concurrency, Connection Pools & WAL Governance (Items 91–100)
- Ensuring `PRAGMA journal_mode=WAL;` and `PRAGMA busy_timeout=5000;` on all SQLite opens.
- Guaranteeing `defer db.Close()` is called immediately after successful database opens.
- Capping connection pool handles to prevent `too many open files` errors.

---

## Mandatory Reporting Schema for Each Codebase Improvement

Every single improvement produced by this prompt MUST adhere strictly to the following 7-field schema:

```markdown
### Item [NNN]: [Refactoring Target Title]
- **File Path & Line Range:** [Relative path and exact line range, e.g. `cli/cmdpull/pull_remediation_ssh.go:105-113`]
- **Function / Symbol Name:** [Exact function signature, e.g. `func updateRepoTransportInDB(repoPath, sshURL string)`]
- **Existing Guideline Violation:** [Exact rule cited: Error Return Sovereignty / Function Length / Nested If / Boolean Naming / Direct SQL / Magic String]
- **Context & Architectural Rationale:** [Detailed explanation of why this code pattern is dangerous, causes unhandled states, breaks audits, or violates repo standards]
- **What to Improve & How to Improve:** [Precise technical refactoring blueprint detailing how to rewrite the function, decompose it, and update its callers]
- **Concrete Code Refactoring (Before vs. After):**
  ```go
  // BEFORE (Violating code):
  func updateRepoTransportInDB(repoPath, sshURL string) {
      db, err := store.OpenDefault()
      if err != nil {
          return
      }
      defer db.Close()
      _ = db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
  }

  // AFTER (Clean, compliant code):
  func updateRepoTransportInDB(repoPath, sshURL string) error {
      db, err := store.OpenDefault()
      if err != nil {
          return err
      }
      defer db.Close()

      return db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
  }
  ```
- **Verification Criterion:** [Specific unit test or linter command, e.g. `python 03-ai-scripts/50-fastgate.py --files <path>` and `go test -C cli -v ./<pkg>`]
```

---

## Concrete Example Item

```markdown
### Item 001: Error Return Sovereignty in `updateRepoTransportInDB`
- **File Path & Line Range:** `cli/cmdpull/pull_remediation_ssh.go:105-113`
- **Function / Symbol Name:** `func updateRepoTransportInDB(repoPath, sshURL string)`
- **Existing Guideline Violation:** Rule 1 (Error Return Sovereignty) — Function is void and swallows `store.OpenDefault()` error via silent return, and swallows `db.UpdateRepoTransport(...)` error via blank assignment `_ = ...`.
- **Context & Architectural Rationale:** When a repository's transport is converted from HTTPS to SSH during automated remediation, failures to persist this change into `gitmap.db` are completely lost. As a result, subsequent pulls will continue failing on credentials because the database remains out of sync with the working tree.
- **What to Improve & How to Improve:**
  1. Change signature to `func updateRepoTransportInDB(repoPath, sshURL string) error`.
  2. Return `err` immediately if `store.OpenDefault()` fails.
  3. Return the result of `db.UpdateRepoTransport(...)` directly.
  4. Update upstream callers in `pull_remediation.go` to handle the error and print a diagnostic notice instead of assuming success.
- **Concrete Code Refactoring (Before vs. After):**
  ```go
  // BEFORE:
  func updateRepoTransportInDB(repoPath, sshURL string) {
      db, err := store.OpenDefault()
      if err != nil {
          return
      }
      defer db.Close()
      _ = db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
  }

  // AFTER:
  func updateRepoTransportInDB(repoPath, sshURL string) error {
      db, err := store.OpenDefault()
      if err != nil {
          return err
      }
      defer db.Close()

      return db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
  }
  ```
- **Verification Criterion:** Run `python 03-ai-scripts/50-fastgate.py --files cli/cmdpull/pull_remediation_ssh.go` and `go test -C cli -v ./cmdpull`.
```

---

## Execution Rules & Guidelines

1. **Strict 100-Item Completeness:** The output MUST contain exactly 100 numbered improvements across the 10 categories.
2. **Real File Paths & Lines:** Every single item must point to an actual, existing source file in the repository with real line numbers.
3. **Surgical Refactoring:** Propose targeted, minimal-diff fixes that preserve business logic while enforcing coding standards.
4. **FastGate Pre-Commit Validation:** Every refactored file must pass `python 03-ai-scripts/50-fastgate.py` with zero violations.
5. **Atomic Commit Format:** On completion, commit atomically using `gitmap cpf "<module> - <summary>"`.
