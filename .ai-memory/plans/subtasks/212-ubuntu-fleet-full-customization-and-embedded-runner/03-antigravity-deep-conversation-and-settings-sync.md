# Subtask 212.3: Antigravity Deep Conversation History, Brain & Settings Synchronization

- **Parent Plan:** [212-ubuntu-fleet-full-customization-and-embedded-runner.md](../../pending/212-ubuntu-fleet-full-customization-and-embedded-runner.md)
- **Spec Reference:** [02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md](../../../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `scripts/sync-antigravity-deep.ps1`, `scripts/normalize-brain-paths.py`, `~/.gemini/antigravity`, `~/.antigravity_tools`

---

## 1. Context & Objective

The primary Windows workstation (`desktop-corei9-direct`) houses extensive Antigravity AI conversation transcripts, artifacts, workspace memory, and SQLite metadata across:
- `$USERPROFILE/.gemini\antigravity\conversation_summaries.db`
- `$USERPROFILE/.antigravity_tools\instances\*\home\.gemini\antigravity\brain\`
- `%APPDATA%\Antigravity\User\globalStorage\`

When migrating or operating across the Ubuntu node `U1` (`ubuntu-fleet-01`), conversations must remain accessible with full history and continuous session context. However, Windows paths (`$WORKSPACE_DIR/`, `$WORKSPACE_DIR/`, `$USERPROFILE/...`) stored in SQLite databases and JSONL transcripts cause broken links, failed workspace resolution, and agent initialization errors on Linux.

The objective is to author an automated deep synchronization script (`scripts/sync-antigravity-deep.ps1`) and path normalization engine that exports conversation histories, streams them via compressed tar over SSH, normalizes paths across SQLite tables and JSON transcripts, and verifies database integrity on `U1`.

---

## 2. Technical Architecture & File Layouts

### 2.1 Cross-OS Path Mapping Matrix

| Entity | Windows Source Representation | Ubuntu Target Representation |
| :--- | :--- | :--- |
| Workspace Root | `$WORKSPACE_DIR/` | `$HOME/git-work/` |
| Workspace URI | `$WORKSPACE_DIR/` or `$WORKSPACE_DIR/` | `file://$HOME/git-work/` |
| Specific Repo URI | `./` | `file://$HOME/git-work/gitmap` |
| User Profile | `$USERPROFILE/` | `$HOME/` |
| Brain Storage | `%USERPROFILE%\.antigravity_tools\instances\*\home\.gemini\antigravity\brain\` | `$HOME/.gemini/antigravity/brain/` |
| Summaries DB | `%USERPROFILE%\.gemini\antigravity\conversation_summaries.db` | `$HOME/.gemini/antigravity/conversation_summaries.db` |
| IDE Global Storage | `%APPDATA%\Antigravity\User\globalStorage\` | `$HOME/.config/Antigravity/User/globalStorage/` |

### 2.2 Tar Streaming Pipeline

To eliminate large intermediate zip files and Windows disk wear, synchronization uses piped tar streaming directly into remote SSH decompression:

```powershell
# Windows PowerShell tar streaming to remote Linux host
tar -czf - -C "$SourceDir" brain conversation_summaries.db | `
    ssh.exe -o BatchMode=yes a@ubuntu-fleet-01 "mkdir -p $HOME/.gemini/antigravity && tar -xzf - -C $HOME/.gemini/antigravity"
```

---

## 3. Path Normalization Subsystem

### 3.1 SQLite Database Schema & SQL Normalization

`conversation_summaries.db` contains metadata tables tracking conversation IDs, workspace folders, and summary snapshots. The script executes targeted SQL updates on Ubuntu using `sqlite3`:

```sql
-- SQLite Path Normalization Script: normalize_summaries.sql
BEGIN TRANSACTION;

-- Normalize Windows workspace paths to Linux workspace directory
UPDATE conversations 
SET workspace_path = REPLACE(workspace_path, '$WORKSPACE_DIR/', '$HOME/git-work/')
WHERE workspace_path LIKE '%$WORKSPACE_DIR/%';

UPDATE conversations 
SET workspace_path = REPLACE(workspace_path, '$WORKSPACE_DIR/', '$HOME/git-work/')
WHERE workspace_path LIKE '%$WORKSPACE_DIR/%';

UPDATE conversations 
SET workspace_path = REPLACE(workspace_path, '$WORKSPACE_DIR/', 'file://$HOME/git-work/')
WHERE workspace_path LIKE '%$WORKSPACE_DIR/%';

UPDATE conversations 
SET workspace_path = REPLACE(workspace_path, '$WORKSPACE_DIR/', 'file://$HOME/git-work/')
WHERE workspace_path LIKE '%$WORKSPACE_DIR/%';

UPDATE conversations 
SET workspace_path = REPLACE(workspace_path, '$USERPROFILE/', '$HOME/')
WHERE workspace_path LIKE '%$USERPROFILE/%';

COMMIT;
VACUUM;
```

### 3.2 JSON Lines (`transcript.jsonl`) Normalization Engine

Transcripts contain serialized JSON payloads per turn. Paths are embedded in `content`, `uri`, and tool invocation parameters. A lightweight Python script (`scripts/normalize-brain-paths.py`) or inline streaming regex transforms lines safely:

- Pattern 1: `$WORKSPACE_DIR/` $\rightarrow$ `file://$HOME/git-work/`
- Pattern 2: `$WORKSPACE_DIR/` $\rightarrow$ `file://$HOME/git-work/`
- Pattern 3: `d:\\\\work\\\\` $\rightarrow$ `$HOME/git-work/`
- Pattern 4: `C:\\\\Users\\\\Administrator\\\\` $\rightarrow$ `$HOME/`

---

## 4. Implementation Details: `scripts/sync-antigravity-deep.ps1`

The PowerShell script provides:
1. **Source Discovery**: Automatically locates the active Antigravity instance brain directory under `%USERPROFILE%\.antigravity_tools\instances\` and the global `conversation_summaries.db`.
2. **Selective Archiving**: Selects conversations modified within the last N days (or `--All` flag) to optimize bandwidth and transfer time.
3. **Piped Tar Streaming**: Transfers the filtered payload over SSH without staging large archives on local disk.
4. **Remote SQLite Execution**: Calls `sqlite3 $HOME/.gemini/antigravity/conversation_summaries.db` over SSH to run atomic normalization transactions.
5. **Remote Transcript Normalizer**: Runs path replacement across all unpacked `transcript.jsonl` and `transcript_full.jsonl` files.
6. **Permission Hardening**: Enforces `chmod 700 $HOME/.gemini` and `chmod 600` on database and credential files.
7. **Verification Gate**: Queries SQLite record counts and validates JSON syntax with `jq` to ensure zero corruption.

---

## 5. Remediation Checklist

- [ ] Create `scripts/sync-antigravity-deep.ps1` supporting parameters: `-TargetHost` (default `ubuntu-fleet-01`), `-User` (default `a`), `-Days` (default `14`), `-All`, and `-DryRun`.
- [ ] Implement local instance discovery scanning `%USERPROFILE%\.antigravity_tools\instances\*\home\.gemini\antigravity\brain`.
- [ ] Implement streaming tar pipeline over SSH directly into target directory `$HOME/.gemini/antigravity/`.
- [ ] Author remote path normalization routine for `conversation_summaries.db` replacing Windows drive and file URI conventions.
- [ ] Author remote JSONL path normalization for `transcript.jsonl` and `transcript_full.jsonl`.
- [ ] Add remote validation step verifying SQLite integrity (`PRAGMA integrity_check;`) and `jq . > /dev/null` on sample transcripts.
- [ ] Document usage, parameters, and rollback procedures in [02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md](../../../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md).
