# Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 150
```

N = total self-loop steps budget that the agents will perform.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Planning & Spec Generation, Steps 1..N/2): Spawn exactly 2 planning subagents (max 2 threads each) to scan the codebase and draft `.lovable/plans/pending/xx-<slug>.md`.
2. [ ] /goal Phase 1 (Subtask Decomposition): Decompose the plan into a lean set of actionable subtasks. Do not over-prompt or generate excessive markdown files. Keep subtasks focused purely on execution and the domain task itself in `.lovable/plans/subtasks/xx-<slug>/*.md`.
3. [ ] /goal Phase 1 (Strict Folder Bounding): Restrict all planning logs, active locks, and status reports strictly within `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`).
4. [ ] /goal Phase 1 (Zero-Stop Transition): Immediately upon completing Phase 1, self-loop and transition directly into Phase 2 execution mode without pausing or stopping.
5. [ ] /goal Phase 2 (Execution & Code Refactoring, Steps N/2+1..N): Spawn exactly 2 execution subagents (max 2 threads each) to execute subtasks on disjoint files in parallel.
6. [ ] /goal Phase 2 (Failure Memory & Error Recovery): If a subagent fails, record the failure log in `.lovable/plan.md` and `.lovable/memory/issues/`; subsequent agents MUST read the failure log first to remediate root causes.
7. [ ] /goal Phase 2 (Quality Gate Verification): Execute local linters and `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring `exit 0` before finishing (test execution is disabled unless explicitly commanded by the repository owner).
8. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
9. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
10. [ ] /learn Ingest `spec/02-coding-guidelines/` for domain-specific architectural specifications.
11. [ ] /learn Ingest `spec/03-error-manage/` for error handling architectures and AppError.
12. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
13. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: 2-Agent Planning & Subtask Generation in .lovable/plans/)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: 2-Agent Parallel Execution, Self-Looping, CI Quality Gates)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, you can rely on progressive disclosure for future runs. Do not keep the entire prompt in your active memory if you don't need it.

---

## 1. 2-Agent Concurrency & Ruthless Orchestration

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- **Strict 2-Agent Limit (Max 2 Threads Each):** When dispatching work in Phase 1 (planning) or Phase 2 (execution), you MUST spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.lovable/`):** Subagents are strictly restricted to writing planning files, subtasks, status reports, and logs inside `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`, `.lovable/memory/issues/`).
- **Context Diet:** When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/xx-slug/01-<subtask-title>.md` and execute it"). The subagent MUST read the necessary files itself.
- **Fail Fast & Kill Stalls:** If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.

---

## 2. Phase 1: Planning Mode & Subtask Generation FIRST (Steps 1 .. N/2)

Before writing any source code changes, you MUST execute Phase 1:

1. **Scan & Discover:** Spawn 2 planning subagents to deeply scan the codebase for target changes or violations.
2. **Master Spec Generation:** Save the master architectural plan into `.lovable/plans/pending/xx-<slug>.md`.
3. **Task-Specific Rule Set:** Write down 3–5 custom rules or constraints unique to this task inside the spec file.
4. **Lean Subtask Decomposition:** Break down the plan into a few highly focused subtask files in `.lovable/plans/subtasks/xx-<slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc. **Task Focus Over Meta-Prompting:** Your goal is to write code and solve the problem, not just generate more AI prompts. Subagent instructions should clearly define the domain task itself.
5. **Strict Relative Git Paths:** All markdown links and file paths in subtasks MUST be strictly relative to the repository root (e.g. `.lovable/spec/...`, `src/...`). Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
6. **MANDATORY AUTO-LOOP (DO NOT STOP):** As soon as Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

---

## 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. **Parallel Dispatch:** Spawn 2 execution subagents (max 2 threads each) assigned to disjoint subtasks from `.lovable/plans/subtasks/xx-<slug>/`.
2. **File Locking:** Verify subagents operate on distinct files using `.lovable/01-index.md`.
3. **Execution & Coding Guidelines:** Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
4. **Failure Memory & Feedback Loop:** If a subagent fails:
   - Rollback dirty changes and write the failure error log to `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`.
   - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
5. **Progress & Completion:** Move completed subtasks to `.lovable/plans/completed/` and update `.lovable/plans/01-index.md`.
6. **Temp & Failure Directory Isolation:** All temporary files, test outputs, and runner caches MUST be isolated within `.lovable/temp/`. Creating `.tmp/` at root is strictly prohibited.
   - Dedicated Failure Directory: `.lovable/temp/failures/` is the dedicated folder where failed tests and failed quality gates write error logs (`<test-or-job-name>.log`).
   - Passing Tests Completely Silent: Passing tests must produce ZERO filesystem artifacts and remain completely silent in output logs.
7. **Local Verification:** Run targeted linters on modified files and ensure code compiles / passes lint checks with exit code 0 (`exit 0`). DO NOT run the full CI/CD runner (`06-cicd-local-runner.py`) during routine task steps unless explicitly commanded by repository owner.
8. **Runner In-Flight ETA Wait Protocol:** When running background commands, the runner writes live status and remaining ETA to `.lovable/temp/runner-eta.json` (emitting in-flight heartbeats strictly every 25 seconds or more). If an agent inspects an active background job, it MUST sleep/wait for **1 minute (60 seconds) each time**, or dynamically sleep for the remaining ETA duration read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling.
9. **Centralized Test Inventory & Incremental Caching:** All unit tests are cataloged in `.lovable/test-inventory.json` with strictly repository-relative paths (`target_file`, `test_file`). First run executes all tests to establish baseline timings; subsequent runs execute incrementally only if target code files or test files change. Slow test threshold defaults to `4.0s` (configurable via `GITMAP_SLOW_TEST_THRESHOLD`).
10. **Dual-Queue Worker Pools:** Slow tests run in a dedicated 4-worker pool running at most 2 tests at a time per batch. Fast tests run in a 4-worker pool running at most 4 tests at a time, pulling in chunks of 100 tests from the test inventory queue until all are complete.

---

## 4. AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** Scanned and learned `03-ai-scripts/01-index.md` before writing temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts executed strictly within the codebase repository root.
- [ ] **Strict .lovable/ Folder Storage:** All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] **Native File Manipulator:** Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] **Go Generate Sync:** If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 5. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: Fully enforced every file in `spec/02-coding-guidelines/` and `.lovable/coding-guidelines.md`.
- [ ] Error Management: Enforced `spec/03-error-manage/` using domain-specific `AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with is or has ONLY (all other prefixes like can, should, was, will, did, must are banned). NO negatives (`!isSuccess` is banned; use `isFail`).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`). Behavior-driven unit test names.
- [ ] Multi-Line Arguments (Rule 9a/9b): Signatures and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths: Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.

---

## 6. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.lovable/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed to update all callers of modified symbols.
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete with local CI green.

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.lovable/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.lovable/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`, `.lovable/memory/issues/`).
- **Context Diet & Task Focus:** Provide subagents with clear, lean instructions that focus on the actual domain task itself. Do not write massive meta-prompts or generate excessive boilerplate markdown. Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.lovable/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- **Lean Subtask Decomposition:** Break down the plan into a few highly focused subtask files in `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc. **Task Focus Over Meta-Prompting:** Your goal is to write code and solve the problem, not just generate more AI prompts. Subagent instructions should clearly define the domain task itself.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.lovable/plans/subtasks/` to `.lovable/plans/completed/` and update `.lovable/plans/01-index.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.
- Isolate all temporary test files, caches, and scratch directories within `.lovable/temp/`. Never create `.tmp/` at root. Failed tests/gates write to `.lovable/temp/failures/`; passing tests remain completely silent and produce zero disk files.
- In-Flight ETA Wait Protocol: When checking background test runners, agents MUST sleep/wait for 1 minute (60s) each time, or dynamically sleep for the remaining ETA duration read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling.
- Record all modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).

## Task Consolidation & File Reduction (End of Loop)

> **CRITICAL:** To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

When all subtasks for a parent task (`.lovable/plans/pending/xx-<slug>.md`) are finished, execute this final cleanup step before ending the run:
1. Combine all the completed granular subtasks from `.lovable/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.lovable/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header that explicitly references how the main task started and documents exactly how many steps/loops it took to complete.
3. Delete the original granular `.md` files in `.lovable/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.lovable/plans/pending/xx-<slug>.md`.
5. Update `.lovable/plans/01-index.md` to point to the newly consolidated completed file.
