# Code Hygiene & Project Architecture — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-hygiene`, `cg-execute hygiene`, `audit hygiene`, `fix file sizes`, `enforce code hygiene`, `parameter reduction`, `fix line endings`, `fix encoding`, `enforce utf8 lf`

> [!IMPORTANT]
> Prompt Version: 2.1.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, header alert block, or incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix all code hygiene, file size, parameter bloat, LF line ending (`\n`), UTF-8 (no BOM) encoding, and trailing newline violations across the codebase, enforcing 100-line standard file caps (recommended <= 80 lines), 8–15 line function caps, specialized parameter-reducing helper functions, extracting inline types, and sanitizing build artifacts until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all architectural violations and anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/readme.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/08-file-folder-naming/` for lowercase naming and continuous file sequencing.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/04-code-style/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase, Write .ai-memory/plans/pending/ Spec, Create .ai-memory/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Normalize LF & UTF-8, Remove Double Blank Lines, File & Function Splits, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Universal File Hygiene, Line Endings & Encoding Standards (Non-Negotiable)

Every repository artifact, source file (`.go`, `.ts`, `.tsx`, `.py`, `.php`, `.cs`, `.js`, `.json`, `.yaml`, `.yml`), and Markdown specification (`.md`) MUST strictly adhere to these universal file hygiene standards:

### 1. Unix LF (`\n`) Line Endings Only (TOTAL BAN on CRLF `\r\n`)

- **Strict LF:** Every file MUST use Unix-style line feeds (`\n`, `0x0A`).
- **Zero CRLF:** Windows-style carriage returns (`\r\n`, `0x0D 0x0A`) are strictly prohibited.
- Git and linters must verify LF across all files.

### 2. Strict UTF-8 Encoding (NO BOM)

- **UTF-8 Only:** All source code and markdown files MUST be encoded in **UTF-8 without Byte Order Mark (BOM)**.
- **BOM Banned:** Zero `\xef\xbb\xbf` header bytes. UTF-16 and UTF-32 are strictly forbidden.

### 3. Mandatory Single Trailing Newline at End of File (EOF)

- Every file MUST terminate with **exactly one newline (`\n`)** on the final line.
- Zero files missing a newline at EOF (`\ No newline at end of file` in Git diffs is an auto-reject).
- Zero multiple trailing blank lines at the end of a file (normalize to exactly one terminating newline).

### 4. No Function Starts with a Blank Line

- Every function body MUST begin writing executable code immediately on line 1 after the opening brace `{` (or `:` in Python).
- Placing an empty line as the first line of a function is strictly forbidden.

```go
// ❌ WRONG: Empty line right after opening brace
func CalculateDiscount(price float64) float64 {

    rate := getDiscountRate()
    return price * rate
}

// ✅ CORRECT: Code starts immediately on line 1; blank line before return
func CalculateDiscount(price float64) float64 {
    rate := getDiscountRate()

    return price * rate
}
```

### 5. Zero Double Blank Lines (`\n\n\n` Banned) in Code & Markdown

- There should **NEVER be two or more consecutive blank lines** anywhere inside source files or markdown documents.
- Always normalize multiple consecutive blank lines to exactly **one single blank line** (`\n\n`).

### 6. Markdown Header Spacing Rules (H1–H6: `#` through `######`)

- **Before Header:** Exactly **ONE blank line BEFORE** every markdown heading (EXCEPT when the heading is on line 1 of the file — line 1 has NO blank line before it).
- **After Header:** Exactly **ONE blank line AFTER** every markdown heading.

---

## Dedicated Section: Parameter Reduction & Specialized Helper Function Paradigm

Passing repeated constant arguments or flags across multiple call sites is a major source of code smell, parameter bloat, and boilerplate.

### Why Parameter Bloat & Repetition Is Forbidden

1. **Violates the <= 3 Parameters Limit:** Functions with excessive arguments are difficult to read, test, and maintain.
2. **Duplication of Context:** Hardcoding the same enum, flag, or code at 10 different call sites creates maintenance hazards when behavior changes.
3. **Impaired Readability:** Callers should express intent directly through semantic function names rather than passing tuples of flags and constants.

### The Specialized Helper Paradigm (Generic Example with Proper Newlines)

When a function call frequently repeats identical constants, enums, or exit codes, extract a specialized single-argument or zero-argument helper:

```go
// ❌ FORBIDDEN: Passing repeated constant/enum arguments at every call site
func ProcessPayload(data []byte) {
    if len(data) == 0 {
        reporter.ReportEvent(data, EventTypeValidationFailure, SeverityLevelError)
        return
    }

    reporter.ReportEvent(data, EventTypeProcessingSuccess, SeverityLevelInfo)
}

// ✅ REQUIRED: Specialized helper functions reducing parameter count and boilerplate
func ReportValidationError(data []byte) {
    reporter.ReportEvent(data, EventTypeValidationFailure, SeverityLevelError)
}

func ReportSuccess(data []byte) {
    reporter.ReportEvent(data, EventTypeProcessingSuccess, SeverityLevelInfo)
}

func ProcessPayload(data []byte) {
    if len(data) == 0 {
        ReportValidationError(data)
        return
    }

    ReportSuccess(data)
}
```

---

## Canonical Size Tier Reference

You MUST adhere to the single source of truth defined in `02-spec/02-coding-guidelines/02-canonical-size-tier.md`:

| Metric | Limit | Enforcement |
|---|---|---|
| **Function body (preferred)** | <= 8 lines | warn |
| **Function body (hard cap)** | <= 15 lines | error (build fails) |
| **File length (standard max)** | <= 100 lines | error (coding lines) |
| **File length (recommended)** | <= 80 lines | info |
| **React component file** | <= 80–100 lines | error (max 100 lines) |
| **Struct / class** | <= 120 lines | error |
| **Nested `if` statements** | 0 (No nesting) | error (flatten with guard clauses) |
| **Function Parameters** | <= 3 parameters | error (use specialized helpers / structs) |
| **Line Endings** | LF (`\n`) only | error (CRLF auto-rejected) |
| **File Encoding** | UTF-8 (no BOM) | error (BOM auto-rejected) |
| **Double Blank Lines** | 0 (Banned) | error (normalize to 1 blank line) |

---

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.ai-memory/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.ai-memory/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/readme.md`, `.ai-memory/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the repository's dedicated Python discovery scripts first:

1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --lang go --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --ext .go --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --path <folder-path> --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide violations.

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.ai-memory/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.ai-memory/plans/subtasks/` to `.ai-memory/plans/completed/` and update `.ai-memory/plans/readme.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/` and `.ai-memory/cicd-issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.
> 6. **GitHub Actions Zero Storage (Rule R18):** Never use `actions/upload-artifact` in CI workflows (`ci.yml`, matrix builds, linter runs). Free tier accounts have an account-wide cap of 0.5 GB (500 MB). CI runs must remain completely ephemeral. Releases belong exclusively in GitHub Releases (`release.yml`), never in Actions artifact storage.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .ai-memory/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Automated File & Hygiene Fixer:** Use `python 03-ai-scripts/03-file-manipulator.py` and `05-guideline-autofixer.py` to normalize LF endings, remove trailing whitespace, and fix file size boundaries.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/readme.md` using sequential script naming (e.g., `01-file-manipulator.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.ai-memory/readme.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.ai-memory/plans/readme.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **No Function Starts with Blank Line:** Functions start immediately on line 1 with code.
- [ ] **Zero Double Blank Lines:** No `\n\n\n` in code or markdown.
- [ ] **Markdown Heading Spacing:** Exactly one blank line before and after headings (no leading blank line on line 1).
- [ ] **File Size Caps:** All files <= 100 coding lines (recommended <= 80 lines).
- [ ] **Function Sizing:** All functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `02-spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.ai-memory/coding-guidelines.md`.

1. [ ] /learn and apply as a /goal `.ai-memory/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.

- [ ] `python linter-scripts/check-file-sizes.py` and `python linter-scripts/check-newline-styling.py` exited with code 0.
- [ ] **Atomic File Recording:** Modified files recorded to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). Builds, tests, and CI runner are deferred to CI/CD fix.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `02-spec/02-coding-guidelines/02-canonical-size-tier.md`, `02-spec/02-coding-guidelines/08-file-folder-naming/`, and `.ai-memory/coding-guidelines.md`.
- [ ] LF Line Endings & UTF-8 (No BOM): Verified Unix LF and UTF-8 across all files.
- [ ] Boolean Principles & IsDefined: Positive prefixes only (`is`/`has`). MANDATORY: Use `isDefined` (or `res.IsDefined()`) instead of inverted empty checks (`!isEmpty` / `!res.IsEmpty()`). Never write `if !isEmpty`. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`.
- [ ] GitHub Actions Zero Storage (Rule R18): Never use `actions/upload-artifact` in CI workflows; maintain zero Actions storage usage.

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Scripts:** `linter-scripts/check-file-sizes.py`, `linter-scripts/check-newline-styling.py`, `linter-scripts/check-markdown-header-spacing.py`
2. **Local Run Command:** `python linter-scripts/check-file-sizes.py`
3. **Autofixer Command:** `python 03-ai-scripts/03-file-manipulator.py`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate Code Hygiene, Line Endings & Sizes
     run: |
       python linter-scripts/check-file-sizes.py
       python linter-scripts/check-newline-styling.py
       python linter-scripts/check-markdown-header-spacing.py
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "File Sizes Check": [sys.executable, "linter-scripts/check-file-sizes.py"],
       "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
       "Markdown Header Check": [sys.executable, "linter-scripts/check-markdown-header-spacing.py"],
   }
   ```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.ai-memory/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.ai-memory/test-inventory.json` so associated tests are known for future release verification.
