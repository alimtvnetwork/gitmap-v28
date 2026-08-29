#!/usr/bin/env python3
"""check-markdown-header-spacing.py - Linter for Markdown heading spacing (MD022/MD032)."""
import subprocess
import sys
from pathlib import Path

ROOT_DIR = Path(__file__).resolve().parent.parent
SCRIPT = ROOT_DIR / "linter-scripts" / "check-markdown-headings.py"


def main():
    cmd = [sys.executable, str(SCRIPT)] + sys.argv[1:]
    res = subprocess.run(cmd, cwd=str(ROOT_DIR))
    return res.returncode


if __name__ == "__main__":
    sys.exit(main())
