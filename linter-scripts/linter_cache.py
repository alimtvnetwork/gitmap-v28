"""Reusable incremental change detection and hash caching engine for repository linters."""
from __future__ import annotations

import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import time
from typing import Any


def get_linter_cache_dir(root_dir: Path) -> Path:
    """Returns the directory used for persistent linter hash caches."""
    target = root_dir / ".ai-memory" / "cicd" / "cache" / "linters"
    target.mkdir(parents=True, exist_ok=True)

    return target


def get_linter_cache_path(linter_name: str, root_dir: Path) -> Path:
    """Returns the file path for a specific linter's cache manifest."""
    return get_linter_cache_dir(root_dir) / f"{linter_name}.json"


def load_linter_cache(linter_name: str, root_dir: Path) -> dict[str, Any]:
    """Loads existing cache data for a linter if present on disk."""
    cache_path = get_linter_cache_path(linter_name, root_dir)
    if cache_path.is_file():
        try:
            return json.loads(cache_path.read_text(encoding="utf-8"))
        except Exception:
            pass

    return {}


def save_linter_cache(linter_name: str, root_dir: Path, cache_data: dict[str, Any]) -> None:
    """Persists updated linter cache data to disk."""
    cache_path = get_linter_cache_path(linter_name, root_dir)
    cache_data["last_run_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ")
    try:
        cache_path.write_text(json.dumps(cache_data, indent=2), encoding="utf-8")
    except Exception:
        pass


def compute_file_hash(file_path: Path) -> str:
    """Computes SHA-256 hash of a file's content for change detection."""
    if not file_path.is_file():
        return ""
    try:
        return hashlib.sha256(file_path.read_bytes()).hexdigest()[:16]
    except Exception:
        return ""


def run_git_lines(args: list[str], root_dir: Path) -> list[str]:
    """Executes a git command and returns stripped stdout lines."""
    git_bin = shutil.which("git") or "git"
    res = subprocess.run([git_bin] + args, cwd=str(root_dir), capture_output=True, text=True)
    if res.returncode != 0:
        return []

    return [line.strip() for line in res.stdout.splitlines() if line.strip()]


def get_git_head_hash(root_dir: Path) -> str:
    """Returns the current 40-character git HEAD commit hash."""
    lines = run_git_lines(["rev-parse", "HEAD"], root_dir)

    return lines[0].strip() if lines else ""


def is_commit_ancestor(ancestor_sha: str, descendant_sha: str, root_dir: Path) -> bool:
    """Checks if ancestor_sha is an ancestor of descendant_sha."""
    git_bin = shutil.which("git") or "git"
    cmd = [git_bin, "merge-base", "--is-ancestor", ancestor_sha, descendant_sha]
    res = subprocess.run(cmd, cwd=str(root_dir), capture_output=True)

    return bool(res.returncode == 0)


def parse_git_status_path(line: str) -> str:
    """Extracts relative file path from a git porcelain status line."""
    path_part = line[2:].strip().strip('"')
    if " -> " in path_part:
        path_part = path_part.split(" -> ")[-1].strip()

    return path_part.replace("\\", "/")


def get_uncommitted_files(root_dir: Path) -> set[str]:
    """Retrieves all uncommitted (staged, unstaged, untracked) file paths."""
    status_lines = run_git_lines(["status", "--porcelain", "-uall"], root_dir)

    return {parse_git_status_path(line) for line in status_lines if parse_git_status_path(line)}


def get_committed_diff_files(last_hash: str, head_hash: str, root_dir: Path) -> set[str]:
    """Retrieves files modified between last_hash and current HEAD."""
    if last_hash == head_hash:
        return set()
    diff_lines = run_git_lines(["diff", "--name-only", f"{last_hash}..{head_hash}"], root_dir)

    return {line.replace("\\", "/").strip() for line in diff_lines if line.strip()}


def is_matching_subdir(rel_posix: str, target_subdir: str) -> bool:
    """Checks if relative path is inside target subdirectory."""
    if target_subdir == "":
        return True

    return bool(rel_posix == target_subdir or rel_posix.startswith(target_subdir + "/"))


def is_path_matching_linter(
    rel_posix: str, target_exts: set[str], exclude_dirs: set[str], target_subdir: str = ""
) -> bool:
    """Checks if a relative path matches extension, directory filters, and subdir."""
    if not is_matching_subdir(rel_posix, target_subdir):
        return False
    ext = os.path.splitext(rel_posix)[1].lower()
    if ext not in target_exts:
        return False
    parts = rel_posix.split("/")
    has_excluded = any(part in exclude_dirs for part in parts)

    return bool(not has_excluded)


def collect_all_repo_targets(
    root_dir: Path, target_exts: set[str], exclude_dirs: set[str], target_subdir: str = ""
) -> list[Path]:
    """Recursively collects all matching files across repository or subdir."""
    scan_dir = root_dir / target_subdir if target_subdir else root_dir
    if not scan_dir.exists():
        return []
    targets: list[Path] = []
    for root, dirs, files in os.walk(scan_dir):
        dirs[:] = [d for d in dirs if d not in exclude_dirs and not any(d.startswith(ex) for ex in exclude_dirs)]
        for f in files:
            p = Path(root) / f
            rel = p.relative_to(root_dir).as_posix()
            if is_path_matching_linter(rel, target_exts, exclude_dirs, target_subdir):
                targets.append(p)

    return sorted(targets)


def is_file_changed_from_cache(p: Path, rel: str, cached_hashes: dict[str, str] | None) -> bool:
    """Checks if file is newly added or modified relative to cached hash."""
    if not p.is_file():
        return False
    if cached_hashes is None:
        return True

    return bool(cached_hashes.get(rel) != compute_file_hash(p))


def filter_eligible_changed_files(
    changed_rel_paths: set[str], root_dir: Path, target_exts: set[str],
    exclude_dirs: set[str], cached_hashes: dict[str, str] | None = None,
    target_subdir: str = ""
) -> list[Path]:
    """Filters changed paths to valid existing files matching linter rules."""
    matched: list[Path] = []
    for rel in sorted(changed_rel_paths):
        if is_path_matching_linter(rel, target_exts, exclude_dirs, target_subdir):
            p = root_dir / rel
            if is_file_changed_from_cache(p, rel, cached_hashes):
                matched.append(p)

    return matched


def resolve_changed_targets_with_cache(
    last_hash: str, head_hash: str, root_dir: Path, target_exts: set[str],
    exclude_dirs: set[str], cached_hashes: dict[str, str] | None = None,
    target_subdir: str = ""
) -> tuple[list[Path], str, bool]:
    """Resolves changed files using ancestor diff and uncommitted status."""
    uncommitted = get_uncommitted_files(root_dir)
    committed = get_committed_diff_files(last_hash, head_hash, root_dir)
    changed = filter_eligible_changed_files(
        uncommitted | committed, root_dir, target_exts, exclude_dirs, cached_hashes, target_subdir
    )
    count = len(changed)
    desc = f"incremental ({count} changed file(s) since {last_hash[:8]})" if count > 0 else f"cached (0 changes since {last_hash[:8]})"

    return changed, desc, True


def is_valid_cache_head(last_hash: str, head_hash: str, root_dir: Path) -> bool:
    """Checks if cached hash is non-empty and an ancestor of current head commit."""
    has_hash = bool(last_hash)

    return has_hash and is_commit_ancestor(last_hash, head_hash, root_dir)


def resolve_linter_targets(
    linter_name: str, root_dir: Path, target_exts: set[str], exclude_dirs: set[str],
    force_all: bool = False, target_subdir: str = ""
) -> tuple[list[Path], str, bool]:
    """Resolves target files for a linter, running on changed files by default."""
    if force_all:
        return collect_all_repo_targets(root_dir, target_exts, exclude_dirs, target_subdir), "forced full scan", False
    cache = load_linter_cache(linter_name, root_dir)
    last_hash, head_hash = cache.get("last_git_hash", ""), get_git_head_hash(root_dir)
    if not is_valid_cache_head(last_hash, head_hash, root_dir):
        return collect_all_repo_targets(root_dir, target_exts, exclude_dirs, target_subdir), "baseline scan (cache seeded)", False

    return resolve_changed_targets_with_cache(
        last_hash, head_hash, root_dir, target_exts, exclude_dirs, cache.get("file_hashes", {}), target_subdir
    )



def record_linter_success(
    linter_name: str, root_dir: Path, target_files: list[Path], is_incremental: bool
) -> None:
    """Updates and saves persistent cache after a successful zero-violation run."""
    cache = load_linter_cache(linter_name, root_dir)
    file_hashes = cache.get("file_hashes", {})
    for p in target_files:
        rel = p.relative_to(root_dir).as_posix()
        file_hashes[rel] = compute_file_hash(p)
    head_hash = get_git_head_hash(root_dir)
    cache["linter"] = linter_name
    cache["last_git_hash"] = head_hash
    cache["file_hashes"] = file_hashes
    save_linter_cache(linter_name, root_dir, cache)
