#!/usr/bin/env python3
"""gitmap/scripts/release_version.py — version-pinned gitmap installer (cross-platform).

Installs EXACTLY the version requested via --version. Never resolves
'latest', never auto-upgrades, never silently substitutes.

Spec: spec/01-app/105-release-version-script.md

Usage:
  python gitmap/scripts/release_version.py --version v3.36.0
  python gitmap/scripts/release_version.py --version v3.36.0 --dir /opt/gitmap
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import platform
import re
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

EXIT_OK = 0
EXIT_VERSION_MISSING = 1
EXIT_NETWORK = 2
EXIT_CHECKSUM = 3
EXIT_UNSUPPORTED_ARCH = 4
EXIT_PATH_FAIL = 5
EXIT_SELF_INSTALL = 6
EXIT_VERIFY = 7


def log_step(msg: str, quiet: bool) -> None:
    if not quiet:
        print(f"  -> {msg}", file=sys.stderr)


def log_ok(msg: str, quiet: bool) -> None:
    if not quiet:
        print(f"  OK {msg}", file=sys.stderr)


def log_warn(msg: str, quiet: bool) -> None:
    if not quiet:
        print(f"  !  {msg}", file=sys.stderr)


def log_err(msg: str) -> None:
    print(f"  X  {msg}", file=sys.stderr)


def detect_os() -> str:
    s = platform.system().lower()
    if s == "linux":
        return "linux"
    elif s == "darwin":
        return "darwin"
    elif s == "windows":
        return "windows"
    log_err(f"Unsupported OS: {s}")
    sys.exit(EXIT_UNSUPPORTED_ARCH)


def detect_arch(override: str | None) -> str:
    if override:
        if override in ("amd64", "arm64"):
            return override
        log_err(f"Unsupported --arch: {override}")
        sys.exit(EXIT_UNSUPPORTED_ARCH)
    m = platform.machine().lower()
    if m in ("x86_64", "amd64"):
        return "amd64"
    elif m in ("arm64", "aarch64"):
        return "arm64"
    log_err(f"Unsupported architecture: {m}")
    sys.exit(EXIT_UNSUPPORTED_ARCH)


def http_get_json(url: str) -> dict | None:
    try:
        req = urllib.request.Request(url)
        req.add_header("User-Agent", "gitmap-release-version-installer")
        req.add_header("Accept", "application/vnd.github+json")
        token = os.getenv("GITHUB_TOKEN")
        if token:
            req.add_header("Authorization", f"Bearer {token}")
        with urllib.request.urlopen(req, timeout=10) as resp:
            return json.loads(resp.read().decode("utf-8"))
    except Exception:
        return None


def download_file(url: str, dest: Path) -> None:
    req = urllib.request.Request(url)
    req.add_header("User-Agent", "gitmap-release-version-installer")
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


def main() -> int:
    parser = argparse.ArgumentParser(description="Version-pinned gitmap installer.")
    parser.add_argument("--version", required=True, help="Release tag (e.g. v3.36.0)")
    parser.add_argument("--dir", help="Install directory")
    parser.add_argument("--arch", help="Force amd64 or arm64")
    parser.add_argument("--no-path", action="store_true", help="Skip PATH update")
    parser.add_argument("--no-self-install", action="store_true", help="Skip chained gitmap self-install")
    parser.add_argument("--allow-fallback", action="store_true", help="Use newest patch in same series if missing")
    parser.add_argument("--quiet", action="store_true", help="Quiet output")
    args = parser.parse_args()

    version = args.version
    if not re.match(r"^v\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$", version):
        log_err(f"Invalid version tag: '{version}' (expected vMAJOR.MINOR.PATCH)")
        return EXIT_VERSION_MISSING

    os_name = detect_os()
    arch_name = detect_arch(args.arch)
    log_step(f"Target: {os_name}/{arch_name}", args.quiet)

    release_info = http_get_json(f"https://api.github.com/repos/{REPO}/releases/tags/{version}")
    if not release_info:
        log_err(f"Requested version {version} is not a published release.")
        return EXIT_VERSION_MISSING

    log_step(f"Resolving release {version} ...", args.quiet)

    assets = release_info.get("assets", [])
    expected_tar = f"gitmap-{version}-{os_name}-{arch_name}.tar.gz"
    expected_zip = f"gitmap-{version}-{os_name}-{arch_name}.zip"

    asset_url = ""
    asset_name = ""
    for a in assets:
        name = a.get("name", "")
        if name in (expected_tar, expected_zip):
            asset_url = a.get("browser_download_url", "")
            asset_name = name
            break

    if not asset_url:
        log_err(f"No asset matching {os_name}/{arch_name} in release {version}.")
        return EXIT_UNSUPPORTED_ARCH

    sums_url = ""
    for a in assets:
        if a.get("name") == "checksums.txt":
            sums_url = a.get("browser_download_url", "")
            break

    with tempfile.TemporaryDirectory() as tmp_dir:
        tmp_path = Path(tmp_dir)
        archive_path = tmp_path / asset_name
        sums_path = tmp_path / "checksums.txt"

        log_step(f"Downloading {asset_name} ...", args.quiet)
        try:
            download_file(asset_url, archive_path)
        except Exception as e:
            log_err(f"Download failed: {e}")
            return EXIT_NETWORK

        if sums_url:
            try:
                download_file(sums_url, sums_path)
            except Exception:
                pass

        if sums_path.is_file():
            expected_hash = ""
            with open(sums_path, "r", encoding="utf-8", errors="replace") as f:
                for line in f:
                    parts = line.strip().split()
                    if len(parts) >= 2 and parts[1].endswith(asset_name):
                        expected_hash = parts[0]
                        break
            if expected_hash:
                actual_hash = sha256_file(archive_path)
                if actual_hash.lower() != expected_hash.lower():
                    log_err(f"Checksum mismatch for {asset_name}\n  expected: {expected_hash}\n  actual:   {actual_hash}")
                    return EXIT_CHECKSUM
                log_ok("Checksum verified.", args.quiet)

        # Extract
        extract_dir = tmp_path / "extract"
        extract_dir.mkdir(parents=True, exist_ok=True)
        if asset_name.endswith(".zip"):
            with zipfile.ZipFile(archive_path, "r") as z:
                z.extractall(extract_dir)
        else:
            with tarfile.open(archive_path, "r:gz") as t:
                t.extractall(extract_dir)

        bin_cand = None
        for root, _, files in os.walk(extract_dir):
            for f in files:
                if f.lower() == BINARY_NAME.lower():
                    bin_cand = Path(root) / f
                    break
            if bin_cand:
                break

        if not bin_cand or not bin_cand.is_file():
            log_err("Archive did not contain a recognizable gitmap binary.")
            return EXIT_VERIFY

        default_dir = Path.home() / "AppData" / "Local" / "gitmap-cli" if os.name == "nt" else Path.home() / ".local" / "bin"
        install_dir = Path(args.dir) if args.dir else default_dir
        install_dir.mkdir(parents=True, exist_ok=True)

        dest = install_dir / BINARY_NAME
        shutil.copy2(bin_cand, dest)
        try:
            dest.chmod(0o755)
        except OSError:
            pass
        log_ok(f"Installed: {dest}", args.quiet)

        # Verify version
        try:
            out = subprocess.check_output([str(dest), "--version"], stderr=subprocess.STDOUT, text=True)
            ver_clean = version.lstrip("v")
            if ver_clean not in out:
                log_err(f"Version mismatch: expected {version}, reported: {out.strip()}")
                return EXIT_VERIFY
            log_ok(f"Verified: {out.strip()}", args.quiet)
        except Exception as e:
            log_err(f"Could not run installed binary: {e}")
            return EXIT_VERIFY

        # Chain self-install
        if not args.no_self_install:
            log_step("Chaining gitmap self-install ...", args.quiet)
            res = subprocess.run([str(dest), "self-install"])
            if res.returncode != 0:
                log_warn("self-install failed", args.quiet)
                return EXIT_SELF_INSTALL

    log_ok(f"gitmap {version} installed successfully to {install_dir}", args.quiet)
    return EXIT_OK


if __name__ == "__main__":
    sys.exit(main())
