# fix-all-permissions.ps1: Elevate and grant full control to current user on D:\wp-work\riseup-asia
param(
    [string]$TargetDir = "D:\wp-work\riseup-asia"
)

$currentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent().Name
$shortUser = $env:USERNAME

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " Repairing Git & File Permissions on: $TargetDir" -ForegroundColor Cyan
Write-Host " Target User: $currentUser ($shortUser)" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Take ownership
Write-Host "[1/3] Taking ownership of files and folders..." -ForegroundColor Yellow
takeown /f $TargetDir /r /d y

# 2. Grant Full Control with inheritance
Write-Host "[2/3] Granting Full Control ACLs to $shortUser..." -ForegroundColor Yellow
icacls $TargetDir /grant "${shortUser}:(OI)(CI)F" /t /c /q

# 3. Strip read-only attributes
Write-Host "[3/3] Stripping read-only attributes..." -ForegroundColor Yellow
attrib -r "$TargetDir\*" /s /d

Write-Host "==========================================================" -ForegroundColor Green
Write-Host " Permissions successfully repaired for $shortUser!" -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green
