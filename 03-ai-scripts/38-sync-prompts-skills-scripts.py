#!/usr/bin/env python3
"""Multi-Repository Prompts, Skills, and AI Scripts Synchronizer.

Synchronizes canonical prompts (01-prompts/), Antigravity skills (.agents/skills/),
and AI scripts (03-ai-scripts/ and .agents/scripts/) from coding-guidelines-v24
across multiple target repositories.

Performs git pull first, safe directory mirroring, git add, git commit, and git push.
"""
from __future__ import annotations

import os
from pathlib import Path
import shutil
import subprocess
import sys

SOURCE_ROOT = Path(__file__).resolve().parent.parent

TARGET_REPOS = [
    SOURCE_ROOT.parent / "antigravity-manager",
    SOURCE_ROOT.parent / "cat-my",
    SOURCE_ROOT.parent / "movie-cli",
    SOURCE_ROOT.parent / "scripts-fixer",
    SOURCE_ROOT.parent / "spec-builder",
    SOURCE_ROOT.parent / "kita-social-media-content-calender",
    SOURCE_ROOT.parent / "laravel-automation",
    SOURCE_ROOT.parent / "wp-exam",
    SOURCE_ROOT.parent / "gitmap",
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


def run_cmd(cmd: str, cwd: Path) -> tuple[int, str, str]:
    res = subprocess.run(
        cmd,
        cwd=str(cwd),
        shell=True,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    return res.returncode, res.stdout.strip(), res.stderr.strip()


def mirror_directory(src: Path, dst: Path) -> tuple[int, int]:
    """Mirror src into dst, removing stale files and copying new/updated ones.
    Returns (copied_count, removed_count).
    """
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
        # Filter out excluded directories
        dirs[:] = [d for d in dirs if d not in EXCLUDE_NAMES]

        rel_root = Path(root).relative_to(src)
        target_dir = dst / rel_root
        target_dir.mkdir(parents=True, exist_ok=True)

        for f in files:
            if f in EXCLUDE_NAMES or Path(f).suffix.lower() in EXCLUDE_EXTS:
                continue

            src_file = Path(root) / f
            dst_file = target_dir / f

            # Check if file changed
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
                        # Compare content bytes if sizes match
                        if src_file.read_bytes() != dst_file.read_bytes():
                            need_copy = True
                except Exception:
                    need_copy = True

            if need_copy:
                shutil.copy2(src_file, dst_file)
                copied += 1

    return copied, removed


def sync_repo(target: Path) -> dict[str, any]:
    print(f"\n=======================================================")
    print(f"Syncing target repository: {target}")
    print(f"=======================================================")

    result = {
        "repo": str(target),
        "pull_status": "ok",
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

    # 1. Git pull
    print("[1/4] Running git pull...")
    code, out, err = run_cmd("git pull", target)
    if code != 0:
        print(f"Git pull warning/error: {err or out}")
        result["pull_status"] = f"error: {err or out}"
    else:
        print(f"Git pull output: {out}")
        result["pull_status"] = out

    # 2. Mirror each directory
    print("[2/4] Mirroring prompts, skills, and AI scripts...")
    total_copied = 0
    total_removed = 0
    for src_rel, dst_rel in SYNC_DIRS:
        src_path = SOURCE_ROOT / src_rel
        dst_path = target / dst_rel
        c, r = mirror_directory(src_path, dst_path)
        print(f"  - {src_rel} -> {dst_rel} (copied: {c}, removed: {r})")
        total_copied += c
        total_removed += r

    result["copied"] = total_copied
    result["removed"] = total_removed

    # 3. Check git status
    print("[3/4] Checking git status...")
    code, out, err = run_cmd("git status --porcelain", target)
    if not out.strip():
        print("No changes detected. Target repository is already up-to-date.")
        result["changed"] = False
        return result

    result["changed"] = True
    print(f"Changes detected in {target}:\n{out[:500]}...")

    # 4. Git add, commit, and push
    print("[4/4] Staging, committing, and pushing...")
    run_cmd("git add -A", target)
    commit_msg = "feat(sync): synchronize prompts, skills, and ai-scripts from coding-guidelines-v24"
    code, out, err = run_cmd(f'git commit -m "{commit_msg}"', target)
    if code != 0:
        print(f"Commit error: {err or out}")
        result["error"] = f"commit failed: {err or out}"
        return result
    print(f"Committed: {out}")
    result["committed"] = True

    code, out, err = run_cmd("git push origin HEAD", target)
    if code != 0:
        print(f"Push error: {err or out}")
        result["error"] = f"push failed: {err or out}"
        return result
    print(f"Pushed to remote: {out}")
    result["pushed"] = True

    return result


def main() -> None:
    targets = [Path(p) for p in sys.argv[1:]] if len(sys.argv) > 1 else TARGET_REPOS

    print(f"Source repository: {SOURCE_ROOT}")
    print(f"Target count: {len(targets)}")

    summary = []
    for t in targets:
        res = sync_repo(t)
        summary.append(res)

    print("\n=======================================================")
    print("SYNCHRONIZATION SUMMARY")
    print("=======================================================")
    for s in summary:
        status = "OK" if not s["error"] else f"FAIL ({s['error']})"
        action = "Pushed" if s["pushed"] else ("Clean (No change)" if not s["changed"] else "Committed (Not pushed)")
        print(f"{s['repo']}: {status} | Action: {action} | Copied: {s['copied']} | Removed: {s['removed']}")


if __name__ == "__main__":
    main()
