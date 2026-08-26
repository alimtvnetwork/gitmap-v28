# Install Script Generation

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Defines the reusable pattern for generating version-pinned install scripts (PowerShell and Bash) during a release pipeline. These scripts are included as GitHub Release assets, enabling one-liner installation for end users.

---

## Placeholder Strategy

Install scripts are written inline in the workflow YAML using placeholder tokens. After the heredoc is written to disk, `sed` replaces the tokens with resolved values:

| Placeholder | Replaced With | Example Value |
|-------------|--------------|---------------|
| `VERSION_PLACEHOLDER` | Resolved version from Git ref | `v1.2.0` |
| `REPO_PLACEHOLDER` | Repository path (owner/repo) | `<owner>/<repo>` |

```bash
sed -i "s|VERSION_PLACEHOLDER|$VERSION|g" dist/install.ps1
sed -i "s|REPO_PLACEHOLDER|$REPO|g" dist/install.ps1
chmod +x dist/install.sh
```

### Why Placeholders Instead of String Interpolation

- Heredocs with `<< 'EOF'` (single-quoted) prevent shell variable expansion, avoiding accidental injection
- `sed` replacement is explicit and auditable
- The same script source can be tested locally without variable substitution

---

## PowerShell Installer (`install.ps1`)

### Structure

```powershell
param(
    [string]$InstallDir = "",
    [string]$Arch = "",
    [switch]$NoPath
)

$ErrorActionPreference = "Stop"
$PinnedVersion = "VERSION_PLACEHOLDER"
$Repo = "REPO_PLACEHOLDER"
$BinaryName = "<binary>.exe"
```

### Required Functions

| Function | Purpose |
|----------|---------|
| `Write-Step` | Cyan status message |
| `Write-OK` | Green success message |
| `Write-Err` | Red error message |
| `Resolve-InstallDir` | Default: `$env:LOCALAPPDATA\<binary>` |
| `Resolve-Arch` | Auto-detect from `$env:PROCESSOR_ARCHITECTURE` (AMD64 → amd64, ARM64 → arm64) |
| `Get-Asset` | Download zip + checksums, verify SHA-256 |
| `Install-Binary` | Extract zip, rename-first upgrade, move binary |
| `Add-ToPath` | Add to user `PATH` via `[Environment]::SetEnvironmentVariable` |

### Checksum Verification Flow

```
1. Download <binary>-<version>-windows-<arch>.zip
2. Download checksums.txt
3. Find matching line in checksums.txt
4. Compare Get-FileHash SHA256 (lowercased) against expected
5. Fail on mismatch — delete temp files and exit 1
```

### Rename-First Upgrade

```powershell
$targetPath = Join-Path $installDir $BinaryName
if (Test-Path $targetPath) {
    $oldPath = "$targetPath.old"
    if (Test-Path $oldPath) { Remove-Item $oldPath -Force }
    Rename-Item $targetPath $oldPath -Force
}

# ... extract and install ...

Remove-Item "$targetPath.old" -Force -ErrorAction SilentlyContinue
```

### PATH Management

```powershell
$currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$parts = if ($currentUserPath) { $currentUserPath -split ";" } else { @() }
$hasDir = $parts | Where-Object { $_.Trim() -ieq $dir }
if (-not $hasDir) {
    $newPath = $currentUserPath.TrimEnd(";") + ";" + $dir
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
}
```

---

## Bash Installer (`install.sh`)

### Structure

```bash
#!/usr/bin/env bash
set -euo pipefail

PINNED_VERSION="VERSION_PLACEHOLDER"
REPO="REPO_PLACEHOLDER"
BINARY_NAME="<binary>"
INSTALL_DIR=""
ARCH=""
NO_PATH=false
```

### CLI Flags

| Flag | Purpose | Default |
|------|---------|---------|
| `--version` | Override pinned version | Embedded value |
| `--dir` | Custom install directory | `$HOME/.local/bin` |
| `--arch` | Force architecture | Auto-detect |
| `--no-path` | Skip PATH modification | `false` |

### Required Functions

| Function | Purpose |
|----------|---------|
| `detect_os` | `uname -s` → `linux` or `darwin` (reject others with install.ps1 suggestion) |
| `detect_arch` | `uname -m` → `amd64` or `arm64` |
| `download` | Try `curl -fsSL`, fall back to `wget -qO` |
| `verify_checksum` | `sha256sum` or `shasum -a 256` fallback |
| `add_to_path` | Shell-aware: `.bashrc`, `.zshrc`, or `config.fish` |

### Shell-Aware PATH Configuration

```bash
shell_name="$(basename "${SHELL:-/bin/bash}")"
case "$shell_name" in
    zsh)  rc_file="$HOME/.zshrc" ;;
    fish) rc_file="$HOME/.config/fish/config.fish" ;;
    *)    rc_file="$HOME/.bashrc" ;;
esac

if [[ "$shell_name" == "fish" ]]; then
    export_line="fish_add_path $dir"
else
    export_line="export PATH=\"$dir:\$PATH\""
fi
```

### Checksum Verification

```bash
if command -v sha256sum &>/dev/null; then
    actual="$(sha256sum "$file" | awk '{print $1}')"
else
    actual="$(shasum -a 256 "$file" | awk '{print $1}')"
fi
```

---

## Common Flow (Both Scripts)

```
1. Print banner with version and repo
2. Detect/resolve platform and architecture
3. Resolve install directory
4. Download binary archive from GitHub Releases
5. Download checksums.txt
6. Verify SHA-256 checksum (fail hard on mismatch)
7. Rename existing binary to .old (if upgrading)
8. Extract and install new binary
9. Clean up .old file and temp directory
10. Add to PATH (unless --no-path / -NoPath)
11. Run <binary> version to confirm installation
12. Print success message
```

---

## Version-Pinned vs. Generic Installers

| Type | Location | Version Source |
|------|----------|---------------|
| **Version-pinned** | GitHub Release asset (`dist/install.ps1`) | Hardcoded via `sed` replacement |
| **Generic** | Repository source (`scripts/install.ps1`) | `$Version` parameter or GitHub API latest release |

The release pipeline generates **version-pinned** scripts. Generic installers (that query the GitHub API for the latest release) live in the repository source and are referenced in the release body.

---

## Cross-References

- [Shared Conventions](./01-shared-conventions.md) — Version resolution, checksums
- [GitHub Release Standard](./02-github-release-standard.md) — Release body template
- [Go Binary Release Pipeline](./02-go-binary-deploy/02-release-pipeline.md) — Full release workflow
- [Self-Update Mechanism](./06-self-update-mechanism.md) — How CLIs use install scripts internally
- [Code Signing](./05-code-signing.md) — Signed binaries that install scripts download
- [Self-Update & App Update (Full Specs)](../14-update/00-overview.md) — Client-side update implementation details

---

*Install script generation — v3.2.0 — 2026-04-10*
