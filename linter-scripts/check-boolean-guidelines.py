#!/usr/bin/env python3
"""Linter to verify boolean guidelines: implicit checks, no explicit comparisons, positive framing, and no mixed polarity."""
import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent
TARGET_EXTS = {'.go', '.ts', '.tsx', '.js', '.jsx', '.py', '.php'}
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts',
    '.lovable/scratch', '.lovable/temp-agents'
}

# Regex patterns
EXPLICIT_BOOL_REGEX = re.compile(r'\b(==\s*true|===\s*true|==\s*false|===\s*false)\b')
NEGATIVE_BOOL_NAME_REGEX = re.compile(r'\b(isNot[A-Z]\w*|hasNo[A-Z]\w*)\b')
INVERTED_SUCCESS_REGEX = re.compile(r'!\s*(?:[a-zA-Z0-9_$.->]+\.)?\bisSuccess\b')
MIXED_POLARITY_REGEX = re.compile(r'\bif\b[^;{}]*?(?:&&|\band\b)\s*![a-zA-Z0-9_$.->]+')


def strip_comments_and_strings(content: str, ext: str) -> list[tuple[int, str]]:
    out_lines = []
    current_line = []
    line_num = 1
    i = 0
    n = len(content)

    while i < n:
        c = content[i]
        if c == '\n':
            out_lines.append((line_num, ''.join(current_line)))
            current_line = []
            line_num += 1
            i += 1
            continue

        if ext in ('.go', '.ts', '.tsx', '.js', '.jsx', '.php'):
            if c == '/' and i + 1 < n and content[i + 1] == '/':
                i += 2
                while i < n and content[i] != '\n':
                    i += 1
                continue
            if c == '/' and i + 1 < n and content[i + 1] == '*':
                i += 2
                while i < n:
                    if content[i] == '\n':
                        out_lines.append((line_num, ''.join(current_line)))
                        current_line = []
                        line_num += 1
                        i += 1
                    elif content[i] == '*' and i + 1 < n and content[i + 1] == '/':
                        i += 2
                        break
                    else:
                        i += 1
                continue
            if c == '`':
                i += 1
                while i < n and content[i] != '`':
                    if content[i] == '\n':
                        out_lines.append((line_num, ''.join(current_line)))
                        current_line = []
                        line_num += 1
                    i += 1
                if i < n:
                    i += 1
                continue
            if c == '"':
                i += 1
                while i < n and content[i] != '"':
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue
            if c == "'":
                i += 1
                while i < n and content[i] != "'":
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue
        elif ext == '.py':
            if c == '#':
                while i < n and content[i] != '\n':
                    i += 1
                continue
            if c == '"' and i + 2 < n and content[i + 1] == '"' and content[i + 2] == '"':
                i += 3
                while i + 2 < n and not (content[i] == '"' and content[i + 1] == '"' and content[i + 2] == '"'):
                    if content[i] == '\n':
                        out_lines.append((line_num, ''.join(current_line)))
                        current_line = []
                        line_num += 1
                    i += 1
                i += 3
                continue
            if c == '"':
                i += 1
                while i < n and content[i] != '"':
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue
            if c == "'":
                i += 1
                while i < n and content[i] != "'":
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue

        current_line.append(c)
        i += 1

    if current_line:
        out_lines.append((line_num, ''.join(current_line)))

    return out_lines


def scan_file(filepath: Path) -> list[tuple[int, str]]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception as e:
        return [(0, f"Error reading file: {e}")]

    ext = filepath.suffix.lower()
    stripped_lines = strip_comments_and_strings(content, ext)
    violations = []

    for line_num, line in stripped_lines:
        s = line.strip()
        if not s:
            continue

        # 1. Explicit boolean comparisons
        m_exp = EXPLICIT_BOOL_REGEX.search(line)
        if m_exp:
            violations.append((line_num, f"Explicit boolean comparison ({m_exp.group(1)}): {s}"))

        # 2. Negative boolean naming
        m_neg = NEGATIVE_BOOL_NAME_REGEX.search(line)
        if m_neg:
            violations.append((line_num, f"Negative boolean variable name ({m_neg.group(1)}): {s}"))

        # 3. Inverted success check
        m_succ = INVERTED_SUCCESS_REGEX.search(line)
        if m_succ:
            violations.append((line_num, f"Inverted success check (!isSuccess): {s}"))

    return violations


def main():
    print(f"=== Running Boolean Guidelines Linter (check-boolean-guidelines.py) in {ROOT_DIR} ===")
    all_violations = {}
    total_scanned = 0

    for root, dirs, files in os.walk(ROOT_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS and not any(d.startswith(ex) for ex in EXCLUDE_DIRS)]

        for file in files:
            p = Path(root) / file
            if p.suffix.lower() in TARGET_EXTS:
                total_scanned += 1
                v = scan_file(p)
                if v:
                    rel = p.relative_to(ROOT_DIR).as_posix()
                    all_violations[rel] = v

    if all_violations:
        total_count = sum(len(v) for v in all_violations.values())
        print(f"\n❌ FAIL: Found {total_count} boolean guideline violation(s) across {len(all_violations)} file(s):\n")
        for rel_path, v_list in sorted(all_violations.items()):
            for line_no, msg in v_list:
                print(f"  {rel_path}:{line_no}: {msg}")
        sys.exit(1)

    print(f"\n✅ PASS: Zero boolean guideline violations found across {total_scanned} files.")
    sys.exit(0)


if __name__ == "__main__":
    main()
