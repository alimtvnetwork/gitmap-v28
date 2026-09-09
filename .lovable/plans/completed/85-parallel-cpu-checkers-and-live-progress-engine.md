# Master Plan: Parallel Multi-Core Quality Checkers & Real-Time Progress Engine

## 1. Executive Summary & Problem Diagnosis

### The Problem
1. **Extremely Low CPU Utilization (0% - 6% Total System Usage)**:
   - In `03-ai-scripts/26-go-code-formatter.py`, `.github/scripts/go-format-check.py`, and individual linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-error-management.py`), files are scanned or formatted **sequentially on a single thread**.
   - With 1,500+ Go source files and 2,500+ repository files, sequential single-core execution leaves 15+ CPU cores completely idle (0% utilization per core, 6% overall system CPU as captured in Task Manager `media_1788920616880.png`).
   - In `go-format-check.py`, `gofmt -l .` executes sequentially on one core, blocking everything.
   - In `26-go-code-formatter.py`, a `for tf in target_files: format_go_file(tf)` loop invokes `subprocess.run(["gofmt", "-w", ...])` sequentially 1,500+ times.

2. **Stuck at 0% Progress**:
   - `03-ai-scripts/06-cicd-local-runner.py`'s `TelemetryTracker` computes percentage as `int(100.0 * completed_count / total_jobs)`.
   - Because long-running sequential gates (like Go format check and linters) take a long time to complete their first item, `completed_count` remains 0, leaving the runner progress locked at `0%` indefinitely.
   - The individual checkers emit zero intermediate progress percentages during their file processing passes.

### The Solution
1. **Upfront File Listing via Shared Engine / File Manipulator**:
   - Pre-discover and cache all target files (`.go`, `.ts`, `.tsx`, `.py`, `.php`) upfront using `03-ai-scripts/02-shared-engine.py` (`stream_directory_files` / `process_repository_files`) or `03-file-manipulator.py`.
2. **Massive Multi-Core Parallelism (100% CPU Speed)**:
   - Partition file lists into chunks across `os.cpu_count()` workers using `concurrent.futures.ProcessPoolExecutor` / `ThreadPoolExecutor`.
   - In `go-format-check.py` and `26-go-code-formatter.py`, dispatch `gofmt` chunks across all available cores concurrently.
   - In linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-error-management.py`), run multi-threaded AST / regex inspection over chunked file batches.
3. **Real-Time Live Progress Percentage & Throughput Tracking**:
   - Implement incremental progress callbacks emitting live percentage updates (e.g., `[150/1500] 10%`, `[750/1500] 50%`, `[1500/1500] 100%`).
   - Update `06-cicd-local-runner.py` `TelemetryTracker` to display active sub-gate execution and dynamic progress so the UI and terminal visibly progress from 0% to 100%.

---

## 2. Architectural Blueprint

### Component 1: Parallel Go Code Formatter & Checker (`26-go-code-formatter.py` & `go-format-check.py`)
- Upfront file acquisition using `stream_directory_files(repo_root, extensions=[".go"])`.
- Chunk files into balanced slices: `chunk_size = max(1, len(files) // (os.cpu_count() * 4))`.
- Run worker threads with `ThreadPoolExecutor(max_workers=os.cpu_count())` executing `gofmt -w` (or `gofmt -l`) on file slices.
- Atomically track completed files and print dynamic progress bar / percentage:
  `Formatting Go files: [ 450/1500 ] 30% | 16 workers | 420 files/sec`

### Component 2: Parallel Linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-error-management.py`)
- List target files upfront.
- Distribute files across `ThreadPoolExecutor(max_workers=os.cpu_count())`.
- Worker functions scan individual files and return violations.
- Thread-safe progress counter updates stdout every $N$ files or $0.5$ seconds with percentage.

### Component 3: Live Telemetry & Progress in `06-cicd-local-runner.py`
- Enhance `TelemetryTracker` to track both gate-level completion and elapsed active durations.
- Stream runner gate output so sub-gate progress (percentages from formatters/linters) is immediately visible in the terminal.
- Ensure `CI_MAX_WORKERS` defaults to `os.cpu_count()` to saturate available processor threads.

---

## 3. Subtask Decomposition

- `01-parallel-gofmt-formatter-and-checker.md`: Refactor `26-go-code-formatter.py` and `.github/scripts/go-format-check.py` for chunked multi-core parallel execution with live progress reporting.
- `02-parallel-linters-multi-core-engine.md`: Refactor `check-nested-ifs.py` and `check-enum-and-boolean.py` to use multi-core worker pools with file pre-listing and percentage progress bars.
- `03-runner-telemetry-and-cpu-saturation.md`: Optimize `06-cicd-local-runner.py` for full CPU core saturation, live progress rendering, and zero-percent prevention.
- `04-benchmarking-and-verification.md`: Benchmark CPU utilization, verify 100% test pass rate, and validate quality gates.
