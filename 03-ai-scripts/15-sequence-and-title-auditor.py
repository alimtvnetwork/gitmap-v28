#!/usr/bin/env python3
"""
Sequence, Numbering & Title Header Auditor
Audits numbered markdown files (e.g. 01-intro.md) to ensure no sequence gaps and
verifies that the primary # H1 title matches the file sequence prefix.
Multi-folder capable, customizable extensions, and thread-safe lazy regex engine.

Performance & Clean Architecture:
1. Flattened Conditionals: Zero deep nesting using early `continue` guard clauses.
2. All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

import argparse
from importlib import import_module
import os
from pathlib import Path
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf
normalize_rel_path = engine.normalize_rel_path
is_ignored_directory = engine.is_ignored_directory
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

def audit_directory_sequences(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = False
) -> tuple[list[str], list[str]]:
    """Audits sequence numbering and # H1 headers in markdown files using flattened guard clauses."""
    re_num = get_compiled_regex(RegexPatternType.FILE_NUM_PREFIX)
    re_h1 = get_compiled_regex(RegexPatternType.H1_HEADER)
    seq_issues = []
    title_issues = []

    for root, dirs, files in os.walk(target_dir):
        dirs[:] = [d for d in dirs if not is_ignored_directory(d)]
        norm_dir = normalize_rel_path(root)

        numbered_files = []
        for f in files:
            m = re_num.match(f)
            if not m:
                continue
            numbered_files.append((int(m.group(1)), f, os.path.join(root, f)))

        if not numbered_files:
            continue

        numbered_files.sort(key=lambda x: x[0])
        first_num = numbered_files[0][0]
        if first_num in (0, 1):
            expected = first_num
            for num_val, f_name, _ in numbered_files:
                if num_val != expected:
                    seq_issues.append(f"{norm_dir}/{f_name} (found {num_val:02d}, expected {expected:02d})")
                expected += 1

        for num_val, f_name, full_path in numbered_files:
            content = read_file_lf(full_path, encoding=DEFAULT_ENCODING)
            m_h1 = re_h1.search(content)
            if not m_h1:
                continue

            h1_num = int(m_h1.group(2))
            if h1_num == num_val:
                continue

            title_issues.append(f"{norm_dir}/{f_name} (file prefix {num_val:02d} != H1 header {h1_num:02d})")
            if is_fix_mode:
                def replacer(match):
                    return f"{match.group(1)}{num_val:02d}{match.group(3)}{match.group(4)}"
                new_content = re_h1.sub(replacer, content, count=1)
                write_file_lf(full_path, new_content, encoding=DEFAULT_ENCODING)

    return seq_issues, title_issues

def run_sequence_auditor(target_dir: str = CURRENT_DIR, is_fix_mode: bool = False) -> int:
    """Executes sequence and title audit across target directory."""
    seq_issues, title_issues = audit_directory_sequences(target_dir, is_fix_mode=is_fix_mode)

    has_errors = False
    has_seq = len(seq_issues) > 0
    if has_seq:
        has_errors = True
        print(f"{LINE_SEPARATOR}❌ Found {len(seq_issues)} sequence gap(s):")
        for issue in seq_issues:
            print(f"  ::error::{issue}")

    has_title = len(title_issues) > 0
    if has_title:
        if not is_fix_mode:
            has_errors = True
        action_word = "Fixed" if is_fix_mode else "Found"
        print(f"{LINE_SEPARATOR}⚠️ {action_word} {len(title_issues)} H1 title mismatch issue(s):")
        for issue in title_issues:
            print(f"  ::warning::{issue}")

    if has_errors:
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ All numbered files and H1 titles in '{target_dir}' are correctly sequenced.")
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Audit file sequence numbering and H1 title alignment")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory to audit")
    parser.add_argument("--dir", "--path", "-p", dest="opt_dir", help="Directory to audit")
    parser.add_argument("--fix", action="store_true", help="Auto-fix H1 title numbers in markdown files")
    args = parser.parse_args()

    target_path = args.opt_dir or args.path or CURRENT_DIR
    sys.exit(run_sequence_auditor(target_dir=target_path, is_fix_mode=args.fix))

if __name__ == "__main__":
    main()
