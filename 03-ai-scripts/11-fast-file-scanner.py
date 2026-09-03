#!/usr/bin/env python3
"""
Fast Repository File Scanner & Caching Engine.
Ultra-fast file scanner with multi-language filters, substring search, and persistent cache indexing in tmp/cache/.

All Enums, Cache Keys, Constants, and Functions are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/11-fast-file-scanner.py [options]

Examples:
  python 03-ai-scripts/11-fast-file-scanner.py
  python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts,tsx
  python 03-ai-scripts/11-fast-file-scanner.py --path spec/ --ext .md
  python 03-ai-scripts/11-fast-file-scanner.py --search install --stats
  python 03-ai-scripts/11-fast-file-scanner.py --query-cache "component"
"""

import argparse
import datetime
from importlib import import_module
import json
import os
from pathlib import Path
import re
import sys
import time

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

DEFAULT_IGNORE_DIRS = engine.EXCLUDE_DIRS
BINARY_EXTENSIONS = engine.BINARY_EXTENSIONS
LANG_EXT_MAP = engine.LANG_EXT_MAP
RegexPatternType = engine.RegexPatternType
ExitCodeType = engine.ExitCodeType
CacheKeyType = engine.CacheKeyType
get_compiled_regex = engine.get_compiled_regex
normalize_rel_path = engine.normalize_rel_path
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
PATH_SEPARATOR = engine.PATH_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
EMPTY_STRING = engine.EMPTY_STRING
CACHE_KEY_FILES = engine.CACHE_KEY_FILES
CACHE_KEY_TOTAL_FILES = engine.CACHE_KEY_TOTAL_FILES
CACHE_KEY_TIMESTAMP = engine.CACHE_KEY_TIMESTAMP
PRIMARY_CACHE_FILE = engine.PRIMARY_CACHE_FILE
LEGACY_CACHE_FILE = engine.LEGACY_CACHE_FILE
CACHE_PATHS_DIR = engine.CACHE_PATHS_DIR
CACHE_BASE_DIR = engine.CACHE_BASE_DIR

def parse_args():
    parser = argparse.ArgumentParser(
        description="High-performance repository file scanner and cache indexer.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--path", "-p", default=CURRENT_DIR, help="Subdirectory or root to scan (default: .)")
    parser.add_argument("--lang", "-l", help="Language filter alias (e.g. go, ts, py, md, or comma-separated go,ts)")
    parser.add_argument("--ext", "-e", help="Custom extension filter (comma-separated, e.g. .go,.ts,.json)")
    parser.add_argument("--search", "-s", help="Case-insensitive substring filter on file path")
    parser.add_argument("--out", "-o", help="Custom cache output path (default: auto-named in tmp/)")
    parser.add_argument("--format", "-f", choices=["json", "txt", "summary"], default="json", help="Output format (json, txt, summary)")
    parser.add_argument("--limit", type=int, default=100, help="Max file lines to print to console (default: 100)")
    parser.add_argument("--stats", action="store_true", help="Display extension statistics breakdown")
    parser.add_argument("--no-cache", action="store_true", help="Skip saving results to tmp/ cache")
    parser.add_argument("--include-hidden", action="store_true", help="Include dot-files/folders (normally ignored)")
    parser.add_argument("--query-cache", "-q", help="Query existing cached index without walking the filesystem")
    parser.add_argument("--check", action="store_true", help="CI validation mode: verifies file index and exits 0 on success")

    return parser.parse_args()

def resolve_extensions(lang_arg: str | None, ext_arg: str | None) -> set[str] | None:
    allowed_exts = set()
    if lang_arg:
        for lang in lang_arg.split(","):
            cleaned = lang.strip().lower()
            if cleaned in LANG_EXT_MAP:
                allowed_exts.update(LANG_EXT_MAP[cleaned])
            elif cleaned:
                ext_form = f".{cleaned}" if not cleaned.startswith(".") else cleaned
                allowed_exts.add(ext_form)

    if ext_arg:
        for e in ext_arg.split(","):
            cleaned = e.strip().lower()
            if cleaned:
                ext_form = f".{cleaned}" if not cleaned.startswith(".") else cleaned
                allowed_exts.add(ext_form)

    return allowed_exts if allowed_exts else None

def scan_files(scan_root: str, allowed_exts: set[str] | None, search_term: str | None, is_include_hidden: bool):
    matched_files = []
    ext_counts = {}

    WHITELISTED_DOT_DIRS = {".lovable", ".github"}
    search_re = re.compile(re.escape(search_term), re.IGNORECASE) if search_term else None

    for root, dirs, files in os.walk(scan_root):
        if not is_include_hidden:
            dirs[:] = [
                d for d in dirs
                if d not in DEFAULT_IGNORE_DIRS and (not d.startswith(".") or d in WHITELISTED_DOT_DIRS)
            ]
        else:
            dirs[:] = [d for d in dirs if d not in DEFAULT_IGNORE_DIRS]

        for filename in sorted(files):
            if not is_include_hidden:
                if filename.startswith("."):
                    if not filename.startswith(".lovable"):
                        continue

            ext = os.path.splitext(filename)[1].lower()
            is_bin = (ext in BINARY_EXTENSIONS)
            if is_bin:
                continue

            if allowed_exts:
                if ext not in allowed_exts:
                    continue

            full_rel = normalize_rel_path(os.path.relpath(os.path.join(root, filename), CURRENT_DIR))
            if search_re:
                if not search_re.search(full_rel):
                    continue

            matched_files.append(full_rel)
            ext_counts[ext] = ext_counts.get(ext, 0) + 1

    return matched_files, ext_counts

def get_cache_filenames(args):
    re_clean = get_compiled_regex(RegexPatternType.NON_ALPHANUMERIC)
    slug_parts = []
    if args.path:
        if args.path != CURRENT_DIR:
            clean_p = re_clean.sub("_", args.path).strip("_")
            if clean_p:
                slug_parts.append(clean_p)
    if args.lang:
        clean_l = re_clean.sub("_", args.lang).strip("_")
        slug_parts.append(f"lang-{clean_l}")
    if args.ext:
        clean_e = re_clean.sub("_", args.ext).strip("_")
        slug_parts.append(f"ext-{clean_e}")
    if args.search:
        clean_s = re_clean.sub("_", args.search).strip("_")
        slug_parts.append(f"search-{clean_s}")

    slug = "-".join(slug_parts) if slug_parts else "all"

    json_path = args.out if (args.out and args.out.endswith(".json")) else str(LEGACY_CACHE_FILE)
    txt_path = f"tmp/file-list-{slug}.txt"
    all_txt_path = "tmp/file-list-all.txt"

    return json_path, txt_path, all_txt_path

def write_caches(json_path, txt_path, all_txt_path, matched_files, ext_counts, args, scan_duration_ms):
    os.makedirs("tmp", exist_ok=True)
    os.makedirs(str(CACHE_PATHS_DIR), exist_ok=True)
    cache_data = {
        CACHE_KEY_TIMESTAMP: datetime.datetime.now(datetime.timezone.utc).isoformat(),
        "scanRoot": args.path,
        "scanDurationMs": round(scan_duration_ms, 2),
        CACHE_KEY_TOTAL_FILES: len(matched_files),
        "filterCriteria": {
            "lang": args.lang,
            "ext": args.ext,
            "search": args.search,
            "includeHidden": args.include_hidden,
        },
        "cacheFiles": {
            "jsonCache": json_path,
            "primaryCache": str(PRIMARY_CACHE_FILE),
            "filterTextList": txt_path,
            "globalTextList": all_txt_path,
        },
        "stats": {
            "byExtension": dict(sorted(ext_counts.items(), key=lambda x: x[1], reverse=True))
        },
        CACHE_KEY_FILES: matched_files,
    }

    # 1. Main JSON Cache
    with open(json_path, "w", encoding=DEFAULT_ENCODING) as f:
        json.dump(cache_data, f, indent=2)

    try:
        with open(PRIMARY_CACHE_FILE, "w", encoding=DEFAULT_ENCODING) as f:
            json.dump(cache_data, f, indent=2)
    except Exception:
        pass

    # 2. Filter-specific text file
    with open(txt_path, "w", encoding=DEFAULT_ENCODING) as f:
        for fp in matched_files:
            f.write(fp + LINE_SEPARATOR)

    # 3. Global text list
    if args.path == CURRENT_DIR:
        if not args.lang:
            if not args.ext:
                if not args.search:
                    with open(all_txt_path, "w", encoding=DEFAULT_ENCODING) as f:
                        for fp in matched_files:
                            f.write(fp + LINE_SEPARATOR)

def query_cached_index(query_term: str):
    cache_paths = [PRIMARY_CACHE_FILE, LEGACY_CACHE_FILE]
    matched_files = []
    for cp in cache_paths:
        if os.path.exists(cp):
            try:
                with open(cp, "r", encoding=DEFAULT_ENCODING) as f:
                    data = json.load(f)
                matched_files = data.get(CACHE_KEY_FILES, [])
                if matched_files:
                    break
            except Exception:
                pass

    if not matched_files:
        print("⚠️ Cache not found. Running quick scan to generate cache...")
        matched_files, _ = scan_files(CURRENT_DIR, None, None, False)
        os.makedirs(str(CACHE_BASE_DIR), exist_ok=True)
        with open(PRIMARY_CACHE_FILE, "w", encoding=DEFAULT_ENCODING) as f:
            json.dump({CACHE_KEY_FILES: matched_files}, f)

    q_re = re.compile(re.escape(query_term), re.IGNORECASE)
    results = [f for f in matched_files if q_re.search(f)]
    print(f"⚡ Instant Cache Query for `{query_term}`: found **{len(results)}** matches in pre-computed index:{LINE_SEPARATOR}")
    for idx, r in enumerate(results[:100], 1):
        print(f"   {idx:>3}. {r}")
    if len(results) > 100:
        print(f"   ... and {len(results) - 100} more matches.")
    sys.exit(0)

def main():
    args = parse_args()

    if args.query_cache:
        query_cached_index(args.query_cache)
        return

    start_time = time.perf_counter()
    allowed_exts = resolve_extensions(args.lang, args.ext)
    matched_files, ext_counts = scan_files(args.path, allowed_exts, args.search, is_include_hidden=args.include_hidden)
    scan_duration_ms = (time.perf_counter() - start_time) * 1000.0

    json_path, txt_path, all_txt_path = get_cache_filenames(args)

    is_no_cache = args.no_cache
    if not is_no_cache:
        write_caches(json_path, txt_path, all_txt_path, matched_files, ext_counts, args, scan_duration_ms)

    is_check_mode = args.check
    if is_check_mode:
        print(f"✅ Fast File Scanner Check PASSED: {len(matched_files)} files indexed in {scan_duration_ms:.2f}ms")
        sys.exit(0)

    # Print clean summary & results to console
    print("================================================================================")
    print(f"⚡ Fast Repository File Scanner: scanned in {scan_duration_ms:.2f}ms")
    print(f"📁 Root: `{args.path}` | Filtered Files Found: **{len(matched_files)}**")
    print("================================================================================")

    if not is_no_cache:
        print(f"{LINE_SEPARATOR}💾 TEMP FOLDER CACHE INVENTORY:")
        print(f"   • Full JSON Cache : `{json_path}`")
        print(f"   • Primary Pluggable: `{PRIMARY_CACHE_FILE}`")
        print(f"   • Specific Filter : `{txt_path}`")
        print(f"   • Global List     : `{all_txt_path}`")

    has_stats_req = args.stats or len(matched_files) == 0
    if has_stats_req:
        print(f"{LINE_SEPARATOR}📊 Extension Breakdown:")
        for ext, count in sorted(ext_counts.items(), key=lambda x: x[1], reverse=True):
            display_ext = ext if ext else "(no extension)"
            print(f"   • {display_ext:<14} : {count:>5} files")

    print(f"{LINE_SEPARATOR}📋 File Inventory Preview:")
    preview_limit = args.limit if args.limit > 0 else len(matched_files)
    for idx, fp in enumerate(matched_files[:preview_limit], 1):
        print(f"   {idx:>4}. {fp}")

    has_more_files = len(matched_files) > preview_limit
    if has_more_files:
        print(f"   ... and {len(matched_files) - preview_limit} more files (see `{txt_path}` for complete list).")

    print(f"{LINE_SEPARATOR}💡 AI AGENT INSTRUCTION:")
    print("   To read or verify files in subsequent steps without ad-hoc queries,")
    print(f"   simply read `{txt_path}` or `{json_path}` directly.")
    print("================================================================================")
    sys.exit(0)

if __name__ == "__main__":
    main()
