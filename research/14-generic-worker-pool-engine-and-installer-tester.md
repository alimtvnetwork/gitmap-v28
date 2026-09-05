# Generic Worker Pool Engine & Parallel Script Modernization Architecture

> **Document:** `research/14-generic-worker-pool-engine-and-installer-tester.md`  
> **Status:** Implemented & Verified  
> **Script Reference:** `03-ai-scripts/02-shared-engine.py`, `03-ai-scripts/16-installer-smoke-tester.py`, `03-ai-scripts/28-go-preflight-ci.py`, `03-ai-scripts/06-cicd-local-runner.py`  
> **Date:** 2026-09-05  

---

## 1. Executive Summary & Problem Context

AI automation scripts and validation utilities across repositories typically face two major limitations when executed serially:
1. **Prolonged Execution Latency**: Testing multiple installer scripts, Go modules, or linters in sequence compounds runtime linearly.
2. **Boilerplate Duplication & Inconsistent CLI UX**: Each script invents its own argument parsing, thread pooling, and error formatting, leading to disjointed behavior.

To solve this across all tooling, `03-ai-scripts/02-shared-engine.py` introduces a centralized **Generic Parallel Worker Pool Engine** (`run_worker_pool` and `add_worker_cli_arguments`). Scripts like `16-installer-smoke-tester.py` and `28-go-preflight-ci.py` now reuse this engine to achieve:
- **Adaptive Concurrency**: Spawns worker groups via `ThreadPoolExecutor` scaled to CPU capacity (`min(8, os.cpu_count())`), configurable via `CI_MAX_WORKERS` or `--workers`.
- **Quiet Success Output vs Detailed Failure Traces**:
  - On 100% success (default mode): Quiet execution ending in a clean tick:
    ```text
    ✔ All passed. (4 installer script(s) in 0.02s)
    ```
  - On any failure: Dumps full error messages, outputs, and stack traces.
- **Unified CLI Ergonomics**: Standard `--all-paths` / `--all`, `--sync`, `--output`, `--json`, and `--filter` flags across all scripts.

---

## 2. Architecture of the Shared Worker Pool Engine

Located in `03-ai-scripts/02-shared-engine.py`:

### 2.1 Core Contracts
```python
@dataclass
class WorkerResult:
    name: str
    is_success: bool
    output: str = ""
    error: str = ""
    elapsed_sec: float = 0.0
    metadata: dict[str, Any] = field(default_factory=dict)

@dataclass
class WorkerPoolSummary:
    total_count: int
    passed_count: int
    failed_count: int
    wall_duration_sec: float
    has_failures: bool
    results: list[WorkerResult]
    exit_code: int
```

### 2.2 Reusable CLI Registrator
`add_worker_cli_arguments(parser: ArgumentParser)` automatically equips any script's parser with:
- `--all-paths`, `--all-passed`, `--all-pass`, `--all`, `-a`: Verbose ticker, summary table, and full logs.
- `--failed`, `-f`: Explicit failure log filtering.
- `--sync`, `--sequential`, `-s`: Serial execution (1 worker).
- `--workers`, `-w`, `--concurrency`: Concurrency thread count.
- `--output`, `-o`, `--file`: Path to write sanitized text audit report.
- `--json`, `--json-output`: Outputs machine-readable JSON to stdout or file.
- `--filter`, `-k`: Case-insensitive substring filtering.

### 2.3 Universal Execution Flow (`run_worker_pool`)
```python
def run_worker_pool(
    items: list[Any],
    worker_fn: Callable[[Any], WorkerResult],
    max_workers: int | None = None,
    is_sync: bool = False,
    show_all: bool = False,
    output_file: str | None = None,
    as_json: bool | str = False,
    title: str = "WORKER POOL EXECUTION",
    item_noun: str = "items",
) -> int:
```

---

## 3. Script Modernization

### 3.1 `16-installer-smoke-tester.py`
- Discovers installer scripts (`install.sh`, `install.ps1`).
- Tests for unreplaced placeholder tokens, SHA256 checksum patterns, and safe non-destructive update logic.
- Executes tests in parallel across worker pool.
- Validated output:
  - Default: `✔ All passed. (4 installer script(s) in 0.02s)`
  - Verbose: Real-time ticker and itemized report.

### 3.2 `28-go-preflight-ci.py`
- Discovers all Go modules containing `go.mod`.
- Concurrently executes `go test ./... -count=1` and `golangci-lint run ./...`.
- Validated output: `✔ All passed. (4 Go preflight task(s) in 191.70s)`.

---

## 4. Verification & Quality Gates

All scripts were verified:
1. `python 03-ai-scripts/16-installer-smoke-tester.py`: PASS (4 scripts in 0.02s).
2. `python 03-ai-scripts/16-installer-smoke-tester.py --all-paths`: PASS (full ticker + summary table).
3. `python 03-ai-scripts/16-installer-smoke-tester.py --json`: PASS (valid JSON payload).
4. `python 03-ai-scripts/16-installer-smoke-tester.py --sync`: PASS (sequential mode).
5. `python 03-ai-scripts/06-cicd-local-runner.py`: PASS (all 16 quality gates in 33.4s).
