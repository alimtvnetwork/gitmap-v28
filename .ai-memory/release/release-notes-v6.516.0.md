# GitMap v6.516.0 Release Notes

## Overview
GitMap `v6.516.0` delivers universal forward-slash path hygiene across all remediation outputs, enforces upstream pull-before-changes invariants across remediation workflows, and introduces an interactive fix menu (`[a] Fix all`, `[s] Fix single`, `[k] Skip`) when multi-repo pull operations (`gitmap pa` / `gitmap pull-all`) encounter dirty or failed repositories.

## Key Enhancements & Fixes
- **Universal Forward-Slash (`/`) Path Hygiene:**
  - Eliminated escaped Windows double-backslashes (`D:\\work\\gitmap`) in all remediation hints and command generators (`formatRepoGitCmd`), formatting clean paths: `git -C "D:/work/gitmap" stash`.
  - Normalized itemized modified and untracked file paths in terminal output using `filepath.ToSlash(f)`.
  - Stored forward-slash normalized paths in `RemediationItem.RepoPath` and prompt status cards.
- **Pull Before Changes Invariant:**
  - Ensured that `commitCmd` (`add -A && commit && pull --rebase`) and `stashCmd` (`stash -u && pull && stash pop`) pull remote updates before finalizing local modifications.
- **Interactive Multi-Repo Pull Remediation Prompt:**
  - Lifted concise output suppression in `handlePullRemediation` so that `gitmap pa` / `gitmap pull-all` interactively prompts the user when dirty repositories exist.
  - Interactive options:
    - `[a/1]` Fix all: Automatically remediates all dirty repositories with stash/re-apply and pre-pull synchronization.
    - `[s/2]` Fix single: Steps interactively through dirty repositories.
    - `[k/q]` Skip / Exit: Cleanly exits and displays CLI remediation guidance.
- **Quality Gates & Linters:**
  - 100% compliant with AST linters (0 nested ifs across 4,432 files, affirmative booleans, gofmt clean, 0 golangci-lint errors).
  - FastGate pre-commit checks pass in <0.15s.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.516.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.516.0/install.sh | sh
```
