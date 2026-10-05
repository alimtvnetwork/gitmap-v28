# Subtask Plan 03: Cursor Profile & AI Memory Export and Ubuntu Loader Subsystem

- **Spec Reference:** [02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/02-component-and-cli-spec.md](../../../../02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/02-component-and-cli-spec.md)
- **Status:** Queued
- **Target Area:** `repo-secrets/05-scripts/sync-cursor-profile.py`, `repo-secrets/05-scripts/sync-cursor-profile.sh`, `repo-secrets/05-scripts/sync-cursor-profile.ps1`, `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md`, `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`

---

## 1. Objective

Architect and implement the automated **Cursor Profile & AI Memory Migration Engine** across Windows and Ubuntu fleet nodes:
1. Export user configuration from `%APPDATA%/Cursor/User` (`settings.json`, `keybindings.json`, `snippets/`, `globalStorage/`).
2. Export AI memory, agents, extensions, and skills from `%USERPROFILE%/.cursor` (`agents/`, `ai-tracking/`, `extensions/`, `plugins/`, `projects/`, `skills-cursor/`, `argv.json`, `cli-config.json`).
3. Implement a deterministic cross-platform path normalization pipeline that translates Windows path separators (`\`), drive letters (`D:\work\*`), and project folder tokens (`d-work-*`) into POSIX equivalents (`/home/a/git-work/*` and `home-a-git-work-*`).
4. Package the sanitized assets into a compressed tar archive, transfer it to Ubuntu node `u1` via SSH/SCP, and unpack it into `~/.config/Cursor/User/` and `~/.cursor/` with correct Linux permissions (`a:a`, `0755`/`0644`).
5. Author complete migration documentation in `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md`.
6. Enforce strict safety invariants: `--dry-run` simulation mode, automatic backup before overwriting, and non-destructive isolation protecting working codebases (`/home/a/git-work/`).

---

## 2. Implementation Details

### Step 1: Migration Script Interface & CLI Flags
Create `repo-secrets/05-scripts/sync-cursor-profile.py` supporting standard CLI flags:
- `--export` (`-e`): Bundle and sanitize local Windows configuration into staging tarball.
- `--import` (`-i <archive>`): Unpack and install tarball into local target directories (used on `u1`).
- `--sync` (`-s`): Orchestrate complete workflow (export -> transfer to `--node` -> remote import).
- `--node` (`-n <alias>`): Remote node alias (default: `u1`).
- `--output` (`-o <path>`): Custom output path for staging archive.
- `--dry-run` (`-d`): Preview files to process and path sanitization changes without disk writes.
- `--backup` (`-b`): Automatically snapshot existing remote directories before modification.
- `--force` (`-f`): Suppress confirmation prompts during live overwrites.

Provide companion runners:
- `repo-secrets/05-scripts/sync-cursor-profile.sh` (POSIX bash wrapper)
- `repo-secrets/05-scripts/sync-cursor-profile.ps1` (PowerShell wrapper)

### Step 2: Source Inventory & Filter Engine
1. Discover source paths dynamically:
   - Windows Config: `%APPDATA%/Cursor/User` (fallback: path under user home).
   - Windows Cursor Home: `%USERPROFILE%/.cursor` (fallback: `.cursor` in root user profile).
2. Apply file and directory filters:
   - **Include:**
     * `settings.json`, `keybindings.json`, `snippets/`
     * `globalStorage/alefragnani.project-manager/projects.json`
     * `skills-cursor/` (all 28 canonical skills)
     * `plugins/` and active manifests
     * `projects/` (agent transcripts, canvases, MCP configs)
     * `argv.json`, `cli-config.json`
   - **Exclude / Sanitize:**
     * Runtime sockets, lockfiles, temporary `.log` files
     * Mass cache: truncate or reset `statsig-cache.json` to prevent payload bloat
     * Machine-specific hardware GPU cache dirs

### Step 3: Deterministic Path Normalization Pipeline
Transform contents in-memory during tarball construction:
1. **Repository Root Translation:**
   - Map `(?i)[dD]:\\work\\([a-zA-Z0-9_\-]+)` to `/home/a/git-work/$1`
   - Map `(?i)[dD]:/work/([a-zA-Z0-9_\-]+)` to `/home/a/git-work/$1`
2. **Project Storage Directory Normalization:**
   - Rename project storage directories under `projects/`:
     `d-work-<repo>` -> `home-a-git-work-<repo>`
3. **User Profile Path Normalization:**
   - Map `(?i)[cC]:\\Users\\[a-zA-Z0-9_\-.]+\\.cursor` to `/home/a/.cursor`
   - Replace backslashes `\` with forward slashes `/` across all JSON path strings.
4. **Project Manager State (`projects.json`):**
   - Ensure all `rootPath` fields point to `/home/a/git-work/<repo>` matching the 49 active repositories.

### Step 4: Archive Generation & Remote Transport
1. Create tar archive with gzip compression: `cursor-profile-u1-transfer.tar.gz`.
2. Format layout inside archive:
   - `config_user/`: Mapped to `~/.config/Cursor/User/`
   - `cursor_home/`: Mapped to `~/.cursor/`
3. Dispatch transfer to node `u1`:
   - Transfer archive to `/tmp/cursor-profile-u1-transfer.tar.gz` over SSH/SCP or using `gitmap ssh exec u1`.
   - Transfer `sync-cursor-profile.py` if not already present on `u1`.

### Step 5: Remote Unpack, Snapshot Backup & Permissions
On node `u1`, the import routine executes:
1. **Pre-flight Snapshot Backup:**
   - If `~/.config/Cursor/User` exists, snapshot to `~/.config/Cursor/User.bak.<timestamp>`.
   - If `~/.cursor` exists, snapshot to `~/.cursor.bak.<timestamp>`.
2. **Unpack & Merge:**
   - Extract `config_user/` into `~/.config/Cursor/User/`.
   - Extract `cursor_home/` into `~/.cursor/`.
3. **Permissions Normalization:**
   - Ensure user ownership: `chown -R a:a ~/.config/Cursor ~/.cursor` (if running with sudo or user privileges).
   - Ensure directory permissions: `chmod 755` for directories.
   - Ensure file permissions: `chmod 644` for JSON/config files.

### Step 6: Migration Notes Documentation
Author `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md` detailing:
1. Executive summary of profile and AI memory migration.
2. Complete inventory table of transferred files, directories, and skill definitions.
3. Exact regex mapping tables used for path and workspace normalization.
4. Backup paths on `u1` and rollback procedure.
5. Verification commands to validate imported state on `u1`.

### Step 7: Non-Destructive Invariant & Dry-Run Mode
1. Provide comprehensive `--dry-run` output:
   - Lists all files discovered in source paths.
   - Displays planned path transformations for sample JSON records.
   - Displays planned archive size and remote destination directories.
   - Exits 0 without creating files or modifying remote state.
2. Safety Guard `isProtectedWorkDirectory`:
   - Explicitly rejects any file deletion or modification targeting `/home/a/git-work/` or repository roots.

---

## 3. Verification & Testing Protocol

1. **Local Dry-Run Verification:**
   ```bash
   python repo-secrets/05-scripts/sync-cursor-profile.py --export --dry-run
   ```
   - Assert: Discovers `%APPDATA%/Cursor/User` and `%USERPROFILE%/.cursor`, outputs planned file tree, exits 0.
2. **Local Staging Export Verification:**
   ```bash
   python repo-secrets/05-scripts/sync-cursor-profile.py --export --output /tmp/cursor-test.tar.gz
   ```
   - Assert: Tarball created, contains `config_user/` and `cursor_home/` with sanitized paths.
3. **Remote End-to-End Sync Verification on Node `u1`:**
   ```bash
   python repo-secrets/05-scripts/sync-cursor-profile.py --sync --node u1
   ```
   - Assert: Archive transferred to `u1`, unpacked into `~/.config/Cursor/User/` and `~/.cursor/`.
4. **Remote Inspection on Node `u1`:**
   - Verify `~/.config/Cursor/User/settings.json` is present and valid JSON.
   - Verify `~/.cursor/skills-cursor/` contains 28 skill folders.
   - Verify `~/.cursor/projects/` contains `home-a-git-work-*` folders.
5. **Migration Notes Verification:**
   - Verify `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md` exists and contains complete inventory.

---

## 4. Acceptance Criteria Checklist

- [ ] `sync-cursor-profile.py` implements `--export`, `--import`, `--sync`, `--node`, `--dry-run`, and `--backup` flags.
- [ ] Windows source paths `%APPDATA%/Cursor/User` and `%USERPROFILE%/.cursor` discovered correctly.
- [ ] Deterministic path transformation converts backslashes, drive prefixes, and project identifiers to POSIX paths.
- [ ] `projects.json` contains valid POSIX paths for all 49 repositories under `/home/a/git-work/*`.
- [ ] Pre-flight snapshot backups created on `u1` before extraction.
- [ ] Correct Linux permissions (`a:a`, `0755` dirs, `0644` files) enforced on `u1`.
- [ ] Comprehensive migration notes authored in `cursor-profile-transfer-notes.md`.
- [ ] Safe non-destructive behavior verified with zero mutation of working repositories.
