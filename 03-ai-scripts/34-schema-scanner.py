#!/usr/bin/env python3
"""34-schema-scanner.py — Scans all SQL table definitions in Go, SQL, and migrations for schema violations."""
from __future__ import annotations

import os
from pathlib import Path
import re
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

ROOT_DIR = Path(__file__).resolve().parent.parent
EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', 'bin', '.gemini', 'release-artifacts'}

RE_CREATE_TABLE = re.compile(r'CREATE TABLE(?:\s+IF NOT EXISTS)?\s+["`]?([A-Za-z0-9_]+)["`]?\s*\((.*?)\)(?:;|\s*`)', re.DOTALL | re.IGNORECASE)
RE_PASCOL_CASE = re.compile(r'^[A-Z][A-Za-z0-9]+$')


def scan_file_tables(file_path: Path) -> list[dict]:
    content = file_path.read_text(encoding='utf-8', errors='replace')
    rel_path = file_path.relative_to(ROOT_DIR).as_posix()
    results = []

    for match in RE_CREATE_TABLE.finditer(content):
        table_name = match.group(1)
        body = match.group(2)
        results.append({
            'file': rel_path,
            'table': table_name,
            'body': body.strip(),
        })

    return results


def is_target_file(file_path: Path) -> bool:
    if file_path.suffix in {'.go', '.sql'}:
        return True

    return False


def collect_all_tables() -> list[dict]:
    all_tables = []
    for root, dirs, files in os.walk(ROOT_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for f in files:
            fp = Path(root) / f
            if is_target_file(fp):
                all_tables.extend(scan_file_tables(fp))

    return all_tables


def inspect_table_violations(table_info: dict) -> list[str]:
    violations = []
    table = table_info['table']
    file = table_info['file']

    # 1. Table name must be PascalCase
    if not RE_PASCOL_CASE.match(table):
        violations.append(f"{file}: Table '{table}' is not PascalCase")

    # 2. Check primary key format
    body = table_info['body']
    expected_pk = f"{table}Id"
    if "PRIMARY KEY" in body.upper():
        if expected_pk.lower() not in body.lower() and not table.endswith("Metadata"):
            violations.append(f"{file}: Table '{table}' primary key does not match {expected_pk}")

    return violations


def main():
    tables = collect_all_tables()
    print(f"Discovered {len(tables)} table definition(s).")
    violations = []
    for t in tables:
        issues = inspect_table_violations(t)
        violations.extend(issues)

    if not violations:
        print("✅ PASS: All database schemas follow PascalCase table and PK rules.")
        sys.exit(0)

    print(f"❌ FAIL: Found {len(violations)} schema violation(s):")
    for v in violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == '__main__':
    main()
