import os
import shutil

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
if os.path.exists(task_dir):
    shutil.rmtree(task_dir)

# Run same script minus the last element
with open("gen_ssh_plan.py", "r") as f:
    lines = f.readlines()

new_lines = []
for line in lines:
    if '{"domain": "Cli", "phase": "Wire+Test", "title": "Test IP rollback logic", "sym": "TestValidateAndRollbackIP", "file": "gitmap/cmd/ipchange_test.go"},' in line:
        continue # skip the 51st item
    new_lines.append(line)

with open("gen_ssh_plan2.py", "w") as f:
    f.writelines(new_lines)
