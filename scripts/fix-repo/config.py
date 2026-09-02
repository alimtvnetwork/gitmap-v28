#!/usr/bin/env python3
"""Config loader and path matching for fix-repo (cross-platform)."""

from __future__ import annotations

import fnmatch
import json
import sys
from pathlib import Path


def load_fixrepo_config(explicit: str | None, repo_root: Path) -> tuple[list[str], list[str]]:
    cfg_file = Path(explicit) if explicit else repo_root / "fix-repo.config.json"
    if not cfg_file.is_file():
        return [], []
    try:
        with open(cfg_file, "r", encoding="utf-8") as f:
            data = json.load(f)
        dirs = [d for d in data.get("ignoreDirs", []) if isinstance(d, str)]
        pats = [p for p in data.get("ignorePatterns", []) if isinstance(p, str)]
        return dirs, pats
    except Exception as e:
        print(f"fix-repo: ERROR reading config {cfg_file}: {e}", file=sys.stderr)
        return [], []


def is_ignored_path(rel_posix: str, ignore_dirs: list[str], ignore_pats: list[str]) -> bool:
    for d in ignore_dirs:
        clean_d = d.rstrip("/")
        if rel_posix == clean_d or rel_posix.startswith(clean_d + "/"):
            return True
    for p in ignore_pats:
        if fnmatch.fnmatch(rel_posix, p):
            return True
    return False
