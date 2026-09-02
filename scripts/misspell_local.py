#!/usr/bin/env python3
"""scripts/misspell_local.py — Run misspell check locally against file filters in config.json.

Usage:
  python scripts/misspell_local.py                 # diff vs origin/main
  python scripts/misspell_local.py --base HEAD~1   # diff vs arbitrary ref
  python scripts/misspell_local.py --all           # scan every tracked file
  python scripts/misspell_local.py --staged        # scan staged (index) files
  python scripts/misspell_local.py --files a.md ...# scan explicit list

Exit codes:
  0 = clean, 1 = misspellings found, 2 = bad usage / missing tool.
"""

from __future__ import annotations

import argparse
import fnmatch
import json
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path

MISSPELL_VERSION = "v0.3.4"
DEFAULT_EXCLUDES = [
    "*.png", "*.jpg", "*.jpeg", "*.gif", "*.ico", "*.svg", "*.webp",
    "*.zip", "*.tar", "*.gz", "*.exe",
    "*/testdata/*", "*/golden/*",
    "*/.gitmap/release/*", "*/.gitmap/release-assets/*",
    "gitmap/completion/allcommands_generated.go",
    ".lovable/*",
]


def ensure_misspell() -> str:
    exe = shutil.which("misspell")
    if exe:
        return exe
    go_exe = shutil.which("go")
    if not go_exe:
        print("misspell not found and Go toolchain missing to install it", file=sys.stderr)
        sys.exit(2)
    print(f"misspell not found — installing pinned {MISSPELL_VERSION}...", file=sys.stderr)
    res = subprocess.run([go_exe, "install", f"github.com/client9/misspell/cmd/misspell@{MISSPELL_VERSION}"])
    if res.returncode != 0:
        print(f"failed to install misspell {MISSPELL_VERSION}", file=sys.stderr)
        sys.exit(2)
    try:
        gopath = subprocess.check_output([go_exe, "env", "GOPATH"], text=True).strip()
        candidate = Path(gopath) / "bin" / ("misspell.exe" if os.name == "nt" else "misspell")
        if candidate.is_file():
            return str(candidate)
    except Exception:
        pass
    exe = shutil.which("misspell")
    if exe:
        return exe
    print("Could not locate installed misspell binary", file=sys.stderr)
    sys.exit(2)


def load_patterns(config_path: Path) -> tuple[list[str], list[str]]:
    if not config_path.is_file():
        return DEFAULT_EXCLUDES, []
    try:
        with open(config_path, "r", encoding="utf-8") as f:
            data = json.load(f)
        m = data.get("misspell", {})
        ex = m.get("exclude", []) or DEFAULT_EXCLUDES
        inc = m.get("include", []) or []
        return ex, inc
    except OSError:
        return DEFAULT_EXCLUDES, []


def matches_any(rel_path: str, patterns: list[str]) -> bool:
    for pat in patterns:
        if fnmatch.fnmatch(rel_path, pat):
            return True
    return False


def main() -> int:
    parser = argparse.ArgumentParser(description="Run misspell check locally.")
    parser.add_argument("--all", action="store_true", help="Scan every tracked file")
    parser.add_argument("--staged", action="store_true", help="Scan staged files")
    parser.add_argument("--base", default="origin/main", help="Base ref for diff mode")
    parser.add_argument("--files", nargs="+", help="Scan explicit list")
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent
    misspell_bin = ensure_misspell()
    config_path = repo_root / "gitmap" / "data" / "config.json"
    excludes, includes = load_patterns(config_path)

    candidates: list[str] = []
    mode = "diff"
    if args.files:
        candidates = args.files
        mode = "files"
    elif args.all:
        mode = "all"
        try:
            out = subprocess.check_output(["git", "ls-files"], cwd=str(repo_root), text=True)
            candidates = [l.strip() for l in out.splitlines() if l.strip()]
        except subprocess.SubprocessError:
            pass
    elif args.staged:
        mode = "staged"
        try:
            out = subprocess.check_output(["git", "diff", "--name-only", "--diff-filter=AM", "--cached"], cwd=str(repo_root), text=True)
            candidates = [l.strip() for l in out.splitlines() if l.strip()]
        except subprocess.SubprocessError:
            pass
    else:
        try:
            out = subprocess.check_output(["git", "diff", "--name-only", "--diff-filter=AM", f"{args.base}...HEAD"], cwd=str(repo_root), text=True)
            candidates = [l.strip() for l in out.splitlines() if l.strip()]
        except subprocess.SubprocessError:
            print(f"base ref '{args.base}' not found or git diff failed", file=sys.stderr)
            return 2

    if not candidates:
        print("no candidate files — nothing to scan")
        return 0

    targets: list[str] = []
    for c in candidates:
        posix = c.replace("\\", "/")
        p = repo_root / c
        if not p.is_file():
            continue
        if matches_any(posix, excludes):
            continue
        if includes and not matches_any(posix, includes):
            continue
        targets.append(str(p))

    if not targets:
        print("no scannable files after filters — nothing to scan")
        return 0

    print(f"misspell-local: mode={mode} | excludes={len(excludes)} includes={len(includes)} | scanning {len(targets)} file(s)")

    # Run misspell in batches if needed
    findings = 0
    batch_size = 100
    for i in range(0, len(targets), batch_size):
        batch = targets[i:i + batch_size]
        res = subprocess.run([misspell_bin, "-locale", "US"] + batch, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
        if res.stdout:
            sys.stdout.write(res.stdout)
            for line in res.stdout.splitlines():
                if re.match(r"^[^:]+:\d+:\d+:", line):
                    findings += 1

    if findings > 0:
        print(f"FAIL: {findings} misspelling(s) found.", file=sys.stderr)
        return 1

    print("OK: no misspellings.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
