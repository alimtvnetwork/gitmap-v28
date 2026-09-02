#!/usr/bin/env python3
"""linter-scripts/run.py — Pull latest changes, run guidelines validator, run all linters (cross-platform).

Usage:
  python linter-scripts/run.py                       # full pipeline
  python linter-scripts/run.py -d                    # git pull only
  python linter-scripts/run.py --path cmd --max-lines 20
  python linter-scripts/run.py --json
  python linter-scripts/run.py --skip-linters        # skip Step 3
  python linter-scripts/run.py --linters-only        # only Step 3
"""

from __future__ import annotations

import argparse
import shutil
import subprocess
import sys
from pathlib import Path

LINTERS = [
    ("tunable-constants", "check-tunable-constants.py"),
    ("mws-error-codes", "check-mws-error-codes.py"),
    ("function-lengths", "check-function-lengths.py"),
    ("forbidden-strings", "check-forbidden-strings.py"),
    ("placeholder-comments", "check-placeholder-comments.py"),
    ("memory-mirror-drift", "check-memory-mirror-drift.py"),
    ("prompts-loaded", "check-prompts-loaded.py"),
    ("readme-canonicals", "check-readme-canonicals.py"),
    ("readme-install", "check-readme-install-section.py"),
    ("root-readme", "check-root-readme.py"),
    ("spec-cross-links", "check-spec-cross-links.py"),
    ("spec-folder-refs", "check-spec-folder-refs.py"),
    ("axios-version", "check-axios-version.py"),
    ("forbidden-spec-paths", "check-forbidden-spec-paths.py"),
    ("runner-dispatch", "check-runner-dispatch-antipatterns.py"),
]


def main() -> int:
    parser = argparse.ArgumentParser(description="Lint runner pipeline.")
    parser.add_argument("-d", action="store_true", help="Skip validation, git pull only")
    parser.add_argument("--path", default="src", help="Path to scan for Go validator")
    parser.add_argument("--max-lines", type=int, default=15, help="Max lines per function")
    parser.add_argument("--json", action="store_true", help="JSON output for validator")
    parser.add_argument("--skip-linters", action="store_true", help="Skip linters step")
    parser.add_argument("--linters-only", action="store_true", help="Run only linters")
    args = parser.parse_args()

    script_dir = Path(__file__).resolve().parent
    repo_root = script_dir.parent

    # Step 1: Git Pull
    if not args.linters_only:
        print("\n═══ Step 1 — git pull ═══")
        res = subprocess.run(["git", "pull"], cwd=str(repo_root))
        if res.returncode == 0:
            print("✅ Repository up to date.")
        else:
            print("⚠️  git pull failed — continuing with local files...")

    if args.d and not args.linters_only:
        print("\n⏭️  Skipping validation (-d flag).")
        return 0

    validator_exit = 0

    # Step 2: Go Validator
    if not args.linters_only:
        print("\n═══ Step 2 — Coding-guidelines validator ═══")
        go_file = script_dir / "validate-guidelines.go"
        go_exe = shutil.which("go")
        if not go_file.is_file():
            print(f"❌ Cannot find {go_file}", file=sys.stderr)
            validator_exit = 1
        elif not go_exe:
            print("❌ Go is not installed or not in PATH.", file=sys.stderr)
            validator_exit = 1
        else:
            cmd = [go_exe, "run", str(go_file), "--path", args.path, "--max-lines", str(args.max_lines)]
            if args.json:
                cmd.append("--json")
            res = subprocess.run(cmd, cwd=str(script_dir))
            validator_exit = res.returncode
            if validator_exit == 0:
                print("✅ Step 2 passed.")
            else:
                print(f"❌ Step 2 failed with CODE RED violations (exit={validator_exit}).", file=sys.stderr)

    # Step 3: Python linters
    linters_exit = 0
    passed: list[str] = []
    failed: list[str] = []

    if not args.skip_linters:
        print("\n═══ Step 3 — Spec / docs linters ═══")
        for label, script_name in LINTERS:
            target_py = script_dir / script_name
            if not target_py.is_file():
                continue
            print(f"\n── {label} ──")
            res = subprocess.run([sys.executable, str(target_py)], cwd=str(repo_root))
            if res.returncode == 0:
                passed.append(label)
            else:
                failed.append(f"{label} (exit={res.returncode})")
                linters_exit = 1

    print("\n═══ Summary ═══")
    print(f"Step 2 (validator): exit={validator_exit}")
    if not args.skip_linters:
        print(f"Step 3 (linters): {len(passed)} passed, {len(failed)} failed")
        for item in failed:
            print(f"  ❌ {item}", file=sys.stderr)

    if validator_exit != 0 or linters_exit != 0:
        return 1

    print("✅ All checks passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
