# Component & Migration Spec: Scripts, Tarball Staging & Verification

## 1. Migration Utilities & Deliverables
All scripts and guides for this task reside under `repo-secrets/`:

1. `repo-secrets/04-ubuntu-migration/step-by-step-transition-guide.md`
   - Detailed conceptual explanation of every stage: source discovery, path conversion rationale, SQLite preservation, tarball compression, remote transport, backup snapshotting, extraction, and post-flight verification.
2. `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py`
   - Python 3 automation script supporting `--export`, `--import`, `--sync`, `--node <alias>`, `--dry-run`, `--backup`.
3. `repo-secrets/05-scripts/migrate-cursor-memories-conversations.sh`
   - Cross-platform shell wrapper for invoking Python migration on Linux or Windows Git Bash.

---

## 2. CLI Command Line Interface & Options
```text
usage: migrate-cursor-memories-conversations.py [-h] [--export] [--import [ARCHIVE]]
                                               [--sync] [--node NODE] [--output OUTPUT]
                                               [--dry-run] [--no-backup] [--force]

Options:
  --export, -e          Stage and export tarball archive only
  --import, -i          Unpack archive onto the current machine
  --sync, -s            Full end-to-end sync: export, dispatch via SSH, and unpack
  --node NODE, -n NODE  Remote target node alias (default: u1)
  --output OUTPUT, -o   Custom output path for archive tarball
  --dry-run, -d         Simulate operations without writing changes
  --no-backup           Skip remote pre-flight snapshot creation
  --force, -f           Force overwrite existing files
```

---

## 3. Remote Verification Protocol (Node u1)
The following automated assertions are run on node `u1`:
1. **Settings & Themes**: `test -f ~/.config/Cursor/User/settings.json` and verification of `Dracula Theme`.
2. **Project Manager**: `test -f ~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json` with 49 entries mapped to `/home/a/git-work/`.
3. **Global AI Chat Database**: `test -f ~/.config/Cursor/User/globalStorage/state.vscdb` (> 100 MB) with SQLite integrity check.
4. **Conversation Search Index**: `test -f ~/.config/Cursor/User/globalStorage/conversation-search.db` (> 1 MB).
5. **Per-Project Transcripts**: Verify `~/.cursor/projects/home-a-git-work-*` directories exist.
6. **Agent Memory & Skills**: Verify `~/.cursor/agents/` and `~/.cursor/skills-cursor/` exist.
7. **Workspace Storage**: Verify `~/.config/Cursor/User/workspaceStorage/` contains workspace folders with valid `workspace.json` pointing to `/home/a/git-work/`.
