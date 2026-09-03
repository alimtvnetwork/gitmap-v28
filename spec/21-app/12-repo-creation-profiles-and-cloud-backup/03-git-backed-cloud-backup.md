# Git-Backed Cloud Backup & Disaster Recovery

## Architecture & Data Flow

GitMap provides transactional cloud backup synchronization using a dedicated private Git repository:

```
[Local Machine]                                  [GitHub Private Repo]
Binary Data Dir:                                 https://github.com/<profile>/gitmap-cloud-backup
  ├── gitmap.db         ─────────────┐
  ├── git_profiles.json ─────────────┼────────>  snapshots/snapshot-YYYY-MM-DD-HHMMSS/
  ├── pipeline_db/      ─────────────┤             ├── gitmap.db
  └── repo_search/      ─────────────┘             ├── git_profiles.json
                                                   ├── pipeline_db/
                                                   ├── repo_search/
                                                   └── manifest.json
```

---

## Transaction Protocol

1. **Remote Repository Verification**:
   - Check if `<profile>/gitmap-cloud-backup` exists via `gh repo view`.
   - If missing, auto-create privately via `gh repo create <slug> --private`.
2. **Local Cache Synchronization**:
   - Maintains clone in `data/cloud-backup/`.
   - Runs `git pull --rebase origin main` before creating new snapshots.
3. **Snapshot Isolation**:
   - Each backup run creates a unique timestamped folder: `snapshots/snapshot-YYYY-MM-DD-HHMMSS/`.
   - Copies SQLite databases, split pipeline DBs, and profiles.
   - Generates `manifest.json` containing timestamp, snapshot ID, file inventory, and optional user note.
4. **Push & Tagging**:
   - Commits snapshot and pushes to `origin main`.
5. **Restoration**:
   - `gitmap backup restore [snapshot-id]` extracts snapshot files and replaces local SQLite databases and configurations after confirmation or `-y` flag.
