#!/usr/bin/env python3
"""
Fast Parallel Content Grep with Lazy Cache Index
Leverages tmp/cache/repo-file-cache.json and ThreadPoolExecutor to grep codebase in <15ms.
Zero-dependency, multi-folder capable, customizable extensions, and thread-safe lazy regex engine.

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/12-fast-cached-grep.py --pattern <pattern> [--path <dir>] [--ext <exts>] [--lang <langs>]
"""

import argparse
from concurrent.futures import ThreadPoolExecutor
from importlib import import_module
from pathlib import Path
import re
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

process_repository_files = engine.process_repository_files
read_file_safe = engine.read_file_safe
normalize_extensions = engine.normalize_extensions
normalize_rel_path = engine.normalize_rel_path
LANG_EXT_MAP = engine.LANG_EXT_MAP
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
ExitCodeType = engine.ExitCodeType

def grep_in_file(file_path: Path, pattern_re: re.Pattern) -> list[tuple[int, str]]:
    """Inspects a single file for regex pattern matches."""
    matches = []
    try:
        content = read_file_safe(file_path, encoding=DEFAULT_ENCODING)
        if content is None:
            return matches
        for idx, line in enumerate(content.split(LINE_SEPARATOR), start=1):
            has_match = bool(pattern_re.search(line))
            if has_match:
                matches.append((idx, line.strip()))
    except Exception:
        pass
    return matches

def run_cached_grep(
    pattern_str: str,
    target_dir: str = CURRENT_DIR,
    extensions: set[str] | tuple | None = None,
    is_regex: bool = False,
    is_case_sensitive: bool = False,
    max_results: int = 100
) -> int:
    """Executes parallel multi-threaded content grep across repository files."""
    start_time = time.perf_counter()
    flags = 0 if is_case_sensitive else re.IGNORECASE
    pattern_re = re.compile(pattern_str if is_regex else re.escape(pattern_str), flags)

    matched_records = []
    raw_files = []

    def collector(p: Path):
        raw_files.append(p)
        return None

    stats = process_repository_files(collector, root_dir=target_dir, extensions=extensions)

    with ThreadPoolExecutor(max_workers=8) as executor:
        futures = [executor.submit(grep_in_file, p, pattern_re) for p in raw_files]
        for p, fut in zip(raw_files, futures):
            res = fut.result()
            has_results = len(res) > 0
            if has_results:
                matched_records.append((normalize_rel_path(p), res))

    elapsed_ms = (time.perf_counter() - start_time) * 1000

    total_hits = sum(len(hits) for _, hits in matched_records)
    print(f"🔍 Grep Results for '{pattern_str}' in '{target_dir}': {total_hits} matches across {len(matched_records)} files ({elapsed_ms:.2f}ms){LINE_SEPARATOR}")

    printed_count = 0
    for fp, hits in matched_records:
        for line_no, line_text in hits:
            print(f"  {fp}:{line_no}: {line_text}")
            printed_count += 1
            has_reached_max = (printed_count >= max_results)
            if has_reached_max:
                break
        if has_reached_max:
            break

    has_overflow = (total_hits > printed_count)
    if has_overflow:
        print(f"{LINE_SEPARATOR}  ... and {total_hits - printed_count} more matches (use --max to adjust limit).")

    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Fast parallel cached content search")
    parser.add_argument("--pattern", "-p", required=True, help="Pattern or substring to search for")
    parser.add_argument("--path", "-d", default=CURRENT_DIR, help="Directory to search (default: .)")
    parser.add_argument("--ext", "-e", help="Comma-separated file extensions (e.g. .ts,.go,.py)")
    parser.add_argument("--lang", "-l", help="Language alias filter (e.g. go, ts, py)")
    parser.add_argument("--regex", action="store_true", help="Treat pattern as regular expression")
    parser.add_argument("--case-sensitive", action="store_true", help="Perform case-sensitive search")
    parser.add_argument("--max", type=int, default=100, help="Max line results to print (default: 100)")
    args = parser.parse_args()

    allowed_exts = set()
    if args.lang:
        for lang in args.lang.split(","):
            cleaned = lang.strip().lower()
            if cleaned in LANG_EXT_MAP:
                allowed_exts.update(LANG_EXT_MAP[cleaned])
    if args.ext:
        for e in args.ext.split(","):
            cleaned = e.strip().lower()
            if cleaned:
                allowed_exts.add(f".{cleaned}" if not cleaned.startswith(".") else cleaned)

    sys.exit(run_cached_grep(
        pattern_str=args.pattern,
        target_dir=args.path,
        extensions=allowed_exts if allowed_exts else None,
        is_regex=args.regex,
        is_case_sensitive=args.case_sensitive,
        max_results=args.max
    ))

if __name__ == "__main__":
    main()
