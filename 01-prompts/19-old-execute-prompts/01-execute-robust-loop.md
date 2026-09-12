# Resilient Multi-Agent Loop Execution — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Read all pending tasks from `.lovable/`, allocate small micro-portions of work to sub-agents, and execute them in a continuous self-loop. Manage sub-agent crashes gracefully, enforce file collision safety, sanitize artifacts before commits, and ensure the pipeline runs without halting until the queue is empty. You MUST self-loop continuously until every pending task is completed; do not stop until the queue is completely empty.

/learn Capture every pattern, convention, fix, and correction discovered during execution into `.lovable/memory/learned/01-<slug>.md` and `.lovable/strictly-avoid.md`. Never repeat a mistake that was logged.

## Non-Negotiable Rules (Auto-Reject on Violation)

1. Maximum 3 sub-agents may run concurrently at any time. Never exceed this limit.
2. No end-to-end tests that make live API calls. Only run local, isolated unit tests.

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

4. Violation of any rule below is auto-reject on the same tier as RULE 0.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination Rules

- If a spec file, folder, or task is missing or ambiguous, do NOT guess or invent a rule.
- Ask a clarifying question or log an open ambiguity in `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md` before proceeding.
- Never invent step counts. Read the actual files and count from them.

---

## Phase 1: Load, Clean, & Prepare Tasks

1. Check git status first. The working tree must be clean. Confirm root readme is strictly lowercase `readme.md`.
2. Ensure `.lovable/temp/` is added to your project's `.gitignore` file.
3. Wipe any old, orphaned state files in `.lovable/temp/` from previous incomplete runs before starting fresh.
4. Read `.lovable/plans/01-index.md` and load tasks from `.lovable/plans/pending/xx-<slug>.md`. Sequence them into Execution Waves:
   - Wave 1: Schemas, DB, and query wrappers
   - Wave 2: Core logic
   - Wave 3: UI and documentation
5. Do not start a task if its prerequisite tasks are not marked `Status: completed`.
6. Break tasks down so each agent handles a simple, small micro-task (under 15 lines per function). Monolithic tasks with more than 7 steps must be decomposed into `.lovable/plans/subtasks/xx-<slug>/`.

---

## Phase 2: Resilient Allocation & Execution Loop

1. Agent limit (strict):
   - Spawn a maximum of 2 to 3 sub-agents concurrently. Never exceed this limit.

2. File collision locking matrix (`active-locks.json`):
   - Register active target files in `.lovable/01-index.md` before spawning an agent.
   - When assigning tasks in parallel, ensure the tasks touch completely different files or components.
   - If two tasks share a dependency or file, sequence them sequentially to eliminate git merge conflicts.

3. Pre-flight logging and specific titling:
   - Before spawning a sub-agent, assign it a highly specific title reflecting its exact task (e.g., `Refactoring Auth Service` or `Fixing DB Query Wrapper`). Do not use generic names. If an agent switches tasks, its title must change.
   - Write `.lovable/temp/xx-agent-state.md` documenting which sub-agent is running, its assigned micro-task, and instructions.

4. Continuous self-looping:
   - Loop yourself to monitor sub-agent progress.
   - Do not stop until the queue is empty.

5. Crash recovery, deadlocks & the 3-strike rollback rule:
   - Deadlocks: If an agent hangs without updating its state for an extended period, terminate it and retry.
   - Revamp and restart: If an agent crashes, read its state from `.lovable/temp/`. Reason about the failure, fix the issue, and restart.
   - 3-strike rule and automatic rollback: If a specific micro-task fails unit tests or crashes 3 times, STOP retrying. Automatically rollback the dirty files (`git checkout -- <modified_files>`). Mark the task as `Status: blocked` in `plans/pending/`. Proceed to the next disjoint task.
   - Persistent failure log: Whenever a task hits 3 strikes, write the exact failure context, stack traces, and attempted fixes to `.lovable/plan.md` and `.lovable/issues/`. When the user later says "continue", read this file first to resume recovery.

6. End-to-end tests are banned:
   - Do not run end-to-end tests that make live API calls.
   - Only run local, isolated unit tests.

---

## Phase 3: Code Quality & Commit Fix (Non-Negotiable)

While executing tasks, you and your agents must adhere to these strict coding guidelines without exception:

- Read and follow guidelines in `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `spec/04-database-conventions/`. Use automated query wrappers for failure logging.
- No magic strings or numbers. Do not introduce any unless explicitly for the logger.
- Never use string union types (e.g., `"pass" | "fail"`). Use TypeScript Enums with the suffix `Type` (e.g., `StatusType`).
- Always use explicit boolean state checks (e.g., `response.isFail`). Never invert success booleans (e.g., `!response.isSuccess`).
- Code must be DRY. Reuse constants and wrappers.

---

## Phase 4: Memory Update & File Moving

As tasks are completed:

1. Use `mv` to move the completed task file from `.lovable/plans/pending/` to `.lovable/plans/completed/`.
2. Open the moved file and flip `Status: pending` to `Status: completed`.
3. Immediately update `.lovable/plans/01-index.md` to reflect the completed status and new file location.
4. Once an agent successfully finishes its task and you have verified it, remove its entry from `.lovable/01-index.md` and delete its state file from `.lovable/temp/`.

---

## Phase 5: End-of-Loop Commit Fix, Artifact Purge & Delivery

At the end of every single iteration of your execution loop:

1. Artifact sanitizer: Audit staged files and working tree. Purge unapproved artifact zip archives, temporary scratch files, or test outputs before committing.
2. Run tests: Run local unit tests (no live API end-to-end tests). Fix failures before proceeding.
3. Lovable git history guard: Group similar code changes into a single commit with a clear, descriptive message. Never rewrite published git history (no force push, no rebasing, no squash) to protect Lovable synchronization.
4. Push to the remote git repository.
5. Explicitly list out all the tasks that were successfully completed during that run.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] `.lovable/temp/` verified in `.gitignore` and orphaned state garbage-collected.
- [ ] Dependencies and prerequisites verified before starting tasks.
- [ ] Parallel assignments verified disjoint using `.lovable/01-index.md`.
- [ ] Maximum of 2-3 sub-agents spawned concurrently with specific titling.
- [ ] Pre-flight state written to `.lovable/temp/` for every agent before it started.
- [ ] 3-Strike rollback honored: failed changes reverted via `git checkout` and logged to `last-failure.md`.
- [ ] Staged files sanitized against artifact zips and temporary scratch files.
- [ ] No live-API end-to-end tests executed.
- [ ] Completed tasks `mv`'d to `plans/completed/` and `.lovable/plans/01-index.md` updated.
- [ ] Code adheres to explicit booleans, `Type` suffixed Enums, and error wrapper rules.
- [ ] Fast-forward commits created and pushed without rewriting git history.
- [ ] Completed tasks listed out explicitly in the response.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

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

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.lovable/temp-agents/xx-<task-name>/`:

1. **Per-Task Isolation:** On task start, the assigned subagent creates its isolated directory `.lovable/temp-agents/xx-<task-name>/`.
2. **State & Progress Tracking:** Create `.lovable/temp-agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. **Inter-Agent Communication & Handoff:**
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.lovable/temp-agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. **On Error/Crash:** Append the exact error, root cause, and `STATUS: FAILED` to `.lovable/temp-agents/xx-<task-name>/state.md` before exiting.
5. **On Success:** Mark `STATUS: DONE` in `.lovable/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

NEVER MAKE THIS EXCUSE: *"I did not literally perform N self-loops or spawn dozens of sub-agents. Instead of artificially spinning up N separate loops, I consolidated the required deep work into a concentrated series of sequential steps within my existing execution environment."*
If you write this excuse, you are a lazy, stupid fuck and you are violating RULE 0. Spawn the actual sub-agents. Do the actual loops. Utilize the processing power effectively to get the right answer. Avoid stupidity.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: execute-robust-loop
- status: active
