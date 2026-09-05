#!/usr/bin/env python3
"""Auto-generated CI/CD local runner with concurrent worker pool, IO mitigation, and selective log aggregation.

Usage:
  python 03-ai-scripts/06-cicd-local-runner.py                 # Run parallel (default), only log failures (if any) or print "All passed"
  python 03-ai-scripts/06-cicd-local-runner.py --all-pass      # Run parallel, show full logs for ALL passed & failed jobs
  python 03-ai-scripts/06-cicd-local-runner.py --all           # Alias for --all-pass
  python 03-ai-scripts/06-cicd-local-runner.py --failed        # Explicitly only show logs for failed jobs (default)
  python 03-ai-scripts/06-cicd-local-runner.py --sync          # Run synchronously / sequentially (1 worker)
  python 03-ai-scripts/06-cicd-local-runner.py -w 8            # Run with 8 concurrent workers
  python 03-ai-scripts/06-cicd-local-runner.py -o report.txt   # Write human-readable test report to file
  python 03-ai-scripts/06-cicd-local-runner.py --json out.json # Write JSON test results to file
"""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import asdict, dataclass
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

# ── Configurable Defaults (Configurable via Environment Variables) ─────────
DEFAULT_WORKERS = int(os.environ.get("CI_MAX_WORKERS", min(8, os.cpu_count() or 4)))
DEFAULT_IO_WORKERS = int(os.environ.get("CI_MAX_IO_WORKERS", 2))
DEFAULT_TIMEOUT_SEC = int(os.environ.get("CI_TIMEOUT_SEC", 300))
DEFAULT_ENCODING = "utf-8"

# ── Environment Configuration ───────────────────────────────────────────────
os.environ.setdefault("CI", "true")
os.environ.setdefault("NODE_ENV", "test")

# ── Job Definitions (Partitioned into Order-Dependent Batches with IO Limits)
JOB_BATCHES: list[dict[str, Any]] = [
    # Batch 1: Linters & AST Checks (Light IO, CPU/AST bound, run highly concurrent)
    {
        "name": "Linters & AST Checks",
        "max_workers": None,  # Inherits global worker limit
        "jobs": {
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
    },
    # Batch 2: Compile & Packaging Gates (Heavy disk IO & RAM, throttled to prevent IO starvation)
    {
        "name": "Compile & Packaging Gates",
        "max_workers": DEFAULT_IO_WORKERS,
        "jobs": {
            "Go Compile Gate": ["go", "build", "-C", "gitmap", "-o", "../bin/gitmap.exe", "."],
            "Web App Build": ["npm", "run", "build"],
        },
    },
    # Batch 3: E2E Smoke Tests (Requires built binary, SQLite single-writer safety)
    {
        "name": "E2E Smoke Tests",
        "max_workers": 1,
        "jobs": {
            "E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],
        },
    },
]

ANSI_ESCAPE_REGEX = re.compile(r"\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])")


def strip_ansi(text: str) -> str:
    """Removes terminal ANSI color escape codes from text."""
    return ANSI_ESCAPE_REGEX.sub("", text)


@dataclass
class JobResult:
    """Encapsulates the execution outcome of an individual quality gate."""
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
    """Executes a single gate subprocess and records duration, return code, and streams."""
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
            encoding=DEFAULT_ENCODING,
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
    except subprocess.TimeoutExpired as exc:
        elapsed = round(time.monotonic() - start, 2)
        out = exc.stdout if isinstance(exc.stdout, str) else ""
        return JobResult(
            name=name,
            cmd=cmd,
            code="timeout",
            out=out,
            err=f"Job timed out after {timeout_sec}s",
            elapsed=elapsed,
        )
    except Exception as exc:
        elapsed = round(time.monotonic() - start, 2)
        return JobResult(
            name=name,
            cmd=cmd,
            code=1,
            out="",
            err=str(exc),
            elapsed=elapsed,
        )


def parse_args() -> argparse.Namespace:
    """Constructs CLI argument parser with comprehensive help and alias support."""
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/06-cicd-local-runner.py",
        description="Concurrent CI/CD local quality gate runner with parallel worker groups, IO protection, and selective log display.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default parallel execution (quiet ticker; only outputs failure logs or "All passed"):
  python 03-ai-scripts/06-cicd-local-runner.py

  # 2. Show detailed logs for all gates (passed and failed):
  python 03-ai-scripts/06-cicd-local-runner.py --all-pass
  python 03-ai-scripts/06-cicd-local-runner.py --all

  # 3. Explicitly only show logs for failed gates:
  python 03-ai-scripts/06-cicd-local-runner.py --failed

  # 4. Run synchronously / sequentially (1 worker):
  python 03-ai-scripts/06-cicd-local-runner.py --sync

  # 5. Tune worker group concurrency and IO throttling:
  python 03-ai-scripts/06-cicd-local-runner.py --workers 8 --io-workers 2

  # 6. Filter by gate name substring:
  python 03-ai-scripts/06-cicd-local-runner.py -k "Linter"

  # 7. Output results to a human-readable text file and JSON file:
  python 03-ai-scripts/06-cicd-local-runner.py -o ci-report.txt --json ci-report.json

Environment Variables:
  CI_MAX_WORKERS      Default max workers for thread pool (default: min(8, cpu_count))
  CI_MAX_IO_WORKERS   Default max workers for heavy IO batch (default: 2)
  CI_TIMEOUT_SEC      Default per-job timeout in seconds (default: 300)
        """,
    )
    parser.add_argument(
        "--all-pass", "--all-passed", "--all-paths", "--all", "-a",
        action="store_true",
        dest="show_all",
        help="Show full stdout and stderr logs for all quality gates (both passed and failed).",
    )
    parser.add_argument(
        "--failed", "-f",
        action="store_true",
        dest="show_failed",
        help="Only display logs for failed quality gates (default behavior).",
    )
    parser.add_argument(
        "--sync", "--sequential", "-s",
        action="store_true",
        dest="sync_mode",
        help="Run quality gates sequentially (1 worker) instead of in parallel.",
    )
    parser.add_argument(
        "--workers", "-w",
        type=int,
        default=DEFAULT_WORKERS,
        help=f"Number of concurrent workers in thread pool (default: {DEFAULT_WORKERS}).",
    )
    parser.add_argument(
        "--io-workers",
        type=int,
        default=DEFAULT_IO_WORKERS,
        help=f"Max workers for heavy IO gates like builds (default: {DEFAULT_IO_WORKERS}).",
    )
    parser.add_argument(
        "--timeout", "-t",
        type=int,
        default=DEFAULT_TIMEOUT_SEC,
        help=f"Per-gate timeout in seconds (default: {DEFAULT_TIMEOUT_SEC}s).",
    )
    parser.add_argument(
        "--filter", "-k",
        type=str,
        default="",
        help="Filter quality gates by case-insensitive name substring.",
    )
    parser.add_argument(
        "--output", "-o", "--output-file",
        type=str,
        default="",
        dest="output_file",
        help="Path to write human-readable test report.",
    )
    parser.add_argument(
        "--json", "--json-output",
        type=str,
        default="",
        dest="json_file",
        help="Path to write machine-readable JSON results report.",
    )
    parser.add_argument(
        "--quiet", "-q",
        action="store_true",
        dest="quiet",
        help="Suppress real-time progress ticker; print only summary or failures.",
    )
    return parser.parse_args()


def filter_job_batches(batches: list[dict[str, Any]], filter_str: str) -> list[dict[str, Any]]:
    """Filters batches to only include gates matching the query substring."""
    if not filter_str:
        return batches
    low_filter = filter_str.lower()
    filtered: list[dict[str, Any]] = []
    for batch in batches:
        matched_jobs = {
            name: cmd
            for name, cmd in batch["jobs"].items()
            if low_filter in name.lower()
        }
        if matched_jobs:
            batch_copy = dict(batch)
            batch_copy["jobs"] = matched_jobs
            filtered.append(batch_copy)
    return filtered


def format_detailed_logs(results: list[JobResult], show_all: bool) -> str:
    """Formats detailed logs for display or file output."""
    lines: list[str] = []
    if show_all:
        lines.append("=" * 60)
        lines.append("                 DETAILED LOGS (--all-pass)")
        lines.append("=" * 60)
        for r in results:
            status_str = "PASS" if r.is_success else ("TIMEOUT" if r.is_timeout else "FAIL")
            lines.append(f"\n[{status_str} LOG] Gate: {r.name} (Duration: {r.elapsed}s) | Exit Code: {r.code}")
            lines.append(f"Command: {' '.join(r.cmd)}")
            if r.out.strip():
                lines.append(f"Stdout:\n{r.out.strip()}")
            if r.err.strip():
                lines.append(f"Stderr:\n{r.err.strip()}")
            lines.append("-" * 60)
        return "\n".join(lines)

    # Default / --failed mode: only format failed or timed-out logs
    failures = [r for r in results if not r.is_success]
    if not failures:
        return ""

    lines.append("=" * 60)
    lines.append("               DETAILED FAILURE LOGS")
    lines.append("=" * 60)
    for r in failures:
        status_str = "TIMEOUT" if r.is_timeout else "FAIL"
        lines.append(f"\n[{status_str} LOG] Gate: {r.name} (Duration: {r.elapsed}s) | Exit Code: {r.code}")
        lines.append(f"Command: {' '.join(r.cmd)}")
        if r.out.strip():
            lines.append(f"Stdout:\n{r.out.strip()}")
        if r.err.strip():
            lines.append(f"Stderr:\n{r.err.strip()}")
        lines.append("-" * 60)

    return "\n".join(lines)


def write_file_reports(
    output_file: str,
    json_file: str,
    summary_text: str,
    detailed_logs: str,
    results: list[JobResult],
    total_elapsed: float,
) -> None:
    """Writes human-readable and JSON reports to disk if requested."""
    if output_file:
        out_path = Path(output_file)
        out_path.parent.mkdir(parents=True, exist_ok=True)
        content = strip_ansi(summary_text)
        if detailed_logs:
            content += "\n\n" + strip_ansi(detailed_logs)
        out_path.write_text(content, encoding=DEFAULT_ENCODING)
        print(f"[INFO] Saved test report to: {output_file}")

    if json_file:
        j_path = Path(json_file)
        j_path.parent.mkdir(parents=True, exist_ok=True)
        payload = {
            "total": len(results),
            "passed": sum(1 for r in results if r.is_success),
            "failed": sum(1 for r in results if not r.is_success and not r.is_timeout),
            "timeouts": sum(1 for r in results if r.is_timeout),
            "duration_sec": total_elapsed,
            "results": [
                {
                    "name": r.name,
                    "code": r.code,
                    "is_success": r.is_success,
                    "is_timeout": r.is_timeout,
                    "elapsed": r.elapsed,
                    "cmd": r.cmd,
                    "out": r.out,
                    "err": r.err,
                }
                for r in results
            ],
        }
        j_path.write_text(json.dumps(payload, indent=2), encoding=DEFAULT_ENCODING)
        print(f"[INFO] Saved JSON report to: {json_file}")


def main() -> None:
    args = parse_args()
    active_batches = filter_job_batches(JOB_BATCHES, args.filter)
    total_jobs = sum(len(b["jobs"]) for b in active_batches)

    if total_jobs == 0:
        print(f"[WARN] No quality gates matched filter: {args.filter!r}")
        sys.exit(0)

    # Determine concurrency
    is_sync = args.sync_mode
    workers = 1 if is_sync else max(1, args.workers)
    io_workers = 1 if is_sync else max(1, args.io_workers)

    mode_label = "Synchronous (Sequential)" if is_sync else f"Parallel ({workers} workers, {io_workers} IO workers)"
    log_mode_label = "All logs (--all-pass)" if args.show_all else "Failed logs only (default)"

    if not args.quiet:
        print(f"[INFO] Enqueued {total_jobs} quality gate(s) across {len(active_batches)} batch(es)")
        print(f"[INFO] Execution Mode : {mode_label}")
        print(f"[INFO] Log Filtering  : {log_mode_label}\n")

    results: list[JobResult] = []
    total_start = time.monotonic()
    job_counter = 0

    # Execute batches sequentially to preserve dependencies & prevent IO lockups
    for batch in active_batches:
        batch_name = batch["name"]
        batch_jobs = batch["jobs"]
        job_items = list(batch_jobs.items())

        # Worker count for current batch
        batch_limit = batch.get("max_workers")
        if is_sync:
            batch_workers = 1
        elif batch_limit is not None:
            batch_workers = min(workers, batch_limit, len(job_items))
        else:
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
                    if not args.quiet:
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
                    if not args.quiet:
                        print(f"  [{job_counter:2d}/{total_jobs}] \033[1;91m✗ FAIL\033[0m [{name}] (Exception: {ex})", flush=True)

    total_elapsed = round(time.monotonic() - total_start, 2)

    passed_count = sum(1 for r in results if r.is_success)
    failed_count = sum(1 for r in results if not r.is_success and not r.is_timeout)
    timeout_count = sum(1 for r in results if r.is_timeout)
    has_failures = (failed_count > 0 or timeout_count > 0)

    # ── Final Consolidated Summary Report ──────────────────────────────────
    summary_lines: list[str] = [
        "=" * 60,
        "           CI/CD EXECUTION SUMMARY REPORT",
        "=" * 60,
        f"Total: {total_jobs} | Passed: {passed_count} | Failed: {failed_count} | Timeouts: {timeout_count} | Time: {total_elapsed}s\n",
    ]
    summary_text = "\n".join(summary_lines)
    print("\n" + summary_text)

    detailed_logs = format_detailed_logs(results, show_all=args.show_all)
    if detailed_logs:
        print(detailed_logs)

    # Write file reports if requested
    write_file_reports(
        args.output_file,
        args.json_file,
        summary_text,
        detailed_logs,
        results,
        total_elapsed,
    )

    if has_failures:
        print(f"\n\033[1;91m[FAILURE]\033[0m CI/CD quality gates failed with {failed_count + timeout_count} error(s).")
        sys.exit(1)
    else:
        print(f"\033[1;92m✓ [SUCCESS] All passed.\033[0m All {total_jobs} quality gates passed successfully! All OK.")
        sys.exit(0)


if __name__ == "__main__":
    main()
