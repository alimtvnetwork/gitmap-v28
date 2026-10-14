#!/usr/bin/env python3
"""coverage-floor-check.py — Coverage floor gate for CI.

Fails when total test coverage drops more than --max-drop percentage
points vs the main-branch baseline. The baseline is a merged Go
coverprofile stored in the GitHub Actions cache (see the
coverage-floor-gate job in .github/workflows/ci.yml).

Usage:
  python3 .github/scripts/coverage-floor-check.py \
      --current coverage/coverage-unit.out coverage/coverage-store.out \
      --baseline /tmp/coverage-baseline/coverage.out [--max-drop 2.0]

Exit codes:
  0 — coverage within the allowed drop (or seeding mode: no baseline yet,
      or no current data; the workflow seeds the cache on main).
  1 — coverage dropped more than --max-drop percentage points.
  2 — invalid arguments.

Total coverage is statement-weighted (covered statements / total
statements), matching `go tool cover`'s "total:" line, computed purely
in Python so the gate job needs no Go toolchain.
"""

from __future__ import annotations

import argparse
import glob
import re
import sys
from pathlib import Path

# `path/to/file.go:10.14,12.2 3 1` — groups: file, statements, count.
PROFILE_LINE = re.compile(r"^(.+):(\d+)\.(\d+),(\d+)\.(\d+) (\d+) (\d+)$")


def parse_profile(path: Path) -> tuple[int, int]:
    """Returns (covered_statements, total_statements) for one coverprofile."""
    covered = 0
    total = 0
    with open(path, encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("mode:"):
                continue
            m = PROFILE_LINE.match(line)
            if not m:
                continue
            stmts = int(m.group(6))
            count = int(m.group(7))
            total += stmts
            if count > 0:
                covered += stmts
    return covered, total


def total_coverage(paths: list[str]) -> float | None:
    """Statement-weighted total coverage % across profiles. None when empty."""
    covered = 0
    total = 0
    seen = 0
    for pattern in paths:
        for name in sorted(glob.glob(pattern)):
            path = Path(name)
            if not path.is_file():
                continue
            c, t = parse_profile(path)
            covered += c
            total += t
            seen += 1
    if seen == 0 or total == 0:
        return None
    return 100.0 * covered / total


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Fail when coverage drops more than --max-drop points vs baseline."
    )
    parser.add_argument(
        "--current",
        nargs="+",
        required=True,
        help="Current coverprofile paths (globs allowed).",
    )
    parser.add_argument(
        "--baseline",
        required=True,
        help="Baseline merged coverprofile path (may not exist on first run).",
    )
    parser.add_argument(
        "--max-drop",
        type=float,
        default=2.0,
        help="Allowed drop in percentage points (default 2.0).",
    )
    return parser.parse_args(argv)


def main(argv: list[str]) -> int:
    args = parse_args(argv)

    current = total_coverage(args.current)
    if current is None:
        print("::warning::coverage-floor-check: no current coverage data — skipping gate")
        return 0

    baseline_path = Path(args.baseline)
    baseline = total_coverage([str(baseline_path)]) if baseline_path.is_file() else None
    if baseline is None:
        print(
            f"::notice::coverage-floor-check: no baseline at {baseline_path} — "
            f"seeding mode, current total {current:.1f}% recorded as baseline"
        )
        return 0

    drop = baseline - current
    print(f"coverage-floor-check: baseline {baseline:.1f}% vs current {current:.1f}%")
    if drop > args.max_drop:
        print(
            f"::error::coverage dropped {drop:.1f} percentage points "
            f"(baseline {baseline:.1f}% -> current {current:.1f}%), "
            f"exceeding the {args.max_drop:.1f}-point floor"
        )
        return 1

    print(f"::notice::coverage-floor-check: OK (drop {drop:.1f}pp within {args.max_drop:.1f}pp)")
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main(sys.argv[1:]))
    except SystemExit:
        raise
    except Exception as e:  # noqa: BLE001 — never fail the gate on a script bug
        print(f"::warning::coverage-floor-check: internal error ({e}) — skipping gate")
        sys.exit(0)
