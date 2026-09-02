#!/usr/bin/env python3
"""fix_repo.py — rewrite prior versioned-repo-name tokens to current (cross-platform).

Spec: spec-authoring/22-fix-repo/01-spec.md

Usage:
  python fix_repo.py                  # default: replace last 2 versions
  python fix_repo.py --2              # explicit
  python fix_repo.py --3              # last 3 versions
  python fix_repo.py --5              # last 5 versions
  python fix_repo.py --all            # every prior version (1..Current-1)
  python fix_repo.py --dry-run        # report only
  python fix_repo.py --verbose        # list each modified file
"""

from __future__ import annotations

import argparse
import fnmatch
import json
import os
import re
import subprocess
import sys
from pathlib import Path

EXIT_OK = 0
EXIT_NOT_A_REPO = 2
EXIT_NO_REMOTE = 3
EXIT_NO_VERSION_SUFFIX = 4
EXIT_BAD_VERSION = 5
EXIT_BAD_FLAG = 6
EXIT_WRITE_FAILED = 7
EXIT_BAD_CONFIG = 8

MAX_FILE_BYTES = 5 * 1024 * 1024
BINARY_EXTENSIONS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip",
    ".tar", ".gz", ".tgz", ".bz2", ".xz", ".7z", ".rar", ".woff",
    ".woff2", ".ttf", ".otf", ".eot", ".mp3", ".mp4", ".mov", ".wav",
    ".ogg", ".webm", ".class", ".jar", ".so", ".dylib", ".dll", ".exe",
    ".pyc", ".db",
}


def get_repo_root() -> Path | None:
    try:
        out = subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"],
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        return Path(out) if out else None
    except subprocess.SubprocessError:
        return None


def get_remote_url(repo_root: Path) -> str:
    try:
        out = subprocess.check_output(
            ["git", "remote", "get-url", "origin"],
            cwd=str(repo_root),
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        if out:
            return out
    except subprocess.SubprocessError:
        pass
    try:
        out = subprocess.check_output(
            ["git", "remote", "-v"],
            cwd=str(repo_root),
            stderr=subprocess.DEVNULL,
            text=True,
        )
        for line in out.splitlines():
            parts = line.split()
            if len(parts) >= 2 and "(fetch)" in line:
                return parts[1]
    except subprocess.SubprocessError:
        pass
    return ""


def parse_remote_url(url: str) -> tuple[str, str, str] | None:
    trimmed = url.rstrip("/")
    if trimmed.endswith(".git"):
        trimmed = trimmed[:-4]
    patterns = [
        r"^https?://([^/:]+)(?::\d+)?/([^/]+)/([^/]+)$",
        r"^git@([^:]+):([^/]+)/([^/]+)$",
        r"^ssh://git@([^/:]+)(?::\d+)?/([^/]+)/([^/]+)$",
    ]
    for pat in patterns:
        m = re.match(pat, trimmed)
        if m:
            return m.group(1), m.group(2), m.group(3)
    return None


def split_repo_version(repo_name: str) -> tuple[str, int] | None:
    m = re.match(r"^(.+)-v(\d+)$", repo_name)
    if m:
        return m.group(1), int(m.group(2))
    return None


def load_config(config_path: str | None, repo_root: Path) -> tuple[list[str], list[str]]:
    cfg_file = Path(config_path) if config_path else repo_root / "fix-repo.config.json"
    if not cfg_file.is_file():
        return [], []
    try:
        with open(cfg_file, "r", encoding="utf-8") as f:
            data = json.load(f)
        dirs = [d for d in data.get("ignoreDirs", []) if isinstance(d, str)]
        pats = [p for p in data.get("ignorePatterns", []) if isinstance(p, str)]
        return dirs, pats
    except Exception as e:
        print(f"fix-repo: ERROR reading config {cfg_file}: {e}", file=sys.stderr)
        sys.exit(EXIT_BAD_CONFIG)


def is_ignored_path(rel_posix: str, ignore_dirs: list[str], ignore_pats: list[str]) -> bool:
    for d in ignore_dirs:
        clean_d = d.rstrip("/")
        if rel_posix == clean_d or rel_posix.startswith(clean_d + "/"):
            return True
    for p in ignore_pats:
        if fnmatch.fnmatch(rel_posix, p):
            return True
    return False


def is_scannable(file_path: Path) -> bool:
    if file_path.is_symlink():
        return False
    if file_path.suffix.lower() in BINARY_EXTENSIONS:
        return False
    try:
        st = file_path.stat()
        if st.st_size > MAX_FILE_BYTES:
            return False
        with open(file_path, "rb") as f:
            chunk = f.read(8192)
            if b"\x00" in chunk:
                return False
    except OSError:
        return False
    return True


def rewrite_content(
    content: str, base: str, targets: list[int], current: int
) -> tuple[str, int]:
    total_reps = 0
    new_content = content
    for n in targets:
        token = f"{base}-v{n}"
        replacement = f"{base}-v{current}"
        # Token must NOT be immediately followed by a digit
        pattern = re.escape(token) + r"(?!\d)"
        matches = len(re.findall(pattern, new_content))
        if matches > 0:
            new_content = re.sub(pattern, replacement, new_content)
            total_reps += matches
    return new_content, total_reps


def main() -> int:
    parser = argparse.ArgumentParser(description="Rewrite prior versioned-repo-name tokens to current.")
    group = parser.add_mutually_exclusive_group()
    group.add_argument("--2", dest="mode", action="store_const", const="2")
    group.add_argument("--3", dest="mode", action="store_const", const="3")
    group.add_argument("--5", dest="mode", action="store_const", const="5")
    group.add_argument("--all", dest="mode", action="store_const", const="all")
    parser.set_defaults(mode="2")
    parser.add_argument("--dry-run", action="store_true", help="Report only, do not modify files")
    parser.add_argument("--verbose", action="store_true", help="List each modified file")
    parser.add_argument("--config", help="Path to config file")
    args = parser.parse_args()

    repo_root = get_repo_root()
    if not repo_root:
        print("fix-repo: ERROR not a git repository (E_NOT_A_REPO)", file=sys.stderr)
        return EXIT_NOT_A_REPO

    remote_url = get_remote_url(repo_root)
    if not remote_url:
        print("fix-repo: ERROR no remote URL found (E_NO_REMOTE)", file=sys.stderr)
        return EXIT_NO_REMOTE

    parsed = parse_remote_url(remote_url)
    if not parsed:
        print(f"fix-repo: ERROR cannot parse remote URL '{remote_url}'", file=sys.stderr)
        return EXIT_NO_REMOTE
    host, owner, repo = parsed

    split = split_repo_version(repo)
    if not split:
        print(f"fix-repo: ERROR no -vN suffix on repo name '{repo}' (E_NO_VERSION_SUFFIX)", file=sys.stderr)
        return EXIT_NO_VERSION_SUFFIX
    split_base, current_version = split
    if current_version < 1:
        print("fix-repo: ERROR version <= 0 (E_BAD_VERSION)", file=sys.stderr)
        return EXIT_BAD_VERSION

    ignore_dirs, ignore_pats = load_config(args.config, repo_root)

    span = current_version - 1 if args.mode == "all" else int(args.mode)
    start_v = max(1, current_version - span)
    targets = list(range(start_v, current_version))

    targets_str = " ".join(f"v{n}" for n in targets) if targets else "(none)"
    print(f"fix-repo  base={split_base}  current=v{current_version}  mode=--{args.mode}")
    print(f"targets:  {targets_str}")
    print(f"host:     {host}  owner={owner}\n")

    if not targets:
        print("scanned: 0 files\nchanged: 0 files (0 replacements)\nmode:    dry-run" if args.dry_run else "mode:    write")
        print("fix-repo: nothing to replace")
        return EXIT_OK

    # Gather tracked files using git ls-files
    try:
        raw_files = subprocess.check_output(
            ["git", "ls-files", "-z"],
            cwd=str(repo_root),
            stderr=subprocess.DEVNULL,
        ).split(b"\x00")
    except subprocess.SubprocessError:
        raw_files = []

    scanned_count = 0
    changed_count = 0
    total_reps = 0
    failed = False

    for b_rel in raw_files:
        if not b_rel:
            continue
        rel = b_rel.decode("utf-8", errors="replace")
        rel_posix = rel.replace("\\", "/")
        if is_ignored_path(rel_posix, ignore_dirs, ignore_pats):
            continue
        full_path = repo_root / rel
        if not full_path.is_file() or not is_scannable(full_path):
            continue

        scanned_count += 1
        try:
            with open(full_path, "r", encoding="utf-8", errors="surrogateescape") as f:
                content = f.read()
            new_content, reps = rewrite_content(content, split_base, targets, current_version)
            if reps > 0:
                changed_count += 1
                total_reps += reps
                if not args.dry_run:
                    with open(full_path, "w", encoding="utf-8", errors="surrogateescape") as f:
                        f.write(new_content)
                if args.verbose:
                    print(f"modified: {rel} ({reps} replacements)")
        except OSError as e:
            print(f"fix-repo: ERROR write failed for {rel}: {e}", file=sys.stderr)
            failed = True

    mode_label = "dry-run" if args.dry_run else "write"
    print(f"\nscanned: {scanned_count} files")
    print(f"changed: {changed_count} files ({total_reps} replacements)")
    print(f"mode:    {mode_label}")

    return EXIT_WRITE_FAILED if failed else EXIT_OK


if __name__ == "__main__":
    sys.exit(main())
