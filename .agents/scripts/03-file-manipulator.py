#!/usr/bin/env python3
"""
Python File Manipulator CLI Tool
Handles mass file renaming (lowercasing), sequence fixing, and encoding normalization.
Zero-dependency, fast, git-aware, with multi-folder scoping and customizable extensions.

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/03-file-manipulator.py lowercase <path> [--except <patterns>]
  python 03-ai-scripts/03-file-manipulator.py fix-seq-files <path> [--order-by-time|--order-by-az] [--pin <mapping>]
  python 03-ai-scripts/03-file-manipulator.py fix-encoding <path> [--ext <extensions>]
"""

import argparse
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

process_repository_files = engine.process_repository_files
read_file_safe = engine.read_file_safe
write_file_lf = engine.write_file_lf
normalize_rel_path = engine.normalize_rel_path
normalize_extensions = engine.normalize_extensions
is_ignored_directory = engine.is_ignored_directory
is_ignored_path = engine.is_ignored_path
is_binary_file = engine.is_binary_file
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
UTF8_SIG_ENCODING = engine.UTF8_SIG_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

def run_git_mv_or_rename(src_path: Path, dst_path: Path) -> bool:
    """Attempts git mv first; falls back to standard filesystem rename with Windows case-hop."""
    if src_path == dst_path:
        return True
    try:
        res = subprocess.run(
            ["git", "mv", str(src_path), str(dst_path)],
            capture_output=True,
            text=True
        )
        if res.returncode == 0:
            return True
    except Exception:
        pass
    try:
        if src_path.name.lower() == dst_path.name.lower() and sys.platform.startswith("win"):
            temp_hop = src_path.with_name(f"{src_path.name}.tmp_case_{os.getpid()}_{int(time.time()*1000)}")
            shutil.move(str(src_path), str(temp_hop))
            shutil.move(str(temp_hop), str(dst_path))
            return True
        shutil.move(str(src_path), str(dst_path))
        return True
    except Exception:
        return False

def parse_except_patterns(except_str: str | None) -> set[str]:
    """Parses comma-separated patterns to ignore."""
    if not except_str:
        return set()
    return {p.strip().lower() for p in except_str.split(",") if p.strip()}

# --- Core Feature 1: Lowercase Renamer ---

def lowercase_directory(target_dir: str = CURRENT_DIR, except_patterns: str | None = None) -> int:
    """Recursively renames files and directories to lowercase."""
    custom_ignores = parse_except_patterns(except_patterns)
    re_upper = get_compiled_regex(RegexPatternType.UPPERCASE)
    renamed_count = 0
    start_time = time.perf_counter()

    for root, dirs, files in os.walk(target_dir, topdown=False):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d, custom_excludes=custom_ignores)]

        for f in files:
            if re_upper.search(f):
                src = Path(root) / f
                if is_ignored_path(src, custom_excludes=custom_ignores):
                    continue
                dst = Path(root) / f.lower()
                if run_git_mv_or_rename(src, dst):
                    print(f"  ✓ Renamed file: {normalize_rel_path(src)} -> {dst.name}")
                    renamed_count += 1

        for d in dirs:
            if re_upper.search(d):
                src = Path(root) / d
                if is_ignored_path(src, custom_excludes=custom_ignores):
                    continue
                dst = Path(root) / d.lower()
                if run_git_mv_or_rename(src, dst):
                    print(f"  ✓ Renamed dir:  {normalize_rel_path(src)} -> {dst.name}")
                    renamed_count += 1

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"{LINE_SEPARATOR}✅ Lowercase check complete: renamed {renamed_count} item(s) in {elapsed_ms:.2f}ms.")
    return ExitCodeType.SUCCESS.value

# --- Core Feature 2: Fix File Sequencing ---

def parse_pin_mappings(pin_str: str | None) -> dict[str, int]:
    """Parses pin mappings like 'readme=00,intro=01' into lookup dict."""
    if not pin_str:
        return {}
    pins = {}
    for part in pin_str.split(","):
        if "=" in part:
            name, num_s = part.split("=", 1)
            try:
                pins[name.strip().lower()] = int(num_s.strip())
            except ValueError:
                pass
    return pins

def fix_sequences_in_folder(
    folder_path: Path,
    is_order_by_time: bool = False,
    is_order_by_az: bool = False,
    is_keep_old_order: bool = True,
    pin_map: dict[str, int] | None = None
) -> int:
    """Re-sequences files inside a single directory sequentially."""
    pin_map = pin_map or {}
    re_seq = get_compiled_regex(RegexPatternType.SEQ_PREFIX)
    files = [f for f in folder_path.iterdir() if f.is_file() and not is_binary_file(f)]
    if not files:
        return 0

    def sort_key(f: Path):
        stem = f.stem.lower()
        for pin_name, pin_seq in pin_map.items():
            if pin_name in stem:
                return (0, pin_seq, stem)
        if is_order_by_time:
            return (1, f.stat().st_mtime, stem)
        if is_order_by_az:
            m = re_seq.match(f.name)
            base = m.group(2) if m else f.name
            return (1, 0, base.lower())
        m = re_seq.match(f.name)
        if m:
            if is_keep_old_order:
                return (1, int(m.group(1)), m.group(2).lower())
        return (1, 999, f.name.lower())

    sorted_files = sorted(files, key=sort_key)
    renamed_count = 0
    seq_idx = 1

    for f in sorted_files:
        stem = f.stem.lower()
        is_pinned = any(pin_name in stem for pin_name in pin_map)
        target_seq = pin_map[next(p for p in pin_map if p in stem)] if is_pinned else seq_idx

        m = re_seq.match(f.name)
        base_name = m.group(2) if m else f.name
        new_name = f"{target_seq:02d}-{base_name}"

        if f.name != new_name:
            dst = f.with_name(new_name)
            if run_git_mv_or_rename(f, dst):
                print(f"  ✓ Re-sequenced: {f.name} -> {new_name}")
                renamed_count += 1

        if not is_pinned:
            seq_idx += 1

    return renamed_count

def fix_sequences_recursive(target_dir: str = CURRENT_DIR, **kwargs) -> int:
    """Recursively fixes file sequencing across directories."""
    total_renamed = 0
    start_time = time.perf_counter()
    p = Path(target_dir)

    if p.is_dir():
        for root, dirs, _ in os.walk(p):
            dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
            total_renamed += fix_sequences_in_folder(Path(root), **kwargs)

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"{LINE_SEPARATOR}✅ Sequencing complete: re-sequenced {total_renamed} file(s) in {elapsed_ms:.2f}ms.")
    return ExitCodeType.SUCCESS.value

# --- Core Feature 3: Fix Encoding & Line Endings ---

def fix_encoding_and_newlines(target_dir: str = CURRENT_DIR, extensions: tuple | set | None = None) -> int:
    """Aggressively normalizes file encodings to UTF-8 without BOM and strict UNIX LF."""
    exts = normalize_extensions(extensions) or DEFAULT_TEXT_EXTENSIONS
    fixed_count = 0
    start_time = time.perf_counter()

    def handler(file_path: Path):
        nonlocal fixed_count
        try:
            with open(file_path, "rb") as f:
                raw_bytes = f.read()
            has_bom = raw_bytes.startswith(b"\xef\xbb\xbf")
            has_crlf = b"\r\n" in raw_bytes
            if has_bom or has_crlf:
                text = raw_bytes.decode(UTF8_SIG_ENCODING, errors="replace")
                if write_file_lf(file_path, text, encoding=DEFAULT_ENCODING):
                    fixed_count += 1
                    return normalize_rel_path(file_path)
        except Exception:
            pass
        return None

    stats = process_repository_files(handler, root_dir=target_dir, extensions=exts)
    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"{LINE_SEPARATOR}✅ Encoding & LF normalization complete: fixed {fixed_count} file(s) in {elapsed_ms:.2f}ms.")
    return ExitCodeType.SUCCESS.value

# --- CLI Entry Point ---

def main():
    parser = argparse.ArgumentParser(
        description="Autonomous file manipulation, lowercasing, sequencing, and encoding tool.",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    subparsers = parser.add_subparsers(dest="command", help="Subcommand to execute")

    # Lowercase sub-command
    p_lower = subparsers.add_parser("lowercase", help="Convert files and folders to lowercase")
    p_lower.add_argument("path", nargs="?", default=CURRENT_DIR, help="Directory to process")
    p_lower.add_argument("--except", dest="except_patterns", help="Comma-separated patterns to ignore")

    # Fix sequences sub-command
    p_seq = subparsers.add_parser("fix-seq-files", help="Re-sequence numbered files")
    p_seq.add_argument("path", nargs="?", default=CURRENT_DIR, help="Directory to process")
    p_seq.add_argument("--order-by-time", action="store_true", help="Order by modification time")
    p_seq.add_argument("--order-by-az", action="store_true", help="Order alphabetically")
    p_seq.add_argument("--keep-old-order", action="store_true", default=True, help="Preserve existing order where possible")
    p_seq.add_argument("--pin", help="Pin mappings, e.g. 'readme=00,intro=01'")

    # Fix encoding sub-command
    p_enc = subparsers.add_parser("fix-encoding", help="Normalize UTF-8 encoding and UNIX LF line endings")
    p_enc.add_argument("path", nargs="?", default=CURRENT_DIR, help="Directory to process")
    p_enc.add_argument("--ext", help="Comma-separated file extensions to process (e.g. .md,.ts,.py)")

    args = parser.parse_args()

    if args.command == "lowercase":
        sys.exit(lowercase_directory(target_dir=args.path, except_patterns=args.except_patterns))
    elif args.command == "fix-seq-files":
        pin_map = parse_pin_mappings(args.pin)
        sys.exit(fix_sequences_recursive(
            target_dir=args.path,
            is_order_by_time=args.order_by_time,
            is_order_by_az=args.order_by_az,
            is_keep_old_order=args.keep_old_order,
            pin_map=pin_map
        ))
    elif args.command == "fix-encoding":
        sys.exit(fix_encoding_and_newlines(target_dir=args.path, extensions=args.ext))
    else:
        parser.print_help()
        sys.exit(0)

if __name__ == "__main__":
    main()
