#!/usr/bin/env python3
"""Multi-Repository V3 Prompts and Skills Synchronizer.

Synchronizes canonical prompts (01-prompts/, including v3/) and
Antigravity skills (.agents/skills/) from coding-guidelines-v24 across
the 43 target repositories requested by the user.

Workflow per repository:
1. Stash any uncommitted working tree changes safely.
2. Checkout primary/base branch (main or develop).
3. Pull latest changes from origin.
4. Create and push timestamped backup branch.
5. Mirror canonical 01-prompts/ and .agents/skills/ without deleting repo-specific skills.
6. Commit changes to base branch and push to origin.
7. Restore any stashed working tree changes.
"""
from __future__ import annotations

import argparse
from datetime import datetime
import os
from pathlib import Path
import shutil
import subprocess
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

SOURCE_ROOT = Path(__file__).resolve().parent.parent
SOURCE_PROMPTS = SOURCE_ROOT / "01-prompts"
SOURCE_SKILLS = SOURCE_ROOT / ".agents" / "skills"

TARGET_REPOS = [
    ("ai-empathy-prompt-tuner-v1", Path(r"D:\work\02-prompts\ai-empathy-prompt-tuner")),
    ("alim-cv-v8", Path(r"D:\work\alim-cv")),
    ("alim-karim-profile-v2", Path(r"D:\work\alim-karim-profile")),
    ("alim.karim.profile", Path(r"D:\work\aukgit\alim.karim.profile")),
    ("Antigravity-Manager", Path(r"D:\work\antigravity-manager")),
    ("cat-my-v12", Path(r"D:\work\cat-my")),
    ("core-v9", Path(r"D:\work\03-aukgo\core")),
    ("digital-name-card", Path(r"D:\work\digital-name-card")),
    ("flat-slide-show", Path(r"D:\work\presentations-repos\flat-slide-show")),
    ("gitlogger-new-v2", Path(r"D:\work\gitlogger-new")),
    ("gitmap-v28", Path(r"D:\work\gitmap")),
    ("global-ppt-v1", Path(r"D:\work\presentations-repos\global-ppt-v1")),
    ("hiltrax-v1", Path(r"D:\work\presentations-repos\hiltrax")),
    ("icon-coding-guidelines", Path(r"D:\work\icon-coding-guidelines")),
    ("img-pdf-v2", Path(r"D:\work\img-pdf")),
    ("ki-health-ppt-v5", Path(r"D:\work\presentations-repos\ki-health-ppt")),
    ("kubernetes-training-v1", Path(r"D:\work\aukgit\kubernetes-training")),
    ("lara-licensing-v4", Path(r"D:\work\lara-licensing")),
    ("lara-publishing-v1", Path(r"D:\work\lara-publishing")),
    ("laravel-automation-v1", Path(r"D:\work\laravel-automation")),
    ("letsmarknow-ui-v2", Path(r"D:\work\letsmarknow-ui")),
    ("letsmarknow-v2", Path(r"D:\work\letsmarknow")),
    ("macro-ahk-v55", Path(r"D:\work\macro-ahk")),
    ("maid-app-spec-presentation-v1", Path(r"D:\work\presentations-repos\maid-app-spec-presentation")),
    ("movie-cli-v8", Path(r"D:\work\movie-cli")),
    ("pathhelper", Path(r"D:\work\03-aukgo\pathhelper")),
    ("presentation-aug-2026-plans-alim", Path(r"D:\work\presentations-repos\presentation-aug-2026-plans-alim")),
    ("prompts-connect-v3", Path(r"D:\work\02-prompts\prompts-connect")),
    ("punam-case-studies-v1", Path(r"D:\work\punam-case-studies-v1")),
    ("rasia-logo", Path(r"D:\work\presentations-repos\rasia-logo")),
    ("scripts-fixer-v20", Path(r"D:\work\scripts-fixer")),
    ("slides-spec", Path(r"D:\work\presentations-repos\slides-spec")),
    ("spec-builder-v11", Path(r"D:\work\spec-builder")),
    ("sweet-digs-finder", Path(r"D:\work\web-system\sweet-digs-finder")),
    ("ui-prompts-cat", Path(r"D:\work\ui-prompts-cat")),
    ("white-presentation-v1", Path(r"D:\work\presentations-repos\white-presentation-v1")),
    ("workflowy-ui-v1", Path(r"D:\work\workflowy-ui")),
    ("workflowy-v3", Path(r"D:\work\workflowy")),
    ("wp-exam-v2", Path(r"D:\work\wp-exam")),
    ("wp-git-log-v3", Path(r"D:\work\wp-git-log")),
    ("wp-html-automate", Path(r"D:\work\wp-html-automate")),
    ("wp-link-manager-v5", Path(r"D:\work\wp-link-manager")),
    ("wp-onboarding-v17", Path(r"D:\work\wp-onboarding")),
]

EXCLUDE_NAMES = {
    "__pycache__",
    ".git",
    ".pytest_cache",
    ".mypy_cache",
    ".DS_Store",
}

EXCLUDE_EXTS = {
    ".pyc",
    ".pyo",
    ".tmp",
}


def run_cmd(cmd: str | list[str], cwd: Path) -> tuple[int, str, str]:
    """Execute command with utf-8 encoding."""
    is_shell = isinstance(cmd, str)
    res = subprocess.run(
        cmd,
        cwd=str(cwd),
        shell=is_shell,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    return res.returncode, res.stdout.strip(), res.stderr.strip()


def detect_base_branch(target: Path) -> str:
    """Detect default branch (develop, main, or master)."""
    code, out, _ = run_cmd("git branch --list develop main master", target)
    lines = [line.strip().replace("*", "").strip() for line in out.splitlines()]
    if target.name == "pathhelper" and "develop" in lines:
        return "develop"
    if "main" in lines:
        return "main"
    if "master" in lines:
        return "master"
    code, out, _ = run_cmd("git branch --show-current", target)
    return out or "main"


def copy_file_if_different(src: Path, dst: Path) -> bool:
    """Copy src to dst if dst does not exist or has different content."""
    dst.parent.mkdir(parents=True, exist_ok=True)
    if dst.exists():
        try:
            if src.stat().st_size == dst.stat().st_size and src.read_bytes() == dst.read_bytes():
                return False
        except Exception:
            pass
    shutil.copy2(src, dst)
    return True


def sync_prompts_and_skills(target: Path) -> tuple[int, int]:
    """Mirror canonical 01-prompts and .agents/skills into target repo.

    Preserves target repo unique skills and unique prompt folders.
    """
    copied = 0
    removed = 0

    # 1. Sync 01-prompts
    target_prompts = target / "01-prompts"
    target_prompts.mkdir(parents=True, exist_ok=True)

    for root, dirs, files in os.walk(SOURCE_PROMPTS):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_NAMES]
        rel_root = Path(root).relative_to(SOURCE_PROMPTS)
        target_dir = target_prompts / rel_root

        for f in files:
            if f in EXCLUDE_NAMES or Path(f).suffix.lower() in EXCLUDE_EXTS:
                continue
            src_file = Path(root) / f
            dst_file = target_dir / f
            if copy_file_if_different(src_file, dst_file):
                copied += 1

    # 2. Sync .agents/skills (canonical skills only, preserve repo-specific)
    target_skills = target / ".agents" / "skills"
    target_skills.mkdir(parents=True, exist_ok=True)

    for skill_dir in SOURCE_SKILLS.iterdir():
        if not skill_dir.is_dir() or skill_dir.name in EXCLUDE_NAMES:
            continue
        dst_skill_dir = target_skills / skill_dir.name
        dst_skill_dir.mkdir(parents=True, exist_ok=True)

        for root, dirs, files in os.walk(skill_dir):
            dirs[:] = [d for d in dirs if d not in EXCLUDE_NAMES]
            rel_root = Path(root).relative_to(skill_dir)
            target_sub = dst_skill_dir / rel_root

            for f in files:
                if f in EXCLUDE_NAMES or Path(f).suffix.lower() in EXCLUDE_EXTS:
                    continue
                src_file = Path(root) / f
                dst_file = target_sub / f
                if copy_file_if_different(src_file, dst_file):
                    copied += 1

    return copied, removed


def sync_repository(name: str, target: Path, dry_run: bool = False, no_push: bool = False) -> dict[str, any]:
    print(f"\n=======================================================")
    print(f"[{name}] Syncing: {target}")
    print(f"=======================================================")

    result = {
        "repo": name,
        "path": str(target),
        "base_branch": "",
        "backup_branch": "",
        "stashed": False,
        "copied": 0,
        "changed": False,
        "committed": False,
        "pushed": False,
        "error": None,
    }

    if not target.exists() or not (target / ".git").exists():
        result["error"] = "Not a valid git repository"
        print(f"ERROR: {result['error']}")
        return result

    base_branch = detect_base_branch(target)
    result["base_branch"] = base_branch
    print(f"1. Detected base branch: {base_branch}")

    # Check dirty status
    code, st_out, _ = run_cmd("git status --porcelain", target)
    is_dirty = bool(st_out.strip())
    if is_dirty:
        print("2. Repository has uncommitted changes. Stashing safely...")
        code, stash_out, err = run_cmd("git stash --include-untracked", target)
        if code == 0 and "No local changes to save" not in stash_out:
            result["stashed"] = True
            print("   Working tree stashed.")
    else:
        print("2. Working tree is clean.")

    try:
        # Checkout base branch
        print(f"3. Checking out {base_branch}...")
        code, _, err = run_cmd(f"git checkout {base_branch}", target)
        if code != 0:
            result["error"] = f"Failed to checkout {base_branch}: {err}"
            print(f"ERROR: {result['error']}")
            return result

        # Pull latest changes
        print(f"4. Pulling latest from origin/{base_branch}...")
        code, pull_out, pull_err = run_cmd(f"git pull origin {base_branch}", target)
        if code != 0:
            print(f"   Pull warning/failure: {pull_err or pull_out}")
            # If fast-forward fails or divergent, log warning but continue
        else:
            print("   Pull successful.")

        # Create and push backup branch
        timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
        backup_branch = f"backup/sync-v3-prompts-skills-{timestamp}"
        result["backup_branch"] = backup_branch
        print(f"5. Creating backup branch: {backup_branch}...")
        code, _, err = run_cmd(f"git branch {backup_branch}", target)
        if code != 0:
            result["error"] = f"Failed to create backup branch: {err}"
            print(f"ERROR: {result['error']}")
            return result

        if not no_push and not dry_run:
            print(f"   Pushing {backup_branch} to origin...")
            code, _, err = run_cmd(f"git push origin {backup_branch}", target)
            if code != 0:
                print(f"   Backup push warning: {err}")
            else:
                print("   Backup branch pushed successfully.")

        if dry_run:
            print("6. Dry run enabled. Skipping file copy and commit.")
            return result

        # Copy prompts and skills
        print("6. Mirroring v3 prompts and canonical skills...")
        copied, removed = sync_prompts_and_skills(target)
        result["copied"] = copied
        print(f"   Updated {copied} files.")

        # Check git status for changes
        code, status_out, _ = run_cmd("git status --porcelain 01-prompts .agents/skills", target)
        if not status_out.strip():
            print("7. No file changes detected. Repository already in sync.")
            result["changed"] = False
            return result

        result["changed"] = True
        print(f"7. Changes detected ({len(status_out.splitlines())} modified/new files).")

        # Stage and commit
        print("8. Staging and committing changes...")
        run_cmd("git add 01-prompts .agents/skills", target)
        commit_msg = "feat(prompts,skills): update v3 prompts and sync canonical skills"
        code, commit_out, commit_err = run_cmd(f'git commit -m "{commit_msg}"', target)
        if code != 0:
            result["error"] = f"Commit failed: {commit_err or commit_out}"
            print(f"ERROR: {result['error']}")
            return result

        result["committed"] = True
        print("   Commit successful.")

        # Push to origin
        if not no_push:
            print(f"9. Pushing {base_branch} to origin...")
            code, push_out, push_err = run_cmd(f"git push origin {base_branch}", target)
            if code != 0:
                result["error"] = f"Push failed: {push_err or push_out}"
                print(f"ERROR: {result['error']}")
                return result
            result["pushed"] = True
            print("   Push successful.")

    finally:
        # Restore stash if needed
        if result["stashed"]:
            print("10. Restoring stashed changes...")
            code, pop_out, pop_err = run_cmd("git stash pop", target)
            if code != 0:
                print(f"    Stash pop warning: {pop_err or pop_out}")
            else:
                print("    Stashed changes restored successfully.")

    print(f"[OK] Completed {name}")
    return result


def main() -> None:
    parser = argparse.ArgumentParser(description="Multi-Repository V3 Prompts and Skills Synchronizer")
    parser.add_argument("--repo", "-r", type=str, help="Specific repository name or slug to sync")
    parser.add_argument("--dry-run", action="store_true", help="Preview without making changes")
    parser.add_argument("--no-push", action="store_true", help="Commit locally without pushing")
    args = parser.parse_args()

    if args.repo:
        matching = [(name, p) for name, p in TARGET_REPOS if name.lower() == args.repo.lower() or p.name.lower() == args.repo.lower()]
        if not matching:
            print(f"Error: Repository '{args.repo}' not found in target list.")
            sys.exit(1)
        targets = matching
    else:
        targets = TARGET_REPOS

    print(f"Source repository : {SOURCE_ROOT.name} ({SOURCE_ROOT})")
    print(f"Target count      : {len(targets)}")
    print(f"Dry run           : {args.dry_run}")
    print(f"No push           : {args.no_push}")

    summary = []
    for name, path in targets:
        res = sync_repository(name, path, dry_run=args.dry_run, no_push=args.no_push)
        summary.append(res)

    print("\n=======================================================")
    print("V3 PROMPTS & SKILLS SYNCHRONIZATION SUMMARY")
    print("=======================================================")
    print(f"{'Repository':<32} | {'Branch':<8} | {'Backup Branch':<38} | {'Status':<15}")
    print("-" * 100)
    for s in summary:
        status = "OK (Pushed)" if s["pushed"] else ("Clean (Sync)" if not s["changed"] and not s["error"] else "FAIL")
        if s["error"]:
            status = f"FAIL: {s['error'][:20]}"
        b_branch = s.get("backup_branch") or "N/A"
        print(f"{s['repo']:<32} | {s['base_branch']:<8} | {b_branch:<38} | {status:<15}")


if __name__ == "__main__":
    main()
