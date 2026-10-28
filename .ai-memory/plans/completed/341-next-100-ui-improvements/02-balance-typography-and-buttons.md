# Subtask 02: Visual Balance, Typography & Button Standards

**Parent Task:** `341-next-100-ui-improvements`  
**Clusters Covered:** Cluster 04, Cluster 05, Cluster 06  
**Target Files:**
- `src/components/docs/DocsLayout.tsx`
- `src/components/docs/CopyPaletteButton.tsx`
- `src/components/docs/troubleshooting/TroubleshootingFilterButtons.tsx`
- `src/components/projects/TypeBadge.tsx`
- `src/pages/Changelog.tsx`

---

## Action Checklist
- [x] Enforce 44px min touch bounds on dark/light radio toggle buttons in `DocsLayout.tsx` (`min-h-[44px] min-w-[44px]` with flex centering).
- [x] Add active depression scaling (`active:scale-[0.98]`) to `CopyPaletteButton.tsx` and filter pill buttons.
- [x] Enforce Von Restorff single-accent isolation across filter buttons and category tabs.
- [x] De-saturate project type badges in `TypeBadge.tsx` to clean, muted neutral pills with soft borders.
- [x] Upgrade "Expand all" / "Collapse all" and version filter buttons in `Changelog.tsx` to 44px touch targets with active tactile feedback.
