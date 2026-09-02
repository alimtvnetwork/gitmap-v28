#!/usr/bin/env python3
"""coverage-report.py — Merge per-package coverage profiles and print package breakdown (cross-platform).

Usage:
  python .github/scripts/coverage-report.py [results_dir]
"""

from __future__ import annotations

import re
import shutil
import subprocess
import sys
from collections import defaultdict
from pathlib import Path


def main() -> int:
    results_dir_name = sys.argv[1] if len(sys.argv) > 1 else "results"
    results_dir = Path(f"../{results_dir_name}")
    if not results_dir.is_dir():
        results_dir = Path(results_dir_name)

    print("\n=========================================")
    print("  COVERAGE BY PACKAGE")
    print("=========================================")

    # Find coverage files
    cover_files: list[Path] = []
    if results_dir.is_dir():
        for d in results_dir.glob("test-results-*"):
            if d.is_dir():
                cover_files.extend(d.glob("coverage-*.out"))

    merged_lines: list[str] = ["mode: atomic\n"]
    for cf in cover_files:
        try:
            with open(cf, "r", encoding="utf-8", errors="replace") as f:
                for line in f:
                    if not line.startswith("mode:"):
                        merged_lines.append(line)
        except OSError:
            pass

    out_file = Path("combined-coverage.out")
    try:
        with open(out_file, "w", encoding="utf-8") as f:
            f.writelines(merged_lines)
    except OSError:
        pass

    if len(merged_lines) <= 1:
        print("No coverage data collected.")
        print("=========================================")
        return 0

    go_exe = shutil.which("go")
    if not go_exe:
        print("go toolchain not found.")
        print("=========================================")
        return 0

    try:
        raw = subprocess.check_output([go_exe, "tool", "cover", f"-func={out_file}"], text=True)
    except subprocess.SubprocessError:
        print("Failed to run go tool cover on combined profile.")
        print("=========================================")
        return 0

    pkg_pcts: dict[str, list[float]] = defaultdict(list)
    total_pct = "0.0%"

    for line in raw.splitlines():
        if line.startswith("total:"):
            parts = line.split()
            if len(parts) >= 2:
                total_pct = parts[-1]
            continue
        parts = line.split()
        if len(parts) < 3:
            continue
        fn_part = parts[0]
        pct_str = parts[-1].rstrip("%")
        try:
            pct_val = float(pct_str)
        except ValueError:
            continue
        # Strip function name
        pkg = re.sub(r"/[^/]+$", "", fn_part)
        pkg_pcts[pkg].append(pct_val)

    for pkg, vals in sorted(pkg_pcts.items()):
        avg = sum(vals) / len(vals) if vals else 0.0
        print(f"  {pkg:<55} {avg:5.1f}%")

    print(f"\n  {'Total':<55} {total_pct}")
    print("=========================================")
    return 0


if __name__ == "__main__":
    sys.exit(main())
