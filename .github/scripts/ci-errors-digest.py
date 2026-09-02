#!/usr/bin/env python3
"""ci-errors-digest.py — Consolidate every CI failure signal into one copy-pasteable block (cross-platform).

Usage:
  python .github/scripts/ci-errors-digest.py [artifacts-root]
"""

from __future__ import annotations

import json
import os
import re
import sys
from pathlib import Path

MAX_PER_SECTION = 80
MAX_TOTAL_LINES = 2000


def main() -> int:
    root_arg = sys.argv[1] if len(sys.argv) > 1 else "results"
    root = Path(root_arg)

    lines: list[str] = [
        "===BEGIN CI ERROR DIGEST===",
        f"Repo:   {os.getenv('GITHUB_REPOSITORY', 'unknown')}",
        f"SHA:    {os.getenv('GITHUB_SHA', 'unknown')}",
        f"Run:    {os.getenv('GITHUB_SERVER_URL', 'https://github.com')}/{os.getenv('GITHUB_REPOSITORY', 'x/x')}/actions/runs/{os.getenv('GITHUB_RUN_ID', '0')}",
        f"Branch: {os.getenv('GITHUB_REF_NAME', 'unknown')}",
        "",
    ]

    total_errors = 0

    # 1. Per-suite go test failures
    test_dirs: list[Path] = []
    if root.is_dir():
        for d in root.iterdir():
            if d.is_dir() and (d.name.startswith("test-results-") or d.name.startswith("e2e-output-") or d.name == "full-suite-outputs"):
                test_dirs.append(d)

    for d in test_dirs:
        log_files = list(d.glob("test-output.txt")) + [p for p in d.glob("*.log") if p.name != "test-output.txt"]
        for log in log_files:
            if not log.is_file():
                continue
            try:
                with open(log, "r", encoding="utf-8", errors="replace") as f:
                    content = f.read()
            except OSError:
                continue

            fail_matches = re.findall(r"^--- FAIL:", content, re.MULTILINE)
            if not fail_matches:
                continue

            suite_name = d.name.replace("test-results-", "").replace("e2e-output-", "")
            lines.append(f"── TEST FAILURES — {suite_name} ({len(fail_matches)} failed) ──")

            # Extract failing tests and context
            current_test = ""
            body_lines: list[str] = []
            capturing = False
            for line in content.splitlines():
                if line.startswith("=== RUN"):
                    parts = line.split()
                    current_test = parts[2] if len(parts) >= 3 else ""
                    body_lines = []
                    capturing = True
                elif line.startswith("--- FAIL:"):
                    if capturing and current_test:
                        lines.append(f"  --- FAIL: {current_test}")
                        for bl in body_lines[:MAX_PER_SECTION]:
                            lines.append(bl)
                    capturing = False
                    current_test = ""
                    body_lines = []
                elif line.startswith("--- PASS:"):
                    capturing = False
                    current_test = ""
                    body_lines = []
                elif capturing:
                    if re.search(r"\.go:\d+:", line) or re.search(r"(expected|got:|want:|Error:|panic:|undefined|mismatch)", line):
                        body_lines.append(f"    {line}")

            lines.append("")
            total_errors += len(fail_matches)

    # 2. golangci-lint strict findings
    strict_log = root / "full-suite-outputs" / "lint-output.txt"
    if strict_log.is_file():
        try:
            with open(strict_log, "r", encoding="utf-8", errors="replace") as f:
                hits = [line.rstrip() for line in f if re.match(r"^[^\s].+:\d+:\d+:", line)]
            if hits:
                lines.append("── GOLANGCI-LINT (strict) ──")
                lines.append(f"(source: {strict_log})")
                for h in hits[:MAX_PER_SECTION]:
                    lines.append(f"  {h}")
                lines.append("")
                total_errors += len(hits)
        except OSError:
            pass

    # 3. golangci-lint baseline-diff JSON
    json_dir = root / "golangci-lint-report"
    if json_dir.is_dir():
        for jf in json_dir.glob("*.json"):
            try:
                with open(jf, "r", encoding="utf-8") as f:
                    data = json.load(f)
                issues = data.get("Issues", [])
                if issues:
                    lines.append("── GOLANGCI-LINT (baseline-diff, NEW findings) ──")
                    for iss in issues[:MAX_PER_SECTION]:
                        pos = iss.get("Pos", {})
                        fn = pos.get("Filename", "")
                        ln = pos.get("Line", "")
                        col = pos.get("Column", "")
                        linter = iss.get("FromLinter", "")
                        txt = iss.get("Text", "")
                        lines.append(f"    {fn}:{ln}:{col}: [{linter}] {txt}")
                    lines.append("")
                    total_errors += len(issues)
            except Exception:
                pass

    # 4. Actionable lint suggestions
    sugg_dir = root / "lint-suggestions"
    if sugg_dir.is_dir():
        for sf in list(sugg_dir.glob("*.txt")) + list(sugg_dir.glob("*.md")):
            try:
                with open(sf, "r", encoding="utf-8", errors="replace") as f:
                    content = f.read().splitlines()
                if content:
                    lines.append(f"── LINT SUGGESTIONS ({sf.name}) ──")
                    for cl in content[:MAX_PER_SECTION]:
                        lines.append(f"  {cl}")
                    lines.append("")
            except OSError:
                pass

    if total_errors == 0:
        lines.append("(no errors detected across collected artifacts)")

    lines.append("===END CI ERROR DIGEST===")

    if len(lines) > MAX_TOTAL_LINES:
        truncated = lines[:MAX_TOTAL_LINES - 2]
        print("\n".join(truncated))
        print(f"... [digest truncated at {MAX_TOTAL_LINES} lines — download artifacts for full output] ...")
        print("===END CI ERROR DIGEST===")
    else:
        print("\n".join(lines))

    return 0


if __name__ == "__main__":
    sys.exit(main())
