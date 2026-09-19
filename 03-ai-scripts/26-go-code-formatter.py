#!/usr/bin/env python3
"""
26-go-code-formatter.py — Cross-platform Go code formatter using gofmt.

Formats Go source files concurrently using worker groups scaled to CPU cores,
or sequentially when requested.

Features:
  - Parallel Execution by default: worker group scaled to CPU cores (configurable via
    variable, CI_MAX_WORKERS env var, or --workers flag).
  - Quiet by default on success: prints only a single clean tick line:
      ✔ All passed. (N Go files verified/formatted in X.XXs)
  - Failure Isolation: prints full stack traces and failure logs only when formatting fails.
  - Comprehensive CLI options:
      --all-paths / --all-passed / --all / -a : Detailed banner, live progress, summary table, and full logs.
      --failed / -f                           : Show logs only for failed files (default behavior).
      --sync / --sequential / -s              : Run sequentially in 1 worker.
      --workers / -w / --concurrency          : Custom concurrency count.
      --output / -o / --file                  : Write execution report to file.
      --json                                  : Machine-readable JSON output for AI agents.
      --staged                                : Format only staged git files.

Usage:
  python 03-ai-scripts/26-go-code-formatter.py                     # format all .go files concurrently (quiet)
  python 03-ai-scripts/26-go-code-formatter.py --all-paths         # format all files with detailed output
  python 03-ai-scripts/26-go-code-formatter.py --staged            # format only staged .go files
  python 03-ai-scripts/26-go-code-formatter.py path/to/file.go     # format specific file(s)
  python 03-ai-scripts/26-go-code-formatter.py --sync              # run sequentially
  python 03-ai-scripts/26-go-code-formatter.py --workers 4         # run with 4 workers
  python 03-ai-scripts/26-go-code-formatter.py --json              # output JSON summary

Exit codes:
  0 — clean or formatted successfully
  1 — formatting error
  2 — tool missing
"""

from __future__ import annotations

import argparse
from dataclasses import asdict
from importlib import import_module
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

normalize_rel_path = engine.normalize_rel_path
stream_directory_files = engine.stream_directory_files
ExitCodeType = engine.ExitCodeType
LINE_SEPARATOR = engine.LINE_SEPARATOR
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS
WorkItemResult = engine.WorkItemResult
WorkGroupSummary = engine.WorkGroupSummary
run_worker_pool = engine.run_worker_pool


def get_staged_go_files(repo_root: Path) -> list[Path]:
    """Retrieves list of staged .go files from git index."""
    git_exe = shutil.which("git")
    if not git_exe:
        return []

    res = subprocess.run(
        [git_exe, "diff", "--cached", "--name-only", "--diff-filter=ACM"],
        cwd=str(repo_root),
        capture_output=True,
        text=True,
    )
    if res.returncode != 0:
        return []

    staged = []
    for line in res.stdout.splitlines():
        rel = line.strip()
        if rel.endswith(".go"):
            file_path = repo_root / rel
            if file_path.is_file():
                staged.append(file_path)

    return staged


def make_format_worker(gofmt_exe: str, repo_root: Path):
    """Creates a worker function bound to the gofmt binary path."""
    def worker_fn(file_path: Path) -> WorkItemResult:
        start_time = time.perf_counter()
        rel_path = normalize_rel_path(file_path.relative_to(repo_root))
        try:
            res = subprocess.run([gofmt_exe, "-w", str(file_path)], capture_output=True, text=True)
            duration_sec = time.perf_counter() - start_time
            is_success = (res.returncode == 0)
            err_output = res.stderr.strip() if res.stderr else ""
            return WorkItemResult(
                name=rel_path,
                is_success=is_success,
                output=err_output,
                duration_sec=duration_sec,
                return_code=res.returncode,
                data={"file": str(file_path)}
            )
        except Exception as exc:
            duration_sec = time.perf_counter() - start_time
            return WorkItemResult(
                name=rel_path,
                is_success=False,
                output=f"Execution error: {exc}",
                duration_sec=duration_sec,
                return_code=-1,
                data={"file": str(file_path)}
            )
    return worker_fn


def build_execution_report(summary: WorkGroupSummary, show_all: bool) -> str:
    """Formats the human-readable text execution report."""
    lines = []

    if show_all:
        lines.append("============================================================")
        lines.append("             GO CODE FORMATTER REPORT                       ")
        lines.append("============================================================")
        for r in summary.results:
            status_icon = "✅" if r.is_success else "❌"
            status_word = "FORMATTED" if r.is_success else "FAILED"
            lines.append(f"{status_icon} [{status_word}] {r.name:<50} ({r.duration_sec:.2f}s)")
        lines.append("------------------------------------------------------------")
        lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
        lines.append(f"Files Passed   : {summary.passed_count}/{summary.total_items}")
        lines.append(f"Files Failed   : {summary.failed_count}/{summary.total_items}")
        lines.append("------------------------------------------------------------")

        lines.append("\n=================== ALL FORMATTING LOGS ===================")
        for r in summary.results:
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"\n--- [{status_word}] {r.name} ({r.duration_sec:.2f}s) ---")
            lines.append(r.output if r.output else "(no warnings)")
            lines.append("------------------------------------------------------------")
    else:
        if summary.has_failures:
            lines.append("\n================ FAILED FORMATTING LOGS ================")
            for r in summary.results:
                if not r.is_success:
                    lines.append(f"\n❌ FAILED: {r.name} ({r.duration_sec:.2f}s)")
                    lines.append(f"--- {r.name} ---")
                    lines.append(r.output if r.output else "(no output)")
                    lines.append(f"--- END {r.name} ---")
            lines.append("------------------------------------------------------------")
            lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
            lines.append(f"Files Passed   : {summary.passed_count}/{summary.total_items}")
            lines.append(f"Files Failed   : {summary.failed_count}/{summary.total_items}")
            lines.append("------------------------------------------------------------")
            lines.append(f"\n❌ Go formatting failed: {summary.failed_count} file(s) reported errors.")

    return LINE_SEPARATOR.join(lines)


def run_go_formatter(
    repo_root: Path,
    target_files: list[Path],
    gofmt_exe: str,
    max_workers: int | None = None,
    show_all: bool = False,
    is_sync: bool = False,
    output_file: str | None = None,
    as_json: bool = False,
) -> int:
    """Formats Go files concurrently using worker group or sequentially."""
    if not target_files:
        if as_json:
            print(json.dumps({"total_items": 0, "passed_count": 0, "failed_count": 0, "results": []}))
        else:
            print("✔ All passed. (0 Go files to format)")
        return ExitCodeType.SUCCESS.value

    worker_count = max_workers or min(len(target_files), DEFAULT_CONCURRENCY_WORKERS)
    if is_sync:
        worker_count = 1

    if not as_json:
        if show_all:
            concurrency_label = "Sequential (1 worker)" if is_sync else f"{worker_count} parallel workers"
            print("================================================================")
            print("              PARALLEL GO CODE FORMATTER                        ")
            print("================================================================")
            print(f"🚀 Execution Mode          : {concurrency_label}")
            print(f"📋 Total Enqueued Files    : {len(target_files)}")
            print("🔍 Display Mode            : SHOW ALL INFORMATION")
            print("----------------------------------------------------------------\n")

    def ticker_callback(res: WorkItemResult, current: int, total: int):
        if not as_json:
            if show_all:
                status_icon = "✅" if res.is_success else "❌"
                status_label = "PASS" if res.is_success else "FAIL"
                print(f"[{current:2d}/{total:2d}] {status_icon} [{status_label}] {res.name} ({res.duration_sec:.2f}s)")

    worker_fn = make_format_worker(gofmt_exe, repo_root)

    summary = run_worker_pool(
        items=target_files,
        worker_fn=worker_fn,
        worker_count=worker_count,
        is_sync=is_sync,
        on_item_complete=ticker_callback
    )

    if as_json:
        payload = {
            "total_items": summary.total_items,
            "passed_count": summary.passed_count,
            "failed_count": summary.failed_count,
            "wall_duration_sec": round(summary.wall_duration_sec, 3),
            "has_failures": summary.has_failures,
            "exit_code": summary.exit_code,
            "results": [asdict(r) for r in summary.results],
        }
        json_str = json.dumps(payload, indent=2)
        print(json_str)
        if output_file:
            try:
                Path(output_file).write_text(json_str, encoding="utf-8")
            except Exception as e:
                print(f"⚠️ Failed to write JSON output to '{output_file}': {e}", file=sys.stderr)
        return summary.exit_code

    report = build_execution_report(summary, show_all=show_all)

    if show_all:
        print(report)
    else:
        if summary.has_failures:
            print(report)
        else:
            print(f"✔ All passed. ({summary.passed_count} Go files verified/formatted in {summary.wall_duration_sec:.2f}s)")

    if output_file:
        try:
            full_report = build_execution_report(summary, show_all=True)
            Path(output_file).write_text(full_report, encoding="utf-8")
        except Exception as e:
            print(f"⚠️ Failed to write report to '{output_file}': {e}", file=sys.stderr)

    return summary.exit_code


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Cross-platform Go code formatter using gofmt concurrently.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python 03-ai-scripts/26-go-code-formatter.py                      # Format all Go files concurrently (quiet on success)
  python 03-ai-scripts/26-go-code-formatter.py --all-paths          # Format all files with full details
  python 03-ai-scripts/26-go-code-formatter.py --staged             # Format only staged git files
  python 03-ai-scripts/26-go-code-formatter.py path/to/file.go      # Format specific file(s)
  python 03-ai-scripts/26-go-code-formatter.py --sync               # Run sequentially in 1 worker
  python 03-ai-scripts/26-go-code-formatter.py -w 4                 # Concurrency capped at 4 workers
  python 03-ai-scripts/26-go-code-formatter.py --json               # Machine-readable JSON summary
        """
    )
    parser.add_argument("paths", nargs="*", help="Specific files or directories to format")
    parser.add_argument("--staged", action="store_true", help="Format only staged git files")
    parser.add_argument(
        "--all-paths", "--all-passed", "--all", "-a",
        dest="all_paths",
        action="store_true",
        help="Display detailed banners, execution tables, and logs for all files."
    )
    parser.add_argument(
        "--failed", "-f",
        dest="failed_only",
        action="store_true",
        default=True,
        help="Show logs only for failed files (default behavior)."
    )
    parser.add_argument(
        "--sync", "--sequential", "-s",
        dest="is_sync",
        action="store_true",
        help="Execute formatting sequentially without threading."
    )
    parser.add_argument(
        "--workers", "-w", "--concurrency",
        dest="max_workers",
        type=int,
        default=None,
        help=f"Number of parallel worker threads (default: {DEFAULT_CONCURRENCY_WORKERS})."
    )
    parser.add_argument(
        "--output", "-o", "--file",
        dest="output_file",
        type=str,
        default=None,
        help="Path to write execution report file."
    )
    parser.add_argument(
        "--json",
        dest="as_json",
        action="store_true",
        help="Output structured JSON summary for automation."
    )

    args = parser.parse_args()

    gofmt_exe = shutil.which("gofmt")
    if not gofmt_exe:
        print("⚠ gofmt not found in PATH — install Go toolchain (https://go.dev/dl/)", file=sys.stderr)
        return int(ExitCodeType.TOOL_ERROR.value)

    repo_root = Path(__file__).resolve().parent.parent

    target_files: list[Path] = []
    if args.staged:
        target_files = get_staged_go_files(repo_root)
    elif args.paths:
        for p_str in args.paths:
            p = Path(p_str).resolve()
            if p.is_file() and p.suffix == ".go":
                target_files.append(p)
            elif p.is_dir():
                target_files.extend(list(p.rglob("*.go")))
    else:
        for f in stream_directory_files(repo_root, extensions=[".go"]):
            target_files.append(f)

    return run_go_formatter(
        repo_root=repo_root,
        target_files=target_files,
        gofmt_exe=gofmt_exe,
        max_workers=args.max_workers,
        show_all=args.all_paths,
        is_sync=args.is_sync,
        output_file=args.output_file,
        as_json=args.as_json,
    )


if __name__ == "__main__":
    sys.exit(main())
