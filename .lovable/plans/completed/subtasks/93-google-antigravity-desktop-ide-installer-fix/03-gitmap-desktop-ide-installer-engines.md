# Subtask 93.03: Gitmap Desktop IDE Installer & Platform Engines

## Goal
Implement cross-platform Google Antigravity Desktop IDE installation in `gitmap/cmd/`.

## Files Impacted
- `gitmap/cmd/installantigravity.go`
- `gitmap/cmd/installantigravity_fetch.go`
- `gitmap/cmd/installantigravity_linux.go` (NEW)
- `gitmap/cmd/installantigravity_windows.go` (NEW)

## Acceptance Criteria
1. `runInstallAntigravityWithOpts` detects existing desktop installation before downloading.
2. Linux: downloads `Antigravity.tar.gz`, unpacks to user/system directory, symlinks binary `antigravity`, sets `0755` permissions, and verifies.
3. Windows: downloads `Antigravity-x64.exe`, runs `/S` silent installer, checks `%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe`.
4. Saves to `InstallationSplitDB` as tool `"antigravity"`.
5. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
