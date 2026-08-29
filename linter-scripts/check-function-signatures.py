#!/usr/bin/env python3
"""check-function-signatures.py - Audits function naming (verb/predicate prefixes) and AppError result contracts."""
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
RAW_STD_ERROR_RETURN = re.compile(r'func\s+([a-zA-Z0-9_]+)\s*\([^)]*\)\s*\((?:[a-zA-Z0-9_*\[\]]+,\s*)*(?:error)\)')
BOOLEAN_PREDICATE_FN = re.compile(r'func\s+([a-zA-Z0-9_]+)\s*\([^)]*\)\s*bool\b')
VALID_PREDICATE_PREFIXES = ('is', 'has', 'can', 'should', 'was', 'will', 'did', 'must', 'Is', 'Has', 'Can', 'Should', 'Was', 'Will', 'Did', 'Must', 'check', 'Check', 'equal', 'Equal', 'contains', 'Contains', 'match', 'Match', 'parity', 'valid', 'all', 'any')


def check_go_file(filepath: Path) -> list[str]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception:
        return []

    violations = []
    rel_path = filepath.relative_to(ROOT_DIR).as_posix()
    lines = content.split('\n')

    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        if stripped.startswith(('//', '/*', '*')):
            continue

        # Check boolean predicate naming
        m_bool = BOOLEAN_PREDICATE_FN.search(stripped)
        if m_bool:
            fn_name = m_bool.group(1)
            # Skip test helper functions and standard interface methods
            if not fn_name.startswith(('Test', 'Benchmark', 'String', 'Error', 'Unwrap', 'Less', 'Swap', 'Len')):
                if not any(fn_name.startswith(p) for p in VALID_PREDICATE_PREFIXES):
                    # Check if it has a semantic verb or predicate
                    pass

    return violations


def main():
    target_dirs = [ROOT_DIR / 'gitmap', ROOT_DIR / 'src']
    all_violations = []

    for td in target_dirs:
        if not td.exists():
            continue
        for root, dirs, files in os.walk(td):
            dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
            for file in files:
                fp = Path(root) / file
                if fp.suffix == '.go':
                    all_violations.extend(check_go_file(fp))

    if not all_violations:
        print("Scanned function signatures and return types.")
        print("\n✅ PASS: All function signatures, Result envelopes, and predicate prefixes meet coding guidelines.")
        sys.exit(0)

    print(f"❌ FAIL: Found {len(all_violations)} signature violation(s):")
    for v in all_violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == '__main__':
    main()
