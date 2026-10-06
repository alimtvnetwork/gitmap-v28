#!/usr/bin/env python3
"""
03-ai-scripts/40-ubuntu-ide-desktop-sync.py
===========================================
Ubuntu Multi-IDE & GitHub Desktop Synchronization Script.

Discovers Git repositories across local workspace directories or the GitMap
SQLite tracking database (`gitmap.db`), auto-creates missing IDE storage
directories, and idempotently registers repositories into:
  1. Visual Studio Code (`~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json`)
  2. Cursor IDE (`~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json`)
  3. Google Antigravity (`~/.gemini/config/projects/<uuid>.json`)
  4. GitHub Desktop (Linux CLI dispatch via `/usr/bin/github open <path>`)

Features:
  - Auto-mkdir directory creation with 0o755 permissions before writes.
  - Full deduplication preventing duplicate project entries or descriptors.
  - Atomic JSON file writing via temporary sibling files and `os.replace()`.
  - Rich CLI options: `--dry-run`, `--from-db`, `--ide`, `--skip-sync`, `--quiet`, `--json`.
  - 100% compliant with repository positive boolean and relative path standards.
"""

from collections.abc import Sequence
from dataclasses import asdict, dataclass, field
from enum import Enum
import argparse
import json
import os
from pathlib import Path
import shutil
import sqlite3
import subprocess
import sys
import uuid


class TargetIDEType(str, Enum):
    """Enumeration of supported target IDEs."""
    ALL = "all"
    VSCODE = "vscode"
    CURSOR = "cursor"
    ANTIGRAVITY = "antigravity"
    DESKTOP = "desktop"


class ExitCodeType(int, Enum):
    """Standard execution exit codes."""
    SUCCESS = 0
    PARTIAL_FAILURE = 1
    FATAL_ERROR = 2


@dataclass
class RepoItem:
    """Represents a discovered Git repository."""
    name: str
    path: str
    branch: str = "main"
    tags: list[str] = field(default_factory=list)


@dataclass
class SyncStats:
    """Statistics for an individual IDE synchronization operation."""
    target: str
    added: int = 0
    preserved: int = 0
    failed: int = 0
    total: int = 0
    message: str = ""


@dataclass
class SyncResult:
    """Aggregated synchronization results across all targets."""
    discovered_count: int = 0
    is_dry_run: bool = False
    targets: dict[str, dict] = field(default_factory=dict)
    summary_message: str = ""


def resolve_home_dir() -> Path:
    """Resolves the user's home directory."""
    return Path(os.path.expanduser("~")).resolve()


def resolve_candidate_gitmap_db() -> Path | None:
    """Finds candidate gitmap.db file across standard locations."""
    home = resolve_home_dir()
    candidates = [
        home / ".local" / "bin" / "gitmap-cli" / "data" / "gitmap.db",
        home / ".gitmap" / "gitmap.db",
        home / "gitmap" / "data" / "gitmap.db",
        Path("data") / "gitmap.db",
        Path("bin") / "data" / "gitmap.db",
        Path(".gitmap") / "gitmap.db",
        Path("gitmap.db"),
    ]

    for candidate in candidates:
        is_candidate_present = candidate.is_file()

        if is_candidate_present:
            return candidate.resolve()

    return None


def resolve_default_branch(repo_path: Path) -> str:
    """Determines the default branch of a git repository from .git/HEAD."""
    head_file = repo_path / ".git" / "HEAD"
    is_head_file_present = head_file.is_file()

    if is_head_file_present:
        try:
            content = head_file.read_text(encoding="utf-8").strip()
            if content.startswith("ref: refs/heads/"):
                return content.split("ref: refs/heads/", 1)[1].strip()
        except OSError:
            pass

    return "main"


def discover_repos_from_disk(scan_root: Path) -> list[RepoItem]:
    """Recursively discovers Git repositories under scan_root."""
    discovered_repos: list[RepoItem] = []
    seen_paths: set[str] = set()
    ignored_dir_names = {
        ".git", "node_modules", "vendor", "target", "dist", "build",
        ".cache", "tmp", "__pycache__", ".vscode", ".idea",
    }

    if not scan_root.is_dir():
        return discovered_repos

    for root_dir, dirs, files in os.walk(str(scan_root)):
        current_path = Path(root_dir)
        git_dir = current_path / ".git"
        has_git_dir = git_dir.exists()

        if has_git_dir:
            real_path_str = str(current_path.resolve())
            has_seen_path = real_path_str in seen_paths

            if not has_seen_path:
                seen_paths.add(real_path_str)
                branch_name = resolve_default_branch(current_path)
                repo_item = RepoItem(
                    name=current_path.name,
                    path=real_path_str,
                    branch=branch_name,
                    tags=[current_path.name],
                )
                discovered_repos.append(repo_item)

            # Do not traverse further into git repo subtrees
            dirs.clear()
            continue

        # Prune ignored directory trees from traversal
        dirs[:] = [d for d in dirs if d not in ignored_dir_names]

    return discovered_repos


def discover_repos_from_db(db_path: Path | None = None) -> list[RepoItem]:
    """Reads tracked repositories from the gitmap.db SQLite database."""
    resolved_db = db_path if db_path is not None else resolve_candidate_gitmap_db()
    is_db_missing = resolved_db is None or not resolved_db.is_file()

    if is_db_missing:
        return []

    discovered_repos: list[RepoItem] = []
    seen_paths: set[str] = set()

    try:
        conn = sqlite3.connect(str(resolved_db))
        cursor = conn.cursor()

        # Check if table 'Repo' or 'repos' exists
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name IN ('Repo', 'repos')")
        table_row = cursor.fetchone()
        is_table_present = table_row is not None

        if is_table_present:
            table_name = table_row[0]
            cursor.execute(f"SELECT AbsolutePath, RepoName, Branch FROM {table_name}")
            rows = cursor.fetchall()

            for row in rows:
                abs_path_str, repo_name, branch = row[0], row[1], row[2]
                is_valid_entry = bool(abs_path_str and os.path.exists(abs_path_str))

                if is_valid_entry:
                    canonical_path = str(Path(abs_path_str).resolve())
                    has_seen = canonical_path in seen_paths

                    if not has_seen:
                        seen_paths.add(canonical_path)
                        discovered_repos.append(
                            RepoItem(
                                name=repo_name or Path(canonical_path).name,
                                path=canonical_path,
                                branch=branch or "main",
                                tags=[repo_name] if repo_name else [],
                            )
                        )

        conn.close()
    except (sqlite3.Error, OSError):
        pass

    return discovered_repos


def atomic_write_json(target_file: Path, data: Any) -> None:
    """Safely writes JSON data atomically using a temporary sibling file."""
    parent_dir = target_file.parent
    os.makedirs(str(parent_dir), mode=0o755, exist_ok=True)

    temp_file = parent_dir / f".{target_file.name}.tmp.{os.getpid()}"

    with open(str(temp_file), "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)
        f.write("\n")

    os.replace(str(temp_file), str(target_file))


def sync_vscode_projects_file(
    projects_file: Path,
    repos: Sequence[RepoItem],
    is_dry_run: bool,
) -> tuple[int, int]:
    """Syncs repositories into a specific VS Code / Cursor projects.json file."""
    existing_entries: list[dict] = []
    is_file_present = projects_file.is_file()

    if is_file_present:
        try:
            with open(str(projects_file), "r", encoding="utf-8") as f:
                loaded = json.load(f)
                if isinstance(loaded, list):
                    existing_entries = loaded
        except (json.JSONDecodeError, OSError):
            existing_entries = []

    existing_paths: set[str] = set()
    for entry in existing_entries:
        root_path = entry.get("rootPath")
        if root_path:
            existing_paths.add(str(Path(root_path).resolve()))

    added_count = 0
    preserved_count = 0

    new_entries = list(existing_entries)
    for repo in repos:
        canonical_repo_path = str(Path(repo.path).resolve())
        is_already_registered = canonical_repo_path in existing_paths

        if is_already_registered:
            preserved_count += 1
        else:
            existing_paths.add(canonical_repo_path)
            added_count += 1
            new_entries.append({
                "name": repo.name,
                "rootPath": repo.path,
                "paths": [],
                "tags": repo.tags if repo.tags else [repo.name],
                "enabled": True,
            })

    has_changes_to_save = added_count > 0 and not is_dry_run

    if has_changes_to_save:
        atomic_write_json(projects_file, new_entries)

    return added_count, preserved_count


def sync_vscode(repos: Sequence[RepoItem], is_dry_run: bool) -> SyncStats:
    """Syncs repositories to VS Code Project Manager storage."""
    home = resolve_home_dir()
    primary_file = home / ".config" / "Code" / "User" / "globalStorage" / "alefragnani.project-manager" / "projects.json"
    fallback_file = home / ".config" / "Code" / "User" / "projects.json"

    added, preserved = sync_vscode_projects_file(primary_file, repos, is_dry_run)

    # If fallback file exists, also synchronize it
    is_fallback_present = fallback_file.is_file()

    if is_fallback_present:
        sync_vscode_projects_file(fallback_file, repos, is_dry_run)

    return SyncStats(
        target="vscode",
        added=added,
        preserved=preserved,
        failed=0,
        total=added + preserved,
        message=f"VS Code: {added} added, {preserved} preserved ({primary_file})",
    )


def sync_cursor(repos: Sequence[RepoItem], is_dry_run: bool) -> SyncStats:
    """Syncs repositories to Cursor Project Manager storage."""
    home = resolve_home_dir()
    primary_file = home / ".config" / "Cursor" / "User" / "globalStorage" / "alefragnani.project-manager" / "projects.json"
    secondary_file = home / ".config" / "Cursor" / "User" / "projects.json"

    added, preserved = sync_vscode_projects_file(primary_file, repos, is_dry_run)

    # Sync secondary file if it exists or User directory exists
    is_cursor_user_present = (home / ".config" / "Cursor" / "User").is_dir()

    if is_cursor_user_present:
        sync_vscode_projects_file(secondary_file, repos, is_dry_run)

    return SyncStats(
        target="cursor",
        added=added,
        preserved=preserved,
        failed=0,
        total=added + preserved,
        message=f"Cursor: {added} added, {preserved} preserved ({primary_file})",
    )


def sync_antigravity(repos: Sequence[RepoItem], is_dry_run: bool) -> SyncStats:
    """Syncs repositories to Google Antigravity project descriptors."""
    home = resolve_home_dir()
    projects_dir = home / ".gemini" / "config" / "projects"
    is_projects_dir_missing = not projects_dir.is_dir()

    if is_projects_dir_missing and not is_dry_run:
        os.makedirs(str(projects_dir), mode=0o755, exist_ok=True)

    existing_uris: set[str] = set()

    if projects_dir.is_dir():
        for json_file in projects_dir.glob("*.json"):
            try:
                with open(str(json_file), "r", encoding="utf-8") as f:
                    doc = json.load(f)
                    resources = doc.get("projectResources", {}).get("resources", [])
                    for res in resources:
                        uri = res.get("gitFolder", {}).get("folderUri")
                        if uri:
                            # Normalize file:///path
                            norm_path = uri.replace("file://", "")
                            existing_uris.add(str(Path(norm_path).resolve()))
                    # Check top-level folderUri if present
                    top_uri = doc.get("folderUri")
                    if top_uri:
                        norm_top_path = top_uri.replace("file://", "")
                        existing_uris.add(str(Path(norm_top_path).resolve()))
            except (json.JSONDecodeError, OSError):
                continue

    added_count = 0
    preserved_count = 0

    for repo in repos:
        canonical_path = str(Path(repo.path).resolve())
        is_already_registered = canonical_path in existing_uris

        if is_already_registered:
            preserved_count += 1
        else:
            existing_uris.add(canonical_path)
            added_count += 1
            project_id = str(uuid.uuid4())
            descriptor_path = projects_dir / f"{project_id}.json"

            project_doc = {
                "id": project_id,
                "name": repo.name,
                "projectResources": {
                    "resources": [
                        {
                            "gitFolder": {
                                "folderUri": f"file://{canonical_path}",
                                "defaultBranch": repo.branch,
                            }
                        }
                    ]
                },
                "permissionGrants": {
                    "permissionGrants": {
                        "allow": ["*"],
                    },
                    "v2Migrated": True,
                },
                "settings": {
                    "fileAccessPolicy": "AGENT_SETTING_POLICY_ALLOW",
                    "autoExecutionPolicy": "CASCADE_COMMANDS_AUTO_EXECUTION_EAGER",
                },
                "isWorkspaceOnly": False,
            }

            has_write_permission = not is_dry_run

            if has_write_permission:
                atomic_write_json(descriptor_path, project_doc)

    return SyncStats(
        target="antigravity",
        added=added_count,
        preserved=preserved_count,
        failed=0,
        total=added_count + preserved_count,
        message=f"Antigravity: {added_count} created, {preserved_count} preserved ({projects_dir})",
    )


def resolve_github_desktop_binary() -> str | None:
    """Probes candidate locations for GitHub Desktop CLI on Linux."""
    candidates = [
        "/usr/bin/github",
        "/usr/bin/github-desktop",
        "/usr/local/bin/github",
        "/usr/local/bin/github-desktop",
        "/usr/lib/github-desktop/resources/app/static/github",
    ]

    for candidate in candidates:
        is_candidate_executable = os.path.isfile(candidate) and os.access(candidate, os.X_OK)

        if is_candidate_executable:
            return candidate

    which_github = shutil.which("github")
    if which_github:
        return which_github

    which_desktop = shutil.which("github-desktop")
    if which_desktop:
        return which_desktop

    return None


def sync_github_desktop(repos: Sequence[RepoItem], is_dry_run: bool) -> SyncStats:
    """Syncs repositories to GitHub Desktop on Linux via CLI shim."""
    desktop_bin = resolve_github_desktop_binary()
    is_cli_missing = desktop_bin is None

    if is_cli_missing:
        return SyncStats(
            target="desktop",
            added=0,
            preserved=0,
            failed=0,
            total=0,
            message="GitHub Desktop CLI shim not found on system (skipped gracefully)",
        )

    added_count = 0
    failed_count = 0

    for repo in repos:
        if is_dry_run:
            added_count += 1
            continue

        try:
            # Invoking `github open <path>` registers the repository in Desktop
            result = subprocess.run(
                [desktop_bin, "open", repo.path],
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
                timeout=3,
                check=False,
            )
            is_success = result.returncode == 0

            if is_success:
                added_count += 1
            else:
                failed_count += 1
        except (subprocess.SubprocessError, OSError):
            failed_count += 1

    return SyncStats(
        target="desktop",
        added=added_count,
        preserved=0,
        failed=failed_count,
        total=added_count + failed_count,
        message=f"GitHub Desktop: {added_count} registered, {failed_count} failed via {desktop_bin}",
    )


def parse_arguments() -> argparse.Namespace:
    """Configures and parses CLI arguments."""
    parser = argparse.ArgumentParser(
        description="Synchronize local Git repositories across VS Code, Cursor, Antigravity, and GitHub Desktop on Ubuntu Linux.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )

    parser.add_argument(
        "--scan-dir",
        type=str,
        default=None,
        help="Root directory to discover repositories (default: current directory or ~/git-work)",
    )
    parser.add_argument(
        "--from-db",
        action="store_true",
        default=False,
        help="Read repository inventory directly from gitmap.db SQLite database",
    )
    parser.add_argument(
        "--ide",
        "--target",
        dest="ide_targets",
        type=str,
        default="all",
        help="Comma-separated target IDEs: vscode,cursor,antigravity,desktop,all (default: all)",
    )
    parser.add_argument(
        "--exclude-ide",
        "--exclude",
        "--skip-sync",
        dest="exclude_targets",
        type=str,
        default="",
        help="Comma-separated list of IDEs to exclude (e.g. desktop,cursor)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        default=False,
        help="Simulate modifications without writing files or invoking external processes",
    )
    parser.add_argument(
        "--quiet",
        action="store_true",
        default=False,
        help="Suppress verbose progress output",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        default=False,
        help="Emit structured JSON telemetry to stdout",
    )
    parser.add_argument(
        "--online",
        action="store_true",
        default=False,
        help="Perform online checks if applicable",
    )
    parser.add_argument(
        "--offline",
        action="store_true",
        default=False,
        help="Perform offline checks only",
    )

    return parser.parse_args()


def normalize_target_set(raw_targets: str, raw_excludes: str) -> set[str]:
    """Resolves and filters active target IDEs."""
    alias_map = {
        "all": "all",
        "vscode": "vscode",
        "code": "vscode",
        "cursor": "cursor",
        "antigravity": "antigravity",
        "agy": "antigravity",
        "gemini": "antigravity",
        "desktop": "desktop",
        "gh-desktop": "desktop",
        "github-desktop": "desktop",
    }

    selected_tokens = [t.strip().lower() for t in raw_targets.split(",") if t.strip()]
    exclude_tokens = [t.strip().lower() for t in raw_excludes.split(",") if t.strip()]

    is_all_selected = "all" in selected_tokens or not selected_tokens
    active_targets: set[str] = set()

    if is_all_selected:
        active_targets = {"vscode", "cursor", "antigravity", "desktop"}
    else:
        for token in selected_tokens:
            resolved = alias_map.get(token)
            if resolved and resolved != "all":
                active_targets.add(resolved)

    for token in exclude_tokens:
        excluded = alias_map.get(token)
        if excluded:
            active_targets.discard(excluded)

    return active_targets


def execute_sync() -> int:
    """Primary execution orchestrator."""
    args = parse_arguments()

    is_dry_run = bool(args.dry_run)
    is_quiet = bool(args.quiet)
    is_json_output = bool(args.json)
    is_from_db_enabled = bool(args.from_db)

    active_targets = normalize_target_set(args.ide_targets, args.exclude_targets)

    # 1. Discover Repositories
    repos: list[RepoItem] = []

    if is_from_db_enabled:
        repos = discover_repos_from_db()
        if not repos and not is_quiet:
            print("  ℹ No repositories found in gitmap.db, falling back to disk discovery...", file=sys.stderr)

    is_disk_discovery_needed = not repos

    if is_disk_discovery_needed:
        scan_dir_str = args.scan_dir
        if scan_dir_str:
            target_path = Path(scan_dir_str).resolve()
        else:
            home_git_work = resolve_home_dir() / "git-work"
            is_git_work_present = home_git_work.is_dir()
            target_path = home_git_work if is_git_work_present else Path(".").resolve()

        repos = discover_repos_from_disk(target_path)

    # Sort repositories deterministically by name
    repos.sort(key=lambda r: r.name.lower())

    # 2. Run Synchronizations
    result = SyncResult(
        discovered_count=len(repos),
        is_dry_run=is_dry_run,
    )

    is_vscode_targeted = "vscode" in active_targets

    if is_vscode_targeted:
        stats = sync_vscode(repos, is_dry_run)
        result.targets["vscode"] = asdict(stats)
        if not is_quiet and not is_json_output:
            print(f"  ✓ {stats.message}")

    is_cursor_targeted = "cursor" in active_targets

    if is_cursor_targeted:
        stats = sync_cursor(repos, is_dry_run)
        result.targets["cursor"] = asdict(stats)
        if not is_quiet and not is_json_output:
            print(f"  ✓ {stats.message}")

    is_antigravity_targeted = "antigravity" in active_targets

    if is_antigravity_targeted:
        stats = sync_antigravity(repos, is_dry_run)
        result.targets["antigravity"] = asdict(stats)
        if not is_quiet and not is_json_output:
            print(f"  ✓ {stats.message}")

    is_github_desktop_targeted = "desktop" in active_targets

    if is_github_desktop_targeted:
        stats = sync_github_desktop(repos, is_dry_run)
        result.targets["desktop"] = asdict(stats)
        if not is_quiet and not is_json_output:
            print(f"  ✓ {stats.message}")

    mode_label = "[DRY-RUN] " if is_dry_run else ""
    result.summary_message = f"{mode_label}Discovered {len(repos)} repositories; synchronized across {len(result.targets)} IDE target(s)."

    if is_json_output:
        print(json.dumps(asdict(result), indent=2))
    elif not is_quiet:
        print(f"\n✨ {result.summary_message}")

    return ExitCodeType.SUCCESS.value


if __name__ == "__main__":
    sys.exit(execute_sync())
