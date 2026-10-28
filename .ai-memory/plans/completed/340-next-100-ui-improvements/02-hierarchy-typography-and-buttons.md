# Subtask 02: Single-Accent Von Restorff Rule, Typography Stacks & Button System

- **Subtask**: 02
- **Task Slug**: `340-next-100-ui-improvements`
- **Status**: PENDING
- **Clusters**: Cluster 04, Cluster 05, Cluster 06 (Items 31–60)

---

## 1. Objectives

1. Add Google Font `@import` for Poppins (`wght@300;400;500;600`) and JetBrains Mono (`wght@400;500`) in `src/index.css`.
2. Update `tailwind.config.ts` font families: `sans: ["Poppins", "system-ui", "sans-serif"]`, `mono: ["JetBrains Mono", "monospace"]`.
3. Standardize `src/components/ui/button.tsx`:
   - 4-size hierarchy: `sm: 36px` (`h-9`), `default: 44px` (`h-11`), `lg: 52px` (`h-13`), `icon: 44px` (`size-11`).
   - Add `glass` variant (`backdrop-blur-md border border-white/10`) and `solid` variant.
   - Add tactile active scale feedback (`active:scale-[0.98]`).
4. Enforce Single-Accent Von Restorff Rule in `src/pages/Index.tsx`:
   - Reserve primary accent exclusively for `Get Started`.
   - Neutralize `View Commands` to `outline` / `ghost`.
5. Isolate focal metrics in `src/pages/Pipeline.tsx` and neutralize multi-chip saturation in `src/pages/Settings.tsx`.
