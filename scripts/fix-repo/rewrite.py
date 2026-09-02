#!/usr/bin/env python3
"""Rewrite engine for fix-repo (cross-platform)."""

from __future__ import annotations

import re
from pathlib import Path


def get_target_versions(current: int, span: int) -> list[int]:
    start = max(1, current - span)
    return list(range(start, current))


def count_token_occurrences(content: str, base: str, n: int) -> int:
    token = f"{base}-v{n}"
    pattern = re.escape(token) + r"(?!\d)"
    return len(re.findall(pattern, content))


def rewrite_file(file_path: Path, base: str, current: int, dry_run: bool, targets: list[int]) -> int:
    try:
        with open(file_path, "r", encoding="utf-8", errors="surrogateescape") as f:
            content = f.read()
    except OSError:
        return 0

    new_content = content
    total_reps = 0
    for n in targets:
        token = f"{base}-v{n}"
        replacement = f"{base}-v{current}"
        pattern = re.escape(token) + r"(?!\d)"
        matches = len(re.findall(pattern, new_content))
        if matches > 0:
            new_content = re.sub(pattern, replacement, new_content)
            total_reps += matches

    if total_reps > 0 and not dry_run:
        try:
            with open(file_path, "w", encoding="utf-8", errors="surrogateescape") as f:
                f.write(new_content)
        except OSError:
            return -1

    return total_reps
