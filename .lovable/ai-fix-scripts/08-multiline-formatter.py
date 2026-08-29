#!/usr/bin/env python3
"""08-multiline-formatter.py - Formats Go function definitions (>2 params) to multi-line layout with trailing comma."""
import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp'
}
FUNC_PATTERN = re.compile(r'^(func(?:\s*\([^)]*\))?\s*[A-Za-z0-9_]+)\(([^)]+)\)(.*)$')


def compute_line_indent(line: str) -> str:
    indent = ""
    for ch in line:
        if ch not in (' ', '\t'):
            break
        indent += ch
    return indent


def reformat_func_line(line: str, stripped: str) -> tuple[str, bool]:
    is_candidate = stripped.startswith('func ') and '(' in stripped and ')' in stripped and len(stripped) > 100
    if not is_candidate:
        return line, False

    m = FUNC_PATTERN.match(stripped)
    if not m:
        return line, False

    fn_head = m.group(1)
    param_str = m.group(2)
    fn_tail = m.group(3)
    params = [p.strip() for p in param_str.split(',') if p.strip()]
    if len(params) <= 2:
        return line, False

    indent = compute_line_indent(line)
    param_indent = indent + "\t"
    formatted_params = ",\n".join(param_indent + p for p in params) + ","
    reformatted = f"{indent}{fn_head}(\n{formatted_params}\n{indent}){fn_tail}"
    return reformatted, True


def format_go_file(filepath: Path) -> int:
    try:
        content = filepath.read_text(encoding='utf-8')
    except Exception:
        return 0

    lines = content.split('\n')
    has_changed = False
    new_lines = []

    for line in lines:
        stripped = line.strip()
        reformatted, was_reformatted = reformat_func_line(line, stripped)
        if was_reformatted:
            has_changed = True
            new_lines.append(reformatted)
            continue
        new_lines.append(line)

    if not has_changed:
        return 0

    filepath.write_text('\n'.join(new_lines), encoding='utf-8')
    return 1


def main():
    target_dirs = [ROOT_DIR / 'gitmap', ROOT_DIR / 'src']
    files_modified = 0

    for td in target_dirs:
        if not td.exists():
            continue
        for root, dirs, files in os.walk(td):
            dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
            for file in files:
                fp = Path(root) / file
                if fp.suffix == '.go':
                    files_modified += format_go_file(fp)

    print(f"Formatted {files_modified} Go file(s) with multi-line parameters.")


if __name__ == '__main__':
    main()
