# Issue 07: Custom Installer Execution Failure ('Unknown Tool' on `gitmap install <custom-tool>`)

## 1. Problem Description

When users created a custom installer via `gitmap install add <name>` (e.g. `gitmap install add alim1`), running `gitmap install alim1` resulted in a fatal execution failure:

```text
  ✖ Unknown tool: 'alim1'

  Use 'gitmap install --list' to see all supported tools.
  Use 'gitmap install --help' for usage examples.

gitmap: [E9000:EXECUTION] unknown tool 'alim1'. Use 'gitmap install --list' to see available tools
  origin: cmd/install.go:113
```

## 2. Root Cause Analysis

1. **Hardcoded Validation Barrier:**
   In `gitmap/cmd/install.go`, `validateToolName` only checked `constants.InstallToolDescriptions` and predefined aliases (`clean-code`, `build-essential`, `ag-m`). It never inspected SQLite (`installer_scripts`) to check if the target tool was an installed user custom script.
2. **Missing Custom Installer Dispatch & Execution:**
   There was no execution pipeline to resolve the stored OS script (`win` with PowerShell on Windows, `ubuntu` or `unix` with bash on Linux/macOS) from `model.InstallerScript.Instructions`.
3. **Cobra Delegation Gap in `installer` Subcommand:**
   `installer` command in `cmd/installer.go` did not forward unrecognized tool names to `runInstall`, causing `gitmap installer <name>` and `gitmap in <name>` to fail with unknown command errors.

## 3. Resolution

1. Implemented `gitmap/cmd/install_custom_exec.go`:
   - `findCustomInstaller(slug)`: Queries SQLite `installer_scripts` via `db.GetInstallerBySlug`.
   - `hasCustomInstaller(slug)`: Used in `isKnownInstallTool` to recognize registered custom tools.
   - `executeCustomInstaller(script, opts)`: Unmarshals JSON instructions, resolves the script for the current OS (`runtime.GOOS`), runs the process via PowerShell (Windows) or bash (Linux/Unix), and records installation success into `installation.db` via `splitDB.SaveInstalledTool`.
2. Updated `gitmap/cmd/install.go`:
   - Added custom script execution routing in `runInstall` before `validateToolName`.
   - Updated `isKnownInstallTool` to recognize custom tools.
3. Updated `gitmap/cmd/installer.go`:
   - Added delegation in `runInstaller` and `RunInstallerCLI` to forward non-Cobra tool names to `runInstall`.
4. Unit Tests in `gitmap/cmd/install_custom_exec_test.go`:
   - `TestParseInstructionsMap`
   - `TestResolveOSScriptForPlatform`
   - `TestResolveOSScriptForPlatform_UnsupportedOS`
   - `TestExecuteCustomInstaller_DryRun`
   - `TestRunInstall_CustomInstallerIntegration`
   - `TestExecuteCustomInstaller_LiveExecution`
