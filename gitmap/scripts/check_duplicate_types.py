#!/usr/bin/env python3
"""gitmap/scripts/check_duplicate_types.py — Detects duplicate type declarations across files in the same Go package (cross-platform).

Exits code 1 if duplicates are found, 0 otherwise.
"""

from __future__ import annotations

import re
import sys
from collections import defaultdict
from pathlib import Path


def main() -> int:
    if sys.platform == "win32":
        try:
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
            sys.stderr.reconfigure(encoding="utf-8", errors="replace")
        except Exception:
            pass

    root = Path(sys.argv[1]) if len(sys.argv) > 1 else Path(".")
    type_decl_re = re.compile(r"^type\s+([A-Z][A-Za-z0-9_]*)\s+")

    # Group go files by package directory
    pkg_files: dict[Path, list[Path]] = defaultdict(list)
    for go_file in root.rglob("*.go"):
        posix = go_file.as_posix()
        if "/vendor/" in posix or go_file.name.endswith("_test.go"):
            continue
        pkg_files[go_file.parent].append(go_file)

    found_duplicate = False

    for pkg_dir, files in sorted(pkg_files.items()):
        # Map: (type_name, build_tag) -> filename
        # If files have mutually exclusive build tags (e.g. windows vs !windows), they are allowed.
        type_tags: dict[str, list[tuple[str, str]]] = defaultdict(list) # type_name -> list of (file.name, tag)

        for f in files:
            build_tag = ""
            try:
                with open(f, "r", encoding="utf-8", errors="replace") as fh:
                    for line in fh:
                        line_s = line.strip()
                        if line_s.startswith("//go:build"):
                            build_tag = line_s.replace("//go:build", "").strip()
                        m = type_decl_re.match(line_s)
                        if m:
                            type_name = m.group(1)
                            type_tags[type_name].append((f.name, build_tag))
            except OSError:
                pass

        pkg_duplicates: list[tuple[str, str, str]] = []
        for tname, occurrences in type_tags.items():
            if len(occurrences) <= 1:
                continue
            # Check if any two occurrences share the same build tag or have no build tag
            for i in range(len(occurrences)):
                for j in range(i + 1, len(occurrences)):
                    f1, tag1 = occurrences[i]
                    f2, tag2 = occurrences[j]
                    if f1 == f2:
                        continue
                    # If tags are mutually exclusive (e.g. windows vs !windows), no conflict
                    is_exclusive = (
                        (tag1 == "windows" and tag2 == "!windows") or
                        (tag1 == "!windows" and tag2 == "windows") or
                        (tag1 and tag2 and (tag1 == f"!{tag2}" or tag2 == f"!{tag1}"))
                    )
                    if not is_exclusive:
                        pkg_duplicates.append((tname, f1, f2))

        if pkg_duplicates:
            found_duplicate = True
            print(f"✗ Duplicate type(s) found in package: {pkg_dir}")
            for tname, first_f, second_f in pkg_duplicates:
                print(f"  DUPLICATE: type {tname}")
                print(f"    -> {first_f}")
                print(f"    -> {second_f}")

    if not found_duplicate:
        print("✓ No duplicate type declarations found.")
        return 0

    print("\nFix: rename one declaration or move it to a shared file.", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
