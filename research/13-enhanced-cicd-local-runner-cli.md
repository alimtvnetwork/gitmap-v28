# Enhanced Multi-Worker CI/CD Local Runner CLI Architecture

> **Document:** `research/13-enhanced-cicd-local-runner-cli.md`  
> **Status:** Implemented & Verified  
> **Script Reference:** `03-ai-scripts/06-cicd-local-runner.py`  
> **Date:** 2026-09-05  

---

## 1. Executive Summary & Context

Local validation of repository quality gates in large monorepos frequently suffers from two competing bottlenecks:
1. **Serial Execution Slowness**: Running 16+ linters, compilers, and end-to-end test suites sequentially takes over a minute.
2. **Uncontrolled Parallel I/O Starvation**: Spawning unrestricted concurrent processes causes disk I/O thrashing, CPU pegging, and transient SQLite database locks (such as WAL checkpoint collisions during smoke testing).

The enhanced `03-ai-scripts/06-cicd-local-runner.py` resolves this dichotomy through:
- **Adaptive Worker Pools**: Auto-detects CPU core count, capped by default at 8 workers and configurable via CLI flags or environment variables.
- **Ordered Batch Partitioning & I/O Throttling**: Partitions quality gates into ordered batches (Linters/AST $\rightarrow$ Compilers $\rightarrow$ E2E Smoke Tests), applying dedicated I/O worker limits on build tasks.
- **Selective Log Display**: Suppresses noise on success (displaying only a clean ticker and final `✓ [SUCCESS] All passed.`), while dumping complete stdout, stderr, and stack traces on failures.
- **Rich CLI Ergonomics**: Full support for `--all-pass` / `--all`, `--failed`, `--sync`, file export (`--output`, `--json`), and query filtering (`--filter`).

---

## 2. CLI Interface & Flag Specifications

### 2.1 Command Usage Matrix

| Flag / Option | Shorthand | Type | Default | Description |
|---|---|---|---|---|
| `--all-pass`, `--all`, `--all-paths` | `-a` | Flag | `False` | Display full stdout and stderr logs for all quality gates (both passed and failed). |
| `--failed` | `-f` | Flag | `True` | Display logs only for failed quality gates (default behavior). |
| `--sync`, `--sequential` | `-s` | Flag | `False` | Run quality gates synchronously in serial sequence (1 worker). |
| `--workers` | `-w` | Integer | `min(8, CPU)` | Concurrency thread count for CPU/AST gate batches. |
| `--io-workers` | - | Integer | `2` | Max concurrency limit for disk-intensive gates (e.g., `go build`, `npm run build`). |
| `--timeout` | `-t` | Integer | `300` | Per-gate execution timeout in seconds. |
| `--filter` | `-k` | String | `""` | Filter quality gates by case-insensitive name substring. |
| `--output`, `--output-file` | `-o` | String | `""` | Output path for formatted human-readable report. |
| `--json`, `--json-output` | - | String | `""` | Output path for structured JSON results report. |
| `--quiet` | `-q` | Flag | `False` | Suppress progress ticker; only display final summary or failure details. |

### 2.2 Environment Variable Overrides
Configuration can be seeded globally or inside containerized CI environments:
- `CI_MAX_WORKERS`: Sets default pool size across batches (e.g., `export CI_MAX_WORKERS=4`).
- `CI_MAX_IO_WORKERS`: Sets maximum concurrency for compile/packaging batches (e.g., `export CI_MAX_IO_WORKERS=1`).
- `CI_TIMEOUT_SEC`: Overrides per-job subprocess timeout.

---

## 3. Architecture & Design Patterns

### 3.1 Batch Partitioning with Dedicated I/O Limits
To prevent concurrency hazards (such as SQLite WAL locks during smoke tests or disk saturation during `npm run build`), gates are grouped into three sequential batches:

```python
JOB_BATCHES = [
    # Batch 1: Linters & AST Checks (Light IO, runs at full worker concurrency)
    {
        "name": "Linters & AST Checks",
        "max_workers": None,  # Uses global worker count
        "jobs": { ... },
    },
    # Batch 2: Compile & Packaging (Heavy IO, throttled to DEFAULT_IO_WORKERS)
    {
        "name": "Compile & Packaging Gates",
        "max_workers": DEFAULT_IO_WORKERS,
        "jobs": { ... },
    },
    # Batch 3: E2E Smoke Tests (Sequential, requires built binary)
    {
        "name": "E2E Smoke Tests",
        "max_workers": 1,
        "jobs": { ... },
    },
]
```

### 3.2 Subprocess Concurrency & GIL Safety
- Concurrency uses Python's `concurrent.futures.ThreadPoolExecutor`.
- Because Python releases the Global Interpreter Lock (GIL) during `subprocess.run`, each thread runs its OS process in true parallel CPU multiprocessing without interpreter contention.
- Outputs are captured into memory (`JobResult`) and printed atomically or selectively aggregated, eliminating interleaved console output.

### 3.3 File Reporting & Machine-Readable Output
- `--output <path>`: Strips terminal ANSI color codes using regex `\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])` and writes a clean text audit report.
- `--json <path>`: Produces a structured JSON document containing execution stats, durations, return codes, and output streams for ingestion by CI pipelines, dashboards, or AI agents.

---

## 4. Verification & Testing Examples

### Default Parallel Run
```bash
python 03-ai-scripts/06-cicd-local-runner.py
```
Output:
```text
[INFO] Enqueued 16 quality gate(s) across 3 batch(es)
[INFO] Execution Mode : Parallel (8 workers, 2 IO workers)
[INFO] Log Filtering  : Failed logs only (default)

  [ 1/16] ✓ PASS [Newline Styling Check] (0.05s)
  ...
✓ [SUCCESS] All passed. All 16 quality gates passed successfully! All OK.
```

### Synchronous Mode with File Export
```bash
python 03-ai-scripts/06-cicd-local-runner.py --sync -o test-report.txt --json test-report.json
```

### Full Verbose Log Mode
```bash
python 03-ai-scripts/06-cicd-local-runner.py --all-pass
```
