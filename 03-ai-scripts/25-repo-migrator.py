#!/usr/bin/env python3
"""
Transactional Repository Layout Migrator & History Engine
Migrates legacy repository layouts (spec/, .lovable/prompts/, .lovable/ai-fix-scripts/)
to the standardized root layout (01-prompts/, 02-spec/, 03-ai-scripts/).

Features:
  - SQLite Transaction Journal (tmp/migrations.db) for 100% reversible operations.
  - Plan Mode First (--plan / --dry-run) previewing all actions.
  - Complete Undo / Redo capabilities with full content snapshots.
  - CLI History & Inspection: `ls`, `history`, `show <id>`, `undo`, `redo`.
  - Non-destructive safe handling: old content snapshots preserved in SQLite and Trash Bin.

Usage Examples:
  # 1. Preview (Plan Mode) migration on target repo:
  python 03-ai-scripts/25-repo-migrator.py migrate --path ./my-project --plan

  # 2. Execute migration with transaction logging:
  python 03-ai-scripts/25-repo-migrator.py migrate --path ./my-project --force

  # 3. List all past migration transactions:
  python 03-ai-scripts/25-repo-migrator.py ls

  # 4. Show details of a specific transaction:
  python 03-ai-scripts/25-repo-migrator.py show tx-20260901-120000

  # 5. Undo the latest (or specific) migration transaction:
  python 03-ai-scripts/25-repo-migrator.py undo

  # 6. Redo a previously reverted transaction:
  python 03-ai-scripts/25-repo-migrator.py redo
"""

from __future__ import annotations

import argparse
import datetime
import fnmatch
from importlib import import_module
import json
import os
from pathlib import Path
import re
import shutil
import sqlite3
import subprocess
import sys
import time
import uuid

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

# Centralized Enums & Literals
ExitCodeType = engine.ExitCodeType
CURRENT_DIR = engine.CURRENT_DIR
LINE_SEPARATOR = engine.LINE_SEPARATOR
DEFAULT_ENCODING = engine.DEFAULT_ENCODING

DB_DIR = Path("tmp")
DB_PATH = DB_DIR / "migrations.db"

EXCLUDE_DIRS = {
    ".git", "node_modules", "dist", "build", ".venv", "tmp", ".tmp",
    ".gemini", "__pycache__", "release-artifacts"
}

TARGET_EXTENSIONS = (
    ".md", ".py", ".js", ".mjs", ".cjs", ".ts", ".tsx",
    ".json", ".yml", ".yaml", ".toml", ".sh", ".ps1", ".go", ".php", ".cs"
)

# Core Patterns for Content Rewrites
PATTERN_SPEC_SLASH = re.compile(r"(?<![a-zA-Z0-9_-])spec/((?:[0-9]{2}-|[a-zA-Z0-9_.-]+/)[a-zA-Z0-9_.-]*)")
PATTERN_SPEC_BACKSLASH = re.compile(r"(?<![a-zA-Z0-9_-])spec\\\\((?:[0-9]{2}-|[a-zA-Z0-9_.-]+\\\\)[a-zA-Z0-9_.-]*)")

DIRECT_STRING_REPLACEMENTS = [
    (".lovable/prompts/01-prompts-category/", "01-prompts/"),
    (".lovable/prompts/", "01-prompts/"),
    (".lovable\\prompts\\01-prompts-category\\", "01-prompts\\\\"),
    (".lovable\\prompts\\", "01-prompts\\\\"),
    (".lovable/ai-fix-scripts/", "03-ai-scripts/"),
    (".lovable\\ai-fix-scripts\\", "03-ai-scripts\\\\"),
    ("lovable/ai-fix-scripts/", "03-ai-scripts/"),
    (".lovable/coding-guidelines/coding-guidelines.md", ".lovable/coding-guidelines.md"),
    (".lovable\\coding-guidelines\\coding-guidelines.md", ".lovable\\coding-guidelines.md"),
    (".lovable/coding-guidelines/", ".lovable/coding-guidelines.md"),
    ("\"spec\"", "\"02-spec\""),
    ("spec/01-index.md", "02-spec/01-index.md"),
    ("spec/spec-index.md", "02-spec/spec-index.md"),
    ("spec/health-dashboard.md", "02-spec/health-dashboard.md"),
    ("spec/dashboard-data.json", "02-spec/dashboard-data.json"),
    ("spec/folder-structure-root.md", "02-spec/folder-structure-root.md"),
    ("spec/99-consistency-report.md", "02-spec/99-consistency-report.md"),
    ("spec/02-_template.md", "02-spec/02-_template.md"),
    ("`spec/`", "`02-spec/`"),
    ("`spec`", "`02-spec`"),
]


def init_db(db_path: Path = DB_PATH) -> sqlite3.Connection:
    """Initializes SQLite database and tables for transaction logging."""
    db_path.parent.mkdir(parents=True, exist_ok=True)
    conn = sqlite3.connect(str(db_path))
    conn.execute("PRAGMA journal_mode = WAL;")
    with conn:
        conn.execute("""
            CREATE TABLE IF NOT EXISTS transactions (
                transaction_id TEXT PRIMARY KEY,
                timestamp TEXT NOT NULL,
                target_repo TEXT NOT NULL,
                operation_count INTEGER NOT NULL,
                status TEXT NOT NULL, -- 'APPLIED', 'REVERTED'
                description TEXT NOT NULL
            );
        """)
        conn.execute("""
            CREATE TABLE IF NOT EXISTS operations (
                op_id INTEGER PRIMARY KEY AUTOINCREMENT,
                transaction_id TEXT NOT NULL,
                seq_num INTEGER NOT NULL,
                op_type TEXT NOT NULL, -- 'MOVE_DIR', 'MOVE_FILE', 'MODIFY_FILE', 'CREATE_FILE', 'TRASH_FILE'
                source_path TEXT NOT NULL,
                target_path TEXT,
                original_content TEXT,
                new_content TEXT,
                status TEXT NOT NULL,
                FOREIGN KEY(transaction_id) REFERENCES transactions(transaction_id)
            );
        """)
    return conn


def normalize_rel_path(path_obj: Path, base_dir: Path) -> str:
    """Returns normalized forward-slash relative path."""
    try:
        rel = path_obj.resolve().relative_to(base_dir.resolve())
        return str(rel).replace("\\", "/")
    except Exception:
        return str(path_obj).replace("\\", "/")


def send_to_trash_bin(target_path: Path) -> bool:
    """Safely moves a file or folder to the OS Recycle/Trash Bin."""
    if not target_path.exists():
        return False
    try:
        import send2trash
        send2trash.send2trash(str(target_path.resolve()))
        return True
    except ImportError:
        pass
    except Exception:
        pass

    if sys.platform == "win32":
        try:
            abs_str = str(target_path.resolve()).replace("\\", "\\\\")
            if target_path.is_dir():
                ps_cmd = f"Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteDirectory('{abs_str}', 'OnlyErrorDialogs', 'SendToRecycleBin')"
            else:
                ps_cmd = f"Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile('{abs_str}', 'OnlyErrorDialogs', 'SendToRecycleBin')"
            cmd = ["powershell", "-NoProfile", "-NonInteractive", "-Command", ps_cmd]
            res = subprocess.run(cmd, capture_output=True, text=True)
            if res.returncode == 0 and not target_path.exists():
                return True
        except Exception:
            pass

    if target_path.is_dir():
        shutil.rmtree(target_path)
    else:
        target_path.unlink()
    return True


class MigrationPlanner:
    """Discovers planned operations required to migrate a repository."""

    def __init__(self, repo_root: Path):
        self.repo_root = repo_root.resolve()
        self.moves: list[tuple[Path, Path, str]] = []  # (src, dst, type)
        self.modifications: list[tuple[Path, str, str]] = []  # (path, old_content, new_content)
        self.stale_files: list[Path] = []

    def plan(self) -> None:
        # 1. Plan Directory & Key File Moves
        prompts_cat = self.repo_root / ".lovable" / "prompts" / "01-prompts-category"
        prompts_root_dir = self.repo_root / ".lovable" / "prompts"
        target_prompts = self.repo_root / "01-prompts"

        if prompts_cat.exists() and not target_prompts.exists():
            self.moves.append((prompts_cat, target_prompts, "MOVE_DIR"))
        elif prompts_root_dir.exists() and not target_prompts.exists():
            self.moves.append((prompts_root_dir, target_prompts, "MOVE_DIR"))

        spec_dir = self.repo_root / "spec"
        target_spec = self.repo_root / "02-spec"
        if spec_dir.exists() and not target_spec.exists():
            self.moves.append((spec_dir, target_spec, "MOVE_DIR"))

        ai_scripts = self.repo_root / ".lovable" / "ai-fix-scripts"
        target_ai_scripts = self.repo_root / "03-ai-scripts"
        if ai_scripts.exists() and not target_ai_scripts.exists():
            self.moves.append((ai_scripts, target_ai_scripts, "MOVE_DIR"))

        nested_guidelines = self.repo_root / ".lovable" / "coding-guidelines" / "coding-guidelines.md"
        target_guidelines = self.repo_root / ".lovable" / "coding-guidelines.md"
        if nested_guidelines.exists() and not target_guidelines.exists():
            self.moves.append((nested_guidelines, target_guidelines, "MOVE_FILE"))

        # 2. Plan Stale Items Cleanup
        nested_cg_dir = self.repo_root / ".lovable" / "coding-guidelines"
        if nested_cg_dir.exists():
            self.stale_files.append(nested_cg_dir)

        # 3. Plan Content Rewrites Across Target Files
        for p in self.repo_root.rglob("*"):
            if p.is_dir() or any(ex in p.parts for ex in EXCLUDE_DIRS):
                continue
            if not any(p.name.endswith(ext) for ext in TARGET_EXTENSIONS):
                continue
            if p.name == "migrations.db":
                continue
            try:
                original_txt = p.read_text(encoding="utf-8")
            except Exception:
                continue

            mod = PATTERN_SPEC_SLASH.sub(r"02-spec/\1", original_txt)
            mod = PATTERN_SPEC_BACKSLASH.sub(r"02-spec\\\\\1", mod)
            for old, new in DIRECT_STRING_REPLACEMENTS:
                if old in mod:
                    mod = mod.replace(old, new)

            if mod != original_txt:
                self.modifications.append((p, original_txt, mod))


def execute_migration(repo_root: Path, is_plan_only: bool = False, is_force: bool = False) -> int:
    """Executes or previews transactional repository migration."""
    planner = MigrationPlanner(repo_root)
    planner.plan()

    has_actions = bool(planner.moves or planner.modifications or planner.stale_files)
    if not has_actions:
        print("[PASS] Repository is already aligned with the latest folder structure. Nothing to migrate.")
        return ExitCodeType.SUCCESS.value

    # Display Plan
    print("=" * 80)
    print(f"[PLAN] Transactional Migration Plan for: {repo_root}")
    print(f"   Directory / File Moves: {len(planner.moves)}")
    print(f"   File Content Rewrites:  {len(planner.modifications)}")
    print(f"   Stale Items Cleanup:    {len(planner.stale_files)}")
    print("=" * 80)

    if planner.moves:
        print("\n[MOVES] Planned Folder/File Moves:")
        for src, dst, kind in planner.moves:
            print(f"  * [{kind}] {normalize_rel_path(src, repo_root)} -> {normalize_rel_path(dst, repo_root)}")

    if planner.stale_files:
        print("\n[CLEAN] Stale Displaced Folders (Will be moved to Trash):")
        for sf in planner.stale_files:
            print(f"  * {normalize_rel_path(sf, repo_root)}")

    if planner.modifications:
        print(f"\n[UPDATES] Planned File Content Updates ({len(planner.modifications)} files):")
        for p, _, _ in planner.modifications[:15]:
            print(f"  * {normalize_rel_path(p, repo_root)}")
        if len(planner.modifications) > 15:
            print(f"  ... and {len(planner.modifications) - 15} more files.")

    if is_plan_only:
        print(f"{LINE_SEPARATOR}[INFO] Plan Mode complete. No files were modified.")
        return ExitCodeType.SUCCESS.value

    if not is_force:
        try:
            resp = input(f"{LINE_SEPARATOR}Proceed with executing this transactional migration? (y/N): ").strip().lower()
            if resp not in {"y", "yes"}:
                print("[ABORT] Migration canceled by user. No files were modified.")
                return ExitCodeType.SUCCESS.value
        except (EOFError, KeyboardInterrupt):
            print("\n[ABORT] Operation canceled.")
            return ExitCodeType.SUCCESS.value

    # Start Transaction
    tx_id = f"tx-{datetime.datetime.now().strftime('%Y%m%d-%H%M%S')}-{uuid.uuid4().hex[:6]}"
    conn = init_db()
    seq = 0
    total_ops = len(planner.moves) + len(planner.modifications) + len(planner.stale_files)

    print(f"\n[EXEC] Executing migration under Transaction ID: {tx_id}...")
    start_time = time.perf_counter()

    try:
        with conn:
            conn.execute("""
                INSERT INTO transactions (transaction_id, timestamp, target_repo, operation_count, status, description)
                VALUES (?, ?, ?, ?, ?, ?);
            """, (tx_id, datetime.datetime.now().isoformat(), str(repo_root), total_ops, "APPLIED",
                  f"Migrated repository to 01-prompts, 02-spec, 03-ai-scripts layout"))

            # 1. Execute Directory & File Moves
            for src, dst, kind in planner.moves:
                seq += 1
                dst.parent.mkdir(parents=True, exist_ok=True)
                shutil.move(str(src), str(dst))
                conn.execute("""
                    INSERT INTO operations (transaction_id, seq_num, op_type, source_path, target_path, original_content, new_content, status)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?);
                """, (tx_id, seq, kind, str(src), str(dst), None, None, "APPLIED"))
                print(f"  * Moved: {normalize_rel_path(src, repo_root)} -> {normalize_rel_path(dst, repo_root)}")

            # 2. Execute File Modifications
            for p, old_txt, new_txt in planner.modifications:
                seq += 1
                p.write_text(new_txt, encoding="utf-8")
                conn.execute("""
                    INSERT INTO operations (transaction_id, seq_num, op_type, source_path, target_path, original_content, new_content, status)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?);
                """, (tx_id, seq, "MODIFY_FILE", str(p), str(p), old_txt, new_txt, "APPLIED"))

            # 3. Clean Stale Items
            for sf in planner.stale_files:
                if sf.exists():
                    seq += 1
                    send_to_trash_bin(sf)
                    conn.execute("""
                        INSERT INTO operations (transaction_id, seq_num, op_type, source_path, target_path, original_content, new_content, status)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?);
                    """, (tx_id, seq, "TRASH_FILE", str(sf), None, None, None, "APPLIED"))
                    print(f"  * Cleaned: {normalize_rel_path(sf, repo_root)}")

        elapsed_ms = (time.perf_counter() - start_time) * 1000
        print(f"\n[PASS] Migration successful in {elapsed_ms:.2f}ms! Transaction logged in {DB_PATH}.")
        print(f"   To revert at any time, run: python 03-ai-scripts/25-repo-migrator.py undo --tx-id {tx_id}")
        return ExitCodeType.SUCCESS.value

    except Exception as err:
        conn.rollback()
        print(f"\n[FAIL] Migration failed: {err}")
        return ExitCodeType.GENERIC_FAILURE.value


def list_transactions() -> int:
    """Lists all past recorded migration transactions."""
    if not DB_PATH.exists():
        print("[INFO] No migration transactions found (tmp/migrations.db does not exist).")
        return ExitCodeType.SUCCESS.value

    conn = init_db()
    cur = conn.cursor()
    cur.execute("SELECT transaction_id, timestamp, operation_count, status, description FROM transactions ORDER BY timestamp DESC;")
    rows = cur.fetchall()

    if not rows:
        print("[INFO] No migration transactions recorded yet.")
        return ExitCodeType.SUCCESS.value

    print("=" * 95)
    print(f"{'TRANSACTION ID':<32} {'TIMESTAMP':<20} {'OPS':<6} {'STATUS':<12} {'DESCRIPTION'}")
    print("=" * 95)
    for tx_id, ts, ops, status, desc in rows:
        ts_clean = ts[:19].replace("T", " ")
        status_badge = "[APPLIED]" if status == "APPLIED" else "[REVERTED]"
        print(f"{tx_id:<32} {ts_clean:<20} {ops:<6} {status_badge:<12} {desc}")
    print("=" * 95)
    return ExitCodeType.SUCCESS.value


def show_transaction(tx_id: str) -> int:
    """Displays detailed file-by-file operations inside a specific transaction."""
    if not DB_PATH.exists():
        print("[ERROR] Database tmp/migrations.db does not exist.")
        return ExitCodeType.GENERIC_FAILURE.value

    conn = init_db()
    cur = conn.cursor()
    cur.execute("SELECT transaction_id, timestamp, target_repo, operation_count, status, description FROM transactions WHERE transaction_id = ?;", (tx_id,))
    tx = cur.fetchone()
    if not tx:
        print(f"[ERROR] Transaction not found: {tx_id}")
        return ExitCodeType.GENERIC_FAILURE.value

    print("=" * 80)
    print(f"[TRANSACTION] Details: {tx[0]}")
    print(f"   Timestamp:   {tx[1]}")
    print(f"   Target Repo: {tx[2]}")
    print(f"   Operations:  {tx[3]}")
    print(f"   Status:      {tx[4]}")
    print(f"   Description: {tx[5]}")
    print("=" * 80)

    cur.execute("SELECT seq_num, op_type, source_path, target_path, status FROM operations WHERE transaction_id = ? ORDER BY seq_num ASC;", (tx_id,))
    ops = cur.fetchall()
    for seq, op_type, src, dst, st in ops:
        dst_str = f" -> {dst}" if dst and dst != src else ""
        print(f"  {seq:>3}. [{op_type}] {src}{dst_str} ({st})")
    return ExitCodeType.SUCCESS.value


def undo_transaction(tx_id: str | None = None) -> int:
    """Reverts an applied migration transaction back to its previous state."""
    if not DB_PATH.exists():
        print("[ERROR] Database tmp/migrations.db does not exist.")
        return ExitCodeType.GENERIC_FAILURE.value

    conn = init_db()
    cur = conn.cursor()

    if not tx_id:
        cur.execute("SELECT transaction_id FROM transactions WHERE status = 'APPLIED' ORDER BY timestamp DESC LIMIT 1;")
        row = cur.fetchone()
        if not row:
            print("[INFO] No active 'APPLIED' transactions available to undo.")
            return ExitCodeType.SUCCESS.value
        tx_id = row[0]

    cur.execute("SELECT transaction_id, status FROM transactions WHERE transaction_id = ?;", (tx_id,))
    tx = cur.fetchone()
    if not tx:
        print(f"[ERROR] Transaction not found: {tx_id}")
        return ExitCodeType.GENERIC_FAILURE.value
    if tx[1] != "APPLIED":
        print(f"[WARN] Transaction {tx_id} is already in state '{tx[1]}'.")
        return ExitCodeType.SUCCESS.value

    cur.execute("SELECT seq_num, op_type, source_path, target_path, original_content, new_content FROM operations WHERE transaction_id = ? ORDER BY seq_num DESC;", (tx_id,))
    ops = cur.fetchall()

    print(f"<< Undoing Transaction: {tx_id} ({len(ops)} operations to revert)...")
    reverted_count = 0

    with conn:
        for seq, op_type, src_str, dst_str, old_txt, new_txt in ops:
            src = Path(src_str)
            dst = Path(dst_str) if dst_str else None

            if op_type in {"MOVE_DIR", "MOVE_FILE"}:
                if dst and dst.exists():
                    src.parent.mkdir(parents=True, exist_ok=True)
                    shutil.move(str(dst), str(src))
                    reverted_count += 1
                    print(f"  * Restored move: {dst} -> {src}")

            elif op_type == "MODIFY_FILE":
                if src.exists() and old_txt is not None:
                    src.write_text(old_txt, encoding="utf-8")
                    reverted_count += 1
                    print(f"  * Restored content: {src}")

            elif op_type == "TRASH_FILE":
                print(f"  * Cleaned item was in trash: {src}")

        conn.execute("UPDATE transactions SET status = 'REVERTED' WHERE transaction_id = ?;", (tx_id,))
        conn.execute("UPDATE operations SET status = 'REVERTED' WHERE transaction_id = ?;", (tx_id,))

    print(f"\n[PASS] Successfully undone transaction {tx_id} ({reverted_count} operations reverted).")
    return ExitCodeType.SUCCESS.value


def redo_transaction(tx_id: str | None = None) -> int:
    """Re-applies a previously reverted migration transaction."""
    if not DB_PATH.exists():
        print("[ERROR] Database tmp/migrations.db does not exist.")
        return ExitCodeType.GENERIC_FAILURE.value

    conn = init_db()
    cur = conn.cursor()

    if not tx_id:
        cur.execute("SELECT transaction_id FROM transactions WHERE status = 'REVERTED' ORDER BY timestamp DESC LIMIT 1;")
        row = cur.fetchone()
        if not row:
            print("[INFO] No 'REVERTED' transactions available to redo.")
            return ExitCodeType.SUCCESS.value
        tx_id = row[0]

    cur.execute("SELECT transaction_id, status FROM transactions WHERE transaction_id = ?;", (tx_id,))
    tx = cur.fetchone()
    if not tx:
        print(f"[ERROR] Transaction not found: {tx_id}")
        return ExitCodeType.GENERIC_FAILURE.value
    if tx[1] != "REVERTED":
        print(f"[WARN] Transaction {tx_id} is already in state '{tx[1]}'.")
        return ExitCodeType.SUCCESS.value

    cur.execute("SELECT seq_num, op_type, source_path, target_path, original_content, new_content FROM operations WHERE transaction_id = ? ORDER BY seq_num ASC;", (tx_id,))
    ops = cur.fetchall()

    print(f">> Redoing Transaction: {tx_id} ({len(ops)} operations to apply)...")
    applied_count = 0

    with conn:
        for seq, op_type, src_str, dst_str, old_txt, new_txt in ops:
            src = Path(src_str)
            dst = Path(dst_str) if dst_str else None

            if op_type in {"MOVE_DIR", "MOVE_FILE"}:
                if src.exists() and dst:
                    dst.parent.mkdir(parents=True, exist_ok=True)
                    shutil.move(str(src), str(dst))
                    applied_count += 1
                    print(f"  * Re-applied move: {src} -> {dst}")

            elif op_type == "MODIFY_FILE":
                if src.exists() and new_txt is not None:
                    src.write_text(new_txt, encoding="utf-8")
                    applied_count += 1
                    print(f"  * Re-applied content: {src}")

        conn.execute("UPDATE transactions SET status = 'APPLIED' WHERE transaction_id = ?;", (tx_id,))
        conn.execute("UPDATE operations SET status = 'APPLIED' WHERE transaction_id = ?;", (tx_id,))

    print(f"\n[PASS] Successfully re-applied transaction {tx_id} ({applied_count} operations applied).")
    return ExitCodeType.SUCCESS.value


def main():
    parser = argparse.ArgumentParser(
        description="Transactional Repository Layout Migrator with SQLite Journal, Undo/Redo, and History",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Subcommands:
  migrate    Preview or execute migration on a repository
  ls         List past migration transactions
  show       Inspect file-by-file details of a transaction
  undo       Roll back an applied migration transaction
  redo       Re-apply a previously reverted migration transaction

Examples:
  python 03-ai-scripts/25-repo-migrator.py migrate --plan
  python 03-ai-scripts/25-repo-migrator.py migrate --path ../other-repo --force
  python 03-ai-scripts/25-repo-migrator.py ls
  python 03-ai-scripts/25-repo-migrator.py show tx-20260901-120000-abc123
  python 03-ai-scripts/25-repo-migrator.py undo
  python 03-ai-scripts/25-repo-migrator.py redo
"""
    )
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")

    # Migrate Subcommand
    migrate_parser = subparsers.add_parser("migrate", help="Migrate repository layout to 01-prompts, 02-spec, 03-ai-scripts")
    migrate_parser.add_argument("--path", "-p", default=CURRENT_DIR, help="Target repository root path (default: .)")
    migrate_parser.add_argument("--plan", "--dry-run", "-d", action="store_true", help="Plan Mode: preview changes without executing")
    migrate_parser.add_argument("--force", "-f", "-y", "--yes", action="store_true", help="Bypass interactive confirmation prompt")

    # LS / History Subcommand
    ls_parser = subparsers.add_parser("ls", help="List past migration transactions")
    history_parser = subparsers.add_parser("history", help="Alias for 'ls'")

    # Show Subcommand
    show_parser = subparsers.add_parser("show", help="Show details of a specific transaction")
    show_parser.add_argument("tx_id", help="Transaction ID to inspect (e.g. tx-20260901-120000-abc123)")

    # Undo Subcommand
    undo_parser = subparsers.add_parser("undo", help="Roll back a migration transaction")
    undo_parser.add_argument("--tx-id", help="Specific transaction ID to undo (default: latest APPLIED)")

    # Redo Subcommand
    redo_parser = subparsers.add_parser("redo", help="Re-apply a reverted migration transaction")
    redo_parser.add_argument("--tx-id", help="Specific transaction ID to redo (default: latest REVERTED)")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(0)

    if args.command == "migrate":
        sys.exit(execute_migration(Path(args.path), is_plan_only=args.plan, is_force=args.force))
    elif args.command in {"ls", "history"}:
        sys.exit(list_transactions())
    elif args.command == "show":
        sys.exit(show_transaction(args.tx_id))
    elif args.command == "undo":
        sys.exit(undo_transaction(args.tx_id))
    elif args.command == "redo":
        sys.exit(redo_transaction(args.tx_id))


if __name__ == "__main__":
    main()
