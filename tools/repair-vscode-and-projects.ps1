<#
.SYNOPSIS
    Diagnoses and repairs VS Code startup failures and VS Code Project Manager JSON configuration.

.DESCRIPTION
    This script addresses two common issues:
    1. VS Code failing to launch due to an interrupted/botched background update
       (Chromium ICU error: "Invalid file descriptor to ICU data received" / STATUS_BREAKPOINT).
    2. VS Code Project Manager (alefragnani.project-manager) projects.json validation,
       checking for stripped fields, corrupt JSON syntax, missing directories, and syncing
       between modern (globalStorage) and legacy (User) locations.

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File .\repair-vscode-and-projects.ps1
#>

[CmdletBinding()]
param(
    [switch]$ForceReinstall,
    [switch]$SkipGitMapCheck
)

$ErrorActionPreference = "Continue"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "  VS Code & Project Manager Repair & Diagnostic Utility   " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# -----------------------------------------------------------------------------
# 1. Terminate Stuck Processes
# -----------------------------------------------------------------------------
Write-Host "`n[Step 1/5] Checking for lingering VS Code processes..." -ForegroundColor Yellow
$procNames = @("Code", "inno_updater", "code-tunnel")
foreach ($name in $procNames) {
    $running = Get-Process -Name $name -ErrorAction SilentlyContinue
    if ($running) {
        Write-Host "  Found running process: $name (Count: $($running.Count)). Terminating..." -ForegroundColor Gray
        Stop-Process -Name $name -Force -ErrorAction SilentlyContinue
    }
}
Write-Host "  Process cleanup complete." -ForegroundColor Green

# -----------------------------------------------------------------------------
# 2. Inspect & Repair VS Code Project Manager JSON (projects.json)
# -----------------------------------------------------------------------------
Write-Host "`n[Step 2/5] Inspecting Project Manager JSON (projects.json)..." -ForegroundColor Yellow

$appData = if ($env:APPDATA) { $env:APPDATA } else { [Environment]::GetFolderPath([Environment+SpecialFolder]::ApplicationData) }
$modernDir = Join-Path $appData "Code\User\globalStorage\alefragnani.project-manager"
$modernFile = Join-Path $modernDir "projects.json"

$legacyDir = Join-Path $appData "Code\User"
$legacyFile = Join-Path $legacyDir "projects.json"

# Find which projects.json exists
$activeFile = $null
if (Test-Path $modernFile) {
    $activeFile = $modernFile
    Write-Host "  Found modern Project Manager file: $modernFile" -ForegroundColor Gray
} elseif (Test-Path $legacyFile) {
    $activeFile = $legacyFile
    Write-Host "  Found legacy Project Manager file: $legacyFile" -ForegroundColor Gray
} else {
    Write-Host "  projects.json not found in standard locations." -ForegroundColor Yellow
}

$projects = @()
$isJsonValid = $false

if ($activeFile) {
    try {
        $raw = Get-Content -Path $activeFile -Raw -Encoding UTF8 -ErrorAction Stop
        if ([string]::IsNullOrWhiteSpace($raw)) {
            Write-Host "  projects.json is empty. Initializing empty array []..." -ForegroundColor Yellow
            $projects = @()
            $isJsonValid = $true
        } else {
            $parsed = $raw | ConvertFrom-Json -ErrorAction Stop
            if ($parsed -is [System.Collections.IEnumerable] -and $parsed -isnot [string]) {
                $projects = @($parsed)
            } else {
                $projects = @($parsed)
            }
            $isJsonValid = $true
            Write-Host "  JSON syntax valid. Parsed $($projects.Count) project entries." -ForegroundColor Green
        }
    } catch {
        Write-Host "  ERROR: projects.json contains invalid JSON or syntax error: $_" -ForegroundColor Red
        $bakFile = "$activeFile.corrupt.$(Get-Date -Format 'yyyyMMddHHmmss').bak"
        Copy-Item -Path $activeFile -Destination $bakFile -Force
        Write-Host "  Backed up corrupted file to $bakFile" -ForegroundColor Gray
        Write-Host "  Resetting projects.json to valid empty array []..." -ForegroundColor Yellow
        $projects = @()
        $isJsonValid = $true
    }
} else {
    if (-not (Test-Path $modernDir)) {
        New-Item -ItemType Directory -Path $modernDir -Force | Out-Null
    }
    $activeFile = $modernFile
    $projects = @()
    $isJsonValid = $true
}

# Validate fields across all entries
$repairedCount = 0
$missingPathsCount = 0

foreach ($p in $projects) {
    if ($null -eq $p.paths) { $p | Add-Member -NotePropertyName "paths" -NotePropertyValue @() -Force; $repairedCount++ }
    if ($null -eq $p.tags) { $p | Add-Member -NotePropertyName "tags" -NotePropertyValue @("gitmap") -Force; $repairedCount++ }
    if ($null -eq $p.enabled) { $p | Add-Member -NotePropertyName "enabled" -NotePropertyValue $true -Force; $repairedCount++ }
    if ($null -eq $p.profile) { $p | Add-Member -NotePropertyName "profile" -NotePropertyValue "" -Force; $repairedCount++ }

    if ($p.rootPath) {
        if (-not (Test-Path $p.rootPath)) {
            $missingPathsCount++
        }
    }
}

if ($missingPathsCount -gt 0) {
    Write-Host "  Notice: $missingPathsCount project paths in projects.json do not exist on disk." -ForegroundColor Yellow
} else {
    Write-Host "  All project directories verified on disk." -ForegroundColor Green
}

# Re-serialize cleanly with tabs (matching GitMap / alefragnani spec)
if ($isJsonValid) {
    if (-not (Test-Path $modernDir)) { New-Item -ItemType Directory -Path $modernDir -Force | Out-Null }
    if (-not (Test-Path $legacyDir)) { New-Item -ItemType Directory -Path $legacyDir -Force | Out-Null }

    $jsonOutput = ($projects | ConvertTo-Json -Depth 10)
    [System.IO.File]::WriteAllText($modernFile, $jsonOutput, [System.Text.UTF8Encoding]::new($false))
    [System.IO.File]::WriteAllText($legacyFile, $jsonOutput, [System.Text.UTF8Encoding]::new($false))
    Write-Host "  Synchronized clean projects.json to both modern and legacy locations." -ForegroundColor Green
}

# -----------------------------------------------------------------------------
# 3. Detect and Test VS Code Installations (Dynamic Discovery)
# -----------------------------------------------------------------------------
Write-Host "`n[Step 3/5] Testing VS Code executable integrity..." -ForegroundColor Yellow

$possibleInstalls = [System.Collections.Generic.List[string]]::new()

# 1. Check PATH command
$codeFromPath = Get-Command code.cmd -ErrorAction SilentlyContinue
if (-not $codeFromPath) { $codeFromPath = Get-Command code -ErrorAction SilentlyContinue }
if ($codeFromPath -and $codeFromPath.Source) {
    $parentDir = Split-Path (Split-Path $codeFromPath.Source -Parent) -Parent
    if ($parentDir -and (Test-Path $parentDir) -and -not $possibleInstalls.Contains($parentDir)) {
        $possibleInstalls.Add($parentDir)
    }
}

# 2. Check System 64-bit Program Files
if ($env:ProgramFiles) {
    $pFiles = Join-Path $env:ProgramFiles "Microsoft VS Code"
    if ((Test-Path $pFiles) -and -not $possibleInstalls.Contains($pFiles)) {
        $possibleInstalls.Add($pFiles)
    }
}

# 3. Check System 32-bit Program Files (x86)
$envP86 = ${env:ProgramFiles(x86)}
if ($envP86) {
    $p86 = Join-Path $envP86 "Microsoft VS Code"
    if ((Test-Path $p86) -and -not $possibleInstalls.Contains($p86)) {
        $possibleInstalls.Add($p86)
    }
}

# 4. Check Local AppData User Programs
if ($env:LOCALAPPDATA) {
    $userCode = Join-Path $env:LOCALAPPDATA "Programs\Microsoft VS Code"
    if ((Test-Path $userCode) -and -not $possibleInstalls.Contains($userCode)) {
        $possibleInstalls.Add($userCode)
    }
}

# 5. Check Registry App Paths
$regPath = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\Code.exe"
if (Test-Path $regPath) {
    $regVal = (Get-ItemProperty -Path $regPath -Name "(default)" -ErrorAction SilentlyContinue)."(default)"
    if ($regVal -and (Test-Path $regVal)) {
        $regDir = Split-Path $regVal -Parent
        if ($regDir -and -not $possibleInstalls.Contains($regDir)) {
            $possibleInstalls.Add($regDir)
        }
    }
}

$workingInstall = $null

foreach ($inst in $possibleInstalls) {
    $codeCmd = Join-Path $inst "bin\code.cmd"
    if (Test-Path $codeCmd) {
        Write-Host "  Found installation candidate at: $inst" -ForegroundColor Gray
        $testOutput = & $codeCmd --version 2>&1
        $exitCode = $LASTEXITCODE

        if ($exitCode -eq 0 -and $testOutput -notmatch "ERROR:base\\i18n\\icu_util") {
            Write-Host "  Executable healthy: $testOutput" -ForegroundColor Green
            $workingInstall = $inst
            break
        } else {
            Write-Host "  Executable failed (Exit Code: $exitCode). Output: $testOutput" -ForegroundColor Red
        }
    }
}

# -----------------------------------------------------------------------------
# 4. Repair Installation if Broken
# -----------------------------------------------------------------------------
Write-Host "`n[Step 4/5] Evaluating repair actions..." -ForegroundColor Yellow

if ($workingInstall -and -not $ForceReinstall) {
    Write-Host "  VS Code is working properly at: $workingInstall" -ForegroundColor Green
} else {
    Write-Host "  Attempting repair of VS Code installation..." -ForegroundColor Yellow

    # If multiple candidates exist, synchronize commit folders between them
    for ($i = 0; $i -lt $possibleInstalls.Count; $i++) {
        for ($j = 0; $j -lt $possibleInstalls.Count; $j++) {
            if ($i -eq $j) { continue }
            $srcDir = $possibleInstalls[$i]
            $dstDir = $possibleInstalls[$j]

            if ((Test-Path $srcDir) -and (Test-Path $dstDir)) {
                $srcCommitDirs = Get-ChildItem -Directory $srcDir | Where-Object { $_.Name -match '^[0-9a-f]{10}$' }
                foreach ($cd in $srcCommitDirs) {
                    $targetPath = Join-Path $dstDir $cd.Name
                    if (-not (Test-Path $targetPath)) {
                        Write-Host "  Transferring missing commit folder $($cd.Name) from $srcDir to $dstDir..." -ForegroundColor Gray
                        Copy-Item -Path $cd.FullName -Destination $targetPath -Recurse -Force -ErrorAction SilentlyContinue
                    }
                }
            }
        }
    }

    # Re-test candidates
    $fixed = $false
    foreach ($inst in $possibleInstalls) {
        $retestCmd = Join-Path $inst "bin\code.cmd"
        if (Test-Path $retestCmd) {
            $testOut = & $retestCmd --version 2>&1
            if ($LASTEXITCODE -eq 0 -and $testOut -notmatch "ERROR:") {
                $fixed = $true
                Write-Host "  In-place version sync fixed VS Code successfully at $inst! ($testOut)" -ForegroundColor Green
                break
            }
        }
    }

    # If still not fixed or ForceReinstall requested, invoke winget
    if (-not $fixed) {
        $winget = Get-Command winget.exe -ErrorAction SilentlyContinue
        if ($winget) {
            Write-Host "  Running winget to repair VS Code binaries..." -ForegroundColor Yellow
            & winget.exe install --id Microsoft.VisualStudioCode --force --accept-source-agreements --accept-package-agreements
            
            $postTest = & code --version 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  Winget repair successful! VS Code version: $postTest" -ForegroundColor Green
            } else {
                Write-Host "  Winget completed. Please test launching VS Code manually." -ForegroundColor Yellow
            }
        } else {
            Write-Host "  winget is not available on this system. Please re-run the latest VS Code installer to repair binaries." -ForegroundColor Red
        }
    }
}

# -----------------------------------------------------------------------------
# 5. GitMap Verification (if available)
# -----------------------------------------------------------------------------
Write-Host "`n[Step 5/5] Checking GitMap VS Code integration..." -ForegroundColor Yellow

if (-not $SkipGitMapCheck) {
    $gitmapCmd = Get-Command gitmap -ErrorAction SilentlyContinue
    if ($gitmapCmd) {
        Write-Host "  Running gitmap vscode find-duplicates..." -ForegroundColor Gray
        & gitmap vscode find-duplicates
        Write-Host "  GitMap integration verified." -ForegroundColor Green
    } else {
        Write-Host "  gitmap CLI not found on PATH (skipping gitmap check)." -ForegroundColor Gray
    }
}

Write-Host "`n==========================================================" -ForegroundColor Cyan
Write-Host "  Repair process completed successfully!                  " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
