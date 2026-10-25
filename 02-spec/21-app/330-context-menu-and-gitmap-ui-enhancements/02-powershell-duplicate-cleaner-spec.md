# 330.2 — PowerShell Duplicate Context Menu Cleaner Specification

> **Status:** Drafted  
> **Subsystem:** `cli/cmdcontextmenu/`  
> **Parent Ledger:** `02-spec/21-app/330-context-menu-and-gitmap-ui-enhancements/00-master-audit-ledger.md`  
> **Rule Compliance:** Leaf Error Sovereignty, Positive Booleans (`is*`/`has*`), Functions <= 15 lines.

---

## 1. Problem Statement

Windows users with PowerShell 7 and Windows PowerShell frequently encounter duplicate context menu entries (as shown in `assets/screenshots/context-menu-and-gitmap-ui-enhancements-01.png` and `03.png`):
- `PowerShell 7 ▸` (`PowerShell7x64` key in `HKLM\Software\Classes\Directory\Background\shell`)
- `PowerShell ▸` (`PowerShellMenu` key in `HKLM\Software\Classes\Directory\Background\shell`)

Users need an automated way in GitMap CLI and Web UI to eliminate the duplicate legacy entry without requiring administrative privileges or breaking system files.

---

## 2. Registry Override Mechanism

In Windows Explorer, keys under `HKCU\Software\Classes\` take precedence over `HKLM\Software\Classes\`.
- Adding the `LegacyDisable` (REG_SZ) or `ProgrammaticAccessOnly` (REG_SZ) value to:
  - `HKCU\Software\Classes\Directory\Background\shell\PowerShellMenu`
  - `HKCU\Software\Classes\Directory\shell\PowerShellMenu`
- Immediately hides the redundant `PowerShell` entry from the Explorer context menu while keeping `PowerShell 7` visible and active.
- To restore: simply delete the override key from `HKCU`.

---

## 3. CLI Command Specification

```bash
# Detect and clean duplicate PowerShell menus:
gitmap context-menu clean-powershell

# Preview without applying:
gitmap context-menu clean-powershell --dry-run

# Restore default (unhide legacy PowerShell):
gitmap context-menu clean-powershell --restore
```
