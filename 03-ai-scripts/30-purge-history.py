import argparse
import ctypes
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import time

STATE_FILE = ".gitmap/purge_state.json"

def run_cmd(cmd, is_check=True, is_capture_output=True, is_text=True):
    try:
        return subprocess.run(cmd, check=is_check, capture_output=is_capture_output, text=is_text)
    except subprocess.CalledProcessError as e:
        if is_check:
            print(f"Command failed: {cmd}")
            if e.stdout:
                print(f"Output: {e.stdout}")
            if e.stderr:
                print(f"Error: {e.stderr}")
            sys.exit(1)
        return e

def send_to_recycle_bin(filepath):
    if not os.path.exists(filepath):
        return False
    if sys.platform != "win32":
        try:
            os.remove(filepath)
            return True
        except OSError:
            return False
    pFrom = ctypes.create_unicode_buffer(os.path.abspath(filepath) + "\0")
    class SHFILEOPSTRUCT(ctypes.Structure):
        _fields_ = [
            ("hwnd", ctypes.c_void_p), ("wFunc", ctypes.c_uint),
            ("pFrom", ctypes.c_wchar_p), ("pTo", ctypes.c_wchar_p),
            ("fFlags", ctypes.c_uint), ("fAnyOperationsAborted", ctypes.c_int),
            ("hNameMappings", ctypes.c_void_p), ("lpszProgressTitle", ctypes.c_wchar_p)
        ]
    op = SHFILEOPSTRUCT(None, 3, ctypes.cast(pFrom, ctypes.c_wchar_p), None, 0x40 | 0x10 | 0x0400 | 0x0004, 0, None, None)
    return ctypes.windll.shell32.SHFileOperationW(ctypes.byref(op)) == 0

def validate_repo_state(norm_pattern):
    status = run_cmd(["git", "status", "--porcelain"]).stdout.strip()
    if status:
        print("Warning: Your working tree is not clean. Please commit or stash changes.")
        sys.exit(1)
    branch_res = run_cmd(["git", "branch", "--show-current"], is_check=False)
    if branch_res.returncode != 0 or not branch_res.stdout.strip():
        print("Error: Could not determine current branch. Are you in a detached HEAD?")
        sys.exit(1)
    ls_res = run_cmd(["git", "ls-files", norm_pattern])
    files = [f for f in ls_res.stdout.split("\n") if f.strip()]
    return branch_res.stdout.strip(), files

def backup_matching_files(files, temp_dir):
    os.makedirs(temp_dir, exist_ok=True)
    backed = []
    for f in files:
        if os.path.exists(f):
            dest = os.path.join(temp_dir, f)
            os.makedirs(os.path.dirname(dest), exist_ok=True)
            shutil.copy2(f, dest)
            backed.append(f)
            send_to_recycle_bin(f)
    return backed

def print_purge_warning(norm_pattern, matching_files, temp_dir, backup_branch):
    print(f"\n--- GIT HISTORY PURGE WARNING ---\nTarget pattern: {norm_pattern}")
    print(f"Matching files in working tree: {len(matching_files)}")
    for f in matching_files[:5]:
        print(f"  - {f}")
    if len(matching_files) > 5:
        print(f"  ... and {len(matching_files) - 5} more.")
    print("\nThis operation will backup files to temp, recycle originals,")
    print("rewrite history with git filter-repo, and update .gitignore.")
    print(f"Backup Temp Directory: {temp_dir}\nBackup Git Branch: {backup_branch}")

def confirm_operation(is_auto_confirm):
    if is_auto_confirm:
        return
    confirm_text = input("\nPlease type 'I confirm' to proceed: ")
    if confirm_text.strip() != "I confirm":
        print("Operation aborted.")
        sys.exit(0)

def execute_purge_rewrite(norm_pattern):
    rem_res = run_cmd(["git", "remote", "get-url", "origin"], is_check=False)
    rem_url = rem_res.stdout.strip() if rem_res.returncode == 0 else ""
    print("Rewriting history with git filter-repo...")
    run_cmd(["git", "filter-repo", "--path-glob", norm_pattern, "--invert-paths", "--force"])
    if rem_url:
        run_cmd(["git", "remote", "add", "origin", rem_url], is_check=False)
    with open(".gitignore", "a") as ig:
        ig.write(f"\n{norm_pattern}\n")
    run_cmd(["git", "add", ".gitignore"])
    run_cmd(["git", "commit", "-m", f"chore: add {norm_pattern} to .gitignore"])

def save_state(current_branch, backup_branch, temp_dir, backed_up_files, norm_pattern):
    os.makedirs(os.path.dirname(STATE_FILE), exist_ok=True)
    state = {
        "original_branch": current_branch, "backup_branch": backup_branch,
        "temp_dir": temp_dir, "files": backed_up_files, "path_pattern": norm_pattern,
    }
    with open(STATE_FILE, "w") as sf:
        json.dump(state, sf, indent=2)

def purge_history(raw_pattern, is_auto_confirm=False):
    if not raw_pattern:
        print("Error: --path is required for purging.")
        sys.exit(1)
    norm_pattern = Path(raw_pattern).as_posix()
    current_branch, matching_files = validate_repo_state(norm_pattern)
    ts = int(time.time())
    backup_branch = f"backup-purge-{ts}"
    temp_dir = os.path.join(tempfile.gettempdir(), f"gitmap_purge_{ts}")
    print_purge_warning(norm_pattern, matching_files, temp_dir, backup_branch)
    confirm_operation(is_auto_confirm)
    run_cmd(["git", "branch", backup_branch])
    backed_files = backup_matching_files(matching_files, temp_dir)
    execute_purge_rewrite(norm_pattern)
    save_state(current_branch, backup_branch, temp_dir, backed_files, norm_pattern)
    print("\n✅ Purge completed successfully!")

def restore_files_from_temp(temp_dir):
    if not (temp_dir and os.path.exists(temp_dir)):
        return
    print(f"Restoring files from temp directory: {temp_dir}...")
    for root, _, files in os.walk(temp_dir):
        for file in files:
            src = os.path.join(root, file)
            rel_path = os.path.relpath(src, temp_dir)
            os.makedirs(os.path.dirname(rel_path), exist_ok=True)
            shutil.copy2(src, rel_path)

def restore_history():
    if not os.path.exists(STATE_FILE):
        print("Error: No purge state found. Cannot restore automatically.")
        sys.exit(1)
    with open(STATE_FILE, "r") as sf:
        state = json.load(sf)
    backup_branch = state.get("backup_branch")
    if not backup_branch:
        print("Error: Invalid state file (missing backup_branch).")
        sys.exit(1)
    print(f"Restoring history from backup branch: {backup_branch}...")
    run_cmd(["git", "reset", "--hard", backup_branch])
    restore_files_from_temp(state.get("temp_dir"))
    print(f"\n✅ Restore completed successfully!\nYou can now safely delete the backup branch: git branch -D {backup_branch}")

def main():
    parser = argparse.ArgumentParser(description="Gitmap automated history purger and file remover.")
    parser.add_argument("--path", type=str, help="Pattern to purge from Git history.")
    parser.add_argument("-y", "--confirm", action="store_true", dest="is_confirm", help="Bypass confirmation prompt.")
    parser.add_argument("--restore", action="store_true", dest="is_restore", help="Restore repository state from last purge.")
    args = parser.parse_args()
    if args.is_restore:
        restore_history()
    else:
        purge_history(args.path, is_auto_confirm=args.is_confirm)

if __name__ == "__main__":
    main()
