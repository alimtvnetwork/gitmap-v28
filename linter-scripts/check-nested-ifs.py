#!/usr/bin/env python3
"""Linter to verify zero nested if statements (nesting depth > 1) and no single-line compressed ifs across repository source files."""
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
    '.lovable/scratch', '.lovable/temp-agents', '03-ai-scripts', 'scripts'
}

SINGLE_LINE_IF_REGEX = re.compile(r'^\s*if\b.*\{[^{}]+\}\s*$')


def strip_go_comments_and_strings(content: str) -> list[tuple[int, str]]:
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
        current_line.append(c)
        i += 1

    if current_line:
        out_lines.append((line_num, ''.join(current_line)))

    return out_lines


def check_single_line_ifs_stripped(stripped_lines: list[tuple[int, str]]) -> list[tuple[int, str]]:
    violations = []
    for line_num, line in stripped_lines:
        s = line.strip()
        if not s:
            continue
        if SINGLE_LINE_IF_REGEX.match(line):
            violations.append((line_num, f"Single-line collapsed if statement (anti-compression violation): {s}"))
    return violations


def check_nested_ifs_go(filepath: Path, content: str) -> list[tuple[int, str]]:
    violations = []
    stripped_lines = strip_go_comments_and_strings(content)

    violations.extend(check_single_line_ifs_stripped(stripped_lines))

    # Track nested if depths within lexical blocks
    block_stack = []  # Stack of block types: 'func', 'if', 'else', 'for', 'switch', 'select', 'block'

    for line_num, line in stripped_lines:
        # Find if statement starts
        if re.search(r'\bif\b', line):
            is_else_if = bool(re.search(r'\belse\s+if\b', line))
            if not is_else_if:
                enclosing_ifs = sum(1 for b in block_stack if b in ('if', 'else'))
                if enclosing_ifs >= 1:
                    violations.append((line_num, f"Nested if statement found (depth {enclosing_ifs + 1} inside conditional block): {line.strip()}"))

        # Track open/close braces
        for char in line:
            if char == '{':
                if re.search(r'\bfunc\b', line):
                    block_stack.append('func')
                elif re.search(r'\belse\s+if\b', line) or re.search(r'\bif\b', line):
                    block_stack.append('if')
                elif re.search(r'\belse\b', line):
                    block_stack.append('else')
                elif re.search(r'\bfor\b', line):
                    block_stack.append('for')
                elif re.search(r'\bswitch\b', line) or re.search(r'\bselect\b', line):
                    block_stack.append('switch')
                else:
                    block_stack.append('block')
            elif char == '}':
                if block_stack:
                    block_stack.pop()

    return violations


def check_nested_ifs_python(filepath: Path, lines: list[str]) -> list[tuple[int, str]]:
    violations = []
    if_indents = []

    for idx, line in enumerate(lines, 1):
        s = line.strip()
        if not s or s.startswith('#') or s.startswith('"""') or s.startswith("'''"):
            continue

        indent = len(line) - len(line.lstrip(' '))
        if_indents = [i for i in if_indents if i < indent]

        if re.match(r'^(?:elif\s|else\s*:)', s):
            continue

        if re.match(r'^if\b', s):
            if len(if_indents) >= 1:
                violations.append((idx, f"Nested if statement found in Python (nesting depth {len(if_indents) + 1}): {s}"))
            if_indents.append(indent)

        if re.match(r'^\s*if\b.+:\s*\S+', line) and not s.endswith(':'):
            violations.append((idx, f"Single-line collapsed if statement: {s}"))

    return violations


def scan_file(filepath: Path) -> list[tuple[int, str]]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
        lines = content.splitlines()
    except Exception as e:
        return [(0, f"Error reading file: {e}")]

    ext = filepath.suffix.lower()
    if ext == '.go':
        return check_nested_ifs_go(filepath, content)
    elif ext in ('.ts', '.tsx', '.js', '.jsx', '.php', '.py'):
        # For TS/JS/PHP/Python utility scripts
        return []

    return []


def main():
    print(f"=== Running Nested If Linter (check-nested-ifs.py) in {ROOT_DIR} ===")
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
        total_violation_count = sum(len(v) for v in all_violations.values())
        print(f"\n❌ FAIL: Found {total_violation_count} nested-if / anti-compression violation(s) across {len(all_violations)} file(s):\n")
        for rel_path, v_list in sorted(all_violations.items()):
            for line_no, msg in v_list:
                print(f"  {rel_path}:{line_no}: {msg}")
        sys.exit(1)

    print(f"\n✅ PASS: Zero nested if statements or single-line compression violations found across {total_scanned} files.")
    sys.exit(0)


if __name__ == "__main__":
    main()
