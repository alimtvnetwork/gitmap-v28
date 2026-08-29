#!/usr/bin/env python3
"""Cross-platform check for bare `fmt.Fprintln(os.Stderr, err)` in gitmap/cmd/."""
import os
import re
import sys


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    cmd_dir = os.path.join(repo_root, "gitmap", "cmd")
    pattern = re.compile(r"fmt\.Fprintln\(os\.Stderr,\s*err\)")

    offenders = []
    for root, _, files in os.walk(cmd_dir):
        for f in sorted(files):
            if f.endswith(".go") and not f.endswith("_test.go"):
                path = os.path.join(root, f)
                with open(path, "r", encoding="utf-8", errors="replace") as fh:
                    for lineno, line in enumerate(fh, start=1):
                        if pattern.search(line):
                            rel_path = os.path.relpath(path, repo_root).replace("\\", "/")
                            offenders.append((rel_path, lineno, line.strip()))

    if not offenders:
        print("OK: no bare 'fmt.Fprintln(os.Stderr, err)' in gitmap/cmd/")
        sys.exit(0)

    print("FAIL: bare error prints found in gitmap/cmd/ — use cliexit.Reportf/Fail", file=sys.stderr)
    print("", file=sys.stderr)
    for rel_path, lineno, line_content in offenders:
        print(f"::error file={rel_path},line={lineno}::[bare-err] use cliexit.Reportf(cmd, op, subject, err) instead", file=sys.stderr)
        print(f"  {rel_path}:{lineno}: {line_content}", file=sys.stderr)

    sys.exit(1)


if __name__ == "__main__":
    main()
