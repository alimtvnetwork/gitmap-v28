# Install — Developer Tool Installer

## Overview

The `install` command automates installation of common developer tools
across platforms. On Windows, it leverages Chocolatey and Winget as
package managers. On Linux/macOS, it uses the native package manager
(apt, brew) or official installation scripts.

The Go implementation is derived from battle-tested PowerShell scripts
(to be provided). Each tool has a dedicated installer function with
pre-check, install, and post-install verification steps.

---

## Commands

### `gitmap-v28 install` (alias: `in`)

Install a developer tool by name.

```bash
gitmap-v28 install vscode
gitmap-v28 install node
gitmap-v28 install go
gitmap-v28 install git
```

---

## Supported Tools

| Tool | Windows (Chocolatey/Winget) | Linux (apt/script) | macOS (brew/script) |
|------|----------------------------|---------------------|---------------------|
| VS Code | `choco install vscode` / `winget install Microsoft.VisualStudioCode` | Official .deb/rpm | `brew install --cask visual-studio-code` |
| Node.js | `choco install nodejs` | `nvm` or NodeSource | `brew install node` |
| Yarn | `choco install yarn` | `npm install -g yarn` | `brew install yarn` |
| Bun | `choco install bun` | Official install script | `brew install oven-sh/bun/bun` |
| pnpm | `choco install pnpm` | `npm install -g pnpm` | `brew install pnpm` |
| Python | `choco install python` | `apt install python3` | `brew install python` |
| Go | `choco install golang` | Official tarball | `brew install go` |
| Git | `choco install git` | `apt install git` | `brew install git` |
| Git LFS | `choco install git-lfs` | `apt install git-lfs` | `brew install git-lfs` |
| GitHub CLI | `choco install gh` | Official repo | `brew install gh` |
| GitHub Desktop | Direct installer | N/A (AppImage) | `brew install --cask github` |
| C++ (MinGW) | `choco install mingw` | `apt install g++` | Xcode CLI tools |
| PHP | `choco install php` | `apt install php` | `brew install php` |
| PowerShell | `winget install Microsoft.PowerShell` | Official repo | `brew install --cask powershell` |
| Notepad++ (npp) | `choco install notepadplusplus` | N/A | N/A |
| Notepad++ Settings (npp-settings) | Settings sync only | Settings sync only | N/A |
| Notepad++ Install Only (install-npp) | `choco install notepadplusplus` | N/A | N/A |

### Notepad++ Variants

| Command | Shortcut | Installs Binary | Syncs Settings |
|---------|----------|-----------------|----------------|
| `npp` | NPP + Settings | ✅ | ✅ |
| `npp-settings` | NPP Settings | ❌ | ✅ |
| `install-npp` | Install NPP | ✅ | ❌ |

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--manager <name>` | auto-detect | Force package manager (choco, winget, apt, brew) |
| `--version <ver>` | latest | Install a specific version |
| `--verbose` | `false` | Show full installer output |
| `--dry-run` | `false` | Show install command without executing |
| `--check` | `false` | Only check if tool is installed |
| `--list` | `false` | List all supported tools |

---

## Installation Flow

For each tool:

```
1. Pre-check: Is the tool already installed? (which/where lookup)
2. If installed: Print version and skip (unless --force)
3. Detect package manager availability
4. Run install command via os/exec
5. Post-verify: Confirm binary is accessible
6. Print success/failure summary
```

---

## Package Manager Detection (Windows)

```
1. Check if `choco` is in PATH -> use Chocolatey
2. Check if `winget` is in PATH -> use Winget
3. Neither found -> offer to install Chocolatey first
```

## Package Manager Detection (Linux)

```
1. Check if `apt` exists -> Debian/Ubuntu
2. Check if `dnf` exists -> Fedora/RHEL
3. Check if `pacman` exists -> Arch
4. Fallback: use official install script
```

## Package Manager Detection (macOS)

```
1. Check if `brew` exists -> use Homebrew
2. Not found -> offer to install Homebrew first
```

---

## Post-Install Actions

Some tools require post-install configuration:

| Tool | Post-Install Action |
|------|---------------------|
| VS Code | Context menu fix (Windows registry), settings sync |
| Git | `git config --global core.longpaths true` |
| Git LFS | `git lfs install` |
| Go | Set GOPATH via `gitmap-v28 env set` |
| Node.js | Verify npm accessible |
| Notepad++ (npp) | Verify exe at expected path, sync settings to AppData |
| npp-settings | Sync settings to AppData only (no binary install) |
| install-npp | Verify exe at expected path, skip settings |

---

## Cross-Platform Notes

### Windows (Primary)
- Chocolatey preferred; Winget as fallback.
- Some tools require admin prompt (detected and warned).
- VS Code context menu fix modifies Windows registry.

### Linux
- apt-based commands for Debian/Ubuntu (primary).
- Official install scripts for tools not in package repos.
- No sudo assumption; prompt user if needed.

### macOS
- Homebrew preferred.
- Xcode CLI tools installed via `xcode-select --install`.
- Cask used for GUI applications.

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Package manager not found | Offer to install it |
| Install command fails | Show error, suggest manual install |
| Tool already installed | Print version, skip |
| No admin/sudo access | Warn and exit for tools that require it |
| Network unavailable | Exit with offline message |

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_install.go` | Tool names, package IDs, messages |
| `cmd/install.go` | Command routing, flag parsing, tool dispatch |
| `cmd/installtools.go` | Per-tool install functions |
| `cmd/installdetect.go` | Package manager and pre-check detection |
| `cmd/installverify.go` | Post-install verification |
| `cmd/installplatform_windows.go` | Windows-specific install logic |
| `cmd/installplatform_unix.go` | Linux/macOS-specific install logic |

---

## Constraints

- All install commands run via `os/exec` (no shell eval).
- Package manager detection is cached per session.
- All files under 200 lines, all functions 8-15 lines.
- No hardcoded paths; use `exec.LookPath` for detection.
- Positive logic only (no `!` or `!=`).

---

## Future Enhancements (Pending)

- REST API integration for remote install manifests.
- `gitmap-v28 install --from <url>` to fetch and execute an install plan.
- Version pinning and update tracking.

---

## Examples

```bash
# Install VS Code
gitmap-v28 install vscode

# Install Node.js with verbose output
gitmap-v28 install node --verbose

# Check if Go is installed
gitmap-v28 install go --check

# Preview install command without running
gitmap-v28 install python --dry-run

# Force Winget on Windows
gitmap-v28 install git --manager winget

# List all supported tools
gitmap-v28 install --list
```

---

## See Also

- env — manage environment variables and PATH
- setup — shell profile and completion
- doctor — verify development environment
