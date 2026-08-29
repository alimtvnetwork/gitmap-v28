#!/usr/bin/env python3
"""Cross-platform check for forbidden legacy deploy layout paths."""
import os
import re
import sys


EXCLUDE_DIRS = {
    ".git", "node_modules", "dist", "build", "bin", ".next",
    ".gitmap", "vendor", "coverage", ".lovable", "spec"
}

EXCLUDE_EXTS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico",
    ".pdf", ".zip", ".gz", ".tar", ".exe", ".dll", ".so", ".dylib",
    ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf",
    ".json", ".md", ".html"
}

EXEMPT_FILES = {
    "gitmap/constants/deploy-manifest.json",
    "gitmap/constants/deploy_manifest.go",
    ".github/scripts/check-deploy-layout.sh",
    ".github/scripts/check-deploy-layout.py",
    ".github/scripts/check-legacy-refs.sh",
    ".github/scripts/check-legacy-refs.py",
    ".github/scripts/smoke-installer.sh",
    ".github/scripts/smoke-installer.ps1"
}

PATTERNS = [
    re.compile(r'Join-Path\s+\$[A-Za-z_.]*[Dd]eploy[A-Za-z_.]*\s+"gitmap"([^-]|$)'),
    re.compile(r'Join-Path\s+\$target\s+"gitmap"([^-]|$)'),
    re.compile(r'\$[A-Za-z_.]*[Dd]eploy[A-Za-z_.]*[\\/]gitmap[\\/]'),
    re.compile(r'[\\/]gitmap[\\/]gitmap\.exe', re.IGNORECASE),
    re.compile(r'filepath\.Join\([^)]*[Dd]eploy[A-Za-z_]*,\s*"gitmap"\s*,'),
    re.compile(r'\$\{?[A-Z_]*DEPLOY[A-Z_]*\}?[/\\]gitmap[/\\]'),
]


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    scan_root = sys.argv[1] if len(sys.argv) > 1 else repo_root
    scan_root = os.path.abspath(scan_root)

    violations = []
    for root, dirs, files in os.walk(scan_root):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for f in files:
            ext = os.path.splitext(f)[1].lower()
            if ext in EXCLUDE_EXTS:
                continue
            full_path = os.path.join(root, f)
            rel_path = os.path.relpath(full_path, repo_root).replace("\\", "/")
            if rel_path in EXEMPT_FILES:
                continue

            try:
                with open(full_path, "r", encoding="utf-8", errors="ignore") as fh:
                    for lineno, line in enumerate(fh, start=1):
                        if "deploy-layout-allow" in line or "legacy" in line.lower() or "legacyappsubdirs" in line.lower():
                            continue
                        for pat in PATTERNS:
                            if pat.search(line):
                                violations.append((rel_path, lineno, line.strip()))
                                break
            except Exception:
                pass

    if violations:
        print(f"::error::Found {len(violations)} forbidden deploy path literal(s):", file=sys.stderr)
        for rel_path, lineno, text in violations:
            print(f"::error file={rel_path},line={lineno}::[deploy-layout] {text}", file=sys.stderr)
        sys.exit(1)

    print("Deploy layout check: OK (all deploy paths use gitmap-cli or deploy-manifest).")
    sys.exit(0)


if __name__ == "__main__":
    main()
