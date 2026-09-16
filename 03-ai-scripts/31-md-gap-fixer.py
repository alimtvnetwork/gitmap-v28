#!/usr/bin/env python3
"""
Fixes multiple consecutive empty lines (3 or more newlines) into exactly 2 newlines (1 empty line).
"""

import argparse
from importlib import import_module
from pathlib import Path
import sys
import re

sys.path.insert(0, str(Path(__file__).parent))
try:
    engine = import_module("02-shared-engine")
    process_repository_files = engine.process_repository_files
    read_file_lf = engine.read_file_lf
    write_file_lf = engine.write_file_lf
    normalize_extensions = engine.normalize_extensions
    normalize_rel_path = engine.normalize_rel_path
    ExitCodeType = engine.ExitCodeType
    LINE_SEPARATOR = engine.LINE_SEPARATOR
    CURRENT_DIR = engine.CURRENT_DIR
    DEFAULT_ENCODING = engine.DEFAULT_ENCODING
except ImportError:
    # Fallback if engine is missing
    print("Could not load 02-shared-engine.py, running standalone mode.", file=sys.stderr)
    sys.exit(1)

def fix_gaps(content: str) -> str:
    """Collapses 3 or more consecutive newlines into exactly 2 newlines."""
    # First, normalize CRLF to LF
    normalized = content.replace('\r\n', '\n')
    # Replace 3 or more newlines with exactly 2
    fixed = re.sub(r'\n{3,}', '\n\n', normalized)
    return fixed

def process_file(file_path: Path, is_fix_mode: bool) -> tuple[str, bool]:
    norm_p = normalize_rel_path(file_path)
    try:
        raw = read_file_lf(file_path, encoding=DEFAULT_ENCODING)
        if not raw:
            return (norm_p, False)

        cleaned = fix_gaps(raw)
        has_changes = (raw != cleaned)

        if not has_changes:
            return (norm_p, False)

        if is_fix_mode:
            write_file_lf(file_path, cleaned, encoding=DEFAULT_ENCODING)
        return (norm_p, True)
    except Exception as e:
        print(f"Error processing {norm_p}: {e}")
    return (norm_p, False)

def main():
    parser = argparse.ArgumentParser(description="Fix multi-line gaps in markdown files")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Target directory")
    parser.add_argument("--fix", action="store_true", help="Apply fixes in-place")
    args = parser.parse_args()

    target_dir = args.path
    is_fix_mode = args.fix
    exts = normalize_extensions(".md")

    def handler(p: Path):
        path_str, has_issue = process_file(p, is_fix_mode=is_fix_mode)
        return path_str if has_issue else None

    stats = process_repository_files(handler, root_dir=target_dir, extensions=exts)
    violations = stats["results"]

    if violations:
        action_word = "Fixed" if is_fix_mode else "Found gaps in"
        print(f"\n{action_word} {len(violations)} file(s) ({stats['elapsed_ms']:.2f}ms):")
        for v in violations[:15]:
            print(f"  - {v}")
        if len(violations) > 15:
            print(f"  ... and {len(violations) - 15} more.")

        if not is_fix_mode:
            sys.exit(ExitCodeType.VIOLATIONS_FOUND.value)
    else:
        print(f"All {stats['total_files']} files in '{target_dir}' have clean gaps ({stats['elapsed_ms']:.2f}ms).")

    sys.exit(ExitCodeType.SUCCESS.value)

if __name__ == "__main__":
    main()
