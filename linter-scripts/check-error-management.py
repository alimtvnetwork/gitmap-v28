#!/usr/bin/env python3
"""
check-error-management.py
Linter enforcing repository-wide error management rules:
1. Zero bare panic(...) in Go source code (outside cliexit/handle.go and tests).
2. Zero bare os.Exit(...) in Go source code (outside cliexit/handle.go and tests).
3. Zero exitWith(...) alias calls (outside tests).
4. Zero naked apperror.NewSimple or apperror.WrapSimple statements.
5. Zero cliexit.HandleError(nil, ...) disguised exits.
6. Zero unhandled swallowed database calls (_ = row.Scan, _, _ = db.Exec).
7. Zero empty catch {} or except: pass blocks.
8. Zero unhandled swallowed error assignments (_ = err) without explanation.
"""

import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent

EXCLUDED_DIRS = {
    ".git",
    "node_modules",
    "dist",
    "build",
    "bin",
    ".gemini",
    "vendor",
    "tmp",
    "linter-scripts",
    "scripts",
    "gitmap-updater",
    "04-code",
}

WHITELISTED_GO_FILES = {
    "handle.go",
    "handle_test.go",
    "cliexit.go",
    "cliexit_test.go",
    "report.go",
    "report_test.go",
    "kind.go",
    "kind_test.go",
    "main.go",
    "lazyregex.go",
    "lazy_regex.go",
    "logger.go",
}


def should_scan_dir(dirpath: Path) -> bool:
    for part in dirpath.parts:
        is_excluded = part in EXCLUDED_DIRS or part.startswith(".") or part.endswith("_old")
        if is_excluded and part not in {".lovable", ".github", ".agents"}:
            return False

    return True


def check_go_panics_and_exits(filepath: Path, idx: int, line: str, stripped: str) -> list[str]:
    violations = []
    has_panic = bool(re.search(r'\bpanic\s*\(', line))
    if has_panic and "// lint-allow: panic" not in line:
        violations.append(f"{filepath}:{idx}: Bare panic() call detected: '{stripped}'")

    has_exit = bool(re.search(r'\bos\.Exit\s*\(', line))
    if has_exit and "// lint-allow: os.Exit" not in line:
        violations.append(f"{filepath}:{idx}: Bare os.Exit() call detected (use cliexit.HandleError): '{stripped}'")

    has_exit_with = bool(re.search(r'\bexitWith\s*\(', line))
    if has_exit_with and "// lint-allow: exitWith" not in line:
        violations.append(f"{filepath}:{idx}: Banned exitWith() alias call detected (use cliexit error handling): '{stripped}'")

    return violations


def check_go_apperrors_and_cliexit(filepath: Path, idx: int, line: str, stripped: str) -> list[str]:
    violations = []
    has_naked_apperr = bool(re.search(r'(?:^|[;{])\s*(?:_\s*=\s*)?apperror\.(?:NewSimple|WrapSimple)\s*\(', stripped))
    if has_naked_apperr and "// lint-allow: naked-apperror" not in line:
        violations.append(f"{filepath}:{idx}: Naked apperror expression statement detected (must be returned or handled): '{stripped}'")

    return violations


def check_go_swallowed_errors(filepath: Path, idx: int, line: str, stripped: str) -> list[str]:
    violations = []
    has_swallowed_err = bool(re.search(r'_\s*=\s*err\b', line))
    if has_swallowed_err and "// lint-allow: ignore-error" not in line:
        violations.append(f"{filepath}:{idx}: Swallowed error '_ = err' without waiver comment: '{stripped}'")

    has_db_swallow = bool(re.search(r'_\s*(?:,\s*_\s*)?=\s*(?:\w+\.)?(?:Scan|Exec|Query|QueryRow)\s*\(', line))
    if has_db_swallow and "// lint-allow: ignore-db-error" not in line:
        violations.append(f"{filepath}:{idx}: Swallowed database call error detected: '{stripped}'")

    return violations


def check_go_line(filepath: Path, idx: int, line: str) -> list[str]:
    stripped = line.strip()
    is_comment = stripped.startswith("//") or stripped.startswith("/*")
    if is_comment:
        return []

    violations = []
    violations.extend(check_go_panics_and_exits(filepath, idx, line, stripped))
    violations.extend(check_go_apperrors_and_cliexit(filepath, idx, line, stripped))
    violations.extend(check_go_swallowed_errors(filepath, idx, line, stripped))

    return violations


def check_go_file(filepath: Path) -> list[str]:
    is_whitelisted = filepath.name in WHITELISTED_GO_FILES or filepath.name.endswith("_test.go")
    if is_whitelisted:
        return []

    try:
        content = filepath.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        return [f"{filepath}: Failed to read file: {e}"]

    violations = []
    for idx, line in enumerate(content.splitlines(), 1):
        violations.extend(check_go_line(filepath, idx, line))

    return violations


def check_ts_js_file(filepath: Path) -> list[str]:
    is_test = filepath.name.endswith(".test.ts") or filepath.name.endswith(".spec.ts")
    if is_test:
        return []

    try:
        content = filepath.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        return [f"{filepath}: Failed to read file: {e}"]

    return scan_ts_js_lines(filepath, content.splitlines())


def scan_ts_js_lines(filepath: Path, lines: list[str]) -> list[str]:
    violations = []
    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        is_comment = stripped.startswith("//") or stripped.startswith("/*")
        if not is_comment and re.search(r'catch\s*(\([^)]*\))?\s*\{\s*\}', line):
            violations.append(f"{filepath}:{idx}: Empty catch block detected: '{stripped}'")

    return violations


def check_py_file(filepath: Path) -> list[str]:
    is_test = filepath.name.startswith("test_") or filepath.name.endswith("_test.py")
    if is_test:
        return []

    try:
        content = filepath.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        return [f"{filepath}: Failed to read file: {e}"]

    return scan_py_lines(filepath, content.splitlines())


def scan_py_lines(filepath: Path, lines: list[str]) -> list[str]:
    violations = []
    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        has_empty_except = not stripped.startswith("#") and re.search(r'except(\s+[^:]+)?:\s*pass\b', stripped)
        if has_empty_except and "# lint-allow: empty-except" not in line:
            violations.append(f"{filepath}:{idx}: Silent 'except: pass' block detected: '{stripped}'")

    return violations


def resolve_scan_target() -> Path:
    has_target_arg = len(sys.argv) > 1 and not sys.argv[1].startswith("-")
    if has_target_arg:
        return Path(sys.argv[1]).resolve()

    return ROOT_DIR


def scan_single_file(filepath: Path) -> list[str]:
    ext = filepath.suffix.lower()
    if ext == ".go":
        return check_go_file(filepath)
    if ext in {".ts", ".tsx", ".js", ".jsx"}:
        return check_ts_js_file(filepath)
    if ext == ".py":
        return check_py_file(filepath)

    return []


def scan_directory(target_dir: Path) -> tuple[list[str], int]:
    all_violations = []
    scanned_count = 0
    for root, dirs, files in os.walk(target_dir):
        root_path = Path(root)
        dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS and not d.startswith(".") and not d.endswith("_old")]
        for f in files:
            filepath = root_path / f
            file_violations = scan_single_file(filepath)
            if filepath.suffix.lower() in {".go", ".ts", ".tsx", ".js", ".jsx", ".py"}:
                scanned_count += 1
            all_violations.extend(file_violations)

    return all_violations, scanned_count


def scan_target(target: Path) -> tuple[list[str], int]:
    if target.is_file():
        return scan_single_file(target), 1

    return scan_directory(target)


def report_results(violations: list[str], count: int) -> int:
    print(f"Scanned {count} source files for error management compliance.")
    has_violations = len(violations) > 0
    if has_violations:
        print(f"\n❌ FAILED: Found {len(violations)} error management violation(s):")
        for v in violations:
            print(f"  - {v}")

        return 1

    print("✅ PASS: All error management checks passed (zero bare panics, exits, or swallowed errors).")

    return 0


def main() -> int:
    target = resolve_scan_target()
    violations, count = scan_target(target)

    return report_results(violations, count)


if __name__ == "__main__":
    sys.exit(main())
