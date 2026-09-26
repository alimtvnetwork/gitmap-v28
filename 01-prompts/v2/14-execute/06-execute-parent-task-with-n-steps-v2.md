# [V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)

> [!IMPORTANT]
> Prompt Version: 2.6.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

- **N = 200:** Total self-loop steps budget that the agents will perform.
- **A = 2:** Count of autonomous subagents running concurrently (`invoke_subagent` launches up to 2 subagents).
- **H = 2 (Hands / Parallel Operations):** Concrete operational capacity defined across two execution axes:
  1. **Workload Hands ($H_{batch} = 2$):** Each subagent is assigned a bounded batch of up to 2 tightly coupled subtasks from `.ai-memory/plans/subtasks/` (e.g., Subagent 1 executes Subtasks 01 & 02; Subagent 2 executes Subtasks 03 & 04). This prevents agent sprawl while doubling throughput.
  2. **Tool-Dispatch Hands ($H_{tool} = 2$):** Within any execution step, each agent or subagent executes up to 2 parallel tool calls in a single response turn (e.g. calling two `view_file` or two `run_command` operations simultaneously).

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

### Multi-Agent Parallel Task Allocation & Orchestration (A = 2, H = 2)

When multiple autonomous agents are present (A >= 2, H >= 2):
1. **Single-Agent Unified Blueprint Mandate:**
   - The initial execution plan, architectural decomposition, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
   - Never allow multiple agents to author disjoint or competing initial plans simultaneously. A single coherent architectural vision must lead.
2. **Most Useful Parallel Tasks (Reading Files & Writing Modular Specs):**
   - Once the unified blueprint is authored by the lead agent, the most effective parallel tasks for subagents (A = 2, H = 2) are:
     - **Reading Files:** Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting types in parallel using GitMap AUM (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary and Python scripts (`03-ai-scripts/17-fast-file-reader.py`, `03-ai-scripts/11-fast-file-scanner.py`) as fallback.
     - **Writing Modular Specs:** Authoring modular, segmented spec sub-files in parallel (e.g. `02-data-contracts.md`, `03-visual-and-ux.md`, `04-verification-gates.md`) adhering to the lead agent's blueprint.
3. **Spec Writing & Spec Audit Emphasis:**
   - **Spec Writing Section:** The lead agent defines the parent spec overview and data boundaries first; subagents are then spawned in parallel to flesh out disjoint spec modules and verification gates concurrently.
   - **Spec Audit Section:** The lead agent establishes the audit methodology and roadmap first; subagents are then dispatched concurrently to inspect disjoint code areas, build the violation ledger, and generate granular subtask files.
4. **Execution Mode (Disjoint Refactoring via Subagent Batches):**
   - Subagents execute parallel disjoint refactoring tasks (each assigned to distinct, non-overlapping target files across its $H=2$ subtasks) and run targeted file-level linters (`exit 0`).

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. Top-Instruction Priority Verification: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable, and strictly incorporated ahead of all other guidelines.
2. Verbatim Prompt Capture: Capture the incoming user request verbatim. It will be recorded losslessly in both the canonical spec under `02-spec/21-app/` and the master execution plan under `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section.
3. Screenshot & Print Screen Base64 Image Ingestion Protocol: If the user request or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Convert that base64 encoding or downloaded image to the file system immediately, saving it as a persistent file under `assets/screenshots/<task-slug>-<NN>.png` or `assets/ui/<task-slug>-<NN>.png`.
   - Never leave raw, massive base64 strings or external ephemeral URLs in the prompt text, specs, or subtasks.
   - Refer back to this saved image file in the master spec, planning document, subtasks, and UI implementation notes strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<task-slug>-<NN>.png)`).
4. Actionable Deliverables Extraction: Break down whatever user requirements were given (regardless of how they were formatted) into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, `Task-03`) under `## Extracted Actionable Task List`.
5. Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):
   - You MUST output the confirmed deliverables list directly in chat, but you MUST NOT end your turn or pause after emitting it!
   - In the EXACT SAME RESPONSE turn as your markdown output, you MUST immediately invoke the first tool call (e.g. `write_to_file` to initialize the planning spec or `run_command` for fast discovery).
   - NEVER emit the breakdown text without invoking a tool call. Emitting visible text without a tool call terminates the AI turn, causing the engine to close or hang the conversation while tasks remain pending.
   - TOTAL BAN ON WAITING FOR APPROVAL: Do not pause, do not wait for the user to reply, and do not ask "Should I proceed?". Override any planning mode stop directives. Proceed immediately into Phase 1B in the same turn.
   - Respect whatever requirements the user has given, parse every request completely, and format each task clearly with proper markdown indentation, vertical blank lines, task state (`State: [IN PROGRESS — EXECUTING IMMEDIATELY]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).
   - TOTAL BAN ON UNFORMATTED RUN-ON TEXT: Never concatenate tasks into a single unformatted line or paragraph block (e.g. NEVER `#1. Task-01: ... #2. Task-02: ...`). Every task must be its own clearly separated markdown item.
   - Line-by-Line Output Format Structure:
     - Line 1: Header `### 📋 Confirmed Task Breakdown & Requirement Ingestion`
     - Line 2: Empty blank line
     - Line 3: Numbered task title `1. **Task-01: [Descriptive Task Title]**`
     - Line 4: Indented state bullet (3 spaces) `   - **State:** [IN PROGRESS — EXECUTING IMMEDIATELY]`
     - Line 5: Indented understanding check (3 spaces) `   - **Understood:** [YES] — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]`
     - Line 6: Indented actionable scope bullet (3 spaces) `   - **Actionable Scope:** [Precise technical deliverable and implementation scope]`
     - Line 7: Indented target files bullet (3 spaces) `   - **Target Files / Area:** [relative/path/or/module]`
     - Line 8: Empty blank line (vertical gap before next task)
     - Concluding Line: `Proceeding directly to Phase 1B: Spec & Subtask Generation (Active Tool Call Running Below).`

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

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call (e.g. `write_to_file` to save the spec or `run_command` to discover files) MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone before execution has begun.
(Clarification: This same-turn chaining rule prevents premature turn closure during planning and synchronous tool loops. It does not apply when awaiting asynchronous background subagents in Phase 2, where yielding control is required for reactive wakeup).

---

### Phase 1B: Planning Mode, Detailed Spec Generation & Lean Subtasks (Steps 1 .. N/2)

#### Step 1: Canonical Application Spec Generation (Folder 21 Standard)
First, write the canonical application specification into `02-spec/21-app/` before creating execution plans or modifying code:
- **Single-Agent Unified Blueprint Mandate (MANDATORY):**
  - The initial plan, architecture overview (`01-overview.md`), data boundaries, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating to subagents. Never allow multiple agents to author disjoint or competing initial plans simultaneously.
  - Once the unified blueprint is established by the single lead agent, subagents (A = 2, H = 2) are spawned to handle parallel modular spec writing (`02-data-contracts.md`, `03-visual-and-ux.md`, `04-verification-gates.md`) adhering strictly to the blueprint.
- **Location & Sizing Standard:**
  - Concise / single-domain specs (<= 150 lines): Write to `02-spec/21-app/xx-<slug>.md`.
  - Large / multi-domain features: Write to a segmented directory `02-spec/21-app/xx-<slug>/` with sequential sub-files:
    - `01-overview.md` (Domain architecture, system context, verbatim user request — authored by lead agent)
    - `02-data-contracts.md` (Models, schemas, interfaces, error types — parallel subagent task)
    - `03-visual-and-ux.md` (Component hierarchy, visual layout, screenshots — parallel subagent task)
    - `04-verification-gates.md` (Quality gates, test invariants, acceptance criteria — parallel subagent task)
- **Lossless Verbatim Capture:** Under `## User Request (Verbatim)`, preserve the exact prompt text and constraints without truncation.
- **Visual Assets & Base64 Screenshots:** If screenshot URLs or base64 images were provided, verify they were decoded and saved to `assets/screenshots/<task-slug>-<NN>.png` and reference them strictly via relative markdown links.
- **Spec Registry Registration:** Register the new spec entry in `02-spec/21-app/readme.md` with status `draft` or `active`.

#### Step 2: Scan & Discover (GitMap AUM Acceleration & Multi-Agent Parallel Reading)
To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, leverage the 2-tier discovery toolchain:

#### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (e.g. `gitmap find "*.go" -ext "go"`, `gitmap find "01*"`)
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

#### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50`
- **Sub-Millisecond Folder Exploration:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Subsystem & Topology Overview:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

**Multi-Agent Parallel Discovery (A = 2, H = 2):** When multiple agents are present, the most useful parallel tasks are reading files and authoring modular specs. Subagents concurrently read disjoint codebase areas, explore dependencies, and trace call sites without merge conflicts.

#### Step 3: Actionable Execution Plan & Lean Subtask Decomposition
With the canonical spec established in `02-spec/21-app/`, initialize the execution plan and decompose it into subtasks in `.ai-memory/plans/`:
- **Parent Plan:** Write `.ai-memory/plans/pending/xx-<slug>.md` containing:
  - Strict relative link to the canonical spec: `Spec Reference: [02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)`
  - Architectural context, custom domain constraints, and blast radius analysis.
  - Complete mapping of deliverables (`Task-01`, `Task-02`, etc.) to subtask files.
- **Subtask Files:** Break down the plan into granular subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`, `02-<subtask>.md`, etc. Complete all spec and subtask writing within 50% of the steps budget (`PHASE_1_STEPS = N / 2`).
- **No Common Boilerplate:** Do not write common repository boilerplate, universal coding rules, banned operations, or generic guidelines inside subtask files. Common rules belong in the parent plan and root guidelines. Subtasks must contain only the unique, non-common items required for that specific subtask.

Subtasks must follow this lean, unique template:
```markdown
# Subtask [01]: [Descriptive Subtask Name]
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)
Target Files: [Strict relative paths from repo root]
Action: [Exact code changes, functions, types, and logic to modify or add]
Acceptance Criteria: [2-4 specific testable conditions proving completion]
Targeted Verification: [Specific file-level linter command or exit 0 check]
```

#### Step 4: Subtask Readiness & Spec Audit Gate

Before transitioning to execution, verify:
- **Unified Audit Blueprint:** The audit methodology, criteria, and roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating audit tasks to subagents.
- **Parallel Subagent Audit Allocation (A = 2, H = 2):** When auditing complex specs or large codebases, subagents execute parallel reading across disjoint modules to verify compliance and build granular subtask items.
- Canonical spec is authored in `02-spec/21-app/` and registered in `02-spec/21-app/readme.md`.
- Every extracted deliverable `Task-xx` has at least one corresponding subtask file.
- All subtask files are non-empty, link back to the spec, and specify disjoint target files.
- All file paths in subtasks use strict relative Git paths (zero absolute paths or `file:///` URIs).

#### Step 5: Unconditional Zero-Question Execution Mandate (Total Ban on Stopping After Spec Writing)

- **Strict 50/50 Time & Step Budget Allocation:** Spec writing and subtask generation MUST strictly complete within the first 50% of the budget (`PHASE_1_STEPS = N / 2`).
- **Zero Questions / Unconditional Execution:** As soon as Phase 1 planning completes, the master orchestrator MUST NOT pause, stop, or ask the user "Should I proceed?", "Would you like me to start execution?", or any confirmation question. There is NO question. It must immediately, unconditionally self-loop and transition directly into Phase 2 execution mode.
- **Spec Writing is Only Half the Task:** Generating specs without executing code changes is an INCOMPLETE FAILURE. The remaining 50% of the budget (`PHASE_2_STEPS = N / 2`) is dedicated strictly to modifying code, running targeted quality linters, consolidating subtasks, and completing the deliverables.

---

### Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. **Parallel Dispatch & Concrete Subagent Schema (A = 2, H = 2):**
   - Unconditionally execute code refactoring across target files in the remaining 50% of the steps budget (`PHASE_2_STEPS = N / 2`).
   - Use `invoke_subagent` to spawn up to A = 2 execution subagents concurrently. Each subagent handles an operational capacity of H = 2 (a bounded batch of up to 2 disjoint subtasks from `.ai-memory/plans/subtasks/xx-<slug>/`).
   - **Mandatory Tool Schema (`TypeName: "self"`):** Subagents executing code changes MUST use `TypeName: "self"` to inherit the parent agent's write and command capabilities (`write_to_file`, `replace_file_content`, `run_command`). (Using `research` is strictly banned for code refactoring).
   - **Invocation Payload Template:**
     ```json
     {
       "Subagents": [
         {
           "TypeName": "self",
           "Role": "Subtask Worker 01: [Feature/Module A]",
           "Model": "inherit",
           "Prompt": "[Subagent Prompt Envelope Below]"
         },
         {
           "TypeName": "self",
           "Role": "Subtask Worker 02: [Feature/Module B]",
           "Model": "inherit",
           "Prompt": "[Subagent Prompt Envelope Below]"
         }
       ]
     }
     ```

2. **Self-Contained Subagent Prompt Envelope (Eliminate Context Blindness):**
   - NEVER spawn subagents with starved one-liner instructions (e.g., NEVER just "Read 01-task.md and execute it"). Subagents spawn with a clean context and must receive a complete, self-contained **Prompt Envelope**:
     ```markdown
     You are Subagent Worker [NN].

     ### Assigned Subtasks (H = 2 Batch Capacity):
     - Subtask 1: `.ai-memory/plans/subtasks/xx-<slug>/01-<name>.md`
     - Subtask 2: `.ai-memory/plans/subtasks/xx-<slug>/02-<name>.md` (if assigned)

     ### Strict Target File Bounding Box (Disjoint Files Only):
     - Allowed Target Files: `[path/to/file1.go, path/to/file2.go]`
     - TOTAL BAN: You are strictly banned from touching or modifying ANY other files in the workspace.

     ### Non-Negotiable Coding Rules:
     1. Positive booleans only (`isReady`, `hasPermission`), NO explicit `== true`.
     2. Structured Go errors: return `*appfault.AppError`, never bare `error`.
     3. Function sizing: <= 8 lines preferred (hard cap 15 lines).
     4. Zero tests or builds: NEVER run `go test`, `pytest`, or build commands.
     5. Targeted quality check: Run only targeted file-level linters (`python 03-ai-scripts/05-guideline-autofixer.py <file>`).

     ### Completion & Reporting Contract:
     When finished, emit a structured completion summary detailing:
     - Exact files modified
     - Verification status
     - Any blockers encountered
     Conclude your response cleanly. Your output will be delivered back to the parent orchestrator via the reactive messaging system.
     ```

3. **Reactive Wakeup & Turn-Yielding Protocol (Deadlock Prevention):**
   - In Google Antigravity, background subagents run asynchronously in the platform runtime. The parent agent receives subagent completions via the **Reactive Wakeup** messaging system (`<SYSTEM_MESSAGE>`).
   - **MANDATORY YIELD RULE:** Immediately after issuing the `invoke_subagent` tool call, the parent orchestrator MUST output a brief progress note to the user (e.g. `Dispatched Subagents [01] and [02] to execute subtasks in parallel. Yielding turn to await completion...`) and **STOP CALLING TOOLS**.
   - **NO BUSY-POLLING (TOTAL BAN):** NEVER run tight-polling loops using `manage_task` or filesystem checks to wait for subagents. Ending the tool-call chain allows the platform scheduler to execute the background subagents and deliver their completion messages into the parent's inbox upon wakeup.
   - When the subagents finish, the engine wakes up the parent agent automatically with a `<SYSTEM_MESSAGE>`. The parent inspects the results, verifies acceptance criteria, and transitions to consolidation.

4. **File Locking & Disjoint Files:** Verify subagents operate on strictly distinct files using `.ai-memory/readme.md`.
5. **Execution & Coding Guidelines:** Subagents refactor code following all coding guidelines (<= 8-15 line functions, single return types, Unix LF line endings).
6. **Failure Memory & Error Recovery:** If a subagent fails or times out, the parent logs the root cause in `.ai-memory/memory/issues/xx-failure.md` and either executes the subtask directly or dispatches a targeted remediation subagent.
7. **Atomic Change Tracking:** Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json` for subsequent CI/CD verification.
8. **Total Ban on Test Running:** Do not run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`), Go (`go test`), or any test runner during routine execution turns. All test execution is strictly deferred to CI/CD pipelines and dedicated fix workflows.
9. **Total Ban on Build Checking:** Do not run build verification commands (`go build`, `npm run build`, compiler invocations). Build compilation is checked later on in CI/CD.
10. **Targeted Quality Linting Only:** Run only targeted, fast file-level linters or autofixers on specifically modified files (`exit 0`). Do not run `06-cicd-local-runner.py` or full test suites.

#### Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI)

When monitoring or checking remote CI/CD pipelines:
1. Mandatory GitMap Authority: Use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`) to query pipeline state and parse `is_running`, `status`, `etaSeconds`, and `nextAiCommand`.
2. Anti-Credit-Waste Waiting Mandate: Never busy-poll or query rapidly (`gh run view` tight loops). When pipeline is running, wait or sleep based on `etaSeconds` using `gitmap pipeline-ai status -t <etaSeconds>`:
   - etaSeconds > 120: wait 20s-30s before querying again.
   - 60 < etaSeconds <= 120: wait 10s-20s before querying again.
   - etaSeconds <= 60: wait 5s-10s before querying again.
3. Targeted Error Diagnostics: Use GitMap automated error extraction to isolate actionable failure lines without pulling verbose passing logs.

#### GitMap High-Efficiency Commands & Tooling (Search, Script Runner, Commit)

To maximize developer velocity, avoid slow generic shells, and ensure atomic consistency, heavily leverage native GitMap commands:
1. Fast Code Search & File Discovery:
   - File finding: `gitmap find "<pattern>" [-ext <ext>]` or `gitmap lf [pattern]`
   - Multi-core regex & text search: `gitmap search "<query>"` or scoped streaming `gitmap aum search "<query>" [dir] --ext <ext>`
   - Inspecting file content: `gitmap cat <filepath>`
2. Semantic Atomic Commits & Pushes (Replaces Raw Multi-Command Git Chains):
   - Feature commits: `gitmap cpf "<summary>"` (stages all files, adds `Feature: `, commits, and pushes).
   - Bug/fix commits: `gitmap cpb "<summary>"` (stages all files, adds `Bug: `, commits, and pushes).
   - Release commits: `gitmap cpr "<summary>"` (stages all files, adds `Release: `, commits, and pushes).
   - Safe pull-commit-push: `gitmap pcp "<summary>"`.
3. Cross-Platform Script Execution:
   - PowerShell runner: `gitmap pwsh "<command>"` or `gitmap ps "<command>"` (executes with `-NoProfile` and auto-fallback).
   - Bash runner: `gitmap bash "<command>"` or `gitmap sh "<command>"`.
4. Pull-All Policy (Do Not Routine-Poll):
   - Never run `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine execution loops. Checking dozens of repos introduces latency.
   - Run pull-all ONLY when explicitly requested by the user, and use JSON mode (`gitmap pae --json`) to suppress visual tabular overhead.

---

### Per-Task Agent Isolation & Workspace Subfolders (`.ai-memory/temp-agents/xx-<task-name>/`)

To prevent cross-task pollution and ensure seamless agent communication, every task must create a dedicated subfolder in `.ai-memory/temp-agents/xx-<task-name>/`:
1. Per-Task Isolation: On task start, the assigned subagent creates its isolated directory `.ai-memory/temp-agents/xx-<task-name>/`.
2. State & Progress Tracking: Create `.ai-memory/temp-agents/xx-<task-name>/state.md` documenting:
   - Task sequence and target deliverables.
   - Files assigned for modification.
   - Current subtask step and completion percentage.
3. Inter-Agent Communication & Scratch Space:
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task must be written inside `.ai-memory/temp-agents/xx-<task-name>/`.
4. On Error/Crash: Append the exact error, root cause, and `STATUS: FAILED` to `.ai-memory/temp-agents/xx-<task-name>/state.md` before exiting.
5. On Success: Mark `STATUS: DONE` in `.ai-memory/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

---

### Phase 3: Task Consolidation & File Reduction (End of Loop)

To reduce markdown file count and bloat, consolidate subtasks when a parent task is 100% complete:

1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, include a header explicitly referencing how the main task started, referencing the canonical spec `[02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)`, and documenting exactly how many steps or loops it took. Note: The canonical specification in `02-spec/21-app/` remains permanently intact in the repository as the architectural source of truth; do not delete it during consolidation.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
6. Final Step Git Commit & Push via GitMap Semantic Commit Commands (Mandatory):
   - Heavily leverage GitMap semantic commit & push commands rather than raw multiline git chains:
     - For features/tasks: `gitmap cpf "<type>(<scope>): <summary>"` (automatically stages all files, prefixes `Feature: `, commits, and pushes to remote).
     - For fixes/bugs: `gitmap cpb "<type>(<scope>): <summary>"` (automatically stages all files, prefixes `Bug: `, commits, and pushes).
     - For releases: `gitmap cpr "<type>(<scope>): <summary>"` (automatically stages all files, prefixes `Release: `, commits, and pushes).
     - For safe pull-commit-push: `gitmap pcp "<type>(<scope>): <summary>"`.
   - If GitMap CLI is unavailable, fallback to raw git: `git add -A && git commit -m "<type>(<scope>): <summary>" && git push origin <branch>`.
   - Under no circumstances commit each file individually.

---

### End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of all tasks and before concluding the turn, you MUST emit this structured verification summary in the chat response.

> [!CRITICAL]
> **STRICT LINE-BY-LINE OUTPUT MANDATE (TOTAL BAN ON HORIZONTAL CONCATENATION):**
> Every single completed task in the `Task Completion Summary` MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
> NEVER concatenate multiple tasks horizontally into a single run-on paragraph or single wrapped line.
> In Markdown, consecutive lines without bullet markers (`- `) collapse into a single run-on horizontal sentence. You MUST format each task as a discrete bullet list item (`- ✅`) followed by an explicit newline!
>
> ❌ **BANNED (Horizontal Run-on Concat):**
> `✅ #1. Task-01: [Title] — Completed ✅ #2. Task-02: [Title] — Completed ✅ #3. Task-03: [Title] — Completed`
>
> ✅ **MANDATORY (Strict Line-by-Line Vertical Markdown List):**
> ```markdown
> ### Task Completion Summary
>
> - ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
> - ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`
> - ✅ **Task-03: [Descriptive Task Title]** — `[Completed]`
> ```

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
- ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`
(If any task failed or was deferred, mark with `- ❌` or `- ⏳` on its own separate line and explain why)

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Implementation Confidence Score

- Confidence: [e.g. 98% or 100%]
- Rationale: [Detailed explanation of verified quality gates, passing linters, contract adherence, and zero regressions]

### 🤖 Independent AI Verification & Audit Prompt

At the conclusion of the turn, emit this copy-pasteable prompt for an independent auditor AI to verify the implementation against the canonical specification and verbatim requirements:

```markdown
### Independent AI Audit & Verification Instructions

You are an Independent AI Verification and Quality Auditor.
Your task is to independently audit, verify, and remediate the implementation against the canonical specification and verbatim requirements.

#### 1. Target Documents & Implemented Code:
- **Canonical Spec & Verbatim Requirements:** [02-spec/21-app/xx-<slug>.md](02-spec/21-app/xx-<slug>.md)
- **Consolidated Plan & Subtasks:** [.ai-memory/plans/completed/xx-<slug>.md](.ai-memory/plans/completed/xx-<slug>.md)
- **Modified & Implemented Code Files:**
  - [relative/path/to/modified/file1.ext](relative/path/to/modified/file1.ext)
  - [relative/path/to/modified/file2.ext](relative/path/to/modified/file2.ext)

#### 2. Verification Protocol:
1. **Strict Verbatim Inspection:** Read the canonical spec file `02-spec/21-app/xx-<slug>.md` completely, focusing on the Verbatim Requirements and Acceptance Criteria.
2. **Line-by-Line Code Comparison:** Inspect the modified code files line-by-line and verify whether every single verbatim requirement is fully implemented.
3. **Gap & Missing Items Identification:** Itemize all missing, incomplete, placeholder, or non-compliant elements.
4. **Autonomous Self-Loop Remediation:** If any gaps or missing items exist, self-loop and modify the code directly until 100% of the verbatim requirements are met. Do not ask for confirmation.
5. **Comparative Score Audit:** Evaluate and report a comparative score between the verbatim specification and the real-life output:
   - **Verbatim Adherence Score:** [X/100] (Degree to which implementation matches exact spec text and constraints)
   - **Completeness Score:** [Y/100] (Percentage of required features and acceptance criteria satisfied)
   - **Guideline Compliance Score:** [Z/100] (Adherence to booleans, error handling, function sizing, type extraction)
   - **Overall Implementation Score:** [(X+Y+Z)/3 / 100]
6. **Final Audit Verdict:** Emit a clear PASS/FAIL verdict with confidence score and rationale.
```
```

---

## De-Duplication & Continuous Loop Architecture

To prevent instruction bloat, context exhaustion, and repetitive failure loops, this prompt enforces strict structural de-duplication:

### 1. Structural De-Duplication Principles

- No Repetitive Sections: The master orchestrator follows a single unified pipeline. Pipeline steps are declared once in chronological order and not duplicated across secondary loop descriptions.
- No Boilerplate Pollution in Subtasks: Universal guidelines (coding standards, line ending rules, banned operations) exist once in the parent spec. Subtask files contain strictly unique, domain-specific requirements.
- Zero Redundant Reads: Subagents do not re-read files that were already scanned and summarized in `.ai-memory/plans/pending/xx-<slug>.md`.

### 2. Continuous 2-Phase Loop Lifecycle
```text
[Phase 1A: Prompt Capture & Task Extraction Gate]
                │
                ▼ (Chained Tool Call — Same Turn)
[Phase 1B: Master Spec & Lean Subtask Generation]
                │
                ▼ (Automatic Transition — Do Not Pause)
[Phase 2: Subagent Dispatch & Execution]
                │
                ├──> invoke_subagent(A=2, H=2)
                │         │
                │         ▼ (Yield Control to Platform Engine)
                │    [Subagents Execute in Background]
                │         │
                │         ▼ (Reactive Wakeup via <SYSTEM_MESSAGE>)
                └──> Ingest Results & Verify Acceptance Gates
                │
                ▼ (On All Subtasks Done)
[Phase 3: Task Consolidation & Atomic Git Push]
                │
                ▼
[End-of-Turn Verification & Confidence Reporting]
```

---

## Issue Destination & Root Cause Analysis (RCA) Routing Mandate

Whenever the task involves fixing an issue, bug, pipeline failure, or performing a fix with RCA (e.g., user reports a failure, provides a CI/CD error log, or commands "fix with RCA"):
1. **CI/CD Issues & Pipeline Failures:**
   - **Target Folder:** `.ai-memory/cicd-issues/`
   - **File Pattern:** `.ai-memory/cicd-issues/NN-<issue-slug>.md`
   - **Registry:** Index the issue in `.ai-memory/cicd-index.md`.
   - **Scope:** CI/CD workflows, GitHub Actions, local runner failures (`06-cicd-local-runner.py`), test runner errors, lint gate failures, or build pipeline failures.
2. **Non-CI/CD Issues (Application Bugs, Feature Defects, Logic/Runtime Errors):**
   - **Target Folder:** `02-spec/22-app-issues/` (canonical spec hierarchy Tier 22)
   - **File Pattern:** `02-spec/22-app-issues/NN-<issue-slug>.md`
   - **Structure & Registry:** Follow the 4-part structure (Reproduction / Cause / Fix / Prevention per AC-AI-001 or Why / How / Root Cause / Code Fix) and index in `02-spec/22-app-issues/readme.md` (cross-referencing in `.ai-memory/memory/issues/` for institutional memory).
   - **Scope:** Application business logic, UI bugs, CLI command errors, API crashes, and domain defects.

---

## 1. AI Fix Scripts Memory (Reusable Tooling)

- [ ] /goal Reuse First: Scanned and learned `03-ai-scripts/readme.md` before writing temporary code.
- [ ] Strict In-Repository Execution: All Python scripts executed strictly within the codebase repository root.
- [ ] Strict .ai-memory/ Folder Storage: All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] Native File Manipulator: Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] Go Generate Sync: If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 2. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in .ai-memory/cicd-issues/NN-<slug>.md (indexed in .ai-memory/cicd-index.md), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in 02-spec/22-app-issues/NN-<slug>.md (indexed in 02-spec/22-app-issues/readme.md).
- [ ] NO TEST RUNNING (TOTAL BAN): Never run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): Never launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): Never bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents must query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait or sleep based on `etaSeconds` to eliminate credit waste.
- [ ] NO PREMATURE TURN CLOSING BEFORE EXECUTION (TOTAL BAN): Never halt execution, conclude the turn, or ask the user for permission after generating specs or subtasks. Planning constitutes only 50% of the task budget; you must proceed unconditionally to Phase 2 code execution. (Note: When dispatching asynchronous background subagents via `invoke_subagent`, yielding control to allow platform reactive wakeup is mandatory and is exempt from this ban).
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpf`, `cpb`, `cpr`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.

---

## 3. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You must verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Concrete Types Centralization (`types.go`): Extracted domain structs, raw generic instantiations, and Result wrappers into dedicated `types.go` files as single reusable named types with follow-through comments (never leak raw generics like `result.Result[*Config]`).
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `*appfault.AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with is or has only (all other prefixes like can, should, was, will, did, must are banned). No negatives (`!isSuccess` is banned; use `isFail`).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`). Behavior-driven unit test names.
- [ ] Multi-Line Arguments (Rule 9a/9b): Signatures and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths: Zero absolute paths (`/absolute/path/to/...`) or `file:///` URIs.

---

## 4. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.ai-memory/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed to update all callers of modified symbols.
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete without running banned test/build commands.
- [ ] Final Step Commit & Push Verified: Staged all changes (`git add -A`), committed everything in a single grouped atomic commit, and pushed to git before ending the turn (no per-file commits).

---

## 5. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, use GitMap semantic commit commands: `gitmap cpf "<summary>"` (features), `gitmap cpb "<summary>"` (bugs), or `gitmap cpr "<summary>"` (releases) which automatically stage, commit with standardized prefixes, and push directly to the remote repository. (Fallback to `git add -A && git commit && git push` only if GitMap CLI is unavailable). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` or `gitmap cpf` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
