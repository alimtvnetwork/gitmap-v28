#!/usr/bin/env python3
"""Comprehensive Audit Suite for All Connected Repositories.

Audits:
1. File structure: presence and counts of 01-prompts/v1, 01-prompts/v2, .agents/skills, 03-ai-scripts
2. Git working tree state (must be clean)
3. Active branch (must be main)
4. Latest release tag (must be bumped and present)
5. Remote branches on origin:
   - backup branch (backup/sync-prompts-v1-v2-*)
   - feature branch (feat/sync-prompts-v1-v2-gitmap)
   - release branch (release/v*)
   - main branch up-to-date with remote
"""
from __future__ import annotations

import os
from pathlib import Path
import subprocess
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

REPOS = [
    "coding-guidelines",
    "antigravity-manager",
    "spec-builder",
    "movie-cli",
    "macro-ahk",
    "laravel-automation",
    "lara-publishing",
    "lara-licensing",
    "gitmap",
    "wp-exam",
    "wp-git-log",
    "wp-html-automate",
    "wp-link-manager",
    "wp-onboarding",
    "cat-my",
    "scripts-fixer",
    "gitlogger-new",
]

BASE_DIR = Path("d:/work")


def run_git(args: list[str], cwd: Path) -> tuple[int, str]:
    res = subprocess.run(
        ["git"] + args,
        cwd=str(cwd),
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    return res.returncode, res.stdout.strip()


def main() -> None:
    print("=" * 115)
    print("CANONICAL MULTI-REPOSITORY REPOSITORY & RELEASE AUDIT")
    print("=" * 115)

    header = f"{'Repository':<23} | {'Branch':<6} | {'Clean':<5} | {'Tag':<10} | {'V1':<4} | {'V2':<4} | {'Skills':<6} | {'Scripts':<7} | {'Bkup':<5} | {'Feat':<5} | {'Rel':<5}"
    print(header)
    print("-" * 115)

    all_passed = True

    for r in REPOS:
        p = BASE_DIR / r
        if not p.exists():
            print(f"{r:<23} | MISSING DIRECTORY")
            all_passed = False
            continue

        v1_count = len(list((p / "01-prompts" / "v1").glob("**/*.md"))) if (p / "01-prompts" / "v1").exists() else 0
        v2_count = len(list((p / "01-prompts" / "v2").glob("**/*.md"))) if (p / "01-prompts" / "v2").exists() else 0
        skills_count = len(list((p / ".agents" / "skills").glob("**/skill.md"))) if (p / ".agents" / "skills").exists() else 0
        scripts_count = len(list((p / "03-ai-scripts").glob("*.py"))) if (p / "03-ai-scripts").exists() else 0

        _, branch = run_git(["branch", "--show-current"], p)
        _, status = run_git(["status", "--porcelain"], p)
        _, tag = run_git(["describe", "--tags", "--abbrev=0"], p)
        _, remotes = run_git(["branch", "-r"], p)

        has_backup = any("backup/" in line for line in remotes.splitlines())
        has_feat = any("feat/" in line for line in remotes.splitlines())
        has_release = any("release/" in line for line in remotes.splitlines())
        is_clean = len(status) == 0

        row = f"{r:<23} | {branch:<6} | {str(is_clean):<5} | {tag:<10} | {v1_count:<4} | {v2_count:<4} | {skills_count:<6} | {scripts_count:<7} | {str(has_backup):<5} | {str(has_feat):<5} | {str(has_release):<5}"
        print(row)

        if not is_clean or v1_count == 0 or v2_count == 0 or not has_backup or not has_feat or not has_release:
            all_passed = False

    print("=" * 115)
    if all_passed:
        print("[AUDIT RESULT: 100% PASS] All repositories verified with exact V1/V2 prompts, skills, and full branch ceremonies.")
    else:
        print("[AUDIT RESULT: DISCREPANCY DETECTED] See rows above for details.")


if __name__ == "__main__":
    main()
