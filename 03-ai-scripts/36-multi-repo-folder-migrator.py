#!/usr/bin/env python3
"""Multi-Repository Folder Structure Migrator
Migrates repositories from legacy layouts (.lovable/, spec/) to the new standards (.ai-memory/, 02-spec/).
Performs git pull first, folder renames, deep content string replacement, redundant folder cleanup,
and atomic git commit & push.
"""
from __future__ import annotations

import os
from pathlib import Path
import re
import shutil
import subprocess
import sys

IGNORE_DIRS = {
    ".git",
    "node_modules",
    "vendor",
    ".venv",
    "venv",
    "dist",
    "build",
    ".idea",
    ".vscode",
    "tmp",
}

BINARY_EXTS = {
    ".exe", ".dll", ".so", ".dylib", ".bin", ".iso", ".zip", ".tar", ".gz",
    ".7z", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".pdf", ".lockb",
    ".sqlite", ".db", ".woff", ".woff2", ".ttf", ".eot", ".mp4", ".mp3",
    ".syso",
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


def replace_in_files(repo_root: Path) -> int:
    changed_count = 0
    for root, dirs, files in os.walk(repo_root):
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
        for fname in files:
            p = Path(root) / fname
            if p.suffix.lower() in BINARY_EXTS:
                continue
            try:
                content = p.read_text(encoding="utf-8")
            except Exception:
                continue

            new_content = content
            if ".lovable" in new_content:
                new_content = new_content.replace(".lovable", ".ai-memory")

            if new_content != content:
                p.write_text(new_content, encoding="utf-8", newline="\n")
                changed_count += 1

    return changed_count


def cleanup_redundant_dirs(repo_root: Path) -> None:
    ai_mem = repo_root / ".ai-memory"
    if not ai_mem.exists():
        return

    prompts = ai_mem / "prompts"
    if prompts.exists() and (repo_root / "01-prompts").exists():
        shutil.rmtree(prompts, ignore_errors=True)

    scripts = ai_mem / "ai-fix-scripts"
    if scripts.exists() and (repo_root / "03-ai-scripts").exists():
        shutil.rmtree(scripts, ignore_errors=True)

    cg = ai_mem / "coding-guidelines"
    if cg.exists():
        shutil.rmtree(cg, ignore_errors=True)


def migrate_single_repo(repo_path: Path) -> bool:
    print(f"\n==========================================")
    print(f"Migrating: {repo_path}")
    print(f"==========================================")

    if not repo_path.exists() or not (repo_path / ".git").exists():
        print(f"Error: {repo_path} is not a valid git repository.")
        return False

    # 1. Git pull
    print("Step 1: Running git pull...")
    code, out, err = run_cmd("git pull", repo_path)
    if code != 0:
        print(f"Git pull warning/error: {err or out}")
    else:
        print(f"Git pull: {out}")

    # 2. Rename .lovable -> .ai-memory
    lovable_dir = repo_path / ".lovable"
    ai_mem_dir = repo_path / ".ai-memory"
    if lovable_dir.exists():
        print("Step 2: Renaming .lovable -> .ai-memory...")
        if ai_mem_dir.exists():
            # Merge if both exist
            for item in lovable_dir.iterdir():
                dest = ai_mem_dir / item.name
                if not dest.exists():
                    shutil.move(str(item), str(dest))
            shutil.rmtree(lovable_dir, ignore_errors=True)
        else:
            shutil.move(str(lovable_dir), str(ai_mem_dir))

    # 3. Clean redundant dirs
    print("Step 3: Cleaning redundant directories...")
    cleanup_redundant_dirs(repo_path)

    # 4. Global replace .lovable -> .ai-memory in file contents
    print("Step 4: Replacing .lovable references across all files...")
    num_changed = replace_in_files(repo_path)
    print(f"Updated {num_changed} files with new .ai-memory paths.")

    # 5. Check and migrate spec -> 02-spec if spec exists
    spec_dir = repo_path / "spec"
    if spec_dir.exists() and spec_dir.is_dir():
        print("Step 5: Migrating spec -> 02-spec using 25-repo-migrator.py...")
        migrator_script = Path("d:/work/coding-guidelines/03-ai-scripts/25-repo-migrator.py")
        code, out, err = run_cmd(f"python {migrator_script} migrate --path {repo_path} --force", repo_path)
        print(f"Migrator output: {out}")
        if code != 0:
            print(f"Migrator error: {err}")

    # 6. Check git status
    code, out, err = run_cmd("git status --porcelain", repo_path)
    if not out.strip():
        print("No changes needed. Repository is already up to date.")
        return True

    # 7. Commit & Push
    print("Step 6: Committing and pushing changes...")
    run_cmd("git add -A", repo_path)
    commit_msg = "refactor(structure): migrate .lovable to .ai-memory, move spec to 02-spec, and update all references"
    code, out, err = run_cmd(f'git commit -m "{commit_msg}"', repo_path)
    print(f"Git commit: {out}")

    code, out, err = run_cmd("git push origin HEAD", repo_path)
    if code == 0:
        print(f"Git push succeeded: {out}")
    else:
        print(f"Git push output: {out} {err}")

    print(f"Successfully finished migration for {repo_path}!\n")
    return True


def main() -> None:
    target_repos = sys.argv[1:]
    if not target_repos:
        print("Usage: python 34-multi-repo-folder-migrator.py <repo_path> [repo_path2...]")
        sys.exit(1)

    for r in target_repos:
        migrate_single_repo(Path(r))


if __name__ == "__main__":
    main()
