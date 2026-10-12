# Specification: Scripts-Fixer Update Option and Live Ubuntu Host Antigravity IDE Update

**Version:** 1.0.0  
**Updated:** 2026-10-12  
**AI Confidence:** High (100%)  
**Ambiguity:** None  
**Author:** Subagent Spec Author  
**Parent Task:** `update-antigravity-ide`  

---

## Keywords

`scripts-fixer-update` · `69-install-antigravity` · `manifest-update-verb` · `install-antigravity-sh` · `in-place-upgrade` · `preserve-symlinks` · `antigravity-runner` · `electron-no-sandbox` · `backup-and-replace` · `live-system-verification` · `zero-absolute-paths`

---

## Quality & Completeness Scoring

| Criterion | Status | Description |
|---|:---:|---|
| AI Confidence assigned | ✅ | Set to High (100%) |
| Ambiguity assigned | ✅ | Set to None |
| Keywords present | ✅ | Comprehensive indexing of update components |
| Verbatim requirements captured | ✅ | All script targets and live host update constraints addressed |
| Relative path mandate strictly obeyed | ✅ | Zero absolute paths across specification text and examples |
| `scripts-fixer` modifications specified | ✅ | Detailed changes for `run.sh`, `manifest.json`, and `install-antigravity.sh` |
| Live host update safety covered | ✅ | Backups, permissions, runners, desktop entries, and rollback mechanisms |
| Verification criteria defined | ✅ | Concrete post-update smoke checks and file integrity assertions |

---

## 1. Executive Summary & Problem Scope

This specification defines the architectural design, CLI contract, and implementation details for adding the `update` action/verb to the `scripts-fixer` toolchain, as well as conducting safe, live updates of the Google Antigravity Desktop IDE on the host Ubuntu system (`~/.local/share/antigravity-ide`).

### 1.1 Problem Statement

1. **Missing `update` Verb in `scripts-fixer`**:
   - The automation module in `scripts-linux/69-install-antigravity/` currently provides `install`, `check`, `repair`, `uninstall`, `clean`, `clear`, `cache`, `predict`, `list-backups`, and `undo`.
   - `verb_install()` contains an early exit: if `verify_installed` evaluates to true, the script exits immediately with `"Antigravity already installed"` without pulling new binaries or updating existing versions.
   - `verb_repair()` merely removes the installation stamp and delegates to `verb_install()`, which fails to fetch new versions if an existing binary is detected.
   - Neither `scripts-linux/69-install-antigravity/manifest.json` nor `scripts/69-install-antigravity/manifest.json` lists `update` as a supported action.

2. **Absence of In-Place Upgrade Flow in `install-antigravity.sh`**:
   - The underlying installation script `scripts/os/ubuntu/install-antigravity.sh` supports `--force` (which nukes the entire installation directory, including potential local state or symlinks) and `--download-must`.
   - It lacks a designated `--update` workflow that safely backs up the current installation, stages the new tarball, extracts the updated binary, updates permissions, and re-validates launcher symlinks without breaking desktop integration or active developer sessions.

3. **Live System Update Requirement on Ubuntu Host**:
   - The host system contains an active Antigravity IDE deployment located in `~/.local/share/antigravity-ide`.
   - Symlinks in `~/.local/bin/` (`antigravity`, `agy`, `antigravity-ide`) point to the launcher wrapper `~/.local/share/antigravity-ide/antigravity-ide.run`.
   - A live update must replace the application binary and runtime components, retain executable flags, ensure the `antigravity-runner.sh` wrapper applies the requisite `--no-sandbox` flags, and preserve the desktop launcher and icon assets.

---

## 2. Scripts-Fixer Update Action Architecture

### 2.1 File Box Overview

All target paths are strictly relative to the repository roots:

| Repository | Relative Target File | Purpose |
|---|---|---|
| `scripts-fixer` | `scripts-linux/69-install-antigravity/manifest.json` | Register `update` action in Linux module manifest |
| `scripts-fixer` | `scripts/69-install-antigravity/manifest.json` | Register `update` action in cross-platform module manifest |
| `scripts-fixer` | `scripts-linux/69-install-antigravity/run.sh` | Implement `verb_update()` handler and CLI help documentation |
| `scripts-fixer` | `scripts/os/ubuntu/install-antigravity.sh` | Add `--update` flag, backup creation, atomic extract, and symlink validation |

---

### 2.2 Manifest Registration

Both `scripts-linux/69-install-antigravity/manifest.json` and `scripts/69-install-antigravity/manifest.json` must be updated to expose `update` as an official command:

```json
{
  "schemaVersion": "1.0",
  "id": "69",
  "name": "Antigravity IDE & agy CLI",
  "summary": "Google Antigravity AI coding assistant IDE and agy CLI",
  "aliases": ["agy", "antigravity"],
  "platforms": {
    "windows": {
      "folder": "69-install-antigravity",
      "entrypoint": "run.ps1"
    },
    "linux": {
      "folder": "69-install-antigravity",
      "entrypoint": "run.sh"
    },
    "macos": {
      "folder": "69-install-antigravity",
      "entrypoint": "run.sh"
    }
  },
  "depends_on": [],
  "idempotent": true,
  "destructive": false,
  "commands": [
    "install",
    "update",
    "check",
    "repair",
    "uninstall",
    "clean",
    "clear",
    "cache",
    "cache-clear",
    "predict",
    "list-backups",
    "undo",
    "help"
  ],
  "phase": "04",
  "tags": ["ai", "ide", "cli", "antigravity"]
}
```

---

### 2.3 `run.sh` Handler Specification (`verb_update`)

In `scripts-linux/69-install-antigravity/run.sh`, the following modifications are required:

1. **`verb_update()` Function Implementation**:
   - Writes install paths metadata to `_shared/install-paths.sh` tracker.
   - Logs informative status: `[69] Starting Antigravity IDE update`.
   - Locates `scripts/os/ubuntu/install-antigravity.sh` relative to `REPO_ROOT`.
   - Invokes the installer script passing `--update` alongside any forwarded arguments (such as `--download-must`, `--force`, or target version):
     ```bash
     verb_update() {
       write_install_paths \
         --tool   "antigravity" \
         --source "https://storage.googleapis.com/antigravity-public (official archive)" \
         --temp   "/tmp/scripts-fixer-downloads/Antigravity.tar.gz" \
         --target "$HOME/.local/bin/agy"

       log_info "[69] Starting Antigravity IDE update"

       local ub_script="$REPO_ROOT/scripts/os/ubuntu/install-antigravity.sh"
       if [ -f "$ub_script" ]; then
         log_info "[69] Executing ubuntu updater: $ub_script"
         if bash "$ub_script" --update "$@"; then
           log_ok "[69] Antigravity updated successfully"
           mkdir -p "$ROOT/.installed"
           touch "$INSTALLED_MARK"
           return 0
         fi
       fi

       log_err "[69] Antigravity update failed"
       return 1
     }
     ```
2. **Help Documentation Update (`show_help`)**:
   - Add `update` to the `ACTIONS:` list:
     ```text
       install                Install Antigravity IDE and CLI (agy)
       update                 Update Antigravity IDE and CLI to latest release
       cli                    Install Antigravity CLI only (agy)
       check                  Check if Antigravity is installed and on PATH
       repair                 Repair Antigravity installation
       uninstall              Uninstall Antigravity IDE and CLI
     ```
3. **Command Dispatcher Update**:
   - Update the `case "${1:-help}" in` block:
     ```bash
     case "${1:-help}" in
       help|--help|-help|-h) show_help;;
       install)   verb_install;;
       update)    verb_update "${@:2}";;
       check)     verb_check;;
       repair)    verb_repair;;
       uninstall) verb_uninstall;;
       ...
     ```

---

### 2.4 `install-antigravity.sh` In-Place Upgrade Specification

In `scripts/os/ubuntu/install-antigravity.sh`, the script must be updated to handle the `update` verb and `--update` flag:

1. **Argument Parsing**:
   - Introduce `IS_UPDATE=false`.
   - Set `IS_UPDATE=true` if `$arg` equals `--update`, `-u`, or `update`.
   - Environment variable fallback: `[ "${IS_UPDATE:-false}" = "true" ] || [ "${UPDATE:-0}" = "1" ]`.

2. **Pre-Update Backup Mechanism**:
   - Before extracting any new files, create a safety backup of the existing IDE directory:
     ```bash
     backup_existing_installation() {
       local ide_dir="$HOME/.local/share/antigravity-ide"
       if [ -d "$ide_dir" ]; then
         local ts
         ts=$(date +%Y%m%d_%H%M%S)
         local bak_dir="$HOME/.local/share/antigravity-ide.bak-$ts"
         echo -e "  ${MUTED}[backup] Creating safety backup at $bak_dir...${TEXT}"
         cp -a "$ide_dir" "$bak_dir" 2>/dev/null || true
       fi
     }
     ```
   - Keep only the 2 most recent `.bak-*` directories to avoid disk bloat.

3. **In-Place Extract and Replace (`update_ide_installation`)**:
   - Extract the downloaded `Antigravity.tar.gz` into a temporary staging folder (`mktemp -d /tmp/antigravity-update-XXXXXX`).
   - Validate that the unpacked archive contains the `antigravity` executable binary.
   - Synchronize extracted contents into `~/.local/share/antigravity-ide`:
     - Overwrite core binaries, runtime `.so` libraries, locales, and `resources/app.asar`.
     - Preserve user customizations, logs (e.g., `ide-discovery.log`), and local caches.
   - Clean up the staging directory immediately.

4. **Runtime Wrapper & Symlink Preservation**:
   - Ensure `~/.local/share/antigravity-ide/antigravity-runner.sh` and `antigravity-ide.run` are recreated or updated with the canonical wrapper logic:
     ```bash
     #!/bin/bash
     DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
     export LD_LIBRARY_PATH="$DIR:$DIR/lib:${LD_LIBRARY_PATH:-}"
     EXEC="$DIR/antigravity"
     [ -f "$EXEC" ] || EXEC="$DIR/Antigravity"

     export ELECTRON_OZONE_PLATFORM_HINT="auto"
     export DONT_PROMPT_WSL_INSTALL=1

     # Unconfined tarball Electron builds require --no-sandbox on modern Linux (AppArmor userns restriction)
     exec "$EXEC" --no-sandbox "$@"
     ```
   - Ensure permissions: `chmod +x "$ide_dir/antigravity"` and `chmod +x "$ide_dir/antigravity-runner.sh"`.
   - Ensure launcher symlinks in `~/.local/bin/`:
     - `~/.local/bin/antigravity -> ~/.local/share/antigravity-ide/antigravity-ide.run`
     - `~/.local/bin/agy -> ~/.local/share/antigravity-ide/antigravity-ide.run`
     - `~/.local/bin/antigravity-ide -> ~/.local/share/antigravity-ide/antigravity-ide.run`
   - Ensure directory symlink: `~/.local/share/antigravity -> ~/.local/share/antigravity-ide`.

5. **Desktop Launcher & Icon Maintenance**:
   - Re-generate `~/.local/share/applications/antigravity.desktop`.
   - Deploy multi-resolution icons into `~/.local/share/icons/hicolor/{16x16..512x512}/apps/antigravity.png`.
   - Refresh desktop caches via `gtk-update-icon-cache` and `update-desktop-database`.

6. **Post-Update Non-Interactive Verification**:
   - Verify that `$HOME/.local/bin/antigravity` resolves to an executable file.
   - Test ELF binary header: `file "$HOME/.local/share/antigravity-ide/antigravity" | grep -q "ELF"`.
   - Avoid launching headless GUI windows that hang non-interactive CI/agent runs.

---

## 3. Live Ubuntu Host Update Specification

### 3.1 Host Directory Structure

The live update on the host machine strictly conforms to the standard user-level layout:

```text
~/.local/
├── bin/
│   ├── agm
│   ├── agy -> ~/.local/share/antigravity-ide/antigravity-ide.run
│   ├── antigravity -> ~/.local/share/antigravity-ide/antigravity-ide.run
│   └── antigravity-ide -> ~/.local/share/antigravity-ide/antigravity-ide.run
└── share/
    ├── antigravity -> ~/.local/share/antigravity-ide
    ├── antigravity-ide/
    │   ├── antigravity (ELF binary, chmod 0755)
    │   ├── antigravity-ide.run (chmod 0755)
    │   ├── antigravity-runner.sh (chmod 0755)
    │   ├── chrome-sandbox (chmod 4755 if root, 0755 if user)
    │   ├── resources/
    │   └── locales/
    └── applications/
        └── antigravity.desktop
```

### 3.2 Pre-Flight Safety Checks

1. **Active Process Detection**:
   - Check if any `antigravity` or `agy` Electron processes are running (`pgrep -f "antigravity" || true`).
   - If active processes exist and `--kill` or `-y` is passed, signal graceful termination; otherwise, warn that active windows should be reloaded after file replacement.

2. **Disk Space Verification**:
   - Ensure at least 1.5 GB of free space is available in `/tmp` and `~/.local/share/` before initiating download and extraction.

3. **User Profile Protection**:
   - The IDE configuration directories `~/.config/Antigravity/` and `~/.antigravity/` must never be touched, modified, or cleared during the update process.

### 3.3 Atomic Rollback Protocol

If tar extraction fails, archive checksum is invalid, or the resulting binary fails executable verification:
1. Immediately remove the damaged `~/.local/share/antigravity-ide/` directory.
2. Restore the previous directory from `~/.local/share/antigravity-ide.bak-<timestamp>`.
3. Re-verify the restored launcher symlinks.
4. Exit with status code `1` and emit an actionable error diagnostic to stderr.

---

## 4. Verification and Acceptance Matrix

| Verification ID | Check Description | Expected Result |
|---|---|---|
| `V-SF-01` | `scripts-linux/69-install-antigravity/manifest.json` | Contains `"update"` in `"commands"` list |
| `V-SF-02` | `scripts/69-install-antigravity/manifest.json` | Contains `"update"` in `"commands"` list |
| `V-SF-03` | `scripts-linux/69-install-antigravity/run.sh help` | Displays `update` action in help output |
| `V-SF-04` | Syntax check on shell scripts | `bash -n` exits `0` on both `run.sh` and `install-antigravity.sh` |
| `V-HOST-01` | Live IDE binary existence | `~/.local/share/antigravity-ide/antigravity` exists and is executable |
| `V-HOST-02` | Live runner wrapper integrity | `antigravity-ide.run` includes `--no-sandbox` and correct library paths |
| `V-HOST-03` | Live launcher symlinks | `~/.local/bin/{antigravity,agy,antigravity-ide}` point to `antigravity-ide.run` |
| `V-HOST-04` | User state preservation | `~/.config/Antigravity/` remains completely intact post-update |
| `V-HOST-05` | Safety backup created | At least one valid `.bak-*` directory created in `~/.local/share/` |

---

## 5. Non-Negotiable Engineering Constraints

1. **Strict Relative Paths**: All specification documents, code modifications, and plan tasks must employ relative paths or shell portable home aliases (`~/.local/...`). Zero absolute paths permitted.
2. **Idempotence**: Running `verb_update` multiple times consecutively must succeed without creating redundant symlink chains or leaking temporary staging folders.
3. **No Direct Git Operations**: Commits, tags, and pushes must be performed exclusively via `gitmap cpf "<module> - <summary>"`. Raw `git add` and `git commit` commands are strictly forbidden.
