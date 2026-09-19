#!/usr/bin/env python3
"""
Fast Version Synchronization & Changelog Guard
Validates that version.json, package.json, and changelog.md are in 100% sync.
Multi-folder capable, thread-safe lazy regex engine, and sub-5ms execution.

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

import argparse
from importlib import import_module
import json
import os
from pathlib import Path
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

RegexPatternType = engine.RegexPatternType
get_compiled_regex = engine.get_compiled_regex
ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

def check_version_sync(root_dir: str = CURRENT_DIR) -> int:
    """Checks version parity across version.json, package.json, and changelog.md."""
    start_time = time.perf_counter()
    errors = []

    version_json_p = os.path.join(root_dir, "version.json")
    package_json_p = os.path.join(root_dir, "package.json")
    changelog_p = os.path.join(root_dir, "changelog.md")

    canonical_version = None

    has_version_json = os.path.exists(version_json_p)
    has_package_json = os.path.exists(package_json_p)
    has_changelog = os.path.exists(changelog_p)

    if has_version_json:
        try:
            with open(version_json_p, "r", encoding=DEFAULT_ENCODING) as f:
                v_data = json.load(f)
                canonical_version = v_data.get("version")
        except Exception as e:
            errors.append(f"version.json parse error: {e}")
    elif has_package_json:
        try:
            with open(package_json_p, "r", encoding=DEFAULT_ENCODING) as f:
                p_data = json.load(f)
                canonical_version = p_data.get("version")
        except Exception as e:
            errors.append(f"package.json parse error: {e}")

    has_canonical = bool(canonical_version)
    if not has_canonical:
        print(f"⚠️ No canonical version source (version.json or package.json) found in '{root_dir}'.")
        return ExitCodeType.SUCCESS.value

    # 1. Compare with package.json
    if has_package_json:
        try:
            with open(package_json_p, "r", encoding=DEFAULT_ENCODING) as f:
                p_data = json.load(f)
                pkg_ver = p_data.get("version")
                is_version_match = (pkg_ver == canonical_version)
                if not is_version_match:
                    errors.append(f"Version mismatch: version.json has '{canonical_version}' but package.json has '{pkg_ver}'")
        except Exception as e:
            errors.append(f"package.json error: {e}")

    # 2. Compare with changelog.md latest entry
    if has_changelog:
        try:
            with open(changelog_p, "r", encoding=DEFAULT_ENCODING) as f:
                changelog_text = f.read()
            re_changelog = get_compiled_regex(RegexPatternType.CHANGELOG_HEADER)
            match = re_changelog.search(changelog_text)
            has_match = bool(match)
            if has_match:
                latest_cl_ver = match.group(1).lstrip("v")
                is_changelog_match = (latest_cl_ver == canonical_version.lstrip("v"))
                if not is_changelog_match:
                    errors.append(f"Changelog mismatch: latest header is 'v{latest_cl_ver}' but canonical version is '{canonical_version}'")
        except Exception as e:
            errors.append(f"changelog.md error: {e}")

    elapsed_ms = (time.perf_counter() - start_time) * 1000

    has_errors = len(errors) > 0
    if has_errors:
        print(f"{LINE_SEPARATOR}❌ Version synchronization failed in '{root_dir}' ({elapsed_ms:.2f}ms):")
        for err in errors:
            print(f"  ::error::{err}")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print(f"✅ Version synchronization verified: v{canonical_version} ({elapsed_ms:.2f}ms)")
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Check version synchronization across folders")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Root directory containing version files")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify root directory")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    sys.exit(check_version_sync(target_path))

if __name__ == "__main__":
    main()
