# Subtask 03: Synchronize Coding-Guidelines Scripts

> **Parent Plan:** [129-smart-test-runner-and-eta-sleep-sync.md](../../pending/129-smart-test-runner-and-eta-sleep-sync.md)
> **Status:** completed

## Requirements
1. Inspect and synchronize `03-ai-scripts/06-cicd-local-runner.py` in `D:\wp-work\riseup-asia\coding-guidelines`:
   - Bring the intelligent dual-queue worker pool architecture:
     - Slow queue: 4 workers, running at most 2 tests at a time per batch.
     - Fast queue: 4 workers, running at most 4 tests at a time, consuming tests in 100-test chunks.
   - Failure folder isolation: `.lovable/temp/failures/` is the sole sink for failing test logs (`<test_id>.log`).
   - Passing tests: 100% silent in terminal output and filesystem.
   - Dynamic ETA calculation and telemetry written to `.lovable/temp/runner-eta.json` for AI agent sleep protocol.
   - Configurable slow test threshold via environment variable / setting (default `4.0s`).
2. Inspect and synchronize `03-ai-scripts/33-test-inventory-generator.py` in `D:\wp-work\riseup-asia\coding-guidelines`:
   - Repository-relative code-to-test mapping.
   - Configurable slow threshold.
   - First-time run executes all tests to establish baseline timings.

## Acceptance Criteria
- [x] Runner and test inventory scripts in `coding-guidelines` updated to match the dual-queue, 100-chunk, failure isolation, and ETA sleep protocol standards.
- [x] Scripts run cleanly without syntax errors or broken imports.
- [x] No hardcoded absolute paths introduced.
