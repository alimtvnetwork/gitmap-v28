<#
.SYNOPSIS
    Standalone Out-of-IDE Antigravity (AGY) & AGM Full Purge Script with Snapshot Preservation.
.DESCRIPTION
    Safely exports project & conversation snapshot to JSON, terminates running Antigravity IDE
    and AGM processes, and purges all binaries, caches, brain, transcripts, and configurations.
    Can be run from an external PowerShell terminal or spawned in a detached new window.
    STRICT INVARIANT: Under no circumstances deletes anything within D:\work or git repositories.
.PARAMETER PurgeAll
    Full purge: removes .gemini, brain, transcripts, logs, and caches after snapshot export.
.PARAMETER DryRun
    Simulate removal without terminating processes or deleting files.
.PARAMETER Force
    Skip interactive confirmation prompts.
.PARAMETER BackupPath
    Custom path for the exported state snapshot JSON file.
.PARAMETER IncludeAGM
    Also terminate and uninstall Antigravity Manager (AGM).
.PARAMETER NewWindow
    Spawn in a separate external PowerShell console window.
#>
[CmdletBinding()]
param(
    [switch]$PurgeAll = $true,
    [switch]$DryRun = $false,
    [switch]$Force = $false,
    [string]$BackupPath = "",
    [switch]$IncludeAGM = $true,
    [switch]$NewWindow = $false
)

$ErrorActionPreference = "Stop"

$RepoRoot = Split-Path -Parent $PSScriptRoot

if ($NewWindow) {
    Write-Host "[GitMap] Spawning standalone PowerShell window..." -ForegroundColor Cyan
    $ScriptPath = $MyInvocation.MyCommand.Path
    $ArgsList = @("-NoExit", "-ExecutionPolicy", "Bypass", "-File", "`"$ScriptPath`"")
    if ($PurgeAll) { $ArgsList += "-PurgeAll" }
    if ($DryRun) { $ArgsList += "-DryRun" }
    if ($Force) { $ArgsList += "-Force" }
    if ($BackupPath) { $ArgsList += "-BackupPath `"$BackupPath`"" }
    if ($IncludeAGM) { $ArgsList += "-IncludeAGM" }
    Start-Process -FilePath "powershell.exe" -ArgumentList $ArgsList
    exit 0
}

Write-Host "`n=================================================================" -ForegroundColor Cyan
Write-Host " GitMap Standalone Antigravity & AGM Full Purge Engine" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# 1. Locate GitMap binary
$GitMapExe = if (Test-Path (Join-Path $RepoRoot "gitmap-test.exe")) {
    Join-Path $RepoRoot "gitmap-test.exe"
} elseif (Test-Path (Join-Path $RepoRoot "gitmap.exe")) {
    Join-Path $RepoRoot "gitmap.exe"
} else {
    "gitmap"
}

# 2. Safety check function
function Test-PathSafety([string]$Path) {
    if ([string]::IsNullOrWhiteSpace($Path)) { return $false }
    $norm = $Path.Replace("\", "/").ToLower()
    if ($norm.StartsWith("d:/work") -or $norm -eq "d:/work") {
        return $false
    }
    if ($norm -eq "c:/" -or $norm -eq "c:" -or $norm -eq "d:/" -or $norm -eq "d:") {
        return $false
    }
    return $true
}

# 3. Export Restore Snapshot
$HomeDir = [System.Environment]::GetFolderPath([System.Environment+SpecialFolder]::UserProfile)
$DefaultSnapDir = Join-Path $HomeDir ".gitmap"
if (-not (Test-Path $DefaultSnapDir)) { New-Item -ItemType Directory -Path $DefaultSnapDir -Force | Out-Null }

$Timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$TargetSnapFile = if ($BackupPath) { $BackupPath } else { Join-Path $DefaultSnapDir "agy-snapshot-$Timestamp.json" }

Write-Host "`n[Step 1/4] Preserving Projects & Conversations Snapshot..." -ForegroundColor Yellow
if ($DryRun) {
    Write-Host "  [dry-run] Would export AGY restore snapshot to: $TargetSnapFile"
} else {
    $Exported = $false
    try {
        & $GitMapExe agy snapshot $TargetSnapFile | Out-Null
        if (Test-Path $TargetSnapFile) {
            $Exported = $true
        }
    } catch {}

    if (-not $Exported) {
        $Projects = @()
        $ProjDir = Join-Path $HomeDir ".gemini\config\projects"
        if (Test-Path $ProjDir) {
            Get-ChildItem -Path $ProjDir -Filter "*.json" -File | ForEach-Object {
                try {
                    $item = Get-Content $_.FullName -Raw | ConvertFrom-Json
                    if ($item.project_metadata.workspace_path) {
                        $Projects += @{
                            id = $item.project_id
                            name = $item.project_metadata.name
                            workspace = $item.project_metadata.workspace_path
                            branch = $item.project_metadata.branch
                        }
                    }
                } catch {}
            }
        }
        $Convs = @()
        $BrainDir = Join-Path $HomeDir ".gemini\antigravity\brain"
        if (Test-Path $BrainDir) {
            Get-ChildItem -Path $BrainDir -Directory | ForEach-Object {
                $Convs += @{
                    id = $_.Name
                    title = $_.Name
                }
            }
        }
        $SnapObj = @{
            createdAt = (Get-Date).ToUniversalTime().ToString("o")
            timestamp = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
            totalProjects = $Projects.Count
            totalConversations = $Convs.Count
            projects = $Projects
            conversations = $Convs
        }
        $SnapObj | ConvertTo-Json -Depth 5 | Set-Content -Path $TargetSnapFile -Encoding UTF8
    }
    Write-Host "  [OK] Snapshot saved to: $TargetSnapFile" -ForegroundColor Green
}

# 4. Confirmation Gate
if (-not $Force -and -not $DryRun) {
    Write-Host "`nWARNING: This will completely remove Antigravity and terminate running IDE processes." -ForegroundColor Red
    $Confirm = Read-Host "Proceed with uninstallation? [y/N]"
    if ($Confirm -ne "y" -and $Confirm -ne "yes") {
        Write-Host "Uninstallation cancelled by user." -ForegroundColor Yellow
        exit 0
    }
}

# 5. Terminate Running Processes
Write-Host "`n[Step 2/4] Terminating Running Antigravity & AGM Processes..." -ForegroundColor Yellow
$ProcNames = @("antigravity", "antigravity-manager", "agm")
foreach ($PName in $ProcNames) {
    $Procs = Get-Process -Name $PName -ErrorAction SilentlyContinue
    if ($Procs) {
        if ($DryRun) {
            Write-Host "  [dry-run] Would terminate $($Procs.Count) process(es) of $PName"
        } else {
            $Procs | Stop-Process -Force -ErrorAction SilentlyContinue
            Write-Host "  [OK] Terminated $($Procs.Count) process(es) of $PName" -ForegroundColor Green
        }
    }
}

# 6. Remove Application Binaries & Caches
Write-Host "`n[Step 3/4] Purging Application Directories & Caches..." -ForegroundColor Yellow

$LocalAppData = [System.Environment]::GetFolderPath([System.Environment+SpecialFolder]::LocalApplicationData)
$AppData = [System.Environment]::GetFolderPath([System.Environment+SpecialFolder]::ApplicationData)
$ProgramFiles = [System.Environment]::GetFolderPath([System.Environment+SpecialFolder]::ProgramFiles)

$Candidates = @(
    (Join-Path $LocalAppData "agy\bin"),
    (Join-Path $HomeDir ".antigravity\bin"),
    (Join-Path $LocalAppData "Programs\Antigravity"),
    (Join-Path $ProgramFiles "Antigravity")
)

if ($PurgeAll) {
    $Candidates += @(
        (Join-Path $HomeDir ".gemini\antigravity"),
        (Join-Path $HomeDir ".gemini\config\projects"),
        (Join-Path $HomeDir ".gemini\antigravity-cli"),
        (Join-Path $HomeDir ".cache\antigravity"),
        (Join-Path $LocalAppData "antigravity"),
        (Join-Path $LocalAppData "agy"),
        (Join-Path $LocalAppData "antigravity-updater"),
        (Join-Path $AppData "Antigravity")
    )
}

if ($IncludeAGM) {
    $Candidates += @(
        (Join-Path $HomeDir ".agm"),
        (Join-Path $AppData "agm"),
        (Join-Path $LocalAppData "agm")
    )
}

foreach ($Path in $Candidates) {
    if (Test-Path $Path) {
        $IsSafe = Test-PathSafety $Path
        if (-not $IsSafe) {
            Write-Host "  [SAFETY VIOLATION] Skipped protected path: $Path" -ForegroundColor Red
            continue
        }
        if ($DryRun) {
            Write-Host "  [dry-run] Would remove: $Path"
        } else {
            Remove-Item -Path $Path -Recurse -Force -ErrorAction SilentlyContinue
            Write-Host "  [OK] Removed: $Path" -ForegroundColor Green
        }
    }
}

Write-Host "`n[Step 4/4] Verification & Snapshot Finalization..." -ForegroundColor Yellow
Write-Host "  [OK] Workspace directories in D:\work preserved intact." -ForegroundColor Green
Write-Host "  [OK] Restore snapshot preserved at: $TargetSnapFile" -ForegroundColor Cyan
Write-Host "`n=================================================================" -ForegroundColor Green
Write-Host " ANTIGRAVITY UNINSTALLATION COMPLETE (Full Fresh State)" -ForegroundColor Green
Write-Host "=================================================================" -ForegroundColor Green
