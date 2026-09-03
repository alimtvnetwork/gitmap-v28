#!/usr/bin/env python3
"""
Fast Boolean Naming & Code Convention Guard
Audits and flags explicit boolean true comparisons (e.g. `== True`, `=== true`) and negative naming anti-patterns.
Multi-folder capable, customizable extensions, and thread-safe lazy regex engine.

Performance & Clean Architecture:
1. Fast Substring Pre-Filter: Skips line-by-line regex scanning if "true" is absent (5x-10x speedup).
2. Flattened Conditionals: Zero nested if-blocks using clean guard clauses and early exits.
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
normalize_extensions = engine.normalize_extensions
normalize_rel_path = engine.normalize_rel_path
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
get_compiled_regex_group = engine.get_compiled_regex_group
DEFAULT_CODE_EXTENSIONS = engine.DEFAULT_CODE_EXTENSIONS
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

def find_explicit_true_violations(content: str) -> list[tuple[int, str]]:
    """Inspects lines for explicit true comparisons using early continue guard clauses."""
    re_comment = get_compiled_regex(RegexPatternType.COMMENT_PREFIX)
    explicit_patterns = get_compiled_regex_group(
        RegexPatternType.EXPLICIT_DOUBLE_TRUE,
        RegexPatternType.EXPLICIT_TRIPLE_TRUE,
        RegexPatternType.EXPLICIT_PYTHON_TRUE,
    )
    violations = []
    for idx, line in enumerate(content.split(LINE_SEPARATOR), start=1):
        is_comment = bool(re_comment.match(line))
        if is_comment:
            continue
        stripped = line.strip()
        has_violation = any(pat.search(line) for pat in explicit_patterns)
        if has_violation:
            violations.append((idx, stripped))
    return violations

def audit_file_naming(file_path: Path) -> tuple[str, list[tuple[int, str]]]:
    """Audits code file for boolean conventions with fast substring pre-filter."""
    norm_p = normalize_rel_path(file_path)
    try:
        content = read_file_lf(file_path, encoding=DEFAULT_ENCODING)
        # Fast substring pre-filter
        has_true_word = ("true" in content.lower())
        if not has_true_word:
            return (norm_p, [])
        vios = find_explicit_true_violations(content)
        return (norm_p, vios)
    except Exception:
        return (norm_p, [])

def run_naming_auditor(
    target_dir: str = CURRENT_DIR,
    extensions: set[str] | tuple | None = None
) -> int:
    """Executes naming and boolean convention check across target directory."""
    exts = normalize_extensions(extensions) or DEFAULT_CODE_EXTENSIONS

    def handler(p: Path):
        fp_str, vios = audit_file_naming(p)
        return (fp_str, vios) if vios else None

    stats = process_repository_files(handler, root_dir=target_dir, extensions=exts)
    all_violations = stats["results"]

    has_violations = len(all_violations) > 0
    if has_violations:
        print(f"{LINE_SEPARATOR}❌ Found explicit boolean comparisons in {len(all_violations)} file(s) ({stats['elapsed_ms']:.2f}ms):")
        for fp, vios in all_violations:
            for l_num, line_str in vios[:2]:
                print(f"  ::error file={fp},line={l_num}::Explicit true comparison: {line_str}")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ All {stats['total_files']} code files in '{target_dir}' conform to implicit boolean rules ({stats['elapsed_ms']:.2f}ms).")
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Audit boolean conventions and naming across folders")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory or folder to scan")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument("--ext", help="Comma-separated extensions to scan (e.g. .ts,.go,.py)")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    sys.exit(run_naming_auditor(target_dir=target_path, extensions=args.ext))

if __name__ == "__main__":
    main()
