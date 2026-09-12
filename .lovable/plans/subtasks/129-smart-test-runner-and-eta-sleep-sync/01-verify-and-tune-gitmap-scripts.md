# Subtask 01: Verify and Fine-Tune Gitmap Scripts

> **Parent Plan:** [129-smart-test-runner-and-eta-sleep-sync.md](../../pending/129-smart-test-runner-and-eta-sleep-sync.md)
> **Status:** completed

## Requirements
1. Verify and ensure `03-ai-scripts/06-cicd-local-runner.py`:
   - Failure folder isolation: `.lovable/temp/failures/` is the only directory where failed test logs are written (`<test_id>.log`).
   - Passing tests: 100% silent in both console output and filesystem (zero log files created).
   - Dual-queue worker execution:
     - Slow queue: 4 workers, running at most 2 tests at a time per batch.
     - Fast queue: 4 workers, running at most 4 tests at a time, consuming from the inventory in chunks of 100 tests.
   - Live ETA calculation and telemetry written to `.lovable/temp/runner-eta.json`:
     - Contains `status`, `total_eta_sec`, `remaining_eta_sec`, `completed`, `passed`, `failed`, `total_tests`.
     - Supports `--smart-tests`, `--run-tests`, and `--pkg` filters.
   - Configurable slow test threshold: `GITMAP_SLOW_TEST_THRESHOLD` environment variable, defaulting to `4.0s`.
2. Verify and ensure `03-ai-scripts/33-test-inventory-generator.py`:
   - All tests mapped to repository-relative `target_file` and `test_file` paths (`cli/...`).
   - Slow test threshold configurable via `GITMAP_SLOW_TEST_THRESHOLD` (default `4.0s`).
   - First-time run runs all tests to establish baseline timings.

## Acceptance Criteria
- [x] Runner and test inventory scripts verified and aligned with exact specifications.
- [x] No hardcoded absolute paths exist in `06-cicd-local-runner.py` and `33-test-inventory-generator.py`.
- [x] Test inventory `.lovable/test-inventory.json` uses strictly relative paths.
