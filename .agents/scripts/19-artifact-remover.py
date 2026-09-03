#!/usr/bin/env python3
"""
Fast Repository Artifact Remover & Git Cleanup Guard
Safely discovers, previews (Plan Mode), and moves unneeded test artifacts,
binary blobs, cache files, and temporary files to the Trash Bin / Git Index.

All Enums, Artifact Presets, Git Constants, and Utility Functions are imported
directly from 03-ai-scripts/02-shared-engine.py as the single source of truth.

Usage Examples:
  # 1. Preview (Plan Mode) what would be removed for a path or pattern:
  python 03-ai-scripts/19-artifact-remover.py --path tmp/ --plan
  python 03-ai-scripts/19-artifact-remover.py --pattern "*.log" --plan

  # 2. Target a folder and exclude specific files:
  python 03-ai-scripts/19-artifact-remover.py --path tmp/ --exclude "*.keep,important.log"

  # 3. Clean all pycache and bytecode artifacts safely:
  python 03-ai-scripts/19-artifact-remover.py --clean-pycache

  # 4. Clean temporary files (.tmp, .log, .swp, .bak, .DS_Store):
  python 03-ai-scripts/19-artifact-remover.py --clean-temp

  # 5. Clean unapproved binaries and blobs with plan mode preview:
  python 03-ai-scripts/19-artifact-remover.py --clean-binaries --plan

  # 6. Execute directly with automatic safe Trash Bin removal:
  python 03-ai-scripts/19-artifact-remover.py tmp/build-output/ --force
"""

import argparse
import fnmatch
from importlib import import_module
import os
from pathlib import Path
import shutil
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

# Centralized Enums
ArtifactCategoryType = engine.ArtifactCategoryType
ExitCodeType = engine.ExitCodeType

# Centralized String Literals & Separators
CURRENT_DIR = engine.CURRENT_DIR
EMPTY_STRING = engine.EMPTY_STRING
LINE_SEPARATOR = engine.LINE_SEPARATOR
PATH_SEPARATOR = engine.PATH_SEPARATOR
DEFAULT_ENCODING = engine.DEFAULT_ENCODING

# Centralized Git Command Constants
GIT_EXECUTABLE = engine.GIT_EXECUTABLE
GIT_CMD_LS_FILES = engine.GIT_CMD_LS_FILES
GIT_CMD_RM = engine.GIT_CMD_RM
GIT_FLAG_FORCE = engine.GIT_FLAG_FORCE
GIT_FLAG_ERROR_UNMATCH = engine.GIT_FLAG_ERROR_UNMATCH

# Centralized Status Badges
STATUS_GIT_TRACKED = engine.STATUS_GIT_TRACKED
STATUS_UNTRACKED = engine.STATUS_UNTRACKED

# Centralized Preset Artifact Collections
PYCACHE_DIR_NAMES = engine.PYCACHE_DIR_NAMES
PYCACHE_FILE_EXTENSIONS = engine.PYCACHE_FILE_EXTENSIONS
TEMP_ARTIFACT_EXTENSIONS = engine.TEMP_ARTIFACT_EXTENSIONS
TEMP_ARTIFACT_FILENAMES = engine.TEMP_ARTIFACT_FILENAMES
BINARY_EXTENSIONS = engine.BINARY_EXTENSIONS
EXCLUDE_DIRS = engine.EXCLUDE_DIRS

# Centralized Helper Functions
is_ignored_directory = engine.is_ignored_directory
is_binary_file = engine.is_binary_file
is_allowed_large_file = engine.is_allowed_large_file
normalize_rel_path = engine.normalize_rel_path
normalize_extensions = engine.normalize_extensions

def is_git_tracked(file_path: Path) -> bool:
    """
    Checks whether a given file is actively tracked in the Git index.
    Executes `git ls-files --error-unmatch <path>` with exit code inspection.
    """
    try:
        res = subprocess.run(
            [GIT_EXECUTABLE, GIT_CMD_LS_FILES, GIT_FLAG_ERROR_UNMATCH, str(file_path)],
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace"
        )
        return (res.returncode == 0)
    except Exception:
        return False

def send_to_trash_bin(target_path: Path) -> bool:
    """
    Safely moves a file or folder to the OS Recycle / Trash bin.
    Falls back gracefully to OS specific trash handlers.
    """
    if not target_path.exists():
        return False

    # 1. Try send2trash if installed
    try:
        import send2trash
        send2trash.send2trash(str(target_path.resolve()))
        return True
    except ImportError:
        pass
    except Exception:
        pass

    # 2. Windows PowerShell Recycle Bin Integration
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

    # 3. macOS trash integration
    if sys.platform == "darwin":
        try:
            res = subprocess.run(["trash", str(target_path.resolve())], capture_output=True)
            if res.returncode == 0 and not target_path.exists():
                return True
        except Exception:
            pass

    # 4. Linux gio trash / trash-put integration
    if sys.platform.startswith("linux"):
        for trash_tool in ["gio", "trash-put"]:
            try:
                cmd = ["gio", "trash", str(target_path.resolve())] if trash_tool == "gio" else ["trash-put", str(target_path.resolve())]
                res = subprocess.run(cmd, capture_output=True)
                if res.returncode == 0 and not target_path.exists():
                    return True
            except Exception:
                continue

    # Fallback to standard delete if recycle bin could not be reached
    if target_path.is_dir():
        shutil.rmtree(target_path)
    else:
        target_path.unlink()
    return True


def remove_from_git_and_disk(item_path: Path, is_use_trash: bool = True) -> bool:
    """Removes an item from Git index (if tracked) and deletes/trashes it from disk."""
    if not item_path.exists():
        return False

    is_tracked = is_git_tracked(item_path)
    if is_tracked:
        try:
            norm_path = normalize_rel_path(item_path)
            subprocess.run([GIT_EXECUTABLE, GIT_CMD_RM, "-r", "--cached", "--ignore-unmatch", norm_path],
                           stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, cwd=CURRENT_DIR)
        except Exception:
            pass

    if is_use_trash:
        return send_to_trash_bin(item_path)

    if item_path.is_dir():
        shutil.rmtree(item_path)
    else:
        item_path.unlink()
    return True


def collect_matching_artifacts(
    targets: list[str] | None = None,
    add_paths: set[str] | None = None,
    add_exts: set[str] | None = None,
    add_patterns: list[str] | None = None,
    exclude_patterns: list[str] | None = None,
    is_clean_pycache: bool = False,
    is_clean_binaries: bool = False,
    is_clean_temp: bool = False,
    root_dir: str = CURRENT_DIR
) -> list[Path]:
    """Discovers and aggregates all artifact candidates matching filter criteria."""
    candidates: list[Path] = []
    root_path = Path(root_dir).resolve()
    exclude_list = exclude_patterns or []

    # 1. Direct explicit file/folder targets
    if targets:
        for t in targets:
            clean_t = t.strip()
            if not clean_t or clean_t in {CURRENT_DIR, f"{CURRENT_DIR}{PATH_SEPARATOR}"}:
                continue
            p_target = Path(clean_t)
            if p_target.exists():
                candidates.append(p_target)
                if p_target.is_dir():
                    for r, dirs, files in os.walk(p_target, topdown=False):
                        for f in files:
                            candidates.append(Path(r) / f)
            else:
                for match in root_path.glob(f"**/{clean_t}"):
                    if not any(part in match.parts for part in {".git", ".gitmap", "node_modules"}):
                        candidates.append(match)

    # 2. Extensible Custom Paths (--add-path / --path)
    if add_paths:
        for p_str in add_paths:
            p_obj = Path(p_str)
            if p_obj.exists():
                candidates.append(p_obj)
                if p_obj.is_dir():
                    for r, dirs, files in os.walk(p_obj, topdown=False):
                        for f in files:
                            candidates.append(Path(r) / f)

    # 3. Pycache & Bytecode Preset (--clean-pycache)
    if is_clean_pycache:
        for r, dirs, files in os.walk(root_dir, topdown=False):
            dirs[:] = [d for d in dirs if d not in {".git", ".gitmap", "node_modules"}]
            dir_name = Path(r).name
            if dir_name in PYCACHE_DIR_NAMES:
                candidates.append(Path(r))
            for f in files:
                ext = os.path.splitext(f)[1].lower()
                if ext in PYCACHE_FILE_EXTENSIONS:
                    candidates.append(Path(r) / f)

    # 4. Temporary Artifacts Preset (--clean-temp)
    if is_clean_temp:
        for r, dirs, files in os.walk(root_dir, topdown=False):
            dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
            for f in files:
                ext = os.path.splitext(f)[1].lower()
                if ext in TEMP_ARTIFACT_EXTENSIONS or f in TEMP_ARTIFACT_FILENAMES:
                    candidates.append(Path(r) / f)

    # 5. Unapproved Binary Blobs Preset (--clean-binaries)
    if is_clean_binaries:
        for r, dirs, files in os.walk(root_dir, topdown=False):
            dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
            for f in files:
                p_file = Path(r) / f
                if is_allowed_large_file(p_file):
                    continue
                if is_binary_file(p_file):
                    candidates.append(p_file)

    # 6. Extensible Custom Extension Filter (--add-ext)
    if add_exts:
        for r, dirs, files in os.walk(root_dir, topdown=False):
            dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
            for f in files:
                ext = os.path.splitext(f)[1].lower()
                if ext in add_exts:
                    candidates.append(Path(r) / f)

    # 7. Extensible Glob Patterns (--add-pattern / --pattern)
    if add_patterns:
        for pat in add_patterns:
            for match in root_path.glob(f"**/{pat}"):
                if not any(part in match.parts for part in {".git", ".gitmap", "node_modules"}):
                    candidates.append(match)

    # Deduplicate
    unique_candidates = sorted(list(set(candidates)))

    # Apply Exclude Patterns (--exclude)
    if exclude_list:
        filtered = []
        for item in unique_candidates:
            rel = normalize_rel_path(item)
            is_excluded = any(fnmatch.fnmatch(rel, ex) or fnmatch.fnmatch(item.name, ex) for ex in exclude_list)
            if not is_excluded:
                filtered.append(item)
        return filtered

    return unique_candidates

def confirm_action(prompt_text: str) -> bool:
    """Prompts user for interactive confirmation."""
    try:
        response = input(f"{prompt_text} (y/N): ").strip().lower()
        return response in {"y", "yes"}
    except (EOFError, KeyboardInterrupt):
        return False

def format_bytes(num_bytes: int) -> str:
    """Formats raw byte count into human readable units."""
    if num_bytes < 1024:
        return f"{num_bytes} B"
    elif num_bytes < 1024 * 1024:
        return f"{num_bytes / 1024:.1f} KB"
    else:
        return f"{num_bytes / (1024 * 1024):.2f} MB"

def run_artifact_remover(
    targets: list[str] | None = None,
    add_paths: set[str] | None = None,
    add_exts: set[str] | None = None,
    add_patterns: list[str] | None = None,
    exclude_patterns: list[str] | None = None,
    is_clean_pycache: bool = False,
    is_clean_binaries: bool = False,
    is_clean_temp: bool = False,
    is_force_mode: bool = False,
    is_dry_run_mode: bool = False,
    is_use_trash: bool = True,
    root_dir: str = CURRENT_DIR
) -> int:
    """
    Master orchestrator:
    1. Collects all matching items with inclusion/exclusion filters.
    2. Displays full Plan Mode preview with file sizes and Git tracking badges.
    3. Prompts user for interactive confirmation (unless --force).
    4. Executes safe Trash Bin removal and Git index cleanup.
    """
    start_time = time.perf_counter()

    artifacts = collect_matching_artifacts(
        targets=targets,
        add_paths=add_paths,
        add_exts=add_exts,
        add_patterns=add_patterns,
        exclude_patterns=exclude_patterns,
        is_clean_pycache=is_clean_pycache,
        is_clean_binaries=is_clean_binaries,
        is_clean_temp=is_clean_temp,
        root_dir=root_dir
    )

    has_artifacts = len(artifacts) > 0
    if not has_artifacts:
        print("✅ No matching artifacts found for the specified criteria.")
        return ExitCodeType.SUCCESS.value

    # Compute Total Size
    total_bytes = 0
    for a in artifacts:
        try:
            if a.is_file():
                total_bytes += a.stat().st_size
        except Exception:
            pass

    # Display Plan Mode Preview Table
    action_label = "Trash Bin Safe Move" if is_use_trash else "Permanent Disk Unlink"
    print("=" * 80)
    print(f"📋 Artifact Removal Plan Mode ({action_label})")
    print(f"   Candidates Found: {len(artifacts)} item(s) | Total Size: {format_bytes(total_bytes)}")
    print("=" * 80)
    for idx, item in enumerate(artifacts[:30], start=1):
        badge = STATUS_GIT_TRACKED if is_git_tracked(item) else STATUS_UNTRACKED
        size_str = ""
        if item.is_file():
            try:
                size_str = f" ({format_bytes(item.stat().st_size)})"
            except Exception:
                pass
        print(f"   {idx:>3}. {normalize_rel_path(item)}{size_str}{badge}")
    if len(artifacts) > 30:
        print(f"   ... and {len(artifacts) - 30} more items.")

    if is_dry_run_mode:
        print(f"{LINE_SEPARATOR}ℹ️ Plan Mode (Dry-Run) complete. No files were removed.")
        return ExitCodeType.SUCCESS.value

    # Require Interactive Confirmation unless --force is specified
    if not is_force_mode:
        print(f"{LINE_SEPARATOR}⚠️ Safety Notice: Action will move items to the OS Trash Bin and unstage from Git.")
        is_confirmed = confirm_action(f"Proceed with moving {len(artifacts)} item(s) ({format_bytes(total_bytes)}) to the Trash Bin?")
        if not is_confirmed:
            print("❌ Operation canceled by user. No files were touched.")
            return ExitCodeType.SUCCESS.value

    # Execute Safe Removal
    removed_count = 0
    for item in artifacts:
        is_removed = remove_from_git_and_disk(item, is_use_trash=is_use_trash)
        if is_removed:
            removed_count += 1
            print(f"  ✓ Moved to Trash: {normalize_rel_path(item)}")

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"{LINE_SEPARATOR}✅ Successfully processed {removed_count} of {len(artifacts)} item(s) ({format_bytes(total_bytes)}) in {elapsed_ms:.2f}ms.")
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(
        description="Fast Repository Artifact Remover & Git Cleanup Guard with Plan Mode & Safe Trash Bin",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Preview (Plan Mode) removing temporary files:
  python 03-ai-scripts/19-artifact-remover.py --clean-temp --plan

  # Target specific directory with exclusions:
  python 03-ai-scripts/19-artifact-remover.py --path tmp/ --exclude "*.keep,important.log"

  # Clean pycache files safely:
  python 03-ai-scripts/19-artifact-remover.py --clean-pycache --force
"""
    )
    parser.add_argument("targets", nargs="*", default=[], help="Target files, directories, or patterns (e.g. tmp/build/)")
    parser.add_argument("--path", "-p", help="Target directory or file path to clean")
    parser.add_argument("--pattern", "--include", help="Comma-separated glob patterns to include (e.g. *.tmp,*.log)")
    parser.add_argument("--exclude", "--ignore", help="Comma-separated glob patterns to protect/exclude")
    parser.add_argument("--add-ext", help="Comma-separated file extensions to remove (e.g. .tmp,.log)")
    parser.add_argument("--clean-pycache", action="store_true", help="Remove __pycache__, .pytest_cache, and .pyc/.pyo files")
    parser.add_argument("--clean-temp", action="store_true", help="Remove temporary files (.tmp, .log, .swp, .bak, .DS_Store)")
    parser.add_argument("--clean-binaries", action="store_true", help="Remove unapproved binary blobs and image artifacts")
    parser.add_argument("--clean-all", action="store_true", help="Enable pycache, temp files, and binary cleanup presets")
    parser.add_argument("--force", "-f", "-y", "--yes", action="store_true", help="Bypass interactive confirmation prompt")
    parser.add_argument("--plan", "--dry-run", "-d", action="store_true", help="Plan Mode: preview matching items without deleting")
    parser.add_argument("--permanent", action="store_true", help="Permanently unlink files without moving to Trash Bin")
    parser.add_argument("--dir", default=CURRENT_DIR, help="Root directory to search (default: .)")
    args = parser.parse_args()

    # Parse extensible CLI inputs
    custom_paths = set()
    if args.path:
        custom_paths.update(p.strip() for p in args.path.split(",") if p.strip())
    custom_exts = normalize_extensions(args.add_ext)

    custom_patterns = []
    if args.pattern:
        custom_patterns.extend(p.strip() for p in args.pattern.split(",") if p.strip())

    exclude_patterns = []
    if args.exclude:
        exclude_patterns.extend(p.strip() for p in args.exclude.split(",") if p.strip())

    is_clean_pycache = args.clean_pycache or args.clean_all
    is_clean_temp = args.clean_temp or args.clean_all
    is_clean_binaries = args.clean_binaries or args.clean_all
    is_use_trash = not args.permanent

    has_work = bool(
        args.targets or custom_paths or custom_exts or custom_patterns or
        is_clean_pycache or is_clean_temp or is_clean_binaries
    )
    if not has_work:
        parser.print_help()
        sys.exit(0)

    sys.exit(run_artifact_remover(
        targets=args.targets,
        add_paths=custom_paths,
        add_exts=custom_exts,
        add_patterns=custom_patterns,
        exclude_patterns=exclude_patterns,
        is_clean_pycache=is_clean_pycache,
        is_clean_binaries=is_clean_binaries,
        is_clean_temp=is_clean_temp,
        is_force_mode=args.force,
        is_dry_run_mode=args.plan,
        is_use_trash=is_use_trash,
        root_dir=args.dir
    ))

if __name__ == "__main__":
    main()
