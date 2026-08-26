import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

def append_hint(match):
    if "(use --help to expand)" in match.group(0):
        return match.group(0)
    return match.group(0)[:-1] + " (use --help to expand)\""

content = re.sub(r'(HelpSSHJoin\s*=\s*".*?")', append_hint, content)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
