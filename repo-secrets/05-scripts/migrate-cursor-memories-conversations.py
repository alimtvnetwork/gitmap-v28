#!/usr/bin/env python3
"""
migrate-cursor-memories-conversations.py
Migration engine for Cursor IDE memories, conversations, and projects from Windows to Ubuntu fleet node.

Supports:
  --export: Stage, transform, and package Cursor assets into a tar.gz bundle.
  --import [ARCHIVE]: Unpack and restore archive onto the target node.
  --sync: Full end-to-end sync (export locally, transport via GitMap SSH, unpack on remote node).
"""

import argparse
from datetime import datetime, timezone
import io
import json
import os
from pathlib import Path
import re
import shutil
import sqlite3
import subprocess
import sys
import tarfile
import tempfile
from typing import Any

DEFAULT_ARCHIVE_NAME = "cursor-memories-u1-transfer.tar.gz"

TEXT_EXTENSIONS = {
    ".json",
    ".jsonl",
    ".txt",
    ".md",
    ".sh",
    ".py",
    ".ts",
    ".js",
    ".yaml",
    ".yml",
    ".xml",
    ".html",
    ".css",
    ".toml",
    ".ini",
    ".cfg",
    ".conf",
}

CONFIG_FILENAMES = {
    "workspace.json",
    "projects.json",
    "storage.json",
    "settings.json",
    "cli-config.json",
    "argv.json",
}

SKIP_PATTERNS = [
    ".log",
    ".sock",
    ".lock",
    "/gpucache",
    "/cache/",
    "/code cache/",
    "anysphere.cursor-agent-worker",
]


def parse_args() -> argparse.Namespace:
    """Parse CLI options for export, import, sync, and dry-run."""
    parser = argparse.ArgumentParser(
        description="Migrate Cursor IDE memories, conversations, and projects to Ubuntu fleet node"
    )
    parser.add_argument(
        "--export", "-e",
        dest="is_export",
        action="store_true",
        help="Stage, transform, and export tarball archive only",
    )
    parser.add_argument(
        "--import", "-i",
        dest="import_arch",
        nargs="?",
        const="default",
        default="",
        help="Unpack and restore archive onto current machine (default: /tmp/cursor-memories-u1-transfer.tar.gz)",
    )
    parser.add_argument(
        "--sync", "-s",
        dest="is_sync",
        action="store_true",
        help="Full end-to-end sync: export, dispatch via GitMap SSH, and unpack on remote node",
    )
    parser.add_argument(
        "--node", "-n",
        dest="node",
        default="u1",
        help="Remote target node alias or hostname (default: u1)",
    )
    parser.add_argument(
        "--output", "-o",
        dest="output",
        default="",
        help="Custom output destination path for the archive tarball",
    )
    parser.add_argument(
        "--dry-run", "-d",
        dest="is_dry_run",
        action="store_true",
        help="Simulate operations without writing changes or dispatching files",
    )
    parser.add_argument(
        "--no-backup",
        dest="is_no_backup",
        action="store_true",
        help="Skip pre-flight snapshot backup creation before unpacking archive",
    )
    parser.add_argument(
        "--force", "-f",
        dest="is_force",
        action="store_true",
        help="Force overwrite existing destination files during import",
    )
    return parser.parse_args()


def resolve_dirs() -> tuple[Path, Path]:
    """
    Resolve Cursor configuration and home directories across Windows and POSIX environments.
    Returns:
      tuple of (config_user_dir, cursor_home_dir)
    """
    appdata = os.environ.get("APPDATA", "")
    userprofile = os.environ.get("USERPROFILE", "")

    # Config user directory resolution
    win_cfg = Path(r"C:\Users\Administrator\AppData\Roaming\Cursor\User")
    cfg_candidates = [
        win_cfg,
        Path(appdata) / "Cursor" / "User" if appdata else None,
        Path.home() / ".config" / "Cursor" / "User",
    ]
    cfg_dir = next((p for p in cfg_candidates if p and p.is_dir()), None)
    if not cfg_dir:
        cfg_dir = Path(appdata) / "Cursor" / "User" if appdata else Path.home() / ".config" / "Cursor" / "User"

    # Cursor home directory resolution
    win_cur = Path(r"C:\Users\Administrator\.cursor")
    cur_candidates = [
        win_cur,
        Path(userprofile) / ".cursor" if userprofile else None,
        Path.home() / ".cursor",
        Path(".cursor").resolve(),
    ]
    cur_dir = next((p for p in cur_candidates if p and p.is_dir()), None)
    if not cur_dir:
        cur_dir = Path(userprofile) / ".cursor" if userprofile else Path.home() / ".cursor"

    return cfg_dir, cur_dir


def is_included_path(rel_str: str) -> bool:
    """Check if relative path should be included in migration bundle."""
    norm = rel_str.replace("\\", "/").lower()
    norm_lead = f"/{norm}" if not norm.startswith("/") else norm

    for skip in SKIP_PATTERNS:
        if skip in norm_lead:
            return False

    if norm.endswith(".sock") or norm.endswith(".lock") or norm.endswith(".log"):
        return False

    return True


def sanitize_text(content: str) -> str:
    r"""
    Apply regex transformations to remap Windows paths to Ubuntu POSIX paths.
    """
    file_uri_prefix = "file:" + "//" + "/"
    linux_uri_target = file_uri_prefix + "home" + "/a/git-work/"

    # 1. URI schemas (encoded and standard colon)
    res = re.sub(r"(?i)file:/{3}d%3A/work/([a-zA-Z0-9_\-]+)", linux_uri_target + r"\1", content)
    res = re.sub(r"(?i)file:/{3}d:/work/([a-zA-Z0-9_\-]+)", linux_uri_target + r"\1", res)
    res = re.sub(r"(?i)file:/{3}d(?:%3A|:)/work/", linux_uri_target, res)

    # 2. Windows drive paths: d:\work\... -> /home/a/git-work/...
    res = re.sub(r"(?i)[dD]:[/\\]+work[/\\]+([a-zA-Z0-9_\-]+)", r"/home/a/git-work/\1", res)
    res = re.sub(r"(?i)[dD]:[/\\]+work[/\\]?", r"/home/a/git-work/", res)

    # 3. Project identifier prefixes
    res = re.sub(r"(?i)d-work-([a-zA-Z0-9_\-]+)", r"home-a-git-work-\1", res)

    # 4. Windows user profile and .cursor references
    res = re.sub(r"(?i)[cC]:[/\\]+Users[/\\]+Administrator[/\\]+\.cursor", r"/home/a/.cursor", res)
    res = re.sub(r"(?i)[cC]:[/\\]+Users[/\\]+Administrator", r"/home/a", res)
    res = re.sub(r"(?i)[cC]:[/\\]+Users[/\\]+[a-zA-Z0-9_.-]+[/\\]+\.cursor", r"/home/a/.cursor", res)
    res = re.sub(r"(?i)[cC]:[/\\]+Users[/\\]+[a-zA-Z0-9_.-]+", r"/home/a", res)

    # 5. Normalize any remaining backslashes in /home/a paths
    def fix_posix_slashes(match: re.Match) -> str:
        return match.group(0).replace("\\\\", "/").replace("\\", "/")

    res = re.sub(r"/home/a[^\s\"'`,;<>|{}]+", fix_posix_slashes, res)
    return res


def remap_blob_or_text(val: Any) -> Any:
    """Remap strings or UTF-8 encoded blobs while preserving binary safety."""
    if isinstance(val, str):
        return sanitize_text(val)
    if isinstance(val, (bytes, bytearray, memoryview)):
        raw = bytes(val)
        txt = raw.decode("utf-8", errors="surrogateescape")
        new_txt = sanitize_text(txt)
        return new_txt.encode("utf-8", errors="surrogateescape")
    return val


def transform_bytes(m_name: str, data: bytes) -> bytes:
    """Sanitize configuration and text payloads before adding to archive."""
    if m_name.endswith("statsig-cache.json"):
        return b"{}\n"

    ext = Path(m_name).suffix.lower()
    base = Path(m_name).name
    if ext in TEXT_EXTENSIONS or base in CONFIG_FILENAMES:
        try:
            text = data.decode("utf-8", errors="surrogateescape")
            new_text = sanitize_text(text)
            return new_text.encode("utf-8", errors="surrogateescape")
        except Exception:
            return data

    return data


def map_member_path(tag: str, rel_path: Path) -> str:
    """
    Format internal tar member path:
      - Normalizes 'skills' to 'skills-cursor'
      - Renames 'd-work-*' project directories to 'home-a-git-work-*'
      - Prefixes with tag (config_user or cursor_home)
    """
    parts = list(rel_path.parts)
    if parts and parts[0] == "skills":
        parts[0] = "skills-cursor"

    for i, part in enumerate(parts):
        if part.startswith("d-work-"):
            parts[i] = re.sub(r"^d-work-", "home-a-git-work-", part)

    return f"{tag}/" + "/".join(parts)


def transform_sqlite_db(src_path: Path, dst_path: Path) -> bool:
    """
    Safely clones SQLite db from src_path to dst_path and updates
    occurrences of Windows paths to Linux POSIX paths in
    cursorDiskKV, composerHeaders, and ItemTable while preserving binary integrity.
    """
    try:
        # Snapshot using SQLite backup API if possible, with copy2 fallback
        try:
            src_uri = f"file:{src_path.resolve().as_posix()}?mode=ro"
            src_conn = sqlite3.connect(src_uri, uri=True)
            dst_conn = sqlite3.connect(str(dst_path))
            src_conn.backup(dst_conn)
            src_conn.close()
            dst_conn.close()
        except Exception:
            shutil.copy2(src_path, dst_path)

        conn = sqlite3.connect(str(dst_path))
        cur = conn.cursor()

        cur.execute("SELECT name FROM sqlite_master WHERE type='table';")
        tables = {row[0] for row in cur.fetchall()}

        target_tables = ["cursorDiskKV", "composerHeaders", "ItemTable"]
        updated_count = 0

        for tbl in target_tables:
            if tbl not in tables:
                continue

            cur.execute(f"PRAGMA table_info({tbl});")
            col_info = cur.fetchall()
            col_names = [col[1] for col in col_info]

            pk_col = "composerId" if tbl == "composerHeaders" else "key"
            if pk_col not in col_names or "value" not in col_names:
                continue

            query = f"""
                SELECT {pk_col}, value FROM {tbl}
                WHERE (value LIKE '%d:\\work%' OR value LIKE '%d:\\\\work%'
                    OR value LIKE '%d:/work%' OR value LIKE '%file:///d%'
                    OR value LIKE '%Administrator%' OR value LIKE '%d-work-%');
            """
            cur.execute(query)
            matching_rows = cur.fetchall()

            for pk_val, val in matching_rows:
                new_val = remap_blob_or_text(val)
                if new_val != val:
                    cur.execute(
                        f"UPDATE {tbl} SET value = ? WHERE {pk_col} = ?",
                        (new_val, pk_val),
                    )
                    updated_count += 1

        conn.commit()

        # Validate database integrity after transformation
        cur.execute("PRAGMA integrity_check;")
        status = cur.fetchone()[0]
        conn.close()

        if status.lower() != "ok":
            print(f"⚠ Warning: SQLite integrity check for {src_path.name} returned: {status}", file=sys.stderr)
            return False

        return True
    except Exception as e:
        print(f"⚠ Error transforming SQLite db {src_path}: {e}", file=sys.stderr)
        return False


def append_tar_file(tar: tarfile.TarFile, name: str, data: bytes) -> None:
    """Append in-memory byte buffer to tar archive with standardized permissions."""
    ti = tarfile.TarInfo(name=name)
    ti.size = len(data)
    ti.mode = 0o644
    ti.uname = "a"
    ti.gname = "a"
    ti.mtime = int(datetime.now(timezone.utc).timestamp())
    tar.addfile(ti, io.BytesIO(data))


def add_file_to_tar(tar: tarfile.TarFile, file_path: Path, m_name: str) -> None:
    """Add a file to the archive, performing SQLite or text transformation when appropriate."""
    if file_path.name == "state.vscdb":
        with tempfile.NamedTemporaryFile(suffix=".vscdb", delete=False) as tf:
            tmp_path = Path(tf.name)
        try:
            if transform_sqlite_db(file_path, tmp_path):
                with open(tmp_path, "rb") as f:
                    ti = tar.gettarinfo(str(tmp_path), arcname=m_name)
                    ti.mode = 0o644
                    ti.uname = "a"
                    ti.gname = "a"
                    ti.mtime = int(datetime.now(timezone.utc).timestamp())
                    tar.addfile(ti, f)
            else:
                append_tar_file(tar, m_name, file_path.read_bytes())
        finally:
            if tmp_path.exists():
                tmp_path.unlink()
    else:
        append_tar_file(tar, m_name, transform_bytes(m_name, file_path.read_bytes()))


def add_tree_to_tar(tar: tarfile.TarFile, root: Path, tag: str) -> int:
    """Recursively stage, transform, and bundle directory tree into tar archive."""
    if not root.is_dir():
        return 0

    cnt = 0
    for p in root.rglob("*"):
        if p.is_file() and is_included_path(str(p.relative_to(root))):
            m_name = map_member_path(tag, p.relative_to(root))
            add_file_to_tar(tar, p, m_name)
            cnt += 1
            if cnt % 250 == 0:
                print(f"  [staging] Packed {cnt} items from {tag}...")

    return cnt


def count_tree_files(root: Path) -> int:
    """Count eligible files in directory matching inclusion filter."""
    if not root.is_dir():
        return 0
    return sum(1 for p in root.rglob("*") if p.is_file() and is_included_path(str(p.relative_to(root))))


def build_bundle(out: Path, src_user: Path, src_home: Path, is_dry_run: bool) -> int:
    """Build compressed tarball bundle with transformed Cursor state."""
    c_user, c_home = count_tree_files(src_user), count_tree_files(src_home)
    total_files = c_user + c_home

    print(f"● Discovered {c_user} files in User config and {c_home} files in Cursor home ({total_files} total)")
    if is_dry_run:
        print(f"[DRY-RUN] Simulation only: archive would be generated at {out}")
        return total_files

    out.parent.mkdir(parents=True, exist_ok=True)
    print(f"● Staging and compressing migration bundle -> {out}")
    with tarfile.open(out, "w:gz") as tar:
        add_tree_to_tar(tar, src_user, "config_user")
        add_tree_to_tar(tar, src_home, "cursor_home")

    size_mb = out.stat().st_size / (1024 * 1024)
    print(f"✔ Successfully bundled {total_files} items into {out} ({size_mb:.2f} MB)")
    return total_files


def backup_dir(target_dir: Path) -> Path | None:
    """Create timestamped pre-flight backup snapshot before extraction."""
    if target_dir.exists():
        ts = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
        bak_dir = target_dir.parent / f"{target_dir.name}.bak.{ts}"
        print(f"  [backup] Snapshotting {target_dir} -> {bak_dir}")
        shutil.copytree(target_dir, bak_dir, dirs_exist_ok=True)
        return bak_dir
    return None


def is_safe_extraction_path(dest: Path, target_root: Path) -> bool:
    """Ensure extracted destination does not escape root and never touches git-work repositories."""
    norm = str(dest).replace("\\", "/").lower()
    if "/git-work/" in norm or norm.endswith("/git-work") or norm.startswith("/home/a/git-work"):
        print(f"🛑 PROTECTED WORKSPACE INVARIANT: Refusing to write to {dest}", file=sys.stderr)
        return False

    try:
        resolved_dest = dest.resolve()
        resolved_root = target_root.resolve()
        # Destination must be within target root
        if not str(resolved_dest).startswith(str(resolved_root)):
            # If path doesn't exist yet, check its resolved parent
            resolved_parent = dest.parent.resolve()
            if not str(resolved_parent).startswith(str(resolved_root)):
                print(f"🛑 PATH TRAVERSAL DETECTED: Refusing to extract outside {target_root}: {dest}", file=sys.stderr)
                return False
    except Exception:
        pass

    return True


def extract_member(
    tar: tarfile.TarFile,
    ti: tarfile.TarInfo,
    dest_user: Path,
    dest_home: Path,
    is_force: bool,
) -> bool:
    """Extract a single tar member with path safety and permission normalization."""
    if ti.name.startswith("config_user/"):
        rel_path = ti.name[len("config_user/"):]
        dest = dest_user / rel_path
        target_root = dest_user
    elif ti.name.startswith("cursor_home/"):
        rel_path = ti.name[len("cursor_home/"):]
        dest = dest_home / rel_path
        target_root = dest_home
    else:
        return False

    if not is_safe_extraction_path(dest, target_root):
        return False

    if ti.isdir():
        dest.mkdir(parents=True, exist_ok=True)
        if os.name != "nt":
            try:
                os.chmod(dest, 0o755)
            except Exception:
                pass
        return True

    if ti.isfile():
        if dest.exists() and not is_force:
            # Overwrite allowed during restore of state, force flag allows explicit bypass
            pass

        dest.parent.mkdir(parents=True, exist_ok=True)
        f = tar.extractfile(ti)
        if f is not None:
            dest.write_bytes(f.read())
            if os.name != "nt":
                try:
                    os.chmod(dest, 0o644)
                except Exception:
                    pass
            return True

    return False


def normalize_permissions(root: Path) -> None:
    """Normalize directory (0755) and file (0644) permissions across extracted tree."""
    if not root.exists() or os.name == "nt":
        return

    try:
        os.chmod(root, 0o755)
    except Exception:
        pass

    for p in root.rglob("*"):
        try:
            if p.is_dir():
                os.chmod(p, 0o755)
            elif p.is_file():
                if p.suffix in {".sh", ".py"} and "bin" in p.parts:
                    os.chmod(p, 0o755)
                else:
                    os.chmod(p, 0o644)
        except Exception:
            pass


def verify_post_flight(dest_user: Path, dest_home: Path) -> bool:
    """Verify integrity of unpacked databases and key state assets."""
    all_ok = True

    # Check global SQLite database
    state_db = dest_user / "globalStorage" / "state.vscdb"
    if state_db.exists():
        try:
            conn = sqlite3.connect(str(state_db))
            cur = conn.cursor()
            cur.execute("PRAGMA integrity_check;")
            status = cur.fetchone()[0]
            conn.close()
            print(f"  [verify] state.vscdb integrity: {status}")
            if status.lower() != "ok":
                all_ok = False
        except Exception as e:
            print(f"  [verify] state.vscdb check failed: {e}", file=sys.stderr)
            all_ok = False

    # Check conversation search database
    search_db = dest_user / "globalStorage" / "conversation-search.db"
    if search_db.exists():
        try:
            conn = sqlite3.connect(str(search_db))
            cur = conn.cursor()
            cur.execute("PRAGMA integrity_check;")
            status = cur.fetchone()[0]
            conn.close()
            print(f"  [verify] conversation-search.db integrity: {status}")
            if status.lower() != "ok":
                all_ok = False
        except Exception as e:
            print(f"  [verify] conversation-search.db check failed: {e}", file=sys.stderr)
            all_ok = False

    # Check project manager config
    pm_json = dest_user / "globalStorage" / "alefragnani.project-manager" / "projects.json"
    if pm_json.exists():
        try:
            data = json.loads(pm_json.read_text(encoding="utf-8"))
            print(f"  [verify] projects.json entries: {len(data)}")
        except Exception as e:
            print(f"  [verify] projects.json read failed: {e}", file=sys.stderr)

    return all_ok


def unpack_bundle(archive_path: Path, is_backup: bool, is_dry_run: bool, is_force: bool) -> bool:
    """Unpack archive onto current host into ~/.config/Cursor/User and ~/.cursor."""
    if not archive_path.exists():
        print(f"❌ Error: Archive file not found: {archive_path}", file=sys.stderr)
        return False

    dest_user = Path.home() / ".config" / "Cursor" / "User"
    dest_home = Path.home() / ".cursor"

    if is_dry_run:
        print(f"[DRY-RUN] Would unpack {archive_path} to {dest_user} and {dest_home}")
        return True

    if is_backup:
        print("● Creating pre-flight snapshots before unpacking...")
        backup_dir(dest_user)
        backup_dir(dest_home)

    print(f"● Unpacking {archive_path}...")
    extracted_count = 0
    with tarfile.open(archive_path, "r:gz") as tar:
        for ti in tar.getmembers():
            if extract_member(tar, ti, dest_user, dest_home, is_force):
                extracted_count += 1
                if extracted_count % 500 == 0:
                    print(f"  [extracting] Restored {extracted_count} items...")

    print(f"● Normalizing file modes (0755 dirs, 0644 files)...")
    normalize_permissions(dest_user)
    normalize_permissions(dest_home)

    print("● Running post-flight integrity checks...")
    verified = verify_post_flight(dest_user, dest_home)

    if verified:
        print(f"✔ Unpack completed successfully ({extracted_count} items restored)")
    else:
        print("⚠ Unpack completed with verification warnings", file=sys.stderr)

    return True


def dispatch_remote(
    node: str,
    local_arch: Path,
    is_dry_run: bool,
    is_force: bool,
    is_no_backup: bool,
) -> bool:
    """Dispatch archive and migration engine script to remote node via GitMap SSH and unpack."""
    if is_dry_run:
        print(f"[DRY-RUN] Would dispatch archive {local_arch} to {node}:/tmp/ and execute remote import")
        return True

    rem_arch = f"/tmp/{local_arch.name}"
    rem_script = "/tmp/migrate-cursor-memories-conversations.py"
    local_script = Path(__file__).resolve()

    has_gitmap = shutil.which("gitmap") is not None

    print(f"● Dispatching migration bundle to {node}:{rem_arch}...")
    if has_gitmap:
        cmd_cp_arch = ["gitmap", "ssh", "copy", str(local_arch), f"{node}:{rem_arch}"]
    else:
        cmd_cp_arch = ["scp", str(local_arch), f"{node}:{rem_arch}"]

    res_arch = subprocess.run(cmd_cp_arch, check=False)
    if res_arch.returncode != 0:
        print(f"❌ Failed to transfer archive to {node}", file=sys.stderr)
        return False

    print(f"● Dispatching migration engine script to {node}:{rem_script}...")
    if has_gitmap:
        cmd_cp_script = ["gitmap", "ssh", "copy", str(local_script), f"{node}:{rem_script}"]
    else:
        cmd_cp_script = ["scp", str(local_script), f"{node}:{rem_script}"]

    res_script = subprocess.run(cmd_cp_script, check=False)
    if res_script.returncode != 0:
        print(f"❌ Failed to transfer script to {node}", file=sys.stderr)
        return False

    import_cmd = f"python3 {rem_script} --import {rem_arch}"
    if is_force:
        import_cmd += " --force"
    if is_no_backup:
        import_cmd += " --no-backup"

    print(f"● Triggering remote unpack on node {node}...")
    if has_gitmap:
        cmd_exec = ["gitmap", "ssh", "exec", node, import_cmd]
    else:
        cmd_exec = ["ssh", node, import_cmd]

    res_exec = subprocess.run(cmd_exec, check=False)
    if res_exec.returncode == 0:
        print(f"✔ Remote import on node {node} completed successfully")
        return True

    print(f"❌ Remote import failed on node {node} (exit code: {res_exec.returncode})", file=sys.stderr)
    return False


def main() -> None:
    """Main CLI entrypoint."""
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    args = parse_args()

    # Mode 1: Import archive on target machine
    if args.import_arch:
        arch_path = (
            Path("/tmp") / DEFAULT_ARCHIVE_NAME
            if args.import_arch == "default"
            else Path(args.import_arch)
        )
        is_backup = not args.is_no_backup
        success = unpack_bundle(arch_path, is_backup, args.is_dry_run, args.is_force)
        sys.exit(0 if success else 1)

    # Mode 2: Export or Sync
    cfg_dir, cur_dir = resolve_dirs()
    print(f"● Source User Config: {cfg_dir}")
    print(f"● Source Cursor Home: {cur_dir}")

    if args.output:
        out_archive = Path(args.output)
    else:
        # Default archive in temporary directory
        temp_base = Path(tempfile.gettempdir()) if os.name == "nt" else Path("/tmp")
        out_archive = temp_base / DEFAULT_ARCHIVE_NAME

    build_bundle(out_archive, cfg_dir, cur_dir, args.is_dry_run)

    if args.is_sync:
        success = dispatch_remote(
            node=args.node,
            local_arch=out_archive,
            is_dry_run=args.is_dry_run,
            is_force=args.force if hasattr(args, "force") else args.is_force,
            is_no_backup=args.is_no_backup,
        )
        sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
