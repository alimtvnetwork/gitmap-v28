<#
.SYNOPSIS
    Version-pinned installer for a specific gitmap release.

.DESCRIPTION
    Installs EXACTLY the version requested via -Version. Never resolves
    "latest", never auto-upgrades, never silently substitutes. Designed
    for use from /release/<version> pages where the URL itself is the
    contract.

    Spec: spec/01-app/105-release-version-script.md

.PARAMETER Version
    Required. Tag of the release to install (e.g. v3.36.0).
    In a snapshot copy this value is baked at the top of the file and
    the parameter is ignored.

.PARAMETER InstallDir
    Target directory. Default: $env:LOCALAPPDATA\gitmap\bin

.PARAMETER Arch
    Force architecture (amd64, arm64). Default: auto-detect.

.PARAMETER NoPath
    Skip adding the install directory to PATH.

.PARAMETER NoSelfInstall
    Skip the chained `gitmap self-install` step (download + extract only).

.PARAMETER AllowFallback
    If the requested version is missing, use the newest patch in the same
    vMAJOR.MINOR series instead of failing. No interactive prompt.

.PARAMETER Quiet
    Suppress prompts and progress output. Non-interactive failure mode:
    a missing version causes immediate exit 1.

.EXAMPLE
    iwr https://gitmap.dev/scripts/release-version.ps1 -OutFile rv.ps1
    .\rv.ps1 -Version v3.36.0

.NOTES
    Repository: https://github.com/alimtvnetwork/gitmap-v27
#>

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [string]$Arch = "",
    [switch]$NoPath,
    [switch]$NoSelfInstall,
    [switch]$AllowFallback,
    [switch]$Quiet
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# --- Repo / asset metadata ---
$Repo = "alimtvnetwork/gitmap-v27"
$BinaryName = "gitmap.exe"

# --- Exit codes (spec 105) ---
$EXIT_OK              = 0
$EXIT_VERSION_MISSING = 1
$EXIT_NETWORK         = 2
$EXIT_CHECKSUM        = 3
$EXIT_UNSUPPORTED_ARCH= 4
$EXIT_PATH_FAIL       = 5
$EXIT_SELF_INSTALL    = 6
$EXIT_VERIFY          = 7

# --- Logging helpers (ASCII only — no Unicode glyphs) ---
function Write-Step([string]$msg) { if (-not $Quiet) { Write-Host "  -> $msg" -ForegroundColor Cyan } }
function Write-OK([string]$msg)   { if (-not $Quiet) { Write-Host "  OK $msg"  -ForegroundColor Green } }
function Write-Warn2([string]$m)  { if (-not $Quiet) { Write-Host "  !  $m"    -ForegroundColor Yellow } }
function Write-Err2([string]$m)   { Write-Host "  X  $m" -ForegroundColor Red }

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
    Write-Err2 "FATAL: $message"

    if ($record) {
        if (-not [string]::IsNullOrWhiteSpace($record.ScriptStackTrace)) {
            Write-Err2 ""
            Write-Err2 "Script stack trace:"
            foreach ($line in ($record.ScriptStackTrace -split "`r?`n")) {
                if (-not [string]::IsNullOrWhiteSpace($line)) {
                    Write-Err2 "  $line"
                }
            }
        }

        if ($record.InvocationInfo) {
            $scriptName = $record.InvocationInfo.ScriptName
            if ([string]::IsNullOrWhiteSpace($scriptName)) {
                $scriptName = $PSCommandPath
            }

            Write-Err2 ""
            Write-Err2 "Failure context:"
            if (-not [string]::IsNullOrWhiteSpace($scriptName)) {
                Write-Err2 "  Script: $scriptName"
            }
            Write-Err2 "  Line: $($record.InvocationInfo.ScriptLineNumber)"
            if (-not [string]::IsNullOrWhiteSpace($record.InvocationInfo.Line)) {
                Write-Err2 "  Code: $($record.InvocationInfo.Line.Trim())"
            }
        }

        if ($record.CategoryInfo) {
            Write-Err2 ""
            Write-Err2 "CategoryInfo: $($record.CategoryInfo)"
        }

        if (-not [string]::IsNullOrWhiteSpace($record.FullyQualifiedErrorId)) {
            Write-Err2 "FullyQualifiedErrorId: $($record.FullyQualifiedErrorId)"
        }

        if ($record.Exception) {
            Write-Err2 ""
            Write-Err2 "Exception:"
            foreach ($line in ($record.Exception.ToString() -split "`r?`n")) {
                if (-not [string]::IsNullOrWhiteSpace($line)) {
                    Write-Err2 "  $line"
                }
            }
        }
    }

    Write-Host ""
    return
}

# ---------------------------------------------------------------------------
# Version validation
# ---------------------------------------------------------------------------
function Test-VersionTag([string]$tag) {
    return $tag -match '^v\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$'
}

# ---------------------------------------------------------------------------
# OS / architecture detection
#
# Windows-only script, but we still validate the OS in case someone runs it
# under a Unix-targeted PowerShell Core (`pwsh` on Linux/macOS) by mistake —
# in that case we point them at release-version.sh.
# ---------------------------------------------------------------------------
function Resolve-OS() {
    if ($IsWindows -or $env:OS -eq 'Windows_NT') {
        return "windows"
    }
    Write-Err2 "release-version.ps1 only runs on Windows. For macOS/Linux use release-version.sh."
    exit $EXIT_UNSUPPORTED_ARCH
}

function Resolve-Arch([string]$override) {
    if (-not [string]::IsNullOrWhiteSpace($override)) {
        $a = $override.ToLower()
        if ($a -in @("amd64","arm64")) { return $a }
        Write-Err2 "Unsupported -Arch override: $override (allowed: amd64, arm64)"
        exit $EXIT_UNSUPPORTED_ARCH
    }
    $cpu = $env:PROCESSOR_ARCHITECTURE
    switch ($cpu) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        "x86"   { return "amd64" }   # 32-bit shell on a 64-bit host
        default {
            Write-Err2 "Unrecognized PROCESSOR_ARCHITECTURE: $cpu"
            exit $EXIT_UNSUPPORTED_ARCH
        }
    }
}

# ---------------------------------------------------------------------------
# GitHub release lookup
# ---------------------------------------------------------------------------
function Invoke-GitHubAPI([string]$path) {
    $url = "https://api.github.com/repos/$Repo$path"
    $headers = @{ "User-Agent" = "gitmap-release-version-installer" }
    if ($env:GITHUB_TOKEN) {
        $headers["Authorization"] = "Bearer $($env:GITHUB_TOKEN)"
    }
    try {
        return Invoke-RestMethod -Uri $url -Headers $headers -UseBasicParsing -ErrorAction Stop
    } catch {
        if ($_.Exception.Response.StatusCode.value__ -eq 404) { return $null }
        Write-Err2 "GitHub API error: $($_.Exception.Message)"
        exit $EXIT_NETWORK
    }
}

function Get-RecentReleases([int]$count = 5) {
    $list = Invoke-GitHubAPI "/releases?per_page=$count"
    if ($null -eq $list) { return @() }
    return @($list | Where-Object { -not $_.draft -and -not $_.prerelease } |
        Select-Object -First $count -ExpandProperty tag_name)
}

function Resolve-FallbackPatch([string]$requested) {
    # Same vMAJOR.MINOR series, newest patch.
    if ($requested -notmatch '^v(\d+)\.(\d+)\.\d+') { return $null }
    $major = $Matches[1]; $minor = $Matches[2]
    $list = Invoke-GitHubAPI "/releases?per_page=100"
    if ($null -eq $list) { return $null }
    $candidates = @($list |
        Where-Object { -not $_.draft -and -not $_.prerelease -and $_.tag_name -match "^v$major\.$minor\.\d+" } |
        ForEach-Object { $_.tag_name } |
        Sort-Object -Descending {
            if ($_ -match '^v\d+\.\d+\.(\d+)') { [int]$Matches[1] } else { 0 }
        })
    if ($candidates.Count -gt 0) { return $candidates[0] }
    return $null
}

function Resolve-RequestedVersion([string]$requested) {
    if (-not (Test-VersionTag $requested)) {
        Write-Err2 "Invalid version tag: '$requested' (expected vMAJOR.MINOR.PATCH)"
        exit $EXIT_VERSION_MISSING
    }

    $rel = Invoke-GitHubAPI "/releases/tags/$requested"
    if ($null -ne $rel) { return $requested }

    Write-Err2 "Requested version $requested is not a published release."

    if ($AllowFallback) {
        $fb = Resolve-FallbackPatch $requested
        if ($fb) {
            Write-Warn2 "Falling back to newest patch in series: $fb"
            return $fb
        }
        Write-Err2 "No same-minor-series patch available for $requested"
        exit $EXIT_VERSION_MISSING
    }

    if (-not (Test-Interactive)) {
        Write-Err2 "Non-interactive run; refusing to substitute. Set -AllowFallback to opt in."
        exit $EXIT_VERSION_MISSING
    }

    $recent = Get-RecentReleases 5
    if ($recent.Count -eq 0) {
        Write-Err2 "Could not list recent releases."
        exit $EXIT_VERSION_MISSING
    }

    Write-Host ""
    Write-Host "  Requested: $requested (not found)" -ForegroundColor Yellow
    Write-Host "  Most recent published releases:" -ForegroundColor Yellow
    for ($i = 0; $i -lt $recent.Count; $i++) {
        Write-Host ("    [{0}] {1}" -f ($i + 1), $recent[$i])
    }
    Write-Host "    [N] Quit (default)"

    $reply = Read-PromptSafe "  Pick a number to install instead, or N to quit"
    if ($null -eq $reply -or [string]::IsNullOrWhiteSpace($reply) -or $reply -match '^[Nn]') {
        exit $EXIT_VERSION_MISSING
    }
    $idx = 0
    if (-not [int]::TryParse($reply, [ref]$idx) -or $idx -lt 1 -or $idx -gt $recent.Count) {
        Write-Err2 "Invalid choice."
        exit $EXIT_VERSION_MISSING
    }
    $chosen = $recent[$idx - 1]
    Write-Warn2 "User selected $chosen as substitute for $requested"
    return $chosen
}

# Test-Interactive returns $true only when we can safely show a prompt and
# read from a real keyboard. The combination matters because each predicate
# alone is wrong in at least one common shell:
#   - `iwr ... | iex` runs with $Host.Name = 'ConsoleHost' but stdin is the
#     piped scriptblock, so Read-Host hangs forever.
#   - CI runners (GitHub Actions, etc.) report UserInteractive=$false but
#     still have a console; honour the $env:CI hint as the override.
function Test-Interactive() {
    if ($Quiet) { return $false }
    if ($env:CI -eq 'true' -or $env:CI -eq '1') { return $false }
    if (-not [Environment]::UserInteractive) { return $false }
    try {
        if ([Console]::IsInputRedirected) { return $false }
    } catch {
        # Older PowerShell hosts lack IsInputRedirected — fall back to a
        # cautious "no" so we never hang.
        return $false
    }
    return $true
}

# Read-PromptSafe wraps Read-Host in a try/catch so a closed stdin (which
# can happen when the script is piped from iwr | iex) returns $null instead
# of throwing an unhelpful pipeline error.
function Read-PromptSafe([string]$prompt) {
    try {
        return Read-Host $prompt
    } catch {
        Write-Err2 "Could not read from stdin: $($_.Exception.Message)"
        return $null
    }
}

# ---------------------------------------------------------------------------
# Asset selection — pick the .zip whose name matches our os/arch.
#
# Naming convention (mirrors install.ps1):
#   gitmap-<version>-<os>-<arch>.zip
# Falls back to scanning all assets if the canonical name isn't present.
# ---------------------------------------------------------------------------
function Select-Asset($release, [string]$os, [string]$arch) {
    $expected = "gitmap-$($release.tag_name)-$os-$arch.zip"
    $hit = $release.assets | Where-Object { $_.name -eq $expected } | Select-Object -First 1
    if ($hit) { return $hit }

    # Loose match: anything ending with -<os>-<arch>.zip
    $loose = $release.assets | Where-Object { $_.name -match "-$os-$arch\.(zip|tar\.gz)$" } | Select-Object -First 1
    if ($loose) {
        Write-Warn2 "Exact asset '$expected' missing; using closest match: $($loose.name)"
        return $loose
    }

    Write-Err2 "No asset matching $os/$arch in release $($release.tag_name)."
    Write-Err2 "Available assets:"
    foreach ($a in $release.assets) { Write-Err2 "  - $($a.name)" }
    exit $EXIT_UNSUPPORTED_ARCH
}

# ---------------------------------------------------------------------------
# Download + checksum
# ---------------------------------------------------------------------------
function Get-Checksums($release) {
    $cs = $release.assets | Where-Object { $_.name -eq "checksums.txt" } | Select-Object -First 1
    if (-not $cs) { return $null }
    try {
        $tmp = New-TemporaryFile
        Invoke-WebRequest -Uri $cs.browser_download_url -OutFile $tmp -UseBasicParsing
        return $tmp
    } catch {
        Write-Warn2 "Could not download checksums.txt: $($_.Exception.Message)"
        return $null
    }
}

function Test-Checksum([string]$archivePath, [string]$assetName, $checksumFile) {
    if (-not $checksumFile) {
        Write-Warn2 "No checksums.txt present; skipping verification."
        return
    }
    $line = Get-Content $checksumFile | Where-Object { $_ -match [regex]::Escape($assetName) } | Select-Object -First 1
    if (-not $line) {
        Write-Warn2 "$assetName not listed in checksums.txt; skipping verification."
        return
    }
    $expected = ($line -split '\s+')[0].ToLower()
    $actual = (Get-FileHash -Path $archivePath -Algorithm SHA256).Hash.ToLower()
    if ($expected -ne $actual) {
        Write-Err2 "Checksum mismatch for $assetName"
        Write-Err2 "  expected: $expected"
        Write-Err2 "  actual:   $actual"
        exit $EXIT_CHECKSUM
    }
    Write-OK "Checksum verified."
}

function Save-Asset($asset) {
    $tmpDir = Join-Path $env:TEMP ("gitmap-rv-" + [guid]::NewGuid().ToString("N").Substring(0,8))
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
    $out = Join-Path $tmpDir $asset.name
    Write-Step "Downloading $($asset.name) ..."
    try {
        Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $out -UseBasicParsing
    } catch {
        Write-Err2 "Download failed: $($_.Exception.Message)"
        exit $EXIT_NETWORK
    }
    return @{ Dir = $tmpDir; Path = $out }
}

# ---------------------------------------------------------------------------
# Extract + install
# ---------------------------------------------------------------------------
function Resolve-InstallDir([string]$override) {
    if (-not [string]::IsNullOrWhiteSpace($override)) { return $override }
    return (Join-Path $env:LOCALAPPDATA "gitmap\bin")
}

function Install-Binary([string]$archivePath, [string]$installDir) {
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }
    $extract = Join-Path (Split-Path $archivePath) "extract"
    New-Item -ItemType Directory -Path $extract -Force | Out-Null
    Expand-Archive -Path $archivePath -DestinationPath $extract -Force

    $candidate = Get-ChildItem -Path $extract -Recurse -File |
        Where-Object { $_.Name -eq $BinaryName -or $_.Name -match '^gitmap(-v[\d.]+)?-windows-(amd64|arm64)\.exe$' } |
        Select-Object -First 1

    if (-not $candidate) {
        Write-Err2 "Archive did not contain a gitmap binary."
        exit $EXIT_VERIFY
    }

    $dest = Join-Path $installDir $BinaryName
    Copy-Item $candidate.FullName $dest -Force
    Write-OK "Installed: $dest"
    return $dest
}

function Add-ToPath([string]$dir) {
    if ($NoPath) { return }
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -and ($userPath.Split(';') -contains $dir)) {
        Write-Step "Already on PATH: $dir"
        return
    }
    try {
        $newPath = if ([string]::IsNullOrEmpty($userPath)) { $dir } else { "$userPath;$dir" }
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-OK "Added to user PATH: $dir"
    } catch {
        Write-Warn2 "Could not update PATH: $($_.Exception.Message)"
        # PATH failure is non-fatal — binary still works via absolute path.
    }
}

function Invoke-SelfInstall([string]$binPath, [string]$expectedVersion) {
    if ($NoSelfInstall) { return }
    Write-Step "Chaining gitmap self-install ..."
    try {
        & $binPath self-install
        if ($LASTEXITCODE -ne 0) {
            Write-Warn2 "self-install exited with code $LASTEXITCODE (continuing)"
            exit $EXIT_SELF_INSTALL
        }
    } catch {
        Write-Warn2 "self-install failed: $($_.Exception.Message)"
        exit $EXIT_SELF_INSTALL
    }
}

function Test-Version([string]$binPath, [string]$expectedVersion) {
    try {
        $reported = & $binPath --version 2>&1 | Select-Object -First 1
    } catch {
        Write-Err2 "Could not run installed binary: $($_.Exception.Message)"
        exit $EXIT_VERIFY
    }
    if ($reported -notmatch [regex]::Escape($expectedVersion.TrimStart('v'))) {
        Write-Err2 "Version mismatch: expected $expectedVersion, binary reported '$reported'"
        exit $EXIT_VERIFY
    }
    Write-OK "Verified: $reported"
}

# ---------------------------------------------------------------------------
# main
# ---------------------------------------------------------------------------
try {
    if ([string]::IsNullOrWhiteSpace($Version)) {
        Write-Err2 "Required: -Version vMAJOR.MINOR.PATCH"
        Write-Err2 "Example:  .\release-version.ps1 -Version v3.36.0"
        exit $EXIT_VERSION_MISSING
    }

    $os = Resolve-OS
    $arch = Resolve-Arch $Arch
    Write-Step "Target: $os/$arch"

    $resolvedVersion = Resolve-RequestedVersion $Version
    Write-Step "Resolving release $resolvedVersion ..."
    $release = Invoke-GitHubAPI "/releases/tags/$resolvedVersion"
    if ($null -eq $release) {
        Write-Err2 "Release vanished after resolution: $resolvedVersion"
        exit $EXIT_VERSION_MISSING
    }

    $asset = Select-Asset $release $os $arch
    $checksumFile = Get-Checksums $release
    $dl = Save-Asset $asset
    Test-Checksum $dl.Path $asset.name $checksumFile

    $installDir = Resolve-InstallDir $InstallDir
    $binPath = Install-Binary $dl.Path $installDir
    Add-ToPath $installDir
    Test-Version $binPath $resolvedVersion
    Invoke-SelfInstall $binPath $resolvedVersion

    # Cleanup temp
    try { Remove-Item $dl.Dir -Recurse -Force -ErrorAction SilentlyContinue } catch {}

    Write-Host ""
    Write-OK "gitmap $resolvedVersion installed to $installDir"
    exit $EXIT_OK
}
catch {
    Write-FatalError $_ $EXIT_NETWORK
    return
}
