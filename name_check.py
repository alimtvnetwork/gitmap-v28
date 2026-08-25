import os
import re

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
bad_names = []
for f in os.listdir(task_dir):
    if re.search(r'[A-Z_ ]', f):
        bad_names.append(f)
if bad_names:
    print(f"FAILED naming: {bad_names}")
else:
    print("naming OK")

seq_bad = []
for f in os.listdir(task_dir):
    if not re.match(r'^[0-9]{3}-', f):
        seq_bad.append(f)
if seq_bad:
    print(f"FAILED sequence: {seq_bad}")
else:
    print("sequence OK")
