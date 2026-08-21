# Subtask: Swallowed Errors (PowerShell)
Status: ✅ Done (fixed empty catch blocks across 8 ps1 files by adding Write-Warning "$_" or Write-Error "$_")

## Steps
1. Edit fix-repo.ps1: line 139, 260, add explicit Write-Error log to empty catch block
2. Edit gitmap/scripts/Get-LastRelease.ps1: lines 47, 69, 93, 112, add log to empty catch blocks
3. Audit and fix all catch {} in PowerShell scripts.
