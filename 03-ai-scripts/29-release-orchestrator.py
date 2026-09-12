#!/usr/bin/env python3
"""Centralized automated release orchestrator and branch lifecycle manager."""
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import json
import os
from pathlib import Path
import re
import subprocess
import sys


def run_cmd(cmd: str) -> str:
    res = subprocess.run(
        cmd,
        shell=True,
        check=True,
        capture_output=True,
        text=True,
    )

    return res.stdout.strip()


def get_current_branch() -> str:
    return run_cmd("git rev-parse --abbrev-ref HEAD")


def bump_version_string(version: str, tier: str) -> str:
    cleaned = version.lstrip("v")
    parts = cleaned.split(".")
    major = int(parts[0]) if len(parts) > 0 else 0
    minor = int(parts[1]) if len(parts) > 1 else 0
    patch = int(parts[2]) if len(parts) > 2 else 0

    if tier == "major":
        return f"{major + 1}.0.0"
    if tier == "minor":
        return f"{major}.{minor + 1}.0"
    if tier == "patch":
        return f"{major}.{minor}.{patch + 1}"

    return f"{major}.{minor}.{patch}"


def update_json_file(filepath: str, keys_to_update: dict[str, str]) -> None:
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        data = json.load(f)

    for key, val in keys_to_update.items():
        for variant in (key, key.capitalize(), key.lower()):
            if variant in data:
                data[variant] = val

    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=4)
        f.write("\n")


def update_constants_go(filepath: str, new_version: str) -> None:
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    pattern = r'(var\s+Version\s*=\s*)"[^"]+"'
    content = re.sub(pattern, rf'\1"{new_version}"', content, count=1)
    with open(filepath, "w", encoding="utf-8", newline="\n") as f:
        f.write(content)

    subprocess.run(["gofmt", "-w", filepath], check=False)


def update_readme(filepath: str, old_version: str, new_version: str) -> None:
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    content = content.replace(f"v{old_version}", f"v{new_version}")
    content = content.replace(old_version, new_version)
    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)


def generate_changelog_entry(new_version: str, today: str, bullets: list[str]) -> str:
    header = f"## [v{new_version}] {today} Release v{new_version}\n\n"
    install_hdr = f"### Install GitMap v{new_version}\n\n"
    pin_text = "To pin your repository to this exact version, run the following one-liner:\n"
    unix_cmd = f'Unix/Bash: `curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.sh | bash -s -- ".lovable/prompts" "v{new_version}"`\n'
    ps_cmd = f'PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.ps1 -OutFile install.ps1; .\\install.ps1 -TargetDir ".lovable/prompts" -Version "v{new_version}"`\n\n'
    changes_hdr = "### Added / Changed / Fixed / Removed\n\n"
    items = "\n".join(f"- {b}" for b in bullets) + "\n\n"

    return header + install_hdr + pin_text + unix_cmd + ps_cmd + changes_hdr + items


def update_changelog(filepath: str, new_version: str, bullets: list[str]) -> None:
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    today = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    entry = generate_changelog_entry(new_version, today, bullets)
    with open(filepath, "w", encoding="utf-8") as f:
        f.write(entry + content)


def update_release_notes(new_version: str, bullets: list[str]) -> str:
    release_dir = Path(".lovable") / "release"
    release_dir.mkdir(parents=True, exist_ok=True)
    rn_path = release_dir / f"release-notes-v{new_version}.md"
    bullet_lines = "\n".join(f"- {b}" for b in bullets)
    body = (
        f"## Quick Install v{new_version}\n\n"
        f"### Windows (PowerShell 5.1+)\n```powershell\n"
        f"irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v{new_version}/install.ps1 | iex\n```\n\n"
        f"### Linux / macOS (Bash)\n```bash\n"
        f"curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v{new_version}/install.sh | bash\n```\n\n"
        f"## Changelog v{new_version}\n\n{bullet_lines}\n"
    )
    with open(rn_path, "w", encoding="utf-8") as f:
        f.write(body)

    return str(rn_path).replace("\\", "/")


def update_user_preferences(filepath: str, new_version: str) -> None:
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    content = re.sub(
        r'Release Mode Active \(v[0-9.]+\)',
        f"Release Mode Active (v{new_version})",
        content,
    )
    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)


def read_canonical_version() -> str:
    if os.path.exists("version.json"):
        with open("version.json", "r", encoding="utf-8") as f:
            data = json.load(f)
            return (data.get("Version") or data.get("version") or "1.0.0").lstrip("v")

    return "1.0.0"


def stage_and_commit_release(new_version: str, scope: str, rn_path: str) -> None:
    manifests = [
        "version.json", "package.json", "readme.md", "changelog.md",
        "cli/constants/constants.go", rn_path
    ]
    if os.path.exists(".lovable/user-preferences"):
        manifests.append(".lovable/user-preferences")
    manifests_str = " ".join(manifests)
    run_cmd(f"git add {manifests_str}")
    status = run_cmd("git status --porcelain")
    if status:
        run_cmd(f'git commit -m "release: v{new_version} {scope}"')


def push_release_artifacts(release_branch: str, tag_name: str, original_branch: str) -> None:
    remotes = run_cmd("git remote")
    if "origin" in remotes.split():
        run_cmd(f"git push origin {original_branch}")
        run_cmd(f"git push origin {release_branch}")
        run_cmd(f"git push origin {tag_name}")


def create_release_branch_and_tag(new_version: str) -> tuple[str, str]:
    release_branch = f"release/v{new_version}"
    branches = run_cmd("git branch")
    if release_branch in branches.split():
        run_cmd(f"git checkout {release_branch}")
    else:
        run_cmd(f"git checkout -b {release_branch}")

    tag_name = f"v{new_version}"
    tags = run_cmd("git tag")
    if tag_name not in tags.split():
        run_cmd(f'git tag -a {tag_name} -m "Release {tag_name}"')

    return release_branch, tag_name


def verify_pre_release_quality_gates(run_tests: bool = True) -> None:
    """Verifies that 100% of CI/CD quality gates and unit tests pass green before cutting a release."""
    if not run_tests:
        print("Skipping pre-release CI/CD runner (--skip-tests active).")
        return

    print("Running pre-release quality gate checks via 03-ai-scripts/06-cicd-local-runner.py...")
    cmd = [sys.executable, "03-ai-scripts/06-cicd-local-runner.py", "--run-tests"]
    res = subprocess.run(cmd, check=False)
    if res.returncode != 0:
        raise RuntimeError(f"Pre-release quality gates failed with exit code {res.returncode}. Release aborted.")
    print("Pre-release quality gates passed 100% green.")


def parse_cli_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Release Orchestrator")
    parser.add_argument("--tier", choices=["major", "minor", "patch"], default="minor")
    parser.add_argument("--scope", default="Automated release orchestration")
    parser.add_argument("--bullet", dest="bullets", action="append", default=[])
    parser.add_argument("--skip-tests", dest="skip_tests", action="store_true", help="Skip running pre-release unit tests.")

    return parser.parse_args()


def execute_release(tier: str, scope: str, bullets: list[str], original_branch: str, skip_tests: bool = False) -> None:
    verify_pre_release_quality_gates(run_tests=not skip_tests)
    cur_ver = read_canonical_version()
    new_ver = bump_version_string(cur_ver, tier)
    print(f"Bumping version from {cur_ver} to {new_ver}")

    actual_bullets = bullets if bullets else [scope]

    today_str = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    keys_to_update = {"version": new_ver, "releaseDate": today_str}
    update_json_file("version.json", keys_to_update)
    update_json_file("package.json", keys_to_update)
    update_constants_go("cli/constants/constants.go", new_ver)
    update_readme("readme.md", cur_ver, new_ver)
    update_user_preferences(".lovable/user-preferences", new_ver)
    update_changelog("changelog.md", new_ver, actual_bullets)
    rn_path = update_release_notes(new_ver, actual_bullets)

    stage_and_commit_release(new_ver, scope, rn_path)
    rel_branch, tag_name = create_release_branch_and_tag(new_ver)
    push_release_artifacts(rel_branch, tag_name, original_branch)
    print(f"Success! Released {tag_name} on branch {rel_branch}.")


def main() -> None:
    args = parse_cli_args()
    orig_branch = get_current_branch()
    print(f"Original Branch: {orig_branch}")
    try:
        execute_release(args.tier, args.scope, args.bullets, orig_branch, skip_tests=args.skip_tests)
    finally:
        run_cmd(f"git checkout {orig_branch}")
        print(f"Reverted to original branch: {orig_branch}")


if __name__ == "__main__":
    main()
