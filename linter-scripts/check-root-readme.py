#!/usr/bin/env python3
"""
check-root-readme.py — Enforce §9 of spec/01-spec-authoring-guide/11-root-readme-conventions.md
against the root readme.md.

Exit codes (stable contract):
  0  PASS — every required element present
  1  FAIL — at least one required element missing or malformed
  2  ERROR — readme.md not found / unreadable

Usage:
  python3 linter-scripts/check-root-readme.py
  python3 linter-scripts/check-root-readme.py --readme path/to/readme.md
"""
from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

EXIT_PASS = 0
EXIT_FAIL = 1
EXIT_ERROR = 2

MAX_LINES = None  # Root README has no length limit (per writing guideline §7).
MIN_BADGES = 8
MAX_BADGES = 16
REQUIRED_AUTHOR = "Md. Alim Ul Karim"
REQUIRED_COMPANY = "Riseup Asia LLC"
REQUIRED_AUTHOR_URL = "https://alimkarim.com/"
REQUIRED_COMPANY_URL = "https://riseup-asia.com/"
REQUIRED_SECTIONS = [
    "What is this",
    "For AI Agents",
    "Bundle Installers",
    "Full-Repo Install",
    "Documentation",
    "Contributing",
]
REQUIRED_STAMPS = [
    "STAMP:BADGES",
    "STAMP:PLATFORM_BADGES",
    "STAMP:VERSION",
    "STAMP:UPDATED",
    "STAMP:FILES",
    "STAMP:FOLDERS",
    "STAMP:LINES",
]


def has_centered_icon(body: str) -> bool:
    pattern = re.compile(
        r'<p\s+align="center">.*?<img[^>]+src="public/images/[^"]+-icon\.png"[^>]+width="160"',
        re.DOTALL,
    )
    return bool(pattern.search(body))


def has_centered_h1(body: str) -> bool:
    return bool(re.search(r'<h1\s+align="center">', body))


def count_badges(body: str) -> int:
    """Count only the hero badges between the two STAMP markers.
    Card-style shields in body sections are excluded by design."""
    total = 0
    for marker in ("BADGES", "PLATFORM_BADGES"):
        m = re.search(rf'<!-- STAMP:{marker} -->([\s\S]*?)<!-- /STAMP:{marker} -->', body)
        if not m:
            continue
        block = m.group(1)
        md = re.findall(r'!\[[^\]]*\]\(https://img\.shields\.io/', block)
        html = re.findall(r'<img\s[^>]*src="https://img\.shields\.io/', block)
        total += len(md) + len(html)
    return total


def has_author_block(body: str) -> bool:
    has_name = REQUIRED_AUTHOR in body
    has_company = REQUIRED_COMPANY in body
    has_author_url = REQUIRED_AUTHOR_URL in body
    has_company_url = REQUIRED_COMPANY_URL in body
    return has_name and has_company and has_author_url and has_company_url


def section_present(body: str, name: str) -> bool:
    md_heading = re.search(rf'^##\s.*{re.escape(name)}', body, re.MULTILINE)
    html_heading = re.search(rf'<h2[^>]*>.*{re.escape(name)}', body, re.IGNORECASE)
    return bool(md_heading or html_heading)


def stamp_present(body: str, key: str) -> bool:
    return f"<!-- {key} -->" in body


def collect_violations(readme: Path) -> list[str]:
    body = readme.read_text(encoding="utf-8")
    violations: list[str] = []

    if MAX_LINES is not None:
        line_count = body.count("\n") + 1
        if line_count > MAX_LINES:
            violations.append(f"length: {line_count} lines exceeds max {MAX_LINES} (extract to docs/)")

    if not has_centered_icon(body):
        violations.append("missing centered brand icon (<p align=\"center\"> + 160px <img>)")

    if not has_centered_h1(body):
        violations.append("missing centered <h1 align=\"center\"> wordmark")

    badges = count_badges(body)
    if badges < MIN_BADGES or badges > MAX_BADGES:
        violations.append(f"badge count {badges} outside required range {MIN_BADGES}-{MAX_BADGES}")

    if not has_author_block(body):
        violations.append(
            f"author block missing or malformed — must contain '{REQUIRED_AUTHOR}', "
            f"'{REQUIRED_COMPANY}', and both canonical URLs"
        )

    for stamp in REQUIRED_STAMPS:
        if not stamp_present(body, stamp):
            violations.append(f"missing stamp marker: {stamp}")

    for section in REQUIRED_SECTIONS:
        if not section_present(body, section):
            violations.append(f"missing required section: '{section}'")

    return violations


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate root readme.md against spec §9.")
    parser.add_argument("--readme", default="readme.md", help="Path to the root README (default: readme.md)")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    readme = Path(args.readme)
    if not readme.is_file():
        sys.stderr.write(f"ERROR: readme not found at {readme}\n")
        return EXIT_ERROR

    violations = collect_violations(readme)
    if violations:
        sys.stderr.write(f"FAIL: {readme} violates root-readme conventions:\n")
        for v in violations:
            sys.stderr.write(f"  - {v}\n")
        sys.stderr.write(
            "\nSee spec/01-spec-authoring-guide/11-root-readme-conventions.md §9 for the full checklist.\n"
        )
        return EXIT_FAIL

    print(f"OK {readme} passes all root-readme conventions (§9 checklist).")
    return EXIT_PASS


if __name__ == "__main__":
    raise SystemExit(main())