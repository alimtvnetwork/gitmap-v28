# Settings Sync

## Overview

The `gitmap-v27 install <tool>-settings` command family synchronises bundled
application settings from the repository's `settings/` directory to the
correct platform-specific config location. This provides a one-command way
to replicate a developer's editor, terminal, and streaming setup on a new
machine.

## Supported Tools

| Command                        | Alias | Source directory              | Description                          |
|--------------------------------|-------|------------------------------|--------------------------------------|
| `gitmap-v27 install npp`           | `in`  | `settings/01 - notepad++`    | Install Notepad++ **and** sync settings |
| `gitmap-v27 install npp-settings`  | `in`  | `settings/01 - notepad++`    | Sync Notepad++ settings only         |
| `gitmap-v27 install install-npp`   | `in`  | --                           | Install Notepad++ only (no settings) |
| `gitmap-v27 install vscode-settings` | `in` | `settings/02 - vscode`      | Sync VS Code settings and extensions |
| `gitmap-v27 install obs-settings`  | `in`  | `settings/03 - obs`          | Sync OBS Studio profiles and scenes  |
| `gitmap-v27 install wt-settings`   | `in`  | `settings/04 - windows-terminal` | Sync Windows Terminal settings.json |

## Settings Directory Layout

```
settings/
  01 - notepad++/
    02. Notepad++ settings.zip    # zipped config (themes, shortcuts, etc.)
    readme.txt
  02 - vscode/
    settings.json                 # VS Code user settings
    keybindings.json              # keyboard shortcuts (optional)
    extensions.txt                # one extension ID per line
    readme.txt
  03 - obs/
    25 - Personal 15 Jul 2025 Malaysia.zip   # bundled OBS settings archive
    readme.txt
  04 - windows-terminal/
    settings.json                 # Windows Terminal settings
    readme.txt
```

Each subfolder is numbered (`01`, `02`, `03`, `04`) for consistent ordering.
A `readme.txt` in each folder documents the contents and target paths.

## Target Paths

### Notepad++ (`npp`, `npp-settings`)

| Platform | Target                                |
|----------|---------------------------------------|
| Windows  | `%APPDATA%\Notepad++`                 |
| Others   | Not supported (Windows-only tool)     |

Settings are extracted from a bundled `.zip` file. If the zip is missing,
the sync falls back to copying loose files from the source directory.

### VS Code (`vscode-settings`)

| Platform | Target                                           |
|----------|--------------------------------------------------|
| Windows  | `%APPDATA%\Code\User`                            |
| macOS    | `~/Library/Application Support/Code/User`        |
| Linux    | `~/.config/Code/User`                            |

Files synced:

- `settings.json` -- user settings
- `keybindings.json` -- keyboard shortcuts (if present)
- Any other files in the source directory (excluding `readme.txt`)

If an `extensions.txt` file is present, each line is treated as a VS Code
extension ID and installed via `code --install-extension <id> --force`.
Lines starting with `#` and blank lines are skipped.

### OBS Studio (`obs-settings`)

| Platform | Target                                           |
|----------|--------------------------------------------------|
| Windows  | `%APPDATA%\obs-studio`                           |
| macOS    | `~/Library/Application Support/obs-studio`       |
| Linux    | `~/.config/obs-studio`                           |

#### Zip extraction and routing

The OBS sync command looks for a `.zip` file in the source directory. When
found, it follows a 4-step extraction process:

1. **Extract** the `.zip` to a temporary directory.
2. **Route `.json` files** to `<target>/basic/scenes/` -- these are OBS
   scene collections. OBS discovers scene files from this directory on
   startup.
3. **Route directories** to `<target>/basic/profiles/` -- these are OBS
   profile folders (each containing `basic.ini`, `streamEncoder.json`,
   etc.). OBS discovers profiles from this directory on startup.
4. **Clean up** the temporary extraction directory.

If no `.zip` file is found, the sync falls back to a recursive directory
copy from the source folder to the target.

The extraction uses `archive/zip` with path traversal protection and a
50 MB per-file size limit (OBS scene files and profile configs can be
larger than typical editor settings).

#### Example

Given a zip containing:

```
My Stream Profile/
  basic.ini
  streamEncoder.json
My Scenes.json
```

After `gitmap-v27 install obs-settings`:

```
%APPDATA%\obs-studio\
  basic\
    profiles\
      My Stream Profile\
        basic.ini
        streamEncoder.json
    scenes\
      My Scenes.json
```

### Windows Terminal (`wt-settings`)

| Platform | Target                                                          |
|----------|-----------------------------------------------------------------|
| Windows  | `%LOCALAPPDATA%\Packages\Microsoft.WindowsTerminal_*\LocalState` |
| Others   | Not supported (Windows-only tool)                               |

The sync copies `settings.json` (and any additional files like theme
fragments) to the Windows Terminal LocalState directory. The target
directory is auto-discovered by scanning `%LOCALAPPDATA%\Packages\` for
a folder matching `Microsoft.WindowsTerminal_*`.

## Path Resolution

All settings sync commands use `resolveSettingsPath()` to locate the
source files. The search order is:

1. **Binary-relative**: `<binary-dir>/settings/<subfolder>`
2. **CWD-relative**: `settings/<subfolder>`
3. **Legacy fallback**: `<binary-dir>/data/<legacy-name>` and `data/<legacy-name>`

This ensures settings are found whether gitmap-v27 is run from the repo root,
from a deployed binary directory, or from an arbitrary working directory.

## Behavior

### Idempotent

Running the same sync command multiple times overwrites target files with
the bundled versions. No merge logic is applied -- the bundled settings
are treated as the canonical source.

### Cross-platform guards

- NPP and Windows Terminal settings sync exit early with a clear error on
  non-Windows platforms.
- VS Code and OBS settings sync work on Windows, macOS, and Linux.

### Error handling

- Missing source directory: prints the searched paths and exits.
- Missing `APPDATA` / `LOCALAPPDATA` (Windows): prints a clear error and exits.
- Individual file copy failures: logged per-file, does not abort the batch.
- Missing `code` CLI (VS Code extensions): skips extension install with a
  warning, settings files are still synced.
- Windows Terminal package not found: prints error if
  `Microsoft.WindowsTerminal_*` folder does not exist.

### Zip extraction

| Tool | Max file size | Traversal protection | Fallback          |
|------|--------------|----------------------|-------------------|
| NPP  | 10 MB        | Yes                  | Loose-file copy   |
| OBS  | 50 MB        | Yes                  | Recursive dir copy |

## Implementation

| File                               | Purpose                                    |
|------------------------------------|--------------------------------------------|
| `gitmap-v27/cmd/installnpp.go`        | NPP settings sync + `resolveSettingsPath`  |
| `gitmap-v27/cmd/installnppextract.go` | NPP zip extraction logic                   |
| `gitmap-v27/cmd/installvscode.go`     | VS Code settings and extension sync        |
| `gitmap-v27/cmd/installobs.go`        | OBS zip extraction + scene/profile routing |
| `gitmap-v27/cmd/installwt.go`         | Windows Terminal settings sync             |
| `gitmap-v27/cmd/install.go`           | Command routing (`executeInstall`)         |
| `gitmap-v27/constants/constants_install.go` | Tool names and messages              |

## Constraints

1. Settings source files must be committed to the repository under
   `settings/`. They are not downloaded at runtime.
2. No merge or diff -- sync always overwrites the target.
3. NPP and Windows Terminal are Windows-only. VS Code and OBS are
   cross-platform.
4. Extension installs require `code` CLI in PATH.
5. The `readme.txt` in each settings subfolder is never copied to the target.
6. NPP zip extraction enforces a 10 MB per-file limit; OBS enforces 50 MB.
7. OBS `.json` files are routed to `basic/scenes/`; directories are routed
   to `basic/profiles/`.

## Related

- `spec/01-app/94-install-script.md` -- installer specification
- `spec/02-app-issues/22-installer-path-not-active-after-install.md` -- PATH fixes
- `gitmap-v27/helptext/install.md` -- CLI help text
