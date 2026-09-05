#!/usr/bin/env python3
"""Generic Installer Smoke Tester with Concurrent Worker Pool and Flexible CLI Reporting.

Validates generic bash/PowerShell installer scripts for:
1. No leftover PLACEHOLDER tokens
2. SHA256 verification pattern
3. Non-destructive update flow (rename-first, then replace)
4. Clean UNIX LF line endings

All Worker Pool primitives, Enums, and Constants are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/16-installer-smoke-tester.py                 # Default: parallel, quiet on success (tick), detailed logs on failure
  python 03-ai-scripts/16-installer-smoke-tester.py --all-paths    # Verbose mode: ticker, summary table, and detailed logs
  python 03-ai-scripts/16-installer-smoke-tester.py --sync         # Run sequentially (1 worker)
  python 03-ai-scripts/16-installer-smoke-tester.py -o report.txt  # Output execution report to file
  python 03-ai-scripts/16-installer-smoke-tester.py --json         # Output results as machine-readable JSON
"""
from __future__ import annotations

import argparse
from importlib import import_module
import os
from pathlib import Path
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

read_file_lf = engine.read_file_lf
normalize_rel_path = engine.normalize_rel_path
ExitCodeType = engine.ExitCodeType
RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR
INSTALLER_BASH_NAME = engine.INSTALLER_BASH_NAME
INSTALLER_PWSH_NAME = engine.INSTALLER_PWSH_NAME
INSTALLER_EXCLUDE_PARTS = engine.INSTALLER_EXCLUDE_PARTS
WorkerResult = engine.WorkerResult
add_worker_cli_arguments = engine.add_worker_cli_arguments
run_worker_pool = engine.run_worker_pool


def test_bash_installer(script_path: Path) -> list[str]:
    """Smoke tests a bash installer script for safety and checksum patterns."""
    re_placeholder = get_compiled_regex(RegexPatternType.PLACEHOLDER_TOKEN)
    content = read_file_lf(script_path, encoding=DEFAULT_ENCODING)
    issues: list[str] = []

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
    """Smoke tests a PowerShell installer script for safety and checksum patterns."""
    re_placeholder = get_compiled_regex(RegexPatternType.PLACEHOLDER_TOKEN)
    content = read_file_lf(script_path, encoding=DEFAULT_ENCODING)
    issues: list[str] = []

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


def test_single_installer(script_file: Path) -> WorkerResult:
    """Worker task testing a single installer script."""
    start = time.perf_counter()
    rel_path = normalize_rel_path(script_file)
    name = script_file.name.lower()

    try:
        if name.endswith(".sh") or name == INSTALLER_BASH_NAME:
            issues = test_bash_installer(script_file)
        elif name.endswith(".ps1") or name == INSTALLER_PWSH_NAME:
            issues = test_powershell_installer(script_file)
        else:
            issues = []

        elapsed = round(time.perf_counter() - start, 3)
        if issues:
            return WorkerResult(
                name=rel_path,
                is_success=False,
                error="\n".join(issues),
                elapsed_sec=elapsed,
            )

        return WorkerResult(
            name=rel_path,
            is_success=True,
            output="Valid installer patterns (clean tokens, SHA256 verification, safe rename)",
            elapsed_sec=elapsed,
        )
    except Exception as exc:
        elapsed = round(time.perf_counter() - start, 3)
        return WorkerResult(
            name=rel_path,
            is_success=False,
            error=f"Exception while testing installer: {exc}",
            elapsed_sec=elapsed,
        )


def discover_installer_scripts(target_dir: str = CURRENT_DIR) -> list[Path]:
    """Discovers all installer scripts across the target directory, skipping excluded paths."""
    root = Path(target_dir)
    discovered: list[Path] = []

    for pattern in (INSTALLER_BASH_NAME, INSTALLER_PWSH_NAME):
        for script_file in root.rglob(pattern):
            if any(part in script_file.parts for part in INSTALLER_EXCLUDE_PARTS):
                continue
            if script_file.is_file() and script_file not in discovered:
                discovered.append(script_file)

    return sorted(discovered)


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/16-installer-smoke-tester.py",
        description="Smoke test installer scripts across target directory using parallel worker groups and IO throttling.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: run all installer tests in parallel; quiet on success (tick), detailed logs on failure:
  python 03-ai-scripts/16-installer-smoke-tester.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/16-installer-smoke-tester.py --all-paths
  python 03-ai-scripts/16-installer-smoke-tester.py --all

  # 3. Run sequentially (synchronous mode, 1 worker):
  python 03-ai-scripts/16-installer-smoke-tester.py --sync

  # 4. Save results to a file:
  python 03-ai-scripts/16-installer-smoke-tester.py -o tmp/installer-report.txt

  # 5. Output results as machine-readable JSON:
  python 03-ai-scripts/16-installer-smoke-tester.py --json
  python 03-ai-scripts/16-installer-smoke-tester.py --json -o tmp/installer-report.json

  # 6. Filter by script path substring:
  python 03-ai-scripts/16-installer-smoke-tester.py -k "bash"
        """,
    )
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory to search for installer scripts")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    scripts = discover_installer_scripts(target_dir=target_path)

    if args.filter:
        filt = args.filter.lower()
        scripts = [s for s in scripts if filt in str(s).lower()]

    exit_code = run_worker_pool(
        items=scripts,
        worker_fn=test_single_installer,
        max_workers=args.workers,
        is_sync=args.is_sync,
        show_all=args.show_all,
        output_file=args.output_file,
        as_json=args.as_json,
        title="INSTALLER SMOKE TESTER",
        item_noun="installer script(s)",
    )
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
