# Subtask 329.2 — Cross-Platform Context Menu Engine & Install Wizard Sub-Item

> **Parent Plan:** `.ai-memory/plans/pending/329-gitmap-context-menu-and-version-fix.md`  
> **Target Packages:** `cli/cmdinstall/`, `cli/cmdcontextmenu/`  
> **Status:** Completed

---

## 1. Objectives
1. Implement `gitmap install context-menu` and `gitmap uninstall context-menu` (supporting aliases `context-menu`, `contextmenu`, `ctx`).
2. Update the Git context menu structure in `cli/cmdinstall/installctxentries.go` to provide:
   - `Commit to main`
   - `Diff`
   - `Diff later`
   - `Diff with previous version`
   - `Show log`
   - `Branch comparison`
3. Add the Install Wizard sub-item for directory background / no-selection right-click:
   - Shows installed tools, active profile, custom installer status, and uninstall button.
4. Support cross-platform registration for Windows (`HKCU\Software\Classes`), Linux (`nautilus-python` and `.desktop`), and macOS (`Automator` workflows).
5. Add unit tests for context menu generation and entry validation.
