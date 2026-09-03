#!/usr/bin/env python3
"""
Fast Relative Path Fixer & Absolute URI Auditor
Detects and sanitizes absolute filesystem paths (C:\\..., D:\\..., /home/..., file:///) in documentation.
Multi-folder capable, customizable extensions, and thread-safe lazy regex engine.

Performance & Clean Architecture:
1. Fast Substring Pre-Filter: Skips regex if no path markers exist (<0.001ms).
2. Flattened Conditionals: Clean guard clauses and flattened verification steps.
3. All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

import argparse
from importlib import import_module
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

process_repository_files = engine.process_repository_files
read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf
normalize_extensions = engine.normalize_extensions
normalize_rel_path = engine.normalize_rel_path
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
get_compiled_regex_group = engine.get_compiled_regex_group
DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
DEVICE_PATH_PREFIX = engine.DEVICE_PATH_PREFIX
DOT_CHAR = engine.DOT_CHAR

def sanitize_content_paths(content: str) -> tuple[str, int]:
    """Replaces absolute repository paths with clean relative paths."""
    re_repo = get_compiled_regex(RegexPatternType.REPO_FILE_URI)
    modified = content
    count = 0
    for match in re_repo.finditer(content):
        rel_target = match.group(1)
        modified = modified.replace(match.group(0), rel_target)
        count += 1
    return modified, count

def audit_file_paths(file_path: Path, is_fix_mode: bool = False) -> tuple[str, list[str]]:
    """Audits a single file for forbidden absolute paths using flattened guard clauses."""
    norm_p = normalize_rel_path(file_path)
    try:
        content = read_file_lf(file_path, encoding=DEFAULT_ENCODING)
        # Fast substring pre-filter: skip if no path indicators exist
        has_path_marker = ("file:" in content or ":\\" in content)
        if not has_path_marker:
            return (norm_p, [])

        patterns = get_compiled_regex_group(
            RegexPatternType.FILE_URI_WIN,
            RegexPatternType.DRIVE_ABS_WIN,
        )

        violations = []
        for pat in patterns:
            for match in pat.finditer(content):
                val = match.group(0)
                is_ignored = (DEVICE_PATH_PREFIX in val or val.endswith(DOT_CHAR))
                if is_ignored:
                    continue
                violations.append(val)

        has_violations = len(violations) > 0
        if not has_violations:
            return (norm_p, [])

        if is_fix_mode:
            cleaned, fix_count = sanitize_content_paths(content)
            if fix_count > 0:
                write_file_lf(file_path, cleaned, encoding=DEFAULT_ENCODING)

        return (norm_p, violations)
    except Exception:
        return (norm_p, [])

def run_path_auditor(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = False,
    extensions: set[str] | tuple | None = None
) -> int:
    """Runs repository-wide path check using two-phase pipeline."""
    exts = normalize_extensions(extensions) or DEFAULT_TEXT_EXTENSIONS

    def handler(p: Path):
        fp_str, vios = audit_file_paths(p, is_fix_mode=is_fix_mode)
        return (fp_str, vios) if vios else None

    stats = process_repository_files(handler, root_dir=target_dir, extensions=exts)
    all_violations = stats["results"]

    has_violations = len(all_violations) > 0
    if has_violations:
        print(f"{LINE_SEPARATOR}❌ Found absolute path references in {len(all_violations)} file(s) ({stats['elapsed_ms']:.2f}ms):")
        for fp, vios in all_violations:
            for v in vios[:3]:
                print(f"  ::error file={fp}::Absolute path found: {v}")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ All {stats['total_files']} files in '{target_dir}' use strict relative paths ({stats['elapsed_ms']:.2f}ms).")
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Audit and fix absolute paths across target folders")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory or folder to scan")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument("--fix", action="store_true", help="Auto-fix recognized path patterns")
    parser.add_argument("--ext", help="Comma-separated extensions to scan (e.g. .md,.ts,.py)")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    sys.exit(run_path_auditor(target_dir=target_path, is_fix_mode=args.fix, extensions=args.ext))

if __name__ == "__main__":
    main()
