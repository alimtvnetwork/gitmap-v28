#!/usr/bin/env python3
"""
31-db-migration-runner.py — Standalone external migration runner for SQLite and relational schemas.

Usage:
  python 03-ai-scripts/31-db-migration-runner.py --db data/gitmap.db --status
  python 03-ai-scripts/31-db-migration-runner.py --db data/pipeline.db --dry-run
  python 03-ai-scripts/31-db-migration-runner.py --db data/pipeline.db --generate-enums
"""

from __future__ import annotations

import argparse
from importlib import import_module
from pathlib import Path
import sqlite3
import subprocess
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType


def inspect_db_schema(db_path: Path) -> dict[str, list[dict]]:
    """Inspects all tables and columns in a SQLite database."""
    if not db_path.is_file():
        return {}

    conn = sqlite3.connect(str(db_path))
    cursor = conn.cursor()

    cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
    tables = [row[0] for row in cursor.fetchall()]

    schema = {}
    for table in tables:
        cursor.execute(f'PRAGMA table_info("{table}");')
        columns = []
        for row in cursor.fetchall():
            columns.append({
                "cid": row[0],
                "name": row[1],
                "type": row[2],
                "notnull": bool(row[3]),
                "dflt_value": row[4],
                "pk": bool(row[5]),
            })
        schema[table] = columns

    conn.close()
    return schema


def print_db_status(db_path: Path) -> None:
    """Prints schema status and tables for the target database."""
    print(f"Database: {db_path}")
    if not db_path.is_file():
        print("  Status: File does not exist yet (pending initialization).")
        return

    schema = inspect_db_schema(db_path)
    if not schema:
        print("  Status: Connected (empty database, no user tables).")
        return

    print(f"  Status: Healthy ({len(schema)} tables found)")
    for table, cols in schema.items():
        pk_cols = [c["name"] for c in cols if c["pk"]]
        pk_str = f" [PK: {', '.join(pk_cols)}]" if pk_cols else ""
        print(f"  - Table: {table} ({len(cols)} columns){pk_str}")
        for c in cols:
            flags = []
            if c["pk"]:
                flags.append("PRIMARY KEY")
            if c["notnull"]:
                flags.append("NOT NULL")
            flag_str = f" ({', '.join(flags)})" if flags else ""
            print(f"      • {c['name']}: {c['type']}{flag_str}")


def run_migration_script(db_path: Path, sql_script: str, dry_run: bool = False) -> bool:
    """Executes a SQL migration script against the target SQLite database."""
    if dry_run:
        print(f"[DRY RUN] Would execute on {db_path}:")
        for line in sql_script.strip().splitlines():
            print(f"  | {line}")
        return True

    db_path.parent.mkdir(parents=True, exist_ok=True)
    conn = sqlite3.connect(str(db_path))
    cursor = conn.cursor()

    try:
        cursor.executescript(sql_script)
        conn.commit()
        print(f"  ✔ Applied migration to {db_path}")
        return True
    except sqlite3.Error as e:
        print(f"  ✖ Migration failed: {e}", file=sys.stderr)
        conn.rollback()
        return False
    finally:
        conn.close()


def main() -> int:
    parser = argparse.ArgumentParser(description="Standalone external SQLite migration runner")
    parser.add_argument("--db", required=True, help="Path to SQLite database file")
    parser.add_argument("--db-type", default="sqlite", choices=["sqlite", "postgres", "mysql", "mariadb", "mssql", "oracle", "mongodb"], help="Target database type (defaults to sqlite)")
    parser.add_argument("--status", action="store_true", help="Inspect and display database status")
    parser.add_argument("--dry-run", action="store_true", help="Preview migrations without writing")
    parser.add_argument("--generate-enums", action="store_true", help="Trigger enum generator after migration")
    parser.add_argument("--sql", default="", help="Inline SQL script to execute")

    args = parser.parse_args()
    repo_root = Path(__file__).resolve().parent.parent

    db_path = Path(args.db)
    if not db_path.is_absolute():
        db_path = repo_root / db_path

    if args.status:
        print(f"Engine Dialect: {args.db_type}")
        print_db_status(db_path)
        return 0

    if args.sql:
        success = run_migration_script(db_path, args.sql, args.dry_run)
        if not success:
            return 1

    if args.generate_enums:
        generator_script = repo_root / "03-ai-scripts" / "30-db-struct-enum-generator.py"
        if generator_script.is_file():
            print("Running enum code generator...")
            cmd = [sys.executable, str(generator_script)]
            if args.dry_run:
                cmd.append("--dry-run")
            res = subprocess.run(cmd, cwd=str(repo_root))
            if res.returncode != 0:
                return res.returncode

    print_db_status(db_path)
    return 0


if __name__ == "__main__":
    sys.exit(main())
