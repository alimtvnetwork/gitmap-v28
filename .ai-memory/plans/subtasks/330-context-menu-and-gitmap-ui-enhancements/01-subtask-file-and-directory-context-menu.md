# Subtask 330.1 — File & Directory Context Menu Architecture

> **Parent Plan:** `.ai-memory/plans/completed/330-context-menu-and-gitmap-ui-enhancements.md`  
> **Target Packages:** `cli/cmdinstall/`  
> **Status:** Completed

---

## 1. Objectives
1. In `cli/cmdinstall/installctx.go`:
   - [x] Add `HKCU\Software\Classes\*\shell\gitmap` so GitMap appears on right-clicking ANY file.
   - [x] Register file-specific actions: `Diff`, `Diff later`, `Diff with previous version`, `Show log`, `Blame`, `Commit to main`.
   - [x] Update directory & background actions to be direct (max 1 level of submenus) so Windows Explorer renders all items without truncation.
2. In `cli/cmdinstall/installctxentries.go`:
   - [x] Change all informational commands (`50_list_releases`, `60_list_versions`, `30_find_next`, `91_docs`) from `Mode: ctxmodetype.Silent` to `Mode: ctxmodetype.Terminal` (`pwsh -NoExit ...`) so terminal output is visible and interactive.
   - [x] Replace complex nested `Repos` cascade with direct commands or Web UI repos launcher so no dropdown is broken.
3. [x] Add unit tests for file entries, background entries, and command generation.

---

## 2. Verification
- FastGate (Relative Paths, Nested Ifs, Boolean Guidelines, File Size): 100% PASS.
- `go test -v ./cmdinstall/...`: 100% PASS.
- Registry installation: `gitmap context-menu install` verified with 190 registry keys across `*`, `Directory`, `Directory\Background`.
