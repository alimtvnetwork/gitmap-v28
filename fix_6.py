import re

# clustercommand.go
p = 'gitmap/cmd/clustercommand.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()
if 'return 0\n}' not in c:
    c = c.replace('return runId\n}', 'return runId\n}') # already has return? wait, `return runId` was there! 
    # Ah! `insertRun` had `return runId` at the end, but I replaced something earlier. Let me just check the file.
with open(p, 'w', encoding='utf8') as f: f.write(c)
