import os
import re

c = 0
miss_file = 0
miss_sec = 0

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
for f in os.listdir(task_dir):
    with open(os.path.join(task_dir, f), 'r', encoding='utf-8') as file:
        content = file.read()
        paths = re.findall(r'(?:spec|\.lovable)/[A-Za-z0-9/._-]+', content)
        c += len(paths)
        for p in paths:
            if not os.path.exists(p):
                print(f"MISSING {p}")
                miss_file += 1

print(f"Citations total: {c}")
print(f"Missing files: {miss_file}")
print(f"Missing sections: {miss_sec}")
