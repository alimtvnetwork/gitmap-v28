#!/usr/bin/env python3
"""
check-enum-and-boolean.py - Linter for Boolean Conventions, Naming, Enums, and Conditionals.

Checks:
1. No explicit boolean comparisons (== true, === true, == false, === false).
2. No inverted success checks (!isSuccess).
3. No single-line if blocks (anti-compression rule).
4. All enum definitions must end in 'Type'.
5. No nested if statements (nesting depth > 1 inside an if/else body in Go source code).
"""

import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent
TARGET_EXTS = {'.go', '.ts', '.tsx', '.php', '.py'}
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts'
}

# Regex patterns
EXPLICIT_BOOL_REGEX = re.compile(r'\b(==\s*true|===\s*true|==\s*false|===\s*false)\b')
INVERTED_SUCCESS_REGEX = re.compile(r'!\s*[a-zA-Z0-9_$.->]*\bisSuccess\b')
SINGLE_LINE_IF_REGEX = re.compile(r'^\s*if\b.*\{[^{}]+\}\s*$')
ENUM_MISSING_TYPE_REGEX = re.compile(r'^\s*(?:export\s+)?enum\s+(?!\w+Type\b)\w+\b')

def is_comment_or_doc(line: str, ext: str) -> bool:
    s = line.strip()
    if not s:
        return True
    if ext in ('.go', '.ts', '.tsx', '.php') and (s.startswith('//') or s.startswith('/*') or s.startswith('*')):
        return True
    if ext == '.py' and (s.startswith('#') or s.startswith('"""') or s.startswith("'''")):
        return True
    return False

def strip_go_code(content: str) -> list[tuple[int, str]]:
    """
    Strips comments, string literals, and rune literals from Go code,
    preserving exact line numbers and code structure.
    """
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
                elif content[i] == '\n':
                    break
                else:
                    i += 1
            if i < n and content[i] == '"':
                i += 1
            continue
        if c == "'":
            i += 1
            while i < n and content[i] != "'":
                if content[i] == '\\':
                    i += 2
                elif content[i] == '\n':
                    break
                else:
                    i += 1
            if i < n and content[i] == "'":
                i += 1
            continue

        current_line.append(c)
        i += 1

    if current_line or content.endswith('\n'):
        out_lines.append((line_num, ''.join(current_line)))

    return out_lines

def check_nested_ifs_in_go(content: str, filepath: Path) -> list[tuple[int, str]]:
    if "gitmap" in filepath.parts and "constants" in filepath.parts:
        return []

    stripped = strip_go_code(content)
    violations = []
    if_stack = []
    brace_depth = 0

    for line_num, code in stripped:
        trimmed = code.strip()
        if not trimmed:
            continue

        is_else_if = 'else if' in trimmed
        is_if = bool(re.search(r'\bif\b', trimmed)) and '{' in trimmed

        closes_at_start = len(re.match(r'^\}+', trimmed).group(0)) if re.match(r'^\}+', trimmed) else 0
        current_depth = brace_depth - closes_at_start

        while if_stack and current_depth <= if_stack[-1]:
            if_stack.pop()

        if is_if:
            if not is_else_if:
                if if_stack:
                    violations.append((line_num, f"Nested 'if' detected (depth {len(if_stack) + 1}): '{trimmed}'"))
                if_stack.append(brace_depth)

        opens = code.count('{')
        closes = code.count('}')
        brace_depth += (opens - closes)

        while if_stack and brace_depth <= if_stack[-1]:
            if_stack.pop()

    return violations

def check_file(filepath: Path) -> list[str]:
    ext = filepath.suffix
    if ext not in TARGET_EXTS:
        return []
    if filepath.name.endswith('_test.go') or '.test.' in filepath.name or '.spec.' in filepath.name:
        return []

    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception:
        return []

    lines = content.splitlines()
    violations = []
    is_constants_file = ("gitmap" in filepath.parts and "constants" in filepath.parts)

    for idx, line in enumerate(lines, start=1):
        if is_comment_or_doc(line, ext):
            continue

        # 1. Explicit boolean comparison
        if EXPLICIT_BOOL_REGEX.search(line):
            violations.append(f"{filepath}:{idx}: Explicit boolean comparison: '{line.strip()}'")

        # 2. Inverted success check
        if INVERTED_SUCCESS_REGEX.search(line):
            violations.append(f"{filepath}:{idx}: Inverted success check (!isSuccess): '{line.strip()}'")

        # 3. Single-line if statement
        if not is_constants_file and SINGLE_LINE_IF_REGEX.search(line) and not line.strip().startswith('//'):
            violations.append(f"{filepath}:{idx}: Single-line 'if' block (must expand to multiple lines): '{line.strip()}'")

        # 4. Enums without Type suffix
        if ext in ('.ts', '.tsx', '.php') and ENUM_MISSING_TYPE_REGEX.search(line):
            violations.append(f"{filepath}:{idx}: Enum missing 'Type' suffix: '{line.strip()}'")

    # 5. Nested ifs in Go files
    if ext == '.go':
        nested = check_nested_ifs_in_go(content, filepath)
        for line_num, msg in nested:
            violations.append(f"{filepath}:{line_num}: {msg}")

    return violations

def main():
    scanned_count = 0
    all_violations = []

    for root, dirs, files in os.walk(ROOT_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS and not any(ex in d for ex in EXCLUDE_DIRS)]
        for file in files:
            p = Path(root) / file
            if p.suffix in TARGET_EXTS and not p.name.endswith('_test.go'):
                scanned_count += 1
                v = check_file(p)
                all_violations.extend(v)

    print(f"Scanned {scanned_count} source files for boolean, enum, and conditional compliance.\n")

    if all_violations:
        print(f"❌ FAILED: Found {len(all_violations)} violation(s):")
        for v in all_violations[:50]:
            print(f"  - {v}")
        if len(all_violations) > 50:
            print(f"  ... and {len(all_violations) - 50} more.")
        sys.exit(1)
    else:
        print("✅ PASS: All boolean, naming, enum, and conditional checks passed (zero explicit booleans, inverted success, or nested ifs).")
        sys.exit(0)

if __name__ == '__main__':
    main()
