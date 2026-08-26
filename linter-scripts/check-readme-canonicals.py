#!/usr/bin/env python3
"""
check-readme-canonicals.py — Verify the root readme.md uses the canonical
repo slug and CDN domain in every badge URL, install one-liner, and workflow
link. Catches stale references after a major-version rename or domain change.

Sources of truth (override with --owner / --slug / --cdn or env vars
README_CANON_OWNER / README_CANON_SLUG / README_CANON_CDN):
  owner = alimtvnetwork
  slug  = coding-guidelines-v24
  cdn   = cdn.riseup.asia

Exit codes:
  0  PASS — every reference matches canonical values
  1  FAIL — at least one stale reference found
  2  ERROR — readme.md not found / unreadable
"""
from __future__ import annotations

import argparse
import os
import re
import sys
from pathlib import Path

EXIT_PASS = 0
EXIT_FAIL = 1
EXIT_ERROR = 2

DEFAULT_OWNER = "alimtvnetwork"
DEFAULT_SLUG = "coding-guidelines-v24"
DEFAULT_CDN = "cdn.riseup.asia"

# Bad CDN variants we explicitly reject.
BAD_CDN_PATTERNS = [
    re.compile(r"cdn\.riseup-asia\.com", re.IGNORECASE),
]

# GitHub host prefixes that should always carry the canonical owner/slug pair.
GITHUB_HOST_PATTERNS = [
    re.compile(r"https?://github\.com/(?P<owner>[^/\s\"'<>]+)/(?P<repo>[^/\s\"'<>#)?]+)"),
    re.compile(r"https?://raw\.githubusercontent\.com/(?P<owner>[^/\s\"'<>]+)/(?P<repo>[^/\s\"'<>]+)"),
    re.compile(r"https?://img\.shields\.io/github/[^/\s]+/(?P<owner>[^/\s\"'<>]+)/(?P<repo>[^/\s\"'<>?]+)"),
]


def load_canonicals(args: argparse.Namespace) -> tuple[str, str, str]:
    owner = args.owner or os.environ.get("README_CANON_OWNER") or DEFAULT_OWNER
    slug = args.slug or os.environ.get("README_CANON_SLUG") or DEFAULT_SLUG
    cdn = args.cdn or os.environ.get("README_CANON_CDN") or DEFAULT_CDN
    return owner, slug, cdn


def find_github_violations(body: str, owner: str, slug: str) -> list[str]:
    violations: list[str] = []
    for line_no, line in enumerate(body.splitlines(), 1):
        for pattern in GITHUB_HOST_PATTERNS:
            for match in pattern.finditer(line):
                found_owner = match.group("owner")
                found_repo = match.group("repo")
                owner_ok = found_owner == owner
                repo_ok = found_repo == slug
                if owner_ok and repo_ok:
                    continue
                violations.append(
                    f"  L{line_no}: {found_owner}/{found_repo}  (expected {owner}/{slug})"
                )
    return violations


def find_cdn_violations(body: str, cdn: str) -> list[str]:
    violations: list[str] = []
    for line_no, line in enumerate(body.splitlines(), 1):
        for bad in BAD_CDN_PATTERNS:
            if bad.search(line):
                violations.append(f"  L{line_no}: stale CDN host found  (expected {cdn})")
    return violations


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Verify root readme.md canonicals.")
    parser.add_argument("--readme", default="readme.md", help="Path to root README (default: readme.md)")
    parser.add_argument("--owner", help=f"Canonical GitHub owner (default: {DEFAULT_OWNER})")
    parser.add_argument("--slug", help=f"Canonical repo slug (default: {DEFAULT_SLUG})")
    parser.add_argument("--cdn", help=f"Canonical CDN host (default: {DEFAULT_CDN})")
    return parser.parse_args()


def main() -> int:
    import sys
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    args = parse_args()
    owner, slug, cdn = load_canonicals(args)
    readme = Path(args.readme)
    if not readme.is_file():
        sys.stderr.write(f"ERROR: readme not found at {readme}\n")
        return EXIT_ERROR

    body = readme.read_text(encoding="utf-8")
    gh_violations = find_github_violations(body, owner, slug)
    cdn_violations = find_cdn_violations(body, cdn)

    print(f"🔎 Canonical owner: {owner}")
    print(f"🔎 Canonical slug : {slug}")
    print(f"🔎 Canonical CDN  : {cdn}")
    print("")

    has_failures = bool(gh_violations) or bool(cdn_violations)
    if gh_violations:
        print(f"❌ GitHub owner/slug mismatches ({len(gh_violations)}):")
        for v in gh_violations:
            print(v)
        print("")
    else:
        print("✅ GitHub owner/slug: all references canonical.")

    if cdn_violations:
        print(f"❌ CDN host mismatches ({len(cdn_violations)}):")
        for v in cdn_violations:
            print(v)
        print("")
    else:
        print("✅ CDN host: no legacy references.")

    if has_failures:
        print("")
        print("Fix hint: update stale references in readme.md, or override the canonical")
        print("values via --owner/--slug/--cdn flags or README_CANON_* env vars.")
        return EXIT_FAIL

    print("")
    print("🎉 readme.md canonicals are clean.")
    return EXIT_PASS


if __name__ == "__main__":
    raise SystemExit(main())
