---
name: smart-test-runner-and-inventory
description: Autonomously orchestrate smart incremental test execution, centralized inventory manifests, dual-queue worker pools, failure isolation, and dynamic ETA sleep protocols.
---

# Smart Test Runner & Centralized Inventory Manifest

> **Scope:** Multi-Repository Test Execution & Caching Engine

## Core Architectural Pillars

### 1. Centralized Test Inventory (`.lovable/test-inventory.json`)

- **Strictly Relative Paths:** Both `target_file` and `test_file` MUST be stored as repository-relative forward-slash paths (e.g., `04-code/golang/examples/converter_and_enum_examples.go`, `04-code/golang/examples/converter_and_enum_examples_test.go`). Total ban on absolute paths or `file:///` URIs.
- **First-Time Profiling:** First-time run executes all cataloged tests to record baseline elapsed timings, file hashes, and classification tiers.
- **Incremental Dirty Tracking:** Subsequent test execution is strictly dirty-driven. A test runs if and only if:
  - Its target source file hash has changed (`target_file`), or
  - Its test file hash has changed (`test_file`), or
  - An explicit filter (`--pkg`, `--file`) matches the test.
- Unchanged tests are skipped instantly as `[cached]` (zero overhead).
- **Configurable Slow Threshold:** Default threshold is `4.0s`. Overridable via `GITMAP_SLOW_TEST_THRESHOLD` environment variable or `--slow-threshold` flag.

### 2. Dual-Queue Concurrency Pools

- **Slow Tests Queue:**
  - Worker Pool: 4 concurrent worker threads.
  - Concurrency Cap: At most 2 tests per worker batch.
  - Purpose: Prevents heavy disk or process contention from long-running integration/heavy tests.
- **Fast Tests Queue:**
  - Worker Pool: 4 concurrent worker threads.
  - Concurrency Cap: At most 4 tests per batch.
  - Chunking Engine: Consumes tests from the test inventory queue in chunks of 100 tests. Once a 100-test chunk completes, the next 100 tests are pulled from the inventory queue.

### 3. Temp & Failure Folder Isolation

- **Root Isolation:** All temporary runner artifacts, caches, and test logs MUST reside strictly within `.lovable/temp/`. Creating `temp/` or `.tmp/` at the repository root is strictly prohibited.
- **Failure Folder:** `.lovable/temp/failures/` is the dedicated repository folder for failing tests. When a test fails, its diagnostic output is written to `.lovable/temp/failures/<test-id>.log`.
- **Silent Passing Tests:** Passing tests MUST be 100% silent in both console output and the filesystem (zero files created).

### 4. Dynamic ETA Sleep Protocol

- **Telemetry File:** The runner calculates expected duration from test inventory timings and writes live progress to `.lovable/temp/runner-eta.json` (`status`, `total_eta_sec`, `remaining_eta_sec`, `passed`, `failed`, `completed`).
- **AI Sleep Rule:** AI agents do NOT spin in active polling loops or burn tokens querying status repeatedly. The agent inspects `runner-eta.json`, sleeps for the estimated wait time (or 60 seconds), wakes up, and if the runner is still active, checks remaining ETA and goes back to sleep until completion.

### 5. CLI Usage Cheatsheet

```bash
# Generate / update test inventory
python 03-ai-scripts/33-test-inventory-generator.py

# Record modified files under lock
python 03-ai-scripts/33-test-inventory-generator.py --record "04-code/golang/examples/converter_and_enum_examples.go"

# Run smart incremental tests
python 03-ai-scripts/06-cicd-local-runner.py --pkg 04-code/golang/examples

# Run full CI suite with tests (explicit command only)
python 03-ai-scripts/06-cicd-local-runner.py --run-tests
```
