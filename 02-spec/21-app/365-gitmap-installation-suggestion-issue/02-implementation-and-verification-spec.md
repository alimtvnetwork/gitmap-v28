# Spec 365: GitMap Installation Suggestion Issue — Implementation & Verification Spec

- **Task Code:** `365-gitmap-installation-suggestion-issue`
- **Slug:** `gitmap-installation-suggestion-issue`
- **Status:** `in-progress`
- **Date:** `31-Oct-26`

---

## 1. Concrete File Modifications

### 1.1 `install.ps1`, `cli/scripts/install.ps1`, `cli/fspath/install.ps1`
1. **Always-Report Profile Suggestions:**
   - In `Configure-PowerShellProfileSuggestions`, initialize `$activeProfiles = @()`.
   - For every target profile containing the valid marker block, append to `$activeProfiles`.
   - If `$configuredProfiles.Count -gt 0`, display configured profiles subtree.
   - Else if `$activeProfiles.Count -gt 0`, display `Write-ProfileSubtree "Active predictive suggestions (already configured):" @($activeProfiles | Select-Object -Unique)`.
2. **Remove Update Mode Suppression:**
   - Change line 2187 from `if (-not $Quiet -and -not $JSON -and $env:GITMAP_UPDATING -ne "1")` to `if (-not $Quiet -and -not $JSON)`.
3. **PSReadLine 2.2.0+ Verification Pre-Flight:**
   - Ensure module discovery and copy happens prior to attempting `Set-PSReadLineOption`.

### 1.2 `local-install.ps1`
1. Add `Resolve-PowerShellProfileTargets`, `Write-ProfileSubtree`, and `Configure-PowerShellProfileSuggestions`.
2. In Stage 5 (Shell Environment Activation), call `Configure-PowerShellProfileSuggestions $targetBin`.

### 1.3 `cli/cmd/root_cobra_completion.go`
1. In `GenerateCobraCompletionScript` for PowerShell:
   - Append `Register-ArgumentCompleter -CommandName 'gm' -ScriptBlock ${__gitmapCompleterBlock}`.

### 1.4 `cli/release/releaseinstallhint.go`
1. In `isVersionedGitmapRepo`:
   - Add early return: if `repoName == "gitmap"` return `true`.

### 1.5 `cli/cmdinstall/install_prompt.go`
1. In `alreadyInstalled`:
   - If `installName == "gitmap"`, invoke `printGitmapInstalledSuggestions()`.
2. Implement `printGitmapInstalledSuggestions()` with Lipgloss cyan styling and actionable commands.

---

## 2. Verification Protocol

1. **Re-run Installer Test:** Run `install.ps1` twice consecutively; verify that `Active predictive suggestions (already configured):` prints on the second run.
2. **Local Installer Test:** Run `local-install.ps1 -DryRun` or inspect output; verify Stage 5 includes suggestions configuration.
3. **Cobra Alias Test:** Run `gitmap completion powershell` and verify output contains `Register-ArgumentCompleter -CommandName 'gm'`.
4. **Release Hint Test:** Verify `isVersionedGitmapRepo("gitmap")` returns `true`.
5. **GitMap Install Test:** Run `gitmap install gitmap` and verify the `💡 GitMap Suggestions:` box is rendered.
