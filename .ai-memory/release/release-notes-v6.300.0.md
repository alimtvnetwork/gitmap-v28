## Quick Install v6.300.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.300.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.300.0/install.sh | bash
```

---

## What's Changed in v6.300.0

- **Native Linux Display Manager (`gitmap os dm`):** Linux display manager detection and configuration parity with LinUtil's `dm.settings`, inspecting active DM (GDM3, LightDM, SDDM), toggling Wayland mode (`gitmap os dm wayland on|off`), and restarting display-manager services without shell scripts.
- **Windows Privacy & Start Menu Tweaks:** Added `gitmap os tweak telemetry off|on` (disables DiagTrack service, sets telemetry opt-out), `activity off|on` (disables Windows Activity Feed collection & cloud upload), and `search clean|default` (removes Bing and web suggestions from Start Menu).
- **Native DNS Switcher & Benchmark (`gitmap os dns`):** Added ultra-fast DNS switcher supporting Cloudflare (1.1.1.1), Google (8.8.8.8), Quad9 (9.9.9.9), and AdGuard (94.140.14.14), automated DHCP restoration (`gitmap os dns dhcp`), and integrated UDP DNS benchmark (`gitmap os dns bench`).
- **Universal System Updater (`gitmap os update / upgrade`):** Cross-platform multi-package manager aggregator supporting winget, apt, dnf, pacman, and brew with dry-run support (`--dry-run`).
- **Desktop Theme Switcher (`gitmap os theme dark|light`):** Native registry theme toggle on Windows (AppsUseLightTheme, SystemUsesLightTheme) and gsettings on Linux/GNOME.
- **Two-Column Styled Help Screens:** Authored comprehensive help documentation and terminal menus for `gitmap help os-dm`, `gitmap help os-dns`, `gitmap help os-theme`, and `gitmap help os-update`.
- **Code Hygiene & Modularity:** Strictly preserved $\le 100$ lines per file, flattened all nested ifs to depth 1, and achieved 100% compliance across all CI/CD policy linters and cross-platform vet.
