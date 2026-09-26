---
name: execute-parent-task
description: Autonomously orchestrate and execute parent tasks by decomposing them into subtasks and running continuous N-step self-loops with strict coding guidelines and error management.
---

# Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)

> [!IMPORTANT]
> Prompt Version: 2.6.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents, default: 2)
H = 2   (Number of hands / parallel operations per agent, default: 2)
```

N = total self-loop steps budget that the agents will perform.
A = count of autonomous subagents running concurrently (default: 2).
H = number of hands / parallel operations per agent (default: 2).

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Preamble Precedence Verification: Whatever is given before this section or prompt (user preamble, header constraints, prior instructions) has been verified as highest priority and non-negotiable, and is strictly incorporated into the task scope ahead of all other guidelines.
2. [ ] /goal Phase 1A (Step 0 - Verbatim Prompt Recording & Task Extraction Gate): Immediately capture the user prompt verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under `## User Request (Verbatim)`. If screenshot URLs or base64 data URIs are provided, decode/save them as image files (`assets/screenshots/<task-slug>-<NN>.png`) and refer back to them via relative paths in specs. Extract actionable deliverables with traceable IDs (`Task-01`, `Task-02`), and output this confirmed task breakdown directly in chat in cleanly indented markdown with vertical blank lines, task state (`State: [PENDING]`), and understanding indicator bracket (`Understood: [YES — ...]`) before any file exploration, scanning, or spec writing.
3. [ ] /goal Phase 1B (Step 1 - Canonical Spec Generation & Single-Agent Blueprint in Folder 21): Write the initial canonical specification overview and blueprint in `02-spec/21-app/xx-<slug>.md` (or directory `02-spec/21-app/xx-<slug>/` for complex features) with a single lead agent, then delegate modular spec sections to parallel subagents (A = 2, H = 2) for lossless verbatim capture and visual assets, register it in `02-spec/21-app/readme.md`, and initialize the execution plan in `.ai-memory/plans/pending/xx-<slug>.md` linking back to the spec.
4. [ ] /goal Phase 1B (Step 2 - Scan & Discover via Parallel Subagents): Use GitMap AUM discovery (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, with parallel subagents (reading files as the primary parallel task) to inventory files and map call sites without tool truncation limits.
5. [ ] /goal Phase 1B (Step 3 - Lean Subtask Decomposition): Decompose the plan into lean subtask files in `.ai-memory/plans/subtasks/xx-<slug>/01-<subslug>.md` cross-referencing the canonical spec. Subtasks must focus purely on unique task deliverables without repeating common repository boilerplate. Complete all spec and subtask writing within 50% of the steps budget (`PHASE_1_STEPS = N / 2`).
6. [ ] /goal Phase 1B (Step 4 - Readiness & Spec Audit Gate): Confirm single-agent audit roadmap is set, canonical spec is registered in `02-spec/21-app/readme.md`, all `Task-xx` deliverables are mapped to subtasks linking back to the spec, and disjoint files are assigned before execution.
7. [ ] /goal Phase 1B (Step 5 - Unconditional Zero-Question Execution Mandate): Immediately upon completing Phase 1, self-loop and transition directly into Phase 2 code execution without pausing, asking questions, or seeking user confirmation. Stopping after spec writing is strictly banned and constitutes an auto-reject failure.
8. [ ] /goal Phase 2 (Execution & Code Refactoring, Steps N/2+1..N): Unconditionally execute the code refactoring in the remaining 50% of the steps budget (`PHASE_2_STEPS = N / 2`). Spawn up to A = 2 execution subagents (H = 2 operations each) to execute subtasks on disjoint files in parallel.
9. [ ] /goal Phase 2 (Failure Memory & Error Recovery): If a subagent fails, record the failure log in `.ai-memory/plan.md` and `.ai-memory/memory/issues/`; subsequent agents must read the failure log first to remediate root causes.
10. [ ] /goal Phase 2 (Change Recording & Quality Linting): Record all modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`) and run targeted file-level linters on specifically modified files (`exit 0`). Do not run `06-cicd-local-runner.py`, unit tests, or build checks (deferred to CI/CD).
11. [ ] /goal Phase 3 (Consolidation & Atomic Push): Consolidate completed subtasks into `.ai-memory/plans/completed/xx-<slug>.md` preserving the canonical spec reference (canonical spec in `02-spec/21-app/` remains permanently intact), delete granular subtasks and pending plan, stage all changes, and push in a single grouped commit.
12. [ ] /goal Phase 3 (Completion & Confidence Reporting): Emit the final Task Completion Summary with strict line-by-line bullet format (- ✅), modified files summary, and implementation confidence score.
13. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
14. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `02-spec/03-error-manage/` for error handling architectures and *appfault.AppError.
17. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## 1. Multi-Agent Concurrency & Orchestration (A = 2, H = 2)

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- **Single-Agent Unified Blueprint Mandate:**
  - The initial execution plan, architectural decomposition, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
  - Never allow multiple agents to author disjoint or competing initial plans simultaneously. A single coherent architectural vision must lead.
- **Most Useful Parallel Tasks (Reading Files & Writing Modular Specs):**
  - When multiple agents are present (A >= 2, H >= 2), the most effective parallel tasks are reading files and authoring modular specs.
  - Subagents execute parallel disjoint tasks: reading dependencies/callers, authoring modular spec sub-files, and refactoring disjoint files.
- **Strict Concurrency Limits (A = 2, H = 2):** When dispatching work in Phase 1 (planning) or Phase 2 (execution), spawn at most A = 2 sub-agents concurrently, with H = 2 parallel operations per agent.
- **Strict Folder Bounding (`.ai-memory/`):** Subagents are strictly restricted to writing planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/readme.md`, `.ai-memory/memory/issues/`).
- **Context Diet:** When spawning a subagent, do not paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.ai-memory/plans/subtasks/xx-slug/01-<subtask-title>.md` and execute it"). The subagent must read the necessary files itself.
- **Fail Fast & Kill Stalls:** If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.

---

## 2. Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. Verbatim Prompt Capture: Capture the incoming user request verbatim. It will be recorded losslessly in both the canonical spec under `02-spec/21-app/` and the master execution plan under `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section.
2. Screenshot & Print Screen Base64 Image Ingestion Protocol: If the user request or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Convert that base64 encoding or downloaded image to the file system immediately, saving it as a persistent file under `assets/screenshots/<task-slug>-<NN>.png` or `assets/ui/<task-slug>-<NN>.png`.
   - Never leave raw, massive base64 strings or external ephemeral URLs in the prompt text, specs, or subtasks.
   - Refer back to this saved image file in the master spec, planning document, subtasks, and UI implementation notes strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<task-slug>-<NN>.png)`).
3. Actionable Deliverables Extraction: Break down whatever user requirements were given (regardless of how they were formatted) into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, `Task-03`) under `## Extracted Actionable Task List`.
4. Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):
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

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call (e.g. `write_to_file` to save the spec or `run_command` to discover files) MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

## 3. Phase 1B: Detailed Spec Generation & Lean Subtasks (Steps 1 .. N/2)

### Step 1: Canonical Application Spec Generation (Folder 21 Standard)
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

### Step 2: Scan & Discover (GitMap AUM Acceleration & Multi-Agent Parallel Reading)
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

### Step 3: Actionable Execution Plan & Lean Subtask Decomposition
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

### Step 4: Subtask Readiness & Spec Audit Gate

Before transitioning to execution, verify:
- **Unified Audit Blueprint:** The audit methodology, criteria, and roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating audit tasks to subagents.
- **Parallel Subagent Audit Allocation (A = 2, H = 2):** When auditing complex specs or large codebases, subagents execute parallel reading across disjoint modules to verify compliance and build granular subtask items.
- Canonical spec is authored in `02-spec/21-app/` and registered in `02-spec/21-app/readme.md`.
- Every extracted deliverable `Task-xx` has at least one corresponding subtask file.
- All subtask files are non-empty, link back to the spec, and specify disjoint target files.
- All file paths in subtasks use strict relative Git paths (zero absolute paths or `file:///` URIs).

### Step 5: Unconditional Zero-Question Execution Mandate (Total Ban on Stopping After Spec Writing)

- **Strict 50/50 Time & Step Budget Allocation:** Spec writing and subtask generation MUST strictly complete within the first 50% of the budget (`PHASE_1_STEPS = N / 2`).
- **Zero Questions / Unconditional Execution:** As soon as Phase 1 planning completes, the master orchestrator MUST NOT pause, stop, or ask the user "Should I proceed?", "Would you like me to start execution?", or any confirmation question. There is NO question. It must immediately, unconditionally self-loop and transition directly into Phase 2 execution mode.
- **Spec Writing is Only Half the Task:** Generating specs without executing code changes is an INCOMPLETE FAILURE. The remaining 50% of the budget (`PHASE_2_STEPS = N / 2`) is dedicated strictly to modifying code, running targeted quality linters, consolidating subtasks, and completing the deliverables.

---

## 4. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. Parallel Dispatch: Spawn up to A = 2 execution subagents (H = 2 operations each) assigned to disjoint subtasks from `.ai-memory/plans/subtasks/xx-<slug>/`. Provide minimal instructions referencing the subtask path.
2. File Locking & Disjoint Files: Verify subagents operate on distinct files using `.ai-memory/readme.md`.
3. Execution & Coding Guidelines: Subagents refactor code following all coding guidelines (<= 8-15 line functions, single return types, universal `*appfault.AppError` wrapping, Unix LF line endings).
4. Failure Memory & Error Recovery: If a subagent fails:
   - Rollback dirty changes and write the failure error log to `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`.
   - The next subagent spawned must read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
5. Atomic Change Tracking: Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json`.
6. Total Ban on Test Running: Do not run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`), Go (`go test`), or any test runner during routine execution turns. All test execution is strictly deferred to CI/CD pipelines.
7. Total Ban on Build Checking: Do not run build verification commands (`go build`, `npm run build`, compiler invocations). Build compilation is checked later on in CI/CD.
8. Targeted Quality Linting Only: Run only targeted, fast file-level linters or autofixers on specifically modified files (`exit 0`). Do not run `06-cicd-local-runner.py` or full test suites.

### Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI)

When monitoring or checking remote CI/CD pipelines:
1. Mandatory GitMap Authority: Use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`) to query pipeline state and parse `is_running`, `status`, `etaSeconds`, and `nextAiCommand`.
2. Anti-Credit-Waste Waiting Mandate: Never busy-poll or query rapidly (`gh run view` tight loops). When a pipeline is in progress, wait or sleep based on `etaSeconds` using `gitmap pipeline-ai status -t <etaSeconds>`:
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

## 5. Per-Task Agent Isolation & Workspace Subfolders (`.ai-memory/temp-agents/xx-<task-name>/`)

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

## 6. Phase 3: Task Consolidation & File Reduction (End of Loop)

To reduce markdown file count and bloat, consolidate subtasks when a parent task is 100% complete:

1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, include a header explicitly referencing how the main task started, referencing the canonical spec `[02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)`, and documenting exactly how many steps or loops it took to complete. Note: The canonical specification in `02-spec/21-app/` remains permanently intact in the repository as the architectural source of truth; do not delete it during consolidation.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
6. Final Step Git Commit & Push via GitMap Semantic Commit Commands (Mandatory):
   - Heavily leverage GitMap semantic commit & push commands rather than raw multiline git chains:
     - For features/tasks: `gitmap cpf "<summary>"` (automatically stages all files, prefixes `Feature: `, commits, and pushes to remote).
     - For fixes/bugs: `gitmap cpb "<summary>"` (automatically stages all files, prefixes `Bug: `, commits, and pushes).
     - For releases: `gitmap cpr "<summary>"` (automatically stages all files, prefixes `Release: `, commits, and pushes).
     - For safe pull-commit-push: `gitmap pcp "<summary>"`.
   - If GitMap CLI is unavailable, fallback to raw git: `git add -A && git commit -m "<summary>" && git push origin <branch>`.
   - Under no circumstances commit each file individually.

---

## 7. End-of-Turn Verification & Confidence Reporting (Mandatory Output)

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

## 8. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in .ai-memory/cicd-issues/NN-<slug>.md (indexed in .ai-memory/cicd-index.md), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in 02-spec/22-app-issues/NN-<slug>.md (indexed in 02-spec/22-app-issues/readme.md).
- [ ] NO TEST RUNNING (TOTAL BAN): Never run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): Never launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): Never bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents must query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait or sleep based on `etaSeconds` to eliminate credit waste.
- [ ] NO STOPPING AFTER SPEC WRITING (TOTAL BAN): Never halt execution, conclude the turn, or ask the user for permission after generating specs or subtasks. Planning constitutes only 50% of the task budget; you must proceed unconditionally to Phase 2 code execution.
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpf`, `cpb`, `cpr`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.

---

## 9. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

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

## 10. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.ai-memory/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed to update all callers of modified symbols.
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete without running banned test/build commands.
- [ ] Final Step Commit & Push Verified: Staged all changes (`git add -A`), committed everything in a single grouped atomic commit, and pushed to git before ending the turn (no per-file commits).

---

## 11. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, use GitMap semantic commit commands: `gitmap cpf "<summary>"` (features), `gitmap cpb "<summary>"` (bugs), or `gitmap cpr "<summary>"` (releases) which automatically stage, commit with standardized prefixes, and push directly to the remote repository. (Fallback to `git add -A && git commit && git push` only if GitMap CLI is unavailable). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` or `gitmap cpf` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.
