# Subtask 02: Antigravity Projects Registration & Dynamic Workspace Integration

> **Task Reference:** `215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager`  
> **Parent Spec:** [01-architecture-spec.md](file:///d:/work/gitmap/02-spec/21-app/215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager/01-architecture-spec.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `192.168.1.22`)  
> **Target User:** `a` (UID 1000, GID 1000)  
> **Execution Mode:** Python3 / SQLite3 Automation over SSH  
> **Status:** READY FOR EXECUTION  

---

## 1. Objective & Scope

Register all 74 synchronized repositories under `/home/a/git-work/` as authoritative first-class projects in the Antigravity IDE configuration registry on Ubuntu node `u1`, link historical conversations in SQLite, and verify dynamic UI dropdown population.

### Key Scope Items:
1. **Repository Discovery & Schema Synthesis:** Enumerate all 74 git workspaces in `/home/a/git-work/`, extract primary git branches (`main`, `master`), and generate RFC 3986-compliant project JSON descriptors.
2. **Deterministic UUIDv5 Generation:** Use UUIDv5 with a fixed DNS namespace based on the canonical folder URI (`file:///home/a/git-work/<repo>`) to guarantee idempotent, stable project GUIDs across re-runs.
3. **Atomic Descriptor Deployment:** Deploy all project descriptors to `/home/a/.gemini/config/projects/<project_id>.json` with correct ownership (`a:a`) and permissions (`0644`).
4. **Auxiliary Projects Synchronization:** Ensure standard auxiliary profiles (`outside-of-project.json` and `default-cli-project.json`) are present and valid.
5. **SQLite Conversation Stitching:** Execute parameterized updates against `/home/a/.gemini/antigravity/conversation_summaries.db` to link every conversation record to its corresponding project GUID, mapping unassociated threads to `outside-of-project`.
6. **Live Ingestion Verification:** Verify that Antigravity's `language_server` detects filesystem events and broadcasts project updates to the Electron UI.

---

## 2. Technical Implementation Specification

### 2.1 Project Descriptor Generator & SQLite Stitching Script

The following Python 3 script performs the discovery, JSON generation, auxiliary profile staging, and SQLite database foreign-key updates in a single automated transaction:

```python
#!/usr/bin/env python3
"""
register_antigravity_projects.py
Discovers all git repositories in /home/a/git-work/, generates standard Antigravity
project JSON descriptors, deploys auxiliary profiles, and stitches conversation_summaries.db.
"""

import os
import sys
import json
import uuid
import sqlite3
import datetime
import subprocess

WORK_DIR = "/home/a/git-work"
PROJECTS_DIR = "/home/a/.gemini/config/projects"
DB_PATH = "/home/a/.gemini/antigravity/conversation_summaries.db"
NAMESPACE = uuid.NAMESPACE_DNS

def get_default_branch(repo_path: str) -> str:
    """Detects default branch from .git/HEAD or falls back to main."""
    head_file = os.path.join(repo_path, ".git", "HEAD")
    if os.path.isfile(head_file):
        try:
            with open(head_file, "r", encoding="utf-8") as f:
                content = f.read().strip()
                if content.startswith("ref: refs/heads/"):
                    return content.replace("ref: refs/heads/", "")
        except Exception:
            pass
    return "main"

def ensure_auxiliary_projects():
    """Generates standard auxiliary descriptors if missing."""
    now_iso = datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z")
    
    # 1. outside-of-project.json
    outside_path = os.path.join(PROJECTS_DIR, "outside-of-project.json")
    outside_payload = {
        "id": "outside-of-project",
        "name": "Outside of Project",
        "settings": {
            "fileAccessPolicy": "AGENT_SETTING_POLICY_ALLOW",
            "sandboxMode": False,
            "autoExecutionPolicy": "CASCADE_COMMANDS_AUTO_EXECUTION_OFF"
        },
        "updatedAt": now_iso
    }
    with open(outside_path, "w", encoding="utf-8") as f:
        json.dump(outside_payload, f, indent=2)
        
    # 2. default-cli-project.json
    cli_path = os.path.join(PROJECTS_DIR, "default-cli-project.json")
    cli_payload = {
        "id": "default-cli-project",
        "name": "CLI Project",
        "projectResources": {}
    }
    with open(cli_path, "w", encoding="utf-8") as f:
        json.dump(cli_payload, f, indent=2)

def main():
    os.makedirs(PROJECTS_DIR, exist_ok=True)
    ensure_auxiliary_projects()
    
    now_iso = datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z")
    project_map = {} # folder_uri -> project_id
    repo_count = 0
    
    # Discover git repositories
    for entry in sorted(os.listdir(WORK_DIR)):
        repo_dir = os.path.join(WORK_DIR, entry)
        if not os.path.isdir(repo_dir) or not os.path.isdir(os.path.join(repo_dir, ".git")):
            continue
            
        folder_uri = f"file:///home/a/git-work/{entry}"
        project_uuid = str(uuid.uuid5(NAMESPACE, folder_uri))
        default_branch = get_default_branch(repo_dir)
        
        descriptor = {
            "id": project_uuid,
            "name": entry,
            "projectResources": {
                "resources": [
                    {
                        "gitFolder": {
                            "folderUri": folder_uri,
                            "defaultBranch": default_branch
                        }
                    }
                ]
            },
            "permissionGrants": {
                "permissionGrants": {
                    "allow": []
                }
            },
            "settings": {},
            "updatedAt": now_iso,
            "isWorkspaceOnly": False
        }
        
        target_json = os.path.join(PROJECTS_DIR, f"{project_uuid}.json")
        with open(target_json, "w", encoding="utf-8") as f:
            json.dump(descriptor, f, indent=2)
            
        project_map[folder_uri] = project_uuid
        repo_count += 1

    print(f"[REGISTRY] Deployed {repo_count} repository project descriptors to {PROJECTS_DIR}")

    # Stitch SQLite conversation summaries
    if os.path.exists(DB_PATH):
        conn = sqlite3.connect(DB_PATH)
        cur = conn.cursor()
        
        updated_conversations = 0
        for folder_uri, proj_id in project_map.items():
            pattern = f"%{folder_uri}%"
            cur.execute("UPDATE conversation_summaries SET project_id = ? WHERE workspace_uris LIKE ?", (proj_id, pattern))
            updated_conversations += cur.rowcount
            
        # Map remaining orphaned or blank records to outside-of-project
        cur.execute("UPDATE conversation_summaries SET project_id = 'outside-of-project' WHERE project_id = '' OR project_id IS NULL")
        orphans_resolved = cur.rowcount
        
        conn.commit()
        conn.close()
        print(f"[SQLITE] Stitched {updated_conversations} conversations to repositories, resolved {orphans_resolved} to outside-of-project.")
    else:
        print(f"[SQLITE] Warning: Database {DB_PATH} not found.")

if __name__ == "__main__":
    main()
```

---

## 3. Remote Execution Pipeline via SSH Streaming

Deploy and execute the generator over SSH using Python 3:

```powershell
# From Windows Management Host:
$scriptContent = Get-Content -Path "register_antigravity_projects.py" -Raw
$scriptContent | ssh.exe -o BatchMode=yes u1 "python3 -"
```

Or invoke the embedded action in `master-embedded-ubuntu-runner.ps1`:
```powershell
pwsh -File d:/work/repo-secrets/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action register-projects
```

---

## 4. Antigravity Dynamic Discovery Verification

Antigravity's `language_server` establishes inotify filesystem watches on `/home/a/.gemini/config/projects/`. When the JSON files are created or modified:
1. `language_server` reads each `<project_id>.json`.
2. Validates schema conformance (`projectResources.resources[].gitFolder.folderUri`).
3. Sends an IPC notification to the Electron browser window.
4. The Workspace Selector in the Antigravity header bar repopulates with all 74 project titles.
5. Clicking any project filters the left sidebar conversation history using the `project_id` key.

---

## 5. Verification Commands & Live Evidence Protocol

Execute the following commands to confirm registration and data integrity:

```powershell
# 1. Verify descriptor file count (74 repos + 2 auxiliary = 76 total)
ssh.exe u1 "ls -1 /home/a/.gemini/config/projects/*.json | wc -l"
# Expected output: >= 76

# 2. Verify auxiliary descriptors exist
ssh.exe u1 "[ -f /home/a/.gemini/config/projects/outside-of-project.json ] && [ -f /home/a/.gemini/config/projects/default-cli-project.json ] && echo 'PASS: Auxiliary descriptors present'"

# 3. Verify 0 orphaned conversation records in SQLite
ssh.exe u1 "sqlite3 /home/a/.gemini/antigravity/conversation_summaries.db \"SELECT count(*) FROM conversation_summaries WHERE project_id = '' OR project_id IS NULL;\""
# Expected output: 0

# 4. Verify conversation count associated with gitmap project
ssh.exe u1 "sqlite3 /home/a/.gemini/antigravity/conversation_summaries.db \"SELECT count(*) FROM conversation_summaries WHERE project_id != 'outside-of-project' AND project_id != '';\""
# Expected output: > 0

# 5. Verify sample project descriptor JSON syntax
ssh.exe u1 "python3 -m json.tool /home/a/.gemini/config/projects/outside-of-project.json >/dev/null && echo 'PASS: Valid JSON'"
```

---

## 6. Acceptance Criteria & Definition of Done

- [ ] **74 Project JSONs Deployed:** Exactly 74 valid repository project descriptors present in `~/.gemini/config/projects/`.
- [ ] **Auxiliary Projects Present:** `outside-of-project.json` and `default-cli-project.json` present and conforming to schema.
- [ ] **Zero Orphaned Conversations:** All rows in `conversation_summaries.db` have non-empty `project_id` matching an existing project JSON.
- [ ] **Dynamic UI Ingestion:** Antigravity IDE UI project switcher reflects the 74 repositories and accurately associates historical conversation threads.
