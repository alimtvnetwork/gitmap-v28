# [V2] Batched Loop & Execution Wave Orchestration — Workflow (must follow)

> [!IMPORTANT]
> Prompt Version: 2.6.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, in a header alert block, or in the incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

/goal Autonomously orchestrate and execute the pending tasks by decomposing them into subtasks and running a continuous N-step self-loop until completion without a single failure.

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

### Multi-Agent Parallel Task Allocation & Orchestration (A = 2, H = 2)

When multiple autonomous agents are present (A >= 2, H >= 2):
1. **Single-Agent Unified Blueprint Mandate:**
   - The initial plan, pending task scoping, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
   - Never allow multiple agents to author disjoint or competing audit plans simultaneously. A single coherent architectural vision must lead.
2. **Most Useful Parallel Tasks (Reading Files & Writing Modular Specs):**
   - Once the unified blueprint is authored by the lead agent, the most effective parallel tasks for subagents (A = 2, H = 2) are:
     - **Reading Files:** Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting codebase context in parallel using GitMap AUM (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary and Python scripts (`03-ai-scripts/17-fast-file-reader.py`, `03-ai-scripts/11-fast-file-scanner.py`) as fallback.
     - **Writing Modular Specs & Lean Subtasks:** Authoring modular spec sections and subtasks in parallel adhering to the lead agent's blueprint.
3. **Spec Writing & Spec Audit Emphasis:**
   - **Spec Writing Section:** The lead agent defines the spec overview and architecture boundaries first; subagents are then spawned in parallel to flesh out disjoint spec modules and acceptance criteria concurrently.
   - **Spec Audit Section:** The lead agent establishes the audit methodology and roadmap first; subagents are then dispatched concurrently to inspect disjoint code areas and generate granular subtask files.
4. **Execution Mode (Disjoint Refactoring):**
   - Subagents execute parallel disjoint refactoring tasks across non-overlapping files and run targeted file-level linters (`exit 0`).

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, you can rely on progressive disclosure for future runs. Do not keep the entire prompt in your active memory if you don't need it.

## Non-Negotiable Rules (Auto-Reject on Violation)

1. Maximum 3 sub-agents may run concurrently at any time. Never exceed this limit.
2. TOTAL BAN on test running and build checking during routine execution: DO NOT run tests using Python scripts, Go (`go test`), or any test runner. DO NOT check builds (`go build`, compiler checks). All test execution and build verification is deferred to CI/CD.
3. At the end of every loop, output explicit task statistics (done, pending, remaining list).

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` Reuse First: I have rigorously scanned and `/learn`ed `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code.
- [ ] `/goal` Learn Error Skills: I have `/learn`ed `02-spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` as an AI skill checklist and ensured that every returned `*appfault.AppError` uses `.WithPath()` and `.WithVar()` context wrappers.
- [ ] Strict In-Repository Execution: All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] Strict .ai-memory/ Folder Storage: All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] Native File Manipulator & Fast Reader: If you need to perform mass file renaming, sequence fixing, rapid file reading, or cached searching, you MUST natively use `python 03-ai-scripts/17-fast-file-reader.py` (for fast reads/listing/searching) and other fast scripts rather than generic tools.
  - How it works: It provides sub-millisecond file reading and directory exploration via `tmp/cache/` for instant lookups (<1ms) and live disk fallback.
  - If Missing (Re-creation): If `03-ai-scripts/17-fast-file-reader.py` is not found, immediately recreate it using Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`). It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.
- [ ] Go Generate Sync: If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] Commit & Track: All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] Index Documentation: I have updated `03-ai-scripts/readme.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

Violation of any rule below is auto-reject on the same tier as RULE 0.

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
> 1. Strictly Relative to Git Root: All file paths, markdown links, citations, and task targets MUST be relative paths starting from the repository root (e.g. `02-spec/03-error-manage/readme.md`, `[SSH Commands](02-spec/13-generic-cli/readme.md)`, `cmd/main.go`).
> 2. Total Ban on Absolute Paths: NEVER write drive letters or absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///absolute/path/to/...`, `file:///absolute/path/to/...`) into ANY file.
>
> Examples:
>
> - ❌ BAD: `[SSH Commands](file:///absolute/path/to/...) — Why: Defines behavior.`
> - ❌ BAD: `Target File: /absolute/path/to/cmd\login.go`
> - ✅ GOOD: `[SSH Commands](02-spec/13-generic-cli/readme.md) — Why: Defines behavior.`
> - ✅ GOOD: `Target File: cmd/login.go`

- If a spec file, folder, or task is missing or ambiguous, do NOT guess or invent a rule.
- Ask a clarifying question or log an open ambiguity in `.ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md` before proceeding.
- Never invent step counts. Read the actual files and count from them.

## Phase 1: Pre-Flight & Gitignore Enforcement (Non-Negotiable)

1. The working tree must be clean. Confirm root readme is strictly lowercase `readme.md`.
2. Verify that `.ai-memory/temp/` is explicitly added to `.gitignore`. This folder is for crash identification and lockfiles and must never be committed.
3. Wipe any orphaned state files in `.ai-memory/temp/` from previous runs.
4. Group pending tasks into Execution Waves:
   - Wave 1: DB schemas and query wrappers
   - Wave 2: Business logic and services
   - Wave 3: UI and documentation
5. Screenshot & Print Screen Base64 Ingestion: If any task or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Convert that base64 encoding or downloaded image to the local file system immediately under `assets/screenshots/<task-slug>-<NN>.png` or `assets/ui/<task-slug>-<NN>.png`.
   - Never leave raw base64 strings or ephemeral external URLs in task plans, specs, or subtasks.
   - Refer back to this saved image file in all specs, plans, and subtasks strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<task-slug>-<NN>.png)`).

## Phase 2: Allocation & Execution (Strict 3x3 Rule & Locking Matrix)

1. Strict limits:
   - Spawn up to 3 sub-agents to run in parallel.
2. Chunking micro-tasks:
   - Each agent is assigned a chunk of simple, small micro-tasks (under 15 lines per function) to complete sequentially in its own context.
   - Tasks exceeding 7 steps must be decomposed into subtasks.
3. File collision locking matrix (`active-locks.json`):
   - Register active target files in `.ai-memory/readme.md`.
   - Ensure parallel tasks touch completely disjoint files to prevent git merge conflicts.
4. Temp folder logging and specific titling (mandatory):
   - Spawn the sub-agent with a highly specific title reflecting its exact task (e.g., `Refactoring Auth Service` or `Fixing DB Query Wrapper`). Do not use generic names. If an agent switches chunks, its title must change.
   - Log its assigned chunk of tasks to `.ai-memory/temp/xx-agent-state.md`.
5. Crash identification and 3-strike rollback:
   - If an agent fails or crashes, inspect its state in `.ai-memory/temp/`.
   - If an agent fails 3 times, automatically revert dirty changes (`git checkout -- <files>`).
   - Log root cause to `.ai-memory/plan.md` and `.ai-memory/issues/`.
   - Restart a new agent for the next disjoint chunk.

## Phase 3: Code Quality (Non-Negotiable)

While executing tasks, you and your agents must adhere to these strict coding guidelines without exception:

- Read and follow guidelines in `02-spec/02-coding-guidelines/`, `02-spec/03-error-manage/`, and `02-spec/04-database-conventions/`.
- No magic strings or numbers. Do not introduce any unless explicitly for the logger.
- Never use string union types (e.g., `"pass" | "fail"`). Use TypeScript Enums with the suffix `Type` (e.g., `StatusType`).
- Always use explicit boolean state checks (e.g., `response.isFail`). Never invert success booleans (e.g., `!response.isSuccess`).
- Code must be DRY. Reuse constants and wrappers.

## Phase 4: Chunked Delivery, Artifact Purge & File Moving

When a chunk of tasks is completed by the agents, do the following before starting the next loop iteration:

1. Use `mv` to move the completed task files from `.ai-memory/plans/pending/` to `.ai-memory/plans/completed/`.
2. Open the moved files and change `Status: pending` to `Status: completed`.
3. Update `.ai-memory/plans/readme.md` to reflect the new file locations.
4. Artifact sanitizer: Audit staged files. Purge unapproved artifact zip archives, temporary scratch files, or test outputs before committing.
5. TOTAL BAN on Routine Test Running & Build Checking: All test running (`go test`, `pytest`, `06-cicd-local-runner.py`) and build checking (`go build`) are strictly banned during routine task execution. Verification is checked later on in CI/CD. Commit code with a clear descriptive message. Never rewrite published git history (no force push, no rebasing, no squash). Push to git cleanly without failure.

## Phase 5: Output Window Stats (Mandatory Every Loop)

Every time you return a response or complete a loop iteration, explicitly output the following statistics:

- Tasks Done (This Chunk): [Number of tasks completed]
- Total Completed: [Total number of tasks in `.ai-memory/plans/completed/`]
- Total Pending: [Number of tasks remaining in `.ai-memory/plans/pending/`]
- Remaining Tasks List: [List the specific filenames/slugs of the tasks remaining]

## Execution Reporting (Mandatory Output Format)

1. Start of Run (Initial Task Breakdown Gate - Clearly Indented Markdown & Understanding Check):
   Before writing any code, scanning, or executing tasks, extract all pending deliverables as discrete, ordered items (`Task-01`, `Task-02`).
   - Respect whatever requirements the user has given, parse every request completely, and format each task clearly with proper markdown indentation, vertical blank lines, task state (`State: [PENDING]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).
   - TOTAL BAN ON UNFORMATTED RUN-ON TEXT: Never concatenate tasks into a single unformatted line or paragraph block (e.g. NEVER `#1. Task-01: ... #2. Task-02: ...`). Every task must be its own clearly separated markdown item.
   - Line-by-Line Output Format Structure:
     - Line 1: Header `### 📋 Confirmed Task Breakdown & Requirement Ingestion`
     - Line 2: Empty blank line
     - Line 3: Numbered task title `1. **Task-01: [Descriptive Task Title]**`
     - Line 4: Indented state bullet (3 spaces) `   - **State:** [PENDING]`
     - Line 5: Indented understanding check (3 spaces) `   - **Understood:** [YES] — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]`
     - Line 6: Indented actionable scope bullet (3 spaces) `   - **Actionable Scope:** [Precise technical deliverable and implementation scope]`
     - Line 7: Indented target files bullet (3 spaces) `   - **Target Files / Area:** [relative/path/or/module]`
     - Line 8: Empty blank line (vertical gap before next task)
     - Concluding Line: `Proceeding directly to execution.`

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[PENDING]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[PENDING]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to execution.
```

2. End-of-Turn Verification & Confidence Reporting: When all tasks are completed (or if the run concludes), you MUST output:

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

3. Lean Subtasks (No Common Boilerplate):
Subtasks must focus purely on unique task deliverables without repeating common repository boilerplate, universal coding rules, banned operations, or generic guidelines.

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] NO TEST RUNNING (TOTAL BAN): NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): NEVER query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents MUST query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait/sleep based on `etaSeconds` to eliminate credit waste.
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpf`, `cpb`, `cpr`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.

## Compliance Checklist (must follow non negociable)

- [x] Coding Guidelines enforced `02-spec/02-coding-guidelines/` and follow explicitly every step `.ai-memory/coding-guidelines.md`.
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

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline. Do not skip steps.

```text
N = 100  (Total self-loop steps budget)
PHASE_1_STEPS = N / 2   (Planning & Subtask Generation)
PHASE_2_STEPS = N / 2   (Parallel Execution & QA)
```

### Phase 1: Planning Mode & Subtask Generation FIRST (Steps 1 .. N/2)

1. Verbatim Prompt Capture & Task Extraction (First Action): Directly write the user's prompt verbatim into the planning spec at `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section. Extract whatever requirements the user has given into discrete actionable items (`Task-01`, `Task-02`) under `## Extracted Actionable Task List`. The AI MUST output this extracted checklist directly in chat before taking further actions or writing specs, using proper markdown indentation, vertical blank lines, task state (`State: [PENDING]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).
   - TOTAL BAN ON UNFORMATTED RUN-ON TEXT: Never concatenate tasks into a single unformatted line or paragraph block (e.g. NEVER `#1. Task-01: ... #2. Task-02: ...`).
   - Line-by-Line Output Format:
     - Line 1: Header `### 📋 Confirmed Task Breakdown & Requirement Ingestion`
     - Line 2: Empty blank line
     - Line 3: Numbered task title `1. **Task-01: [Descriptive Task Title]**`
     - Line 4: Indented state bullet (3 spaces) `   - **State:** [PENDING]`
     - Line 5: Indented understanding check (3 spaces) `   - **Understood:** [YES] — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]`
     - Line 6: Indented actionable scope bullet (3 spaces) `   - **Actionable Scope:** [Precise technical deliverable and implementation scope]`
     - Line 7: Indented target files bullet (3 spaces) `   - **Target Files / Area:** [relative/path/or/module]`
     - Line 8: Empty blank line (vertical gap before next task)
     - Concluding Line: `Proceeding directly to execution.`

   ```markdown
   ### 📋 Confirmed Task Breakdown & Requirement Ingestion

   1. **Task-01: [Descriptive Task Title]**
      - **State:** `[PENDING]`
      - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
      - **Actionable Scope:** [Precise technical deliverable and implementation scope]
      - **Target Files / Area:** `[relative/path/or/module]`

   2. **Task-02: [Descriptive Task Title]**
      - **State:** `[PENDING]`
      - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
      - **Actionable Scope:** [Precise technical deliverable and implementation scope]
      - **Target Files / Area:** `[relative/path/or/module]`

   Proceeding directly to execution.
   ```
2. Scan & Discover: Use the `invoke_subagent` tool to spawn exactly 2 planning subagents. Their role is to deeply scan the codebase for target changes.
3. Canonical Spec & Execution Plan: First, write the canonical application specification into `02-spec/21-app/xx-<slug>.md` (or directory `02-spec/21-app/xx-<slug>/` for complex features) with lossless verbatim prompt capture, and register it in `02-spec/21-app/readme.md`. Then initialize the master execution plan in `.ai-memory/plans/pending/xx-<slug>.md` linking back to the canonical spec. Write down 3–5 custom rules or constraints unique to this task inside the spec file.
4. Lean Subtask Decomposition: Break down the master plan into granular, single-responsibility subtask files in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`. Do not write common repository boilerplate, universal coding rules, banned operations, or generic guidelines inside subtask files.
   Subtasks MUST follow this lean template to prevent bloat:
   ```markdown
   # Subtask: [Name]
   Traceability ID: Task-01
   Spec Reference: [02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)
   Target Files: [Relative paths]
   Action: [Exact code changes required]
   Acceptance Criteria: [2-4 specific testable conditions]
   Targeted Verification: [Specific file-level linter command or exit 0 check]
   ```
5. UNCONDITIONAL ZERO-QUESTION EXECUTION MANDATE (TOTAL BAN ON STOPPING AFTER SPEC WRITING): Spec writing and subtask generation MUST strictly complete within the first 50% of the budget. As soon as Phase 1 planning completes, the master orchestrator MUST NOT STOP, pause, or ask the user "Should I proceed?", "Would you like me to start execution?", or any confirmation question. There is NO question. It MUST immediately, unconditionally self-loop and transition directly into Phase 2 execution mode. Generating specs without executing code changes is an INCOMPLETE FAILURE. The remaining 50% of the budget is dedicated strictly to modifying code, running targeted quality linters, consolidating subtasks, and completing the deliverables.

### Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. Parallel Dispatch: Use the `invoke_subagent` tool to spawn exactly 2 execution subagents (max 2 threads each) assigned to disjoint subtasks from `.ai-memory/plans/subtasks/xx-<slug>/`. Provide subagents with minimal instructions (e.g., "Read `.ai-memory/plans/subtasks/xx-slug/01-task.md` and execute it").
2. Execution & Coding Guidelines: Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, Unix LF line endings).
3. Failure Memory & Error Recovery: If a subagent fails, record the failure log in `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`; subsequent agents MUST read the failure log first to remediate root causes.
4. Local Verification: Run targeted linters on modified files ensuring `exit 0`.
5. TOTAL BAN on Routine Test Running & Build Checking: DO NOT run tests (`go test`, `pytest`, python runners) or check builds (`go build`, compiler checks) during routine execution. All test execution and build checks are deferred to CI/CD pipelines.
6. Atomic Change Tracking: Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json` for subsequent CI/CD verification.

### Phase 3: Task Consolidation & File Reduction (End of Loop)

> [!CRITICAL]
> To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

1. Combine all the completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header explicitly referencing how the main task started, referencing the canonical spec `[02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)`, and documenting exactly how many steps/loops it took. Note: The canonical specification in `02-spec/21-app/` remains permanently intact in the repository as the architectural source of truth; do not delete it during consolidation.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
6. Final Step Git Commit & Push (MANDATORY): Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, do NOT run builds or tests during routine turns (build and test verification deferred to CI/CD), group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, use GitMap semantic commit commands: `gitmap cpf "<summary>"` (features), `gitmap cpb "<summary>"` (bugs), or `gitmap cpr "<summary>"` (releases) which automatically stage, commit with standardized prefixes, and push directly to the remote repository. (Fallback to `git add -A && git commit && git push` only if GitMap CLI is unavailable). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` or `gitmap cpf` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.
