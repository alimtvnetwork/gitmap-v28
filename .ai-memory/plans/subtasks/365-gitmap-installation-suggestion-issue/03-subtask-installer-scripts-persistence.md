# Subtask 03: Installer Scripts Fix — Profile Suggestion Persistence & Always-Report Logic

- **Task Code:** `TASK-03`
- **Slug:** `gitmap-installation-suggestion-issue - Task 03`
- **Parent Plan:** [365-gitmap-installation-suggestion-issue.md](../pending/365-gitmap-installation-suggestion-issue.md)
- **Status:** `pending`

---

## Scope & Deliverables
1. Update `install.ps1`, `cli/scripts/install.ps1`, and `cli/fspath/install.ps1`.
2. Track `$activeProfiles` alongside `$configuredProfiles` in `Configure-PowerShellProfileSuggestions`.
3. Emit `Write-ProfileSubtree "Active predictive suggestions (already configured):"` if files are unchanged.
4. Remove the `-and $env:GITMAP_UPDATING -ne "1"` guard so suggestion checks run on updates.
