# Subtask 02: GNOME Ergonomics & Windows Keybindings

> **Task Reference:** `214-ubuntu-fleet-automation-and-workstation-governance`  
> **Parent Plan:** [.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md](./.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `ubuntu-fleet-01`)  

---

## 1. Objective & Requirements
- Programmatically set GNOME desktop text scaling to 140% (`text-scaling-factor 1.4`).
- Configure Windows muscle-memory shortcuts:
  - `Win+D`: `show-desktop`
  - `Win+Tab`: `toggle-overview`
  - `Ctrl+Win+Left` / `Right`: Workspace navigation
  - `Ctrl+Shift+D` / `Ctrl+Super+D`: Last workspace
  - `Alt+Tab`: Ungrouped individual window switcher (`switch-windows`)
- Programmatically set desktop wallpaper.

---

## 2. Implementation & Commands
- Implemented inside `run_desktop_setup()` and `run_wallpaper_setup()` in `master-embedded-ubuntu-runner.ps1`.
- D-Bus user session bridge dynamically resolved: `export DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$(id -u)/bus"`.
- Execution command:
  ```powershell
  pwsh -File $SECRETS_DIR/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action desktop
  pwsh -File $SECRETS_DIR/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action wallpaper
  ```

---

## 3. Verification & Live Evidence
- Live `gsettings get org.gnome.desktop.interface text-scaling-factor` returns `1.3999999999999999` (1.40).
- Keybindings active and verified in GNOME session.
- Status: **DONE**
