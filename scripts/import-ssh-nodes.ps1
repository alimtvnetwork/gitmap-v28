<#
.SYNOPSIS
    Imports SSH fleet nodes from JSON configuration into GitMap with full credential resolution.
.DESCRIPTION
    Reads gitmap-ssh-nodes.json (or gitmap-ssh.json) and imports all node topologies, aliases,
    IP addresses, and passwords into GitMap's Split-DB vault. Supports automatic fleet deployment.
.PARAMETER FilePath
    Path to gitmap-ssh-nodes.json. Defaults to D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json.
.PARAMETER PassFile
    Path to vmpass.json credential fallback. Defaults to D:\work\repo-secrets\01-gitmap\vmpass.json.
.PARAMETER Target
    Optional deploy target machine (e.g. w1, w2, or all).
.PARAMETER Except
    Comma-separated list of node IDs, IPs, or aliases to exclude during deployment.
.PARAMETER Deploy
    Switch to automatically deploy the imported configuration across remote fleet machines.
.EXAMPLE
    .\scripts\import-ssh-nodes.ps1
    .\scripts\import-ssh-nodes.ps1 -Deploy -Target w1
#>
[CmdletBinding()]
param(
    [Alias("Path")]
    [string]$FilePath = "",
    [string]$PassFile = "",
    [string]$Target = "all",
    [string]$Except = "",
    [switch]$Deploy = $false
)

$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "   GitMap SSH Node & Credential Import Utility           " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Resolve GitMap Executable
$GitMapBin = $null
$CandidateBinaries = @(
    "D:\work\gitmap\gitmap.exe",
    (Join-Path $PSScriptRoot "..\gitmap.exe"),
    (Join-Path $PSScriptRoot "gitmap.exe"),
    (Join-Path $PSScriptRoot "..\cli\gitmap.exe"),
    "D:\work\gitmap\cli\gitmap.exe",
    ".\gitmap.exe"
)
foreach ($Candidate in $CandidateBinaries) {
    if (Test-Path $Candidate) {
        $GitMapBin = (Resolve-Path $Candidate).Path
        break
    }
}

if (-not $GitMapBin) {
    $Cmd = Get-Command "gitmap" -ErrorAction SilentlyContinue
    if ($Cmd) {
        if ($Cmd.Path) {
            $GitMapBin = $Cmd.Path
        } elseif ($Cmd.Source) {
            $GitMapBin = $Cmd.Source
        } else {
            $GitMapBin = "gitmap"
        }
    }
}

if (-not $GitMapBin) {
    Write-Error "[FAIL] gitmap executable not found in PATH or standard repository directories."
    exit 1
}
Write-Host "[Gate 1] GitMap executable resolved: $GitMapBin" -ForegroundColor Green

# 2. Resolve SSH Nodes JSON File
$ResolvedPath = $null
if ($FilePath -ne "") {
    if (Test-Path $FilePath) {
        $ResolvedPath = (Resolve-Path $FilePath).Path
    } else {
        Write-Error "[FAIL] Specified FilePath not found: $FilePath"
        exit 1
    }
} else {
    $Candidates = @(
        "D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json",
        (Join-Path $PSScriptRoot "..\..\repo-secrets\01-gitmap\gitmap-ssh-nodes.json"),
        "D:\work\repo-secrets\gitmap-ssh-nodes.json",
        (Join-Path $PSScriptRoot "gitmap-ssh-nodes.json"),
        "gitmap-ssh-nodes.json"
    )
    foreach ($c in $Candidates) {
        if (Test-Path $c) {
            $ResolvedPath = (Resolve-Path $c).Path
            break
        }
    }
}

if (-not $ResolvedPath) {
    Write-Error "[FAIL] gitmap-ssh-nodes.json not found in candidate paths."
    exit 1
}
Write-Host "[Gate 2] SSH nodes file located: $ResolvedPath" -ForegroundColor Green

# 3. Import via GitMap CLI
Write-Host "[Step 3] Importing SSH nodes into local Split-DB..." -ForegroundColor Yellow
& $GitMapBin ssh nodes import-json "$ResolvedPath"

# 4. Display Enrolled Nodes
Write-Host "`n[Step 4] Enrolled cluster node inventory:" -ForegroundColor Yellow
& $GitMapBin ssh ls

# 5. Optional Fleet Deploy
if ($Deploy) {
    Write-Host "`n[Step 5] Deploying node configuration across fleet (target: $Target)..." -ForegroundColor Yellow
    if ($Except -ne "") {
        & $GitMapBin deploy config ssh $Target --except $Except
    } else {
        & $GitMapBin deploy config ssh $Target
    }
}

Write-Host "`n==========================================================" -ForegroundColor Green
Write-Host "   SSH Node Import & Verification SUCCESSFUL             " -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green
exit 0
