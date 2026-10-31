# Spec 365: GitMap Installation Suggestion Issue — Root Cause Analysis & Architecture

- **Task Code:** `365-gitmap-installation-suggestion-issue`
- **Slug:** `gitmap-installation-suggestion-issue`
- **Status:** `in-progress`
- **Date:** `31-Oct-26`

---

## 1. Executive Summary

During GitMap installation and command execution, terminal suggestions (including PSReadLine predictive IntelliSense, shell completions, and post-installation CLI hint readouts) intermittently appeared or disappeared completely across test sessions. This specification details the exhaustive Root Cause Analysis (RCA) covering all 6 identified failure modes and defines the architectural model guaranteeing persistent, deterministic suggestion delivery.

---

## 2. Four-Part Grounded Root Cause Analysis (RCA)

### 2.1 Failure Mode 1: Profile Readout Suppression on Re-runs and Updates
- **Symptom:** In `install.ps1`, running the installer a second time or upgrading an existing GitMap installation causes the message `Configured predictive suggestions:` to vanish completely from terminal output.
- **Root Cause:** In function `Configure-PowerShellProfileSuggestions`, the target profiles were only added to `$configuredProfiles` if `$newContent -ne $content`. On any existing installation where the block was already present, `$configuredProfiles.Count` evaluated to `0`, causing the summary call `Write-ProfileSubtree "Configured predictive suggestions:"` to be completely bypassed.
- **Remediation:** Track `$activeProfiles` alongside `$configuredProfiles`. If no files required modification, emit `Write-ProfileSubtree "Active predictive suggestions (already configured):"` so the user and AI agents receive continuous visual verification of the suggestion state.

### 2.2 Failure Mode 2: Update Mode Skip Flag
- **Symptom:** Running `gitmap update` or `install.ps1` with `$env:GITMAP_UPDATING = "1"` skips suggestion verification entirely.
- **Root Cause:** Line 2187 guarded execution with `if (-not $Quiet -and -not $JSON -and $env:GITMAP_UPDATING -ne "1")`. This suppressed `Configure-PowerShellProfileSuggestions` during self-updates, preventing repairs to broken or newly added completion hooks.
- **Remediation:** Remove the `$env:GITMAP_UPDATING -ne "1"` guard so suggestion configuration and verification execute uniformly across both initial installs and updates.

### 2.3 Failure Mode 3: Local Installer Parity Gap (`local-install.ps1`)
- **Symptom:** Running `local-install.ps1` installs the binary and updates PATH, but predictive suggestions and completions are never configured in PowerShell profiles.
- **Root Cause:** `local-install.ps1` completely omitted `Configure-PowerShellProfileSuggestions` and profile discovery helpers in Stage 5 (Shell Environment Activation).
- **Remediation:** Integrate `Configure-PowerShellProfileSuggestions`, `Resolve-PowerShellProfileTargets`, and `Write-ProfileSubtree` into `local-install.ps1`, ensuring 100% parity with remote install workflows.

### 2.4 Failure Mode 4: Cobra Completion Alias Omission (`gm`)
- **Symptom:** Typing `gitmap <tab>` displayed suggestions, but typing `gm <tab>` (the primary 2-letter alias) showed nothing.
- **Root Cause:** In `cli/cmd/root_cobra_completion.go`, Cobra's PowerShell completion generator registered `Register-ArgumentCompleter -CommandName 'gitmap'`, but omitted `gm`.
- **Remediation:** Append `Register-ArgumentCompleter -CommandName 'gm' -ScriptBlock ${__gitmapCompleterBlock}` in `GenerateCobraCompletionScript` and the generated completion script.

### 2.5 Failure Mode 5: Non-Versioned Repo Remote URL Suppression
- **Symptom:** During release operations or clone-level checks, `printInstallHint` produced no output in repositories cloned without a numeric version suffix (e.g. `alimtvnetwork/gitmap`).
- **Root Cause:** In `cli/release/releaseinstallhint.go`, `isVersionedGitmapRepo` required `strings.HasPrefix(repoName, "gitmap-v")` and numeric suffix verification. Unversioned canonical repos evaluated to `false`.
- **Remediation:** Update `isVersionedGitmapRepo` to explicitly accept `repoName == "gitmap"` in addition to `gitmap-v<N>`.

### 2.6 Failure Mode 6: Silent Exit on Already-Installed `gitmap`
- **Symptom:** Running `gitmap install gitmap` in the terminal when GitMap is installed outputs `✓ gitmap is already installed` and immediately terminates without suggestions.
- **Root Cause:** In `cli/cmdinstall/install_prompt.go`, `alreadyInstalled` lacked a suggestions hook for `gitmap`.
- **Remediation:** Add `printGitmapInstalledSuggestions()` to display actionable next-step commands (`completion install`, `status`, `update`, `deploy-bin all`).

---

## 3. Architecture & Data Flow

```
[Installer Execution: install.ps1 / local-install.ps1]
             │
             ▼
   Resolve Profile Targets
  (CurrentUser + WindowsPowerShell + PowerShell 7)
             │
             ▼
   Evaluate Block State
             │
    ┌────────┴────────┐
    ▼                 ▼
[New / Outdated]  [In Sync]
    │                 │
    ▼                 ▼
Write Block       Mark Active
$configured++     $active++
    │                 │
    └────────┬────────┘
             │
             ▼
   Always-Report Engine:
   - Configured -> Write-ProfileSubtree "Configured predictive suggestions:"
   - Verified   -> Write-ProfileSubtree "Active predictive suggestions (already configured):"
             │
             ▼
   Cobra Engine (gm + gitmap completers active)
```
