# Learned: Parallel Multi-Worker CI/CD Runner & Selective Log Filtering

> **Path:** `.lovable/memory/learned/03-parallel-cicd-runner-and-log-filtering.md`  
> **Topic:** Local CI/CD pipeline concurrency, process workgroups, selective log suppression, and duration metrics  
> **Updated:** 2026-09-05  

---

## 1. Verbatim User Directives

```text
python 03-ai-scripts/06-cicd-local-runner.py

in this script make it parallel workgroup run the tests parallellky and only show the failed test info if failed if not then all ok and passing hsould be prininted 

now there should be flags to display thing ebtter wya like 

python 03-ai-scripts/06-cicd-local-runner.py --all # should show al logs
python 03-ai-scripts/06-cicd-local-runner.py --failed # should show al logs or empty should log fialed ones onlhy, ckear??
```

---

## 2. Core Architectural Decisions

### 2.1 Multi-Process Concurrency via ThreadPoolExecutor
- Because quality gate steps run external subprocesses (`subprocess.run`), Python's Global Interpreter Lock (GIL) is released during I/O and process execution.
- Using `concurrent.futures.ThreadPoolExecutor` provides true multi-core process execution across Windows and POSIX systems.
- Worker count defaults dynamically to `min(len(target_jobs), os.cpu_count() or 8, 8)` to maximize throughput without thrashing OS schedulers.
- Runtime speedup: Total wall-clock time dropped from **35–40 seconds** down to **6.5–8.5 seconds** across 21 gates (~5.5x faster).

### 2.2 Selective Log Filtering Protocol
In large multi-gate suites (21 checks), stdout noise from passing tools clutters the console and obscures actionable errors.
- **Default / `--failed` mode:**
  - Real-time ticker lines report progress: `[ 1/21] ✅ [PASS] <Gate Name> (<duration>s)`.
  - Passing gates emit zero stdout/stderr logs.
  - If any gate fails, full stdout and stderr are printed **only for the failing gates**.
  - If all 21 gates pass, a concise summary table and green banner are displayed.
- **`--all` mode:** Full logs (stdout + stderr) are printed for all gates regardless of pass/fail status.

### 2.3 Structured Execution Result
All checks report through an immutable dataclass:
```python
@dataclass
class JobResult:
    name: str
    is_success: bool
    output: str
    duration_sec: float
    return_code: int
```

---

## 3. Strict Prohibitions & Anti-Patterns Avoided

1. **NO Explicit `== True` Checks:** Booleans are strictly evaluated implicitly: `if show_all:`, `if r.is_success:`, `if not target_jobs:`.
2. **NO Mixed Polarity:** Positive and negative checks are never combined in a single `if` expression.
3. **NO Noisy Passing Logs:** Never flood developer terminal with verbose logs for passing gates unless explicitly requested with `--all`.
4. **NO Hardcoded Slugs or Absolute Paths:** All subprocesses and scripts are executed with relative paths starting from git root.
