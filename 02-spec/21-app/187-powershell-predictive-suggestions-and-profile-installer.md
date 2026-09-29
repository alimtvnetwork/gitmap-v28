# Specification 187: PowerShell Predictive Suggestions & Profile Auto-Configuration

## 1. Executive Summary & Problem Root Cause

### 1.1 Problem Statement
Users frequently observed that interactive command suggestions (the PSReadLine dropdown `ListView` menu showing `<History(10)>` predictions such as `gitmap ssh join ...`, `gitmap update`, etc.) were missing after running GitMap updates or installing GitMap on fresh machines and worker nodes.

### 1.2 Root Cause Analysis (4-Part RCA)

1. **PSReadLine Version Disparity on Windows PowerShell 5.1**:
   - Out of the box, Windows 10/11 and Windows Server ship Windows PowerShell 5.1 with built-in **PSReadLine 2.0.0**.
   - The parameters `-PredictionSource` and `-PredictionViewStyle ListView` were only introduced in **PSReadLine 2.2.0**.
   - Invoking `Set-PSReadLineOption -PredictionSource History` or `Set-PSReadLineOption -PredictionViewStyle ListView` in Windows PowerShell 5.1 throws `NamedParameterNotFound`.
   - Previous scripts silently caught and swallowed this error without notifying the user or offering the required upgrade command (`Install-Module PSReadLine -Scope CurrentUser -Force -SkipPublisherCheck`).

2. **Profile Multiplicity & Absence by Default**:
   - Windows does **not** create `$PROFILE` files by default for new user accounts or fresh systems.
   - PowerShell maintains separate profile paths for:
     - PowerShell 7 Core (`pwsh`): `$HOME\Documents\PowerShell\Microsoft.PowerShell_profile.ps1` and `profile.ps1`
     - Windows PowerShell 5.1 (`powershell.exe`): `$HOME\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1` and `profile.ps1`
   - Neither `install.ps1` nor `cli/scripts/install.ps1` contained logic to discover, initialize, or populate these profile directories.

3. **Output Redirection Guard & Console Handle Checks**:
   - In earlier versions, PSReadLine configuration was guarded with `-not [Console]::IsOutputRedirected`. In many terminal hosts (VS Code integrated terminals, redirected install runners, SSH sessions, non-standard console handles), this test evaluated to false, skipping configuration.
   - Non-interactive calls to `Set-PSReadLineOption -PredictionViewStyle ListView` throw "The handle is invalid" because screen buffer APIs require interactive console handles. Profiling code must guard this dynamically so shell startup never crashes.

4. **Line Ending Disparity & Non-Idempotent Profile Pollution**:
   - In `cli/completion/install.go`, profile strings were compared using Go's `\n` line endings against Windows files containing `\r\n`.
   - As a result, `strings.Contains` always returned false, causing repeated calls to `gitmap setup` or `gitmap completion install` to append duplicate legacy blocks (observed as 7+ redundant blocks).
   - Old single-line completion lines (`. '...completions.ps1'`) were never reconciled or upgraded in-place.

5. **Update Disconnect**:
   - Post-update lifecycle routines (`cli/cmdupdate/updateremoteinstall.go`) replaced binary executables but never refreshed user profiles or verified that PSReadLine predictive suggestions were configured across all shells.

---

## 2. Architectural Solution & Implementation

### 2.1 Standardized Marker Block Architecture
Both the Go CLI completion engine (`cli/completion/install.go`) and the PowerShell installer (`install.ps1` & `cli/scripts/install.ps1`) deploy an identical, idempotent marker block:

```powershell
# >>> gitmap shell completion & predictive suggestions >>>
if (-not $global:__gitmap_suggestions_configured) {
    $global:__gitmap_suggestions_configured = $true
    if (Get-Command Set-PSReadLineOption -ErrorAction SilentlyContinue) {
        try {
            Set-PSReadLineOption -PredictionSource HistoryAndPlugin -ErrorAction SilentlyContinue
        } catch {
            try {
                Set-PSReadLineOption -PredictionSource History -ErrorAction SilentlyContinue
            } catch {}
        }
        try {
            Set-PSReadLineOption -PredictionViewStyle ListView -ErrorAction SilentlyContinue
        } catch {}
    }
}
$compPath = "$env:APPDATA\gitmap\completions.ps1"
if (Test-Path -LiteralPath $compPath) {
    . $compPath
}
# <<< gitmap shell completion & predictive suggestions <<<
```

### 2.2 Global Execution Guard
The `$global:__gitmap_suggestions_configured` guard ensures that even when both `CurrentUserAllHosts` (`profile.ps1`) and `CurrentUserCurrentHost` (`Microsoft.PowerShell_profile.ps1`) are sourced in a single session, the PSReadLine options and prediction hooks are executed exactly once.

### 2.3 Resilient PSReadLine Fallback Chain
1. Probe command existence via `Get-Command Set-PSReadLineOption -ErrorAction SilentlyContinue`.
2. Attempt `Set-PSReadLineOption -PredictionSource HistoryAndPlugin`.
3. If `HistoryAndPlugin` is unavailable (PSReadLine 2.1), fall back to `Set-PSReadLineOption -PredictionSource History`.
4. If running on PSReadLine 2.0.0, the inner `catch` silently prevents terminal startup errors.
5. Attempt `Set-PSReadLineOption -PredictionViewStyle ListView`, safely caught if running in a headless or unsupported host.

### 2.4 Automatic Profile Reconciliation & Legacy Strip
When writing profiles:
1. Normalize line endings (`\r\n` vs `\n`).
2. Search for existing `# >>> gitmap shell completion & predictive suggestions >>>` blocks and replace in-place.
3. Automatically detect and strip legacy `# gitmap shell completion` blocks, removing stale paths and duplicate lines.
4. Preserve all non-GitMap profile contents (e.g., custom functions, environment variables, PATH modifications).
5. Only write to disk when contents have actually changed, emitting `MsgCompAlreadyDone` when clean.

### 2.5 Multi-Target Profile Provisioning
Both `install.ps1` and `completion.Install` resolve and provision all user profile targets:
- `Documents\PowerShell\Microsoft.PowerShell_profile.ps1` (pwsh CurrentHost)
- `Documents\PowerShell\profile.ps1` (pwsh AllHosts)
- `Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1` (WinPS 5.1 CurrentHost)
- `Documents\WindowsPowerShell\profile.ps1` (WinPS 5.1 AllHosts)
- `.config\powershell\` (cross-platform pwsh)

### 2.6 Installer Guidance for Windows PowerShell 5.1
When `install.ps1` runs on Windows, it inspects `Get-Module -ListAvailable PSReadLine`. If the highest version is `< 2.2.0`, it prints an actionable recommendation:
```
  ◦ Note: Windows PowerShell built-in PSReadLine is 2.0.0 (requires >= 2.2.0 for predictive dropdown).
    To enable ListView predictive IntelliSense in Windows PowerShell 5.1, run:
      Install-Module PSReadLine -Scope CurrentUser -Force -SkipPublisherCheck
```

### 2.7 Post-Update Lifecycle Hook
In `cli/cmdupdate/updateremoteinstall.go`, `finishRemoteUpdate` invokes `ensurePostUpdateCompletions()`, which calls `completion.Install(completion.DetectShell())`. This ensures that every `gitmap update` execution refreshes completion scripts and keeps profile predictive suggestions active.

---

## 3. Verification & Compliance
- **Unit Testing**: `cli/completion/install_test.go` verifies:
  - Idempotent double-installation (`TestAddSourceLinePowerShellIdempotent`).
  - Automatic removal of legacy single-line and multi-line blocks (`TestAddSourceLinePowerShellStripsLegacyBlocks`).
  - Preservation of user-defined functions and variables.
  - Windows CRLF and Unix LF line ending stability.
- **Interactive Verification**:
  - `gitmap completion install powershell` idempotently configured all 4 profile files on Windows.
  - Verified toggle key `F2` between `ListView` (popup menu) and `InlineView` (ghost text).
