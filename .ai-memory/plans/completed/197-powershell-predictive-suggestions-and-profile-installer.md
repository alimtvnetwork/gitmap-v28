# Plan 197: PowerShell Predictive Suggestions & Profile Auto-Configuration

## 1. Intent & Context
- The user requested:
  `Why this suggestions are not availle everywhere after updates, how to get those and why don't you configure this in the install script, pleasee`
- Investigated missing PSReadLine predictive IntelliSense (the dropdown `<History(10)>` ListView menu) across machines and sessions post-update.
- Conducted 4-part root cause analysis (RCA):
  1. Windows PowerShell 5.1 ships with PSReadLine 2.0.0 by default, which lacks `-PredictionSource` and `-PredictionViewStyle ListView` (introduced in PSReadLine 2.2.0).
  2. Windows does not create `$PROFILE` files by default; neither `install.ps1` nor `cli/scripts/install.ps1` configured PSReadLine or initialized profiles.
  3. `-not [Console]::IsOutputRedirected` caused completion setup to be skipped in redirected, nested, or SSH terminal hosts.
  4. CRLF (`\r\n`) vs LF (`\n`) comparison mismatch in Go's `addSourceLine` caused `strings.Contains` to always fail, resulting in 7+ duplicate completion blocks in user profiles.
  5. `gitmap update` did not refresh shell profile completions after replacing the binary.

## 2. Tasks Executed
1. **PowerShell Installers (`install.ps1` & `cli/scripts/install.ps1`)**:
   - Implemented `Configure-PowerShellProfileSuggestions` targeting all 4 PowerShell profile files (`PowerShell` & `WindowsPowerShell`, `Microsoft.PowerShell_profile.ps1` & `profile.ps1`).
   - Added automatic legacy block stripping and replacement with a clean, idempotent marker block (`# >>> gitmap shell completion & predictive suggestions >>>`).
   - Added friendly upgrade guidance for Windows PowerShell 5.1 when built-in PSReadLine is `< 2.2.0`.
   - Wired `Configure-PowerShellProfileSuggestions` immediately after setup in `install.ps1` and synchronized `cli/scripts/install.ps1`.
2. **Go Completion Engine (`cli/completion/install.go`)**:
   - Implemented idempotent profile reconciliation with native EOL preservation (`\r\n` on Windows, `\n` on Unix).
   - Added `stripLegacyPowerShellBlocks` to clean up duplicate legacy blocks without affecting user custom functions/variables.
   - Added global execution guard `$global:__gitmap_suggestions_configured`.
   - Added resilient PSReadLine fallback chain (`HistoryAndPlugin` -> `History`).
3. **Cobra Completion Template (`cli/cmd/root_cobra_completion.go`)**:
   - Updated PowerShell completion script generation to include the resilient PSReadLine prediction block.
4. **Post-Update Hook (`cli/cmdupdate/updateremoteinstall.go`)**:
   - Added `ensurePostUpdateCompletions()` inside `finishRemoteUpdate` so every `gitmap update` automatically updates user profiles and ensures predictive IntelliSense is active.
5. **Testing & Verification**:
   - Added unit tests in `cli/completion/install_test.go`:
     - `TestAddSourceLinePowerShellIdempotent` (verifies exactly 1 marker block after multiple installs).
     - `TestAddSourceLinePowerShellStripsLegacyBlocks` (verifies legacy block cleanup while preserving custom user functions).
   - Rebuilt binary and tested `gitmap completion install powershell`. Verified clean reconciliation and removal of 7 duplicate blocks in active profiles.
   - Tested idempotence: repeated runs report `Shell completion already configured for powershell`.
6. **Documentation**:
   - Created Specification `02-spec/21-app/187-powershell-predictive-suggestions-and-profile-installer.md`.
   - Updated `02-spec/21-app/readme.md`.
