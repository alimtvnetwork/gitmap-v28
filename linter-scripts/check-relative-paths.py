#!/usr/bin/env python3
"""Linter to verify that no absolute filesystem paths or file:/// URIs exist in repository files."""
import os
import re
import subprocess
import sys

# Patterns that indicate absolute or non-portable file paths in repository documentation & code
FORBIDDEN_PATTERNS = [
    (re.compile(r"file:///[a-zA-Z]:[/\\]?", re.IGNORECASE), "Absolute file:/// URI with drive letter"),
    (re.compile(r"file:///work/", re.IGNORECASE), "Absolute file:/// URI"),
    (re.compile(r"file:///(?:Users|home|root)/", re.IGNORECASE), "Absolute file:/// URI to user directory"),
    (re.compile(r"\b[dD]:[/\\]work[/\\]gitmap\b", re.IGNORECASE), "Hardcoded absolute repo path (D:\\work\\gitmap)"),
    (re.compile(r"\b[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+[/\\]\.gemini\b", re.IGNORECASE), "Hardcoded user agent directory"),
]

EXCLUDE_DIRS = {".git", "node_modules", "dist", "build", "vendor", ".cache", ".next"}
EXCLUDE_EXTS = {".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip", ".gz", ".tar", ".exe", ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf"}

# Allowlist for files that legitimately mention drive letters in documentation / regex specs
ALLOWLIST_FILES = {
    ".github/workflows/goreleaser-smoke.yml",  # contains comment explaining Windows file:/// URL format
    "linter-scripts/check-relative-paths.py",
    ".lovable/ai-fix-scripts/04-relative-path-fixer.py",
}


def scan_file(filepath, repo_root):
    rel_path = os.path.relpath(filepath, repo_root).replace("\\", "/")
    if rel_path in ALLOWLIST_FILES:
        return []

    try:
        with open(filepath, "r", encoding="utf-8", errors="replace") as fh:
            lines = fh.readlines()
    except Exception as e:
        return [(rel_path, 0, f"Failed to read file: {e}")]

    violations = []
    for line_idx, line in enumerate(lines, 1):
        for pattern, desc in FORBIDDEN_PATTERNS:
            if pattern.search(line):
                # Skip comments in this linter or fixer itself
                if "FORBIDDEN_PATTERNS" in line or "check-relative-paths.py" in line:
                    continue
                violations.append((rel_path, line_idx, f"{desc}: {line.strip()[:100]}"))

    return violations


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    print(f"=== Checking for Absolute Paths and file:/// URIs in {repo_root} ===")

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
    all_violations = []

    for f in files:
        ext = os.path.splitext(f)[1].lower()
        if ext in EXCLUDE_EXTS:
            continue
        full_path = os.path.join(repo_root, f)
        if os.path.isfile(full_path):
            violations = scan_file(full_path, repo_root)
            all_violations.extend(violations)

    if all_violations:
        print(f"\n❌ FAIL: Found {len(all_violations)} absolute path / URI violation(s):\n")
        for file_path, line_no, msg in all_violations:
            print(f"  {file_path}:{line_no}: {msg}")
        sys.exit(1)

    print(f"\n✅ PASS: No absolute filesystem paths or file:/// URIs found across {len(files)} tracked files.")
    sys.exit(0)


if __name__ == "__main__":
    main()
