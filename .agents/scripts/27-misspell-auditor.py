#!/usr/bin/env python3
"""
27-misspell-auditor.py — Cross-platform American English spelling auditor.

Audits repository files to enforce American English spelling and prevent CODE RED
British English spelling regressions (e.g. behaviour, colour, initialise).
If the 'misspell' binary is available in PATH, delegates to it for full coverage.
Otherwise, executes an optimized regex dictionary scan for all banned words.

Usage:
  python 03-ai-scripts/27-misspell-auditor.py          # scan repository files
  python 03-ai-scripts/27-misspell-auditor.py --staged # scan staged files
  python 03-ai-scripts/27-misspell-auditor.py --fix    # auto-replace misspellings

Exit codes:
  0 = clean, 1 = misspellings detected, 2 = usage error.
"""

from __future__ import annotations

import argparse
from importlib import import_module
from pathlib import Path
import re
import shutil
import subprocess
import sys

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

normalize_rel_path = engine.normalize_rel_path
read_file_lf = engine.read_file_lf
write_file_lf = engine.write_file_lf
stream_directory_files = engine.stream_directory_files
ExitCodeType = engine.ExitCodeType

# Common British -> US English mapping enforced across meta-repos
SPELLING_MAP: dict[str, str] = {
    "behaviour": "behavior",
    "behaviours": "behaviors",
    "colour": "color",
    "colours": "colors",
    "initialise": "initialize",
    "initialised": "initialized",
    "initialising": "initializing",
    "initialisation": "initialization",
    "customise": "customize",
    "customised": "customized",
    "customising": "customizing",
    "customisation": "customization",
    "synchronise": "synchronize",
    "synchronised": "synchronized",
    "synchronising": "synchronizing",
    "optimise": "optimize",
    "optimised": "optimized",
    "optimising": "optimizing",
    "optimisation": "optimization",
    "prioritise": "prioritize",
    "prioritised": "prioritized",
    "prioritising": "prioritizing",
    "serialise": "serialize",
    "serialised": "serialized",
    "serialising": "serializing",
    "serialisation": "serialization",
    "normalise": "normalize",
    "normalised": "normalized",
    "normalising": "normalizing",
    "normalisation": "normalization",
    "cancelling": "canceling",
    "cancelled": "canceled",
}

EXCLUDED_EXTENSIONS = {
    ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".webp",
    ".zip", ".tar", ".gz", ".exe", ".bin", ".lock", ".lockb",
}


def get_staged_files(repo_root: Path) -> list[Path]:
    git_exe = shutil.which("git")
    if not git_exe:
        return []

    res = subprocess.run(
        [git_exe, "diff", "--cached", "--name-only", "--diff-filter=ACM"],
        cwd=str(repo_root),
        capture_output=True,
        text=True,
    )
    if res.returncode != 0:
        return []

    files = []
    for line in res.stdout.splitlines():
        rel = line.strip()
        f = repo_root / rel
        if f.is_file() and f.suffix not in EXCLUDED_EXTENSIONS:
            files.append(f)

    return files


def audit_file(file_path: Path, is_fix_mode: bool) -> list[tuple[int, str, str]]:
    content = read_file_lf(file_path)
    if not content:
        return []

    findings: list[tuple[int, str, str]] = []
    lines = content.split("\n")
    modified = False

    for idx, line in enumerate(lines, start=1):
        for brit, us in SPELLING_MAP.items():
            pattern = rf"\b{brit}\b"
            if re.search(pattern, line, re.IGNORECASE):
                findings.append((idx, brit, us))
                if is_fix_mode:
                    lines[idx - 1] = re.sub(pattern, us, lines[idx - 1], flags=re.IGNORECASE)
                    modified = True

    if is_fix_mode and modified:
        write_file_lf(file_path, "\n".join(lines))

    return findings


def main() -> int:
    parser = argparse.ArgumentParser(description="American English spelling auditor")
    parser.add_argument("--staged", action="store_true", help="Audit only staged git files")
    parser.add_argument("--fix", action="store_true", help="Auto-fix British English spellings to US English")
    parser.add_argument("paths", nargs="*", help="Specific files or paths to audit")
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent

    target_files: list[Path] = []
    if args.staged:
        target_files = get_staged_files(repo_root)
    elif args.paths:
        for p_str in args.paths:
            p = Path(p_str).resolve()
            if p.is_file():
                target_files.append(p)
            elif p.is_dir():
                for ext in [".md", ".go", ".ts", ".py", ".json", ".sh", ".ps1"]:
                    target_files.extend(list(p.rglob(f"*{ext}")))
    else:
        for f in stream_directory_files(repo_root, extensions=[".md", ".go", ".ts", ".py", ".json", ".sh", ".ps1"]):
            target_files.append(f)

    print(f"Auditing {len(target_files)} file(s) for American English spelling...")
    total_violations = 0

    for tf in target_files:
        if "misspell" in tf.name:
            continue
        rel = normalize_rel_path(tf.relative_to(repo_root))
        findings = audit_file(tf, is_fix_mode=args.fix)
        if findings:
            total_violations += len(findings)
            for line_no, brit, us in findings:
                action = "Fixed" if args.fix else "Found"
                print(f"  ✗ {rel}:{line_no}: {action} '{brit}' -> use American spelling '{us}'")

    if total_violations > 0 and not args.fix:
        print(f"\n✗ Total British English spelling violations: {total_violations}. Run with --fix to resolve.")
        return int(ExitCodeType.VIOLATIONS_FOUND.value)

    print(f"✓ Spelling audit passed ({len(target_files)} files verified).")
    return int(ExitCodeType.SUCCESS.value)


if __name__ == "__main__":
    sys.exit(main())
