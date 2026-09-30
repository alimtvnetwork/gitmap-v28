#!/usr/bin/env python3
"""
Cross-Platform Work Directory Build Artifact & System Cache Cleaner
Scans, plans (Dry-Run / Plan Mode), and reclaims disk space across 8 layers:
  1. work-artifacts      : Work directory (e.g. d:/work) build folders (dist, build, target, tmp, .cache,
                           node_modules/.cache, node_modules/.vite, __pycache__) & untracked binaries
                           (strictly preserving node_modules packages and .venv environments).
  2. go-cache            : Go build cache (GOCACHE), module cache (GOMODCACHE), testcache, fuzzcache.
  3. npm-pnpm-node-cache : npm cache, pnpm CAS store (.pnpm-store, dev-tool/pnpm/store), Yarn, Bun,
                           node-gyp, and V8 compile caches (keeping node_modules intact).
  4. devtools-cache      : Browser DevTools/GPU caches (Chrome, Edge, Brave), VS Code Cache/CachedData,
                           and Antigravity scratch/crash/temp media caches.
  5. temp-dirs           : OS & user temp folders (%TEMP%, %LOCALAPPDATA%/Temp, Windows/Temp, /tmp).
  6. windows-update      : Windows SoftwareDistribution/Download & DeliveryOptimization (or apt/brew cache).
  7. recycle-bin         : Windows Recycle Bin ($Recycle.Bin / Clear-RecycleBin), macOS ~/.Trash, Linux Trash.
  8. git-cache           : Git cache folders (~/.gitcache), .gitmap temp/sandbox/purge, and stale tmp_pack_* files.

Usage Examples:
  # 1. Plan / Dry-Run mode (preview everything that would be removed and total space):
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py --dry-run

  # 2. Interactive mode (shows plan first, then prompts [y/N] before removing):
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py

  # 3. Auto-approve mode (shows plan first, removes automatically with -y, prints space saved summary):
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py -y

  # 4. Target specific categories or custom work directory:
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py --only work-artifacts,go-cache,npm-pnpm-node-cache -y
"""

from __future__ import annotations

import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import time


def configure_utf8_console() -> None:
    if hasattr(sys.stdout, "reconfigure"):
        try:
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
            sys.stderr.reconfigure(encoding="utf-8", errors="replace")
        except Exception:
            pass


configure_utf8_console()

SKIP_REPO_DIRS = {
    "node_modules",
    ".git",
    ".gitmap",
    ".venv",
    "venv",
    "vendor",
}

WORK_BUILD_DIR_NAMES = {
    "dist",
    "build",
    "target",
    ".next",
    ".nuxt",
    ".turbo",
    ".parcel-cache",
    ".vite",
    ".cache",
    "__pycache__",
    ".pytest_cache",
    ".mypy_cache",
    ".ruff_cache",
    "tmp",
    ".tmp",
    "coverage",
    ".nyc_output",
}

WORK_ARTIFACT_EXTENSIONS = {
    ".exe",
    ".dll",
    ".so",
    ".dylib",
    ".test",
    ".out",
    ".pdb",
    ".pyc",
    ".pyo",
}

WORK_ARTIFACT_FILENAMES = {
    "coverage.out",
    "coverage.html",
    "tsconfig.tsbuildinfo",
    ".eslintcache",
    ".stylelintcache",
}

PROTECTED_AI_FILENAMES = {
    "antigravity_state.pbtxt",
    "installation_id",
    "config.json",
    ".gitkeep",
}


def format_bytes(num_bytes: int) -> str:
    if num_bytes < 1024:
        return f"{num_bytes} B"
    kb = num_bytes / 1024.0
    if kb < 1024:
        return f"{kb:.2f} KB"
    mb = kb / 1024.0
    if mb < 1024:
        return f"{mb:.2f} MB"
    gb = mb / 1024.0
    return f"{gb:.2f} GB"


def detect_default_work_dir() -> Path:
    env_work = os.environ.get("WORK_DIR", "").strip()
    if env_work:
        p_env = Path(env_work)
        if p_env.is_dir():
            return p_env

    cwd = Path.cwd().resolve()
    for parent in [cwd, *cwd.parents]:
        if parent.name.lower() == "work":
            return parent

    for drive in ("D:", "C:", "E:"):
        candidate = Path(f"{drive}{os.sep}work")
        if candidate.is_dir():
            return candidate

    home_work = Path.home() / "work"
    if home_work.is_dir():
        return home_work

    return cwd.parent


def compute_path_stats(target_path: Path) -> tuple[int, int]:
    if not target_path.exists():
        return 0, 0
    if target_path.is_file() or target_path.is_symlink():
        try:
            return 1, target_path.stat().st_size
        except OSError:
            return 0, 0

    file_count = 0
    total_bytes = 0
    try:
        for root, _, files in os.walk(target_path):
            for fname in files:
                fpath = Path(root) / fname
                try:
                    total_bytes += fpath.stat().st_size
                    file_count += 1
                except OSError:
                    continue
    except OSError:
        pass

    if file_count == 0:
        file_count = 1
    return file_count, total_bytes


def is_git_tracked_file(file_path: Path) -> bool:
    try:
        res = subprocess.run(
            ["git", "ls-files", "--error-unmatch", str(file_path.name)],
            cwd=str(file_path.parent),
            capture_output=True,
            text=True,
            timeout=3,
        )
        return res.returncode == 0
    except Exception:
        return False


def is_git_tracked_directory(dir_path: Path) -> bool:
    try:
        res = subprocess.run(
            ["git", "ls-files", str(dir_path)],
            cwd=str(dir_path.parent),
            capture_output=True,
            text=True,
            timeout=4,
        )
        has_tracked_output = bool(res.stdout.strip())
        return res.returncode == 0 and has_tracked_output
    except Exception:
        return False


def discover_work_artifacts(work_dir: Path) -> dict:
    targets: list[dict] = []
    if not work_dir.is_dir():
        return build_category_record(
            "work-artifacts",
            f"Work Directory Build & Binary Artifacts ({work_dir})",
            targets,
        )

    for root, dirs, files in os.walk(work_dir):
        root_path = Path(root)

        # Check node_modules internal caches while keeping node_modules packages intact
        if "node_modules" in dirs:
            nm_path = root_path / "node_modules"
            for cache_sub in (".cache", ".vite", ".prisma"):
                sub_path = nm_path / cache_sub
                if sub_path.is_dir():
                    cnt, sz = compute_path_stats(sub_path)
                    targets.append({"path": str(sub_path), "kind": "dir", "files": cnt, "bytes": sz})

        # Check .tanstack/tmp
        if ".tanstack" in dirs:
            ts_tmp = root_path / ".tanstack" / "tmp"
            if ts_tmp.is_dir():
                cnt, sz = compute_path_stats(ts_tmp)
                targets.append({"path": str(ts_tmp), "kind": "dir", "files": cnt, "bytes": sz})

        # Filter skipped directories from recursion
        dirs[:] = [d for d in dirs if d not in SKIP_REPO_DIRS]

        for dname in list(dirs):
            if dname in WORK_BUILD_DIR_NAMES:
                cand_dir = root_path / dname
                if not is_git_tracked_directory(cand_dir):
                    cnt, sz = compute_path_stats(cand_dir)
                    targets.append({"path": str(cand_dir), "kind": "dir", "files": cnt, "bytes": sz})
                dirs.remove(dname)

        for fname in files:
            ext = os.path.splitext(fname)[1].lower()
            is_artifact_ext = ext in WORK_ARTIFACT_EXTENSIONS
            is_artifact_name = fname in WORK_ARTIFACT_FILENAMES
            if is_artifact_ext or is_artifact_name:
                cand_file = root_path / fname
                if not is_git_tracked_file(cand_file):
                    cnt, sz = compute_path_stats(cand_file)
                    if sz > 0 or is_artifact_name:
                        targets.append({"path": str(cand_file), "kind": "file", "files": 1, "bytes": sz})

    return build_category_record(
        "work-artifacts",
        f"Work Dir Build & Binary Artifacts ({work_dir.name})",
        targets,
    )


def get_go_env_path(var_name: str) -> str:
    try:
        res = subprocess.run(["go", "env", var_name], capture_output=True, text=True, timeout=5)
        if res.returncode == 0:
            return res.stdout.strip()
    except Exception:
        pass
    return ""


def discover_go_caches() -> dict:
    targets: list[dict] = []
    seen: set[str] = set()

    candidate_dirs: list[Path] = []
    gocache = get_go_env_path("GOCACHE")
    gomodcache = get_go_env_path("GOMODCACHE")
    if gocache:
        candidate_dirs.append(Path(gocache))
    if gomodcache:
        candidate_dirs.append(Path(gomodcache))

    local_app = os.environ.get("LOCALAPPDATA", "")
    if local_app:
        candidate_dirs.append(Path(local_app) / "go-build")
        candidate_dirs.append(Path(local_app) / "golangci-lint")

    for drive in ("C:", "D:"):
        candidate_dirs.append(Path(f"{drive}{os.sep}dev-tool{os.sep}go{os.sep}cache"))
        candidate_dirs.append(Path(f"{drive}{os.sep}dev-tool{os.sep}go{os.sep}pkg{os.sep}mod"))

    candidate_dirs.append(Path.home() / "go" / "pkg" / "mod")
    candidate_dirs.append(Path.home() / ".cache" / "go-build")
    candidate_dirs.append(Path.home() / "Library" / "Caches" / "go-build")

    for cdir in candidate_dirs:
        if not cdir.exists():
            continue
        resolved = str(cdir.resolve()).lower()
        if any(resolved.startswith(s) or s.startswith(resolved) for s in seen):
            continue
        cnt, sz = compute_path_stats(cdir)
        if sz > 0:
            seen.add(resolved)
            targets.append({"path": str(cdir), "kind": "go-cache", "files": cnt, "bytes": sz})

    return build_category_record(
        "go-cache",
        "Go Build, Test, Fuzz & Module Caches (GOCACHE / GOMODCACHE)",
        targets,
    )


def discover_npm_pnpm_node_caches(work_dir: Path) -> dict:
    targets: list[dict] = []
    seen: set[str] = set()
    local_app = os.environ.get("LOCALAPPDATA", "")
    app_data = os.environ.get("APPDATA", "")

    candidate_dirs: list[Path] = [
        work_dir / ".pnpm-store",
        Path.home() / ".npm",
        Path.home() / ".pnpm-store",
        Path.home() / ".local" / "share" / "pnpm" / "store",
        Path.home() / "Library" / "pnpm" / "store",
        Path.home() / ".cache" / "yarn",
        Path.home() / ".cache" / "node-gyp",
        Path.home() / ".bun" / "install" / "cache",
    ]

    if local_app:
        candidate_dirs.extend([
            Path(local_app) / "npm-cache",
            Path(local_app) / "pnpm" / "store",
            Path(local_app) / "pnpm-cache",
            Path(local_app) / "node-gyp",
            Path(local_app) / "Yarn" / "Cache",
            Path(local_app) / "bun" / "install" / "cache",
        ])
    if app_data:
        candidate_dirs.append(Path(app_data) / "npm-cache")

    for drive in ("C:", "D:"):
        candidate_dirs.append(Path(f"{drive}{os.sep}dev-tool{os.sep}pnpm{os.sep}store"))
        candidate_dirs.append(Path(f"{drive}{os.sep}.pnpm-store"))

    for cdir in candidate_dirs:
        if not cdir.exists():
            continue
        resolved = str(cdir.resolve()).lower()
        if resolved in seen:
            continue
        cnt, sz = compute_path_stats(cdir)
        if sz > 0:
            seen.add(resolved)
            targets.append({"path": str(cdir), "kind": "node-cache", "files": cnt, "bytes": sz})

    return build_category_record(
        "npm-pnpm-node-cache",
        "npm, pnpm Store, Yarn, Bun & Node.js Caches (node_modules kept)",
        targets,
    )


def discover_devtools_caches() -> dict:
    targets: list[dict] = []
    seen: set[str] = set()
    local_app = os.environ.get("LOCALAPPDATA", "")
    app_data = os.environ.get("APPDATA", "")

    candidate_dirs: list[Path] = []
    if local_app:
        chrome_base = Path(local_app) / "Google" / "Chrome" / "User Data" / "Default"
        edge_base = Path(local_app) / "Microsoft" / "Edge" / "User Data" / "Default"
        brave_base = Path(local_app) / "BraveSoftware" / "Brave-Browser" / "User Data" / "Default"
        for b_base in (chrome_base, edge_base, brave_base):
            for sub in ("Cache", "Code Cache", "GPUCache", "DawnCache", "ShaderCache"):
                candidate_dirs.append(b_base / sub)

    if app_data:
        vscode_base = Path(app_data) / "Code"
        for sub in ("Cache", "CachedData", "CachedExtensionVSIXs", "GPUCache", "DawnCache"):
            candidate_dirs.append(vscode_base / sub)

    # macOS / Linux browser & VSCode DevTools caches
    for posix_sub in (
        Path.home() / ".config" / "Code" / "Cache",
        Path.home() / ".config" / "Code" / "CachedData",
        Path.home() / ".cache" / "google-chrome",
        Path.home() / "Library" / "Application Support" / "Code" / "Cache",
        Path.home() / "Library" / "Application Support" / "Code" / "CachedData",
        Path.home() / ".gemini" / "antigravity" / "crashes",
        Path.home() / ".gemini" / "antigravity" / "scratch",
        Path.home() / ".gemini" / "antigravity" / "brain" / "tempmediaStorage",
    ):
        candidate_dirs.append(posix_sub)

    for cdir in candidate_dirs:
        if not cdir.exists():
            continue
        resolved = str(cdir.resolve()).lower()
        if resolved in seen:
            continue
        cnt, sz = compute_path_stats(cdir)
        if sz > 0:
            seen.add(resolved)
            targets.append({"path": str(cdir), "kind": "devtools-cache", "files": cnt, "bytes": sz})

    return build_category_record(
        "devtools-cache",
        "DevTools, Browser GPU/Code Cache, VS Code & Antigravity Caches",
        targets,
    )


def discover_temp_dirs() -> dict:
    targets: list[dict] = []
    seen: set[str] = set()

    candidate_dirs: list[Path] = []
    sys_tmp = Path(tempfile.gettempdir())
    candidate_dirs.append(sys_tmp)

    local_app = os.environ.get("LOCALAPPDATA", "")
    if local_app:
        candidate_dirs.append(Path(local_app) / "Temp")

    if sys.platform == "win32":
        windir = os.environ.get("WINDIR", r"C:\Windows")
        candidate_dirs.append(Path(windir) / "Temp")
    else:
        candidate_dirs.append(Path("/tmp"))
        candidate_dirs.append(Path("/var/tmp"))

    for cdir in candidate_dirs:
        if not cdir.is_dir():
            continue
        try:
            resolved = str(cdir.resolve()).lower()
        except OSError:
            resolved = str(cdir).lower()
        if resolved in seen:
            continue
        seen.add(resolved)

        try:
            for child in cdir.iterdir():
                # Skip active Antigravity conversation task logs in current session
                if "antigravity" in child.name.lower():
                    continue
                cnt, sz = compute_path_stats(child)
                if sz > 0:
                    targets.append({"path": str(child), "kind": "temp-item", "files": cnt, "bytes": sz})
        except OSError:
            continue

    return build_category_record(
        "temp-dirs",
        "OS & User Temporary Directories (%TEMP%, Windows/Temp, /tmp)",
        targets,
    )


def discover_windows_update_caches() -> dict:
    targets: list[dict] = []
    if sys.platform == "win32":
        windir = os.environ.get("WINDIR", r"C:\Windows")
        wu_dirs = [
            Path(windir) / "SoftwareDistribution" / "Download",
            Path(windir) / "SoftwareDistribution" / "DeliveryOptimization",
        ]
        for wdir in wu_dirs:
            if wdir.is_dir():
                cnt, sz = compute_path_stats(wdir)
                targets.append({"path": str(wdir), "kind": "wu-cache", "files": cnt, "bytes": sz})
    elif sys.platform == "darwin":
        brew_cache = Path.home() / "Library" / "Caches" / "Homebrew"
        if brew_cache.is_dir():
            cnt, sz = compute_path_stats(brew_cache)
            targets.append({"path": str(brew_cache), "kind": "os-pkg-cache", "files": cnt, "bytes": sz})
    else:
        apt_cache = Path("/var/cache/apt/archives")
        if apt_cache.is_dir():
            cnt, sz = compute_path_stats(apt_cache)
            targets.append({"path": str(apt_cache), "kind": "os-pkg-cache", "files": cnt, "bytes": sz})

    return build_category_record(
        "windows-update",
        "Windows SoftwareDistribution\\Download & OS Update Caches",
        targets,
    )


def discover_recycle_bin() -> dict:
    targets: list[dict] = []
    if sys.platform == "win32":
        for drive in ("C:", "D:", "E:"):
            rbin = Path(f"{drive}{os.sep}$Recycle.Bin")
            if rbin.is_dir():
                cnt, sz = compute_path_stats(rbin)
                targets.append({"path": str(rbin), "kind": "recycle-bin", "files": cnt, "bytes": sz})
    elif sys.platform == "darwin":
        trash_dir = Path.home() / ".Trash"
        if trash_dir.is_dir():
            cnt, sz = compute_path_stats(trash_dir)
            targets.append({"path": str(trash_dir), "kind": "recycle-bin", "files": cnt, "bytes": sz})
    else:
        linux_trash = Path.home() / ".local" / "share" / "Trash"
        if linux_trash.is_dir():
            cnt, sz = compute_path_stats(linux_trash)
            targets.append({"path": str(linux_trash), "kind": "recycle-bin", "files": cnt, "bytes": sz})

    return build_category_record(
        "recycle-bin",
        "System Recycle Bin / Trash ($Recycle.Bin, ~/.Trash)",
        targets,
    )


def discover_git_caches(work_dir: Path) -> dict:
    targets: list[dict] = []
    candidate_paths = [
        Path.home() / ".gitcache",
        Path.home() / ".cache" / "git",
        Path.home() / ".gitmap" / "downloads",
        Path.home() / ".gitmap" / "build",
        Path.home() / ".gitmap" / "sandbox",
        Path.home() / ".gitmap" / "purge",
        Path.home() / ".gitmap" / "temp",
        Path.home() / ".gitmap" / "logs",
        Path.home() / ".gitmap" / "last_error.log",
        work_dir / ".gitmap" / "temp",
        work_dir / ".gitmap" / "purge",
        work_dir / ".gitmap" / "downloads",
        work_dir / ".gitmap" / "logs",
        work_dir / ".gitmap" / "clone-cache.json",
        work_dir / ".gitmap" / "last_error.log",
    ]
    local_app = os.environ.get("LOCALAPPDATA", "")
    if local_app:
        candidate_paths.append(Path(local_app) / "GitHubDesktop" / "Cache")
        candidate_paths.append(Path(local_app) / "GitCredentialManager" / "cache")

    for cpath in candidate_paths:
        if cpath.exists():
            cnt, sz = compute_path_stats(cpath)
            if sz > 0:
                targets.append({"path": str(cpath), "kind": "git-cache", "files": cnt, "bytes": sz})

    # Scan for orphaned .git/objects/pack/tmp_* and .git/lfs/tmp files in work_dir repositories
    if work_dir.is_dir():
        try:
            for entry in work_dir.iterdir():
                git_dir = entry / ".git"
                git_pack_dir = git_dir / "objects" / "pack"
                if git_pack_dir.is_dir():
                    for pack_file in git_pack_dir.glob("tmp_*"):
                        cnt, sz = compute_path_stats(pack_file)
                        targets.append({"path": str(pack_file), "kind": "git-pack-tmp", "files": cnt, "bytes": sz})
                git_lfs_tmp = git_dir / "lfs" / "tmp"
                if git_lfs_tmp.is_dir():
                    cnt, sz = compute_path_stats(git_lfs_tmp)
                    if sz > 0:
                        targets.append({"path": str(git_lfs_tmp), "kind": "git-lfs-tmp", "files": cnt, "bytes": sz})
        except OSError:
            pass

    return build_category_record(
        "git-cache",
        "Git Cache Folders (~/.gitcache, .gitmap cache/logs, stale tmp_pack_*)",
        targets,
    )


def build_category_record(cat_id: str, title: str, targets: list[dict]) -> dict:
    total_files = sum(t["files"] for t in targets)
    total_bytes = sum(t["bytes"] for t in targets)
    return {
        "id": cat_id,
        "title": title,
        "target_count": len(targets),
        "file_count": total_files,
        "total_bytes": total_bytes,
        "human_size": format_bytes(total_bytes),
        "targets": targets,
    }


def discover_all_layers(work_dir: Path, only_ids: set[str], skip_ids: set[str]) -> list[dict]:
    builders = [
        ("work-artifacts", lambda: discover_work_artifacts(work_dir)),
        ("go-cache", discover_go_caches),
        ("npm-pnpm-node-cache", lambda: discover_npm_pnpm_node_caches(work_dir)),
        ("devtools-cache", discover_devtools_caches),
        ("temp-dirs", discover_temp_dirs),
        ("windows-update", discover_windows_update_caches),
        ("recycle-bin", discover_recycle_bin),
        ("git-cache", lambda: discover_git_caches(work_dir)),
    ]

    categories: list[dict] = []
    for cid, fn in builders:
        if only_ids and cid not in only_ids:
            continue
        if cid in skip_ids:
            continue
        categories.append(fn())
    return categories


def print_plan_report(categories: list[dict], work_dir: Path) -> None:
    total_targets = sum(c["target_count"] for c in categories)
    total_files = sum(c["file_count"] for c in categories)
    total_bytes = sum(c["total_bytes"] for c in categories)

    divider = "=" * 96
    row_sep = "-" * 96
    print(divider)
    print("📋 MULTI-LAYER ARTIFACT & SYSTEM CACHE CLEANUP — PLAN PREVIEW")
    print(f"   Work Directory : {work_dir}")
    print(f"   Platform       : {sys.platform}")
    print(f"   Total Found    : {total_targets} target path(s) | {total_files} file(s) | {format_bytes(total_bytes)} reclaimable")
    print(divider)
    print(f"  {'ID':<21} | {'Category Description':<44} | {'Items':>8} | {'Size':>12}")
    print(row_sep)
    for cat in categories:
        print(
            f"  {cat['id']:<21} | {cat['title'][:44]:<44} | {cat['file_count']:>8} | {cat['human_size']:>12}"
        )
    print(row_sep)
    print(f"  {'TOTAL PLAN':<21} | {'All Selected Layers':<44} | {total_files:>8} | {format_bytes(total_bytes):>12}")
    print(divider)

    print("\n🔍 Top Discovered Targets to Remove (Preview):")
    shown = 0
    all_targets: list[tuple[str, dict]] = []
    for cat in categories:
        for t in cat["targets"]:
            all_targets.append((cat["id"], t))

    for cid, t in sorted(all_targets, key=lambda item: -item[1]["bytes"])[:25]:
        if t["bytes"] > 0 or shown < 10:
            shown += 1
            print(f"   {shown:>2}. [{cid:<19}] {format_bytes(t['bytes']):>10}  {t['path']}")
    if len(all_targets) > shown:
        print(f"   ... and {len(all_targets) - shown} additional target path(s).")
    print("")


def remove_target_contents_or_item(target_path: Path, is_preserve_root_dir: bool = False) -> tuple[int, int]:
    if not target_path.exists():
        return 0, 0

    cnt_before, bytes_before = compute_path_stats(target_path)

    if target_path.is_file() or target_path.is_symlink():
        try:
            target_path.unlink()
            return 1, bytes_before
        except OSError:
            return 0, 0

    if is_preserve_root_dir:
        freed_files = 0
        freed_bytes = 0
        try:
            for child in target_path.iterdir():
                c_cnt, c_bytes = remove_target_contents_or_item(child, is_preserve_root_dir=False)
                freed_files += c_cnt
                freed_bytes += c_bytes
        except OSError:
            pass
        return freed_files, freed_bytes

    try:
        shutil.rmtree(target_path, ignore_errors=False)
        return cnt_before, bytes_before
    except OSError:
        # Fallback: best-effort per-file unlink when some files are locked by OS/processes
        freed_files = 0
        freed_bytes = 0
        for root, dirs, files in os.walk(target_path, topdown=False):
            for fname in files:
                fpath = Path(root) / fname
                try:
                    fsize = fpath.stat().st_size
                    fpath.unlink()
                    freed_files += 1
                    freed_bytes += fsize
                except OSError:
                    continue
            for dname in dirs:
                dpath = Path(root) / dname
                try:
                    dpath.rmdir()
                except OSError:
                    continue
        try:
            target_path.rmdir()
        except OSError:
            pass
        return freed_files, freed_bytes


def run_native_cache_commands(active_ids: set[str]) -> None:
    if "go-cache" in active_ids:
        try:
            subprocess.run(
                ["go", "clean", "-cache", "-testcache", "-fuzzcache", "-modcache"],
                capture_output=True,
                timeout=60,
            )
        except Exception:
            pass

    if "npm-pnpm-node-cache" in active_ids:
        npm_bin = "npm.cmd" if sys.platform == "win32" else "npm"
        pnpm_bin = "pnpm.cmd" if sys.platform == "win32" else "pnpm"
        try:
            subprocess.run([npm_bin, "cache", "clean", "--force"], capture_output=True, timeout=45)
        except Exception:
            pass
        try:
            subprocess.run([pnpm_bin, "store", "prune"], capture_output=True, timeout=45)
        except Exception:
            pass

    if "recycle-bin" in active_ids and sys.platform == "win32":
        try:
            subprocess.run(
                [
                    "powershell",
                    "-NoProfile",
                    "-NonInteractive",
                    "-Command",
                    "Clear-RecycleBin -Force -ErrorAction SilentlyContinue",
                ],
                capture_output=True,
                timeout=30,
            )
        except Exception:
            pass


def execute_cleanup_plan(categories: list[dict]) -> list[dict]:
    active_ids = {c["id"] for c in categories}
    print("⚙️  Running native toolchain & OS cache commands (Go, npm, pnpm, Recycle Bin)...")
    run_native_cache_commands(active_ids)

    results: list[dict] = []
    for cat in categories:
        cid = cat["id"]
        planned_bytes = cat["total_bytes"]
        planned_files = cat["file_count"]
        preserve_root = cid in {"windows-update", "recycle-bin", "go-cache", "npm-pnpm-node-cache", "devtools-cache"}

        for t in cat["targets"]:
            t_path = Path(t["path"])
            remove_target_contents_or_item(t_path, is_preserve_root_dir=preserve_root)

        # Measure residual bytes after cleanup to report exact space reclaimed
        residual_bytes = 0
        residual_files = 0
        for t in cat["targets"]:
            t_path = Path(t["path"])
            if t_path.exists():
                r_cnt, r_sz = compute_path_stats(t_path)
                # Empty directory placeholder has 0 bytes
                if r_sz > 0:
                    residual_files += r_cnt
                    residual_bytes += r_sz

        saved_bytes = max(0, planned_bytes - residual_bytes)
        saved_files = max(0, planned_files - residual_files)
        results.append({
            "id": cid,
            "title": cat["title"],
            "planned_files": planned_files,
            "planned_bytes": planned_bytes,
            "residual_bytes": residual_bytes,
            "saved_files": saved_files,
            "saved_bytes": saved_bytes,
            "saved_human": format_bytes(saved_bytes),
        })

    return results


def print_final_summary(results: list[dict], elapsed_sec: float) -> None:
    total_planned = sum(r["planned_bytes"] for r in results)
    total_saved = sum(r["saved_bytes"] for r in results)
    total_residual = sum(r["residual_bytes"] for r in results)
    total_files_saved = sum(r["saved_files"] for r in results)

    divider = "=" * 96
    row_sep = "-" * 96
    print("\n" + divider)
    print("✅ FINAL CLEANUP SUMMARY — DISK SPACE RECLAIMED")
    print(divider)
    print(f"  {'ID':<21} | {'Before':>12} | {'After (Locked)':>14} | {'Space Saved':>14} | {'Files Removed':>13}")
    print(row_sep)
    for r in results:
        print(
            f"  {r['id']:<21} | {format_bytes(r['planned_bytes']):>12} | "
            f"{format_bytes(r['residual_bytes']):>14} | {r['saved_human']:>14} | {r['saved_files']:>13}"
        )
    print(row_sep)
    print(
        f"  {'TOTAL RECLAIMED':<21} | {format_bytes(total_planned):>12} | "
        f"{format_bytes(total_residual):>14} | {format_bytes(total_saved):>14} | {total_files_saved:>13}"
    )
    print(divider)
    print(f"🎉 Completed in {elapsed_sec:.2f}s — Reclaimed {format_bytes(total_saved)} of disk space!\n")


def confirm_cleanup_prompt(total_bytes: int, total_files: int) -> bool:
    try:
        prompt_msg = (
            f"Proceed with cleaning {total_files} file(s) and reclaiming up to "
            f"{format_bytes(total_bytes)}? [y/N]: "
        )
        ans = input(prompt_msg).strip().lower()
        return ans in {"y", "yes"}
    except (EOFError, KeyboardInterrupt):
        return False


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Multi-Layer Work Directory Build Artifact & System Cache Cleaner (Windows, macOS, Linux)",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Categories (--only / --skip):
  work-artifacts      : Work directory (d:/work) build folders (dist, build, target, tmp, .cache) & binaries
  go-cache            : Go build cache (GOCACHE), module cache (GOMODCACHE), testcache, fuzzcache
  npm-pnpm-node-cache : npm cache, pnpm store (.pnpm-store, dev-tool/pnpm/store), Yarn, Bun, node-gyp
  devtools-cache      : Chrome/Edge/Brave DevTools & GPU caches, VS Code Cache, Antigravity scratch
  temp-dirs           : OS & User Temp folders (%TEMP%, %LOCALAPPDATA%/Temp, Windows/Temp, /tmp)
  windows-update      : Windows SoftwareDistribution\\Download & DeliveryOptimization (or apt/brew cache)
  recycle-bin         : System Recycle Bin / Trash ($Recycle.Bin, ~/.Trash, ~/.local/share/Trash)
  git-cache           : Git caches (~/.gitcache), .gitmap temp/purge/downloads, stale .git tmp_pack_*

Examples:
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py --dry-run
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py -y
  python 03-ai-scripts/44-work-and-system-cache-cleaner.py --only work-artifacts,go-cache -y
""",
    )
    parser.add_argument(
        "--work-dir",
        "-w",
        default="",
        help="Work directory root to scan for build elements (default: auto-detect d:/work or ~/work)",
    )
    parser.add_argument(
        "--dry-run",
        "--plan",
        "-n",
        "-d",
        action="store_true",
        help="Plan Mode (Dry-Run): preview all items and reclaimable space without removing",
    )
    parser.add_argument(
        "--yes",
        "-y",
        "--force",
        "-f",
        action="store_true",
        help="Automatically approve removal after displaying the Plan table (no interactive prompt)",
    )
    parser.add_argument(
        "--only",
        default="",
        help="Comma-separated list of category IDs to run (e.g. work-artifacts,go-cache,npm-pnpm-node-cache)",
    )
    parser.add_argument(
        "--skip",
        "--exclude",
        default="",
        help="Comma-separated list of category IDs to skip",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Emit machine-readable JSON plan/summary report",
    )

    args = parser.parse_args()
    start_ts = time.perf_counter()

    work_dir = Path(args.work_dir).resolve() if args.work_dir else detect_default_work_dir()
    only_ids = {x.strip().lower() for x in args.only.split(",") if x.strip()}
    skip_ids = {x.strip().lower() for x in args.skip.split(",") if x.strip()}

    categories = discover_all_layers(work_dir, only_ids, skip_ids)
    total_files = sum(c["file_count"] for c in categories)
    total_bytes = sum(c["total_bytes"] for c in categories)

    is_dry_run = bool(args.dry_run)
    is_auto_yes = bool(args.yes) or (os.environ.get("SCRIPTS_FIXER_YES") == "1")

    if args.json and is_dry_run:
        print(json.dumps({
            "mode": "dry-run",
            "work_dir": str(work_dir),
            "total_files": total_files,
            "total_bytes": total_bytes,
            "human_size": format_bytes(total_bytes),
            "categories": categories,
        }, indent=2))
        return 0

    # Always show the Plan Table first before any action
    print_plan_report(categories, work_dir)

    if total_bytes == 0 and total_files == 0:
        print("✨ All selected cache and artifact layers are already clean (0 B).")
        return 0

    if is_dry_run:
        print("ℹ️  Plan Mode (--dry-run) complete. Re-run with `-y` to execute cleanup.")
        return 0

    if not is_auto_yes:
        is_confirmed = confirm_cleanup_prompt(total_bytes, total_files)
        if not is_confirmed:
            print("❌ Cleanup canceled by user. No files were removed.")
            return 0

    results = execute_cleanup_plan(categories)
    elapsed = time.perf_counter() - start_ts

    if args.json:
        print(json.dumps({
            "mode": "executed",
            "work_dir": str(work_dir),
            "elapsed_seconds": round(elapsed, 2),
            "total_saved_bytes": sum(r["saved_bytes"] for r in results),
            "total_saved_human": format_bytes(sum(r["saved_bytes"] for r in results)),
            "results": results,
        }, indent=2))
        return 0

    print_final_summary(results, elapsed)
    return 0


if __name__ == "__main__":
    sys.exit(main())
