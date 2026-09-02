#!/usr/bin/env python3
"""Path combining helpers (cross-platform)."""

from pathlib import Path
import os

def combine_path(path1: str, path2: str) -> str:
    p1 = path1.rstrip("/\\")
    p2 = path2.lstrip("/\\")
    return f"{p1}/{p2}"

def combine_with_base_path(source_path: str, dest_path: str) -> str:
    base_dir = os.path.basename(dest_path.rstrip("/\\"))
    return combine_path(source_path, base_dir)

if __name__ == "__main__":
    import sys
    if len(sys.argv) >= 3:
        print(combine_path(sys.argv[1], sys.argv[2]))
