# 127-smart-test-runner-and-inventory-v2.md: Smart Test Runner & Centralized Inventory V2

**Status: completed**

> **Task Initiation & Execution:**
> Completed in response to the user's directive regarding test directory isolation (`.lovable/temp/` vs root `.tmp/`), dedicated failure output folder (`.lovable/temp/failures/`), zero-trace passing tests, relative path code-to-test mapping in `.lovable/test-inventory.json`, configurable slow test thresholds, dual worker queue architecture (slow: 4 workers x 2 tests; fast: 4 workers x 4 tests in chunks of 100), dynamic in-flight ETA synchronization (`.lovable/temp/runner-eta.json`), AI sleep/wait protocol, and dual-repo synchronization across `gitmap` and `coding-guidelines`.

## 1. Problem Diagnosis & Architecture Resolutions

1. **Temp Directory & Failure Log Isolation**:
   - Total ban on OS temp directories (`AppData/Local/Temp` or `/tmp`) and repository root `.tmp/`.
   - `06-cicd-local-runner.py` explicitly binds `TMP_CACHE_DIR = REPO_ROOT / ".lovable" / "temp"` and `FAILURES_DIR = TMP_CACHE_DIR / "failures"`.
   - `execute_subprocess`, `run_package_tests_worker`, and module-level executions enforce `GOTMPDIR`, `TMPDIR`, `TEMP`, and `TMP` pointing strictly to `.lovable/temp/`.
   - Only failed tests write logs to `.lovable/temp/failures/{test_id}.log`.

2. **Zero-Trace Passing Tests**:
   - Passing tests write ZERO files to disk anywhere in the filesystem.
   - Passing test names are NEVER enumerated in output logs.
   - Output summary produces a single aggregate line: `Passed X tests (Y slow [4w x 2], Z fast [4w x 4 in 100-chunks]) in Ts (N tests cached)`.

3. **Unified Test Discovery & 100% Relative Paths**:
   - Connected `06-cicd-local-runner.py`'s `build_or_update_test_inventory` to centralized `33-test-inventory-generator.py`.
   - Indexes all 3,534 tests across 134 packages across the entire repository (not just `cli/`).
   - 100% of tests have non-empty, repository-relative paths for both `target_file` (code) and `test_file` (test) — 0 absolute paths, 0 empty strings.
   - Preserves all summary metrics: `slow_tests`, `fast_tests`, `heavy_tests`, `unit_tests`, `slow_threshold_sec`, `estimated_slow_sec`, `estimated_fast_sec`.

4. **Configurable Slow vs Fast Groups**:
   - Tests taking >4.0s (configurable via repo/env `GITMAP_SLOW_TEST_THRESHOLD` and `--slow-threshold`) are categorized as slow/heavy.
   - 50 slow tests (>4.0s) strictly isolated to `cli/tests/heavy_test/`.
   - 3,484 fast unit tests (<=4.0s) running purely in-memory in <0.01s.
   - First-time runs profile all tests via `--force-run-all`.
   - Subsequent turns execute strictly dirty tests where code or test files changed; otherwise skipped as cached.

5. **Dual Worker Queue Architecture**:
   - **Slow Test Pool**: 4 dedicated workers running 2 tests at a time (batch size = 2).
   - **Fast Test Pool**: 4 dedicated workers running 4 tests at a time, processing chunks of 100 tests from the test inventory queue.

6. **Dynamic ETA Calculation & AI Sleep Protocol**:
   - Dynamic ETA computed from duration estimates stored in test inventory.
   - Live telemetry written to `.lovable/temp/runner-eta.json` with status, total_eta_sec, elapsed_sec, remaining_eta_sec, slow/fast totals, completed, passed, and failed counts.
   - In-flight heartbeat minimum 25s.
   - AI agents inspecting background test runs read `remaining_eta_sec` from `.lovable/temp/runner-eta.json` and sleep/wait for that duration instead of busy-polling.

7. **Cross-Workspace Synchronization**:
   - `03-ai-scripts/33-test-inventory-generator.py` identical across `gitmap` and `coding-guidelines`.
   - `03-ai-scripts/06-cicd-local-runner.py` uses `.lovable/temp/failures/` and `runner-eta.json` in both repositories.
   - `01-prompts/` catalog in 100% parity across both workspaces.

## 2. Subtask Consolidation & Deliverables

- **Subtask 1: Verify Temp Failures Directory & Silent Passing Tests**: Completed. `.lovable/temp/failures/` isolation, env redirection (`GOTMPDIR`, `TMPDIR`, `TEMP`, `TMP`), zero-trace passing tests.
- **Subtask 2: Verify Test Inventory Relative Paths & Groups**: Completed. Integrated `06-cicd-local-runner.py` with `33-test-inventory-generator.py`. 3,534 tests indexed with 100% relative paths and configurable threshold.
- **Subtask 3: Verify Dual Worker Queue & ETA Sleep Protocol**: Completed. Dual worker queues (4w x 2 slow; 4w x 4 fast in 100-chunks), dynamic ETA calculation, and `runner-eta.json` updates.
- **Subtask 4: Synchronize Scripts, Prompts & Skills to Coding-Guidelines**: Completed. Prompts synchronized, scripts aligned, and working tree verified clean.
- **Subtask 5: Quality Gate Verification & Consolidation**: Completed. All quality linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-relative-paths.py`, `test_ci_scripts.py`, `gofmt -l cli/`) passed cleanly.
