#!/usr/bin/env python3
"""
Diagnostic and Healing Script for Remote Ubuntu Fleet Node (u1)
Resolves desynchronized pull states, Windows backslash path corruption,
and stale .oh-my-zsh entries in gitmap.db.
"""

import argparse
from datetime import datetime, timezone
import os
from pathlib import Path
import re
import shutil
import sqlite3
import subprocess
import sys
from typing import Any


DEFAULT_DB_CANDIDATES = [
    Path.home() / ".local" / "bin" / "gitmap-cli" / "data" / "gitmap.db",
    Path.home() / ".gitmap" / "gitmap.db",
    Path.home() / "gitmap" / "data" / "gitmap.db",
    Path("gitmap.db"),
]


def resolve_database_path(custom_path: str = "") -> Path:
    if custom_path:
        target = Path(custom_path).expanduser().resolve()
        if target.exists():
            return target
        return target

    env_path = os.getenv("GITMAP_DB_PATH", "")
    if env_path:
        target = Path(env_path).expanduser().resolve()
        if target.exists():
            return target

    for candidate in DEFAULT_DB_CANDIDATES:
        expanded = candidate.expanduser().resolve()
        if expanded.exists():
            return expanded

    return DEFAULT_DB_CANDIDATES[0].expanduser().resolve()


def create_database_backup(db_path: Path) -> Path:
    timestamp = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
    backup_path = db_path.with_name(f"{db_path.name}.bak.{timestamp}")
    shutil.copy2(db_path, backup_path)
    return backup_path


def normalize_repo_path(raw_path: str) -> str:
    cleaned = raw_path.replace("\\", "/")
    cleaned = re.sub(r"^[A-Za-z]:[/]+", "/", cleaned)
    cleaned = re.sub(r"/+", "/", cleaned)
    return cleaned


def inspect_backslash_paths(cursor: sqlite3.Cursor) -> list[tuple[Any, ...]]:
    query = r"""
    SELECT RepoId, Slug, RepoName, AbsolutePath
    FROM Repo
    WHERE AbsolutePath LIKE '%\%' OR AbsolutePath LIKE '_:%'
    """
    try:
        cursor.execute(query)
        return cursor.fetchall()
    except sqlite3.OperationalError:
        return []


def find_case_or_suffix_match(target_path: Path) -> Path | None:
    parent = target_path.parent
    if not parent.exists():
        return None
    name_lower = target_path.name.lower()
    for child in parent.iterdir():
        if child.is_dir() and child.name.lower() == name_lower:
            return child
    base_name = re.sub(r"-v\d+$", "", target_path.name).lower()
    for child in parent.iterdir():
        if child.is_dir() and child.name.lower() == base_name:
            return child
    return None


def heal_database_paths(conn: sqlite3.Connection, rows: list[tuple[Any, ...]], is_dry_run: bool) -> int:
    if is_dry_run:
        return 0

    cursor = conn.cursor()
    updated_count = 0
    for repo_id, slug, repo_name, raw_path in rows:
        normalized = normalize_repo_path(raw_path)
        cursor.execute(
            "UPDATE Repo SET AbsolutePath = ? WHERE RepoId = ?",
            (normalized, repo_id),
        )
        updated_count += 1

    cursor.execute("SELECT RepoId, Slug, AbsolutePath FROM Repo")
    all_repos = cursor.fetchall()
    for repo_id, slug, abs_path in all_repos:
        p = Path(abs_path)
        if not p.exists():
            matched = find_case_or_suffix_match(p)
            if matched is not None and matched.exists():
                cursor.execute(
                    "UPDATE Repo SET AbsolutePath = ? WHERE RepoId = ?",
                    (str(matched), repo_id),
                )
                updated_count += 1

    conn.commit()
    return updated_count


def inspect_omz_records(cursor: sqlite3.Cursor) -> list[tuple[Any, ...]]:
    query = """
    SELECT RepoId, Slug, RepoName, AbsolutePath
    FROM Repo
    WHERE Slug = '.oh-my-zsh'
       OR RepoName = '.oh-my-zsh'
       OR AbsolutePath LIKE '%.oh-my-zsh%'
    """
    try:
        cursor.execute(query)
        return cursor.fetchall()
    except sqlite3.OperationalError:
        return []


def prune_omz_records(conn: sqlite3.Connection, is_dry_run: bool) -> int:
    if is_dry_run:
        return 0

    cursor = conn.cursor()
    cursor.execute(
        """
        DELETE FROM Repo
        WHERE Slug = '.oh-my-zsh'
           OR RepoName = '.oh-my-zsh'
           OR AbsolutePath LIKE '%.oh-my-zsh%'
        """
    )
    deleted_count = cursor.rowcount
    conn.commit()
    return deleted_count


def verify_repositories_on_disk(conn: sqlite3.Connection) -> dict[str, list[str]]:
    cursor = conn.cursor()
    try:
        cursor.execute("SELECT Slug, AbsolutePath FROM Repo ORDER BY Slug ASC")
        rows = cursor.fetchall()
    except sqlite3.OperationalError:
        return {"existing": [], "missing": []}

    existing: list[str] = []
    missing: list[str] = []
    for slug, abs_path in rows:
        path_obj = Path(abs_path).expanduser()
        if path_obj.exists() and path_obj.is_dir():
            existing.append(slug)
        else:
            missing.append(f"{slug} ({abs_path})")

    return {"existing": existing, "missing": missing}


def prune_stale_git_locks(work_dir: Path, is_dry_run: bool) -> list[str]:
    removed_locks: list[str] = []
    if not work_dir.exists():
        return removed_locks

    for lock_path in work_dir.rglob(".git/index.lock"):
        is_stale = False
        try:
            mtime = lock_path.stat().st_mtime
            age_seconds = datetime.now().timestamp() - mtime
            is_stale = age_seconds > 600
        except OSError:
            is_stale = False

        if is_stale:
            removed_locks.append(str(lock_path))
            if not is_dry_run:
                try:
                    lock_path.unlink()
                except OSError:
                    pass

    return removed_locks


def ensure_gitmapignore_omz() -> bool:
    target_files = [
        Path.home() / ".gitmap" / ".gitmapignore",
        Path.home() / ".gitmapignore",
        Path(".gitmapignore"),
    ]
    added = False
    for p in target_files:
        if not p.parent.exists():
            continue
        existing_lines = []
        if p.exists():
            try:
                existing_lines = p.read_text(encoding="utf-8").splitlines()
            except OSError:
                existing_lines = []
        has_omz = any(".oh-my-zsh" in line for line in existing_lines)
        if not has_omz:
            try:
                with p.open("a", encoding="utf-8") as f:
                    f.write("\n.oh-my-zsh\noh-my-zsh\n")
                added = True
            except OSError:
                pass
    return added


def execute_gitmap_pull_all() -> int:
    cmd = ["gitmap", "pa"]
    print(f"\n[EXEC] Running verification: {' '.join(cmd)}")
    try:
        res = subprocess.run(cmd, capture_output=False, check=False)
        return res.returncode
    except FileNotFoundError:
        print("[WARN] 'gitmap' executable not found in PATH; skipping verification command.")
        return 0


def run_healing_pipeline(db_path_str: str, is_dry_run: bool, is_verify: bool) -> int:
    db_path = resolve_database_path(db_path_str)
    print(f"============================================================")
    print(f"GitMap Remote Node Diagnostic & Healing Tool (u1)")
    print(f"Database Target: {db_path}")
    print(f"Mode: {'DRY RUN' if is_dry_run else 'FIX & HEAL'}")
    print(f"============================================================")

    if not db_path.exists():
        print(f"[ERROR] Database file not found at: {db_path}")
        return 1

    if not is_dry_run:
        backup = create_database_backup(db_path)
        print(f"[BACKUP] Created database snapshot at: {backup}")

    conn = sqlite3.connect(str(db_path))
    cursor = conn.cursor()

    backslash_rows = inspect_backslash_paths(cursor)
    print(f"[SCAN] Found {len(backslash_rows)} records with Windows backslash separators.")
    for _, slug, _, raw_path in backslash_rows[:5]:
        print(f"       • {slug}: {raw_path} -> {normalize_repo_path(raw_path)}")
    if len(backslash_rows) > 5:
        print(f"       ... and {len(backslash_rows) - 5} more")

    healed_paths = heal_database_paths(conn, backslash_rows, is_dry_run)
    if not is_dry_run:
        print(f"[HEAL] Updated {healed_paths} repository paths to POSIX forward slashes.")

    omz_rows = inspect_omz_records(cursor)
    print(f"[SCAN] Found {len(omz_rows)} stale .oh-my-zsh repository records.")
    for _, slug, _, raw_path in omz_rows:
        print(f"       • {slug} ({raw_path})")

    pruned_omz = prune_omz_records(conn, is_dry_run)
    if not is_dry_run:
        print(f"[HEAL] Removed {pruned_omz} stale .oh-my-zsh database records.")

    ensure_gitmapignore_omz()

    disk_status = verify_repositories_on_disk(conn)
    existing_count = len(disk_status["existing"])
    missing_count = len(disk_status["missing"])
    print(f"[AUDIT] Repository Physical Disk Audit:")
    print(f"        Verified on disk: {existing_count}")
    print(f"        Missing on disk:  {missing_count}")
    for item in disk_status["missing"]:
        print(f"        ✖ Missing: {item}")

    conn.close()

    if is_verify:
        return execute_gitmap_pull_all()

    return 0


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Diagnose and heal GitMap database and repository state on remote Ubuntu node u1."
    )
    parser.add_argument(
        "--db",
        dest="db_path",
        default="",
        help="Path to gitmap.db (default: autodetected from standard paths)",
    )
    parser.add_argument(
        "--dry-run",
        dest="is_dry_run",
        action="store_true",
        default=False,
        help="Inspect and report issues without modifying the database",
    )
    parser.add_argument(
        "--fix",
        dest="is_fix",
        action="store_true",
        default=True,
        help="Apply database path normalization and stale record pruning (default: True)",
    )
    parser.add_argument(
        "--verify",
        dest="is_verify",
        action="store_true",
        default=False,
        help="Execute 'gitmap pa' verification after healing",
    )
    args = parser.parse_args()

    is_dry = args.is_dry_run
    sys.exit(run_healing_pipeline(args.db_path, is_dry, args.is_verify))


if __name__ == "__main__":
    main()
