#!/usr/bin/env python3
"""
Script 23: Coding Guideline Path Consolidator
Autonomously consolidates all references from the old nested path
(.lovable/coding-guidelines/coding-guidelines.md) to the single canonical path
(.lovable/coding-guidelines.md) across all documentation, code, specs, and linters.
"""

import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent

EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', '.venv', 'tmp',
    '.gemini', '__pycache__', 'release-artifacts'
}

TARGET_EXTENSIONS = (
    '.md', '.py', '.js', '.mjs', '.cjs', '.ts', '.tsx',
    '.json', '.yml', '.yaml', '.toml', '.sh', '.ps1', '.go', '.php', '.cs'
)

REPLACEMENTS = [
    ('.lovable/coding-guidelines/coding-guidelines.md', '.lovable/coding-guidelines.md'),
    ('.lovable\\coding-guidelines\\coding-guidelines.md', '.lovable\\coding-guidelines.md'),
    ('.lovable/coding-guidelines/', '.lovable/coding-guidelines.md'),
    ('lovable/coding-guidelines/coding-guidelines.md', '.lovable/coding-guidelines.md'),
]


def consolidate_paths() -> int:
    modified_files = []
    for p in REPO_ROOT.rglob('*'):
        if p.is_dir() or any(ex in p.parts for ex in EXCLUDE_DIRS):
            continue
        if not any(p.name.endswith(ext) for ext in TARGET_EXTENSIONS):
            continue
        if p == Path(__file__).resolve():
            continue
        try:
            content = p.read_text(encoding='utf-8')
        except Exception:
            continue

        mod = content
        for old, new in REPLACEMENTS:
            if old in mod:
                mod = mod.replace(old, new)
        if mod != content:
            p.write_text(mod, encoding='utf-8')
            modified_files.append(p.relative_to(REPO_ROOT).as_posix())

    print(f'Consolidated {len(modified_files)} file(s) to .lovable/coding-guidelines.md')
    for f in modified_files:
        print(f'  - {f}')
    return 0


if __name__ == '__main__':
    sys.exit(consolidate_paths())
