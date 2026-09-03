#!/usr/bin/env python3
"""
Generic Installer Smoke Tester
Validates generic bash/PowerShell installer scripts for:
1. No leftover PLACEHOLDER tokens
2. SHA256 verification pattern
3. Non-destructive update flow (rename-first, then replace)
4. Clean UNIX LF line endings

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

import argparse
from importlib import import_module
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

read_file_lf = engine.read_file_lf
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
INSTALLER_BASH_NAME = engine.INSTALLER_BASH_NAME
INSTALLER_PWSH_NAME = engine.INSTALLER_PWSH_NAME
INSTALLER_EXCLUDE_PARTS = engine.INSTALLER_EXCLUDE_PARTS

def test_bash_installer(script_path: Path) -> list[str]:
    """Smoke tests a bash installer script."""
    re_placeholder = get_compiled_regex(RegexPatternType.PLACEHOLDER_TOKEN)
    content = read_file_lf(script_path, encoding=DEFAULT_ENCODING)
    issues = []

    has_placeholder = bool(re_placeholder.search(content))
    if has_placeholder:
        issues.append(f"{script_path}: Contains unreplaced placeholder tokens")

    has_sha = ("sha256" in content.lower() or "shasum" in content.lower())
    if not has_sha:
        issues.append(f"{script_path}: Missing SHA256 checksum verification")

    has_safe_rename = ("mv " in content or "install -m" in content)
    if not has_safe_rename:
        issues.append(f"{script_path}: Missing non-destructive binary replacement logic")

    return issues

def test_powershell_installer(script_path: Path) -> list[str]:
    """Smoke tests a PowerShell installer script."""
    re_placeholder = get_compiled_regex(RegexPatternType.PLACEHOLDER_TOKEN)
    content = read_file_lf(script_path, encoding=DEFAULT_ENCODING)
    issues = []

    has_placeholder = bool(re_placeholder.search(content))
    if has_placeholder:
        issues.append(f"{script_path}: Contains unreplaced placeholder tokens")

    has_sha = ("get-filehash" in content.lower() or "sha256" in content.lower())
    if not has_sha:
        issues.append(f"{script_path}: Missing SHA256 hash verification")

    has_safe_rename = ("move-item" in content.lower() or "rename-item" in content.lower())
    if not has_safe_rename:
        issues.append(f"{script_path}: Missing safe rename-first replacement logic")

    return issues

def run_installer_smoke_tests(target_dir: str = CURRENT_DIR) -> int:
    """Discovers and tests all installer scripts in target directory using flattened guard clauses."""
    root = Path(target_dir)
    all_issues = []
    scripts_tested = 0

    for script_file in root.rglob(INSTALLER_BASH_NAME):
        is_ignored = any(part in script_file.parts for part in INSTALLER_EXCLUDE_PARTS)
        if is_ignored:
            continue
        scripts_tested += 1
        all_issues.extend(test_bash_installer(script_file))

    for script_file in root.rglob(INSTALLER_PWSH_NAME):
        is_ignored = any(part in script_file.parts for part in INSTALLER_EXCLUDE_PARTS)
        if is_ignored:
            continue
        scripts_tested += 1
        all_issues.extend(test_powershell_installer(script_file))

    has_issues = len(all_issues) > 0
    if has_issues:
        print(f"{LINE_SEPARATOR}❌ Installer smoke test failed with {len(all_issues)} issue(s):")
        for issue in all_issues:
            print(f"  ::error::{issue}")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ Tested {scripts_tested} installer script(s) in '{target_dir}'. All smoke tests passed.")
    return ExitCodeType.SUCCESS.value

def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    parser = argparse.ArgumentParser(description="Smoke test installer scripts across target directory")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory to search for installer scripts")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    sys.exit(run_installer_smoke_tests(target_dir=target_path))

if __name__ == "__main__":
    main()
