# Custom Installer Execution Engine & Installation Database Synchronization

**Updated:** 2026-09-11  
**Specification:** `spec/20-custom-installer-registry/`  
**Reference Issue:** `.lovable/issues/07-custom-installer-execution-unknown-tool.md`  
**Reference Plan:** `.lovable/plans/completed/84-custom-installer-execution-engine.md`

---

## 1. Problem & Root Cause

When a user creates a custom installer via `gitmap install add <name>` (or `gitmap installer add <name>`), the script instructions are saved to the `installer_scripts` table in SQLite with target OS instructions (`win`, `unix`, `ubuntu`, or `all`).

Previously, executing `gitmap install <name>` failed with:
```text
  ✖ Unknown tool: '<name>'
gitmap: [E9000:EXECUTION] unknown tool '<name>'. Use 'gitmap install --list' to see available tools
```

### Root Causes:
1. `validateToolName` in `cmd/install.go` strictly consulted static map `constants.InstallToolDescriptions` and hardcoded aliases without querying SQLite for user-defined custom installers.
2. `cmd/install.go` lacked an execution dispatcher to parse the multi-OS JSON instructions and run the appropriate platform script via PowerShell (Windows) or bash (Linux/macOS).
3. Tool installation status was not persisted in `installation.db` after successful execution, causing `install ls` to always show custom tools as uninstalled (`○`).

---

## 2. Architecture & Resolution

### A. Dynamic Custom Installer Lookup
In `cmd/install_custom_exec.go`:
- `findCustomInstaller(slug)`: Queries SQLite `installer_scripts` using `db.GetInstallerBySlug(cleanSlug)`.
- `hasCustomInstaller(slug)`: Integrated into `isKnownInstallTool(tool)` so custom installers pass tool validation across `install` and `uninstall`.

### B. OS Script Resolution & Execution
`resolveOSScriptForPlatform(script *model.InstallerScript)` unmarshals `script.Instructions`:
- **Windows (`runtime.GOOS == "windows"`):** Resolves `win` or `all` script; executes using `powershell -NoProfile -ExecutionPolicy Bypass -Command <script>`.
- **Linux (`runtime.GOOS == "linux"`):** Checks `/etc/os-release` for Ubuntu/Debian host. If Ubuntu, runs `ubuntu` script; otherwise runs `unix` or `all` script using `bash -c <script>`.
- **macOS (`runtime.GOOS == "darwin"`):** Resolves `unix` or `all` script; executes using `bash -c <script>`.
- Unsupported OS combinations return clean, domain-specific `apperror.NewValidationError` without panicking.

### C. Installation State Tracking
Upon successful script execution:
- Calls `splitDB.SaveInstalledTool(script.Slug, script.Version, "custom")` on `installation.db`.
- `gitmap install ls` immediately updates the tool's status indicator to `●` (installed) in the `Custom Tools` section.

### D. Dual CLI Parity (`install` & `installer`)
- `installer` command in `cmd/installer.go` checks `isDelegatedInstallArg(args)` and forwards any tool name or unhandled command to `runInstall(args)`.
- Guarantees 100% parity: `gitmap install <name>` and `gitmap installer <name>` execute identically.
