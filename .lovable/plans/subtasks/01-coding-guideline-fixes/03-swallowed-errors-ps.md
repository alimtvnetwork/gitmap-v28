# Subtask: Swallowed Errors in PowerShell
Status: pending

## Steps
1. `fix-repo.ps1:139`: Replace silent catch with `catch { Write-Warning "fix-repo: failed to process branch: $_" }`
2. `fix-repo.ps1:261`: Replace silent catch with `catch { Write-Warning "fix-repo: operation failed: $_" }`
3. `gitmap/scripts/Get-LastRelease.ps1:47`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: API error: $_" }`
4. `gitmap/scripts/Get-LastRelease.ps1:69`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: JSON parse error: $_" }`
5. `gitmap/scripts/Get-LastRelease.ps1:93`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: fallback query error: $_" }`
6. `gitmap/scripts/Get-LastRelease.ps1:112`: Replace silent catch with `catch { Write-Warning "Get-LastRelease: unexpected error: $_" }`
7. `gitmap/scripts/install.ps1:130`: Replace silent catch with `catch { Write-Warning "install: elevate failed: $_" }`
8. `gitmap/scripts/install.ps1:175`: Replace silent catch with `catch { Write-Warning "install: temp cleanup failed: $_" }`
9. `gitmap/scripts/install.ps1:196`: Replace silent catch with `catch { Write-Warning "install: asset download failed: $_" }`
10. `gitmap/scripts/install.ps1:359`: Replace silent catch with `catch { Write-Warning "install: checksum failed: $_" }`
11. `gitmap/scripts/install.ps1:371`: Replace silent catch with `catch { Write-Warning "install: signature check failed: $_" }`
12. `gitmap/scripts/install.ps1:419`: Replace silent catch with `catch { Write-Warning "install: unpack failed: $_" }`
13. `gitmap/scripts/install.ps1:532`: Replace silent catch with `catch { Write-Warning "install: link creation failed: $_" }`
14. `gitmap/scripts/install.ps1:654`: Replace silent catch with `catch { Write-Warning "install: profile registration failed: $_" }`
15. `gitmap/scripts/install.ps1:687`: Replace silent catch with `catch { Write-Warning "install: PATH update failed: $_" }`
16. `gitmap/scripts/release-version.ps1:203`: Replace silent catch with `catch { Write-Warning "release-version: tag query failed: $_" }`
17. `gitmap/scripts/release-version.ps1:300`: Replace silent catch with `catch { Write-Warning "release-version: git status failed: $_" }`
18. `gitmap/scripts/release-version.ps1:314`: Replace silent catch with `catch { Write-Warning "release-version: branch validation failed: $_" }`
19. `gitmap/scripts/release-version.ps1:355`: Replace silent catch with `catch { Write-Warning "release-version: changelog sync failed: $_" }`
20. `gitmap/scripts/release-version.ps1:389`: Replace silent catch with `catch { Write-Warning "release-version: tag creation failed: $_" }`
