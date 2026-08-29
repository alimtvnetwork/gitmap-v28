#!/usr/bin/env python3
"""Cross-platform check for forbidden legacy version strings."""
import os
import re
import sys


EXCLUDE_DIRS = {
    ".git", "node_modules", "dist", "build", "bin", ".next",
    ".gitmap", "vendor", "coverage"
}

EXCLUDE_EXTS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico",
    ".pdf", ".zip", ".gz", ".tar", ".exe", ".dll", ".so", ".dylib",
    ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf"
}

EXEMPT_FILES = {
    ".github/scripts/check-legacy-refs.sh",
    ".github/scripts/check-legacy-refs.py",
}


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    pattern_str = os.environ.get("LEGACY_PATTERN", r"gitmap-v[567]\b")
    pattern = re.compile(pattern_str)

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    scan_root = sys.argv[1] if len(sys.argv) > 1 else repo_root
    scan_root = os.path.abspath(scan_root)

    print(f"  [legacy-refs] scanning '{scan_root}' for pattern: {pattern_str}")

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
                        if "gitmap-legacy-ref-allow" in line:
                            continue
                        if pattern.search(line):
                            violations.append((rel_path, lineno, line.strip()))
            except Exception:
                pass

    if not violations:
        print("  [legacy-refs] OK — no forbidden legacy refs found.")
        sys.exit(0)

    print(f"::error::Found {len(violations)} legacy reference(s) matching /{pattern_str}/", file=sys.stderr)
    for rel, lineno, text in violations:
        print(f"::error file={rel},line={lineno}::[legacy-ref] {text}", file=sys.stderr)
        print(f"    {rel}:{lineno}: {text}", file=sys.stderr)

    sys.exit(1)


if __name__ == "__main__":
    main()
