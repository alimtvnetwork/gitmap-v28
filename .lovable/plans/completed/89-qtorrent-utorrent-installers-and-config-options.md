# Plan 89: qBittorrent & uTorrent Installers and JSON Config Export/Import Engine

**Title:** qBittorrent & uTorrent Installers and JSON Config Export/Import Engine  
**Status:** Completed  
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.1.0)  
**Budget (N):** 50 steps  
**Target Directory:** `.lovable/plans/completed/`  

---

## 1. Context & Objectives

1. **qBittorrent & uTorrent Cross-Platform Installers**:
   - Installable on Windows (Chocolatey: `qbittorrent`, `utorrent`; Winget: `qBittorrent.qBittorrent`, `BitTorrent.uTorrent`).
   - Installable on Ubuntu/Debian Linux (APT: `qbittorrent`, `utorrent`).
   - Installable on macOS (Homebrew: `qbittorrent`, `utorrent`).
   - Canonical constants: `ToolQBittorrent = "qbittorrent"`, `ToolUTorrent = "utorrent"`.
   - Category: `ToolCategoryUtilities` ("Terminal & Utilities").
   - Aliases: `qtorrent`, `qbittorrent`, `qbit`, `utorrent`, `u-torrent`, `uttorrent`.
   - Probing: `qbittorrent --version`, `qbittorrent-nox --version`, `utorrent --version`, `utserver --version`.
   - Display in `gitmap install --list` / `gitmap in` table.

2. **Portable JSON Config Export/Import Engine**:
   - Commands:
     - `gitmap export-config <tool|all> [path]` (alias: `gitmap config-export`)
     - `gitmap import-config <tool|all> [path]` (aliases: `gitmap config-import`, `gitmap improt-config`)
   - Default file output naming convention:
     - `vscode` -> `vscode.json`
     - `qtorrent` -> `qtorrent.json`
     - `uttorrent` / `utorrent` -> `uttorrent.json` / `utorrent.json`
   - Target path resolution:
     - If path is omitted, export/import `<tool>.json` in current directory.
     - If path is a folder (or ends with slash), export/import `<folder>/<tool>.json`.
     - If `all` is specified, batch export or batch import all configurations to/from folder.
   - Cross-Platform Configuration directories & files:
     - **VS Code**:
       - Windows: `%APPDATA%\Code\User\settings.json`, `keybindings.json`
       - Linux: `~/.config/Code/User/settings.json`, `keybindings.json`
       - macOS: `~/Library/Application Support/Code/User/settings.json`, `keybindings.json`
       - Extensions: `code --list-extensions` saved in bundle
     - **qBittorrent**:
       - Windows: `%APPDATA%\qBittorrent\qBittorrent.ini`
       - Linux: `~/.config/qBittorrent/qBittorrent.conf`
       - macOS: `~/Library/Application Support/qBittorrent/qBittorrent.conf`
     - **uTorrent**:
       - Windows: `%APPDATA%\uTorrent\settings.dat`
       - Linux: `~/.config/utorrent/settings.dat`
       - macOS: `~/Library/Application Support/uTorrent/settings.dat`
   - Universal JSON Schema (v1):
     - Text files encoded as UTF-8 string.
     - Binary files (e.g. `settings.dat`) encoded as Base64.
     - Starter template fallback when local tool is not yet configured.

3. **Help Text, UI Help & Verification**:
   - `gitmap/helptext/export-config.md` (<= 120 lines, fenced blocks, golden tested).
   - `gitmap/helptext/import-config.md` (<= 120 lines, fenced blocks, golden tested).
   - `gitmap/helptext/install.md` updated with supported torrent tools table.
   - `src/data/commands.ts` updated with UI command definitions and examples.
   - AST parity verification: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1`.
   - Comprehensive unit tests in `gitmap/cmd/config_export_import_test.go` and `gitmap/cmd/install_packages_test.go`.

---

## 2. Subtasks Ledger

- **Subtask 89.01 (DONE):** Define tool constants, category grouping, descriptions, and package manager mappings across Chocolatey, Winget, Apt, and Brew in `gitmap/constants/constants_install.go`, `gitmap/cmd/install_packages.go`, `gitmap/cmd/install_packages_extra.go`, and `gitmap/cmd/installprobe.go`.
- **Subtask 89.02 (DONE):** Implement portable JSON configuration export/import engine in `gitmap/cmd/config_model.go`, `gitmap/cmd/config_tool_paths.go`, `gitmap/cmd/config_export.go`, `gitmap/cmd/config_export_writers.go`, `gitmap/cmd/config_import.go`, `gitmap/cmd/config_import_readers.go`, and wire into `gitmap/cmd/roottooling.go`.
- **Subtask 89.03 (DONE):** Register CLI constants in `gitmap/constants/constants_cli.go`, regenerate code completion in `gitmap/completion/allcommands_generated.go`, verify AST parity in `gitmap/constants/cmd_constants_test.go`.
- **Subtask 89.04 (DONE):** Author help text markdown files in `gitmap/helptext/export-config.md` & `gitmap/helptext/import-config.md`, update `gitmap/helptext/install.md`, and update web UI command registry in `src/data/commands.ts`.
- **Subtask 89.05 (DONE):** Implement unit test suites in `gitmap/cmd/config_export_import_test.go` and update `gitmap/cmd/install_packages_test.go`. Verify all tests pass, run custom quality linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-error-management.py`), and verify clean Git tree.
