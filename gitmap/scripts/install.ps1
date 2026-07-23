<#
.SYNOPSIS
    One-liner installer for gitmap CLI.

.DESCRIPTION
    Downloads the latest gitmap release from GitHub, verifies checksums,
    extracts to a local directory, and adds it to PATH.

.PARAMETER Version
    Install a specific version (e.g. v2.48.0). Default: latest.

.PARAMETER InstallDir
    Target directory. Default: $env:LOCALAPPDATA\gitmap

.PARAMETER NoPath
    Skip adding to PATH.

.PARAMETER Arch
    Force architecture (amd64, arm64). Default: auto-detect.

.PARAMETER DryRun
    Resolve the asset URL + filename, run the pre-flight HEAD probe to
    confirm the release asset exists, print a machine-parseable
    `dryrun.*=` report, and exit 0 without downloading or installing.
    Used by CI to validate the install.ps1 / release-pipeline naming
    contract.

.EXAMPLE
    irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.ps1 | iex

.EXAMPLE
    & ./install.ps1 -Version v2.48.0

.EXAMPLE
    & ./install.ps1 -Version v4.1.0 -DryRun -NoDiscovery

.NOTES
    Repository: https://github.com/alimtvnetwork/gitmap-v28
#>

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [string]$Arch = "",
    [switch]$NoPath,
    [switch]$Uninstall,
    [switch]$NoDiscovery,
    [int]$ProbeCeiling = 30,
    # --- Uninstall safety knobs ---
    # -Force: skip the "is this really a gitmap install?" guard AND
    #         skip the keep-data prompt (defaults to keep).
    # -KeepData / -PurgeData: explicit data-folder choice; when both
    #         are absent the user is prompted interactively.
    [switch]$Force,
    [switch]$KeepData,
    [switch]$PurgeData,
    # -DryRun: resolve the asset name + URL, run the pre-flight HEAD
    # probe to confirm the release asset exists, print a structured
    # report, and exit 0 — without downloading, extracting, or
    # touching the install dir / PATH. Used by CI to validate that
    # the install.ps1 + release-pipeline asset-naming contract is
    # honored before users hit it. Spec: spec/07-generic-release/
    # 09-generic-install-script-behavior.md §6.
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Force UTF-8 console I/O so non-ASCII glyphs (warning signs, arrows, box-
# drawing chars) and any future localized strings render correctly instead
# of cp437/Windows-1252 mojibake. Also defensive for `irm | iex`, where the
# active console code page may still be 437/1252 even though the script
# body is decoded as UTF-8 by Invoke-RestMethod.
#
# Order:
#   1. Flip the OS console code page to UTF-8 (65001) via Win32 so conhost
#      stops transcoding our bytes.
#   2. Set .NET stdout/stderr/stdin encodings to UTF-8.
#   3. Set PowerShell pipeline encoding so external processes (curl, tar,
#      gitmap.exe) also receive UTF-8.
#   4. Disable PSStyle ANSI rendering on PS 7+ so legacy hosts do not echo
#      raw escape sequences as literal text.
try {
    Add-Type -Namespace GitmapInstaller -Name NativeConsole -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("kernel32.dll")]
public static extern bool SetConsoleOutputCP(uint codePageID);
[System.Runtime.InteropServices.DllImport("kernel32.dll")]
public static extern bool SetConsoleCP(uint codePageID);
'@ -ErrorAction SilentlyContinue
    [GitmapInstaller.NativeConsole]::SetConsoleOutputCP(65001) | Out-Null
    [GitmapInstaller.NativeConsole]::SetConsoleCP(65001) | Out-Null
} catch {}
try {
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    [Console]::InputEncoding  = [System.Text.Encoding]::UTF8
    $OutputEncoding           = [System.Text.Encoding]::UTF8
    if ($PSStyle) { $PSStyle.OutputRendering = 'PlainText' }
} catch {}

$Repo = "alimtvnetwork/gitmap-v28"
$BinaryName = "gitmap.exe"
$BinaryAlias = "gm.exe"
$InstallerVersion = "1.0.0"

class InstallerFailure : System.Exception {
    [int]$ExitCode

    InstallerFailure([string]$message, [int]$exitCode) : base($message) {
        $this.ExitCode = $exitCode
    }
}

# ---------------------------------------------------------------------------
# Versioned repo discovery (spec/01-app/95-installer-script-find-latest-repo.md)
# ---------------------------------------------------------------------------

function Split-RepoSuffix([string]$repo) {
    if ($repo -match '^([^/]+)/(.+)-v(\d+)$') {
        return @{ Owner = $Matches[1]; Stem = $Matches[2]; N = [int]$Matches[3] }
    }
    return $null
}

function Test-RepoExists([string]$url) {
    try {
        $resp = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 5 `
            -UseBasicParsing -ErrorAction Stop
        return ($resp.StatusCode -eq 200)
    } catch {
        return $false
    }
}

function Resolve-EffectiveRepo([string]$repo, [int]$ceiling) {
    $parts = Split-RepoSuffix $repo
    if ($null -eq $parts) {
        Write-Host "  [discovery] no -v<N> suffix on '$repo'; installing baseline as-is"
        return $repo
    }

    $owner = $parts.Owner; $stem = $parts.Stem; $baseline = $parts.N
    $effective = $baseline

    Write-Host "  [discovery] baseline: $owner/$stem-v$baseline"
    Write-Host "  [discovery] probe ceiling: $ceiling"

    for ($m = $baseline + 1; $m -le $ceiling; $m++) {
        $url = "https://github.com/$owner/$stem-v$m"
        if (Test-RepoExists $url) {
            Write-Host "  [discovery] HEAD $url ... HIT"
            $effective = $m
        } else {
            Write-Host "  [discovery] HEAD $url ... MISS (fail-fast)"
            break
        }
    }

    if ($effective -eq $baseline) {
        Write-Host "  [discovery] no higher version found; using baseline -v$baseline"
        return $repo
    }

    Write-Host "  [discovery] effective: $owner/$stem-v$effective (was -v$baseline)"
    return "$owner/$stem-v$effective"
}

function Invoke-DelegatedFullInstaller([string]$effectiveRepo) {
    $delegatedUrl = "https://raw.githubusercontent.com/$effectiveRepo/main/gitmap/scripts/install.ps1"
    Write-Host "  [discovery] delegating to $delegatedUrl"

    $env:INSTALLER_DELEGATED = "1"
    try {
        $script = (Invoke-WebRequest -Uri $delegatedUrl -UseBasicParsing -TimeoutSec 15).Content
    } catch {
        Write-Host "  [discovery] [WARN] could not fetch delegated installer: $_" -ForegroundColor Yellow
        Write-Host "  [discovery] falling back to baseline installer" -ForegroundColor Yellow
        Remove-Item Env:INSTALLER_DELEGATED -ErrorAction SilentlyContinue
        return $false
    }

    $block = [ScriptBlock]::Create($script)
    $passArgs = @{ ProbeCeiling = $ProbeCeiling }
    if (-not [string]::IsNullOrWhiteSpace($Version))    { $passArgs.Version    = $Version }
    if (-not [string]::IsNullOrWhiteSpace($InstallDir)) { $passArgs.InstallDir = $InstallDir }
    if (-not [string]::IsNullOrWhiteSpace($Arch))       { $passArgs.Arch       = $Arch }
    if ($NoPath)    { $passArgs.NoPath    = $true }
    if ($Uninstall) { $passArgs.Uninstall = $true }
    if ($Force)     { $passArgs.Force     = $true }
    if ($KeepData)  { $passArgs.KeepData  = $true }
    if ($PurgeData) { $passArgs.PurgeData = $true }

    try {
        & $block @passArgs
    }
    catch {
        throw
    }
    return $true
}

if ($env:INSTALLER_DELEGATED -eq "1") {
    Write-Host "  [discovery] INSTALLER_DELEGATED=1; skipping discovery (loop guard)"
} elseif ($NoDiscovery) {
    Write-Host "  [discovery] -NoDiscovery set; skipping probe"
} elseif (-not [string]::IsNullOrWhiteSpace($Version)) {
    # Pinned-version contract (spec/07-generic-release/08-pinned-version-install-snippet.md):
    # When -Version is supplied, install EXACTLY that version from the embedded $Repo.
    # Skip versioned-repo discovery so a snippet copied from a v3.x release page
    # never silently jumps to the v4 repo's latest tag.
    Write-Host "  [discovery] -Version $Version pinned; skipping repo probe (exact-version install)"
} else {
    $effective = Resolve-EffectiveRepo $Repo $ProbeCeiling
    if ($effective -ne $Repo) {
        if (Invoke-DelegatedFullInstaller $effective) { return }
    }
}


# --- Logging helpers ---

function Write-Step([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Cyan
}

function Write-OK([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Green
}

function Write-Err([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Red
}

# Get-Sha256Hex computes a lowercase SHA256 hex digest. Uses Get-FileHash
# when available (PowerShell 3.0+) and falls back to .NET SHA256 so the
# installer works on stock Windows 7 / PowerShell 2.0 hosts where
# Get-FileHash is not a recognized cmdlet.
function Get-Sha256Hex([string]$path) {
    if (Get-Command Get-FileHash -ErrorAction SilentlyContinue) {
        return (Get-FileHash -Path $path -Algorithm SHA256).Hash.ToLower()
    }
    $sha = [System.Security.Cryptography.SHA256]::Create()
    try {
        $stream = [System.IO.File]::OpenRead($path)
        try {
            $bytes = $sha.ComputeHash($stream)
        } finally { $stream.Dispose() }
    } finally { $sha.Dispose() }
    return -join ($bytes | ForEach-Object { $_.ToString('x2') })
}

function Set-InstallerExitCode([int]$exitCode) {
    $global:LASTEXITCODE = $exitCode
    [System.Environment]::ExitCode = $exitCode
}

function Write-FatalError($record, [int]$exitCode = 1) {
    Set-InstallerExitCode $exitCode
    $message = "Unknown PowerShell error"
    if ($record) {
        if ($record.Exception -and -not [string]::IsNullOrWhiteSpace($record.Exception.Message)) {
            $message = $record.Exception.Message
        }
        elseif (-not [string]::IsNullOrWhiteSpace($record.ToString())) {
            $message = $record.ToString()
        }
    }

    Write-Host ""
    Write-Err "FATAL: $message"

    if ($record) {
        if (-not [string]::IsNullOrWhiteSpace($record.ScriptStackTrace)) {
            Write-Err ""
            Write-Err "Script stack trace:"
            foreach ($line in ($record.ScriptStackTrace -split "`r?`n")) {
                if (-not [string]::IsNullOrWhiteSpace($line)) {
                    Write-Err "  $line"
                }
            }
        }

        if ($record.InvocationInfo) {
            $scriptName = $record.InvocationInfo.ScriptName
            if ([string]::IsNullOrWhiteSpace($scriptName)) {
                $scriptName = $PSCommandPath
            }

            Write-Err ""
            Write-Err "Failure context:"
            if (-not [string]::IsNullOrWhiteSpace($scriptName)) {
                Write-Err "  Script: $scriptName"
            }
            Write-Err "  Line: $($record.InvocationInfo.ScriptLineNumber)"
            if (-not [string]::IsNullOrWhiteSpace($record.InvocationInfo.Line)) {
                Write-Err "  Code: $($record.InvocationInfo.Line.Trim())"
            }
        }

        if ($record.CategoryInfo) {
            Write-Err ""
            Write-Err "CategoryInfo: $($record.CategoryInfo)"
        }

        if (-not [string]::IsNullOrWhiteSpace($record.FullyQualifiedErrorId)) {
            Write-Err "FullyQualifiedErrorId: $($record.FullyQualifiedErrorId)"
        }

        if ($record.Exception) {
            Write-Err ""
            Write-Err "Exception:"
            foreach ($line in ($record.Exception.ToString() -split "`r?`n")) {
                if (-not [string]::IsNullOrWhiteSpace($line)) {
                    Write-Err "  $line"
                }
            }
        }
    }

    Write-Host ""
    return
}

# --- Resolve install directory ---

function Resolve-InstallDir([string]$dir) {
    if ($dir -ne "") { return $dir }
    return Join-Path $env:LOCALAPPDATA "gitmap"
}

# --- Detect architecture ---

function Resolve-Arch([string]$arch) {
    if ($arch -ne "") { return $arch }

    $cpu = $env:PROCESSOR_ARCHITECTURE
    switch ($cpu) {
        "AMD64"   { return "amd64" }
        "ARM64"   { return "arm64" }
        "x86"     { return "amd64" }
        default   { return "amd64" }
    }
}

# --- Resolve version (latest or pinned) ---

function Resolve-Version([string]$version) {
    if ($version -ne "") { return $version }

    $url = "https://api.github.com/repos/$Repo/releases/latest"
    Write-Step "Fetching latest release..."
    Write-Step "  URL: $url"

    try {
        $response = Invoke-WebRequest -Uri $url -UseBasicParsing -ErrorAction Stop
        $release = $response.Content | ConvertFrom-Json
        return $release.tag_name
    }
    catch {
        $statusCode = "unknown"
        $body = ""

        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
            try {
                $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
                $body = $reader.ReadToEnd()
                $reader.Close()
            }
            catch {
                $body = $_.Exception.Message
            }
        }
        else {
            $body = $_.Exception.Message
        }

        Write-Err "Failed to fetch latest release"
        Write-Err "  HTTP $statusCode -- $url"
        if ($body) {
            Write-Err "  Response: $body"
        }
        Write-Err ""
        Write-Err "  Possible causes:"
        Write-Err "    - No published releases in the repository"
        Write-Err "    - Repository is private (needs authentication)"
        Write-Err "    - Repository name has changed"
        Write-Err ""
        Write-Err "  Try: https://github.com/$Repo/releases"
        throw [InstallerFailure]::new("Failed to fetch latest release", 1)
    }
}

# --- Strict-tag failure (spec/07-generic-release/09 section 3) ---
# Print the canonical no-fallback message and exit 1. Called from
# Get-Asset whenever -Version was supplied explicitly and the
# requested release asset cannot be downloaded or verified.
function Stop-Strict([string]$detail) {
    Write-Err ""
    Write-Err "Error: requested release $Version not found in $Repo;"
    Write-Err "       refusing to fall back per strict-tag contract."
    Write-Err "       See spec/07-generic-release/09-generic-install-script-behavior.md `$3."
    if ($detail) { Write-Err "       Detail: $detail" }
    throw [InstallerFailure]::new("Strict version install failed", 1)
}

# --- Pre-flight asset existence ---

# Test-AssetExists: HEAD-probe the asset URL. Returns $true on
# HTTP 200/301/302, $false otherwise. Used to fail fast with a clear
# error before any download/extract work.
function Test-AssetExists([string]$url) {
    try {
        $resp = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 10 `
            -UseBasicParsing -ErrorAction Stop
        return ($resp.StatusCode -eq 200)
    } catch {
        return $false
    }
}

# Write-DryRunReport: emit a machine-parseable, key=value report of
# the resolved release-asset coordinates. Validates that the
# resolved name matches the documented contract
# `gitmap-vX.Y.Z-windows-{amd64|arm64}.zip` and exits non-zero with
# a clear error if it does not. CI parses the `dryrun.*=` lines.
function Write-DryRunReport([string]$version, [string]$arch,
                            [string]$assetName, [string]$assetUrl,
                            [string]$checksumUrl) {
    $pattern = '^gitmap-v\d+\.\d+\.\d+-windows-(amd64|arm64)\.zip$'
    $hasValidName = $assetName -match $pattern

    Write-Host ""
    Write-Host "================================================================"
    Write-Host " INSTALL.PS1 DRY-RUN REPORT"
    Write-Host "================================================================"
    Write-Host "dryrun.version=$version"
    Write-Host "dryrun.arch=$arch"
    Write-Host "dryrun.asset_name=$assetName"
    Write-Host "dryrun.asset_url=$assetUrl"
    Write-Host "dryrun.checksum_url=$checksumUrl"
    Write-Host "dryrun.expected_pattern=$pattern"
    Write-Host "dryrun.name_matches_contract=$hasValidName"
    Write-Host "dryrun.preflight_head=ok"
    Write-Host "================================================================"
    Write-Host ""

    if (-not $hasValidName) {
        Write-Err "::error::Resolved asset name '$assetName' does not match release contract '$pattern'"
        exit 5
    }
    Write-Host "OK install.ps1 dry-run passed for $version ($arch)"
}


function Write-MissingAssetError([string]$version, [string]$arch,
                                  [string]$assetName, [string]$assetUrl) {
    $releasePage = "https://github.com/$Repo/releases/tag/$version"
    Write-Err ""
    Write-Err "================================================================"
    Write-Err " RELEASE ASSET NOT FOUND"
    Write-Err "================================================================"
    Write-Err " Expected pattern: gitmap-$version-windows-$arch.zip"
    Write-Err " Resolved name:    $assetName"
    Write-Err " Probed URL:       $assetUrl"
    Write-Err " Release page:     $releasePage"
    Write-Err ""
    Write-Err " The release tag '$version' exists but does not publish the"
    Write-Err " expected Windows/$arch asset. This usually means:"
    Write-Err "   - the cross-compile job for this platform failed in CI,"
    Write-Err "   - the release was cut from a branch missing this target, or"
    Write-Err "   - you pinned a pre-release/draft tag with partial assets."
    Write-Err ""
    Write-Err " Inspect the release page above to see which assets ARE"
    Write-Err " present, then either pin a different -Version or wait for"
    Write-Err " the release pipeline to re-publish."
    Write-Err "================================================================"
    Write-Err ""
}

# --- Download asset ---

function Get-Asset([string]$version, [string]$arch) {
    $assetName = "gitmap-${version}-windows-${arch}.zip"
    $baseUrl = "https://github.com/$Repo/releases/download/$version"
    $assetUrl = "$baseUrl/$assetName"
    $checksumUrl = "$baseUrl/checksums.txt"

    # Strict mode: -Version was supplied explicitly. Any failure here
    # MUST exit 1 with the canonical message and MUST NOT fall back.
    $strict = -not [string]::IsNullOrWhiteSpace($Version)
    if ($strict) {
        Write-Step "  [strict] download: $assetUrl"
    }

    $tmpDir = Join-Path $env:TEMP "gitmap-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

    $zipPath = Join-Path $tmpDir $assetName
    $checksumPath = Join-Path $tmpDir "checksums.txt"

    # Pre-flight asset-existence check. Emits a clearly-formatted
    # error block when the expected release-asset pattern is missing,
    # so users see WHAT was expected, WHERE we looked, and the release
    # page to inspect — not just a generic "Download failed".
    if (-not (Test-AssetExists $assetUrl)) {
        Write-MissingAssetError $version $arch $assetName $assetUrl
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        if ($strict) {
            Stop-Strict "expected asset $assetName not found at $assetUrl"
        }
        throw [InstallerFailure]::new("Release asset not found", 1)
    }

    # Dry-run short-circuit: the URL exists, the naming contract is
    # honored — print a machine-parseable report and exit before any
    # download. CI greps these lines to assert correctness.
    if ($DryRun) {
        Write-DryRunReport $version $arch $assetName $assetUrl $checksumUrl
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 0
    }

    Write-Step "Downloading $assetName ($version)..."

    try {
        Invoke-WebRequest -Uri $assetUrl -OutFile $zipPath -UseBasicParsing
        Invoke-WebRequest -Uri $checksumUrl -OutFile $checksumPath -UseBasicParsing
    }
    catch {
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        if ($strict) {
            Stop-Strict "download failed: $($_.Exception.Message)"
        }
        Write-Err "Download failed: $_"
        throw [InstallerFailure]::new("Download failed", 1)
    }

    # Verify checksum
    Write-Step "Verifying checksum..."
    $expectedLine = (Get-Content $checksumPath | Where-Object { $_ -match $assetName })
    if (-not $expectedLine) {
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        if ($strict) {
            Stop-Strict "asset $assetName not listed in checksums.txt for $version"
        }
        Write-Err "Asset not found in checksums.txt"
        throw [InstallerFailure]::new("Asset not found in checksums.txt", 1)
    }

    $expectedHash = ($expectedLine -split '\s+')[0]
    $actualHash = Get-Sha256Hex $zipPath

    if ($actualHash -ne $expectedHash) {
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        if ($strict) {
            Stop-Strict "checksum mismatch for $assetName (expected $expectedHash, got $actualHash)"
        }
        Write-Err "Checksum mismatch!"
        Write-Err "  Expected: $expectedHash"
        Write-Err "  Got:      $actualHash"
        throw [InstallerFailure]::new("Checksum mismatch", 1)
    }

    Write-OK "Checksum verified."
    return @{ ZipPath = $zipPath; TmpDir = $tmpDir }
}

# --- Extract and install ---

function Install-Binary([string]$zipPath, [string]$installDir) {
    Write-Step "Installing to $installDir..."

    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    $targetPath = Join-Path $installDir $BinaryName
    $extractDir = Join-Path $installDir ".install-extract"

    # Rename-first strategy for running binary
    if (Test-Path $targetPath) {
        $oldPath = "$targetPath.old"
        if (Test-Path $oldPath) { Remove-Item $oldPath -Force }
        Rename-Item $targetPath $oldPath -Force
    }

    if (Test-Path $extractDir) {
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
    }

    New-Item -ItemType Directory -Path $extractDir -Force | Out-Null
    Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

    # Match exact names OR versioned patterns like gitmap-v28.54.6-windows-amd64.exe
    $candidateNames = @(
        $BinaryName,
        [System.IO.Path]::GetFileNameWithoutExtension($BinaryName),
        "gitmap-windows-amd64.exe",
        "gitmap-windows-arm64.exe"
    )

    $extracted = Get-ChildItem -Path $extractDir -File -Recurse |
        Where-Object {
            ($candidateNames -icontains $_.Name) -or
            ($_.Name -match "^gitmap-v[\d.]+-windows-(amd64|arm64)\.exe$")
        } |
        Select-Object -First 1

    if (-not $extracted) {
        $availableFiles = Get-ChildItem -Path $extractDir -File -Recurse | Select-Object -ExpandProperty FullName
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Err "Installed archive did not contain $BinaryName"
        if ($availableFiles) {
            Write-Err "Archive files:"
            foreach ($file in $availableFiles) {
                Write-Err "  $file"
            }
        }
        throw [InstallerFailure]::new("Installed archive did not contain $BinaryName", 1)
    }

    Move-Item $extracted.FullName $targetPath -Force

    Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue

    if (-not (Test-Path $targetPath)) {
        Write-Err "Install failed: $BinaryName was not written to $installDir"
        throw [InstallerFailure]::new("Install failed: $BinaryName was not written to $installDir", 1)
    }

    # Cleanup .old
    $oldPath = "$targetPath.old"
    if (Test-Path $oldPath) {
        Remove-Item $oldPath -Force -ErrorAction SilentlyContinue
    }

    Write-OK "Installed $BinaryName to $installDir"

    # Install short alias `gm.exe` alongside `gitmap.exe` so both commands
    # resolve to the same binary on PATH. Windows symlinks require admin or
    # Developer Mode, so a straight Copy-Item is the only reliable approach
    # for HKCU installs. Non-fatal: the primary install already succeeded.
    $aliasPath = Join-Path $installDir $BinaryAlias
    try {
        if (Test-Path $aliasPath) {
            Remove-Item $aliasPath -Force -ErrorAction Stop
        }
        Copy-Item -Path $targetPath -Destination $aliasPath -Force -ErrorAction Stop
        Write-OK "Installed alias $BinaryAlias -> $BinaryName (copy)"
    } catch {
        Write-Err "  ! Could not create '$BinaryAlias' alias next to ${BinaryName}: $($_.Exception.Message); only the full '$BinaryName' command will be available."
    }
}

# --- Download seed data files (downloader-config.json, etc.) ---
# The release zip only contains the binary, but the binary expects
# ./data/downloader-config.json to exist next to it. Fetch the seed
# files from the pinned tag on GitHub so the first run does not
# print "Could not read downloader seed".
function Install-SeedData([string]$version, [string]$installDir) {
    $dataDir = Join-Path $installDir "data"
    if (-not (Test-Path $dataDir)) {
        New-Item -ItemType Directory -Path $dataDir -Force | Out-Null
    }

    $seedFiles = @(
        "downloader-config.json",
        "config.json",
        "git-setup.json",
        "seo-templates.json"
    )

    Write-Step "Downloading seed data files ($version)..."

    $installed = 0
    foreach ($name in $seedFiles) {
        $rawUrl = "https://raw.githubusercontent.com/$Repo/$version/gitmap/data/$name"
        $dest = Join-Path $dataDir $name
        try {
            Invoke-WebRequest -Uri $rawUrl -OutFile $dest -UseBasicParsing -ErrorAction Stop
            $installed++
        }
        catch {
            # Best-effort: some files may not exist in older tags.
            Write-Host ("    skip  {0} (not in {1})" -f $name, $version) -ForegroundColor DarkGray
        }
    }

    if ($installed -gt 0) {
        Write-OK ("Installed {0} seed file(s) to {1}" -f $installed, $dataDir)
    }
    else {
        Write-Host ("    WARN  No seed files downloaded; gitmap will use built-in defaults") -ForegroundColor Yellow
    }
}

# --- Post-install self-check ---
# Verifies that the critical downloader-config.json seed file landed on
# disk after Install-SeedData. The binary refuses to read its downloader
# settings without this file and prints a confusing "Could not read
# downloader seed" warning on every subsequent run, so we fail loudly
# here with an actionable message instead of letting the user discover
# the problem 10 minutes later.
function Assert-InstallSelfCheck([string]$installDir) {
    Write-Step "Running install self-check..."

    $seedPath = Join-Path (Join-Path $installDir "data") "downloader-config.json"

    if (-not (Test-Path $seedPath)) {
        Write-Err ""
        Write-Err "Install self-check FAILED: required seed file is missing."
        Write-Err "  Expected: $seedPath"
        Write-Err ""
        Write-Err "  Why this matters:"
        Write-Err "    gitmap reads downloader-config.json on every run to"
        Write-Err "    configure parallel downloads, split sizes, and the"
        Write-Err "    Aria2C/GoDownloader fallback chain. Without it the"
        Write-Err "    binary will warn 'Could not read downloader seed'"
        Write-Err "    and operate with reduced functionality."
        Write-Err ""
        Write-Err "  How to fix:"
        Write-Err "    1) Re-run the installer with a working network connection, OR"
        Write-Err "    2) Manually download the file from:"
        Write-Err "         https://raw.githubusercontent.com/$Repo/main/gitmap/data/downloader-config.json"
        Write-Err "       and save it to:"
        Write-Err "         $seedPath"
        Write-Err ""
        throw [InstallerFailure]::new("Install self-check failed: downloader-config.json missing at $seedPath", 1)
    }

    # Sanity-check the file is non-empty JSON; a 0-byte file would
    # parse as a different "missing seed" error inside the binary.
    $info = Get-Item $seedPath
    if ($info.Length -le 0) {
        Write-Err "Install self-check FAILED: downloader-config.json is empty at $seedPath"
        throw [InstallerFailure]::new("Install self-check failed: downloader-config.json is empty at $seedPath", 1)
    }

    Write-OK "Self-check passed: downloader-config.json present ($($info.Length) bytes)."
}

# --- Download and extract docs-site.zip release asset ---
# Required for `gitmap help-dashboard` (hd). Best-effort: skip silently
# if the release does not bundle docs-site.zip (older versions).
function Install-DocsSite([string]$version, [string]$installDir) {
    $assetName = "docs-site.zip"
    $assetUrl = "https://github.com/$Repo/releases/download/$version/$assetName"
    $tmpZip = Join-Path $env:TEMP "gitmap-docs-site-$(Get-Random).zip"

    Write-Step "Downloading docs-site.zip ($version)..."

    try {
        Invoke-WebRequest -Uri $assetUrl -OutFile $tmpZip -UseBasicParsing -ErrorAction Stop
    }
    catch {
        Write-Step "  docs-site.zip not available for $version - skipping (gitmap hd will fall back to hosted docs)"
        Remove-Item $tmpZip -Force -ErrorAction SilentlyContinue
        return
    }

    # Remove any existing docs-site/ before extracting fresh.
    $docsDir = Join-Path $installDir "docs-site"
    if (Test-Path $docsDir) {
        Remove-Item $docsDir -Recurse -Force -ErrorAction SilentlyContinue
    }

    try {
        # The zip's internal layout is docs-site/dist/... so it extracts directly.
        Expand-Archive -Path $tmpZip -DestinationPath $installDir -Force
        Write-OK "Installed docs-site to $docsDir"
    }
    catch {
        Write-Err "Failed to extract docs-site.zip: $_"
    }
    finally {
        Remove-Item $tmpZip -Force -ErrorAction SilentlyContinue
    }
}

# --- Add to PATH ---

function Test-PathEntry([string]$pathValue, [string]$dir) {
    if ([string]::IsNullOrWhiteSpace($pathValue)) {
        return $false
    }

    $parts = $pathValue -split ";"

    foreach ($part in $parts) {
        if ($part.Trim() -ieq $dir) {
            return $true
        }
    }

    return $false
}

function Rebuild-SessionPath([string]$dir) {
    # Rebuild session PATH from registry (Machine + User) to pick up any changes
    $machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $parts = @()
    if ($machinePath) { $parts += $machinePath.TrimEnd(";") }
    if ($userPath) { $parts += $userPath.TrimEnd(";") }
    $rebuilt = $parts -join ";"

    # Ensure install dir is present even if not yet persisted
    if (-not (Test-PathEntry $rebuilt $dir)) {
        $rebuilt = $rebuilt.TrimEnd(";") + ";" + $dir
    }

    return $rebuilt
}

function Get-GitmapCommandWrapperBlock([string]$dir) {
	$safeDir = $dir.Replace("'", "''")
    $template = @'
# gitmap command wrapper v1
function global:Get-GitmapCommand { $candidate = Join-Path -Path '__GITMAP_DIR__' -ChildPath 'gitmap.exe'; if (Test-Path -LiteralPath $candidate) { return $candidate }; return (Get-Command gitmap.exe -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1).Source }
function global:Invoke-GitmapAndSetLocation { param([string[]]$GitMapArgs); $real = Get-GitmapCommand; if (-not $real) { Write-Error "gitmap executable not found"; return }; if ($GitMapArgs.Count -gt 0 -and ($GitMapArgs[0] -eq 'cd' -or $GitMapArgs[0] -eq 'go')) { $env:GITMAP_WRAPPER = "1"; $env:GITMAP_COMMAND_WRAPPER = "1"; $dest = [string](& $real @GitMapArgs | Out-String); if ($LASTEXITCODE -ne 0) { return }; $dest = $dest.Trim(); if ($dest -and (Test-Path -LiteralPath ([string]$dest))) { Set-Location -LiteralPath ([string]$dest) }; return }; $handoff = [IO.Path]::Combine([IO.Path]::GetTempPath(), "gitmap-handoff-$([Guid]::NewGuid().ToString('N')).txt"); try { $env:GITMAP_HANDOFF_FILE = $handoff; $env:GITMAP_WRAPPER = "1"; $env:GITMAP_COMMAND_WRAPPER = "1"; & $real @GitMapArgs; if ((Test-Path -LiteralPath $handoff) -and ((Get-Item -LiteralPath $handoff).Length -gt 0)) { $target = [string](Get-Content -LiteralPath $handoff -Raw); $target = $target.Trim(); if ($target -and (Test-Path -LiteralPath ([string]$target))) { Set-Location -LiteralPath ([string]$target) } } } finally { Remove-Item -LiteralPath $handoff -ErrorAction SilentlyContinue; Remove-Item Env:\GITMAP_HANDOFF_FILE -ErrorAction SilentlyContinue } }
function global:gcd { Invoke-GitmapAndSetLocation -GitMapArgs (@('cd') + $args) }
function global:gitmap { Invoke-GitmapAndSetLocation $args }
# gitmap command wrapper v1 end
'@
    return $template.Replace('__GITMAP_DIR__', $safeDir)
}

function Get-GitmapPowerShellShimContent([string]$dir) {
    $safeDir = $dir.Replace("'", "''")
    $template = @'
$real = Join-Path -Path '__GITMAP_DIR__' -ChildPath 'gitmap.exe'
if (-not (Test-Path -LiteralPath $real)) { Write-Error "gitmap executable not found: $real"; return }
if ($args.Count -gt 0 -and ($args[0] -eq 'cd' -or $args[0] -eq 'go')) {
  $env:GITMAP_WRAPPER = "1"; $env:GITMAP_COMMAND_WRAPPER = "1"
  $dest = [string](& $real @args | Out-String)
  if ($LASTEXITCODE -ne 0) { $global:LASTEXITCODE = $LASTEXITCODE; return }
  $dest = $dest.Trim()
  if ($dest -and (Test-Path -LiteralPath ([string]$dest))) { Set-Location -LiteralPath ([string]$dest) }
  return
}
$handoff = [IO.Path]::Combine([IO.Path]::GetTempPath(), "gitmap-handoff-$([Guid]::NewGuid().ToString('N')).txt")
try {
  $env:GITMAP_HANDOFF_FILE = $handoff; $env:GITMAP_WRAPPER = "1"; $env:GITMAP_COMMAND_WRAPPER = "1"
  & $real @args
  $exitCode = $LASTEXITCODE
  if ((Test-Path -LiteralPath $handoff) -and ((Get-Item -LiteralPath $handoff).Length -gt 0)) {
    $target = [string](Get-Content -LiteralPath $handoff -Raw); $target = $target.Trim()
    if ($target -and (Test-Path -LiteralPath ([string]$target))) { Set-Location -LiteralPath ([string]$target) }
  }
  $global:LASTEXITCODE = $exitCode
  return
}
finally {
  Remove-Item -LiteralPath $handoff -ErrorAction SilentlyContinue
  Remove-Item Env:\GITMAP_HANDOFF_FILE -ErrorAction SilentlyContinue
}
'@
    return $template.Replace('__GITMAP_DIR__', $safeDir)
}

function Install-PowerShellShim([string]$dir) {
    $shimPath = Join-Path $dir "gitmap.ps1"
    Set-Content -Path $shimPath -Value (Get-GitmapPowerShellShimContent $dir) -Encoding UTF8
    Write-OK "Installed PowerShell command shim: $shimPath"
}

function Add-CommandWrapperToProfile([string]$profilePath, [string]$dir) {
    $marker = "# gitmap command wrapper v1"
    $endMarker = "# gitmap command wrapper v1 end"
    $block = Get-GitmapCommandWrapperBlock $dir
    $profileDir = Split-Path $profilePath -Parent
    if ($profileDir -and -not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }
    if (-not (Test-Path $profilePath)) { Add-Content -Path $profilePath -Value $block -Encoding UTF8; return $true }
    $content = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
    if ($content -and ($content -match [regex]::Escape($marker))) {
        $pattern = "(?s)" + [regex]::Escape($marker) + ".*?" + [regex]::Escape($endMarker)
        $without = [regex]::Replace($content, $pattern, "", 1).TrimEnd()
        $rewritten = if ([string]::IsNullOrWhiteSpace($without)) { "$block`r`n" } else { "$without`r`n`r`n$block`r`n" }
        Set-Content -Path $profilePath -Value $rewritten -Encoding UTF8
        return $true
    }
    Add-Content -Path $profilePath -Value $block -Encoding UTF8
    return $true
}

function Install-CommandWrapperForSession([string]$dir) {
    Invoke-Expression (Get-GitmapCommandWrapperBlock $dir)
}

function Resolve-PowerShellProfileTargets {
    $targets = @()
    if ($PROFILE) { $targets += $PROFILE }
    if ($env:USERPROFILE) {
        $targets += (Join-Path $env:USERPROFILE "Documents\PowerShell\profile.ps1")
        $targets += (Join-Path $env:USERPROFILE "Documents\PowerShell\Microsoft.PowerShell_profile.ps1")
        $targets += (Join-Path $env:USERPROFILE "Documents\WindowsPowerShell\profile.ps1")
        $targets += (Join-Path $env:USERPROFILE "Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1")
    }
    return @($targets | Where-Object { -not [string]::IsNullOrWhiteSpace($_) } | Select-Object -Unique)
}

function Update-PowerShellProfilePathLine([string]$profilePath, [string]$dir) {
    $marker = "# gitmap-path"
    $exportLine = "`$env:PATH = `"$dir;`$env:PATH`" $marker"
    $profileDir = Split-Path $profilePath -Parent
    if ($profileDir -and -not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }
    if (-not (Test-Path $profilePath)) {
        Add-Content -Path $profilePath -Value "`n$exportLine" -Encoding UTF8
        return $true
    }
    $lines = Get-Content $profilePath -ErrorAction SilentlyContinue
    $filtered = @($lines | Where-Object { $_ -notmatch [regex]::Escape($marker) })
    $filtered += $exportLine
    $next = ($filtered -join "`r`n") + "`r`n"
    $current = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
    if ($current -eq $next) { return $false }
    Set-Content -Path $profilePath -Value $next -Encoding UTF8
    return $true
}

function Broadcast-EnvironmentChange {
    Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

public static class GitMapEnvNative {
    [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd,
        uint Msg,
        UIntPtr wParam,
        string lParam,
        uint fuFlags,
        uint uTimeout,
        out UIntPtr lpdwResult
    );
}
"@ -ErrorAction SilentlyContinue | Out-Null

    $HWND_BROADCAST = [IntPtr]0xffff
    $WM_SETTINGCHANGE = 0x001A
    $SMTO_ABORTIFHUNG = 0x0002
    [UIntPtr]$result = [UIntPtr]::Zero

    [void][GitMapEnvNative]::SendMessageTimeout(
        $HWND_BROADCAST,
        $WM_SETTINGCHANGE,
        [UIntPtr]::Zero,
        "Environment",
        $SMTO_ABORTIFHUNG,
        5000,
        [ref]$result
    )
}

function Add-ToPath([string]$dir) {
    $modified = @()
    Install-PowerShellShim $dir

    # --- 1. Windows User PATH (registry) -- covers CMD + new PowerShell windows ---
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $userHasDir = Test-PathEntry $currentUserPath $dir

    if (-not $userHasDir) {
        if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
            $newPath = $dir
        }
        else {
            $newPath = $currentUserPath.TrimEnd(";") + ";" + $dir
        }

        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Broadcast-EnvironmentChange
        Write-OK "Added to User PATH (registry)."
        $modified += "User PATH (registry)"
    }
    else {
        Write-Step "Already in User PATH (registry)."
    }

    # --- 2. PowerShell $PROFILE -- ensures PATH in all PowerShell sessions ---
    $psProfilePath = $PROFILE
    if ($psProfilePath) {
        if (Update-PowerShellProfilePathLine $psProfilePath $dir) {
            Write-OK "Added to PowerShell profile: $psProfilePath"
            $modified += "PowerShell `$PROFILE"
        }
        else {
            Write-Step "Already in PowerShell profile."
        }

        $wrapperProfiles = Resolve-PowerShellProfileTargets
        foreach ($wrapperProfile in $wrapperProfiles) {
            if (Add-CommandWrapperToProfile $wrapperProfile $dir) {
                Write-OK "Added gitmap cd command wrapper to PowerShell profile: $wrapperProfile"
            }
        }
        if ($wrapperProfiles.Count -gt 0) { $modified += "PowerShell command wrapper" }
    }

    Install-CommandWrapperForSession $dir

    # --- 3. Git Bash profiles (~/.bashrc, ~/.bash_profile) ---
    $homeDir = $env:USERPROFILE
    if ($homeDir) {
        $bashExportLine = "export PATH=`"$($dir -replace '\\','/'):`$PATH`""
        $bashMarker = "# gitmap-path"
        $bashProfiles = @(
            (Join-Path $homeDir ".bashrc"),
            (Join-Path $homeDir ".bash_profile")
        )

        foreach ($bashProfile in $bashProfiles) {
            $bashAlreadyPresent = $false
            $bashProfileName = Split-Path $bashProfile -Leaf

            if (Test-Path $bashProfile) {
                $bashContent = Get-Content $bashProfile -Raw -ErrorAction SilentlyContinue
                if ($bashContent -and ($bashContent -match [regex]::Escape($bashMarker))) {
                    $bashAlreadyPresent = $true
                }
            }

            if (-not $bashAlreadyPresent) {
                Add-Content -Path $bashProfile -Value "`n$bashExportLine $bashMarker" -Encoding UTF8
                Write-OK "Added to Git Bash profile: ~/$bashProfileName"
                $modified += "~/$bashProfileName"
            }
            else {
                Write-Step "Already in ~/$bashProfileName."
            }
        }
    }

    if ($modified.Count -gt 0) {
        return @{ Target = ($modified -join ", "); Status = "added" }
    }

    return @{ Target = "All profiles"; Status = "already present" }
}

# --- Remove from PATH (uninstall) ---

function Remove-FromPath([string]$dir) {
    $removed = @()
    $marker = "# gitmap-path"

    # --- 1. Windows User PATH (registry) ---
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentUserPath -and (Test-PathEntry $currentUserPath $dir)) {
        $parts = ($currentUserPath -split ";") | Where-Object { $_.Trim() -ine $dir -and $_.Trim() -ne "" }
        $newPath = $parts -join ";"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Broadcast-EnvironmentChange
        Write-OK "Removed from User PATH (registry)."
        $removed += "User PATH (registry)"
    }

    # --- 2. PowerShell $PROFILE ---
    $psProfilePath = $PROFILE
    if ($psProfilePath -and (Test-Path $psProfilePath)) {
        $lines = Get-Content $psProfilePath
        $filtered = $lines | Where-Object { $_ -notmatch [regex]::Escape($marker) }
        if ($filtered.Count -lt $lines.Count) {
            $filtered | Set-Content $psProfilePath -Encoding UTF8
            Write-OK "Removed marker lines from PowerShell profile: $psProfilePath"
            $removed += "PowerShell `$PROFILE"
        }
    }

    # --- 3. Git Bash profiles ---
    $homeDir = $env:USERPROFILE
    if ($homeDir) {
        $bashProfiles = @(
            (Join-Path $homeDir ".bashrc"),
            (Join-Path $homeDir ".bash_profile")
        )

        foreach ($bashProfile in $bashProfiles) {
            if (Test-Path $bashProfile) {
                $lines = Get-Content $bashProfile
                $filtered = $lines | Where-Object { $_ -notmatch [regex]::Escape($marker) }
                if ($filtered.Count -lt $lines.Count) {
                    $filtered | Set-Content $bashProfile -Encoding UTF8
                    $name = Split-Path $bashProfile -Leaf
                    Write-OK "Removed marker lines from ~/$name"
                    $removed += "~/$name"
                }
            }
        }
    }

    if ($removed.Count -gt 0) {
        Write-Host ""
        Write-OK "PATH entries removed from: $($removed -join ', ')"
    }
    else {
        Write-Step "No gitmap PATH entries found in any profile."
    }

    return $removed
}

function Write-InstallSummary([string]$version, [string]$binPath, [string]$installDir, [hashtable]$pathResult, [bool]$isNoPath) {
    Write-Host ""
    Write-Host "  -----------------------------------------------" -ForegroundColor DarkGray
    Write-Host "  gitmap install summary" -ForegroundColor White
    Write-Host "  -----------------------------------------------" -ForegroundColor DarkGray
    Write-Host "    Version    : $version"
    Write-Host "    Binary     : $binPath"
    Write-Host "    Install Dir: $installDir"

    if ($isNoPath) {
        Write-Host "    PATH       : skipped (-NoPath)"
        return
    }

    Write-Host "    PATH target: $($pathResult.Target) ($($pathResult.Status))"
    Write-Host "    Session    : refreshed for current PowerShell session"

    Write-Host ""
    Write-Host "  Profiles modified:" -ForegroundColor White
    Write-Host "    - User PATH (registry)  : CMD, new PowerShell windows"
    Write-Host "    - PowerShell `$PROFILE    : all PowerShell sessions"
    Write-Host "    - ~/.bashrc             : Git Bash interactive shells"
    Write-Host "    - ~/.bash_profile       : Git Bash login shells"

    Write-Host ""
    Write-Host "  If gitmap is not found in a new terminal, run:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "    PowerShell:  `$env:PATH = `"$installDir;`$env:PATH`"" -ForegroundColor Cyan
    Write-Host "    CMD:         set PATH=$installDir;%PATH%" -ForegroundColor Cyan
    Write-Host "    Git Bash:    export PATH=`"$($installDir -replace '\\','/'):`$PATH`"" -ForegroundColor Cyan
    Write-Host ""
}

# --- Post-install verification ---

# Invoke-InstallVerification runs the three post-install checks the
# user asked for: (1) print installed version by invoking the binary
# directly, (2) confirm `gitmap` resolves via Get-Command in the
# refreshed session PATH, (3) ensure the per-install data folder
# exists (create on miss). All checks emit PASS/WARN; none throw,
# because the binary is already on disk and the user can recover.
function Invoke-InstallVerification([string]$binPath, [string]$installDir, [bool]$isNoPath) {
    $dataDir = Join-Path $installDir "data"

    Write-Host ""
    Write-Step "Verifying installation"

    # 1. Version
    if (Test-Path $binPath) {
        try {
            $verLine = (& $binPath version 2>&1 | Out-String).Trim().Split("`n")[0]
            Write-Host ("    PASS  Version: {0}" -f $verLine) -ForegroundColor Green
        }
        catch {
            Write-Host ("    WARN  Could not run {0} version: {1}" -f $binPath, $_) -ForegroundColor Yellow
        }
    }
    else {
        Write-Host ("    WARN  Binary missing: {0}" -f $binPath) -ForegroundColor Yellow
    }

    # 2. PATH active in this session
    $resolved = Get-Command $BinaryName -ErrorAction SilentlyContinue
    if ($resolved) {
        Write-Host ("    PASS  PATH active: {0} -> {1}" -f $BinaryName, $resolved.Source) -ForegroundColor Green
    }
    elseif ($isNoPath) {
        Write-Host ("    WARN  PATH skipped (-NoPath); invoke with full path: {0}" -f $binPath) -ForegroundColor Yellow
    }
    else {
        Write-Host ("    WARN  {0} not on PATH yet - open a new terminal or reload `$PROFILE." -f $BinaryName) -ForegroundColor Yellow
    }

    # 3. Data folder
    if (Test-Path $dataDir) {
        Write-Host ("    PASS  Data folder exists: {0}" -f $dataDir) -ForegroundColor Green
    }
    else {
        try {
            New-Item -ItemType Directory -Path $dataDir -Force | Out-Null
            Write-Host ("    PASS  Data folder created: {0}" -f $dataDir) -ForegroundColor Green
        }
        catch {
            Write-Host ("    WARN  Could not create data folder: {0}" -f $dataDir) -ForegroundColor Yellow
        }
    }
}

# --- Main ---

function Main {
    Write-Host ""
    Write-Host "  gitmap installer v$InstallerVersion" -ForegroundColor White
    Write-Host "  github.com/$Repo" -ForegroundColor DarkGray
    Write-Host ""

    try {
        $resolvedArch = Resolve-Arch $Arch
        $resolvedVersion = Resolve-Version $Version
        $resolvedDir = Resolve-InstallDir $InstallDir

        $result = Get-Asset $resolvedVersion $resolvedArch

        try {
            Install-Binary $result.ZipPath $resolvedDir
        }
        finally {
            Remove-Item $result.TmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }

        # Seed data files (downloader-config.json, etc.) so the binary
        # does not warn about a missing seed on first run.
        Install-SeedData $resolvedVersion $resolvedDir

        # Verify the seed actually landed before we declare success.
        Assert-InstallSelfCheck $resolvedDir

        # Bundle the docs site so `gitmap help-dashboard` works after install.
        Install-DocsSite $resolvedVersion $resolvedDir

        Install-PowerShellShim $resolvedDir

        $pathResult = @{ Target = "-NoPath"; Status = "skipped" }
        if (-not $NoPath) {
            $pathResult = Add-ToPath $resolvedDir

            # NOTE: We intentionally do NOT call Chocolatey's `refreshenv` here.
            # On Windows 11 24H2 / Server 2025 `wmic.exe` was removed from the
            # OS, and older `refreshenv` implementations shell out to
            # `wmic process get parentprocessid ...` which prints
            #   'wmic' is not recognized as an internal or external command
            # to the user's console. PowerShell `try/catch` cannot suppress
            # native-stderr from a `cmd.exe` child, so the noise leaks
            # through. `Rebuild-SessionPath` below already reads the
            # Machine + User PATH from the registry, which is exactly what
            # `refreshenv` would have done -- so the call is redundant.

            # Force-rebuild $env:PATH in this scope so gitmap is usable immediately
            $script:NewPath = Rebuild-SessionPath $resolvedDir
        }
        else {
            $script:NewPath = $env:PATH
        }

        return @{ InstallDir = $resolvedDir; NewPath = $script:NewPath; Version = $resolvedVersion; PathResult = $pathResult }
    }
    catch {
        throw
    }
}

# --- Uninstall safety helpers ---

# Confirm-IsGitmapInstall verifies the binary at $binPath actually
# IS our gitmap CLI before any destructive action. We invoke
# `<binary> version` and look for the literal token "gitmap" in the
# output. A mismatch means the user pointed -InstallDir at the wrong
# folder (or a stale path collides with another tool) -- bail loudly
# rather than rip PATH entries belonging to something else. -Force
# overrides the guard for advanced users.
function Confirm-IsGitmapInstall([string]$binPath, [bool]$isForce) {
    if (-not (Test-Path $binPath)) {
        # Nothing to verify against -- also nothing to wreck. Allow
        # PATH cleanup so a half-installed/broken state can be removed.
        Write-Step "No binary at $binPath; skipping identity check."
        return $true
    }
    try {
        $out = (& $binPath version 2>&1 | Out-String)
    }
    catch {
        $out = ""
    }
    if ($out -match '(?i)\bgitmap\b') {
        Write-OK "Verified gitmap binary at $binPath."
        return $true
    }
    if ($isForce) {
        Write-Step "Identity check failed but -Force set; continuing."
        return $true
    }
    Write-Err "Refusing to uninstall: $binPath does not look like gitmap."
    Write-Err "  Output: $($out.Trim())"
    Write-Err "  Re-run with -Force to override, or pass the correct -InstallDir."
    return $false
}

# Resolve-DataChoice returns 'keep' or 'purge'. Honors -KeepData /
# -PurgeData unconditionally; -Force defaults to 'keep' (safer);
# otherwise prompts the user. When stdin is not interactive (e.g.
# `iex` in a non-tty pipeline) we also default to 'keep'.
function Resolve-DataChoice([string]$dataDir, [bool]$isKeep, [bool]$isPurge, [bool]$isForce) {
    if (-not (Test-Path $dataDir)) { return 'keep' }
    if ($isPurge) { return 'purge' }
    if ($isKeep) { return 'keep' }
    if ($isForce) { return 'keep' }
    if (-not [Environment]::UserInteractive) { return 'keep' }
    Write-Host ""
    Write-Host ("  Data folder found: {0}" -f $dataDir) -ForegroundColor Yellow
    $reply = Read-Host "  Delete user data too? [y/N]"
    if ($reply -match '^(y|yes)$') { return 'purge' }
    return 'keep'
}

# --- Uninstall mode ---
if ($Uninstall) {
    Write-Host ""
    Write-Host "  gitmap uninstaller" -ForegroundColor White
    Write-Host ""

    $resolvedDir = Resolve-InstallDir $InstallDir
    $binPath = Join-Path $resolvedDir $BinaryName
    $dataDir = Join-Path $resolvedDir "data"

    # Safety guard: confirm this is actually gitmap before we touch
    # PATH entries or delete anything. -Force bypasses the check.
    if (-not (Confirm-IsGitmapInstall $binPath $Force.IsPresent)) {
        Write-Host ""
        return
    }

    $dataChoice = Resolve-DataChoice $dataDir $KeepData.IsPresent $PurgeData.IsPresent $Force.IsPresent

    # Remove PATH entries from all profiles
    $removedProfiles = Remove-FromPath $resolvedDir

    # Remove binary and install directory
    if (Test-Path $binPath) {
        Remove-Item $binPath -Force -ErrorAction SilentlyContinue
        Write-OK "Removed binary: $binPath"
    }

    $oldPath = "$binPath.old"
    if (Test-Path $oldPath) {
        Remove-Item $oldPath -Force -ErrorAction SilentlyContinue
    }

    # Honor the data-folder choice BEFORE the empty-dir sweep so a
    # 'keep' run does not get the install dir removed out from under
    # the surviving data folder.
    if ($dataChoice -eq 'purge' -and (Test-Path $dataDir)) {
        Remove-Item $dataDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-OK "Removed data folder: $dataDir"
    }
    elseif (Test-Path $dataDir) {
        Write-Step "Kept data folder: $dataDir"
    }

    # Remove install dir if empty
    if ((Test-Path $resolvedDir) -and @(Get-ChildItem $resolvedDir).Count -eq 0) {
        Remove-Item $resolvedDir -Force -ErrorAction SilentlyContinue
        Write-OK "Removed empty directory: $resolvedDir"
    }

    # Rebuild session PATH without the dir
    $machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $parts = @()
    if ($machinePath) { $parts += ($machinePath -split ";") }
    if ($userPath) { $parts += ($userPath -split ";") }
    $env:PATH = ($parts | Where-Object { $_.Trim() -ine $resolvedDir -and $_.Trim() -ne "" }) -join ";"

    Write-Host ""
    Write-OK "gitmap has been uninstalled."
    Write-Host ""
    return
}


try {
    $installResult = Main

    if (-not $installResult) {
        Set-InstallerExitCode 1
        return
    }

    # Set $env:PATH at the TOP-LEVEL script scope (not inside a function)
    # This ensures the change persists in the caller's session when run via iex
    $env:PATH = $installResult.NewPath

    # Verify the binary works.
    #
    # NOTE: We deliberately discard the binary's stderr (`2>$null`) and only
    # parse stdout. Reasons:
    #   1. First-run side effects like `SeedDownloaderConfig` print info
    #      lines (e.g. "◦ Downloader config already customized …") to
    #      stderr. With `2>&1` and `$ErrorActionPreference='Stop'`, PowerShell
    #      wraps those into ErrorRecords and throws, falsely tripping the
    #      "Binary found but failed to run" branch.
    #   2. Even when not throwing, those bytes were rendered as cp1252
    #      mojibake ("Γùª", "ΓÇö") in the user's session because PowerShell
    #      decodes child stderr via `[Console]::OutputEncoding` which is not
    #      always honored across hosts (`irm | iex`, ConEmu, ISE).
    # We then filter to the canonical `gitmap vX.Y.Z` line only.
    $binPath = Join-Path $installResult.InstallDir $BinaryName
    $installedVersion = $installResult.Version
    if (Test-Path $binPath) {
        Write-Host ""
        try {
            $versionOutput = & $binPath version 2>$null
            $versionLine = $versionOutput |
                ForEach-Object { $_.ToString() } |
                Where-Object { $_ -match '^gitmap v[0-9]' } |
                Select-Object -First 1
            if ($versionLine) {
                $installedVersion = $versionLine.Trim()
                Write-OK $installedVersion
            } else {
                Write-OK "gitmap $installedVersion"
            }
        }
        catch {
            Write-Err "Binary found but failed to run: $_"
        }
    }
    else {
        Write-Err "Binary not found at $binPath"
    }

    Write-InstallSummary $installedVersion $binPath $installResult.InstallDir $installResult.PathResult $NoPath.IsPresent

    Invoke-InstallVerification $binPath $installResult.InstallDir $NoPath.IsPresent

    # v5.18.0+: auto-run `gitmap setup` so the shell wrapper (gcd /
    # gitmap function) + completions are installed without a second
    # command. Setup is idempotent (marker `# gitmap shell wrapper v2`).
    # Non-fatal: install itself already succeeded.
    if (Test-Path -LiteralPath $binPath) {
        Write-Host ""
        Write-Host "  -> Running 'gitmap setup' to install shell wrapper + completions..." -ForegroundColor Cyan
        try {
            & $binPath setup
            if ($LASTEXITCODE -ne 0) {
                Write-Host "  (setup auto-run exited $LASTEXITCODE; install already succeeded)" -ForegroundColor Yellow
            }
        } catch { Write-Host "  (setup auto-run skipped: $_)" -ForegroundColor Yellow }
    }

    Write-Host ""
    Write-OK "Done! Run 'gitmap --help' to get started."
    Write-Host ""
    Set-InstallerExitCode 0
}
catch {
    Write-FatalError $_ 1
    return
}
