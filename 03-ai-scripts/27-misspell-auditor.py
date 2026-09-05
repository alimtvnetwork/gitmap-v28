#!/usr/bin/env python3
"""
27-misspell-auditor.py — Cross-platform American English spelling auditor.

Audits repository files to enforce American English spelling and prevent
non-US English spelling regressions.
Executes an optimized parallel worker pool with fast pre-filtering, quiet tick
on success, and detailed traces on failure.

Usage:
  python 03-ai-scripts/27-misspell-auditor.py              # scan repository files in parallel
  python 03-ai-scripts/27-misspell-auditor.py --all-paths  # show full logs and report
  python 03-ai-scripts/27-misspell-auditor.py --staged     # scan staged files
  python 03-ai-scripts/27-misspell-auditor.py --fix        # auto-replace misspellings
"""

from __future__ import annotations

import argparse
from importlib import import_module
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

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

normalize_rel_path = engine.normalize_rel_path
read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf
stream_directory_files = engine.stream_directory_files
ExitCodeType = engine.ExitCodeType
WorkerResult = engine.WorkerResult
run_worker_pool = engine.run_worker_pool
add_worker_cli_arguments = engine.add_worker_cli_arguments
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS

# Common British -> US English mapping enforced across meta-repos
# Stored with split string concatenation so static spell-checkers do not trigger on the lookup table
_SPELLING_PAIRS: list[tuple[str, str]] = [
    ("behavi" + "our", "behavior"),
    ("behavi" + "ours", "behaviors"),
    ("col" + "our", "color"),
    ("col" + "ours", "colors"),
    ("initial" + "ise", "initialize"),
    ("initial" + "ised", "initialized"),
    ("initial" + "ising", "initializing"),
    ("initial" + "isation", "initialization"),
    ("custom" + "ise", "customize"),
    ("custom" + "ised", "customized"),
    ("custom" + "ising", "customizing"),
    ("custom" + "isation", "customization"),
    ("synchron" + "ise", "synchronize"),
    ("synchron" + "ised", "synchronized"),
    ("synchron" + "ising", "synchronizing"),
    ("optim" + "ise", "optimize"),
    ("optim" + "ised", "optimized"),
    ("optim" + "ising", "optimizing"),
    ("optim" + "isation", "optimization"),
    ("priorit" + "ise", "prioritize"),
    ("priorit" + "ised", "prioritized"),
    ("priorit" + "ising", "prioritizing"),
    ("serial" + "ise", "serialize"),
    ("serial" + "ised", "serialized"),
    ("serial" + "ising", "serializing"),
    ("serial" + "isation", "serialization"),
    ("normal" + "ise", "normalize"),
    ("normal" + "ised", "normalized"),
    ("normal" + "ising", "normalizing"),
    ("normal" + "isation", "normalization"),
    ("cancell" + "ing", "canceling"),
    ("cancell" + "ed", "canceled"),
]
SPELLING_MAP: dict[str, str] = dict(_SPELLING_PAIRS)
COMPILED_PATTERNS: list[tuple[re.Pattern, str, str]] = [
    (re.compile(rf"\b{brit}\b", re.IGNORECASE), brit, us) for brit, us in SPELLING_MAP.items()
]

EXCLUDED_EXTENSIONS = {
    ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".webp",
    ".zip", ".tar", ".gz", ".exe", ".bin", ".lock", ".lockb",
}


def get_staged_files(repo_root: Path) -> list[Path]:
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

    files = []
    for line in res.stdout.splitlines():
        rel = line.strip()
        f = repo_root / rel
        if f.is_file() and f.suffix not in EXCLUDED_EXTENSIONS:
            files.append(f)

    return files


def audit_file(file_path: Path, is_fix_mode: bool) -> list[tuple[int, str, str]]:
    content = read_file_lf(file_path)
    if not content:
        return []

    # Fast substring pre-filter before regex execution
    content_lower = content.lower()
    if not any(k in content_lower for k in SPELLING_MAP):
        return []

    findings: list[tuple[int, str, str]] = []
    lines = content.split("\n")
    modified = False

    for idx, line in enumerate(lines, start=1):
        line_lower = line.lower()
        for pattern, brit, us in COMPILED_PATTERNS:
            if brit in line_lower and pattern.search(line):
                findings.append((idx, brit, us))
                if is_fix_mode:
                    lines[idx - 1] = pattern.sub(us, lines[idx - 1])
                    modified = True

    if is_fix_mode and modified:
        write_file_lf(file_path, "\n".join(lines))

    return findings

def execute_file_spelling_task(task_tuple: tuple[Path, bool]) -> WorkerResult:
    file_path, is_fix_mode = task_tuple
    start_time = time.perf_counter()
    try:
        rel_path = file_path.resolve().relative_to(Path.cwd().resolve())
        norm_path = normalize_rel_path(rel_path)
    except Exception:
        norm_path = normalize_rel_path(file_path)

    if "misspell" in file_path.name.lower():
        return WorkerResult(
            name=norm_path,
            is_success=True,
            output="Excluded self misspell script",
            elapsed_sec=round(time.perf_counter() - start_time, 4),
        )

    try:
        findings = audit_file(file_path, is_fix_mode=is_fix_mode)
        elapsed = round(time.perf_counter() - start_time, 4)
        if findings:
            action = "Fixed" if is_fix_mode else "Found"
            err_lines = [f"  ✗ {norm_path}:{line_no}: {action} '{brit}' -> use American spelling '{us}'" for line_no, brit, us in findings]
            is_pass = bool(is_fix_mode)
            return WorkerResult(
                name=norm_path,
                is_success=is_pass,
                error="\n".join(err_lines),
                output="\n".join(err_lines) if is_fix_mode else "",
                elapsed_sec=elapsed,
            )
        return WorkerResult(
            name=norm_path,
            is_success=True,
            output=f"Spelling verified: {norm_path}",
            elapsed_sec=elapsed,
        )
    except Exception as exc:
        return WorkerResult(
            name=norm_path,
            is_success=False,
            error=f"Exception auditing {norm_path}: {exc}",
            elapsed_sec=round(time.perf_counter() - start_time, 4),
        )


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/27-misspell-auditor.py",
        description="Cross-platform American English spelling auditor using parallel worker pool.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: scan repository files in parallel; quiet on success (tick), detailed on error:
  python 03-ai-scripts/27-misspell-auditor.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/27-misspell-auditor.py --all-paths
  python 03-ai-scripts/27-misspell-auditor.py --all

  # 3. Run sequentially (1 worker):
  python 03-ai-scripts/27-misspell-auditor.py --sync

  # 4. Save report to a file:
  python 03-ai-scripts/27-misspell-auditor.py -o tmp/misspell-report.txt

  # 5. Output results as machine-readable JSON:
  python 03-ai-scripts/27-misspell-auditor.py --json

  # 6. Automatically fix detected British spellings across all files:
  python 03-ai-scripts/27-misspell-auditor.py --fix
        """,
    )
    parser.add_argument("--staged", action="store_true", help="Audit only staged git files")
    parser.add_argument("--fix", action="store_true", help="Auto-replace British English spellings to US English")
    parser.add_argument("paths", nargs="*", help="Specific files or paths to audit")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent

    target_files: list[Path] = []
    if args.staged:
        target_files = get_staged_files(repo_root)
    elif args.paths:
        for p_str in args.paths:
            p = Path(p_str).resolve()
            if p.is_file():
                target_files.append(p)
            elif p.is_dir():
                for ext in [".md", ".go", ".ts", ".py", ".json", ".sh", ".ps1"]:
                    target_files.extend(list(p.rglob(f"*{ext}")))
    else:
        for f in stream_directory_files(repo_root, extensions=[".md", ".go", ".ts", ".py", ".json", ".sh", ".ps1"]):
            target_files.append(f)

    if args.filter:
        filt = args.filter.lower().replace("\\", "/")
        target_files = [f for f in target_files if filt in normalize_rel_path(f.relative_to(repo_root)).lower()]

    if not target_files:
        print(f"✔ All passed. (0 file(s) in 0.00s)")
        sys.exit(0)

    task_items = [(f, args.fix) for f in target_files]

    exit_code = run_worker_pool(
        items=task_items,
        worker_fn=execute_file_spelling_task,
        max_workers=args.workers,
        is_sync=args.is_sync,
        show_all=args.show_all,
        output_file=args.output_file,
        as_json=args.as_json,
        title="AMERICAN ENGLISH SPELLING AUDITOR",
        item_noun="file(s)",
    )
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
