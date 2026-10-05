# Subtask 06: Remote Ubuntu Node Healing Script in repo-secrets

> **Parent Plan:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`  
> **Status:** READY  
> **Target Files:**  
> - `repo-secrets/05-scripts/heal-u1-pull-errors.py`  
> - `repo-secrets/05-scripts/heal-u1-pull-errors.sh`  

---

## 1. Objectives

1. **Remote Healing Script Architecture:**
   Author an automated diagnostic and repair tool `repo-secrets/05-scripts/heal-u1-pull-errors.py` and companion wrapper `heal-u1-pull-errors.sh` designed to resolve desynchronized pull states on the remote Ubuntu fleet node `u1`.
2. **Database Path Normalization:**
   Connect to `u1` (or run directly on the host) and inspect SQLite database `~/.gitmap/gitmap.db`:
   - Identify repository records storing Windows backslashes (`\`) or drive letter prefixes (`D:\work\...`).
   - Create a timestamped backup (`gitmap.db.bak.<timestamp>`).
   - Convert all backslashes to POSIX forward slashes (`/`).
   - Re-map Windows root prefixes to Ubuntu workspace paths (e.g., `/home/<user>/work/` or `$GITMAP_WORKSPACE`).
   - Flag or prune records for repositories that do not physically exist on disk.
3. **Stale `.oh-my-zsh` Pruning:**
   - Detect and clean orphaned `.oh-my-zsh` installations or submodules causing slow git status checks and lock contention on `u1`.
   - Purge stale zsh locks and temporary git locks (`.git/index.lock`).
   - Ensure `.oh-my-zsh` is added to `.gitmapignore` so repository discovery ignores shell customization folders.
4. **Verification via `gitmap pa`:**
   - Execute `gitmap pa --json` (or `gitmap pa`) on node `u1`.
   - Validate that all repositories sync cleanly with zero failures.
   - Output a structured terminal summary report.

---

## 2. Implementation Steps

### Step 2.1: Python Diagnostic & Healing Script (`repo-secrets/05-scripts/heal-u1-pull-errors.py`)

1. **Structure & Pre-Flight Checks:**
   - Detect execution environment: Direct execution on `u1` or remote execution over SSH via `paramiko` / standard `ssh` subprocess.
   - Locate remote database: Default `~/.gitmap/gitmap.db`.
2. **Database Backup & Sanitation:**
   ```python
   def sanitize_gitmap_db(db_path: str, dry_run: bool = False) -> dict:
       conn = sqlite3.connect(db_path)
       cursor = conn.cursor()
       
       # Fetch records with Windows path separators
       cursor.execute("SELECT repo_name, path FROM repositories WHERE path LIKE '%\\\\%' OR path LIKE '_:%'")
       rows = cursor.fetchall()
       
       updated = 0
       for name, path in rows:
           normalized = path.replace("\\", "/")
           # Strip Windows drive letters (e.g., D:/work/ -> /home/user/work/)
           normalized = re.sub(r'^[A-Za-z]:[/\\]', '/', normalized)
           if not dry_run:
               cursor.execute("UPDATE repositories SET path = ? WHERE repo_name = ?", (normalized, name))
               updated += 1
               
       if not dry_run:
           conn.commit()
       conn.close()
       return {"found": len(rows), "updated": updated}
   ```
3. **Stale `.oh-my-zsh` & Lock Pruning:**
   - Search for `.oh-my-zsh` git repositories outside expected dotfile directories.
   - Remove stale `.git/index.lock` files older than 10 minutes.
   - Ensure `~/.gitmap/.gitmapignore` contains `.oh-my-zsh`.
4. **Automated Verification:**
   - Run `gitmap pa` remotely.
   - Check return code and verify output indicates 0 failed repositories.

### Step 2.2: Executable Bash Wrapper (`repo-secrets/05-scripts/heal-u1-pull-errors.sh`)

1. Provide a shell wrapper with usage flags:
   - `--dry-run`: Inspect and report issues without writing changes.
   - `--fix`: Apply database sanitation and cleanup.
   - `--verify`: Run `gitmap pa` after healing.
2. Ensure POSIX compliance (`#!/usr/bin/env bash`) with strict error modes (`set -euo pipefail`).

### Step 2.3: Verification & Execution Checklist

1. Run syntax verification:
   ```bash
   python3 -m py_compile repo-secrets/05-scripts/heal-u1-pull-errors.py
   bash -n repo-secrets/05-scripts/heal-u1-pull-errors.sh
   ```
2. Execute dry run against target node:
   ```bash
   python3 repo-secrets/05-scripts/heal-u1-pull-errors.py --target u1 --dry-run
   ```
3. Confirm backup creation and path normalization logic.

---

## 3. Acceptance Criteria

- [ ] `heal-u1-pull-errors.py` created with robust error handling and SQLite path translation.
- [ ] Automated database backup is performed before any SQL update is executed.
- [ ] Backslashes and Windows drive prefixes in `gitmap.db` on `u1` are converted to Linux POSIX paths.
- [ ] Stale `.oh-my-zsh` directories and dead lockfiles are pruned or ignored.
- [ ] Post-healing verification executes `gitmap pa` and asserts 0 failures.
- [ ] Scripts pass linting and syntax checks without errors.
