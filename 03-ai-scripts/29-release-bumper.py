#!/usr/bin/env python3
"""
03-ai-scripts/29-release-bumper.py — Canonical Release & Version Bumper.
Autonomously bumps version across all SSoT manifests:
- version.json
- package.json
- gitmap/constants/constants.go
- changelog.md
- .lovable/user-preferences
- .lovable/release/release-notes-v<version>.md

Usage:
  python 03-ai-scripts/29-release-bumper.py --bump minor
  python 03-ai-scripts/29-release-bumper.py --bump patch
  python 03-ai-scripts/29-release-bumper.py --bump major
  python 03-ai-scripts/29-release-bumper.py --version 6.184.0
"""

import argparse
from datetime import datetime, timezone
from importlib import import_module
import json
import os
from pathlib import Path
import re
import subprocess
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
CURRENT_DIR = engine.CURRENT_DIR


def parse_args():
    parser = argparse.ArgumentParser(description="Canonical Release & Version Bumper")
    parser.add_argument(
        "--bump",
        choices=["minor", "patch", "major"],
        default="minor",
        help="Type of bump to perform (default: minor)",
    )
    parser.add_argument(
        "--version",
        dest="target_version",
        default=None,
        help="Explicit version to bump to (e.g. 6.184.0)",
    )
    parser.add_argument(
        "--bullet",
        dest="bullets",
        action="append",
        default=[],
        help="Changelog bullet to include (can be specified multiple times)",
    )
    return parser.parse_args()


def compute_next_version(current_ver: str, bump_type: str) -> str:
    cleaned = current_ver.lstrip("v")
    parts = cleaned.split(".")
    if len(parts) != 3:
        raise ValueError(f"Invalid semver version: {current_ver}")

    major, minor, patch = int(parts[0]), int(parts[1]), int(parts[2])
    if bump_type == "major":
        return f"{major + 1}.0.0"
    elif bump_type == "minor":
        return f"{major}.{minor + 1}.0"
    elif bump_type == "patch":
        return f"{major}.{minor}.{patch + 1}"
    else:
        raise ValueError(f"Unknown bump type: {bump_type}")


def update_version_json(path: Path, new_version: str) -> None:
    with open(path, "r", encoding=DEFAULT_ENCODING) as f:
        content = f.read()

    content = re.sub(
        r'("Version"\s*:\s*)"[^"]+"',
        rf'\1"{new_version}"',
        content,
        count=1,
    )
    content = re.sub(
        r'("version"\s*:\s*)"[^"]+"',
        rf'\1"{new_version}"',
        content,
    )

    with open(path, "w", encoding=DEFAULT_ENCODING) as f:
        f.write(content)


def update_package_json(path: Path, new_version: str) -> None:
    with open(path, "r", encoding=DEFAULT_ENCODING) as f:
        content = f.read()

    content = re.sub(
        r'("version"\s*:\s*)"[^"]+"',
        rf'\1"{new_version}"',
        content,
        count=1,
    )

    with open(path, "w", encoding=DEFAULT_ENCODING) as f:
        f.write(content)


def update_constants_go(path: Path, new_version: str) -> None:
    with open(path, "r", encoding=DEFAULT_ENCODING) as f:
        content = f.read()

    content = re.sub(
        r'(var\s+Version\s*=\s*)"[^"]+"',
        rf'\1"{new_version}"',
        content,
        count=1,
    )

    with open(path, "w", encoding=DEFAULT_ENCODING) as f:
        f.write(content)


def update_user_preferences(path: Path, new_version: str) -> None:
    if not path.exists():
        return
    with open(path, "r", encoding=DEFAULT_ENCODING) as f:
        content = f.read()

    content = re.sub(
        r'Release Mode Active \(v[0-9.]+\)',
        f"Release Mode Active (v{new_version})",
        content,
    )

    with open(path, "w", encoding=DEFAULT_ENCODING) as f:
        f.write(content)


def update_changelog_and_release_notes(
    repo_root: Path,
    new_version: str,
    bullets: list[str],
) -> None:
    today_utc = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    default_bullets = [
        "Enhanced GitMap CLI footer when inside a git repository to display repo name, remote Git URL, active branch, latest branch, open PR count, and comprehensive branch status",
        "Resolved real-time git branch tracking and upstream sync counters (ahead/behind/up to date)",
        "Automated SSoT version bumping pipeline via 03-ai-scripts/29-release-bumper.py",
        f"Bumped version to v{new_version} across all Single Source of Truth manifests",
    ]
    actual_bullets = bullets if bullets else default_bullets

    bullet_lines = "\n".join(f"- {b}" for b in actual_bullets)

    changelog_path = repo_root / "changelog.md"
    new_entry = f"""## [v{new_version}] {today_utc} Release v{new_version}

### Install GitMap v{new_version}

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.sh | bash -s -- ".lovable/prompts" "v{new_version}"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.ps1 -OutFile install.ps1; .\\install.ps1 -TargetDir ".lovable/prompts" -Version "v{new_version}"`

### Added / Changed / Fixed / Removed

{bullet_lines}

"""
    with open(changelog_path, "r", encoding=DEFAULT_ENCODING) as f:
        existing_cl = f.read()

    with open(changelog_path, "w", encoding=DEFAULT_ENCODING) as f:
        f.write(new_entry + existing_cl)

    # Write .lovable/release/release-notes-v<version>.md
    release_dir = repo_root / ".lovable" / "release"
    release_dir.mkdir(parents=True, exist_ok=True)
    rn_path = release_dir / f"release-notes-v{new_version}.md"

    rn_content = f"""## Quick Install v{new_version}

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v{new_version}/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v{new_version}/install.sh | bash
```

## Changelog v{new_version}

{bullet_lines}
"""
    with open(rn_path, "w", encoding=DEFAULT_ENCODING) as f:
        f.write(rn_content)


def main():
    args = parse_args()
    repo_root = Path(CURRENT_DIR).resolve()

    version_json_p = repo_root / "version.json"
    package_json_p = repo_root / "package.json"
    constants_go_p = repo_root / "gitmap" / "constants" / "constants.go"
    user_prefs_p = repo_root / ".lovable" / "user-preferences"

    with open(version_json_p, "r", encoding=DEFAULT_ENCODING) as f:
        v_data = json.load(f)
        current_version = v_data.get("Version") or v_data.get("version")

    if args.target_version:
        new_version = args.target_version.lstrip("v")
    else:
        new_version = compute_next_version(current_version, args.bump)

    print("============================================================")
    print(f"Bumping GitMap Version: v{current_version} -> v{new_version}")
    print("============================================================")

    print(f"-> Updating {version_json_p.name}...")
    update_version_json(version_json_p, new_version)

    print(f"-> Updating {package_json_p.name}...")
    update_package_json(package_json_p, new_version)

    print(f"-> Updating {constants_go_p.name}...")
    update_constants_go(constants_go_p, new_version)

    print(f"-> Updating {user_prefs_p.name}...")
    update_user_preferences(user_prefs_p, new_version)

    print("-> Updating changelog.md and generating release notes...")
    update_changelog_and_release_notes(repo_root, new_version, args.bullets)

    print(f"SUCCESS: All SSoT manifests bumped successfully to v{new_version}")
    return ExitCodeType.SUCCESS.value


if __name__ == "__main__":
    sys.exit(main())
