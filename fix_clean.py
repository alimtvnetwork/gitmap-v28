import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    lines = f.readlines()

new_lines = []
for line in lines:
    if "HelpInstall " in line: continue
    if "HelpInstall=" in line: continue
    if "HelpInstaller " in line: continue
    if "HelpInstaller=" in line: continue
    if "HelpMacro " in line: continue
    if "HelpMacro=" in line: continue
    if "HelpSchedule " in line: continue
    if "HelpSchedule=" in line: continue
    if "HelpVSCode " in line: continue
    if "HelpVSCode=" in line: continue
    if "HelpAgy " in line: continue
    if "HelpAgy=" in line: continue
    new_lines.append(line)

with open(path, "w", encoding="utf-8") as f:
    f.writelines(new_lines)
print("done")
