# Cursor Profile & AI Memory Migration Transfer Notes

## 1. Executive Summary

This document records the migration specification, inventory schema, path transformation pipeline, and verification procedures for transferring Cursor IDE user profiles, settings, and AI memory assets from local development workstations to Ubuntu fleet node `u1`.

The migration engine is implemented in `repo-secrets/05-scripts/sync-cursor-profile.py` and wrapped by `repo-secrets/05-scripts/sync-cursor-profile.sh`. It enforces a strict non-destructive invariant:
- Active development workspaces under `/home/a/git-work/` are protected from modification or deletion.
- Existing user profiles on destination nodes are automatically snapshotted into timestamped backup folders (`*.bak.<timestamp>`) prior to archive extraction.
- Path normalization systematically maps Windows drive letters, path separators, and workspace identifiers to standard POSIX conventions.

---

## 2. Inventory of Transferred Assets

| Source Category | Windows Source Path | Ubuntu Target Path | Contents & Purpose |
|---|---|---|---|
| User Configuration | `%APPDATA%/Cursor/User` | `~/.config/Cursor/User/` | `settings.json` (editor configs, Dracula theme, font specs), `keybindings.json`, `snippets/`. |
| Project Manager State | `%APPDATA%/Cursor/User/globalStorage` | `~/.config/Cursor/User/globalStorage/` | `alefragnani.project-manager/projects.json` mapping all 49 repositories under `/home/a/git-work/`. |
| AI Agent Memory | `%USERPROFILE%/.cursor/agents` | `~/.cursor/agents/` | Custom agent instructions and persona prompt configurations. |
| AI Interaction Tracking | `%USERPROFILE%/.cursor/ai-tracking` | `~/.cursor/ai-tracking/` | Interaction telemetry and conversation tracking metadata. |
| Extensions & Plugins | `%USERPROFILE%/.cursor/extensions`, `plugins` | `~/.cursor/extensions/`, `~/.cursor/plugins/` | Installed Cursor plugins and active runtime manifests. |
| Workspace State | `%USERPROFILE%/.cursor/projects` | `~/.cursor/projects/` | Per-project agent transcripts and canvases renamed to `home-a-git-work-*`. |
| Canonical Skills | `%USERPROFILE%/.cursor/skills` | `~/.cursor/skills-cursor/` | 28+ canonical AI skills (`automate`, `autopilot`, `canvas`, `loop`, `goal`, etc.). |
| CLI & Runtime | `%USERPROFILE%/.cursor/argv.json`, `cli-config.json` | `~/.cursor/argv.json`, `~/.cursor/cli-config.json` | Chromium flags and Cursor CLI credentials. |

---

## 3. Path Normalization & Transformation Pipeline

During tarball staging, text and JSON records undergo deterministic transformation:

| Source Pattern | Replacement Target | Context |
|---|---|---|
| `(?i)[dD]:[/\\]work[/\\]([a-zA-Z0-9_\-]+)` | `/home/a/git-work/\1` | Workspace repository roots across configuration files. |
| `(?i)d-work-([a-zA-Z0-9_\-]+)` | `home-a-git-work-\1` | Project storage folder identifiers inside `~/.cursor/projects/`. |
| `(?i)[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+[/\\]\.cursor` | `/home/a/.cursor` | Cursor home directory references. |
| `(?i)[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+` | `/home/a` | User home directory references. |
| `statsig-cache.json` | `{}` | Reset feature flags cache to prevent payload bloat. |

---

## 4. Backup & Rollback Protocol on Node u1

### 4.1 Automatic Pre-Flight Snapshots
Before extracting incoming archives on node `u1`, the migration engine creates timestamped snapshots:
- `~/.config/Cursor/User.bak.<timestamp>/`
- `~/.cursor.bak.<timestamp>/`

### 4.2 Rollback Procedure
If rollback is required, execute the following commands on node `u1`:
```bash
# 1. Identify latest backup timestamp
ls -ld ~/.config/Cursor/User.bak.* ~/.cursor.bak.*

# 2. Restore User configuration
LATEST_USER_BAK=$(ls -td ~/.config/Cursor/User.bak.* | head -1)
rm -rf ~/.config/Cursor/User
cp -a "${LATEST_USER_BAK}" ~/.config/Cursor/User

# 3. Restore Cursor Home configuration
LATEST_HOME_BAK=$(ls -td ~/.cursor.bak.* | head -1)
rm -rf ~/.cursor
cp -a "${LATEST_HOME_BAK}" ~/.cursor
```

---

## 5. Linux Permissions Enforcement

Extracted directories and files are normalized to standard Linux POSIX permissions:
- Directory mode: `0755` (`drwxr-xr-x`)
- File mode: `0644` (`-rw-r--r--`)
- Ownership: User `a:a` via `chown -R a:a ~/.config/Cursor ~/.cursor`

---

## 6. Remote Verification Commands

Run the following checks on remote Ubuntu node `u1`:

```bash
# 1. Verify User settings and Dracula theme
test -f ~/.config/Cursor/User/settings.json && grep -q 'Dracula Theme' ~/.config/Cursor/User/settings.json && echo 'SETTINGS_OK'

# 2. Verify keybindings and snippets directory
test -f ~/.config/Cursor/User/keybindings.json && test -d ~/.config/Cursor/User/snippets && echo 'CONFIG_USER_OK'

# 3. Verify Project Manager repository indexing
grep -c 'rootPath' ~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json

# 4. Verify AI skills count
ls -1 ~/.cursor/skills-cursor/ | wc -l

# 5. Verify renamed workspace project folders
ls -d ~/.cursor/projects/home-a-git-work-* | wc -l
```
