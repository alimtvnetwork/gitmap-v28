# Subtask 03: Repo-Secrets Gitmap JSON Harmonization and Standalone PowerShell Script

> **Parent Plan:** `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md`
> **Status:** `PENDING`
> **Target Files:**
> - `D:\work\repo-secrets\07-final-network-machine\gitmap.json`
> - `D:\work\repo-secrets\07-final-network-machine\clone-gitmap.ps1`
> - `scripts/clone-gitmap.ps1`

---

## Technical Specification

1. **`D:\work\repo-secrets\07-final-network-machine\gitmap.json`**:
   - Set `"version": "2.0"` in `attributes`.
   - Update `workDirectory`:
     ```json
     "workDirectory": {
       "path": "${workDir}",
       "defaultPath": "${workDir}",
       "isApplied": true,
       "isEnforced": false,
       "variables": {
         "workDir": "D:\\work"
       }
     }
     ```
   - Remove redundant top-level `defaultWorkDirectory`, `isWorkDirectoryApplied`, `isWorkDirectoryEnforced`.
   - Set `help` object:
     ```json
     "help": {
       "description": "Master GitMap repositories manifest with cross-platform clone instructions and variable resolution.",
       "importCommand": "gitmap import-all-json gitmap.json",
       "cloneMissingCommand": "gitmap clone-only-missing gitmap.json",
       "cloneOnlyMissingCommand": "gitmap cfr gitmap.json",
       "mergeCommand": "gitmap merge-json gitmap.json",
       "whichFormatCommand": "gitmap which-format gitmap.json",
       "examples": [
         "gitmap cfr gitmap.json                      # Clone only missing repositories with automated fix-repo",
         "gitmap clone gitmap.json                    # Clone repositories declared in manifest",
         "gitmap clone-only-missing gitmap.json       # Clone only repositories not yet present on disk",
         "gitmap merge-json target.json source.json   # Merge and synchronize repository manifests",
         "pwsh ./clone-gitmap.ps1                     # Standalone PowerShell clone without gitmap CLI"
       ]
     }
     ```
   - Top-level variables:
     `"workDir": "D:\\work"`, `"repoDir": "${workDir}\\gitmap"`, `"secretsDir": "."`.
   - `data`:
     - 1-based indexing: `id` starts at 1 (1 to 66).
     - Replace hardcoded `D:\work\` with `${workDir}\` in `absolutePath`.

2. **Standalone `clone-gitmap.ps1`**:
   - Works with PowerShell 5.1 and Core (pwsh 7+).
   - Reads `gitmap.json` (defaults to same folder if `-Manifest` not passed).
   - Expands `${workDir}` (defaults to parent folder or current folder).
   - Iterates through `data` array:
     - If destination folder exists (or has `.git`): skips with `[skip] ... already exists`.
     - If missing: executes `git clone -b $branch $url $destFolder`.
   - Summarizes total cloned, skipped, and failed.
