# Subtask 212-01: GNOME Desktop Wallpaper CLI & Remote GUI Application Launching

> **Parent Plan:** [82-ubuntu-fleet-full-customization-and-embedded-runner.md](../../82-ubuntu-fleet-full-customization-and-embedded-runner.md)  
> **Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md)  
> **Status:** PENDING  
> **Target Node:** Ubuntu U1 (`ubuntu-fleet-01`)  
> **Target Files:**  
> - `$SECRETS_DIR/04-ubuntu-migration/set-desktop-wallpaper.sh`  
> - `$SECRETS_DIR/04-ubuntu-migration/launch-gui-app.sh`  

---

## 1. Technical Context & Scope

In the hybrid Windows 11 and Ubuntu 24.04 LTS development environment, headless SSH sessions cannot directly interact with GNOME desktop settings or launch graphical applications without connecting to the active Wayland/XWayland display session and user D-Bus socket.

This subtask delivers two hardened, modular bash utilities:
1. `set-desktop-wallpaper.sh`: Connects to user session D-Bus, normalizes file paths into standard `file://` URIs, sets both `picture-uri` (light) and `picture-uri-dark` (dark) keys in `org.gnome.desktop.background`, sets `picture-options "zoom"`, and verifies active settings.
2. `launch-gui-app.sh`: Leverages systemd user slices via `systemd-run --user` to execute GUI applications (Antigravity IDE, GNOME Text Editor, VS Code) directly into the user's active desktop session from headless SSH, detaching from the SSH shell lifecycle and preventing authorization errors.

---

## 2. Technical Specification & Component Contracts

### 2.1 GNOME Desktop Wallpaper Utility (`set-desktop-wallpaper.sh`)

#### CLI Interface
```bash
./set-desktop-wallpaper.sh [OPTIONS] <wallpaper-file-path>
```
- **Flags:**
  - `-f, --file <path>`: Explicit path to the wallpaper image (PNG, JPG, SVG, WebP).
  - `-o, --options <mode>`: Wallpaper display mode (`zoom`, `scaled`, `centered`, `stretched`, `spanned`). Default: `zoom`.
  - `-m, --mode <theme>`: Target theme mode (`both`, `dark`, `light`). Default: `both`.
  - `-h, --help`: Display usage guide and examples.

#### Environment & D-Bus Connection
- Locates user runtime D-Bus session socket:
  ```bash
  export DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$(id -u)/bus"
  ```
- Validates socket presence prior to issuing `gsettings` commands.

#### GSettings Schema & Keys
- Schema: `org.gnome.desktop.background`
  - `picture-uri`: `"file://$HOME/..."`
  - `picture-uri-dark`: `"file://$HOME/..."`
  - `picture-options`: `"zoom"`

#### Fallback & Path Normalization
- Converts relative paths to canonical absolute paths (`realpath`).
- Formats path into a compliant URI (`file://` prefix, percent-encoded spaces if necessary).
- Validates file existence and readability (`[ -f "$path" ] && [ -r "$path" ]`).

---

### 2.2 Remote GUI Application Launcher (`launch-gui-app.sh`)

#### CLI Interface
```bash
./launch-gui-app.sh [OPTIONS] <app-name-or-path> [args...]
```
- **Supported Application Aliases:**
  - `antigravity` / `agy`: Resolves `/usr/local/bin/antigravity` (or `$HOME/bin/antigravity`).
  - `editor` / `text`: Resolves `/usr/bin/gnome-text-editor`.
  - `code` / `vscode`: Resolves `/usr/bin/code` or `/snap/bin/code`.
  - `terminal`: Resolves `/usr/bin/gnome-terminal`.
  - Full binary path: Executes any specified executable.

#### Execution Architecture
Headless SSH processes lack direct access to `WAYLAND_DISPLAY` or `XAUTHORITY`. By executing via `systemd-run --user`, the application launches as a transient user unit inside the active graphical systemd session slice:
```bash
systemd-run --user \
  --unit="gui-app-${app_name}-${timestamp}" \
  --description="Remote GUI Launch: ${app_name}" \
  ${binary_path} "${app_args[@]}"
```
- Detaches from SSH connection without blocking or killing on SSH disconnect.
- Automatically inherits user desktop display variables (`DISPLAY=:0`, `WAYLAND_DISPLAY=wayland-0`).

---

## 3. Function Decomposition Plan (<= 15 Lines Per Function)

### `set-desktop-wallpaper.sh`
1. `init_dbus_session()` (<= 12 lines): Verifies `/run/user/$(id -u)/bus` and exports `DBUS_SESSION_BUS_ADDRESS`.
2. `has_dbus_socket(path)` (<= 6 lines): Returns 0 if socket exists, 1 otherwise.
3. `resolve_absolute_path(path)` (<= 8 lines): Canonicalizes path using `realpath` and asserts file readability.
4. `format_file_uri(abs_path)` (<= 8 lines): Prepends `file://` scheme to canonical path.
5. `apply_wallpaper_setting(key, uri)` (<= 10 lines): Calls `gsettings set org.gnome.desktop.background "$key" "$uri"`.
6. `apply_picture_options(option)` (<= 8 lines): Calls `gsettings set org.gnome.desktop.background picture-options "$option"`.
7. `verify_wallpaper_settings()` (<= 12 lines): Queries `gsettings get` for `picture-uri`, `picture-uri-dark`, and `picture-options`.
8. `show_wallpaper_help()` (<= 14 lines): Prints CLI usage and flag documentation.
9. `parse_wallpaper_args("$@")` (<= 15 lines): Parses positional and optional arguments.
10. `main("$@")` (<= 12 lines): Orchestrates initialization, application, and verification.

### `launch-gui-app.sh`
1. `resolve_app_binary(alias)` (<= 14 lines): Maps friendly aliases (`antigravity`, `editor`, `code`) to absolute executable paths.
2. `has_user_systemd()` (<= 8 lines): Verifies `systemctl --user is-system-running` is active or degraded.
3. `generate_unit_name(app_name)` (<= 6 lines): Creates unique transient unit name `gui-app-${app_name}-$(date +%s)`.
4. `dispatch_systemd_gui(unit, binary, args)` (<= 12 lines): Executes `systemd-run --user --unit="$unit" "$binary" "$@"`.
5. `verify_unit_spawn(unit)` (<= 10 lines): Probes transient unit status via `systemctl --user status "$unit"`.
6. `show_launch_help()` (<= 14 lines): Prints CLI usage and available aliases.
7. `main("$@")` (<= 12 lines): Validates input, resolves binary, dispatches launcher, and reports status.

---

## 4. Coding Guidelines & Invariant Rules

- **Positive Booleans Only:** Use `has_file`, `is_valid_uri`, `has_dbus_socket`, `has_running_user_slice`. Avoid negative flag naming.
- **Zero Raw Magic Strings:** Centralize gsettings schemas, unit prefixes, and default modes as script variables.
- **Strict Error Handling:** Scripts must begin with `set -euo pipefail`.
- **Zero Git Commands:** Scripts must strictly avoid modifying or executing any git commands.
- **Relative Path Hygiene:** Markdown references must use strict relative paths.

---

## 5. Verification & Acceptance Protocol

### Test Case 1: Remote Wallpaper Setting
```bash
ssh u1 "bash -s" < set-desktop-wallpaper.sh --file /usr/share/backgrounds/warty-final-ubuntu.png --options zoom
```
- **Expected Output:**
  - `DBUS_SESSION_BUS_ADDRESS set to unix:path=/run/user/1000/bus`
  - `Setting picture-uri: file:///usr/share/backgrounds/warty-final-ubuntu.png`
  - `Setting picture-uri-dark: file:///usr/share/backgrounds/warty-final-ubuntu.png`
  - `Verification: picture-uri-dark = 'file:///usr/share/backgrounds/warty-final-ubuntu.png'`
- **Exit Code:** `0`

### Test Case 2: Remote Antigravity Launch
```bash
ssh u1 "bash -s" < launch-gui-app.sh antigravity $HOME/git-work/gitmap
```
- **Expected Output:**
  - `Resolved application: /usr/local/bin/antigravity`
  - `Running as unit: gui-app-antigravity-...`
  - Antigravity window appears on active Ubuntu desktop displaying repository `$HOME/git-work/gitmap`.
- **Exit Code:** `0`

### Test Case 3: Remote Text Editor Launch
```bash
ssh u1 "bash -s" < launch-gui-app.sh editor
```
- **Expected Output:**
  - `Resolved application: /usr/bin/gnome-text-editor`
  - `Running as unit: gui-app-editor-...`
  - GNOME Text Editor launches immediately on Ubuntu display session.
- **Exit Code:** `0`
