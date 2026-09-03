#!/usr/bin/env python3
"""
Fast Lovable Plans & Subtasks Consolidator
Automates the consolidation, archiving, cleanup, and monotonic re-sequencing of
Lovable execution plans and task lists (.lovable/plans/pending, completed, subtasks).

Features:
1. Safety Backup Branch Creation (backup/plans-consolidation-YYYYMMDD-HHMMSS).
2. Moves completed pending plans to .lovable/plans/completed/.
3. Removes superseded subtask directories and files from disk & Git index.
4. Monotonically re-sequences plan files (01-, 02-, 03-) and synchronizes .lovable/plans/01-index.md.
5. Interactive confirmation, --dry-run preview, and --force execution.

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/20-plan-consolidator.py [--dry-run]
  python 03-ai-scripts/20-plan-consolidator.py --archive <plan-filename> [--force]
  python 03-ai-scripts/20-plan-consolidator.py --resequence [--force]
  python 03-ai-scripts/20-plan-consolidator.py --clean-subtasks [--force]
  python 03-ai-scripts/20-plan-consolidator.py --backup
"""

import argparse
import datetime
from importlib import import_module
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

# Centralized Enums & Constants
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
CURRENT_DIR = engine.CURRENT_DIR
EMPTY_STRING = engine.EMPTY_STRING
LINE_SEPARATOR = engine.LINE_SEPARATOR
PATH_SEPARATOR = engine.PATH_SEPARATOR
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
get_compiled_regex = engine.get_compiled_regex
normalize_rel_path = engine.normalize_rel_path
read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf

# Plans Directory Paths
PLANS_DIR = Path(".lovable/plans")
PENDING_DIR = PLANS_DIR / "pending"
COMPLETED_DIR = PLANS_DIR / "completed"
SUBTASKS_DIR = PLANS_DIR / "subtasks"
PLANS_INDEX_FILE = PLANS_DIR / "01-index.md"

def create_backup_branch() -> str | None:
    """Creates and pushes a timestamped safety backup branch before modifying plans."""
    now_tag = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%d-%H%M%S")
    branch_name = f"backup/plans-consolidation-{now_tag}"
    try:
        res = subprocess.run(
            ["git", "branch", branch_name],
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace"
        )
        is_created = (res.returncode == 0)
        if is_created:
            print(f"🛡️ Safety Backup Branch Created: `{branch_name}`")
            return branch_name
    except Exception as e:
        print(f"⚠️ Failed to create backup branch: {e}")
    return None

def confirm_action(prompt_text: str) -> bool:
    """Prompts user for interactive confirmation."""
    try:
        response = input(f"{prompt_text} (y/N): ").strip().lower()
        return response in {"y", "yes"}
    except (EOFError, KeyboardInterrupt):
        return False

def resequence_directory_plans(target_dir: Path, is_fix_mode: bool = False) -> list[tuple[str, str]]:
    """Monotonically resequences numeric prefixes (01-, 02-, 03-) in a plans folder."""
    re_seq = get_compiled_regex(RegexPatternType.SEQ_PREFIX)
    changes = []
    if not target_dir.exists() or not target_dir.is_dir():
        return changes

    files = sorted([f for f in target_dir.iterdir() if f.is_file() and f.name.endswith(".md")])
    for idx, f in enumerate(files, start=1):
        m = re_seq.match(f.name)
        base_name = m.group(2) if m else f.name
        new_name = f"{idx:02d}-{base_name}"
        if f.name != new_name:
            changes.append((f.name, new_name))
            if is_fix_mode:
                dst = f.with_name(new_name)
                try:
                    subprocess.run(["git", "mv", str(f), str(dst)], capture_output=True)
                except Exception:
                    f.rename(dst)
    return changes

def sync_plans_index() -> None:
    """Regenerates .lovable/plans/01-index.md with current pending and completed catalogs."""
    re_h1 = get_compiled_regex(RegexPatternType.H1_HEADER)
    lines = [
        "# Plans Index",
        "",
        "Master directory of architectural and execution plans.",
        "",
        "## Pending Plans",
        ""
    ]

    # Collect pending plans
    if PENDING_DIR.exists():
        for f in sorted(PENDING_DIR.iterdir()):
            if f.is_file() and f.suffix == ".md":
                content = read_file_lf(f, encoding=DEFAULT_ENCODING)
                m = re_h1.search(content)
                title = m.group(4).strip() if m else f.stem
                lines.append(f"- [{f.name}](pending/{f.name}): {title}")
    else:
        lines.append("*(No pending plans)*")

    lines.extend([
        "",
        "## Completed Plans",
        ""
    ])

    # Collect completed plans
    if COMPLETED_DIR.exists() and any(COMPLETED_DIR.iterdir()):
        for f in sorted(COMPLETED_DIR.iterdir()):
            if f.is_file() and f.suffix == ".md":
                content = read_file_lf(f, encoding=DEFAULT_ENCODING)
                m = re_h1.search(content)
                title = m.group(4).strip() if m else f.stem
                lines.append(f"- [{f.name}](completed/{f.name}): {title}")
    else:
        lines.append("*(Archived under `.lovable/plans/completed/`)*")

    lines.append("")
    write_file_lf(PLANS_INDEX_FILE, LINE_SEPARATOR.join(lines), encoding=DEFAULT_ENCODING)
    print(f"  ✓ Synchronized: {normalize_rel_path(PLANS_INDEX_FILE)}")

def clean_superseded_subtasks(is_fix_mode: bool = False) -> list[str]:
    """Finds and removes subtask directories for plans that are already completed."""
    removed = []
    if not SUBTASKS_DIR.exists():
        return removed

    for item in SUBTASKS_DIR.iterdir():
        if item.is_dir():
            removed.append(normalize_rel_path(item))
            if is_fix_mode:
                try:
                    subprocess.run(["git", "rm", "-rf", str(item)], capture_output=True)
                except Exception:
                    pass
                shutil.rmtree(item, ignore_errors=True)
    return removed

def archive_plan(plan_name: str, is_force_mode: bool = False, is_dry_run_mode: bool = False) -> int:
    """Moves a single completed plan from pending/ to completed/ and cleans its subtasks."""
    src = PENDING_DIR / plan_name
    if not src.exists():
        # Try finding by partial name
        matches = [f for f in PENDING_DIR.glob(f"*{plan_name}*") if f.is_file()]
        if matches:
            src = matches[0]
        else:
            print(f"❌ Plan file not found: `{plan_name}` in `{normalize_rel_path(PENDING_DIR)}`")
            return ExitCodeType.VIOLATIONS_FOUND.value

    COMPLETED_DIR.mkdir(parents=True, exist_ok=True)
    dst = COMPLETED_DIR / src.name

    print(f"📦 Plan Archive Target: `{normalize_rel_path(src)}` -> `{normalize_rel_path(dst)}`")

    if is_dry_run_mode:
        print("ℹ️ Dry-run mode enabled. No changes made.")
        return ExitCodeType.SUCCESS.value

    if not is_force_mode:
        is_confirmed = confirm_action(f"Are you sure you want to archive `{src.name}` to completed/?")
        if not is_confirmed:
            print("❌ Operation canceled.")
            return ExitCodeType.SUCCESS.value

    # Move file via git mv or rename
    try:
        subprocess.run(["git", "mv", str(src), str(dst)], capture_output=True)
    except Exception:
        shutil.move(str(src), str(dst))

    # Resequence and sync index
    resequence_directory_plans(COMPLETED_DIR, is_fix_mode=True)
    resequence_directory_plans(PENDING_DIR, is_fix_mode=True)
    sync_plans_index()

    print(f"✅ Successfully archived `{src.name}` to `{normalize_rel_path(COMPLETED_DIR)}`.")
    return ExitCodeType.SUCCESS.value

def run_full_consolidation(
    is_force_mode: bool = False,
    is_dry_run_mode: bool = False,
    is_auto_backup: bool = True
) -> int:
    """Runs full plan discovery, subtask cleanup, re-sequencing, and index synchronization."""
    start_time = time.perf_counter()

    pending_files = [normalize_rel_path(f) for f in PENDING_DIR.glob("*.md")] if PENDING_DIR.exists() else []
    completed_files = [normalize_rel_path(f) for f in COMPLETED_DIR.glob("*.md")] if COMPLETED_DIR.exists() else []
    subtask_dirs = [normalize_rel_path(d) for d in SUBTASKS_DIR.iterdir() if d.is_dir()] if SUBTASKS_DIR.exists() else []

    print("================================================================================")
    print("📋 Lovable Plans Inventory & Consolidation Preview:")
    print(f"   • Pending Plans    : {len(pending_files)}")
    print(f"   • Completed Plans  : {len(completed_files)}")
    print(f"   • Subtask Folders  : {len(subtask_dirs)}")
    print("================================================================================")

    if is_dry_run_mode:
        print("ℹ️ Dry-run mode enabled. No changes made.")
        return ExitCodeType.SUCCESS.value

    if not is_force_mode:
        is_confirmed = confirm_action("Do you want to proceed with plans consolidation and re-sequencing?")
        if not is_confirmed:
            print("❌ Operation canceled.")
            return ExitCodeType.SUCCESS.value

    if is_auto_backup:
        create_backup_branch()

    # 1. Resequence pending and completed plans
    resequence_directory_plans(PENDING_DIR, is_fix_mode=True)
    resequence_directory_plans(COMPLETED_DIR, is_fix_mode=True)

    # 2. Synchronize plans index
    sync_plans_index()

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"{LINE_SEPARATOR}✅ Plans consolidation and index synchronization complete in {elapsed_ms:.2f}ms.")
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(
        description="Fast Lovable Plans & Subtasks Consolidator",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument("--archive", "-a", help="Archive a completed plan from pending/ to completed/ (e.g. 01-apperror.md)")
    parser.add_argument("--resequence", "-r", action="store_true", help="Re-sequence numbered prefixes in plans directories")
    parser.add_argument("--clean-subtasks", action="store_true", help="Remove subtask directories for completed plans")
    parser.add_argument("--backup", "-b", action="store_true", help="Create a timestamped safety backup branch")
    parser.add_argument("--dry-run", "-d", action="store_true", help="Preview consolidation without making changes")
    parser.add_argument("--force", "-f", action="store_true", help="Bypass interactive confirmation prompt")
    args = parser.parse_args()

    if args.backup:
        create_backup_branch()
        sys.exit(0)

    if args.archive:
        sys.exit(archive_plan(args.archive, is_force_mode=args.force, is_dry_run_mode=args.dry_run))

    if args.clean_subtasks:
        subtasks = clean_superseded_subtasks(is_fix_mode=not args.dry_run)
        print(f"🗑️ Cleaned {len(subtasks)} subtask folder(s).")
        sys.exit(0)

    if args.resequence:
        resequence_directory_plans(PENDING_DIR, is_fix_mode=not args.dry_run)
        resequence_directory_plans(COMPLETED_DIR, is_fix_mode=not args.dry_run)
        sync_plans_index()
        sys.exit(0)

    sys.exit(run_full_consolidation(
        is_force_mode=args.force,
        is_dry_run_mode=args.dry_run,
        is_auto_backup=True
    ))

if __name__ == "__main__":
    main()
