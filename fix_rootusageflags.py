import re

path = "gitmap/cmd/rootusageflags.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace('\tfmt.Println(constants.HelpSetup)\n', '')

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
