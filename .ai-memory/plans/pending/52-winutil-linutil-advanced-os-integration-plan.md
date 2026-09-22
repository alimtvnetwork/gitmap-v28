# Plan 52: Advanced WinUtil & LinUtil Native Go Integration (DM Settings, OS Tweaks, DNS Switcher & Bubbletea TUI)

## 1. Context & Motivation
Following the successful initial integration of native Go OS auto-login and core Windows desktop tweaks (Plan 51, released in `v6.299.0` and `v6.299.1`), a deep exploration of Chris Titus Tech's repositories in `D:\work\chris` (`winutil` and `linutil`) reveals several high-value capabilities that can be brought into GitMap in pure native Go:
1. **Display Manager (`dm.settings`) & Session Management (Linux / Ubuntu):**
   - LinUtil inspects and configures Display Managers (GDM3, LightDM, SDDM) and Desktop Environments (GNOME, KDE Plasma, XFCE).
   - GitMap currently handles auto-login and display blanking timeouts, but lacks a dedicated `gitmap os dm` subsystem to inspect active display managers, toggle Wayland vs X11 session modes (critical for headless VMs, remote desktop, and GPU stability), switch active DMs, and restart greeters.
2. **Windows Privacy, Performance & Debloating Tweaks (`winutil`):**
   - Telemetry opt-out (`DiagTrack` service disable, `AllowTelemetry=0`, `POWERSHELL_TELEMETRY_OPTOUT=1`).
   - Activity history publication disable (`PublishUserActivities=0`, `UploadUserActivities=0`).
   - Bing and cloud search in Start Menu disable (`BingSearchEnabled=0`).
   - Dark / Light mode system theme toggle (`AppsUseLightTheme`, `SystemUsesLightTheme`).
3. **High-Performance Native DNS Switcher (`dns.json`):**
   - Instant adapter switching between pre-configured, low-latency secure DNS providers (Google, Cloudflare, Quad9, AdGuard, DHCP default).
   - Native DNS ping / latency benchmarking.
4. **Universal Multi-Distro System Updater (`linutil` `system-update.sh`):**
   - Universal package and system upgrade command (`gitmap os update` / `gitmap os upgrade`): APT, DNF, Pacman, Flatpak, Snap, Homebrew, Winget.
5. **Interactive Cross-Platform Bubbletea TUI (inspired by `linutil`'s ratatui & `winutil`'s WPF):**
   - While `winutil` relies on Windows-only PowerShell WPF and `linutil` relies on Rust Ratatui, GitMap already contains Charm's `bubbletea`, `lipgloss`, and `bubbles` in `cli/go.mod`.
   - An interactive `gitmap os tui` keyboard-navigable dashboard for browsing, toggling, and applying OS tweaks, auto-login credentials, display settings, and cleanups in any terminal on Windows, Linux, and macOS.

---

## 2. Architecture & Subsystem Design

```
                     ┌──────────────────────────────────────────────┐
                     │          gitmap os (CLI Dispatcher)          │
                     └──────────────────────┬───────────────────────┘
                                            │
        ┌───────────────────┬───────────────┴───────────────┬───────────────────┐
        │                   │                               │                   │
        ▼                   ▼                               ▼                   ▼
┌──────────────┐    ┌──────────────┐                ┌──────────────┐    ┌──────────────┐
│ gitmap os dm │    │gitmap os dns │                │gitmap os tweak│   │gitmap os tui │
│(Linux Display│    │ (Fast Secure │                │ (Windows OS  │    │(Interactive  │
│   Manager)   │    │  Resolvers)  │                │   Tweaks)    │    │  Bubbletea)  │
└───────┬──────┘    └───────┬──────┘                └───────┬──────┘    └───────┬──────┘
        │                   │                               │                   │
        ▼                   ▼                               ▼                   ▼
- Active DM detect  - Cloudflare (1.1.1.1)          - Telemetry opt-out - Categorized
- Wayland on/off    - Google (8.8.8.8)              - Activity history    tabs & items
- DM Switch         - Quad9 (9.9.9.9)               - Start search      - Multi-select
- Greeter restart   - AdGuard (94.140.14.14)        - Dark/Light theme    checkboxes
- Autologin sync    - DHCP auto-revert              - Right-click menu  - In-place run
```

---

## 3. Subsystem Specifications

### Part A: Linux Display Manager Subsystem (`gitmap os dm`)
- **CLI Commands:**
  - `gitmap os dm status`: Inspects active Display Manager (GDM3, LightDM, SDDM), session type (Wayland / X11), desktop environment, and greeter configuration.
  - `gitmap os dm wayland <enable|disable>`:
    - On GDM3: Safely edits `/etc/gdm3/custom.conf` (`WaylandEnable=false|true`).
    - On LightDM/SDDM: Configures session type in configuration manifests.
  - `gitmap os dm restart`: Restarts display manager service via `systemctl restart display-manager`.
  - `gitmap os dm switch <gdm3|lightdm|sddm>`: Configures default systemd display-manager target.
- **Implementation Strategy:**
  - Pure Go file read/write with atomic tempfile replacement and permission preservation (`0644`).
  - Cross-platform build tagging: `cli/cmdos/os_dm_linux.go` and `cli/cmdos/os_dm_other.go` (stubs on Windows/macOS).

### Part B: Windows Privacy & Performance Tweaks (`gitmap os tweak`)
- **CLI Commands:**
  - `gitmap os tweak telemetry <off|on>`:
    - Writes `HKLM\SOFTWARE\Policies\Microsoft\Windows\DataCollection\AllowTelemetry = 0|1`.
    - Stops and disables or starts `DiagTrack` (Connected User Experiences and Telemetry).
    - Sets machine environment variable `POWERSHELL_TELEMETRY_OPTOUT = 1|0`.
  - `gitmap os tweak activity <off|on>`:
    - Writes `HKLM\SOFTWARE\Policies\Microsoft\Windows\System\PublishUserActivities = 0|1`.
    - Writes `HKLM\SOFTWARE\Policies\Microsoft\Windows\System\UploadUserActivities = 0|1`.
  - `gitmap os tweak search <clean|default>`:
    - Writes `HKCU\Software\Microsoft\Windows\CurrentVersion\Search\BingSearchEnabled = 0|1`.
  - `gitmap os tweak theme <dark|light>`:
    - Windows: `AppsUseLightTheme` and `SystemUsesLightTheme` in `HKCU\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`.
    - Linux: `gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark'|'default'`.

### Part C: Fast Native DNS Switcher (`gitmap os dns`)
- **CLI Commands:**
  - `gitmap os dns status`: Prints active DNS servers per network interface.
  - `gitmap os dns set <cloudflare|google|quad9|adguard>`: Sets primary and secondary IPv4 & IPv6 addresses on active adapters.
  - `gitmap os dns dhcp`: Reverts adapter DNS configuration to automatic DHCP assignment.
  - `gitmap os dns benchmark`: Measures ping and DNS lookup latency to all supported providers.

### Part D: Universal System Updater (`gitmap os update` / `gitmap os upgrade`)
- **CLI Commands:**
  - `gitmap os update`: Checks and updates package manager repositories without applying upgrades.
  - `gitmap os upgrade`: Executes full package and toolchain upgrades across all detected package managers on the host:
    - Linux: APT (`apt-get dist-upgrade -y`), DNF (`dnf upgrade -y`), Pacman (`pacman -Syu --noconfirm`), Flatpak (`flatpak update -y`), Snap (`snap refresh`).
    - Windows: Winget (`winget upgrade --all --include-unknown`).
    - macOS: Homebrew (`brew update && brew upgrade`).

### Part E: Cross-Platform Interactive Bubbletea TUI (`gitmap os tui`)
- Interactive dual-pane terminal interface:
  - Tab 1: **Tweaks** (Context Menu, Start Menu, Telemetry, Activity Feed, Hibernation, Power).
  - Tab 2: **Auto-Login** (Interactive status, credential inputs, enable/disable).
  - Tab 3: **Display & DM** (Display timeouts, Wayland on/off, active DM).
  - Tab 4: **Network & DNS** (Current IPs, fast DNS provider switcher).
  - Tab 5: **Maintenance** (Package cleanup, journal vacuum, system update).
- Real-time status badges (`[Active]`, `[Disabled]`, `[Default]`).
- Keyboard controls: Arrow keys, Tab / Shift-Tab, Space to toggle, Enter to execute, `q` to quit.

---

## 4. Execution Subtasks

- [ ] **Subtask 52.1**: Implement Display Manager (`gitmap os dm`) subsystem in native Go (`os_dm_types.go`, `os_dm_linux.go`, `os_dm_other.go`, `os_dm_cmd.go`).
- [ ] **Subtask 52.2**: Extend Windows desktop tweaks with telemetry, activity feed, and dark/light theme switching (`os_tweak_windows_privacy.go`, `os_tweak_theme.go`).
- [ ] **Subtask 52.3**: Implement Native DNS Switcher (`os_dns_types.go`, `os_dns_windows.go`, `os_dns_linux.go`, `os_dns_cmd.go`).
- [ ] **Subtask 52.4**: Implement Universal Multi-Distro Updater (`os_update_types.go`, `os_update_engine.go`, `os_update_cmd.go`).
- [ ] **Subtask 52.5**: Author Two-Column Styled Help Screens (`cli/helptext/os-dm.md`, `cli/helptext/os-dns.md`, `cli/helptext/os-update.md`).
- [ ] **Subtask 52.6**: Add Unit Tests for new OS subsystems in `cli/cmdos/`.
- [ ] **Subtask 52.7**: Verify all quality gates, linters, and release ceremony via `06-ci-cd-fix-with-release.md`.
