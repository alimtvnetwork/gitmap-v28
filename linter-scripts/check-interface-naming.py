#!/usr/bin/env python3
"""
check-interface-naming.py - CI/CD Quality Gate Linter for Go Interface Naming.

Rule:
All named Go interface declarations MUST have 'er' as a suffix (e.g. BaseEnumer,
NumberEnumer, LogSinker, Writer, Streamer, WrappedByter, WrappedJsoner, StringCompiler).
"""

import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent

EXCLUDE_DIRS = {
    ".git", ".gitmap", "node_modules", "dist", "build", "bin", ".next",
    "vendor", "coverage", ".gemini", ".system_generated", "tests/fixtures",
    "scratch", "temp-scripts", "temp-agents", "temp", "tmp", "release-artifacts",
}

# Regex to detect interface declarations
# 1. Standalone: type Fooer interface { ... } or type Bar[T any] interface { ... }
PATTERN_STANDALONE = re.compile(r'^\s*type\s+([A-Za-z0-9_]+)(?:\[[^\]]+\])?\s+interface\s*\{')
# 2. Inside a type (...) block: Fooer interface { ... } or Bar[T any] interface { ... }
PATTERN_BLOCK_MEMBER = re.compile(r'^\s*([A-Za-z0-9_]+)(?:\[[^\]]+\])?\s+interface\s*\{')


def check_go_file(filepath: Path, root_dir: Path) -> list[str]:
    try:
        content = filepath.read_text(encoding="utf-8", errors="replace")
    except Exception as err:
        return [f"{filepath.as_posix()}:0 Cannot read file: {err}"]

    violations = []
    rel_path = filepath.relative_to(root_dir).as_posix()
    lines = content.splitlines()

    in_type_block = False

    for idx, line in enumerate(lines, 1):
        stripped = line.strip()

        # Skip comments
        if stripped.startswith(("//", "/*", "*")):
            continue

        # Track type ( ... ) blocks
        if stripped.startswith("type (") or stripped.startswith("type("):
            in_type_block = True
            continue
        if in_type_block and stripped.startswith(")"):
            in_type_block = False
            continue

        matched = PATTERN_STANDALONE.match(line)
        if not matched and in_type_block:
            matched = PATTERN_BLOCK_MEMBER.match(line)

        if matched:
            name = matched.group(1)
            if not name.endswith("er"):
                suggested = name + "er"
                violations.append(
                    f"{rel_path}:{idx} Interface '{name}' must have 'er' as suffix (suggested: '{suggested}')"
                )

    return violations


def main() -> int:
    target_root = ROOT_DIR
    if len(sys.argv) > 1 and not sys.argv[1].startswith("-"):
        target_root = Path(sys.argv[1]).resolve()

    all_violations: list[str] = []

    for root, dirs, files in os.walk(target_root):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for file in files:
            if file.endswith(".go") and not file.endswith("_test.go"):
                fp = Path(root) / file
                all_violations.extend(check_go_file(fp, target_root))

    if not all_violations:
        print(f"Scanned Go interfaces across {target_root.as_posix()}.")
        print("PASS: All Go interface definitions end with mandatory 'er' suffix.")
        return 0

    print(f"FAIL: Found {len(all_violations)} interface naming violation(s):")
    for v in all_violations:
        print(f"  {v}")
    return 1


if __name__ == "__main__":
    sys.exit(main())
