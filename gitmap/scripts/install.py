#!/usr/bin/env python3
"""gitmap/scripts/install.py — Canonical cross-platform installer entrypoint.

Delegates directly to repository root install.py.
"""

import subprocess
import sys
from pathlib import Path

if __name__ == "__main__":
    root_installer = Path(__file__).resolve().parent.parent.parent / "install.py"
    if root_installer.is_file():
        sys.exit(subprocess.run([sys.executable, str(root_installer)] + sys.argv[1:]).returncode)
    else:
        print("Root install.py not found.", file=sys.stderr)
        sys.exit(1)
