# [V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, you can rely on progressive disclosure for future runs. Do not keep the entire prompt in your active memory if you don't need it.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline. Do not skip steps.

```text
N = 100  (Total self-loop steps budget)
PHASE_1_STEPS = N / 2   (Planning & Subtask Generation)
PHASE_2_STEPS = N / 2   (Parallel Execution & QA)
```

### Phase 1: Planning Mode & Subtask Generation FIRST (Steps 1 .. N/2)

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips during planning, the planning subagents and master orchestrator MUST use the repository's dedicated Python discovery scripts first:

1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide files.

1. **Verbatim Prompt Capture & Task Extraction (First Action):** Directly write the user's prompt verbatim into the planning spec at `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section. Extract the specific task list from this prompt as actionable bullet points / checklist items under `## Extracted Actionable Task List`. The AI MUST output this extracted checklist directly in chat confirming: *"Confirmed Task Deliverables: 1. [task 1], 2. [task 2]..."* before taking further actions.
2. **Scan & Discover:** Use the `invoke_subagent` tool to spawn exactly 2 planning subagents. Their role is to deeply scan the codebase for target changes.
3. **Master Spec Generation:** Save the master architectural plan into `.ai-memory/plans/pending/xx-<slug>.md`. Write down 3–5 custom rules or constraints unique to this task inside the spec file.
4. **Lean Subtask Decomposition:** Break down the master plan into granular, single-responsibility subtask files in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`.
   *Subtasks MUST follow this lean template to prevent bloat:*
   > `# Subtask: [Name]`
   > `**Target Files:** [Relative paths]`
   > `**Action:** [Exact code changes required]`
   > `**Constraints:** [Key rules to follow]`
5. **MANDATORY AUTO-LOOP (DO NOT STOP):** As soon as Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. **Parallel Dispatch:** Use the `invoke_subagent` tool to spawn exactly 2 execution subagents (max 2 threads each) assigned to disjoint subtasks from `.ai-memory/plans/subtasks/xx-<slug>/`. Provide subagents with minimal instructions (e.g., "Read `.ai-memory/plans/subtasks/xx-slug/01-task.md` and execute it").
2. **Execution & Coding Guidelines:** Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, Unix LF line endings).
3. **Failure Memory & Error Recovery:** If a subagent fails, record the failure log in `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`; subsequent agents MUST read the failure log first to remediate root causes.
4. **Atomic Change Tracking:** Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json` for subsequent CI/CD verification.
5. **TOTAL BAN on Test Running:** DO NOT run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`), Go (`go test`), or any test runner during routine execution turns. All test execution is strictly deferred to CI/CD pipelines and dedicated fix workflows.
6. **TOTAL BAN on Build Checking:** DO NOT run build verification commands (`go build`, `npm run build`, compiler invocations). Build compilation is checked later on in CI/CD.
7. **Targeted Quality Linting Only:** Run only targeted, fast file-level linters/autofixers on specifically modified files (`exit 0`). DO NOT run `06-cicd-local-runner.py` or full test suites.

8. **Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI):**
   When monitoring or checking remote CI/CD pipelines (e.g., following git push or during pipeline audits):
   - **Mandatory GitMap Authority:** Use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`) to query pipeline state and parse `is_running`, `status`, `etaSeconds`, and `nextAiCommand`.
   - **Anti-Credit-Waste Waiting Mandate:** NEVER busy-poll or query rapidly (`gh run view` tight loops). When pipeline is running, wait/sleep based on `etaSeconds` using `gitmap pipeline-ai status -t <etaSeconds>` (ETA > 120s: wait 20s-30s; 60s < ETA <= 120s: wait 10s-20s; ETA <= 60s: wait 5s-10s).
   - **Targeted Error Diagnostics:** Use GitMap to isolate actionable failure lines (`##[error]`, `FAIL:`) without pulling verbose logs.

### Per-Task Agent Isolation & Workspace Subfolders (`.ai-memory/temp-agents/xx-<task-name>/`)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.ai-memory/temp-agents/xx-<task-name>/`:
1. **Per-Task Isolation:** On task start, the assigned subagent creates its isolated directory `.ai-memory/temp-agents/xx-<task-name>/`.
2. **State & Progress Tracking:** Create `.ai-memory/temp-agents/xx-<task-name>/state.md` documenting:
   - Task sequence and target deliverables.
   - Files assigned for modification.
   - Current subtask step and completion percentage.
3. **Inter-Agent Communication & Scratch Space:**
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.ai-memory/temp-agents/xx-<task-name>/`.
4. **On Error/Crash:** Append the exact error, root cause, and `STATUS: FAILED` to `.ai-memory/temp-agents/xx-<task-name>/state.md` before exiting.
5. **On Success:** Mark `STATUS: DONE` in `.ai-memory/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

### Phase 3: Task Consolidation & File Reduction (End of Loop)

> **CRITICAL:** To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

1. Combine all the completed granular subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.ai-memory/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header explicitly referencing how the main task started and documenting exactly how many steps/loops it took.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/01-index.md` to point to the newly consolidated completed file.
6. **Final Step Git Commit & Push (MANDATORY):** Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## 1. AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** Scanned and learned `03-ai-scripts/01-index.md` before writing temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts executed strictly within the codebase repository root.
- [ ] **Strict .ai-memory/ Folder Storage:** All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] **Native File Manipulator:** Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] **Go Generate Sync:** If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 2. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] **NO RAPID CI/CD POLLING (TOTAL BAN):** NEVER query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents MUST query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait/sleep based on `etaSeconds` to eliminate credit waste.

---

## 3. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with is or has ONLY (all other prefixes like can, should, was, will, did, must are banned). NO negatives (`!isSuccess` is banned; use `isFail`).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`). Behavior-driven unit test names.
- [ ] Multi-Line Arguments (Rule 9a/9b): Signatures and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths: Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.

---

## 4. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.ai-memory/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed to update all callers of modified symbols.
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete without running banned test/build commands (all testing and build verification deferred to CI/CD).
- [ ] Final Step Commit & Push Verified: Staged all changes (`git add -A`), committed everything in a single grouped atomic commit, and pushed to git before ending the turn (no per-file commits).

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!
