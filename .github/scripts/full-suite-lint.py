#!/usr/bin/env python3
"""Cross-platform runner for full suite golangci-lint with real-time streaming."""
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

    issue_count = 0
    issue_regex = re.compile(r'^\S+:[0-9]+:[0-9]+:')

    try:
        proc = subprocess.Popen(
            cmd,
            cwd=gitmap_dir,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            encoding="utf-8",
            errors="replace",
            bufsize=1
        )
    except FileNotFoundError:
        print("ERROR: golangci-lint not on PATH", file=sys.stderr)
        sys.exit(2)

    if proc.stdout:
        for line in proc.stdout:
            sys.stdout.write(line)
            sys.stdout.flush()
            if issue_regex.match(line.strip()):
                issue_count += 1

    proc.wait()

    github_output = os.environ.get("GITHUB_OUTPUT")
    if github_output and os.path.isfile(github_output):
        with open(github_output, "a", encoding="utf-8") as fh:
            fh.write(f"exit_code={proc.returncode}\n")
            fh.write(f"issue_count={issue_count}\n")

    sys.exit(proc.returncode)


if __name__ == "__main__":
    main()
