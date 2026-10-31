# Subtask 04: Local Installer Parity — Integrate Suggestions into local-install.ps1

- **Task Code:** `TASK-04`
- **Slug:** `gitmap-installation-suggestion-issue - Task 04`
- **Parent Plan:** [365-gitmap-installation-suggestion-issue.md](../pending/365-gitmap-installation-suggestion-issue.md)
- **Status:** `pending`

---

## Scope & Deliverables
1. Port `Configure-PowerShellProfileSuggestions`, `Resolve-PowerShellProfileTargets`, and `Write-ProfileSubtree` into `local-install.ps1`.
2. Hook `Configure-PowerShellProfileSuggestions $targetBin` into Stage 5 (Shell Environment Activation).
3. Ensure offline/local installations report active predictive suggestions identically to remote installer.
