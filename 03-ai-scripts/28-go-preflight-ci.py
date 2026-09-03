#!/usr/bin/env python3
"""
28-go-preflight-ci.py — Cross-platform Go test and linter preflight runner.

Executes local test suite and golangci-lint before submitting commits or PRs.

Usage:
  python 03-ai-scripts/28-go-preflight-ci.py          # run both tests and lint
  python 03-ai-scripts/28-go-preflight-ci.py test     # run tests only
  python 03-ai-scripts/28-go-preflight-ci.py lint     # run linters only

Exit codes:
  0 = all checks passed, 1 = test or lint failures.
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

ExitCodeType = engine.ExitCodeType


def find_go_modules(repo_root: Path) -> list[Path]:
    """Finds all directories containing a go.mod file."""
    mod_files = list(repo_root.rglob("go.mod"))
    modules = []
    for mf in mod_files:
        if "node_modules" not in str(mf) and ".git" not in str(mf):
            modules.append(mf.parent)

    return modules


def run_tests(mod_dir: Path) -> bool:
    print(f"=== [Go Test] {mod_dir} (go test ./... -count=1) ===")
    go_exe = shutil.which("go")
    if not go_exe:
        print("✗ Go toolchain missing from PATH", file=sys.stderr)
        return False

    res = subprocess.run([go_exe, "test", "./...", "-count=1"], cwd=str(mod_dir))
    if res.returncode != 0:
        print(f"✗ Tests failed in {mod_dir}", file=sys.stderr)
        return False

    print(f"✓ All tests passed in {mod_dir}")
    return True


def run_lint(mod_dir: Path) -> bool:
    print(f"=== [Go Lint] {mod_dir} (golangci-lint run) ===")
    lint_exe = shutil.which("golangci-lint")
    if not lint_exe:
        print("⚠ golangci-lint not installed — skipping Go linter check", file=sys.stderr)
        print("  Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest", file=sys.stderr)
        return True

    res = subprocess.run([lint_exe, "run", "./..."], cwd=str(mod_dir))
    if res.returncode != 0:
        print(f"✗ golangci-lint reported issues in {mod_dir}", file=sys.stderr)
        return False

    print(f"✓ golangci-lint clean in {mod_dir}")
    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="Go Preflight CI Runner")
    parser.add_argument("phase", nargs="?", choices=["all", "test", "lint"], default="all", help="Phase to execute")
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent
    modules = find_go_modules(repo_root)

    if not modules:
        print("✓ No Go modules found in repository.")
        return int(ExitCodeType.SUCCESS.value)

    has_failure = False
    for mod in modules:
        if args.phase in ("all", "test"):
            is_test_ok = run_tests(mod)
            if not is_test_ok:
                has_failure = True

        if args.phase in ("all", "lint"):
            is_lint_ok = run_lint(mod)
            if not is_lint_ok:
                has_failure = True

    if has_failure:
        return int(ExitCodeType.VIOLATIONS_FOUND.value)

    print("\n✓ All Go Preflight checks passed successfully.")
    return int(ExitCodeType.SUCCESS.value)


if __name__ == "__main__":
    sys.exit(main())
