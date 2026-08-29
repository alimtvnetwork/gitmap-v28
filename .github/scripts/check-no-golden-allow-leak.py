#!/usr/bin/env python3
"""Cross-platform check to prevent GITMAP_ALLOW_GOLDEN_UPDATE leaks in CI and scripts."""
import os
import re
import sys


WHITELIST_PREFIX = "gitmap/goldenguard/"

PATTERNS_BY_EXT = {
    (".sh", ".bash"): [
        re.compile(r"(?:^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE\s*="),
        re.compile(r"export\s+GITMAP_ALLOW_GOLDEN_UPDATE"),
    ],
    (".ps1", ".psm1"): [
        re.compile(r"\$env:GITMAP_ALLOW_GOLDEN_UPDATE"),
        re.compile(r"Set-Item.*GITMAP_ALLOW_GOLDEN_UPDATE"),
    ],
    (".yml", ".yaml"): [
        re.compile(r"^\s*GITMAP_ALLOW_GOLDEN_UPDATE\s*:"),
        re.compile(r"GITMAP_ALLOW_GOLDEN_UPDATE\s*="),
    ],
    (".mk", "makefile"): [
        re.compile(r"(?:^|[^A-Z_])GITMAP_ALLOW_GOLDEN_UPDATE\s*="),
        re.compile(r"export\s+GITMAP_ALLOW_GOLDEN_UPDATE"),
    ],
    ("dockerfile",): [
        re.compile(r"^\s*ENV\s+GITMAP_ALLOW_GOLDEN_UPDATE"),
    ],
    (".go",): [
        re.compile(r'(?:os|t|tt|child)\.Setenv\(\s*"GITMAP_ALLOW_GOLDEN_UPDATE"'),
        re.compile(r'\.Setenv\(\s*(?:goldenguard\.)?AllowUpdateEnv'),
    ],
}


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    scan_root = sys.argv[1] if len(sys.argv) > 1 else repo_root
    scan_root = os.path.abspath(scan_root)

    violations = []
    for root, _, files in os.walk(scan_root):
        for f in files:
            full_path = os.path.join(root, f)
            rel_path = os.path.relpath(full_path, repo_root).replace("\\", "/")

            if rel_path.startswith(WHITELIST_PREFIX) or rel_path.endswith("check-no-golden-allow-leak.py") or rel_path.endswith("check-no-golden-allow-leak.sh"):
                continue

            fname_lower = f.lower()
            ext = os.path.splitext(fname_lower)[1]

            applicable_patterns = []
            for exts, pats in PATTERNS_BY_EXT.items():
                if ext in exts or fname_lower in exts or (fname_lower.startswith("dockerfile") and "dockerfile" in exts) or (fname_lower.startswith("makefile") and "makefile" in exts):
                    applicable_patterns.extend(pats)

            if not applicable_patterns:
                continue

            try:
                with open(full_path, "r", encoding="utf-8", errors="ignore") as fh:
                    for lineno, line in enumerate(fh, start=1):
                        for pat in applicable_patterns:
                            if pat.search(line):
                                violations.append((rel_path, lineno, line.strip()))
                                break
            except Exception:
                pass

    if not violations:
        print("GITMAP_ALLOW_GOLDEN_UPDATE check: OK (no leaked update overrides).")
        sys.exit(0)

    print(f"::error::Found {len(violations)} illegal assignment(s) to GITMAP_ALLOW_GOLDEN_UPDATE outside goldenguard:", file=sys.stderr)
    for rel, lineno, text in violations:
        print(f"::error file={rel},line={lineno}::[golden-allow-leak] {text}", file=sys.stderr)

    sys.exit(1)


if __name__ == "__main__":
    main()
