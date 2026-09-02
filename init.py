#!/usr/bin/env python3
"""init.sh / init.py — one-shot repo init: ensure repo is public, then rewrite stale tokens via fix-repo.

Spec: spec/03-general/11-init-pipeline.md

Usage:
  python init.py             # ensure public, then rewrite stale version tokens
  python init.py --dry-run   # preview both steps
"""

from __future__ import annotations

import argparse
import os
import shutil
import subprocess
import sys
from pathlib import Path


def run_step(label: str, script_name: str, subcmd: str, script_args: list[str], gitmap_args: list[str]) -> int:
    repo_root = Path(__file__).resolve().parent
    py_script = repo_root / script_name

    print()
    if py_script.is_file():
        cmd = [sys.executable, str(py_script)] + script_args
        print(f"==> [{label}] {py_script.name} {' '.join(script_args)}")
        res = subprocess.run(cmd)
        return res.returncode

    gitmap_bin = shutil.which("gitmap")
    if gitmap_bin:
        cmd = [gitmap_bin, subcmd] + gitmap_args
        print(f"==> [{label}] gitmap {subcmd} {' '.join(gitmap_args)}")
        res = subprocess.run(cmd)
        return res.returncode

    print(f"==> [{label}] SKIP — neither {script_name} nor 'gitmap' binary found on PATH", file=sys.stderr)
    return 127


def main() -> int:
    parser = argparse.ArgumentParser(description="One-shot repo init: ensure repo is public, then rewrite stale tokens.")
    parser.add_argument("--dry-run", action="store_true", help="Preview both steps without modifying")
    args = parser.parse_args()

    vis_args = ["--visible", "pub", "--yes"]
    vis_gargs = ["--yes"]
    if args.dry_run:
        vis_args.append("--dry-run")
        vis_gargs.append("--dry-run")

    vis_rc = run_step("visibility", "visibility_change.py", "make-public", vis_args, vis_gargs)

    fix_args = ["--all"]
    fix_gargs = ["--all"]
    if args.dry_run:
        fix_args.append("--dry-run")
        fix_gargs.append("--dry-run")

    fix_rc = run_step("fix-repo", "fix_repo.py", "fix-repo", fix_args, fix_gargs)

    print("\n==> init summary")
    print(f"    visibility-change : exit {vis_rc}")
    print(f"    fix-repo          : exit {fix_rc}")

    if vis_rc != 0:
        return vis_rc
    if fix_rc != 0:
        return fix_rc
    return 0


if __name__ == "__main__":
    sys.exit(main())
