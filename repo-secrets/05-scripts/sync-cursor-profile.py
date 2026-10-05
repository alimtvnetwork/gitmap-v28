#!/usr/bin/env python3
"""sync-cursor-profile.py: Cursor profile & AI memory migration utility."""
import argparse
from datetime import datetime, timezone
import io
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tarfile


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Cursor Profile & AI Memory Migration")
    p.add_argument("--export", "--export-only", "-e", dest="is_export", action="store_true")
    p.add_argument("--import", "--import-only", "-i", dest="import_arch", nargs="?", const="default", default="")
    p.add_argument("--sync", "-s", dest="is_sync", action="store_true")
    p.add_argument("--node", "-n", default="u1")
    p.add_argument("--output", "-o", default="")
    p.add_argument("--dry-run", "-d", dest="is_dry_run", action="store_true")
    p.add_argument("--backup", "-b", dest="is_backup", action="store_true", default=True)
    p.add_argument("--force", "-f", dest="is_force", action="store_true")
    return p.parse_args()


def resolve_dirs() -> tuple[Path, Path]:
    appdata, userprofile = os.environ.get("APPDATA", ""), os.environ.get("USERPROFILE", "")
    cfg = Path(appdata) / "Cursor" / "User" if appdata else Path.home() / ".config" / "Cursor" / "User"
    cur = Path(userprofile) / ".cursor" if userprofile else Path.home() / ".cursor"
    if not cur.exists() and Path(".cursor").exists():
        cur = Path(".cursor").resolve()
    return cfg, cur


def is_included_path(rel_str: str) -> bool:
    norm = rel_str.replace("\\", "/").lower()
    for skip in [".log", ".sock", ".lock", "/gpucache", "/cache/", "/code cache/"]:
        if skip in norm:
            return False
    return True


def sanitize_text(content: str) -> str:
    res = re.sub(r"(?i)[dD]:[/\\]work[/\\]([a-zA-Z0-9_\-]+)", r"/home/a/git-work/\1", content)
    res = re.sub(r"(?i)d-work-([a-zA-Z0-9_\-]+)", r"home-a-git-work-\1", res)
    res = re.sub(r"(?i)[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+[/\\]\.cursor", r"/home/a/.cursor", res)
    return re.sub(r"(?i)[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+", r"/home/a", res)


def transform_bytes(m_name: str, data: bytes) -> bytes:
    if m_name.endswith("statsig-cache.json"):
        return b"{}\n"
    for ext in [".json", ".txt", ".md", ".sh", ".py", ".ts", ".js"]:
        if m_name.endswith(ext):
            return sanitize_text(data.decode("utf-8", errors="replace")).encode("utf-8")
    return data


def map_member_path(tag: str, rel_path: Path) -> str:
    parts = list(rel_path.parts)
    if parts and parts[0] == "skills":
        parts[0] = "skills-cursor"
    if len(parts) > 1 and parts[0] == "projects":
        parts[1] = re.sub(r"^d-work-", "home-a-git-work-", parts[1])
    return f"{tag}/" + "/".join(parts)


def append_tar_file(tar: tarfile.TarFile, name: str, data: bytes) -> None:
    ti = tarfile.TarInfo(name=name)
    ti.size, ti.mode, ti.uname, ti.gname = len(data), 0o644, "a", "a"
    ti.mtime = int(datetime.now(timezone.utc).timestamp())
    tar.addfile(ti, io.BytesIO(data))


def add_tree_to_tar(tar: tarfile.TarFile, root: Path, tag: str) -> int:
    if not root.is_dir():
        return 0
    cnt = 0
    for p in root.rglob("*"):
        if p.is_file() and is_included_path(str(p.relative_to(root))):
            m_name = map_member_path(tag, p.relative_to(root))
            append_tar_file(tar, m_name, transform_bytes(m_name, p.read_bytes()))
            cnt += 1
    return cnt


def count_tree_files(root: Path) -> int:
    if not root.is_dir():
        return 0
    return sum(1 for p in root.rglob("*") if p.is_file() and is_included_path(str(p.relative_to(root))))


def build_bundle(out: Path, src_user: Path, src_home: Path, is_dry_run: bool) -> int:
    c1, c2 = count_tree_files(src_user), count_tree_files(src_home)
    if is_dry_run:
        print(f"[DRY-RUN] Discovered {c1} user files, {c2} home items ({c1 + c2} total)")
        return c1 + c2
    out.parent.mkdir(parents=True, exist_ok=True)
    with tarfile.open(out, "w:gz") as tar:
        add_tree_to_tar(tar, src_user, "config_user")
        add_tree_to_tar(tar, src_home, "cursor_home")
    return c1 + c2


def backup_dir(target_dir: Path) -> None:
    if target_dir.exists():
        ts = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
        shutil.copytree(target_dir, target_dir.parent / f"{target_dir.name}.bak.{ts}", dirs_exist_ok=True)


def extract_member(tar: tarfile.TarFile, ti: tarfile.TarInfo, du: Path, dh: Path) -> None:
    dest = du / ti.name[12:] if ti.name.startswith("config_user/") else dh / ti.name[12:]
    norm = str(dest).replace("\\", "/").lower()
    if "/git-work/" in norm or norm.endswith("/git-work"):
        return
    dest.parent.mkdir(parents=True, exist_ok=True)
    f = tar.extractfile(ti)
    if f is not None:
        dest.write_bytes(f.read())


def unpack_bundle(archive_path: Path, is_backup: bool, is_dry_run: bool) -> bool:
    if not archive_path.exists() or is_dry_run:
        return archive_path.exists()
    du, dh = Path.home() / ".config" / "Cursor" / "User", Path.home() / ".cursor"
    if is_backup:
        backup_dir(du)
        backup_dir(dh)
    with tarfile.open(archive_path, "r:gz") as tar:
        for ti in tar.getmembers():
            if ti.isfile():
                extract_member(tar, ti, du, dh)
    return True


def dispatch_remote(node: str, local_arch: Path, is_dry_run: bool) -> bool:
    if is_dry_run:
        return True
    rem_arch, rem_script = f"/tmp/{local_arch.name}", "/tmp/sync-cursor-profile.py"
    cmd_cp = ["gitmap", "ssh", "copy", str(local_arch), f"{node}:{rem_arch}"] if shutil.which("gitmap") else ["scp", str(local_arch), f"{node}:{rem_arch}"]
    subprocess.run(cmd_cp, check=False)
    cmd_sc = ["gitmap", "ssh", "copy", str(Path(__file__).resolve()), f"{node}:{rem_script}"] if shutil.which("gitmap") else ["scp", str(Path(__file__).resolve()), f"{node}:{rem_script}"]
    subprocess.run(cmd_sc, check=False)
    cmd_exec = ["gitmap", "ssh", "exec", node, f"python3 {rem_script} --import {rem_arch}"] if shutil.which("gitmap") else ["ssh", node, f"python3 {rem_script} --import {rem_arch}"]
    return subprocess.run(cmd_exec, check=False).returncode == 0


def main() -> None:
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    a = parse_args()
    if a.import_arch:
        arch = Path("/tmp/cursor-profile-u1-transfer.tar.gz" if a.import_arch == "default" else a.import_arch)
        sys.exit(0 if unpack_bundle(arch, a.is_backup, a.is_dry_run) else 1)
    u_dir, h_dir = resolve_dirs()
    out = Path(a.output) if a.output else Path("/tmp/cursor-profile-u1-transfer.tar.gz")
    cnt = build_bundle(out, u_dir, h_dir, a.is_dry_run)
    print(f"Bundled {cnt} items to {out}")
    if a.is_sync:
        dispatch_remote(a.node, out, a.is_dry_run)


if __name__ == "__main__":
    main()
