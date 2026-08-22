import re

with open("gitmap/cmd/commitin/orchestrator/exclude.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isExcluded == false(f, rules)", "!isExcluded(f, rules)")
with open("gitmap/cmd/commitin/orchestrator/exclude.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cluster/exec_proj.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("isNotFound := foundPath == emptyString", "isMissing := foundPath == emptyString")
c = c.replace("if isNonFound {", "if isMissing {")
with open("gitmap/cluster/exec_proj.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("gitmap/cluster/ui.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("hasNoClients := numClients == DefaultProgress", "isEmpty := numClients == DefaultProgress")
c = c.replace("if isEmptyClients {", "if isEmpty {")
with open("gitmap/cluster/ui.go", "w", encoding="utf-8") as f:
    f.write(c)

print("Fixed boolean errors")
