# Subtask Plan 01: Cursor Start Menu App Picker Layout Integration & Desktop Database Refresh on Ubuntu

- **Spec Reference:** [02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/01-architecture-spec.md](../../../../02-spec/21-app/225-ubuntu-cursor-start-menu-dock-position-and-profile-migration/01-architecture-spec.md)
- **Status:** Queued
- **Target Area:** `repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `repo-secrets/05-scripts/setup-cursor-ubuntu.sh`, `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`

---

## 1. Objective

Integrate Cursor into the Ubuntu Show Applications / Start menu (GNOME App Grid) on remote node `u1`:
1. Query and parse the existing GNOME App Grid layout (`org.gnome.shell app-picker-layout`) over the active DBus session.
2. Locate the maximum existing application position index on page 0 (index 17 ending with `Antigravity Manager Tools.desktop`).
3. Safely append `'cursor.desktop': <{'position': <18>}>` to the page 0 dictionary in `app-picker-layout` so Cursor is immediately pinned to Page 1 of the application launcher.
4. Enforce standard Freedesktop file permissions (`chmod 755`) on both user and system `cursor.desktop` launchers.
5. Rebuild desktop application and MIME caches using `update-desktop-database` and `gtk-update-icon-cache`.
6. Update the fleet ledger (`cursor-fleet-status.json`) with verification evidence.

---

## 2. Technical Context & Isolated Facts

- **GSettings Key:** `org.gnome.shell app-picker-layout`
- **GVariant Type:** `aa{sv}` (array of dictionaries of string to variant).
- **Existing State on Node `u1`:**
  Page 0 contains 18 items with indices 0 through 17:
  `[{'org.gnome.Calculator.desktop': <{'position': <0>}>, ..., 'Antigravity Manager Tools.desktop': <{'position': <17>}>}]`.
- **Root Cause of Missing Icon in Start Menu:**
  GNOME Shell organizes the App Grid strictly by `app-picker-layout`. Because `cursor.desktop` was omitted from this layout, GNOME Shell did not place it on Page 1 of the application launcher grid.
- **Remediation Action:**
  Appending `'cursor.desktop': <{'position': <18>}>` directly assigns Cursor slot 18 on Page 1.

---

## 3. Implementation Details

### Step 1: Active DBus Session Resolution
When executing via SSH or background automation scripts, `DBUS_SESSION_BUS_ADDRESS` is not set by default.
1. Resolve bus address dynamically:
   - Primary: Check `$DBUS_SESSION_BUS_ADDRESS`.
   - Fallback: Check `/run/user/$(id -u)/bus` or `/run/user/1000/bus`.
2. Construct execution environment:
   ```python
   env_map = os.environ.copy()
   bus_addr = resolve_dbus_bus_address()
   if bus_addr:
       env_map["DBUS_SESSION_BUS_ADDRESS"] = bus_addr
   ```

### Step 2: App Picker Layout GVariant Parsing & Injection Engine
In `repo-secrets/05-scripts/setup-cursor-ubuntu.py`:
1. **Query Existing Layout:**
   ```python
   def read_app_picker_layout(env_map: dict[str, str]) -> str:
       cmd = ["gsettings", "get", "org.gnome.shell", "app-picker-layout"]
       res = subprocess.run(cmd, env=env_map, capture_output=True, text=True, check=False)
       return res.stdout.strip()
   ```
2. **Detect Existing Presence:**
   If `'cursor.desktop'` is already present in the raw layout string, skip modification (idempotency guarantee).
3. **Calculate Target Position:**
   - Use regular expressions or string tokenization to find all patterns `<{'position': <(\d+)>}>`.
   - Calculate `max_pos = max(positions)` (expected 17).
   - Set `target_pos = max_pos + 1` (expected 18).
4. **Append Entry to Page 0:**
   - If the layout has standard structure `[{...}]`:
     Locate the closing brace `}]` of the first page array.
     Inject `, 'cursor.desktop': <{'position': <18>}>` before the closing `}]`.
   - If the layout is empty `@aa{sv} []` or `[]`:
     Set layout to `[{'cursor.desktop': <{'position': <0>}>}]`.
5. **Commit Layout via GSettings:**
   ```python
   def save_app_picker_layout(layout_str: str, env_map: dict[str, str], is_dry_run: bool) -> bool:
       if is_dry_run:
           return True
       cmd = ["gsettings", "set", "org.gnome.shell", "app-picker-layout", layout_str]
       res = subprocess.run(cmd, env=env_map, capture_output=True, text=True, check=False)
       return res.returncode == 0
   ```

### Step 3: Freedesktop Permissions Hardening
1. Ensure executable and readable permissions on desktop launchers:
   - User desktop file: `~/.local/share/applications/cursor.desktop` -> `chmod 755` (or `0644` with readable permissions).
   - System desktop file: `/usr/share/applications/cursor.desktop` -> `sudo -n chmod 755 /usr/share/applications/cursor.desktop || true`.
2. Ensure user wrapper script permissions:
   - `chmod 755 ~/.local/bin/cursor` and `sudo -n chmod 755 /usr/local/bin/cursor || true`.

### Step 4: Rebuild Desktop Application & Icon Databases
1. Execute `update-desktop-database`:
   ```python
   def refresh_desktop_database(is_dry_run: bool) -> bool:
       if is_dry_run:
           return True
       user_apps = Path.home() / ".local" / "share" / "applications"
       if user_apps.exists() and shutil.which("update-desktop-database"):
           subprocess.run(["update-desktop-database", str(user_apps)], check=False)
       if shutil.which("update-desktop-database"):
           subprocess.run(["sudo", "-n", "update-desktop-database", "/usr/share/applications"], check=False)
       if shutil.which("gtk-update-icon-cache"):
           subprocess.run(["gtk-update-icon-cache", "-f", "-t", str(Path.home() / ".local" / "share" / "icons" / "hicolor")], check=False)
       return True
   ```

### Step 5: Update Fleet Ledger Status
Extend `build_node_status_dict` in `setup-cursor-ubuntu.py`:
```python
"appPickerPinned": True,
"appPickerPosition": 18,
"desktopDatabaseRefreshed": True
```
Persist to `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`.

---

## 4. Verification Commands

Run the following checks on remote Ubuntu node `u1`:

```bash
# 1. Verify app-picker-layout includes cursor.desktop at position 18
DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/1000/bus" gsettings get org.gnome.shell app-picker-layout | grep "cursor.desktop"

# 2. Check desktop file permissions
ls -la ~/.local/share/applications/cursor.desktop
ls -la /usr/share/applications/cursor.desktop

# 3. Validate desktop file compliance
desktop-file-validate ~/.local/share/applications/cursor.desktop

# 4. Verify desktop database mimeinfo cache exists
ls -la ~/.local/share/applications/mimeinfo.cache /usr/share/applications/mimeinfo.cache
```

---

## 5. Acceptance Criteria

- [ ] `org.gnome.shell app-picker-layout` contains `'cursor.desktop': <{'position': <18>}>` (or sequential slot) on node `u1`.
- [ ] `cursor.desktop` is executable (`0755`) and passes `desktop-file-validate`.
- [ ] `update-desktop-database` runs cleanly without error on user and system application directories.
- [ ] `cursor-fleet-status.json` records `appPickerPinned: true`.
- [ ] `--dry-run` flag in `setup-cursor-ubuntu.py` simulates changes without modifying system state.
