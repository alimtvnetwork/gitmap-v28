#!/usr/bin/env python3
"""Cross-platform Go test and linter preflight runner with concurrent worker pool and flexible reporting.

Executes local test suite and golangci-lint before submitting commits or PRs.

All Worker Pool primitives, Enums, and Constants are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/28-go-preflight-ci.py                 # Run both tests and lint in parallel (quiet on success)
  python 03-ai-scripts/28-go-preflight-ci.py --all-paths    # Show detailed ticker, summary table, and full logs
  python 03-ai-scripts/28-go-preflight-ci.py --sync         # Run sequentially (1 worker)
  python 03-ai-scripts/28-go-preflight-ci.py test           # Run tests only
  python 03-ai-scripts/28-go-preflight-ci.py lint           # Run linters only
  python 03-ai-scripts/28-go-preflight-ci.py -o report.txt  # Output execution report to file
  python 03-ai-scripts/28-go-preflight-ci.py --json         # Output results as machine-readable JSON
"""
from __future__ import annotations

import argparse
from dataclasses import dataclass
from importlib import import_module
import os
from pathlib import Path
import shutil
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

normalize_rel_path = engine.normalize_rel_path
ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR
WorkerResult = engine.WorkerResult
add_worker_cli_arguments = engine.add_worker_cli_arguments
run_worker_pool = engine.run_worker_pool


@dataclass
class PreflightTask:
    name: str
    task_type: str  # "test" or "lint"
    mod_dir: Path


def find_go_modules(repo_root: Path) -> list[Path]:
    """Finds all directories containing a go.mod file, excluding node_modules and .git."""
    mod_files = list(repo_root.rglob("go.mod"))
    modules: list[Path] = []
    for mf in mod_files:
        if "node_modules" not in str(mf) and ".git" not in str(mf):
            modules.append(mf.parent)

    return sorted(modules)


def execute_preflight_task(task: PreflightTask) -> WorkerResult:
    """Worker task executing go test or golangci-lint on a specific Go module."""
    start = time.perf_counter()
    mod_rel = normalize_rel_path(task.mod_dir)

    if task.task_type == "test":
        go_exe = shutil.which("go")
        if not go_exe:
            return WorkerResult(
                name=task.name,
                is_success=False,
                error="Go toolchain missing from PATH",
                elapsed_sec=round(time.perf_counter() - start, 3),
            )
        res = subprocess.run(
            [go_exe, "test", "./...", "-count=1"],
            cwd=str(task.mod_dir),
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace",
        )
        elapsed = round(time.perf_counter() - start, 3)
        return WorkerResult(
            name=task.name,
            is_success=(res.returncode == 0),
            output=res.stdout,
            error=res.stderr,
            elapsed_sec=elapsed,
        )

    # Lint task
    lint_exe = shutil.which("golangci-lint")
    if not lint_exe:
        return WorkerResult(
            name=task.name,
            is_success=True,
            output="golangci-lint not installed — skipped check",
            elapsed_sec=round(time.perf_counter() - start, 3),
        )

    res = subprocess.run(
        [lint_exe, "run", "./..."],
        cwd=str(task.mod_dir),
        capture_output=True,
        text=True,
        encoding=DEFAULT_ENCODING,
        errors="replace",
    )
    elapsed = round(time.perf_counter() - start, 3)
    return WorkerResult(
        name=task.name,
        is_success=(res.returncode == 0),
        output=res.stdout,
        error=res.stderr,
        elapsed_sec=elapsed,
    )


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/28-go-preflight-ci.py",
        description="Go Preflight CI Runner with parallel worker pool and flexible reporting.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: run all Go tests and linters in parallel; quiet on success (tick), detailed on error:
  python 03-ai-scripts/28-go-preflight-ci.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/28-go-preflight-ci.py --all-paths
  python 03-ai-scripts/28-go-preflight-ci.py --all

  # 3. Run tests only:
  python 03-ai-scripts/28-go-preflight-ci.py test

  # 4. Run linters only:
  python 03-ai-scripts/28-go-preflight-ci.py lint

  # 5. Run sequentially (1 worker):
  python 03-ai-scripts/28-go-preflight-ci.py --sync

  # 6. Save results to a file or JSON:
  python 03-ai-scripts/28-go-preflight-ci.py -o tmp/go-preflight.txt
  python 03-ai-scripts/28-go-preflight-ci.py --json
        """,
    )
    parser.add_argument("phase", nargs="?", choices=["all", "test", "lint"], default="all", help="Phase to execute")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent
    modules = find_go_modules(repo_root)

    tasks: list[PreflightTask] = []
    for mod in modules:
        mod_rel = normalize_rel_path(mod)
        if args.phase in ("all", "test"):
            tasks.append(PreflightTask(name=f"Go Test: {mod_rel}", task_type="test", mod_dir=mod))
        if args.phase in ("all", "lint"):
            tasks.append(PreflightTask(name=f"Go Lint: {mod_rel}", task_type="lint", mod_dir=mod))

    if args.filter:
        filt = args.filter.lower()
        tasks = [t for t in tasks if filt in t.name.lower()]

    exit_code = run_worker_pool(
        items=tasks,
        worker_fn=execute_preflight_task,
        max_workers=args.workers,
        is_sync=args.is_sync,
        show_all=args.show_all,
        output_file=args.output_file,
        as_json=args.as_json,
        title="GO PREFLIGHT CI RUNNER",
        item_noun="Go preflight task(s)",
    )
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
