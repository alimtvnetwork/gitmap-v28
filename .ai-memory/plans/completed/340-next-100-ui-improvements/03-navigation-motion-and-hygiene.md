# Subtask 03: Navigation Physics, Hardware Motion & Component Deconstruction

- **Subtask**: 03
- **Task Slug**: `340-next-100-ui-improvements`
- **Status**: PENDING
- **Clusters**: Cluster 07, Cluster 08, Cluster 09, Cluster 10 (Items 61–100)

---

## 1. Objectives

1. Standardize `DocsLayout.tsx` header to 72px glassmorphic geometry (`h-[72px] sticky top-0 z-50 backdrop-blur-xl`) with 12px scroll trigger elevation.
2. Upgrade `cli/cmdui/ui_assets_theme.go`:
   - Add `@media (prefers-reduced-motion: reduce)` block to clamp animations.
   - Refactor `.skip-link` transition from layout-shifting `top` to `transform: translateY(...)`.
   - Upgrade `.btn` classes to 44px hit bounds and add active tactile feedback.
3. Align main documentation container to canonical 1280px (`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8`).
4. Deconstruct `src/components/docs/DocsSidebar.tsx` into modular navigation sub-components (`DocsSidebarNav.tsx`, `DocsSidebarFooter.tsx`) adhering to the $\le 80$–100 lines limit.
