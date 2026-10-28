# Subtask 01: Tokens, Contrast & Depth Materiality

**Parent Task:** `341-next-100-ui-improvements`  
**Clusters Covered:** Cluster 01, Cluster 02, Cluster 03  
**Target Files:**
- `src/components/terminal/TerminalDrawer.tsx`
- `src/components/terminal/TerminalView.tsx`
- `src/pages/Watch.tsx`
- `src/pages/Database.tsx`
- `src/pages/Projects.tsx`

---

## Action Checklist
- [x] Eliminate hardcoded hex codes in terminal themes, substituting semantic tokens or CSS variables.
- [ ] Calibrate table cell badges in `Database.tsx` (`emerald`, `rose`, `purple`, `sky`, `amber`) to WCAG 2.2 AA compliant tokens (`text-emerald-700 dark:text-emerald-400`, etc.).
- [ ] Eliminate double opacity suppression (`text-muted-foreground/60 italic` on `NULL` indicators, `opacity-50`).
- [ ] Upgrade database table wrapper to Plane 1 recessed well base (`bg-plane-well`).
- [ ] Upgrade modal inspection containers to Plane 3 elevation tokens (`dark:bg-plane-elevated`) with top hairline highlight `inset 0 1px 0 rgba(255,255,255,0.12)`.
