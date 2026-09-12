# [V2] Batched Loop & Execution Wave Orchestration — Workflow (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute the pending tasks by decomposing them into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 100
```

N = total self-loop steps budget that the agents will perform.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, you can rely on progressive disclosure for future runs. Do not keep the entire prompt in your active memory if you don't need it.

## Non-Negotiable Rules (Auto-Reject on Violation)

1. Maximum 3 sub-agents may run concurrently at any time. Never exceed this limit.
2. No end-to-end tests that make live API calls. Only run local, isolated unit tests.
3. At the end of every loop, output explicit task statistics (done, pending, remaining list).

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] `/goal` **Learn Error Skills:** I have `/learn`ed `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` as an AI skill checklist and ensured that every returned `*appfault.AppError` uses `.WithPath()` and `.WithVar()` context wrappers.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator & Fast Reader:** If you need to perform mass file renaming, sequence fixing, rapid file reading, or cached searching, you MUST natively use `python 03-ai-scripts/17-fast-file-reader.py` (for fast reads/listing/searching) and other fast scripts rather than generic tools.
  - **How it works:** It provides sub-millisecond file reading and directory exploration via `tmp/cache/` for instant lookups (<1ms) and live disk fallback.
  - **If Missing (Re-creation):** If `03-ai-scripts/17-fast-file-reader.py` is not found, immediately recreate it using Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`). It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

5. Violation of any rule below is auto-reject on the same tier as RULE 0.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination & Relative Path Rules

> [!CAUTION]
> **STRICT RELATIVE GIT PATHS ONLY — NO ABSOLUTE PATHS / NO `file:///` URIs:**
>
> When generating plans, subtasks (`.lovable/plans/subtasks/`), memory issue logs (`.lovable/memory/issues/`), specs, code comments, or citations:
>
> 1. **Strictly Relative to Git Root:** All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `spec/03-error-manage/01-index.md`, `[SSH Commands](spec/13-generic-cli/01-index.md)`, `cmd/main.go`).
> 2. **Total Ban on Absolute Paths:** NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///absolute/path/to/...`, `file:///absolute/path/to/...`) into ANY file.
>
> **Examples:**
>
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

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline. Do not skip steps.

```text
N = 100  (Total self-loop steps budget)
PHASE_1_STEPS = N / 2   (Planning & Subtask Generation)
PHASE_2_STEPS = N / 2   (Parallel Execution & QA)
```

### Phase 1: Planning Mode & Subtask Generation FIRST (Steps 1 .. N/2)

1. **Scan & Discover:** Use the `invoke_subagent` tool to spawn exactly 2 planning subagents. Their role is to deeply scan the codebase for target changes.
2. **Master Spec Generation:** Save the master architectural plan into `.lovable/plans/pending/xx-<slug>.md`. Write down 3–5 custom rules or constraints unique to this task inside the spec file.
3. **Lean Subtask Decomposition:** Break down the master plan into granular, single-responsibility subtask files in `.lovable/plans/subtasks/xx-<slug>/01-<subtask>.md`.
   *Subtasks MUST follow this lean template to prevent bloat:*
   > `# Subtask: [Name]`
   > `**Target Files:** [Relative paths]`
   > `**Action:** [Exact code changes required]`
   > `**Constraints:** [Key rules to follow]`
4. **MANDATORY AUTO-LOOP (DO NOT STOP):** As soon as Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. **Parallel Dispatch:** Use the `invoke_subagent` tool to spawn exactly 2 execution subagents (max 2 threads each) assigned to disjoint subtasks from `.lovable/plans/subtasks/xx-<slug>/`. Provide subagents with minimal instructions (e.g., "Read `.lovable/plans/subtasks/xx-slug/01-task.md` and execute it").
2. **Execution & Coding Guidelines:** Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, Unix LF line endings).
3. **Failure Memory & Error Recovery:** If a subagent fails, record the failure log in `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`; subsequent agents MUST read the failure log first to remediate root causes.
4. **Local Verification:** Run targeted linters on modified files ensuring `exit 0`. (Full CI/CD runner `--run-tests` runs ONLY at the final release ceremony).
5. **Atomic Change Tracking:** Append all modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.lovable/test-inventory.json`.

### Phase 3: Task Consolidation & File Reduction (End of Loop)

> **CRITICAL:** To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

1. Combine all the completed granular subtasks from `.lovable/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.lovable/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header explicitly referencing how the main task started and documenting exactly how many steps/loops it took.
3. Delete the original granular `.md` files in `.lovable/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.lovable/plans/pending/xx-<slug>.md`.
5. Update `.lovable/plans/01-index.md` to point to the newly consolidated completed file.

---

## 1. AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] `/goal` **Learn Error Skills:** I have `/learn`ed `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` as an AI skill checklist and ensured that every returned `*appfault.AppError` uses `.WithPath()` and `.WithVar()` context wrappers.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator & Fast Reader:** If you need to perform mass file renaming, sequence fixing, rapid file reading, or cached searching, you MUST natively use `python 03-ai-scripts/17-fast-file-reader.py` (for fast reads/listing/searching) and other fast scripts rather than generic tools.
  - **How it works:** It provides sub-millisecond file reading and directory exploration via `tmp/cache/` for instant lookups (<1ms) and live disk fallback.
  - **If Missing (Re-creation):** If `03-ai-scripts/17-fast-file-reader.py` is not found, immediately recreate it using Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`). It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

5. Violation of any rule below is auto-reject on the same tier as RULE 0.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination & Relative Path Rules

> [!CAUTION]
> **STRICT RELATIVE GIT PATHS ONLY — NO ABSOLUTE PATHS / NO `file:///` URIs:**
>
> When generating plans, subtasks (`.lovable/plans/subtasks/`), memory issue logs (`.lovable/memory/issues/`), specs, code comments, or citations:
>
> 1. **Strictly Relative to Git Root:** All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `spec/03-error-manage/01-index.md`, `[SSH Commands](spec/13-generic-cli/01-index.md)`, `cmd/main.go`).
> 2. **Total Ban on Absolute Paths:** NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///absolute/path/to/...`, `file:///absolute/path/to/...`) into ANY file.
>
> **Examples:**
>
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

---

## 2. MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
