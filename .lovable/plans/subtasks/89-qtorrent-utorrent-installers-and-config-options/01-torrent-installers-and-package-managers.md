# Subtask 89.01: Torrent Installers & Package Managers

**Plan:** 89-qtorrent-utorrent-installers-and-config-options.md  
**Status:** Completed  

## 1. Description & Requirements
- Register `ToolQBittorrent = "qbittorrent"` and `ToolUTorrent = "utorrent"` in `gitmap/constants/constants_install.go`.
- Assign to `ToolCategoryUtilities` ("Terminal & Utilities") in `InstallToolCategories`.
- Add tool descriptions:
  - `ToolQBittorrent`: "qBittorrent free and open-source BitTorrent client"
  - `ToolUTorrent`: "uTorrent lightweight BitTorrent client"
- Define package manager mapping constants:
  - Chocolatey: `ChocoPkgQBittorrent = "qbittorrent"`, `ChocoPkgUTorrent = "utorrent"`
  - Winget: `WingetPkgQBittorrent = "qBittorrent.qBittorrent"`, `WingetPkgUTorrent = "BitTorrent.uTorrent"`
  - Apt: `AptPkgQBittorrent = "qbittorrent"`, `AptPkgUTorrent = "utorrent"`
  - Brew: `BrewPkgQBittorrent = "qbittorrent"`, `BrewPkgUTorrent = "utorrent"`
- Register aliases in `gitmap/cmd/install_packages.go`: `qtorrent`, `qbittorrent`, `qbit`, `utorrent`, `u-torrent`, `uttorrent`.
- Register version probing in `gitmap/cmd/installprobe.go`:
  - `ToolQBittorrent`: `bins: ["qbittorrent", "qbittorrent-nox"]`, `args: ["--version"]`
  - `ToolUTorrent`: `bins: ["utorrent", "uTorrent", "utserver"]`, `args: ["--version"]`
