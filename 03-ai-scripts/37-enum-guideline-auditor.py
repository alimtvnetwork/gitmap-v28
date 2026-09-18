#!/usr/bin/env python3
"""
37-enum-guideline-auditor.py - Audits repository for:
1. Enums missing '*Type' suffix (Go, TypeScript, Python, PHP).
2. Raw rune numerical casts (e.g. rune(10), rune(13), rune(0)).
3. TypeScript string unions used for status / state.
"""

from importlib import import_module
import os
from pathlib import Path
import re
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
LINE_SEPARATOR = engine.LINE_SEPARATOR
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR

EXCLUDED_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts',
    '.ai-memory/scratch', '.ai-memory/temp-agents', '03-ai-scripts', 'scripts',
    '04-code'
}

RAW_RUNE_CAST = re.compile(r'\brune\s*\(\s*(?:\d+|\'[^\']+\'\s*\+\s*[^)]+)\s*\)')
GO_TYPE_DEF = re.compile(r'^\s*type\s+([A-Za-z]\w*?)(?<!Type)\s+(?:string|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|byte)\b', re.MULTILINE)
TS_ENUM_DEF = re.compile(r'^\s*(?:export\s+)?enum\s+([A-Z]\w*?)(?<!Type)\s*\{', re.MULTILINE)
TS_CONST_OBJ = re.compile(r'^\s*(?:export\s+)?const\s+([A-Z]\w*?)(?<!Type)\s*=\s*\{.*?\}\s*as\s+const', re.DOTALL | re.MULTILINE)
PY_ENUM_DEF = re.compile(r'^\s*class\s+([A-Z]\w*?)(?<!Type)\s*\((?:StrEnum|IntEnum|Enum)\)\s*:', re.MULTILINE)
TS_STRING_UNION = re.compile(r'^\s*(?:export\s+)?type\s+(\w*(?:Status|State|Kind|Mode|Action))\s*=\s*(?:[\'"][^\'"]+[\'"]\s*\|\s*)+', re.MULTILINE)

def audit_go_file(file_path: Path, content: str) -> list[tuple[int, str, str]]:
    violations = []
    lines = content.split('\n')

    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        if stripped.startswith(('//', '/*', '*')):
            continue
        if RAW_RUNE_CAST.search(line):
            violations.append((idx, "RAW_RUNE_CAST", f"Raw numeric rune cast: {stripped[:80]}"))

    for m in GO_TYPE_DEF.finditer(content):
        tname = m.group(1)
        if re.search(rf'\b{tname}\s*=\s*(?:iota|")', content) or re.search(rf'\bconst\s+.*?\b{tname}\b', content, re.DOTALL):
            line_no = content[:m.start()].count('\n') + 1
            violations.append((line_no, "GO_ENUM_MISSING_TYPE", f"Go enum '{tname}' missing mandatory 'Type' suffix"))

    return violations

def audit_ts_file(file_path: Path, content: str) -> list[tuple[int, str, str]]:
    violations = []
    for m in TS_ENUM_DEF.finditer(content):
        line_no = content[:m.start()].count('\n') + 1
        violations.append((line_no, "TS_ENUM_MISSING_TYPE", f"TypeScript enum '{m.group(1)}' missing 'Type' suffix"))

    for m in TS_STRING_UNION.finditer(content):
        line_no = content[:m.start()].count('\n') + 1
        violations.append((line_no, "TS_STRING_UNION", f"TypeScript string union enum '{m.group(1)}' should be a typed enum"))

    return violations

def audit_py_file(file_path: Path, content: str) -> list[tuple[int, str, str]]:
    violations = []
    for m in PY_ENUM_DEF.finditer(content):
        line_no = content[:m.start()].count('\n') + 1
        violations.append((line_no, "PY_ENUM_MISSING_TYPE", f"Python Enum class '{m.group(1)}' missing 'Type' suffix"))

    return violations

def run_auditor(root_dir: str = CURRENT_DIR) -> int:
    root = Path(root_dir).resolve()
    target_dirs = [root / 'cli', root / 'src', root / 'scripts']
    all_violations = []

    for td in target_dirs:
        if not td.exists():
            continue
        for r, dirs, files in os.walk(td):
            dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS]
            for file in files:
                fp = Path(r) / file
                try:
                    content = fp.read_text(encoding='utf-8', errors='replace')
                except Exception:
                    continue

                rel_path = fp.relative_to(root).as_posix()
                vios = []
                if fp.suffix == '.go':
                    vios = audit_go_file(fp, content)
                elif fp.suffix in ('.ts', '.tsx'):
                    vios = audit_ts_file(fp, content)
                elif fp.suffix == '.py':
                    vios = audit_py_file(fp, content)

                for lno, kind, desc in vios:
                    all_violations.append((rel_path, lno, kind, desc))

    if all_violations:
        print(f"{LINE_SEPARATOR}❌ Found {len(all_violations)} enum/constant guideline violation(s):")
        for rel_path, lno, kind, desc in all_violations:
            print(f"  {rel_path}:{lno} [{kind}] {desc}")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ All code files in '{root_dir}' conform to enum and constant guidelines.")
    return ExitCodeType.SUCCESS.value

if __name__ == '__main__':
    sys.exit(run_auditor())
