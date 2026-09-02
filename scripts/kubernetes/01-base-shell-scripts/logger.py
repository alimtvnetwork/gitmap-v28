#!/usr/bin/env python3
"""Colorful timestamped logging for Kubernetes shell scripts (cross-platform)."""

from pathlib import Path
import sys

helpers_dir = Path(__file__).resolve().parent.parent / "01-base-helpers"
if str(helpers_dir) not in sys.path:
    sys.path.insert(0, str(helpers_dir))

from logger import *
