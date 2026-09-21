#!/usr/bin/env python3
"""Linter to verify boolean guidelines: implicit checks, no explicit comparisons, positive framing, and no mixed polarity."""
import argparse
import os
import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from linter_cache import load_linter_cache, record_linter_success, resolve_linter_targets

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent
TARGET_EXTS = {'.go', '.ts', '.tsx', '.js', '.jsx', '.py', '.php'}
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts',
    '.ai-memory/scratch', '.ai-memory/temp-agents', '03-ai-scripts', 'scripts',
    '04-code', '.tmp'
}

# Regex patterns
EXPLICIT_BOOL_REGEX = re.compile(r'\b(==\s*true|===\s*true|==\s*false|===\s*false)\b')
NEGATIVE_BOOL_NAME_REGEX = re.compile(r'\b(isNot[A-Z]\w*|hasNo[A-Z]\w*)\b')
INVERTED_SUCCESS_REGEX = re.compile(r'!\s*(?:[a-zA-Z0-9_$.->]+\.)?[iI]sSuccess\b')
BANNED_FUNC_PREFIX_REGEX = re.compile(r'\bfunc\s+(?:(?:\([a-zA-Z0-9_*]+\)\s+)?)(can[A-Z]\w*|should[A-Z]\w*|was[A-Z]\w*|will[A-Z]\w*|did[A-Z]\w*|must[A-Z]\w*)\s*\([^)]*\)\s*bool\b')
MIXED_POLARITY_REGEX = re.compile(r'\bif\b[^;{}]*?(?:&&|\band\b)\s*![a-zA-Z0-9_$.->]+')


def strip_comments_and_strings(content: str, ext: str) -> list[tuple[int, str]]:
    out_lines = []
    current_line = []
    line_num = 1
    i = 0
    n = len(content)

    while i < n:
        c = content[i]
        if c == '\n':
            out_lines.append((line_num, ''.join(current_line)))
            current_line = []
            line_num += 1
            i += 1
            continue

        if ext in ('.go', '.ts', '.tsx', '.js', '.jsx', '.php'):
            if c == '/' and i + 1 < n and content[i + 1] == '/':
                i += 2
                while i < n and content[i] != '\n':
                    i += 1
                continue
            if c == '/' and i + 1 < n and content[i + 1] == '*':
                i += 2
                while i < n:
                    if content[i] == '\n':
                        out_lines.append((line_num, ''.join(current_line)))
                        current_line = []
                        line_num += 1
                        i += 1
                    elif content[i] == '*' and i + 1 < n and content[i + 1] == '/':
                        i += 2
                        break
                    else:
                        i += 1
                continue
            if c == '`':
                i += 1
                while i < n and content[i] != '`':
                    if content[i] == '\n':
                        out_lines.append((line_num, ''.join(current_line)))
                        current_line = []
                        line_num += 1
                    i += 1
                if i < n:
                    i += 1
                continue
            if c == '"':
                i += 1
                while i < n and content[i] != '"':
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue
            if c == "'":
                i += 1
                while i < n and content[i] != "'":
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue
        elif ext == '.py':
            if c == '#':
                while i < n and content[i] != '\n':
                    i += 1
                continue
            if c == '"' and i + 2 < n and content[i + 1] == '"' and content[i + 2] == '"':
                i += 3
                while i + 2 < n and not (content[i] == '"' and content[i + 1] == '"' and content[i + 2] == '"'):
                    if content[i] == '\n':
                        out_lines.append((line_num, ''.join(current_line)))
                        current_line = []
                        line_num += 1
                    i += 1
                i += 3
                continue
            if c == '"':
                i += 1
                while i < n and content[i] != '"':
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue
            if c == "'":
                i += 1
                while i < n and content[i] != "'":
                    if content[i] == '\\':
                        i += 2
                    else:
                        i += 1
                if i < n:
                    i += 1
                continue

        current_line.append(c)
        i += 1

    if current_line:
        out_lines.append((line_num, ''.join(current_line)))

    return out_lines


def scan_file(filepath: Path) -> list[tuple[int, str]]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception as e:
        return [(0, f"Error reading file: {e}")]

    ext = filepath.suffix.lower()
    stripped_lines = strip_comments_and_strings(content, ext)
    violations = []
    in_fail_method = False

    for line_num, line in stripped_lines:
        s = line.strip()
        if not s:
            continue

        is_fail_def = False
        if ext == '.go':
            if re.search(r'\bfunc\s+(?:\([^)]+\)\s+)?(?:IsFail|IsFailed)\b', line):
                in_fail_method = True
                is_fail_def = True

        # 1. Explicit boolean comparisons
        m_exp = EXPLICIT_BOOL_REGEX.search(line)
        if m_exp:
            violations.append((line_num, f"Explicit boolean comparison ({m_exp.group(1)}): {s}"))

        # 2. Negative boolean naming
        m_neg = NEGATIVE_BOOL_NAME_REGEX.search(line)
        if m_neg:
            violations.append((line_num, f"Negative boolean variable name ({m_neg.group(1)}): {s}"))

        # 3. Inverted success check
        if not in_fail_method:
            m_succ = INVERTED_SUCCESS_REGEX.search(line)
            if m_succ:
                violations.append((line_num, f"Inverted success check (!isSuccess): {s}"))

        # 4. Banned function prefixes returning bool
        m_func = BANNED_FUNC_PREFIX_REGEX.search(line)
        if m_func:
            violations.append((line_num, f"Banned function prefix returning bool ({m_func.group(1)}): {s}"))

        if ext == '.go' and in_fail_method and '}' in line and (not is_fail_def or '{' in line):
            in_fail_method = False

    return violations


def parse_cli_args() -> argparse.Namespace:
    """Parses command line arguments for boolean guidelines linter."""
    parser = argparse.ArgumentParser(description="Check boolean guidelines across repository.")
    parser.add_argument("--all", "--force", dest="force_all", action="store_true", help="Scan all repository files")

    return parser.parse_args()


def scan_target_files(targets: list[Path]) -> dict[str, list[tuple[int, str]]]:
    """Scans all target files and collects boolean guideline violations."""
    all_violations: dict[str, list[tuple[int, str]]] = {}
    for p in targets:
        v = scan_file(p)
        if v:
            rel = p.relative_to(ROOT_DIR).as_posix()
            all_violations[rel] = v

    return all_violations


def report_boolean_violations(all_violations: dict[str, list[tuple[int, str]]], desc: str, total_files: int) -> int:
    """Reports boolean violations and exits with appropriate status code."""
    if all_violations:
        total = sum(len(v) for v in all_violations.values())
        print(f"\n❌ FAIL: Found {total} boolean guideline violation(s) across {len(all_violations)} file(s):\n")
        for rel_path, v_list in sorted(all_violations.items()):
            for line_no, msg in v_list:
                print(f"  {rel_path}:{line_no}: {msg}")
        return 1
    file_info = f" across {total_files:,} files" if total_files > 0 else ""
    print(f"\n✅ PASS (0 boolean guideline violations{file_info}) [{desc}]")

    return 0


def report_cached_boolean_pass(desc: str) -> None:
    """Reports clean pass when zero changed files need scanning."""
    cache = load_linter_cache("check-boolean-guidelines", ROOT_DIR)
    cached_count = len(cache.get("file_hashes", {}))
    file_info = f" across {cached_count:,} files" if cached_count > 0 else ""
    print(f"\n✅ PASS (0 boolean guideline violations{file_info}) [{desc}]")
    sys.exit(0)


def execute_boolean_run(targets: list[Path], desc: str, is_incremental: bool) -> None:
    """Executes target scanning, reporting, and records cache on success."""
    all_violations = scan_target_files(targets)
    code = report_boolean_violations(all_violations, desc, len(targets))
    if code != 0:
        sys.exit(code)
    record_linter_success("check-boolean-guidelines", ROOT_DIR, targets, is_incremental)
    sys.exit(0)


def main() -> None:
    """Main execution function for boolean guidelines linter."""
    args = parse_cli_args()
    targets, desc, is_incremental = resolve_linter_targets(
        "check-boolean-guidelines", ROOT_DIR, TARGET_EXTS, EXCLUDE_DIRS, force_all=args.force_all
    )
    print(f"=== Running Boolean Guidelines Linter (check-boolean-guidelines.py) [{desc}] in {ROOT_DIR} ===")
    if not targets:
        report_cached_boolean_pass(desc)
    execute_boolean_run(targets, desc, is_incremental)


if __name__ == "__main__":
    main()

