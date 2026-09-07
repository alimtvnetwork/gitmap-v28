#!/usr/bin/env python3
"""Fast Parallel Multi-Worker Local CI/CD Runner with Real-Time Failure Streaming & Terminal Feedback.

AI AGENT REAL-TIME CI/CD INSTRUCTIONS:
  1. Live Error Streaming: All failures are streamed to .lovable/temp/cicd/errors.log immediately as they occur.
  2. Immediate Terminal Output: Full stack traces and failing commands print to stdout instantly upon failure.
  3. Parallel Remediation: You do NOT have to wait for the entire runner to finish. Inspect .lovable/temp/cicd/errors.log
     and start fixing detected failures while background gates continue running.
  4. Quiet Passes: Passing gates stay quiet by default (single summary tick) so terminal context stays clean.
  5. Machine-Readable Failures: Real-time JSON error list is maintained at .lovable/temp/cicd/errors.json.

Usage:
  python 03-ai-scripts/06-cicd-local-runner.py                    # Quiet on success (tick "✔ All passed."), live failure stream
  python 03-ai-scripts/06-cicd-local-runner.py --all-paths       # Show all gates (passed and failed) with detailed summary
  python 03-ai-scripts/06-cicd-local-runner.py --sync            # Run sequentially (synchronous mode, 1 worker)
  python 03-ai-scripts/06-cicd-local-runner.py -w 4              # Custom worker concurrency
  python 03-ai-scripts/06-cicd-local-runner.py -o report.txt     # Save execution report to file
  python 03-ai-scripts/06-cicd-local-runner.py --json            # Output machine-readable JSON
"""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import asdict, dataclass
from importlib import import_module
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import threading
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

# ── Configurable Defaults (Configurable via Environment Variables) ─────────
DEFAULT_WORKERS = int(os.environ.get("CI_MAX_WORKERS", min(8, os.cpu_count() or 4)))
DEFAULT_IO_WORKERS = int(os.environ.get("CI_MAX_IO_WORKERS", 2))
DEFAULT_TIMEOUT_SEC = int(os.environ.get("CI_TIMEOUT_SEC", 1200))
DEFAULT_ENCODING = "utf-8"

# ── Environment Configuration ───────────────────────────────────────────────
os.environ.setdefault("CI", "true")
os.environ.setdefault("NODE_ENV", "test")
TMP_CACHE_DIR = Path(__file__).resolve().parent.parent.parent / ".tmp"
TMP_CACHE_DIR.mkdir(parents=True, exist_ok=True)
os.environ.setdefault("GOTMPDIR", str(TMP_CACHE_DIR))

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
            "Interface Naming Check": [sys.executable, "linter-scripts/check-interface-naming.py"],
            "CLI Help Parity Check": [sys.executable, "03-ai-scripts/09-cli-help-auditor.py"],
            "Constants Registry AST Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdRegistryMatchesAST", "-count=1"],
            "Constants Collision Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdConstantsAreUnique", "-count=1"],
            "Helptext Parity Check": ["go", "test", "-C", "gitmap", "./helptext/...", "-count=1"],
            "Lint Script Unit Tests": [sys.executable, ".github/scripts/tests/test_ci_scripts.py"],
            "govulncheck": [sys.executable, ".github/scripts/check-vulncheck.py"],
            "Startup Build-Tags (linux)": {"cmd": ["go", "build", "./startup/..."], "cwd": "gitmap", "env": {"GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
            "Startup Build-Tags (darwin)": {"cmd": ["go", "build", "./startup/..."], "cwd": "gitmap", "env": {"GOOS": "darwin", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
            "Startup Build-Tags (windows)": {"cmd": ["go", "build", "./startup/..."], "cwd": "gitmap", "env": {"GOOS": "windows", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
                                    "golangci-lint (strict)": {"cmd": ["golangci-lint", "run", "--issues-exit-code=1", "--timeout=10m", "-c", ".golangci.yml", "--path-prefix", "gitmap", "./..."], "cwd": "gitmap"},
            "Cross-OS Vet (Windows)": {"cmd": ["go", "vet", "-C", "gitmap", "./..."], "env": {"GOOS": "windows", "GOARCH": "amd64"}},
            "Cross-OS Vet (Darwin)": {"cmd": ["go", "vet", "-C", "gitmap", "./..."], "env": {"GOOS": "darwin", "GOARCH": "amd64"}},
        },
    },
    # Batch 2: Compile & Packaging Gates (Heavy disk IO & RAM, throttled to prevent IO starvation)
    {
        "name": "Compile & Packaging Gates",
        "max_workers": DEFAULT_IO_WORKERS,
        "jobs": {
            "Go Compile Gate": ["go", "build", "-C", "gitmap", "-o", "../bin/gitmap.exe", "."],
            "Web App Build": ["npm", "run", "build"],
            "GoReleaser Snapshot Build": {"cmd": ["go", "run", "github.com/goreleaser/goreleaser/v2@latest", "release", "--snapshot", "--clean", "--parallelism=1"], "cwd": "gitmap"},
        },
    },
    # Batch 3: E2E Smoke Tests (Requires built binary, SQLite single-writer safety)
    {
        "name": "E2E Smoke Tests",
        "max_workers": 1,
        "jobs": {
            "E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],
            "Installer Smoke (source)": [sys.executable, ".github/scripts/smoke-installer.py", "source"],
            "Installer Smoke (release)": [sys.executable, ".github/scripts/smoke-installer.py", "release"],
            "History Purge Smoke": [sys.executable, ".github/scripts/smoke-history-purge.py", "bin/gitmap.exe"],
            "History Pin Smoke": [sys.executable, ".github/scripts/smoke-history-pin.py", "bin/gitmap.exe"],
            "Go Test Coverage Profile": ["go", "test", "-C", "gitmap", "-count=1", "-coverprofile=../coverage.out", "./..."],
            "Coverage Floor Guard": [sys.executable, ".github/scripts/coverage-floor.py", "coverage.out"],
            "Go Test Race (Hot Packages)": ["go", "test", "-C", "gitmap", "-count=1", "-timeout=15m", "./cmd/...", "./cloneconcurrency/...", "./visibility/...", "./store/...", "./uipref/..."],
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



GLOBAL_TIMINGS: dict[str, float] = {}


def record_job_timing(job_name: str, elapsed_sec: float) -> None:
    """Stores the execution duration of a job for persistence."""
    GLOBAL_TIMINGS[job_name] = elapsed_sec


def run_job(name: str, cmd: list[str], timeout_sec: int, env: dict[str, str] = None, cwd: str = None) -> JobResult:
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
            env=env,
            cwd=cwd,
        )
        elapsed = round(time.monotonic() - start, 2)
        record_job_timing(name, elapsed)
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
        record_job_timing(name, elapsed)
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
        record_job_timing(name, elapsed)
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
        "--eta-interval",
        type=int,
        default=120,
        help="Print remaining ETA every N seconds (0 to disable). Default: 120.",
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


TIMING_FILE_PATH = Path(".lovable/cicd-timing.json")
DEFAULT_JOB_ESTIMATE_SEC = 10.0


def handle_json_output(
    args: argparse.Namespace,
    results: list[JobResult],
    counts: tuple[int, int, int, int],
    wall_duration_sec: float,
    has_failures: bool,
) -> int:
    """Handles JSON serialization and writing to stdout or target file."""
    total_jobs, passed_count, failed_count, timeout_count = counts
    payload = {
        "total_jobs": total_jobs,
        "passed_count": passed_count,
        "failed_count": failed_count,
        "timeout_count": timeout_count,
        "wall_duration_sec": wall_duration_sec,
        "has_failures": has_failures,
        "exit_code": 1 if has_failures else 0,
        "gates": [asdict(r) for r in results],
    }
    json_content = json.dumps(payload, indent=2, ensure_ascii=False)
    target_path = args.json_mode if isinstance(args.json_mode, str) else args.output_file
    if target_path:
        p = Path(target_path)
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(json_content, encoding=DEFAULT_ENCODING)
        print(f"📄 JSON results saved to: {target_path}")
    else:
        print(json_content)

    return 1 if has_failures else 0


def handle_text_output(
    args: argparse.Namespace,
    results: list[JobResult],
    counts: tuple[int, int, int, int],
    total_elapsed: float,
) -> int:
    """Handles text report formatting, output file saving, and console printing."""
    total_jobs, passed_count, failed_count, timeout_count = counts
    has_failures = bool(failed_count > 0 or timeout_count > 0)
    if args.output_file:
        file_report = format_full_report(
            results, total_jobs, passed_count, failed_count, timeout_count, total_elapsed, show_all=True
        )
        p = Path(args.output_file)
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(strip_ansi(file_report), encoding=DEFAULT_ENCODING)
        print(f"📄 Execution report saved to: {args.output_file}")

    if has_failures:
        failure_text = format_full_report(
            results, total_jobs, passed_count, failed_count, timeout_count, total_elapsed, show_all=False
        )
        print(failure_text)
        print(f"\n\033[1;91m[FAILURE]\033[0m CI/CD quality gates failed with {failed_count + timeout_count} error(s).")
        return 1

    if args.show_all:
        all_text = format_full_report(
            results, total_jobs, passed_count, failed_count, timeout_count, total_elapsed, show_all=True
        )
        print("\n" + all_text)
        print(f"\n\033[1;92m🎉 All quality gates passed successfully! Codebase is 100% green.\033[0m")
    else:
        print(f"✔ All passed. ({passed_count} gates in {total_elapsed:.2f}s)")

    return 0


CICD_TEMP_DIR = Path(".lovable/temp/cicd")
CICD_ERRORS_LOG = CICD_TEMP_DIR / "errors.log"
CICD_ERRORS_JSON = CICD_TEMP_DIR / "errors.json"
CICD_RUN_LOG = CICD_TEMP_DIR / "run.log"
CICD_SUMMARY_JSON = CICD_TEMP_DIR / "summary.json"


def init_cicd_temp_stream() -> None:
    """Initializes streaming directories and empty log files under .lovable/temp/cicd/."""
    CICD_TEMP_DIR.mkdir(parents=True, exist_ok=True)
    ts = time.strftime("%Y-%m-%d %H:%M:%S")
    CICD_ERRORS_LOG.write_text(f"# CI/CD Real-Time Failure Stream — {ts}\n\n", encoding=DEFAULT_ENCODING)
    CICD_ERRORS_JSON.write_text("[]\n", encoding=DEFAULT_ENCODING)
    CICD_RUN_LOG.write_text(f"# CI/CD Full Run Log — {ts}\n\n", encoding=DEFAULT_ENCODING)
    init_meta = {
        "status": "running",
        "started_at": ts,
        "active_failures": [],
        "errors_log": str(CICD_ERRORS_LOG),
        "errors_json": str(CICD_ERRORS_JSON),
    }
    CICD_SUMMARY_JSON.write_text(json.dumps(init_meta, indent=2), encoding=DEFAULT_ENCODING)


def extract_stack_or_error(res: JobResult) -> str:
    """Combines and returns non-empty stdout/stderr from a failed job."""
    out = (res.out or "").strip()
    err = (res.err or "").strip()
    if out and err:
        return f"{out}\n\n{err}"

    return err or out or "No output captured."


def extract_failing_files(text: str) -> list[str]:
    """Extracts suspected source file paths and locations from error text."""
    pattern = re.compile(r'(?:[a-zA-Z0-9_\-./\\]+\.(?:go|py|ts|tsx|js|json|sh|ps1|yml|yaml)(?::\d+(?::\d+)?)?)')
    matches = pattern.findall(text)
    seen: set[str] = set()
    files: list[str] = []
    for m in matches:
        if m not in seen and not m.startswith("http"):
            seen.add(m)
            files.append(m)

    return files[:10]


def append_run_log(res: JobResult) -> None:
    """Appends gate execution outcome to full chronological run.log."""
    ts = time.strftime("%Y-%m-%d %H:%M:%S")
    status = "PASS" if res.is_success else "FAIL"
    cmd_str = " ".join(res.cmd) if isinstance(res.cmd, list) else str(res.cmd)
    entry = f"[{ts}] [{status}] {res.name} (code={res.code}, {res.elapsed}s) | Cmd: {cmd_str}\n"
    if not res.is_success:
        err = extract_stack_or_error(res)
        entry += f"  Error: {strip_ansi(err)}\n"
    try:
        with open(CICD_RUN_LOG, "a", encoding=DEFAULT_ENCODING) as fh:
            fh.write(entry)
    except OSError:
        pass


def update_cicd_summary(state: dict[str, Any], is_finished: bool = False) -> None:
    """Updates summary.json metadata with real-time status and counts."""
    total = state["total"]
    results = state.get("results", [])
    passed = sum(1 for r in results if r.is_success)
    failed = sum(1 for r in results if not r.is_success)
    remaining = max(0, total - passed - failed)
    status = "completed" if is_finished and failed == 0 else ("failed" if is_finished else "running")
    meta = {
        "status": status,
        "total_gates": total,
        "passed_gates": passed,
        "failed_gates": failed,
        "remaining_gates": remaining,
        "errors_count": len(state.get("errors_list", [])),
        "errors_log": str(CICD_ERRORS_LOG),
        "errors_json": str(CICD_ERRORS_JSON),
        "run_log": str(CICD_RUN_LOG),
        "summary_json": str(CICD_SUMMARY_JSON),
        "active_failures": [e["name"] for e in state.get("errors_list", [])],
    }
    try:
        CICD_SUMMARY_JSON.write_text(json.dumps(meta, indent=2), encoding=DEFAULT_ENCODING)
    except OSError:
        pass


def stream_failure_to_disk(res: JobResult, state: dict[str, Any]) -> None:
    """Appends failure information to real-time logs in .lovable/temp/cicd/."""
    ts = time.strftime("%Y-%m-%d %H:%M:%S")
    err_text = extract_stack_or_error(res)
    cmd_str = " ".join(res.cmd) if isinstance(res.cmd, list) else str(res.cmd)
    suspect_files = extract_failing_files(err_text)
    files_line = f"- **Suspect Files**: {', '.join(suspect_files)}\n" if suspect_files else ""
    entry = (
        f"### [{ts}] FAIL: {res.name}\n"
        f"- **Command**: `{cmd_str}`\n"
        f"- **Exit Code**: `{res.code}` ({res.elapsed}s)\n"
        f"{files_line}"
        f"```text\n{strip_ansi(err_text)}\n```\n\n"
    )
    with open(CICD_ERRORS_LOG, "a", encoding=DEFAULT_ENCODING) as fh:
        fh.write(entry)

    errors_list = state.setdefault("errors_list", [])
    errors_list.append({
        "name": res.name,
        "cmd": res.cmd,
        "code": res.code,
        "elapsed": res.elapsed,
        "suspect_files": suspect_files,
        "error": strip_ansi(err_text),
    })
    try:
        CICD_ERRORS_JSON.write_text(json.dumps(errors_list, indent=2), encoding=DEFAULT_ENCODING)
    except OSError:
        pass
    update_cicd_summary(state, is_finished=False)


def print_immediate_failure_report(res: JobResult, idx: int, total: int) -> None:
    """Immediately prints full failure stack trace to terminal without waiting for suite completion."""
    err_text = extract_stack_or_error(res)
    cmd_str = " ".join(res.cmd) if isinstance(res.cmd, list) else str(res.cmd)
    suspect_files = extract_failing_files(err_text)
    files_str = "\n".join(f"    • {f}" for f in suspect_files) if suspect_files else "    (None detected in output)"
    banner = (
        f"\n\033[1;91m================================================================\n"
        f"🚨 [IMMEDIATE FAILURE DETECTED] [{idx}/{total}] {res.name}\n"
        f"================================================================\033[0m\n"
        f"  Command       : {cmd_str}\n"
        f"  Exit Code     : {res.code} ({res.elapsed}s)\n"
        f"  Failing Files :\n{files_str}\n"
        f"  Stream Log    : {CICD_ERRORS_LOG}\n"
        f"  Stream JSON   : {CICD_ERRORS_JSON}\n\n"
        f"\033[1mStack Trace / Failure Output:\033[0m\n"
        f"----------------------------------------------------------------\n"
        f"{err_text}\n"
        f"\033[1;91m================================================================\033[0m\n"
    )
    print(banner, flush=True)


def print_runner_banner(concurrency_label: str, total_jobs: int) -> None:
    """Prints informational execution banner for verbose mode."""
    print("================================================================")
    print("           PARALLEL LOCAL CI/CD QUALITY GATE RUNNER             ")
    print("================================================================")
    print(f"🚀 Execution Mode          : {concurrency_label}")
    print(f"📋 Total Enqueued Gates    : {total_jobs}")
    print("🔍 Display Mode            : SHOW ALL INFORMATION (--all-paths)")
    print("----------------------------------------------------------------\n")


def execute_job_batch(
    batch: dict[str, Any],
    workers: int,
    is_sync: bool,
    args: argparse.Namespace,
    state: dict[str, Any],
) -> None:
    """Executes all jobs within a single batch with worker pool."""
    job_items = list(batch["jobs"].items())
    batch_limit = batch.get("max_workers")
    batch_workers = 1 if is_sync else min(workers, batch_limit or workers, len(job_items))

    with ThreadPoolExecutor(max_workers=batch_workers) as executor:
        future_to_name = {
            executor.submit(
                run_job,
                name,
                cmd.get("cmd") if isinstance(cmd, dict) else cmd,
                args.timeout,
                {**os.environ, **cmd.get("env")} if isinstance(cmd, dict) and "env" in cmd else None,
                cmd.get("cwd") if isinstance(cmd, dict) else None,
            ): name
            for name, cmd in job_items
        }
        for future in as_completed(future_to_name):
            state["counter"] += 1
            idx = state["counter"]
            total = state["total"]
            try:
                res = future.result()
                state["results"].append(res)
                append_run_log(res)
                if not res.is_success:
                    stream_failure_to_disk(res, state)
                    if not state["is_json"]:
                        print_immediate_failure_report(res, idx, total)
                elif args.show_all and not state["is_json"]:
                    print(f"  [{idx:2d}/{total}] \033[1;92m✓ PASS\033[0m [{res.name}] ({res.elapsed}s)", flush=True)
            except Exception as ex:
                name = future_to_name[future]
                failed_res = JobResult(
                    name=name,
                    cmd=batch["jobs"].get(name, []),
                    code=1,
                    out="",
                    err=str(ex),
                    elapsed=0.0,
                )
                state["results"].append(failed_res)
                append_run_log(failed_res)
                stream_failure_to_disk(failed_res, state)
                if not state["is_json"]:
                    print_immediate_failure_report(failed_res, idx, total)


def execute_runner(args: argparse.Namespace, active_batches: list[dict[str, Any]]) -> int:
    """Orchestrates test batch execution and report output generation."""
    total_jobs = sum(len(b["jobs"]) for b in active_batches)
    if total_jobs == 0:
        print(f"[WARN] No quality gates matched filter: {args.filter!r}")
        return 0

    init_cicd_temp_stream()
    is_sync = args.sync_mode
    workers = 1 if is_sync else max(1, args.workers)
    io_workers = 1 if is_sync else max(1, args.io_workers)
    is_json = bool(args.json_mode)

    if not is_json:
        print(f"📡 Real-Time Error Stream : {CICD_ERRORS_LOG}")
        print(f"📊 Real-Time JSON Status  : {CICD_SUMMARY_JSON}")
        if args.show_all:
            label = "Sequential (1 worker)" if is_sync else f"Parallel ({workers} workers, {io_workers} IO workers)"
            print_runner_banner(label, total_jobs)

    state: dict[str, Any] = {"counter": 0, "total": total_jobs, "results": [], "is_json": is_json}
    start_time = time.monotonic()

    for batch in active_batches:
        execute_job_batch(batch, workers, is_sync, args, state)

    total_elapsed = round(time.monotonic() - start_time, 2)
    results: list[JobResult] = state["results"]
    passed_count = sum(1 for r in results if r.is_success)
    failed_count = sum(1 for r in results if not r.is_success and not r.is_timeout)
    timeout_count = sum(1 for r in results if r.is_timeout)
    has_failures = bool(failed_count > 0 or timeout_count > 0)
    counts = (total_jobs, passed_count, failed_count, timeout_count)
    update_cicd_summary(state, is_finished=True)

    if is_json:
        return handle_json_output(args, results, counts, total_elapsed, has_failures)

    return handle_text_output(args, results, counts, total_elapsed)


def load_cicd_timings(path: Path) -> dict[str, float]:
    """Loads historical CI/CD job timings from JSON file if available."""
    if not path.exists():
        return {}

    try:
        data = json.loads(path.read_text(encoding=DEFAULT_ENCODING))
        if isinstance(data, dict):
            return {k: float(v) for k, v in data.items()}
    except Exception:
        pass

    return {}


def save_cicd_timings(path: Path, timings: dict[str, float]) -> None:
    """Persists recorded CI/CD job timings to the timing file."""
    path.parent.mkdir(parents=True, exist_ok=True)
    try:
        path.write_text(json.dumps(timings, indent=2), encoding=DEFAULT_ENCODING)
    except Exception:
        pass


def calculate_total_eta(active_batches: list[dict[str, Any]], timings: dict[str, float]) -> int:
    """Computes total estimated execution time across all active batches."""
    total_sec = 0.0
    for batch in active_batches:
        batch_max = DEFAULT_JOB_ESTIMATE_SEC
        for job_name in batch.get("jobs", {}):
            job_est = timings.get(job_name, DEFAULT_JOB_ESTIMATE_SEC)
            batch_max = max(batch_max, job_est)
        total_sec += batch_max

    return int(total_sec)


def run_eta_worker(
    interval_sec: int,
    total_est_sec: int,
    start_time: float,
    stop_event: threading.Event,
) -> None:
    """Background worker reporting remaining ETA every interval."""
    while not stop_event.is_set():
        if stop_event.wait(timeout=interval_sec):
            break
        elapsed = time.time() - start_time
        remaining = max(0, int(total_est_sec - elapsed))
        print(f"\n[ETA] Estimated remaining time: {remaining} seconds\n", flush=True)


def start_eta_reporter(
    interval_sec: int,
    total_est_sec: int,
    stop_event: threading.Event,
) -> threading.Thread | None:
    """Spawns background ETA daemon thread if interval > 0."""
    if interval_sec <= 0:
        return None

    worker = threading.Thread(
        target=run_eta_worker,
        args=(interval_sec, total_est_sec, time.time(), stop_event),
        daemon=True,
    )
    worker.start()

    return worker


def main() -> None:
    """Primary entry point for local CI/CD quality gate runner."""
    args = parse_args()
    active_batches = filter_job_batches(JOB_BATCHES, args.filter)
    timings = load_cicd_timings(TIMING_FILE_PATH)
    total_est = calculate_total_eta(active_batches, timings)

    stop_event = threading.Event()
    start_eta_reporter(args.eta_interval, total_est, stop_event)

    exit_code = 0
    try:
        exit_code = execute_runner(args, active_batches)
    finally:
        stop_event.set()
        timings.update(GLOBAL_TIMINGS)
        save_cicd_timings(TIMING_FILE_PATH, timings)

    sys.exit(exit_code)


if __name__ == "__main__":
    main()
