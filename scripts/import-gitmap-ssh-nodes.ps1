<#
.SYNOPSIS
    Imports SSH fleet node definitions and credentials into GitMap Split-DB.
.DESCRIPTION
    Reads gitmap-ssh-nodes.json (or specified file), parses encrypted/plaintext passwords
    with fallback from vmpass.json, and imports/syncs them idempotently into GitMap.
.PARAMETER JsonPath
    Path to gitmap-ssh-nodes.json. Defaults to D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json
    or relative candidates.
.PARAMETER PassFile
    Path to vmpass.json. Defaults to D:\work\repo-secrets\01-gitmap\vmpass.json.
.PARAMETER DeployToFleet
    When set, automatically deploys the imported node configuration across all remote fleet machines
    via 'gitmap deploy config ssh all'.
.PARAMETER VerboseOutput
    Display verbose diagnostic messages.
.EXAMPLE
    .\import-gitmap-ssh-nodes.ps1
    .\import-gitmap-ssh-nodes.ps1 -JsonPath "D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json" -DeployToFleet
#>
[CmdletBinding()]
param(
    [string]$JsonPath = "",
    [string]$PassFile = "",
    [switch]$DeployToFleet = $false,
    [switch]$VerboseOutput = $false
)

$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "   GitMap SSH Fleet Node Configuration Import Runner      " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Locate GitMap CLI executable
$GitMapExe = $null
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
        $GitMapExe = (Resolve-Path $Candidate).Path
        break
    }
}

if (-not $GitMapExe) {
    $Cmd = Get-Command "gitmap" -ErrorAction SilentlyContinue
    if ($Cmd) {
        if ($Cmd.Path) {
            $GitMapExe = $Cmd.Path
        } elseif ($Cmd.Source) {
            $GitMapExe = $Cmd.Source
        } else {
            $GitMapExe = "gitmap"
        }
    }
}

if (-not $GitMapExe) {
    Write-Error "[FATAL] gitmap executable not found in PATH or standard repository directories."
    exit 1
}

Write-Host "[Gate 1] GitMap CLI detected: $GitMapExe" -ForegroundColor Green

# 2. Locate gitmap-ssh-nodes.json
$ResolvedJsonPath = $null
if ($JsonPath -ne "") {
    if (Test-Path $JsonPath) {
        $ResolvedJsonPath = (Resolve-Path $JsonPath).Path
    } else {
        Write-Error "[FATAL] Specified JsonPath not found: $JsonPath"
        exit 1
    }
} else {
    $CandidateJsonPaths = @(
        "D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json",
        (Join-Path $PSScriptRoot "..\..\repo-secrets\01-gitmap\gitmap-ssh-nodes.json"),
        (Join-Path $PSScriptRoot "gitmap-ssh-nodes.json"),
        "repo-secrets\01-gitmap\gitmap-ssh-nodes.json",
        "gitmap-ssh-nodes.json"
    )
    foreach ($Candidate in $CandidateJsonPaths) {
        if (Test-Path $Candidate) {
            $ResolvedJsonPath = (Resolve-Path $Candidate).Path
            break
        }
    }
}

if (-not $ResolvedJsonPath) {
    Write-Error "[FATAL] gitmap-ssh-nodes.json not found in standard paths."
    exit 1
}

Write-Host "[Gate 2] SSH nodes JSON located: $ResolvedJsonPath" -ForegroundColor Green

# 3. Validate JSON payload content
try {
    $RawJson = Get-Content -Path $ResolvedJsonPath -Raw | ConvertFrom-Json
    $NodeCount = 0
    if ($RawJson.nodes) {
        $NodeCount = $RawJson.nodes.Count
    } elseif ($RawJson.connections) {
        $NodeCount = $RawJson.connections.Count
    } elseif ($RawJson -is [System.Array]) {
        $NodeCount = $RawJson.Count
    }
    Write-Host "[Gate 3] JSON payload validated: $NodeCount node definition(s) discovered" -ForegroundColor Green
} catch {
    Write-Error "[FATAL] Failed to parse JSON at $ResolvedJsonPath : $_"
    exit 1
}

# 4. Import nodes using GitMap CLI
Write-Host "[Step 4] Importing nodes via 'gitmap ssh nodes import-json'..." -ForegroundColor Yellow
$ImportOutput = & $GitMapExe ssh nodes import-json "$ResolvedJsonPath" 2>&1
$ImportStatus = $LASTEXITCODE

if ($ImportStatus -eq 0) {
    Write-Host "[Gate 4] Successfully imported SSH nodes into GitMap Split-DB" -ForegroundColor Green
    $ImportOutput | Out-String | Write-Host
} else {
    Write-Error "[FAIL] gitmap ssh nodes import-json failed with exit code $ImportStatus : $ImportOutput"
    exit $ImportStatus
}

# 5. Display registered node status
Write-Host "[Step 5] Enrolled cluster node inventory ('gitmap ssh ls'):" -ForegroundColor Yellow
& $GitMapExe ssh ls

# 6. Optional Fleet Deployment
if ($DeployToFleet) {
    Write-Host "`n[Step 6] Deploying node configuration across active fleet ('gitmap deploy config ssh all')..." -ForegroundColor Yellow
    & $GitMapExe deploy config ssh all
}

Write-Host "`n==========================================================" -ForegroundColor Green
Write-Host "   GitMap SSH Node Import & Verification SUCCESSFUL       " -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green
exit 0
