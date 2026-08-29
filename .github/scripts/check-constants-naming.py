#!/usr/bin/env python3
"""Cross-platform check for constants naming conventions in gitmap/constants/."""
import argparse
import os
import re
import sys


def extract_constants_from_file(filepath):
    """Extracts top-level constant names defined in a Go source file.

    Handles single-line consts, block consts, and skips raw string literals / comments.
    """
    with open(filepath, "r", encoding="utf-8", errors="replace") as fh:
        lines = fh.readlines()

    const_names = []
    in_const = False
    in_rawstr = False

    for line in lines:
        raw_line = line
        num_backticks = raw_line.count("`")
        if num_backticks % 2 == 1:
            in_rawstr = not in_rawstr
            continue
        if in_rawstr:
            continue

        clean = raw_line.strip()
        if re.match(r"^const\s*\(", clean):
            in_const = True
            continue
        if in_const and clean.startswith(")"):
            in_const = False
            continue

        if clean.startswith("const "):
            rest = clean[6:].strip()
            m = re.match(r"^([A-Z][A-Za-z0-9_]*)", rest)
            if m:
                const_names.append(m.group(1))
            continue

        if in_const:
            if clean.startswith("//") or not clean:
                continue
            m = re.match(r"^([A-Z][A-Za-z0-9_]*)", clean)
            if m:
                name = m.group(1)
                after = clean[len(name):].strip()
                # Must be followed by =, a type, comment, or end of line (iota)
                if not after or after.startswith("=") or after.startswith("//") or re.match(r"^[A-Za-z0-9_.]+\s*(?:=|$|//)", after):
                    const_names.append(name)

    return const_names


def extract_all_constants(const_dir):
    all_names = set()
    for root, _, files in os.walk(const_dir):
        for f in sorted(files):
            if f.endswith(".go") and not f.endswith("_test.go"):
                full_path = os.path.join(root, f)
                names = extract_constants_from_file(full_path)
                all_names.update(names)
    return sorted(all_names)


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    parser = argparse.ArgumentParser(description="Check constants naming against baseline")
    parser.add_argument("--const-dir", default=os.environ.get("CONST_DIR", "gitmap/constants"))
    parser.add_argument("--baseline", default=os.environ.get("BASELINE_FILE", ".github/scripts/constants-baseline.txt"))
    parser.add_argument("--regenerate-baseline", action="store_true", help="Regenerate the baseline file")
    args = parser.parse_args()

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    const_dir = os.path.abspath(os.path.join(repo_root, args.const_dir))
    baseline_path = os.path.abspath(os.path.join(repo_root, args.baseline))

    if not os.path.isdir(const_dir):
        print(f"::error::constants directory not found: {const_dir}", file=sys.stderr)
        sys.exit(1)

    current_constants = extract_all_constants(const_dir)

    if args.regenerate_baseline:
        os.makedirs(os.path.dirname(baseline_path), exist_ok=True)
        with open(baseline_path, "w", encoding="utf-8") as fh:
            for name in current_constants:
                fh.write(name + "\n")
        print(f"Regenerated baseline at {baseline_path} ({len(current_constants)} constants).")
        sys.exit(0)

    if not os.path.isfile(baseline_path):
        print(f"WARN: baseline file {baseline_path} not found; treating all current constants as baseline.")
        sys.exit(0)

    with open(baseline_path, "r", encoding="utf-8", errors="replace") as fh:
        baseline_constants = set(line.strip() for line in fh if line.strip() and not line.startswith("#"))

    allowed_re = re.compile(r"^[A-Z][A-Za-z0-9_]*$")
    violations = []

    for name in current_constants:
        if name in baseline_constants:
            continue
        if not allowed_re.match(name):
            violations.append(name)

    if violations:
        print(f"::error::Found {len(violations)} newly added constant(s) violating PascalCase identifier naming:", file=sys.stderr)
        for v in violations:
            print(f"::error::  {v}", file=sys.stderr)
        sys.exit(1)

    print(f"All {len(current_constants)} constants pass naming checks.")
    sys.exit(0)


if __name__ == "__main__":
    main()
