---
name: execute-coding-guideline-fix
description: Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations across the target codebase in bounded 5-8 file micro-batches.
---

# Execute Coding Guideline Fix

> [!IMPORTANT]
> Skill Version: 2.6.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

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
     - **Reading Files:** Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting AST violations in parallel using GitMap AUM (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary and Python scripts (`03-ai-scripts/17-fast-file-reader.py`, `03-ai-scripts/11-fast-file-scanner.py`) as fallback.
     - **Writing Modular Specs & Violation Ledgers:** Authoring modular spec sections and populating the violation ledger in parallel adhering to the lead agent's blueprint.
3. **Spec Writing & Spec Audit Emphasis:**
   - **Spec Writing Section:** The lead agent defines the guideline audit overview and rule boundaries first; subagents are then spawned in parallel to flesh out disjoint spec modules and verification gates concurrently.
   - **Spec Audit Section:** The lead agent establishes the audit methodology and roadmap first; subagents are then dispatched concurrently to inspect disjoint code areas, build the violation ledger, and generate granular subtask files.
4. **Execution Mode (Disjoint Refactoring):**
   - Subagents execute parallel disjoint refactoring tasks across non-overlapping files and run targeted file-level linters (`exit 0`).

---

## Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable, and strictly incorporated ahead of all other guidelines.
2. **Verbatim Prompt Capture:** Capture the incoming user request verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section.
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
- Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting AST violations in parallel across disjoint subtrees using GitMap AUM (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary and Python scripts (`03-ai-scripts/17-fast-file-reader.py`, `03-ai-scripts/11-fast-file-scanner.py`) as fallback.

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

2. **Core Refactoring Rules:**
   - **No Line Compression:** Maintain mandatory blank lines before `if`, after `}`, before `return`.
   - **Implicit Booleans:** Never write `== true`. Replace with implicit checks.
   - **Affirmative Boolean Parameters & Fields:** Ban single-letter (`v bool`, `b bool`), bare identifiers (`stop bool`, `pause bool`), and negative identifiers (`isNot*`, `isUndefined`, `hasNo*`). Rename to `isStopOnFail bool`, `isStopped bool`, `isPaused bool`, `isDefined bool`. Try `isDefined` / `IsDefined` instead of negatives (e.g., use `isDefined` instead of `isUndefined` or `isNotDefined`, and invert at callsite with `!isDefined` if testing for absence; use `isValid` instead of `isNotValid`, `hasValue` instead of `hasNoValue`).
   - **IsDefined over `!isEmpty`:** MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Never write `if !isEmpty`. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
   - **Guard Clauses:** Invert early checks to return immediately and flatten nested blocks.
   - **Go Errors & Result Containers:** Return `*appfault.AppError` and replace multi-value error tuples with `appfault.Result[T]`, `appfault.ResultSlice[T]`, `appfault.ResultMap[K, V]`.
   - **Mandatory types.go Single Reusable Type Definition:** Centralize all domain payload structs and repeated Result aliases into a dedicated `types.go` file within the package as a single reusable named type. Never leave unexported structs or raw generic Result declarations scattered inline.
   - **Pointer Null Safety & Method Composition:** All Result inspection methods must attach to pointer receivers (`*Result[T]`) with line-1 `if r == nil` guards, composing existing methods (`r.IsFailure()`, `r.IsSuccess()`, `r.Count()`) rather than repeating raw pointer/error checks.
   - **Fluent Predicates:** Enforce `IsCountOtherThan(N)`, `IsEmpty()`, `HasRecord()`, `IsDefined()` at call sites instead of compound checks (`err != nil || len(...) != N`).
   - **Atomic Change Tracking:** Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json` for subsequent CI/CD verification.
   - **Targeted Batch Verification:** Run targeted linter / autofixer on the modified files in the batch (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <files>`). DO NOT run `06-cicd-local-runner.py`, test runners, or full builds (`npm run build`, `go build ./...`) during routine batch fixes.

---

## Phase 3: Task Consolidation & File Reduction (End of Loop)

When all subtasks for a parent task (`.ai-memory/plans/pending/xx-<slug>.md`) are finished, execute this final cleanup step before ending the run:
1. Combine all the completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header that explicitly references how the main task started and documents exactly how many steps/loops it took to complete.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
6. Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] **NO RAPID CI/CD POLLING (TOTAL BAN):** NEVER query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents MUST query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait/sleep based on `etaSeconds` to eliminate credit waste.

---

## Fast File Discovery & Reading Toolchain (GitMap AUM Primary, Python Fallback)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (e.g. `gitmap find "*.go" -ext "go"`, `gitmap find "01*"`)
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<pattern>" --lang go --limit 50`
- **Sub-Millisecond Folder Explorer & Reader:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> --ext .go --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Codebase Topology:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!
