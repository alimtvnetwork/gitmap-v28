#!/usr/bin/env python3
"""scripts/preflight_ci.py — Run test + lint suite locally matching CI (cross-platform).

Usage:
  python scripts/preflight_ci.py           # run both phases
  python scripts/preflight_ci.py test      # tests only
  python scripts/preflight_ci.py lint      # lint only

Exit 0 = clean, exit 1 = failures.
"""

from __future__ import annotations

import shutil
import subprocess
import sys
from pathlib import Path

GOLANGCI_LINT_VERSION = "v1.64.8"


def run_tests(gitmap_dir: Path) -> bool:
    print("=== [1/2] go test ./... (every package, no cache) ===")
    go_exe = shutil.which("go")
    if not go_exe:
        print("✗ preflight-ci: Go toolchain missing", file=sys.stderr)
        return False
    res = subprocess.run([go_exe, "test", "./...", "-count=1"], cwd=str(gitmap_dir))
    if res.returncode != 0:
        print("\n✗ preflight-ci: go test failed", file=sys.stderr)
        return False
    print("  ✓ tests passed")
    return True


def ensure_golangci_lint() -> bool:
    lint_exe = shutil.which("golangci-lint")
    if not lint_exe:
        print("✗ preflight-ci: golangci-lint not installed", file=sys.stderr)
        print(f"  Install pinned version with:\n    go install github.com/golangci/golangci-lint/cmd/golangci-lint@{GOLANGCI_LINT_VERSION}", file=sys.stderr)
        return False
    try:
        out = subprocess.check_output([lint_exe, "version", "--format", "short"], text=True, stderr=subprocess.DEVNULL).strip()
        expected = GOLANGCI_LINT_VERSION.lstrip("v")
        if out != expected:
            print("⚠ preflight-ci: golangci-lint version mismatch", file=sys.stderr)
            print(f"  installed: {out}\n  expected:  {expected}", file=sys.stderr)
    except Exception:
        pass
    return True


def run_lint(gitmap_dir: Path) -> bool:
    print("=== [2/2] golangci-lint (strict, full suite) ===")
    if not ensure_golangci_lint():
        return False
    lint_exe = shutil.which("golangci-lint") or "golangci-lint"
    cmd = [
        lint_exe,
        "run",
        "./...",
        "--timeout=5m",
        "--max-issues-per-linter=0",
        "--max-same-issues=0",
    ]
    res = subprocess.run(cmd, cwd=str(gitmap_dir))
    if res.returncode != 0:
        print("\n✗ preflight-ci: golangci-lint failed", file=sys.stderr)
        return False
    print("  ✓ lint clean")
    return True


def main() -> int:
    repo_root = Path(__file__).resolve().parent.parent
    gitmap_dir = repo_root / "gitmap"
    if not gitmap_dir.is_dir():
        print(f"✗ preflight-ci: 'gitmap/' not found at {repo_root}", file=sys.stderr)
        return 1

    phase = sys.argv[1] if len(sys.argv) > 1 else "all"
    if phase not in ("all", "test", "tests", "lint"):
        print("Usage: python preflight_ci.py [all|test|lint]", file=sys.stderr)
        return 1

    tests_ok = True
    lint_ok = True

    if phase in ("all", "test", "tests"):
        tests_ok = run_tests(gitmap_dir)
        print()

    if phase in ("all", "lint"):
        lint_ok = run_lint(gitmap_dir)

    print("\n=== preflight-ci summary ===")
    print(f"  tests: {'PASS' if tests_ok else 'FAIL'}")
    print(f"  lint:  {'PASS' if lint_ok else 'FAIL'}")

    if tests_ok and lint_ok:
        print("\n✓ preflight-ci: ready to push")
        return 0

    print("\n✗ preflight-ci: fix failures above before pushing", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
