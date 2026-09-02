#!/usr/bin/env python3
"""linter-scripts/check-forbidden-spec-paths.py — Forbidden Spec Paths Guard (cross-platform).

Fails CI on:
1. Deprecated update folders under spec/: spec/14-generic-update, spec/15-self-update-app-update
2. Transient merge-proposal.md files under spec/
3. Any uppercase-letter .md filenames under spec/ or release-artifacts/
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

FORBIDDEN_DIRS = [
    Path("spec/14-generic-update"),
    Path("spec/15-self-update-app-update"),
]


def check_uppercase_md(root: Path) -> list[Path]:
    if not root.is_dir():
        return []
    hits: list[Path] = []
    for p in root.rglob("*.md"):
        if p.is_file() and any(c.isupper() for c in p.name):
            hits.append(p)
    return hits


def main() -> int:
    if sys.platform == "win32":
        try:
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
            sys.stderr.reconfigure(encoding="utf-8", errors="replace")
        except Exception:
            pass

    exit_code = 0
    print("🔍 Checking for forbidden spec paths and uppercase .md filenames...")

    # 1. Forbidden folders
    for d in FORBIDDEN_DIRS:
        if d.exists():
            print(f"::error file={d}::Forbidden folder present: {d} (merged into spec/14-update/, must not re-appear)")
            exit_code = 1

    # 2. Forbidden files
    spec_dir = Path("spec")
    if spec_dir.is_dir():
        for p in spec_dir.rglob("*"):
            if p.is_file() and p.name.lower() == "merge-proposal.md":
                print(f"::error file={p}::Forbidden file: MERGE-PROPOSAL.md must not be committed under spec/")
                exit_code = 1

    # 3. Uppercase .md filenames
    for root_name in ("spec", "release-artifacts"):
        for hit in check_uppercase_md(Path(root_name)):
            print(f"::error file={hit}::Uppercase letters in .md filename — rename to lowercase: {hit.name}")
            exit_code = 1

    print()
    if exit_code == 0:
        print("✅ No forbidden paths or uppercase .md filenames detected.")
    else:
        print("❌ Violations detected. See errors above.", file=sys.stderr)
        print("   - Consolidated update home: spec/14-update/", file=sys.stderr)
        print("   - Markdown filenames must be all lowercase (e.g. readme.md).", file=sys.stderr)

    return exit_code


if __name__ == "__main__":
    sys.exit(main())
