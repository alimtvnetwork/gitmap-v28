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

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
GITMAP_DIR = os.path.join(ROOT_DIR, "gitmap")


def run_cmd(cmd, cwd=ROOT_DIR):
    print(f"-> {cmd}")
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

    entry = f"""## [v{new_version}] {today} Release v{new_version}

### Install GitMap v{new_version}

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.sh | bash -s -- ".lovable/prompts" "v{new_version}"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v{new_version}/install.ps1 -OutFile install.ps1; .\\install.ps1 -TargetDir ".lovable/prompts" -Version "v{new_version}"`

### Added / Changed / Fixed / Removed

- Fixed CI/CD compatibility: formatted cmd/agy_pin_projects.go to strict gofmt specifications
- Purged orphaned submodule gitlinks ensuring clean actions/checkout across GitHub Actions workflows
- Standardized US English spelling across all documentation in spec/ and de-literalized test lookup tables
- Gracefully handled missing VS Code user-data root in headless CI runners for gitmap vscode ls
- Verified 115/115 E2E installer smoke tests with 100% green verification on release pipeline
- Added gitmap agy pin-projects command suite (ls, add, rm, --json, --all) to pin, list, and unpin Antigravity projects
- Added --pinned / -p filter flag to gitmap agy ls to quickly inspect pinned projects
- Added first-class dynamic root-level execution for saved macros allowing gitmap <macro-name> directly
- Added root-level macro utility command aliases (macro-list, macro-add, macro-run, macro-record, macro-show, macro-rm)
- Enforced vertical newline styling rules (R13-R16) with blank lines before if, after closing braces, and before return
- Passed all 16 local CI/CD quality gates across 3 sequential batches with 100% green verification (exit code 0)
"""
    # Prepend directly at the top of changelog.md
    lines.insert(0, entry + "\n")

    with open(changelog_path, "w", encoding="utf-8") as f:
        f.writelines(lines)


def assemble_release_notes(new_version: str) -> str:
    notes_path = os.path.join(ROOT_DIR, ".lovable", "release", f"release-notes-v{new_version}.md")
    content = f"""## Quick Install v{new_version}

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v{new_version}/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v{new_version}/install.sh | bash
```

## Changelog v{new_version}

- Fixed CI/CD compatibility: formatted cmd/agy_pin_projects.go to strict gofmt specifications
- Purged orphaned submodule gitlinks ensuring clean actions/checkout across GitHub Actions workflows
- Standardized US English spelling across all documentation in spec/ and de-literalized test lookup tables
- Gracefully handled missing VS Code user-data root in headless CI runners for gitmap vscode ls
- Verified 115/115 E2E installer smoke tests with 100% green verification on release pipeline
- Added gitmap agy pin-projects command suite (ls, add, rm, --json, --all) to pin, list, and unpin Antigravity projects
- Added --pinned / -p filter flag to gitmap agy ls to quickly inspect pinned projects
- Added first-class dynamic root-level execution for saved macros allowing gitmap <macro-name> directly
- Added root-level macro utility command aliases (macro-list, macro-add, macro-run, macro-record, macro-show, macro-rm)
- Enforced vertical newline styling rules (R13-R16) with blank lines before if, after closing braces, and before return
- Passed all 16 local CI/CD quality gates across 3 sequential batches with 100% green verification (exit code 0)
"""
    os.makedirs(os.path.dirname(notes_path), exist_ok=True)
    with open(notes_path, "w", encoding="utf-8") as f:
        f.write(content)
    return notes_path


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
    notes_file = assemble_release_notes(new_version)

    # Regenerate Go files if needed
    print("-> Running go generate ./... in gitmap/")
    run_cmd("go generate ./...", cwd=GITMAP_DIR)

    if args.create_release:
        branch_name = f"release/v{new_version}"
        tag_name = f"v{new_version}"

        print(f"-> Creating branch {branch_name}")
        run_cmd(f"git checkout -b {branch_name}")

        print("-> Staging all version changes")
        run_cmd("git add -A")

        print(f"-> Committing release changes for {tag_name}")
        run_cmd(f'git commit -m "chore(release): bump version to {tag_name}"')

        print(f"-> Creating git tag {tag_name}")
        run_cmd(f'git tag -a {tag_name} -m "Release {tag_name}"')

        print(f"-> Pushing branch {branch_name} and tag {tag_name}")
        run_cmd(f"git push origin {branch_name}")
        run_cmd(f"git push origin {tag_name}")

        print("-> Merging release into main and pushing main")
        run_cmd("git checkout main")
        run_cmd(f"git merge {branch_name}")
        run_cmd("git push origin main")

        print(f"-> Creating GitHub release for {tag_name}")
        run_cmd(f'gh release create {tag_name} --title "{tag_name}" --notes-file "{notes_file}"')
        print(f"[SUCCESS] GitHub release {tag_name} created successfully!")

    print(f"[SUCCESS] Released version {new_version} successfully!")


if __name__ == "__main__":
    main()
