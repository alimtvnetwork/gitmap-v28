#!/usr/bin/env python3
"""Cross-platform misspell baseline-diff runner and gate."""
import os
import sys
import subprocess

def main():
    os.environ.setdefault("LINTER", "misspell")
    if "CURRENT_OUT" not in os.environ:
        os.environ["CURRENT_OUT"] = os.path.join(os.environ.get("TEMP", "/tmp"), "lint-misspell-current", "report.json")
    script_path = os.path.join(os.path.dirname(__file__), "check-single-linter-diff.py")
    cmd = [sys.executable, script_path] + sys.argv[1:]
    res = subprocess.run(cmd)
    sys.exit(res.returncode)

if __name__ == "__main__":
    main()
