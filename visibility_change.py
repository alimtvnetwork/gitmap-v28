#!/usr/bin/env python3
"""visibility_change.py — toggle/set GitHub/GitLab repo visibility (cross-platform).

Usage:
  python visibility_change.py                   # toggle current visibility
  python visibility_change.py --visible pub     # force public
  python visibility_change.py --visible pri     # force private
  python visibility_change.py --yes             # skip private->public prompt
  python visibility_change.py --dry-run         # preview, no API call
"""

from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys

EXIT_OK = 0
EXIT_NOT_A_REPO = 2
EXIT_NO_ORIGIN = 3
EXIT_BAD_PROVIDER = 4
EXIT_AUTH_FAILED = 5
EXIT_BAD_FLAG = 6
EXIT_CONFIRM_REQ = 7
EXIT_VERIFY_FAILED = 8


def get_origin_url() -> str:
    try:
        out = subprocess.check_output(
            ["git", "remote", "get-url", "origin"], stderr=subprocess.DEVNULL, text=True
        )
        return out.strip()
    except subprocess.SubprocessError:
        return ""


def resolve_provider(url: str) -> str:
    if not url:
        return ""
    m = re.match(r"^(https?://|ssh://[^@]+@|[^@]+@)([^/:]+)", url)
    if not m:
        return ""
    host = m.group(2).lower()
    if host in ("github.com", "ssh.github.com"):
        return "github"
    if host == "gitlab.com":
        return "gitlab"
    allowed_gl = os.getenv("VISIBILITY_GITLAB_HOSTS", "")
    for h in allowed_gl.split(","):
        if h.strip() and h.strip().lower() == host:
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
    except (subprocess.SubprocessError, json.JSONDecodeError):
        pass
    return ""


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


def main() -> int:
    parser = argparse.ArgumentParser(description="Toggle or set GitHub/GitLab repo visibility.")
    parser.add_argument("--visible", choices=["pub", "public", "pri", "private"], help="Target visibility")
    parser.add_argument("--yes", "-y", action="store_true", help="Skip confirmation when making public")
    parser.add_argument("--dry-run", action="store_true", help="Preview without making API calls")
    args = parser.parse_args()

    # Check git repo root
    try:
        subprocess.check_call(
            ["git", "rev-parse", "--show-toplevel"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
    except subprocess.SubprocessError:
        print("visibility-change: ERROR not a git repository", file=sys.stderr)
        return EXIT_NOT_A_REPO

    url = get_origin_url()
    if not url:
        print("visibility-change: ERROR no origin remote", file=sys.stderr)
        return EXIT_NO_ORIGIN

    provider = resolve_provider(url)
    if not provider:
        print(f"visibility-change: ERROR unsupported host in '{url}'", file=sys.stderr)
        return EXIT_BAD_PROVIDER

    slug = resolve_owner_repo(url)
    if not slug:
        print(f"visibility-change: ERROR cannot parse owner/repo from '{url}'", file=sys.stderr)
        return EXIT_BAD_PROVIDER

    cli = "gh" if provider == "github" else "glab"
    if not shutil.which(cli):
        print(f"visibility-change: ERROR '{cli}' not found on PATH", file=sys.stderr)
        return EXIT_AUTH_FAILED

    current = get_current_visibility(provider, slug)
    if not current:
        print("visibility-change: ERROR cannot read current visibility (auth?)", file=sys.stderr)
        return EXIT_AUTH_FAILED

    forced = ""
    if args.visible in ("pub", "public"):
        forced = "public"
    elif args.visible in ("pri", "private"):
        forced = "private"

    target = forced if forced else ("private" if current == "public" else "public")

    if current == target:
        print(f"visibility: already {current} ({provider})")
        return EXIT_OK

    if target == "public" and current == "private" and not args.yes:
        if not confirm_public_change(slug, provider):
            print("visibility-change: ERROR confirmation required (pass --yes for non-interactive)", file=sys.stderr)
            return EXIT_CONFIRM_REQ

    if args.dry_run:
        print(f"[dry-run] visibility: {current} → {target} ({provider})")
        return EXIT_OK

    if not apply_visibility(provider, slug, target):
        print("visibility-change: ERROR apply failed", file=sys.stderr)
        return EXIT_AUTH_FAILED

    actual = get_current_visibility(provider, slug)
    if actual != target:
        print("visibility-change: ERROR verification failed (visibility did not change)", file=sys.stderr)
        return EXIT_VERIFY_FAILED

    print(f"visibility: {current} → {target} ({provider})")
    return EXIT_OK


if __name__ == "__main__":
    sys.exit(main())
