# [V2] Parent Task N-Step Continuous Loop & Multi-Agent Orchestration — Workflow (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

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
4. **Atomic Change Tracking:** Append all modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.lovable/test-inventory.json`.
7. **Local Verification:** Run targeted linters on modified files and ensure code compiles / passes lint checks with exit code 0 (`exit 0`). DO NOT run the full CI/CD runner (`06-cicd-local-runner.py`) during routine task steps.

### Phase 3: Task Consolidation & File Reduction (End of Loop)

> **CRITICAL:** To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

1. Combine all the completed granular subtasks from `.lovable/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.lovable/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header explicitly referencing how the main task started and documenting exactly how many steps/loops it took.
3. Delete the original granular `.md` files in `.lovable/plans/subtasks/xx-<slug>/`.
4. Delete the original parent plan `.lovable/plans/pending/xx-<slug>.md`.
5. Update `.lovable/plans/01-index.md` to point to the newly consolidated completed file.

---

## 1. AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** Scanned and learned `03-ai-scripts/01-index.md` before writing temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts executed strictly within the codebase repository root.
- [ ] **Strict .lovable/ Folder Storage:** All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] **Native File Manipulator:** Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] **Go Generate Sync:** If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 2. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

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

## 3. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.lovable/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed to update all callers of modified symbols.
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete with local CI green.
