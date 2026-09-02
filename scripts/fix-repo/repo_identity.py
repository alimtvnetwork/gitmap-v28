#!/usr/bin/env python3
"""Repo-identity helpers for fix-repo (cross-platform)."""

from __future__ import annotations

import re
import subprocess
from pathlib import Path


def get_repo_root() -> Path | None:
    try:
        out = subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"],
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        return Path(out) if out else None
    except subprocess.SubprocessError:
        return None


def get_remote_url(repo_root: Path | None = None) -> str:
    cwd = str(repo_root) if repo_root else None
    try:
        out = subprocess.check_output(
            ["git", "remote", "get-url", "origin"],
            cwd=cwd,
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        if out:
            return out
    except subprocess.SubprocessError:
        pass
    try:
        out = subprocess.check_output(
            ["git", "remote", "-v"],
            cwd=cwd,
            stderr=subprocess.DEVNULL,
            text=True,
        )
        for line in out.splitlines():
            parts = line.split()
            if len(parts) >= 2 and "(fetch)" in line:
                return parts[1]
    except subprocess.SubprocessError:
        pass
    return ""


def parse_remote_url(url: str) -> tuple[str, str, str] | None:
    trimmed = url.rstrip("/")
    if trimmed.endswith(".git"):
        trimmed = trimmed[:-4]
    patterns = [
        r"^https?://([^/:]+)(?::\d+)?/([^/]+)/([^/]+)$",
        r"^git@([^:]+):([^/]+)/([^/]+)$",
        r"^ssh://git@([^/:]+)(?::\d+)?/([^/]+)/([^/]+)$",
    ]
    for pat in patterns:
        m = re.match(pat, trimmed)
        if m:
            return m.group(1), m.group(2), m.group(3)
    return None


def split_repo_version(repo_name: str) -> tuple[str, int] | None:
    m = re.match(r"^(.+)-v(\d+)$", repo_name)
    if m:
        return m.group(1), int(m.group(2))
    return None
