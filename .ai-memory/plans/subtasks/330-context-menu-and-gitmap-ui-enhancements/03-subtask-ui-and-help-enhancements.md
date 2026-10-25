# Subtask 330.3 — UI & Help Enhancements

> **Parent Plan:** `.ai-memory/plans/pending/330-context-menu-and-gitmap-ui-enhancements.md`  
> **Target Packages:** `cli/cmd/`, `cli/cmdui/`, `cli/cmdinstall/`  
> **Status:** Pending

---

## 1. Objectives
1. In `cli/cmd/roothelp_clusters.go`:
   - Add `context-menu (ctx)` under System/OS cluster so it shows in root `gitmap --help`.
2. In `cli/cmdinstall/install_help_sections.go` and `install_help_menu.go`:
   - Add `context-menu (ctx)` under install tools.
3. In `cli/cmdui/`:
   - Under Settings (`ui_assets_markup_settings_b.go`), add the Context Menu Integration card with Status badge, Install/Uninstall buttons, and a 1-click `Clean Duplicate PowerShell Menu` button.
   - Wire backend API `/api/context-menu/clean-powershell` and `/api/context-menu/status` in `ui_server.go`.
   - In `ui_assets_markup_ops.go`, highlight SSH nodes fleet commands (`gitmap nodes`, `gitmap pa --ssh`, `gitmap pe`).
   - Add direct link to the Show Log dual-pane UI.
4. Add unit tests verifying help registration and UI routes.
