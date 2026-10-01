#!/usr/bin/env python3
"""Multi-Repository Prompts, Skills, Coding Guidelines, and AI Scripts Synchronizer & Release Orchestrator.

Synchronizes canonical prompts (01-prompts/), Antigravity & Cursor skills (.agents/skills/, .cursor/skills/),
AI scripts (03-ai-scripts/ and .agents/scripts/), and coding guidelines / design system specs
(02-spec/02-coding-guidelines/, 02-spec/07-design-system/, 02-spec/17-consolidated-guidelines/,
.ai-memory/coding-guidelines.md, .ai-memory/prompts.md) from coding-guidelines across all 43 connected repositories.

IMPORTANT: Never synchronizes 06-old-prompts/ (archived v1/v2/v3 prompts).

Performs the complete, safe multi-branch backup and release ceremony per repo:
1. Detect base/current branch and pull latest changes (`git pull origin <base_branch> --no-rebase`).
2. Create and push pre-change backup branch (`backup/pre-v3-nsteps-sync-<timestamp>`).
3. Ensure a pre-change release tag (`vX.Y.Z`) and release branch (`release/vX.Y.Z`) exist and are pushed before changes.
4. Mirror 01-prompts/, .agents/skills/, .cursor/skills/, 03-ai-scripts/, .agents/scripts/, and conditional 02-spec / .ai-memory guidelines cleanly.
5. Commit and push changes on the main/base branch.
6. Create post-change release branch (`release/v<next_ver>`) and tag (`v<next_ver>`) and push to origin.
7. Merge release branch back into base branch with `[skip ci]` and push.
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
    (".cursor/skills", ".cursor/skills"),
    ("03-ai-scripts", "03-ai-scripts"),
    (".agents/scripts", ".agents/scripts"),
    ("02-spec/02-coding-guidelines", "02-spec/02-coding-guidelines"),
]

CONDITIONAL_SPEC_DIRS = [
    "02-spec/07-design-system",
    "02-spec/17-consolidated-guidelines",
]

CONDITIONAL_MEMORY_FILES = [
    ".ai-memory/coding-guidelines.md",
    ".ai-memory/prompts.md",
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
    _, out, _ = run_cmd("git branch --show-current", target)
    curr = out.strip()
    is_ephemeral = curr.startswith(("release/", "backup/", "feat/"))

    if curr and not is_ephemeral:
        return curr

    _, out, _ = run_cmd("git branch --list main master", target)
    lines = [line.strip().replace("*", "").strip() for line in out.splitlines()]

    if "main" in lines:
        return "main"

    if "master" in lines:
        return "master"

    return curr or "main"


def detect_latest_version(target: Path) -> str:
    """Detect current SemVer string from git tags or version files."""
    code, out, _ = run_cmd("git tag --list \"v[0-9]*\" --sort=-v:refname", target)
    has_tag_output = code == 0 and bool(out.strip())

    if has_tag_output:
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
    has_three_parts = len(parts) >= 3 and parts[2].isdigit()

    if has_three_parts:
        parts[2] = str(int(parts[2]) + 1)

        return ".".join(parts[:3])

    return f"{ver_str}.1"


def copy_single_file(src_file: Path, dst_file: Path, is_dry_run: bool = False) -> int:
    """Copy a single file if missing or changed. Returns 1 if copied, 0 otherwise."""
    if not src_file.exists():
        return 0

    is_copy_needed = False

    if not dst_file.exists():
        is_copy_needed = True
    else:
        try:
            if src_file.stat().st_size != dst_file.stat().st_size:
                is_copy_needed = True
            elif src_file.read_bytes() != dst_file.read_bytes():
                is_copy_needed = True
        except Exception:
            is_copy_needed = True

    if is_copy_needed:
        if not is_dry_run:
            dst_file.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(src_file, dst_file)

        return 1

    return 0


def mirror_directory(src: Path, dst: Path, is_dry_run: bool = False) -> tuple[int, int]:
    """Mirror src into dst, removing stale files and copying new/updated ones."""
    copied = 0
    removed = 0

    if not src.exists():
        return 0, 0

    if not is_dry_run:
        dst.mkdir(parents=True, exist_ok=True)

    # 1. Clean stale files/dirs in dst that are no longer in src
    if dst.exists():
        for root, dirs, files in os.walk(dst, topdown=False):
            rel_root = Path(root).relative_to(dst)
            src_root = src / rel_root

            for f in files:
                dst_file = Path(root) / f
                src_file = src_root / f
                is_excluded = f in EXCLUDE_NAMES or dst_file.suffix.lower() in EXCLUDE_EXTS

                if is_excluded or not src_file.exists():
                    if not is_dry_run:
                        dst_file.unlink(missing_ok=True)
                    removed += 1

            for d in dirs:
                dst_dir = Path(root) / d
                src_dir = src_root / d

                if d in EXCLUDE_NAMES or not src_dir.exists():
                    if not is_dry_run:
                        shutil.rmtree(dst_dir, ignore_errors=True)
                    removed += 1

    # 2. Copy files from src to dst
    for root, dirs, files in os.walk(src):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_NAMES]

        rel_root = Path(root).relative_to(src)
        target_dir = dst / rel_root

        if not is_dry_run:
            target_dir.mkdir(parents=True, exist_ok=True)

        for f in files:
            is_excluded = f in EXCLUDE_NAMES or Path(f).suffix.lower() in EXCLUDE_EXTS

            if is_excluded:
                continue

            src_file = Path(root) / f
            dst_file = target_dir / f
            copied += copy_single_file(src_file, dst_file, is_dry_run=is_dry_run)

    return copied, removed


def mirror_conditional_guidelines(target: Path, is_dry_run: bool = False) -> tuple[int, int]:
    """Conditionally sync other spec guidelines and .ai-memory guideline indices when present in target."""
    copied = 0
    removed = 0

    for spec_rel in CONDITIONAL_SPEC_DIRS:
        if (target / spec_rel).exists():
            c, r = mirror_directory(
                SOURCE_ROOT / spec_rel,
                target / spec_rel,
                is_dry_run=is_dry_run,
            )
            copied += c
            removed += r

    if (target / ".ai-memory").exists():
        for mem_rel in CONDITIONAL_MEMORY_FILES:
            copied += copy_single_file(
                SOURCE_ROOT / mem_rel,
                target / mem_rel,
                is_dry_run=is_dry_run,
            )

    return copied, removed


def ensure_pre_release(target: Path, current_ver: str, is_dry_run: bool, is_no_push: bool) -> str:
    """Ensure a pre-change release tag and release branch exist and are pushed for current HEAD."""
    _, head_tags, _ = run_cmd("git tag --points-at HEAD", target)
    existing_ver_tags = [
        t.strip() for t in head_tags.splitlines() if re.match(r"^v\d+\.\d+\.\d+$", t.strip())
    ]

    if existing_ver_tags:
        pre_tag = existing_ver_tags[0]
    else:
        _, tag_lookup, _ = run_cmd(f"git tag --list v{current_ver}", target)

        if tag_lookup.strip():
            pre_ver = bump_patch_version(current_ver)
        else:
            pre_ver = current_ver

        pre_tag = f"v{pre_ver}"

        if not is_dry_run:
            run_cmd(f'git tag -a {pre_tag} -m "Pre-sync release {pre_tag}"', target)

    pre_rel_branch = f"release/{pre_tag}"

    if not is_dry_run:
        run_cmd(f"git branch -f {pre_rel_branch} HEAD", target)

        if not is_no_push:
            run_cmd(f"git push -u origin {pre_rel_branch}", target)
            run_cmd(f"git push origin {pre_tag}", target)

    return pre_tag


def sync_repo(target: Path, is_dry_run: bool = False, is_no_push: bool = False) -> dict[str, object]:
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

    has_git_repo = target.exists() and (target / ".git").exists()

    if not has_git_repo:
        err = f"Directory {target} is not a valid git repository."
        result["error"] = err

        return result

    base_branch = detect_base_branch(target)
    result["base_branch"] = base_branch

    # 1. Checkout base branch and pull latest
    if not is_dry_run:
        run_cmd(f"git checkout {base_branch}", target)

    is_live_push = not is_dry_run and not is_no_push

    if is_live_push:
        run_cmd(f"git pull origin {base_branch} --no-rebase", target)

    # 2. Create and push pre-change backup branch from HEAD, then return to base branch
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    backup_branch = f"backup/sync-{timestamp}"
    result["backup_branch"] = backup_branch

    if not is_dry_run:
        run_cmd(f"git checkout -b {backup_branch}", target)

        if not is_no_push:
            run_cmd(f"git push -u origin {backup_branch}", target)

        run_cmd(f"git checkout {base_branch}", target)

    # 3. Ensure pre-change release tag and release branch exist before modifying
    current_ver = detect_latest_version(target)
    pre_tag = ensure_pre_release(
        target,
        current_ver,
        is_dry_run=is_dry_run,
        is_no_push=is_no_push,
    )
    result["pre_release_tag"] = pre_tag
    current_ver = pre_tag.lstrip("v")

    # Ensure 06-old-prompts is never present in target repository
    old_prompts_dir = target / "06-old-prompts"

    if old_prompts_dir.exists() and not is_dry_run:
        shutil.rmtree(old_prompts_dir, ignore_errors=True)

    # 4. Mirror directories and conditional coding guidelines cleanly
    total_copied = 0
    total_removed = 0

    for src_rel, dst_rel in SYNC_DIRS:
        src_path = SOURCE_ROOT / src_rel
        dst_path = target / dst_rel
        c, r = mirror_directory(src_path, dst_path, is_dry_run=is_dry_run)
        total_copied += c
        total_removed += r

    c_spec, r_spec = mirror_conditional_guidelines(target, is_dry_run=is_dry_run)
    total_copied += c_spec
    total_removed += r_spec

    result["copied"] = total_copied
    result["removed"] = total_removed

    if is_dry_run:
        has_dry_changes = (total_copied + total_removed) > 0
        result["changed"] = has_dry_changes
        result["release_tag"] = f"v{bump_patch_version(current_ver)}" if has_dry_changes else pre_tag

        return result

    # 5. Check git status
    _, out, _ = run_cmd("git status --porcelain", target)

    if not out.strip():
        result["changed"] = False
        result["release_tag"] = pre_tag
        result["pushed"] = not is_no_push

        return result

    result["changed"] = True

    # 6. Commit changes on base branch and push
    run_cmd("git add -A", target)
    commit_msg = (
        "feat(sync): sync v6 prompts, sqlite task manager, skills, and coding guidelines"
    )
    code, out, err = run_cmd(f'git commit -m "{commit_msg}"', target)

    if code != 0:
        result["error"] = f"Commit failed: {err or out}"

        return result

    result["committed"] = True

    if not is_no_push:
        run_cmd(f"git push -u origin {base_branch}", target)

    # 7. Post-change release ceremony: Bump patch version, create release branch & tag, and push
    next_ver = bump_patch_version(current_ver)
    release_tag = f"v{next_ver}"
    release_branch = f"release/{release_tag}"
    result["release_tag"] = release_tag

    run_cmd(f"git checkout -B {release_branch}", target)

    # Update version files if present
    is_ver_updated = False

    bump_script = target / "03-ai-scripts" / "37-bump-version.py"
    if bump_script.exists():
        b_code, b_out, b_err = run_cmd(
            f'python "{bump_script}" --tier patch --scope "Synchronize prompts, skills, AI scripts, and coding guidelines"',
            target,
        )
        if b_code == 0:
            is_ver_updated = True
            bumped_ver = detect_latest_version(target)
            if bumped_ver and bumped_ver != current_ver:
                next_ver = bumped_ver
                release_tag = f"v{next_ver}"
                release_branch = f"release/{release_tag}"
                result["release_tag"] = release_tag
    elif (target / "scripts" / "bump-version.mjs").exists():
        b_code, b_out, b_err = run_cmd("node scripts/bump-version.mjs patch", target)
        if b_code == 0:
            is_ver_updated = True
            bumped_ver = detect_latest_version(target)
            if bumped_ver and bumped_ver != current_ver:
                next_ver = bumped_ver
                release_tag = f"v{next_ver}"
                release_branch = f"release/{release_tag}"
                result["release_tag"] = release_tag

    if not is_ver_updated:
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
        f'git tag -a {release_tag} -m "Release {release_tag} - Synchronize prompts, skills, AI scripts, and coding guidelines"',
        target,
    )

    if not is_no_push:
        run_cmd(f"git push -u origin {release_branch}", target)
        run_cmd(f"git push origin {release_tag}", target)

    # Merge back into base branch and push
    run_cmd(f"git checkout {base_branch}", target)
    run_cmd(f'git merge {release_branch} -m "chore(release): merge {release_tag} [skip ci]"', target)

    if not is_no_push:
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

    is_dry_run = bool(args.dry_run)
    is_no_push = bool(args.no_push)

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
    print(f"Dry run           : {is_dry_run}")
    print(f"No push           : {is_no_push}")
    print("=" * 95)

    summary = []
    has_parallel_pool = args.workers > 1 and len(targets) > 1

    if has_parallel_pool:
        with ThreadPoolExecutor(max_workers=args.workers) as pool:
            future_map = {
                pool.submit(sync_repo, t, is_dry_run, is_no_push): t for t in targets
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
            res = sync_repo(t, is_dry_run=is_dry_run, is_no_push=is_no_push)
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
