# Subtask: Swallowed Errors (PowerShell)
Status: 🔄 In Progress

## Steps
1. Edit fix-repo.ps1: line 139, 260, add explicit Write-Error log to empty catch block
2. Edit gitmap/scripts/Get-LastRelease.ps1: lines 47, 69, 93, 112, add log to empty catch blocks
3. Audit and fix all catch {} in PowerShell scripts.
