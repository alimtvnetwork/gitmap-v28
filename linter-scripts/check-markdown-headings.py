#!/usr/bin/env python3
"""
check-markdown-headings.py
Linter rule: every Markdown heading (# through ######) must have
exactly one blank line before it (unless it is the very first content line)
and exactly one blank line after it (unless it is the very last line).

Exit 0 = all clean.
Exit 1 = violations found.

Usage:
    python3 linter-scripts/check-markdown-headings.py [--fix] [path ...]

    --fix  auto-fix files in place (idempotent)
    path   one or more files or directories (default: repo root .)
"""

import os
import re
import sys
from pathlib import Path

HEADING_RE = re.compile(r'^#{1,6}(\s|$)')

IGNORE_DIRS = {
    '.git', 'node_modules', '.github', '.ci-out', 'dist', 'build',
    'tmp', '.lovable/temp-scripts',
}


def is_heading(line: str) -> bool:
    return bool(HEADING_RE.match(line))


def check_file(filepath: str) -> list[tuple[int, str]]:
    """Return list of (line_number, message) violations."""
    try:
        text = Path(filepath).read_text(encoding='utf-8')
    except (OSError, UnicodeDecodeError):
        return []

    lines = text.split('\n')
    violations: list[tuple[int, str]] = []
    total = len(lines)

    for idx, line in enumerate(lines):
        if not is_heading(line):
            continue

        line_num = idx + 1  # 1-indexed for output

        # --- blank line BEFORE (skip if first non-empty line in file) ---
        is_first_content = all(l.strip() == '' for l in lines[:idx])
        if not is_first_content:
            prev_line = lines[idx - 1] if idx > 0 else ''
            if prev_line.strip() != '':
                violations.append((
                    line_num,
                    f'Heading must be preceded by a blank line (got: {repr(prev_line[:60])})',
                ))

        # --- blank line AFTER (skip if last line or heading is last content) ---
        is_last_line = idx + 1 >= total
        if not is_last_line:
            next_line = lines[idx + 1]
            if next_line.strip() != '':
                violations.append((
                    line_num,
                    f'Heading must be followed by a blank line (got: {repr(next_line[:60])})',
                ))

    return violations


def fix_file(filepath: str) -> bool:
    """Auto-fix: ensure blank line before and after every heading. Returns True if changed."""
    try:
        text = Path(filepath).read_text(encoding='utf-8')
    except (OSError, UnicodeDecodeError):
        return False

    lines = text.split('\n')
    result: list[str] = []
    changed = False

    for idx, line in enumerate(lines):
        is_first_content = all(l.strip() == '' for l in lines[:idx])

        if is_heading(line):
            # Ensure blank line before (unless first content)
            if not is_first_content and result and result[-1].strip() != '':
                result.append('')
                changed = True
            result.append(line)
            # Ensure blank line after (unless last line)
            if idx + 1 < len(lines) and lines[idx + 1].strip() != '':
                result.append('')
                changed = True
        else:
            result.append(line)

    # Collapse 3+ consecutive blank lines to 2 (avoid over-spacing)
    cleaned: list[str] = []
    blank_run = 0
    for ln in result:
        if ln.strip() == '':
            blank_run += 1
            if blank_run <= 2:
                cleaned.append(ln)
        else:
            blank_run = 0
            cleaned.append(ln)

    new_text = '\n'.join(cleaned)
    if new_text != text:
        Path(filepath).write_text(new_text, encoding='utf-8')
        return True

    return False


def collect_markdown_files(paths: list[str]) -> list[str]:
    """Collect all .md files from given paths (files or directories)."""
    collected: list[str] = []
    for raw_path in paths:
        p = Path(raw_path)
        if p.is_file() and p.suffix == '.md':
            collected.append(str(p))
        elif p.is_dir():
            for md in p.rglob('*.md'):
                parts = md.parts
                if any(part in IGNORE_DIRS for part in parts):
                    continue
                collected.append(str(md))
    return sorted(collected)


def main() -> None:
    args = sys.argv[1:]
    is_fix_mode = '--fix' in args
    paths_raw = [a for a in args if not a.startswith('--')]

    if not paths_raw:
        paths_raw = ['.']

    md_files = collect_markdown_files(paths_raw)

    if is_fix_mode:
        fixed_count = 0
        for filepath in md_files:
            was_fixed = fix_file(filepath)
            if was_fixed:
                rel = os.path.relpath(filepath)
                print(f'  fixed  {rel}')
                fixed_count += 1
        print(f'\nMarkdown heading fixer: {fixed_count} file(s) updated out of {len(md_files)} scanned.')
        sys.exit(0)
    else:
        total_violations = 0
        for filepath in md_files:
            file_violations = check_file(filepath)
            for line_num, msg in file_violations:
                rel = os.path.relpath(filepath)
                print(f'{rel}:{line_num}: [MD-H001] {msg}')
                total_violations += 1

        if total_violations > 0:
            print(f'\nMarkdown heading linter: {total_violations} violation(s) found.')
            sys.exit(1)
        else:
            print('Markdown heading linter: all files OK.')
            sys.exit(0)


if __name__ == '__main__':
    main()
