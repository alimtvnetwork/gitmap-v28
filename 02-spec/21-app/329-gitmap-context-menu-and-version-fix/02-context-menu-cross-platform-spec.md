# 329.2 — Cross-Platform Context Menu & Install Wizard Specification

> **Status:** Drafted  
> **Subsystem:** `cli/cmdinstall/`, `cli/cmdcontextmenu/`  
> **Parent Ledger:** `02-spec/21-app/329-gitmap-context-menu-and-version-fix/00-master-audit-ledger.md`  
> **Rule Compliance:** Leaf Error Sovereignty, Positive Booleans (`is*`/`has*`), Functions <= 15 lines.

---

## 1. Overview & Operational Requirements

The user requested a seamless right-click context menu integration for Windows, Linux/Ubuntu, and macOS, installable via:
```bash
gitmap install context-menu
# Also supported:
gitmap install ctx
gitmap context-menu install
gitmap context-menu uninstall
```

### Context Menu Feature Hierarchy (TortoiseGit-Inspired)
When clicking on a file or directory within a Git repository:
- **Git Section:**
  1. `Commit to main`: Quick commit / push dialog or terminal execution (`gitmap commit -b main`).
  2. `Diff`: Show unstaged / working-tree diff against HEAD.
  3. `Diff later`: Stashes / pins the selected file for later two-file or multi-revision diffing.
  4. `Diff with previous version`: Compares current file against `HEAD~1` or previous commit.
  5. `Show log`: Opens the GitMap Show Log UI for the selected file or folder.
  6. `Branch comparison`: Compares current branch / file against other local or remote branches.

When right-clicking on an empty space / folder background (or when no file is selected):
- **Install Section (Sub-item):**
  - Opens the **GitMap Install Wizard**:
    - Displays currently installed components (e.g. Git, Go, Node, Beyond Compare, VS Code).
    - Lists active profiles (e.g. default, work-env, minimal).
    - Highlights custom installer status.
    - Provides a one-click **Uninstall** button capable of cleaning up GitMap components and context menu registrations across all operating systems.

---

## 2. Platform Architecture & Registration

### 2.1 Windows Registry Architecture
Registrations are scoped to `HKCU\Software\Classes` (no elevation or Administrator prompt required):
- `HKCU\Software\Classes\*\shell\GitMap` (All files)
- `HKCU\Software\Classes\Directory\shell\GitMap` (Directories)
- `HKCU\Software\Classes\Directory\Background\shell\GitMap` (Directory Background / White Space)

Submenu items use the Windows Explorer cascading verb pattern:
- `MUIVerb`: Display name (e.g. `GitMap ▸`)
- `SubCommands`: Blank (indicates nested `shell\` hierarchy)
- `Icon`: `<gitmap_exe>,0`

Command templates:
- For GUI/Visual tools: `"<gitmap_exe>" <subcommand> "%V"` (or `"%1"`)
- For interactive terminal actions: `pwsh -NoExit -NoProfile -Command "Set-Location '%V'; gitmap <subcommand>"`

### 2.2 Linux / Ubuntu Architecture
- **Nautilus Extension:** Written to `~/.local/share/nautilus-python/extensions/gitmap-nautilus.py`.
- **Desktop Actions:** Written to `~/.local/share/file-manager/actions/gitmap-*.desktop` and KDE Dolphin service menus in `~/.local/share/kio/servicemenus/gitmap.desktop`.
- Supports both background and file-selection triggers.

### 2.3 macOS Architecture
- **Automator Quick Actions:** Workflows installed in `~/Library/Services/GitMap - <Action>.workflow`.
- Configured to accept "files or folders" in "Finder.app".

---

## 3. Install Wizard Sub-Item & UI Architecture

### CLI & Webview Entrypoints
- `gitmap install --wizard` (or `gitmap context-menu wizard`)
- `gitmap ui install`

### Features of the Install Wizard
1. **Installed Tool Grid:** Queries installed status of developer tools (Git, Go, Node, Beyond Compare, VS Code, Antigravity) with green/gray status badges.
2. **Profile Manager:** Displays active profile and available profiles.
3. **Context Menu Status:** Shows whether Explorer/Nautilus/Finder context menus are active with an Install/Reinstall button.
4. **Unified Uninstall:** Cleanly uninstalls the context menu and removes registered shell integrations with confirmation.
