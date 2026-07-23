<#
.SYNOPSIS
    One-liner uninstaller for gitmap-v28 on Windows.

.DESCRIPTION
    Removes the gitmap-v28 binary, deploy folder, user PATH entry, and (optionally)
    the per-user data folder. Works whether gitmap-v28 was installed via:

      - install-quick.ps1 (one-liner)
      - gitmap-v28/scripts/install.ps1 (canonical installer)
      - manual `run.ps1` build-and-deploy

    Strategy:
      1. If `gitmap-v28` is on PATH and reports its install dir, prefer that.
         Delegate to `gitmap-v28 self-uninstall -y` (best path — uses the binary's
         own knowledge of marker-block PATH cleanup, scheduled-task removal,
         etc.).
      2. If `gitmap-v28` is NOT on PATH (already partially removed, broken install),
         fall back to a manual sweep:
           - delete <root>/gitmap-cli/ AND legacy <root>/gitmap-v28/
           - strip the deploy root from User PATH
           - prompt before deleting %APPDATA%/gitmap-v28

    Run via one-liner:
      irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.ps1 | iex

    Or locally:
      ./uninstall-quick.ps1
      ./uninstall-quick.ps1 -KeepData
      ./uninstall-quick.ps1 -InstallDir "E:\bin-run"
      ./uninstall-quick.ps1 -Yes
#>

param(
    [string]$InstallDir = "",
    [switch]$KeepData,
    [switch]$Yes
)

$ErrorActionPreference = "Continue"
$ProgressPreference    = "SilentlyContinue"

function Write-Step($msg) { Write-Host "  $msg" -ForegroundColor Cyan }
function Write-Info($msg) { Write-Host "    $msg" -ForegroundColor DarkGray }
function Write-Ok($msg)   { Write-Host "    $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    $msg" -ForegroundColor Yellow }
function Write-Err($msg)  { Write-Host "    $msg" -ForegroundColor Red }

function Confirm-Or-Exit([string]$prompt) {
    if ($Yes) { return $true }
    Write-Host ""
    Write-Host "  $prompt [y/N]: " -ForegroundColor Yellow -NoNewline
    $answer = Read-Host
    return ($answer -match '^(y|yes)$')
}

# ---------------------------------------------------------------------------
# Step 1 — try the canonical `gitmap-v28 self-uninstall` (best path).
# ---------------------------------------------------------------------------

function Try-SelfUninstall {
    # Get-Command can return MULTIPLE ApplicationInfo entries when stale
    # binaries from a prior install (e.g. drive-root shim + nested
    # gitmap-cli/) are both on PATH. Interpolating the array into a
    # string would join the .Source paths with a space and PowerShell
    # would later treat that joined string as a command name —
    # producing the error:
    #   The term 'E:\gitmap-v28\gitmap.exe E:\bin-run\gitmap-cli\gitmap.exe'
    #   is not recognized as a name of a cmdlet ...
    # Always pick the first entry (the one PATH would actually invoke).
    $cmd = Get-Command gitmap-v28 -ErrorAction SilentlyContinue | Select-Object -First 1
    if (-not $cmd) {
        Write-Info "gitmap-v28 not found on PATH, skipping self-uninstall (will sweep manually)"
        return $false
    }

    $activeBinary = [string]$cmd.Source
    if ([string]::IsNullOrWhiteSpace($activeBinary)) {
        Write-Info "Active binary path was empty; skipping self-uninstall (will sweep manually)"
        return $false
    }

    Write-Info "Active binary: $activeBinary"
    Write-Info "Delegating to: $activeBinary self-uninstall -y"
    Write-Host ""
    try {
        # Invoke by absolute path (not by 'gitmap-v28') so PATH-resolution
        # quirks can't pick a different binary than the one we logged.
        & $activeBinary self-uninstall -y
        if ($LASTEXITCODE -eq 0) {
            Write-Ok "self-uninstall completed cleanly"
            return $true
        }
        Write-Warn "self-uninstall exited with code $LASTEXITCODE; falling back to manual sweep"
        return $false
    } catch {
        Write-Warn "self-uninstall threw: $_; falling back to manual sweep"
        return $false
    }
}

# ---------------------------------------------------------------------------
# Step 2 — manual sweep fallback.
# ---------------------------------------------------------------------------

function Resolve-DeployRoot {
    if ($InstallDir.Length -gt 0) { return $InstallDir }

    # Check the active binary's grandparent (deployRoot/gitmap-cli/gitmap.exe).
    # Same multi-binary defense as Try-SelfUninstall: pick the FIRST
    # entry so a stale shim doesn't poison Split-Path.
    $cmd = Get-Command gitmap-v28 -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($cmd -and (Test-Path $cmd.Source)) {
        $parent = Split-Path (Resolve-Path $cmd.Source).Path -Parent
        $grand  = Split-Path $parent -Parent
        return $grand
    }

    # Common defaults to probe.
    foreach ($candidate in @("E:\bin-run", "D:\gitmap-v28", "$env:LOCALAPPDATA\gitmap-v28", "$env:USERPROFILE\gitmap-v28")) {
        if (Test-Path (Join-Path $candidate "gitmap-cli\gitmap.exe")) { return $candidate }
        if (Test-Path (Join-Path $candidate "gitmap-v28\gitmap.exe"))     { return $candidate }
        if (Test-Path (Join-Path $candidate "gitmap.exe"))            { return $candidate }
    }
    return ""
}

function Remove-DeployFolders([string]$root) {
    if (-not $root) {
        Write-Warn "could not locate a gitmap-v28 deploy root; skipping folder removal"
        return
    }

    foreach ($sub in @("gitmap-cli", "gitmap-v28")) {
        $dir = Join-Path $root $sub
        if (Test-Path $dir) {
            try {
                Remove-Item $dir -Recurse -Force -ErrorAction Stop
                Write-Ok "removed $dir"
            } catch {
                Write-Err "could not remove $dir : $_"
                Write-Warn "close any terminal/process using gitmap.exe and retry"
            }
        }
    }

    # Also remove a flat deploy (no subfolder).
    $flatBin = Join-Path $root "gitmap.exe"
    if (Test-Path $flatBin) {
        try {
            Remove-Item $flatBin -Force -ErrorAction Stop
            Write-Ok "removed $flatBin"
        } catch {
            Write-Err "could not remove $flatBin : $_"
        }
    }
}

function Remove-FromUserPath([string]$root) {
    if (-not $root) { return }

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not $userPath) { return }

    $entries = $userPath -split ';' | Where-Object { $_ -ne "" }
    $toRemove = @(
        $root,
        (Join-Path $root "gitmap-cli"),
        (Join-Path $root "gitmap-v28")
    )

    $kept = $entries | Where-Object {
        $entry = $_.TrimEnd('\')
        $match = $false
        foreach ($r in $toRemove) {
            if ($entry -ieq $r.TrimEnd('\')) { $match = $true; break }
        }
        -not $match
    }

    $newPath = ($kept -join ';')
    if ($newPath -ne $userPath) {
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Ok "stripped gitmap-v28 entries from User PATH (restart shells to take effect)"
    } else {
        Write-Info "no gitmap-v28 entries found in User PATH"
    }
}

function Remove-DataFolder {
    $appdata = Join-Path $env:APPDATA "gitmap-v28"
    if (-not (Test-Path $appdata)) { return }

    if ($KeepData) {
        Write-Info "keeping data folder: $appdata"
        return
    }

    if (-not (Confirm-Or-Exit "Also delete user data at $appdata?")) {
        Write-Info "kept: $appdata"
        return
    }

    try {
        Remove-Item $appdata -Recurse -Force -ErrorAction Stop
        Write-Ok "removed $appdata"
    } catch {
        Write-Err "could not remove $appdata : $_"
    }
}

# ---------------------------------------------------------------------------
# Exhaustive PATH sweep — find EVERY gitmap.exe still reachable and nuke it.
# This is the safety net that catches stale binaries the canonical
# self-uninstall and deploy-folder sweep can miss (drive-root copies,
# manually-placed shims, old installs in unusual paths, etc.).
# ---------------------------------------------------------------------------

function Get-AllGitmapOnPath {
    # Get-Command -All returns every match across PATH, not just the first.
    $found = @()
    try {
        $cmds = Get-Command gitmap-v28 -All -ErrorAction SilentlyContinue
        foreach ($c in $cmds) {
            if ($c.Source -and (Test-Path $c.Source)) {
                $found += (Resolve-Path $c.Source).Path
            }
        }
    } catch {}

    # Belt-and-suspenders: also scan every PATH entry directly for a
    # gitmap.exe / gitmap-v28 file in case Get-Command missed something.
    $pathDirs = @()
    foreach ($scope in @("Machine", "User")) {
        $p = [Environment]::GetEnvironmentVariable("Path", $scope)
        if ($p) { $pathDirs += ($p -split ';') }
    }
    $pathDirs += ($env:PATH -split ';')

    foreach ($d in $pathDirs) {
        if (-not $d) { continue }
        $d = $d.TrimEnd('\')
        foreach ($name in @("gitmap.exe", "gitmap-v28")) {
            $candidate = Join-Path $d $name
            if (Test-Path $candidate) {
                $resolved = (Resolve-Path $candidate).Path
                if ($found -notcontains $resolved) { $found += $resolved }
            }
        }
    }

    return $found | Select-Object -Unique
}

function Remove-StrayBinaries {
    $all = Get-AllGitmapOnPath
    if (-not $all -or $all.Count -eq 0) {
        Write-Info "no stray gitmap-v28 binaries found on PATH"
        return @()
    }

    Write-Info "found $($all.Count) gitmap-v28 binary location(s):"
    foreach ($b in $all) { Write-Info "  - $b" }

    $cleanedDirs = @()
    foreach ($bin in $all) {
        try {
            Remove-Item $bin -Force -ErrorAction Stop
            Write-Ok "removed $bin"
            $cleanedDirs += (Split-Path $bin -Parent)
        } catch {
            Write-Err "could not remove $bin : $_"
            Write-Warn "close any terminal/process using gitmap.exe and retry"
        }
    }
    return $cleanedDirs | Select-Object -Unique
}

function Remove-DirsFromUserPath([string[]]$dirs) {
    if (-not $dirs -or $dirs.Count -eq 0) { return }

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not $userPath) { return }

    $entries = $userPath -split ';' | Where-Object { $_ -ne "" }
    $normTargets = $dirs | ForEach-Object { $_.TrimEnd('\').ToLower() }

    $kept = $entries | Where-Object {
        $entry = $_.TrimEnd('\').ToLower()
        -not ($normTargets -contains $entry)
    }

    $newPath = ($kept -join ';')
    if ($newPath -ne $userPath) {
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Ok "stripped stray gitmap-v28 dirs from User PATH"
    }
}

function Remove-CompletionSourceLines {
    $userHomeDir = $env:USERPROFILE
    $profiles = @(
        (Join-Path $userHomeDir "Documents\PowerShell\profile.ps1"),
        (Join-Path $userHomeDir "Documents\PowerShell\Microsoft.PowerShell_profile.ps1"),
        (Join-Path $userHomeDir "Documents\WindowsPowerShell\profile.ps1"),
        (Join-Path $userHomeDir "Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1")
    )

    foreach ($p in $profiles) {
        if (-not (Test-Path $p)) { continue }
        $lines = Get-Content $p -Raw
        if (-not $lines) { continue }

        $cleaned = ($lines -split "`n" | Where-Object {
            $trimmed = $_.Trim()
            -not ($trimmed -eq '# gitmap-v28 shell completion') -and
            -not ($trimmed -match "^\.\s+'.*completions\.ps1'")
        }) -join "`n"

        if ($cleaned -ne $lines) {
            Set-Content -Path $p -Value $cleaned -NoNewline -Encoding UTF8
            Write-Ok "Removed completion source line from $p"
        }
    }
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

Write-Host ""
Write-Host "  gitmap-v28 quick uninstaller" -ForegroundColor Cyan
Write-Host "  ------------------------" -ForegroundColor DarkGray
Write-Host ""

Write-Step "[1/4] Trying canonical self-uninstall"
$ok = Try-SelfUninstall

if (-not $ok) {
    Write-Host ""
    Write-Step "[2/4] Manual sweep — locating deploy root"
    $root = Resolve-DeployRoot
    if ($root) { Write-Info "Deploy root: $root" } else { Write-Warn "no deploy root found" }

    Write-Host ""
    Write-Step "[3/4] Removing deploy folders"
    Remove-DeployFolders $root

    Write-Host ""
    Write-Step "[4/4] Cleaning User PATH"
    Remove-FromUserPath $root
}

Write-Host ""
Write-Step "Cleaning PowerShell profile completion lines"
Remove-CompletionSourceLines

Write-Host ""
Write-Step "Exhaustive PATH sweep — removing any remaining gitmap-v28 binaries"
$strayDirs = Remove-StrayBinaries
Remove-DirsFromUserPath $strayDirs

Write-Host ""
Write-Step "User data"
Remove-DataFolder

Write-Host ""
Write-Host "  Done. Open a new terminal to refresh PATH." -ForegroundColor Green
Write-Host ""
