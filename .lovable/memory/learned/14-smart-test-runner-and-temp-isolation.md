# Learned: Smart Test Runner, Dual Worker Queues, Temp Isolation & In-Flight ETA Synchronization

## 1. Context & Background
Previously, tests and local runners defaulted `TMP_CACHE_DIR` to `.tmp/` at the repository root, creating uncommitted directories, pollutive artifacts, and cache leakage. Furthermore, passing tests were generating filesystem traces and noisy logs, test inventory entries were missing target file relative paths, slow tests (> 4.0s) were mixed with fast tests, and the CI runner reported premature `[ETA] Estimated remaining time: 0 seconds` while long jobs were actively in flight.

## 2. Key Architecture & Learned Solutions

### A. Temp & Failure Directory Isolation (.lovable/temp/failures/)
- All temporary directories, runner caches, and test run environments are restricted strictly to `.lovable/temp/`.
- Creating or writing to `.tmp/` at the repository root is strictly prohibited.
- `GOTMPDIR`, `TMPDIR`, `TEMP`, and `TMP` environment variables are mapped to `.lovable/temp/`.
- Dedicated failure directory: `.lovable/temp/failures/` is established. Only failing tests/gates write `<test-or-job>.log` there.
- Silent passing tests: Passing tests produce zero filesystem artifacts and remain completely silent in output logs.

### B. Test Inventory Schema & Non-Empty Relative Paths
- In `03-ai-scripts/33-test-inventory-generator.py`:
  - `resolve_target_file` maps each test file to a valid non-empty relative source file path with fallback to package-level Go files or domain directories (e.g. `tests/heavy_test/` -> `gitmap/<domain>/`).
  - Strict forward-slash relative paths guaranteed repository-wide (0 empty target files, 0 absolute paths).
  - Configurable slow test threshold (`--slow-threshold`, default 4.0s, env `GITMAP_SLOW_TEST_THRESHOLD`).
  - Tests segregated into `tier: "slow"` (> 4.0s) and `tier: "fast"` (<= 4.0s).
  - Flags `--clear` and `--force-run-all` support full baseline profiling.

### C. Dual Worker Queue Architecture
- In `03-ai-scripts/06-cicd-local-runner.py`:
  - Slow Test Pool: 4 dedicated workers running 2 tests at a time (batch size = 2).
  - Fast Test Pool: 4 dedicated workers running 4 tests at a time per worker, dispatched in chunks of 100 tests from inventory until drained.
  - Strict conditional execution: tests only run when corresponding code file or test file has changed (`needs_run == True`).

### D. In-Flight Heartbeat (25s) & AI Sleep/Wait Protocol (1 Minute / ETA)
- 25-Second Heartbeat Interval: In `03-ai-scripts/06-cicd-local-runner.py`, `TelemetryTracker` emits in-flight progress heartbeats strictly every 25 seconds or more (default `heartbeat_interval=25.0`, CLI `--heartbeat-interval`, env `RUNNER_HEARTBEAT_INTERVAL`), eliminating terminal spam and fast-loop chatter.
- Dynamic ETA evaluation: computes aggregate duration of dirty tests from `.lovable/test-inventory.json` rather than static 5.0s constants.
- ETA clamping: remaining ETA is clamped to `>= 1s` while jobs are in flight, preventing premature 0s display.
- Live progress persistence: status, elapsed, total estimated, remaining ETA, and total/completed counts streamed to `.lovable/temp/runner-eta.json`.
- AI Sleep Discipline (1 Minute / Dynamic ETA): When an AI agent launches or checks an active background runner, the agent MUST sleep/wait for **1 minute (60 seconds) each time**, or dynamically sleep for the remaining ETA duration read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling.

### E. Cross-Workspace Synchronization
- Changes synchronized between `d:\wp-work\riseup-asia\gitmap` and `D:\wp-work\riseup-asia\coding-guidelines`.
- Prompts (`01-prompts/`) and skills (`.agents/skills/`) updated in both repositories to enforce the temp isolation, failure folder, and ETA wait protocol.
