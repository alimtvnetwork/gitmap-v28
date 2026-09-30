<#
.SYNOPSIS
    Standalone GitMap Repository Cloner (v2.0 Manifest Compatible).
.DESCRIPTION
    Clones all missing repositories defined in a GitMap JSON manifest (e.g. gitmap.json)
    without requiring the gitmap binary to be installed. Expands ${workDir} dynamically,
    skips repositories already present on disk, and supports branch/transport selection.
.PARAMETER Manifest
    Path to the gitmap.json manifest file. Defaults to ./gitmap.json.
.PARAMETER WorkDir
    Target directory for cloning. Defaults to the workDir defined in the manifest, or current directory.
.PARAMETER UseSSH
    Force cloning over SSH (git@github.com:...).
.PARAMETER UseHTTPS
    Force cloning over HTTPS (https://github.com/...).
.PARAMETER DryRun
    Preview repositories that would be cloned without executing git clone.
.EXAMPLE
    pwsh ./clone-gitmap.ps1
    pwsh ./clone-gitmap.ps1 -Manifest ./gitmap.json -UseSSH
    pwsh ./clone-gitmap.ps1 -DryRun
#>
[CmdletBinding()]
param(
    [string]$Manifest = "",
    [string]$WorkDir = "",
    [switch]$UseSSH,
    [switch]$UseHTTPS,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

if (-not $Manifest) {
    if (Test-Path "$PSScriptRoot/gitmap.json") {
        $Manifest = "$PSScriptRoot/gitmap.json"
    } elseif (Test-Path "./gitmap.json") {
        $Manifest = "./gitmap.json"
    } else {
        Write-Error "No gitmap.json manifest found in current or script directory. Specify -Manifest <path>."
        exit 1
    }
}

if (-not (Test-Path $Manifest)) {
    Write-Error "Manifest file not found: $Manifest"
    exit 1
}

Write-Host "Reading manifest: $Manifest" -ForegroundColor Cyan
$rawJson = Get-Content -Path $Manifest -Raw -Encoding UTF8
$doc = $rawJson | ConvertFrom-Json

$resolvedWorkDir = $WorkDir
if (-not $resolvedWorkDir) {
    if ($doc.variables -and $doc.variables.workDir) {
        $resolvedWorkDir = $doc.variables.workDir
    } elseif ($doc.attributes -and $doc.attributes.workDirectory -and $doc.attributes.workDirectory.variables -and $doc.attributes.workDirectory.variables.workDir) {
        $resolvedWorkDir = $doc.attributes.workDirectory.variables.workDir
    } else {
        $resolvedWorkDir = (Get-Location).Path
    }
}

Write-Host "Target work directory: $resolvedWorkDir" -ForegroundColor Cyan
if (-not (Test-Path $resolvedWorkDir) -and -not $DryRun) {
    New-Item -ItemType Directory -Path $resolvedWorkDir -Force | Out-Null
}

$repos = @()
if ($doc.data) {
    $repos = $doc.data
} elseif ($doc -is [array]) {
    $repos = $doc
}

if ($repos.Count -eq 0) {
    Write-Warning "No repository records found in $Manifest."
    exit 0
}

Write-Host "Discovered $($repos.Count) repository record(s) in manifest." -ForegroundColor Green

$clonedCount = 0
$skippedCount = 0
$failedCount = 0

foreach ($repo in $repos) {
    $repoName = if ($repo.repoName) { $repo.repoName } else { $repo.slug }
    $relPath = if ($repo.relativePath) { $repo.relativePath } else { $repoName }
    
    # Expand ${workDir}
    $targetPath = Join-Path $resolvedWorkDir $relPath
    $targetPath = [System.IO.Path]::GetFullPath($targetPath)

    $gitDir = Join-Path $targetPath ".git"
    if (Test-Path $gitDir) {
        Write-Host "  [skip] $repoName already exists: $targetPath" -ForegroundColor DarkGray
        $skippedCount++
        continue
    }

    # Resolve URL based on switches and manifest
    $cloneUrl = ""
    if ($UseSSH -and $repo.sshUrl) {
        $cloneUrl = $repo.sshUrl
    } elseif ($UseHTTPS -and $repo.httpsUrl) {
        $cloneUrl = $repo.httpsUrl
    } elseif ($repo.transport -eq "ssh" -and $repo.sshUrl) {
        $cloneUrl = $repo.sshUrl
    } elseif ($repo.httpsUrl) {
        $cloneUrl = $repo.httpsUrl
    } elseif ($repo.sshUrl) {
        $cloneUrl = $repo.sshUrl
    } else {
        $cloneUrl = $repo.discoveredUrl
    }

    $branch = if ($repo.branch) { $repo.branch } else { "main" }

    if ($DryRun) {
        Write-Host "  [dry-run] git clone -b $branch $cloneUrl `"$targetPath`"" -ForegroundColor Yellow
        $clonedCount++
        continue
    }

    $parentDir = Split-Path -Parent $targetPath
    if (-not (Test-Path $parentDir)) {
        New-Item -ItemType Directory -Path $parentDir -Force | Out-Null
    }

    Write-Host "  [clone] $repoName (branch: $branch) -> $targetPath" -ForegroundColor White
    $cloneArgs = @("clone", "-b", $branch, $cloneUrl, $targetPath)
    
    & git $cloneArgs
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  [done]  $repoName cloned successfully." -ForegroundColor Green
        $clonedCount++
    } else {
        Write-Host "  [error] Failed to clone $repoName (Exit code: $LASTEXITCODE)" -ForegroundColor Red
        $failedCount++
    }
}

Write-Host "`nSummary:" -ForegroundColor Cyan
Write-Host "  Total:   $($repos.Count)" -ForegroundColor White
Write-Host "  Cloned:  $clonedCount" -ForegroundColor Green
Write-Host "  Skipped: $skippedCount" -ForegroundColor DarkGray
Write-Host "  Failed:  $failedCount" -ForegroundColor $(if ($failedCount -gt 0) { "Red" } else { "Green" })

if ($failedCount -gt 0) {
    exit 1
}
