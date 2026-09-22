# WinUtil & LinUtil Native Integration: OS Auto-Login (Windows & Ubuntu), System Tweaks & Maintenance

- **Number:** 32
- **Date:** 2026-09-22
- **Topic:** Native Go integration of Windows/Ubuntu auto-login, desktop shell tweaks, power schemes, and system cleanups.
- **Reference:** `D:\work\chris\winutil` and `D:\work\chris\linutil`

## 1. Context & Architecture

Analyzed Chris Titus's repositories (`winutil` and `linutil`) to extract high-utility workstation maintenance features and port them directly into GitMap:
- Instead of downloading Sysinternals `Autologon.exe` or executing slow PowerShell scripts (`Invoke-WPFPanelAutologin.ps1`), GitMap implements Windows Auto-Login natively via `golang.org/x/sys/windows/registry` against `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`.
- For Ubuntu/Linux, GitMap detects Display Managers (`gdm3` via `/etc/gdm3/custom.conf`, `lightdm` via `/etc/lightdm/lightdm.conf.d/50-autologin.conf`) and configures automatic login atomically with file backup and restore.
- Zero external runtime dependencies, zero PowerShell overhead.

## 2. Windows Tweaks Engine

Implemented under `cli/cmdos/`:
- **Classic Context Menu:** Toggles `HKCU\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32` default value `""` to bypass Windows 11 XAML menu and restore classic Windows 10 menu without third-party software.
- **Start Menu Previous Layout:** Sets `EnabledState` DWORD in `HKLM\SYSTEM\ControlSet001\Control\FeatureManagement\Overrides\8\3036241548`.
- **Ultimate Performance Power Scheme:** Activates hidden Windows scheme `e9a42b02-d5df-448d-aa00-03f14749eb61` and reverts to Balanced `381b4222-f694-41f0-9685-ff5bb260df2e`.
- **Hibernation Toggle:** Controls `powercfg /hibernate on|off` to purge `C:\hiberfil.sys`, freeing RAM-sized disk space on SSDs.

## 3. Linux Maintenance Parity

- **System Clean:** Package manager cache purging (`apt-get clean && autoremove -y`, `dnf clean all`, `pacman -Sc`) and systemd journal vacuuming (`journalctl --vacuum-time=3d`) via `gitmap os clean sys`.

## 4. Coding Guideline Strict Compliance

- All 22 new Go files are strictly <= 100 lines each.
- Hermetic test injection via `SetAutoLoginEngine(...)` and `SetTweakEngine(...)` ensuring unit tests never mutate host OS or trigger system modifications.
