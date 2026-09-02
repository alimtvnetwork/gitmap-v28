#!/usr/bin/env python3
"""Apply and verify visibility change via gh / glab (cross-platform)."""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path

# Support relative import if run from package or standalone
try:
    from .provider import get_current_visibility
except ImportError:
    from provider import get_current_visibility


def apply_visibility(provider: str, slug: str, target: str) -> bool:
    try:
        if provider == "github":
            res = subprocess.run(
                [
                    "gh",
                    "repo",
                    "edit",
                    slug,
                    "--visibility",
                    target,
                    "--accept-visibility-change-consequences",
                ]
            )
            return res.returncode == 0
        elif provider == "gitlab":
            res = subprocess.run(["glab", "repo", "edit", slug, "--visibility", target])
            return res.returncode == 0
    except OSError:
        return False
    return False


def visibility_matches(provider: str, slug: str, target: str) -> bool:
    actual = get_current_visibility(provider, slug)
    return actual == target


def confirm_public_change(slug: str, provider: str) -> bool:
    if not sys.stdin.isatty():
        return False
    print(f"\n⚠  About to make {slug} PUBLIC on {provider}.")
    print("   Type 'yes' to continue, anything else aborts:")
    try:
        answer = input().strip().lower()
        return answer == "yes"
    except (EOFError, KeyboardInterrupt):
        return False
