#!/usr/bin/env python3
"""Git pull with optional permissions and user assignment (cross-platform)."""

import os
import subprocess
import sys
from pathlib import Path

def main() -> int:
    if len(sys.argv) > 1 and sys.argv[1] in ("-h", "--help"):
        print("Usage: python git_pull.py [permissions] [username]")
        return 0

    try:
        repo_root = subprocess.check_output(["git", "rev-parse", "--show-toplevel"], text=True).strip()
    except subprocess.SubprocessError:
        print("Not a git repository.", file=sys.stderr)
        return 1

    print(f"repo root: {repo_root}\n")
    res = subprocess.run(["git", "pull", "--ff-only"], cwd=repo_root)
    return res.returncode

if __name__ == "__main__":
    sys.exit(main())
