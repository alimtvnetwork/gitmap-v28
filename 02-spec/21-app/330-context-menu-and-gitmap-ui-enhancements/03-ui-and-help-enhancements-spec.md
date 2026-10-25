# 330.3 — UI & Help Enhancements Specification

> **Status:** Drafted  
> **Subsystem:** `cli/cmd/`, `cli/cmdui/`, `cli/cmdinstall/`  
> **Parent Ledger:** `02-spec/21-app/330-context-menu-and-gitmap-ui-enhancements/00-master-audit-ledger.md`  
> **Rule Compliance:** Leaf Error Sovereignty, Positive Booleans (`is*`/`has*`), Functions <= 15 lines.

---

## 1. CLI Help Integration

1. **Root CLI Help (`cli/cmd/roothelp_clusters.go`):**
   Add `context-menu (ctx)` under `5. System, OS & Developer Tooling`:
   ```go
   {Command: "context-menu (ctx)", Description: "OS right-click context menu (install, uninstall, clean-powershell)"},
   ```
2. **Installer Help (`cli/cmdinstall/install_help_sections.go`):**
   Add `context-menu (ctx)` under `Custom Scripts & Archives` or dedicated `OS Integration` section:
   ```go
   {Command: "context-menu (ctx)", Description: "Register Explorer, Nautilus & Finder context menus"},
   ```

---

## 2. Web UI Enhancements (`cli/cmdui/`)

1. **Context Menu & PowerShell Cleaner in Settings:**
   Under the Settings tab (`ui_assets_markup_settings_b.go`):
   - Add **Context Menu Integration Card**:
     - Status badge: Active / Inactive
     - Buttons: `Install Context Menu`, `Uninstall Context Menu`
     - Duplicate Cleaner button: `Clean Duplicate PowerShell Menu` (triggers `/api/context-menu/clean-powershell`)
2. **SSH Fleet Commands Card:**
   Ensure `ui_assets_markup_ops.go` prominently displays the fleet commands:
   - `gitmap pa --ssh`: Enqueue pull-all across SSH fleet
   - `gitmap nodes`: Inspect all cluster nodes
   - `gitmap nodes ping`: Dual-stack ICMP/TCP ping test
3. **Log Viewer Navigation:**
   Add `Git Log & Diff` navigation button linking directly to the dual-pane log viewer (`gitmap show-log`).
