# Subtask 330.3 — UI & Help Enhancements

> **Parent Plan:** `.ai-memory/plans/completed/330-context-menu-and-gitmap-ui-enhancements.md`  
> **Target Packages:** `cli/cmd/`, `cli/cmdui/`, `cli/cmdinstall/`  
> **Status:** Completed

---

## 1. Objectives
1. In `cli/cmd/roothelp_clusters.go`:
   - [x] Add `context-menu (ctx)` under System/OS cluster so it shows in root `gitmap --help`.
2. In `cli/cmdinstall/install_help_sections.go` and `install_help_menu.go`:
   - [x] Add `context-menu (ctx)` under install tools.
3. In `cli/cmdui/`:
   - [x] Under Settings (`ui_assets_markup_settings_b.go`), add the Context Menu Integration card with Status badge, Install/Uninstall buttons, and a 1-click `Clean Duplicate PowerShell Menu` button.
   - [x] Wire backend API `/api/context-menu/clean-powershell` and `/api/context-menu/status` in `ui_server_contextmenu.go`.
   - [x] In `ui_assets_markup_ops.go`, highlight SSH nodes fleet commands (`gitmap nodes`, `gitmap pa --ssh`, `gitmap pe`).
   - [x] Add direct link to the Show Log dual-pane UI.
4. [x] Add unit tests verifying help registration and UI routes.

---

## 2. Verification
- FastGate (Relative Paths, Nested Ifs, Boolean Guidelines, File Size): 100% PASS.
- `go test -v ./cmd/... ./cmdinstall/... ./cmdui/...`: 100% PASS.
- Manual CLI verification: `gitmap --help` and `gitmap install --help` display `context-menu (ctx)` correctly.
