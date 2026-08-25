import os
import shutil

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
if os.path.exists(task_dir):
    shutil.rmtree(task_dir)
os.makedirs(task_dir, exist_ok=True)
