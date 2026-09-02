#!/usr/bin/env python3
"""test-summary.py — Aggregates test results from CI matrix jobs (cross-platform).

Usage:
  python .github/scripts/test-summary.py <results-dir>
"""

from __future__ import annotations

import re
import sys
from pathlib import Path


def main() -> int:
    if len(sys.argv) < 2:
        print("Usage: test-summary.py <results-dir>", file=sys.stderr)
        return 1

    results_dir = Path(sys.argv[1])

    print("=========================================")
    print("  ALL TEST RESULTS")
    print("=========================================")

    overall = 0
    all_failures: list[str] = []

    if results_dir.is_dir():
        for d in sorted(results_dir.glob("test-results-*")):
            if not d.is_dir():
                continue
            suite = d.name.replace("test-results-", "")
            test_out = d / "test-output.txt"
            if not test_out.is_file():
                continue

            try:
                with open(test_out, "r", encoding="utf-8", errors="replace") as f:
                    content = f.read()
            except OSError:
                continue

            pass_count = len(re.findall(r"^--- PASS:", content, re.MULTILINE))
            fail_count = len(re.findall(r"^--- FAIL:", content, re.MULTILINE))

            if fail_count > 0:
                print(f"\n❌ {suite}: {fail_count} failed, {pass_count} passed")
                overall = 1

                suite_failures: list[str] = []
                # Extract failing test details
                failing_names = re.findall(r"^--- FAIL:\s+([^\s(]+)", content, re.MULTILINE)
                for tname in failing_names:
                    suite_failures.append(f"    --- FAIL: {tname}")
                    # Look for assertion reason in the test block
                    pattern = re.compile(rf"=== RUN\s+{re.escape(tname)}(.*?)(?:--- FAIL:\s+{re.escape(tname)}|$)", re.DOTALL)
                    m = pattern.search(content)
                    if m:
                        block = m.group(1)
                        reasons = []
                        for bline in block.splitlines():
                            if re.search(r"\.go:\d+:", bline) or (re.search(r"(expected|got|Error|FAIL|panic|undefined|mismatch)", bline) and not bline.startswith("=== RUN")):
                                reasons.append(f"        {bline.strip()}")
                        for r in reasons[:10]:
                            suite_failures.append(r)

                all_failures.append(
                    f"\n-----------------------------------------\n"
                    f"  Suite: {suite} ({fail_count} failed)\n"
                    f"-----------------------------------------\n" + "\n".join(suite_failures)
                )
            else:
                print(f"✅ {suite}: {pass_count} passed")

    print("\n=========================================")

    if overall != 0:
        print("\n=========================================")
        print("  FAILURE REPORT (copy-paste ready)")
        print("=========================================")
        print("\n".join(all_failures))
        print("=========================================\n")
        print("::error::Some test suites failed — see failure report above.", file=sys.stderr)
        return 1

    print("All test suites passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
