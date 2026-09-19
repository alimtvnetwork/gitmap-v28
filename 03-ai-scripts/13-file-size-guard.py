#!/usr/bin/env python3
"""
Fast Repository File Size & Blob Guard
Scans tracked repository files to ensure no accidental massive binary files exceed thresholds.
Multi-folder capable, customizable extensions, and nested ignore pruning (.git, .gitmap, node_modules).

Performance & Clean Architecture:
1. Flattened Conditionals: Zero nested if-blocks using clean guard clauses.
2. All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

import argparse
from importlib import import_module
import os
from pathlib import Path
import sys
import time

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

DEFAULT_MAX_FILE_KB = engine.DEFAULT_MAX_FILE_KB
EXCLUDE_DIRS = engine.EXCLUDE_DIRS
is_allowed_large_file = engine.is_allowed_large_file
is_ignored_directory = engine.is_ignored_directory
normalize_rel_path = engine.normalize_rel_path
normalize_extensions = engine.normalize_extensions
ExitCodeType = engine.ExitCodeType
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

def audit_file_sizes(
    max_kb: int = DEFAULT_MAX_FILE_KB,
    target_dir: str = CURRENT_DIR,
    allowed_exts: set[str] | None = None
) -> int:
    """Scans files and checks sizes against threshold across target directory using flattened guard clauses."""
    start_time = time.perf_counter()
    violations = []
    total_files = 0
    total_bytes = 0
    max_bytes = max_kb * 1024

    for root, dirs, files in os.walk(target_dir):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
        for f in files:
            fp = os.path.join(root, f)
            norm_fp = normalize_rel_path(fp)
            is_allowed = is_allowed_large_file(norm_fp)
            if is_allowed:
                continue

            is_filtered_ext = (allowed_exts is not None and os.path.splitext(f)[1].lower() not in allowed_exts)
            if is_filtered_ext:
                continue

            try:
                sz = os.path.getsize(fp)
                total_files += 1
                total_bytes += sz
                is_oversized = (sz > max_bytes)
                if is_oversized:
                    violations.append((norm_fp, sz))
            except Exception:
                pass

    elapsed_ms = (time.perf_counter() - start_time) * 1000
    print(f"📊 Scanned {total_files:,} files in '{target_dir}' ({total_bytes / (1024*1024):.2f} MB) in {elapsed_ms:.2f}ms")

    has_violations = len(violations) > 0
    if has_violations:
        print(f"{LINE_SEPARATOR}❌ Found {len(violations)} oversized file(s) exceeding {max_kb} KB:")
        for fp, sz in sorted(violations, key=lambda x: x[1], reverse=True):
            print(f"  ::error file={fp}::{fp} ({sz / 1024:.1f} KB > {max_kb} KB)")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ All files in '{target_dir}' within {max_kb} KB limit.")
    return ExitCodeType.SUCCESS.value

def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    parser = argparse.ArgumentParser(description="Audit repository file sizes across folders")
    parser.add_argument("--max-kb", type=int, default=DEFAULT_MAX_FILE_KB, help="Maximum allowed file size in KB")
    parser.add_argument("--path", "-p", default=CURRENT_DIR, help="Root path or folder to audit")
    parser.add_argument("--ext", help="Optional comma-separated extension filter (e.g. .json,.zip)")
    args = parser.parse_args()

    exts = normalize_extensions(args.ext)
    sys.exit(audit_file_sizes(max_kb=args.max_kb, target_dir=args.path, allowed_exts=exts))

if __name__ == "__main__":
    main()
