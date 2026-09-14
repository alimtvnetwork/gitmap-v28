---
name: os-fix-clean-profiles-clone-suite
description: Autonomously implement and verify OS IP auto-revert, OS fix registry, OS clean temp cache, OS user/group import-export, VSCode profiles management, and clone sequential table and filtering across GitMap.
---

# OS Fix, Clean, User, User-Group, VSCode Profiles & Clone Suite

## Core Capabilities

1. **OS IP Configuration & Revert**:
   - `gitmap os ip`: Displays active network interfaces, IP addresses, netmasks, gateways, DHCP status, and link state.
   - `gitmap os ip help`: Prints structured usage help for `os ip`.
   - `gitmap os ip set <ip> [gateway]`: Sets static IPv4 address with pre-change snapshot, validates connectivity to Google (`8.8.8.8`), and automatically reverts to original snapshot if ping fails.

2. **OS Fix Registry & Runner**:
   - `gitmap os fix add <name> <command>`: Register new system fix/script.
   - `gitmap os fix edit <name> <command>`: Edit existing registered fix.
   - `gitmap os fix import <file>`: Import fixes from JSON/YAML file.
   - `gitmap os fix export <name> [file]`: Export specific fix to file or stdout.
   - `gitmap os fix import-all <file>`: Bulk import fixes.
   - `gitmap os fix export-all [file]`: Bulk export all fixes.
   - `gitmap os fix run <name | command>`: Execute a registered fix by name or run new fix commands directly.
   - `gitmap os fix ls`: List registered fixes.

3. **OS Clean & Temp Cache**:
   - `gitmap os clean`: Clean temporary and cached directories across the host operating system.
   - `gitmap os clean temp`: Clean `%TEMP%` / `%TMP%` / `%LOCALAPPDATA%\Temp` on Windows, and `/tmp` / `/var/tmp` / `~/.cache` on Linux/macOS.
   - `gitmap os clear temp`: Alias for `os clean temp`.

4. **OS User & Group Suite**:
   - `gitmap os user add/edit/import/export/import-all/export-all`: Full user account lifecycle management with portable JSON/YAML schemas.
   - `gitmap os user-group add/edit/import/export/import-all/export-all` (alias `gitmap os group`): Portable group management with membership tracking.

5. **VSCode Profiles Management**:
   - `gitmap vscode profiles ls`: List all installed VSCode profiles (Default and named custom profiles) with configuration details in a formatted table.
   - `gitmap vscode profiles export <name> [file]`: Export a specific profile's settings, extensions, and keybindings.
   - `gitmap vscode profiles export-all [dir|file]`: Export all detected profiles.
   - `gitmap vscode profiles import <file>`: Import a profile configuration into VSCode.
   - `gitmap vscode profiles import-all <dir|file>`: Import multiple profiles.

6. **Clone Sequence IDs & Filtering**:
   - `gitmap clone -ls file.json` (also `--ls`, `-l`): List repositories in candidate file in a formatted terminal table with 1-based sequential integer IDs (`1`, `2`, ...).
   - `gitmap clone file.json --only 1,3`: Filter targets by sequential IDs (`1,3`), slug (`owner/repo`), or prefix patterns (`alim*`).
