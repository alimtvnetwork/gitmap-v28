# Subtask 02: Antigravity Tools Update Script & Command Execution Plan

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/update-antigravity-ide.md`  
> **Target Path:** `.ai-memory/plans/subtasks/update-antigravity-ide/02-antigravity-manager-tools.md`  
> **Status:** READY  
> **Estimated Effort:** 3 Engineering Units  
> **Dependencies:** `02-spec/21-app/update-antigravity-ide/01-architecture-and-gitmap-command.md`

---

## 1. Objectives & Scope

1. Author and integrate a robust, unattended shell script: `scripts/update-antigravity-ide.sh` for Linux and macOS environments.
2. Provide synchronized copies in both `gitmap` (`scripts/update-antigravity-ide.sh`) and companion repository `antigravity-manager` (`scripts/update-antigravity-ide.sh`).
3. Connect `antigravity-manager` CLI and toolchain to invoke this updater via `agm update ide` or `npm run update:ide`.
4. Implement automatic pre-flight detection of active `antigravity` processes with graceful SIGTERM termination or `--kill-running` support.
5. Implement atomic staging and directory swapping to prevent partial binary corruption if network or decompression fails.
6. Reconcile the user launcher script (`~/.local/bin/antigravity`) with Electron AppArmor sandbox bypass (`--no-sandbox`) and display server flags.
7. Reconcile XDG desktop entries (`~/.local/share/applications/antigravity.desktop`) and icon resources.
8. Validate execution with strict relative paths and zero absolute paths across all scripts and tests.

---

## 2. Detailed Technical Plan

### 2.1 Script Design: `scripts/update-antigravity-ide.sh`

The update script must be standalone POSIX/Bash-compatible with no external language runtime dependencies (Go, Node.js, Python not required).

#### 2.1.1 Command-Line Flags
- `-f, --force`: Force reinstallation and overwrite existing files even if version matches.
- `-n, --dry-run`: Preview download URL, target paths, and extraction steps without making modifications.
- `-v, --version <ver>`: Explicit Antigravity IDE version string (default: `2.19.1` or dynamically fetched latest).
- `-b, --build-id <id>`: Explicit build ID (default: `6046815158665216`).
- `-k, --kill-running`: Automatically terminate active Antigravity IDE processes before starting extraction.
- `-q, --quiet`: Suppress non-error output.
- `-h, --help`: Display usage and flag descriptions.

#### 2.1.2 Pre-Flight Inspection & Process Safeguards
1. **Architecture Detection**:
   ```bash
   OS="$(uname -s)"
   ARCH="$(uname -m)"
   case "${ARCH}" in
       x86_64) ARCH_TAG="x64" ;;
       aarch64|arm64) ARCH_TAG="arm64" ;;
       *) echo "Error: Unsupported architecture ${ARCH}" >&2; exit 1 ;;
   esac
   ```
2. **Tooling Prerequisites**:
   - Ensure `tar`, `curl` or `wget`, and `pgrep` are available in system `$PATH`.
3. **Active Process Guard**:
   - Inspect for running processes matching `antigravity`:
   ```bash
   if pgrep -f "antigravity" >/dev/null 2>&1; then
       if [ "${KILL_RUNNING}" = "true" ]; then
           pkill -TERM -f "antigravity" 2>/dev/null || true
           sleep 1
           pkill -KILL -f "antigravity" 2>/dev/null || true
       else
           echo "Warning: Antigravity IDE process is running. Close it or use --kill-running."
           exit 2
       fi
   fi
   ```

#### 2.1.3 Download & Atomic Decompression
1. **Target Paths**:
   - Install root: `$HOME/.local/share/antigravity-ide`
   - Staging root: `$HOME/.local/share/antigravity-ide.tmp`
   - Backup root: `$HOME/.local/share/antigravity-ide.bak`
   - Bin directory: `$HOME/.local/bin`
   - Desktop directory: `$HOME/.local/share/applications`
2. **Download Execution**:
   - Construct download URL using hub pattern:
     `https://storage.googleapis.com/antigravity-public/antigravity-hub/${VERSION}-${BUILD_ID}/${PLATFORM}/${ARTIFACT}`
   - Stream download to a temporary archive file in a secured scratch directory.
3. **Extraction & Verification**:
   - Clear staging directory `antigravity-ide.tmp`.
   - Extract archive into staging directory:
     `tar -xzf "${ARCHIVE_FILE}" -C "${STAGING_DIR}" --strip-components=1`
   - Validate presence and permissions of executable:
     `chmod 0755 "${STAGING_DIR}/antigravity"`
     `test -x "${STAGING_DIR}/antigravity"`
4. **Atomic Swap & Rollback**:
   - Move current `${INSTALL_DIR}` to `${BACKUP_DIR}` if it exists.
   - Move `${STAGING_DIR}` to `${INSTALL_DIR}`.
   - If move succeeds, remove `${BACKUP_DIR}` and archive.
   - If move fails, restore `${BACKUP_DIR}` back to `${INSTALL_DIR}`.

#### 2.1.4 Launcher Script & Desktop Integration
1. **Regenerate Launcher (`~/.local/bin/antigravity`)**:
   ```bash
   mkdir -p "$HOME/.local/bin"
   cat << 'EOF' > "$HOME/.local/bin/antigravity"
   #!/bin/bash
   DIR="$HOME/.local/share/antigravity-ide"
   export LD_LIBRARY_PATH="${DIR}:${DIR}/lib:${LD_LIBRARY_PATH:-}"
   export ELECTRON_OZONE_PLATFORM_HINT="auto"
   export DONT_PROMPT_WSL_INSTALL=1
   exec "${DIR}/antigravity" --no-sandbox "$@"
   EOF
   chmod +x "$HOME/.local/bin/antigravity"
   ```
2. **Regenerate Desktop Entry (`~/.local/share/applications/antigravity.desktop`)**:
   - Ensure `Exec` points to `~/.local/bin/antigravity %F`.
   - Ensure `Icon` points to `~/.local/share/antigravity-ide/resources/app/resources/linux/code.png` or bundled SVG icon.
   - Update desktop database if `update-desktop-database` is present.

### 2.2 Integration with Antigravity Tools (`antigravity-manager`)

1. **Companion Script Synchronization**:
   - Mirror `scripts/update-antigravity-ide.sh` into `../antigravity-manager/scripts/update-antigravity-ide.sh`.
2. **Package Script Entry**:
   - In `antigravity-manager/package.json`:
     Add `"update:ide": "bash ./scripts/update-antigravity-ide.sh"` under scripts.
3. **CLI Command Binding**:
   - In `antigravity-manager` CLI dispatcher:
     Add subcommand `agm update ide` invoking the update script with passed options.

---

## 3. Verification Criteria & Acceptance Checklist

- [ ] Script passes shell syntax checking: `bash -n scripts/update-antigravity-ide.sh`.
- [ ] `--dry-run` prints planned operations without modifying any directories or files.
- [ ] Script successfully downloads and extracts the IDE archive when run with `--force`.
- [ ] Atomic directory swap preserves existing installation in event of extraction errors.
- [ ] Active process detection catches running `antigravity` process and respects `--kill-running`.
- [ ] Launcher script in `~/.local/bin/antigravity` includes `--no-sandbox` and launches successfully.
- [ ] Desktop launcher entry is generated with correct relative/home path expansion.
- [ ] Zero absolute paths verified across script and plan documentation.

---

## 4. File Bounding Box & Impact Summary

| File Path (Strictly Relative) | Repository | Action | Description |
| :--- | :--- | :--- | :--- |
| `scripts/update-antigravity-ide.sh` | `gitmap` | Create | Standalone POSIX/Bash IDE updater script |
| `scripts/update-antigravity-ide.sh` | `antigravity-manager` | Create | Mirrored updater script in companion repo |
| `package.json` | `antigravity-manager` | Modify | Add `update:ide` npm script target |
| `.ai-memory/plans/subtasks/update-antigravity-ide/02-antigravity-manager-tools.md` | `gitmap` | Create | Subtask execution plan for Antigravity tools updater |
