# 12 — Ubuntu Workstation Governance and Customization

- **Subsystem:** Linux System Configuration & Desktop Governance
- **Status:** Authoritative Reference

## 1. High-DPI Display Scaling
- Configures GNOME fractional scaling (140%) via `gsettings` for 4K display readability.
- Dark theme presets, dock layout configuration, and terminal keyboard shortcuts.

## 2. Virtualization Mount Persistence
- Manages systemd mount units for VMware shared folders with user-specific ownership (`uid=1000,gid=1000`).

## 3. Multi-OS Case Renamer
- Two-phase atomic Git renames resolve case-sensitivity desynchronizations across NTFS and ext4.
