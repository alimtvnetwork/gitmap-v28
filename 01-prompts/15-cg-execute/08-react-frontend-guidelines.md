# React & Frontend Architecture — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-react`, `cg-execute react`, `audit react`, `fix frontend guidelines`, `enforce react standards`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all React and frontend architecture violations across the codebase, decomposing oversized components (<= 80–100 lines), enforcing 8–15 line function caps, converting custom hook tuple returns to named objects, and eliminating redundant `useEffect` hooks until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase to inventory all architectural violations and anti-patterns.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` with exit code 0 (`exit 0`).
9. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `spec/02-coding-guidelines/02-typescript/` for domain-specific architectural specifications.
15. [ ] /learn Ingest `spec/07-design-system/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase, Write .lovable/plans/pending/ Spec, Create .lovable/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Component Refactoring, Linter Verification, Local CI Runner Verification, Plan Completion)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Canonical Sizing & Nesting Hierarchy

- **React Components:** Recommended <= 80 lines; standard max <= 100 lines.
- **Functions & Hooks:** Preferred <= 8 lines; hard cap <= 15 lines.
- **Nested `if` Statements:** Zero tolerance (flatten with guard clauses and early returns).
- **Standard File Sizing:** Max 100 coding lines per file (recommended <= 80 lines).

---

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents strict, microscopic instructions.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet: When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/XX-react-frontend/01-<subtask-title>.md` and execute it"). The subagent MUST read the necessary files itself.

---

## 2. Phase 1: Scan Codebase & Write Implementation Spec First (Steps 1 to PHASE_1_STEPS)

Before modifying application code, you MUST thoroughly scan the repository and write an actionable execution spec.

- **Actionable Scan:** Use search/grep and line-count tools across all `.tsx` and `.jsx` files to identify:
  1. Component files exceeding 80–100 lines of code.
  2. Functions and hooks exceeding 8 lines (hard cap 15 lines).
  3. Nested `if` conditionals (nesting depth > 1).
  4. Custom hooks returning array literals (tuples) instead of named objects (`return [val, setVal]` vs `return { val, onUpdate }`).
  5. Redundant `useEffect` hooks used for derived calculations or local state sync.
  6. Direct state mutations (`state.items.push()` or `state.prop = x`).
  7. Missing Mermaid component tree diagrams in UI documentation.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-react-frontend-audit.md` listing every affected file, exact line numbers, and the extraction plan.
- **Create a Task-Specific Rule Set:** Analyze the specific domain and write 3-5 custom rules inside the spec file.
- **Subtasks:** Break the plan down into granular subtask files inside `.lovable/plans/subtasks/XX-react-frontend/` (e.g. `01-decompose-oversized-components.md`, `02-hook-named-objects.md`).

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/02-canonical-size-tier.md`**
  - **Why:** Universal size limits across all languages.
  - **How:** Components <= 80 lines recommended (max 100 lines). Functions <= 8 lines preferred (max 15 lines). Max 100 lines for any file. Zero line-compression cheating.
- [ ] **`spec/02-coding-guidelines/06-ai-optimization/01-index.md`**
  - **Why:** Comprehensive catalog of forbidden vs required generation patterns.
  - **How:** Strictly follow AH-N1 to AH-T2 rules. Zero ghost diffs, zero truncation stubs (`// ...`), zero unverified claims.
- [ ] **`spec/02-coding-guidelines/06-ai-optimization/06-citation-requirement.md`**
  - **Why:** Grounded rule enforcement and traceability.
  - **How:** Cite authoritative spec files for every code modification made.
- [ ] **`spec/02-coding-guidelines/01-cross-language/04-code-style/02-braces-and-nesting.md`**
  - **Why:** Zero tolerance for nested conditionals.
  - **How:** Flatten all nested `if` statements with guard clauses and early returns.
- [ ] **`spec/02-coding-guidelines/02-typescript/09-typescript-standards-reference.md`**
  - **Why:** Authoritative TypeScript standards.
  - **How:** Strict typing, zero `any`, explicit function return types, and discriminated unions for polymorphic data.
- [ ] **`spec/02-coding-guidelines/02-typescript/15-state-management.md`**
  - **Why:** Immutable state & hook conventions.
  - **How:** Custom React hooks MUST NOT return tuples (`[state, setState]` is banned). Custom hooks MUST return named property objects (`{ userProfile, isLoading, onUpdate }`).
- [ ] **`spec/02-coding-guidelines/02-typescript/10-promise-await-patterns.md`**
  - **Why:** Async lifecycle safety in React.
  - **How:** Handle all async flows with `try/catch` and `AppError` wrappers; avoid unhandled promises in `useEffect`.
- [ ] **`spec/07-design-system/02-design-principles.md`**
  - **Why:** Modular UI architecture & sizing caps.
  - **How:** Hard 100-line cap per React component file (target <= 80 lines). Decompose large UI blocks into single-responsibility child components.
- [ ] **`spec/07-design-system/03-theme-variable-architecture.md`**
  - **Why:** Theme token usage.
  - **How:** Zero hardcoded hex color codes in components. Consume CSS theme custom properties and semantic tokens.
- [ ] **`spec/03-error-manage/02-error-architecture/04-error-modal/02-react-components/03-error-store.md`**
  - **Why:** Centralized UI error modal.
  - **How:** Route all unhandled UI exceptions and API error responses to the global error store; display via standard error modal.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-enum-and-boolean.mjs` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If no dedicated frontend linter exists, create `linter-scripts/check-enum-and-boolean.mjs` that AST-scans `.tsx` and `.jsx` files for:
  1. Component file line counts exceeding 80–100 lines.
  2. Functions and hooks exceeding 8–15 lines.
  3. Nested `if` statements.
  4. Single-line `if/else` violations.
  5. Custom hooks returning array literals (tuples) instead of object literals.
  6. Redundant `useEffect` calls syncing derived state.
  7. Direct state mutations (`state.items.push()` or `state.prop = x`).
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  node linter-scripts/check-enum-and-boolean.mjs
  # And verify standard ESLint:
  npm run lint
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `03-ai-scripts/06-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:frontend"] = ["node", "linter-scripts/check-enum-and-boolean.mjs"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running the frontend linter.

---

## 5. Phase 2: Active Code Refactoring & Autonomous Fix Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Open the offending component files and directly rewrite the code to eliminate violations. Maintain continuous self-looping until all checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-react-frontend/
    2. Open and modify the actual source code files:
       - Decompose oversized components into modular sub-components <= 80–100 lines.
       - Break functions/hooks > 8 lines into concise helpers (max 15 lines).
       - Flatten nested if statements with guard clauses.
       - NEVER collapse if/else onto a single line to cheat line caps.
       - Convert custom hook return tuples to named property objects and update callers.
       - Eliminate redundant useEffect by computing derived state inline.
       - Enforce immutable state updates across all handlers.
    3. Run the frontend linter:
          node linter-scripts/check-enum-and-boolean.mjs
    4. Run frontend type checks and linters (unit test execution is disabled unless explicitly commanded by the repository owner).
    5. Run the local CI runner:
          python 03-ai-scripts/06-cicd-local-runner.py --no-tests
    6. IF any check fails:
          - Diagnose failure, fix component code directly, and re-run immediately.
       IF all checks pass (exit code 0):
          - Mark subtask completed and proceed to next subtask.

    7. When all subtasks are finished and local CI is 100% green:
          - Move .lovable/plans/pending/XX-react-frontend-audit.md to .lovable/plans/completed/
          - Update .lovable/plans/01-index.md
          - Run plan consolidator:
            python 03-ai-scripts/20-plan-consolidator.py
          - Stage modified files with git add and create semantic commit:
            git commit -m "refactor(ui): decompose component hierarchy, flatten nested ifs, and enforce named hook objects"
          - BREAK and finish turn.
```

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.lovable/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.lovable/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`, `.lovable/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.lovable/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.lovable/plans/subtasks/` to `.lovable/plans/completed/` and update `.lovable/plans/01-index.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute local linters and `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring `exit 0` before concluding.

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/01-index.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.lovable/plans/01-index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] All `.tsx`/`.jsx` component files are <= 80–100 lines.
- [ ] All functions/hooks are <= 8 lines preferred, hard cap 15 lines.
- [ ] Zero Nested Ifs: Flattened with guard clauses.
- [ ] NO Line-Compression Cheating: No single-line `if/else`, no deleted blank lines (R13-R16).
- [ ] Custom hooks return named objects `{ data, isLoading }`, not tuples.
- [ ] No redundant `useEffect` for derived computations.
- [ ] `node linter-scripts/check-enum-and-boolean.mjs` exited with code 0.
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced every file in `spec/02-coding-guidelines/02-typescript/` and `.lovable/coding-guidelines.md`.
- [ ] Component Sizing: All component files are <= 80–100 lines (hard cap 100 lines).
- [ ] Function Sizing: All functions <= 8–15 lines.
- [ ] Zero Nested Ifs: Flat flow with early returns.
- [ ] Anti-Compression: Zero single-line `if/else` or compressed whitespace tricks.
- [ ] Hook Return Signatures: Named property objects only, zero tuple returns.
- [ ] Semantic Naming: Absolutely NO generic garbage names (`temp`, `data`, `obj`, `comp_100`).

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] Echo Back the Spec: I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] Exhaustive Violation Ledger: I have maintained an exact markdown table ledger in `.lovable/plans/pending/` tracking every single violation `| Id | File | Line | Snippet | Planned Fix | Status |` and reconciled every item.
- [ ] Pre-Commit Diff Proof (Disk Reality Check): I have executed `git status --porcelain` and `git diff --stat` and verified that every file I claim to have modified is actually listed as modified in the terminal output before committing.
- [ ] Zero Truncation / No Placeholder Search: I ran a regex search for `TODO`, `FIXME`, `\[.*\]`, `// ...`, and `/* ... */` in my modified files and confirmed I left zero placeholders or truncated stubs behind. I actually wrote the complete implementation.
- [ ] Verifiable Tool Execution: I did not fabricate test/linter passes. I executed the actual linter script and test runner via tool calls and captured `exit code 0`.
- [ ] Spec Citation Grounding: Every refactoring action cites the exact authoritative rule in `spec/` (e.g. `spec/02-coding-guidelines/01-cross-language/01-index.md`).
- [ ] Index Sync Deadman Switch: I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.lovable/what-to-read.md`. I did not leave any orphaned files.
- [ ] Blast Radius Acknowledgment: Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.

---

## Strictly Avoid: No Automatic Releases & No Test Running (Strict Policy)

- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **CI/CD Quality Gates:** Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring all quality gates exit with code 0 (`exit 0`).
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and linters with test execution disabled (--no-tests) unless explicitly commanded by the repository owner, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## STRICT AVOIDANCE: Anti-Compression & Formatting Integrity (No Cheating)

> [!CAUTION]
> **TOTAL BAN ON LINE-COMPRESSION CHEATING:**
> When enforcing file size (<= 100 coding lines) and function size (8–15 lines), AI agents frequently attempt to "cheat" the line counter by destroying formatting. This is strictly forbidden and results in immediate rejection.

- **NO Single-Line If/Else:** NEVER collapse `if/else`, return statements, or blocks into a single line (e.g. `if (x) return y;` or `if (x) { y(); }` are strictly forbidden). Every statement requires its own line and curly braces.
- **NO Deleting Required Blank Lines (R13-R16):** NEVER delete blank lines before `return`/`throw` or after closing `}` to artificially reduce file size.
- **NO Stripping Types or Comments:** NEVER remove TypeScript types, docstrings, or clean indentation to cram code into fewer lines.
- **Mandatory Solution:** The ONLY acceptable way to satisfy line limits is **legitimate modular decomposition** — extracting helper functions into separate files and breaking large components into child components.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

---

## Metadata

- slug: cg-react-frontend
- priority: medium
- status: active
