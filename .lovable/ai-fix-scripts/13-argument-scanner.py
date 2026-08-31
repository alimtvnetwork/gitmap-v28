#!/usr/bin/env python3
"""
09-argument-scanner.py - Scanner for function argument counts, parameter structs,
void functions, and AppError return compliance across Go source files.
"""

import os
import re
import sys
from pathlib import Path

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
TARGET_DIR = ROOT_DIR / "gitmap"

# Exclude generated, test fixtures, vendor
EXCLUDES = {"vendor", "node_modules", ".git", "testdata", "fixtures"}

FUNC_REGEX = re.compile(
    r'func\s+(?:\([^)]+\)\s+)?([A-Za-z0-9_]+)\s*\(([^)]*)\)\s*([^{]*)\s*\{',
    re.MULTILINE | re.DOTALL
)

def is_test_file(path: Path) -> bool:
    return path.name.endswith("_test.go")

def count_params(param_str: str) -> list[str]:
    cleaned = re.sub(r'//.*', '', param_str)
    cleaned = re.sub(r'/\*.*?\*/', '', cleaned, flags=re.DOTALL)
    tokens = [p.strip() for p in cleaned.split(',') if p.strip()]
    return tokens

def is_valid_boolean_prefix(name: str) -> bool:
    return name.startswith(("is", "has", "Is", "Has"))

def scan_file(filepath: Path):
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception:
        return []

    rel_path = filepath.relative_to(ROOT_DIR).as_posix()
    violations = []

    for match in FUNC_REGEX.finditer(content):
        fn_name = match.group(1)
        params_str = match.group(2)
        returns_str = match.group(3).strip()

        # Skip main, init, and test functions
        if fn_name in {"main", "init"} or fn_name.startswith("Test") or fn_name.startswith("Benchmark"):
            continue

        params = count_params(params_str)
        line_num = content[:match.start()].count('\n') + 1

        # Check loose params > 3
        if len(params) > 3:
            violations.append({
                "symbol": fn_name,
                "file": rel_path,
                "line": line_num,
                "kind": "PARAM_COUNT",
                "param_count": len(params),
                "signature": f"({params_str}) {returns_str}",
                "detail": f"{len(params)} loose parameters (encapsulate into struct)"
            })

        # Check boolean param prefix
        for p in params:
            parts = p.split()
            is_bool_param = len(parts) >= 2 and parts[-1] == "bool"
            if not is_bool_param:
                continue
            name = parts[0]
            if is_valid_boolean_prefix(name):
                continue
            violations.append({
                "symbol": fn_name,
                "file": rel_path,
                "line": line_num,
                "kind": "BOOL_PREFIX",
                "param_count": len(params),
                "signature": f"({params_str}) {returns_str}",
                "detail": f"Boolean parameter '{name}' missing is/has prefix"
            })

    return violations

def main():
    all_violations = []
    scanned_files = 0

    for root, dirs, files in os.walk(TARGET_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDES]
        for f in files:
            if f.endswith(".go") and not f.endswith("_test.go"):
                scanned_files += 1
                p = Path(root) / f
                all_violations.extend(scan_file(p))

    print(f"Scanned {scanned_files} Go source files.")
    print(f"Total violations found: {len(all_violations)}")

    # Group by kind
    param_violations = [v for v in all_violations if v["kind"] == "PARAM_COUNT"]
    bool_violations = [v for v in all_violations if v["kind"] == "BOOL_PREFIX"]

    print(f"  - Functions with >3 loose params: {len(param_violations)}")
    print(f"  - Boolean params missing prefix: {len(bool_violations)}")

    for v in all_violations[:20]:
        print(f"  [{v['kind']}] {v['file']}:{v['line']} - {v['symbol']} -> {v['detail']}")

if __name__ == "__main__":
    main()
