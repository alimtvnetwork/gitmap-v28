#!/usr/bin/env python3
"""
26-go-code-formatter.py — High-Performance Parallel Go Code Formatter using gofmt.

Modes:
  python 03-ai-scripts/26-go-code-formatter.py                 # format every .go file in repo
  python 03-ai-scripts/26-go-code-formatter.py --staged        # format only staged .go files
  python 03-ai-scripts/26-go-code-formatter.py path/to/file.go # format specific file(s)

Exit codes:
  0 — clean or formatted successfully
  1 — tool missing or formatting error
"""

from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from importlib import import_module
import os
from pathlib import Path
import shutil
import subprocess
import sys
import threading
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

normalize_rel_path = engine.normalize_rel_path
stream_directory_files = engine.stream_directory_files
ExitCodeType = engine.ExitCodeType

CHUNK_SIZE = 30


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


def format_chunk(gofmt_exe: str, chunk: list[Path]) -> tuple[bool, list[str]]:
    """Executes gofmt -w on a chunk of Go files."""
    args = [gofmt_exe, "-w"] + [str(p) for p in chunk]
    res = subprocess.run(args, capture_output=True, text=True)
    if res.returncode != 0:
        err = res.stderr.strip() or "gofmt execution failure"
        return False, [err]

    return True, []


def chunk_list(items: list[Path], size: int) -> list[list[Path]]:
    """Splits a list into chunks of at most `size` items."""
    return [items[i:i + size] for i in range(0, len(items), size)]


def print_progress(completed: int, total: int, workers: int, start_time: float) -> None:
    """Emits live formatted progress percentage and throughput."""
    pct = (completed / total * 100.0) if total > 0 else 100.0
    elapsed = max(0.001, time.time() - start_time)
    fps = completed / elapsed
    msg = f"\rFormatting Go files: [ {completed:4d}/{total:4d} ] {pct:5.1f}% | {workers} workers | {fps:5.1f} files/sec"
    sys.stdout.write(msg)
    sys.stdout.flush()


def run_parallel_formatting(gofmt_exe: str, target_files: list[Path]) -> bool:
    """Formats all files in parallel using ThreadPoolExecutor across all CPU cores."""
    total_files = len(target_files)
    cpu_cores = os.cpu_count() or 16
    chunks = chunk_list(target_files, CHUNK_SIZE)
    completed_count = 0
    lock = threading.Lock()
    has_error = False
    start_time = time.time()

    print(f"▸ Discovered {total_files} Go file(s). Formatting with {cpu_cores} worker threads...")
    print_progress(0, total_files, cpu_cores, start_time)

    with ThreadPoolExecutor(max_workers=cpu_cores) as pool:
        futures = {pool.submit(format_chunk, gofmt_exe, chunk): len(chunk) for chunk in chunks}
        for fut in as_completed(futures):
            chunk_len = futures[fut]
            is_success, errors = fut.result()
            with lock:
                if not is_success:
                    has_error = True
                    for err in errors:
                        sys.stderr.write(f"\n✗ Error: {err}\n")
                completed_count += chunk_len
                print_progress(completed_count, total_files, cpu_cores, start_time)

    elapsed = time.time() - start_time
    sys.stdout.write("\n")
    if not has_error:
        print(f"✓ Successfully formatted {total_files} Go file(s) across {cpu_cores} CPU cores in {elapsed:.2f}s.")
    return not has_error


def collect_target_files(args: argparse.Namespace, repo_root: Path) -> list[Path]:
    """Resolves target Go files based on command-line flags."""
    if args.staged:
        files = get_staged_go_files(repo_root)
        print(f"Formatting {len(files)} staged Go file(s)...")
        return files

    if args.paths:
        target_files = []
        for p_str in args.paths:
            p = Path(p_str).resolve()
            if p.is_file() and p.suffix == ".go":
                target_files.append(p)
            elif p.is_dir():
                target_files.extend(list(p.rglob("*.go")))
        return target_files

    gitmap_dir = repo_root / "gitmap"
    search_root = gitmap_dir if gitmap_dir.is_dir() else repo_root
    excludes = {".lovable", ".git", ".tmp", "temp-scripts", "scratch", "node_modules", "dist", "bin"}
    return list(stream_directory_files(search_root, extensions=[".go"], custom_excludes=excludes))


def main() -> int:
    parser = argparse.ArgumentParser(description="Parallel Cross-Platform Go Code Formatter")
    parser.add_argument("paths", nargs="*", help="Specific files or directories to format")
    parser.add_argument("--staged", action="store_true", help="Format only staged git files")
    args = parser.parse_args()

    gofmt_exe = shutil.which("gofmt")
    if not gofmt_exe:
        print("⚠ gofmt not found in PATH — install Go toolchain (https://go.dev/dl/)", file=sys.stderr)
        return int(ExitCodeType.TOOL_ERROR.value)

    repo_root = Path(__file__).resolve().parent.parent
    target_files = collect_target_files(args, repo_root)

    if not target_files:
        print("✓ No Go files to format.")
        return int(ExitCodeType.SUCCESS.value)

    is_clean = run_parallel_formatting(gofmt_exe, target_files)
    if not is_clean:
        return int(ExitCodeType.VIOLATIONS_FOUND.value)

    return int(ExitCodeType.SUCCESS.value)


if __name__ == "__main__":
    sys.exit(main())
