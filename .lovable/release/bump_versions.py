import sys
import json
import re
import os

def bump_version(current, bump_type):
    parts = current.replace("v", "").split(".")
    major, minor, patch = int(parts[0]), int(parts[1]), int(parts[2])
    if bump_type == "major":
        return f"{major + 1}.0.0"
    elif bump_type == "minor":
        return f"{major}.{minor + 1}.0"
    elif bump_type == "patch":
        return f"{major}.{minor}.{patch + 1}"
    return current

if len(sys.argv) < 3 or sys.argv[1] != "--type":
    print("Usage: bump_versions.py --type <major|minor|patch>")
    sys.exit(1)

bump_type = sys.argv[2]

with open("version.json", "r", encoding="utf-8") as f:
    vdata = json.load(f)

current = vdata.get("Version", "0.0.0")
new_version = bump_version(current, bump_type)
vdata["Version"] = new_version

with open("version.json", "w", encoding="utf-8") as f:
    json.dump(vdata, f, indent=4)

print(f"Bumped version from {current} to {new_version}")

# Update constants
const_path = "gitmap/constants/constants.go"
with open(const_path, "r", encoding="utf-8") as f:
    content = f.read()

content = re.sub(r'var Version = "[^"]+"', f'var Version = "{new_version}"', content)
content = re.sub(r'const Version = "[^"]+"', f'var Version = "{new_version}"', content)

with open(const_path, "w", encoding="utf-8") as f:
    f.write(content)

# Update readme
readme_path = "readme.md"
with open(readme_path, "r", encoding="utf-8") as f:
    content = f.read()

content = re.sub(r'Pinned version: v\d+\.\d+\.\d+', f'Pinned version: v{new_version}', content)

with open(readme_path, "w", encoding="utf-8") as f:
    f.write(content)

what_to_read_path = "what-to-read.md"
if os.path.exists(what_to_read_path):
    with open(what_to_read_path, "r", encoding="utf-8") as f:
        content = f.read()
    content = re.sub(r'Pinned version: v\d+\.\d+\.\d+', f'Pinned version: v{new_version}', content)
    with open(what_to_read_path, "w", encoding="utf-8") as f:
        f.write(content)

print("Updated version.json, constants.go, readme.md, and what-to-read.md")
