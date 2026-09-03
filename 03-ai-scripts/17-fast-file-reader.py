#!/usr/bin/env python3
"""
AI Agent Fast File Reader & Directory Explorer
Sub-millisecond file reader, folder explorer, and pattern searcher designed specifically for AI agents.
Leverages tmp/cache/ paths for instant lookups (<1ms) with automatic fallback to live disk scanning.

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder_path> [--ext .md,.ts]
  python 03-ai-scripts/17-fast-file-reader.py --read-file <file_path> [--max-bytes 50000]
  python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<term>" [--path <dir>]
"""

import argparse
from importlib import import_module
import json
import os
from pathlib import Path
import re
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

load_repo_cache = engine.load_repo_cache
read_file_safe = engine.read_file_safe
normalize_rel_path = engine.normalize_rel_path
normalize_extensions = engine.normalize_extensions
is_ignored_directory = engine.is_ignored_directory
is_ignored_path = engine.is_ignored_path
is_binary_file = engine.is_binary_file
ExitCodeType = engine.ExitCodeType
DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
PATH_SEPARATOR = engine.PATH_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
EMPTY_STRING = engine.EMPTY_STRING
CACHE_KEY_FILES = engine.CACHE_KEY_FILES

def list_folder_contents(folder_path: str = CURRENT_DIR, extensions: tuple | set | None = None) -> None:
    """Lists files and child directories within a folder."""
    start_time = time.perf_counter()
    norm_folder = normalize_rel_path(folder_path).rstrip(PATH_SEPARATOR)
    ext_set = normalize_extensions(extensions)
    cache = load_repo_cache()
    files_in_folder = []
    subfolders = set()

    has_cache = bool(cache and CACHE_KEY_FILES in cache)
    if has_cache:
        for f_path in cache[CACHE_KEY_FILES]:
            norm_f = normalize_rel_path(f_path)
            prefix = norm_folder + PATH_SEPARATOR if norm_folder and norm_folder != CURRENT_DIR else EMPTY_STRING
            if prefix and not norm_f.startswith(prefix):
                continue
            rel_to_folder = norm_f[len(prefix):] if prefix else norm_f
            if PATH_SEPARATOR in rel_to_folder:
                subfolders.add(rel_to_folder.split(PATH_SEPARATOR)[0])
            else:
                p = Path(norm_f)
                if ext_set and p.suffix.lower() not in ext_set:
                    continue
                files_in_folder.append(norm_f)
    else:
        p_folder = Path(norm_folder if norm_folder else CURRENT_DIR)
        if p_folder.exists() and p_folder.is_dir():
            for item in p_folder.iterdir():
                if is_ignored_directory(item.name):
                    continue
                if item.is_dir():
                    subfolders.add(item.name)
                elif item.is_file():
                    if ext_set and item.suffix.lower() not in ext_set:
                        continue
                    files_in_folder.append(normalize_rel_path(item))

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"📁 Folder Listing for `{folder_path}` ({elapsed_ms:.2f}ms):")
    has_subfolders = bool(subfolders)
    if has_subfolders:
        print("  📂 Subdirectories:")
        for sf in sorted(subfolders):
            print(f"    • {sf}/")
    has_files = bool(files_in_folder)
    if has_files:
        print(f"  📄 Files ({len(files_in_folder)}):")
        for f in sorted(files_in_folder):
            print(f"    • {f}")
    if not has_subfolders and not has_files:
        print("  (Empty directory or no matching files)")

def read_single_file(file_path: str, max_bytes: int = 100000) -> None:
    """Safely reads and displays the contents of a target file."""
    start_time = time.perf_counter()
    p = Path(file_path)
    if not p.exists() or not p.is_file():
        print(f"❌ File not found: {file_path}")
        sys.exit(ExitCodeType.VIOLATIONS_FOUND.value)

    content = read_file_safe(p, max_bytes=max_bytes, encoding=DEFAULT_ENCODING)
    elapsed_ms = (time.perf_counter() - start_time) * 1000

    if content is None:
        print(f"⚠️ Unable to read file: {file_path}")
        sys.exit(ExitCodeType.VIOLATIONS_FOUND.value)

    lines = content.split(LINE_SEPARATOR)
    print(f"📄 [{file_path}] — {len(lines)} lines, {len(content)} bytes ({elapsed_ms:.2f}ms):")
    print("-" * 80)
    print(content)
    print("-" * 80)

def search_files_by_pattern(pattern: str, target_dir: str = CURRENT_DIR, is_regex: bool = False) -> None:
    """Fast regex or literal content pattern search across target folder."""
    start_time = time.perf_counter()
    cache = load_repo_cache()
    compiled_re = re.compile(pattern if is_regex else re.escape(pattern), re.IGNORECASE)
    norm_root = normalize_rel_path(target_dir).rstrip(PATH_SEPARATOR)

    matches = []
    file_list = cache.get(CACHE_KEY_FILES, []) if (cache and CACHE_KEY_FILES in cache) else []
    if not file_list:
        p_root = Path(norm_root if norm_root else CURRENT_DIR)
        if p_root.exists():
            for root, dirs, files in os.walk(p_root):
                dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
                for f in files:
                    file_list.append(normalize_rel_path(os.path.join(root, f)))

    for f_str in file_list:
        norm_f = normalize_rel_path(f_str)
        prefix = norm_root + PATH_SEPARATOR if norm_root and norm_root != CURRENT_DIR else EMPTY_STRING
        if prefix and not norm_f.startswith(prefix):
            continue
        p = Path(norm_f)
        if is_binary_file(p):
            continue
        try:
            content = read_file_safe(p, encoding=DEFAULT_ENCODING)
            if content and compiled_re.search(content):
                matches.append(norm_f)
        except Exception:
            pass

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"🔍 Pattern `{pattern}` matched in {len(matches)} file(s) ({elapsed_ms:.2f}ms):")
    for m in matches[:25]:
        print(f"  • {m}")
    if len(matches) > 25:
        print(f"  ... and {len(matches) - 25} more files.")

def main():
    parser = argparse.ArgumentParser(
        description="Fast File Reader & Directory Explorer for AI Agents",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument("--list-folder", "-l", help="List directory contents and subfolders")
    parser.add_argument("--read-file", "-r", help="Read target text file contents")
    parser.add_argument("--search-pattern", "-s", help="Search content pattern across files")
    parser.add_argument("--path", "-p", default=CURRENT_DIR, help="Root directory for search or listing")
    parser.add_argument("--ext", help="Comma-separated extensions filter (e.g. .md,.ts,.py)")
    parser.add_argument("--max-bytes", type=int, default=100000, help="Maximum bytes to read (default 100KB)")
    parser.add_argument("--regex", action="store_true", help="Treat search pattern as regular expression")
    args = parser.parse_args()

    if args.list_folder:
        list_folder_contents(args.list_folder, extensions=args.ext)
    elif args.read_file:
        read_single_file(args.read_file, max_bytes=args.max_bytes)
    elif args.search_pattern:
        search_files_by_pattern(args.search_pattern, target_dir=args.path, is_regex=args.regex)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
