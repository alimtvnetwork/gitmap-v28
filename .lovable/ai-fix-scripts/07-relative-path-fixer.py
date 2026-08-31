#!/usr/bin/env python3
"""Autonomously converts all absolute filesystem paths and file:/// URIs across the codebase to strictly relative Git root paths.

Usage:
    python .lovable/ai-fix-scripts/04-relative-path-fixer.py [target_directory]
"""
import os
import re
import subprocess
import sys

EXCLUDE_DIRS = {".git", "node_modules", "dist", "build", "vendor", ".cache", ".next"}
EXCLUDE_EXTS = {".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip", ".gz", ".tar", ".exe", ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf"}

ALLOWLIST_FILES = {
    ".github/workflows/goreleaser-smoke.yml",
    "linter-scripts/check-relative-paths.py",
    ".lovable/ai-fix-scripts/07-relative-path-fixer.py",
}


def clean_file_content(content: str, repo_root: str, file_relpath: str) -> str:
    # 1. file:///d:/work/gitmap/ (or D:, with optional percent-encoding like %3A or %3a)
    content = re.sub(r"file:///[dD](?:%3[aA]|:)/work/gitmap/", "", content)
    content = re.sub(r"file:///[dD]:[/\\]work[/\\]gitmap[/\\]?", "", content)
    content = re.sub(r"file:///[dD]:", "", content)

    # 2. Markdown links: [text](file:///.../gitmap/rel/path) -> [text](rel/path)
    content = re.sub(r"\[([^\]]+)\]\(file:///[a-zA-Z]:[/\\]work[/\\]gitmap[/\\]([^)]+)\)", r"[\1](\2)", content)
    content = re.sub(r"\[([^\]]+)\]\(file:///[^)]*gitmap/([^)]+)\)", r"[\1](\2)", content)

    # 3. Direct hardcoded repo drive paths: D:\work\gitmap\ or d:/work/gitmap/
    # In Windows backslash form
    content = re.sub(r"\b[dD]:\\work\\gitmap\\", "", content)
    # In forward slash form
    content = re.sub(r"\b[dD]:/work/gitmap/", "", content)

    # 4. Agent log / brain path references
    content = re.sub(r"[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+[/\\]\.gemini[/\\]antigravity[/\\]brain[/\\][a-f0-9-]+[/\\]", ".lovable/", content)

    # 5. Fix scripts referencing hardcoded paths
    if file_relpath.endswith(".py"):
        content = re.sub(r'r?["\'][dD]:[/\\]work[/\\]gitmap[/\\]gitmap["\']', 'os.path.join(repo_root, "gitmap")', content)
        content = re.sub(r'r?["\'][dD]:[/\\]work[/\\]gitmap["\']', 'repo_root', content)

    # 6. Normalize backslashes inside markdown link targets [title](foo\bar.md) -> [title](foo/bar.md)
    def normalize_link_slashes(match):
        title = match.group(1)
        target = match.group(2).replace("\\", "/")
        return f"[{title}]({target})"

    content = re.sub(r"\[([^\]]+)\]\(([^)]+\.md)\)", normalize_link_slashes, content)

    return content


def process_file(filepath: str, repo_root: str) -> bool:
    rel_path = os.path.relpath(filepath, repo_root).replace("\\", "/")
    if rel_path in ALLOWLIST_FILES:
        return False

    try:
        with open(filepath, "r", encoding="utf-8", errors="replace") as fh:
            original = fh.read()
    except Exception:
        return False

    updated = clean_file_content(original, repo_root, rel_path)
    if updated != original:
        with open(filepath, "w", encoding="utf-8") as fh:
            fh.write(updated)
        return True
    return False


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    print(f"=== Running Relative Path Fixer on {repo_root} ===")

    res = subprocess.run(
        ["git", "ls-files"],
        cwd=repo_root,
        capture_output=True,
        text=True,
        encoding="utf-8"
    )
    if res.returncode != 0:
        print("Failed to run git ls-files", file=sys.stderr)
        sys.exit(1)

    files = [f.strip() for f in res.stdout.splitlines() if f.strip()]
    modified_count = 0

    for f in files:
        ext = os.path.splitext(f)[1].lower()
        if ext in EXCLUDE_EXTS:
            continue
        full_path = os.path.join(repo_root, f)
        if os.path.isfile(full_path) and process_file(full_path, repo_root):
            rel = os.path.relpath(full_path, repo_root).replace("\\", "/")
            print(f"  Fixed absolute paths in: {rel}")
            modified_count += 1

    print(f"\n[SUCCESS] Fixed absolute paths in {modified_count} file(s).")


if __name__ == "__main__":
    main()
