# Notepad++ Installation and Settings — Shortcut Standardization

## Overview

This spec standardizes the three Notepad++ installation variants and their
shortcut naming across all CLI commands, help text, specs, and documentation.
NPP always means Notepad++.

---

## Installation Variants

| Command        | Shortcut Label   | Installs Binary | Syncs Settings |
|----------------|------------------|-----------------|----------------|
| `npp`          | NPP + Settings   | ✅               | ✅              |
| `npp-settings` | NPP Settings     | ❌               | ✅              |
| `install-npp`  | Install NPP      | ✅               | ❌              |

---

## Shortcut Naming Standard

| Label            | Meaning                             |
|------------------|-------------------------------------|
| NPP + Settings   | Notepad++ with settings             |
| NPP Settings     | Notepad++ settings only             |
| Install NPP      | Notepad++ install only              |

**Rules:**
- NPP always means Notepad++.
- Use these exact labels in all specs, commands, help text, and UI.
- Do not mix or abbreviate variant labels.

---

## Settings Package

The Notepad++ settings are distributed as a zip archive:
`data/npp-settings/npp-settings.zip`

### Zip Contents

The archive contains standard Notepad++ configuration files:
- `config.xml` — main configuration
- `contextMenu.xml` — right-click menu
- `functionList/*.xml` — language function definitions
- `backup/` — session backup files
- Theme files (e.g., `Dracula.xml`)

### Extraction Path

The zip is extracted to the user-specific roaming AppData directory:

| Variable    | Path                                          |
|-------------|-----------------------------------------------|
| `%APPDATA%` | `C:\Users\{user}\AppData\Roaming`             |
| Target      | `C:\Users\{user}\AppData\Roaming\Notepad++`   |

- `{user}` is replaced by the active Windows username at runtime.
- The `APPDATA` environment variable is used for resolution.
- If `APPDATA` is not set, the command exits with an error.

### Extraction Behavior

- Creates the target directory if it does not exist.
- Extracts all files from the zip, preserving directory structure.
- **Overwrites** existing files with the same name.
- Falls back to copying loose files from the source directory if the zip is missing.

---

## CLI Routing

```
executeInstall(opts):
  if opts.Tool == "npp-settings" → runNppSettingsOnly()
  if opts.Tool == "install-npp"  → install binary only, skip settings
  if opts.Tool == "npp"          → install binary, then sync settings
```

---

## Post-Install Verification

After binary installation (`npp` or `install-npp`):
1. Run `--version` check via `toolBinaryName` lookup.
2. Verify exe exists at `C:\Program Files\Notepad++\notepad++.exe`.
3. Report success or failure.

---

## File Layout

| File                               | Purpose                           |
|------------------------------------|-----------------------------------|
| `constants/constants_install.go`   | Tool names, descriptions, messages|
| `cmd/install.go`                   | Command routing, flag parsing     |
| `cmd/installnpp.go`               | NPP-specific install and settings |
| `cmd/installverify.go`            | Post-install exe verification     |
| `data/npp-settings/npp-settings.zip` | Bundled settings archive        |
| `helptext/install.md`             | CLI help text                     |

---

## Constraints

- Settings sync is Windows-only (`runtime.GOOS == "windows"`).
- All functions under 15 lines.
- All files under 200 lines.
- No hardcoded usernames; resolve via `APPDATA` env var.
- Positive logic only.

---

## Ambiguities Noted

1. **Overwrite behavior**: Settings extraction overwrites existing files.
   No backup of prior settings is made. A rollback mechanism is not yet implemented.
2. **Multiple zip versions**: Only one settings zip is currently supported.
   Version tagging for settings packages is a future enhancement.

---

## See Also

- [81-install.md](../01-app/81-install.md) — Install command specification
- [install.md](../../gitmap-v28/helptext/install.md) — CLI help text
