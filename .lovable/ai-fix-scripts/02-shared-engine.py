#!/usr/bin/env python3
"""02-shared-engine.py - Shared Engine & Core Utilities for AI Fix Scripts.

Provides high-performance process execution, worker pooling, path normalization,
Git CLI helpers, and structured terminal reporting used across all ai-fix-scripts.
"""
import concurrent.futures
import os
import platform
import subprocess
import sys
import time
from typing import Callable, List, Optional, Tuple, Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
DEFAULT_IGNORES = [".git", "node_modules", "vendor", "__pycache__", "dist"]


def normalize_path(path_str: str) -> str:
    """Returns normalized absolute path with Windows long-path prefix if needed."""
    p = os.path.abspath(path_str)
    if platform.system() == "Windows" and len(p) > 248 and not p.startswith("\\\\?\\"):
        return "\\\\?\\" + p
    return p


def to_relative_path(abs_path: str, base_dir: Optional[str] = None) -> str:
    """Converts an absolute path to forward-slash relative path."""
    base = base_dir or ROOT_DIR
    try:
        rel = os.path.relpath(abs_path, base)
        return rel.replace("\\", "/")
    except ValueError:
        return abs_path.replace("\\", "/")


def run_cmd(
    name: str,
    cmd_str: str,
    cwd: Optional[str] = None,
    timeout_sec: int = 300,
) -> Tuple[str, str, int, str, str, float]:
    """Runs a shell command with real-time timing, timeout protection, and utf-8 decoding."""
    work_dir = cwd or ROOT_DIR
    start = time.monotonic()
    try:
        res = subprocess.run(
            cmd_str,
            cwd=work_dir,
            shell=True,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            timeout=timeout_sec,
        )
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, res.returncode, res.stdout, res.stderr, elapsed
    except subprocess.TimeoutExpired as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, -1, e.stdout or "", f"Timeout after {timeout_sec}s", elapsed
    except Exception as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, 1, "", str(e), elapsed


def run_git_mv(src: str, dst: str, cwd: Optional[str] = None) -> bool:
    """Attempts to rename a file using git mv to preserve version history."""
    work_dir = cwd or ROOT_DIR
    try:
        res = subprocess.run(
            ["git", "mv", src, dst],
            cwd=work_dir,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
        return res.returncode == 0
    except Exception:
        return False


def run_in_pool(
    fn: Callable[[Any], Any],
    items: List[Any],
    max_workers: int = 4,
) -> List[Any]:
    """Executes a worker function across items concurrently with bounded parallelism."""
    results = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(fn, item) for item in items]
        for future in concurrent.futures.as_completed(futures):
            try:
                results.append(future.result())
            except Exception as e:
                results.append(e)
    return results


if __name__ == "__main__":
    print(f"✅ 02-shared-engine loaded successfully. Root directory: {ROOT_DIR}")
