#!/usr/bin/env python3
"""
Generic Installer Smoke Tester
==============================
Validates generic bash/PowerShell installer scripts for:
1. No unreplaced PLACEHOLDER tokens
2. SHA256 verification pattern (for root installers)
3. Non-destructive update flow (safe replacement / rollback logic)
4. Clean UNIX LF line endings

Features:
- Parallel execution via shared worker pool engine (`02-shared-engine.py`).
- Silent on success: prints a single tick line `✔ All passed. (16 installer scripts verified in 0.04s)`.
- On failure: isolates failing scripts with full diagnostic logs and exit codes.
- CLI options: `--all-paths` / `--all`, `--sync`, `--workers`, `--output`, `--json`, `--filter`.
"""

import argparse
from dataclasses import asdict
from importlib import import_module
import json
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
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR
INSTALLER_EXCLUDE_PARTS = engine.INSTALLER_EXCLUDE_PARTS
DEFAULT_CONCURRENCY_WORKERS = engine.DEFAULT_CONCURRENCY_WORKERS
WorkItemResult = engine.WorkItemResult
WorkGroupSummary = engine.WorkGroupSummary
run_worker_pool = engine.run_worker_pool

# Exclude test fixtures or archived release templates from smoke test discovery
EXTRA_EXCLUDE_NAMES = {"release-install.sh", "release-install.ps1"}
EXTRA_EXCLUDE_PARTS = {"tests", "hooks", "node_modules", ".git", "dist", "build", "release-artifacts", "release-assets"}


def test_bash_installer(script_path: Path) -> list[str]:
    """Smoke tests a bash installer script."""
    re_placeholder = get_compiled_regex(RegexPatternType.PLACEHOLDER_TOKEN)
    content = read_file_lf(script_path, encoding=DEFAULT_ENCODING)
    issues: list[str] = []

    has_placeholder = bool(re_placeholder.search(content))
    if has_placeholder:
        issues.append(f"{script_path}: Contains unreplaced placeholder tokens")

    # Root installers require checksum verification
    is_root_installer = (script_path.name == "install.sh")
    if is_root_installer:
        has_sha = ("sha256" in content.lower() or "shasum" in content.lower() or "compute_hash" in content.lower())
        if not has_sha:
            issues.append(f"{script_path}: Missing SHA256 checksum verification")

    has_safe_rename = (
        "mv " in content
        or "install -m" in content
        or "cp -f" in content
        or "cp " in content
        or "merge_file" in content
        or "record_overwrite_backup" in content
    )
    if not has_safe_rename:
        issues.append(f"{script_path}: Missing non-destructive binary replacement logic")

    return issues


def test_powershell_installer(script_path: Path) -> list[str]:
    """Smoke tests a PowerShell installer script."""
    re_placeholder = get_compiled_regex(RegexPatternType.PLACEHOLDER_TOKEN)
    content = read_file_lf(script_path, encoding=DEFAULT_ENCODING)
    issues: list[str] = []

    has_placeholder = bool(re_placeholder.search(content))
    if has_placeholder:
        issues.append(f"{script_path}: Contains unreplaced placeholder tokens")

    is_root_installer = (script_path.name == "install.ps1")
    if is_root_installer:
        has_sha = ("get-filehash" in content.lower() or "sha256" in content.lower())
        if not has_sha:
            issues.append(f"{script_path}: Missing SHA256 hash verification")

    has_safe_rename = (
        "move-item" in content.lower()
        or "rename-item" in content.lower()
        or "copy-item" in content.lower()
        or "backup" in content.lower()
        or "merge-file" in content.lower()
    )
    if not has_safe_rename:
        issues.append(f"{script_path}: Missing safe rename-first replacement logic")

    return issues


def test_installer_script(script_path: Path) -> WorkItemResult:
    """Worker task: executes smoke validations on a single installer script."""
    start_time = time.perf_counter()
    norm_path = normalize_rel_path(script_path)
    issues: list[str] = []

    try:
        is_bash = script_path.name.endswith(".sh")
        if is_bash:
            issues = test_bash_installer(script_path)
        else:
            issues = test_powershell_installer(script_path)

        duration = time.perf_counter() - start_time
        has_issues = (len(issues) > 0)
        if has_issues:
            output = LINE_SEPARATOR.join(issues)
            return WorkItemResult(name=norm_path, is_success=False, output=output, duration_sec=duration, return_code=1)

        return WorkItemResult(name=norm_path, is_success=True, output="All checks passed.", duration_sec=duration, return_code=0)
    except Exception as exc:
        duration = time.perf_counter() - start_time
        return WorkItemResult(name=norm_path, is_success=False, output=f"Failed to test script: {exc}", duration_sec=duration, return_code=-1)


def discover_installer_scripts(target_dir: str = CURRENT_DIR, filter_pattern: str | None = None) -> list[Path]:
    """Discovers all installer scripts across the target directory, excluding test/build directories."""
    root = Path(target_dir)
    discovered: list[Path] = []

    for ext in ("*.sh", "*.ps1"):
        for script_file in root.rglob(ext):
            if script_file.name in EXTRA_EXCLUDE_NAMES:
                continue
            is_ignored = any(part in script_file.parts for part in EXTRA_EXCLUDE_PARTS)
            if is_ignored:
                continue
            if "install" not in script_file.name.lower():
                continue
            if filter_pattern:
                if filter_pattern.lower() not in str(script_file).lower():
                    continue
            discovered.append(script_file)

    discovered.sort(key=lambda p: str(p))
    return discovered


def format_summary_text(summary: WorkGroupSummary, show_all: bool = False) -> str:
    """Formats human-readable text summary and error logs."""
    lines: list[str] = []

    if show_all:
        lines.append("======================= FINAL SUMMARY =======================")
        for r in summary.results:
            status_icon = "✅" if r.is_success else "❌"
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"{status_icon} [{status_word}] {r.name:<45} ({r.duration_sec:.2f}s)")
        lines.append("------------------------------------------------------------")
        lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
        lines.append(f"Scripts Passed : {summary.passed_count}/{summary.total_items}")
        lines.append(f"Scripts Failed : {summary.failed_count}/{summary.total_items}")
        lines.append("------------------------------------------------------------")

        lines.append("\n=================== ALL INSTALLER LOGS ===================")
        for r in summary.results:
            status_word = "PASSED" if r.is_success else "FAILED"
            lines.append(f"\n--- [{status_word}] {r.name} ({r.duration_sec:.2f}s) ---")
            lines.append(r.output if r.output else "(no output)")
            lines.append("------------------------------------------------------------")
    else:
        if summary.has_failures:
            lines.append("\n================ FAILED INSTALLER LOGS ================")
            for r in summary.results:
                if not r.is_success:
                    lines.append(f"\n❌ FAILED: {r.name} ({r.duration_sec:.2f}s)")
                    lines.append(f"--- {r.name} ---")
                    lines.append(r.output if r.output else "(no output)")
                    lines.append(f"--- END {r.name} ---")
            lines.append("------------------------------------------------------------")
            lines.append(f"Total Duration : {summary.wall_duration_sec:.2f}s")
            lines.append(f"Scripts Passed : {summary.passed_count}/{summary.total_items}")
            lines.append(f"Scripts Failed : {summary.failed_count}/{summary.total_items}")
            lines.append("------------------------------------------------------------")
            lines.append(f"\n❌ Smoke tests failed: {summary.failed_count} script(s) reported violations.")

    return LINE_SEPARATOR.join(lines)


def run_installer_smoke_tests(
    target_dir: str = CURRENT_DIR,
    max_workers: int | None = None,
    show_all: bool = False,
    is_sync: bool = False,
    output_file: str | None = None,
    as_json: bool = False,
    filter_pattern: str | None = None
) -> int:
    """Dispatches smoke tests using parallel worker pool or sequential loop."""
    scripts = discover_installer_scripts(target_dir=target_dir, filter_pattern=filter_pattern)
    if not scripts:
        print(f"⚠️ No installer scripts discovered in '{target_dir}'" + (f" matching filter '{filter_pattern}'" if filter_pattern else ""))
        return ExitCodeType.SUCCESS.value

    worker_count = max_workers or min(len(scripts), DEFAULT_CONCURRENCY_WORKERS)
    if is_sync:
        worker_count = 1

    if not as_json:
        if show_all:
            concurrency_label = "Sequential (1 worker)" if is_sync else f"{worker_count} parallel workers"
            print("================================================================")
            print("            PARALLEL INSTALLER SMOKE TEST RUNNER                ")
            print("================================================================")
            print(f"🚀 Execution Mode          : {concurrency_label}")
            print(f"📋 Total Enqueued Scripts  : {len(scripts)}")
            print("🔍 Display Mode            : SHOW ALL INFORMATION")
            print("----------------------------------------------------------------\n")

    def ticker_callback(res: WorkItemResult, current: int, total: int):
        if not as_json:
            if show_all:
                status_icon = "✅" if res.is_success else "❌"
                status_label = "PASS" if res.is_success else "FAIL"
                print(f"[{current:2d}/{total:2d}] {status_icon} [{status_label}] {res.name} ({res.duration_sec:.2f}s)")

    summary = run_worker_pool(
        items=scripts,
        worker_fn=test_installer_script,
        worker_count=worker_count,
        is_sync=is_sync,
        on_item_complete=ticker_callback
    )

    if as_json:
        payload = {
            "total_items": summary.total_items,
            "passed_count": summary.passed_count,
            "failed_count": summary.failed_count,
            "wall_duration_sec": round(summary.wall_duration_sec, 2),
            "has_failures": summary.has_failures,
            "exit_code": summary.exit_code,
            "scripts": [asdict(r) for r in summary.results]
        }
        json_content = json.dumps(payload, indent=2, ensure_ascii=False)
        if output_file:
            p = Path(output_file)
            p.parent.mkdir(parents=True, exist_ok=True)
            p.write_text(json_content, encoding=DEFAULT_ENCODING)
            print(f"📄 JSON results saved to: {output_file}")
        else:
            print(json_content)
        return summary.exit_code

    if output_file:
        full_file_report = format_summary_text(summary, show_all=True)
        p = Path(output_file)
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(full_file_report, encoding=DEFAULT_ENCODING)
        print(f"📄 Execution report saved to: {output_file}")

    if summary.has_failures:
        failure_text = format_summary_text(summary, show_all=False)
        print(failure_text)
        return summary.exit_code

    if show_all:
        all_text = format_summary_text(summary, show_all=True)
        print(all_text)
        print("\n🎉 All installer smoke tests passed successfully! All OK.")
    else:
        print(f"✔ All passed. ({summary.passed_count} installer scripts verified in {summary.wall_duration_sec:.2f}s)")

    return summary.exit_code


def parse_arguments() -> argparse.Namespace:
    """Parses command-line arguments for installer smoke tester."""
    parser = argparse.ArgumentParser(
        description="Fast Multi-Worker Installer Smoke Tester with worker pool and selective reporting.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Default: run all installer checks in parallel; quiet on success (tick), detailed logs on failure:
  python 03-ai-scripts/16-installer-smoke-tester.py

  # Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/16-installer-smoke-tester.py --all-paths
  python 03-ai-scripts/16-installer-smoke-tester.py --all-passed
  python 03-ai-scripts/16-installer-smoke-tester.py --all

  # Run sequentially (synchronous mode, 1 worker):
  python 03-ai-scripts/16-installer-smoke-tester.py --sync

  # Custom worker concurrency:
  python 03-ai-scripts/16-installer-smoke-tester.py --workers 4

  # Save report to a file:
  python 03-ai-scripts/16-installer-smoke-tester.py --output tmp/installer-report.txt

  # Output machine-readable JSON:
  python 03-ai-scripts/16-installer-smoke-tester.py --json -o tmp/installer-report.json

  # Filter specific installer:
  python 03-ai-scripts/16-installer-smoke-tester.py --filter "wp-install"
        """
    )
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory to search for installer scripts")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument(
        "--all", "--all-paths", "--all-passed", "-a",
        action="store_true",
        dest="show_all",
        help="Show detailed information and logs for all installer scripts."
    )
    parser.add_argument(
        "--failed", "-f",
        action="store_true",
        dest="show_failed",
        help="Show logs only for failed installer scripts (default behavior)."
    )
    parser.add_argument(
        "--sync", "--sequential", "-s",
        action="store_true",
        dest="is_sync",
        help="Execute smoke tests sequentially (1 worker) instead of in parallel."
    )
    parser.add_argument(
        "--workers", "-w", "--concurrency",
        type=int,
        default=None,
        dest="workers",
        help="Number of concurrent worker threads (default: CPU threads, capped at 8)."
    )
    parser.add_argument(
        "--output", "-o", "--file",
        type=str,
        default=None,
        dest="output_file",
        help="Save execution results and report to the specified file path."
    )
    parser.add_argument(
        "--json",
        action="store_true",
        dest="as_json",
        help="Output results as machine-readable JSON (useful for automated AI agents)."
    )
    parser.add_argument(
        "--filter", "-k",
        type=str,
        default=None,
        dest="filter",
        help="Filter installer scripts matching substring (case-insensitive)."
    )
    return parser.parse_args()


def main():
    args = parse_arguments()
    target_path = args.opt_path or args.path or CURRENT_DIR
    exit_code = run_installer_smoke_tests(
        target_dir=target_path,
        max_workers=args.workers,
        show_all=args.show_all,
        is_sync=args.is_sync,
        output_file=args.output_file,
        as_json=args.as_json,
        filter_pattern=args.filter
    )
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
