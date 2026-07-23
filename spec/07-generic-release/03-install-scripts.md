# 03 — Install Scripts

## Purpose

Generate cross-platform, one-liner installer scripts that download the
correct binary for the user's platform, verify its checksum, and install
it to a standard location. These scripts are published as release assets.

---

## Overview

Two scripts are generated per release:

| Script | Platform | Invocation |
|--------|----------|------------|
| `install.ps1` | Windows (PowerShell 5.1+) | `irm <url>/install.ps1 \| iex` |
| `install.sh` | Linux / macOS (Bash 4+) | `curl -fsSL <url>/install.sh \| bash` |

Both scripts are **version-pinned** at generation time — they download
a specific release version, not "latest".

---

## PowerShell Installer (`install.ps1`)

### Parameters

```powershell
param(
    [string]$Version    = "",      # Pin to specific version (default: latest)
    [string]$InstallDir = "",      # Override install location
    [string]$Arch       = "",      # Override architecture detection
    [switch]$NoPath                # Skip PATH modification
)
```

### Pipeline

```
1. Resolve version (param or GitHub API → latest tag)
2. Detect architecture (amd64 or arm64)
3. Download archive (.zip) from release assets
4. Download checksums.txt
5. Verify SHA-256 hash
6. Extract binary to install directory
7. Add install directory to user PATH (unless -NoPath)
8. Verify installation: <binary> version
9. Print summary
```

### Architecture Detection

```powershell
function Resolve-Arch {
    param([string]$arch)
    if ($arch -ne "") { return $arch }

    $cpuArch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    if ($cpuArch -eq "Arm64") { return "arm64" }

    return "amd64"
}
```

### Checksum Verification

```powershell
$expectedHash = (Get-Content checksums.txt |
    Where-Object { $_ -match $archiveName } |
    ForEach-Object { ($_ -split '\s+')[0] })

$actualHash = (Get-FileHash $archivePath -Algorithm SHA256).Hash

if ($actualHash -ne $expectedHash) {
    Write-Error "Checksum mismatch! Expected: $expectedHash, Got: $actualHash"
    exit 1
}
```

### Default Install Location

```powershell
$defaultDir = Join-Path $env:LOCALAPPDATA "<binary>"
```

Falls back to `$HOME\.<binary>` if `LOCALAPPDATA` is not available.

### PATH Registration

```powershell
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable(
        "Path", "$currentPath;$installDir", "User"
    )
}
```

After modifying PATH, broadcast the change so other processes pick it up:

```powershell
# SendMessageTimeout to notify Explorer of environment change
Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
    [DllImport("user32.dll", SetLastError = true)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
        uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
$HWND_BROADCAST = [IntPtr]0xffff
$WM_SETTINGCHANGE = 0x1a
$result = [UIntPtr]::Zero
[Win32.NativeMethods]::SendMessageTimeout(
    $HWND_BROADCAST, $WM_SETTINGCHANGE,
    [UIntPtr]::Zero, "Environment", 2, 5000, [ref]$result
) | Out-Null
```

### Progress Bar Suppression

When running via `irm | iex`, suppress the download progress bar for
cleaner output:

```powershell
$ProgressPreference = "SilentlyContinue"
```

### Post-Install Summary

The script must print a visible summary after installation:

```
  ✓ <binary> v1.2.0 installed successfully

  Binary:  C:\Users\<user>\AppData\Local\<binary>\<binary>.exe
  Version: v1.2.0

  PATH: Added C:\Users\<user>\AppData\Local\<binary> to user PATH
  Note: Restart your terminal for PATH changes to take effect
```

---

## Bash Installer (`install.sh`)

### Parameters

```bash
--version <ver>    # Pin to specific version
--dir <path>       # Override install location
--arch <arch>      # Override architecture detection
--no-path          # Skip PATH modification
```

### Shell Compatibility Guard

The script must work when piped via `curl | sh`. Since `sh` on many
systems is `dash` (not `bash`), add a self-re-exec guard:

```bash
#!/usr/bin/env bash

# Re-exec under bash if running under a different shell
if [ -z "$BASH_VERSION" ]; then
    if command -v bash >/dev/null 2>&1; then
        exec bash -s -- "$@" < /dev/stdin
    fi
    echo "Error: bash is required but not found" >&2
    exit 1
fi

set -euo pipefail
```

### OS and Architecture Detection

```bash
detect_os() {
    local os
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$os" in
        linux*)  echo "linux" ;;
        darwin*) echo "darwin" ;;
        mingw*|msys*|cygwin*) echo "windows" ;;
        *)
            echo "Unsupported OS: $os" >&2
            exit 1
            ;;
    esac
}

detect_arch() {
    local arch
    arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)
            echo "Unsupported architecture: $arch" >&2
            exit 1
            ;;
    esac
}
```

### Checksum Verification

```bash
verify_checksum() {
    local archive="$1"
    local checksums_file="$2"
    local archive_name
    archive_name="$(basename "$archive")"

    local expected
    expected="$(grep "$archive_name" "$checksums_file" | awk '{print $1}')"

    local actual
    actual="$(sha256sum "$archive" | awk '{print $1}')"

    if [[ "$actual" != "$expected" ]]; then
        echo "Checksum mismatch!" >&2
        echo "  Expected: $expected" >&2
        echo "  Got:      $actual" >&2
        exit 1
    fi
}
```

### Default Install Location

```bash
if [[ -d "$HOME/.local/bin" ]]; then
    default_dir="$HOME/.local/bin"
elif [[ -d "/usr/local/bin" ]] && [[ -w "/usr/local/bin" ]]; then
    default_dir="/usr/local/bin"
else
    default_dir="$HOME/bin"
fi
```

### PATH Registration

Detect the user's active shell and append to the correct profile:

```bash
add_to_path() {
    local dir="$1"
    local shell_name
    shell_name="$(basename "$SHELL")"

    local profile=""
    case "$shell_name" in
        zsh)  profile="$HOME/.zshrc" ;;
        bash)
            if [[ -f "$HOME/.bash_profile" ]]; then
                profile="$HOME/.bash_profile"
            else
                profile="$HOME/.bashrc"
            fi
            ;;
        fish) profile="$HOME/.config/fish/config.fish" ;;
    esac

    if [[ -z "$profile" ]]; then
        echo "  !! Could not detect shell profile. Add manually:"
        echo "     export PATH=\"$dir:\$PATH\""
        return
    fi

    # Check if already in profile
    if grep -q "$dir" "$profile" 2>/dev/null; then
        return
    fi

    if [[ "$shell_name" == "fish" ]]; then
        echo "set -gx PATH \"$dir\" \$PATH" >> "$profile"
    else
        echo "export PATH=\"$dir:\$PATH\"" >> "$profile"
    fi

    echo "  Added $dir to $profile"
    echo "  Run: source $profile"
}
```

### Cleanup

Use a trap to clean up temporary files on exit:

```bash
cleanup() {
    [[ -n "${TMP_DIR:-}" ]] && rm -rf "$TMP_DIR"
}
trap cleanup EXIT
```

### Post-Install Summary

```
  ✓ <binary> v1.2.0 installed successfully

  Binary:  /home/<user>/.local/bin/<binary>
  Version: v1.2.0

  PATH: Added /home/<user>/.local/bin to ~/.zshrc
  Note: Run 'source ~/.zshrc' or restart your terminal
```

---

## Version Pinning via Placeholder Substitution

During CI, the install scripts are generated from templates using
placeholder substitution:

```bash
sed -i "s/VERSION_PLACEHOLDER/$VERSION/g" dist/install.ps1
sed -i "s/REPO_PLACEHOLDER/$REPO/g" dist/install.ps1
```

Placeholders:

| Placeholder | Replaced With |
|-------------|---------------|
| `VERSION_PLACEHOLDER` | The release version (e.g., `1.2.0`) |
| `REPO_PLACEHOLDER` | The repository identifier (e.g., `org/repo`) |

This ensures each release's install scripts always download that
specific version, not "latest".

---

## Constraints

- Scripts must work without any pre-installed dependencies beyond
  the shell itself (PowerShell 5.1+ or Bash 4+).
- Checksum verification is mandatory — never skip it.
- PATH modifications must be idempotent (don't add duplicates).
- Always print a visible post-install summary with binary path and version.
- Suppress progress bars when running in piped/non-interactive mode.
- Bash scripts must handle `curl | sh` execution via the self-re-exec guard.

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Version-pinned install scripts in release pipeline |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
