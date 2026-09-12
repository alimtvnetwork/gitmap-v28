---
name: execute-batched-loop
description: >-
  Use this skill when the user asks you to execute a batch of pending tasks using a batched loop and 3 sub-agents.
---

Run again if said: go, continue, or next

# Instruction (must follow): Execute Batched Loop (3 Agents, Chunked Commits)

/goal Execute pending tasks from `.lovable/plans/pending/` using a strictly batched multi-agent loop. Use exactly 3 sub-agents, assign small micro-task chunks per agent, enforce file collision safety through a locking matrix, sanitize artifacts before commits, handle crashes via `.lovable/temp/`, and push chunked commits to git without failure. At the end of every loop, explicitly list task statistics in your output window. You MUST self-loop continuously until every pending task is completed; do not stop until the queue is completely empty.

/learn Capture every pattern, convention, fix, and correction discovered during execution into `.lovable/memory/learned/01-<slug>.md` and `.lovable/strictly-avoid.md`. Never repeat a mistake that was logged.

## Non-Negotiable Rules (Auto-Reject on Violation)

1. Maximum 3 sub-agents may run concurrently at any time. Never exceed this limit.
2. No end-to-end tests that make live API calls. Only run local, isolated unit tests.
3. At the end of every loop, output explicit task statistics (done, pending, remaining list).

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

5. Violation of any rule below is auto-reject on the same tier as RULE 0.

## Anti-Hallucination & Relative Path Rules

> [!CAUTION]
> **STRICT RELATIVE GIT PATHS ONLY — NO ABSOLUTE PATHS / NO `file:///` URIs:**
>
> When generating plans, subtasks (`.lovable/plans/subtasks/`), memory issue logs (`.lovable/memory/issues/`), specs, code comments, or citations:
> 1. **Strictly Relative to Git Root:** All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `spec/03-error-manage/01-index.md`, `[SSH Commands](spec/13-generic-cli/01-index.md)`, `cmd/main.go`).
> 2. **Total Ban on Absolute Paths:** NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///absolute/path/to/...`, `file:///absolute/path/to/...`) into ANY file.
>
> **Examples:**
> - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...) — Why: Defines behavior.`
> - ❌ **BAD:** `Target File: /absolute/path/to/cmd\login.go`
> - ✅ **GOOD:** `[SSH Commands](spec/13-generic-cli/01-index.md) — Why: Defines behavior.`
> - ✅ **GOOD:** `Target File: cmd/login.go`

- If a spec file, folder, or task is missing or ambiguous, do NOT guess or invent a rule.
- Ask a clarifying question or log an open ambiguity in `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md` before proceeding.
- Never invent step counts. Read the actual files and count from them.

## Phase 1: Pre-Flight & Gitignore Enforcement (Non-Negotiable)

1. The working tree must be clean. Confirm root readme is strictly lowercase `readme.md`.
2. Verify that `.lovable/temp/` is explicitly added to `.gitignore`. This folder is for crash identification and lockfiles and must never be committed.
3. Wipe any orphaned state files in `.lovable/temp/` from previous runs.
4. Group pending tasks into Execution Waves:
   - Wave 1: DB schemas and query wrappers
   - Wave 2: Business logic and services
   - Wave 3: UI and documentation

## Phase 2: Allocation & Execution (Strict 3x3 Rule & Locking Matrix)

1. Strict limits:
   - Spawn up to 3 sub-agents to run in parallel.
2. Chunking micro-tasks:
   - Each agent is assigned a chunk of simple, small micro-tasks (under 15 lines per function) to complete sequentially in its own context.
   - Tasks exceeding 7 steps must be decomposed into subtasks.
3. File collision locking matrix (`active-locks.json`):
   - Register active target files in `.lovable/01-index.md`.
   - Ensure parallel tasks touch completely disjoint files to prevent git merge conflicts.
4. Temp folder logging and specific titling (mandatory):
   - Spawn the sub-agent with a highly specific title reflecting its exact task (e.g., `Refactoring Auth Service` or `Fixing DB Query Wrapper`). Do not use generic names. If an agent switches chunks, its title must change.
   - Log its assigned chunk of tasks to `.lovable/temp/xx-agent-state.md`.
5. Crash identification and 3-strike rollback:
   - If an agent fails or crashes, inspect its state in `.lovable/temp/`.
   - If an agent fails 3 times, automatically revert dirty changes (`git checkout -- <files>`).
   - Log root cause to `.lovable/plan.md` and `.lovable/issues/`.
   - Restart a new agent for the next disjoint chunk.

## Phase 3: Code Quality (Non-Negotiable)

While executing tasks, you and your agents must adhere to these strict coding guidelines without exception:

- Read and follow guidelines in `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `spec/04-database-conventions/`.
- No magic strings or numbers. Do not introduce any unless explicitly for the logger.
- Never use string union types (e.g., `"pass" | "fail"`). Use TypeScript Enums with the suffix `Type` (e.g., `StatusType`).
- Always use explicit boolean state checks (e.g., `response.isFail`). Never invert success booleans (e.g., `!response.isSuccess`).
- Code must be DRY. Reuse constants and wrappers.

## Phase 4: Chunked Delivery, Artifact Purge & File Moving

When a chunk of tasks is completed by the agents, do the following before starting the next loop iteration:

1. Use `mv` to move the completed task files from `.lovable/plans/pending/` to `.lovable/plans/completed/`.
2. Open the moved files and change `Status: pending` to `Status: completed`.
3. Update `.lovable/plans/01-index.md` to reflect the new file locations.
4. Artifact sanitizer: Audit staged files. Purge unapproved artifact zip archives, temporary scratch files, or test outputs before committing.
5. Lovable git history guard: Run local tests (no live API calls). Commit code with a clear descriptive message. Never rewrite published git history (no force push, no rebasing, no squash). Push to git cleanly without failure.

## Phase 5: Output Window Stats (Mandatory Every Loop)

Every time you return a response or complete a loop iteration, explicitly output the following statistics:

- Tasks Done (This Chunk): [Number of tasks completed]
- Total Completed: [Total number of tasks in `.lovable/plans/completed/`]
- Total Pending: [Number of tasks remaining in `.lovable/plans/pending/`]
- Remaining Tasks List: [List the specific filenames/slugs of the tasks remaining]

## Execution Reporting (Mandatory Output Format)

1. Start of Run (Initial Output): Before writing any code, explicitly list out all pending tasks in your output window.
2. End of Run Summary: When all tasks are completed (or if the run concludes), you MUST output a comprehensive final summary containing:
   - Completed Tasks: Explicit list of what was successfully completed.
   - Pending Tasks Left: Explicit list of any tasks still remaining.
   - Quality Assessment: A brief summary of how well the execution went.
   - Compliance Checklist: A markdown checklist explicitly verifying that you followed the rules:

## Compliance Checklist (must follow non negociable)

- [x] Coding Guidelines enforced (spec/02-coding-guidelines/ and follow explicitly every steps .lovable/coding-guidelines.md).
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

## End of Tunnel Release (Anti-Hallucination Checklist)

Past execution turns were sloppy and failed to pin READMEs or bump versions. To prevent this hallucination, when EVERYTHING is completely finished (at the very end of the tunnel), you MUST trigger a release and physically check off these items in your final report:

- [ ] **Full Unit Test & CI/CD Verification (MANDATORY):** I have executed `python 03-ai-scripts/06-cicd-local-runner.py --run-tests` and verified that 100% of all unit tests and quality gates pass green (`exit 0`).
- [ ] **Test Inventory Validation:** I have checked `.lovable/temp/recent-file-changes.json` against `.lovable/test-inventory.json` and verified all tests associated with modified files pass.
- [ ] Minor Bump: I have bumped the MINOR version in the canonical `version.json` file.
- [ ] Test File Ban: I have strictly excluded all test files (`*test*`, `*.spec.*`) from version scanning.
- [ ] Root readme.md (lowercase always) Pinning (FATAL): I have pinned the latest release version into the root `readme.md` file! I have verified badges and install snippets match the new version.
- [ ] Changelog Formatting: I have updated the changelog exactly according to the `version.json` format.
- [ ] Release Architecture Map: I have maintained `.lovable/memory/01-index.md`, enqueued it in `what-to-read.md`, and linked it in the root `readme.md`.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] `.gitignore` verified to exclude `.lovable/temp/` and garbage collection executed.
- [ ] Strictly up to 3 agents spawned, each assigned disjoint files tracked in `.lovable/01-index.md`.
- [ ] Pre-flight state written to `.lovable/temp/` for every agent.
- [ ] 3-Strike rollback honored with `git checkout` and logged to `last-failure.md`.
- [ ] Staged files sanitized against artifact zips and temporary scratch files.
- [ ] No end-to-end live API tests executed.
- [ ] Completed task files `mv`'d and `.lovable/plans/01-index.md` updated.
- [ ] Fast-forward commit created and pushed without rewriting git history.
- [ ] Output window explicitly lists "Done", "Pending", and remaining task names.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Self Loop until all pending tasks are done.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.

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

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## Task Consolidation & File Reduction (End of Loop)

> **CRITICAL:** To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

When all subtasks for a parent task (`.lovable/plans/pending/xx-<slug>.md`) are finished, execute this final cleanup step before ending the run:
1. Combine all the completed granular subtasks from `.lovable/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.lovable/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header that explicitly references how the main task started and documents exactly how many steps/loops it took to complete.
3. Delete the original granular `.md` files in `.lovable/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.lovable/plans/pending/xx-<slug>.md`.
5. Update `.lovable/plans/01-index.md` to point to the newly consolidated completed file.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: execute-batched-loop
- status: active
