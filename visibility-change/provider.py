#!/usr/bin/env python3
"""Provider and auth helpers for visibility-change (cross-platform)."""

from pathlib import Path
import sys

scripts_vc = Path(__file__).resolve().parent.parent / "scripts" / "visibility-change"
if str(scripts_vc) not in sys.path:
    sys.path.insert(0, str(scripts_vc))

from provider import *
