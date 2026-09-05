# Transaction Log 19: Enhanced Multi-Worker CI/CD Local Runner CLI with I/O Throttling and File Export

> **Directory:** `05-changes-history/19-enhanced-cicd-local-runner-cli/`  
> **Date:** 2026-09-05  
> **Author/Agent:** Antigravity AI  
> **Module Affected:** `03-ai-scripts/06-cicd-local-runner.py`, `research/`, `05-changes-history/`  
> **Status:** Completed & Verified  

---

## 1. Context & User Directives

The user requested significant enhancements to `03-ai-scripts/06-cicd-local-runner.py`:
```text
python 03-ai-scripts/06-cicd-local-runner.py

I think I requested this on the coding guideline. You could do the same here. So this coding file should, should run by default parallelly with work group all the tests, uh, parallelly as much as possible based on the CPU parallelism and threads. So it would-- There should be variable to change it by default. Um, so it would only print out the failure test if there is any with all the stack traces. However, if there is no failure, it will finally just do a tick and then say, "All passed." Okay? So please make sure that i-it is improved, okay? And also you can follow the coding guideline. If that is already improved, you can use the worker pattern to make sure the coding is good, so that it does not actually creates a, creates a too much pressure on IO. Remember that. And also it should be flexible in both parts. Okay. Do you understand what I'm saying? Can you please help me

Also make it like a CLI. That means we can say all paths, hyphen, hyphen, all paths, then it would show all path information as well. But by default, it would only give the failed information along with the parallel running. Okay? We could also do a sync running that run in sequentially. Uh, we could also output the results to a file. So these are the help screens to be there and be file-based also. Okay. Uh, so make sure that these are all done and any AI can implement or reuse this as well. So this has documentation and help regarding this. Is it clear
```

### Core Requirements
1. **Parallel Execution by Default**:
   - Run tests concurrently across a worker group based on available CPU cores and threads.
   - Configurable defaults (`DEFAULT_WORKERS`, `DEFAULT_IO_WORKERS`, `CI_MAX_WORKERS` environment variable, `--workers` / `-w`).
2. **I/O Pressure Mitigation**:
   - Throttling mechanisms to avoid overwhelming disk I/O and RAM during intensive compilation/bundling steps.
   - Partitioning into batches with dedicated concurrency limits (e.g. `DEFAULT_IO_WORKERS = 2` for `go build` and `npm run build`; 1 worker for SQLite E2E smoke tests).
3. **Selective Log Display**:
   - Default mode: real-time ticker; only output failure logs and full stack traces if a test fails. If all pass, output a clean tick and `"All passed."` (`✓ [SUCCESS] All passed.`).
   - Verbose mode: support `--all-pass`, `--all-passed`, `--all-paths`, `--all`, `-a` to print full stdout and stderr for all passed and failed quality gates.
   - `--failed` / `-f`: explicitly filter to failed logs only.
4. **Synchronous / Sequential Mode**:
   - Support `--sync`, `--sequential`, `-s` to run gates serially (1 worker).
5. **File-Based Output**:
   - Support `--output` / `-o` / `--output-file` to save human-readable text audit reports (ANSI color codes cleanly stripped).
   - Support `--json` / `--json-output` to save structured machine-readable JSON reports.
6. **Comprehensive CLI Help & Reusability**:
   - Rich `--help` screen with practical examples and environment variable documentation.
   - Dedicated architectural documentation in `research/` and audit logs in `05-changes-history/`.

---

## 2. Architectural & Engineering Design

### 2.1 Concurrency Architecture & I/O Protection
The pipeline implements an ordered batch worker pool pattern:
- **Batch 1 (Linters & AST Checks)**: CPU/AST bound with light read-only I/O. Inherits global concurrency limit (`DEFAULT_WORKERS = min(8, cpu_count)`).
- **Batch 2 (Compile & Packaging Gates)**: Heavy disk I/O and RAM footprint (`go build`, `npm run build`). Throttled to `DEFAULT_IO_WORKERS` (default: 2, configurable via `CI_MAX_IO_WORKERS` or `--io-workers`).
- **Batch 3 (E2E Smoke Tests)**: Depends on the compiled binary and executes end-to-end CLI commands on SQLite databases. Bound to 1 worker to ensure single-writer WAL safety and prevent `database is locked (261)` errors.

### 2.2 CLI Ergonomics & Aliases
`parse_args()` in `06-cicd-local-runner.py` provides intuitive options:
- `--all-pass`, `--all-passed`, `--all-paths`, `--all`, `-a`: Verbose logging for all gates.
- `--failed`, `-f`: Explicit failure-only logs.
- `--sync`, `--sequential`, `-s`: Serial execution.
- `--output`, `-o`: File-based report export.
- `--json`, `--json-output`: JSON export.
- `--quiet`, `-q`: Quiet mode suppressing real-time ticker.

### 2.3 ANSI-Free File Logging
Terminal colors look great in console output but corrupt text files. The runner uses `strip_ansi()` with regex `\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])` to sanitize output prior to file writing.

---

## 3. Files Created & Modified

### Modified
1. `03-ai-scripts/06-cicd-local-runner.py` - Complete refactor with batch I/O limits, all-pass aliases, sync mode, file export, and comprehensive CLI help.
2. `gitmap/store/store.go` - Reordered SQLite pragmas so `SQLPragmaBusyTimeout5s` runs before `SQLPragmaJournalWAL`, preventing transient `database is locked (261)` errors.
3. `research/01-index.md` - Registered document 13.
4. `05-changes-history/01-index.md` - Registered transaction log 19.

### Created
1. `research/13-enhanced-cicd-local-runner-cli.md` - Architectural guide, specifications, and reference patterns.
2. `05-changes-history/19-enhanced-cicd-local-runner-cli/01-transaction-log.md` - This transaction log.

---

## 4. Verification & Quality Gate Results

1. **CLI Help Screen**:
   Executed `python 03-ai-scripts/06-cicd-local-runner.py --help` - Verified clean formatting, option flags, examples, and environment variable documentation.

2. **Default Filtered Execution**:
   Executed `python 03-ai-scripts/06-cicd-local-runner.py -k "Newline"`:
   - Output: `✓ [SUCCESS] All passed. All 1 quality gates passed successfully! All OK.`

3. **Verbose `--all-pass` and `--all-paths` Execution**:
   Executed `python 03-ai-scripts/06-cicd-local-runner.py -k "Newline" --all-pass` and `--all-paths`:
   - Output: Printed full command, duration, exit code, and stdout/stderr details for passed gate.

4. **Sequential Mode (`--sync`)**:
   Executed `python 03-ai-scripts/06-cicd-local-runner.py -k "Check" --sync`:
   - Output: Ran 9 gates serially (1 worker) in 7.02s with 100% pass.

5. **File Output (`--output` and `--json`)**:
   Executed `python 03-ai-scripts/06-cicd-local-runner.py -k "Newline" -o tmp/ci-report.txt --json tmp/ci-report.json`:
   - Verified clean ANSI-stripped text report and valid JSON payload.
   - Cleaned up temporary files.

6. **Full Suite Parallel Run**:
   All 16 quality gates pass cleanly across the parallel worker pool.

---

## 5. Next Steps & Hand-off Context

- The local CI/CD runner is fully flexible, robust against I/O pressure, and provides self-documenting CLI flags for both human engineers and autonomous AI agents.
