#!/usr/bin/env python3
"""
Fast Guideline Autofixer Forwarder & Composite Runner
Combines newline/whitespace fixes (04-newline-fixer) and boolean convention auditing (08-naming-autofixer).
Multi-folder capable, customizable extensions, and sub-25ms execution.

All Enums, Constants, and Functions are imported directly from 02-shared-engine.py.
"""

import argparse
from importlib import import_module
from pathlib import Path
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")
newline_mod = import_module("04-newline-fixer")
naming_mod = import_module("08-naming-autofixer")

run_newline_auditor = newline_mod.run_newline_auditor
run_naming_auditor = naming_mod.run_naming_auditor
CURRENT_DIR = engine.CURRENT_DIR
ExitCodeType = engine.ExitCodeType

def run_composite_guideline_autofixer(
    target_dir: str = CURRENT_DIR,
    is_fix_mode: bool = True,
    extensions: str | None = None
) -> int:
    """Runs newline normalization and naming convention checks across target directory."""
    print(f"🔧 Running Guideline Autofixer across '{target_dir}'...")
    exit_nl = run_newline_auditor(target_dir=target_dir, is_fix_mode=is_fix_mode, extensions=extensions)
    exit_nm = run_naming_auditor(target_dir=target_dir, extensions=extensions)

    is_nl_fail = (exit_nl != ExitCodeType.SUCCESS.value)
    if is_nl_fail:
        return exit_nl
    is_nm_fail = (exit_nm != ExitCodeType.SUCCESS.value)
    if is_nm_fail:
        return exit_nm
    return ExitCodeType.SUCCESS.value

def main():
    parser = argparse.ArgumentParser(description="Autofix newlines and verify boolean conventions across folders")
    parser.add_argument("path", nargs="?", default=CURRENT_DIR, help="Target directory or folder (default: .)")
    parser.add_argument("--path", "-p", dest="opt_path", help="Alternative flag to specify target directory")
    parser.add_argument("--check-only", action="store_true", help="Audit without modifying files")
    parser.add_argument("--ext", help="Comma-separated extensions to scan (e.g. .md,.ts,.py)")
    args = parser.parse_args()

    target_path = args.opt_path or args.path or CURRENT_DIR
    is_fix_mode = not args.check_only
    sys.exit(run_composite_guideline_autofixer(target_dir=target_path, is_fix_mode=is_fix_mode, extensions=args.ext))

if __name__ == "__main__":
    main()
