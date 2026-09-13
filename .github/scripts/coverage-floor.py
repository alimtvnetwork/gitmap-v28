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


def parse_floor_line(line: str) -> tuple[str, float] | None:
    """Extracts package import path and floor percentage from a line."""
    clean = line.strip()
    if not clean or clean.startswith("#"):
        return None

    parts = clean.split()
    if len(parts) >= 2:
        try:
            return parts[0], float(parts[1])
        except ValueError:
            return None

    return None


def load_floors() -> dict[str, float]:
    """Loads configured package coverage floors from .github/coverage.floor."""
    if not FLOORS_FILE.is_file():
        return {}

    floors: dict[str, float] = {}
    content = FLOORS_FILE.read_text(encoding="utf-8")
    for line in content.splitlines():
        parsed = parse_floor_line(line)
        if parsed is not None:
            floors[parsed[0]] = parsed[1]

    return floors


def parse_coverage_line(line: str) -> tuple[str, float] | None:
    """Parses a single go tool cover -func line into package path and percentage."""
    if ".go:" not in line:
        return None

    parts = line.split()
    if len(parts) < 3:
        return None

    file_part = parts[0]
    pct_part = parts[-1].rstrip("%")
    try:
        pct = float(pct_part)
    except ValueError:
        return None

    pkg = re.sub(r"/[^/]+\.go:\d+:?$", "", file_part)

    return pkg, pct


def aggregate_coverage(out: str) -> tuple[dict[str, float], dict[str, int]]:
    """Aggregates percentage sums and function counts per package."""
    pkg_totals: dict[str, float] = defaultdict(float)
    pkg_counts: dict[str, int] = defaultdict(int)

    for line in out.splitlines():
        res = parse_coverage_line(line)
        if res is not None:
            pkg, pct = res
            pkg_totals[pkg] += pct
            pkg_counts[pkg] += 1

    return pkg_totals, pkg_counts


def resolve_pkg_stats(
    pkg: str, pkg_totals: dict[str, float], pkg_counts: dict[str, int]
) -> tuple[float, int]:
    """Looks up total percent and statement count for a package, supporting module path aliases."""
    if pkg in pkg_counts:
        return pkg_totals[pkg], pkg_counts[pkg]

    alt_cli = pkg.replace("/gitmap-v28/gitmap/", "/gitmap-v28/cli/")
    if alt_cli in pkg_counts:
        return pkg_totals[alt_cli], pkg_counts[alt_cli]

    alt_gitmap = pkg.replace("/gitmap-v28/cli/", "/gitmap-v28/gitmap/")
    if alt_gitmap in pkg_counts:
        return pkg_totals[alt_gitmap], pkg_counts[alt_gitmap]

    return 0.0, 0


def check_coverage_floors(
    pkg_totals: dict[str, float],
    pkg_counts: dict[str, int],
    floors: dict[str, float],
) -> bool:
    """Validates packages against configured floors and prints violations."""
    is_failed = False
    for pkg, floor in floors.items():
        tot, count = resolve_pkg_stats(pkg, pkg_totals, pkg_counts)
        avg = tot / count if count > 0 else 0.0
        if avg < floor:
            print(f"coverage-floor: {pkg} below floor (avg={avg:.1f}%, floor={floor:.1f}%)", file=sys.stderr)
            is_failed = True

    return is_failed


def run_cover_tool(go_exe: str, cover_file: Path) -> str | None:
    """Runs go tool cover -func inside the gitmap module directory."""
    go_dir = Path("cli").resolve() if (Path("cli") / "go.mod").is_file() else Path.cwd()
    try:
        return subprocess.check_output(
            [go_exe, "tool", "cover", f"-func={cover_file.resolve()}"],
            cwd=go_dir,
            text=True,
        )
    except subprocess.SubprocessError as err:
        print(f"coverage-floor: go tool cover failed: {err}", file=sys.stderr)

        return None


def validate_input_args() -> Path | None:
    """Validates command-line arguments and returns the cover profile path."""
    if len(sys.argv) < 2:
        print("usage: coverage-floor.py coverage.out", file=sys.stderr)

        return None

    path = Path(sys.argv[1])
    if not path.is_file() or path.stat().st_size == 0:
        print(f"coverage-floor: empty or missing coverage profile at {path} — skipping")

        return None

    return path


def main() -> int:
    cover_file = validate_input_args()
    if cover_file is None:
        return 0 if len(sys.argv) >= 2 else 1

    go_exe = shutil.which("go")
    if not go_exe:
        print("coverage-floor: go toolchain not found on PATH", file=sys.stderr)

        return 0

    out = run_cover_tool(go_exe, cover_file)
    if out is None:
        return 1

    pkg_totals, pkg_counts = aggregate_coverage(out)
    floors = load_floors()
    has_failure = check_coverage_floors(pkg_totals, pkg_counts, floors)

    return 1 if has_failure else 0


if __name__ == "__main__":
    sys.exit(main())
