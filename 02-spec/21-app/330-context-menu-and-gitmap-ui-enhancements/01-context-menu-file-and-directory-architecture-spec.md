# 330.1 — File & Directory Context Menu Architecture Specification

> **Status:** Drafted  
> **Subsystem:** `cli/cmdinstall/`  
> **Parent Ledger:** `02-spec/21-app/330-context-menu-and-gitmap-ui-enhancements/00-master-audit-ledger.md`  
> **Rule Compliance:** Leaf Error Sovereignty, Positive Booleans (`is*`/`has*`), Functions <= 15 lines.

---

## 1. Problem Diagnosis & Registry Deficiencies

### Issue 1: Missing File Context Menu
In `assets/screenshots/context-menu-and-gitmap-ui-enhancements-04.png`, right-clicking on `.prettierignore` displays TortoiseGit but no GitMap menu.
- **Cause:** `cli/cmdinstall/installctx.go` registered only `Directory\Background` and `Directory` keys. `HKCU\Software\Classes\*\shell\gitmap` was omitted.
- **Solution:** Add `HKCU\Software\Classes\*\shell\gitmap` with file-specific commands.

### Issue 2: Explorer Cascade Truncation & Repos Dropdown Failure
In `assets/screenshots/context-menu-and-gitmap-ui-enhancements-01.png` and `02.png`, the background menu only rendered `Scan`, `Clone`, `Release`, `Repos`, and `Repos` failed to expand.
- **Cause:** Windows Explorer static cascades fail when nesting 3 levels deep (`gitmap -> Repos -> Go repos`), and Explorer limits the total number of subkeys parsed when `SubCommands` is empty.
- **Solution:** Flatten the hierarchy into direct high-value entries under `gitmap ▸`, keeping sub-cascades strictly 1 level deep.

### Issue 3: Silent Mode Invisible Execution
Commands like `List releases` and `List versions` had `Mode: ctxmodetype.Silent`, executing `pwsh ... | msg.exe`.
- **Cause:** `msg.exe` fails on Windows 10/11 or is disabled; `WindowStyle Hidden` prevents any terminal from showing.
- **Solution:** All output-producing commands must use `Mode: ctxmodetype.Terminal` (`pwsh -NoExit ...`) so terminal output is visible and interactive.

---

## 2. Context Menu Structure Contract

### 2.1 File Context Menu (`HKCU\Software\Classes\*\shell\gitmap`)
Applied when right-clicking on ANY file:
- `Diff`: Compare file against HEAD using Beyond Compare (`gitmap diff "%1"`)
- `Diff later`: Remember/stash this file path for diffing (`gitmap diff --stash-path "%1"`)
- `Diff with previous version`: Compare file against `HEAD~1` (`gitmap diff HEAD~1 "%1"`)
- `Show log`: Open Tri-pane Show Log UI focused on this file (`gitmap show-log "%1"`)
- `Blame`: View git blame for this file in terminal (`gitmap blame "%1"`)
- `Commit to main`: Quick commit (`gitmap commit -b main`)

### 2.2 Directory & Background Context Menu (`Directory\shell\gitmap` & `Directory\Background\shell\gitmap`)
Applied when right-clicking on a folder or folder background:
- `Show log`: Open Tri-pane Show Log UI for current repo/folder (`gitmap show-log "%V"`)
- `Diff (Working Tree)`: Launch Beyond Compare for entire working tree (`gitmap diff "%V"`)
- `Branch comparison`: Open branch comparison UI (`gitmap show-log --branch-compare "%V"`)
- `Commit to main`: Quick commit (`gitmap commit -b main`)
- `Pull`: Pull current repo (`gitmap pull`)
- `Push`: Push current repo (`gitmap push`)
- `Make Public / Make Private`: Repo visibility controls
- `GitMap Web UI`: Open local dashboard (`gitmap ui`)
- `Install Wizard`: Developer tools and profiles installer (`gitmap install --wizard`)
- `Clean Duplicate PowerShell`: Context menu customizer (`gitmap context-menu clean-powershell`)
- Submenus (Max 1 level deep):
  - `Release ▸`: `Release current`, `Release next (minor)`, `List releases (Terminal)`, `List versions (Terminal)`
  - `Scan ▸`: `Scan here`, `Rescan`
