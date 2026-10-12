# Subtask 03 — scripts-fixer-update (task update-antigravity-ide, scripts-fixer)

## File Box
All paths relative to repository root (`scripts-fixer`):
- `scripts-linux/69-install-antigravity/manifest.json`
- `scripts/69-install-antigravity/manifest.json`
- `scripts-linux/69-install-antigravity/run.sh`
- `scripts/os/ubuntu/install-antigravity.sh`

## Objective & Scope
Implement the `update` verb in the `scripts-fixer` toolchain for Google Antigravity IDE & agy CLI automation module (ID: 69). This enables developers and automated pipelines to trigger in-place upgrades of the Antigravity desktop IDE and CLI companion on Linux/Ubuntu platforms without completely nuking existing settings or failing on early-exit checks.

## Specification Reference
- `02-spec/21-app/update-antigravity-ide/02-scripts-fixer-and-live-system-update.md`

---

## Detailed Implementation Steps

### Step 1: Update Manifest Files
1. Edit `scripts-linux/69-install-antigravity/manifest.json`:
   - Append `"update"` to the `"commands"` array.
   - Ensure alphabetical or logical placement immediately following `"install"`.
2. Edit `scripts/69-install-antigravity/manifest.json`:
   - Append `"update"` to the `"commands"` array to maintain parity between cross-platform and Linux-specific manifests.
3. Validate JSON formatting:
   - Run python JSON validator to ensure valid syntax:
     ```bash
     python3 -m json.tool scripts-linux/69-install-antigravity/manifest.json > /dev/null
     python3 -m json.tool scripts/69-install-antigravity/manifest.json > /dev/null
     ```

### Step 2: Implement `verb_update()` in `run.sh`
1. Edit `scripts-linux/69-install-antigravity/run.sh`:
   - Define `verb_update()`:
     - Register install path telemetry via `write_install_paths`.
     - Log initiation: `log_info "[69] Starting Antigravity IDE update"`.
     - Locate `scripts/os/ubuntu/install-antigravity.sh` via relative path `$REPO_ROOT/scripts/os/ubuntu/install-antigravity.sh`.
     - Execute the installer script passing `--update "$@"`.
     - On success: log `log_ok "[69] Antigravity updated successfully"`, touch `$INSTALLED_MARK`, and return 0.
     - On failure: log `log_err "[69] Antigravity update failed"` and return 1.
2. Update `show_help()`:
   - Under `ACTIONS:`, add:
     ```text
       update                 Update Antigravity IDE and CLI to latest release
     ```
3. Update CLI argument dispatcher:
   - In `case "${1:-help}" in`, add route:
     ```bash
     update) verb_update "${@:2}";;
     ```

### Step 3: Implement In-Place `--update` Mode in `install-antigravity.sh`
1. Edit `scripts/os/ubuntu/install-antigravity.sh`:
   - Parse `--update`, `-u`, or `update` in CLI arguments loop, setting `IS_UPDATE=true`.
   - Update `--help` text to document the `--update` option.
2. Add safe pre-upgrade backup routine:
   - Create timestamped backup of existing `~/.local/share/antigravity-ide` directory:
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
3. Update main installation lifecycle:
   - When `IS_UPDATE` is true:
     - Skip destructive wipe (`clean_existing_installation`).
     - Trigger `backup_existing_installation`.
     - Fetch target download URL and download archive (using aria2c with 16 connections or curl fallback).
     - Extract into temporary staging directory and sync new binaries into `~/.local/share/antigravity-ide`.
     - Recreate or ensure launcher wrapper `antigravity-runner.sh` and `antigravity-ide.run` with `--no-sandbox`.
     - Verify and refresh launcher symlinks in `~/.local/bin/` (`antigravity`, `agy`, `antigravity-ide`).
     - Re-generate `~/.local/share/applications/antigravity.desktop` and refresh icon caches.
     - Execute non-interactive post-update smoke check.

---

## Validation & Verification Criteria

1. **Syntax Integrity**:
   - `bash -n scripts-linux/69-install-antigravity/run.sh` exits `0`.
   - `bash -n scripts/os/ubuntu/install-antigravity.sh` exits `0`.
2. **Help Documentation Check**:
   - `scripts-linux/69-install-antigravity/run.sh help` displays the `update` verb in output.
   - `scripts/os/ubuntu/install-antigravity.sh --help` displays `--update` option in output.
3. **Manifest Conformance**:
   - Both `manifest.json` files parse cleanly and contain `"update"` in their commands list.
4. **Dry-Run / Non-Destructive Check**:
   - Invoking `./run.sh 69 check` continues to verify current installation correctly.

---

## Non-Negotiable Rules

- **Zero Absolute Paths**: Strictly use relative paths (`scripts-linux/69-install-antigravity/...`, `scripts/os/ubuntu/...`) or portable home paths (`~/.local/...`). Never hardcode system-specific absolute paths.
- **Zero Raw Git Commits**: Never run `git add` or `git commit`. All changes are committed and pushed via `gitmap cpf "<module> - <summary>"`.
- **Error Propagation**: Ensure all functions return explicit exit codes (`0` or `1`) with `set -u` compatibility.

## Deliverables
- Clean file modifications across the 4 listed files in `scripts-fixer`.
- Verification command outputs demonstrating passing syntax checks and manifest validation.
