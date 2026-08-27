#!/usr/bin/env python3
"""
check-readme-install-section.py — Enforce two install-section invariants on
the root `readme.md`:

  1. The "Install in One Line" section is the FIRST content section after the
     badges block (i.e. the first `<h2>` / `##` heading after the
     `<!-- /STAMP:PLATFORM_BADGES -->` marker).
  2. The "Bundle Installers" section appears immediately after the install
     section and before the Table of Contents, mirroring the UI order.
  3. Every install code fence (powershell / bash / sh / pwsh) inside that
     section, AND every fence inside any "Bundle Installers" section,
     contains EXACTLY ONE non-empty command line, with NO inline `#`
     comments, NO blank lines, and NO multi-line `\\` continuations.

Spec: `.lovable/memory/constraints/install-command-formatting.md`
      `spec/01-spec-authoring-guide/11-root-readme-conventions.md`

Exit codes:
  0  PASS — install section position is correct AND every install fence is a
            single bare command.
  1  FAIL — at least one violation found (each printed with file:line).
  2  ERROR — readme.md not found or unreadable.

Usage:
  python3 linter-scripts/check-readme-install-section.py
  python3 linter-scripts/check-readme-install-section.py --readme readme.md
"""
from __future__ import annotations

import argparse
import re
import sys
from dataclasses import dataclass
from pathlib import Path

EXIT_PASS = 0
EXIT_FAIL = 1
EXIT_ERROR = 2

INSTALL_HEADING_PATTERN = re.compile(r"install\s+in\s+one\s+line", re.IGNORECASE)
BUNDLE_HEADING_PATTERN = re.compile(r"bundle\s+installers?", re.IGNORECASE)
TOC_HEADING_PATTERN = re.compile(r"table\s+of\s+contents", re.IGNORECASE)
HEADING_PATTERN = re.compile(
    r"^\s*(?:##\s+(?P<md>.+?)|<h2\b[^>]*>(?P<html>.+?)</h2>)\s*$"
)
PLATFORM_BADGES_END = "<!-- /STAMP:PLATFORM_BADGES -->"
INSTALL_FENCE_LANGS = {"powershell", "pwsh", "bash", "sh", "shell", "console"}


@dataclass(frozen=True)
class Violation:
    line: int
    rule: str
    detail: str


def main() -> int:
    args = parse_args()
    readme_path = Path(args.readme)
    text = read_readme(readme_path)
    if text is None:
        print(f"::error file={readme_path}::readme.md not found or unreadable")
        return EXIT_ERROR

    lines = text.splitlines()
    violations: list[Violation] = []
    violations.extend(check_section_position(lines))
    violations.extend(check_bundle_section_position(lines))
    violations.extend(check_install_fences(lines))

    if not violations:
        print(f"OK {readme_path}: install section position + fence formatting valid")
        return EXIT_PASS

    print_violations(readme_path, violations)
    return EXIT_FAIL


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--readme", default="readme.md", help="Path to root readme.md")
    return parser.parse_args()


def read_readme(path: Path) -> str | None:
    try:
        return path.read_text(encoding="utf-8")
    except OSError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return None


def check_section_position(lines: list[str]) -> list[Violation]:
    """Rule 1 — Install heading must be the first H2 after platform badges."""
    badges_end_line = find_platform_badges_end(lines)
    if badges_end_line is None:
        return [Violation(line=1, rule="MISSING_BADGES",
                          detail=f"missing marker {PLATFORM_BADGES_END!r}")]

    first_heading = find_first_h2_after(lines, badges_end_line)
    if first_heading is None:
        return [Violation(line=badges_end_line + 1, rule="MISSING_INSTALL_HEADING",
                          detail="no H2 heading found after badges block")]

    heading_line, heading_text = first_heading
    if INSTALL_HEADING_PATTERN.search(heading_text):
        return []

    return [Violation(
        line=heading_line,
        rule="INSTALL_NOT_FIRST",
        detail=(f"first section after badges is {heading_text!r}; "
                "expected 'Install in One Line'"),
    )]


def find_platform_badges_end(lines: list[str]) -> int | None:
    for index, line in enumerate(lines, start=1):
        if PLATFORM_BADGES_END in line:
            return index
    return None


def find_first_h2_after(lines: list[str], start_line: int) -> tuple[int, str] | None:
    for index in range(start_line, len(lines)):
        match = HEADING_PATTERN.match(lines[index])
        if match is None:
            continue
        return index + 1, strip_heading_decoration(heading_text_from(match))
    return None


def find_all_h2(lines: list[str]) -> list[tuple[int, str]]:
    headings: list[tuple[int, str]] = []
    for index, raw in enumerate(lines, start=1):
        match = HEADING_PATTERN.match(raw)
        if match is None:
            continue
        headings.append((index, strip_heading_decoration(heading_text_from(match))))
    return headings


def check_bundle_section_position(lines: list[str]) -> list[Violation]:
    headings = find_all_h2(lines)
    install_index = find_heading_index(headings, INSTALL_HEADING_PATTERN)
    if install_index is None:
        return []

    bundle_index = find_heading_index(headings, BUNDLE_HEADING_PATTERN)
    if bundle_index is None:
        install_line, _ = headings[install_index]
        return [Violation(
            line=install_line,
            rule="MISSING_BUNDLE_SECTION",
            detail="missing 'Bundle Installers' section after 'Install in One Line'",
        )]

    toc_index = find_heading_index(headings, TOC_HEADING_PATTERN)
    bundle_line, _ = headings[bundle_index]

    if bundle_index != install_index + 1:
        return [Violation(
            line=bundle_line,
            rule="BUNDLE_SECTION_ORDER",
            detail="'Bundle Installers' must appear immediately after 'Install in One Line'",
        )]

    if toc_index is not None and bundle_index > toc_index:
        return [Violation(
            line=bundle_line,
            rule="BUNDLE_AFTER_TOC",
            detail="'Bundle Installers' must appear before 'Table of Contents'",
        )]

    return []


def find_heading_index(headings: list[tuple[int, str]], pattern: re.Pattern[str]) -> int | None:
    for index, (_, heading) in enumerate(headings):
        if pattern.search(heading):
            return index
    return None


def strip_heading_decoration(text: str) -> str:
    cleaned = re.sub(r"<[^>]+>", "", text)
    cleaned = re.sub(r"[\u2600-\u27BF\U0001F300-\U0001FAFF]", "", cleaned)
    return cleaned.strip()


def heading_text_from(match: re.Match[str]) -> str:
    return match.group("md") or match.group("html") or ""


def check_install_fences(lines: list[str]) -> list[Violation]:
    """Rule 2 — every install fence must be a single bare command."""
    violations: list[Violation] = []
    in_install_zone = False
    fence_open_line: int | None = None
    fence_lang: str | None = None
    fence_body: list[str] = []

    for index, raw in enumerate(lines, start=1):
        if fence_open_line is None:
            in_install_zone = update_zone(raw, in_install_zone)
            opened = try_open_fence(raw, in_install_zone)
            if opened is None:
                continue
            fence_open_line, fence_lang = index, opened
            fence_body = []
            continue

        if is_fence_close(raw):
            violations.extend(validate_fence_body(
                fence_open_line, fence_lang or "", fence_body,
            ))
            fence_open_line, fence_lang, fence_body = None, None, []
            continue

        fence_body.append(raw)

    return violations


def update_zone(raw: str, in_install_zone: bool) -> bool:
    match = HEADING_PATTERN.match(raw)
    if match is None:
        return in_install_zone
    heading = strip_heading_decoration(heading_text_from(match))
    if INSTALL_HEADING_PATTERN.search(heading):
        return True
    if BUNDLE_HEADING_PATTERN.search(heading):
        return True
    return False


def try_open_fence(raw: str, in_install_zone: bool) -> str | None:
    if not in_install_zone:
        return None
    stripped = raw.lstrip()
    if not stripped.startswith("```"):
        return None
    lang = stripped[3:].strip().lower()
    if lang not in INSTALL_FENCE_LANGS:
        return None
    return lang


def is_fence_close(raw: str) -> bool:
    stripped = raw.strip()
    return stripped == "```"


def validate_fence_body(open_line: int, lang: str,
                        body: list[str]) -> list[Violation]:
    violations: list[Violation] = []
    command_count = 0

    for offset, raw in enumerate(body, start=1):
        line_no = open_line + offset
        stripped = raw.strip()
        if stripped == "":
            violations.append(Violation(
                line=line_no, rule="BLANK_LINE_IN_FENCE",
                detail=f"blank line inside {lang} install fence",
            ))
            continue
        if stripped.startswith("#"):
            violations.append(Violation(
                line=line_no, rule="COMMENT_IN_FENCE",
                detail=f"inline comment {stripped!r} inside {lang} install fence",
            ))
            continue
        if stripped.endswith("\\"):
            violations.append(Violation(
                line=line_no, rule="LINE_CONTINUATION_IN_FENCE",
                detail="multi-line `\\` continuation is forbidden in install fence",
            ))
            continue
        command_count += 1

    if command_count != 1:
        violations.append(Violation(
            line=open_line, rule="WRONG_COMMAND_COUNT",
            detail=(f"{lang} install fence has {command_count} commands; "
                    "exactly 1 is required"),
        ))

    return violations


def print_violations(path: Path, violations: list[Violation]) -> None:
    print(f"FAIL {path}: {len(violations)} install-section violation(s)")
    for v in violations:
        print(f"::error file={path},line={v.line}::[{v.rule}] {v.detail}")


if __name__ == "__main__":
    sys.exit(main())
