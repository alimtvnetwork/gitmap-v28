#!/usr/bin/env python3
"""Post-rewrite audit: catches paired-literal desync (cross-platform)."""

from __future__ import annotations

import re
import sys
from pathlib import Path

PAIRED_AUDIT_LOOKAHEAD = 2


def is_test_file(path: str) -> bool:
    return path.endswith("_test.go")


def find_paired_literal_hits(file_path: Path, base: str, current: int) -> list[tuple[int, str]]:
    prev = current - 1
    if prev < 1:
        return []

    try:
        with open(file_path, "r", encoding="utf-8", errors="replace") as f:
            lines = f.readlines()
    except OSError:
        return []

    needle = f"{base}-v{current}"
    qrx = re.compile(rf'"{prev}"')
    brx = re.compile(rf"(^|[^v0-9]){prev}($|[^0-9])")

    hits: list[tuple[int, str]] = []
    for i, line in enumerate(lines):
        if needle not in line:
            continue
        end = min(len(lines), i + 1 + PAIRED_AUDIT_LOOKAHEAD)
        for j in range(i, end):
            target_line = lines[j]
            if qrx.search(target_line) or brx.search(target_line):
                hits.append((j + 1, target_line.rstrip()))
                break
    return hits


def run_paired_literal_audit(base: str, current: int, dry_run: bool, changed_files: list[str]) -> bool:
    if dry_run:
        print("audit:   skipped (dry-run)")
        return True

    total = 0
    files_with_hits = 0
    for f in changed_files:
        if not is_test_file(f):
            continue
        p = Path(f)
        hits = find_paired_literal_hits(p, base, current)
        if not hits:
            continue
        files_with_hits += 1
        for lno, content in hits:
            total += 1
            print(
                f"fix-repo: AUDIT paired-literal at {f}:{lno}: '{base}-v{current}' on/near sibling literal '{current - 1}'\n  -> line: {content}",
                file=sys.stderr,
            )

    if total == 0:
        print("audit:   no paired-literal desync detected")
        return True

    print(
        f"fix-repo: ERROR paired-literal audit failed: {total} hit(s) in {files_with_hits} file(s) (E_PAIRED_LITERAL)",
        file=sys.stderr,
    )
    return False
