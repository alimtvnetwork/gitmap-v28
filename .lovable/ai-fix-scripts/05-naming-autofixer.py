#!/usr/bin/env python3
"""05-naming-autofixer.py - Scans and refactors variable/boolean naming violations, anti-ok patterns, and negative booleans."""
import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', '.next', '.gitmap', 'vendor', 'brain', '.gemini'}


def check_bare_ok(is_go: bool, stripped: str, rel_path: str, idx: int) -> dict | None:
    if not is_go:
        return None
    if not re.search(r',\s*ok\s*(:=|=)\s*', stripped):
        return None
    return {
        "file": rel_path,
        "line": idx,
        "type": "bare_ok",
        "snippet": stripped[:80],
        "fix": "Rename bare ok to semantic affirmative boolean (e.g. isAppErr, isFound, hasValue, isSuccess)"
    }


def check_negative_bool(stripped: str, rel_path: str, idx: int) -> dict | None:
    m_neg = re.search(r'\b(hasNo[A-Z]\w*|isNot[A-Z]\w*)\b', stripped)
    if not m_neg:
        return None
    if stripped.startswith(('"', "'")):
        return None
    neg_name = m_neg.group(1)
    return {
        "file": rel_path,
        "line": idx,
        "type": "negative_bool",
        "snippet": stripped[:80],
        "fix": f"Convert negative boolean '{neg_name}' to positive framing"
    }


def check_explicit_bool(stripped: str, rel_path: str, idx: int) -> dict | None:
    if not re.search(r'\b(==|===|!=|!==)\s*(true|false)\b', stripped):
        return None
    return {
        "file": rel_path,
        "line": idx,
        "type": "explicit_bool",
        "snippet": stripped[:80],
        "fix": "Use implicit boolean evaluation"
    }


def scan_file_naming(filepath: Path) -> list[dict]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception:
        return []

    lines = content.split('\n')
    violations = []
    rel_path = filepath.relative_to(ROOT_DIR).as_posix()
    is_go = filepath.suffix == '.go'

    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        if stripped.startswith(('//', '/*', '*', '#')):
            continue

        v_ok = check_bare_ok(is_go, stripped, rel_path, idx)
        if v_ok:
            violations.append(v_ok)

        v_neg = check_negative_bool(stripped, rel_path, idx)
        if v_neg:
            violations.append(v_neg)

        v_exp = check_explicit_bool(stripped, rel_path, idx)
        if v_exp:
            violations.append(v_exp)

    return violations


def process_target_path(td: Path) -> list[dict]:
    if td.is_file():
        return scan_file_naming(td)

    violations = []
    for root, dirs, files in os.walk(td):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for file in files:
            fp = Path(root) / file
            if fp.suffix.lower() not in {'.go', '.ts', '.tsx', '.js', '.jsx', '.py'}:
                continue
            violations.extend(scan_file_naming(fp))
    return violations


def main():
    target_dirs = [ROOT_DIR / "gitmap", ROOT_DIR / "src", ROOT_DIR / "linter-scripts", ROOT_DIR / "scripts"]
    if len(sys.argv) > 1:
        target_dirs = [Path(p) for p in sys.argv[1:]]

    all_violations = []
    for td in target_dirs:
        all_violations.extend(process_target_path(td))

    print(f"=== Naming & Boolean Audit: Scanned files, found {len(all_violations)} violation(s) ===")
    for v in all_violations:
        print(f"  {v['file']}:{v['line']} [{v['type']}] {v['snippet']}")

    return 0 if len(all_violations) == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
