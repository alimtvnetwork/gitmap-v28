import os
import sys
import argparse
import subprocess
import tempfile
import time
import shutil
import json
import ctypes
from pathlib import Path

STATE_FILE = ".gitmap/purge_state.json"

def run_cmd(cmd, check=True, capture_output=True, text=True):
    try:
        return subprocess.run(cmd, check=check, capture_output=capture_output, text=text, shell=True)
    except subprocess.CalledProcessError as e:
        if check:
            print(f"Command failed: {cmd}")
            if e.stdout: print(f"Output: {e.stdout}")
            if e.stderr: print(f"Error: {e.stderr}")
            sys.exit(1)
        return e

def send_to_recycle_bin(filepath):
    if not os.path.exists(filepath):
        return False
    filepath = os.path.abspath(filepath)
    FO_DELETE = 3
    FOF_ALLOWUNDO = 0x40
    FOF_NOCONFIRMATION = 0x10
    FOF_NOERRORUI = 0x0400
    FOF_SILENT = 0x0004
    
    pFrom = ctypes.create_unicode_buffer(filepath + '\0')
    
    class SHFILEOPSTRUCT(ctypes.Structure):
        _fields_ = [
            ("hwnd", ctypes.c_void_p),
            ("wFunc", ctypes.c_uint),
            ("pFrom", ctypes.c_wchar_p),
            ("pTo", ctypes.c_wchar_p),
            ("fFlags", ctypes.c_uint),
            ("fAnyOperationsAborted", ctypes.c_int),
            ("hNameMappings", ctypes.c_void_p),
            ("lpszProgressTitle", ctypes.c_wchar_p)
        ]
        
    op = SHFILEOPSTRUCT()
    op.hwnd = None
    op.wFunc = FO_DELETE
    op.pFrom = ctypes.cast(pFrom, ctypes.c_wchar_p)
    op.pTo = None
    op.fFlags = FOF_ALLOWUNDO | FOF_NOCONFIRMATION | FOF_NOERRORUI | FOF_SILENT
    op.fAnyOperationsAborted = 0
    op.hNameMappings = None
    op.lpszProgressTitle = None
    
    ret = ctypes.windll.shell32.SHFileOperationW(ctypes.byref(op))
    return ret == 0

def purge_history(args):
    if not args.path:
        print("Error: --path is required for purging.")
        sys.exit(1)
        
    path_pattern = args.path
    
    status = run_cmd("git status --porcelain").stdout.strip()
    if status:
        print("Warning: Your working tree is not clean. Please commit or stash changes before running.")
        sys.exit(1)
        
    branch_result = run_cmd("git branch --show-current", check=False)
    if branch_result.returncode != 0 or not branch_result.stdout.strip():
        print("Error: Could not determine current branch. Are you in a detached HEAD?")
        sys.exit(1)
    current_branch = branch_result.stdout.strip()
    
    ls_result = run_cmd(f"git ls-files \"{path_pattern}\"")
    matching_files = [f for f in ls_result.stdout.split('\n') if f.strip()]
    
    remote_url_result = run_cmd("git remote get-url origin", check=False)
    remote_url = remote_url_result.stdout.strip() if remote_url_result.returncode == 0 else ""
    
    timestamp = int(time.time())
    backup_branch = f"backup-purge-{timestamp}"
    temp_dir = os.path.join(tempfile.gettempdir(), f"gitmap_purge_{timestamp}")
    
    print("\n--- GIT HISTORY PURGE WARNING ---")
    print(f"Target pattern: {path_pattern}")
    print(f"Matching files in working tree: {len(matching_files)}")
    for f in matching_files[:5]:
        print(f"  - {f}")
    if len(matching_files) > 5:
        print(f"  ... and {len(matching_files) - 5} more.")
        
    print("\nThis operation will:")
    print("1. Copy the current matching files to your Windows Temp directory.")
    print("2. Send the original matching files in your repo to the Recycle Bin.")
    print("3. Create a backup branch of your current git history.")
    print("4. Rewrite your ENTIRE Git history using git filter-repo to remove the files.")
    print("5. Add the pattern to .gitignore and commit it.")
    print(f"\nBackup Temp Directory will be: {temp_dir}")
    print(f"Backup Git Branch will be: {backup_branch}")
    print("---------------------------------")
    
    if not args.confirm:
        confirm_text = input("\nPlease type 'I confirm' to proceed: ")
        if confirm_text.strip() != "I confirm":
            print("Operation aborted.")
            sys.exit(0)
            
    print("\nStarting purge process...")
    
    run_cmd(f"git branch {backup_branch}")
    print(f"Created backup branch: {backup_branch}")
    
    os.makedirs(temp_dir, exist_ok=True)
    backed_up_files = []
    
    for f in matching_files:
        if os.path.exists(f):
            dest = os.path.join(temp_dir, f)
            os.makedirs(os.path.dirname(dest), exist_ok=True)
            shutil.copy2(f, dest)
            backed_up_files.append(f)
            send_to_recycle_bin(f)
            
    if backed_up_files:
        print(f"Backed up and recycled {len(backed_up_files)} files.")
        
    print("Rewriting history with git filter-repo (this may take a moment)...")
    filter_cmd = f"git filter-repo --path-glob \"{path_pattern}\" --invert-paths --force"
    run_cmd(filter_cmd)
    
    if remote_url:
        run_cmd(f"git remote add origin {remote_url}", check=False)
        
    with open(".gitignore", "a") as ig:
        ig.write(f"\n{path_pattern}\n")
        
    run_cmd(f"git add .gitignore")
    run_cmd(f"git commit -m \"chore: add {path_pattern} to .gitignore\"")
    
    os.makedirs(os.path.dirname(STATE_FILE), exist_ok=True)
    state = {
        "original_branch": current_branch,
        "backup_branch": backup_branch,
        "temp_dir": temp_dir,
        "files": backed_up_files,
        "path_pattern": path_pattern
    }
    with open(STATE_FILE, "w") as sf:
        json.dump(state, sf, indent=2)
        
    print("\n✅ Purge completed successfully!")
    print("\nIf everything looks good, force push your changes:")
    print(f"  git push origin --force --all")
    print(f"  git push origin --force --tags")
    print(f"\nIf you need to revert this operation, run:")
    print(f"  python {sys.argv[0]} --restore")

def restore_history():
    if not os.path.exists(STATE_FILE):
        print("Error: No purge state found. Cannot restore automatically.")
        sys.exit(1)
        
    with open(STATE_FILE, "r") as sf:
        state = json.load(sf)
        
    backup_branch = state.get("backup_branch")
    temp_dir = state.get("temp_dir")
    
    if not backup_branch:
        print("Error: Invalid state file (missing backup_branch).")
        sys.exit(1)
        
    print(f"Restoring history from backup branch: {backup_branch}...")
    run_cmd(f"git reset --hard {backup_branch}")
    
    if temp_dir and os.path.exists(temp_dir):
        print(f"Restoring files from temp directory: {temp_dir}...")
        for root, dirs, files in os.walk(temp_dir):
            for file in files:
                src = os.path.join(root, file)
                rel_path = os.path.relpath(src, temp_dir)
                os.makedirs(os.path.dirname(rel_path), exist_ok=True)
                shutil.copy2(src, rel_path)
                
    print(f"\n✅ Restore completed successfully!")
    print(f"You can now safely delete the backup branch if desired: git branch -D {backup_branch}")

def main():
    parser = argparse.ArgumentParser(description="Gitmap automated history purger and file remover.")
    parser.add_argument("--path", type=str, help="The file, folder, or extension pattern (e.g., '*.db') to purge from Git history.")
    parser.add_argument("-y", "--confirm", action="store_true", help="Bypass the confirmation prompt.")
    parser.add_argument("--restore", action="store_true", help="Restore the repository state from the last purge operation.")
    
    args = parser.parse_args()
    
    if args.restore:
        restore_history()
    else:
        purge_history(args)

if __name__ == "__main__":
    main()
