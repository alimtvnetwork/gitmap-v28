#!/usr/bin/env python3
"""
33-git-history-tracer-and-purger.py: Traces historically deleted files in Git,
provides pre-flight inspection with selective exclusion, restores files to workspace,
and deeply purges deleted files across all commits, branches, and tags.
"""

import argparse
from dataclasses import dataclass
import datetime
import fnmatch
from importlib import import_module
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
normalize_rel_path = engine.normalize_rel_path
read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf

@dataclass
class DeletedFileEntry:
    path: str
    commit_hash: str
    author: str
    date: str
    subject: str

def run_git_raw(args: list[str]) -> subprocess.CompletedProcess[str]:
    """Executes a git command with UTF-8 encoding and output capture."""
    git_bin = shutil.which("git") or "git"
    return subprocess.run([git_bin] + args, capture_output=True, text=True, encoding="utf-8", errors="replace")

def get_current_tracked_files() -> set[str]:
    """Returns all paths currently tracked in Git HEAD."""
    res = run_git_raw(["ls-files"])
    if res.returncode != 0:
        return set()
    return {normalize_rel_path(line.strip()) for line in res.stdout.splitlines() if line.strip()}

def parse_commit_header(line: str) -> tuple[str, str, str, str]:
    """Parses COMMIT:<hash>|<author>|<date>|<subject> header line."""
    parts = line[7:].split("|", 3)
    c_hash = parts[0] if len(parts) > 0 else ""
    author = parts[1] if len(parts) > 1 else ""
    date_str = parts[2] if len(parts) > 2 else ""
    subj = parts[3] if len(parts) > 3 else ""
    return c_hash, author, date_str, subj

def parse_deleted_entries(lines: list[str], tracked_files: set[str]) -> list[DeletedFileEntry]:
    """Parses git log output into DeletedFileEntry list, skipping existing files."""
    entries: list[DeletedFileEntry] = []
    seen_paths: set[str] = set()
    cur_hash, cur_author, cur_date, cur_subj = "", "", "", ""
    for line in lines:
        if line.startswith("COMMIT:"):
            cur_hash, cur_author, cur_date, cur_subj = parse_commit_header(line)
            continue
        rel_p = normalize_rel_path(line.strip())
        is_candidate = bool(rel_p and rel_p not in tracked_files and rel_p not in seen_paths)
        if is_candidate:
            seen_paths.add(rel_p)
            entries.append(DeletedFileEntry(rel_p, cur_hash, cur_author, cur_date, cur_subj))
    return entries

def trace_deleted_files(path_filter: str = "") -> list[DeletedFileEntry]:
    """Traverses Git history to catalog files removed from current HEAD."""
    cmd = ["log", "--diff-filter=D", "--name-only", "--pretty=format:COMMIT:%H|%an|%ad|%s", "--date=short"]
    if path_filter:
        cmd.extend(["--", path_filter])
    res = run_git_raw(cmd)
    if res.returncode != 0:
        return []
    tracked = get_current_tracked_files()
    return parse_deleted_entries(res.stdout.splitlines(), tracked)

def is_entry_matching(entry: DeletedFileEntry, path_prefix: str, ext: str, pattern: str) -> bool:
    """Evaluates whether an entry satisfies path, extension, and pattern filters."""
    norm_entry = entry.path.lower()
    if path_prefix and not norm_entry.startswith(path_prefix.lower().rstrip("/") + "/"):
        if norm_entry != path_prefix.lower():
            return False
    if ext and not norm_entry.endswith(ext.lower()):
        return False
    if pattern and not fnmatch.fnmatch(norm_entry, pattern.lower()):
        if not re.search(pattern, entry.path, re.IGNORECASE):
            return False
    return True

def filter_entries(entries: list[DeletedFileEntry], path_pfx: str, ext: str, pat: str) -> list[DeletedFileEntry]:
    """Filters entries by path prefix, extension, and pattern."""
    return [e for e in entries if is_entry_matching(e, path_pfx, ext, pat)]

def parse_single_range(token: str, result_set: set[int]) -> None:
    """Parses a single number or range token (e.g. '3' or '5-8') into result set."""
    if "-" in token:
        parts = token.split("-", 1)
        if parts[0].isdigit() and parts[1].isdigit():
            start, end = int(parts[0]), int(parts[1])
            result_set.update(range(min(start, end), max(start, end) + 1))
        return
    if token.isdigit():
        result_set.add(int(token))

def parse_exclusion_indices(raw_input: str) -> set[int]:
    """Parses comma-separated indices and dash ranges into a set of 1-based indices."""
    result: set[int] = set()
    if not raw_input:
        return result
    for token in raw_input.replace(" ", "").split(","):
        parse_single_range(token, result)
    return result

def prompt_interactive_exclusion(entries: list[DeletedFileEntry]) -> set[int]:
    """Interactively prompts user for exclusion indices."""
    try:
        ans = input("\n👉 Enter indices or ranges to EXCLUDE (e.g. 1, 3, 5-8) [Enter to keep all]: ").strip()
        return parse_exclusion_indices(ans)
    except (EOFError, KeyboardInterrupt):
        return set()

def render_entry_row(idx: int, entry: DeletedFileEntry, is_excluded: bool) -> None:
    """Renders a single row in the pre-flight inspection table."""
    status_label = "🚫 EXCLUDED" if is_excluded else "🎯 TARGET  "
    short_hash = entry.commit_hash[:7] if len(entry.commit_hash) >= 7 else entry.commit_hash
    short_author = (entry.author[:14] + "..") if len(entry.author) > 16 else entry.author
    print(f"[{idx:3d}] {status_label} | {entry.date:10s} | {short_hash} | {short_author:16s} | {entry.path}")

def render_preflight_table(entries: list[DeletedFileEntry], excluded_indices: set[int]) -> None:
    """Renders formatted candidate table with header, items, and counts."""
    print("\n" + "=" * 90)
    print("              GIT HISTORICAL DELETED FILES PRE-FLIGHT INSPECTION              ")
    print("=" * 90)
    print("Idx   Status      | Date       | Commit  | Author           | Relative Path")
    print("-" * 90)
    for i, entry in enumerate(entries, start=1):
        is_excl = i in excluded_indices
        render_entry_row(i, entry, is_excl)
    print("-" * 90)
    active_count = len(entries) - len(excluded_indices.intersection(range(1, len(entries) + 1)))
    print(f"📊 Summary: {len(entries)} total traced | {len(excluded_indices)} excluded | {active_count} active targets\n")

def restore_single_file(entry: DeletedFileEntry, restore_root: Path | None) -> bool:
    """Restores a single file from commit~1:<path> to disk."""
    rev_spec = f"{entry.commit_hash}~1:{entry.path}"
    res = run_git_raw(["show", rev_spec])
    if res.returncode != 0:
        return False
    target_path = Path(entry.path) if restore_root is None else (restore_root / entry.path)
    target_path.parent.mkdir(parents=True, exist_ok=True)
    target_path.write_text(res.stdout, encoding="utf-8", errors="replace")
    return True

def restore_selected_files(entries: list[DeletedFileEntry], excluded_indices: set[int], restore_root: Path | None) -> int:
    """Restores all non-excluded entries to workspace."""
    success_count = 0
    for i, entry in enumerate(entries, start=1):
        if i in excluded_indices:
            continue
        is_restored = restore_single_file(entry, restore_root)
        if is_restored:
            success_count += 1
            dst = entry.path if restore_root is None else str(restore_root / entry.path)
            print(f"  ✓ Restored: {dst}")
        else:
            print(f"  ✗ Failed to restore: {entry.path}")
    return success_count

def create_safety_backup_branch() -> str:
    """Creates a timestamped safety backup branch before modifying history."""
    tag = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%d-%H%M%S")
    branch_name = f"backup/history-purge-{tag}"
    res = run_git_raw(["branch", branch_name])
    if res.returncode != 0:
        print(f"⚠️ Warning: Could not create local branch {branch_name}: {res.stderr.strip()}")
    else:
        print(f"🛡️ Safety Backup Branch Created: `{branch_name}`")
    return branch_name

def verify_clean_working_tree() -> bool:
    """Verifies that git working tree has no uncommitted changes."""
    res = run_git_raw(["status", "--porcelain"])
    return bool(res.returncode == 0 and not res.stdout.strip())

def confirm_purge_action(active_count: int, is_auto_confirm: bool) -> bool:
    """Prompts for explicit confirmation before rewriting git history."""
    if is_auto_confirm:
        return True
    print(f"⚠️ DANGER: You are about to permanently purge {active_count} files from ALL commits, branches, and tags!")
    try:
        val = input("Type 'I confirm' to execute irreversible history rewrite: ").strip()
        return bool(val == "I confirm")
    except (EOFError, KeyboardInterrupt):
        return False

def execute_filter_repo(paths_file: Path) -> bool:
    """Executes git filter-repo with paths file."""
    res = run_git_raw(["filter-repo", "--paths-from-file", str(paths_file), "--invert-paths", "--force"])
    if res.returncode != 0:
        print(f"❌ git filter-repo failed: {res.stderr.strip()}")
        return False
    return True

def purge_selected_files(entries: list[DeletedFileEntry], excluded_indices: set[int], is_auto_confirm: bool) -> bool:
    """Permanently purges non-excluded historical files using git filter-repo."""
    targets = [e.path for i, e in enumerate(entries, start=1) if i not in excluded_indices]
    if not targets:
        print("ℹ️ No target files selected for purge.")
        return True
    if not verify_clean_working_tree():
        print("❌ Working tree is dirty. Commit or stash changes before running history purge.")
        return False
    if not confirm_purge_action(len(targets), is_auto_confirm):
        print("❌ Purge aborted by user.")
        return False
    backup_branch = create_safety_backup_branch()
    remote_res = run_git_raw(["remote", "get-url", "origin"])
    origin_url = remote_res.stdout.strip() if remote_res.returncode == 0 else ""
    with tempfile.NamedTemporaryFile("w", delete=False, encoding="utf-8") as tf:
        tf.write("\n".join(targets) + "\n")
        tmp_name = tf.name
    is_ok = execute_filter_repo(Path(tmp_name))
    Path(tmp_name).unlink(missing_ok=True)
    if origin_url:
        run_git_raw(["remote", "add", "origin", origin_url])
    if is_ok:
        print(f"\n✅ History purge complete! Eradicated {len(targets)} files.")
        print(f"Rollback recipe (if needed): git reset --hard {backup_branch}")
    return is_ok

def apply_presets(args: argparse.Namespace) -> tuple[str, str]:
    """Resolves path and extension filters from CLI presets."""
    path_val = args.path or ""
    ext_val = args.ext or ""
    if args.is_lovable_md:
        path_val = ".lovable"
        ext_val = ".md"
    elif args.is_spec_md:
        path_val = "spec"
        ext_val = ".md"
    return path_val, ext_val

def resolve_exclusion_set(args: argparse.Namespace, entries: list[DeletedFileEntry]) -> set[int]:
    """Resolves exclusion indices from flags, patterns, or interactive prompt."""
    excl_set = parse_exclusion_indices(args.exclude)
    if args.exclude_pattern:
        for i, e in enumerate(entries, start=1):
            if fnmatch.fnmatch(e.path, args.exclude_pattern) or re.search(args.exclude_pattern, e.path):
                excl_set.add(i)
    if not args.is_confirm and not args.exclude and not args.exclude_pattern:
        if args.is_restore or args.is_purge:
            excl_set = prompt_interactive_exclusion(entries)
    return excl_set

def parse_cli_args() -> argparse.Namespace:
    """Configures and parses CLI arguments."""
    parser = argparse.ArgumentParser(
        description="Git Historical Deleted Files Tracer, Restorer & Deep Purger.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""Examples:
  # 1. Preview all deleted files under .lovable/
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --lovable-md

  # 2. Preview all deleted specs under spec/
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --spec-md

  # 3. Restore all deleted markdown files excluding items 1 to 5
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --lovable-md --restore --exclude 1-5

  # 4. Deeply purge deleted files across entire Git history
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --lovable-md --purge
"""
    )
    parser.add_argument("--list", "-l", action="store_true", dest="is_list", help="Preview deleted files (default).")
    parser.add_argument("--restore", "-r", action="store_true", dest="is_restore", help="Restore deleted files to workspace.")
    parser.add_argument("--purge", "-p", action="store_true", dest="is_purge", help="Purge files from all Git history.")
    parser.add_argument("--lovable-md", action="store_true", dest="is_lovable_md", help="Preset: deleted .md in .lovable/")
    parser.add_argument("--spec-md", action="store_true", dest="is_spec_md", help="Preset: deleted .md in spec/")
    parser.add_argument("--path", type=str, default="", help="Filter by folder path prefix.")
    parser.add_argument("--ext", type=str, default="", help="Filter by file extension.")
    parser.add_argument("--pattern", type=str, default="", help="Filter by glob or regex.")
    parser.add_argument("--exclude", type=str, default="", help="Comma-separated indices or ranges to exclude.")
    parser.add_argument("--exclude-pattern", type=str, default="", help="Exclude paths matching pattern.")
    parser.add_argument("--restore-to", type=str, default="", help="Destination directory for restoration.")
    parser.add_argument("-y", "--confirm", action="store_true", dest="is_confirm", help="Bypass confirmation prompt.")
    parser.add_argument("--backup-only", action="store_true", dest="is_backup_only", help="Create safety backup branch only.")
    return parser.parse_args()

def main() -> int:
    """Main CLI entrypoint."""
    args = parse_cli_args()
    if args.is_backup_only:
        create_safety_backup_branch()
        return ExitCodeType.SUCCESS.value
    path_val, ext_val = apply_presets(args)
    entries = trace_deleted_files(path_val)
    filtered = filter_entries(entries, path_val, ext_val, args.pattern)
    if not filtered:
        print("ℹ️ No historically deleted files found matching criteria.")
        return ExitCodeType.SUCCESS.value
    excl_set = resolve_exclusion_set(args, filtered)
    render_preflight_table(filtered, excl_set)
    if args.is_restore:
        dest_dir = Path(args.restore_to) if args.restore_to else None
        count = restore_selected_files(filtered, excl_set, dest_dir)
        print(f"\n🎉 Successfully restored {count} files.")
        return ExitCodeType.SUCCESS.value
    if args.is_purge:
        is_purged = purge_selected_files(filtered, excl_set, args.is_confirm)
        return ExitCodeType.SUCCESS.value if is_purged else ExitCodeType.TOOL_ERROR.value
    return ExitCodeType.SUCCESS.value

if __name__ == "__main__":
    sys.exit(main())
