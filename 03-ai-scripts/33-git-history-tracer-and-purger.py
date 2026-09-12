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

def get_existing_files_as_entries(path_filter: str) -> list[DeletedFileEntry]:
    """Scans existing workspace files matching path_filter when no historical deletions exist."""
    if not path_filter or not Path(path_filter).exists():
        return []
    p = Path(path_filter)
    candidates = [p] if p.is_file() else list(p.rglob("*"))
    entries = []
    for c in candidates:
        if c.is_file():
            rel = normalize_rel_path(str(c))
            entries.append(DeletedFileEntry(rel, "HEAD", "Current", "Active", "Workspace File"))
    return entries

def trace_deleted_files(path_filter: str = "") -> list[DeletedFileEntry]:
    """Traverses Git history and active tree to catalog target files."""
    cmd = ["log", "--diff-filter=D", "--name-only", "--pretty=format:COMMIT:%H|%an|%ad|%s", "--date=short"]
    if path_filter:
        cmd.extend(["--", path_filter])
    res = run_git_raw(cmd)
    tracked = get_current_tracked_files()
    hist_entries = parse_deleted_entries(res.stdout.splitlines(), tracked) if res.returncode == 0 else []
    if hist_entries:
        return hist_entries
    return get_existing_files_as_entries(path_filter)

def is_pattern_matched(pattern: str, text: str) -> bool:
    """Safely evaluates pattern against text using fnmatch and regex."""
    if not pattern:
        return True
    norm_text = text.lower()
    norm_pat = pattern.lower()
    if fnmatch.fnmatch(norm_text, norm_pat) or fnmatch.fnmatch(norm_text, f"*{norm_pat}*"):
        return True
    try:
        return bool(re.search(pattern, text, re.IGNORECASE))
    except (re.error, ValueError):
        return False

def is_entry_matching(entry: DeletedFileEntry, path_prefix: str, ext: str, pattern: str) -> bool:
    """Evaluates whether an entry satisfies path, extension, and pattern filters."""
    norm_entry = entry.path.lower()
    clean_pfx = path_prefix.strip("/\\").lower()
    if clean_pfx and not (norm_entry.startswith(clean_pfx + "/") or norm_entry == clean_pfx):
        return False
    if ext and not norm_entry.endswith(ext.lower()):
        return False
    if not is_pattern_matched(pattern, entry.path):
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

def create_temp_backup_dir() -> Path:
    """Creates a timestamped backup directory in the OS temp directory."""
    tag = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%d-%H%M%S")
    backup_path = Path(tempfile.gettempdir()) / "gitmap" / "purge" / f"gitmap-deleted-backup-{tag}"
    backup_path.mkdir(parents=True, exist_ok=True)
    return backup_path

def backup_file_to_temp(src_path: Path, backup_root: Path) -> bool:
    """Copies an existing workspace file to the temp backup directory."""
    if not src_path.is_file():
        return False
    dst_path = backup_root / src_path
    dst_path.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(src_path, dst_path)
    return True

def backup_historical_blobs_to_temp(entries: list[DeletedFileEntry], targets: set[str], backup_root: Path) -> int:
    """Extracts historical blobs into temp backup directory before purge."""
    count = 0
    for e in entries:
        if e.path not in targets or e.commit_hash == "HEAD":
            continue
        res = run_git_raw(["show", f"{e.commit_hash}~1:{e.path}"])
        if res.returncode == 0:
            dst = backup_root / e.path
            dst.parent.mkdir(parents=True, exist_ok=True)
            dst.write_text(res.stdout, encoding="utf-8", errors="replace")
            count += 1
    return count

def send_file_to_recycle_bin(path: Path) -> bool:
    """Deletes a file or directory using the OS Recycle Bin or Trash."""
    if not path.exists():
        return False
    try:
        import send2trash
        send2trash.send2trash(str(path))
        return True
    except Exception:
        return False

def delete_workspace_files(targets: list[str], backup_root: Path) -> int:
    """Copies workspace files to temp backup and moves them to Recycle Bin."""
    deleted_count = 0
    for rel_p in targets:
        p = Path(rel_p)
        if p.exists():
            backup_file_to_temp(p, backup_root)
            is_recycled = send_file_to_recycle_bin(p)
            if is_recycled:
                deleted_count += 1
                print(f"  🗑️ Recycled: {rel_p}")
    return deleted_count

def validate_purge_preconditions(targets: list[str], is_confirm: bool) -> bool:
    """Validates targets, working tree cleanliness, and user confirmation."""
    if not targets:
        print("ℹ️ No target files selected for purge.")
        return False
    if not verify_clean_working_tree():
        print("❌ Working tree is dirty. Commit or stash changes before running history purge.")
        return False
    if not confirm_purge_action(len(targets), is_confirm):
        print("❌ Purge aborted by user.")
        return False
    return True

def run_purge_execution(targets: list[str], origin_url: str, backup_branch: str, backup_root: Path) -> bool:
    """Executes filter-repo and outputs recovery and rollback details."""
    with tempfile.NamedTemporaryFile("w", delete=False, encoding="utf-8") as tf:
        tf.write("\n".join(targets) + "\n")
        tmp_name = tf.name
    is_ok = execute_filter_repo(Path(tmp_name))
    Path(tmp_name).unlink(missing_ok=True)
    if origin_url:
        run_git_raw(["remote", "add", "origin", origin_url])
    if is_ok:
        print(f"\n✅ History purge complete! Eradicated {len(targets)} files.")
        print(f"🛡️ Local OS Temp Backup Path: {backup_root}")
        print(f"Rollback recipe (from temp): copy files from {backup_root} to repository")
        print(f"Rollback recipe (from Git):  git reset --hard {backup_branch}")
    return is_ok

def purge_selected_files(entries: list[DeletedFileEntry], targets: list[str], is_confirm: bool) -> bool:
    """Permanently purges historical files with temp backup and Recycle Bin."""
    if not validate_purge_preconditions(targets, is_confirm):
        return False
    backup_branch = create_safety_backup_branch()
    backup_root = create_temp_backup_dir()
    delete_workspace_files(targets, backup_root)
    backup_historical_blobs_to_temp(entries, set(targets), backup_root)
    remote_res = run_git_raw(["remote", "get-url", "origin"])
    origin_url = remote_res.stdout.strip() if remote_res.returncode == 0 else ""
    return run_purge_execution(targets, origin_url, backup_branch, backup_root)

def resolve_preset_target(args: argparse.Namespace) -> tuple[str, str]:
    """Resolves specific folder preset shortcuts to (path, ext) tuples."""
    preset_map = {
        "is_spec_25_audit": ("spec/21-app/25-app-spec-audit", ""),
        "is_spec_audit": ("spec/19-main-worker-service/audit", ""),
        "is_lovable_subtasks": (".lovable/plans/subtasks", ""),
        "is_lovable_audits": (".lovable/audits", ""),
        "is_lovable_md": (".lovable", ".md"),
        "is_lovable": (".lovable", ""),
        "is_spec_md": ("spec", ".md"),
        "is_spec": ("spec", ""),
    }
    for flag, target in preset_map.items():
        if getattr(args, flag, False):
            return target
    return "", ""

def apply_presets(args: argparse.Namespace) -> tuple[str, str, str]:
    """Resolves path, extension, and pattern filters from CLI presets, targets, or flags."""
    preset_path, preset_ext = resolve_preset_target(args)
    raw_path = preset_path or args.target or args.path or ""
    clean_path = "" if raw_path in (".", "/", "\\", "root") else raw_path
    ext_val = preset_ext or args.ext or ""
    pat_val = "*audit*" if getattr(args, "is_audit", False) else (args.pattern or "")
    return clean_path, ext_val, pat_val

def resolve_exclusion_set(args: argparse.Namespace, entries: list[DeletedFileEntry]) -> set[int]:
    """Resolves exclusion indices from flags, patterns, or interactive prompt."""
    excl_set = parse_exclusion_indices(args.exclude)
    if args.exclude_pattern:
        for i, e in enumerate(entries, start=1):
            if is_pattern_matched(args.exclude_pattern, e.path):
                excl_set.add(i)
    if not args.is_confirm and not args.exclude and not args.exclude_pattern:
        if args.is_restore or args.is_delete or args.is_purge:
            excl_set = prompt_interactive_exclusion(entries)
    return excl_set

def add_action_args(parser: argparse.ArgumentParser) -> None:
    """Registers core operational action flags."""
    parser.add_argument("target", nargs="?", default="", help="Target folder or file path (or . for root).")
    parser.add_argument("--list", "-l", action="store_true", dest="is_list", help="Preview deleted files (default).")
    parser.add_argument("--restore", "-r", action="store_true", dest="is_restore", help="Restore deleted files.")
    parser.add_argument("--delete", "-d", action="store_true", dest="is_delete", help="Recycle existing files to Recycle Bin.")
    parser.add_argument("--purge", "-p", action="store_true", dest="is_purge", help="Purge files from all Git history.")

def add_preset_args(parser: argparse.ArgumentParser) -> None:
    """Registers folder and scope preset flags."""
    parser.add_argument("--spec-25-audit", action="store_true", dest="is_spec_25_audit", help="Preset: spec/21-app/25-app-spec-audit")
    parser.add_argument("--spec-audit", action="store_true", dest="is_spec_audit", help="Preset: spec/19-main-worker-service/audit")
    parser.add_argument("--audit", action="store_true", dest="is_audit", help="Preset: all audit files repo-wide (*audit*)")
    parser.add_argument("--lovable-subtasks", action="store_true", dest="is_lovable_subtasks", help="Preset: .lovable/plans/subtasks/")
    parser.add_argument("--lovable-audits", action="store_true", dest="is_lovable_audits", help="Preset: .lovable/audits/")
    parser.add_argument("--lovable", action="store_true", dest="is_lovable", help="Preset: all files in .lovable/")
    parser.add_argument("--lovable-md", action="store_true", dest="is_lovable_md", help="Preset: .md files in .lovable/")
    parser.add_argument("--spec", action="store_true", dest="is_spec", help="Preset: all files in spec/")
    parser.add_argument("--spec-md", action="store_true", dest="is_spec_md", help="Preset: .md files in spec/")

def add_filter_and_option_args(parser: argparse.ArgumentParser) -> None:
    """Registers filter and operational option flags."""
    parser.add_argument("--path", type=str, default="", help="Filter by folder path prefix.")
    parser.add_argument("--ext", type=str, default="", help="Filter by file extension.")
    parser.add_argument("--pattern", type=str, default="", help="Filter by glob or regex.")
    parser.add_argument("--exclude", type=str, default="", help="Indices or ranges to exclude (e.g. 1, 3-5).")
    parser.add_argument("--exclude-pattern", type=str, default="", help="Exclude paths matching pattern.")
    parser.add_argument("--restore-to", type=str, default="", help="Destination directory for restoration.")
    parser.add_argument("-y", "--confirm", action="store_true", dest="is_confirm", help="Bypass confirmation prompt.")
    parser.add_argument("--backup-only", action="store_true", dest="is_backup_only", help="Create safety backup branch only.")

CLI_EPILOG = """Examples:
  # 1. Preview files in spec folder 25 (spec/21-app/25-app-spec-audit)
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --spec-25-audit

  # 2. Preview deleted files in spec/19-main-worker-service/audit
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --spec-audit

  # 3. Preview all deleted audit files repo-wide
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --audit

  # 4. Delete existing audit files to Recycle Bin with temp backup
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --spec-25-audit --delete

  # 5. Deep purge files from Git history with temp backup & safety branch
  python 03-ai-scripts/33-git-history-tracer-and-purger.py --spec-25-audit --purge
"""

def parse_cli_args() -> argparse.Namespace:
    """Configures and parses CLI arguments."""
    parser = argparse.ArgumentParser(
        description="Git Historical Deleted Files Tracer, Restorer & Deep Purger.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=CLI_EPILOG,
    )
    add_action_args(parser)
    add_preset_args(parser)
    add_filter_and_option_args(parser)
    return parser.parse_args()

def handle_restore_or_delete(args: argparse.Namespace, filtered: list[DeletedFileEntry], targets: list[str], excl_set: set[int]) -> int:
    """Handles restore or recycle bin deletion action."""
    if args.is_restore:
        dest_dir = Path(args.restore_to) if args.restore_to else None
        count = restore_selected_files(filtered, excl_set, dest_dir)
        print(f"\n🎉 Successfully restored {count} files.")
        return ExitCodeType.SUCCESS.value
    backup_root = create_temp_backup_dir()
    count = delete_workspace_files(targets, backup_root)
    print(f"\n🗑️ Recycled {count} files to Recycle Bin.")
    print(f"🛡️ Local OS Temp Backup Path: {backup_root}")
    print(f"Rollback recipe: copy files from {backup_root} to workspace")
    return ExitCodeType.SUCCESS.value

def dispatch_action(args: argparse.Namespace, filtered: list[DeletedFileEntry], excl_set: set[int]) -> int:
    """Executes selected CLI action (restore, delete, purge, or list)."""
    targets = [e.path for i, e in enumerate(filtered, start=1) if i not in excl_set]
    if args.is_restore or args.is_delete:
        return handle_restore_or_delete(args, filtered, targets, excl_set)
    if args.is_purge:
        is_purged = purge_selected_files(filtered, targets, args.is_confirm)
        return ExitCodeType.SUCCESS.value if is_purged else ExitCodeType.TOOL_ERROR.value
    return ExitCodeType.SUCCESS.value

def main() -> int:
    """Main CLI entrypoint."""
    args = parse_cli_args()
    if args.is_backup_only:
        create_safety_backup_branch()
        return ExitCodeType.SUCCESS.value
    path_val, ext_val, pat_val = apply_presets(args)
    entries = trace_deleted_files(path_val)
    filtered = filter_entries(entries, path_val, ext_val, pat_val)
    if not filtered:
        print("ℹ️ No historically deleted or matching files found.")
        return ExitCodeType.SUCCESS.value
    excl_set = resolve_exclusion_set(args, filtered)
    render_preflight_table(filtered, excl_set)
    return dispatch_action(args, filtered, excl_set)

if __name__ == "__main__":
    sys.exit(main())
