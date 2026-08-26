#!/usr/bin/env python3
"""
Self-test for check-file-sizes.py.

Locks the contract that guards Hard Rule #6:
  - --check exits 0 when only baseline-pinned violations exist
  - --check exits 1 when a pinned file GROWS past its baseline
  - --check exits 1 when a NEW file over cap appears
  - --strict exits 1 on any violation regardless of baseline
  - Waivers on the first 5 lines relax the cap up to ABSOLUTE_CEILING
  - Malformed waivers (over ceiling) exit 2
"""
from __future__ import annotations

import json
import subprocess
import sys
import tempfile
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[2]
SCRIPT = REPO_ROOT / "linter-scripts" / "check-file-sizes.py"


def run(args, cwd):
    # Use the COPY inside `cwd` so its REPO_ROOT resolves into the sandbox.
    local_script = Path(cwd) / "linter-scripts" / "check-file-sizes.py"
    return subprocess.run(
        [sys.executable, str(local_script), *args],
        cwd=cwd, capture_output=True, text=True,
    )


def assert_eq(label, got, want):
    if got != want:
        print(f"FAIL {label}: got {got!r} want {want!r}")
        sys.exit(1)
    print(f"OK   {label}")


def scaffold(tmp: Path):
    (tmp / "linter-scripts").mkdir()
    # Copy script so relative REPO_ROOT resolves to tmp
    src_text = SCRIPT.read_text()
    (tmp / "linter-scripts" / "check-file-sizes.py").write_text(src_text)
    (tmp / "src").mkdir()


def make_file(path: Path, lines: int, waiver: str | None = None):
    path.parent.mkdir(parents=True, exist_ok=True)
    body = []
    if waiver:
        body.append(f"// {waiver}")
    body.extend(f"// line {i}" for i in range(lines - (1 if waiver else 0)))
    path.write_text("\n".join(body) + "\n")


def test_all():
    with tempfile.TemporaryDirectory() as td:
        tmp = Path(td)
        scaffold(tmp)

        # Case A: only baseline-pinned violation
        big = tmp / "src" / "legacy.ts"
        make_file(big, 400)
        r = run(["--write-baseline"], tmp)
        assert_eq("baseline write exit", r.returncode, 0)
        r = run(["--check"], tmp)
        assert_eq("check with only pinned (exit 0)", r.returncode, 0)

        # Case B: pinned file grows
        make_file(big, 450)
        r = run(["--check"], tmp)
        assert_eq("check grown pinned (exit 1)", r.returncode, 1)
        if "GREW from 400" not in r.stdout:
            print("FAIL grown message missing")
            print(r.stdout)
            sys.exit(1)
        print("OK   grown message present")

        # Reset baseline to new size then add NEW violation
        r = run(["--write-baseline"], tmp)
        assert_eq("re-baseline exit", r.returncode, 0)
        make_file(tmp / "src" / "fresh.tsx", 250)  # tsx cap 100
        r = run(["--check"], tmp)
        assert_eq("check NEW tsx over cap (exit 1)", r.returncode, 1)
        if "NEW" not in r.stdout or "fresh.tsx" not in r.stdout:
            print("FAIL new-file message")
            print(r.stdout)
            sys.exit(1)
        print("OK   new-file message present")

        # Case C: waiver clears violation
        make_file(tmp / "src" / "fresh.tsx", 250,
                  waiver='lint-allow: file-size reason="registry table" max=300')
        r = run(["--check"], tmp)
        assert_eq("waiver clears violation (exit 0)", r.returncode, 0)

        # Case D: waiver over ceiling
        make_file(tmp / "src" / "bad.ts", 50,
                  waiver='lint-allow: file-size reason="x" max=9999')
        r = run(["--check"], tmp)
        assert_eq("waiver over ceiling (exit 2)", r.returncode, 2)

        # Cleanup bad waiver, then test --strict
        (tmp / "src" / "bad.ts").unlink()
        r = run(["--strict"], tmp)
        # legacy.ts still 450, still over 300 => strict fails even with baseline
        assert_eq("strict ignores baseline (exit 1)", r.returncode, 1)

    print("\nAll check-file-sizes self-tests passed.")


if __name__ == "__main__":
    test_all()
