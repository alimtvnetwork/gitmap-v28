#!/usr/bin/env python3
"""
check-memory-mirror-drift.py — Detect drift between .lovable/memory/index.md Core
section and spec/17-consolidated-guidelines/21-lovable-folder-structure.md §X mirror.

Exit codes:
  0  No drift — all Core keywords present in mirror
  1  Drift detected — Core rules added/changed without updating mirror
  2  Structural error — file missing or section markers absent

Strategy: Extract distinctive keywords from each Core bullet (multi-word capitalized
phrases, identifier strings, numeric thresholds) and assert each appears in §X.
This is a presence check, not a verbatim diff — the mirror reformats into tables.
"""
from __future__ import annotations
import re
import sys
from pathlib import Path

MEMORY = Path(".lovable/memory/index.md")
MIRROR = Path("spec/17-consolidated-guidelines/21-lovable-folder-structure.md")
SECTION_MARKER = "## §X Project Memory"

# Distinctive tokens to search for. Each tuple is (description, keyword).
# Keywords must be specific enough that absence = real drift.
EXPECTED_TOKENS: list[tuple[str, str]] = [
    ("Code-Red zero-nesting rule", "Zero-nesting"),
    ("Function size metric", "8\u201315"),  # en-dash 8–15
    ("File size metric", "300 lines"),
    ("React component size", "100 lines"),
    ("Sync block list",            "01-app"),
    ("Sync block: app-issues",     "02-app-issues"),
    ("Sync block: consolidated",   "12-consolidated-guidelines"),
    ("PascalCase rule",            "PascalCase"),
    ("Rust exception",             "snake_case"),
    ("Primary key pattern",        "{TableName}Id"),
    ("No UUIDs rule",              "No UUIDs"),
    ("Approved inverse: Disabled", "IsDisabled"),
    ("Approved inverse: Unread",   "IsUnread"),
    ("Approved inverse: Locked",   "IsLocked"),
    ("Description column rule",    "Description TEXT NULL"),
    ("Notes/Comments rule",        "Notes"),
    ("Spec-First workflow",        "Spec-First"),
    ("Issue-First workflow",       "Issue-First"),
    ("Repo identity",              "alimtvnetwork/coding-guidelines-v24"),
    ("Axios pin: safe",            "1.14.0"),
    ("Axios pin: blocked",         "1.14.1"),
]

def fail(msg: str, code: int = 1) -> None:
    print(f"[memory-mirror-drift] FAIL: {msg}", file=sys.stderr)
    sys.exit(code)

def extract_mirror_section(text: str) -> str:
    idx = text.find(SECTION_MARKER)
    if idx == -1:
        fail(f"section marker not found: {SECTION_MARKER!r}", code=2)
    return text[idx:]

def main() -> int:
    if not MEMORY.exists():
        fail(f"missing {MEMORY}", code=2)
    if not MIRROR.exists():
        fail(f"missing {MIRROR}", code=2)

    mirror_text = MIRROR.read_text(encoding="utf-8")
    mirror_section = extract_mirror_section(mirror_text)

    missing: list[tuple[str, str]] = []
    for desc, token in EXPECTED_TOKENS:
        if token not in mirror_section:
            missing.append((desc, token))

    if missing:
        print("[memory-mirror-drift] DRIFT DETECTED", file=sys.stderr)
        print(f"[memory-mirror-drift] {len(missing)} expected token(s) missing from §X mirror:", file=sys.stderr)
        for desc, token in missing:
            print(f"  - {desc!r}: token {token!r} not found", file=sys.stderr)
        print("", file=sys.stderr)
        print("Action: update §X in spec/17-consolidated-guidelines/21-lovable-folder-structure.md", file=sys.stderr)
        print("        to mirror the latest .lovable/memory/index.md Core rules.", file=sys.stderr)
        return 1

    print(f"[memory-mirror-drift] OK — all {len(EXPECTED_TOKENS)} core tokens present in §X mirror")
    return 0

if __name__ == "__main__":
    sys.exit(main())
