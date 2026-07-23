# Post-Mortem: Installer Crashes — Progress Bar & Binary Detection

## Issues Fixed (v2.55.0)

### 1. PowerShell Progress Bar Crash

**Symptom**: Running `irm ... | iex` froze or crashed the terminal with cyan progress bar overlay blocks.

**Root Cause**: `Invoke-WebRequest` renders a progress bar by default. When piped through `iex`, the progress rendering conflicts with the terminal session and can hang or crash PowerShell.

**Fix**: Added `$ProgressPreference = "SilentlyContinue"` at the top of `install.ps1` to suppress progress bars entirely.

```powershell
# ✅ Required at script top for irm | iex compatibility
$ProgressPreference = "SilentlyContinue"
```

### 2. Versioned Binary Name Not Detected

**Symptom**: Installer reported "Installed archive did not contain gitmap.exe" and listed the actual file (e.g., `gitmap-v28.54.6-windows-amd64.exe`).

**Root Cause**: The candidate name list only checked for `gitmap.exe` and `gitmap-windows-amd64.exe`, but the CI release pipeline names binaries with the version embedded (e.g., `gitmap-v28.54.6-windows-amd64.exe`).

**Fix**: Added regex pattern matching for versioned filenames:

```powershell
# ✅ Match versioned binary names from CI
$_.Name -match "^gitmap-v[\d.]+-windows-(amd64|arm64)\.exe$"
```

### 3. Missing Error Recovery

**Symptom**: Any failure during installation caused `$ErrorActionPreference = "Stop"` to terminate the script with an unhandled exception, crashing the terminal.

**Fix**: Wrapped the entire `Main` function body in `try/catch` with a friendly error message and manual download link fallback.

## Prevention Rules

1. **Always set `$ProgressPreference = "SilentlyContinue"`** in any PowerShell script intended for `irm | iex` usage.
2. **Binary detection must use flexible pattern matching** — never hardcode exact filenames when CI may embed versions, architectures, or build numbers.
3. **All installer scripts must have top-level error handling** — `$ErrorActionPreference = "Stop"` requires a corresponding `try/catch` to prevent terminal crashes.
4. **Test installers on clean machines** — `irm | iex` behaves differently from direct script execution (no `param()` binding, different scope rules).

## Related

- `gitmap-v28/scripts/install.ps1` — One-liner installer
- `spec/02-app-issues/13-release-pipeline-dist-directory.md` — CI pipeline issues
- CHANGELOG.md v2.55.0 — Installer Fix
