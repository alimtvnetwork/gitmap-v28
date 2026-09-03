# Instruction (must follow): Execute Coding Guidelines — Error Management & Architecture

Trigger Keywords & Aliases: `cg-error`, `cg-execute error`, `audit error`, `fix error guidelines`, `enforce error management`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute error management compliance across the entire repository by decomposing violations into subtasks, verifying/creating error linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all active files, listing all error management spec files with why and how, creating the Antigravity skill, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, refactoring error handling to follow `AppError` context wrapping and universal envelopes, running the error linter, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action and also create agent rules in the repo if required to or missing from rules set of agent memory.
- [ ] /learn `.lovable/coding-guidelines/coding-guidelines.md` and it is must and /goal apply the guidelines in coding every aspect.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, AppError Refactoring, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-error-management/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this error management prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-error-management
   description: >-
     Autonomously audits, refactors, and validates repository-wide error management against spec/03-error-manage/ using AppError wrappers, universal response envelopes, and CI linters.
   ---
   ```
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents strict, microscopic instructions.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet: When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/XX-error-management/01-task.md` and execute it"). The subagent MUST read the necessary files itself.

---

## 2. Phase 1: Write the Implementation Spec & Subtasks FIRST (Steps 1 to PHASE_1_STEPS)

Before doing anything else, you MUST write a highly detailed execution spec.

- **What to write:** Break down the parent task into a detailed architectural plan, complete error violation inventory, code review guides, and embedded error management guidelines.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-error-management-audit.md`. Do not hallucinate folders.
- **Create a Task-Specific Rule Set:** Before executing, analyze the specific task domain and explicitly write down 3-5 custom rules or constraints unique to this task inside the spec file. This prevents domain-specific regressions and forces sub-agents to follow exact architectures.
- **Subtasks:** You MUST break the plan down and create detailed subtask files inside `.lovable/plans/subtasks/XX-error-management/`. Every subtask file must contain actionable, microscopic instructions.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/03-error-manage/00-overview.md`**
  - **Why:** Authoritative error management foundation across all services.
  - **How:** Never swallow errors; every `catch` logs with operation name and key inputs, then rethrows or returns a typed error.
- [ ] **`spec/03-error-manage/02-error-architecture/01-error-handling-reference.md`**
  - **Why:** Universal cross-language `AppError` and `AppException` structure.
  - **How:** Implement `apperror.Wrap(err, "OpName", ctx)` in Go, `throw new AppError(cause, { op, ctx })` in TS, and `AppException` in C#/PHP; preserve the root cause and causal stack.
- [ ] **`spec/03-error-manage/02-error-architecture/02-go-delegation-fix.md`**
  - **Why:** Prevents nil pointer panics and raw error leaks in Go routines.
  - **How:** Never delegate errors to uninitialized handlers; use explicit, typed error delegation channels with mutex guards.
- [ ] **`spec/03-error-manage/02-error-architecture/03-notification-colors.md`**
  - **Why:** Standardized error severity and UI feedback mapping.
  - **How:** Map log levels strictly: `debug` (trace), `info` (lifecycle), `warn` (recoverable/amber), `error` (user-visible failure/red), `fatal` (process exit).
- [ ] **`spec/03-error-manage/02-error-architecture/04-error-modal/01-copy-formats/07-envelope-error-response.md`**
  - **Why:** Universal API response contract across all endpoints.
  - **How:** Every HTTP/RPC response MUST return the standard envelope: `{ "data": T, "errors": [AppError], "meta": Meta }`. Never return raw un-enveloped error text.
- [ ] **`spec/03-error-manage/02-error-architecture/04-error-modal/02-react-components/02-error-store.md`**
  - **Why:** Centralized UI error presentation.
  - **How:** Frontend errors flow exclusively through a single global error store and universal error modal; no per-component alert boxes or unhandled promise rejections.
- [ ] **`spec/03-error-manage/03-error-code-registry/`**
  - **Why:** Stable error code registry and catalog.
  - **How:** All error codes must be registered constants (e.g. `ErrCodeNotFound`, `INVALID_PAYLOAD`). No ad-hoc string literals invented at the throw site.
- [ ] **`spec/03-error-manage/01-error-resolution/02-debugging-cheat-sheet.md`**
  - **Why:** Rapid triage and systematic root-cause discovery.
  - **How:** Follow the 4-part RCA pattern: Symptoms, Root Cause (1 sentence), Fix Applied, and Regression Prevention.
- [ ] **`spec/03-error-manage/01-error-resolution/04-verification-patterns/01-frontend-backend-sync.md`**
  - **Why:** Bidirectional integration verification.
  - **How:** Before claiming an integration works, verify both directions: inspect backend response payloads and test frontend error rendering. One side is not enough.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-error-management.py` or `linter-scripts/validate-guidelines.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If no dedicated error linter exists, create `linter-scripts/check-error-management.py` that AST-scans for:
  1. Empty `catch` or `except` blocks.
  2. Bare un-wrapped error returns (`return err` instead of `apperror.Wrap`).
  3. Bare panics/hard exits (`panic()`, `process.exit()`, `os.Exit()`).
  4. Non-standard API responses lacking the `{ data, errors, meta }` envelope.
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-error-management.py
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:errors"] = ["python", "linter-scripts/check-error-management.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running `python linter-scripts/check-error-management.py`.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask, applying surgical refactoring until all error management checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-error-management/
    2. Apply surgical refactoring (wrap errors with AppError, inject operation context, fix response envelopes).
    3. Run the error management linter:
          python linter-scripts/check-error-management.py
    4. Run the universal guideline autofixer:
          python .lovable/ai-fix-scripts/02-guideline-autofixer.py <modified-files>
    5. Run the local CI runner:
          python .lovable/ai-fix-scripts/03-cicd-local-runner.py
    6. IF any check fails:
          - Diagnose failure, fix code, and re-run immediately.
       IF all checks pass (exit code 0):
          - Mark subtask completed and proceed to next subtask.

    7. When all subtasks are finished and local CI is 100% green:
          - BREAK and proceed to End of Tunnel.
```

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `.lovable/ai-fix-scripts/index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python .lovable/ai-fix-scripts/01-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `.lovable/ai-fix-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `.lovable/ai-fix-scripts/index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/temp/active-locks.json`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `plans/index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] Staged files sanitized of artifact zips and temporary scratch files.
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.lovable/coding-guidelines/coding-guidelines.md`.
- [ ] /learn and apply as a /goal `.lovable/coding-guidelines/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.
- [ ] Error Manage Checklist: I have fully read and enforced the error management files at `spec/03-error-manage/`. I understand which files to follow (architecture, response envelopes) and how to follow them (never swallow errors, always wrap with context).
- [ ] Boolean Examples & Fixations: All boolean variables MUST begin with `is`, `has`, `can`, or `should` (e.g., `isReady`, `hasData`). NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`). NEVER use negative booleans (e.g., `isNotReady`, `disableCache`). NEVER invert success checks (e.g., `!response.isSuccess` is banned; use `response.isFail`).
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Semantic Tests: All unit test names are strictly semantic and behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`). `TestHandleComp100` is an immediate failure.
- [ ] Function Size: No function exceeds 15 lines. Long arguments are split across lines (max 100 chars).
- [ ] Error Handling (AppError): Errors use domain-specific `AppError` or custom `AppException` (for C#/OOP), not generic base `Error`.
- [ ] Code adheres to explicit booleans, `Type` suffixed Enums, and error wrapper rules.
- [ ] Formatting & Acronyms: Spacing rules are strictly followed. Acronyms are strictly PascalCase (`SwapIpWindows` not `SwapIPWindows`).
- [ ] Fast-forward commits created and pushed without rewriting published git history.
- [ ] Continuous loop maintained; only pausing to ask for "continue" on critical unrecoverable failures.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: I have fully read and strictly enforced every file in `spec/02-coding-guidelines/` and `.lovable/coding-guidelines/coding-guidelines.md`.
- [ ] Error Management: I have read and enforced `spec/03-error-manage/`. I used `AppError`/`AppException` and did not swallow errors.
- [ ] Boolean Conventions: All booleans begin with `is`, `has`, `can`, or `should` (e.g., `isFail`, `hasData`). NO negatives (`!isSuccess` is banned, use `isFail`).
- [ ] Semantic Naming: Absolutely NO generic garbage names (`temp`, `data`, `obj`, `comp_100`). All unit tests are behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`).
- [ ] Formatting: Signatures > 3 parameters or > 100 chars are split to one parameter per line. Newlines around every Markdown header (MD022) and lists are surrounded by blank lines (MD032).
- [ ] Acronyms & Magic Strings: Acronyms are PascalCase (`UserId` not `UserID`). Magic strings/numbers are extracted to constants.
- [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)
- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] Echo Back the Spec: I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] Pre-Commit Diff Proof: I have executed `git status` or `git diff --stat` and verified that the files I claim to have modified are actually listed as modified in the terminal output before committing.
- [ ] No Placeholder Search: I ran a regex search for `TODO` and `\[.*\]` in my modified files and confirmed I left zero placeholders behind. I actually wrote the implementation.
- [ ] Index Sync Deadman Switch: I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.lovable/what-to-read.md`. I did not leave any orphaned files.
- [ ] Blast Radius Acknowledgment: Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.

---

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memories/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

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

- slug: cg-error-management
- priority: high
- status: active
