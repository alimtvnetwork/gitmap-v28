#!/usr/bin/env python3
"""
26-go-code-formatter.py — Cross-platform Go code formatter using gofmt.

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
from importlib import import_module
from pathlib import Path
import shutil
import subprocess
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

normalize_rel_path = engine.normalize_rel_path
stream_directory_files = engine.stream_directory_files
ExitCodeType = engine.ExitCodeType


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


def format_go_file(gofmt_exe: str, file_path: Path) -> bool:
    """Executes gofmt -w on a single Go file."""
    res = subprocess.run([gofmt_exe, "-w", str(file_path)], capture_output=True, text=True)
    if res.returncode != 0:
        print(f"✗ Failed to format {file_path}: {res.stderr.strip()}", file=sys.stderr)
        return False

    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="Cross-platform Go code formatter")
    parser.add_argument("paths", nargs="*", help="Specific files or directories to format")
    parser.add_argument("--staged", action="store_true", help="Format only staged git files")
    args = parser.parse_args()

    gofmt_exe = shutil.which("gofmt")
    if not gofmt_exe:
        print("⚠ gofmt not found in PATH — install Go toolchain (https://go.dev/dl/)", file=sys.stderr)
        return int(ExitCodeType.TOOL_ERROR.value)

    repo_root = Path(__file__).resolve().parent.parent

    target_files: list[Path] = []
    if args.staged:
        target_files = get_staged_go_files(repo_root)
        print(f"Formatting {len(target_files)} staged Go file(s)...")
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

    if not target_files:
        print("✓ No Go files to format.")
        return int(ExitCodeType.SUCCESS.value)

    has_error = False
    for tf in target_files:
        rel = normalize_rel_path(tf.relative_to(repo_root))
        is_success = format_go_file(gofmt_exe, tf)
        if is_success:
            print(f"  ✓ Formatted: {rel}")
        else:
            has_error = True

    if has_error:
        return int(ExitCodeType.VIOLATIONS_FOUND.value)

    print(f"✓ Successfully processed {len(target_files)} Go file(s).")
    return int(ExitCodeType.SUCCESS.value)


if __name__ == "__main__":
    sys.exit(main())
