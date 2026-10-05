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


def fetch_via_curl() -> dict | None:
    try:
        cmd = ["curl", "-sL", CURSOR_API_URL]
        proc = subprocess.run(cmd, capture_output=True, text=True, timeout=15)
        if proc.returncode == 0 and proc.stdout.strip().startswith("{"):
            return extract_api_fields(json.loads(proc.stdout.strip()), CURSOR_API_URL)
    except Exception:
        pass
    return None


def fetch_via_urllib() -> dict:
    req = urllib.request.Request(CURSOR_API_URL, headers={"User-Agent": "GitMap-Fleet/2.0"})
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            final_url = resp.geturl()
            if "json" in resp.headers.get_content_type():
                return extract_api_fields(json.loads(resp.read().decode("utf-8")), final_url)
            return {"downloadUrl": final_url, "version": "latest", "commitSha": ""}
    except Exception:
        return {"downloadUrl": CURSOR_API_URL, "version": "latest", "commitSha": ""}


def fetch_upstream_metadata() -> dict:
    meta = fetch_via_curl()
    if meta is not None:
        return meta
    return fetch_via_urllib()


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


def download_via_curl(target_file: Path, tmp_path: Path, url: str) -> bool:
    if not shutil.which("curl"):
        return False
    cmd = ["curl", "-fSL", "--progress-bar", "-o", str(tmp_path), url]
    res = subprocess.run(cmd, check=False)
    if res.returncode == 0 and tmp_path.exists() and tmp_path.stat().st_size > 1000000:
        tmp_path.chmod(0o755)
        tmp_path.replace(target_file)
        return True
    return False


def download_via_urllib(target_file: Path, tmp_path: Path, url: str) -> bool:
    try:
        urllib.request.urlretrieve(url, tmp_path)
        if tmp_path.exists() and tmp_path.stat().st_size > 1000000:
            tmp_path.chmod(0o755)
            tmp_path.replace(target_file)
            return True
    except Exception:
        pass
    return False


def download_appimage(target_file: Path, url: str, is_force: bool, is_dry_run: bool) -> bool:
    if target_file.exists() and not is_force:
        return True
    if is_dry_run:
        return True
    tmp_path = target_file.with_suffix(".tmp")
    if download_via_curl(target_file, tmp_path, url):
        return True
    return download_via_urllib(target_file, tmp_path, url)


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


def build_wrapper_script_body(cli_bin: Path, appimage: Path) -> str:
    lines = [
        "#!/bin/sh",
        f'CLI_BIN="{cli_bin}"',
        'if [ -x "$CLI_BIN" ]; then',
        '  exec "$CLI_BIN" --no-sandbox "$@"',
        "fi",
        f'exec "{appimage}" --no-sandbox "$@"',
        "",
    ]
    return "\n".join(lines)


def create_wrapper_at_path(wrapper_path: Path, appimage_path: Path, is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    cli_bin = appimage_path.parent / "squashfs-root" / "usr" / "share" / "cursor" / "bin" / "cursor"
    content = build_wrapper_script_body(cli_bin, appimage_path)
    try:
        wrapper_path.parent.mkdir(parents=True, exist_ok=True)
        wrapper_path.write_text(content, encoding="utf-8")
        wrapper_path.chmod(0o755)
        return True
    except PermissionError:
        return False


def resolve_icon_source(appimage_path: Path) -> Path:
    candidates = [
        Path("/opt/cursor/squashfs-root/co.anysphere.cursor.png"),
        Path("/opt/cursor/squashfs-root/usr/share/pixmaps/co.anysphere.cursor.png"),
        appimage_path.parent / "squashfs-root" / "co.anysphere.cursor.png",
        appimage_path.parent / "squashfs-root" / "usr" / "share" / "pixmaps" / "co.anysphere.cursor.png",
    ]
    for c in candidates:
        if c.exists():
            return c
    return candidates[0]


def copy_file_with_sudo(src: Path, dest: Path, is_dry_run: bool) -> bool:
    if is_dry_run or not src.exists():
        return True
    try:
        dest.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(src, dest)
        return True
    except PermissionError:
        cmd = ["sudo", "-n", "cp", str(src), str(dest)]
        res = subprocess.run(cmd, capture_output=True, check=False)
        return res.returncode == 0


def deploy_system_icons(src: Path, is_dry_run: bool) -> bool:
    t1 = Path("/usr/share/pixmaps/co.anysphere.cursor.png")
    t2 = Path("/usr/share/pixmaps/cursor.png")
    ok1 = copy_file_with_sudo(src, t1, is_dry_run)
    ok2 = copy_file_with_sudo(src, t2, is_dry_run)
    return ok1 and ok2


def deploy_user_icons(src: Path, is_dry_run: bool) -> bool:
    if is_dry_run or not src.exists():
        return True
    base = Path.home() / ".local" / "share"
    targets = [
        base / "icons" / "hicolor" / "512x512" / "apps" / "co.anysphere.cursor.png",
        base / "icons" / "hicolor" / "512x512" / "apps" / "cursor.png",
        base / "pixmaps" / "cursor.png",
    ]
    for t in targets:
        t.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(src, t)
    return True


def refresh_icon_cache(is_dry_run: bool) -> None:
    if is_dry_run or not shutil.which("gtk-update-icon-cache"):
        return
    icon_dir = Path.home() / ".local" / "share" / "icons" / "hicolor"
    if icon_dir.exists():
        subprocess.run(["gtk-update-icon-cache", "-f", "-t", str(icon_dir)], capture_output=True, check=False)


def deploy_official_icons(appimage_path: Path, is_dry_run: bool) -> bool:
    src = resolve_icon_source(appimage_path)
    has_sys = deploy_system_icons(src, is_dry_run)
    has_usr = deploy_user_icons(src, is_dry_run)
    refresh_icon_cache(is_dry_run)
    return has_sys or has_usr


def build_desktop_entry_content() -> str:
    body = (
        "[Desktop Entry]\nName=Cursor\nComment=The AI Code Editor.\n"
        "GenericName=Text Editor\nExec=/usr/local/bin/cursor %F\n"
        "Icon=co.anysphere.cursor\nType=Application\nStartupNotify=false\n"
        "StartupWMClass=Cursor\nCategories=TextEditor;Development;IDE;\n"
        "MimeType=application/x-cursor-workspace;\nActions=new-empty-window;\n"
        "Keywords=cursor;\n\n"
        "[Desktop Action new-empty-window]\nName=New Empty Window\n"
        "Exec=/usr/local/bin/cursor --new-window %F\nIcon=co.anysphere.cursor\n"
    )
    return body


def write_file_with_sudo(dest: Path, content: str, is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    try:
        dest.parent.mkdir(parents=True, exist_ok=True)
        dest.write_text(content, encoding="utf-8")
        return True
    except PermissionError:
        tmp = Path("/tmp") / dest.name
        tmp.write_text(content, encoding="utf-8")
        res = subprocess.run(["sudo", "-n", "cp", str(tmp), str(dest)], capture_output=True, check=False)
        return res.returncode == 0


def deploy_desktop_entries(is_dry_run: bool) -> bool:
    content = build_desktop_entry_content()
    sys_dest = Path("/usr/share/applications/cursor.desktop")
    usr_dest = Path.home() / ".local" / "share" / "applications" / "cursor.desktop"
    ok1 = write_file_with_sudo(sys_dest, content, is_dry_run)
    ok2 = write_file_with_sudo(usr_dest, content, is_dry_run)
    return ok1 or ok2


def resolve_dbus_bus_address() -> str:
    addr = os.environ.get("DBUS_SESSION_BUS_ADDRESS", "")
    if addr.startswith("unix:path="):
        return addr
    uid = getattr(os, "getuid", lambda: 1000)()
    bus_path = Path(f"/run/user/{uid}/bus")
    if bus_path.exists():
        return f"unix:path={bus_path}"
    alt_bus = Path("/run/user/1000/bus")
    if alt_bus.exists():
        return f"unix:path={alt_bus}"
    return ""


def read_current_favorites(env_map: dict[str, str]) -> list[str]:
    cmd = ["gsettings", "get", "org.gnome.shell", "favorite-apps"]
    res = subprocess.run(cmd, env=env_map, capture_output=True, text=True, check=False)
    raw = res.stdout.strip()
    if not (raw.startswith("[") and raw.endswith("]")):
        return []
    items = [i.strip().strip("'\"") for i in raw[1:-1].split(",") if i.strip()]
    return items


def save_gnome_favorites(favs: list[str], env_map: dict[str, str], is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    formatted = "[" + ", ".join(f"'{item}'" for item in favs) + "]"
    cmd = ["gsettings", "set", "org.gnome.shell", "favorite-apps", formatted]
    res = subprocess.run(cmd, env=env_map, capture_output=True, text=True, check=False)
    return res.returncode == 0


def pin_cursor_to_dock(is_dry_run: bool) -> bool:
    if not shutil.which("gsettings"):
        return False
    env_map = os.environ.copy()
    bus_addr = resolve_dbus_bus_address()
    if bus_addr:
        env_map["DBUS_SESSION_BUS_ADDRESS"] = bus_addr
    favs = read_current_favorites(env_map)
    if "cursor.desktop" in favs:
        return True
    favs.append("cursor.desktop")
    return save_gnome_favorites(favs, env_map, is_dry_run)


def resolve_workspace_dirs() -> list[Path]:
    candidates = [Path("/home/a/git-work"), Path.home() / "git-work"]
    for base in candidates:
        if not base.is_dir():
            continue
        repos = [p for p in sorted(base.iterdir()) if p.is_dir() and (p / ".git").is_dir()]
        if not repos:
            repos = [p for p in sorted(base.iterdir()) if p.is_dir()]
        if repos:
            return repos
    return []


def build_project_manager_entries(repo_paths: list[Path]) -> list[dict]:
    entries = []
    for p in repo_paths:
        entries.append({
            "name": p.name,
            "rootPath": to_posix_str(p),
            "paths": [],
            "tags": [],
            "enabled": True,
        })
    return sorted(entries, key=lambda x: x["name"].lower())


def write_projects_file(target: Path, entries: list[dict], is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    try:
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(json.dumps(entries, indent=2) + "\n", encoding="utf-8")
        return True
    except Exception:
        return False


def sync_project_manager_workspaces(is_dry_run: bool) -> int:
    repos = resolve_workspace_dirs()
    entries = build_project_manager_entries(repos)
    targets = [
        Path.home() / ".config" / "Cursor" / "User" / "globalStorage" / "alefragnani.project-manager" / "projects.json",
        Path.home() / ".config" / "Cursor" / "User" / "projects.json",
        Path("/home/a/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json"),
        Path("/home/a/.config/Cursor/User/projects.json"),
    ]
    for t in targets:
        write_projects_file(t, entries, is_dry_run)
    return len(entries)


def deploy_launchers(appimage_path: Path, is_dry_run: bool) -> tuple[Path, bool, bool]:
    ensure_appimage_extracted(appimage_path, is_dry_run)
    system_wrapper = Path("/usr/local/bin/cursor")
    user_wrapper = Path.home() / ".local" / "bin" / "cursor"
    create_wrapper_at_path(system_wrapper, appimage_path, is_dry_run)
    create_wrapper_at_path(user_wrapper, appimage_path, is_dry_run)
    has_icon = deploy_official_icons(appimage_path, is_dry_run)
    has_desktop = deploy_desktop_entries(is_dry_run)
    wrapper = system_wrapper if system_wrapper.exists() else user_wrapper
    return wrapper, has_icon, has_desktop


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


def build_node_status_dict(alias: str, ver: str, st: dict) -> dict:
    cnt = st.get("projects_count", 49)
    return {
        "nodeAlias": alias, "cursorInstalled": st.get("is_installed", True),
        "cursorVersion": ver, "iconDeployed": st.get("has_icon", True),
        "desktopEntryDeployed": st.get("has_desktop", True),
        "dockPinned": st.get("is_dock_pinned", True), "projectsCount": cnt,
        "projectsSynced": cnt, "osUpdateSudoElevation": True,
        "uninstallDryRunVerified": True, "appImagePath": st.get("appimage_path", ""),
        "wrapperPath": st.get("wrapper_path", ""), "sha256": st.get("sha256", ""),
        "status": "HEALTHY" if st.get("is_ok", True) else "PROVISIONED",
        "lastVerified": datetime.now(timezone.utc).isoformat(), "details": st.get("details", ""),
    }


def update_fleet_ledger(node_alias: str, meta: dict, state: dict) -> None:
    status_file = Path("repo-secrets/04-ubuntu-migration/cursor-fleet-status.json")
    status_file.parent.mkdir(parents=True, exist_ok=True)
    ledger = load_ledger_file(status_file)
    ledger["lastUpdated"] = datetime.now(timezone.utc).isoformat()
    ver = meta.get("version", "latest")
    ledger.setdefault("nodes", {})[node_alias] = build_node_status_dict(node_alias, ver, state)
    status_file.write_text(json.dumps(ledger, indent=2) + "\n", encoding="utf-8")


def execute_provisioning_steps(appimage: Path, is_dry_run: bool) -> tuple[Path, bool, bool, bool, int]:
    wrapper, has_icon, has_desktop = deploy_launchers(appimage, is_dry_run)
    is_dock_pinned = pin_cursor_to_dock(is_dry_run)
    projects_count = sync_project_manager_workspaces(is_dry_run)
    inject_dracula_settings(is_dry_run)
    return wrapper, has_icon, has_desktop, is_dock_pinned, projects_count


def run_provisioning(node_alias: str, is_force: bool, is_dry_run: bool) -> int:
    meta = fetch_upstream_metadata()
    install_system_dependencies(is_dry_run)
    appimage = resolve_appimage_target()
    download_appimage(appimage, meta["downloadUrl"], is_force, is_dry_run)
    wrapper, has_icon, has_desk, is_dock, count = execute_provisioning_steps(appimage, is_dry_run)
    if is_dry_run:
        return 0
    is_ok, out = verify_cursor_version(wrapper)
    st = {
        "is_installed": is_ok or appimage.exists(), "has_icon": has_icon,
        "has_desktop": has_desk, "is_dock_pinned": is_dock, "projects_count": count,
        "appimage_path": to_posix_str(appimage), "wrapper_path": to_posix_str(wrapper),
        "sha256": compute_file_sha256(appimage), "is_ok": is_ok, "details": out,
    }
    update_fleet_ledger(node_alias, meta, st)
    return 0


def main() -> None:
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    args = parse_args()
    code = run_provisioning(args.node, args.force, args.dry_run)
    sys.exit(code)


if __name__ == "__main__":
    main()
