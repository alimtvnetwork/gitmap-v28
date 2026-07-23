# Environment Variable Setup

## Overview

This document specifies an `env` command that manages persistent, cross-platform environment variables and PATH entries. The command allows users to define a custom drive or directory (e.g., `E:\tools` or `/opt/tools`) where the tool is installed, and ensures the environment variable is always set — automatically registering it if missing.

This is particularly useful for portable installations where the binary lives on a non-standard drive or path, and the user needs the system to "just work" without manual `PATH` or environment variable configuration.

---

## Core Concept

```
User specifies:    E:\gitmap-v27           (or /opt/gitmap-v27)
Tool ensures:      GITMAP_HOME=E:\gitmap-v27   is set persistently
                   E:\gitmap-v27 is in PATH    if not already present
```

If the environment variable is already set and valid → no action taken.
If the environment variable is missing or points to a stale path → auto-register it.

---

## Command Interface

```
<tool> env                         # Show all managed variables
<tool> env set <KEY> <VALUE>       # Set a persistent environment variable
<tool> env remove <KEY>            # Remove a managed environment variable
<tool> env path add <DIR>          # Add a directory to PATH persistently
<tool> env path remove <DIR>       # Remove a directory from PATH
<tool> env home <DIR>              # Set <TOOL>_HOME and add to PATH
<tool> env doctor                  # Verify all managed variables are active
```

### Alias: `ev`

---

## `env home` — Drive/Directory Setup

This is the primary power feature. The user specifies where the tool lives:

```
$ <tool> env home E:\gitmap-v27
```

This does:

1. Validates the directory exists (or creates it with confirmation)
2. Sets `<TOOL>_HOME=E:\gitmap-v27` persistently
3. Adds `E:\gitmap-v27` to PATH if not already present
4. Records the registration in `env-registry.json`
5. Prints activation command

### Terminal Output

```
$ <tool> env home E:\gitmap-v27

  Setting <TOOL>_HOME...

  [+] <TOOL>_HOME = E:\gitmap-v27

  Registering PATH...

  [+] Windows Registry (User PATH)
  [+] PowerShell profile
  [=] Git Bash profile (already registered)

  ============================================
  Environment configured!

  To activate in this session:

    $env:GITMAP_HOME = "E:\gitmap-v27"
    $env:Path = "E:\gitmap-v27;" + $env:Path

  Or restart your terminal.
  ============================================
```

### Auto-Registration on Startup

When the tool starts, it checks if `<TOOL>_HOME` is set. If not, and if the tool can determine its own location, it auto-registers:

```go
func ensureHomeEnv() {
    home := os.Getenv("GITMAP_HOME")
    if home != "" && dirExists(home) {
        return // Already configured and valid
    }

    // Resolve from binary location
    binaryDir := resolveBinaryDir()
    if binaryDir == "" {
        return
    }

    // Auto-register
    setEnvPersistent("GITMAP_HOME", binaryDir)
    addToPath(binaryDir)
    fmt.Printf("  Auto-configured GITMAP_HOME=%s\n", binaryDir)
}
```

---

## Platform-Specific Persistence

### Windows

Environment variables are set via the Windows Registry:

```go
// User-level variable
key, _ := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
key.SetStringValue("GITMAP_HOME", value)

// Notify the system of the change
syscall.SendMessage(syscall.HWND_BROADCAST, syscall.WM_SETTINGCHANGE, 0, "Environment")
```

PATH is updated in both:
- Registry (`HKCU\Environment\Path`)
- PowerShell profile (`$PROFILE`)
- Git Bash profiles (`~/.bashrc`, `~/.bash_profile`)

### Unix (Linux / macOS)

Environment variables are persisted by writing to shell profiles:

```bash
# Appended to ~/.bashrc, ~/.zshrc, etc.
export GITMAP_HOME="/opt/gitmap-v27"   # <tool>-env
```

The marker comment (`# <tool>-env`) enables idempotent updates and clean removal.

### Shell Override Flag

```
<tool> env set KEY VALUE --shell bash    # Only write to .bashrc
<tool> env set KEY VALUE --shell zsh     # Only write to .zshrc
```

---

## Environment Registry

The tool maintains an `env-registry.json` file to track all managed variables:

```json
{
  "variables": [
    {
      "key": "GITMAP_HOME",
      "value": "E:\\gitmap-v27",
      "createdAt": "2026-04-09T14:30:00Z",
      "platforms": ["registry", "powershell-profile", "git-bash"]
    },
    {
      "key": "GITMAP_DATA",
      "value": "E:\\gitmap-v27\\data",
      "createdAt": "2026-04-09T14:30:00Z",
      "platforms": ["registry", "powershell-profile"]
    }
  ],
  "pathEntries": [
    {
      "directory": "E:\\gitmap-v27",
      "createdAt": "2026-04-09T14:30:00Z"
    }
  ]
}
```

This registry enables:
- `env remove` to know which profiles to clean
- `env doctor` to verify all registrations are still active
- Uninstall to remove all managed variables

---

## `env doctor` — Verification

```
$ <tool> env doctor

  Checking managed environment variables...

  [OK]   GITMAP_HOME = E:\gitmap-v27 (directory exists)
  [OK]   E:\gitmap-v27 is in PATH

  Checking shell profiles...

  [OK]   PowerShell profile: GITMAP_HOME registered
  [OK]   Git Bash profile: GITMAP_HOME registered
  [WARN] Zsh profile: not found (not applicable on Windows)

  All checks passed.
```

### Failure Output

```
$ <tool> env doctor

  Checking managed environment variables...

  [FAIL] GITMAP_HOME = E:\gitmap-v27 (directory does NOT exist)
  [WARN] E:\gitmap-v27 is NOT in PATH

  Suggested fix:

    <tool> env home E:\gitmap-v27    # Re-register with valid path

  1 failure, 1 warning.
```

---

## `env set` / `env remove`

### Set

```
$ <tool> env set EDITOR "code --wait"

  [+] EDITOR = code --wait

  Registered in:
    [+] Windows Registry (User)
    [+] PowerShell profile
    [+] Git Bash profile

  To activate: restart your terminal or run:
    $env:EDITOR = "code --wait"
```

### Remove

```
$ <tool> env remove EDITOR

  [-] EDITOR removed from:
    [-] Windows Registry (User)
    [-] PowerShell profile
    [-] Git Bash profile

  To deactivate in this session:
    Remove-Item Env:\EDITOR
```

---

## Integration with Install Scripts

The install scripts (`install.ps1`, `install.sh`) should call `env home` logic after installation:

```powershell
# install.ps1 — after binary extraction
& "$InstallDir\<tool>.exe" env home "$InstallDir" --quiet
```

```bash
# install.sh — after binary extraction
"$INSTALL_DIR/<tool>" env home "$INSTALL_DIR" --quiet 2>/dev/null || true
```

The `--quiet` flag suppresses the detailed output since the installer already prints its own summary.

---

## Integration with CI/CD Pipeline

### Release Workflow

The release pipeline should verify that the installed binary can set its own home:

```yaml
- name: Verify env home
  run: |
    ./dist/<binary>-*-linux-amd64 env home /tmp/test-install --quiet
    test -n "$(<binary upper>_HOME)" || echo "::warning::env home did not persist (expected in CI)"
```

### Pipeline Spec File

This feature should be documented in the pipeline README alongside other specs so any AI implementing CI/CD knows to:

1. Include `env home` in the install script flow
2. Test that `env doctor` passes after installation
3. Use `<TOOL>_HOME` as the canonical way to find the tool

---

## File Organization

```
cmd/
  env.go              # Subcommand routing (env, env set, env remove, env home, env doctor)
  envops.go           # CRUD operations for variables and PATH
  envplatform_windows.go  # Registry + profile writes (Windows)
  envplatform_unix.go     # Shell profile writes (Linux/macOS)

constants/
  constants_env.go    # All env-related messages, defaults, and SQL

model/
  envregistry.go      # EnvRegistry struct and JSON serialization
```

---

## Constraints

- All environment variable writes must be idempotent (marker-based)
- Registry writes on Windows must broadcast `WM_SETTINGCHANGE`
- Shell profile writes use `# <tool>-env` markers for clean removal
- `env home` validates the directory exists before registering
- Auto-registration at startup is silent (no error on failure)
- `env doctor` never modifies the system — read-only verification
- PowerShell 5.1 compatibility (no `??`, no multi-arg `Join-Path`)
- Bash 3.2+ compatibility (macOS ships old bash)
- The `env-registry.json` file lives in the tool's data directory
