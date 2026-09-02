#!/usr/bin/env python3
"""scripts/smoke_act.py — Run installer-smoke CI jobs locally via nektos/act (cross-platform).

Usage:
  python scripts/smoke_act.py
  python scripts/smoke_act.py <job-name>
"""

from __future__ import annotations

import shutil
import subprocess
import sys
from pathlib import Path


def main() -> int:
    job = sys.argv[1] if len(sys.argv) > 1 else "installer-smoke"

    if not shutil.which("act"):
        print("✗ act not found on PATH.", file=sys.stderr)
        print("  Install: https://github.com/nektos/act#installation", file=sys.stderr)
        return 2

    try:
        subprocess.check_call(["docker", "info"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except subprocess.SubprocessError:
        print("✗ Docker daemon is not reachable. act requires Docker to run.", file=sys.stderr)
        return 2

    repo_root = Path(__file__).resolve().parent.parent
    print(f"→ Running CI job '{job}' locally via act...")
    print("  (Windows-only jobs cannot run under act; use CI for those.)\n")

    cmd = [
        "act",
        "-j",
        job,
        "--pull=false",
        "-P",
        "ubuntu-latest=catthehacker/ubuntu:act-latest",
    ]
    res = subprocess.run(cmd, cwd=str(repo_root))
    return res.returncode


if __name__ == "__main__":
    sys.exit(main())
