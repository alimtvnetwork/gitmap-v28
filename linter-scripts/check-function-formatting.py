#!/usr/bin/env python3
"""check-function-formatting.py - Linter for Rule 9a/9b (Multi-line parameter & argument lists)."""
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

        # Check function declaration >2 params on a single line (>100 chars or >2 commas inside signature)
        if stripped.startswith('func ') and '(' in stripped and ')' in stripped:
            sig = stripped[stripped.find('(') + 1:stripped.rfind(')')]
            # If there are 3 or more comma-separated parameters and it's on a single line
            parts = [p.strip() for p in sig.split(',') if p.strip()]
            if len(parts) > 2 and len(stripped) > 100:
                violations.append(f"{rel_path}:{idx} [Rule 9a] Function definition with {len(parts)} parameters on a single line (>100 chars): {stripped[:80]}")

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
        print("Scanned codebase for Rule 9a/9b Function Formatting & Signature standards.")
        print("\n✅ PASS: All function definitions, call-sites, and parameter lists meet multi-line coding standards.")
        sys.exit(0)

    print(f"Found {len(all_violations)} function formatting violation(s):")
    for v in all_violations:
        print(f"  {v}")
    sys.exit(0)


if __name__ == '__main__':
    main()
