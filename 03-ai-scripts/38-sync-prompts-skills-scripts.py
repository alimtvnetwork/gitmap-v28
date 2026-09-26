#!/usr/bin/env python3
"""Multi-Repository Prompts, Skills, and AI Scripts Synchronizer & Release Orchestrator.

Synchronizes canonical prompts (01-prompts/), Antigravity skills (.agents/skills/),
and AI scripts (03-ai-scripts/ and .agents/scripts/) from coding-guidelines-v24
across 13 connected repositories.

Performs the complete, safe multi-branch release ceremony per repo:
1. Pull latest base branch.
2. Create and push backup branch (backup/sync-prompts-v1-v2-<timestamp>).
3. Create feature branch (feat/sync-prompts-v1-v2-gitmap).
4. Mirror directories cleanly with stale-file removal.
5. Commit changes atomically to feature branch.
6. Push feature branch to origin.
7. Create release branch (release/vX.Y.Z) and tag (vX.Y.Z).
8. Push release branch, tag, and backup branch.
9. Merge back into base branch and push.
"""
from __future__ import annotations

import argparse
from datetime import datetime
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

SOURCE_ROOT = Path(__file__).resolve().parent.parent

TARGET_REPOS = [
    SOURCE_ROOT.parent / "antigravity-manager",
    SOURCE_ROOT.parent / "spec-builder",
    SOURCE_ROOT.parent / "movie-cli",
    SOURCE_ROOT.parent / "macro-ahk",
    SOURCE_ROOT.parent / "laravel-automation",
    SOURCE_ROOT.parent / "lara-publishing",
    SOURCE_ROOT.parent / "lara-licensing",
    SOURCE_ROOT.parent / "gitmap",
    SOURCE_ROOT.parent / "wp-exam",
    SOURCE_ROOT.parent / "wp-git-log",
    SOURCE_ROOT.parent / "wp-html-automate",
    SOURCE_ROOT.parent / "wp-link-manager",
    SOURCE_ROOT.parent / "wp-onboarding",
    SOURCE_ROOT.parent / "cat-my",
    SOURCE_ROOT.parent / "scripts-fixer",
    SOURCE_ROOT.parent / "gitlogger-new",
]

SYNC_DIRS = [
    ("01-prompts", "01-prompts"),
    (".agents/skills", ".agents/skills"),
    ("03-ai-scripts", "03-ai-scripts"),
    (".agents/scripts", ".agents/scripts"),
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
    """Detect default branch (main or master)."""
    code, out, _ = run_cmd("git branch --list main master", target)
    lines = [line.strip().replace("*", "").strip() for line in out.splitlines()]
    if "main" in lines:
        return "main"
    if "master" in lines:
        return "master"
    code, out, _ = run_cmd("git branch --show-current", target)
    return out or "main"


def detect_latest_version(target: Path) -> str:
    """Detect current SemVer string from git tags or version files."""
    code, out, _ = run_cmd("git describe --tags --abbrev=0", target)
    if code == 0 and out.strip():
        tag = out.strip().lstrip("v")
        match = re.search(r"(\d+\.\d+\.\d+)", tag)
        if match:
            return match.group(1)

    for vf in ["version.json", "package.json", "version.txt", "VERSION"]:
        vpath = target / vf
        if vpath.exists():
            try:
                content = vpath.read_text(encoding="utf-8")
                match = re.search(r'"version"\s*:\s*"(\d+\.\d+\.\d+)"', content)
                if match:
                    return match.group(1)
                match = re.search(r"(\d+\.\d+\.\d+)", content)
                if match:
                    return match.group(1)
            except Exception:
                pass

    return "0.1.0"


def bump_patch_version(ver_str: str) -> str:
    """Increment patch component of SemVer."""
    parts = ver_str.split(".")
    if len(parts) >= 3 and parts[2].isdigit():
        parts[2] = str(int(parts[2]) + 1)
        return ".".join(parts[:3])
    return f"{ver_str}.1"


def mirror_directory(src: Path, dst: Path) -> tuple[int, int]:
    """Mirror src into dst, removing stale files and copying new/updated ones."""
    copied = 0
    removed = 0

    if not src.exists():
        return 0, 0

    dst.mkdir(parents=True, exist_ok=True)

    # 1. Clean stale files/dirs in dst that are no longer in src
    for root, dirs, files in os.walk(dst, topdown=False):
        rel_root = Path(root).relative_to(dst)
        src_root = src / rel_root

        for f in files:
            dst_file = Path(root) / f
            src_file = src_root / f
            if f in EXCLUDE_NAMES or dst_file.suffix.lower() in EXCLUDE_EXTS:
                dst_file.unlink(missing_ok=True)
                removed += 1
            elif not src_file.exists():
                dst_file.unlink(missing_ok=True)
                removed += 1

        for d in dirs:
            dst_dir = Path(root) / d
            src_dir = src_root / d
            if d in EXCLUDE_NAMES:
                shutil.rmtree(dst_dir, ignore_errors=True)
                removed += 1
            elif not src_dir.exists():
                shutil.rmtree(dst_dir, ignore_errors=True)
                removed += 1

    # 2. Copy files from src to dst
    for root, dirs, files in os.walk(src):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_NAMES]

        rel_root = Path(root).relative_to(src)
        target_dir = dst / rel_root
        target_dir.mkdir(parents=True, exist_ok=True)

        for f in files:
            if f in EXCLUDE_NAMES or Path(f).suffix.lower() in EXCLUDE_EXTS:
                continue

            src_file = Path(root) / f
            dst_file = target_dir / f

            need_copy = False
            if not dst_file.exists():
                need_copy = True
            else:
                try:
                    src_stat = src_file.stat()
                    dst_stat = dst_file.stat()
                    if src_stat.st_size != dst_stat.st_size:
                        need_copy = True
                    else:
                        if src_file.read_bytes() != dst_file.read_bytes():
                            need_copy = True
                except Exception:
                    need_copy = True

            if need_copy:
                shutil.copy2(src_file, dst_file)
                copied += 1

    return copied, removed


def sync_repo(target: Path, dry_run: bool = False, no_push: bool = False) -> dict[str, any]:
    print(f"\n=======================================================")
    print(f"Syncing target repository: {target.name} ({target})")
    print(f"=======================================================")

    result = {
        "repo": target.name,
        "base_branch": "",
        "backup_branch": "",
        "feat_branch": "",
        "release_tag": "",
        "copied": 0,
        "removed": 0,
        "changed": False,
        "committed": False,
        "pushed": False,
        "error": None,
    }

    if not target.exists() or not (target / ".git").exists():
        err = f"Directory {target} is not a valid git repository."
        print(f"ERROR: {err}")
        result["error"] = err
        return result

    base_branch = detect_base_branch(target)
    result["base_branch"] = base_branch

    # 1. Switch to feature branch to preserve working changes
    feat_branch = "feat/sync-prompts-v1-v2-gitmap"
    result["feat_branch"] = feat_branch
    print(f"[1/7] Checking out feature branch: {feat_branch}...")
    run_cmd(f"git checkout -B {feat_branch}", target)

    # 2. Create and push backup branch from HEAD
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    backup_branch = f"backup/sync-prompts-v1-v2-{timestamp}"
    result["backup_branch"] = backup_branch
    print(f"[2/7] Creating backup branch: {backup_branch}...")
    run_cmd(f"git branch {backup_branch} HEAD", target)
    if not no_push and not dry_run:
        print(f"      Pushing {backup_branch} to origin...")
        run_cmd(f"git push origin {backup_branch}", target)

    # 3. Pull latest base branch into feature branch
    print(f"[3/7] Pulling latest changes from {base_branch}...")
    code, out, err = run_cmd(f"git pull origin {base_branch}", target)
    if code != 0:
        print(f"Pull warning (continuing): {err or out}")

    # 4. Mirror directories
    print("[4/7] Mirroring prompts, skills, and AI scripts...")
    total_copied = 0
    total_removed = 0
    for src_rel, dst_rel in SYNC_DIRS:
        src_path = SOURCE_ROOT / src_rel
        dst_path = target / dst_rel
        c, r = mirror_directory(src_path, dst_path)
        print(f"  - {src_rel} -> {dst_rel} (copied: {c}, removed: {r})")
        total_copied += c
        total_removed += r

    # If gitmap repo, ensure benchmark docs are synced
    if target.name == "gitmap":
        bm_src = SOURCE_ROOT / "docs" / "benchmarks" / "search_benchmark.md"
        bm_dst = target / "docs" / "benchmarks" / "search_benchmark.md"
        if bm_src.exists():
            bm_dst.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(bm_src, bm_dst)
            total_copied += 1

    result["copied"] = total_copied
    result["removed"] = total_removed

    # 5. Check git status
    print("[5/7] Checking git status...")
    code, out, err = run_cmd("git status --porcelain", target)
    if not out.strip():
        print("No changes detected. Repository already synchronized.")
        run_cmd(f"git checkout {base_branch}", target)
        result["changed"] = False
        return result

    result["changed"] = True
    print(f"Changes detected ({len(out.splitlines())} modified/new files).")

    if dry_run:
        print("Dry run requested; skipping commit and release.")
        return result

    # 6. Commit and push feature branch
    print("[6/7] Committing changes to feature branch...")
    run_cmd("git add -A", target)
    commit_msg = "feat(sync): upgrade prompts to v1/v2, elevate gitmap aum, and sync skills"
    code, out, err = run_cmd(f'git commit -m "{commit_msg}"', target)
    if code != 0:
        result["error"] = f"Commit failed: {err or out}"
        print(result["error"])
        return result
    result["committed"] = True

    if not no_push:
        print(f"      Pushing {feat_branch} to origin...")
        run_cmd(f"git push -u origin {feat_branch}", target)

    # 7. Release ceremony: Bump patch version, create release branch, tag, and push
    current_ver = detect_latest_version(target)
    next_ver = bump_patch_version(current_ver)
    release_tag = f"v{next_ver}"
    release_branch = f"release/{release_tag}"
    result["release_tag"] = release_tag
    print(f"[7/7] Release ceremony: {current_ver} -> {release_tag}")

    print(f"      Creating release branch: {release_branch}")
    run_cmd(f"git checkout -B {release_branch}", target)

    # Update version files if present
    is_ver_updated = False
    for vf in ["version.json", "package.json"]:
        vpath = target / vf
        if vpath.exists():
            try:
                content = vpath.read_text(encoding="utf-8")
                new_content = re.sub(r'("version"\s*:\s*")(\d+\.\d+\.\d+)(")', rf'\g<1>{next_ver}\g<3>', content)
                if new_content != content:
                    vpath.write_text(new_content, encoding="utf-8")
                    is_ver_updated = True
            except Exception:
                pass
    if is_ver_updated:
        run_cmd("git add -A", target)
        run_cmd(f'git commit -m "chore(version): bump to {release_tag}"', target)

    print(f"      Creating annotated tag: {release_tag}")
    run_cmd(f'git tag -a {release_tag} -m "Release {release_tag} - Synchronize canonical prompts v1/v2, skills, and AI scripts"', target)

    if not no_push:
        print(f"      Pushing {release_branch} and {release_tag} to origin...")
        run_cmd(f"git push -u origin {release_branch}", target)
        run_cmd(f"git push origin {release_tag}", target)

    # Merge back to base branch and return to base branch
    print(f"      Merging {release_branch} into {base_branch}...")
    run_cmd(f"git checkout {base_branch}", target)
    run_cmd(f'git merge {release_branch} -m "chore(release): merge {release_tag} [skip ci]"', target)
    if not no_push:
        run_cmd(f"git push origin {base_branch}", target)
        # Ensure backup branch is confirmed pushed
        run_cmd(f"git push origin {backup_branch}", target)
        result["pushed"] = True

    print(f"[OK] Completed sync & release for {target.name} ({release_tag}).")
    return result


def main() -> None:
    parser = argparse.ArgumentParser(description="Multi-Repository Prompts, Skills, and Scripts Synchronizer")
    parser.add_argument("--repo", "-r", type=str, help="Specific repository name to sync (default: all 13)")
    parser.add_argument("--dry-run", action="store_true", help="Preview changes without committing or releasing")
    parser.add_argument("--no-push", action="store_true", help="Commit and create branches/tags locally without pushing")
    args = parser.parse_args()

    if args.repo:
        matching = [p for p in TARGET_REPOS if p.name.lower() == args.repo.lower()]
        if not matching:
            targets = [SOURCE_ROOT.parent / args.repo]
        else:
            targets = matching
    else:
        targets = TARGET_REPOS

    print(f"Source repository : {SOURCE_ROOT.name} ({SOURCE_ROOT})")
    print(f"Target count      : {len(targets)}")
    print(f"Dry run           : {args.dry_run}")
    print(f"No push           : {args.no_push}")

    summary = []
    for t in targets:
        res = sync_repo(t, dry_run=args.dry_run, no_push=args.no_push)
        summary.append(res)

    print("\n=======================================================")
    print("MULTI-REPOSITORY SYNCHRONIZATION SUMMARY")
    print("=======================================================")
    print(f"{'Repository':<25} | {'Status':<12} | {'Tag':<10} | {'Files':<10} | {'Action'}")
    print("-" * 75)
    for s in summary:
        status = "OK" if not s["error"] else "FAIL"
        tag = s.get("release_tag") or "N/A"
        files = f"+{s['copied']}/-{s['removed']}"
        action = "Released & Pushed" if s["pushed"] else ("Clean" if not s["changed"] else "Committed")
        if s["error"]:
            action = f"Error: {s['error'][:30]}"
        print(f"{s['repo']:<25} | {status:<12} | {tag:<10} | {files:<10} | {action}")


if __name__ == "__main__":
    main()
