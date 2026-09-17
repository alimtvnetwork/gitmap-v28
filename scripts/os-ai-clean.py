#!/usr/bin/env python3
"""scripts/os-ai-clean.py — Standalone cross-platform AI cache cleaner.

Scans and purges Antigravity brain caches, system generated tasks,
OS temp AI dumps, and GitMap installer temp caches.
"""

from __future__ import annotations

import argparse
import glob
import json
import os
import sys
import tempfile
from pathlib import Path

PROTECTED_NAMES = {
    "antigravity_state.pbtxt",
    "installation_id",
    "config.json",
    ".gitkeep",
}


def configure_utf8_output() -> None:
    try:
        if hasattr(sys.stdout, "reconfigure"):
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
            sys.stderr.reconfigure(encoding="utf-8", errors="replace")
    except Exception:
        pass


configure_utf8_output()


def is_protected_file(file_path: Path) -> bool:
    return file_path.name in PROTECTED_NAMES


def scan_directory_tree(dir_path: Path) -> tuple[list[str], int, int]:
    if not dir_path.is_dir():
        return [], 0, 0
    paths: list[str] = []
    total_bytes = 0
    for root, _, files in os.walk(dir_path):
        for f in files:
            p = Path(root) / f
            if is_protected_file(p):
                continue
            try:
                sz = p.stat().st_size
                paths.append(str(p))
                total_bytes += sz
            except OSError:
                continue
    return paths, len(paths), total_bytes


def make_category(name: str, desc: str, paths: list[str], total_bytes: int) -> dict:
    return {
        "name": name,
        "description": desc,
        "paths": paths,
        "file_count": len(paths),
        "total_bytes": total_bytes,
    }


def resolve_home_gemini() -> Path:
    return Path.home() / ".gemini" / "antigravity"


def scan_brain_subdirs(brain_dir: Path) -> tuple[list[str], int]:
    if not brain_dir.is_dir():
        return [], 0
    paths: list[str] = []
    total_bytes = 0
    for item in brain_dir.iterdir():
        if not item.is_dir() or item.name == "tempmediaStorage":
            continue
        scratch_dir = item / "scratch"
        sub_paths, _, sub_bytes = scan_directory_tree(scratch_dir)
        paths.extend(sub_paths)
        total_bytes += sub_bytes
    return paths, total_bytes


def discover_antigravity_brain_caches() -> dict:
    base_dir = resolve_home_gemini()
    paths: list[str] = []
    total_bytes = 0
    target_subdirs = ["crashes", "scratch", str(Path("brain") / "tempmediaStorage")]
    for sub in target_subdirs:
        p_list, _, b_sz = scan_directory_tree(base_dir / sub)
        paths.extend(p_list)
        total_bytes += b_sz
    brain_paths, brain_bytes = scan_brain_subdirs(base_dir / "brain")
    paths.extend(brain_paths)
    total_bytes += brain_bytes
    return make_category(
        "Antigravity Brain Caches",
        "Antigravity scratch files, crash logs, and temp media",
        paths,
        total_bytes,
    )


def scan_brain_task_subdirs(brain_dir: Path) -> tuple[list[str], int]:
    if not brain_dir.is_dir():
        return [], 0
    paths: list[str] = []
    total_bytes = 0
    for item in brain_dir.iterdir():
        if not item.is_dir():
            continue
        task_dir = item / ".system_generated" / "tasks"
        sub_paths, _, sub_bytes = scan_directory_tree(task_dir)
        paths.extend(sub_paths)
        total_bytes += sub_bytes
    return paths, total_bytes


def discover_system_generated_tasks() -> dict:
    paths: list[str] = []
    total_bytes = 0
    ws_tasks = Path.cwd() / ".system_generated" / "tasks"
    ws_paths, _, ws_bytes = scan_directory_tree(ws_tasks)
    paths.extend(ws_paths)
    total_bytes += ws_bytes
    base_dir = resolve_home_gemini()
    brain_paths, brain_bytes = scan_brain_task_subdirs(base_dir / "brain")
    paths.extend(brain_paths)
    total_bytes += brain_bytes
    return make_category(
        "System Generated Tasks",
        "System-generated background execution tasks and worker logs",
        paths,
        total_bytes,
    )


def process_glob_entry(p_str: str) -> tuple[list[str], int]:
    p = Path(p_str)
    if p.is_dir():
        paths, _, b_sz = scan_directory_tree(p)
        return paths, b_sz
    if is_protected_file(p):
        return [], 0
    try:
        return [str(p)], p.stat().st_size
    except OSError:
        return [], 0


def scan_glob_patterns(patterns: list[str]) -> tuple[list[str], int]:
    paths: list[str] = []
    total_bytes = 0
    for pat in patterns:
        for match in glob.glob(pat):
            match_paths, match_bytes = process_glob_entry(match)
            paths.extend(match_paths)
            total_bytes += match_bytes
    return paths, total_bytes


def resolve_updater_dir() -> Path:
    local_app_data = os.environ.get("LOCALAPPDATA", "")
    if local_app_data:
        return Path(local_app_data) / "antigravity-updater"
    return Path.home() / ".local" / "share" / "antigravity-updater"


def discover_os_temp_ai_dumps() -> dict:
    tmp = tempfile.gettempdir()
    patterns = [
        os.path.join(tmp, "antigravity*"),
        os.path.join(tmp, "gemini*"),
        os.path.join(tmp, "agent*"),
        "/tmp/antigravity*",
    ]
    paths, total_bytes = scan_glob_patterns(patterns)
    cache_dir = Path.home() / ".cache" / "antigravity"
    c_paths, _, c_bytes = scan_directory_tree(cache_dir)
    paths.extend(c_paths)
    total_bytes += c_bytes
    up_paths, _, up_bytes = scan_directory_tree(resolve_updater_dir())
    paths.extend(up_paths)
    total_bytes += up_bytes
    return make_category(
        "OS Temp AI Dumps",
        "Operating system temporary AI caches and agent crash dumps",
        paths,
        total_bytes,
    )


def discover_gitmap_installer_caches() -> dict:
    paths: list[str] = []
    total_bytes = 0
    sub_dirs = ["downloads", "build", "sandbox", "purge"]
    bases = [Path.home() / ".gitmap", Path(tempfile.gettempdir()) / "gitmap"]
    for base in bases:
        for sub in sub_dirs:
            p_list, _, b_sz = scan_directory_tree(base / sub)
            paths.extend(p_list)
            total_bytes += b_sz
    return make_category(
        "GitMap Installer Temp Caches",
        "GitMap installer downloads, build sandboxes, and purge staging",
        paths,
        total_bytes,
    )


def discover_all_categories() -> list[dict]:
    return [
        discover_antigravity_brain_caches(),
        discover_system_generated_tasks(),
        discover_os_temp_ai_dumps(),
        discover_gitmap_installer_caches(),
    ]


def format_bytes(total_bytes: int) -> str:
    if total_bytes < 1024:
        return f"{total_bytes} B"
    kb = total_bytes / 1024.0
    if kb < 1024:
        return f"{kb:.2f} KB"
    mb = kb / 1024.0
    if mb < 1024:
        return f"{mb:.2f} MB"
    gb = mb / 1024.0
    return f"{gb:.2f} GB"


def render_preflight_table(categories: list[dict], total_files: int, total_bytes: int) -> str:
    divider = "+------------------------------+-------------+-------------+\n"
    lines = [
        divider,
        "| Category                     | Files       | Size        |\n",
        divider,
    ]
    for cat in categories:
        sz_str = format_bytes(cat["total_bytes"])
        lines.append(f"| {cat['name']:<28} | {cat['file_count']:>11} | {sz_str:>11} |\n")
    lines.append(divider)
    total_sz = format_bytes(total_bytes)
    lines.append(f"| {'Total':<28} | {total_files:>11} | {total_sz:>11} |\n")
    lines.append(divider)
    return "".join(lines)


def remove_file_safely(file_path: str) -> tuple[int, int]:
    p = Path(file_path)
    if is_protected_file(p):
        return 0, 0
    try:
        sz = p.stat().st_size
        p.unlink()
        return 1, sz
    except OSError:
        return 0, 0


def purge_all_files(categories: list[dict]) -> tuple[int, int]:
    total_freed_files = 0
    total_freed_bytes = 0
    for cat in categories:
        for file_path in cat["paths"]:
            cnt, sz = remove_file_safely(file_path)
            total_freed_files += cnt
            total_freed_bytes += sz
    return total_freed_files, total_freed_bytes


def prompt_user_confirmation() -> bool:
    try:
        response = input("Proceed with AI cache cleanup? [y/N]: ").strip().lower()
        return response in ("y", "yes")
    except (EOFError, KeyboardInterrupt):
        return False


def run_json_mode(categories: list[dict], total_files: int, total_bytes: int, is_dry_run: bool, is_yes: bool) -> int:
    status = "clean" if total_files == 0 else "ready"
    freed_files = 0
    freed_bytes = 0
    if not is_dry_run and total_files > 0 and is_yes:
        freed_files, freed_bytes = purge_all_files(categories)
        status = "cleaned"
    payload = {
        "status": status,
        "is_dry_run": is_dry_run,
        "total_files": total_files,
        "total_bytes": total_bytes,
        "human_bytes": format_bytes(total_bytes),
        "freed_files": freed_files,
        "freed_bytes": freed_bytes,
        "categories": categories,
    }
    print(json.dumps(payload, indent=2))
    return 0


def run_terminal_mode(categories: list[dict], total_files: int, total_bytes: int, is_dry_run: bool, is_yes: bool) -> int:
    sys.stdout.write(render_preflight_table(categories, total_files, total_bytes))
    if total_files == 0:
        print("✔ No AI cache files found. System is clean.")
        return 0
    if is_dry_run:
        print("ℹ [dry-run] Preview mode only; no files removed.")
        return 0
    if not is_yes:
        is_confirmed = prompt_user_confirmation()
        if not is_confirmed:
            print("Cleanup aborted.")
            return 0
    freed_files, freed_bytes = purge_all_files(categories)
    print(f"✔ Successfully cleaned {freed_files} AI cache files ({format_bytes(freed_bytes)} freed).")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Clean Antigravity and GitMap AI caches.")
    parser.add_argument("-n", "--dry-run", action="store_true", help="Preview mode without deleting")
    parser.add_argument("-y", "--yes", action="store_true", help="Bypass interactive confirmation")
    parser.add_argument("--json", action="store_true", help="Output results in JSON format")
    args = parser.parse_args()

    cats = discover_all_categories()
    total_files = sum(c["file_count"] for c in cats)
    total_bytes = sum(c["total_bytes"] for c in cats)

    if args.json:
        return run_json_mode(cats, total_files, total_bytes, args.dry_run, args.yes)
    return run_terminal_mode(cats, total_files, total_bytes, args.dry_run, args.yes)


if __name__ == "__main__":
    sys.exit(main())
