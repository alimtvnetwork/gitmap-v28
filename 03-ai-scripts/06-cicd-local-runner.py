#!/usr/bin/env python3
"""Fast Parallel Multi-Worker Local CI/CD Runner with IO Throttling and Flexible CLI Reporting.

Usage:
  python 03-ai-scripts/06-cicd-local-runner.py                    # Default: quiet on success (tick "✔ All passed."), detailed logs on failure
  python 03-ai-scripts/06-cicd-local-runner.py --all-paths       # Show all information: ticker, summary table, full logs for all gates
  python 03-ai-scripts/06-cicd-local-runner.py --all-passed      # Alias for --all-paths
  python 03-ai-scripts/06-cicd-local-runner.py --all             # Alias for --all-paths
  python 03-ai-scripts/06-cicd-local-runner.py --sync            # Run sequentially (synchronous mode, 1 worker)
  python 03-ai-scripts/06-cicd-local-runner.py -w 4              # Custom worker concurrency
  python 03-ai-scripts/06-cicd-local-runner.py -o report.txt     # Save execution report to file
  python 03-ai-scripts/06-cicd-local-runner.py --json            # Output machine-readable JSON
  python 03-ai-scripts/06-cicd-local-runner.py --json -o out.json# Save JSON results to file
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
        description="Fast Multi-Worker Local CI/CD Runner with parallel worker pool, IO throttling, and flexible reporting.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: run all gates in parallel; quiet on success (tick "✔ All passed."), detailed logs on failure:
  python 03-ai-scripts/06-cicd-local-runner.py

  # 2. Show all information (ticker, summary table, full logs for all gates):
  python 03-ai-scripts/06-cicd-local-runner.py --all-paths
  python 03-ai-scripts/06-cicd-local-runner.py --all-passed
  python 03-ai-scripts/06-cicd-local-runner.py --all

  # 3. Explicitly only show logs for failed gates:
  python 03-ai-scripts/06-cicd-local-runner.py --failed

  # 4. Run sequentially (synchronous mode, 1 worker):
  python 03-ai-scripts/06-cicd-local-runner.py --sync

  # 5. Custom worker concurrency and IO throttling:
  python 03-ai-scripts/06-cicd-local-runner.py --workers 4 --io-workers 2

  # 6. Save report to a file:
  python 03-ai-scripts/06-cicd-local-runner.py --output tmp/cicd-report.txt

  # 7. Output machine-readable JSON (to stdout, or to file with -o):
  python 03-ai-scripts/06-cicd-local-runner.py --json
  python 03-ai-scripts/06-cicd-local-runner.py --json -o tmp/cicd-report.json

  # 8. Filter by gate name substring:
  python 03-ai-scripts/06-cicd-local-runner.py --filter "Linter"

Environment Variables:
  CI_MAX_WORKERS      Default max workers for thread pool (default: min(8, cpu_count))
  CI_MAX_IO_WORKERS   Default max workers for heavy IO batch (default: 2)
  CI_TIMEOUT_SEC      Default per-job timeout in seconds (default: 300)
        """,
    )
    parser.add_argument(
        "--all-paths", "--all-passed", "--all-pass", "--all", "-a",
        action="store_true",
        dest="show_all",
        help="Show detailed information and full logs for all quality gates (both passed and failed).",
    )
    parser.add_argument(
        "--failed", "-f",
        action="store_true",
        dest="show_failed",
        help="Show logs only for failed quality gates (default behavior).",
    )
    parser.add_argument(
        "--sync", "--sequential", "-s",
        action="store_true",
        dest="sync_mode",
        help="Execute quality gates sequentially (1 worker) instead of in parallel.",
    )
    parser.add_argument(
        "--workers", "-w", "--concurrency",
        type=int,
        default=DEFAULT_WORKERS,
        dest="workers",
        help=f"Number of concurrent worker threads (default: {DEFAULT_WORKERS}).",
    )
    parser.add_argument(
        "--io-workers",
        type=int,
        default=DEFAULT_IO_WORKERS,
        dest="io_workers",
        help=f"Max workers for heavy IO gates like builds (default: {DEFAULT_IO_WORKERS}).",
    )
    parser.add_argument(
        "--timeout", "-t",
        type=int,
        default=DEFAULT_TIMEOUT_SEC,
        dest="timeout",
        help=f"Per-gate timeout in seconds (default: {DEFAULT_TIMEOUT_SEC}s).",
    )
    parser.add_argument(
        "--filter", "-k",
        type=str,
        default="",
        dest="filter",
        help="Filter quality gates by case-insensitive name substring.",
    )
    parser.add_argument(
        "--output", "-o", "--file", "--output-file",
        type=str,
        default="",
        dest="output_file",
        help="Save execution results and report to the specified file path.",
    )
    parser.add_argument(
        "--json", "--json-output",
        nargs="?",
        const=True,
        default=False,
        dest="json_mode",
        help="Output results as machine-readable JSON (to stdout, or to file if specified).",
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


def format_full_report(
    results: list[JobResult],
    total_jobs: int,
    passed_count: int,
    failed_count: int,
    timeout_count: int,
    total_elapsed: float,
    show_all: bool,
) -> str:
    """Formats full human-readable summary and logs."""
    lines: list[str] = [
        "=" * 60,
        "           CI/CD EXECUTION SUMMARY REPORT",
        "=" * 60,
    ]
    for r in results:
        status_icon = "✅" if r.is_success else ("⏳" if r.is_timeout else "❌")
        status_word = "PASSED" if r.is_success else ("TIMEOUT" if r.is_timeout else "FAILED")
        lines.append(f"{status_icon} [{status_word}] {r.name:<40} ({r.elapsed:.2f}s)")

    lines.append("-" * 60)
    lines.append(f"Total Duration : {total_elapsed:.2f}s")
    lines.append(f"Gates Passed   : {passed_count}/{total_jobs}")
    lines.append(f"Gates Failed   : {failed_count}/{total_jobs}")
    if timeout_count > 0:
        lines.append(f"Gates Timed Out: {timeout_count}/{total_jobs}")
    lines.append("-" * 60)

    if show_all:
        lines.append("\n" + "=" * 60)
        lines.append("                 ALL QUALITY GATE LOGS")
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
    else:
        failures = [r for r in results if not r.is_success]
        if failures:
            lines.append("\n" + "=" * 60)
            lines.append("               FAILED QUALITY GATE LOGS")
            lines.append("=" * 60)
            for r in failures:
                status_str = "TIMEOUT" if r.is_timeout else "FAIL"
                lines.append(f"\n❌ FAILED: {r.name} (exit code: {r.code}, duration: {r.elapsed}s)")
                lines.append(f"Command: {' '.join(r.cmd)}")
                if r.out.strip():
                    lines.append(f"Stdout:\n{r.out.strip()}")
                if r.err.strip():
                    lines.append(f"Stderr:\n{r.err.strip()}")
                lines.append("-" * 60)

    return "\n".join(lines)


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
    is_json = bool(args.json_mode)

    concurrency_label = "Sequential (1 worker)" if is_sync else f"Parallel ({workers} workers, {io_workers} IO workers)"

    # In --all-paths / --all-passed mode, show banner upfront
    if args.show_all and not is_json:
        print("================================================================")
        print("           PARALLEL LOCAL CI/CD QUALITY GATE RUNNER             ")
        print("================================================================")
        print(f"🚀 Execution Mode          : {concurrency_label}")
        print(f"📋 Total Enqueued Gates    : {total_jobs}")
        print("🔍 Display Mode            : SHOW ALL INFORMATION (--all-paths)")
        print("----------------------------------------------------------------\n")

    results: list[JobResult] = []
    total_start = time.monotonic()
    job_counter = 0

    # Execute batches sequentially to preserve dependencies & prevent IO lockups
    for batch in active_batches:
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
                    if args.show_all and not is_json:
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
                    if args.show_all and not is_json:
                        print(f"  [{job_counter:2d}/{total_jobs}] \033[1;91m✗ FAIL\033[0m [{name}] (Exception: {ex})", flush=True)

    total_elapsed = round(time.monotonic() - total_start, 2)

    passed_count = sum(1 for r in results if r.is_success)
    failed_count = sum(1 for r in results if not r.is_success and not r.is_timeout)
    timeout_count = sum(1 for r in results if r.is_timeout)
    has_failures = (failed_count > 0 or timeout_count > 0)

    # ── Handle JSON Output Mode ────────────────────────────────────────────
    if is_json:
        payload = {
            "total_jobs": total_jobs,
            "passed_count": passed_count,
            "failed_count": failed_count,
            "timeout_count": timeout_count,
            "wall_duration_sec": total_elapsed,
            "has_failures": has_failures,
            "exit_code": 1 if has_failures else 0,
            "gates": [asdict(r) for r in results],
        }
        json_content = json.dumps(payload, indent=2, ensure_ascii=False)
        target_json_path = args.json_mode if isinstance(args.json_mode, str) else args.output_file

        if target_json_path:
            p = Path(target_json_path)
            p.parent.mkdir(parents=True, exist_ok=True)
            p.write_text(json_content, encoding=DEFAULT_ENCODING)
            print(f"📄 JSON results saved to: {target_json_path}")
        else:
            print(json_content)

        sys.exit(1 if has_failures else 0)

    # ── Handle File Output (Human-readable text) ───────────────────────────
    if args.output_file:
        full_file_report = format_full_report(
            results, total_jobs, passed_count, failed_count, timeout_count, total_elapsed, show_all=True
        )
        p = Path(args.output_file)
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(strip_ansi(full_file_report), encoding=DEFAULT_ENCODING)
        print(f"📄 Execution report saved to: {args.output_file}")

    # ── Terminal Output Presentation ───────────────────────────────────────
    if has_failures:
        failure_text = format_full_report(
            results, total_jobs, passed_count, failed_count, timeout_count, total_elapsed, show_all=False
        )
        print(failure_text)
        print(f"\n\033[1;91m[FAILURE]\033[0m CI/CD quality gates failed with {failed_count + timeout_count} error(s).")
        sys.exit(1)

    if args.show_all:
        all_text = format_full_report(
            results, total_jobs, passed_count, failed_count, timeout_count, total_elapsed, show_all=True
        )
        print("\n" + all_text)
        print(f"\n\033[1;92m🎉 All quality gates passed successfully! Codebase is 100% green.\033[0m")
    else:
        # Default mode: clean tick and "All passed."
        print(f"✔ All passed. ({passed_count} gates in {total_elapsed:.2f}s)")

    sys.exit(0)


if __name__ == "__main__":
    main()
