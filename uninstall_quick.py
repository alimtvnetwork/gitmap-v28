#!/usr/bin/env python3
"""uninstall_quick.py — Cross-platform uninstaller for gitmap.

Removes gitmap binary, install folders, shell PATH entries, and (optionally)
the per-user data folder.

Usage:
  python uninstall_quick.py
  python uninstall_quick.py --keep-data
  python uninstall_quick.py --dir <path>
  python uninstall_quick.py --yes
"""

from __future__ import annotations

import argparse
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path


def info(msg: str) -> None:
    print(f"    \033[90m{msg}\033[0m")


def ok(msg: str) -> None:
    print(f"    \033[32m{msg}\033[0m")


def warn(msg: str) -> None:
    print(f"    \033[33m{msg}\033[0m")


def err(msg: str) -> None:
    print(f"    \033[31m{msg}\033[0m")


def confirm(prompt: str, assume_yes: bool) -> bool:
    if assume_yes:
        return True
    try:
        ans = input(f"\n  \033[33m{prompt} [y/N]: \033[0m").strip().lower()
        return ans in ("y", "yes")
    except (EOFError, KeyboardInterrupt):
        return False


def try_self_uninstall() -> bool:
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    gitmap_bin = shutil.which(bin_name) or shutil.which("gitmap")
    if not gitmap_bin:
        info("gitmap not found on PATH, skipping self-uninstall (will sweep manually)")
        return False

    info(f"Active binary: {gitmap_bin}")
    info("Delegating to: gitmap self-uninstall -y\n")
    try:
        res = subprocess.run([gitmap_bin, "self-uninstall", "-y"])
        if res.returncode == 0:
            ok("self-uninstall completed cleanly")
            return True
    except OSError:
        pass
    warn("self-uninstall exited non-zero; falling back to manual sweep")
    return False


def resolve_install_dir(explicit: str | None) -> Path | None:
    if explicit:
        return Path(explicit)

    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    active = shutil.which(bin_name) or shutil.which("gitmap")
    if active:
        active_path = Path(active).resolve()
        parent = active_path.parent
        grand = parent.parent
        if parent.name in ("gitmap-cli", "gitmap"):
            return grand
        return parent

    home = Path.home()
    candidates = [
        home / ".local" / "bin",
        home / "bin",
        Path("/opt/gitmap"),
        Path("/usr/local/bin"),
        home / "AppData" / "Local" / "gitmap-cli",
        home / "AppData" / "Local" / "gitmap",
    ]
    for c in candidates:
        if (c / bin_name).is_file() or (c / "gitmap-cli" / bin_name).is_file():
            return c
    return None


def remove_install_files(root: Path | None) -> None:
    if not root or not root.exists():
        warn("could not locate a gitmap install dir; skipping file removal")
        return

    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    for sub in ("gitmap-cli", "gitmap"):
        d = root / sub
        if d.is_dir():
            try:
                shutil.rmtree(d)
                ok(f"removed {d}")
            except OSError as e:
                err(f"could not remove {d}: {e}")

    flat = root / bin_name
    if flat.is_file():
        try:
            flat.unlink()
            ok(f"removed {flat}")
        except OSError as e:
            err(f"could not remove {flat}: {e}")


def clean_rc_files(root: Path | None) -> None:
    if not root:
        return
    home = Path.home()
    rc_files = [
        home / ".bashrc",
        home / ".zshrc",
        home / ".profile",
        home / ".bash_profile",
    ]
    root_str = str(root)
    for rc in rc_files:
        if not rc.is_file():
            continue
        try:
            with open(rc, "r", encoding="utf-8", errors="replace") as f:
                lines = f.readlines()
            new_lines = [
                l for l in lines
                if root_str not in l and "gitmap-cli" not in l and "gitmap/gitmap" not in l
            ]
            if len(new_lines) != len(lines):
                backup = rc.with_suffix(rc.suffix + ".gitmap-uninstall.bak")
                shutil.copy2(rc, backup)
                with open(rc, "w", encoding="utf-8") as f:
                    f.writelines(new_lines)
                ok(f"cleaned PATH entries from {rc} (backup: {backup})")
        except OSError:
            pass


def remove_data_folder(keep_data: bool, assume_yes: bool) -> None:
    home = Path.home()
    data_dir = home / ".config" / "gitmap"
    if os.name == "nt":
        app_data = os.getenv("APPDATA")
        if app_data:
            data_dir = Path(app_data) / "gitmap"

    if not data_dir.is_dir():
        return

    if keep_data:
        info(f"keeping data folder: {data_dir}")
        return

    if confirm(f"Also delete user data at {data_dir}?", assume_yes):
        try:
            shutil.rmtree(data_dir)
            ok(f"removed {data_dir}")
        except OSError as e:
            err(f"could not remove {data_dir}: {e}")
    else:
        info(f"kept: {data_dir}")


def remove_stray_binaries() -> None:
    path_env = os.environ.get("PATH", "")
    sep = ";" if os.name == "nt" else ":"
    bin_names = ["gitmap.exe", "gitmap"] if os.name == "nt" else ["gitmap", "gm"]
    seen = set()

    for p in path_env.split(sep):
        if not p.strip():
            continue
        dir_path = Path(p.strip())
        for name in bin_names:
            target = dir_path / name
            if target.is_file() and target not in seen:
                seen.add(target)
                try:
                    target.unlink()
                    ok(f"removed {target}")
                except OSError as e:
                    err(f"could not remove {target}: {e}")


def main() -> int:
    parser = argparse.ArgumentParser(description="Cross-platform gitmap quick uninstaller.")
    parser.add_argument("--dir", help="Explicit install directory")
    parser.add_argument("--keep-data", action="store_true", help="Preserve data folder")
    parser.add_argument("-y", "--yes", action="store_true", help="Assume yes to all prompts")
    args = parser.parse_args()

    print("\n  \033[36mgitmap quick uninstaller\033[0m")
    print("  \033[90m------------------------\033[0m\n")

    print("  \033[36m[1/4] Trying canonical self-uninstall\033[0m")
    if not try_self_uninstall():
        print("\n  \033[36m[2/4] Manual sweep — locating install dir\033[0m")
        root = resolve_install_dir(args.dir)
        if root:
            info(f"Install dir: {root}")
        else:
            warn("no install dir found")

        print("\n  \033[36m[3/4] Removing install files\033[0m")
        remove_install_files(root)

        print("\n  \033[36m[4/4] Cleaning shell rc files\033[0m")
        clean_rc_files(root)

    print("\n  \033[36mExhaustive PATH sweep — removing any remaining gitmap binaries\033[0m")
    remove_stray_binaries()

    print("\n  \033[36mUser data\033[0m")
    remove_data_folder(args.keep_data, args.yes)

    print("\n  \033[32mDone. Open a new shell to refresh PATH.\033[0m\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
