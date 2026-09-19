#!/usr/bin/env python3
"""
Fast Parallel Multi-Worker Local CI/CD Runner
==============================================
Runs repository quality gates concurrently via a controlled ThreadPoolExecutor worker group.
Provides flexible execution (parallel worker pool vs synchronous sequential),
smart log filtering (silent tick on success, full stack traces on failure),
file-based result logging, machine-readable JSON output, and CLI / AI programmatic interfaces.

Default Behavior:
- Runs in parallel using a worker pool bounded to CPU parallelism (capped at 8 to avoid I/O thrashing).
- Silent on success: prints a single tick line `✔ All passed. (21 gates in 2.34s)`.
- On failure: prints detailed error logs, exit codes, durations, and stack traces.

CLI Options:
- `--all-paths` / `--all-passed` / `--all` / `-a`: Detailed progress ticker, summary table, and full logs.
- `--failed` / `-f`: Show logs only for failed quality gates (default behavior).
- `--sync` / `--sequential` / `-s`: Run quality gates sequentially in a single thread.
- `--workers` / `-w` / `--concurrency`: Worker pool concurrency (default: CPU threads, capped at 8).
- `--output` / `-o` / `--file`: Save formatted execution report (or JSON) to a file.
- `--json`: Output results as machine-readable JSON for automated AI agent workflows.
- `--filter` / `-k`: Run only quality gates matching a substring (case-insensitive).
"""

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import asdict, dataclass
import datetime
import hashlib
from importlib import import_module
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import time
from typing import Any, Callable

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
DEFAULT_MAX_WORKERS = engine.DEFAULT_MAX_WORKERS
CI_JOBS_MATRIX = engine.CI_JOBS_MATRIX

DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS


@dataclass
class JobResult:
    """Represents the execution outcome of an individual CI quality gate."""
    name: str
    is_success: bool
    output: str
    duration_sec: float
    return_code: int


@dataclass
class PipelineSummary:
    """Complete summary of a quality gate pipeline run."""
    total_jobs: int
    passed_count: int
    failed_count: int
    wall_duration_sec: float
    results: list[JobResult]
    has_failures: bool
    exit_code: int


CI_JOB_DEFAULT_ARGS: dict[str, list[str]] = {
    "Sequence & Title Check": ["01-prompts"],
    "Boolean Naming Check": ["04-code"],
    "Misspell Check": ["--staged"],
}


def resolve_job_command(job_name: str, command: list[str]) -> list[str]:
    """Resolves command with safe default targets if not explicitly specified."""
    if job_name in CI_JOB_DEFAULT_ARGS and len(command) <= 2:
        return [*command, *CI_JOB_DEFAULT_ARGS[job_name]]
    return list(command)


REPO_ROOT = Path(__file__).resolve().parent.parent
TMP_CACHE_DIR = REPO_ROOT / ".ai-memory" / "temp"
FAILURES_DIR = TMP_CACHE_DIR / "failures"
RUNNER_ETA_FILE = TMP_CACHE_DIR / "runner-eta.json"
TEST_INVENTORY_PATH = REPO_ROOT / ".ai-memory" / "test-inventory.json"


def compute_file_hash(filepath: Path) -> str:
    """Computes a 16-character SHA-256 hash for a file."""
    if not filepath.is_file():
        return ""
    try:
        data = filepath.read_bytes()
        return hashlib.sha256(data).hexdigest()[:16]
    except OSError:
        return ""


def build_or_update_test_inventory(repo_root: Path, force: bool = False) -> dict[str, Any]:
    """Ensures test inventory exists and is up to date."""
    if TEST_INVENTORY_PATH.is_file() and not force:
        try:
            return json.loads(TEST_INVENTORY_PATH.read_text(encoding="utf-8"))
        except Exception:
            pass
    inv_script = repo_root / "03-ai-scripts" / "33-test-inventory-generator.py"
    if inv_script.is_file():
        cmd = [sys.executable, str(inv_script)]
        if force:
            cmd.append("--force-run-all")
        subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, encoding="utf-8")
        if TEST_INVENTORY_PATH.is_file():
            try:
                return json.loads(TEST_INVENTORY_PATH.read_text(encoding="utf-8"))
            except Exception:
                pass
    return {"tests": {}, "summary": {}}


def run_package_tests_worker(
    pkg: str, pkg_tests: list[dict[str, Any]], repo_root: Path, timeout_sec: int = 120
) -> tuple[int, int, str, dict[str, dict[str, Any]]]:
    """Worker function executing a batch of tests within a package using go test -json."""
    rel_in_go = pkg
    if rel_in_go.startswith("04-code/golang/"):
        rel_in_go = "./" + rel_in_go[len("04-code/golang/"):]
    elif rel_in_go == "04-code/golang":
        rel_in_go = "."
    elif rel_in_go.startswith("cli/"):
        rel_in_go = "./" + rel_in_go[len("cli/"):]
    elif rel_in_go == "cli":
        rel_in_go = "."
    else:
        rel_in_go = f"./{rel_in_go}"

    cmd = ["go", "test", "-json", rel_in_go, "-count=1"]
    test_funcs = [t["test_func"] for t in pkg_tests]
    if len(test_funcs) <= 25:
        run_regex = "^(" + "|".join(test_funcs) + ")$"
        cmd.extend(["-run", run_regex])

    cwd = repo_root / "04-code" / "golang"
    if not cwd.is_dir():
        cwd = repo_root / "cli"
    if not cwd.is_dir():
        cwd = repo_root

    test_env = dict(os.environ)
    test_env["GOTMPDIR"] = str(TMP_CACHE_DIR)
    test_env["TMPDIR"] = str(TMP_CACHE_DIR)
    test_env["TEMP"] = str(TMP_CACHE_DIR)
    test_env["TMP"] = str(TMP_CACHE_DIR)

    try:
        proc = subprocess.run(
            cmd, cwd=cwd, capture_output=True, text=True,
            encoding="utf-8", errors="replace", timeout=timeout_sec,
            env=test_env
        )
    except subprocess.TimeoutExpired:
        for t in pkg_tests:
            fail_log = FAILURES_DIR / f"{t['id'].replace('/', '_')}.log"
            fail_log.write_text(f"Timeout expired after {timeout_sec}s for test {t['id']}", encoding="utf-8")
        return 0, len(pkg_tests), f"Timeout expired after {timeout_sec}s", {}
    except Exception as exc:
        for t in pkg_tests:
            fail_log = FAILURES_DIR / f"{t['id'].replace('/', '_')}.log"
            fail_log.write_text(f"Execution error: {exc}", encoding="utf-8")
        return 0, len(pkg_tests), str(exc), {}

    test_results: dict[str, dict[str, Any]] = {}
    test_output_map: dict[str, list[str]] = {}
    passed = 0
    failed = 0
    raw_stdout = proc.stdout or ""

    for line in raw_stdout.splitlines():
        line_str = line.strip()
        if not line_str:
            continue
        try:
            data = json.loads(line_str)
            action = data.get("Action")
            tname = data.get("Test")
            if tname:
                tid = f"{pkg}.{tname}"
                test_output_map.setdefault(tid, []).append(data.get("Output", ""))
                if action == "pass":
                    passed += 1
                    test_results[tid] = {"status": "passed", "elapsed": float(data.get("Elapsed", 0.0))}
                elif action == "fail":
                    failed += 1
                    test_results[tid] = {"status": "failed", "elapsed": float(data.get("Elapsed", 0.0))}
        except Exception:
            pass

    for t in pkg_tests:
        tid = t["id"]
        if tid not in test_results:
            if proc.returncode == 0:
                passed += 1
                test_results[tid] = {"status": "passed", "elapsed": 0.0}
            else:
                failed += 1
                test_results[tid] = {"status": "failed", "elapsed": 0.0}

    # Write failure logs only for failing tests; passing tests are completely silent
    failure_snippets: list[str] = []
    for tid, res_info in test_results.items():
        if res_info["status"] == "failed":
            fail_log = FAILURES_DIR / f"{tid.replace('/', '_')}.log"
            err_content = "".join(test_output_map.get(tid, [])) or f"Test {tid} failed with exit code {proc.returncode}"
            fail_log.write_text(err_content, encoding="utf-8")
            failure_snippets.append(f"[{tid}] {err_content.strip()}")

    out_summary = "\n".join(failure_snippets) if failed > 0 else ""
    return passed, failed, out_summary, test_results


def filter_tests_by_package_or_file(tests: dict[str, Any], queries: list[str], repo_root: Path) -> list[dict[str, Any]]:
    """Filters inventory tests based on code file paths, code file names, or Go package names."""
    matched_ids: set[str] = set()
    for q_raw in queries:
        q = q_raw.strip().replace("\\", "/").rstrip("/")
        if not q:
            continue
        q_base = os.path.basename(q)
        for tid, t in tests.items():
            pkg = t.get("package", "").replace("\\", "/")
            tf = t.get("target_file", "").replace("\\", "/")
            test_f = t.get("test_file", "").replace("\\", "/")

            if pkg == q or pkg.endswith("/" + q) or pkg.split("/")[-1] == q:
                matched_ids.add(tid)
                continue
            if tf == q or test_f == q or tf.endswith("/" + q) or test_f.endswith("/" + q):
                matched_ids.add(tid)
                continue
            if os.path.basename(tf) == q_base or os.path.basename(test_f) == q_base:
                matched_ids.add(tid)
                continue
    return [tests[tid] for tid in matched_ids if tid in tests]


def run_smart_go_tests(
    name: str = "Go Base Test Suite",
    timeout_sec: int = 120,
    force: bool = False,
    package_filter: list[str] | str | None = None
) -> JobResult:
    """Executes Go tests with dual worker queues (slow: 4w x 2 tests; fast: 4w x 4 tests in 100-chunks)."""
    start_time = time.monotonic()
    repo_root = REPO_ROOT
    inventory = build_or_update_test_inventory(repo_root, force=force)
    all_tests = inventory.get("tests", {})
    tests = {
        k: v for k, v in all_tests.items()
        if v.get("test_file", "").endswith(".go")
    }

    if package_filter:
        queries = [package_filter] if isinstance(package_filter, str) else list(package_filter)
        target_tests = filter_tests_by_package_or_file(tests, queries, repo_root)
        if not target_tests:
            elapsed = round(time.monotonic() - start_time, 2)
            out_msg = f"[WARN] No unit tests found matching package/file query: {', '.join(queries)}"
            return JobResult(
                name=name, is_success=True, output=out_msg, duration_sec=elapsed, return_code=0
            )
        dirty_tests = target_tests
    else:
        dirty_tests = [t for t in tests.values() if t.get("needs_run", True) or force]

    if not dirty_tests:
        elapsed = round(time.monotonic() - start_time, 2)
        total_tests = len(tests)
        out_msg = f"[CACHED] All {total_tests} Go unit tests skipped (0 target functions or tests modified)"
        return JobResult(
            name=name, is_success=True, output=out_msg, duration_sec=elapsed, return_code=0
        )

    slow_threshold = float(os.environ.get("GITMAP_SLOW_TEST_THRESHOLD") or os.environ.get("CG_SLOW_TEST_THRESHOLD", "4.0"))
    slow_tests = [
        t for t in dirty_tests
        if t.get("tier") in ("slow", "heavy")
        or t.get("is_slow", False)
        or float(t.get("duration_sec", 0.0)) >= slow_threshold
    ]
    fast_tests = [t for t in dirty_tests if t not in slow_tests]

    total_dirty = len(dirty_tests)
    passed_count = 0
    failed_count = 0
    error_outputs: list[str] = []

    slow_dur = sum(float(t.get("duration_sec", 4.0)) for t in slow_tests)
    fast_dur = sum(float(t.get("duration_sec", 0.005)) for t in fast_tests)
    slow_eta = slow_dur / 8.0   # 4 workers * 2 tests
    fast_eta = fast_dur / 16.0  # 4 workers * 4 tests
    total_test_eta = max(1.0, round(slow_eta + fast_eta, 1))

    # Live ETA file update
    update_runner_eta("running", 0.0, total_test_eta, total_test_eta, total_dirty, 0)

    # Queue 1: Slow Tests Pool (4 workers, 2 tests per batch)
    if slow_tests:
        slow_batches = [slow_tests[i:i + 2] for i in range(0, len(slow_tests), 2)]
        worker_limit = min(4, len(slow_batches))
        with ThreadPoolExecutor(max_workers=worker_limit) as executor:
            futures = {}
            for batch in slow_batches:
                batch_pkg_map: dict[str, list[dict[str, Any]]] = {}
                for t in batch:
                    batch_pkg_map.setdefault(t["package"], []).append(t)
                for pkg, b_tests in batch_pkg_map.items():
                    fut = executor.submit(run_package_tests_worker, pkg, b_tests, repo_root, timeout_sec)
                    futures[fut] = (pkg, b_tests)

            for fut in as_completed(futures):
                pkg, b_tests = futures[fut]
                try:
                    pkg_passed, pkg_failed, pkg_out, test_results = fut.result()
                    passed_count += pkg_passed
                    failed_count += pkg_failed
                    if pkg_failed > 0:
                        error_outputs.append(f"[{pkg}] {pkg_out}")
                    for tid, res_info in test_results.items():
                        if tid in tests:
                            tests[tid]["duration_sec"] = res_info["elapsed"]
                            tests[tid]["last_status"] = res_info["status"]
                            tests[tid]["last_run_at"] = time.strftime("%Y-%m-%dT%H:%M:%S")
                            tests[tid]["needs_run"] = (res_info["status"] != "passed")
                            tests[tid]["code_hash"] = compute_file_hash(repo_root / tests[tid].get("target_file", ""))
                            tests[tid]["test_hash"] = compute_file_hash(repo_root / tests[tid].get("test_file", ""))
                            if res_info["elapsed"] >= slow_threshold:
                                tests[tid]["is_slow"] = True
                                tests[tid]["tier"] = "slow"
                except Exception as ex:
                    failed_count += len(b_tests)
                    error_outputs.append(f"[{pkg}] Slow worker exception: {ex}")

    # Queue 2: Fast Tests Pool (4 workers, 4 tests per batch, chunks of 100 tests)
    if fast_tests:
        chunk_size = 100
        for chunk_start in range(0, len(fast_tests), chunk_size):
            chunk = fast_tests[chunk_start:chunk_start + chunk_size]
            sub_batches = [chunk[i:i + 4] for i in range(0, len(chunk), 4)]
            worker_limit = min(4, len(sub_batches))
            with ThreadPoolExecutor(max_workers=worker_limit) as executor:
                futures = {}
                for sbatch in sub_batches:
                    batch_pkg_map = {}
                    for t in sbatch:
                        batch_pkg_map.setdefault(t["package"], []).append(t)
                    for pkg, b_tests in batch_pkg_map.items():
                        fut = executor.submit(run_package_tests_worker, pkg, b_tests, repo_root, timeout_sec)
                        futures[fut] = (pkg, b_tests)

                for fut in as_completed(futures):
                    pkg, b_tests = futures[fut]
                    try:
                        pkg_passed, pkg_failed, pkg_out, test_results = fut.result()
                        passed_count += pkg_passed
                        failed_count += pkg_failed
                        if pkg_failed > 0:
                            error_outputs.append(f"[{pkg}] {pkg_out}")
                        for tid, res_info in test_results.items():
                            if tid in tests:
                                tests[tid]["duration_sec"] = res_info["elapsed"]
                                tests[tid]["last_status"] = res_info["status"]
                                tests[tid]["last_run_at"] = time.strftime("%Y-%m-%dT%H:%M:%S")
                                tests[tid]["needs_run"] = (res_info["status"] != "passed")
                                tests[tid]["code_hash"] = compute_file_hash(repo_root / tests[tid].get("target_file", ""))
                                tests[tid]["test_hash"] = compute_file_hash(repo_root / tests[tid].get("test_file", ""))
                                if res_info["elapsed"] >= slow_threshold:
                                    tests[tid]["is_slow"] = True
                                    tests[tid]["tier"] = "slow"
                    except Exception as ex:
                        failed_count += len(b_tests)
                        error_outputs.append(f"[{pkg}] Fast worker exception: {ex}")

            # Update live telemetry after each 100-test chunk
            cur_elapsed = round(time.monotonic() - start_time, 1)
            rem = max(1.0, round(total_test_eta - cur_elapsed, 1))
            update_runner_eta("running", cur_elapsed, total_test_eta, rem, total_dirty, passed_count + failed_count)

    elapsed = round(time.monotonic() - start_time, 2)
    update_runner_eta(
        "completed" if failed_count == 0 else "failed",
        elapsed, total_test_eta, 0.0, total_dirty, passed_count + failed_count
    )

    inventory["summary"]["dirty"] = failed_count
    inventory["summary"]["cached"] = len(tests) - failed_count
    try:
        TEST_INVENTORY_PATH.write_text(json.dumps(inventory, indent=2), encoding="utf-8")
    except Exception:
        pass

    if failed_count > 0:
        err_text = "\n".join(error_outputs)
        return JobResult(
            name=name, is_success=False, output=err_text, duration_sec=elapsed, return_code=1
        )

    out_msg = f"Passed {passed_count} tests ({len(slow_tests)} slow [4w x 2], {len(fast_tests)} fast [4w x 4 in 100-chunks]) in {elapsed}s ({len(tests) - total_dirty} tests cached)"
    return JobResult(
        name=name, is_success=True, output=out_msg, duration_sec=elapsed, return_code=0
    )


def update_runner_eta(status: str, elapsed: float, total_est: float, remaining: float, total_jobs: int, completed: int) -> None:
    """Updates real-time runner status and ETA for AI sleep/wait protocol."""
    try:
        TMP_CACHE_DIR.mkdir(parents=True, exist_ok=True)
        data = {
            "status": status,
            "updated_at": datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            "elapsed_sec": round(elapsed, 1),
            "total_estimated_sec": round(total_est, 1),
            "remaining_eta_sec": max(1 if status == "running" else 0, int(round(remaining))),
            "total_jobs": total_jobs,
            "completed_jobs": completed,
        }
        RUNNER_ETA_FILE.write_text(json.dumps(data, indent=2), encoding=DEFAULT_ENCODING)
    except Exception:
        pass


def execute_ci_job(job_name: str, command: list[str]) -> JobResult:
    """Executes a single validation check asynchronously and records output and duration."""
    if job_name == "Go Base Test Suite":
        return run_smart_go_tests(name=job_name)
    TMP_CACHE_DIR.mkdir(parents=True, exist_ok=True)
    FAILURES_DIR.mkdir(parents=True, exist_ok=True)
    start_time = time.perf_counter()
    effective_cmd = resolve_job_command(job_name, command)
    test_env = os.environ.copy()
    abs_temp = str(TMP_CACHE_DIR.resolve())
    test_env["GOTMPDIR"] = abs_temp
    test_env["TMPDIR"] = abs_temp
    test_env["TEMP"] = abs_temp
    test_env["TMP"] = abs_temp
    try:
        res = subprocess.run(
            effective_cmd,
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace",
            env=test_env,
        )
        duration = time.perf_counter() - start_time
        is_success = (res.returncode == 0)

        output_parts: list[str] = []
        stdout_clean = res.stdout.strip()
        if stdout_clean:
            output_parts.append(stdout_clean)
        stderr_clean = res.stderr.strip()
        if stderr_clean:
            output_parts.append(stderr_clean)

        combined_output = LINE_SEPARATOR.join(output_parts)
        safe_name = re.sub(r'[^a-zA-Z0-9_-]', '_', job_name)
        fail_log = FAILURES_DIR / f"{safe_name}.log"
        if is_success:
            if fail_log.exists():
                try:
                    fail_log.unlink()
                except Exception:
                    pass
        else:
            try:
                fail_log.write_text(combined_output, encoding=DEFAULT_ENCODING, errors="replace")
            except Exception:
                pass

        return JobResult(
            name=job_name,
            is_success=is_success,
            output=combined_output,
            duration_sec=duration,
            return_code=res.returncode
        )
    except Exception as exc:
        duration = time.perf_counter() - start_time
        err_msg = f"Failed to execute process: {exc}"
        safe_name = re.sub(r'[^a-zA-Z0-9_-]', '_', job_name)
        try:
            (FAILURES_DIR / f"{safe_name}.log").write_text(err_msg, encoding=DEFAULT_ENCODING, errors="replace")
        except Exception:
            pass
        return JobResult(
            name=job_name,
            is_success=False,
            output=err_msg,
            duration_sec=duration,
            return_code=-1
        )


def execute_pipeline(
    target_jobs: dict[str, list[str]],
    worker_count: int = DEFAULT_CONCURRENCY_WORKERS,
    is_sync: bool = False,
    on_job_complete: Callable[[JobResult, int, int], None] | None = None
) -> PipelineSummary:
    """Executes target CI jobs in parallel or sequentially, returning structured summary."""
    effective_workers = 1 if is_sync else max(1, worker_count)
    start_wall_time = time.perf_counter()
    est_total_duration = max(5.0, len(target_jobs) * 0.8)
    update_runner_eta("running", 0.0, est_total_duration, est_total_duration, len(target_jobs), 0)
    results: list[JobResult] = []

    if is_sync:
        for name, cmd in target_jobs.items():
            res = execute_ci_job(name, cmd)
            results.append(res)
            update_runner_eta(
                "running",
                time.perf_counter() - start_wall_time,
                est_total_duration,
                max(1.0, est_total_duration - (time.perf_counter() - start_wall_time)),
                len(target_jobs),
                len(results)
            )
            if on_job_complete:
                on_job_complete(res, len(results), len(target_jobs))
    else:
        with ThreadPoolExecutor(max_workers=effective_workers) as executor:
            future_map = {
                executor.submit(execute_ci_job, name, cmd): name
                for name, cmd in target_jobs.items()
            }
            for future in as_completed(future_map):
                res = future.result()
                results.append(res)
                update_runner_eta(
                    "running",
                    time.perf_counter() - start_wall_time,
                    est_total_duration,
                    max(1.0, est_total_duration - (time.perf_counter() - start_wall_time)),
                    len(target_jobs),
                    len(results)
                )
                if on_job_complete:
                    on_job_complete(res, len(results), len(target_jobs))

    total_wall_duration = time.perf_counter() - start_wall_time

    job_order = list(target_jobs.keys())
    results.sort(key=lambda r: job_order.index(r.name) if r.name in job_order else 999)

    passed_count = sum(1 for r in results if r.is_success)
    failed_count = sum(1 for r in results if not r.is_success)
    has_failures = (failed_count > 0)
    exit_code = ExitCodeType.VIOLATIONS_FOUND.value if has_failures else ExitCodeType.SUCCESS.value

    update_runner_eta(
        "completed" if not has_failures else "failed",
        total_wall_duration,
        est_total_duration,
        0.0,
        len(target_jobs),
        len(results)
    )

    return PipelineSummary(
        total_jobs=len(target_jobs),
        passed_count=passed_count,
        failed_count=failed_count,
        wall_duration_sec=total_wall_duration,
        results=results,
        has_failures=has_failures,
        exit_code=exit_code
    )


def format_summary_text(summary: PipelineSummary, show_all: bool = False) -> str:
    """Formats human-readable text summary and error logs."""
    lines: list[str] = []

    if show_all:
        lines.append("======================= FINAL SUMMARY =======================")
        for r in summary.results:
            status_icon = "✅" if r.is_success else "❌"
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"{status_icon} [{status_word}] {r.name:<40} ({r.duration_sec:.2f}s)")
        lines.append("------------------------------------------------------------")
        lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
        lines.append(f"Gates Passed   : {summary.passed_count}/{summary.total_jobs}")
        lines.append(f"Gates Failed   : {summary.failed_count}/{summary.total_jobs}")
        lines.append("------------------------------------------------------------")

        lines.append("\n===================== ALL QUALITY GATE LOGS =====================")
        for r in summary.results:
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"\n--- [{status_word}] {r.name} ({r.duration_sec:.2f}s, exit code: {r.return_code}) ---")
            if r.output:
                lines.append(r.output)
            else:
                lines.append("(no output)")
            lines.append("-----------------------------------------------------------------")
    else:
        if summary.has_failures:
            lines.append("\n================== FAILED QUALITY GATE LOGS ==================")
            for r in summary.results:
                if not r.is_success:
                    lines.append(f"\n❌ FAILED: {r.name} (exit code: {r.return_code}, duration: {r.duration_sec:.2f}s)")
                    lines.append(f"--- {r.name} LOG ---")
                    if r.output:
                        lines.append(r.output)
                    else:
                        lines.append("(no output)")
                    lines.append(f"--- END {r.name} LOG ---")
            lines.append("-----------------------------------------------------------------")
            lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
            lines.append(f"Gates Passed   : {summary.passed_count}/{summary.total_jobs}")
            lines.append(f"Gates Failed   : {summary.failed_count}/{summary.total_jobs}")
            lines.append("------------------------------------------------------------")
            lines.append(f"\n❌ Pipeline failed: {summary.failed_count} quality gate(s) reported violations.")

    return LINE_SEPARATOR.join(lines)


def write_output_file(output_path: str, content: str) -> None:
    """Safely writes report content to output file, creating directories if needed."""
    p = Path(output_path)
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(content, encoding=DEFAULT_ENCODING)


def is_test_job(name: str, command: list[str]) -> bool:
    """Identifies whether a CI quality gate is a test execution suite."""
    name_lower = name.lower()
    if "test" in name_lower or "smoke" in name_lower:
        return True
    for part in command:
        part_lower = part.lower()
        if "test" in part_lower or "pytest" in part_lower:
            return True
    return False


def run_pipeline(
    jobs: dict[str, list[str]] | None = None,
    max_workers: int | None = None,
    show_all: bool = False,
    is_sync: bool = False,
    output_file: str | None = None,
    as_json: bool = False,
    filter_pattern: str | None = None,
    no_tests: bool = False,
    run_tests: bool = False,
) -> int:
    """Dispatches CI jobs using worker group or sequential runner, with selective reporting."""
    all_jobs = jobs or CI_JOBS_MATRIX

    if no_tests and not run_tests:
        all_jobs = {k: v for k, v in all_jobs.items() if not is_test_job(k, v)}

    if filter_pattern:
        pattern_lower = filter_pattern.lower()
        target_jobs = {k: v for k, v in all_jobs.items() if pattern_lower in k.lower()}
        if not target_jobs:
            print(f"⚠️ No CI quality gates matched filter: '{filter_pattern}'")
            return ExitCodeType.TOOL_ERROR.value
    else:
        target_jobs = all_jobs

    worker_count = max_workers or min(len(target_jobs), DEFAULT_CONCURRENCY_WORKERS)
    if is_sync:
        worker_count = 1

    if not as_json:
        if show_all:
            concurrency_label = "Sequential (1 worker)" if is_sync else f"{worker_count} parallel workers"
            print("================================================================")
            print("           PARALLEL LOCAL CI/CD QUALITY GATE RUNNER             ")
            print("================================================================")
            print(f"🚀 Execution Mode          : {concurrency_label}")
            print(f"📋 Total Enqueued Gates    : {len(target_jobs)}")
            print("🔍 Display Mode            : SHOW ALL INFORMATION")
            print("----------------------------------------------------------------\n")

    def ticker_callback(res: JobResult, current: int, total: int):
        if not as_json:
            if show_all:
                status_icon = "✅" if res.is_success else "❌"
                status_label = "PASS" if res.is_success else "FAIL"
                print(f"[{current:2d}/{total:2d}] {status_icon} [{status_label}] {res.name} ({res.duration_sec:.2f}s)")

    summary = execute_pipeline(
        target_jobs=target_jobs,
        worker_count=worker_count,
        is_sync=is_sync,
        on_job_complete=ticker_callback
    )

    if as_json:
        payload = {
            "total_jobs": summary.total_jobs,
            "passed_count": summary.passed_count,
            "failed_count": summary.failed_count,
            "wall_duration_sec": round(summary.wall_duration_sec, 2),
            "has_failures": summary.has_failures,
            "exit_code": summary.exit_code,
            "gates": [asdict(r) for r in summary.results]
        }
        json_content = json.dumps(payload, indent=2, ensure_ascii=False)
        if output_file:
            write_output_file(output_file, json_content)
            print(f"📄 JSON results saved to: {output_file}")
        else:
            print(json_content)
        return summary.exit_code

    if output_file:
        full_file_report = format_summary_text(summary, show_all=True)
        write_output_file(output_file, full_file_report)
        print(f"📄 Execution report saved to: {output_file}")

    if summary.has_failures:
        failure_text = format_summary_text(summary, show_all=False)
        print(failure_text)
        return summary.exit_code

    if show_all:
        all_text = format_summary_text(summary, show_all=True)
        print(all_text)
        print("\n🎉 All quality gates passed successfully! Codebase is 100% green.")
    else:
        print(f"✔ All passed. ({summary.passed_count} gates in {summary.wall_duration_sec:.2f}s)")

    return summary.exit_code


def parse_arguments() -> argparse.Namespace:
    """Parses command-line arguments for CI runner."""
    parser = argparse.ArgumentParser(
        description="Fast Multi-Worker Local CI/CD Runner with parallel worker pool and flexible reporting.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Default: run all gates in parallel; quiet on success (tick), detailed logs on failure:
  python 03-ai-scripts/06-cicd-local-runner.py

  # Show all information (ticker, summary table, full logs for all gates):
  python 03-ai-scripts/06-cicd-local-runner.py --all-paths
  python 03-ai-scripts/06-cicd-local-runner.py --all-passed
  python 03-ai-scripts/06-cicd-local-runner.py --all

  # Run sequentially (synchronous mode, 1 worker):
  python 03-ai-scripts/06-cicd-local-runner.py --sync

  # Custom worker concurrency:
  python 03-ai-scripts/06-cicd-local-runner.py --workers 4

  # Save report to a file:
  python 03-ai-scripts/06-cicd-local-runner.py --output tmp/cicd-report.txt

  # Output machine-readable JSON:
  python 03-ai-scripts/06-cicd-local-runner.py --json -o tmp/cicd-report.json

  # Filter specific gate:
  python 03-ai-scripts/06-cicd-local-runner.py --filter "Go Base"
        """
    )
    parser.add_argument(
        "--all", "--all-paths", "--all-passed", "-a",
        action="store_true",
        dest="show_all",
        help="Show detailed information and logs for all quality gates (both passed and failed)."
    )
    parser.add_argument(
        "--failed", "-f",
        action="store_true",
        dest="show_failed",
        help="Show logs only for failed quality gates (default behavior)."
    )
    parser.add_argument(
        "--sync", "--sequential", "-s",
        action="store_true",
        dest="is_sync",
        help="Execute quality gates sequentially (1 worker) instead of in parallel."
    )
    parser.add_argument(
        "--workers", "-w", "--concurrency",
        type=int,
        default=None,
        dest="workers",
        help="Number of concurrent worker threads (default: CPU threads, capped at 8 to prevent I/O thrashing)."
    )
    parser.add_argument(
        "--output", "-o", "--file",
        type=str,
        default=None,
        dest="output_file",
        help="Save execution results and report to the specified file path."
    )
    parser.add_argument(
        "--json",
        action="store_true",
        dest="as_json",
        help="Output results as machine-readable JSON (useful for automated AI agents)."
    )
    parser.add_argument(
        "--filter", "-k",
        type=str,
        default=None,
        dest="filter",
        help="Filter jobs matching substring (case-insensitive)."
    )
    parser.add_argument(
        "--no-tests", "--skip-tests",
        action="store_true",
        dest="no_tests",
        help="Strictly exclude all test suites and self-tests from execution (for non-release and prompt tasks)."
    )
    parser.add_argument(
        "--run-tests", "--with-tests",
        action="store_true",
        dest="run_tests",
        help="Explicitly enable test suite execution (for release verification or when explicitly requested by owner)."
    )
    parser.add_argument(
        "--changed-only",
        action="store_true",
        dest="changed_only",
        help="Run quality gates scoped strictly to files modified in .ai-memory/temp/recent-file-changes.json."
    )
    parser.add_argument(
        "--pkg", "--package", "-p", "--target-file",
        nargs="*",
        dest="package_filter",
        help="Run specific Go test package or file based on relative path or package name."
    )
    return parser.parse_args()


CICD_LAST_RUN_CACHE = Path(".ai-memory/cicd/last_run_cache.json")


def check_recent_run_cache(cache_file: Path, signature: str, normal_ttl: float = 15.0) -> int | None:
    """Returns cached exit code if an identical runner run occurred within the debounce window."""
    if not cache_file.exists():
        return None
    try:
        data = json.loads(cache_file.read_text(encoding=DEFAULT_ENCODING))
        last_time = float(data.get("timestamp", 0.0))
        elapsed = time.time() - last_time
        if elapsed < normal_ttl and data.get("signature") == signature:
            print("================================================================")
            print("Here is the result from the previous run.")
            print("================================================================")
            print(f"⏱️  Cached from previous run {elapsed:.1f}s ago (debounce TTL: {normal_ttl:.0f}s).")
            status_text = "PASSED (exit 0)" if data.get("exit_code") == 0 else f"FAILED (exit {data.get('exit_code')})"
            print(f"📋 Status: {status_text}")
            summary = data.get("summary")
            if summary:
                print(f"📊 Summary: {summary}")
            print("================================================================")
            return int(data.get("exit_code", 0))
    except Exception:
        return None
    return None


def save_recent_run_cache(cache_file: Path, signature: str, exit_code: int, summary: str = "") -> None:
    """Persists recent run result to debounce cache."""
    try:
        cache_file.parent.mkdir(parents=True, exist_ok=True)
        data = {
            "timestamp": time.time(),
            "signature": signature,
            "exit_code": exit_code,
            "summary": summary,
        }
        cache_file.write_text(json.dumps(data, indent=2), encoding=DEFAULT_ENCODING)
    except Exception:
        pass


def clean_stale_locks(max_age_sec: float = 300.0) -> None:
    """Removes stale lock files older than max_age_sec to prevent pipeline deadlocks."""
    lock_paths = [
        Path(".ai-memory/temp/recent-file-changes.lock"),
        Path(".ai-memory/temp/active-locks.json"),
    ]
    for lock_file in lock_paths:
        if not lock_file.exists():
            continue
        try:
            mtime = lock_file.stat().st_mtime
            if (time.time() - mtime) > max_age_sec:
                if lock_file.suffix == ".lock":
                    lock_file.unlink(missing_ok=True)
                elif lock_file.suffix == ".json":
                    try:
                        data = json.loads(lock_file.read_text(encoding=DEFAULT_ENCODING))
                        if isinstance(data, dict):
                            now = time.time()
                            filtered = {
                                k: v for k, v in data.items()
                                if isinstance(v, dict) and (now - float(v.get("timestamp", now))) < max_age_sec
                            }
                            if len(filtered) != len(data):
                                lock_file.write_text(json.dumps(filtered, indent=2), encoding=DEFAULT_ENCODING)
                    except Exception:
                        pass
        except Exception:
            pass


RECENT_CHANGES_FILE = Path(".ai-memory/temp/recent-file-changes.json")


def filter_jobs_for_changed_files(all_jobs: dict[str, list[str]]) -> dict[str, list[str]] | None:
    """Filters CI jobs based on modified files tracked in .ai-memory/temp/recent-file-changes.json."""
    if not RECENT_CHANGES_FILE.exists():
        return None
    try:
        data = json.loads(RECENT_CHANGES_FILE.read_text(encoding=DEFAULT_ENCODING))
        changed_files: list[str] = []
        if isinstance(data, list):
            changed_files = [str(f) for f in data]
        elif isinstance(data, dict):
            changed_files = [str(f) for f in data.get("files", [])]

        if not changed_files:
            return {}

        has_go = any(f.endswith(".go") or "golang" in f for f in changed_files)
        has_py = any(f.endswith(".py") or "linter-scripts" in f or "03-ai-scripts" in f for f in changed_files)
        has_md = any(f.endswith(".md") or "02-spec" in f for f in changed_files)
        has_web = any("slides-app" in f or "src" in f or f.endswith((".ts", ".tsx", ".jsx", ".js")) for f in changed_files)

        filtered_jobs: dict[str, list[str]] = {}
        for name, cmd in all_jobs.items():
            name_lower = name.lower()
            if "markdown" in name_lower or "spec" in name_lower or "link" in name_lower:
                if has_md:
                    filtered_jobs[name] = cmd
            elif "go " in name_lower or "golang" in name_lower:
                if has_go:
                    filtered_jobs[name] = cmd
            elif "boolean" in name_lower or "error" in name_lower or "nested if" in name_lower:
                if has_go or has_py or has_web:
                    filtered_jobs[name] = cmd
            elif "web" in name_lower or "slides" in name_lower:
                if has_web:
                    filtered_jobs[name] = cmd
            elif "python" in name_lower or "linter" in name_lower:
                if has_py:
                    filtered_jobs[name] = cmd
            else:
                filtered_jobs[name] = cmd

        return filtered_jobs
    except Exception:
        return None


def main():
    args = parse_arguments()

    clean_stale_locks()

    sig = f"no_tests={getattr(args, 'no_tests', False)},run_tests={getattr(args, 'run_tests', False)},filter={getattr(args, 'filter', '') or ''},changed_only={getattr(args, 'changed_only', False)}"
    cached_code = check_recent_run_cache(CICD_LAST_RUN_CACHE, sig)
    if cached_code is not None:
        sys.exit(cached_code)

    if getattr(args, "package_filter", None):
        res = run_smart_go_tests(name="Go Package Tests", package_filter=args.package_filter)
        if res.is_success:
            print(f"✔ {res.output}")
            sys.exit(0)
        else:
            print(f"❌ {res.output}")
            sys.exit(1)

    target_jobs = None
    if args.changed_only:
        scoped = filter_jobs_for_changed_files(CI_JOBS_MATRIX)
        if scoped is not None:
            if not scoped:
                print("================================================================")
                print("ℹ️  --changed-only: No modified files requiring CI validation.")
                print("================================================================")
                save_recent_run_cache(CICD_LAST_RUN_CACHE, sig, 0, "No relevant files modified")
                sys.exit(0)
            target_jobs = scoped

    exit_code = run_pipeline(
        jobs=target_jobs,
        max_workers=args.workers,
        show_all=args.show_all,
        is_sync=args.is_sync,
        output_file=args.output_file,
        as_json=args.as_json,
        filter_pattern=args.filter,
        no_tests=args.no_tests,
        run_tests=args.run_tests,
    )
    summary_str = f"Result code: {exit_code}"
    save_recent_run_cache(CICD_LAST_RUN_CACHE, sig, exit_code, summary_str)
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
