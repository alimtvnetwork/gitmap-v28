# Subtask 03: Navigation Physics, Motion & Component Hygiene

**Parent Task:** `341-next-100-ui-improvements`  
**Clusters Covered:** Cluster 07, Cluster 08, Cluster 09, Cluster 10  
**Target Files:**
- `cli/cmdui/ui_assets_theme.go`
- `cli/cmdui/ui_assets_markup_misc.go`
- `src/components/ui/dialog.tsx`
- `src/components/ui/sheet.tsx`
- `src/components/ui/sidebar.tsx`

---

## Action Checklist
- [ ] Upgrade embedded web UI `header` in `ui_assets_theme.go` to canonical 72px min-height with glassmorphic `backdrop-filter: blur(12px)`.
- [ ] Add `will-change: transform` and GPU composite acceleration (`transform-gpu`) to modal dialogs, sheets, and sidebar transitions.
- [ ] Add `env(safe-area-inset-bottom)` and `env(safe-area-inset-top)` support across CLI UI headers and floating toasts.
- [ ] Ensure all modified component modules adhere to strict length limits (<= 100 lines per function/component).
