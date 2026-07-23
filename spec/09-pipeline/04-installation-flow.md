# Installation Flow

## Overview

This document describes the complete installation experience for a CLI tool — from one-liner install commands to terminal output, version display, and post-install verification. Any AI or engineer implementing a new CLI tool should follow these patterns to provide a polished, production-ready installation experience.

---

## One-Liner Install Commands

### PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/<owner>/<repo>/main/install.ps1 | iex
```

### Bash (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/<owner>/<repo>/main/install.sh | bash
```

### Version-Pinned Variants

Always provide version-pinned alternatives alongside the "latest" install:

```powershell
# PowerShell — specific version
irm https://github.com/<owner>/<repo>/releases/download/v1.2.3/install.ps1 | iex
```

```bash
# Bash — specific version
curl -fsSL https://github.com/<owner>/<repo>/releases/download/v1.2.3/install.sh | bash
```

---

## Install Script Structure

### `install.ps1` (PowerShell)

The script must follow this exact flow:

```
1. Print banner with tool name and version
2. Detect architecture (amd64 / arm64)
3. Construct download URL from version + arch
4. Download binary archive (.zip) + checksums.txt
5. Verify SHA-256 checksum (fail on mismatch with clear error)
6. If upgrading: rename existing binary to .old (rename-first strategy)
7. Extract binary from archive to install directory
8. Clean up .old file on success
9. Register PATH:
   a. Windows Registry (User PATH)
   b. PowerShell profile ($PROFILE)
   c. Git Bash profiles (~/.bash_profile, ~/.bashrc)
10. Print post-install summary with activation commands
11. Print installed version
```

### `install.sh` (Bash)

```
1. Print banner with tool name and version
2. Detect OS (linux / darwin) and architecture (amd64 / arm64)
3. Select download tool (curl preferred, wget fallback)
4. Construct download URL from version + OS + arch
5. Download binary archive (.tar.gz) + checksums.txt
6. Verify SHA-256 checksum (sha256sum or shasum -a 256 fallback)
7. If upgrading: rename existing binary to .old
8. Extract binary to install directory
9. Clean up .old file on success
10. Register PATH in shell profiles:
    a. Bash → ~/.bashrc (and ~/.bash_profile if exists)
    b. Zsh → ~/.zshrc
    c. Fish → ~/.config/fish/config.fish
11. Use marker comments (e.g., "# <tool>-path") for idempotent PATH entries
12. Print post-install summary with activation commands
13. Print installed version
```

### CLI Flags for `install.sh`

| Flag | Description | Default |
|------|-------------|---------|
| `--version <ver>` | Install a specific version | Latest |
| `--dir <path>` | Custom install directory | `~/.local/bin` or platform default |
| `--arch <arch>` | Override architecture detection | Auto-detected |
| `--no-path` | Skip PATH registration | Register PATH |

---

## Terminal Output — Installation Sample

### Fresh Install (PowerShell)

```
PS> irm https://raw.githubusercontent.com/owner/repo/main/install.ps1 | iex

  ╔══════════════════════════════════════╗
  ║     <tool> Installer v1.2.3         ║
  ╚══════════════════════════════════════╝

  Platform:       windows/amd64
  Install path:   C:\Users\Admin\AppData\Local\<tool>\
  Downloading:    <tool>-v1.2.3-windows-amd64.zip

  Verifying checksum... OK
  Extracting binary...  OK
  Registering PATH...

    [+] Windows Registry (User PATH)
    [+] PowerShell profile
    [+] Git Bash profile

  ============================================
  Installation complete!

  To start using <tool> right now, run:

    $env:Path = [System.Environment]::GetEnvironmentVariable('Path','User') + ';' + $env:Path

  Or restart your terminal.
  ============================================

  <tool> v1.2.3
```

### Fresh Install (Bash)

```
$ curl -fsSL https://raw.githubusercontent.com/owner/repo/main/install.sh | bash

  <tool> Installer v1.2.3

  Platform:       linux/amd64
  Install path:   /home/user/.local/bin
  Downloading:    <tool>-v1.2.3-linux-amd64.tar.gz

  Verifying checksum... OK
  Extracting binary...  OK
  Registering PATH...

    [+] ~/.bashrc
    [+] ~/.zshrc

  ============================================
  Installation complete!

  To start using <tool> right now, run:

    source ~/.bashrc

  Or restart your terminal.
  ============================================

  <tool> v1.2.3
```

### Upgrade (PowerShell)

```
PS> irm https://raw.githubusercontent.com/owner/repo/main/install.ps1 | iex

  <tool> Installer v1.3.0

  Platform:       windows/amd64
  Install path:   C:\Users\Admin\AppData\Local\<tool>\
  Existing:       v1.2.3 -> renaming to <tool>.exe.old
  Downloading:    <tool>-v1.3.0-windows-amd64.zip

  Verifying checksum... OK
  Extracting binary...  OK
  Cleaning up old binary... OK
  PATH already registered.

  ============================================
  Upgrade complete! v1.2.3 -> v1.3.0

  To start using the new version, run:

    $env:Path = [System.Environment]::GetEnvironmentVariable('Path','User') + ';' + $env:Path

  Or restart your terminal.
  ============================================

  <tool> v1.3.0
```

---

## Version Display at Install Time

Both install scripts MUST print the tool version at two points:

1. **Banner** — at the very top, so the user immediately knows which version is being installed
2. **Final line** — after installation completes, by running `<tool> version` (or `<tool> --version`) to confirm the installed binary works

If the binary cannot execute (e.g., wrong architecture), the final version check should print a clear error:

```
  WARNING: Could not verify installed version.
  Try running: <tool> version
```

---

## Post-Install Verification

The install script should suggest verification commands:

```
  Verify your installation:

    <tool> version        # Print version
    <tool> doctor         # Run health checks (if available)
    <tool> help           # Show available commands
```

---

## Uninstall

### PowerShell

```powershell
irm https://raw.githubusercontent.com/owner/repo/main/install.ps1 | iex -Uninstall
```

Or the CLI itself:

```
<tool> uninstall
```

### Uninstall Flow

```
1. Locate install directory
2. Remove binary
3. Remove PATH entries:
   a. Windows Registry
   b. PowerShell profile (remove marker lines)
   c. Git Bash profiles (remove marker lines)
4. Remove data directory (with --purge flag only)
5. Print summary
```

### Terminal Output — Uninstall

```
  Uninstalling <tool>...

  Removing binary:   C:\Users\Admin\AppData\Local\<tool>\<tool>.exe
  Cleaning PATH:

    [-] Windows Registry (User PATH)
    [-] PowerShell profile
    [-] Git Bash profile

  <tool> has been uninstalled.
  Data directory preserved at: C:\Users\Admin\AppData\Local\<tool>\data\
  To remove all data, run: <tool> uninstall --purge
```

---

## Constraints

- Scripts must use ASCII-only output (no Unicode) unless saved with UTF-8 BOM
- All download URLs must be constructed from version + platform + arch — never hardcoded
- Checksum verification is mandatory — never skip
- Rename-first upgrade strategy — never delete-then-write (avoids Windows file locks)
- PATH registration must be idempotent (use marker comments, check before adding)
- Install scripts must work on PowerShell 5.1+ (no modern syntax like `??` or `Join-Path` with 3+ args)
- Bash scripts must work on bash 3.2+ (macOS ships old bash)
