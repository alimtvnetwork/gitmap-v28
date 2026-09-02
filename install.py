#!/usr/bin/env python3
"""install.py — Universal cross-platform installer for gitmap (Windows, macOS, Linux).

Usage:
  python install.py
  python install.py --version v6.164.0
  python install.py --dir /path/to/bin
  python install.py --arch amd64
  python install.py --no-path
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import platform
import shutil
import subprocess
import sys
import tarfile
import tempfile
import urllib.request
import zipfile
from pathlib import Path

REPO = "alimtvnetwork/gitmap-v28"
BINARY_NAME = "gitmap.exe" if os.name == "nt" else "gitmap"
APP_SUBDIR = "gitmap-cli"


def step(msg: str) -> None:
    print(f"  \033[36m{msg}\033[0m")


def ok(msg: str) -> None:
    print(f"  \033[32m{msg}\033[0m")


def warn(msg: str) -> None:
    print(f"  \033[33m{msg}\033[0m")


def err(msg: str) -> None:
    print(f"  \033[31m{msg}\033[0m", file=sys.stderr)


def detect_os() -> str:
    system = platform.system().lower()
    if system == "linux":
        return "linux"
    elif system == "darwin":
        return "darwin"
    elif system == "windows":
        return "windows"
    err(f"Unsupported OS: {system}")
    sys.exit(1)


def detect_arch(override: str | None) -> str:
    if override:
        return override
    machine = platform.machine().lower()
    if machine in ("x86_64", "amd64"):
        return "amd64"
    elif machine in ("arm64", "aarch64"):
        return "arm64"
    err(f"Unsupported architecture: {machine}")
    sys.exit(1)


def resolve_version(tag: str | None) -> str:
    if tag:
        return tag if tag.startswith("v") else f"v{tag}"
    step("Fetching latest release...")
    url = f"https://api.github.com/repos/{REPO}/releases/latest"
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "gitmap-installer"})
        token = os.getenv("GITHUB_TOKEN")
        if token:
            req.add_header("Authorization", f"Bearer {token}")
        with urllib.request.urlopen(req, timeout=10) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            return str(data.get("tag_name", ""))
    except Exception as e:
        err(f"Failed to fetch latest release: {e}")
        sys.exit(1)


def download_file(url: str, dest: Path) -> None:
    req = urllib.request.Request(url, headers={"User-Agent": "gitmap-installer"})
    token = os.getenv("GITHUB_TOKEN")
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    with urllib.request.urlopen(req, timeout=30) as resp, open(dest, "wb") as f:
        shutil.copyfileobj(resp, f)


def sha256_file(path: Path) -> str:
    h = hashlib.sha256()
    with open(path, "rb") as f:
        while chunk := f.read(65536):
            h.update(chunk)
    return h.hexdigest()


def verify_checksum(archive_path: Path, asset_name: str, checksum_path: Path) -> bool:
    if not checksum_path.is_file():
        warn("No checksums.txt available; skipping verification.")
        return True
    expected = ""
    with open(checksum_path, "r", encoding="utf-8", errors="replace") as f:
        for line in f:
            parts = line.strip().split()
            if len(parts) >= 2 and parts[1].endswith(asset_name):
                expected = parts[0]
                break
    if not expected:
        warn(f"{asset_name} not found in checksums.txt; skipping verification.")
        return True

    actual = sha256_file(archive_path)
    if actual.lower() != expected.lower():
        err(f"Checksum mismatch for {asset_name}!\n  Expected: {expected}\n  Got:      {actual}")
        return False
    ok("Checksum verified.")
    return True


def extract_archive(archive_path: Path, dest_dir: Path) -> Path | None:
    dest_dir.mkdir(parents=True, exist_ok=True)
    if archive_path.suffix == ".zip":
        with zipfile.ZipFile(archive_path, "r") as z:
            z.extractall(dest_dir)
    elif archive_path.name.endswith((".tar.gz", ".tgz")):
        with tarfile.open(archive_path, "r:gz") as t:
            t.extractall(dest_dir)

    bin_name = BINARY_NAME
    for root, _, files in os.walk(dest_dir):
        for f in files:
            if f.lower() == bin_name.lower():
                return Path(root) / f
    return None


def add_to_path(install_dir: Path) -> None:
    install_str = str(install_dir)
    if os.name == "nt":
        # Check current user PATH
        cur_path = os.environ.get("PATH", "")
        if install_str.lower() not in cur_path.lower():
            try:
                # Add to User Environment PATH via powershell
                ps_cmd = f'[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";{install_str}", "User")'
                subprocess.run(["powershell", "-NoProfile", "-Command", ps_cmd], check=False)
                ok(f"Added {install_dir} to User PATH.")
            except Exception:
                warn(f"Please add {install_dir} to your system PATH manually.")
    else:
        home = Path.home()
        rc_files = [home / ".bashrc", home / ".zshrc", home / ".profile"]
        for rc in rc_files:
            if rc.is_file():
                try:
                    with open(rc, "r", encoding="utf-8") as f:
                        content = f.read()
                    if install_str not in content:
                        with open(rc, "a", encoding="utf-8") as f:
                            f.write(f'\nexport PATH="$PATH:{install_str}"\n')
                        ok(f"Added to {rc}")
                except OSError:
                    pass


def main() -> int:
    parser = argparse.ArgumentParser(description="Cross-platform gitmap installer.")
    parser.add_argument("--version", help="Specific release tag to install")
    parser.add_argument("--dir", help="Install root directory")
    parser.add_argument("--arch", help="Architecture override (amd64, arm64)")
    parser.add_argument("--no-path", action="store_true", help="Skip modifying PATH")
    args = parser.parse_args()

    os_name = detect_os()
    arch_name = detect_arch(args.arch)
    version = resolve_version(args.version)

    default_root = Path.home() / "AppData" / "Local" if os.name == "nt" else Path.home() / ".local" / "bin"
    root_dir = Path(args.dir) if args.dir else default_root
    app_dir = root_dir / APP_SUBDIR
    app_dir.mkdir(parents=True, exist_ok=True)

    step(f"Target: {os_name}/{arch_name} · Release: {version}")

    ext = ".zip" if os_name == "windows" else ".tar.gz"
    asset_name = f"gitmap-{version}-{os_name}-{arch_name}{ext}"
    base_url = f"https://github.com/{REPO}/releases/download/{version}"
    asset_url = f"{base_url}/{asset_name}"
    checksum_url = f"{base_url}/checksums.txt"

    with tempfile.TemporaryDirectory() as tmp_dir:
        tmp_path = Path(tmp_dir)
        archive_path = tmp_path / asset_name
        checksum_path = tmp_path / "checksums.txt"

        step(f"Downloading {asset_name}...")
        try:
            download_file(asset_url, archive_path)
        except Exception as e:
            err(f"Download failed: {e}")
            return 1

        try:
            download_file(checksum_url, checksum_path)
        except Exception:
            pass

        if not verify_checksum(archive_path, asset_name, checksum_path):
            return 1

        step(f"Extracting to {app_dir}...")
        bin_found = extract_archive(archive_path, tmp_path / "extract")
        if not bin_found or not bin_found.is_file():
            err("Could not find gitmap binary in downloaded archive.")
            return 1

        dest_bin = app_dir / BINARY_NAME
        shutil.copy2(bin_found, dest_bin)
        try:
            dest_bin.chmod(0o755)
        except OSError:
            pass
        ok(f"Installed {BINARY_NAME} to {app_dir}")

        # Install alias `gm`
        alias_name = "gm.exe" if os.name == "nt" else "gm"
        alias_dest = app_dir / alias_name
        try:
            shutil.copy2(dest_bin, alias_dest)
            ok(f"Installed alias {alias_name}")
        except OSError:
            pass

    if not args.no_path:
        add_to_path(app_dir)

    # Run gitmap setup
    if dest_bin.is_file():
        step("Running initial setup...")
        subprocess.run([str(dest_bin), "setup"], check=False)

    ok(f"Installation of gitmap {version} complete!")
    return 0


if __name__ == "__main__":
    sys.exit(main())
