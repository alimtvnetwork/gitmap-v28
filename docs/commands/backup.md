# Git-Backed Cloud Backup & Disaster Recovery (`gitmap backup`)

Safeguard GitMap SQLite databases, catalog configurations, profiles, and scheduled macros by snapshotting and pushing them to a private Git backup repository on GitHub or GitLab.

---

## Command Overview

```bash
gitmap backup [subcommand] [flags]
```

### Subcommands Summary

| Subcommand | Alias | Purpose |
|------------|-------|---------|
| `create` | `push`, `now`, `run` | Create a new snapshot of all databases and push to cloud backup repo |
| `ls` | `list` | List all historical snapshots stored in the cloud backup repository |
| `restore` | - | Restore databases, profiles, and schemas from a cloud snapshot |
| `rm` | `remove`, `delete` | Delete a specific snapshot from the cloud backup repository |
| `status` | - | Inspect backup repository URL, credentials, and snapshot counts |
| `prune` | - | Prune local fix-repo snapshots by count (`--keep`) or age (`--older-than`) |

---

## How It Works

1. **Dedicated Private Repository**:
   GitMap uses the active default profile (e.g. `alimtvnetwork`) to target a dedicated private repository named `gitmap-cloud-backup` (or custom name specified via `--repo`).
2. **Auto-Provisioning**:
   If the remote repository does not exist on GitHub, GitMap provisions it as private automatically via GitHub CLI.
3. **Incremental Sync**:
   If the backup repository already exists, GitMap pulls the latest commits, writes the new timestamped snapshot under `snapshots/snapshot-YYYY-MM-DD-HHMMSS/`, creates a metadata `manifest.json` with file checksums, commits, and pushes back to remote.
4. **Disaster Recovery**:
   On a new machine or after data loss, running `gitmap backup restore` pulls the latest (or selected) snapshot and restores all SQLite databases and configurations directly into place.

---

## Practical Examples

### 1. Create and Push a Cloud Backup Snapshot

```bash
gitmap backup create --note "Full backup before system reinstall"
```

Output:
```text
  ✓ Cloud backup snapshot created and pushed successfully!
  ● Snapshot ID: snapshot-2026-09-03-184512
  ● Remote Repo: https://github.com/alimtvnetwork/gitmap-cloud-backup
  ● Timestamp:   2026-09-03 18:45:12
```

### 2. List Historical Cloud Snapshots

```bash
gitmap backup ls
```

Output:
```text
  ● Cloud Backup Snapshots
  --------------------------------------------------------------------------------
  #    Snapshot ID                Date                 Note
  --------------------------------------------------------------------------------
  [1]  snapshot-2026-09-03-184512 2026-09-03           Full backup before reinstall
  [2]  snapshot-2026-09-02-120000 2026-09-02           Weekly auto-backup
  --------------------------------------------------------------------------------
```

List as JSON:
```bash
gitmap backup ls --json
```

### 3. Restore Databases from Snapshot

```bash

# Restore by sequence number

gitmap backup restore 1 -y

# Or interactive selection:

gitmap backup restore
```

### 4. Remove Old Cloud Snapshot

```bash
gitmap backup rm snapshot-2026-09-02-120000 -y
```

### 5. Check Backup Health & Status

```bash
gitmap backup status
```

Output:
```text
  ● Cloud Backup Status
  ● Remote Repository: https://github.com/alimtvnetwork/gitmap-cloud-backup
  ● Active Profile:    alimtvnetwork (github)
  ● Total Snapshots:   2
  ● Local Cache Path:  C:\Users\...\bin\data\cloud-backup
```
