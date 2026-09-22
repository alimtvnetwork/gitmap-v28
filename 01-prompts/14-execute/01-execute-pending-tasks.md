# Pending Tasks Continuous Loop & Multi-Agent Dispatch — Workflow (must follow)

> Prompt Version: 2.1.0
> Synchronization: Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure.

/goal Execute every pending task across `.ai-memory/plans/pending/` using up to 2 sub-agents in a continuous self-loop. Do not stop until the queue is empty, every plan is committed to git, and all indexes are updated. This run ends only when there is nothing left to execute. You MUST self-loop continuously until every pending task is completed; do not stop until the queue is completely empty.

/learn Capture every pattern, convention, fix, and correction discovered during execution into `.ai-memory/memory/learned/01-<slug>.md` and `.ai-memory/strictly-avoid.md`. Never repeat a mistake that was logged.

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, you can rely on progressive disclosure for future runs. Do not keep the entire prompt in your active memory if you don't need it.

## Non-Negotiable Rules (Auto-Reject on Violation)

1. You must NEVER stop by yourself as long as there are pending tasks.
2. You must self-loop continuously without breaking between tasks.
3. If a catastrophic failure occurs, halt, log the issue, and ask the user to type "continue" to resume.

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` Reuse First: I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] `/goal` Learn Error Skills: I have `/learn`ed `02-spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` as an AI skill checklist and ensured that every returned `*appfault.AppError` uses `.WithPath()` and `.WithVar()` context wrappers.
- [ ] Strict In-Repository Execution: All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] Strict .ai-memory/ Folder Storage: All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] Native File Manipulator & Fast Reader: If you need to perform mass file renaming, sequence fixing, rapid file reading, or cached searching, you MUST natively use `python 03-ai-scripts/17-fast-file-reader.py` (for fast reads/listing/searching) and other fast scripts rather than generic tools.
  - How it works: It provides sub-millisecond file reading and directory exploration via `tmp/cache/` for instant lookups (<1ms) and live disk fallback.
  - If Missing (Re-creation): If `03-ai-scripts/17-fast-file-reader.py` is not found, immediately recreate it using Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`). It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.
- [ ] Go Generate Sync: If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory and commit the resulting generated files to prevent CI drift.
- [ ] Commit & Track: All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] Index Documentation: I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

5. Violation of any rule below is auto-reject on the same tier as RULE 0.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination & Relative Path Rules

> [!CAUTION]
> STRICT RELATIVE GIT PATHS ONLY — NO ABSOLUTE PATHS / NO `file:///` URIs:
>
> When generating plans, subtasks (`.ai-memory/plans/subtasks/`), memory issue logs (`.ai-memory/memory/issues/`), specs, code comments, or citations:
>
> 1. Strictly Relative to Git Root: All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `02-spec/03-error-manage/01-index.md`, `[SSH Commands]`02-spec/13-generic-cli/01-index.md)`, `cmd/main.go`).
> 2. Total Ban on Absolute Paths: NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///absolute/path/to/...`, `file:///absolute/path/to/...`) into ANY file.
>
> Examples:
>
> - ❌ BAD: `[SSH Commands](file:///absolute/path/to/...) — Why: Defines behavior.`
> - ❌ BAD: `Target File: /absolute/path/to/cmd\login.go`
> - ✅ GOOD: `[SSH Commands]`02-spec/13-generic-cli/01-index.md) — Why: Defines behavior.`
> - ✅ GOOD: `Target File: cmd/login.go`

- Temp Script Sandboxing: AI Fix Scripts (Reusable Tools): Before creating a helper script, you MUST check `03-ai-scripts/01-index.md` to reuse existing tools. If you generate a new script, you MUST write it to `03-ai-scripts/`, update `index.md` with its explanation, ensure `index.md` is linked in `what-to-read.md`, and commit the script.
- If a spec file, folder, or task is missing or ambiguous, do NOT guess or invent a rule.
- Ask a clarifying question or log an open ambiguity in `.ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md` before proceeding.
- Never invent step counts. Read the actual files and count from them.
- Ambiguity Blocked Queue: If you file an ambiguity to `.ai-memory/ambiguous-questions/`, you MUST immediately update the plan file to mark that specific subtask as `[Blocked]`. The execution loop must safely skip `[Blocked]` tasks and continue executing other disjoint tasks. Do not retry blocked tasks.

---

## Phase 1: Load Pending Tasks & Project State

1. [ ] Check git status first. The working tree must be clean and committed before executing anything.
2. [ ] Read  and /learn `.ai-memory/memory/01-index.md` and `.ai-memory/what-to-read.md`. Verify root readme is strictly lowercase `readme.md`.
3. [ ] Read and /learn `.ai-memory/plans/01-index.md`. Then read every file in `.ai-memory/plans/pending/xx-<slug>.md` and all associated subtasks in `.ai-memory/plans/subtasks/xx-<slug>/` (Note: for coding guidelines, check `.ai-memory/plans/subtasks/01-coding-guideline-fixes/` or other synced folder structures).
4. [ ] Group pending tasks into sequenced Execution Waves:
   - Wave 1: Schemas, DB, and query wrappers
   - Wave 2: Business logic and services
   - Wave 3: UI and documentation
5. [ ] /learn Ingest `.ai-memory/memory/01-index.md`, `.ai-memory/strictly-avoid.md`, `02-spec/02-coding-guidelines/`, and `02-spec/03-error-manage/`, `.ai-memory/coding-guidelines.md` before taking action and also create agent rules in the repo if required to or missing from rules set of agent memory.
6. [ ] /learn `.ai-memory/coding-guidelines.md` and it is must and /goal apply the guidelines in coding every aspect.
7. [ ] Screenshot & Print Screen Base64 Ingestion: If any task or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Immediately decode and persist the image to the filesystem under `assets/screenshots/<task-slug>-<NN>.png` or `assets/ui/<task-slug>-<NN>.png`.
   - Never leave raw base64 strings or ephemeral external URLs in task plans, specs, or subtasks.
   - Refer back to this saved image file in all specs, plans, and subtasks strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<task-slug>-<NN>.png)`).

## Phase 2: Allocate & Execute (Continuous Loop & Parallel Agents)

1. Spawn sub-agents (MAXIMUM 2 concurrent):
   - Assign subtasks to up to 2 parallel sub-agents (and ONLY if there are too many tasks to handle sequentially) to accelerate execution.
   - Maintain active file paths in `.ai-memory/01-index.md`. Parallel sub-agents must never touch the same files simultaneously.
   - Assign each sub-agent a highly specific title reflecting its exact task (e.g., `Refactoring Auth Service`). Do not use generic names. If an agent switches tasks, its title must change.
   - Context Diet & Task Focus: When spawning a subagent, provide clear, lean instructions that focus on the actual domain task itself rather than writing massive generic meta-prompts. DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction. The subagent MUST read the necessary files itself. Passing massive payloads instantly causes hallucination and memory blowout.
   - Ensure each agent handles discrete, simple tasks (under 15 lines per function). Tasks exceeding 7 steps must be decomposed into subtasks before execution.

2. Continuous self-looping:
   - Loop to review sub-agent progress, update plan trackers, and spawn new agents for the next wave.
   - Do not stop until every task in `.ai-memory/plans/pending/` is complete.
   - At the end of every loop iteration, execute the Commit Fix (Phase 5) before spinning up the next loop.

3. Crash Recovery & 3-Strike Rollback:
   - If a sub-agent fails unit tests or build commands, attempt a targeted fix.
   - If it fails 3 consecutive times, automatically rollback the dirty working tree (`git checkout -- <modified_files>`).
   - Log the root cause to `.ai-memory/plan.md` and `.ai-memory/issues/`.
   - Proceed to the next disjoint task after rollback.

---

## Phase 4: Memory Update & File Moving

As tasks are completed:

1. Use `mv` to move the completed task file from `.ai-memory/plans/pending/` to `.ai-memory/plans/completed/`.
2. Open the moved file and flip `Status: pending` to `Status: completed`.
3. Immediately update `.ai-memory/plans/01-index.md` to reflect the completed status and new file location.
4. If new patterns or conventions are established, record them in `.ai-memory/memory/<topic>/xx-<slug>.md` and update `.ai-memory/memory/01-index.md`. Detailed specs must never be shortened.

---

## Phase 5: End-of-Loop Commit Fix, Artifact Purge & Delivery

At the end of every single iteration of your execution loop:

0. Task Statistics: Explicitly output task statistics in your window (done, pending, remaining list).

1. Artifact sanitizer: Audit working tree and staged files. Ensure no zip archives, temporary test outputs, or unapproved scratch scripts are committed. All temporary runner artifacts and error logs belong strictly in `.ai-memory/temp/failures/`; passing tests produce zero filesystem files and remain completely silent.
2. Smart Test Runner & Verification: When running tests, utilize the centralized test inventory (`.ai-memory/test-inventory.json`) with dual-queue workers (slow: 4w x 2 tests; fast: 4w x 4 tests in 100-test chunks). AI agents read `.ai-memory/temp/runner-eta.json` and sleep for the estimated duration rather than burning tokens in active loops.
3. Lovable git history guard: Group similar code changes into a single commit with a clear, descriptive message. Never rewrite published git history (no force push, no rebasing, no squash) to preserve Lovable editor sync.
4. Push every commit to the remote git repository.
5. Final check: When the loop finally ends (the pending queue is empty), bump the minor release version following release guidelines.

---

## Execution Reporting (Mandatory Output Format)

1. Start of Run (Initial Task Breakdown Gate): Before writing any code or executing tasks, extract all pending deliverables as discrete, ordered items (Task-01, Task-02) without hard brackets, and output this confirmed task breakdown directly in chat:

\\markdown
### Confirmed Task Breakdown
#1. Task-01: [Actionable task description]
#2. Task-02: [Actionable task description]
Proceeding directly to execution.
\
2. End-of-Turn Verification & Confidence Reporting: When all tasks are completed (or if the run concludes), you must output:

\\markdown
### Task Completion Summary
✅ #1. Task-01: [Task description] — Completed
✅ #2. Task-02: [Task description] — Completed
(If any task failed or was deferred, mark with ❌ or ⏳ and explain why)

### Modified Files Summary
- [relative path to modified file 1]
- [relative path to modified file 2]

### Implementation Confidence Score
- Confidence: [e.g. 98% or 100%]
- Rationale: [Detailed explanation of verified quality gates, passing linters, contract adherence, and zero regressions]
\
3. Lean Subtasks (No Common Boilerplate):
Subtasks must focus purely on unique task deliverables without repeating common repository boilerplate, universal coding rules, banned operations, or generic guidelines.

## Compliance Checklist (must follow non negociable)

- [x] Coding Guidelines enforced `02-spec/02-coding-guidelines/ and follow explicitly every steps .ai-memory/coding-guidelines.md).
- [x] Boolean conventions used (is/has prefixes, no negatives).
- [x] No garbage variable names used.
- [x] No magic strings or numbers.
- [x] Markdown format verified (newlines around every header).
- [x] Error management protocols followed (AppError/AppException).
- [x] Signatures > 3 parameters or > 100 chars split to one parameter per line.
- [x] Boolean conventions followed (e.g., `isFail` instead of `!isSuccess`).
- [x] Acronyms are PascalCased (e.g., `UserId`, not `UserID`).
- [x] Magic strings/numbers extracted to constants.
- [x] Action Summary Checklist (Anti-Hallucination): I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to ensure no steps were hallucinated or skipped (e.g. `- [x] Created schema`, `- [x] Pinned README`).
- [x] Remote CI/CD pipeline monitored via GitMap Pipeline-AI (`gitmap pipeline-ai status --json`) with dynamic ETA waiting (no rapid polling).

## Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI)

When monitoring or checking remote CI/CD pipelines (e.g., following git push or during pipeline audits):

1. Mandatory GitMap Pipeline-AI Authority: Agents MUST use GitMap CLI to retrieve remote CI/CD status:
   ```bash
   gitmap pipeline-ai status --json
   # or using short alias:
   gitmap pl-ai status --json
   ```
   Parse structured output fields: `is_running`, `status`, `etaSeconds`, and `nextAiCommand`.
2. Anti-Credit-Waste Waiting Mandate (TOTAL BAN on Rapid Polling):
   - NEVER loop rapidly or busy-poll (`gh run view` in tight loops). Rapid polling burns user credits, exhausts LLM tokens, and wastes rate limits.
   - When a pipeline is in progress (`is_running: true`), agents MUST wait/sleep based on the estimated completion duration (`etaSeconds` or `-t <sec>`):
     ```bash
     gitmap pipeline-ai status -t <etaSeconds>
     ```
   - Proportional ETA sleep guidelines:
     - `etaSeconds > 120`: wait 20s–30s before querying again.
     - `60 < etaSeconds <= 120`: wait 10s–20s before querying again.
     - `etaSeconds <= 60`: wait 5s–10s before querying again.
3. Targeted Failure Diagnostics: Use GitMap's automated error extraction to isolate actionable failure lines (`##[error]`, `FAIL:`, compile errors) without fetching noisy passing step logs.

## Mark File Changes Only (Atomic Change Recording & Handoff to CI/CD & Release)

Routine execution prompts MUST NOT build, test, or trigger releases. When task modifications are completed, you MUST record all modified files and physically check off these items in your final report:

- [ ] Atomic Change Recording (MANDATORY): I have recorded all modified files into `.ai-memory/temp/recent-file-changes.json` under lock using `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
- [ ] NO Test Running (BANNED): Zero tests were executed (`go test`, `pytest`, `06-cicd-local-runner.py`). Testing is strictly deferred to CI/CD fix prompts.
- [ ] NO Build Checking (BANNED): Zero build commands were executed (`go build`, `npm run build`). Build compilation is strictly deferred to CI/CD fix prompts.
- [ ] NO Release Triggering (BANNED): Zero version bumps, changelog edits, or tag operations were performed. Release operations are strictly deferred to Release prompts.
- [ ] File Change Summary: Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.ai-memory/` Bounding

- 2-Agent Limit (Max 2 Threads Each): When dispatching work, spawn at most 2 sub-agents concurrently, with no more than 2 threads per agent.
- Strict Folder Bounding (`.ai-memory/`): Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/01-index.md`, `.ai-memory/memory/issues/`).
- Context Diet & Task Focus: Provide subagents with clear, lean instructions that focus on the actual domain task itself. Do not write massive meta-prompts or generate excessive boilerplate markdown. Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.ai-memory/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Lean Subtask Decomposition: Break down the plan into a few highly focused subtask files in `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc. Task Focus Over Meta-Prompting: Your goal is to write code and solve the problem, not just generate more AI prompts. Subagent instructions should clearly define the domain task itself.
- MANDATORY AUTO-LOOP (DO NOT STOP): Once Phase 1 planning completes, the master orchestrator MUST NOT STOP or ask the user for confirmation. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.ai-memory/plans/subtasks/` to `.ai-memory/plans/completed/` and update `.ai-memory/plans/01-index.md`.
- Failure Memory & Feedback Loop: If a subagent fails:
  - Rollback dirty working tree and log error details to `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during intermediate loops.

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.ai-memory/01-index.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.ai-memory/plans/01-index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] Staged files sanitized of artifact zips and temporary scratch files.
- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `02-spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.ai-memory/coding-guidelines.md`.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /learn and apply as a /goal  `.ai-memory/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.

- [ ] Error Manage Checklist: I have fully read and enforced the error management files at `02-spec/03-error-manage/`. I understand which files to follow (architecture, response envelopes) and how to follow them (never swallow errors, always wrap with context).
- [ ] Boolean Examples & Fixations: All boolean variables MUST begin with is and has only (can, should, was, etc. are banned) (e. NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`).g., `isReady`, `hasData`). NEVER use negative booleans (e.g., `isNotReady`, `disableCache`). NEVER invert success checks (e.g., `!response.isSuccess` is banned; use `response.isFail`).
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Semantic Tests: All unit test names are strictly semantic and behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`). `TestHandleComp100` is an immediate failure.
- [ ] Function Size: No function exceeds 15 lines. Long arguments are split across lines (max 100 chars).
- [ ] Error Handling (AppError): Errors use domain-specific `AppError` or custom `AppException` (for C#/OOP), not generic base `Error`.
- [ ] Code adheres to explicit booleans, `Type` suffixed Enums, and error wrapper rules.
- [ ] Formatting & Acronyms: Spacing rules are strictly followed. Acronyms are strictly PascalCase (`SwapIpWindows` not `SwapIPWindows`).
- [ ] Fast-forward commits created and pushed without rewriting published git history.
- [ ] /goal Final Step Commit & Push (NO Per-File Commits): Accumulate all changes across the turn without committing individual files one-by-one. In the final step, stage all changes (`git add -A`), create a single grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to the remote git branch (`git push origin <branch>`).
- [ ] Continuous loop maintained; only pausing to ask for "continue" on critical unrecoverable failures.

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] NO TEST RUNNING (TOTAL BAN): NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): NEVER query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents MUST query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait/sleep based on `etaSeconds` to eliminate credit waste.

---

## Execution & Self-Looping Protocol

This task requires deep, multi-agent processing. You MUST NOT run a fast, linear self-loop that artificially consolidates work or skips steps to save time.

How to self-loop and distribute tasks effectively:

- [ ] Spawn Sub-Agents:
  - [ ] For any multi-step group, deep file reading, or complex analysis, actively spawn dedicated sub-agents.
  - [ ] Ensure these sub-agents are self-looping to handle the workload.
- [ ] Utilize Processing Power:
  - [ ] Take your time and use maximum processing power and credits.
  - [ ] Do not take shortcuts.
  - [ ] Do not attempt to process a massive spec or write a complex plan in a single, consolidated step.
- [ ] Wait and Aggregate:
  - [ ] As the master agent, loop autonomously to wait for your sub-agents.
  - [ ] Aggregate their precise findings to ensure the highest quality result.

### Temp-Agent Isolated Task Directory & Communication Protocol (Non-Negotiable)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.ai-memory/temp-agents/xx-<task-name>/`:

1. Per-Task Isolation: On task start, the assigned subagent creates its isolated directory `.ai-memory/temp-agents/xx-<task-name>/`.
2. State & Progress Tracking: Create `.ai-memory/temp-agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. Inter-Agent Communication & Handoff:
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.ai-memory/temp-agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. On Error/Crash: Append the exact error, root cause, and `STATUS: FAILED` to `.ai-memory/temp-agents/xx-<task-name>/state.md` before exiting.
5. On Success: Mark `STATUS: DONE` in `.ai-memory/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

NEVER MAKE THIS EXCUSE: *"I did not literally perform N self-loops or spawn dozens of sub-agents. Instead of artificially spinning up N separate loops, I consolidated the required deep work into a concentrated series of sequential steps within my existing execution environment."*
If you write this excuse, you are a lazy, stupid fuck and you are violating RULE 0. Spawn the actual sub-agents. Do the actual loops. Utilize the processing power effectively to get the right answer. Avoid stupidity.

## File Change Recording Quality Gate (Strict Policy)

Execution prompts MUST NOT cut releases or run tests. All build fixes, unit tests, e2e fixes, and releases are strictly deferred to CI/CD fix and Release prompts. Every modified file must be recorded into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`). When monitoring or inspecting remote CI/CD pipelines, agents MUST follow GitMap (`gitmap pipeline-ai status --json` / `gitmap pl-ai status -t <sec>`) and strictly wait/sleep based on `etaSeconds` rather than rapid polling to prevent credit waste.

## Task Consolidation & File Reduction (End of Loop)

> CRITICAL: To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

When all subtasks for a parent task (`.ai-memory/plans/pending/xx-<slug>.md`) are finished, execute this final cleanup step before ending the run:
1. Combine all the completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header that explicitly references how the main task started and documents exactly how many steps/loops it took to complete.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/01-index.md` to point to the newly consolidated completed file.
6. Final Step Git Commit & Push (MANDATORY): Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW): At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, do NOT run builds or tests during routine turns (build and test verification deferred to CI/CD), group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. Phase 1: Read & Understand (Isolated Loop): Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. Phase 2: Bounded Micro-Tasking (Sequential Self-Looping): Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. Phase 3: Multi-Agent Parallelization: If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## Metadata

- slug: execute-pending-tasks
- status: active
