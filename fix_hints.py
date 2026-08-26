import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

def append_hint(match):
    return match.group(0) + " (use --help to expand)"

content = re.sub(r'(HelpSJ\s*=\s*".*?")', append_hint, content)
content = re.sub(r'(HelpInstall\s*=\s*".*?")', append_hint, content)
content = re.sub(r'(HelpInstaller\s*=\s*".*?")', append_hint, content)
content = re.sub(r'(HelpMacro\s*=\s*".*?")', append_hint, content)
content = re.sub(r'(HelpSchedule\s*=\s*".*?")', append_hint, content)
content = re.sub(r'(HelpVSCode\s*=\s*".*?")', append_hint, content)
content = re.sub(r'(HelpAgy\s*=\s*".*?")', append_hint, content)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
