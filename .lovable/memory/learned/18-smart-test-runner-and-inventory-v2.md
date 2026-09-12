# 18-smart-test-runner-and-inventory-v2.md: Smart Test Runner & Centralized Inventory V2

**Category: Runner & Testing Architecture**
**Updated: 2026-09-12**

## 1. Directory Isolation & Failure Logging Standards
1. **Total Ban on Root `.tmp/` and OS Temp**:
   - All temporary build and test directories must be isolated to `.lovable/temp/`.
   - `03-ai-scripts/06-cicd-local-runner.py` explicitly sets:
     - `TMP_CACHE_DIR = REPO_ROOT / ".lovable" / "temp"`
     - `FAILURES_DIR = TMP_CACHE_DIR / "failures"`
     - Environment variables: `GOTMPDIR`, `TMPDIR`, `TEMP`, and `TMP` bound to `.lovable/temp/`.
2. **Dedicated Failure Folder (`.lovable/temp/failures/`)**:
   - When any test fails, its complete error log and stack trace is written to `.lovable/temp/failures/{test_id}.log`.
   - The directory is created automatically with parent directory creation (`mkdir(parents=True, exist_ok=True)`).
3. **Zero-Trace Passing Tests**:
   - Passing tests must NEVER write files to disk anywhere in the filesystem.
   - Passing test names are NEVER enumerated in output logs.
   - The runner emits only a clean single-line aggregate summary:
     `Passed {passed} tests ({slow} slow [4w x 2], {fast} fast [4w x 4 in 100-chunks]) in {time}s ({cached} tests cached)`.

## 2. Test Inventory Unified Discovery & Relative Path Rules
1. **Centralized Manifest (`.lovable/test-inventory.json`)**:
   - Single Source of Truth maintained by `03-ai-scripts/33-test-inventory-generator.py`.
   - `build_or_update_test_inventory` in `06-cicd-local-runner.py` delegates full discovery to `33-test-inventory-generator.py` and provides high-speed incremental hash checking for subsequent runs.
2. **100% Repository-Relative Paths**:
   - Every entry in `.lovable/test-inventory.json` must have non-empty, forward-slash relative paths for both `target_file` (code) and `test_file` (test).
   - Strict ban on absolute paths or `file:///` URIs across all inventory entries.
3. **Configurable Slow Test Threshold**:
   - Tests taking >4.0s (default, configurable via `--slow-threshold` or `GITMAP_SLOW_TEST_THRESHOLD`) are grouped as slow/heavy.
   - Heavy tests are isolated to `cli/tests/heavy_test/`.
   - Full baseline execution runs with `--force-run-all`.
   - Incremental execution runs only dirty tests where code or test files changed; otherwise skipped as cached.

## 3. Dual Worker Queue Architecture
1. **Queue 1: Slow Tests Pool**:
   - 4 workers running 2 tests at a time (batch size = 2).
   - Prevents I/O and process contention for heavy subprocess and network listener tests.
2. **Queue 2: Fast Tests Pool**:
   - 4 workers running 4 tests at a time.
   - Dispatched in chunks of 100 tests from the test inventory queue.
   - High-throughput execution across thousands of in-memory unit tests.

## 4. Dynamic ETA & AI Agent Sleep Protocol
1. **Dynamic ETA Telemetry (`.lovable/temp/runner-eta.json`)**:
   - Runner calculates total ETA dynamically from duration estimates in the test inventory.
   - Live progress updated on startup, after each 100-test chunk, and at completion.
   - In-flight runner heartbeat cadence is maintained at >=25 seconds.
2. **AI Sleep/Wait Discipline**:
   - When an AI agent inspects an active background test run, it must read `remaining_eta_sec` from `.lovable/temp/runner-eta.json` and sleep/wait for that duration rather than busy-polling.
