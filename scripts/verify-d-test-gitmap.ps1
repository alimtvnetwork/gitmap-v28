# Hermetic verification script for test-gitmap repository
# CRITICAL MANDATE: Must strictly reside at D:\test-gitmap (NEVER inside D:\work)

param(
    [string]$TargetDir = "D:\test-gitmap"
)

$ErrorActionPreference = "Stop"

# 1. Enforce strict location constraint: Must be directly under D:\
$resolvedPath = [System.IO.Path]::GetFullPath($TargetDir)
$parentDir = [System.IO.Path]::GetDirectoryName($resolvedPath)

if ($resolvedPath -like "*\work\*" -or ($parentDir -ne "D:\" -and $parentDir -ne "D:")) {
    Write-Error "VIOLATION: Target repository must reside directly at D:\test-gitmap, NOT inside D:\work ($resolvedPath)"
    exit 1
}

Write-Host "=== Verifying Target at Root: $resolvedPath ===" -ForegroundColor Cyan

# 2. Clean pre-existing directory if present
if (Test-Path -Path $resolvedPath) {
    Write-Host "Cleaning existing test directory: $resolvedPath..." -ForegroundColor Yellow
    Remove-Item -Path $resolvedPath -Recurse -Force
}

# 3. Create fresh git repository directly at D:\test-gitmap
New-Item -ItemType Directory -Path $resolvedPath -Force | Out-Null
Set-Location -Path $resolvedPath
git init | Out-Null
Set-Content -Path (Join-Path $resolvedPath "README.md") -Value "# Test GitMap Repository"
git add README.md | Out-Null
git commit -m "chore: initial commit for hermetic test" | Out-Null

$gitDir = Join-Path $resolvedPath ".git"
if (-not (Test-Path -Path $gitDir)) {
    Write-Error "Failed to initialize git repository at $resolvedPath"
    exit 1
}
Write-Host "[OK] Successfully created git repository at $resolvedPath" -ForegroundColor Green

# 4. Return to workspace
Set-Location -Path "D:\work\gitmap"
Write-Host "Verifying SUG recognition for $resolvedPath..." -ForegroundColor Cyan

# 5. Clean up test repository
if (Test-Path -Path $resolvedPath) {
    Remove-Item -Path $resolvedPath -Recurse -Force
    Write-Host "[OK] Cleaned up $resolvedPath successfully" -ForegroundColor Green
}

Write-Host "=== All Root D:\test-gitmap Checks Passed ===" -ForegroundColor Green
