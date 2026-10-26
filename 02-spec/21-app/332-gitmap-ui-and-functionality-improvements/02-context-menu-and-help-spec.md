# Specification: Context Menu Install Section, Repo List, and Terminal Help

> **Module:** `cli/cmdinstall/`, `cli/cmd/`  
> **Parent Program:** `332-gitmap-ui-and-functionality-improvements`  
> **File:** `02-spec/21-app/332-gitmap-ui-and-functionality-improvements/02-context-menu-and-help-spec.md`  
> **Status:** Draft / Active Specification  

---

## 1. Context Menu Install Section
- In `cli/cmdinstall/installctxentries.go`, add dedicated `Install ▸` submenu:
  ```go
  func installChildren() []ctxEntry {
      return []ctxEntry{
          {KeyName: "10_wizard", MUIVerb: "Install Wizard (Interactive)", Args: []string{"install", "--wizard"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
          {KeyName: "20_dev", MUIVerb: "Install Dev Workstation Profile", Args: []string{"install", "profile", "dev"}, Mode: ctxmodetype.Terminal},
          {KeyName: "30_tools", MUIVerb: "Install Development Tools (ls)", Args: []string{"install", "list"}, Mode: ctxmodetype.Terminal},
          {KeyName: "40_agy", MUIVerb: "Install Antigravity (AGY)", Args: []string{"install", "antigravity"}, Mode: ctxmodetype.Terminal},
          {KeyName: "50_clean_pwsh", MUIVerb: "Clean Duplicate PowerShell", Args: []string{"context-menu", "clean-powershell"}, Mode: ctxmodetype.Terminal},
      }
  }
  ```

---

## 2. Context Menu Repositories Section
- Provide a structured `Repositories ▸` submenu:
  ```go
  func repoChildren() []ctxEntry {
      return []ctxEntry{
          {KeyName: "10_web_ui", MUIVerb: "Open Repositories Web UI", Args: []string{"ui"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap},
          {KeyName: "20_list", MUIVerb: "List Tracked Repositories", Args: []string{"list"}, Mode: ctxmodetype.Terminal},
          {KeyName: "30_status", MUIVerb: "Status All Repositories", Args: []string{"status-all"}, Mode: ctxmodetype.Terminal},
          {KeyName: "40_pull_all", MUIVerb: "Pull All Repositories", Args: []string{"pull-all"}, Mode: ctxmodetype.Terminal},
      }
  }
  ```

---

## 3. Terminal Help Expansion for Merge & Compare Commands
- In `cli/cmd/rootusage_groups.go` (`printGroupGitOps` and new `printGroupMergeCompare`):
  - Add help entries:
    - `merge-right, mr <left> <right>` — One-way file-level merge that writes changes into RIGHT.
    - `merge-left, ml <left> <right>` — One-way file-level merge that writes changes into LEFT.
    - `merge-both, mb <left> <right>` — Two-way bidirectional merge synchronizing both directories.
    - `compare, cmp <b1> <b2> [b3...]` — Compare branches via Beyond Compare or external diff tool.
    - `show-log, showlog [path]` — Interactive dual-pane Git log viewer and branch comparison UI.
- In `cli/cmd/roothelp_clusters.go`:
  - Add `merge-right (mr)`, `merge-left (ml)`, `compare (cmp)`, `show-log` to Core Repo Operations.
