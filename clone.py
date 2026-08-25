import os
import hashlib
from collections import defaultdict

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
buckets = defaultdict(list)

for f in os.listdir(task_dir):
    if not f.endswith(".md"): continue
    p = os.path.join(task_dir, f)
    with open(p, "r", encoding="utf-8") as file:
        lines = file.readlines()[3:] # tail -n +4 equivalent
        
    filtered = []
    for line in lines:
        if not line.startswith("**Plan") and not line.startswith("**Domain") and not line.startswith("**Target Files") and not line.startswith("**Depends On"):
            filtered.append(line)
            
    content = "".join(filtered)
    h = hashlib.sha256(content.encode()).hexdigest()[:12]
    buckets[h].append(f)

for h, files in sorted(buckets.items(), key=lambda x: len(x[1]), reverse=True):
    print(f"{len(files)} {h} {files[0]}")
    if len(files) > 1:
        print(f"FAILED! Bucket > 1 for {h}: {files}")
        
print(f"Max buckets: {max(len(x) for x in buckets.values())}")
