# Subtask 330.1 — File & Directory Context Menu Architecture

> **Parent Plan:** `.ai-memory/plans/pending/330-context-menu-and-gitmap-ui-enhancements.md`  
> **Target Packages:** `cli/cmdinstall/`  
> **Status:** Pending

---

## 1. Objectives
1. In `cli/cmdinstall/installctx.go`:
   - Add `HKCU\Software\Classes\*\shell\gitmap` so GitMap appears on right-clicking ANY file.
   - Register file-specific actions: `Diff`, `Diff later`, `Diff with previous version`, `Show log`, `Blame`, `Commit to main`.
   - Update directory & background actions to be direct (max 1 level of submenus) so Windows Explorer renders all items without truncation.
2. In `cli/cmdinstall/installctxentries.go`:
   - Change all informational commands (`50_list_releases`, `60_list_versions`, `30_find_next`, `91_docs`) from `Mode: ctxmodetype.Silent` to `Mode: ctxmodetype.Terminal` (`pwsh -NoExit ...`) so terminal output is visible and interactive.
   - Replace complex nested `Repos` cascade with direct commands or Web UI repos launcher so no dropdown is broken.
3. Add unit tests for file entries, background entries, and command generation.
