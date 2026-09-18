# Install Scripts

**Version:** 3.2.0
**Updated:** 2026-04-17

---

## Purpose

Generate cross-platform, one-liner installer scripts that download the correct binary for the user's platform, verify its checksum, and install it to a standard location. These scripts are published as release assets.

---

## Overview

Two scripts are generated per release:

| Script | Platform | Invocation |
|--------|----------|------------|
| `install.ps1` | Windows (PowerShell 5.1+) | `irm <url>/install.ps1 \| iex` |
| `install.sh` | Linux / macOS (Bash 4+) | `curl -fsSL <url>/install.sh \| bash` |

Both scripts are **version-pinned** at generation time — they download a specific release version, not "latest".

---

## PowerShell Installer (`install.ps1`)

### Parameters

```powershell
param(
  [string]$Version = "",     # Pin to specific version (default: latest)
  [string]$InstallDir = "",  # Override install location
  [string]$Arch = "",        # Override architecture detection
  [switch]$NoPath            # Skip PATH modification
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

When running via `irm | iex`, suppress the download progress bar for cleaner output:

```powershell
$ProgressPreference = "SilentlyContinue"
```

### Post-Install Summary

```
 ✓ <binary> v1.2.0 installed successfully

 Binary:  /Users/<user>\AppData\Local\<binary>\<binary>.exe
 Version: v1.2.0

 PATH:    Added /Users/<user>\AppData\Local\<binary> to user PATH
 Note:    Restart your terminal for PATH changes to take effect
```

---

## Bash Installer (`install.sh`)

### Parameters

```bash
--version <ver>  # Pin to specific version
--dir <path>     # Override install location
--arch <arch>    # Override architecture detection
--no-path        # Skip PATH modification
```

### Shell Compatibility Guard

The script must work when piped via `curl | sh`. Since `sh` on many systems is `dash` (not `bash`), add a self-re-exec guard:

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
    zsh) profile="$HOME/.zshrc" ;;
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

 PATH:    Added /home/<user>/.local/bin to ~/.zshrc
 Note:    Run 'source ~/.zshrc' or restart your terminal
```

---

## Version Pinning via Placeholder Substitution

During CI, the install scripts are generated from templates using placeholder substitution:

```bash
sed -i "s/VERSION_PLACEHOLDER/$VERSION/g" dist/install.ps1
sed -i "s/REPO_PLACEHOLDER/$REPO/g" dist/install.ps1
```

Placeholders:

| Placeholder | Replaced With |
|-------------|---------------|
| `VERSION_PLACEHOLDER` | The release version (e.g., `1.2.0`) |
| `REPO_PLACEHOLDER` | The repository identifier (e.g., `org/repo`) |

This ensures each release's install scripts always download that specific version, not "latest".

---

## Spec Repository Cloning

When the installer also clones a **spec / coding-guidelines repository**
(not just a binary), it MUST display the cloned version after the
`git clone` completes. The version is read from a `version.json` file
at the **repository root**.

### `version.json` (repository root)

Every spec repo MUST publish a `version.json` at its root:

```json
{
  "version": "1.21.0",
  "updated": "2026-04-17",
  "name": "coding-guidelines",
  "description": "Cross-language coding standards, error handling, CI/CD, and self-update specifications."
}
```

| Field | Required | Source of truth |
|-------|----------|-----------------|
| `version` | ✅ | Mirrors `package.json#version` (auto-synced) |
| `updated` | ✅ | ISO date (UTC+8) of last sync |
| `name` | ✅ | Short slug — appears in installer output |
| `description` | ⚠️ optional | One-line human summary |

`version.json` is generated by `scripts/sync-version.mjs` and MUST be
committed. The `prebuild` npm hook keeps it in lockstep with
`package.json`, so it never drifts.

### Clone Step (PowerShell)

```powershell
function Clone-SpecRepo {
    param(
        [string]$RepoUrl,
        [string]$TargetDir
    )

    Write-Info "Cloning $RepoUrl into $TargetDir"
    git clone --depth 1 $RepoUrl $TargetDir
    if ($LASTEXITCODE -ne 0) { throw "git clone failed (exit $LASTEXITCODE)" }

    $verPath = Join-Path $TargetDir "version.json"
    if (-not (Test-Path $verPath)) {
        Write-Warn "version.json missing at repo root — cannot report version"
        return
    }

    $ver = Get-Content $verPath | ConvertFrom-Json
    Write-Success "Cloned $($ver.name) v$($ver.version) (updated $($ver.updated))"
}
```

### Clone Step (Bash)

```bash
clone_spec_repo() {
    local repo_url="$1"
    local target_dir="$2"

    write_info "Cloning $repo_url into $target_dir"
    git clone --depth 1 "$repo_url" "$target_dir" || {
        write_fail "git clone failed"; exit 1;
    }

    local ver_path="$target_dir/version.json"
    if [[ ! -f "$ver_path" ]]; then
        write_warn "version.json missing at repo root — cannot report version"
        return
    fi

    # Read with python3 (always present), fall back to jq, then grep.
    local name version updated
    if command -v python3 &>/dev/null; then
        name=$(python3 -c "import json;print(json.load(open('$ver_path'))['name'])")
        version=$(python3 -c "import json;print(json.load(open('$ver_path'))['version'])")
        updated=$(python3 -c "import json;print(json.load(open('$ver_path'))['updated'])")
    elif command -v jq &>/dev/null; then
        name=$(jq -r .name "$ver_path")
        version=$(jq -r .version "$ver_path")
        updated=$(jq -r .updated "$ver_path")
    else
        # Last-resort grep — only safe because version.json is a flat object.
        name=$(grep -oE '"name"[^,}]*' "$ver_path" | sed -E 's/.*"([^"]+)"$/\1/')
        version=$(grep -oE '"version"[^,}]*' "$ver_path" | sed -E 's/.*"([^"]+)"$/\1/')
        updated=$(grep -oE '"updated"[^,}]*' "$ver_path" | sed -E 's/.*"([^"]+)"$/\1/')
    fi

    write_success "Cloned $name v$version (updated $updated)"
}
```

### Output Example

```
  -> Cloning https://github.com/<org>/coding-guidelines into ./coding-guidelines
  OK Cloned coding-guidelines v1.21.0 (updated 2026-04-17)
```

### Constraints

- `version.json` MUST live at the **repository root**, not inside any subfolder.
- The installer MUST **read after clone**, not probe GitHub before — keeps the
  flow offline-friendly and reports exactly what landed on disk.
- If `version.json` is missing, warn but do NOT fail the clone — the user
  still gets the spec content; only the version banner is missing.
- The `version` field MUST be a single dotted string (no `v` prefix, no
  build suffix). Installer prepends `v` for display.
- `version.json` is **never hand-edited** — always generated from
  `package.json` via `scripts/sync-version.mjs`. Hand-edits will be
  overwritten by the next `prebuild`.

---

## Constraints

- Scripts must work without any pre-installed dependencies beyond the shell itself (PowerShell 5.1+ or Bash 4+).
- Checksum verification is mandatory — never skip it.
- PATH modifications must be idempotent (don't add duplicates).
- Always print a visible post-install summary with binary path and version.
- Suppress progress bars when running in piped/non-interactive mode.
- Bash scripts must handle `curl | sh` execution via the self-re-exec guard.
- When cloning a spec repo, ALWAYS report the version from `version.json`.

---

## Cross-References

- [Checksums & Verification](./15-checksums-verification.md) — SHA-256 verification patterns
- [Release Assets](./14-release-assets.md) — Asset naming conventions
- [Release Pipeline](./18-release-pipeline.md) — Where install scripts are generated
- [CI/CD Install Script Generation](../12-cicd-pipeline-workflows/07-install-script-generation.md) — Server-side generation patterns

---

*Install scripts — v3.2.0 — 2026-04-17*
