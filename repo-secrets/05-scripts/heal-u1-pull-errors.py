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
    Path("data") / "gitmap.db",
    Path("gitmap.db"),
]


def resolve_database_path(custom_path: str = "") -> Path:
    if custom_path:
        return Path(custom_path).expanduser().resolve()

    env_path = os.getenv("GITMAP_DB_PATH", "")
    target_env = resolve_existing_path_or_none(env_path)
    if target_env is not None:
        return target_env

    for candidate in DEFAULT_DB_CANDIDATES:
        expanded = candidate.expanduser().resolve()
        if expanded.exists():
            return expanded

    return DEFAULT_DB_CANDIDATES[0].expanduser().resolve()


def resolve_existing_path_or_none(raw_path: str) -> Path | None:
    if not raw_path:
        return None
    target = Path(raw_path).expanduser().resolve()
    if target.exists():
        return target
    return None


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
    expanded = target_path.expanduser()
    parent = expanded.parent
    if not parent.exists():
        return None

    name_lower = expanded.name.lower()
    for child in parent.iterdir():
        is_case_match = child.is_dir() and child.name.lower() == name_lower
        if is_case_match:
            return child

    base_name = re.sub(r"-v\d+$", "", expanded.name).lower()
    for child in parent.iterdir():
        is_suffix_match = child.is_dir() and child.name.lower() == base_name
        if is_suffix_match:
            return child

    return None


def heal_database_paths(conn: sqlite3.Connection, rows: list[tuple[Any, ...]], is_dry_run: bool) -> int:
    cursor = conn.cursor()
    updated_count = 0

    if not is_dry_run:
        for repo_id, slug, repo_name, raw_path in rows:
            normalized = normalize_repo_path(raw_path)
            cursor.execute(
                "UPDATE Repo SET AbsolutePath = ? WHERE RepoId = ?",
                (normalized, repo_id),
            )
            updated_count += 1
    else:
        updated_count = len(rows)

    cursor.execute("SELECT RepoId, Slug, AbsolutePath FROM Repo")
    all_repos = cursor.fetchall()
    case_suffix_count = 0

    for repo_id, slug, abs_path in all_repos:
        p = Path(abs_path).expanduser()
        if p.exists():
            continue
        matched = find_case_or_suffix_match(p)
        is_valid_match = matched is not None and matched.exists()
        if not is_valid_match:
            continue
        case_suffix_count += 1
        if not is_dry_run:
            cursor.execute(
                "UPDATE Repo SET AbsolutePath = ? WHERE RepoId = ?",
                (str(matched), repo_id),
            )

    if not is_dry_run:
        conn.commit()

    return updated_count + case_suffix_count


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
        cursor = conn.cursor()
        rows = inspect_omz_records(cursor)
        return len(rows)

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
        is_stale = is_lock_older_than_seconds(lock_path, 600)
        if not is_stale:
            continue
        removed_locks.append(str(lock_path))
        if not is_dry_run:
            unlink_path_safely(lock_path)

    return removed_locks


def is_lock_older_than_seconds(lock_path: Path, max_age_seconds: float) -> bool:
    try:
        mtime = lock_path.stat().st_mtime
        age_seconds = datetime.now().timestamp() - mtime
        return age_seconds > max_age_seconds
    except OSError:
        return False


def unlink_path_safely(p: Path) -> bool:
    try:
        p.unlink()
        return True
    except OSError:
        return False


def collect_lock_search_directories(conn: sqlite3.Connection) -> list[Path]:
    candidates: set[Path] = {
        Path.home() / "work",
        Path.home() / "git-work",
        Path.cwd(),
    }
    cursor = conn.cursor()
    try:
        cursor.execute("SELECT AbsolutePath FROM Repo")
        for (raw_path,) in cursor.fetchall():
            p = Path(raw_path).expanduser()
            if p.parent.exists():
                candidates.add(p.parent)
    except sqlite3.OperationalError:
        pass

    return [d for d in candidates if d.exists()]


def ensure_gitmapignore_omz() -> bool:
    target_files = [
        Path.home() / ".gitmap" / ".gitmapignore",
        Path.home() / ".gitmapignore",
        Path(".gitmapignore"),
    ]
    has_added = False
    for p in target_files:
        if not p.parent.exists():
            continue
        existing_lines = read_ignore_file_lines(p)
        has_omz = any(".oh-my-zsh" in line for line in existing_lines)
        if has_omz:
            continue
        has_written = append_omz_entries_to_file(p)
        if has_written:
            has_added = True
    return has_added


def read_ignore_file_lines(p: Path) -> list[str]:
    if not p.exists():
        return []
    try:
        return p.read_text(encoding="utf-8").splitlines()
    except OSError:
        return []


def append_omz_entries_to_file(p: Path) -> bool:
    try:
        with p.open("a", encoding="utf-8") as f:
            f.write("\n.oh-my-zsh\noh-my-zsh\n")
        return True
    except OSError:
        return False


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
    print("============================================================")
    print("GitMap Remote Node Diagnostic & Healing Tool (u1)")
    print(f"Database Target: {db_path}")
    print(f"Mode: {'DRY RUN' if is_dry_run else 'FIX & HEAL'}")
    print("============================================================")

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
    action_label = "Would update" if is_dry_run else "Updated"
    print(f"[HEAL] {action_label} {healed_paths} repository paths (POSIX slashes & casing/suffix alignment).")

    omz_rows = inspect_omz_records(cursor)
    print(f"[SCAN] Found {len(omz_rows)} stale .oh-my-zsh repository records.")
    for _, slug, _, raw_path in omz_rows:
        print(f"       • {slug} ({raw_path})")

    pruned_omz = prune_omz_records(conn, is_dry_run)
    omz_label = "Would remove" if is_dry_run else "Removed"
    print(f"[HEAL] {omz_label} {pruned_omz} stale .oh-my-zsh database records.")

    ensure_gitmapignore_omz()

    search_dirs = collect_lock_search_directories(conn)
    total_pruned_locks: list[str] = []
    for s_dir in search_dirs:
        locks = prune_stale_git_locks(s_dir, is_dry_run)
        total_pruned_locks.extend(locks)

    lock_label = "Found" if is_dry_run else "Pruned"
    print(f"[LOCKS] {lock_label} {len(total_pruned_locks)} stale .git/index.lock files older than 600s.")
    for l_path in total_pruned_locks[:3]:
        print(f"        • {l_path}")

    disk_status = verify_repositories_on_disk(conn)
    existing_count = len(disk_status["existing"])
    missing_count = len(disk_status["missing"])
    print("[AUDIT] Repository Physical Disk Audit:")
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
