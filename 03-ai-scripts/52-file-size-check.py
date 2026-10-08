#!/usr/bin/env python3
"""
Go File Size Check (spec 243.4 enforcement)
Fails when any .go file under cli/ crosses the line limit (default 300).
Run via: gitmap py 03-ai-scripts/52-file-size-check.py

Modes:
  default        scan every .go file under cli/
  --staged       check only staged .go files (pre-commit ratchet: changed files only,
                 so the grandfathered 300-1000 band never blocks unrelated commits)
  --diff REF     check .go files changed versus REF (e.g. --diff HEAD)
  --files ...    check explicit files
  --full         same as default (explicit full-repo scan)
  --limit N      line limit override (default 300)
  --root DIR     scope root for .go discovery (default cli)

Exit code: 0 when all files are within the limit, 1 on any violation.
The check_file_sizes() entry point is imported by 03-ai-scripts/50-fastgate.py
so the pre-commit fast-gate enforces the cap on staged files.
"""

from __future__ import annotations

import argparse
from importlib import import_module
import os
from pathlib import Path
import shutil
import subprocess
import sys
from typing import Sequence

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).resolve().parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR

DEFAULT_LINE_LIMIT = 300
DEFAULT_SCOPE_DIR = "cli"
GO_EXTENSION = ".go"


def get_repo_root() -> Path:
    """Returns canonical repository root path."""
    return Path(__file__).resolve().parent.parent


def run_git_lines(cmd_args: list[str], root_dir: Path) -> list[str]:
    """Executes git and returns stripped output lines (empty on any failure)."""
    git_bin = shutil.which("git") or "git"
    res = subprocess.run(
        [git_bin] + cmd_args, cwd=str(root_dir),
        capture_output=True, text=True, encoding=DEFAULT_ENCODING,
    )
    has_output = bool(res.returncode == 0 and res.stdout)
    if not has_output:
        return []
    return [line.strip().replace("\\", "/") for line in res.stdout.splitlines() if line.strip()]


def get_staged_go_files(root_dir: Path) -> list[str]:
    """Returns staged .go files from the git index."""
    args = ["diff", "--cached", "--name-only", "--diff-filter=ACMR"]
    return [f for f in run_git_lines(args, root_dir) if f.endswith(GO_EXTENSION)]


def get_diff_go_files(root_dir: Path, ref: str) -> list[str]:
    """Returns .go files changed versus the given git ref."""
    args = ["diff", ref, "--name-only", "--diff-filter=ACMR"]
    return [f for f in run_git_lines(args, root_dir) if f.endswith(GO_EXTENSION)]


def discover_go_files(scope_dir: Path) -> list[str]:
    """Discovers every .go file under the scope directory."""
    found: list[str] = []
    if not scope_dir.is_dir():
        return found
    for dirpath, dirnames, filenames in os.walk(scope_dir):
        dirnames[:] = [d for d in dirnames if d != ".git"]
        for name in filenames:
            if name.endswith(GO_EXTENSION):
                found.append(str(Path(dirpath) / name))
    return found


def count_file_lines(path: Path) -> int:
    """Counts lines in a file (0 when unreadable)."""
    try:
        with path.open("r", encoding=DEFAULT_ENCODING, errors="replace") as handle:
            return sum(1 for _ in handle)
    except OSError:
        return 0


def check_go_file_sizes(
    files: Sequence[str | Path],
    root_dir: Path,
    limit: int = DEFAULT_LINE_LIMIT,
) -> list[tuple[str, int]]:
    """Returns [(display_path, line_count)] for files strictly over the limit."""
    violations: list[tuple[str, int]] = []
    for entry in files:
        path = Path(entry)
        if not path.is_absolute():
            path = root_dir / path
        if not path.is_file():
            continue
        line_count = count_file_lines(path)
        is_over = line_count > limit
        if is_over:
            try:
                display = str(path.relative_to(root_dir))
            except ValueError:
                display = str(path)
            violations.append((display.replace("\\", "/"), line_count))
    violations.sort(key=lambda item: item[1], reverse=True)
    return violations


def resolve_target_files(args: argparse.Namespace, root_dir: Path) -> list[str]:
    """Resolves target files from CLI flags (explicit > staged/diff > full scan)."""
    has_explicit = bool(args.files)
    if has_explicit:
        return [f.replace("\\", "/") for f in args.files if f.strip()]
    if args.staged:
        return get_staged_go_files(root_dir)
    if args.diff:
        return get_diff_go_files(root_dir, args.diff)
    scope = root_dir / args.root
    return [str(Path(p).relative_to(root_dir)).replace("\\", "/") for p in discover_go_files(scope)]


def print_report(violations: list[tuple[str, int]], limit: int, scanned: int) -> None:
    """Prints the per-file violation report."""
    if not violations:
        print(f"✔ Go file-size check passed: {scanned} file(s) scanned, none over {limit} lines.")
        return
    print(f"❌ Go file-size check failed: {len(violations)} file(s) over {limit} lines ({scanned} scanned):")
    for rel_path, line_count in violations:
        print(f"   • {rel_path}: {line_count} lines (over by {line_count - limit})")
    print("Split oversized files by concern (spec 243.4) before committing.")


def parse_cli_arguments() -> argparse.Namespace:
    """Parses command line arguments for the go file-size check."""
    parser = argparse.ArgumentParser(description="Go file line-count enforcement (spec 243.4)")
    parser.add_argument("--limit", type=int, default=DEFAULT_LINE_LIMIT, help="Max allowed lines per .go file")
    parser.add_argument("--staged", action="store_true", help="Check only staged .go files (pre-commit ratchet)")
    parser.add_argument("--diff", default="", help="Check .go files changed versus this git ref (e.g. HEAD)")
    parser.add_argument("--files", nargs="*", default=[], help="Explicit files to check")
    parser.add_argument("--full", action="store_true", help="Full scan of the scope dir (default behavior)")
    parser.add_argument("--root", default=DEFAULT_SCOPE_DIR, help="Scope dir for .go discovery (default cli)")
    return parser.parse_args()


def execute_file_size_check(args: argparse.Namespace) -> int:
    """Executes the file-size check and returns the process exit code."""
    root_dir = get_repo_root()
    files = resolve_target_files(args, root_dir)
    violations = check_go_file_sizes(files, root_dir, args.limit)
    print_report(violations, args.limit, len(files))
    is_clean = not violations
    return ExitCodeType.SUCCESS.value if is_clean else ExitCodeType.VIOLATIONS_FOUND.value


def main() -> None:
    """Main execution entry point."""
    sys.exit(execute_file_size_check(parse_cli_arguments()))


if __name__ == "__main__":
    main()
