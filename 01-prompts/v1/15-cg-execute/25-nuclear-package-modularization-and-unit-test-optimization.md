# Nuclear Package Modularization & Unit Test Optimization — Coding Guideline Execution (must follow)

Trigger Keywords & Aliases: `cg-nuclear-packages`, `cg-package-modularization`, `cg-test-optimization`, `cg-nuclear`, `nuclear-packages`, `isolate-heavy-tests`, `optimize-unit-tests`, `split-packages`

> [!IMPORTANT]
> Prompt Version: 2.2.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, header alert block, or incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform (configurable per run).

/goal Autonomously scan, audit, plan, and modularize monolithic packages and optimize unit test execution across the codebase. Enforce a strict Directed Acyclic Graph (DAG) architecture, extract reusable zero-dependency leaf packages, isolate slow or destructive tests (`exec.Command`, git CLI subprocesses, network sockets, `time.Sleep`) into external blackbox test packages (`tests/heavy_test/` under `package heavy_test`), maintain the centralized test inventory manifest (`.ai-memory/test-inventory.json`), and adhere to the 5-day cache freshness decision engine to achieve ultra-fast sub-0.05s unit test execution loops without circular dependencies.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

You MUST execute this task via a strict 3-Phase pipeline governed by the N-step budget. Do not skip steps.

```text
N = 200  (Total self-loop steps budget, read-only after initialization)
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Packages, Cache Freshness Audit, DAG Mapping, Master Plan in .ai-memory/plans/pending/)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Leaf Package Extraction, Domain Modularization, Heavy Test Isolation, Inventory Sync)
```

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step 0 - Verbatim Prompt Recording & Deliverables Extraction Gate): Directly capture the user's prompt verbatim into `.ai-memory/plans/pending/xx-nuclear-packages.md` under a dedicated `## User Request (Verbatim)` section. Extract whatever requirements were given into actionable deliverables with traceable IDs (`Task-01`, `Task-02`), and output this confirmed task breakdown directly in chat in cleanly indented markdown with vertical blank lines, task state (`State: [PENDING]`), and understanding indicator bracket (`Understood: [YES — ...]`) before any file exploration, scanning, or spec writing.
2. [ ] /goal Phase 1 (Step A - Test Inventory Freshness & Duration Audit): Audit test inventory freshness by running `python 03-ai-scripts/33-test-inventory-generator.py --check-age --max-age-days 5`. If `.ai-memory/test-inventory.json` is missing or older than 5 days (or unprofiled), run an initial profiling pass to record durations. If fresh (`<= 5 days old`), skip re-running all tests and use cached durations directly to make modularization decisions.
3. [ ] /goal Phase 1 (Step B - Monolithic Package & Dependency Topology Discovery): Deeply scan the target codebase using fast discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `18-codebase-topology-discoverer.py`) to map package dependency graphs, import cycles, bloated packages (>10 files or test runtimes > 2.0s), and candidate leaf packages.
4. [ ] /goal Phase 1 (Step C - Master Plan Generation & Violation Ledger): Write the master architectural plan into `.ai-memory/plans/pending/xx-nuclear-packages.md` with an exhaustive Violation Ledger table (Package, Current File Count, Monolithic Anti-Patterns, Target DAG Subpackages, Slow Tests to Isolate, Status).
5. [ ] /goal Phase 1 (Step D - Lean Subtask Decomposition): Decompose into granular subtasks in `.ai-memory/plans/subtasks/xx-nuclear-packages/01-<subtask>.md`, etc., with strictly relative Git paths.
6. [ ] /goal Phase 1 (Step E - Mandatory Auto-Loop): As soon as Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for permission**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.
7. [ ] /goal Phase 2 (Step A - Leaf Package Extraction): Extract zero-dependency leaf packages (`pkg/constants`, `pkg/model`, `pkg/appfault`, `pkg/fsutil`, `pkg/cliexit`). Ensure leaf packages NEVER import parent packages or domain packages.
8. [ ] /goal Phase 2 (Step B - Domain Subpackage Segregation & DAG Enactment): Decompose monolithic command/service packages into cohesive domain subpackages (e.g. `cmdprompt`, `cmdpurge`, `cloner`). Verify zero circular dependencies (`import cycle not allowed`).
9. [ ] /goal Phase 2 (Step C - Heavy Test Isolation): Segregate heavy tests invoking `exec.Command`, git CLI processes, network sockets, or `time.Sleep` into `tests/heavy_test/` (or `cli/tests/heavy_test/`) under `package heavy_test`. Ensure routine package unit tests contain only fast in-memory unit tests (< 0.05s).
10. [ ] /goal Phase 2 (Step D - Test Inventory Manifest Synchronization): Synchronize `.ai-memory/test-inventory.json` and atomically record all modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
11. [ ] /goal Phase 2 (Step E - Function & File Sizing Compliance): Enforce functions <= 8–15 lines, files <= 80–100 lines, affirmative booleans (`is*`, `has*`), zero explicit `== true`, and guard clauses.
12. [ ] /goal Phase 2 (Step F - Banned Intermediate Verification): DO NOT run unit test suites or builds during intermediate file edits.
13. [ ] /goal Phase 2 (Step G - Final Step Build Verification): At the conclusion of all subtasks, execute targeted build checks (`go vet ./...`, `go build ./...`) to verify 0 compiler errors or circular imports.
14. [ ] /goal Phase 3 (Step A - Task Consolidation): Consolidate completed subtasks into `.ai-memory/plans/completed/xx-nuclear-packages.md`, delete subtask files, and update `.ai-memory/plans/readme.md`.
15. [ ] /goal Phase 3 (Step B - Final Step Git Commit & Push): Stage all modified files, inventory updates, and plans (`git add -A`), commit in a single atomic commit, and push to git. Never commit per-file.
16. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
17. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
18. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical size tiers.
19. [ ] /learn Ingest `02-spec/02-coding-guidelines/03-golang/readme.md` for Go coding standards and package architecture.
20. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.

---

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate exploratory roundtrips, the AI agent MUST use the repository's dedicated Python discovery scripts first:

1. **Check Test Inventory Freshness & Timing Profiles:**
   ```bash
   # Check if inventory is <= 5 days old and contains duration profiles
   python 03-ai-scripts/33-test-inventory-generator.py --check-age --max-age-days 5

   # Inspect test inventory in JSON format
   python 03-ai-scripts/33-test-inventory-generator.py --check-age --json
   ```

2. **Discover Monolithic Packages & Subprocess Test Invocations:**
   ```bash
   # Find exec.Command invocations in tests (candidates for heavy test isolation)
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "exec.Command" --lang go --limit 100

   # Find time.Sleep calls in test files
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "time.Sleep" --lang go --limit 50

   # Map topology and package imports
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```

3. **Record Atomic Changes Under File Lock:**
   ```bash
   # Record refactored files into recent-file-changes.json cache
   python 03-ai-scripts/33-test-inventory-generator.py --record <file1> <file2>
   ```

---

## 1. The Monolithic Package Anti-Pattern & Test Latency Bottlenecks

In Go (and modular polyglot systems), package structure directly governs compilation speed, test isolation, and dependency hygiene:

### The Problem
When 20–40 files reside in a single monolithic package (e.g., `package cmd` or `package cli`):
- **Monolithic Invalidation:** Modifying a tiny string formatting helper in `utils.go` forces the Go compiler to recompile the entire package and re-link all associated tests.
- **Co-Mingled Test Latency:** If that same package contains slow integration tests (e.g. running `git clone`, spawning subprocesses, spinning up HTTP servers), `go test ./pkg/cmd` blocks for 10–30+ seconds every single test run.
- **Circular Dependency Trap:** Monolithic packages tempt developers into tight circular dependencies, making extraction harder over time.

### The Solution: Nuclear Modularization
Decompose the monolith into clean, acyclic single-responsibility packages organized in a strict Directed Acyclic Graph (DAG):
1. **Leaf Packages:** Have zero external domain dependencies; imported by everyone, import nobody else.
2. **Domain Packages:** Cohesive units implementing one clear subsystem (e.g. `cloner`, `cmdprompt`, `cmdpurge`).
3. **Blackbox Heavy Tests:** Slow subprocess tests are segregated into `tests/heavy_test/` (`package heavy_test`), leaving package unit tests lightning fast (< 0.05s).

---

## 2. Centralized Test Inventory Manifest (`.ai-memory/test-inventory.json`)

The test inventory manifest serves as the single source of truth for repository test suite health, categorization, and execution durations. Modeled after the GitMap architecture, it enables intelligent test execution, selective running, and duration tracking.

### Structural JSON Example:
```json
{
  "version": "1.0.0",
  "updated_at": "2026-09-17T02:00:00Z",
  "total_tests": 45,
  "summary": {
    "total_packages": 6,
    "slow_tests": 3,
    "fast_tests": 42,
    "avg_duration_sec": 0.04
  },
  "tests": {
    "tests/heavy_test/TestGitClone_Integration": {
      "id": "TestGitClone_Integration",
      "package": "tests/heavy_test",
      "test_file": "tests/heavy_test/clone_heavy_test.go",
      "target_file": "pkg/cloner/cloner.go",
      "duration_sec": 4.85,
      "tier": "heavy",
      "is_slow": true,
      "last_status": "passed",
      "needs_run": false
    },
    "pkg/fsutil/TestSanitizePath_Success": {
      "id": "TestSanitizePath_Success",
      "package": "pkg/fsutil",
      "test_file": "pkg/fsutil/path_test.go",
      "target_file": "pkg/fsutil/path.go",
      "duration_sec": 0.002,
      "tier": "fast",
      "is_slow": false,
      "last_status": "passed",
      "needs_run": false
    },
    "pkg/cmdprompt/TestFormatPrompt_Valid": {
      "id": "TestFormatPrompt_Valid",
      "package": "pkg/cmdprompt",
      "test_file": "pkg/cmdprompt/prompt_test.go",
      "target_file": "pkg/cmdprompt/prompt.go",
      "duration_sec": 0.008,
      "tier": "fast",
      "is_slow": false,
      "last_status": "passed",
      "needs_run": false
    }
  }
}
```

### Key Field Definitions:

- `id`: The unique test function or suite identifier (e.g. `TestGitClone_Integration`).
- `package`: Relative package import path (e.g. `pkg/fsutil`, `tests/heavy_test`).
- `test_file`: Exact relative path to the test implementation file.
- `target_file`: Primary source file under test.
- `duration_sec`: Profiled execution duration in seconds.
- `tier`: Classification tier (`fast` for in-memory tests < 0.1s; `heavy` for subprocess/network tests).
- `is_slow`: Boolean flag set when `duration_sec >= slow_threshold` (default 4.0s).
- `last_status`: Outcome of last test run (`passed`, `failed`, `skipped`).
- `needs_run`: Dirty flag set when source or test files are modified.

---

## 3. The 5-Day Freshness Decision Engine

To prevent redundant full-test profiling runs that consume precious tokens and CPU cycles, agents MUST evaluate the freshness of `.ai-memory/test-inventory.json`:

```text
                            Check Inventory Freshness
           [python 03-ai-scripts/33-test-inventory-generator.py --check-age]
                                      │
                   ┌──────────────────┴──────────────────┐
                   ▼                                     ▼
        Exit Code 0: FRESH                   Exit Code 1: STALE / MISSING
     (Age <= 5 days & Profiled)              (Age > 5 days or 0.0s Durations)
                   │                                     │
                   ▼                                     ▼
      FAST PATH: Read Cached Timings           PROFILING PASS: Run Baseline
      - Ingest durations directly              - Profile tests once to measure
      - Zero redundant test runs               - Save durations to manifest
      - Plan modularization immediately        - Proceed with cached timings
```

### Execution Rules:

1. **Always Audit First:** Run `python 03-ai-scripts/33-test-inventory-generator.py --check-age --max-age-days 5`.
2. **Fresh Inventory (Exit 0):** The AI agent is STRICTLY FORBIDDEN from running all tests. It MUST read `.ai-memory/test-inventory.json` directly and use existing test durations to identify slow tests and modularization targets.
3. **Stale or Missing Inventory (Exit 1):** The AI agent executes a single baseline inventory generation pass to populate duration metrics, commits the updated manifest, and uses those metrics for subsequent decisions.

---

## 4. Strict Directed Acyclic Graph (DAG) Architecture

To prevent circular dependency errors (`import cycle not allowed`), all packages MUST follow a strict multi-tier hierarchy where packages only import downward:

```text
┌─────────────────────────────────────────────────────────────┐
│                    Layer 4: Entrypoint                      │
│                  main.go, cmd/root.go                       │
└──────────────────────────────┬──────────────────────────────┘
                               │ (imports downward only)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│             Layer 3: Dispatchers & Domain CLI               │
│          cmd/cmdinit, cmd/cmdclone, cmd/cmdpurge            │
└──────────────────────────────┬──────────────────────────────┘
                               │ (imports downward only)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│             Layer 2: Domain Services & Engines              │
│       pkg/cloner, pkg/gitrunner, pkg/config, pkg/audit      │
└──────────────────────────────┬──────────────────────────────┘
                               │ (imports downward only)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│            Layer 1: Pure Leaf Packages (Zero Domain)        │
│    pkg/constants, pkg/model, pkg/appfault, pkg/fsutil       │
└─────────────────────────────────────────────────────────────┘
```

### Isolation Rules:

- **Leaf Packages (Layer 1):** NEVER import Layer 2, Layer 3, or Layer 4 packages. They depend only on standard library or generic utility packages.
- **Domain Packages (Layer 2):** Import Layer 1 leaf packages freely. NEVER import each other if it creates a cycle. Use interfaces for cross-domain communication.
- **Dispatchers (Layer 3):** Wire domain engines and UI dispatchers together. Never contain low-level business logic.
- **Blackbox Heavy Tests (`tests/heavy_test/`):** Declared as `package heavy_test`. Because it is an external test package, it can import any package from Layer 1 through Layer 3 without creating cyclic dependencies inside production packages!

---

## 5. Heavy Test Isolation Architecture (`tests/heavy_test/`)

Integration tests that execute external system processes, spawn CLI subprocesses, access network sockets, or use `time.Sleep` MUST NOT reside in routine unit test files.

### What Qualifies as a Heavy Test?

- Spawns `exec.Command` (e.g. `git`, `bash`, `powershell`, Docker).
- Interacts with disk fixtures creating real Git repositories or directory trees.
- Starts live HTTP / TCP listeners.
- Uses `time.Sleep` with durations > 50ms.
- Runs longer than 0.5s per test case.

### Segregation Protocol:

1. **Move to Dedicated Directory:** Relocate heavy test functions to `tests/heavy_test/<domain>_heavy_test.go` (or `cli/tests/heavy_test/`).
2. **Package Name:** Set the package declaration to `package heavy_test` (not `package <domain>`).
3. **Public API Assertion:** Test packages import the target domain package as an external caller (e.g. `import "coding-guidelines/pkg/cloner"`), validating public contracts cleanly.
4. **Routine Package Fast Tests:** Ensure the in-package `*_test.go` files contain only in-memory, hermetic tests using mocks, test doubles, and parameter structs that execute in < 0.01s.

---

## 6. Per-Task Agent Isolation & Workspace Subfolders (`.ai-memory/temp-agents/xx-<task-name>/`)

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

---

## 7. Function & File Sizing Rules

During modularization and test refactoring:
- **Function Size Cap:** Target <= 8 lines of body logic; hard maximum of <= 15 lines.
- **File Size Cap:** Target <= 80 lines; hard maximum of <= 100 lines (excluding allowed JSON, types, vars, consts, maps exceptions).
- **Blank Lines Before Branching & Returns:** Insert a blank line before `if` statements and before `return` statements.
- **Affirmative Boolean Naming:** All booleans prefixed with `is*` or `has*`. Zero negative checks (`!isSuccess` is banned; use `isFail`). Implicit positive checks only (`if isReady { ... }`).
- **Control Flow Flattening:** Nesting depth MUST remain <= 1. Use early guard returns.

---

## 8. Execution & Build Policy (Strict Ban on Intermediate Runs)

To maintain high throughput and prevent distracting CI noise during modularization:

1. **NO INTERMEDIATE TEST RUNNING (TOTAL BAN):**
   - NEVER run unit test suites (`go test ./...`, `npm test`, `pytest`, runner scripts) during routine file editing turns.
   - All unit test verification is strictly deferred to dedicated QA and CI/CD pipelines.
2. **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):**
   - DO NOT execute build commands (`go build`, `npm run build`, compiler invocations) after editing individual files.
3. **FINAL STEP BUILD VERIFICATION (END OF STEPS ONLY):**
   - Syntax and compilation verification is performed strictly at the **final step** of the modularization run using targeted commands (`go vet ./...`, `go build ./...`) to ensure 0 import errors or circular dependencies before concluding.

---

## 9. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) after individual file edits. Build verification is checked ONLY at the final step.
- [ ] **NO CIRCULAR DEPENDENCIES (TOTAL BAN):** NEVER allow packages to import each other cyclically. Enforce leaf-only downward imports.
- [ ] **NO UNISOLATED SUBPROCESS TESTS (TOTAL BAN):** NEVER leave `exec.Command` or `time.Sleep` integration tests in routine package unit test files.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## 10. Non-Negotiable Coding Guidelines Checklist

- [ ] Strict DAG Architecture: Package dependencies flow strictly downward; zero circular imports.
- [ ] Leaf Packages Segregated: `pkg/constants`, `pkg/model`, `pkg/appfault`, `pkg/fsutil` have zero domain dependencies.
- [ ] Heavy Tests Isolated: Tests with `exec.Command`, git CLI, sockets, or `time.Sleep` placed in `tests/heavy_test/` (`package heavy_test`).
- [ ] Fast Unit Tests Retained: Routine in-package tests run in < 0.05s using mocks and test doubles.
- [ ] Test Inventory Synchronized: `.ai-memory/test-inventory.json` updated with test IDs, packages, and duration categories.
- [ ] Cache Freshness Checked: `python 03-ai-scripts/33-test-inventory-generator.py --check-age` evaluated first; zero redundant full runs on fresh cache.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] File Sizing: Files <= 80 lines preferred (hard cap 100 lines).
- [ ] Whitespace & Style Preserved: Blank line before return, blank line after closing brace `}`.
- [ ] Boolean Conventions: All booleans prefixed with `is` or `has`. Zero negative checks (`!isSuccess` banned; use `isFail`). Implicit positive checks only.
- [ ] Control Flow Flattened: Zero nested `if` statements (nesting depth <= 1). Early guard returns used throughout.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Strict Relative Git Paths: All citations and links use relative paths from repository root.

---

## 11. Task Consolidation & Final Step Git Commit & Push Mandate

### Task Consolidation & File Reduction (End of Loop)
When all subtasks for the parent task (`.ai-memory/plans/pending/xx-nuclear-packages.md`) are finished:
1. Combine all completed granular subtasks from `.ai-memory/plans/subtasks/xx-nuclear-packages/*.md` into `.ai-memory/plans/completed/xx-nuclear-packages.md`.
2. Include a header explicitly documenting initial and optimized package structures, test duration reductions, and loop step metrics.
3. Delete the original granular `.md` files in `.ai-memory/plans/subtasks/xx-nuclear-packages/`.
4. Delete the original parent plan `.ai-memory/plans/pending/xx-nuclear-packages.md`.
5. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.

### Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Metadata

- slug: cg-nuclear-packages
- version: 2.2.0
- category: package-architecture-and-test-performance
- priority: high
- autoloop: true
- languages: [golang, polyglot]
