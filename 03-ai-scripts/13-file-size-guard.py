#!/usr/bin/env python3
"""
Fast Repository File Size & Blob Guard
Scans tracked repository files to ensure no accidental massive binary files exceed thresholds.
Multi-folder capable, customizable extensions, nested ignore pruning (.git, .gitmap, node_modules),
and concurrent execution via the shared worker pool engine.

Performance & Clean Architecture:
1. Flattened Conditionals: Zero nested if-blocks using clean guard clauses.
2. Concurrent Worker Pool: Audits files in parallel across CPU cores, quiet on success, detailed on failure.
3. All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

from __future__ import annotations

import argparse
from importlib import import_module
import os
from pathlib import Path
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

DEFAULT_MAX_FILE_KB = engine.DEFAULT_MAX_FILE_KB
EXCLUDE_DIRS = engine.EXCLUDE_DIRS
is_allowed_large_file = engine.is_allowed_large_file
is_ignored_directory = engine.is_ignored_directory
normalize_rel_path = engine.normalize_rel_path
normalize_extensions = engine.normalize_extensions
ExitCodeType = engine.ExitCodeType
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS
WorkerResult = engine.WorkerResult
run_worker_pool = engine.run_worker_pool
add_worker_cli_arguments = engine.add_worker_cli_arguments


def check_file_size(
    file_path: Path,
    max_bytes: int,
) -> WorkerResult:
    """Audits a single file's size against threshold using flattened guard clauses."""
    start_time = time.perf_counter()
    norm_fp = normalize_rel_path(file_path)

    if is_allowed_large_file(norm_fp):
        return WorkerResult(
            name=norm_fp,
            is_success=True,
            output=f"Allowed large file waiver: {norm_fp}",
            elapsed_sec=round(time.perf_counter() - start_time, 4),
        )

    try:
        sz = file_path.stat().st_size
        elapsed = round(time.perf_counter() - start_time, 4)
        if sz > max_bytes:
            max_kb = max_bytes // 1024
            return WorkerResult(
                name=norm_fp,
                is_success=False,
                error=f"::error file={norm_fp}::{norm_fp} ({sz / 1024:.1f} KB > {max_kb} KB)",
                elapsed_sec=elapsed,
            )
        return WorkerResult(
            name=norm_fp,
            is_success=True,
            output=f"{norm_fp} ({sz / 1024:.1f} KB)",
            elapsed_sec=elapsed,
        )
    except Exception as exc:
        return WorkerResult(
            name=norm_fp,
            is_success=False,
            error=f"Exception checking {norm_fp}: {exc}",
            elapsed_sec=round(time.perf_counter() - start_time, 4),
        )


def audit_file_sizes(
    max_kb: int = DEFAULT_MAX_FILE_KB,
    target_dir: str = CURRENT_DIR,
    allowed_exts: set[str] | None = None,
    max_workers: int | None = None,
    is_sync: bool = False,
    show_all: bool = False,
    output_file: str | None = None,
    as_json: bool | str = False,
    filter_pattern: str | None = None,
) -> int:
    """Scans files and checks sizes against threshold across target directory using parallel worker pool."""
    target_path = Path(target_dir).resolve()
    max_bytes = max_kb * 1024

    files_to_check: list[Path] = []
    if target_path.is_file():
        files_to_check.append(target_path)
    else:
        for root, dirs, files in os.walk(str(target_path)):
            dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
            for f in files:
                fp = Path(root) / f
                if allowed_exts is not None and fp.suffix.lower() not in allowed_exts:
                    continue
                files_to_check.append(fp)

    if filter_pattern:
        filt = filter_pattern.lower().replace("\\", "/")
        files_to_check = [f for f in files_to_check if filt in normalize_rel_path(f).lower()]

    if not files_to_check:
        print(f"✔ All passed. (0 file(s) in 0.00s)")
        return ExitCodeType.SUCCESS.value

    def worker_adapter(p: Path) -> WorkerResult:
        return check_file_size(p, max_bytes=max_bytes)

    return run_worker_pool(
        items=files_to_check,
        worker_fn=worker_adapter,
        max_workers=max_workers,
        is_sync=is_sync,
        show_all=show_all,
        output_file=output_file,
        as_json=as_json,
        title=f"REPOSITORY FILE SIZE GUARD (Max {max_kb} KB)",
        item_noun="file(s)",
    )


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/13-file-size-guard.py",
        description="Audit repository file sizes across folders using parallel worker pool.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: run all file size audits in parallel; quiet on success (tick), detailed on error:
  python 03-ai-scripts/13-file-size-guard.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/13-file-size-guard.py --all-paths
  python 03-ai-scripts/13-file-size-guard.py --all

  # 3. Run sequentially (1 worker):
  python 03-ai-scripts/13-file-size-guard.py --sync

  # 4. Save report to a file:
  python 03-ai-scripts/13-file-size-guard.py -o tmp/file-size-report.txt

  # 5. Output results as machine-readable JSON:
  python 03-ai-scripts/13-file-size-guard.py --json
        """,
    )
    parser.add_argument("--max-kb", type=int, default=DEFAULT_MAX_FILE_KB, help=f"Maximum allowed file size in KB (default: {DEFAULT_MAX_FILE_KB})")
    parser.add_argument("--path", "-p", default=CURRENT_DIR, help="Root path or folder to audit")
    parser.add_argument("--ext", help="Optional comma-separated extension filter (e.g. .json,.zip)")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    exts = normalize_extensions(args.ext)
    exit_code = audit_file_sizes(
        max_kb=args.max_kb,
        target_dir=args.path,
        allowed_exts=exts,
        max_workers=args.workers,
        is_sync=args.is_sync,
        show_all=args.show_all,
        output_file=args.output_file,
        as_json=args.as_json,
        filter_pattern=args.filter,
    )
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
