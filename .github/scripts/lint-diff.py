#!/usr/bin/env python3
"""
lint-diff.py — diff golangci-lint JSON reports against a baseline.

Reads the current run's JSON report (--current) and an optional baseline
JSON (--baseline, missing/empty file is treated as "no baseline yet").
Prints a per-file summary of NEW findings (present in current, not in
baseline), FIXED findings (present in baseline, not in current), and
UNCHANGED findings (carried over). Exits 1 when there are any NEW
findings; exits 0 in every other case.

A "finding" is normalized to a tuple of (file, line, linter, message)
so trivial reordering or unrelated edits in adjacent files don't trip
the diff. Severity, source-line snippet, and column are intentionally
ignored — they're noisy across golangci-lint versions.

This script is the heart of the lint-baseline-diff CI job in ci.yml:
the job fails ONLY on NEW findings, so the lint debt is bent toward
zero without forcing every PR to fix pre-existing issues at once.
"""

from __future__ import annotations

import argparse
import json
import os
import sys
from collections import defaultdict
from typing import Iterable


Finding = tuple[str, int, str, str]  # (file, line, linter, message)


def main() -> int:
    args = parse_args()
    current = load_findings(args.current)
    baseline_present = bool(args.baseline) \
        and os.path.exists(args.baseline) \
        and os.path.getsize(args.baseline) > 0
    baseline = load_findings(args.baseline) if baseline_present else set()

    added = sorted(current - baseline)
    fixed = sorted(baseline - current)
    unchanged = sorted(current & baseline)

    print_summary(added, fixed, unchanged, baseline_present=baseline_present,
                  current_path=args.current, baseline_path=args.baseline)

    # Seeding run (no baseline yet): never gate. Surface the findings as
    # warnings so the next run has something to diff against, but exit 0
    # so the very first push to main on a brand-new repo doesn't fail.
    if not baseline_present:
        for f in added:
            file, line, linter, message = f
            print(f"::warning file={file},line={line}::"
                  f"[{linter}] {message} (seeding baseline; "
                  f"report={args.current})")
        return 0

    if added:
        # GitHub Actions error annotations — surface in the PR check UI.
        # Each annotation includes the source report path so log readers
        # can locate the raw JSON entry without grepping the whole job.
        for f in added:
            file, line, linter, message = f
            print(f"::error file={file},line={line}::"
                  f"[{linter}] {message} (NEW in {args.current})")
        return 1

    return 0


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--current", required=True,
                        help="Path to current golangci-lint JSON report")
    parser.add_argument("--baseline", default="",
                        help="Path to baseline JSON (optional)")
    return parser.parse_args()


def load_findings(path: str) -> set[Finding]:
    """Parse a golangci-lint JSON report into a set of normalized findings.

    A missing, empty, or malformed file yields an empty set — callers
    treat that as "no baseline" and skip the gate. We never want a
    transient cache miss to fail an otherwise-clean PR.
    """
    if not path or not os.path.exists(path):
        return set()
    if os.path.getsize(path) == 0:
        return set()

    try:
        with open(path, encoding="utf-8") as fh:
            data = json.load(fh)
    except (json.JSONDecodeError, OSError) as err:
        print(f"::warning::could not parse {path}: {err}", file=sys.stderr)
        return set()

    return set(extract_findings(data.get("Issues") or []))


def extract_findings(issues: Iterable[dict]) -> Iterable[Finding]:
    for issue in issues:
        pos = issue.get("Pos") or {}
        file = pos.get("Filename", "")
        line = int(pos.get("Line", 0) or 0)
        linter = issue.get("FromLinter", "")
        message = (issue.get("Text") or "").strip()
        if not file or not linter:
            continue
        yield (file, line, linter, message)


def print_summary(added: list[Finding], fixed: list[Finding],
                  unchanged: list[Finding], baseline_present: bool,
                  current_path: str = "", baseline_path: str = "") -> None:
    """Print a human-readable diff so log readers don't need to scroll
    through raw JSON. Mirrors the test-summary.sh layout for visual
    consistency across CI jobs.

    Source report paths are echoed up front so failures in CI logs
    always cite the JSON file that produced them — useful when multiple
    lint jobs (root vs. gitmap submodule) run in the same workflow.
    """
    print()
    print("=" * 72)
    print("  GOLANGCI-LINT DIFF vs LAST SUCCESSFUL MAIN")
    print("=" * 72)
    print()
    if current_path:
        print(f"  current  : {current_path}")
    if baseline_path:
        print(f"  baseline : {baseline_path}")
    if current_path or baseline_path:
        print()

    if not baseline_present:
        print("  ⚠ No baseline cached yet — this run will seed the cache.")
        print("    The diff gate is disabled until a baseline exists.")
        print()

    print(f"  + NEW       : {len(added):4d}")
    print(f"  - FIXED     : {len(fixed):4d}")
    print(f"  = UNCHANGED : {len(unchanged):4d}")
    print()

    if added:
        print_per_file("NEW findings", added, sigil="+")
    if fixed:
        print_per_file("FIXED findings", fixed, sigil="-")

    print("=" * 72)


def print_per_file(label: str, findings: list[Finding], sigil: str) -> None:
    print(f"  {label}:")
    grouped: dict[str, list[Finding]] = defaultdict(list)
    for f in findings:
        grouped[f[0]].append(f)
    for file in sorted(grouped):
        print(f"    {file}")
        for _, line, linter, message in grouped[file]:
            print(f"      {sigil} L{line} [{linter}] {message}")
    print()


if __name__ == "__main__":
    sys.exit(main())
