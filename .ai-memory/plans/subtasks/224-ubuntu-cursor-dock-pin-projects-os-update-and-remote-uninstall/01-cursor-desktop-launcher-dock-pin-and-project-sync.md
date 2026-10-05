# Subtask Plan 01: Cursor Desktop Launcher, Official Icon Deployment, GNOME Dock Pinning & Project Workspace Synchronization on Ubuntu

- **Spec Reference:** [02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/01-architecture-spec.md](../../../../02-spec/21-app/224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall/01-architecture-spec.md)
- **Status:** Queued
- **Target Area:** `repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `cli/cmdcursor/cursor_install.go`, `cli/cmdcursor/cursor_sync.go`

---

## 1. Objective

Integrate the Cursor IDE on Ubuntu node `u1` into the native desktop environment:
1. Deploy the official Cursor 512x512 icon from the extracted AppImage squashfs root to system and user icon directories.
2. Install standardized `.desktop` application launchers to system and user application paths with `StartupWMClass=Cursor`.
3. Pin Cursor to the active GNOME Shell favorites dock over the live DBus session (`DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus`).
4. Synchronize all 49 Git repositories located in `/home/a/git-work/*` into the Cursor Project Manager extension registry (`projects.json`).

---

## 2. Implementation Details

### Step 1: Official Icon Deployment
1. Source detection:
   - Primary: `/opt/cursor/squashfs-root/co.anysphere.cursor.png`
   - Secondary: `/opt/cursor/squashfs-root/usr/share/pixmaps/co.anysphere.cursor.png`
2. Destination directories:
   - System Pixmaps: `/usr/share/pixmaps/co.anysphere.cursor.png` and `/usr/share/pixmaps/cursor.png` (using `sudo -n cp` if available).
   - User Icon Theme: `~/.local/share/icons/hicolor/512x512/apps/cursor.png` and `~/.local/share/icons/hicolor/512x512/apps/co.anysphere.cursor.png`.
3. Create target directory `~/.local/share/icons/hicolor/512x512/apps` if it does not exist.
4. Update icon cache via `gtk-update-icon-cache -f -t ~/.local/share/icons/hicolor` if present.

### Step 2: Desktop Launcher Deployment (`cursor.desktop`)
1. Extract or construct desktop file adhering to official template:
   ```ini
   [Desktop Entry]
   Name=Cursor
   GenericName=AI Code Editor
   Comment=AI-powered code editor built on VS Code
   Exec=/usr/local/bin/cursor %F
   Icon=co.anysphere.cursor
   Type=Application
   StartupNotify=true
   StartupWMClass=Cursor
   Categories=Development;IDE;TextEditor;
   MimeType=text/plain;inode/directory;
   Actions=new-empty-window;

   [Desktop Action new-empty-window]
   Name=New Empty Window
   Exec=/usr/local/bin/cursor --new-window %F
   Icon=co.anysphere.cursor
   ```
2. Write launcher to:
   - User path: `~/.local/share/applications/cursor.desktop` (mode `0644`).
   - System path: `/usr/share/applications/cursor.desktop` (using `sudo -n` if available).
3. Ensure executable permission on binary wrapper (`/usr/local/bin/cursor` and `~/.local/bin/cursor`) invoking `--no-sandbox "$@"`.

### Step 3: GNOME Dock Favorite-Apps Pinning
1. Target DBus session bus address:
   `export DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$(id -u)/bus"` (or `unix:path=/run/user/1000/bus` for user `a`).
2. Read current favorites:
   `CURRENT_FAVS=$(DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/1000/bus" gsettings get org.gnome.shell favorite-apps)`
3. Parse list; if `'cursor.desktop'` is not present:
   - Append `'cursor.desktop'` to the list.
   - Commit change:
     `DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/1000/bus" gsettings set org.gnome.shell favorite-apps "$NEW_FAVS"`
4. Verify by re-reading the GSettings key.

### Step 4: Workspace Synchronization (49 Repositories)
1. Scan directory `/home/a/git-work/` for all subfolders containing a `.git` directory.
2. Resolve target configuration file:
   `~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json`.
3. If `projects.json` already exists:
   - Load existing entries and index by `rootPath`.
4. For each discovered repository:
   - Name: directory base name (e.g., `gitmap`).
   - RootPath: absolute directory path (e.g., `/home/a/git-work/gitmap`).
   - Paths: `[]`.
   - Tags: `["git-work", "remote-fleet"]`.
   - Enabled: `true`.
5. Sort array alphabetically by `name`.
6. Write out cleanly formatted JSON with 2 spaces indentation.

---

## 3. Verification Commands

Run following checks on Ubuntu node `u1` to verify complete execution:

```bash
# 1. Verify icons exist
ls -la ~/.local/share/icons/hicolor/512x512/apps/cursor.png
ls -la /usr/share/pixmaps/co.anysphere.cursor.png

# 2. Validate desktop entry format
desktop-file-validate ~/.local/share/applications/cursor.desktop
grep "StartupWMClass=Cursor" ~/.local/share/applications/cursor.desktop

# 3. Verify GNOME dock favorite-apps pinning
DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/1000/bus" gsettings get org.gnome.shell favorite-apps | grep "cursor.desktop"

# 4. Verify Project Manager JSON contains all 49 repositories
python3 -c "
import json
with open('/home/a/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json') as f:
    repos = json.load(f)
print(f'Total configured projects: {len(repos)}')
assert len(repos) >= 49, 'Expected at least 49 repositories'
"
```
