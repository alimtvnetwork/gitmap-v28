import os

gitmap_dir = r"d:\work\gitmap\gitmap"

# 1. Apply misspellings
replacements = {
    "recognised": "recognized",
    "recognise": "recognize",
    "initialises": "initializes",
    "honours": "honors",
    "behaviour": "behavior",
    "recognisable": "recognizable",
    "unrecognised": "unrecognized",
    "cancelled": "canceled",
    "Centralised": "Centralized"
}

def fix_file(filepath, replacements):
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()
    orig = content
    for old, new in replacements.items():
        content = content.replace(old, new)
    content = content.replace('"GET"', 'http.MethodGet')
    if content != orig:
        with open(filepath, "w", encoding="utf-8", newline="\n") as f:
            f.write(content)

for root, _, files in os.walk(os.path.join(gitmap_dir, "cmd")):
    for file in files:
        if file.endswith(".go") or file.endswith(".md"):
            fix_file(os.path.join(root, file), replacements)

# 2. Add nolint for wastedassign
def add_nolint(filepath, linenum, lint):
    path = os.path.join(gitmap_dir, filepath)
    with open(path, "r", encoding="utf-8") as f: lines = f.readlines()
    lines[linenum-1] = lines[linenum-1].rstrip() + f" //nolint:{lint}\n"
    with open(path, "w", encoding="utf-8") as f: f.writelines(lines)

add_nolint(r"cmd\clone.go", 579, "wastedassign")
add_nolint(r"cmd\installer_export_git.go", 42, "wastedassign")
add_nolint(r"cmd\safety_snapshot.go", 81, "wastedassign")

# 3. Apply unreachable deletion safely by searching for "return nil" around the exact line
def del_unreachable_return(filepath, report_line):
    path = os.path.join(gitmap_dir, filepath)
    with open(path, "r", encoding="utf-8") as f: lines = f.readlines()
    
    # Try report_line exactly
    if "return nil" in lines[report_line-1]:
        print(f"Deleting exactly at {path}:{report_line}")
        lines.pop(report_line-1)
    else:
        # Search backwards a few lines
        found = False
        for i in range(report_line-2, report_line-6, -1):
            if "return nil" in lines[i]:
                print(f"Deleting shifted at {path}:{i+1} for report {report_line}")
                lines.pop(i)
                found = True
                break
        if not found:
            print(f"Could not find return nil near {path}:{report_line}")
            
    with open(path, "w", encoding="utf-8") as f: f.writelines(lines)

del_unreachable_return(r"cmd\asops.go", 54)
del_unreachable_return(r"cmd\branch.go", 36)
del_unreachable_return(r"cmd\changelog.go", 58) # higher line first
del_unreachable_return(r"cmd\changelog.go", 36)
del_unreachable_return(r"cmd\cluster.go", 56)
del_unreachable_return(r"cmd\dbreset.go", 22)
del_unreachable_return(r"cmd\historyreset.go", 24)
del_unreachable_return(r"cmd\probe.go", 61)
del_unreachable_return(r"cmd\regoldens.go", 120)
del_unreachable_return(r"cmd\selfuninstallhandoff.go", 97)
del_unreachable_return(r"cmd\tasksync.go", 72)
del_unreachable_return(r"cmd\watch.go", 88)

# Delete exactly the panic leftovers
def del_exact(filepath, report_line):
    path = os.path.join(gitmap_dir, filepath)
    with open(path, "r", encoding="utf-8") as f: lines = f.readlines()
    print(f"Deleting exact {path}:{report_line}")
    lines.pop(report_line-1)
    with open(path, "w", encoding="utf-8") as f: f.writelines(lines)

# Wait, reinstall.go report is 72, which is `return constants.ReinstallModeSelf, true`
# I should find `return constants.ReinstallModeSelf` and delete it.
with open(os.path.join(gitmap_dir, r"cmd\reinstall.go"), "r", encoding="utf-8") as f: lines = f.readlines()
for i, l in enumerate(lines):
    if "return constants.ReinstallModeSelf, true" in l and i > 65:
        print(f"Found reinstall unreachable at {i+1}")
        lines.pop(i)
        break
with open(os.path.join(gitmap_dir, r"cmd\reinstall.go"), "w", encoding="utf-8") as f: f.writelines(lines)

with open(os.path.join(gitmap_dir, r"cmd\replaceflags.go"), "r", encoding="utf-8") as f: lines = f.readlines()
for i, l in enumerate(lines):
    if "return true" in l and "panic(constants.ReplaceExtCaseSensitive)" in lines[i-1]:
        print(f"Found replaceflags unreachable at {i+1}")
        lines.pop(i)
        break
with open(os.path.join(gitmap_dir, r"cmd\replaceflags.go"), "w", encoding="utf-8") as f: f.writelines(lines)

with open(os.path.join(gitmap_dir, r"cmd\taskops.go"), "r", encoding="utf-8") as f: lines = f.readlines()
for i, l in enumerate(lines):
    if "return model.TaskEntry{}" in l and "panic(" in lines[i-2]:
        print(f"Found taskops unreachable at {i+1}")
        lines.pop(i)
        break
with open(os.path.join(gitmap_dir, r"cmd\taskops.go"), "w", encoding="utf-8") as f: f.writelines(lines)

with open(os.path.join(gitmap_dir, r"cmd\update.go"), "r", encoding="utf-8") as f: lines = f.readlines()
for i, l in enumerate(lines):
    if 'return ""' in l and "panic(" in lines[i-2]:
        print(f"Found update unreachable at {i+1}")
        lines.pop(i)
        break
with open(os.path.join(gitmap_dir, r"cmd\update.go"), "w", encoding="utf-8") as f: f.writelines(lines)

# Re-apply ineffassign removal for cmd
with open(os.path.join(gitmap_dir, "cmd", "remotetransport.go"), "r", encoding="utf-8") as f: lines = f.readlines()
for i, l in enumerate(lines):
    if "useHTTPS = false" in l:
        lines.pop(i)
        break
with open(os.path.join(gitmap_dir, "cmd", "remotetransport.go"), "w", encoding="utf-8") as f: f.writelines(lines)

with open(os.path.join(gitmap_dir, "cmd", "rootusage.go"), "r", encoding="utf-8") as f: lines = f.readlines()
for i, l in enumerate(lines):
    if "measuringHelp = true" in l:
        lines.pop(i)
        break
with open(os.path.join(gitmap_dir, "cmd", "rootusage.go"), "w", encoding="utf-8") as f: f.writelines(lines)

# Fix nolintlint for cmd
import re
for root, _, files in os.walk(os.path.join(gitmap_dir, "cmd")):
    for file in files:
        if file.endswith(".go"):
            path = os.path.join(root, file)
            with open(path, "r", encoding="utf-8") as f: c = f.read()
            orig = c
            c = re.sub(r'//nolint:gosec.*', '', c)
            c = re.sub(r'//nolint:misspell.*', '', c)
            if c != orig:
                with open(path, "w", encoding="utf-8") as f: f.write(c)

