---
name: uninstall-quick-scripts
description: Root-level uninstall-quick.ps1 / uninstall-quick.sh one-liner uninstallers that wrap `gitmap self-uninstall` with a manual-sweep fallback. Mirror install-quick.* layout.
type: feature
---
# Uninstall quick scripts

Two new root-level scripts mirror `install-quick.ps1` / `install-quick.sh`:

- `uninstall-quick.ps1` (Windows)
- `uninstall-quick.sh` (Linux / macOS)

Both are designed to be piped from a one-liner:

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/uninstall-quick.ps1 | iex
```

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/uninstall-quick.sh | bash
```

## Strategy (both scripts)

1. **Try canonical `gitmap self-uninstall -y` first.** This is the best path — the binary itself knows about marker-block PATH cleanup, scheduled-task removal, etc. (See `gitmap/cmd/selfuninstall.go`.)
2. **Manual sweep fallback** (when `gitmap` is no longer on PATH):
   - Auto-detect deploy root by walking the active binary's grandparent, then probing common defaults (`E:\bin-run`, `D:\gitmap`, `$LOCALAPPDATA\gitmap` on Windows; `~/.local/bin`, `~/bin`, `/opt/gitmap` on Unix).
   - Delete BOTH `<root>/gitmap-cli/` (current v3.6+ layout) AND `<root>/gitmap/` (legacy pre-rename layout) AND any flat `<root>/gitmap.exe` for very old installs.
   - Strip the deploy root from User PATH (Windows) or shell rc files (Unix: `~/.bashrc`, `~/.zshrc`, `~/.profile`, `~/.bash_profile` — backed up to `*.gitmap-uninstall.bak`).
3. **Always prompt before deleting user data** (`%APPDATA%\gitmap` on Windows, `${XDG_CONFIG_HOME:-$HOME/.config}/gitmap` on Unix). `-KeepData`/`--keep-data` skips the prompt and keeps; `-Yes`/`-y --yes` skips the prompt and deletes.

## Flags

| Windows | Unix | Effect |
|---|---|---|
| `-Yes` | `-y`, `--yes` | Skip the "delete user data?" prompt, assume yes |
| `-KeepData` | `--keep-data` | Always keep user data folder (overrides `-Yes`) |
| `-InstallDir` | `--dir` | Override auto-detected deploy root |

## Why both layout names?

The deploy folder rename `gitmap/` → `gitmap-cli/` (v3.6+) means the uninstaller MUST handle both layouts to clean up legacy installs and current installs alike. The sweep deletes both unconditionally — neither existing is a no-op.

## Documentation

- README "Quick Start" section now has an "Uninstall — Quick (one-liner)" subsection right under the install commands.
- The React UI on the docs site (`src/pages/Index.tsx`) was simplified at the same time: the Windows install command went from the long `Set-ExecutionPolicy Bypass ... iex ((New-Object ...))` form to the short `irm install-quick.ps1 | iex` form. This matches the README and the Linux side, and looks cleaner in the terminal screenshot on the homepage.
