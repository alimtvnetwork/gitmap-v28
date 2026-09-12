# Plan 129: Smart Test Runner, Dual-Queue Worker Pools, and Dynamic ETA Sleep Protocol Synchronization

> **Version:** 1.0.0
> **Scope:** Multi-Repository Synchronization (`gitmap` + `coding-guidelines`)
> **Status:** completed

## Objectives & High-Level Architecture
1. **Directory Isolation**:
   - Test temporary root is strictly `.lovable/temp/` (NOT global system `/tmp` or root `temp/`).
   - Failure output is isolated strictly to `.lovable/temp/failures/<test_id>.log`.
   - Passing tests are 100% silent in both console output and the filesystem (zero logs created).
2. **Centralized Test Inventory (`.lovable/test-inventory.json`)**:
   - Complete repository-relative paths (`cli/...`) for both `target_file` and `test_file`.
   - Tests execute only if target code file or test file has changed (incremental hash dirty-tracking).
   - Configurable slow test threshold: `GITMAP_SLOW_TEST_THRESHOLD` environment variable / repo configuration, defaulting to `4.0s`.
   - First-time run executes all tests to establish baseline timings.
3. **Intelligent Dual-Queue Concurrency**:
   - **Slow Queue**: 4 workers, running at most 2 tests at a time per worker batch.
   - **Fast Queue**: 4 workers, running at most 4 tests at a time, pulling in 100-test chunks from the inventory queue. Once 100 tests complete, the next chunk is enqueued.
4. **Dynamic ETA Sleep Protocol**:
   - `06-cicd-local-runner.py` writes live telemetry and ETA estimates to `.lovable/temp/runner-eta.json` calculated from past runs.
   - AI agent checks `runner-eta.json`, identifies total wait time, and uses sleep/wait protocol rather than spinning in active polling loops.
   - If upon waking the runner is still in progress, the agent re-checks `remaining_eta_sec` and goes back to sleep.
5. **Cross-Repo Prompts, Skills & Python Scripts Synchronization**:
   - Synchronize across `gitmap` and `coding-guidelines` (`D:\wp-work\riseup-asia\coding-guidelines`).

## Task-Specific Constraints
1. **Strict No Releases**: Do NOT bump versions, touch changelog files, or cut git tags.
2. **Strict Relative Paths**: Total ban on hardcoded absolute paths or `file:///` URIs.
3. **Strict No Routine Test Running**: Do NOT run full unit test suites during routine turns unless explicitly instructed.

## Subtask Directory
- [01-verify-and-tune-gitmap-scripts.md](subtasks/129-smart-test-runner-and-eta-sleep-sync/01-verify-and-tune-gitmap-scripts.md)
- [02-update-gitmap-prompts-and-skills.md](subtasks/129-smart-test-runner-and-eta-sleep-sync/02-update-gitmap-prompts-and-skills.md)
- [03-sync-coding-guidelines-scripts.md](subtasks/129-smart-test-runner-and-eta-sleep-sync/03-sync-coding-guidelines-scripts.md)
- [04-sync-coding-guidelines-prompts-and-skills.md](subtasks/129-smart-test-runner-and-eta-sleep-sync/04-sync-coding-guidelines-prompts-and-skills.md)
- [05-quality-gate-verification-and-consolidation.md](subtasks/129-smart-test-runner-and-eta-sleep-sync/05-quality-gate-verification-and-consolidation.md)
