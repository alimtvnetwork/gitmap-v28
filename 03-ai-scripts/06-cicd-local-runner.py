#!/usr/bin/env python3
"""Auto-generated CI/CD local runner with concurrent worker pool and selective log aggregation.

Usage:
  python 03-ai-scripts/06-cicd-local-runner.py          # Run parallel, only log failures (if any) or print all OK
  python 03-ai-scripts/06-cicd-local-runner.py --all    # Run parallel, show full logs for ALL jobs
  python 03-ai-scripts/06-cicd-local-runner.py --failed # Explicitly only show logs for failed jobs (default)
  python 03-ai-scripts/06-cicd-local-runner.py -w 8     # Run with 8 concurrent workers
"""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
import os
import shutil
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

# ── Configurable Defaults ───────────────────────────────────────────────────
DEFAULT_WORKERS = min(8, os.cpu_count() or 4)
DEFAULT_TIMEOUT_SEC = 300

# ── Environment Configuration ───────────────────────────────────────────────
os.environ.setdefault("CI", "true")
os.environ.setdefault("NODE_ENV", "test")

# ── Job Definitions (Partitioned into Order-Dependent Sequential Batches) ───
JOB_BATCHES: list[dict[str, list[str]]] = [
    # Batch 1: Linters & AST Checks (order-independent, run concurrently)
    {
        "Spell Check (misspell)": [sys.executable, ".github/scripts/misspell-changed.py"],
        "Nested If Linter": [sys.executable, "linter-scripts/check-nested-ifs.py"],
        "Boolean & Enum Linter": [sys.executable, "linter-scripts/check-enum-and-boolean.py"],
        "Boolean Guidelines Linter": [sys.executable, "linter-scripts/check-boolean-guidelines.py"],
        "Enum Guidelines Linter": [sys.executable, "linter-scripts/check-enum-guidelines.py"],
        "Error Management Check": [sys.executable, "linter-scripts/check-error-management.py"],
        "Relative Path Check": [sys.executable, "linter-scripts/check-relative-paths.py"],
        "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
        "MWS Error Codes Check": [sys.executable, "linter-scripts/check-mws-error-codes.py"],
        "CLI Help Parity Check": [sys.executable, "03-ai-scripts/09-cli-help-auditor.py"],
        "Constants Registry AST Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdRegistryMatchesAST", "-count=1"],
        "Constants Collision Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdConstantsAreUnique", "-count=1"],
        "Helptext Parity Check": ["go", "test", "-C", "gitmap", "./helptext/...", "-count=1"],
    },
    # Batch 2: Compile & Packaging Gates
    {
        "Go Compile Gate": ["go", "build", "-C", "gitmap", "-o", "../bin/gitmap.exe", "."],
        "Web App Build": ["npm", "run", "build"],
    },
    # Batch 3: E2E Smoke Tests (depends on Go Compile Gate binary)
    {
        "E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],
    },
]


@dataclass
class JobResult:
    name: str
    cmd: list[str]
    code: int | str
    out: str
    err: str
    elapsed: float

    @property
    def is_success(self) -> bool:
        return self.code == 0

    @property
    def is_timeout(self) -> bool:
        return self.code == "timeout"


def run_job(name: str, cmd: list[str], timeout_sec: int) -> JobResult:
    start = time.monotonic()
    resolved_cmd = list(cmd)
    binary_path = shutil.which(cmd[0])
    if binary_path is not None:
        resolved_cmd[0] = binary_path

    try:
        result = subprocess.run(
            resolved_cmd,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            timeout=timeout_sec,
        )
        elapsed = round(time.monotonic() - start, 2)
        return JobResult(
            name=name,
            cmd=cmd,
            code=result.returncode,
            out=result.stdout,
            err=result.stderr,
            elapsed=elapsed,
        )
    except subprocess.TimeoutExpired as e:
        elapsed = round(time.monotonic() - start, 2)
        out = e.stdout if isinstance(e.stdout, str) else ""
        return JobResult(
            name=name,
            cmd=cmd,
            code="timeout",
            out=out,
            err=f"Job timed out after {timeout_sec}s",
            elapsed=elapsed,
        )
    except Exception as e:
        elapsed = round(time.monotonic() - start, 2)
        return JobResult(
            name=name,
            cmd=cmd,
            code=1,
            out="",
            err=str(e),
            elapsed=elapsed,
        )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Concurrent CI/CD local quality gate runner with parallel worker groups and selective log filtering.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python 03-ai-scripts/06-cicd-local-runner.py          # Run parallel, show ticker + summary; only log failed jobs (if any)
  python 03-ai-scripts/06-cicd-local-runner.py --all    # Run parallel, print full logs for ALL jobs (passed & failed)
  python 03-ai-scripts/06-cicd-local-runner.py --failed # Explicitly only show logs for failed jobs (default)
  python 03-ai-scripts/06-cicd-local-runner.py -w 8     # Run with 8 concurrent workers
        """,
    )
    parser.add_argument(
        "--all", "-a",
        action="store_true",
        dest="show_all",
        help="Show full stdout and stderr logs for all jobs (passed and failed).",
    )
    parser.add_argument(
        "--failed", "-f",
        action="store_true",
        dest="show_failed",
        help="Only show logs for failed jobs (default behavior).",
    )
    parser.add_argument(
        "--workers", "-w",
        type=int,
        default=DEFAULT_WORKERS,
        help=f"Number of concurrent workers in the thread pool (default: {DEFAULT_WORKERS}).",
    )
    parser.add_argument(
        "--timeout", "-t",
        type=int,
        default=DEFAULT_TIMEOUT_SEC,
        help=f"Per-job timeout in seconds (default: {DEFAULT_TIMEOUT_SEC}s).",
    )
    parser.add_argument(
        "--filter", "-k",
        type=str,
        default="",
        help="Filter quality gates by name substring.",
    )
    return parser.parse_args()


def filter_batches(batches: list[dict[str, list[str]]], filter_str: str) -> list[dict[str, list[str]]]:
    if not filter_str:
        return batches
    low_filter = filter_str.lower()
    filtered: list[dict[str, list[str]]] = []
    for batch in batches:
        matched = {name: cmd for name, cmd in batch.items() if low_filter in name.lower()}
        if matched:
            filtered.append(matched)
    return filtered


def print_detailed_logs(results: list[JobResult], show_all: bool) -> None:
    if show_all:
        print("\n" + "=" * 60)
        print("                 DETAILED LOGS (--all)")
        print("=" * 60)
        for r in results:
            status_str = "PASS" if r.is_success else ("TIMEOUT" if r.is_timeout else "FAIL")
            print(f"\n[{status_str} LOG] Job: {r.name} (Duration: {r.elapsed}s) | Exit Code: {r.code}")
            print(f"Command: {' '.join(r.cmd)}")
            if r.out.strip():
                print(f"Stdout:\n{r.out.strip()}")
            if r.err.strip():
                print(f"Stderr:\n{r.err.strip()}")
            print("-" * 60)
        return

    # Default / --failed mode: only print failed or timed-out logs
    failures = [r for r in results if not r.is_success]
    if not failures:
        return

    print("\n" + "=" * 60)
    print("               DETAILED FAILURE LOGS")
    print("=" * 60)
    for r in failures:
        status_str = "TIMEOUT" if r.is_timeout else "FAIL"
        print(f"\n[{status_str} LOG] Job: {r.name} (Duration: {r.elapsed}s) | Exit Code: {r.code}")
        print(f"Command: {' '.join(r.cmd)}")
        if r.out.strip():
            print(f"Stdout:\n{r.out.strip()}")
        if r.err.strip():
            print(f"Stderr:\n{r.err.strip()}")
        print("-" * 60)


def main() -> None:
    args = parse_args()
    active_batches = filter_batches(JOB_BATCHES, args.filter)
    total_jobs = sum(len(b) for b in active_batches)

    if total_jobs == 0:
        print(f"[WARN] No quality gates matched filter: {args.filter!r}")
        sys.exit(0)

    workers = max(1, args.workers)
    mode_desc = "all logs" if args.show_all else "failed logs only"
    print(f"[INFO] Enqueued {total_jobs} quality gate(s) across {workers} concurrent worker(s) ({len(active_batches)} sequential batch(es)) [Mode: {mode_desc}]...\n")

    results: list[JobResult] = []
    total_start = time.monotonic()
    job_counter = 0

    # Execute each batch; jobs within a batch run concurrently with worker group
    for batch_idx, batch_jobs in enumerate(active_batches, start=1):
        job_items = list(batch_jobs.items())
        batch_workers = min(workers, len(job_items))

        with ThreadPoolExecutor(max_workers=batch_workers) as executor:
            future_to_name = {
                executor.submit(run_job, name, cmd, args.timeout): name
                for name, cmd in job_items
            }

            for future in as_completed(future_to_name):
                job_counter += 1
                try:
                    res = future.result()
                    results.append(res)
                    if res.is_success:
                        print(f"  [{job_counter:2d}/{total_jobs}] \033[1;92m✓ PASS\033[0m [{res.name}] ({res.elapsed}s)", flush=True)
                    elif res.is_timeout:
                        print(f"  [{job_counter:2d}/{total_jobs}] \033[1;93m⏳ TIMEOUT\033[0m [{res.name}] ({res.elapsed}s)", flush=True)
                    else:
                        print(f"  [{job_counter:2d}/{total_jobs}] \033[1;91m✗ FAIL\033[0m [{res.name}] ({res.elapsed}s)", flush=True)
                except Exception as ex:
                    name = future_to_name[future]
                    failed_res = JobResult(
                        name=name,
                        cmd=batch_jobs.get(name, []),
                        code=1,
                        out="",
                        err=str(ex),
                        elapsed=0.0,
                    )
                    results.append(failed_res)
                    print(f"  [{job_counter:2d}/{total_jobs}] \033[1;91m✗ FAIL\033[0m [{name}] (Exception: {ex})", flush=True)

    total_elapsed = round(time.monotonic() - total_start, 2)

    # ── Final Consolidated Summary Report ──────────────────────────────────
    print("\n" + "=" * 60)
    print("           CI/CD EXECUTION SUMMARY REPORT")
    print("=" * 60)

    passed_count = sum(1 for r in results if r.is_success)
    failed_count = sum(1 for r in results if not r.is_success and not r.is_timeout)
    timeout_count = sum(1 for r in results if r.is_timeout)

    print(f"Total: {total_jobs} | Passed: {passed_count} | Failed: {failed_count} | Timeouts: {timeout_count} | Time: {total_elapsed}s\n")

    print_detailed_logs(results, show_all=args.show_all)

    if failed_count > 0 or timeout_count > 0:
        print(f"\n\033[1;91m[FAILURE]\033[0m CI/CD quality gates failed with {failed_count + timeout_count} error(s).")
        sys.exit(1)
    else:
        print(f"\033[1;92m✓ [SUCCESS]\033[0m All {total_jobs} quality gates passed successfully! All OK.")
        sys.exit(0)


if __name__ == "__main__":
    main()

