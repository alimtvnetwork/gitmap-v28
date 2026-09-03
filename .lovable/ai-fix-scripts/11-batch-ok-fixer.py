#!/usr/bin/env python3
"""07-batch-ok-fixer.py - Comprehensively refactors remaining bare ok patterns to semantic affirmative booleans."""
import os
import re
import sys
from pathlib import Path

ROOT_DIR = Path(__file__).resolve().parent.parent.parent

EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', '.next', '.gitmap', 'vendor', 'brain', '.gemini'}


def choose_semantic_name(line: str, var_name: str) -> str:
    line_lower = line.lower()
    var_lower = var_name.lower()

    if 'err' in var_lower or 'exiterr' in var_lower or 'apperr' in var_lower:
        return 'isAppErr' if 'apperr' in line_lower else 'isExitErr' if 'exiterr' in line_lower else 'isErr'
    if 'convert' in line_lower or 'url' in line_lower or 'ssh' in line_lower or 'https' in line_lower:
        return 'isConverted' if 'convert' in line_lower else 'hasURL'
    if 'resolve' in line_lower or 'target' in var_lower or 'path' in var_lower or 'dest' in var_lower:
        return 'isResolved'
    if 'match' in line_lower or 'regex' in line_lower or 'suggest' in line_lower:
        return 'isMatch'
    if 'parse' in line_lower or 'extract' in line_lower or 'cut' in line_lower or 'split' in line_lower:
        return 'isParsed'
    if 'cache' in line_lower or 'lookup' in line_lower:
        return 'isCached' if 'cache' in line_lower else 'isFound'
    if 'probe' in line_lower or 'detect' in line_lower or 'check' in line_lower:
        return 'isDetected' if 'detect' in line_lower else 'isFound'
    if 'profile' in line_lower:
        return 'isProfileFound' if 'profile' in var_lower else 'isResolved'
    if 'map' in line_lower or '[' in line:
        return 'isFound'
    if 'read' in line_lower or 'load' in line_lower or 'get' in line_lower or 'fetch' in line_lower:
        return 'isLoaded' if 'load' in line_lower else 'isRead' if 'read' in line_lower else 'isFound'
    if 'clone' in line_lower:
        return 'isCloned' if 'run' in line_lower else 'isResolved'
    if 'tree' in line_lower or 'sha' in line_lower or 'hash' in line_lower:
        return 'hasSHA' if 'sha' in line_lower else 'hasHash' if 'hash' in line_lower else 'hasTree'

    return 'isOk'


def refactor_line_content(content: str) -> str:
    lines = content.split('\n')
    new_lines = []
    i = 0
    while i < len(lines):
        line = lines[i]

        # Pattern 1: if v, ok := expr; ok {
        m_if_ok = re.search(r'(\bif\s+(?:(?:\w+|_,?)\s*,\s*)*)(\w+),\s*ok\s*(:=|=)\s*(.+?);\s*ok(\s*\{.*)', line)
        if m_if_ok:
            prefix, var_name, assign_op, expr, suffix = m_if_ok.groups()
            sem_name = choose_semantic_name(line, var_name)
            line = f"{prefix}{var_name}, {sem_name} {assign_op} {expr}; {sem_name}{suffix}"
            lines[i] = line
            new_lines.append(line)
            i += 1
            continue

        # Pattern 2: if v, ok := expr; !ok {
        m_if_not_ok = re.search(r'(\bif\s+(?:(?:\w+|_,?)\s*,\s*)*)(\w+),\s*ok\s*(:=|=)\s*(.+?);\s*!ok(\s*\{.*)', line)
        if m_if_not_ok:
            prefix, var_name, assign_op, expr, suffix = m_if_not_ok.groups()
            sem_name = choose_semantic_name(line, var_name)
            line = f"{prefix}{var_name}, {sem_name} {assign_op} {expr}; !{sem_name}{suffix}"
            lines[i] = line
            new_lines.append(line)
            i += 1
            continue

        # Pattern 3: if _, ok := expr; !ok {
        m_if_blank_not_ok = re.search(r'(\bif\s+_,\s*)ok\s*(:=|=)\s*(.+?);\s*!ok(\s*\{.*)', line)
        if m_if_blank_not_ok:
            prefix, assign_op, expr, suffix = m_if_blank_not_ok.groups()
            sem_name = choose_semantic_name(line, 'item')
            line = f"{prefix}{sem_name} {assign_op} {expr}; !{sem_name}{suffix}"
            lines[i] = line
            new_lines.append(line)
            i += 1
            continue

        # Pattern 4: if _, ok := expr; ok {
        m_if_blank_ok = re.search(r'(\bif\s+_,\s*)ok\s*(:=|=)\s*(.+?);\s*ok(\s*\{.*)', line)
        if m_if_blank_ok:
            prefix, assign_op, expr, suffix = m_if_blank_ok.groups()
            sem_name = choose_semantic_name(line, 'item')
            line = f"{prefix}{sem_name} {assign_op} {expr}; {sem_name}{suffix}"
            lines[i] = line
            new_lines.append(line)
            i += 1
            continue

        res_line, is_matched = handle_assign_ok(line, lines, i)
        if is_matched:
            new_lines.append(res_line)
            i += 1
            continue

        res_line, is_blank_matched = handle_assign_blank_ok(line, lines, i)
        if is_blank_matched:
            new_lines.append(res_line)
            i += 1
            continue

        new_lines.append(line)
        i += 1

    return '\n'.join(new_lines)


def update_next_line(lines: list[str], i: int, sem_name: str) -> None:
    if i + 1 >= len(lines):
        return
    next_line = lines[i + 1]
    if re.search(r'\bif\s+!ok\b', next_line):
        lines[i + 1] = re.sub(r'\bif\s+!ok\b', f'if !{sem_name}', next_line)
        return
    if re.search(r'\bif\s+ok\b', next_line):
        lines[i + 1] = re.sub(r'\bif\s+ok\b', f'if {sem_name}', next_line)


def handle_assign_ok(line: str, lines: list[str], i: int) -> tuple[str, bool]:
    m_assign_ok = re.search(r'^(\s*)(?:(?:\w+|_,?)\s*,\s*)*(\w+),\s*ok\s*(:=|=)\s*(.+)$', line)
    if not m_assign_ok or line.strip().startswith(('//', '/*', '*')):
        return line, False
    indent, var_name, assign_op, expr = m_assign_ok.groups()
    sem_name = choose_semantic_name(line, var_name)
    line = f"{indent}{line.strip()[:-len('ok' + assign_op + expr)].rstrip(', ')}, {sem_name} {assign_op} {expr}"
    lines[i] = line
    update_next_line(lines, i, sem_name)
    return line, True


def handle_assign_blank_ok(line: str, lines: list[str], i: int) -> tuple[str, bool]:
    m_assign_blank_ok = re.search(r'^(\s*)_,\s*ok\s*(:=|=)\s*(.+)$', line)
    if not m_assign_blank_ok or line.strip().startswith(('//', '/*', '*')):
        return line, False
    indent, assign_op, expr = m_assign_blank_ok.groups()
    sem_name = choose_semantic_name(line, 'item')
    line = f"{indent}_, {sem_name} {assign_op} {expr}"
    lines[i] = line
    update_next_line(lines, i, sem_name)
    return line, True


def process_all_go_files():
    count = 0
    for root, dirs, files in os.walk(ROOT_DIR / 'gitmap'):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for f in files:
            if not f.endswith('.go'):
                continue
            fp = Path(root) / f
            content = fp.read_text(encoding='utf-8')
            new_content = refactor_line_content(content)
            if new_content != content:
                fp.write_text(new_content, encoding='utf-8')
                print(f"Batch fixed {fp.relative_to(ROOT_DIR).as_posix()}")
                count += 1
    print(f"Total files batch fixed: {count}")


if __name__ == '__main__':
    process_all_go_files()
