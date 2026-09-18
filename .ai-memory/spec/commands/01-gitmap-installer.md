# Gitmap Installer System Specification

## Core Requirements

- A script installation creation, export, and import system.
- Cross-platform portability: Paths must be stored as forward slashes and relative to a root "work directory" so the system can transition seamlessly between Windows, Ubuntu, CentOS, Debian, Arch Linux, etc.
- Persistence: All installer information must be stored in the SQLite database.
- Reset safety: `gitmap reset` must NOT clear out the installer scripts unless the user explicitly provides a parameter to do so.

## CLI Commands

### Creation & Management

- `gitmap installer create "<name>"`: Prompts for a description if not provided. Must prompt for target OS (win, ubuntu, centos, debian, arch-linux, etc.) providing a clear way to detect and assign new OSes.
- `gitmap installer update "<name>"`: Auto version update.
- `gitmap installer update-win "<name>"`: Auto version update, prompts only for Windows version.
- `gitmap installer install-win "<name or slug>" "<description>"`: Installs the Windows version. Variations should exist for other OSes (ubuntu, unix, centos, etc.).

### Import & Export

- `gitmap installer export-all`: Exports a `.zip` file (default name `gitmap-export.zip`) containing JSON files of the installers to the current directory or given path.
- `gitmap installer export "<slug>"`: Exports the specific installer as a `.zip` (containing JSON) so it can be imported elsewhere.
- `gitmap installer import`: Without arguments, looks for `gitmap-export.zip` in the current directory and imports it.
- `gitmap installer import "<filepath or name as .zip>"`: Imports from a specific ZIP file, single JSON file, or URL. If the installer already exists, it imports and updates to the latest one.

### Global Import & Export (Full Gitmap State)

- `gitmap export-all [<.zip file>]`: Exports the entire Gitmap state (installers, macros, git repos, etc.) while preserving the relative folder structure.
- `gitmap import-all [<.zip file>]`: Imports the full exported structure from another machine.
- `gitmap export-only installer,macro,repos [<.zip file>]`: Selectively exports specific domains.
- `gitmap export-only installer [<.zip file>]`: Only exports installers (identical to `installer export-all`).

### Version Control

- `gitmap reset installer --all`: Resets all installers.
- `gitmap reset installer "<slug>"`: Resets a specific installer.
- `gitmap installer undo-version "<slug>"`: Reverts to the previous version.
- `gitmap installer redo-version "<slug>"`: Reverts to the redone version.
- `gitmap installer revert-version "<slug>" v<semantic version>`: Reverts to a specific version.
- `gitmap installer revert-version "<slug>" <semantic version>`: Alternative syntax for revert.

### Listing & UI

- `gitmap installer ls`: Lists all created installers in a table.
- `gitmap installer ls win`: Lists only the Windows installers.
- Interactive UI/Help: Running `gitmap installer` should show detailed help text, explanations, and examples.
- UI Dashboards: The system should provide HD UI versions of this, showcasing examples of capabilities (e.g., Composer example).
