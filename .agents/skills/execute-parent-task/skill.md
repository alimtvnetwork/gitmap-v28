---
name: execute-parent-task
description: Autonomously orchestrate and execute parent tasks by decomposing them into subtasks and running continuous N-step self-loops with strict coding guidelines and error management.
---

# Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)

> [!IMPORTANT]
> Prompt Version: 2.4.0
> Synchronization: Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 150
```

N = total self-loop steps budget that the agents will perform.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1A (Step 0 - Verbatim Prompt Recording & Task Extraction Gate): Immediately capture the user prompt verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under `## User Request (Verbatim)`. If screenshot URLs or base64 data URIs are provided, decode/save them as image files (`assets/screenshots/<task-slug>-<NN>.png`) and refer back to them via relative paths in specs. Extract actionable bullet-point deliverables with traceable IDs (`Task-01`, `Task-02`), and output this confirmed task list directly in chat formatted as `#1. Task-01: <details>` without hard brackets before any file exploration or scanning.
2. [ ] /goal Phase 1B (Step 1 - Master Spec Generation): Write the master architectural plan in `.ai-memory/plans/pending/xx-<slug>.md`, documenting each task in a detailed manner with architectural specifications, constraints, and blast radius.
3. [ ] /goal Phase 1B (Step 2 - Scan & Discover): Use fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) to inventory files and map call sites without tool truncation limits.
4. [ ] /goal Phase 1B (Step 3 - Lean Subtask Decomposition): Decompose the plan into lean subtask files in `.ai-memory/plans/subtasks/xx-<slug>/01-<subslug>.md`. Subtasks must focus purely on unique task deliverables without repeating common repository boilerplate.
5. [ ] /goal Phase 1B (Step 4 - Readiness Audit Gate): Confirm all `Task-xx` deliverables are mapped to subtasks and disjoint files before execution.
6. [ ] /goal Phase 1B (Zero-Stop Transition): Immediately upon completing Phase 1, self-loop and transition directly into Phase 2 execution mode without pausing or asking for permission.
7. [ ] /goal Phase 2 (Execution & Code Refactoring, Steps N/2+1..N): Spawn at most 2 execution subagents (max 2 threads each) to execute subtasks on disjoint files in parallel.
8. [ ] /goal Phase 2 (Failure Memory & Error Recovery): If a subagent fails, record the failure log in `.ai-memory/plan.md` and `.ai-memory/memory/issues/`; subsequent agents must read the failure log first to remediate root causes.
9. [ ] /goal Phase 2 (Change Recording & Quality Linting): Record all modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`) and run targeted file-level linters on specifically modified files (`exit 0`). Do not run `06-cicd-local-runner.py`, unit tests, or build checks (deferred to CI/CD).
10. [ ] /goal Phase 3 (Consolidation & Atomic Push): Consolidate completed subtasks into `.ai-memory/plans/completed/xx-<slug>.md`, delete granular subtasks and pending plan, stage all changes, and push in a single grouped commit.
11. [ ] /goal Phase 3 (Completion & Confidence Reporting): Emit the final Task Completion Summary with green check mark emojis, modified files summary, and implementation confidence score.
12. [ ] /learn Ingest `.ai-memory/memory/01-index.md` for project memory index and past learnings.
13. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
15. [ ] /learn Ingest `02-spec/03-error-manage/` for error handling architectures and AppError.
16. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/<slug>/SKILL.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `SKILL.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## 1. 2-Agent Concurrency & Ruthless Orchestration

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- Strict 2-Agent Limit (Max 2 Threads Each): When dispatching work in Phase 1 (planning) or Phase 2 (execution), spawn at most 2 sub-agents concurrently, with no more than 2 threads per agent.
- Strict Folder Bounding (`.ai-memory/`): Subagents are strictly restricted to writing planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/01-index.md`, `.ai-memory/memory/issues/`).
- Context Diet: When spawning a subagent, do not paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.ai-memory/plans/subtasks/xx-slug/01-<subtask-title>.md` and execute it"). The subagent must read the necessary files itself.
- Fail Fast & Kill Stalls: If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.

---

## 2. Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, or code changes, you must execute Phase 1A:

1. Verbatim Prompt Capture: Directly write the user prompt verbatim into the planning spec at `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section.
2. Screenshot & Print Screen Base64 Image Ingestion Protocol: If the user request or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Convert that base64 encoding or downloaded image to the file system immediately, saving it as a persistent file under `assets/screenshots/<task-slug>-<NN>.png` or `assets/ui/<task-slug>-<NN>.png`.
   - Never leave raw, massive base64 strings or external ephemeral URLs in the prompt text, specs, or subtasks.
   - Refer back to this saved image file in the master spec, planning document, subtasks, and UI implementation notes strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<task-slug>-<NN>.png)`).
3. Actionable Deliverables Extraction: Break down the user prompt into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, `Task-03`) under `## Extracted Actionable Task List`. Do not use hard bracket notation like `[T-01]:`.
4. Mandatory Chat Output Gate: Output the confirmed deliverables list directly in chat before performing any scans or tool calls:

```markdown
### Confirmed Task Breakdown
#1. Task-01: [Actionable deliverable description]
#2. Task-02: [Actionable deliverable description]
#3. Task-03: [Actionable deliverable description]
Proceeding directly to Phase 1B: Spec & Subtask Generation.
```

No tool calls, codebase searches, or subagent spawning may proceed until this breakdown is emitted.

---

## 3. Phase 1B: Detailed Spec Generation & Lean Subtasks (Steps 1 .. N/2)

### Step 1: Master Spec Generation
Save the master architectural plan into `.ai-memory/plans/pending/xx-<slug>.md`. Provide a detailed breakdown of each task:
- Architectural context, domain logic, and module interactions.
- Visual specification references: If screenshots were provided, embed the relative markdown links to the saved image files and specify visual layout, typography, and UX requirements.
- Input and output data contracts.
- 3 to 5 custom rules or constraints unique to this task domain.
- Blast radius analysis identifying all downstream callers.

### Step 2: Scan & Discover (Python Toolchain Acceleration)
To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, use the repository's dedicated Python discovery scripts:
- Inventory Target Files: `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- Fast Cached Grep (<15ms): `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50`
- Sub-Millisecond Folder & File Exploration: `python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --limit 50`
- Read Target File: `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- Fast Pattern Search: `python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --limit 50`
- Subsystem & Topology Overview: `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

### Step 3: Lean Subtask Decomposition (No Common Boilerplate in Subtasks)
Break down the master plan into granular subtask files in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`, `02-<subtask>.md`, etc.

Important Rule: Do not write common repository boilerplate, universal coding rules, banned operations, or generic guidelines inside subtask files. Common rules belong in the parent plan and root guidelines. Subtasks must contain only the unique, non-common items required for that specific subtask.

Subtasks must follow this lean, unique template:
```markdown
# Subtask [01]: [Descriptive Subtask Name]
Traceability ID: Task-01
Target Files: [Strict relative paths from repo root]
Action: [Exact code changes, functions, types, and logic to modify or add]
Acceptance Criteria: [2-4 specific testable conditions proving completion]
Targeted Verification: [Specific file-level linter command or exit 0 check]
```

### Step 4: Subtask Readiness Audit Gate
Before transitioning to execution, verify:
- Every extracted deliverable `Task-xx` has at least one corresponding subtask file.
- All subtask files are non-empty and specify disjoint target files.
- All file paths in subtasks use strict relative Git paths (zero absolute paths or `file:///` URIs).

### Step 5: Mandatory Auto-Loop (Do Not Stop)
As soon as Phase 1 planning completes, the master orchestrator must not stop or ask the user for permission. It must immediately self-loop and transition directly into Phase 2 execution mode.

---

## 4. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. Parallel Dispatch: Spawn at most 2 execution subagents (max 2 threads each) assigned to disjoint subtasks from `.ai-memory/plans/subtasks/xx-<slug>/`. Provide minimal instructions referencing the subtask path.
2. File Locking & Disjoint Files: Verify subagents operate on distinct files using `.ai-memory/01-index.md`.
3. Execution & Coding Guidelines: Subagents refactor code following all coding guidelines (<= 8-15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
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
3. Targeted Failure Diagnostics: Use GitMap automated error extraction to isolate actionable failure lines without pulling verbose passing logs.

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
2. In this single consolidated file, include a header that explicitly references how the main task started and documents exactly how many steps or loops it took to complete.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-<slug>.md`.
5. Update `.ai-memory/plans/01-index.md` to point to the newly consolidated completed file.
6. Final Step Git Commit & Push (Mandatory): Stage all modified files, consolidated plans, and memory records (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<type>(<scope>): <summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## 7. End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of all tasks and before concluding the turn, you must emit this structured verification summary in the chat response:

```markdown
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
```

---

## 8. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] NO TEST RUNNING (TOTAL BAN): Never run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): Never launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): Never bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents must query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait or sleep based on `etaSeconds` to eliminate credit waste.

---

## 9. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You must verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `AppError`, never generic error.
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

- [ ] MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.
