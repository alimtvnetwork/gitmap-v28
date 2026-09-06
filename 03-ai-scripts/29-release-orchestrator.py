import argparse
import json
import os
import re
import subprocess
import sys
from datetime import datetime, timezone

def run_cmd(cmd):
    try:
        result = subprocess.run(cmd, shell=True, check=True, capture_output=True, text=True)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Command failed: {cmd}\nOutput: {e.output}\nStderr: {e.stderr}")
        sys.exit(1)

def get_current_branch():
    return run_cmd("git rev-parse --abbrev-ref HEAD")

def update_json_file(filepath, keys_to_update):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    modified = False
    for k, v in keys_to_update.items():
        if k in data:
            data[k] = v
            modified = True
        elif k.capitalize() in data:
            data[k.capitalize()] = v
            modified = True
        elif k.lower() in data:
            data[k.lower()] = v
            modified = True

    if modified:
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=4)
            f.write("\n") # standard ending

def bump_version_string(version, tier):
    # expect format like X.Y.Z
    parts = version.split('.')
    if len(parts) != 3:
        # fallback
        parts = [parts[0] if len(parts)>0 else "0", "0", "0"]
    major, minor, patch = int(parts[0]), int(parts[1]), int(parts[2])
    
    if tier == 'major':
        major += 1
        minor = 0
        patch = 0
    elif tier == 'minor':
        minor += 1
        patch = 0
    elif tier == 'patch':
        patch += 1
    
    return f"{major}.{minor}.{patch}"

def update_readme(filepath, old_version, new_version):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Simple replacement for badges and headers
    content = content.replace(old_version, new_version)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

def update_changelog(filepath, new_version, scope):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    today = datetime.now(timezone.utc).strftime('%Y-%m-%d')
    new_entry = f"## [{new_version}] - {today}\n\n### Changed\n- {scope}\n\n"
    
    # Try to find # Changelog or similar header
    if "# Changelog" in content:
        content = content.replace("# Changelog\n", f"# Changelog\n\n{new_entry}")
    elif "# CHANGELOG" in content:
        content = content.replace("# CHANGELOG\n", f"# CHANGELOG\n\n{new_entry}")
    else:
        # Just prepend if no header
        content = new_entry + content
        
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

def main():
    parser = argparse.ArgumentParser(description="Release Orchestrator")
    parser.add_argument("--tier", choices=['major', 'minor', 'patch'], default='minor')
    parser.add_argument("--scope", default="Automated release")
    args = parser.parse_args()

    # 1. Identify starting branch
    original_branch = get_current_branch()
    print(f"Original Branch: {original_branch}")

    try:
        # 2. Bump SemVer
        # Find current version
        current_version = "1.0.0"
        if os.path.exists("version.json"):
            with open("version.json", "r", encoding="utf-8") as f:
                d = json.load(f)
                current_version = d.get("Version", d.get("version", "1.0.0"))
        
        new_version = bump_version_string(current_version, args.tier)
        print(f"Bumping version from {current_version} to {new_version}")
        
        today_str = datetime.now(timezone.utc).strftime('%Y-%m-%d')
        keys_to_update = {"version": new_version, "releaseDate": today_str}
        
        update_json_file("version.json", keys_to_update)
        update_json_file("package.json", keys_to_update)
        update_readme("readme.md", current_version, new_version)
        update_changelog("changelog.md", new_version, args.scope)
        
        # 3. Stage & Commit
        run_cmd("git add version.json package.json readme.md changelog.md")
        # Check if there is anything to commit
        status = run_cmd("git status --porcelain")
        if status:
            run_cmd(f'git commit -m "release: v{new_version} {args.scope}"')
        
        # 4. Create Release Branch
        release_branch = f"release/v{new_version}"
        run_cmd(f"git checkout -b {release_branch}")
        
        # 5. Create Annotated Tag
        tag_name = f"v{new_version}"
        run_cmd(f'git tag -a {tag_name} -m "Release {tag_name}"')
        
        # 6. Push
        # Check if remote exists before pushing
        remotes = run_cmd("git remote")
        if remotes:
            run_cmd(f"git push origin {release_branch}")
            run_cmd(f"git push origin {tag_name}")
        else:
            print("No remote 'origin' found. Skipping push.")

        print(f"Success! Released {tag_name} on branch {release_branch}.")
        
    finally:
        # 7. MANDATORY REVERT
        run_cmd(f"git checkout {original_branch}")
        print(f"Reverted to original branch: {original_branch}")

if __name__ == '__main__':
    main()
