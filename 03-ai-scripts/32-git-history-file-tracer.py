#!/usr/bin/env python3
"""
Git History Removed File Tracer, Pre-Flight Inspector, Restorer & Purger
Autonomously traces deleted files across Git history, generates numbered
pre-flight reports, supports selective inclusion/exclusion, restores files
from their pre-deletion commits, or permanently purges them from Git history.

All Enums, Constants, and Base Path Utilities are imported directly from
03-ai-scripts/02-shared-engine.py.

Usage Examples:
  # 1. Pre-flight inspection of removed markdown files in .ai-memory/ (default mode):
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable

  # 2. Pre-flight inspection of removed markdown files in 02-spec/:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-spec

  # 3. Custom path and extension filter:
  python 03-ai-scripts/32-git-history-file-tracer.py --path 04-code/ --ext .go,.ts

  # 4. Filter by numbering: only include items 1 to 5 and exclude item 3:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1-5 --exclude 3

  # 5. Restore selected files to a separate staging directory:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1,2 --restore-to tmp/recovery/

  # 6. Restore selected files in-place:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1 --restore

  # 7. Permanently purge removed files from all Git history (branches, tags):
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1,2 --purge --confirm-purge
"""

import argparse
from dataclasses import dataclass
import fnmatch
from importlib import import_module
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

# Centralized imports from shared engine
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR
PATH_SEPARATOR = engine.PATH_SEPARATOR
LINE_SEPARATOR = engine.LINE_SEPARATOR
EMPTY_STRING = engine.EMPTY_STRING
GIT_EXECUTABLE = engine.GIT_EXECUTABLE
ExitCodeType = engine.ExitCodeType
CACHE_KEY_FILES = engine.CACHE_KEY_FILES

load_repo_cache = engine.load_repo_cache
normalize_rel_path = engine.normalize_rel_path
normalize_extensions = engine.normalize_extensions
is_ignored_directory = engine.is_ignored_directory

RECORD_DELIMITER = "TRACE_COMMIT_RECORD"
PRESET_LOVABLE_PATH = ".ai-memory"
PRESET_PLANS_PATH = ".ai-memory/plans"
PRESET_SUBTASKS_PATH = ".ai-memory/plans/subtasks"
PRESET_SPEC_PATH = "02-spec"
PRESET_AUDIT_PATH = "02-spec/25-app-spec-audit"
DEFAULT_MD_EXT = {".md"}


@dataclass
class DeletedFileRecord:
    """Stores metadata of a historically deleted file."""
    index: int
    file_path: str
    deletion_commit_full: str
    deletion_commit_short: str
    parent_commit: str
    commit_date: str
    author: str
    commit_subject: str
    is_on_disk: bool = False
    size_bytes: int = 0


def parse_index_spec(spec_str: str | None, max_bound: int) -> set[int]:
    """Parses a comma-separated range string (e.g. '1,3,5-8') into a set of 1-based indices."""
    result: set[int] = set()
    has_spec = bool(spec_str and spec_str.strip())
    if not has_spec:
        return result

    for part in spec_str.split(","):
        cleaned = part.strip()
        has_dash = ("-" in cleaned)
        if has_dash:
            bounds = cleaned.split("-", 1)
            is_valid_range = bool(len(bounds) == 2 and bounds[0].isdigit() and bounds[1].isdigit())
            if is_valid_range:
                start = max(1, int(bounds[0]))
                end = min(max_bound, int(bounds[1]))
                if start <= end:
                    result.update(range(start, end + 1))
        elif cleaned.isdigit():
            val = int(cleaned)
            is_in_bounds = bool(1 <= val <= max_bound)
            if is_in_bounds:
                result.add(val)

    return result


def resolve_scope_and_extensions(args: argparse.Namespace) -> tuple[str, set[str]]:
    """Resolves target directory path and file extension filter based on presets or arguments."""
    is_audit = bool(args.preset_audit)
    if is_audit:
        return PRESET_AUDIT_PATH, DEFAULT_MD_EXT

    is_plans = bool(args.preset_plans)
    if is_plans:
        return PRESET_PLANS_PATH, DEFAULT_MD_EXT

    is_subtasks = bool(args.preset_subtasks)
    if is_subtasks:
        return PRESET_SUBTASKS_PATH, DEFAULT_MD_EXT

    is_lovable = bool(args.preset_lovable)
    if is_lovable:
        return PRESET_LOVABLE_PATH, DEFAULT_MD_EXT

    is_spec = bool(args.preset_spec)
    if is_spec:
        return PRESET_SPEC_PATH, DEFAULT_MD_EXT

    raw_path = args.path.strip() if args.path else CURRENT_DIR
    norm_path = normalize_rel_path(raw_path).rstrip(PATH_SEPARATOR)
    resolved_path = norm_path if norm_path else CURRENT_DIR

    ext_set = normalize_extensions(args.ext) if args.ext else set()
    return resolved_path, ext_set


def query_deleted_log_entries(path_scope: str) -> list[str]:
    """Executes git log with --diff-filter=D to retrieve deletion commits and files."""
    cmd = [
        GIT_EXECUTABLE,
        "log",
        "--diff-filter=D",
        f"--pretty=format:{RECORD_DELIMITER}|%H|%h|%an|%ad|%s",
        "--date=short",
        "--name-only",
    ]
    is_targeted_scope = bool(path_scope and path_scope != CURRENT_DIR)
    if is_targeted_scope:
        cmd.extend(["--", path_scope])

    try:
        proc = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace",
            check=True,
        )
        return proc.stdout.splitlines()
    except subprocess.CalledProcessError as err:
        print(f"[ERROR] Failed to query Git log: {err}")
        return []


def parse_deletion_records(
    log_lines: list[str],
    path_scope: str,
    extensions: set[str],
    pattern: str | None,
) -> list[DeletedFileRecord]:
    """Parses raw git log output into a deduplicated list of DeletedFileRecord objects."""
    records: list[DeletedFileRecord] = []
    seen_paths: set[str] = set()
    curr_commit: dict[str, str] = {}

    scope_prefix = path_scope.rstrip(PATH_SEPARATOR) + PATH_SEPARATOR if path_scope != CURRENT_DIR else EMPTY_STRING

    for raw_line in log_lines:
        line = raw_line.strip()
        is_empty = bool(not line)
        if is_empty:
            continue

        is_header = line.startswith(RECORD_DELIMITER + "|")
        if is_header:
            parts = line.split("|", 5)
            curr_commit = {
                "full_hash": parts[1] if len(parts) > 1 else EMPTY_STRING,
                "short_hash": parts[2] if len(parts) > 2 else EMPTY_STRING,
                "author": parts[3] if len(parts) > 3 else EMPTY_STRING,
                "date": parts[4] if len(parts) > 4 else EMPTY_STRING,
                "subject": parts[5] if len(parts) > 5 else EMPTY_STRING,
            }
            continue

        file_path = normalize_rel_path(line)
        is_duplicate = (file_path in seen_paths)
        if is_duplicate:
            continue

        is_in_scope = bool(not scope_prefix or file_path.startswith(scope_prefix))
        if not is_in_scope:
            continue

        p_obj = Path(file_path)
        is_ext_match = bool(not extensions or p_obj.suffix.lower() in extensions)
        if not is_ext_match:
            continue

        is_pattern_match = bool(not pattern or fnmatch.fnmatch(file_path, pattern) or fnmatch.fnmatch(p_obj.name, pattern))
        if not is_pattern_match:
            continue

        seen_paths.add(file_path)
        deletion_full = curr_commit.get("full_hash", EMPTY_STRING)
        record = DeletedFileRecord(
            index=len(records) + 1,
            file_path=file_path,
            deletion_commit_full=deletion_full,
            deletion_commit_short=curr_commit.get("short_hash", EMPTY_STRING),
            parent_commit=f"{deletion_full}^" if deletion_full else EMPTY_STRING,
            commit_date=curr_commit.get("date", EMPTY_STRING),
            author=curr_commit.get("author", EMPTY_STRING),
            commit_subject=curr_commit.get("subject", EMPTY_STRING),
        )
        records.append(record)

    return records


def populate_disk_and_size_info(records: list[DeletedFileRecord]) -> None:
    """Populates on-disk existence and historical file size for each record."""
    cache = load_repo_cache()
    cached_files = set(cache[CACHE_KEY_FILES]) if cache and CACHE_KEY_FILES in cache else None

    for rec in records:
        if cached_files is not None:
            rec.is_on_disk = (rec.file_path in cached_files)
        else:
            rec.is_on_disk = os.path.exists(rec.file_path)

        has_parent = bool(rec.parent_commit)
        if has_parent:
            rec.size_bytes = query_blob_size(rec.parent_commit, rec.file_path)


def query_blob_size(parent_commit: str, file_path: str) -> int:
    """Queries blob size in bytes from git cat-file."""
    try:
        res = subprocess.run(
            [GIT_EXECUTABLE, "cat-file", "-s", f"{parent_commit}:{file_path}"],
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace",
        )
        has_size = bool(res.returncode == 0 and res.stdout.strip().isdigit())
        if has_size:
            return int(res.stdout.strip())
    except Exception:
        pass
    return 0


def filter_records_by_selection(
    records: list[DeletedFileRecord],
    include_set: set[int],
    exclude_set: set[int],
    exclude_pattern: str | None,
) -> list[DeletedFileRecord]:
    """Filters records based on include set, exclude set, and exclude glob pattern."""
    filtered: list[DeletedFileRecord] = []
    has_include = bool(include_set)

    for rec in records:
        if has_include:
            is_included = (rec.index in include_set)
            if not is_included:
                continue

        is_excluded_num = (rec.index in exclude_set)
        if is_excluded_num:
            continue

        has_exclude_pat = bool(exclude_pattern)
        if has_exclude_pat:
            is_pat_match = bool(fnmatch.fnmatch(rec.file_path, exclude_pattern) or fnmatch.fnmatch(Path(rec.file_path).name, exclude_pattern))
            if is_pat_match:
                continue

        filtered.append(rec)

    return filtered


def format_bytes_human(num_bytes: int) -> str:
    """Formats byte counts into human-readable strings."""
    if num_bytes >= 1048576:
        return f"{num_bytes / 1048576:.1f} MB"
    if num_bytes >= 1024:
        return f"{num_bytes / 1024:.1f} KB"
    return f"{num_bytes} B"


def display_preflight_report(
    records: list[DeletedFileRecord],
    path_scope: str,
    ext_set: set[str],
    elapsed_ms: float,
) -> None:
    """Renders a formatted pre-flight report table to the terminal."""
    ext_str = ", ".join(sorted(ext_set)) if ext_set else "all"
    print("=" * 80)
    print("🔍 GIT HISTORY REMOVED FILE TRACER — PRE-FLIGHT REPORT")
    print("=" * 80)
    print(f"📁 Target Scope : `{path_scope}`")
    print(f"🏷️  Extensions   : {ext_str}")
    print(f"⏱️  Scan Time    : {elapsed_ms:.2f}ms")
    print(f"📊 Total Traced : {len(records)} removed file(s)")
    print("-" * 80)

    is_empty = bool(not records)
    if is_empty:
        print("  (No deleted files found matching the criteria in Git history)")
        print("=" * 80)
        return

    total_bytes = 0
    on_disk_count = 0

    for rec in records:
        total_bytes += rec.size_bytes
        if rec.is_on_disk:
            on_disk_count += 1
            status_badge = "⚠️  ACTIVE ON DISK (Re-introduced)"
        else:
            status_badge = "❌ NOT ON DISK (Eligible for restore/purge)"

        size_str = format_bytes_human(rec.size_bytes)
        print(f"[{rec.index:>3}] {rec.file_path}")
        print(f"      • Deletion Commit : {rec.deletion_commit_short} ({rec.commit_date}) by {rec.author}")
        print(f"      • Commit Subject  : {rec.commit_subject}")
        print(f"      • Recovery Source : {rec.parent_commit} ({size_str})")
        print(f"      • Status          : {status_badge}")
        print()

    print("-" * 80)
    print(f"📈 Summary: {len(records)} removed files | {format_bytes_human(total_bytes)} recoverable | {on_disk_count} currently on disk")
    print("💡 Next Steps:")
    print("  • To restore files to working tree:   python 03-ai-scripts/32-git-history-file-tracer.py <flags> --restore")
    print("  • To restore to a staging directory:  python 03-ai-scripts/32-git-history-file-tracer.py <flags> --restore-to tmp/recovery/")
    print("  • To permanently purge from Git:      python 03-ai-scripts/32-git-history-file-tracer.py <flags> --purge --confirm-purge")
    print("=" * 80)


def create_safety_backup_branch() -> str:
    """Creates and pushes a timestamped safety backup branch before destructive actions."""
    timestamp = time.strftime("%Y%m%d-%H%M%S")
    backup_name = f"backup/git-purge-{timestamp}"
    try:
        subprocess.run([GIT_EXECUTABLE, "branch", backup_name], check=True, capture_output=True)
        print(f"🛡️  Created local safety backup branch: {backup_name}")
        return backup_name
    except subprocess.CalledProcessError as err:
        print(f"[WARNING] Could not create local backup branch: {err}")
        return EMPTY_STRING


def execute_restore_action(records: list[DeletedFileRecord], restore_dir: str | None) -> int:
    """Restores selected files from their parent commit prior to deletion."""
    restored_count = 0
    target_base = Path(restore_dir) if restore_dir else None

    if target_base:
        target_base.mkdir(parents=True, exist_ok=True)
        print(f"📥 Restoring {len(records)} file(s) into staging directory: `{target_base}`")
    else:
        print(f"📥 Restoring {len(records)} file(s) in-place into repository working tree...")

    for rec in records:
        has_parent = bool(rec.parent_commit)
        if not has_parent:
            print(f"[SKIP] No parent commit available for `{rec.file_path}`")
            continue

        try:
            res = subprocess.run(
                [GIT_EXECUTABLE, "cat-file", "-p", f"{rec.parent_commit}:{rec.file_path}"],
                capture_output=True,
                check=True,
            )
            raw_data = res.stdout

            out_path = (target_base / rec.file_path) if target_base else Path(rec.file_path)
            out_path.parent.mkdir(parents=True, exist_ok=True)
            with open(out_path, "wb") as f:
                f.write(raw_data)

            print(f"  ✅ Restored: {out_path} ({format_bytes_human(len(raw_data))})")
            restored_count += 1
        except subprocess.CalledProcessError as err:
            print(f"  ❌ Failed to restore `{rec.file_path}`: {err}")

    print(f"\n🎉 Restoration complete: {restored_count} of {len(records)} files successfully recovered.")
    return restored_count


def _trash_windows(target_path: Path) -> bool:
    """Uses Windows PowerShell VisualBasic FileSystem to send to Recycle Bin."""
    try:
        abs_str = str(target_path.resolve()).replace("\\", "\\\\")
        method = "DeleteDirectory" if target_path.is_dir() else "DeleteFile"
        ps_cmd = f"Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::{method}('{abs_str}', 'OnlyErrorDialogs', 'SendToRecycleBin')"
        res = subprocess.run(["powershell", "-NoProfile", "-NonInteractive", "-Command", ps_cmd], capture_output=True, text=True)

        is_success = bool(res.returncode == 0 and not target_path.exists())
        return is_success
    except Exception:
        return False


def _trash_posix(target_path: Path) -> bool:
    """Uses macOS trash or Linux gio/trash-put to move files to OS trash."""
    resolved = str(target_path.resolve())
    tools = ["trash"] if sys.platform == "darwin" else ["gio", "trash-put"]
    for tool in tools:
        try:
            cmd = ["gio", "trash", resolved] if tool == "gio" else [tool, resolved]
            res = subprocess.run(cmd, capture_output=True)

            is_moved = bool(res.returncode == 0 and not target_path.exists())
            if is_moved:
                return True
        except Exception:
            continue

    return False


def send_to_trash_bin(target_path: Path) -> bool:
    """Safely moves a file or directory to the OS Recycle / Trash Bin with fallbacks."""
    is_missing = bool(not target_path.exists())
    if is_missing:
        return False

    try:
        import send2trash
        send2trash.send2trash(str(target_path.resolve()))
        return True
    except Exception:
        pass

    is_windows = bool(sys.platform == "win32")
    if is_windows and _trash_windows(target_path):
        return True

    is_posix = bool(sys.platform == "darwin" or sys.platform.startswith("linux"))
    if is_posix and _trash_posix(target_path):
        return True

    is_dir = target_path.is_dir()
    if is_dir:
        shutil.rmtree(target_path)
    else:
        target_path.unlink()

    return True


def _backup_single_record(rec: DeletedFileRecord, backup_base: Path) -> bool:
    """Backs up a single record to the OS temp directory from disk or git history."""
    dest_file = backup_base / rec.file_path
    dest_file.parent.mkdir(parents=True, exist_ok=True)
    local_path = Path(rec.file_path)

    is_local_file = local_path.is_file()
    if is_local_file:
        shutil.copy2(local_path, dest_file)
        return True

    has_parent = bool(rec.parent_commit)
    if has_parent:
        try:
            res = subprocess.run(
                [GIT_EXECUTABLE, "cat-file", "-p", f"{rec.parent_commit}:{rec.file_path}"],
                capture_output=True,
                check=True,
            )
            with open(dest_file, "wb") as f:
                f.write(res.stdout)
            return True
        except Exception:
            pass

    return False


def create_temp_directory_backup(records: list[DeletedFileRecord]) -> Path:
    """Creates a full physical backup of all targeted files in the OS temp directory."""
    timestamp = time.strftime("%Y%m%d-%H%M%S")
    backup_base = Path(tempfile.gettempdir()) / f"git-history-tracer-backup-{timestamp}"
    backup_base.mkdir(parents=True, exist_ok=True)

    backed_up = sum(1 for rec in records if _backup_single_record(rec, backup_base))

    print(f"🛡️  Created OS temp directory safety backup ({backed_up} files):")
    print(f"    📂 {backup_base}")
    return backup_base


def execute_delete_action(records: list[DeletedFileRecord]) -> int:
    """Safely deletes on-disk files using the OS Recycle Bin after temp backup."""
    backup_base = create_temp_directory_backup(records)
    trashed_count = 0

    print(f"🗑️  Moving {len(records)} targeted file(s) to the OS Recycle / Trash Bin...")
    for rec in records:
        target_path = Path(rec.file_path)
        is_on_disk = target_path.exists()
        if is_on_disk:
            is_recycled = send_to_trash_bin(target_path)
            if is_recycled:
                print(f"  🗑️ Sent to Recycle Bin: {rec.file_path}")
                trashed_count += 1
            else:
                print(f"  ❌ Failed to recycle: {rec.file_path}")
        else:
            print(f"  ⏭️ Already removed from disk: {rec.file_path}")

    print("-" * 80)
    print(f"🎉 Working tree cleanup complete: {trashed_count} file(s) sent to OS Recycle Bin.")
    print(f"🛡️  OS Temp Backup Location: {backup_base}")
    print("=" * 80)
    return trashed_count


def execute_purge_action(records: list[DeletedFileRecord], force: bool) -> int:
    """Permanently purges specified files from Git history with OS temp backup & recycle bin cleanup."""
    is_confirmed = bool(force)
    if not is_confirmed:
        print("[ERROR] History purge requires explicit confirmation via `--confirm-purge`.")
        print("        Purging rewrites Git history across commits, branches, and tags.")
        return 0

    backup_base = create_temp_directory_backup(records)
    backup_branch = create_safety_backup_branch()

    recycled_count = 0
    for rec in records:
        target_path = Path(rec.file_path)
        is_on_disk = target_path.exists()
        if is_on_disk and send_to_trash_bin(target_path):
            recycled_count += 1

    paths_to_purge = [rec.file_path for rec in records]
    print(f"🔥 Initiating full Git history purge for {len(paths_to_purge)} file(s)...")

    cmd = [GIT_EXECUTABLE, "filter-repo", "--invert-paths", "--force"]
    for p in paths_to_purge:
        cmd.extend(["--path", p])

    try:
        proc = subprocess.run(cmd, capture_output=True, text=True, encoding=DEFAULT_ENCODING, check=True)
        print(proc.stdout)
        print("🧹 Reclaiming repository storage via reflog expire and aggressive gc...")
        subprocess.run([GIT_EXECUTABLE, "reflog", "expire", "--expire=now", "--all"], check=True)
        subprocess.run([GIT_EXECUTABLE, "gc", "--prune=now", "--aggressive"], check=True)
        print(f"🎉 History purge complete. {len(paths_to_purge)} file(s) erased across all branches and tags.")
        has_recycled = bool(recycled_count > 0)
        if has_recycled:
            print(f"🗑️  Working tree files sent to OS Recycle Bin: {recycled_count}")
        print(f"🛡️  OS Temp Directory Backup: {backup_base}")
        has_backup_branch = bool(backup_branch)
        if has_backup_branch:
            print(f"ℹ️  To roll back Git history if needed: `git reset --hard {backup_branch}`")
        return len(paths_to_purge)
    except FileNotFoundError:
        print("[ERROR] `git-filter-repo` executable not found. Install via `pip install git-filter-repo`.")
        print(f"🛡️  Your files remain safely backed up at: {backup_base}")
        return 0
    except subprocess.CalledProcessError as err:
        print(f"[ERROR] History purge failed: {err.stderr if err.stderr else err}")
        print(f"🛡️  Your files remain safely backed up at: {backup_base}")
        return 0


def build_arg_parser() -> argparse.ArgumentParser:
    """Builds the command-line argument parser with rich help text."""
    parser = argparse.ArgumentParser(
        description="Git History Removed File Tracer, Pre-Flight Inspector, Restorer & Purger",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Pre-flight inspection of removed markdown files in .ai-memory/ (Default):
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable

  # Pre-flight inspection of removed markdown files in 02-spec/:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-spec

  # Scan specific path and file extensions:
  python 03-ai-scripts/32-git-history-file-tracer.py --path 04-code/ --ext .go,.ts

  # Include only items 1 to 5 and exclude item 3:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1-5 --exclude 3

  # Restore selected files to a staging directory:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1,2 --restore-to tmp/recovery/

  # Restore selected files in-place into repository:
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1 --restore

  # Move selected files from working tree to OS Recycle Bin (with OS temp backup):
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete

  # Permanently purge selected files from all Git history (branches, tags):
  python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1,2 --purge --confirm-purge
""",
    )

    group_scope = parser.add_argument_group("Scope & Presets")
    group_scope.add_argument("--path", default=CURRENT_DIR, help="Target directory path to inspect from repository root (default: '.')")
    group_scope.add_argument("--ext", help="Comma-separated file extensions to filter (e.g. '.md', '.go,.ts')")
    group_scope.add_argument("--pattern", help="Glob pattern to match file names/paths (e.g. '*plan*', 'ss-*')")
    group_scope.add_argument("--preset-audit", "--audit", "--spec-audit", dest="preset_audit", action="store_true", help=f"Shortcut: target spec audit folder '{PRESET_AUDIT_PATH}/' for removed '.md' files")
    group_scope.add_argument("--preset-plans", "--plans", dest="preset_plans", action="store_true", help=f"Shortcut: target plans directory '{PRESET_PLANS_PATH}/' for removed '.md' files")
    group_scope.add_argument("--preset-subtasks", "--subtasks", dest="preset_subtasks", action="store_true", help=f"Shortcut: target subtasks directory '{PRESET_SUBTASKS_PATH}/' for removed '.md' files")
    group_scope.add_argument("--preset-lovable", "--lovable", dest="preset_lovable", action="store_true", help=f"Shortcut: target '{PRESET_LOVABLE_PATH}/' for removed '.md' files")
    group_scope.add_argument("--preset-spec", "--spec", dest="preset_spec", action="store_true", help=f"Shortcut: target '{PRESET_SPEC_PATH}/' for removed '.md' files")

    group_select = parser.add_argument_group("Filtering & Selection")
    group_select.add_argument("--include", help="Comma-separated list or ranges of numbers to target (e.g. '1,3,5-8')")
    group_select.add_argument("--exclude", help="Comma-separated list or ranges of numbers to exclude (e.g. '2,4')")
    group_select.add_argument("--exclude-pattern", help="Glob pattern to exclude files (e.g. 'ss-*')")

    group_action = parser.add_argument_group("Actions")
    group_action.add_argument("--preflight", "--plan", action="store_true", default=True, help="(Default) Pre-flight dry run")
    group_action.add_argument("--restore", action="store_true", help="Restore selected removed files in-place")
    group_action.add_argument("--restore-to", help="Restore selected files into a specific directory")
    group_action.add_argument("--delete", "--remove", "--trash", dest="delete", action="store_true", help="Move selected files from working tree to OS Recycle Bin (after OS temp backup)")
    group_action.add_argument("--purge", action="store_true", help="Permanently purge selected files from Git history")
    group_action.add_argument("--confirm-purge", action="store_true", help="Required safety confirmation flag for --purge")

    group_output = parser.add_argument_group("Output Formatting")
    group_output.add_argument("--json", action="store_true", help="Output pre-flight results in JSON format")

    return parser


def main() -> int:
    """Main CLI entrypoint."""
    start_time = time.perf_counter()
    parser = build_arg_parser()
    args = parser.parse_args()

    path_scope, ext_set = resolve_scope_and_extensions(args)
    raw_lines = query_deleted_log_entries(path_scope)
    records = parse_deletion_records(raw_lines, path_scope, ext_set, args.pattern)
    include_set = parse_index_spec(args.include, len(records))
    exclude_set = parse_index_spec(args.exclude, len(records))
    selected_records = filter_records_by_selection(records, include_set, exclude_set, args.exclude_pattern)
    populate_disk_and_size_info(selected_records)

    elapsed_ms = (time.perf_counter() - start_time) * 1000

    if args.json:
        output_data = [
            {
                "index": r.index,
                "file_path": r.file_path,
                "deletion_commit": r.deletion_commit_short,
                "parent_commit": r.parent_commit,
                "commit_date": r.commit_date,
                "author": r.author,
                "subject": r.commit_subject,
                "size_bytes": r.size_bytes,
                "is_on_disk": r.is_on_disk,
            }
            for r in selected_records
        ]
        print(json.dumps(output_data, indent=2))
        return ExitCodeType.SUCCESS.value

    is_delete_requested = bool(args.delete)
    if is_delete_requested:
        execute_delete_action(selected_records)
        return ExitCodeType.SUCCESS.value

    is_purge_requested = bool(args.purge)
    if is_purge_requested:
        execute_purge_action(selected_records, args.confirm_purge)
        return ExitCodeType.SUCCESS.value

    is_restore_requested = bool(args.restore or args.restore_to)
    if is_restore_requested:
        execute_restore_action(selected_records, args.restore_to)
        return ExitCodeType.SUCCESS.value

    display_preflight_report(selected_records, path_scope, ext_set, elapsed_ms)
    return ExitCodeType.SUCCESS.value


if __name__ == "__main__":
    sys.exit(main())
