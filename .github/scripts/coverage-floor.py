#!/usr/bin/env python3
"""coverage-floor.py — Coverage-floor enforcer (cross-platform).

Usage:
  python .github/scripts/coverage-floor.py <coverage.out>
"""

from __future__ import annotations

import re
import shutil
import subprocess
import sys
from collections import defaultdict
from pathlib import Path

FLOORS_FILE = Path(".github/coverage.floor")
FLOOR_DEFAULT = 70.0


def load_floors() -> dict[str, float]:
    floors = {}
    if FLOORS_FILE.is_file():
        try:
            with open(FLOORS_FILE, "r", encoding="utf-8") as f:
                for line in f:
                    line = line.strip()
                    if not line or line.startswith("#"):
                        continue
                    parts = line.split()
                    if len(parts) >= 2:
                        try:
                            floors[parts[0]] = float(parts[1])
                        except ValueError:
                            pass
        except OSError:
            pass
    return floors


def main() -> int:
    if len(sys.argv) < 2:
        print("usage: coverage-floor.py coverage.out", file=sys.stderr)
        return 1

    cover_file = Path(sys.argv[1])
    if not cover_file.is_file() or cover_file.stat().st_size == 0:
        print(f"coverage-floor: empty or missing coverage profile at {cover_file} — skipping")
        return 0

    go_exe = shutil.which("go")
    if not go_exe:
        print("coverage-floor: go toolchain not found on PATH", file=sys.stderr)
        return 0

    try:
        out = subprocess.check_output([go_exe, "tool", "cover", f"-func={cover_file}"], text=True)
    except subprocess.SubprocessError as e:
        print(f"coverage-floor: go tool cover failed: {e}", file=sys.stderr)
        return 1

    # Aggregate per package
    pkg_totals: dict[str, float] = defaultdict(float)
    pkg_counts: dict[str, int] = defaultdict(int)

    for line in out.splitlines():
        if ".go:" not in line:
            continue
        parts = line.split()
        if len(parts) < 3:
            continue
        file_part = parts[0]
        pct_part = parts[-1].rstrip("%")
        try:
            pct = float(pct_part)
        except ValueError:
            continue
        # Strip filename and line number to get package path
        pkg = re.sub(r"/[^/]+\.go:\d+$", "", file_part)
        pkg_totals[pkg] += pct
        pkg_counts[pkg] += 1

    floors = load_floors()
    failed = False

    for pkg, total in pkg_totals.items():
        count = pkg_counts[pkg]
        avg = total / count if count > 0 else 0.0
        floor = floors.get(pkg, FLOOR_DEFAULT)
        if avg < floor:
            print(f"coverage-floor: {pkg} below floor (avg={avg:.1f}%, floor={floor:.1f}%)", file=sys.stderr)
            failed = True

    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main())
