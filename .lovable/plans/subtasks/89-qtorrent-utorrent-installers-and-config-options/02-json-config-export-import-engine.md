# Subtask 89.02: JSON Config Export & Import Engine

**Plan:** 89-qtorrent-utorrent-installers-and-config-options.md  
**Status:** Completed  

## 1. Description & Requirements
- Implement portable configuration bundle model (`ConfigBundle`, `ConfigFilePayload`) in `gitmap/cmd/config_model.go`.
- Implement tool canonical naming and default filename resolver:
  - `vscode` -> `vscode.json`
  - `qtorrent` -> `qtorrent.json`
  - `uttorrent` -> `uttorrent.json` (also `utorrent` -> `utorrent.json`)
  - Target directory path resolution: if directory specified, join with default filename.
- Implement cross-platform directory resolution in `gitmap/cmd/config_tool_paths.go`:
  - `resolveOSConfigPath(winSub, macSub, linuxSub string)`
  - qBittorrent: `%APPDATA%\qBittorrent` (Windows), `~/.config/qBittorrent` (Linux), `~/Library/Application Support/qBittorrent` (macOS).
  - uTorrent: `%APPDATA%\uTorrent` (Windows), `~/.config/utorrent` (Linux), `~/Library/Application Support/uTorrent` (macOS).
  - VS Code: `%APPDATA%\Code\User` (Windows), `~/.config/Code/User` (Linux), `~/Library/Application Support/Code/User` (macOS).
- Implement export engine in `gitmap/cmd/config_export.go` and `gitmap/cmd/config_export_writers.go`:
  - Single export: `gitmap export-config <tool> [path]`
  - Batch export: `gitmap export-config all [path]`
  - UTF-8 text encoding for ini/conf/json, Base64 encoding for binary (`settings.dat`).
  - Starter fallback templates if local files are missing.
- Implement import engine in `gitmap/cmd/config_import.go` and `gitmap/cmd/config_import_readers.go`:
  - Single import: `gitmap import-config <tool> [path]` (also `improt-config` and `config-import`).
  - Batch import: `gitmap import-config all [path]` discovering all `*.json` bundles.
  - OS filename translation (e.g. `qBittorrent.ini` on Windows <-> `qBittorrent.conf` on Linux).
  - Base64 decoding for binary restoration.
- Wire CLI dispatch in `gitmap/cmd/roottooling.go`:
  - `export-config`, `config-export`
  - `import-config`, `config-import`, `improt-config`
