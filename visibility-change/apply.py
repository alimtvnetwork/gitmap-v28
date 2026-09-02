#!/usr/bin/env python3
"""Apply and verify visibility change via gh / glab (cross-platform)."""

from pathlib import Path
import sys

scripts_vc = Path(__file__).resolve().parent.parent / "scripts" / "visibility-change"
if str(scripts_vc) not in sys.path:
    sys.path.insert(0, str(scripts_vc))

from apply import *
