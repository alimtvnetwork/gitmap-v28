#!/usr/bin/env python3
"""
42-clean-test-and-build-caches.py — Fast Cache & Test Artifact Purge Guard

Autonomously cleans:
1. Repository temporary directories and test artifacts:
   - tmp/, .tmp/, __pycache__/, *.test, coverage.out, coverage_func.txt
2. OS temporary directories:
   - Stale go-build* and gitmap* directories under $TEMP / /tmp
3. Golang build, test, and fuzz caches:
   - go clean -cache -testcache -fuzzcache
   - golangci-lint cache clean (if available)

Usage:
  python 03-ai-scripts/42-clean-test-and-build-caches.py
  python 03-ai-scripts/42-clean-test-and-build-caches.py --plan
  python 03-ai-scripts/42-clean-test-and-build-caches.py --quiet
"""

from __future__ import annotations

import argparse
from importlib import import_module
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

REPO_ROOT = Path(__file__).resolve().parent.parent

REPO_TEMP_DIRS = [
    REPO_ROOT / "tmp",
    REPO_ROOT / ".tmp",
    REPO_ROOT / ".ai-memory" / "temp",
    REPO_ROOT / ".ai-memory" / "temp-scripts",
    REPO_ROOT / ".ai-memory" / "temp-agents",
]

REPO_ARTIFACT_PATTERNS = [
    "*.test",
    "*.test.exe",
    "coverage.out",
    "coverage_func.txt",
    "*.tmp",
]

def clean_repo_temp_dirs(is_plan: bool) -> int:
    removed_count = 0
    for temp_dir in REPO_TEMP_DIRS:
        if temp_dir.exists() and temp_dir.is_dir():
            if is_plan:
                print(f"[PLAN] Would remove repo temp dir: {temp_dir}")
            else:
                try:
                    shutil.rmtree(temp_dir, ignore_errors=True)
                    removed_count += 1
                except Exception as ex:
                    print(f"  ⚠ Failed to remove {temp_dir}: {ex}")

    # Remove accidental nested .ai-memory dirs outside root
    for nested_mem in REPO_ROOT.rglob(".ai-memory"):
        if nested_mem.resolve() == (REPO_ROOT / ".ai-memory").resolve():
            continue
        if ".git" in nested_mem.parts:
            continue
        if is_plan:
            print(f"[PLAN] Would remove nested .ai-memory dir: {nested_mem}")
        else:
            try:
                shutil.rmtree(nested_mem, ignore_errors=True)
                removed_count += 1
            except Exception as ex:
                print(f"  ⚠ Failed to remove nested {nested_mem}: {ex}")

    # Remove __pycache__ trees
    for pycache in REPO_ROOT.rglob("__pycache__"):
        if ".git" in pycache.parts:
            continue
        if is_plan:
            print(f"[PLAN] Would remove pycache: {pycache}")
        else:
            shutil.rmtree(pycache, ignore_errors=True)
            removed_count += 1

    # Remove repo artifact files
    for pattern in REPO_ARTIFACT_PATTERNS:
        for file_path in REPO_ROOT.glob(pattern):
            if file_path.is_file():
                if is_plan:
                    print(f"[PLAN] Would remove artifact file: {file_path}")
                else:
                    try:
                        file_path.unlink(missing_ok=True)
                        removed_count += 1
                    except Exception as ex:
                        print(f"  ⚠ Failed to unlink {file_path}: {ex}")

    return removed_count


def clean_os_temp_dirs(is_plan: bool) -> int:
    removed_count = 0
    os_temp = Path(tempfile.gettempdir())
    if not os_temp.exists():
        return 0

    prefixes = ("go-build", "gitmap-", "gitmap_")
    try:
        for entry in os_temp.iterdir():
            if any(entry.name.startswith(p) for p in prefixes):
                if is_plan:
                    print(f"[PLAN] Would remove OS temp item: {entry}")
                else:
                    try:
                        if entry.is_dir():
                            shutil.rmtree(entry, ignore_errors=True)
                        else:
                            entry.unlink(missing_ok=True)
                        removed_count += 1
                    except Exception:
                        pass
    except Exception as ex:
        print(f"  ⚠ Could not iterate OS temp dir {os_temp}: {ex}")

    return removed_count


def clean_go_caches(is_plan: bool) -> bool:
    if is_plan:
        print("[PLAN] Would execute: go clean -cache -testcache -fuzzcache")
        return True

    cmd = ["go", "clean", "-cache", "-testcache", "-fuzzcache"]
    try:
        proc = subprocess.run(cmd, cwd=str(REPO_ROOT), capture_output=True, text=True, check=False)
        if proc.returncode != 0 and proc.stderr:
            print(f"  ⚠ go clean notice: {proc.stderr.strip()}")
    except FileNotFoundError:
        print("  ⚠ go executable not found in PATH")
        return False

    if shutil.which("golangci-lint"):
        try:
            subprocess.run(["golangci-lint", "cache", "clean"], cwd=str(REPO_ROOT), capture_output=True, check=False)
        except Exception:
            pass

    return True


def run_full_cleanup(is_plan: bool, is_quiet: bool) -> int:
    if not is_quiet:
        print("🧹 Cleaning repository temporary directories, OS temp caches, and Go caches...")

    repo_count = clean_repo_temp_dirs(is_plan)
    os_count = clean_os_temp_dirs(is_plan)
    clean_go_caches(is_plan)

    if not is_quiet:
        mode_label = "[PLAN] " if is_plan else ""
        print(f"✔ {mode_label}Cleanup complete: {repo_count} repo item(s), {os_count} OS temp item(s) purged, Go cache cleared.")

    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Clean repository temp folders, OS temp caches, and Go build/test caches")
    parser.add_argument("--plan", action="store_true", help="Preview items to be removed without deleting")
    parser.add_argument("--quiet", "-q", action="store_true", help="Suppress informational messages")
    args = parser.parse_args()

    return run_full_cleanup(args.plan, args.quiet)


if __name__ == "__main__":
    sys.exit(main())
