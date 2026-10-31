#!/usr/bin/env python3
"""Multi-Repository Prompts, Skills, and Scripts Synchronizer — V7 Master Engine Version.

Executes fleet synchronization adhering to V7 Master Engine time restrictions,
morning commit cutoffs (< 08:00 AM on weekdays), autonomous bot identity verification,
and auto-window scheduling.

When running under V7 parent tasks:
  python 03-ai-scripts/38-sync-prompts-skills-scripts-v7.py [--tz <tz>] [--repo <name>]

When running under V6 parent tasks:
  python 03-ai-scripts/38-sync-prompts-skills-scripts.py [--repo <name>]
"""
from __future__ import annotations

import os
from pathlib import Path
import subprocess
import sys

SCRIPT_DIR = Path(__file__).resolve().parent
BASE_SYNC_SCRIPT = SCRIPT_DIR / "38-sync-prompts-skills-scripts.py"


def main() -> None:
    # Ensure --v7, --auto-window, and --schedule are active by default for V7 execution
    passthrough_args = list(sys.argv[1:])

    has_mode = any(arg in passthrough_args for arg in ("--v7", "--mode"))
    if not has_mode:
        passthrough_args.insert(0, "--v7")

    has_window = any(arg in passthrough_args for arg in ("--auto-window", "--schedule"))
    if not has_window:
        passthrough_args.extend(["--auto-window", "--schedule"])

    cmd = [sys.executable, str(BASE_SYNC_SCRIPT)] + passthrough_args
    print(f"[V7 SYNC WRAPPER] Executing: {' '.join(cmd)}")
    result = subprocess.run(cmd, check=False)
    sys.exit(result.returncode)


if __name__ == "__main__":
    main()
