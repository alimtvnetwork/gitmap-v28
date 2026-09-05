#!/usr/bin/env python3
"""
Fast UTF-8 & UNIX LF Encoding Normalizer
Recursively audits and standardizes all text files to UTF-8 without BOM and strict UNIX LF (\\n).
Multi-folder capable, customizable extensions, and concurrent execution via shared worker pool engine.

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

process_repository_files = engine.process_repository_files
write_file_lf = engine.write_file_lf
is_binary_file = engine.is_binary_file
normalize_extensions = engine.normalize_extensions
normalize_rel_path = engine.normalize_rel_path
ExitCodeType = engine.ExitCodeType
stream_directory_files = engine.stream_directory_files
WorkerResult = engine.WorkerResult
run_worker_pool = engine.run_worker_pool
add_worker_cli_arguments = engine.add_worker_cli_arguments
DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
UTF8_SIG_ENCODING = engine.UTF8_SIG_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
UTF8_BOM_BYTES = engine.UTF8_BOM_BYTES
CRLF_BYTES = engine.CRLF_BYTES
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS


def normalize_single_file(file_path: Path, is_fix_mode: bool = False) -> tuple[str, bool, str]:
    """Audits and converts CRLF/BOM in a file to clean UTF-8 LF using flattened guard clauses."""
    try:
        rel_path = file_path.resolve().relative_to(Path.cwd().resolve())
        norm_p = normalize_rel_path(rel_path)
    except Exception:
        norm_p = normalize_rel_path(file_path)

    if is_binary_file(file_path):
        return (norm_p, False, "")

    try:
        with open(file_path, "rb") as f:
            raw_bytes = f.read()

        has_bom = raw_bytes.startswith(UTF8_BOM_BYTES)
        has_crlf = CRLF_BYTES in raw_bytes
        has_issue = has_bom or has_crlf
        if not has_issue:
            return (norm_p, False, "")

        issue_desc = []
        if has_bom:
            issue_desc.append("UTF-8 BOM detected")
        if has_crlf:
            issue_desc.append("CRLF (\\r\\n) line endings detected")

        detail = ", ".join(issue_desc)

        if is_fix_mode:
            text = raw_bytes.decode(UTF8_SIG_ENCODING, errors="replace")
            write_file_lf(file_path, text, encoding=DEFAULT_ENCODING)
            return (norm_p, True, f"Normalized to UTF-8 LF: {detail}")

        return (norm_p, True, detail)
    except Exception as exc:
        return (norm_p, False, f"Read error: {exc}")


def execute_encoding_task(task_tuple: tuple[Path, bool]) -> WorkerResult:
    file_path, is_fix_mode = task_tuple
    start_time = time.perf_counter()
    norm_p, has_issue, detail = normalize_single_file(file_path, is_fix_mode=is_fix_mode)
    elapsed = round(time.perf_counter() - start_time, 4)

    if has_issue and not is_fix_mode:
        return WorkerResult(
            name=norm_p,
            is_success=False,
            error=f"::notice file={norm_p}::{norm_p} ({detail})",
            elapsed_sec=elapsed,
        )
    elif has_issue and is_fix_mode:
        return WorkerResult(
            name=norm_p,
            is_success=True,
            output=f"Fixed {norm_p} ({detail})",
            elapsed_sec=elapsed,
        )

    return WorkerResult(
        name=norm_p,
        is_success=True,
        output=f"UTF-8 LF clean: {norm_p}",
        elapsed_sec=elapsed,
    )


def run_encoding_normalizer(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = False,
    extensions: set[str] | tuple | None = None,
    max_workers: int | None = None,
    is_sync: bool = False,
    show_all: bool = False,
    output_file: str | None = None,
    as_json: bool | str = False,
    filter_pattern: str | None = None,
) -> int:
    """Runs repository encoding check and normalizer across target directory using parallel worker pool."""
    exts = normalize_extensions(extensions) or DEFAULT_TEXT_EXTENSIONS
    target_path = Path(target_dir).resolve()

    files: list[Path] = []
    if target_path.is_file():
        files.append(target_path)
    else:
        for f in stream_directory_files(root_dir=str(target_path), extensions=exts):
            files.append(f)

    if filter_pattern:
        filt = filter_pattern.lower().replace("\\", "/")
        files = [f for f in files if filt in normalize_rel_path(f).lower()]

    if not files:
        print(f"✔ All passed. (0 text file(s) in 0.00s)")
        return ExitCodeType.SUCCESS.value

    tasks = [(f, is_fix_mode) for f in files]

    return run_worker_pool(
        items=tasks,
        worker_fn=execute_encoding_task,
        max_workers=max_workers,
        is_sync=is_sync,
        show_all=show_all,
        output_file=output_file,
        as_json=as_json,
        title="UTF-8 & UNIX LF ENCODING NORMALIZER",
        item_noun="text file(s)",
    )


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/10-encoding-normalizer.py",
        description="Normalize files to UTF-8 UNIX LF across folders using parallel worker pool.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: check files in parallel; quiet on success (tick), detailed on error:
  python 03-ai-scripts/10-encoding-normalizer.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/10-encoding-normalizer.py --all-paths
  python 03-ai-scripts/10-encoding-normalizer.py --all

  # 3. Run sequentially (1 worker):
  python 03-ai-scripts/10-encoding-normalizer.py --sync

  # 4. Save report to a file:
  python 03-ai-scripts/10-encoding-normalizer.py -o tmp/encoding-report.txt

  # 5. Output results as machine-readable JSON:
  python 03-ai-scripts/10-encoding-normalizer.py --json

  # 6. Fix CRLF and BOM in-place:
  python 03-ai-scripts/10-encoding-normalizer.py --fix
        """,
    )
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory or subfolder")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument("--fix", action="store_true", help="Fix CRLF and BOM in-place")
    parser.add_argument("--ext", help="Comma-separated extensions to scan (e.g. .md,.ts,.py)")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    exit_code = run_encoding_normalizer(
        target_dir=target_path,
        is_fix_mode=args.fix,
        extensions=args.ext,
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
