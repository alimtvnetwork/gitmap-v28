#!/usr/bin/env python3
"""File-traversal and binary-detection helpers for fix-repo (cross-platform)."""

from __future__ import annotations

from pathlib import Path

MAX_FILE_BYTES = 5 * 1024 * 1024
BINARY_EXTENSIONS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip",
    ".tar", ".gz", ".tgz", ".bz2", ".xz", ".7z", ".rar", ".woff",
    ".woff2", ".ttf", ".otf", ".eot", ".mp3", ".mp4", ".mov", ".wav",
    ".ogg", ".webm", ".class", ".jar", ".so", ".dylib", ".dll", ".exe",
    ".pyc", ".db",
}


def is_scannable_file(path: Path) -> bool:
    if path.is_symlink():
        return False
    if path.suffix.lower() in BINARY_EXTENSIONS:
        return False
    try:
        st = path.stat()
        if st.st_size > MAX_FILE_BYTES:
            return False
        with open(path, "rb") as f:
            chunk = f.read(8192)
            if b"\x00" in chunk:
                return False
    except OSError:
        return False
    return True
