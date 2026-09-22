## Quick Install v6.299.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.299.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.299.0/install.sh | bash
```

## Changelog v6.299.0

- Native OS Auto-Login Subsystem: Added `gitmap os autologin` with interactive 3-parameter credentials (username, domain, masked password) and flags (`-u`, `-d`, `-p`), natively writing Windows Winlogon registry keys and configuring Ubuntu display managers (GDM3, LightDM) with zero PowerShell or external executable dependencies.
- Windows Desktop Tweaks: Added `gitmap os tweak context-menu classic|modern` (CLSID InprocServer32 restore for Windows 10 right-click menu) and `start-menu classic|default` layout override.
- Ultimate Performance & Hibernation: Added `gitmap os tweak power ultimate|balanced` to duplicate and activate Windows Ultimate Performance scheme and `gitmap os tweak hibernate off|on` to delete `C:\hiberfil.sys` and reclaim RAM-sized disk space.
- Linux System Maintenance: Added `gitmap os clean sys` for multi-distro package manager cache purging (APT, DNF, Pacman) and systemd journal vacuuming (`journalctl --vacuum-time=3d`).
- Two-Column Styled Help Screens: Authored comprehensive help documentation and terminal menus for `gitmap help os-autologin` and `gitmap help os-tweak`.
- CI/CD Self-Healing: Restored local runner helpers in `03-ai-scripts/06-cicd-local-runner.py`, fixed 18/18 Python CI test cases, removed unused symbols, and verified cross-platform compilation across Linux, Darwin, and Windows.
