#!/usr/bin/env python3
"""
Sequence, Numbering & Title Header Auditor
Audits numbered markdown files (e.g. 01-intro.md) to ensure no sequence gaps and
verifies that the primary # H1 title matches the file sequence prefix.
Multi-folder capable, customizable extensions, thread-safe regex engine, and
concurrent execution via the shared worker pool engine.

Performance & Clean Architecture:
1. Flattened Conditionals: Zero deep nesting using early `continue` guard clauses.
2. Concurrent Worker Pool: Audits directories in parallel across CPU cores, quiet on success, detailed on failure.
3. All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

from __future__ import annotations

import argparse
from dataclasses import dataclass
from importlib import import_module
import os
from pathlib import Path
import re
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf
normalize_rel_path = engine.normalize_rel_path
is_ignored_directory = engine.is_ignored_directory
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
WorkerResult = engine.WorkerResult
run_worker_pool = engine.run_worker_pool
add_worker_cli_arguments = engine.add_worker_cli_arguments
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS


@dataclass
class DirectorySequenceTask:
    dir_path: Path
    files: list[str]
    is_fix_mode: bool


def audit_single_directory(task: DirectorySequenceTask) -> WorkerResult:
    start_time = time.perf_counter()
    re_num = get_compiled_regex(RegexPatternType.FILE_NUM_PREFIX)
    re_h1 = get_compiled_regex(RegexPatternType.H1_HEADER)
    seq_issues = []
    title_issues = []

    try:
        rel_dir = normalize_rel_path(task.dir_path.resolve().relative_to(Path.cwd().resolve()))
    except Exception:
        rel_dir = normalize_rel_path(task.dir_path)

    numbered_files = []
    for f in task.files:
        m = re_num.match(f)
        if not m:
            continue
        numbered_files.append((int(m.group(1)), f, str(task.dir_path / f)))

    if not numbered_files:
        return WorkerResult(
            name=rel_dir,
            is_success=True,
            output=f"No numbered markdown files in {rel_dir}",
            elapsed_sec=round(time.perf_counter() - start_time, 4),
        )

    numbered_files.sort(key=lambda x: x[0])
    first_num = numbered_files[0][0]
    if first_num in (0, 1):
        expected = first_num
        for num_val, f_name, _ in numbered_files:
            if num_val != expected:
                seq_issues.append(f"{rel_dir}/{f_name} (found {num_val:02d}, expected {expected:02d})")
            expected += 1

    for num_val, f_name, full_path in numbered_files:
        try:
            content = read_file_lf(full_path, encoding=DEFAULT_ENCODING)
            m_h1 = re_h1.search(content)
            if not m_h1:
                continue

            h1_num = int(m_h1.group(2))
            if h1_num == num_val:
                continue

            title_issues.append(f"{rel_dir}/{f_name} (file prefix {num_val:02d} != H1 header {h1_num:02d})")
            if task.is_fix_mode:
                def replacer(match: re.Match) -> str:
                    return f"{match.group(1)}{num_val:02d}{match.group(3)}{match.group(4)}"
                new_content = re_h1.sub(replacer, content, count=1)
                write_file_lf(full_path, new_content, encoding=DEFAULT_ENCODING)
        except Exception:
            pass

    elapsed = round(time.perf_counter() - start_time, 4)
    has_errors = bool(seq_issues) or (bool(title_issues) and not task.is_fix_mode)

    error_lines = []
    if seq_issues:
        error_lines.extend([f"  ::error::{issue}" for issue in seq_issues])
    if title_issues:
        prefix = "Fixed" if task.is_fix_mode else "Warning"
        error_lines.extend([f"  ::{prefix.lower()}::{issue}" for issue in title_issues])

    if has_errors:
        return WorkerResult(
            name=rel_dir,
            is_success=False,
            error="\n".join(error_lines),
            elapsed_sec=elapsed,
        )

    return WorkerResult(
        name=rel_dir,
        is_success=True,
        output=f"Verified {len(numbered_files)} numbered file(s) in {rel_dir}",
        elapsed_sec=elapsed,
    )


def run_sequence_auditor(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = False,
    max_workers: int | None = None,
    is_sync: bool = False,
    show_all: bool = False,
    output_file: str | None = None,
    as_json: bool | str = False,
    filter_pattern: str | None = None,
) -> int:
    """Executes sequence and title audit across target directory using parallel worker pool."""
    target_path = Path(target_dir).resolve()
    re_num = get_compiled_regex(RegexPatternType.FILE_NUM_PREFIX)

    tasks: list[DirectorySequenceTask] = []
    for root, dirs, files in os.walk(str(target_path)):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
        md_files = [f for f in files if f.endswith(".md") and re_num.match(f)]
        if md_files:
            tasks.append(
                DirectorySequenceTask(
                    dir_path=Path(root),
                    files=md_files,
                    is_fix_mode=is_fix_mode,
                )
            )

    if filter_pattern:
        filt = filter_pattern.lower().replace("\\", "/")
        tasks = [t for t in tasks if filt in normalize_rel_path(t.dir_path).lower()]

    if not tasks:
        print(f"✔ All passed. (0 directory sequence group(s) in 0.00s)")
        return ExitCodeType.SUCCESS.value

    return run_worker_pool(
        items=tasks,
        worker_fn=audit_single_directory,
        max_workers=max_workers,
        is_sync=is_sync,
        show_all=show_all,
        output_file=output_file,
        as_json=as_json,
        title="SEQUENCE NUMBERING & H1 HEADER AUDITOR",
        item_noun="directory sequence group(s)",
    )


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/15-sequence-and-title-auditor.py",
        description="Audit file sequence numbering and H1 title alignment using parallel worker pool.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: run all audits in parallel; quiet on success (tick), detailed on error:
  python 03-ai-scripts/15-sequence-and-title-auditor.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/15-sequence-and-title-auditor.py --all-paths
  python 03-ai-scripts/15-sequence-and-title-auditor.py --all

  # 3. Run sequentially (1 worker):
  python 03-ai-scripts/15-sequence-and-title-auditor.py --sync

  # 4. Save report to a file:
  python 03-ai-scripts/15-sequence-and-title-auditor.py -o tmp/seq-report.txt

  # 5. Output results as machine-readable JSON:
  python 03-ai-scripts/15-sequence-and-title-auditor.py --json

  # 6. Auto-fix H1 title numbers in markdown files:
  python 03-ai-scripts/15-sequence-and-title-auditor.py --fix
        """,
    )
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory to audit")
    parser.add_argument("--dir", "--path", "-p", dest="opt_dir", help="Directory to audit")
    parser.add_argument("--fix", action="store_true", help="Auto-fix H1 title numbers in markdown files")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    target_path = args.opt_dir or args.path or CURRENT_DIR
    exit_code = run_sequence_auditor(
        target_dir=target_path,
        is_fix_mode=args.fix,
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
