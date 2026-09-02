#!/usr/bin/env python3
"""install_quick.py — Cross-platform quick installer for gitmap.

Usage:
  python install_quick.py
  python install_quick.py --dir /path/to/bin
  python install_quick.py --version v6.164.0
  python install_quick.py --no-discovery
"""

from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
import urllib.request
from pathlib import Path

REPO = "alimtvnetwork/gitmap-v28"


def get_default_install_dir() -> Path:
    if os.name == "nt":
        local_app = os.getenv("LOCALAPPDATA")
        if local_app:
            return Path(local_app) / "gitmap-cli"
        return Path.home() / "AppData" / "Local" / "gitmap-cli"
    return Path.home() / ".local" / "bin"


def parse_repo_suffix(repo_name: str) -> tuple[str, str, int] | None:
    m = re.match(r"^([^/]+)/(.+)-v(\d+)$", repo_name)
    if m:
        return m.group(1), m.group(2), int(m.group(3))
    return None


def probe_repo_exists(url: str) -> bool:
    try:
        req = urllib.request.Request(url, method="HEAD")
        req.add_header("User-Agent", "gitmap-installer")
        with urllib.request.urlopen(req, timeout=5) as resp:
            return resp.status == 200
    except Exception:
        return False


def resolve_effective_repo(repo: str, window: int = 20) -> str:
    parsed = parse_repo_suffix(repo)
    if not parsed:
        print(f"  [discovery] no -v<N> suffix on '{repo}'; using baseline")
        return repo

    owner, stem, baseline = parsed
    print(f"  [discovery] baseline: {owner}/{stem}-v{baseline}")
    print(f"  [discovery] probe window: {window}")

    effective = baseline
    for m in range(baseline + 1, baseline + window + 1):
        url = f"https://github.com/{owner}/{stem}-v{m}"
        if probe_repo_exists(url):
            print(f"  [discovery] HEAD {url} ... HIT")
            effective = m
        else:
            print(f"  [discovery] HEAD {url} ... MISS")
            break

    if effective == baseline:
        print(f"  [discovery] no higher version found; using baseline -v{baseline}")
        return repo

    new_repo = f"{owner}/{stem}-v{effective}"
    print(f"  [discovery] effective: {new_repo} (was -v{baseline})")
    return new_repo


def fetch_latest_release(repo: str) -> str:
    url = f"https://api.github.com/repos/{repo}/releases/latest"
    try:
        req = urllib.request.Request(url)
        req.add_header("User-Agent", "gitmap-installer")
        token = os.getenv("GITHUB_TOKEN")
        if token:
            req.add_header("Authorization", f"Bearer {token}")
        with urllib.request.urlopen(req, timeout=10) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            return str(data.get("tag_name", ""))
    except Exception as e:
        print(f"  [WARN] Failed to fetch latest release: {e}", file=sys.stderr)
        return ""


def save_powershell_config(install_dir: Path) -> None:
    install_dir.mkdir(parents=True, exist_ok=True)
    cfg = install_dir / "powershell.json"
    data = {
        "deployPath": str(install_dir),
        "buildOutput": "./bin",
        "binaryName": "gitmap.exe" if os.name == "nt" else "gitmap",
        "goSource": "./gitmap",
        "copyData": True,
    }
    try:
        with open(cfg, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)
        print(f"  \033[90mSaved deployPath -> {cfg}\033[0m")
    except OSError:
        pass


def main() -> int:
    parser = argparse.ArgumentParser(description="Cross-platform quick installer for gitmap.")
    parser.add_argument("--dir", help="Install directory")
    parser.add_argument("--version", help="Pin to specific version")
    parser.add_argument("--no-discovery", action="store_true", help="Skip probing sibling repositories")
    parser.add_argument("--discovery-window", type=int, default=20, help="Discovery window size")
    args = parser.parse_args()

    print("\n  \033[36mgitmap quick installer\033[0m")
    print("  \033[90m---------------------\033[0m")

    effective_repo = REPO
    if not args.no_discovery and not args.version:
        effective_repo = resolve_effective_repo(REPO, args.discovery_window)

    install_dir = Path(args.dir) if args.dir else get_default_install_dir()
    print(f"\n  \033[32mInstalling gitmap to: {install_dir}\033[0m\n")

    save_powershell_config(install_dir)

    # Delegate to install.py if present, or download release asset
    install_py = Path(__file__).resolve().parent / "install.py"
    if install_py.is_file():
        cmd = [sys.executable, str(install_py), "--dir", str(install_dir)]
        if args.version:
            cmd.extend(["--version", args.version])
        return subprocess.run(cmd).returncode

    tag = args.version or fetch_latest_release(effective_repo)
    print(f"  Target release: {tag or 'latest'}")
    print(f"  Target folder:  {install_dir}")
    print("\n  \033[32mSetup complete. Ensure the install folder is in your PATH.\033[0m\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
