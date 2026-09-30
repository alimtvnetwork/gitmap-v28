#!/usr/bin/env python3
"""Multi-Repository Prompts, Skills, and AI Scripts Synchronizer & Release Orchestrator.

Synchronizes canonical flat prompts (01-prompts/), Antigravity skills (.agents/skills/),
and AI scripts (03-ai-scripts/ and .agents/scripts/) from coding-guidelines
across all 43 connected repositories.

IMPORTANT: Never synchronizes 06-old-prompts/ (archived v1/v2/v3 prompts).

Performs the complete, safe multi-branch backup and release ceremony per repo:
1. Detect base/current branch and pull latest changes.
2. Create and push pre-change backup branch (backup/pre-flat-v4-<timestamp>).
3. Ensure a pre-change release tag and release branch exist and are pushed before changes.
4. Mirror 01-prompts/ (flat V4), .agents/skills/, 03-ai-scripts/, and .agents/scripts/ cleanly.
5. Commit and push changes on the main/base branch.
6. Create post-change release branch (release/vX.Y.Z) and tag (vX.Y.Z) and push to origin.
7. Merge release branch back into base branch and push.
"""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
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
WORK_ROOT = SOURCE_ROOT.parent

TARGET_REPOS = [
    WORK_ROOT / "02-prompts" / "ai-empathy-prompt-tuner",
    WORK_ROOT / "alim-cv",
    WORK_ROOT / "alim-karim-profile",
    WORK_ROOT / "aukgit" / "alim.karim.profile",
    WORK_ROOT / "antigravity-manager",
    WORK_ROOT / "cat-my",
    WORK_ROOT / "03-aukgo" / "core",
    WORK_ROOT / "digital-name-card",
    WORK_ROOT / "presentations-repos" / "flat-slide-show",
    WORK_ROOT / "gitlogger-new",
    WORK_ROOT / "gitmap",
    WORK_ROOT / "presentations-repos" / "global-ppt-v1",
    WORK_ROOT / "presentations-repos" / "hiltrax",
    WORK_ROOT / "icon-coding-guidelines",
    WORK_ROOT / "img-pdf",
    WORK_ROOT / "presentations-repos" / "ki-health-ppt",
    WORK_ROOT / "aukgit" / "kubernetes-training",
    WORK_ROOT / "lara-licensing",
    WORK_ROOT / "lara-publishing",
    WORK_ROOT / "laravel-automation",
    WORK_ROOT / "letsmarknow-ui",
    WORK_ROOT / "letsmarknow",
    WORK_ROOT / "macro-ahk",
    WORK_ROOT / "presentations-repos" / "maid-app-spec-presentation",
    WORK_ROOT / "movie-cli",
    WORK_ROOT / "03-aukgo" / "pathhelper",
    WORK_ROOT / "presentations-repos" / "presentation-aug-2026-plans-alim",
    WORK_ROOT / "02-prompts" / "prompts-connect",
    WORK_ROOT / "punam-case-studies-v1",
    WORK_ROOT / "presentations-repos" / "rasia-logo",
    WORK_ROOT / "scripts-fixer",
    WORK_ROOT / "presentations-repos" / "slides-spec",
    WORK_ROOT / "spec-builder",
    WORK_ROOT / "web-system" / "sweet-digs-finder",
    WORK_ROOT / "ui-prompts-cat",
    WORK_ROOT / "presentations-repos" / "white-presentation-v1",
    WORK_ROOT / "workflowy-ui",
    WORK_ROOT / "workflowy",
    WORK_ROOT / "wp-exam",
    WORK_ROOT / "wp-git-log",
    WORK_ROOT / "wp-html-automate",
    WORK_ROOT / "wp-link-manager",
    WORK_ROOT / "wp-onboarding",
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
    "06-old-prompts",
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
    """Detect default branch (main or master or current)."""
    code, out, _ = run_cmd("git branch --show-current", target)
    curr = out.strip()
    if curr and not curr.startswith(("release/", "backup/", "feat/")):
        return curr
    code, out, _ = run_cmd("git branch --list main master", target)
    lines = [line.strip().replace("*", "").strip() for line in out.splitlines()]
    if "main" in lines:
        return "main"
    if "master" in lines:
        return "master"
    return curr or "main"


def detect_latest_version(target: Path) -> str:
    """Detect current SemVer string from git tags or version files."""
    code, out, _ = run_cmd("git tag --list \"v[0-9]*\" --sort=-v:refname", target)
    if code == 0 and out.strip():
        for line in out.splitlines():
            tag = line.strip().lstrip("v")
            match = re.match(r"^(\d+\.\d+\.\d+)$", tag)
            if match:
                return match.group(1)

    for vf in ["version.json", "package.json", "version.txt", "VERSION"]:
        vpath = target / vf
        if vpath.exists():
            try:
                content = vpath.read_text(encoding="utf-8")
                match = re.search(r'"[Vv]ersion"\s*:\s*"(\d+\.\d+\.\d+)"', content)
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
                    elif src_file.read_bytes() != dst_file.read_bytes():
                        need_copy = True
                except Exception:
                    need_copy = True

            if need_copy:
                shutil.copy2(src_file, dst_file)
                copied += 1

    return copied, removed


def ensure_pre_release(target: Path, current_ver: str, dry_run: bool, no_push: bool) -> str:
    """Ensure a pre-change release tag and release branch exist for current HEAD."""
    code, head_tags, _ = run_cmd("git tag --points-at HEAD", target)
    existing_ver_tags = [
        t.strip() for t in head_tags.splitlines() if re.match(r"^v\d+\.\d+\.\d+$", t.strip())
    ]
    if existing_ver_tags:
        pre_tag = existing_ver_tags[0]
    else:
        # Check if v{current_ver} already exists on an older commit
        code, tag_Lookup, _ = run_cmd(f"git tag --list v{current_ver}", target)
        if tag_Lookup.strip():
            pre_ver = bump_patch_version(current_ver)
        else:
            pre_ver = current_ver
        pre_tag = f"v{pre_ver}"
        if not dry_run:
            run_cmd(f'git tag -a {pre_tag} -m "Pre-sync release {pre_tag}"', target)

    pre_rel_branch = f"release/{pre_tag}"
    if not dry_run:
        run_cmd(f"git branch -f {pre_rel_branch} HEAD", target)
        if not no_push:
            run_cmd(f"git push origin {pre_rel_branch}", target)
            run_cmd(f"git push origin {pre_tag}", target)
    return pre_tag


def sync_repo(target: Path, dry_run: bool = False, no_push: bool = False) -> dict[str, object]:
    result: dict[str, object] = {
        "repo": target.name,
        "path": str(target),
        "base_branch": "",
        "backup_branch": "",
        "pre_release_tag": "",
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
        result["error"] = err
        return result

    base_branch = detect_base_branch(target)
    result["base_branch"] = base_branch

    # 1. Checkout base branch and pull latest
    run_cmd(f"git checkout {base_branch}", target)
    if not dry_run and not no_push:
        run_cmd(f"git pull origin {base_branch} --no-rebase", target)

    # 2. Create and push pre-change backup branch from HEAD
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    backup_branch = f"backup/pre-flat-v4-{timestamp}"
    result["backup_branch"] = backup_branch
    if not dry_run:
        run_cmd(f"git branch {backup_branch} HEAD", target)
        if not no_push:
            run_cmd(f"git push origin {backup_branch}", target)

    # 3. Ensure pre-change release tag and release branch exist before modifying
    current_ver = detect_latest_version(target)
    pre_tag = ensure_pre_release(target, current_ver, dry_run=dry_run, no_push=no_push)
    result["pre_release_tag"] = pre_tag
    current_ver = pre_tag.lstrip("v")

    # Ensure 06-old-prompts is never present in target repository
    old_prompts_dir = target / "06-old-prompts"
    if old_prompts_dir.exists():
        shutil.rmtree(old_prompts_dir, ignore_errors=True)

    # 4. Mirror directories cleanly
    total_copied = 0
    total_removed = 0
    for src_rel, dst_rel in SYNC_DIRS:
        src_path = SOURCE_ROOT / src_rel
        dst_path = target / dst_rel
        c, r = mirror_directory(src_path, dst_path)
        total_copied += c
        total_removed += r

    result["copied"] = total_copied
    result["removed"] = total_removed

    # 5. Check git status
    code, out, err = run_cmd("git status --porcelain", target)
    if not out.strip():
        result["changed"] = False
        result["release_tag"] = pre_tag
        result["pushed"] = not no_push
        return result

    result["changed"] = True

    if dry_run:
        return result

    # 6. Commit changes on base branch and push
    run_cmd("git add -A", target)
    commit_msg = "feat(sync): flatten v4 prompts into 01-prompts, fix above/below mandates, and sync skills"
    code, out, err = run_cmd(f'git commit -m "{commit_msg}"', target)
    if code != 0:
        result["error"] = f"Commit failed: {err or out}"
        return result
    result["committed"] = True

    if not no_push:
        run_cmd(f"git push -u origin {base_branch}", target)

    # 7. Post-change release ceremony: Bump patch version, create release branch & tag, and push
    next_ver = bump_patch_version(current_ver)
    release_tag = f"v{next_ver}"
    release_branch = f"release/{release_tag}"
    result["release_tag"] = release_tag

    run_cmd(f"git checkout -B {release_branch}", target)

    # Update version files if present
    is_ver_updated = False
    for vf in ["version.json", "package.json"]:
        vpath = target / vf
        if vpath.exists():
            try:
                content = vpath.read_text(encoding="utf-8")
                new_content = re.sub(
                    r'("[Vv]ersion"\s*:\s*")(\d+\.\d+\.\d+)(")',
                    rf"\g<1>{next_ver}\g<3>",
                    content,
                )
                if new_content != content:
                    vpath.write_text(new_content, encoding="utf-8")
                    is_ver_updated = True
            except Exception:
                pass
    if is_ver_updated:
        run_cmd("git add -A", target)
        run_cmd(f'git commit -m "chore(version): bump to {release_tag}"', target)

    run_cmd(
        f'git tag -a {release_tag} -m "Release {release_tag} - Synchronize flat V4 prompts, skills, and AI scripts"',
        target,
    )

    if not no_push:
        run_cmd(f"git push -u origin {release_branch}", target)
        run_cmd(f"git push origin {release_tag}", target)

    # Merge back into base branch and push
    run_cmd(f"git checkout {base_branch}", target)
    run_cmd(f'git merge {release_branch} -m "chore(release): merge {release_tag} [skip ci]"', target)
    if not no_push:
        p_code, p_out, p_err = run_cmd(f"git push origin {base_branch}", target)
        if p_code == 0:
            result["pushed"] = True
        else:
            result["error"] = f"Push failed: {p_err or p_out}"

    print(
        f"[OK] {target.name:<32} | Pre: {pre_tag:<9} | Post: {release_tag:<9} | +{total_copied}/-{total_removed}"
    )
    return result


def main() -> None:
    parser = argparse.ArgumentParser(description="Multi-Repository Prompts, Skills, and Scripts Synchronizer")
    parser.add_argument("--repo", "-r", type=str, help="Specific repository name to sync (default: all 43)")
    parser.add_argument("--workers", "-w", type=int, default=6, help="Parallel worker count (default: 6)")
    parser.add_argument("--dry-run", action="store_true", help="Preview changes without committing or releasing")
    parser.add_argument("--no-push", action="store_true", help="Commit and create branches/tags locally without pushing")
    args = parser.parse_args()

    if args.repo:
        matching = [p for p in TARGET_REPOS if p.name.lower() == args.repo.lower()]
        if not matching:
            targets = [WORK_ROOT / args.repo]
        else:
            targets = matching
    else:
        targets = TARGET_REPOS

    print(f"Source repository : {SOURCE_ROOT.name} ({SOURCE_ROOT})")
    print(f"Target count      : {len(targets)}")
    print(f"Workers           : {args.workers}")
    print(f"Dry run           : {args.dry_run}")
    print(f"No push           : {args.no_push}")
    print("=" * 95)

    summary = []
    if args.workers > 1 and len(targets) > 1:
        with ThreadPoolExecutor(max_workers=args.workers) as pool:
            future_map = {
                pool.submit(sync_repo, t, args.dry_run, args.no_push): t for t in targets
            }
            for fut in as_completed(future_map):
                t = future_map[fut]
                try:
                    res = fut.result()
                except Exception as exc:
                    res = {
                        "repo": t.name,
                        "path": str(t),
                        "pre_release_tag": "N/A",
                        "release_tag": "N/A",
                        "copied": 0,
                        "removed": 0,
                        "changed": False,
                        "pushed": False,
                        "error": str(exc),
                    }
                summary.append(res)
    else:
        for t in targets:
            res = sync_repo(t, dry_run=args.dry_run, no_push=args.no_push)
            summary.append(res)

    summary.sort(key=lambda x: str(x["repo"]).lower())

    print("\n===============================================================================================")
    print("MULTI-REPOSITORY SYNCHRONIZATION SUMMARY")
    print("===============================================================================================")
    print(
        f"{'Repository':<34} | {'Status':<8} | {'Pre-Tag':<10} | {'Post-Tag':<10} | {'Files':<12} | {'Action'}"
    )
    print("-" * 105)
    for s in summary:
        status = "OK" if not s["error"] else "FAIL"
        pre_tag = str(s.get("pre_release_tag") or "N/A")
        tag = str(s.get("release_tag") or "N/A")
        files = f"+{s['copied']}/-{s['removed']}"
        action = "Released & Pushed" if s["pushed"] else ("Clean" if not s["changed"] else "Committed")
        if s["error"]:
            action = f"Error: {str(s['error'])[:35]}"
        print(f"{str(s['repo']):<34} | {status:<8} | {pre_tag:<10} | {tag:<10} | {files:<12} | {action}")


if __name__ == "__main__":
    main()
