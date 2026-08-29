#!/usr/bin/env python3
"""Cross-platform check for collision-prone helper names in gitmap/cmd/."""
import os
import re
import sys


PATTERNS = [
    (re.compile(r"^func\s+(invoke|persist)\s*\("), "bare verb (invoke/persist) without domain noun"),
    (re.compile(r"^func\s+runOne\s*\("), "bare runOne without domain noun"),
    (re.compile(r"^func\s+(invoke|persist|runOne)(Release|Task|Job|Item|All|One|Cmd)\s*\("), "over-generic noun"),
]


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    cmd_dir = sys.argv[1] if len(sys.argv) > 1 else os.path.join(repo_root, "gitmap", "cmd")

    if not os.path.isdir(cmd_dir):
        print(f"::error::cmd directory not found: {cmd_dir}", file=sys.stderr)
        sys.exit(1)

    violations = []
    for root, _, files in os.walk(cmd_dir):
        for f in sorted(files):
            if f.endswith(".go") and not f.endswith("_test.go"):
                path = os.path.join(root, f)
                rel_path = os.path.relpath(path, repo_root).replace("\\", "/")
                with open(path, "r", encoding="utf-8", errors="replace") as fh:
                    for lineno, line in enumerate(fh, start=1):
                        stripped = line.strip()
                        for pat, desc in PATTERNS:
                            if pat.search(stripped):
                                violations.append((rel_path, lineno, stripped, desc))

    if violations:
        print(f"::error::Found {len(violations)} collision-prone naming pattern(s) in {cmd_dir}", file=sys.stderr)
        print("::error::Rename helpers with a domain prefix that narrows the scope, e.g.:", file=sys.stderr)
        print("::error::  invokeRelease       -> invokeAliasRelease", file=sys.stderr)
        print("::error::  runOne              -> runOnePullJob, runOneScanRelease", file=sys.stderr)
        print("::error::  persistAll          -> persistReleaseToDB", file=sys.stderr)
        for rel, lineno, text, desc in violations:
            print(f"::error file={rel},line={lineno}::[{desc}] {text}", file=sys.stderr)
        sys.exit(1)

    print("All cmd/ helper names are domain-qualified.")
    sys.exit(0)


if __name__ == "__main__":
    main()
