# Plan: WinUtil & LinUtil Native Integration: OS Auto-Login (Windows & Ubuntu), System Tweaks & Maintenance

- **Slug:** `51-winutil-and-linutil-os-autologin-and-tweaks`
- **Status:** Pending
- **Author:** Antigravity (Pair Programming with Alim Karim)
- **Source Repositories Analyzed:** `D:\work\chris\winutil` and `D:\work\chris\linutil`
- **Target Architecture:** Native Go subsystem under `cli/cmdos/` with zero PowerShell dependency, hermetic testing, and two-column terminal menus.

---

## 1. Executive Summary & Insights from Chris Titus Repositories

After analyzing `D:\work\chris\winutil` and `D:\work\chris\linutil`, we identified high-value system management components that can be natively ported to GitMap. Rather than relying on external executables (`autologon.exe`), slow PowerShell cmdlets (`Invoke-WPFPanelAutologin`), or shell scripts, GitMap will implement these capabilities in **pure, high-performance Go**:

1. **OS Auto-Login (Windows & Ubuntu):**
   - **WinUtil Approach:** Downloads Sysinternals `Autologon.exe` (or sets Winlogon registry keys). Requires 3 parameters: `username`, `domain`, and `password`.
   - **GitMap Native Go Implementation:** Interacts directly with Windows Registry (`HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`) using `golang.org/x/sys/windows/registry`. For Ubuntu/Debian, detects Display Manager (`gdm3`, `lightdm`, `sddm`, or systemd console `agetty`) and updates configuration files atomically.
2. **Windows Desktop & Explorer Tweaks:**
   - **Classic Context Menu:** Toggles Windows 11 full context menu via CLSID `{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}` in `HKCU\Software\Classes\CLSID`. Restarts explorer gracefully.
   - **Old Start Menu:** Configures Windows Feature Management override `HKLM\SYSTEM\ControlSet001\Control\FeatureManagement\Overrides\8\3036241548`.
   - **Power Settings & Ultimate Performance:** Activates the hidden Windows "Ultimate Performance" power scheme (`e9a42b02-d5df-448d-aa00-03f14749eb61`) and manages hibernation (`powercfg /hibernate off`) to reclaim RAM-sized disk space.
3. **Linux Maintenance & System Updates (LinUtil Parity):**
   - **Deep Cleanup:** Package manager cache purging (`apt-get clean`, `autoremove -y`, `dnf clean all`, `pacman -Sc`), journal vacuuming (`journalctl --vacuum-time=3d`), and log truncation.
   - **Unified Multi-Distro Update:** Coordinated package and repository updates across APT, DNF, Pacman, Snap, and Flatpak.
4. **Terminal UI & Interactive Checkbox Selection:**
   - Adapts LinUtil's category-based TUI into GitMap's `termhelp` / `termtable` interactive terminal interface.

---

## 2. Command Architecture & CLI Specifications

```bash
# ── OS Auto-Login Commands ──────────────────────────────────────────────────
gitmap os autologin                        # Interactive prompt for 3 parameters: User, Domain, Password
gitmap os autologin enable -u <user> -d <domain> -p <pass>  # Direct parameter execution
gitmap os autologin disable                # Disable auto-login and revert to secure login
gitmap os autologin status                 # Inspect current auto-login configuration (masked password)

# ── Windows System Tweaks ───────────────────────────────────────────────────
gitmap os tweak context-menu classic       # Restore Windows 10 classic right-click menu
gitmap os tweak context-menu modern        # Restore Windows 11 default context menu
gitmap os tweak start-menu classic         # Bring back classic Start Menu layout
gitmap os tweak start-menu default         # Revert Start Menu layout to Windows default
gitmap os tweak power ultimate             # Activate Ultimate Performance power scheme
gitmap os tweak power balanced             # Revert to Balanced power scheme
gitmap os tweak hibernate off              # Disable hibernation and delete hiberfil.sys
gitmap os tweak hibernate on               # Enable hibernation

# ── Linux / Ubuntu Maintenance ──────────────────────────────────────────────
gitmap os clean --system                   # Deep system cleanup (APT/DNF cache, logs, journal vacuum)
gitmap os update                           # Unified package manager and repository update
```

---

## 3. Detailed Component Decomposition & Technical Strategy

### Phase A: Native Go OS Auto-Login (`cli/cmdos/`)

#### 1. Windows Auto-Login Engine (`os_autologin_windows.go`)
- Uses `golang.org/x/sys/windows/registry`:
  - Opens `registry.LOCAL_MACHINE` under `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`.
  - Sets `AutoAdminLogon` = `"1"` (or `"0"` on disable).
  - Sets `DefaultUserName` = `<username>`.
  - Sets `DefaultDomainName` = `<domain>` (defaults to local machine name if empty or `.` is passed).
  - Sets `DefaultPassword` = `<password>`.
  - Sets `ForceAutoLogon` = `"0"` or `"1"`.
- Zero subprocess calls to PowerShell or `autologon.exe`.
- Hermetic validation: checks for elevated administrative privileges (`windows.OpenCurrentProcessToken()`), returning clear `apperror` instructions if not elevated.

#### 2. Ubuntu / Linux Auto-Login Engine (`os_autologin_linux.go`)
- Detects the active Display Manager or desktop environment:
  - **GDM3 (Ubuntu default):** Parses and writes `/etc/gdm3/custom.conf` under `[daemon]` block:
    ```ini
    [daemon]
    AutomaticLoginEnable=true
    AutomaticLogin=<username>
    ```
  - **LightDM (Xubuntu/LXD):** Writes `/etc/lightdm/lightdm.conf.d/50-autologin.conf`:
    ```ini
    [Seat:*]
    autologin-user=<username>
    autologin-user-timeout=0
    ```
  - **Headless / Console (agetty):** Configures `/etc/systemd/system/getty@tty1.service.d/autologin.conf` with `-a <username>`.
- Reversion logic cleanly restores original configuration files from backups (`.gitmap.bak`).

#### 3. 3-Parameter Prompt Handler (`os_autologin_prompt.go`)
- Prompts user interactively if flags are omitted:
  1. `Username` (defaults to current OS user).
  2. `Domain` (defaults to `.` for local machine or current `USERDOMAIN`).
  3. `Password` (secure terminal entry with hidden characters via `golang.org/x/term.ReadPassword`).

---

### Phase B: Windows Tweaks Engine (`os_tweak_windows.go`)

#### 1. Classic Right-Click Context Menu
- Registry path: `HKCU\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`.
- Setting empty string `""` as `(Default)` value forces Windows Explorer to bypass the Windows 11 XAML context menu and load the classic Win32 menu.
- Explorer refresh: sends `WM_SETTINGCHANGE` or gracefully restarts `explorer.exe`.

#### 2. Start Menu Previous Layout
- Registry path: `HKLM\SYSTEM\ControlSet001\Control\FeatureManagement\Overrides\8\3036241548`.
- Sets `EnabledState` DWORD to `1` (Enable) or deletes override entry on revert.

#### 3. Power Schemes & Hibernation
- Ultimate Performance: executes `powercfg /duplicatescheme e9a42b02-d5df-448d-aa00-03f14749eb61` to duplicate scheme if absent, then sets active scheme via `powercfg /setactive <guid>`.
- Hibernation toggle: `powercfg /hibernate off` deletes `C:\hiberfil.sys`, immediately saving 8 GB to 64 GB of disk space on NVMe/SSD drives.

---

### Phase C: Linux System Maintenance & Cleanup (`os_clean_linux.go`)

- **Package Cleanup:** Executes native system package clean commands based on detected package manager (`apt-get clean && apt-get autoremove -y`, `dnf clean all`, `pacman -Sc --noconfirm`).
- **Journal Vacuuming:** Prunes systemd journal logs older than 3 days (`journalctl --vacuum-time=3d`).
- **Log Truncation:** Safely truncates archived `.log` files in `/var/log` to 0 bytes without removing handles.

---

## 4. Execution & Implementation Tasks

- [ ] **Task 1 (Data Types & Interfaces):** Create `cli/cmdos/os_autologin_types.go` defining `AutoLoginConfig`, `TweakTarget`, and `SystemCleanOptions`.
- [ ] **Task 2 (Windows Auto-Login):** Implement `cli/cmdos/os_autologin_windows.go` with native registry read/write operations.
- [ ] **Task 3 (Ubuntu Auto-Login):** Implement `cli/cmdos/os_autologin_linux.go` with GDM3, LightDM, and getty support.
- [ ] **Task 4 (CLI Prompts & Options):** Implement `cli/cmdos/os_autologin_prompt.go` for 3-parameter interactive input with masked password entry.
- [ ] **Task 5 (Windows Shell Tweaks):** Implement `cli/cmdos/os_tweak_windows.go` (Context Menu, Start Menu, Power Schemes, Hibernation).
- [ ] **Task 6 (Linux System Clean):** Implement `cli/cmdos/os_clean_linux.go` for package cache, journal vacuuming, and log cleanup.
- [ ] **Task 7 (Subcommand Dispatch & Help Menus):** Add `gitmap os autologin` and `gitmap os tweak` to `cli/cmdos/os_cmd.go` and author two-column styled help menus in `cli/helptext/os-autologin.md` and `cli/helptext/os-tweak.md`.
- [ ] **Task 8 (Hermetic Unit Testing):** Author unit tests with mock registry and mock filesystem providers in `os_autologin_test.go` and `os_tweak_test.go`.

---

## 5. Verification Plan

1. **Unit Tests:** `go test ./cmdos/... -v` must pass 100% on both Windows and Linux CI runners.
2. **Coding Guidelines:** `python linter-scripts/check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`, `check-enum-and-boolean.py` must pass with 0 errors.
3. **Interactive Verification:**
   - Run `gitmap os autologin status` to check existing state.
   - Run `gitmap os tweak context-menu classic` and verify File Explorer context menu layout.
   - Run `gitmap os tweak power ultimate` and verify active power scheme.
