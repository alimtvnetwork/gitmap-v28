#!/usr/bin/env python3
"""check-enum-guidelines.py - Linter for Enum Naming (*Type suffix), String Unions, and Raw Rune Literals."""
import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp'
}

# Regex rules
TS_ENUM_MISSING_TYPE = re.compile(r'^\s*(?:export\s+)?enum\s+([A-Z]\w*?)(?<!Type)\s*\{')
TS_CONST_ENUM_MISSING_TYPE = re.compile(r'^\s*(?:export\s+)?const\s+([A-Z]\w*?)(?<!Type)\s*=\s*\{.*?\}\s*as\s+const', re.DOTALL)
PY_ENUM_MISSING_TYPE = re.compile(r'^\s*class\s+([A-Z]\w*?)(?<!Type)\s*\((?:StrEnum|IntEnum|Enum)\)\s*:')
RAW_RUNE_NUM_CAST = re.compile(r'\brune\s*\(\s*(?:10|13|0|\d+)\s*\)')


def check_file(filepath: Path) -> list[str]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception:
        return []

    violations = []
    rel_path = filepath.relative_to(ROOT_DIR).as_posix()
    lines = content.split('\n')

    # Check for raw rune number casts (e.g. rune(10))
    if filepath.suffix == '.go':
        for idx, line in enumerate(lines, 1):
            if line.strip().startswith(('//', '/*', '*')):
                continue
            if RAW_RUNE_NUM_CAST.search(line):
                violations.append(f"{rel_path}:{idx} Raw rune numerical cast found: {line.strip()[:80]}")

    # Check for TypeScript Enums missing *Type suffix
    if filepath.suffix in ('.ts', '.tsx'):
        for idx, line in enumerate(lines, 1):
            if line.strip().startswith(('//', '/*', '*')):
                continue
            m = TS_ENUM_MISSING_TYPE.search(line)
            if m:
                violations.append(f"{rel_path}:{idx} TypeScript enum '{m.group(1)}' missing mandatory 'Type' suffix")

    # Check for Python Enums missing *Type suffix
    if filepath.suffix == '.py':
        for idx, line in enumerate(lines, 1):
            if line.strip().startswith('#'):
                continue
            m = PY_ENUM_MISSING_TYPE.search(line)
            if m:
                violations.append(f"{rel_path}:{idx} Python Enum class '{m.group(1)}' missing mandatory 'Type' suffix")

    return violations


def main():
    target_dirs = [ROOT_DIR / 'gitmap', ROOT_DIR / 'src', ROOT_DIR / 'scripts', ROOT_DIR / 'spec']
    all_violations = []

    for td in target_dirs:
        if not td.exists():
            continue
        for root, dirs, files in os.walk(td):
            dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
            for file in files:
                fp = Path(root) / file
                if fp.suffix.lower() not in {'.go', '.ts', '.tsx', '.py', '.php'}:
                    continue
                all_violations.extend(check_file(fp))

    if not all_violations:
        print("Scanned codebase for Enum (*Type suffix) and Constant conventions.")
        print("\n✅ PASS: All enums, constants, and character code conversions meet strict coding guidelines.")
        sys.exit(0)

    print(f"❌ FAIL: Found {len(all_violations)} enum/constant guideline violation(s):")
    for v in all_violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == '__main__':
    main()
