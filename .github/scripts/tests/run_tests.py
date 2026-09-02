#!/usr/bin/env python3
"""Run CI scripts test suite (cross-platform)."""

import subprocess
import sys
from pathlib import Path

if __name__ == "__main__":
    test_script = Path(__file__).resolve().parent / "test_ci_scripts.py"
    sys.exit(subprocess.run([sys.executable, str(test_script)] + sys.argv[1:]).returncode)
