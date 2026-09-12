# Cross-Platform CI/CD & Run Scripts Fix — Workflow (must follow)

Trigger Keywords & Aliases: `fix with RCA`, `FRCA : Fix with RCA`, `fix`, `fix, fix`, `CI/CD fix`, `fix run scripts`, `force rebuild`, `force create`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

- **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script `03-ai-scripts/06-cicd-local-runner.py` from scratch**, regardless of whether the file already exists on disk.

```text
N = 200
```

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal First `N/2` steps (Phase 1): Review the central CI/CD pipeline definitions (`.github/workflows`, `.gitlab-ci.yml`, etc.) and cross-reference them with the local Python runner (`03-ai-scripts/06-cicd-local-runner.py`).
   - **Condition:** If `03-ai-scripts/06-cicd-local-runner.py` does not exist, you must create it immediately.
   - **Condition:** You must ensure that **every single CI/CD case** that needs to run in the pipeline can also be run locally from this Python script (with Docker stripped for native host execution). Improve the Python script to cover all cases if any are missing.
2. [ ] /goal Second `N/2` steps (Phase 2): Run the local runner script (`python 03-ai-scripts/06-cicd-local-runner.py --all`) to catch all errors. Singly execute the script in an autonomous self-loop, zeroing in on one failing error per turn (4-part RCA -> surgical fix -> guideline autofixer -> re-verify).
3. [ ] /goal Finalize CI/CD: Your ultimate goal is to fix and finalize the CI/CD. You must loop until the Python local runner script executes flawlessly with **no errors** (exit code 0) for all registered cases. Do not stop until this goal is met.

/goal Perform a Root Cause Analysis (RCA) on all failing run scripts and CI/CD workflows, update `03-ai-scripts/06-cicd-local-runner.py` with any newly added pipeline steps from screenshots or workflow files, zero in on each error singly using self-looping, persist the RCA into `.lovable/cicd-issues/` and `.lovable/strictly-avoid.md`, implement universal query wrappers with explicit success/failure boolean results and automated error logging, verify clean builds, commit logically, and push to git.

/learn Ingest recent Root Cause Analysis (RCA) records from `.lovable/cicd-issues/`, `.lovable/issues/`, `spec/02-coding-guidelines/02-canonical-size-tier.md`, `spec/02-coding-guidelines/01-cross-language/01-index.md`, `spec/02-coding-guidelines/01-cross-language/01-index.md`, and `spec/03-error-manage/` so previous mistakes and anti-patterns are never repeated.

---

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without cancelling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands](spec/13-generic-cli/01-index.md)`
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## Screenshot & Pipeline Discovery Protocol (Execute First When Any Image Is Provided)

> [!IMPORTANT]
> **If the user provides any image or screenshot showing a CI/CD pipeline name, failing workflow, or error log:**
>
> 1. **FIRST ACTION — Update Python Runner:** Locate the pipeline/job in `.github/workflows/*.yml` (or repo CI configs) to find whatever new jobs, steps, or linters were added, and **immediately update `03-ai-scripts/06-cicd-local-runner.py`** to include them in the `JOBS` dictionary.
> 2. **SECOND ACTION — Singly-Done Self-Loop Execution:** Run the Python script iteratively, zeroing in on one failing error at a time using strictly bounded self-loop turns until all checks exit with code 0 (`exit 0`).

### Bounded Single-Step Self-Loop Sequence (Singly Done — No Overloaded Steps)

Every step must be **singly done** using bounded self-looping turns:

- **Self-Loop Step 1 (Extract Pipeline Name from Screenshot):**
  1. Read image to extract the pipeline name, failing job name, and error snippet.
  2. Scan `.github/workflows/*.yml` to identify the corresponding shell commands and dependencies.

- **Self-Loop Step 2 (FIRST ACTION: Update Python Runner Script):**
  1. Open `03-ai-scripts/06-cicd-local-runner.py`.
  2. If new jobs/steps are found in workflows, strip Docker wrappers and register the new commands in the `JOBS` dictionary.
  3. Save the runner script and verify syntax.

- **Self-Loop Step 3 (Execute Runner & Baseline Failures):**
  1. Run `python 03-ai-scripts/06-cicd-local-runner.py`.
  2. If exit code = 0, proceed to End of Tunnel. If exit code != 0, zero in on the first specific failure.

- **Self-Loop Step 4 (RCA & Zero In on Error):**
  1. Write 4-part RCA in `.lovable/memory/issues/xx-<slug>.md`.
  2. Register in `.lovable/01-index.md` and `.lovable/strictly-avoid.md`.

- **Self-Loop Step 5 (Surgical Code Fix):**
  1. Open the specific file and line, apply minimal surgical fix.
  2. Run `python 03-ai-scripts/05-guideline-autofixer.py <modified-files>`.

- **Self-Loop Step 6 (Re-Verify & Loop):**
  1. Re-run `python 03-ai-scripts/06-cicd-local-runner.py`.
  2. If resolved and more errors remain, self-loop to Step 4 to zero in on the next error until exit code = 0.

---

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## Action Items — Must Follow (Non-Negotiable)

- [ ] Ingest past RCAs from `.lovable/cicd-issues/` and `.lovable/issues/` before coding.
- [ ] Fix CI/CD and run all scripts.
- [ ] Find the root cause of the issue and write it into the 'avoid' part of the `.lovable` memory (`.lovable/strictly-avoid.md` and `.lovable/cicd-issues/01-<slug>.md`).
- [ ] Make Git commits properly.
- [ ] Check the CI/CD, run the tests, and build the code; fix any issues found.
- [ ] Create a query wrapper for PHP/Python/TS that automatically logs failures to reduce code duplication.
- [ ] Ensure the wrapper explicitly returns success or failure states (e.g., `isSuccess`, `isFail`).
- [ ] Identify everywhere this logging wrapper pattern was missed or messed up and fix those places.
- [ ] Update the memory inside the `.lovable` folder regarding this wrapper pattern so future AI agents do not make the same mistake.
- [ ] Make a plan for the required fixes and self-loop to execute it.
- [ ] Group similar code changes together into single commits (do not commit one file at a time) and include a nice commit message.
- [ ] Push the code to the repository before ending the job.
- [ ] Fix any remaining issues that arise before completion.

## Before Writing Code

Read and follow spec folders `02`, `03` and `04` before writing any code. Error management must be followed. Code must be DRY.

## The 4-Part RCA Requirement (Mandatory Memory File)

Before you write any code to fix the problem, you MUST document the issue in `.lovable/memory/issues/xx-<slug>.md` (where XX is the next available sequential number). The file MUST contain these exact four sections:

1. **Why it happened:** The high-level business, logical, or architectural breakdown of the failure.
2. **How it happened:** The technical execution flow that triggered the bug.
3. **Root Cause:** The exact file, line, and dependency responsible for the failure.
4. **Code Fix:** The exact code snippets showing what needed to be changed to fix the root cause.

## STRICT AVOIDANCE: Never Disable CLI Linting, Static Analysis, or CI/CD Checks (No Shortcut Cheating)

> [!CAUTION]
> **TOTAL BAN ON DISABLING, SKIPPING, OR BYPASSING CLI LINTERS AND CI/CD GATES:**
>
> - **NEVER** disable, comment out, delete, or skip any CLI linting command (`golangci-lint`, `eslint`, `markdownlint`, `tsc`, `pytest`, `phpstan`, `mypy`, `check-*.py`), build step, or test suite.
> - **NEVER** add `|| true`, `continue-on-error: true`, `# nolint`, `// eslint-disable`, or ignore flags to "quickly win the race" or fake a pipeline pass.
> - **Your job is to legitimately fix the underlying source code.** If resolving complex lint errors or test failures requires multiple sub-steps, sub-agents, or nested self-looping turns, you MUST execute all necessary turns until the code is 100% clean and compliant.
> - Disabling or bypassing any CI/CD or CLI lint check is an automatic and immediate rejection.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## Actionable Items & Checklist (must follow)

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /learn previous RCAs in `.lovable/cicd-issues/` and `.lovable/strictly-avoid.md`.

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.

## Core Rules & Non-Negotiable Checklist for AI (Must Verify Before Completing Task)

Before finalizing any code modification, you MUST manually verify the following:

- [ ] **No Disabling CLI Linting (Zero Bypassing):** All CLI linters and CI/CD quality gates executed fully without `|| true`, `continue-on-error`, or suppression comments. Code was legitimately fixed.
- [ ] **Legitimate Multi-Step Self-Looping:** If complex errors occurred, I performed dedicated, single-step self-loop iterations to resolve each underlying failure instead of taking shortcuts.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.lovable/coding-guidelines.md`.
- [ ] Error Manage Checklist: I have fully read and enforced the error management files at `spec/03-error-manage/`. I understand which files to follow (architecture, response envelopes) and how to follow them (never swallow errors, always wrap with context).
- [ ] Boolean Examples & Fixations: All boolean variables MUST begin with is and has only (can, should, was, etc. are banned) (e.g., `isReady`, `hasData`). NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`). NEVER use negative booleans (e.g., `isNotReady`, `disableCache`). NEVER invert success checks (e.g., `!response.isSuccess` is banned; use `response.isFail`).
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Function Signatures (R4, R5, R9): If a function has `> 3 parameters` or the signature is `> 100 chars`, you MUST split it so there is exactly one parameter per line.
- [ ] Error Handling (R7): No silent failures or swallowed errors. Use explicit boolean states (e.g., `isFail`). Never invert success booleans (e.g., avoid `!isSuccess`).
- [ ] Magic Strings/Numbers (R8): Extract all magic strings/numbers into named constants.
- [ ] Enums: TypeScript string unions are banned. All Enums must end with the `Type` suffix.
- [ ] Naming & Casing (R1, R2): PascalCase everywhere. Acronyms (Id, Json, Url) are Pascal case, never all-caps (e.g., `UserId`, not `UserID`).
- [ ] Blank Lines (R13-R20): One blank line before every `return`/`throw`. One blank line after closing `}`. Never two blank lines in a row.

## Language-Specific Requirements

- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.

## End of Tunnel Checklist

- [ ] **Zero Linting/CI/CD Bypass:** Confirmed that NO CLI linters, static analysis tools, or test scripts were disabled, commented out, skipped, or bypassed with `|| true`.
- [ ] **Local CI Runner Clean:** `python 03-ai-scripts/06-cicd-local-runner.py` exited with code 0.
- [ ] **All Scripts & Workflows Verified:** All tests, builds, and query wrappers run without errors.
- [ ] **RCA Documented:** Memory files written to `.lovable/cicd-issues/` and `.lovable/memory/issues/`.
- [ ] **Antigravity Skill Updated:** Verified `.agents/skills/ci-cd-fix/skill.md` is present and synchronized with the latest rules.
- [ ] **Commit & Push:** Group changes into clean development commit and push to remote. No automatic release.
- [ ] **File Change Summary:** Provide a detailed summary in chat listing exactly which files were changed, what was changed inside them, and why they were changed.
