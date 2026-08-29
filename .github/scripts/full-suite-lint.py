#!/usr/bin/env python3
"""Cross-platform runner for full suite golangci-lint."""
import os
import re
import subprocess
import sys


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    gitmap_dir = sys.argv[1] if len(sys.argv) > 1 else os.path.join(repo_root, "gitmap")

    cmd = [
        "golangci-lint", "run", "./...",
        "--timeout=5m",
        "--max-issues-per-linter=0",
        "--max-same-issues=0"
    ]

    try:
        res = subprocess.run(
            cmd,
            cwd=gitmap_dir,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace"
        )
    except FileNotFoundError:
        print("ERROR: golangci-lint not on PATH", file=sys.stderr)
        sys.exit(2)

    output = (res.stdout or "") + "\n" + (res.stderr or "")
    print(output.strip())

    issue_count = sum(1 for line in output.splitlines() if re.match(r'^\S+:[0-9]+:[0-9]+:', line.strip()))

    github_output = os.environ.get("GITHUB_OUTPUT")
    if github_output and os.path.isfile(github_output):
        with open(github_output, "a", encoding="utf-8") as fh:
            fh.write(f"exit_code={res.returncode}\n")
            fh.write(f"issue_count={issue_count}\n")

    sys.exit(res.returncode)


if __name__ == "__main__":
    main()
