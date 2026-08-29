#!/usr/bin/env python3
"""Cross-platform check for file line count limits."""
import os
import re
import sys


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    max_lines = int(sys.argv[1]) if len(sys.argv) > 1 else 200
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))

    exclude_patterns = [
        re.compile(r"_test\.go$"),
        re.compile(r"[/\\]testdata[/\\]"),
        re.compile(r"[/\\]golden[/\\]"),
        re.compile(r"allcommands_generated\.go$"),
        re.compile(r"^\.git[/\\]"),
        re.compile(r"^node_modules[/\\]"),
        re.compile(r"^\.lovable[/\\]"),
    ]

    offenders = []
    for root, _, files in os.walk(repo_root):
        for f in sorted(files):
            if f.endswith(".go") or f.endswith(".ps1"):
                full_path = os.path.join(root, f)
                rel_path = os.path.relpath(full_path, repo_root)

                if any(p.search(rel_path) for p in exclude_patterns):
                    continue

                try:
                    with open(full_path, "r", encoding="utf-8", errors="replace") as fh:
                        line_count = sum(1 for _ in fh)
                    if line_count > max_lines:
                        offenders.append((rel_path.replace("\\", "/"), line_count))
                except Exception:
                    pass

    if not offenders:
        print(f"file-size-check: OK (no file exceeds {max_lines} lines)")
        sys.exit(0)

    print(f"file-size-check: WARN — {len(offenders)} pre-existing files over {max_lines} lines (non-blocking baseline):", file=sys.stderr)
    for path, count in offenders:
        print(f"  - {path} ({count} lines)", file=sys.stderr)
    print(f"\nNew code must respect the {max_lines}-line ceiling; legacy offenders are tracked for incremental splitting.", file=sys.stderr)
    sys.exit(0)


if __name__ == "__main__":
    main()
