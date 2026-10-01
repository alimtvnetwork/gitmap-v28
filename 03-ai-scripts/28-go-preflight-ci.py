#!/usr/bin/env python3
"""
28-go-preflight-ci.py — Cross-platform Go test and linter preflight runner.

Executes local Go test suites and golangci-lint checks concurrently using
worker groups scaled to CPU cores, or sequentially when requested.

Features:
  - Parallel Execution by default: worker group scaled to CPU cores (configurable via
    variable, CI_MAX_WORKERS env var, or --workers flag).
  - Quiet by default on success: prints only a single clean tick line:
      ✔ All passed. (N Go preflight gates verified in X.XXs)
  - Failure Isolation: prints full stack traces and failure logs only when checks fail.
  - Comprehensive CLI options:
      --all-paths / --all-passed / --all / -a : Detailed banner, live progress, summary table, and full logs.
      --failed / -f                           : Show logs only for failed checks (default behavior).
      --sync / --sequential / -s              : Run sequentially in 1 worker.
      --workers / -w / --concurrency          : Custom concurrency count.
      --output / -o / --file                  : Write execution report to file.
      --json                                  : Machine-readable JSON output for AI agents.
      --filter / -k                           : Filter checks by module name or phase substring.
      --phase                                 : Target check phase (all, test, lint).

Usage:
  python 03-ai-scripts/28-go-preflight-ci.py                     # run all gates concurrently (quiet)
  python 03-ai-scripts/28-go-preflight-ci.py --all-paths         # run all gates with detailed output
  python 03-ai-scripts/28-go-preflight-ci.py test                # run tests only
  python 03-ai-scripts/28-go-preflight-ci.py lint                # run linters only
  python 03-ai-scripts/28-go-preflight-ci.py --sync              # run sequentially
  python 03-ai-scripts/28-go-preflight-ci.py --workers 4         # run with 4 workers
  python 03-ai-scripts/28-go-preflight-ci.py --json              # output JSON summary

Exit codes:
  0 = all checks passed, 1 = test or lint failures.
"""

from __future__ import annotations

import argparse
from dataclasses import asdict, dataclass
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
ExitCodeType = engine.ExitCodeType
LINE_SEPARATOR = engine.LINE_SEPARATOR
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS
WorkItemResult = engine.WorkItemResult
WorkGroupSummary = engine.WorkGroupSummary
run_worker_pool = engine.run_worker_pool


@dataclass
class GoPreflightItem:
    """Represents an individual Go test or lint gate for a module."""
    name: str
    phase: str
    mod_dir: Path
    rel_mod: str
    command: list[str]


def find_go_modules(repo_root: Path) -> list[Path]:
    """Finds all directories containing a go.mod file, excluding node_modules and .git."""
    mod_files = list(repo_root.rglob("go.mod"))
    modules = []
    for mf in mod_files:
        mf_str = str(mf)
        is_ignored = ("node_modules" in mf_str or ".git" in mf_str)
        if is_ignored:
            continue
        modules.append(mf.parent)

    return sorted(modules)


def execute_go_check(item: GoPreflightItem) -> WorkItemResult:
    """Executes a single Go preflight check item and captures standard streams and timing."""
    start_time = time.perf_counter()
    exe_name = item.command[0]
    exe_path = shutil.which(exe_name)

    if not exe_path:
        is_lint = (item.phase == "lint")
        if is_lint:
            note = (
                f"golangci-lint not installed — skipping Go linter check for {item.rel_mod}.\n"
                "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest\n"
            )
            return WorkItemResult(
                name=item.name,
                is_success=True,
                output=note,
                duration_sec=0.0,
                return_code=0,
                data={"module": item.rel_mod, "phase": item.phase, "skipped": True}
            )

        err_msg = f"✗ Go toolchain '{exe_name}' missing from PATH"
        return WorkItemResult(
            name=item.name,
            is_success=False,
            output=err_msg,
            duration_sec=0.0,
            return_code=-1,
            data={"module": item.rel_mod, "phase": item.phase, "skipped": False}
        )

    try:
        res = subprocess.run(
            [exe_path] + item.command[1:],
            cwd=str(item.mod_dir),
            capture_output=True,
            text=True
        )
        duration_sec = time.perf_counter() - start_time
        combined_output = (res.stdout or "") + (res.stderr or "")
        is_success = (res.returncode == 0)

        return WorkItemResult(
            name=item.name,
            is_success=is_success,
            output=combined_output.strip(),
            duration_sec=duration_sec,
            return_code=res.returncode,
            data={"module": item.rel_mod, "phase": item.phase, "skipped": False}
        )
    except Exception as exc:
        duration_sec = time.perf_counter() - start_time
        return WorkItemResult(
            name=item.name,
            is_success=False,
            output=f"Execution error: {exc}",
            duration_sec=duration_sec,
            return_code=-1,
            data={"module": item.rel_mod, "phase": item.phase, "skipped": False}
        )


def build_execution_report(summary: WorkGroupSummary, show_all: bool) -> str:
    """Formats the human-readable text execution report."""
    lines = []

    if show_all:
        lines.append("============================================================")
        lines.append("             GO PREFLIGHT VERIFICATION REPORT               ")
        lines.append("============================================================")
        for r in summary.results:
            status_icon = "✅" if r.is_success else "❌"
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"{status_icon} [{status_word}] {r.name:<45} ({r.duration_sec:.2f}s)")
        lines.append("------------------------------------------------------------")
        lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
        lines.append(f"Gates Passed   : {summary.passed_count}/{summary.total_items}")
        lines.append(f"Gates Failed   : {summary.failed_count}/{summary.total_items}")
        lines.append("------------------------------------------------------------")

        lines.append("\n=================== ALL GATE LOGS ===================")
        for r in summary.results:
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"\n--- [{status_word}] {r.name} ({r.duration_sec:.2f}s) ---")
            lines.append(r.output if r.output else "(no output)")
            lines.append("------------------------------------------------------------")
    else:
        if summary.has_failures:
            lines.append("\n================ FAILED GATE LOGS ================")
            for r in summary.results:
                if not r.is_success:
                    lines.append(f"\n❌ FAILED: {r.name} ({r.duration_sec:.2f}s)")
                    lines.append(f"--- {r.name} ---")
                    lines.append(r.output if r.output else "(no output)")
                    lines.append(f"--- END {r.name} ---")
            lines.append("------------------------------------------------------------")
            lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
            lines.append(f"Gates Passed   : {summary.passed_count}/{summary.total_items}")
            lines.append(f"Gates Failed   : {summary.failed_count}/{summary.total_items}")
            lines.append("------------------------------------------------------------")
            lines.append(f"\n❌ Go Preflight failed: {summary.failed_count} gate(s) reported violations.")

    return LINE_SEPARATOR.join(lines)


def run_go_preflight_suite(
    repo_root: Path,
    phase: str = "all",
    max_workers: int | None = None,
    show_all: bool = False,
    is_sync: bool = False,
    output_file: str | None = None,
    as_json: bool = False,
    filter_pattern: str | None = None,
) -> int:
    """Dispatches Go preflight checks concurrently via worker group or sequentially."""
    modules = find_go_modules(repo_root)
    if not modules:
        if as_json:
            print(json.dumps({"total_items": 0, "passed_count": 0, "failed_count": 0, "results": []}))
        else:
            print("✔ All passed. (0 Go modules discovered in repository)")
        return ExitCodeType.SUCCESS.value

    items: list[GoPreflightItem] = []
    for mod in modules:
        rel_mod = normalize_rel_path(mod.relative_to(repo_root))
        if phase in ("all", "test"):
            items.append(GoPreflightItem(
                name=f"{rel_mod} [test]",
                phase="test",
                mod_dir=mod,
                rel_mod=rel_mod,
                command=["go", "test", "./...", "-count=1"]
            ))

        if phase in ("all", "lint"):
            items.append(GoPreflightItem(
                name=f"{rel_mod} [lint]",
                phase="lint",
                mod_dir=mod,
                rel_mod=rel_mod,
                command=["golangci-lint", "run", "./..."]
            ))

    if filter_pattern:
        filter_lower = filter_pattern.lower()
        items = [i for i in items if filter_lower in i.name.lower()]

    if not items:
        if as_json:
            print(json.dumps({"total_items": 0, "passed_count": 0, "failed_count": 0, "results": []}))
        else:
            print(f"⚠️ No Go preflight checks matched filter '{filter_pattern}'")
        return ExitCodeType.SUCCESS.value

    worker_count = max_workers or min(len(items), DEFAULT_CONCURRENCY_WORKERS)
    if is_sync:
        worker_count = 1

    if not as_json:
        if show_all:
            concurrency_label = "Sequential (1 worker)" if is_sync else f"{worker_count} parallel workers"
            print("================================================================")
            print("             PARALLEL GO PREFLIGHT CI RUNNER                   ")
            print("================================================================")
            print(f"🚀 Execution Mode          : {concurrency_label}")
            print(f"📋 Total Enqueued Gates    : {len(items)}")
            print("🔍 Display Mode            : SHOW ALL INFORMATION")
            print("----------------------------------------------------------------\n")

    def ticker_callback(res: WorkItemResult, current: int, total: int):
        if not as_json:
            if show_all:
                status_icon = "✅" if res.is_success else "❌"
                status_label = "PASS" if res.is_success else "FAIL"
                print(f"[{current:2d}/{total:2d}] {status_icon} [{status_label}] {res.name} ({res.duration_sec:.2f}s)")

    summary = run_worker_pool(
        items=items,
        worker_fn=execute_go_check,
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
            print(f"✔ All passed. ({summary.passed_count} Go preflight gates verified in {summary.wall_duration_sec:.2f}s)")

    if output_file:
        try:
            full_report = build_execution_report(summary, show_all=True)
            Path(output_file).write_text(full_report, encoding="utf-8")
        except Exception as e:
            print(f"⚠️ Failed to write report to '{output_file}': {e}", file=sys.stderr)

    return summary.exit_code


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Go Preflight CI Runner — executes Go test and linter suites concurrently.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python 03-ai-scripts/28-go-preflight-ci.py                      # Run all checks concurrently (quiet on success)
  python 03-ai-scripts/28-go-preflight-ci.py --all-paths          # Run all checks with full output and progress
  python 03-ai-scripts/28-go-preflight-ci.py test                 # Run tests only
  python 03-ai-scripts/28-go-preflight-ci.py lint                 # Run linters only
  python 03-ai-scripts/28-go-preflight-ci.py --sync               # Run sequentially in 1 worker
  python 03-ai-scripts/28-go-preflight-ci.py -w 4                 # Concurrency capped at 4 workers
  python 03-ai-scripts/28-go-preflight-ci.py --json               # Machine-readable JSON summary
  python 03-ai-scripts/28-go-preflight-ci.py -o preflight.log     # Save full report to file
        """
    )
    parser.add_argument(
        "phase_pos",
        nargs="?",
        choices=["all", "test", "lint"],
        default=None,
        help="Positional check phase to execute (all, test, lint). Default: all."
    )
    parser.add_argument(
        "--phase",
        choices=["all", "test", "lint"],
        default=None,
        help="Explicit flag for check phase (all, test, lint)."
    )
    parser.add_argument(
        "--all-paths", "--all-passed", "--all", "-a",
        dest="all_paths",
        action="store_true",
        help="Display detailed banners, execution tables, and logs for all gates."
    )
    parser.add_argument(
        "--failed", "-f",
        dest="failed_only",
        action="store_true",
        default=True,
        help="Show logs only for failed checks (default behavior)."
    )
    parser.add_argument(
        "--sync", "--sequential", "-s",
        dest="is_sync",
        action="store_true",
        help="Execute checks sequentially without threading."
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
    parser.add_argument(
        "--filter", "-k",
        dest="filter_pattern",
        type=str,
        default=None,
        help="Filter checks by module path or phase substring."
    )

    args = parser.parse_args()

    selected_phase = args.phase or args.phase_pos or "all"
    repo_root = Path(__file__).resolve().parent.parent

    return run_go_preflight_suite(
        repo_root=repo_root,
        phase=selected_phase,
        max_workers=args.max_workers,
        show_all=args.all_paths,
        is_sync=args.is_sync,
        output_file=args.output_file,
        as_json=args.as_json,
        filter_pattern=args.filter_pattern
    )


if __name__ == "__main__":
    sys.exit(main())
