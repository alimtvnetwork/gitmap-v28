#!/usr/bin/env python3
"""setup.py — One-time project setup after cloning (cross-platform).

Usage:
  python setup.py
"""

from __future__ import annotations

import os
import shutil
import subprocess
import sys
from pathlib import Path


def main() -> int:
    repo_root = Path(__file__).resolve().parent
    hook_src = repo_root / "hooks" / "pre-commit"
    hook_dst = repo_root / ".git" / "hooks" / "pre-commit"

    print("Setting up gitmap development environment...")

    # 1. Install pre-commit hook
    if hook_src.is_file():
        hook_dst.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(hook_src, hook_dst)
        try:
            hook_dst.chmod(0o755)
        except OSError:
            pass
        print("  ✓ Pre-commit hook installed")
    else:
        print(f"  ✗ Hook source not found: {hook_src}")

    # 2. Verify Go toolchain
    go_exe = shutil.which("go")
    if go_exe:
        try:
            ver_out = subprocess.check_output([go_exe, "version"], text=True).strip()
            parts = ver_out.split()
            ver_str = parts[2] if len(parts) >= 3 else ver_out
            print(f"  ✓ Go {ver_str}")
        except subprocess.SubprocessError:
            print("  ✓ Go found on PATH")
    else:
        print("  ⚠ Go not found — install from https://go.dev/dl/")

    # 3. Check or install golangci-lint
    lint_exe = shutil.which("golangci-lint")
    if lint_exe:
        print("  ✓ golangci-lint available")
    elif go_exe:
        print("  → Installing golangci-lint...")
        res = subprocess.run(
            [go_exe, "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8"]
        )
        if res.returncode == 0:
            print("  ✓ golangci-lint installed")
        else:
            print("  ⚠ Failed to install golangci-lint")
    else:
        print("  ⚠ golangci-lint missing and Go not available to install it")

    # 4. Download Go dependencies
    go_mod = repo_root / "gitmap" / "go.mod"
    if go_mod.is_file() and go_exe:
        print("  → Downloading Go dependencies...")
        res = subprocess.run([go_exe, "mod", "download"], cwd=str(repo_root / "gitmap"))
        if res.returncode == 0:
            print("  ✓ Dependencies ready")
        else:
            print("  ⚠ go mod download failed")

    print("\nDone! Run 'cd gitmap && go test ./...' to verify.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
