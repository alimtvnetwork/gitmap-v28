#!/usr/bin/env python3
"""Provider and auth helpers for visibility-change (cross-platform)."""

from __future__ import annotations

import json
import os
import re
import shutil
import subprocess
import sys


def get_origin_url() -> str:
    try:
        out = subprocess.check_output(["git", "remote", "get-url", "origin"], stderr=subprocess.DEVNULL, text=True)
        return out.strip()
    except subprocess.SubprocessError:
        return ""


def host_is_allowlisted_gitlab(host: str) -> bool:
    allow = os.getenv("VISIBILITY_GITLAB_HOSTS", "")
    for h in allow.split(","):
        clean = h.strip()
        if clean and clean.lower() == host.lower():
            return True
    return False


def resolve_provider(url: str) -> str:
    if not url:
        return ""
    m = re.match(r"^(https?://|ssh://[^@]+@|[^@]+@)([^/:]+)", url)
    if not m:
        return ""
    host = m.group(2).lower()
    if host in ("github.com", "ssh.github.com"):
        return "github"
    if host == "gitlab.com" or host_is_allowlisted_gitlab(host):
        return "gitlab"
    return ""


def resolve_owner_repo(url: str) -> str:
    trimmed = url.rstrip("/")
    if trimmed.endswith(".git"):
        trimmed = trimmed[:-4]
    patterns = [
        r"^https?://[^/]+/([^/]+)/([^/]+)$",
        r"^[^@]+@[^:]+:([^/]+)/([^/]+)$",
        r"^ssh://[^@]+@[^/]+/([^/]+)/([^/]+)$",
    ]
    for pat in patterns:
        m = re.match(pat, trimmed)
        if m:
            return f"{m.group(1)}/{m.group(2)}"
    return ""


def is_cli_available(cli: str) -> bool:
    return bool(shutil.which(cli))


def get_current_visibility(provider: str, slug: str) -> str:
    try:
        if provider == "github":
            out = subprocess.check_output(
                ["gh", "repo", "view", slug, "--json", "visibility", "-q", ".visibility"],
                stderr=subprocess.DEVNULL,
                text=True,
            )
            return out.strip().lower()
        elif provider == "gitlab":
            out = subprocess.check_output(
                ["glab", "repo", "view", slug, "-F", "json"],
                stderr=subprocess.DEVNULL,
                text=True,
            )
            data = json.loads(out)
            return str(data.get("visibility", "")).lower()
    except Exception:
        pass
    return ""
