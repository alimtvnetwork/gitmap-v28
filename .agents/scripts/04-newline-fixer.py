#!/usr/bin/env python3
"""
Fast Newline & Trailing Whitespace Fixer
Enforces clean UNIX LF line endings, trims trailing spaces, and ensures a single trailing newline.
Multi-folder capable, customizable extensions, and thread-safe lazy regex engine.

Performance & Clean Architecture:
1. Flattened Conditionals: Zero nested if-blocks using clean guard clauses.
2. All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
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
DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

def clean_file_content(content: str) -> str:
    """Strips trailing whitespace per line and guarantees a single final newline."""
    re_crlf = get_compiled_regex(RegexPatternType.CRLF)
    normalized = re_crlf.sub(LINE_SEPARATOR, content)
    lines = [line.rstrip() for line in normalized.split(LINE_SEPARATOR)]
    while lines:
        if lines[-1]:
            break
        lines.pop()
    return LINE_SEPARATOR.join(lines) + LINE_SEPARATOR

def process_file_newlines(file_path: Path, is_fix_mode: bool = False) -> tuple[str, bool]:
    """Checks and optionally fixes newlines and trailing whitespace in a single file."""
    norm_p = normalize_rel_path(file_path)
    try:
        raw = read_file_lf(file_path, encoding=DEFAULT_ENCODING)
        if not raw:
            return (norm_p, False)
        cleaned = clean_file_content(raw)
        has_changes = (raw != cleaned)
        if not has_changes:
            return (norm_p, False)
        if is_fix_mode:
            write_file_lf(file_path, cleaned, encoding=DEFAULT_ENCODING)
        return (norm_p, True)
    except Exception:
        pass
    return (norm_p, False)

def run_newline_auditor(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = False,
    extensions: set[str] | tuple | None = None
) -> int:
    """Executes two-phase repository scan to audit/fix newlines across any target directory."""
    exts = normalize_extensions(extensions) or DEFAULT_TEXT_EXTENSIONS

    def handler(p: Path):
        path_str, has_issue = process_file_newlines(p, is_fix_mode=is_fix_mode)
        return path_str if has_issue else None

    stats = process_repository_files(handler, root_dir=target_dir, extensions=exts)
    violations = stats["results"]

    if violations:
        action_word = "Fixed" if is_fix_mode else "Found issues in"
        print(f"{LINE_SEPARATOR}⚠️ {action_word} {len(violations)} file(s) ({stats['elapsed_ms']:.2f}ms):")
        for v in violations[:10]:
            print(f"  ::notice file={v}::{v}")
        if not is_fix_mode:
            return ExitCodeType.VIOLATIONS_FOUND.value
    else:
        print(f"✅ All {stats['total_files']} files in '{target_dir}' have clean newlines ({stats['elapsed_ms']:.2f}ms).")

    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Fix trailing whitespace and newlines across folders")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory or subfolder to scan")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument("--fix", action="store_true", help="Auto-fix whitespace issues in-place")
    parser.add_argument("--ext", help="Comma-separated extensions to scan (e.g. .md,.ts,.py)")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    sys.exit(run_newline_auditor(target_dir=target_path, is_fix_mode=args.fix, extensions=args.ext))

if __name__ == "__main__":
    main()
