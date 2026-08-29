#!/usr/bin/env python3
"""02-guideline-autofixer.py - Automatically formats code to comply with Return New Line (R13-R16), Blank Line before if/control structures, and styling guidelines."""
import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
TARGET_EXTS = {'.go', '.ts', '.tsx', '.js', '.jsx'}
EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', '.next', '.gitmap', 'vendor'}


def is_return_needing_newline(stripped: str, prev: str) -> bool:
    is_return = stripped.startswith('return ') or stripped == 'return' or stripped.startswith('throw ')
    if not is_return:
        return False
    if prev == '':
        return False
    if prev.endswith(('{', '}', ':')):
        return False
    if prev.startswith(('//', '/*', '*')):
        return False
    return True


def is_brace_needing_newline(prev_line: str, stripped: str, out_len: int, out_list: list[str]) -> bool:
    if prev_line != '}':
        return False
    if stripped == '':
        return False
    if stripped.startswith(('}', 'else', 'catch', 'finally', ')', ']', ',', ';', '//', '/*', '</')):
        return False
    if out_len >= 2 and out_list[-2].strip() == '':
        return False
    return True


def is_control_needing_newline(stripped: str, prev: str) -> bool:
    is_ctrl = stripped.startswith(('if ', 'if(', 'for ', 'for(', 'switch ', 'switch(', 'while ', 'while('))
    if not is_ctrl:
        return False
    if prev == '':
        return False
    if prev.endswith(('{', ':', '//', '/*', '*')):
        return False
    if prev == '}':
        return False
    return True


def autofix_newlines(content: str) -> str:
    lines = content.split('\n')
    out = []

    for line in lines:
        stripped = line.strip()

        # Rule 1: No double empty lines (\n\n\n) -> collapse to single empty line
        if stripped == '' and out and out[-1].strip() == '':
            continue

        # Rule 2: No empty line at the start of a function/block (after '{')
        if stripped == '' and out and out[-1].strip().endswith('{'):
            continue

        prev = out[-1].strip() if out else ''

        # Rule 3: Blank line required before return/throw
        if is_return_needing_newline(stripped, prev):
            out.append('')
            prev = ''

        # Rule 4: Blank line required after '}' if followed by more code
        if is_brace_needing_newline(prev, stripped, len(out), out):
            out.append('')
            prev = ''

        # Rule 5: Blank line required before if / control structures
        if is_control_needing_newline(stripped, prev):
            out.append('')

        out.append(line)

    result = '\n'.join(out)
    result = result.rstrip(' \t\r\n') + '\n'
    return result


def fix_file(filepath: Path) -> bool:
    try:
        content = filepath.read_text(encoding='utf-8')
    except Exception:
        return False

    fixed = autofix_newlines(content)
    if fixed == content:
        return False

    filepath.write_text(fixed, encoding='utf-8', newline='\n')
    return True


def process_path(p: Path) -> int:
    if p.is_file():
        is_fixed = fix_file(p)
        return 1 if is_fixed else 0

    count = 0
    for root, dirs, files in os.walk(p):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for file in files:
            fp = Path(root) / file
            if fp.suffix.lower() not in TARGET_EXTS:
                continue
            if fix_file(fp):
                count += 1
    return count


def main():
    paths = sys.argv[1:]
    if not paths:
        paths = [str(ROOT_DIR / 'src')]

    count = 0
    for p_str in paths:
        count += process_path(Path(p_str))

    print(f"Autofixer finished: modified {count} file(s).")


if __name__ == "__main__":
    main()
