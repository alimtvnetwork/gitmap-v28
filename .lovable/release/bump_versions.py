#!/usr/bin/env python3
"""Auto-release and version bumping script for gitmap.
Usage:
    python .lovable/release/bump_versions.py --type minor [--create-release]
"""
import argparse
import datetime
import json
import os
import re
import subprocess
import sys

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
GITMAP_DIR = os.path.join(ROOT_DIR, "gitmap")


def run_cmd(cmd, cwd=ROOT_DIR):
    print(f"▶ {cmd}")
    res = subprocess.run(cmd, cwd=cwd, shell=True, capture_output=True, text=True, encoding="utf-8")
    if res.returncode != 0:
        print(f"Error executing '{cmd}':\n{res.stderr or res.stdout}", file=sys.stderr)
        sys.exit(res.returncode)
    return res.stdout.strip()


def bump_version(current: str, bump_type: str) -> str:
    parts = current.replace("v", "").split(".")
    major, minor, patch = int(parts[0]), int(parts[1]), int(parts[2])
    if bump_type == "major":
        return f"{major + 1}.0.0"
    elif bump_type == "minor":
        return f"{major}.{minor + 1}.0"
    elif bump_type == "patch":
        return f"{major}.{minor}.{patch + 1}"
    return current


def update_version_json(new_version: str) -> str:
    v_path = os.path.join(ROOT_DIR, "version.json")
    with open(v_path, "r", encoding="utf-8") as f:
        vdata = json.load(f)

    current = vdata.get("Version", "0.0.0")
    vdata["Version"] = new_version
    vdata["version"] = new_version
    vdata["updated"] = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

    with open(v_path, "w", encoding="utf-8") as f:
        json.dump(vdata, f, indent=4)
        f.write("\n")

    return current


def update_constants_go(new_version: str):
    const_path = os.path.join(GITMAP_DIR, "constants", "constants.go")
    with open(const_path, "r", encoding="utf-8") as f:
        content = f.read()

    content = re.sub(r'var Version = "[^"]+"', f'var Version = "{new_version}"', content)
    content = re.sub(r'const Version = "[^"]+"', f'var Version = "{new_version}"', content)

    with open(const_path, "w", encoding="utf-8") as f:
        f.write(content)


def update_readme_md(current: str, new_version: str):
    readme_path = os.path.join(ROOT_DIR, "readme.md")
    with open(readme_path, "r", encoding="utf-8") as f:
        content = f.read()

    content = content.replace(f"v{current}", f"v{new_version}")
    content = content.replace(current, new_version)

    # Clean any legacy pin blocks
    repl_text = (
        f'### 📌 Pinned version (v{new_version})\n\n'
        f'To pin your repository to this exact version, run the following one-liner:\n'
        f'**Unix/Bash:** `curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.sh | bash -s -- ".lovable/prompts" "v{new_version}"`\n'
        f'**PowerShell:** `Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.ps1 -OutFile install.ps1; .\\install.ps1 -TargetDir ".lovable/prompts" -Version "v{new_version}"`'
    )
    content = re.sub(
        r'### 📌 Pinned version \(v[^\)]+\)\s*\n\s*To pin your repository to this exact version, run the following one-liner:\s*\n\s*\*\*Unix/Bash:\*\*[^\n]+\n\s*\*\*PowerShell:\*\*[^\n]+',
        lambda m: repl_text,
        content
    )

    with open(readme_path, "w", encoding="utf-8") as f:
        f.write(content)


def update_changelog_md(new_version: str):
    changelog_path = os.path.join(ROOT_DIR, "changelog.md")
    today = datetime.date.today().isoformat()

    with open(changelog_path, "r", encoding="utf-8") as f:
        lines = f.readlines()

    header_index = -1
    for idx, line in enumerate(lines):
        if line.strip() == "# Changelog":
            header_index = idx
            break

    entry = f"""
## [v{new_version}] {today} Release v{new_version}

### Install gitmap v{new_version}

Unix/Bash:
`curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.sh | bash -s -- ".lovable/prompts" "v{new_version}"`

PowerShell:
`Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.ps1 -OutFile install.ps1; .\\install.ps1 -TargetDir ".lovable/prompts" -Version "v{new_version}"`

### Added / Changed / Fixed / Removed

- Fix smoke-installer.sh version extraction to support var Version alongside const Version
- Fix Git LFS binary zip false-positive in generate drift check
- Synchronize constants.Version with changelog.md and version.json
- Sanitize jq --argjson line parsing in check-single-linter-diff.sh and check-misspell-diff.sh
- Synchronize web VERSION export from version.json in src/constants/index.ts
"""
    if header_index != -1:
        lines.insert(header_index + 1, entry + "\n")
    else:
        lines.insert(0, "# Changelog\n" + entry + "\n")

    with open(changelog_path, "w", encoding="utf-8") as f:
        f.writelines(lines)


def main():
    parser = argparse.ArgumentParser(description="Bump gitmap release version")
    parser.add_argument("--type", choices=["major", "minor", "patch"], default="minor", help="Semver bump type")
    parser.add_argument("--create-release", action="store_true", help="Create branch, commit, tag, and GitHub release")
    args = parser.parse_args()

    v_path = os.path.join(ROOT_DIR, "version.json")
    with open(v_path, "r", encoding="utf-8") as f:
        vdata = json.load(f)

    current = vdata.get("Version", "0.0.0")
    new_version = bump_version(current, args.type)

    print(f"=== Bumping version from {current} to {new_version} (type: {args.type}) ===")
    update_version_json(new_version)
    update_constants_go(new_version)
    update_readme_md(current, new_version)
    update_changelog_md(new_version)

    # Regenerate Go files if needed
    print("▶ Running go generate ./... in gitmap/")
    run_cmd("go generate ./...", cwd=GITMAP_DIR)

    if args.create_release:
        branch_name = f"release/v{new_version}"
        tag_name = f"v{new_version}"

        print(f"▶ Creating branch {branch_name}")
        run_cmd(f"git checkout -b {branch_name}")

        print("▶ Staging all version changes")
        run_cmd("git add -A")

        print(f"▶ Committing release changes for {tag_name}")
        run_cmd(f'git commit -m "chore(release): bump version to {tag_name}"')

        print(f"▶ Creating git tag {tag_name}")
        run_cmd(f'git tag -a {tag_name} -m "Release {tag_name}"')

        print(f"▶ Pushing branch {branch_name} and tag {tag_name}")
        run_cmd(f"git push origin {branch_name}")
        run_cmd(f"git push origin {tag_name}")

        print("▶ Merging release into main and pushing main")
        run_cmd("git checkout main")
        run_cmd(f"git merge {branch_name}")
        run_cmd("git push origin main")

        print(f"▶ Creating GitHub release for {tag_name}")
        release_notes = f"Release {tag_name}\n\nAutomated release cut from {branch_name}."
        run_cmd(f'gh release create {tag_name} --title "gitmap {tag_name}" --notes "{release_notes}"')
        print(f"✅ GitHub release {tag_name} created successfully!")

    print(f"[SUCCESS] Released version {new_version} successfully!")


if __name__ == "__main__":
    main()
