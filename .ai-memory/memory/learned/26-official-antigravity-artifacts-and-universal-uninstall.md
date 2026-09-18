# Official Antigravity Artifacts & Universal Uninstall Engine

- **Date**: 2026-09-13
- **Plan Reference**: `.ai-memory/plans/completed/152-fix-antigravity-and-universal-uninstall.md`

## Overview
This memory record captures the architectural requirements and bug fixes for installing Google Antigravity Desktop IDE and Antigravity CLI (`agy`), along with the universal uninstallation engine across all tracked tools and platforms.

## 1. Official Google Cloud Storage Endpoints
Google Antigravity Desktop IDE artifacts must be downloaded from official Google Cloud Storage endpoints:
- Linux x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-x64/Antigravity.tar.gz`
- Linux ARM: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/linux-arm/Antigravity.tar.gz`
- Windows x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/windows-x64/Antigravity-x64.exe`
- macOS ARM: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-arm/Antigravity.dmg`
- macOS x64: `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.13.0-6362815968182272/darwin-x64/Antigravity.dmg`

## 2. Broken Symlink & Dead Desktop File Traps
- **The Issue**: Go's `os.Stat(path)` follows symlinks. When a symlink is broken/dead, `os.Stat` fails with `os.ErrNotExist`, causing naive cleanup routines (`if _, err := os.Stat(path); err == nil { os.Remove(path) }`) to skip deleting the broken symlink.
- **The Fix**: Call `os.Remove(path)` directly, checking `if err == nil || errors.Is(err, os.ErrNotExist)`.
- **Desktop Entry Gating**: Never write `~/.local/share/applications/antigravity.desktop` or symlink `~/.local/bin/antigravity` until the executable binary has been verified on disk (exists, regular file, size > 0, executable bit `0755`).
- **Cache Invalidation**: On Linux, always run `update-desktop-database` after adding or removing desktop entries to invalidate the GNOME/XDG application launcher cache.

## 3. Universal Uninstall Coverage
- **Every tool MUST be uninstallable**: Custom standalone tools (`agy`, `antigravity`, `ag-manager`, `scripts-fixer`, `coding-guidelines`, `macro-ahk`), context menus (`ctx`, `vscode-ctx`, `pwsh-ctx`, `ag-ctx`), cloned repos (`scripts`), and self-uninstall (`gitmap uninstall gitmap`).
- **Flag Reordering**: Standard Go `flag.FlagSet.Parse(args)` stops at the first positional argument. Trailing flags (e.g. `gitmap uninstall agy --force`) will be ignored unless args are normalized with `reorderFlagsBeforeArgs(args)`.
- **Dual DB Deletion**: Uninstallation must purge records from both `gitmap.db` and `installation.db`.

## 4. AppError Stack Traces
- All uninstallation and installation errors must be wrapped in `*apperror.AppError` and output full stack traces via `cliexit.HandleError`. Never swallow or truncate error traces.
