#!/usr/bin/env python3
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent

EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', '.venv', 'tmp', 'release-artifacts'}
TARGET_EXTENSIONS = (
    '.md', '.py', '.js', '.mjs', '.cjs', '.ts', '.tsx',
    '.json', '.yml', '.yaml', '.toml', '.sh', '.ps1', '.go', '.php', '.cs'
)

PATTERN_SLASH = re.compile(r'(?<![a-zA-Z0-9_-])spec/((?:[0-9]{2}-|[a-zA-Z0-9_.-]+/)[a-zA-Z0-9_.-]*)')
PATTERN_BACKSLASH = re.compile(r'(?<![a-zA-Z0-9_-])spec\\\\((?:[0-9]{2}-|[a-zA-Z0-9_.-]+\\\\)[a-zA-Z0-9_.-]*)')

REPLACEMENTS = [
    ('"spec"', '"02-spec"'),
    ('spec/01-index.md', '02-spec/01-index.md'),
    ('spec/spec-index.md', '02-spec/spec-index.md'),
    ('spec/health-dashboard.md', '02-spec/health-dashboard.md'),
    ('spec/dashboard-data.json', '02-spec/dashboard-data.json'),
    ('spec/folder-structure-root.md', '02-spec/folder-structure-root.md'),
    ('spec/99-consistency-report.md', '02-spec/99-consistency-report.md'),
    ('spec/02-_template.md', '02-spec/02-_template.md'),
]

def migrate() -> int:
    modified_files = []
    for p in REPO_ROOT.rglob('*'):
        if p.is_dir() or any(ex in p.parts for ex in EXCLUDE_DIRS):
            continue
        if not any(p.name.endswith(ext) for ext in TARGET_EXTENSIONS):
            continue
        if p == Path(__file__).resolve():
            continue
        try:
            txt = p.read_text(encoding='utf-8')
        except Exception:
            continue
        mod = PATTERN_SLASH.sub(r'02-spec/\1', txt)
        mod = PATTERN_BACKSLASH.sub(r'02-spec\\\\\1', mod)
        for old, new in REPLACEMENTS:
            if old in mod:
                mod = mod.replace(old, new)
        if mod != txt:
            p.write_text(mod, encoding='utf-8')
            modified_files.append(p.relative_to(REPO_ROOT).as_posix())

    print(f'Migrated spec/ -> 02-spec/ in {len(modified_files)} files.')
    for f in modified_files[:25]:
        print(f'  - {f}')
    return 0

if __name__ == '__main__':
    sys.exit(migrate())
