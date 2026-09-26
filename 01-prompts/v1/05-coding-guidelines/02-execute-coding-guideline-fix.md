# Guideline Fix Execution & Linter Remediation — Workflow (must follow)

> [!IMPORTANT]
> Prompt Version: 2.6.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Apply concrete, surgical fixes for all coding guideline violations listed in the pending tasks. Strictly adhere to all style rules, boolean principles, function size limits (< 8–15 lines), and type-safety standards in bounded 5-8 file micro-batches without introducing regressions.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents, default: 2)
H = 2   (Number of hands / parallel operations per agent, default: 2)
```

N = total self-loop steps budget that the agents will perform.
A = count of autonomous subagents running concurrently (default: 2).
H = number of hands / parallel operations per agent (default: 2).

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec with Violation Ledger in .ai-memory/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Active Code Refactoring, Disk Verification, Linter Verification, Local CI Runner Verification, Plan Completion)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

### Multi-Agent Parallel Task Allocation & Orchestration (A = 2, H = 2)

When multiple autonomous agents are present (A >= 2, H >= 2):
1. **Single-Agent Unified Blueprint Mandate:**
   - The initial plan, violation scoping, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
   - Never allow multiple agents to author disjoint or competing audit plans simultaneously. A single coherent architectural vision must lead.
2. **Most Useful Parallel Tasks (Reading Files & Writing Modular Specs):**
   - Once the unified blueprint is authored by the lead agent, the most effective parallel tasks for subagents (A = 2, H = 2) are:
     - **Reading Files:** Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting AST violations in parallel using `03-ai-scripts/17-fast-file-reader.py` and `03-ai-scripts/11-fast-file-scanner.py`.
     - **Writing Modular Specs & Violation Ledgers:** Authoring modular spec sections and populating the violation ledger in parallel adhering to the lead agent's blueprint.
3. **Spec Writing & Spec Audit Emphasis:**
   - **Spec Writing Section:** The lead agent defines the guideline audit overview and rule boundaries first; subagents are then spawned in parallel to flesh out disjoint spec modules and verification gates concurrently.
   - **Spec Audit Section:** The lead agent establishes the audit methodology and roadmap first; subagents are then dispatched concurrently to inspect disjoint code areas, build the violation ledger, and generate granular subtask files.
4. **Execution Mode (Disjoint Refactoring):**
   - Subagents execute parallel disjoint refactoring tasks across non-overlapping files and run targeted file-level linters (`exit 0`).

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/execute-coding-guideline-fix/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## Non-Negotiable Core Rules

- **DRY:** Eliminate duplicated logic.
- **Typed Enums:** Use typed enums with `Type` suffix instead of magic strings or numbers.
- **Function Size Caps:** Functions must be strictly < 8–15 lines. Extract sub-logic into private helpers.
- **Source File Caps:** Source files should be ≤ 80–100 lines; split into logical modules if larger.
- **No Code Mutation:** Only apply fixes, never introduce regressions or alter external business contracts.
- **Affirmative Positive Booleans (`is` / `has`):** No `isNot`, `hasNo`, `isUndefined`, or negative prefixes. Always use affirmative positive framing: try `isDefined` / `IsDefined` instead of negatives (e.g., use `isDefined` instead of `isUndefined` or `isNotDefined`, and invert with `!isDefined` in guard clauses; use `isValid` instead of `isNotValid`, `hasValue` instead of `hasNoValue`, `isReady` instead of `isNotReady`). No nested if statements, no magic values.
- **IsDefined over `!isEmpty`:** MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Never write `if !isEmpty`. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
- **Vertical Line Gaps:** Mandatory blank lines before `if`, after `}`, before `return`, and around multiline struct calls.
- **Mandatory File Path & Variable Context:** Any error created or returned when operating on files, paths, or variables must embed the target path via `.WithPath(path)` / `WrapWithPath(..., path)` and variable context via `.WithVar(name, value)` (Rule R7).
- **Golang Structured Error & Result Containers:** Return `*appfault.AppError` and replace multi-value error tuples with `appfault.Result[T]`, `appfault.ResultSlice[T]`, `appfault.ResultMap[K, V]`.
- **Mandatory types.go Single Reusable Type Definition:** Centralize all domain payload structs and repeated Result aliases into a dedicated `types.go` file within the package as a single reusable named type. Never leave unexported structs or raw generic Result declarations scattered inline.
- **Pointer Null Safety & Method Composition:** All Result inspection methods must attach to pointer receivers (`*Result[T]`) with line-1 `if r == nil` guards, composing existing methods (`r.IsFailure()`, `r.IsSuccess()`, `r.Count()`) rather than repeating raw pointer/error checks.
- **TOTAL BAN on Test Running & Build Checking during Routine Execution:** Test execution (`go test`, `pytest`, python runners) and build checks (`go build`, compiler verification) are strictly banned. Testing and builds are verified later on in CI/CD.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## Anti-Hallucination & Relative Path Rules

> [!CAUTION]
> **STRICT RELATIVE GIT PATHS ONLY — NO ABSOLUTE PATHS / NO `file:///` URIs:**
>
> When generating plans, subtasks (`.ai-memory/plans/subtasks/`), memory issue logs (`.ai-memory/memory/issues/`), specs, code comments, or citations:
> 1. Strictly Relative to Git Root: All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `02-spec/02-coding-guidelines/readme.md`, `cmd/main.go`).
> 2. Total Ban on Absolute Paths: NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `C:\...`, `/home/...`) or absolute file URIs (`file:///...`) into ANY file.

---

## Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI)

When monitoring or checking remote CI/CD pipelines (e.g., after pushing commits or inspecting CI status):

1. **Mandatory GitMap Pipeline-AI Authority:** Agents MUST use GitMap CLI to inspect pipeline state:
   ```bash
   gitmap pipeline-ai status --json
   # or alias:
   gitmap pl-ai status --json
   ```
   Parse structured output fields: `is_running`, `status`, `etaSeconds`, and `nextAiCommand`.
2. **Anti-Credit-Waste Waiting Mandate (TOTAL BAN on Rapid Polling):**
   - NEVER loop rapidly or busy-poll (`gh run view` in tight loops). This burns user credits, exhausts LLM tokens, and wastes rate limits.
   - When a pipeline is in progress (`is_running: true`), agents MUST wait/sleep based on the estimated completion duration (`etaSeconds` or `-t <sec>`):
     ```bash
     gitmap pipeline-ai status -t <etaSeconds>
     ```
   - Proportional ETA sleep guidelines:
     - `etaSeconds > 120`: wait 20s–30s before querying again.
     - `60 < etaSeconds <= 120`: wait 10s–20s before querying again.
     - `etaSeconds <= 60`: wait 5s–10s before querying again.
3. **Targeted Failure Diagnostics:** Use GitMap's automated error extraction to isolate actionable failure lines (`##[error]`, `FAIL:`, compile errors) without fetching noisy passing step logs.

---

## Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Verify incoming directives and checklists as highest priority.
2. **Verbatim Prompt Capture:** Capture the incoming user request verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under `## User Request (Verbatim)`.
3. **Actionable Deliverables Extraction:** Break down user requirements into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, `Task-03`) under `## Extracted Actionable Task List`.
4. **Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):**
   - You MUST output the confirmed deliverables list directly in chat, but you MUST NOT end your turn or pause after emitting it!
   - In the EXACT SAME RESPONSE turn as your markdown output, you MUST immediately invoke the first tool call (e.g. `write_to_file` to initialize the planning spec or `run_command` for fast discovery).
   - NEVER emit the breakdown text without invoking a tool call. Emitting visible text without a tool call terminates the AI turn, causing the engine to close or hang the conversation while tasks remain pending.
   - TOTAL BAN ON WAITING FOR APPROVAL: Do not pause, do not wait for the user to reply, and do not ask "Should I proceed?". Override any planning mode stop directives. Proceed immediately into Phase 1B in the same turn.
   - Format each task clearly with proper markdown indentation, vertical blank lines, task state (`State: [IN PROGRESS — EXECUTING IMMEDIATELY]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).

Line-by-Line Output Format Structure:
```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to Phase 1B: Spec & Subtask Generation (Active Tool Call Running Below).
```

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

## Phase 1B: Scan, Spec Generation & Lean Subtasks (Steps 1 .. N/2)

### Step 1: Single-Agent Unified Blueprint & Master Spec
- Author `.ai-memory/plans/pending/xx-<slug>-audit.md` with an exhaustive Violation Ledger tracking every violation:
  `| Violation Id | File Path | Line Number | Exact Snippet | Planned Fix | Status (PENDING/DONE) |`
- The lead agent authors the complete overarching plan and violation scoping as a unified blueprint before delegating subtasks.
- Once the blueprint is established, subagents (A = 2, H = 2) may be spawned in parallel to read disjoint files and write modular spec sections.

### Step 2: Parallel Reading & Violation Discovery (A = 2, H = 2)
- Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting AST violations in parallel across disjoint subtrees using `03-ai-scripts/17-fast-file-reader.py` and `03-ai-scripts/11-fast-file-scanner.py`.

### Step 3: Granular Subtask Decomposition
- Break down the violation ledger into micro-batches of strictly 5-8 files in `.ai-memory/plans/subtasks/xx-<slug>/<batch-id>-subtask.md`.
- Subtasks follow a lean template:
  ```markdown
  # Subtask: [Batch Name]
  Traceability ID: Task-01
  Target Files: [5-8 relative paths]
  Violations: [IDs from Violation Ledger]
  Action: [Concrete refactoring steps: implicit booleans, affirmative names, guard clauses, Result wrappers]
  Targeted Verification: python 03-ai-scripts/05-guideline-autofixer.py <files>
  ```

### Step 4: Spec Audit Gate (Zero Gaps Mandate)
- Audit the plan against `02-spec/02-coding-guidelines/` and `02-spec/03-error-manage/`.
- Verify every item in the violation ledger maps 1:1 to an actionable subtask.
- As soon as Step 4 finishes, UNCONDITIONALLY transition into Phase 2 in the same turn without asking user confirmation!

---

## Phase 2: Active Code Refactoring & Micro-Batching (Steps N/2+1 .. N)

1. **Parallel Dispatch (A = 2, H = 2):**
   - Spawn subagents to execute disjoint 5-8 file micro-batches concurrently.
   - Maintain active file paths in `.ai-memory/readme.md`. Parallel subagents must never touch the same files simultaneously.

2. **Execution Steps:**
   - Locate the affected source files for the assigned micro-batch.
   - Refactor duplicated code into shared helpers.
   - Invert guard clauses and flatten nested `if` statements.
   - Enforce positive affirmative booleans (`is` and `has`). Replace `== true` with implicit checks.
   - Record modified files under lock: `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
   - Run targeted file-level linters/autofixers on modified files (`exit 0`). DO NOT run banned full test suites or builds during routine execution turns (deferred to CI/CD).
   - Update the pending task file to mark completed violations.

---

## Phase 3: Task Consolidation & File Reduction (End of Loop)

When all subtasks for a parent task (`.ai-memory/plans/pending/xx-<slug>.md`) are finished, execute this final cleanup step before ending the run:
1. Combine all the completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header that explicitly references how the main task started and documents exactly how many steps/loops it took to complete.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
6. Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## The 4-Part RCA Requirement (Mandatory Memory File)

Before you write any code to fix a bug or unexpected error, you MUST document the issue in `.ai-memory/memory/issues/xx-<slug>.md` (where XX is the next available sequential number). The file MUST contain these exact four sections:

1. **Why it happened:** The high-level business, logical, or architectural breakdown of the failure.
2. **How it happened:** The technical execution flow that triggered the bug.
3. **Root Cause:** The exact file, line, and dependency responsible for the failure.
4. **Code Fix:** The exact code snippets showing what needed to be changed to fix the root cause.

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] **NO RAPID CI/CD POLLING (TOTAL BAN):** NEVER query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents MUST query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait/sleep based on `etaSeconds` to eliminate credit waste.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!
