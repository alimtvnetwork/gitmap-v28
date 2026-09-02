#!/usr/bin/env python3
"""linter-scripts/check-axios-version.py — Axios Version Safeguard (cross-platform).

Validates that Axios is pinned to an approved safe version and not using range symbols.
Blocked: 1.14.1, 0.30.4
Approved: 1.14.0, 0.30.3
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

BLOCKED_VERSIONS = {"1.14.1", "0.30.4"}
APPROVED_VERSIONS = {"1.14.0", "0.30.3"}


def main() -> int:
    if sys.platform == "win32":
        try:
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
            sys.stderr.reconfigure(encoding="utf-8", errors="replace")
        except Exception:
            pass

    pkg_file = Path("package.json")
    if not pkg_file.is_file():
        # Axios not declared if package.json does not exist
        print("⚠️  package.json not found — skipping Axios version check")
        return 0

    try:
        with open(pkg_file, "r", encoding="utf-8") as f:
            pkg = json.load(f)
    except Exception as e:
        print(f"⚠️  Could not parse package.json: {e}")
        return 0

    deps = pkg.get("dependencies", {})
    dev_deps = pkg.get("devDependencies", {})
    current = deps.get("axios") or dev_deps.get("axios")

    if not current:
        print("⚠️  Axios is not declared in package.json")
        return 0

    print(f"📦 Axios version in package.json: {current}")

    # Check: range symbols
    if any(current.startswith(prefix) for prefix in ("^", "~", ">=", "*")) or current == "latest":
        print(f"❌ FAIL: Axios version uses a range symbol or tag: {current}", file=sys.stderr)
        print('   Fix: Use an exact version like "axios": "1.14.0"', file=sys.stderr)
        return 1

    # Check: blocked versions
    if current in BLOCKED_VERSIONS:
        print(f"❌ FAIL: Axios version {current} is BLOCKED (known security vulnerability)", file=sys.stderr)
        print(f"   Approved versions: {', '.join(sorted(APPROVED_VERSIONS))}", file=sys.stderr)
        return 1

    # Check: approved versions
    if current in APPROVED_VERSIONS:
        print(f"✅ PASS: Axios {current} is an approved safe version")
        return 0

    print(f"⚠️  WARNING: Axios {current} is not in the approved list ({', '.join(sorted(APPROVED_VERSIONS))})", file=sys.stderr)
    print("   This version has not been verified. Please review spec/01-app/axios-version-control/", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
