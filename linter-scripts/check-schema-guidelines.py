#!/usr/bin/env python3
"""check-schema-guidelines.py — Linter for database schemas, table casing, PK convention, and affirmative booleans."""
from __future__ import annotations

import os
from pathlib import Path
import re
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent
EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', 'bin', '.gemini', 'release-artifacts'}
ALLOWLIST_TABLES = {'token_service', 'gitmap_metadata', 'chrome_profiles', 'chrome_preferences', 'chrome_bookmarks', 'chrome_extensions', 'chrome_blobs', 'chrome_tokens', 'steps', 'trajectory_metadata_blob', 'test_job'}

RE_CREATE_TABLE = re.compile(r'CREATE TABLE(?:\s+IF NOT EXISTS)?\s+["`]?([A-Za-z0-9_]+)["`]?\s*\((.*?)\)(?:;|\s*`)', re.DOTALL | re.IGNORECASE)
RE_PASCAL_CASE = re.compile(r'^[A-Z][A-Za-z0-9]+$')
RE_BANNED_BOOL = re.compile(r'\b(isNot|hasNo|is_not|has_no)[A-Za-z0-9_]*\b', re.IGNORECASE)


EXEMPT_PK_TABLES = {
    'Setting': 'Key',
    'ClusterNode': 'NodeId',
    'SSHConnection': 'Alias',
    'RepoCGVersion': 'RepoAlias',
    'CloneInteractiveSelection': 'SelectionId',
}


def check_table_casing(rel_path: str, table: str) -> str | None:
    if table in ALLOWLIST_TABLES or table.endswith(('_test', '_v15', '_v15new', '_v15phase5')):
        return None
    if not RE_PASCAL_CASE.match(table):
        return f"{rel_path}: Table '{table}' must be PascalCase singular."

    return None


def check_table_pk(rel_path: str, table: str, body: str) -> str | None:
    if table in ALLOWLIST_TABLES or table.endswith(('_test', '_v15', '_v15new', '_v15phase5', 'Cache', 'Index', 'Log', 'Metadata', 'Item', 'Repo')):
        return None
    if expected_custom_pk := EXEMPT_PK_TABLES.get(table):
        if expected_custom_pk.lower() not in body.lower():
            return f"{rel_path}: Table '{table}' primary key should match {expected_custom_pk}."
        return None

    expected_pk = f"{table}Id"
    if "PRIMARY KEY" in body.upper() and expected_pk.lower() not in body.lower():
        return f"{rel_path}: Table '{table}' primary key should match {expected_pk}."

    return None


def check_boolean_columns(rel_path: str, table: str, body: str) -> list[str]:
    violations = []
    for match in RE_BANNED_BOOL.finditer(body):
        violations.append(f"{rel_path}: Table '{table}' has negative boolean column '{match.group(0)}'.")

    return violations


def audit_sql_content(file_path: Path) -> list[str]:
    content = file_path.read_text(encoding='utf-8', errors='replace')
    rel_path = file_path.relative_to(ROOT_DIR).as_posix()
    violations = []

    for match in RE_CREATE_TABLE.finditer(content):
        table = match.group(1)
        body = match.group(2)
        if case_err := check_table_casing(rel_path, table):
            violations.append(case_err)
        if pk_err := check_table_pk(rel_path, table, body):
            violations.append(pk_err)
        violations.extend(check_boolean_columns(rel_path, table, body))

    return violations


def scan_target_files() -> list[str]:
    target_dirs = [ROOT_DIR / "cli" / "constants", ROOT_DIR / "cli" / "db"]
    violations = []
    for target in target_dirs:
        for root, dirs, files in os.walk(target):
            dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
            for f in files:
                if f.endswith(('.go', '.sql')):
                    violations.extend(audit_sql_content(Path(root) / f))

    return violations


def main() -> None:
    violations = scan_target_files()
    if not violations:
        print("Scanned SQL schemas and table definitions for coding guidelines.")
        print("\n✅ PASS: All database tables, primary keys, and column names meet strict schema guidelines.")
        sys.exit(0)

    print(f"❌ FAIL: Found {len(violations)} schema guideline violation(s):")
    for v in violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == '__main__':
    main()
