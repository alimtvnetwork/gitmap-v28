#!/usr/bin/env python3
"""Fast Version Synchronization & Manifest Guard with Concurrent Worker Pool and Flexible Reporting.

Validates that version.json, package.json, changelog.md, and constants.go are in 100% sync.
All Worker Pool primitives, Enums, and Constants are imported directly from 02-shared-engine.py.

Usage:
  python 03-ai-scripts/14-version-sync-checker.py                 # Default: parallel, quiet on success (tick), detailed on error
  python 03-ai-scripts/14-version-sync-checker.py --all-paths    # Show all information: ticker, summary table, full logs
  python 03-ai-scripts/14-version-sync-checker.py --sync         # Run sequentially (1 worker)
  python 03-ai-scripts/14-version-sync-checker.py -o report.txt  # Output execution report to file
  python 03-ai-scripts/14-version-sync-checker.py --json         # Output results as machine-readable JSON
"""
from __future__ import annotations

import argparse
from dataclasses import dataclass
from importlib import import_module
import json
import os
from pathlib import Path
import re
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR
WorkerResult = engine.WorkerResult
add_worker_cli_arguments = engine.add_worker_cli_arguments
run_worker_pool = engine.run_worker_pool


@dataclass
class VersionCheckTask:
    name: str
    check_type: str
    target_path: Path
    canonical_version: str


def get_canonical_version(root_dir: str = CURRENT_DIR) -> str:
    """Extracts canonical version from version.json or package.json."""
    v_path = os.path.join(root_dir, "version.json")
    p_path = os.path.join(root_dir, "package.json")

    if os.path.exists(v_path):
        with open(v_path, "r", encoding=DEFAULT_ENCODING) as f:
            data = json.load(f)
            return (data.get("Version") or data.get("version") or "").lstrip("v")

    if os.path.exists(p_path):
        with open(p_path, "r", encoding=DEFAULT_ENCODING) as f:
            data = json.load(f)
            return (data.get("version") or "").lstrip("v")

    return ""


def execute_version_check(task: VersionCheckTask) -> WorkerResult:
    """Executes a single version synchronization check."""
    start = time.perf_counter()
    p = task.target_path
    canonical = task.canonical_version

    if not p.exists():
        return WorkerResult(
            name=task.name,
            is_success=False,
            error=f"Target file does not exist: {p}",
            elapsed_sec=round(time.perf_counter() - start, 3),
        )

    try:
        if task.check_type == "package_json":
            with open(p, "r", encoding=DEFAULT_ENCODING) as f:
                data = json.load(f)
            pkg_ver = (data.get("version") or "").lstrip("v")
            if pkg_ver != canonical:
                return WorkerResult(
                    name=task.name,
                    is_success=False,
                    error=f"package.json version '{pkg_ver}' does not match canonical 'v{canonical}'",
                    elapsed_sec=round(time.perf_counter() - start, 3),
                )
            return WorkerResult(
                name=task.name,
                is_success=True,
                output=f"package.json version matches canonical v{canonical}",
                elapsed_sec=round(time.perf_counter() - start, 3),
            )

        elif task.check_type == "changelog":
            with open(p, "r", encoding=DEFAULT_ENCODING) as f:
                content = f.read()
            re_changelog = get_compiled_regex(RegexPatternType.CHANGELOG_HEADER)
            match = re_changelog.search(content)
            if not match:
                return WorkerResult(
                    name=task.name,
                    is_success=False,
                    error="Missing or unparseable release header in changelog.md",
                    elapsed_sec=round(time.perf_counter() - start, 3),
                )
            latest_ver = match.group(1).lstrip("v")
            if latest_ver != canonical:
                return WorkerResult(
                    name=task.name,
                    is_success=False,
                    error=f"changelog.md latest header is 'v{latest_ver}', expected 'v{canonical}'",
                    elapsed_sec=round(time.perf_counter() - start, 3),
                )
            return WorkerResult(
                name=task.name,
                is_success=True,
                output=f"changelog.md latest header matches canonical v{canonical}",
                elapsed_sec=round(time.perf_counter() - start, 3),
            )

        elif task.check_type == "constants_go":
            with open(p, "r", encoding=DEFAULT_ENCODING) as f:
                content = f.read()
            match = re.search(r'var\s+Version\s*=\s*"([^"]+)"', content)
            if not match:
                return WorkerResult(
                    name=task.name,
                    is_success=False,
                    error=f"Could not locate var Version in {p}",
                    elapsed_sec=round(time.perf_counter() - start, 3),
                )
            const_ver = match.group(1).lstrip("v")
            if const_ver != canonical:
                return WorkerResult(
                    name=task.name,
                    is_success=False,
                    error=f"constants.go Version '{const_ver}' does not match canonical 'v{canonical}'",
                    elapsed_sec=round(time.perf_counter() - start, 3),
                )
            return WorkerResult(
                name=task.name,
                is_success=True,
                output=f"constants.go Version matches canonical v{canonical}",
                elapsed_sec=round(time.perf_counter() - start, 3),
            )

        return WorkerResult(
            name=task.name,
            is_success=True,
            output="Unknown check skipped",
            elapsed_sec=round(time.perf_counter() - start, 3),
        )
    except Exception as exc:
        return WorkerResult(
            name=task.name,
            is_success=False,
            error=f"Exception while validating {p}: {exc}",
            elapsed_sec=round(time.perf_counter() - start, 3),
        )


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="python 03-ai-scripts/14-version-sync-checker.py",
        description="Check version synchronization across manifests using parallel worker pool.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # 1. Default: run all version checks in parallel; quiet on success (tick), detailed on error:
  python 03-ai-scripts/14-version-sync-checker.py

  # 2. Show all information (ticker, summary table, full logs):
  python 03-ai-scripts/14-version-sync-checker.py --all-paths
  python 03-ai-scripts/14-version-sync-checker.py --all

  # 3. Run sequentially (1 worker):
  python 03-ai-scripts/14-version-sync-checker.py --sync

  # 4. Save report to a file:
  python 03-ai-scripts/14-version-sync-checker.py -o tmp/version-report.txt

  # 5. Output results as machine-readable JSON:
  python 03-ai-scripts/14-version-sync-checker.py --json
        """,
    )
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory containing version files")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify root directory")
    add_worker_cli_arguments(parser)
    args = parser.parse_args()

    target_path = Path(args.opt_path or args.path or CURRENT_DIR).resolve()
    canonical = get_canonical_version(str(target_path))

    if not canonical:
        print(f"⚠️ No canonical version found in '{target_path}'.")
        sys.exit(0)

    tasks = [
        VersionCheckTask(
            name="Package Manifest Sync (package.json)",
            check_type="package_json",
            target_path=target_path / "package.json",
            canonical_version=canonical,
        ),
        VersionCheckTask(
            name="Changelog Header Sync (changelog.md)",
            check_type="changelog",
            target_path=target_path / "changelog.md",
            canonical_version=canonical,
        ),
    ]

    constants_go = target_path / "gitmap" / "constants" / "constants.go"
    if constants_go.exists():
        tasks.append(
            VersionCheckTask(
                name="Go Constants Sync (gitmap/constants/constants.go)",
                check_type="constants_go",
                target_path=constants_go,
                canonical_version=canonical,
            )
        )

    if args.filter:
        filt = args.filter.lower()
        tasks = [t for t in tasks if filt in t.name.lower()]

    exit_code = run_worker_pool(
        items=tasks,
        worker_fn=execute_version_check,
        max_workers=args.workers,
        is_sync=args.is_sync,
        show_all=args.show_all,
        output_file=args.output_file,
        as_json=args.as_json,
        title=f"VERSION SYNCHRONIZATION GUARD (v{canonical})",
        item_noun="version sync check(s)",
    )
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
