#!/usr/bin/env python3
"""
Script 22: Doc & Prompt Path Integrity Linter
Autonomously verifies all markdown links and path references across 01-prompts/,
spec/, .lovable/, and .agents/ to guarantee zero dead path references and zero hallucinations.
"""

import os
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(REPO_ROOT))


def run() -> int:
    linter_path = REPO_ROOT / 'linter-scripts' / 'check-prompt-and-spec-paths.py'
    if not linter_path.exists():
        print(f'Error: Linter not found at {linter_path}')
        return 1

    import subprocess
    res = subprocess.run([sys.executable, str(linter_path)], cwd=str(REPO_ROOT))
    return res.returncode


if __name__ == '__main__':
    sys.exit(run())
