# Subtask 330.2 — PowerShell Duplicate Context Menu Cleaner

> **Parent Plan:** `.ai-memory/plans/completed/330-context-menu-and-gitmap-ui-enhancements.md`  
> **Target Packages:** `cli/cmdcontextmenu/`  
> **Status:** Completed

---

## 1. Objectives
1. Implement `cli/cmdcontextmenu/clean_powershell.go`:
   - Detect duplicate PowerShell keys in `HKCR\Directory\Background\shell` and `HKCR\Directory\shell` (`PowerShellMenu` vs `PowerShell7x64`).
   - Command `gitmap context-menu clean-powershell [--dry-run|--restore]`.
   - When cleaning, set `LegacyDisable` and `ProgrammaticAccessOnly` in `HKCU\Software\Classes\Directory\Background\shell\PowerShellMenu` and `Directory\shell\PowerShellMenu`.
   - When restoring, remove the `HKCU` override keys.
2. Expose the functionality in `RunContextMenu`:
   - `gitmap context-menu clean-powershell`
   - `gitmap context-menu customize`
3. Add unit tests for registry key generation and clean/restore logic in `clean_powershell_test.go` and `contextmenu_test.go`.

---

## 2. Verification
- FastGate (Relative Paths, Nested Ifs, Boolean Guidelines, File Size): 100% PASS in 0.015s.
- `go test -v ./cmdcontextmenu/...`: 100% PASS (14 tests passed).
