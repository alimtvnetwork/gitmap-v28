<#
.SYNOPSIS
    Uninstaller for gitmap CLI.

.DESCRIPTION
    Removes the gitmap binary, cleans up the install directory,
    and removes the directory from the user PATH.

.PARAMETER InstallDir
    Directory where gitmap is installed. Default: $env:LOCALAPPDATA\gitmap

.EXAMPLE
    irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/uninstall.ps1 | iex

.NOTES
    Repository: https://github.com/alimtvnetwork/gitmap-v27
#>

param(
    [string]$InstallDir = "",
    # -Force: skip the gitmap identity guard AND skip the keep-data
    #         prompt (defaults to keep). Use only when you know the
    #         install dir is correct.
    # -KeepData / -PurgeData: explicit data-folder choice; when both
    #         are absent the user is prompted interactively (when a
    #         tty is attached — non-interactive runs default to keep).
    [switch]$Force,
    [switch]$KeepData,
    [switch]$PurgeData
)

$ErrorActionPreference = "Stop"

$BinaryName = "gitmap.exe"
$BinaryAlias = "gm.exe"

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

# --- Resolve install directory ---

function Resolve-InstallDir([string]$dir) {
    if ($dir -ne "") { return $dir }
    return Join-Path $env:LOCALAPPDATA "gitmap"
}

# --- Remove from PATH ---

function Test-PathEntry([string]$pathValue, [string]$dir) {
    if ([string]::IsNullOrWhiteSpace($pathValue)) { return $false }
    $parts = $pathValue -split ";"
    foreach ($part in $parts) {
        if ($part.Trim() -ieq $dir) { return $true }
    }
    return $false
}

function Remove-FromPath([string]$dir) {
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")

    if (-not (Test-PathEntry $currentUserPath $dir)) {
        Write-Step "Install directory not found in PATH (nothing to remove)."
        return
    }

    $parts = ($currentUserPath -split ";") | Where-Object { $_.Trim() -ine $dir -and $_.Trim() -ne "" }
    $newPath = $parts -join ";"

    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")

    # Broadcast change
    Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

public static class GitMapUninstallNative {
    [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
        uint fuFlags, uint uTimeout, out UIntPtr lpdwResult
    );
}
"@ -ErrorAction SilentlyContinue | Out-Null

    $HWND_BROADCAST = [IntPtr]0xffff
    $WM_SETTINGCHANGE = 0x001A
    $SMTO_ABORTIFHUNG = 0x0002
    [UIntPtr]$result = [UIntPtr]::Zero

    [void][GitMapUninstallNative]::SendMessageTimeout(
        $HWND_BROADCAST, $WM_SETTINGCHANGE, [UIntPtr]::Zero,
        "Environment", $SMTO_ABORTIFHUNG, 5000, [ref]$result
    )

    # Update current session
    $machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $freshPath = @($machinePath, $userPath) | Where-Object { $_ } | ForEach-Object { $_.TrimEnd(";") }
    $env:PATH = ($freshPath -join ";")

    Write-OK "Removed from PATH."
}

# --- Safety guards ---

# Confirm-IsGitmapInstall verifies $binPath actually IS our gitmap
# CLI before any destructive action. Mirrors the helper in
# install.ps1 so both uninstall entry points share one contract.
function Confirm-IsGitmapInstall([string]$binPath, [bool]$isForce) {
    if (-not (Test-Path $binPath)) {
        Write-Step "No binary at $binPath; skipping identity check."
        return $true
    }
    try { $out = (& $binPath version 2>&1 | Out-String) } catch { $out = "" }
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

# Resolve-DataChoice returns 'keep' or 'purge'. Same precedence as
# install.ps1: explicit flags > -Force (keep) > non-interactive (keep)
# > interactive prompt.
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

# --- Remove files ---

function Remove-InstallDir([string]$dir, [string]$dataChoice) {
    if (-not (Test-Path $dir)) {
        Write-Err "Install directory not found: $dir"
        return $false
    }

    $binPath = Join-Path $dir $BinaryName
    if (Test-Path $binPath) {
        Remove-Item $binPath -Force
        Write-OK "Removed $BinaryName"
    }

    # Also remove the `gm.exe` alias copy installed by install.ps1 so the
    # uninstall is symmetric and doesn't leave a stale short command on PATH.
    $aliasPath = Join-Path $dir $BinaryAlias
    if (Test-Path $aliasPath) {
        Remove-Item $aliasPath -Force -ErrorAction SilentlyContinue
        Write-OK "Removed $BinaryAlias"
    }

    $dataDir = Join-Path $dir "data"
    if ($dataChoice -eq 'purge' -and (Test-Path $dataDir)) {
        Remove-Item $dataDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-OK "Removed data folder: $dataDir"
    }
    elseif (Test-Path $dataDir) {
        Write-Step "Kept data folder: $dataDir"
    }

    # Remove remaining files (old binaries, etc.) iff the dir is empty.
    $remaining = Get-ChildItem -Path $dir -Force -ErrorAction SilentlyContinue
    if ($remaining.Count -eq 0) {
        Remove-Item $dir -Force
        Write-OK "Removed install directory: $dir"
    }
    else {
        Write-Step "Directory not empty, keeping: $dir"
        Write-Step "  Remaining: $($remaining.Name -join ', ')"
    }

    return $true
}

# --- Main ---

function Main {
    Write-Host ""
    Write-Host "  gitmap uninstaller" -ForegroundColor White
    Write-Host ""

    $resolvedDir = Resolve-InstallDir $InstallDir
    $binPath = Join-Path $resolvedDir $BinaryName
    $dataDir = Join-Path $resolvedDir "data"

    Write-Step "Uninstalling from $resolvedDir..."

    if (-not (Confirm-IsGitmapInstall $binPath $Force.IsPresent)) {
        Write-Host ""
        return
    }

    $dataChoice = Resolve-DataChoice $dataDir $KeepData.IsPresent $PurgeData.IsPresent $Force.IsPresent

    $removed = Remove-InstallDir $resolvedDir $dataChoice
    Remove-FromPath $resolvedDir

    Write-Host ""
    if ($removed) {
        Write-OK "gitmap has been uninstalled."
    }
    else {
        Write-Err "gitmap was not found at $resolvedDir."
    }
    Write-Host ""
}

Main

