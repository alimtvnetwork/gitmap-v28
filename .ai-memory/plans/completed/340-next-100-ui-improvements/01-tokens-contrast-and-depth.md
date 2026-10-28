# Subtask 01: Tokens, WCAG Contrast Calibration & 4-Plane Depth Hierarchy

- **Subtask**: 01
- **Task Slug**: `340-next-100-ui-improvements`
- **Status**: PENDING
- **Clusters**: Cluster 01, Cluster 02, Cluster 03 (Items 01–30)

---

## 1. Objectives

1. Adjust light mode `--muted-foreground` to `0 0% 38%` in `src/index.css` to achieve $\ge 4.5:1$ contrast against card surfaces.
2. Remove double opacity suppression (`opacity-50`, `opacity-40`) on muted text and icons across `AGYPrompts.tsx`, `SpecIndex.tsx`, `Projects.tsx`, and `DesignSystem.tsx`.
3. Configure 4-plane elevation classes in `src/index.css` and `tailwind.config.ts`:
   - Plane 0: Canvas Base (`#0A0F1D` dark / `#FBF9F6` light)
   - Plane 1: Raised Wells (`#10172A` dark / `#F3EFEA` light)
   - Plane 2: Surface Cards (`#172036` dark / `#FFFFFF` light)
   - Plane 3: Elevated Modals (`#1E2945` dark / `#FFFFFF` light)
4. Upgrade `src/components/ui/card.tsx` with top hairline lighting: `inset 0 1px 0 rgba(255, 255, 255, 0.12), 0 4px 20px rgba(0, 0, 0, 0.25)`.
5. Upgrade `src/components/ui/dialog.tsx` with Plane 3 elevated popover styling and calibrated backdrop blur.
