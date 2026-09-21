#!/usr/bin/env python3
"""Linter to verify newline styling and strict Unix LF line endings across repository source files."""
import argparse
from pathlib import Path
import sys
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

sys.path.insert(0, str(Path(__file__).resolve().parent))
from linter_cache import resolve_linter_targets, record_linter_success

ROOT_DIR = Path(__file__).resolve().parent.parent
TARGET_SUBDIR = "src" if (ROOT_DIR / "src").is_dir() else ""
TARGET_EXTS = {'.go', '.ts', '.tsx', '.js', '.jsx'}
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts',
    '.ai-memory/scratch', '.ai-memory/temp-agents', '03-ai-scripts', 'scripts',
    '04-code', '.tmp'
}


def check_empty_streak(stripped: str, empty_streak: int, line_idx: int) -> tuple[int, list[tuple[int, str]]]:
    """Checks for consecutive double empty lines."""
    if stripped == "":
        new_streak = empty_streak + 1
        if new_streak == 2:
            return new_streak, [(line_idx + 1, "No double empty lines (\\n\\n\\n) allowed")]
        return new_streak, []

    return 0, []


def check_block_start(stripped: str, next_line: str, line_idx: int) -> list[tuple[int, str]]:
    """Checks that a block starting with { is not followed by an empty line."""
    if stripped.endswith("{") and next_line.strip() == "":
        return [(line_idx + 2, "No empty line at the start of a function or block")]

    return []


def is_valid_prev_return_line(prev: str) -> bool:
    """Checks if previous line is a valid preceding line for return statement."""
    has_empty = prev == ""
    has_brace = prev.endswith("{") or prev.endswith("}") or prev.endswith(":")
    has_comment = prev.startswith("//") or prev.startswith("/*") or prev.startswith("*")

    return has_empty or has_brace or has_comment


def check_return_spacing(stripped: str, prev_line: str, line_idx: int) -> list[tuple[int, str]]:
    """Checks that return statement is preceded by a blank line or structural delimiter."""
    is_return = stripped.startswith("return ") or stripped == "return"
    if is_return and line_idx > 0:
        is_valid = is_valid_prev_return_line(prev_line.strip())
        if is_valid:
            return []
        return [(line_idx + 1, "Blank line required before return")]

    return []


def is_continuation_line(next_line: str) -> bool:
    """Checks if next line continues the block statement."""
    prefixes = ('}', 'else', 'catch', 'finally', ')', ']', ',', ';', '//', '/*', '</', '|', '&')

    return next_line.startswith(prefixes)


def check_closing_brace_spacing(stripped: str, next_line: str, line_idx: int) -> list[tuple[int, str]]:
    """Checks that closing brace is followed by a blank line if more code follows."""
    if stripped == "}":
        next_stripped = next_line.strip()
        if next_stripped != "":
            if is_continuation_line(next_stripped):
                return []
            return [(line_idx + 1, "Blank line required after '}' if followed by more code")]

    return []


def check_go_newline_literal(line: str, is_go_file: bool, line_idx: int) -> list[tuple[int, str]]:
    """Checks for literal newline string in Go files."""
    if is_go_file and '"\\n"' in line:
        return [(line_idx + 1, 'Use constants.NewLineUnix instead of "\\n"')]

    return []


def inspect_line_violations(line: str, prev_line: str, next_line: str, is_go: bool, idx: int) -> list[tuple[int, str]]:
    """Inspects block start, brace spacing, return spacing, and go newline literal."""
    stripped = line.strip()
    violations: list[tuple[int, str]] = []
    if next_line != "":
        violations.extend(check_block_start(stripped, next_line, idx))
        violations.extend(check_closing_brace_spacing(stripped, next_line, idx))
    violations.extend(check_return_spacing(stripped, prev_line, idx))
    violations.extend(check_go_newline_literal(line, is_go, idx))

    return violations


def check_file_content_lines(filepath: Path, lines: list[str]) -> list[tuple[int, str]]:
    """Checks newline styling rules across individual lines of a file."""
    violations: list[tuple[int, str]] = []
    empty_streak = 0
    total = len(lines)
    is_go = filepath.suffix == ".go"
    for i, line in enumerate(lines):
        empty_streak, errs = check_empty_streak(line.strip(), empty_streak, i)
        violations.extend(errs)
        prev = lines[i - 1] if i > 0 else ""
        nxt = lines[i + 1] if i + 1 < total else ""
        violations.extend(inspect_line_violations(line, prev, nxt, is_go, i))

    return violations


def check_file(filepath: Path) -> list[tuple[int, str]]:
    """Reads and validates line endings and newline styling of a file."""
    try:
        raw_bytes = filepath.read_bytes()
    except (OSError, IOError) as exc:
        print(f"Error reading {filepath}: {exc}", file=sys.stderr)
        return []
    if b'\r' in raw_bytes:
        return [(1, "File contains Windows CRLF or CR line endings; strictly Unix LF (\\n) required")]
    content = raw_bytes.decode('utf-8', errors='replace')
    lines = content.split('\n')

    return check_file_content_lines(filepath, lines)


def parse_cli_args() -> argparse.Namespace:
    """Parses command line arguments for newline styling linter."""
    parser = argparse.ArgumentParser(description="Check newline styling and Unix LF across repository.")
    parser.add_argument("--all", "--force", dest="force_all", action="store_true", help="Scan all repository files")

    return parser.parse_args()


def scan_target_files(targets: list[Path]) -> dict[str, list[tuple[int, str]]]:
    """Scans all target files and collects newline styling violations."""
    all_violations: dict[str, list[tuple[int, str]]] = {}
    for p in targets:
        v = check_file(p)
        if v:
            rel = p.relative_to(ROOT_DIR).as_posix()
            all_violations[rel] = v

    return all_violations


def report_newline_violations(all_violations: dict[str, list[tuple[int, str]]], desc: str) -> int:
    """Reports newline styling violations and exits with appropriate status code."""
    if all_violations:
        total = sum(len(v) for v in all_violations.values())
        print(f"\n❌ FAIL: Found {total} newline styling violation(s) across {len(all_violations)} file(s):\n")
        for rel_path, v_list in sorted(all_violations.items()):
            for line_no, msg in v_list:
                print(f"  {rel_path}:{line_no}: {msg}")
        return 1
    print(f"\n✅ PASS (All files Unix LF) [{desc}]")

    return 0


def execute_linter_run(targets: list[Path], desc: str, is_incremental: bool) -> None:
    """Executes target scanning, reporting, and records cache on success."""
    all_violations = scan_target_files(targets)
    code = report_newline_violations(all_violations, desc)
    if code != 0:
        sys.exit(code)
    record_linter_success("check-newline-styling", ROOT_DIR, targets, is_incremental)
    sys.exit(0)


def main() -> None:
    """Main execution function for newline styling linter."""
    args = parse_cli_args()
    targets, desc, is_incremental = resolve_linter_targets(
        "check-newline-styling", ROOT_DIR, TARGET_EXTS, EXCLUDE_DIRS,
        force_all=args.force_all, target_subdir=TARGET_SUBDIR
    )
    print(f"=== Running Newline Styling Linter (check-newline-styling.py) [{desc}] in {ROOT_DIR} ===")
    if not targets:
        print(f"\n✅ PASS (All files Unix LF) [{desc}]")
        sys.exit(0)
    execute_linter_run(targets, desc, is_incremental)




if __name__ == "__main__":
    main()
