import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# Replace `" (use --help to expand)` with ` (use --help to expand)"`
content = content.replace('" (use --help to expand)', ' (use --help to expand)"')

with open(path, "w", encoding="utf-8") as f:
    f.write(content)
print("done")
