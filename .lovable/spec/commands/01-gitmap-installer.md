# Gitmap Installer Commands

## Goal
A comprehensive installer management system in Gitmap allowing users to create, update, export, import, and version-control installation scripts across various OS platforms (win, ubuntu, centos, debian, arc-linux). Information will be securely tracked via SQLite database, and the entire system must use a consistent relative work directory strategy.

## Context
Gitmap should always work from a work/root directory. After setting the root directory, all other paths should remain relative so that switching OS (e.g. Windows to Linux) is completely seamless.

## Command Specifications

### Create & Update
- `gitmap installer create "<name>"`: Prompts for description if not given. Prompts to define script behavior for `win`, `ubuntu`, `centos`, `debian`, `arc-linux`. Cleanly detects OS and maps commands.
- `gitmap installer update "<name>"`: Auto version updates it.
- `gitmap installer update-<os> "<name>"`: (e.g. `update-win`, `update-ubuntu`). Auto version updates it, but only prompts for the specific OS version.
- `gitmap installer install-<os> "<name|slug>" "<description>"`: Installs the specific OS version of the installer.

### Export & Import
- `gitmap installer export-all`: Exports all installers to a `.zip` file in the current directory (or given path). Contains JSON profiles of the installers.
- `gitmap installer export "<name>"`: Exports a specific installer as a JSON/zip.
- `gitmap installer import`: Automatically searches for `gitmap-export.zip` in the current directory by default and imports the JSON payloads.
- `gitmap installer import "<filepath|name>"`: Imports from a specified zip file, JSON file, or URL. If the installer exists, it imports as the latest version.

### Reset & Version Control
- `gitmap reset installer --all`: Resets all installers.
- `gitmap reset installer "<name>"`: Resets a specific installer. Note: `gitmap reset` (global) should NOT clear installer scripts unless explicitly requested by the user.
- `gitmap installer undo-version "<slug>"`: Reverts to previous version.
- `gitmap installer redo-version "<slug>"`: Reverts to redo version.
- `gitmap installer revert-version "<slug>" v<semantic-version>`: Reverts to a specific version string.
- `gitmap installer revert-version "<slug>" <semantic-version>`: Same as above.

### Listing
- `gitmap installer ls`: Lists all created installers in a table format.
- `gitmap installer ls <os>`: Lists installers that have a script available for the specified OS (e.g. `ls win`).

### Global Export/Import
- `gitmap export-all [<.zip file>]`: Exports the entire Gitmap environment (git repos, installers, macros) preserving the relative directory structure. Default file is `gitmap-export.zip`.
- `gitmap import-all [<.zip file>]`: Imports the entire environment. Default file is `gitmap-export.zip`.
- `gitmap export-only installer,macro,repos [<.zip file>]`: Selectively exports components.
- `gitmap export-only installer [<.zip file>]`: Same as `export-all` but strictly limited to installers.

## UI and Help Text Requirements
- Executing `gitmap installer` without arguments must show detailed help text with explanations.
- UI integration with the `HD` (Help Desk / Help UI) command is required to showcase examples of what these commands can do (including a Composer example) and how they can assist the user.
