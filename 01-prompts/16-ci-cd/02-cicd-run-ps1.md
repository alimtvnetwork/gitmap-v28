# PowerShell CI/CD Pipeline & Runner Creation — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

You are responsible for creating or fixing CI/CD pipelines (e.g., GitHub Actions workflows) and dynamic execution scripts (e.g., `run.ps1`, `run.sh`, `run.config.json`).

## Strictly Avoid: No Automatic Releases & Run All Tests (Strict Policy)

> [!CAUTION]
> **NO AUTOMATIC RELEASES:** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
> **RUN ALL TESTS & QUALITY GATES (FULL PIPELINE FIDELITY):** All CI/CD repair and runner tasks MUST run all unit tests, integration tests, and quality gates properly (`python 03-ai-scripts/06-cicd-local-runner.py`). NEVER skip or disable tests with `--no-tests` in CI/CD fix workflows.
> **ATOMIC CHANGE TRACKING:** Record all modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

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
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `02-spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands]`02-spec/13-generic-cli/01-index.md)`
> 6. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.
> 7. **Temp & Failure Folder Isolation:** All temporary directories, runner caches, and test artifacts MUST be strictly placed in `.lovable/temp/`. Creating `.tmp/` at the repository root or outside `.lovable/` is strictly forbidden.
>    - Dedicated Failure Directory: `.lovable/temp/failures/` is the dedicated folder where failed tests and failed quality gates write error logs (`<test-or-job-name>.log`).
>    - Passing Tests Completely Silent: Passing tests must produce ZERO filesystem artifacts (zero files written) and remain completely silent in output logs.
> 8. **Runner In-Flight ETA Wait Protocol:** When running background commands, the runner dynamically writes live status and remaining ETA to `.lovable/temp/runner-eta.json` (emitting in-flight heartbeats strictly every 25 seconds or more). If an agent inspects an active background job and it is still running, the agent MUST sleep/wait for **1 minute (60 seconds) each time**, or dynamically sleep for the remaining ETA duration read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling or querying in loops.
> 9. **Centralized Test Inventory & Incremental Caching:** All unit tests are cataloged in `.lovable/test-inventory.json` with strictly repository-relative paths (`target_file`, `test_file`). First run executes all tests to establish baseline timings; subsequent runs execute incrementally only if target code files or test files change. Slow test threshold defaults to `4.0s` (configurable via `GITMAP_SLOW_TEST_THRESHOLD`).
> 10. **Dual-Queue Worker Pools:** Slow tests run in a dedicated 4-worker pool running at most 2 tests at a time per batch. Fast tests run in a 4-worker pool running at most 4 tests at a time, pulling in chunks of 100 tests from the test inventory queue until all are complete.
> 11. **Dynamic ETA Sleep Protocol:** The AI agent reads `.lovable/temp/runner-eta.json`, sleeps for the estimated wait time rather than looping, and if still active upon waking, re-checks remaining ETA and sleeps again to avoid burning tokens.
> 12. **Zero-Storage GitHub Actions Mandate (Total Ban on CI Artifact Uploads):** CI workflows MUST NOT upload test outputs, coverage files, Playwright reports, or drift summaries via `actions/upload-artifact`. Free-tier accounts have a strict 0.5 GB shared quota across all repositories. All reports, failures, and summaries MUST be emitted directly to `$GITHUB_STEP_SUMMARY`, console stdout (`cat log.txt`), or sticky PR comments with zero storage consumption. Only true GitHub release assets (binaries/tarballs on tagged releases) are permitted.

---

## Generic CI/CD & Automation Guidelines Checklist

This document serves as a strict, universal checklist and specification for setting up and fixing CI/CD pipelines, as well as the standard run scripts (e.g., `run.ps1`/`run.sh`), for any project. Any AI agent operating on DevOps or CI/CD tasks MUST read and follow these guidelines to ensure consistency across different tech stacks and repositories.

### 1. Specification Files to Read & Maintain

Before making any changes to `.github/workflows` or automation scripts, you must read the following architecture documents. These contain the foundational constraints and mechanisms for deployment, automation, and CI/CD pipelines.

#### PowerShell & Orchestration (`02-spec/11-powershell-integration`)

- [ ] `02-spec/11-powershell-integration/01-index.md`
- [ ] `02-spec/11-powershell-integration/04-script-reference.md`
- [ ] `02-spec/11-powershell-integration/05-integration-guide.md`

#### CI/CD Pipeline Workflows (`02-spec/12-cicd-pipeline-workflows`)

- [ ] `02-spec/12-cicd-pipeline-workflows/01-index.md`
- [ ] `02-spec/12-cicd-pipeline-workflows/02-ci-pipeline.md`
- [ ] `02-spec/12-cicd-pipeline-workflows/05-release-pipeline.md`
- [ ] `02-spec/12-cicd-pipeline-workflows/07-install-script-generation.md`
- [ ] `02-spec/12-cicd-pipeline-workflows/09-changelog-integration.md`

#### CLI & Build (`02-spec/13-generic-cli`)

- [ ] `02-spec/13-generic-cli/01-index.md`
- [ ] `02-spec/13-generic-cli/11-build-deploy.md`
- [ ] `02-spec/13-generic-cli/18-batch-execution.md`

#### Update Mechanisms (`02-spec/14-update`)

- [ ] `02-spec/14-update/05-build-scripts.md`
- [ ] `02-spec/14-update/18-release-pipeline.md`
- [ ] `02-spec/14-update/19-install-scripts.md`

#### Release Engineering (`02-spec/16-generic-release`)

- [ ] `02-spec/16-generic-release/01-index.md`
- [ ] `02-spec/16-generic-release/04-install-scripts.md`

#### Context / Issue Logging

- [ ] `cicd-issues/<issue-name>.md`
  - Any time a CI/CD pipeline fails, an AI must log the failure here (including error traces and environment context) before attempting a fix.

---

### 2. CI/CD Implementation Checklist

When building or fixing the CI/CD pipelines (e.g., `.github/workflows/ci.yml`), enforce the following steps universally:

- [ ] Dependency Caching: The pipeline MUST cache dependencies based on the project's lockfiles (e.g., `package-lock.json`, `bun.lockb`, `poetry.lock`, `go.sum`). This minimizes build times.
- [ ] Toolchain Matching: Ensure the CI runner uses the exact tool versions specified in the project's configuration (e.g., `.nvmrc`, `.python-version`, or the generic `run.config.json`).
- [ ] Linting & Formatting: Run the project's defined linting and formatting commands first. The pipeline must fail immediately if there are style violations, preventing bad code from proceeding to tests.
- [ ] Type Checking / Static Analysis: If the language supports it (e.g., TypeScript, Python with mypy, Go), run static analysis as a parallel job before or alongside testing.
- [ ] Test Execution: Execute the project's testing suites (unit, integration, e2e) as defined in the configuration file.
- [ ] Artifact Generation & Zero CI Storage: CI workflows MUST NOT upload test reports, logs, or coverage artifacts via `actions/upload-artifact`. Compile/zip release binaries ONLY on tagged release jobs (attached directly to GitHub Releases, not Actions artifact storage). All CI summaries and failure diagnostics MUST stream to `$GITHUB_STEP_SUMMARY` or console stdout.

---

### 3. Dynamic Script Architecture (e.g., `run.ps1` & `run.config.json`)

To prevent hardcoded commands and ports, the execution architecture MUST be dynamic. The local run script (e.g., `run.ps1` or `run.sh`) must act as a generic orchestrator, while a JSON configuration file serves as the single source of truth for both local development and CI execution.

#### Expected Configuration Structure (Reference: `run.config.json`)

The JSON file should define all services, ports, and lifecycle commands.

```json
{
  "projectName": "Generic Project Name",
  "frontend": {
    "dir": "path/to/frontend",
    "port": 3000
  },
  "backend": {
    "dir": "path/to/backend",
    "port": 8080
  },
  "host": "127.0.0.1",
  "commands": {
    "install": "package-manager install",
    "dev:frontend": "command to start frontend dev server --port {fePort}",
    "dev:backend": "command to start backend server --port {bePort}",
    "build:frontend": "command to build frontend",
    "test:frontend": "command to run frontend tests",
    "test:backend": "command to run backend tests",
    "lint": "command to run linter"
  }
}
```

#### Expected Run Script Implementation Rules (Reference: `run.ps1` / `run.sh`)

When writing or modifying the orchestration scripts, any AI must adhere to the following logic:

1. Parse Configuration First:
   The script must read the configuration file (e.g., `run.config.json`) into memory and extract variables (ports, directories, commands).
2. Dynamic Command Injection:
   Never hardcode commands (e.g., `npm run dev`). Instead, read the command string from the JSON and dynamically substitute any necessary placeholders (like `{fePort}`).
3. Graceful Error Handling & Waiting:
   When orchestrating multiple services, the script must wait for dependencies to become healthy before proceeding. Use generic HTTP or TCP polling mechanisms to ensure a backend is fully online before a frontend attempts to connect to it.
4. CI Mode Toggle:
   The script must accept a `-CI` (or `--ci`) flag. When running in CI mode:
   - Skip launching local browsers or interactive shells.
   - Run the build, lint, and test commands from the JSON instead of the dev server commands.
   - Exit with code `1` immediately if any step fails.
5. Process Cleanup (Trap/Finally):
   Ensure all spawned processes are captured in a job/PID array. The script must aggressively kill these processes in the `finally {}` or `trap` block on exit. There should be no zombie processes left blocking ports.

---

## 4. Commit Fix Actionable Items & Checklist (Non-negotiable)

After completing the pipeline and run script creation, you MUST follow this checklist for committing and verifying the changes.

### 4.1. Pre-flight & Planning

- [ ] Ensure the git repository starts completely clean. If dirty, commit, stash, or fix git issues before writing any new code.
- [ ] Read the overarching main task plan from `.lovable/plans/pending/xx-<slug>.md` to understand what needs to be executed.
- [ ] Derive the `<slug>` from the plan filename itself (e.g., plan file `03-auth-refactor.md` → slug is `03-auth-refactor`). Never invent a slug.
- [ ] Confirm subtask files exist under `.lovable/plans/subtasks/xx-<slug>/SS-<subslug>.md` for each step that needs parallel execution. Create them if missing, following the plan prompt structure.
- [ ] Ensure the plan is highly extensive, explicitly detailing where and how to make changes so sub-agents can easily execute tasks (Non-negotiable).
- [ ] Direct subtasks to `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md` and update plans in `.lovable/plans/pending/01-<slug>.md`. Do not write randomly into `.lovable`.
- [ ] Read the memory files and the spec folder coding guidelines + error manage guidelines before touching code.
- [ ] Anti-Hallucination: If referenced files are missing or ambiguous, stop and ask clarifying questions.

### 4.2. Ruthless Management & Subtask Looping

- [ ] Map out the subtasks from the big plan and spawn sub-agents for all independent tasks simultaneously (MAXIMUM 2 sub-agents concurrently to avoid RAM and caching issues).
- [ ] Each sub-agent may only run a MAXIMUM of 2-3 async operations at a time.
- [ ] Enforce lifecycle: sub-agent reads subtask file → marks `🔄 In Progress` → works → marks `✅ Done` with file list and summary → updates parent plan step → signals completion.
- [ ] Track queue state by counting total subtasks spawned vs. total `✅ Done` entries in the subtask files. Proceed to commit only when the counts match.
- [ ] If a sub-agent fails to update its status file or gives garbage, kill it immediately and restart it.

### 4.3. Root Cause

- [ ] Find the root cause of the problem first, before applying any fix.
- [ ] Record the root cause strictly into the `.lovable` memory structure per the write protocols.

### 4.4. File System Writes & Main Agent Commit

- [ ] Sub-agents write to the file system and update their task entries. They do NOT commit.
- [ ] Wait until all sub-agents have signaled completion and updated their subtask files under `.lovable/plans/subtasks/`.
- [ ] Ensure `.gitignore` explicitly excludes test reports, test data, artifacts, and compiled binaries (Non-negotiable).
- [ ] RED FLAG: Verify absolutely NO test results or binaries are staged before making the commit.
- [ ] Group all completed work into a single logical commit.
- [ ] If issues arise during the commit, fix them immediately and retry.
- [ ] Push the commit to the remote repository. Pushing is non-negotiable.

### 4.5. Code Standards (non-negotiable)

- [ ] Follow the code review guidelines from the aspect folder.
- [ ] Ensure every try-catch block explicitly logs the error according to the error manage folder.
- [ ] Create a query wrapper for PHP/Python/TS that handles automatic failure logging, so logging is not scattered.
- [ ] Use explicit `isFail` properties; NEVER use inverted success checks (use `response.isFail`, not `!response.isSuccess`).
- [ ] Remove all magic strings and magic numbers unless used directly for logging — and state that logger exception in the typing.
- [ ] Replace TypeScript string union types (e.g. `"pass" | "fail" | "fallback"`) with Enums.
- [ ] Ensure every Enum name ends with the `Type` suffix (e.g. `StatusType`, never `Status` or `Status7`).
- [ ] Ensure all Enum values are written in PascalCase (e.g., `enum StatusType { ActiveState = "ACTIVE" }`), avoiding `_camelCase` or `camelCase`, unless the specific language (like Rust) conventionally dictates otherwise.
- [ ] Reuse constants — never duplicate them. Code must always be DRY; never repeat code. This is high priority.

### 4.6. End-of-Loop Final Verification (Once only, at the very end)

- [ ] Check the full build. Fix every build failure, commit, and push.
- [ ] Run all unit tests. Fix every failing test, commit, and push.
- [ ] Check CI/CD status and ensure pipelines pass.
- [ ] Audit that coding guidelines from the aspect folder and error manage folder have been followed across all changed files.
- [ ] Finish the job only when everything is green, pushed, and fully verified.

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

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.

---

## Metadata

- slug: cicd-run-ps1
- status: active

## Core Rules & Non-Negotiable Checklist for AI (Must Verify Before Completing Task)

Before finalizing any code modification, you MUST manually verify the following:

- [ ] **No Disabling CLI Linting (Zero Bypassing):** All CLI linters and CI/CD quality gates executed fully without `|| true`, `continue-on-error`, or suppression comments. Code was legitimately fixed.
- [ ] **Zero Actions Storage (Total Ban on CI Artifacts):** Confirmed that NO `actions/upload-artifact` steps exist in CI workflows; all diagnostic outputs stream to `$GITHUB_STEP_SUMMARY` or console logs.
- [ ] **Legitimate Multi-Step Self-Looping:** If complex errors occurred, I performed dedicated, single-step self-loop iterations to resolve each underlying failure instead of taking shortcuts.
- [ ] **Boolean Standards & IsDefined (Rule R3):** Positive prefixes only (`is` and `has`). MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Never write `if !isEmpty`. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
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
- [ ] **Zero Actions Storage:** Confirmed that NO `actions/upload-artifact` steps exist in CI workflows, eliminating quota depletion.
- [ ] **Local CI Runner Clean:** `python 03-ai-scripts/06-cicd-local-runner.py` exited with code 0.
- [ ] **All Scripts & Workflows Verified:** All tests, builds, and query wrappers run without errors.
- [ ] **RCA Documented:** Memory files written to `.lovable/cicd-issues/` and `.lovable/memory/issues/`.
- [ ] **Antigravity Skill Updated:** Verified `.agents/skills/ci-cd-fix/skill.md` is present and synchronized with the latest rules.
- [ ] **Commit & Push:** Group changes into clean development commit and push to remote. No automatic release.
- [ ] **File Change Summary:** Provide a detailed summary in chat listing exactly which files were changed, what was changed inside them, and why they were changed.
