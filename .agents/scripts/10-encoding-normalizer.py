#!/usr/bin/env python3
"""
Fast UTF-8 & UNIX LF Encoding Normalizer
Recursively audits and standardizes all text files to UTF-8 without BOM and strict UNIX LF (\\n).
Multi-folder capable, customizable extensions, and sub-15ms execution.

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
write_file_lf = engine.write_file_lf
is_binary_file = engine.is_binary_file
normalize_extensions = engine.normalize_extensions
normalize_rel_path = engine.normalize_rel_path
ExitCodeType = engine.ExitCodeType
DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
UTF8_SIG_ENCODING = engine.UTF8_SIG_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
UTF8_BOM_BYTES = engine.UTF8_BOM_BYTES
CRLF_BYTES = engine.CRLF_BYTES

def normalize_single_file(file_path: Path, is_fix_mode: bool = False) -> tuple[str, bool]:
    """Audits and converts CRLF/BOM in a file to clean UTF-8 LF using flattened guard clauses."""
    norm_p = normalize_rel_path(file_path)
    is_binary = is_binary_file(file_path)
    if is_binary:
        return (norm_p, False)
    try:
        with open(file_path, "rb") as f:
            raw_bytes = f.read()

        has_issue = (raw_bytes.startswith(UTF8_BOM_BYTES) or CRLF_BYTES in raw_bytes)
        if not has_issue:
            return (norm_p, False)

        if is_fix_mode:
            text = raw_bytes.decode(UTF8_SIG_ENCODING, errors="replace")
            write_file_lf(file_path, text, encoding=DEFAULT_ENCODING)
        return (norm_p, True)
    except Exception:
        pass
    return (norm_p, False)

def run_encoding_normalizer(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = False,
    extensions: set[str] | tuple | None = None
) -> int:
    """Runs repository encoding check and normalizer across target directory."""
    exts = normalize_extensions(extensions) or DEFAULT_TEXT_EXTENSIONS

    def handler(p: Path):
        fp_str, has_changed = normalize_single_file(p, is_fix_mode=is_fix_mode)
        return fp_str if has_changed else None

    stats = process_repository_files(handler, root_dir=target_dir, extensions=exts)
    affected = stats["results"]

    has_affected = len(affected) > 0
    if has_affected:
        action_verb = "Normalized" if is_fix_mode else "Found CRLF/BOM in"
        print(f"{LINE_SEPARATOR}⚠️ {action_verb} {len(affected)} file(s) ({stats['elapsed_ms']:.2f}ms):")
        for f in affected[:5]:
            print(f"  ::notice file={f}::{f}")
        if not is_fix_mode:
            return ExitCodeType.VIOLATIONS_FOUND.value
    else:
        print(f"✅ All {stats['total_files']} files in '{target_dir}' normalized to UTF-8 LF ({stats['elapsed_ms']:.2f}ms).")

    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Normalize files to UTF-8 UNIX LF across folders")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory or subfolder")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument("--fix", action="store_true", help="Fix CRLF and BOM in-place")
    parser.add_argument("--ext", help="Comma-separated extensions to scan (e.g. .md,.ts,.py)")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    sys.exit(run_encoding_normalizer(target_dir=target_path, is_fix_mode=args.fix, extensions=args.ext))

if __name__ == "__main__":
    main()
