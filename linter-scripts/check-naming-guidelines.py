#!/usr/bin/env python3
"""check-naming-guidelines.py - Linter for bare 'ok' identifiers, negative booleans, and explicit comparisons."""
import subprocess
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent
SCRIPT = ROOT_DIR / ".lovable" / "ai-fix-scripts" / "05-naming-autofixer.py"


def main():
    if not SCRIPT.exists():
        print(f"Error: {SCRIPT} not found.")
        sys.exit(1)

    result = subprocess.run([sys.executable, str(SCRIPT)], capture_output=True, text=True, cwd=str(ROOT_DIR))
    print(result.stdout)
    if result.stderr:
        print(result.stderr, file=sys.stderr)
    sys.exit(result.returncode)


if __name__ == "__main__":
    main()
