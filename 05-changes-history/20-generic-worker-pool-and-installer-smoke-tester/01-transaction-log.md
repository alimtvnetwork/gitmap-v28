# Transaction Log 20: Generic Worker Pool Base Engine and Installer Smoke Tester Modernization

> **Directory:** `05-changes-history/20-generic-worker-pool-and-installer-smoke-tester/`  
> **Date:** 2026-09-05  
> **Author/Agent:** Antigravity AI  
> **Module Affected:** `03-ai-scripts/02-shared-engine.py`, `03-ai-scripts/16-installer-smoke-tester.py`, `03-ai-scripts/28-go-preflight-ci.py`, `research/`, `05-changes-history/`  
> **Status:** Completed & Verified  

---

## 1. Context & User Directives

The user requested making the worker pool pattern reusable across AI scripts and refactoring `16-installer-smoke-tester.py` and other test runner scripts:
```text
python 03-ai-scripts\16-installer-smoke-tester.py or 03-ai-scripts\*.* (not all but as many py files which can run with worker process and also the shared file should have the worker process base function so that it can be reused?
 or 03-ai-scripts\ (test runner file) ..

I think I requested this on the coding guideline. You could do the same here. So this coding file should, should run by default parallelly with work group all the tests, uh, parallelly as much as possible based on the CPU parallelism and threads. So it would-- There should be variable to change it by default. Um, so it would only print out the failure test if there is any with all the stack traces. However, if there is no failure, it will finally just do a tick and then say, "All passed." Okay? So please make sure that i-it is improved, okay? And also you can follow the coding guideline. If that is already improved, you can use the worker pattern to make sure the coding is good, so that it does not actually creates a, creates a too much pressure on IO. Remember that. And also it should be flexible in both parts. Okay. Do you understand what I'm saying? Can you please help me

Also make it like a CLI. That means we can say all paths, hyphen, hyphen, all paths, then it would show all path information as well. But by default, it would only give the failed information along with the parallel running. Okay? We could also do a sync running that run in sequentially. Uh, we could also output the results to a file. So these are the help screens to be there and be file-based also. Okay. Uh, so make sure that these are all done and any AI can implement or reuse this as well. So this has documentation and help regarding this. Is it clear
```

### Core Requirements
1. **Shared Worker Pool Engine in `02-shared-engine.py`**:
   - Provide a generic, reusable worker group execution function (`run_worker_pool`) and argument registrar (`add_worker_cli_arguments`).
   - Concurrency auto-detects CPU capacity (`min(8, os.cpu_count() or 4)`) and is configurable via `CI_MAX_WORKERS` or `--workers`.
   - Default quiet mode: silent on passes, prints failures immediately or upon summary, and on 100% pass prints a clean tick: `✔ All passed. (<count> <noun> in <time>s)`.
   - Verbose mode: `--all-paths` / `--all-passed` / `--all` showing banner, real-time ticker, summary table, and itemized logs.
   - Synchronous mode: `--sync` / `--sequential` (1 worker).
   - File output: `--output FILE` for text reports and `--json` for machine-readable JSON.
2. **Modernize `16-installer-smoke-tester.py`**:
   - Refactor to discover installer scripts and execute validations in parallel via the shared worker pool.
   - Support all standard CLI flags.
3. **Modernize `28-go-preflight-ci.py`**:
   - Refactor to run Go tests and linters across modules via the shared worker pool.

---

## 2. Architectural & Engineering Design

### 2.1 Centralized Worker Engine (`03-ai-scripts/02-shared-engine.py`)
- **`WorkerResult`**: Container for task name, success flag, stdout/output, stderr/error, and duration.
- **`WorkerPoolSummary`**: Aggregated run metrics (total, passed, failed, wall duration, failure flag, item list, exit code).
- **`add_worker_cli_arguments`**: Standardized argument registration into any `argparse.ArgumentParser`.
- **`run_worker_pool`**: High-level execution lifecycle managing thread pool workers, future consumption, selective log presentation, and file export.

### 2.2 Reusable Script Modernization
- **`16-installer-smoke-tester.py`**:
  - Encapsulated smoke checks (`test_bash_installer` and `test_powershell_installer`) inside `test_single_installer(script_file: Path) -> WorkerResult`.
  - Dispatches to `run_worker_pool`.
  - Clean quiet tick on success: `✔ All passed. (4 installer script(s) in 0.02s)`.
- **`28-go-preflight-ci.py`**:
  - Dynamically builds `PreflightTask` objects per Go module.
  - Concurrently dispatches tests and linter tasks.
  - Clean quiet tick on success: `✔ All passed. (4 Go preflight task(s) in 191.70s)`.

---

## 3. Files Created & Modified

### Modified
1. `03-ai-scripts/02-shared-engine.py` - Added worker pool engine, data classes, formatting helpers, and CLI argument binder.
2. `03-ai-scripts/16-installer-smoke-tester.py` - Refactored to use shared worker pool.
3. `03-ai-scripts/28-go-preflight-ci.py` - Refactored to use shared worker pool.
4. `research/01-index.md` - Registered document 14.
5. `05-changes-history/01-index.md` - Registered transaction log 20.

### Created
1. `research/14-generic-worker-pool-engine-and-installer-tester.md` - Architectural specification for generic worker pool.
2. `05-changes-history/20-generic-worker-pool-and-installer-smoke-tester/01-transaction-log.md` - This transaction log.

---

## 4. Verification & Quality Gate Results

1. **`16-installer-smoke-tester.py`**:
   - Default: `✔ All passed. (4 installer script(s) in 0.02s)`
   - `--all-paths`: Full ticker, summary table, and detailed output.
   - `--json`: Valid JSON output structure.
   - `--sync`: Sequential run (0.04s).

2. **`28-go-preflight-ci.py`**:
   - `python 03-ai-scripts/28-go-preflight-ci.py test`: `✔ All passed. (4 Go preflight task(s) in 191.70s)`.

3. **CI/CD Quality Gates (`03-ai-scripts/06-cicd-local-runner.py`)**:
   - All 16 quality gates passed successfully.

---

## 5. Next Steps & Hand-off Context

- Any new or existing Python script under `03-ai-scripts/` can now easily import `add_worker_cli_arguments` and `run_worker_pool` from `02-shared-engine.py` to get out-of-the-box parallel execution, I/O protection, quiet tick output, and JSON/file exports.
