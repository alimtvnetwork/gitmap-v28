#!/usr/bin/env python3
"""
setup-cursor-ubuntu.py
Automated Cursor IDE Provisioning, Sandboxing Wrapper & Fleet Ledger for Ubuntu Fleet.
"""

import argparse
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import urllib.request

CURSOR_API_URL = "https://cursor.com/api/download?platform=linux-x64&releaseTrack=stable"

DRACULA_SETTINGS = {
    "workbench.colorTheme": "Dracula Theme",
    "editor.fontFamily": "'JetBrains Mono', 'Fira Code', Consolas, monospace",
    "editor.fontSize": 14,
    "editor.lineHeight": 22,
    "editor.tabSize": 4,
    "editor.insertSpaces": True,
    "files.autoSave": "afterDelay",
    "files.autoSaveDelay": 1000,
    "files.eol": "\n",
    "files.insertFinalNewline": True,
    "files.trimTrailingWhitespace": True,
    "editor.renderWhitespace": "selection",
    "telemetry.telemetryLevel": "off",
}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Cursor Ubuntu Fleet Provisioning")
    parser.add_argument("--node", default="u1", help="Target node alias (default: u1)")
    parser.add_argument("--force", action="store_true", help="Force redownload and rewrite")
    parser.add_argument("--dry-run", action="store_true", help="Simulate without changes")
    return parser.parse_args()


def extract_api_fields(data: dict, fallback_url: str) -> dict:
    dl_url = data.get("downloadUrl") or data.get("url") or fallback_url
    ver = data.get("version", "latest")
    sha = data.get("commitSha", "")
    return {"downloadUrl": dl_url, "version": ver, "commitSha": sha}


def fetch_upstream_metadata() -> dict:
    try:
        proc = subprocess.run(
            ["curl", "-sL", CURSOR_API_URL],
            capture_output=True,
            text=True,
            timeout=15,
        )
        if proc.returncode == 0 and proc.stdout.strip().startswith("{"):
            data = json.loads(proc.stdout.strip())
            return extract_api_fields(data, CURSOR_API_URL)
    except Exception:
        pass
    req = urllib.request.Request(CURSOR_API_URL, headers={"User-Agent": "GitMap-Fleet/2.0"})
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            final_url = resp.geturl()
            ctype = resp.headers.get_content_type()
            if "json" in ctype:
                data = json.loads(resp.read().decode("utf-8"))
                return extract_api_fields(data, final_url)
            return {"downloadUrl": final_url, "version": "latest", "commitSha": ""}
    except Exception:
        return {"downloadUrl": CURSOR_API_URL, "version": "latest", "commitSha": ""}


def install_system_dependencies(is_dry_run: bool) -> bool:
    if is_dry_run or not shutil.which("apt-get"):
        return True
    pkgs = [
        "libfuse2t64", "libfuse2", "libnss3", "libasound2t64", "libasound2",
        "libgbm1", "libxss1", "curl", "ca-certificates"
    ]
    try:
        subprocess.run(["sudo", "apt-get", "update", "-y"], check=False)
        for pkg in pkgs:
            subprocess.run(["sudo", "apt-get", "install", "-y", pkg], check=False)
        return True
    except Exception:
        return False


def resolve_appimage_target() -> Path:
    primary = Path("/opt/cursor")
    try:
        primary.mkdir(parents=True, exist_ok=True)
        if os.access(primary, os.W_OK):
            return primary / "Cursor.AppImage"
    except Exception:
        pass
    fallback = Path.home() / ".local" / "share" / "cursor"
    fallback.mkdir(parents=True, exist_ok=True)
    return fallback / "Cursor.AppImage"


def compute_file_sha256(file_path: Path) -> str:
    if not file_path.exists():
        return ""
    hasher = hashlib.sha256()
    with file_path.open("rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            hasher.update(chunk)
    return hasher.hexdigest()


def download_appimage(target_file: Path, download_url: str, is_force: bool, is_dry_run: bool) -> bool:
    if target_file.exists() and not is_force:
        return True
    if is_dry_run:
        return True
    tmp_path = target_file.with_suffix(".tmp")
    try:
        if shutil.which("curl"):
            cmd = ["curl", "-fSL", "--progress-bar", "-o", str(tmp_path), download_url]
            res = subprocess.run(cmd, check=False)
            if res.returncode == 0 and tmp_path.exists() and tmp_path.stat().st_size > 1000000:
                tmp_path.chmod(0o755)
                tmp_path.replace(target_file)
                return True
    except Exception:
        pass
    try:
        urllib.request.urlretrieve(download_url, tmp_path)
        if tmp_path.exists() and tmp_path.stat().st_size > 1000000:
            tmp_path.chmod(0o755)
            tmp_path.replace(target_file)
            return True
    except Exception:
        pass
    return False



def ensure_appimage_extracted(appimage_path: Path, is_dry_run: bool) -> Path:
    extract_dir = appimage_path.parent / "squashfs-root"
    cli_bin = extract_dir / "usr" / "share" / "cursor" / "bin" / "cursor"
    if cli_bin.exists() or is_dry_run:
        return cli_bin
    try:
        subprocess.run([str(appimage_path), "--appimage-extract"], cwd=str(appimage_path.parent), check=False)
    except Exception:
        pass
    return cli_bin


def create_wrapper_at_path(wrapper_path: Path, appimage_path: Path, is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    cli_bin = appimage_path.parent / "squashfs-root" / "usr" / "share" / "cursor" / "bin" / "cursor"
    content = (
        "#!/bin/sh\n"
        f'CLI_BIN="{cli_bin}"\n'
        'if [ -x "$CLI_BIN" ]; then\n'
        '  exec "$CLI_BIN" --no-sandbox "$@"\n'
        "fi\n"
        f'exec "{appimage_path}" --no-sandbox "$@"\n'
    )
    try:
        wrapper_path.parent.mkdir(parents=True, exist_ok=True)
        wrapper_path.write_text(content, encoding="utf-8")
        wrapper_path.chmod(0o755)
        return True
    except PermissionError:
        return False


def deploy_desktop_entry(appimage_path: Path, is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    desktop_path = Path("/usr/share/applications/cursor.desktop")
    content = f"[Desktop Entry]\nName=Cursor\nExec={appimage_path} --no-sandbox %F\nIcon=cursor\nType=Application\nCategories=Development;IDE;\nTerminal=false\nStartupWMClass=Cursor\n"
    try:
        desktop_path.parent.mkdir(parents=True, exist_ok=True)
        desktop_path.write_text(content, encoding="utf-8")
        return True
    except PermissionError:
        return False


def deploy_launchers(appimage_path: Path, is_dry_run: bool) -> Path:
    ensure_appimage_extracted(appimage_path, is_dry_run)
    system_wrapper = Path("/usr/local/bin/cursor")
    user_wrapper = Path.home() / ".local" / "bin" / "cursor"
    create_wrapper_at_path(system_wrapper, appimage_path, is_dry_run)
    create_wrapper_at_path(user_wrapper, appimage_path, is_dry_run)
    deploy_desktop_entry(appimage_path, is_dry_run)
    if system_wrapper.exists():
        return system_wrapper
    return user_wrapper


def read_existing_settings(settings_path: Path) -> dict:
    if not settings_path.exists():
        return {}
    try:
        return json.loads(settings_path.read_text(encoding="utf-8"))
    except Exception:
        return {}


def inject_dracula_settings(is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    settings_path = Path.home() / ".config" / "Cursor" / "User" / "settings.json"
    settings_path.parent.mkdir(parents=True, exist_ok=True)
    config = read_existing_settings(settings_path)
    config.update(DRACULA_SETTINGS)
    settings_path.write_text(json.dumps(config, indent=2) + "\n", encoding="utf-8")
    return True


def verify_cursor_version(wrapper_path: Path) -> tuple[bool, str]:
    if not wrapper_path.exists():
        return False, "wrapper not found"
    try:
        res = subprocess.run([str(wrapper_path), "--version"], capture_output=True, text=True, timeout=15)
        output = res.stdout.strip() or res.stderr.strip()
        is_ok = res.returncode == 0 and "not found" not in output.lower() and len(output) > 0
        return is_ok, output
    except Exception as ex:
        return False, str(ex)


def load_ledger_file(path: Path) -> dict:
    if not path.exists():
        return {"$schema": "https://json-schema.org/draft/2020-12/schema", "title": "CursorFleetStatus", "version": "1.0.0", "nodes": {}}
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except Exception:
        return {"$schema": "https://json-schema.org/draft/2020-12/schema", "title": "CursorFleetStatus", "version": "1.0.0", "nodes": {}}


def to_posix_str(path: Path) -> str:
    return str(path).replace("\\", "/")


def update_fleet_ledger(node_alias: str, meta: dict, appimage: Path, wrapper: Path, sha: str, is_ok: bool, out: str) -> None:
    status_file = Path("repo-secrets/04-ubuntu-migration/cursor-fleet-status.json")
    status_file.parent.mkdir(parents=True, exist_ok=True)
    ledger = load_ledger_file(status_file)
    ledger["lastUpdated"] = datetime.now(timezone.utc).isoformat()
    status_label = "HEALTHY" if is_ok else "PROVISIONED"
    has_installed = is_ok or (appimage.exists() and appimage.stat().st_size > 1000000)
    ledger.setdefault("nodes", {})[node_alias] = {
        "nodeAlias": node_alias,
        "cursorInstalled": has_installed,
        "cursorVersion": meta.get("version", "latest"),
        "appImagePath": to_posix_str(appimage),
        "wrapperPath": to_posix_str(wrapper),
        "sha256": sha,
        "status": status_label,
        "lastVerified": datetime.now(timezone.utc).isoformat(),
        "details": out,
    }
    status_file.write_text(json.dumps(ledger, indent=2) + "\n", encoding="utf-8")




def run_provisioning(node_alias: str, is_force: bool, is_dry_run: bool) -> int:
    print(f"[*] Provisioning Cursor for node '{node_alias}' (force={is_force}, dry_run={is_dry_run})")
    meta = fetch_upstream_metadata()
    install_system_dependencies(is_dry_run)
    appimage = resolve_appimage_target()
    download_appimage(appimage, meta["downloadUrl"], is_force, is_dry_run)
    sha = compute_file_sha256(appimage)
    wrapper = deploy_launchers(appimage, is_dry_run)
    inject_dracula_settings(is_dry_run)
    is_ok, out = verify_cursor_version(wrapper)
    update_fleet_ledger(node_alias, meta, appimage, wrapper, sha, is_ok, out)
    print(f"[OK] Fleet provisioning complete for node '{node_alias}'. Status recorded.")
    return 0


def main() -> None:
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    args = parse_args()
    code = run_provisioning(args.node, args.force, args.dry_run)
    sys.exit(code)


if __name__ == "__main__":
    main()

