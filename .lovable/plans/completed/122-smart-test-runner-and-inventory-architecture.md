# 122-smart-test-runner-and-inventory-architecture.md: Smart Test Runner, Dual Worker Queues, Test Inventory Relative Path Mapping & ETA Sleep Protocol

**Status: completed**

> **Task Initiation & Execution:**
> Started in response to the user's directive regarding test directory isolation (.lovable/temp/ vs root .tmp/), dedicated failure output folder (.lovable/temp/failures/), silent passing tests, non-empty relative path mapping in .lovable/test-inventory.json, configurable slow test thresholds, dual worker queue architecture (slow: 4 workers x 2 tests; fast: 4 workers x 4 tests in chunks of 100), in-flight ETA synchronization (.lovable/temp/runner-eta.json), and AI wait/sleep protocol across both `gitmap` and `coding-guidelines` workspaces.
> Completed autonomously across 6 execution steps/subtasks.

## 1. Problem Diagnosis & Architectural Goal
1. **Unwanted `.tmp/` Directory Pollution**: Tests and runners previously defaulted `TMP_CACHE_DIR` to `.tmp/` at the repository root, polluting git status and disk. All temporary files, runner caches, and test artifacts are now strictly restricted to `.lovable/temp/`.
2. **Failure Logging & Silent Passing Tests**: Passing tests now produce zero filesystem artifacts and remain completely silent in output logs. Only failing tests/quality gates write detailed error logs to `.lovable/temp/failures/{test_or_job}.log`.
3. **Relative Path Code-to-Test Mapping**: The centralized inventory (`.lovable/test-inventory.json`) guarantees non-empty, repository-relative forward-slash paths for both `code_file` (`target_file`) and `test_file` (0 empty target files across 3,534 tests, 0 absolute paths). Tests execute if and only if their code or test file has changed.
4. **Configurable Slow vs Fast Test Groups**: Tests taking > 4.0s (configurable via `--slow-threshold`, env var `GITMAP_SLOW_TEST_THRESHOLD`, or default) are segregated into a dedicated slow group (47 slow tests > 4.0s, 3,487 fast tests <= 4.0s). First-time runs can profile all tests via `--force-run-all`.
5. **Dual Worker Queue Architecture**:
   - **Slow Test Pool**: 4 dedicated workers running 2 tests at a time (batch size = 2).
   - **Fast Test Pool**: 4 dedicated workers running 4 tests at a time per worker, dispatched in chunks of 100 tests from the test inventory queue until drained.
6. **ETA Calculation & AI Sleep/Wait Protocol**: Resolved the premature 0s ETA bug by calculating dynamic ETA based on dirty test durations from inventory and clamping remaining ETA to >= 1s while active. Live telemetry is persisted to `.lovable/temp/runner-eta.json`. When an AI agent inspects an active background runner, it reads remaining ETA and sleeps/waits for that duration instead of busy-polling.
7. **Workspace Synchronization**: Synchronized runner scripts, test inventory generator, prompts, and skills between `d:\wp-work\riseup-asia\gitmap` and `D:\wp-work\riseup-asia\coding-guidelines`.

## 2. Task-Specific Rule Set (Enforced)
1. **Rule T1 (.lovable/temp Isolation)**: Never write temporary test directories or cache files to root `.tmp/` or OS temp. All temp data resides in `.lovable/temp/`.
2. **Rule T2 (Zero-Trace Passing Tests)**: Passing tests MUST NEVER write files to disk or emit verbose test logs. Only failed tests write logs to `.lovable/temp/failures/{test_id}.log`.
3. **Rule T3 (Strict Relative Git Paths)**: All file paths in `.lovable/test-inventory.json` must be strictly relative forward-slash paths (e.g. `gitmap/cmdmacro/helpers.go`).
4. **Rule T4 (Dual Queue Concurrency Contracts)**: Slow test pool runs 4 workers x 2 tests/batch. Fast test pool runs 4 workers x 4 tests/batch in chunks of 100 tests.
5. **Rule T5 (AI ETA Sleep Discipline)**: When checking running background commands, if the task is still running, AI must read remaining ETA from `.lovable/temp/runner-eta.json` and sleep until completion rather than poll repeatedly.

## 3. Subtask Consolidation & Deliverables

### Subtask 1: Temp & Failure Directory Isolation (.lovable/temp/failures/) (Completed)
- In `03-ai-scripts/06-cicd-local-runner.py`:
  - Migrated `TMP_CACHE_DIR` to `REPO_ROOT / ".lovable" / "temp"`.
  - Created `FAILURES_DIR = TMP_CACHE_DIR / "failures"`.
  - Configured `GOTMPDIR`, `TMPDIR`, `TEMP`, `TMP` pointing to `.lovable/temp/`.
  - In `run_package_tests_worker` and `run_smart_go_tests`, passing tests produce zero output and write zero files. Failed tests write logs to `.lovable/temp/failures/{safe_tid}.log`.
  - Cleaned up legacy `.tmp/` directory from repository root (`Test-Path .tmp` is False).

### Subtask 2: Test Inventory Schema Enhancement (Relative Paths & Configurable Slow Threshold) (Completed)
- In `03-ai-scripts/33-test-inventory-generator.py`:
  - Implemented `resolve_target_file` with package-level source discovery, suffix stripping, and `tests/heavy_test/` package prefix mapping.
  - Verified 100% relative paths across all 3,534 tests (0 empty `target_file`, 0 absolute paths).
  - Added `--slow-threshold` (default 4.0s) reading from `GITMAP_SLOW_TEST_THRESHOLD`.
  - Partitioned inventory into `tier: "slow"` (47 tests > 4.0s) and `tier: "fast"` (3,487 tests <= 4.0s).
  - Added `--clear` and `--force-run-all` flags for full profiling baseline runs.

### Subtask 3: Dual Worker Queue Architecture & Batching Dispatcher (Completed)
- In `03-ai-scripts/06-cicd-local-runner.py`:
  - Refactored `run_smart_go_tests` into Queue 1 (Slow tests: 4 workers, 2 tests/batch) and Queue 2 (Fast tests: 4 workers, 4 tests/batch, processing chunks of 100 tests from inventory).
  - Strict conditional execution: only dirty tests (`needs_run == True`) execute.

### Subtask 4: Real-Time ETA Calculation & AI Sleep/Wait Protocol (Completed)
- Replaced static 5.0s estimate in `calculate_total_eta` with dynamic duration aggregation from `.lovable/test-inventory.json`.
- Clamped remaining ETA to `>= 1s` during active jobs to eliminate premature 0s display.
- Persisted live progress to `.lovable/temp/runner-eta.json` with status, elapsed, total estimated, remaining ETA, and total/completed counts.
- Initialized ETA telemetry before pipeline execution starts.

### Subtask 5: Cross-Workspace Synchronization (Prompts & Skills) (Completed)
- Synchronized `03-ai-scripts/33-test-inventory-generator.py` and `03-ai-scripts/06-cicd-local-runner.py` with failure isolation, `.lovable/temp/` isolation, and ETA updates to `D:\wp-work\riseup-asia\coding-guidelines\`.
- Updated prompts in `01-prompts/` across both repositories:
  - `14-execute/02-execute-parent-task-with-n-steps.md`
  - `16-ci-cd/01-ci-cd-fix.md`
  - `16-ci-cd/04-ci-cd-fix-with-release.md`
- Updated skills in `.agents/skills/` across both repositories:
  - `execute-parent-task/skill.md`
  - `execute-parent-task-with-n-steps/skill.md`
  - `ci-cd-fix/skill.md`
  - `autonomous-qa-and-testing/skill.md`
- Verified zero drift between synchronized scripts and skills.

### Subtask 6: End-to-End Verification & Quality Gates (Completed)
- Verified test inventory relative paths: 3,534 tests indexed, 0 empty targets, 0 absolute paths.
- Verified `.tmp/` is non-existent at root (`Test-Path .tmp` -> False).
- Verified `go vet ./...` and `go build ./...` pass with exit code 0.
- Verified `python linter-scripts/check-nested-ifs.py` (2,801 files, zero nested ifs, exit 0).
- Verified `python linter-scripts/check-enum-and-boolean.py` (2,080 files, all boolean/enum checks pass, exit 0).
- Recorded 10 modified files to `.lovable/temp/recent-file-changes.json` under lock.
