# Plan 330 — Context Menu and GitMap UI Enhancements

> **Status:** Completed  
> **Sequence ID:** 330  
> **Slug:** `context-menu-and-gitmap-ui-enhancements`  
> **Specification:** `02-spec/21-app/330-context-menu-and-gitmap-ui-enhancements/`  
> **Subtasks Directory:** `.ai-memory/plans/subtasks/330-context-menu-and-gitmap-ui-enhancements/`

---

## 1. Executive Summary & Problem Context
1. **File Context Menu Gap:** `gitmap` was omitted from `HKCU\Software\Classes\*\shell\gitmap`, preventing right-click file actions. Fixed by registering `*\shell\gitmap` with 6 file commands (`Diff`, `Diff later`, `Diff with previous`, `Show log`, `Blame`, `Commit to main`).
2. **Explorer Cascade Truncation:** Deep nesting prevented `Repos` and downstream items from displaying; silent commands failed invisibly. Fixed by flattening directory/background menus to direct top-level entries and converting all informational actions to `Terminal` mode (`pwsh -NoExit`).
3. **PowerShell Duplicate Remediation:** Windows users with both PowerShell 7 and legacy PowerShell need a 1-click cleaner in GitMap CLI and Web UI. Implemented `gitmap context-menu clean-powershell` writing `LegacyDisable` to `PowerShellMenu` in `HKCU`.
4. **Help & UI Visibility:** `context-menu` was omitted from root/install help menus; Web UI lacked Settings customizer and SSH command cards. Fixed by registering in `roothelp_clusters.go`, `install_help_sections.go`, and adding the Context Menu & Shell Settings card and endpoints in `cmdui`.

---

## 2. Subtasks Roadmap

| Subtask ID | File | Description | Owner | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | `01-subtask-file-and-directory-context-menu.md` | Register `*\shell\gitmap`, flatten directory menus, convert silent actions to interactive terminal mode (`pwsh -NoExit`). | Worker 01 | Completed |
| **Subtask-02** | `02-subtask-powershell-duplicate-cleaner.md` | Implement `gitmap context-menu clean-powershell` to detect and suppress duplicate `PowerShellMenu` keys via user registry overrides. | Worker 02 | Completed |
| **Subtask-03** | `03-subtask-ui-and-help-enhancements.md` | Add `context-menu` to root & install help; add Settings context menu customizer card, SSH fleet cards, and log viewer linkage in Web UI. | Worker 03 | Completed |

---

## 3. Verification & Acceptance Criteria
- [x] Right-clicking any file displays GitMap with `Diff`, `Diff later`, `Diff with previous`, `Show log`, `Blame`.
- [x] Right-clicking folder/background displays clean, un-truncated direct actions without broken dropdowns.
- [x] `gitmap context-menu clean-powershell` successfully suppresses duplicate `PowerShellMenu`.
- [x] `gitmap --help` and `gitmap install --help` prominently document `context-menu (ctx)`.
- [x] All informational buttons (`list-releases`, `list-versions`) open an interactive terminal.
- [x] Zero absolute paths; FastGate passes 100%.
- [x] All unit test suites pass 100% (`cmdcontextmenu`, `cmdinstall`, `cmdui`, `cmd`).
