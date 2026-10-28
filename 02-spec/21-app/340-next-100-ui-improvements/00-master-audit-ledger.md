# Master Audit Ledger: Next 100 UI Improvements

- **Sequence**: 340
- **Task Slug**: `next-100-ui-improvements`
- **Specification Status**: IN_PROGRESS
- **Execution Target**: Frontend React Web Application (`src/`) & Fleet Web UI (`cli/cmdui/`)

---

## 1. Executive Summary & Improvement Catalog

This audit operationalizes the **100 concrete UI improvements** defined across the 10 core design system clusters:
- **Cluster 01**: Color Tokens, HSL Architecture & Zero Hardcoded Hex (Items 01–10)
- **Cluster 02**: Contrast Ratios & WCAG 2.2 AA/AAA Compliance (Items 11–20)
- **Cluster 03**: 4-Plane Depth Hierarchy & Materiality (Items 21–30)
- **Cluster 04**: 60/30/10 Visual Balance & Single-Accent Von Restorff Rule (Items 31–40)
- **Cluster 05**: Typography Stacks & Pure Live DOM Text Mandate (Items 41–50)
- **Cluster 06**: Button System Standards & Hit Target Geometry (Items 51–60)
- **Cluster 07**: Header, Navigation & Safe-Region Physics (Items 61–70)
- **Cluster 08**: Spacing Scales, Containers & Band Rhythm (Items 71–80)
- **Cluster 09**: CSS3 Hardware-Accelerated Motion & Accessibility (Items 81–90)
- **Cluster 10**: Frontend Code Hygiene & React Component Deconstruction (Items 91–100)

---

## 2. Specification Structure

- [`01-tokens-contrast-and-depth-spec.md`](01-tokens-contrast-and-depth-spec.md): Clusters 01, 02, and 03.
- [`02-hierarchy-typography-and-buttons-spec.md`](02-hierarchy-typography-and-buttons-spec.md): Clusters 04, 05, and 06.
- [`03-navigation-motion-and-hygiene-spec.md`](03-navigation-motion-and-hygiene-spec.md): Clusters 07, 08, 09, and 10.

---

## 3. The 100 UI Improvement Master Catalog

### Cluster 01: Color Tokens, HSL Architecture & Zero Hardcoded Hex
1. `UI-001`: Purge hardcoded hex values in `TerminalDrawer.tsx` (`#5e2750`, `#300a24`, `#25071c`, `#e95420`, `#0b3c7b`, `#012456`).
2. `UI-002`: Purge hardcoded hex values in `TerminalView.tsx` theme mappings into CSS variable tokens.
3. `UI-003`: Eliminate arbitrary inline HSL literals `style={{ backgroundColor: "hsl(220 25% 8%)" }}` in `DesignSystem.tsx`.
4. `UI-004`: Clean leftover template hex styling in `App.css` (`#646cffaa`, `#61dafbaa`, `#888`).
5. `UI-005`: Replace `text-emerald-400` / `text-green-400` with semantic token `text-accent-success` across `ProjectDetection.tsx` and `Database.tsx`.
6. `UI-006`: Replace `bg-purple-500/15 text-purple-400` with semantic brand tokens in `Database.tsx`.
7. `UI-007`: Replace `bg-amber-500/15 text-amber-400` with semantic warning tokens in `Database.tsx`.
8. `UI-008`: Map syntax highlighting tokens in `src/index.css` to semantic CSS custom properties.
9. `UI-009`: Standardize all CSS variable declarations in `:root` and `.dark` into space-separated `H S% L%` values.
10. `UI-010`: Eliminate ad-hoc raw Tailwind palette utilities (`zinc-400`, `zinc-500`) in favor of `muted-foreground` and `border`.

### Cluster 02: Contrast Ratios & WCAG 2.2 AA/AAA Compliance
11. `UI-011`: Adjust light mode `--muted-foreground` from `0 0% 42%` to `0 0% 38%` to guarantee $\ge 4.5:1$ contrast against card backgrounds (`0 0% 97%`).
12. `UI-012`: Remove double opacity suppression (`text-muted-foreground opacity-50`) in `AGYPrompts.tsx` line 359.
13. `UI-013`: Remove `opacity-40` suppression on muted labels in `SpecIndex.tsx` line 128.
14. `UI-014`: Remove `opacity-50` suppression on workspace info in `Projects.tsx` line 172.
15. `UI-015`: Remove `opacity-50` multiplier on color preview swatch text in `DesignSystem.tsx`.
16. `UI-016`: Refactor low-contrast status badge `text-zinc-400` on light background in `Settings.tsx` to `text-muted-foreground`.
17. `UI-017`: Provide dual light/dark syntax highlighting variable classes so light mode code blocks achieve $\ge 4.5:1$ contrast.
18. `UI-018`: Enforce $\ge 3:1$ contrast on all UI icons, badge borders, and input boundaries against their immediate parent surface.
19. `UI-019`: Guarantee dark high-contrast mode achieves $\ge 7:1$ (AAA) on all primary reading surfaces.
20. `UI-020`: Verify fast relative luminance formula ($L = 0.299 R + 0.587 G + 0.114 B$) across all dynamic badge pairings.

### Cluster 03: 4-Plane Depth Hierarchy & Materiality
21. `UI-021`: Establish Plane 0 Canvas Base token in `index.css` (`#0A0F1D` dark / `#FBF9F6` light).
22. `UI-022`: Establish Plane 1 Raised Wells token for sidebars and section bands (`#10172A` dark / `#F3EFEA` light).
23. `UI-023`: Establish Plane 2 Surface Cards token for standard content cards and tables (`#172036` dark / `#FFFFFF` light).
24. `UI-024`: Establish Plane 3 Elevated Modals token for floating dialogs and dropdowns (`#1E2945` dark / `#FFFFFF` light).
25. `UI-025`: Upgrade `card.tsx` with top-edge hairline lighting: `box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12), 0 4px 20px rgba(0, 0, 0, 0.25)`.
26. `UI-026`: Upgrade `dialog.tsx` `DialogContent` with Plane 3 elevation and hairline refraction.
27. `UI-027`: Calibrate `DialogOverlay` backdrop blur to `bg-black/60 backdrop-blur-sm`.
28. `UI-028`: Fix Plane 1 lightness inversion in `DocsSidebar.tsx` where sidebar was darker than canvas base.
29. `UI-029`: Add layered ambient radial glow to hero terminal container in `Index.tsx`.
30. `UI-030`: Add restrained hover elevation lift (`hover:-translate-y-0.5 hover:shadow-[0_8px_24px_rgba(0,0,0,0.3)]`) to `FeatureCard.tsx`.

### Cluster 04: 60/30/10 Visual Balance & Single-Accent Von Restorff Rule
31. `UI-031`: Enforce Single-Accent Von Restorff rule on `Index.tsx` hero: exactly one primary accent button (`Get Started`).
32. `UI-032`: De-saturate secondary hero action `View Commands` in `Index.tsx` to `outline` or `ghost`.
33. `UI-033`: De-saturate competing tab header text in `Index.tsx` hero fold to `text-muted-foreground`.
34. `UI-034`: Isolate focal metric in `Pipeline.tsx`: elevate `Pipeline Status` card with primary accent; neutralize secondary cards.
35. `UI-035`: Remove saturated primary chips across all section headers in `Settings.tsx`, preserving primary accent solely for `Save Changes`.
36. `UI-036`: Replace full-tint callout background washes in `CommandCard.tsx` with neutral cards and a 2px left hairline accent.
37. `UI-037`: Bound saturated accent color coverage to strictly $\le 10\%$ of visible viewport pixels across all dashboard pages.
38. `UI-038`: Reserve 60% dominant neutral background across main layouts.
39. `UI-039`: Allocate 30% structural surface density for cards, dividers, and typography containers.
40. `UI-040`: Eliminate multi-accent badge collisions on list view items in `Projects.tsx`.

### Cluster 05: Typography Stacks & Pure Live DOM Text Mandate
41. `UI-041`: Load `Poppins` (`wght@300;400;500;600`) font family in `index.css`.
42. `UI-042`: Load `JetBrains Mono` (`wght@400;500`) font family in `index.css`.
43. `UI-043`: Configure `tailwind.config.ts` font families: `heading` (Ubuntu), `sans`/body (Poppins), `mono` (JetBrains Mono), `terminal` (Ubuntu Mono).
44. `UI-044`: Apply fluid typography clamp scaling on main display titles: `text-[clamp(1.85rem,3.5vw,3rem)]`.
45. `UI-045`: Apply optical tracking calibration (`tracking-[-0.02em]`) on display headings in `Index.tsx` and `Dashboard.tsx`.
46. `UI-046`: Constrain body prose reading measure to `max-w-[68ch]` with `leading-[1.62]` in `FeatureCard.tsx` and docs layouts.
47. `UI-047`: Enforce copy length budgets: lead paragraphs $\le 28$ words; card descriptions $\le 32$ words.
48. `UI-048`: Differentiate table monospace pills (JetBrains Mono) from descriptions (Poppins) in `Install.tsx` and `Dashboard.tsx`.
49. `UI-049`: Standardize uppercase tracked eyebrow styling: `text-[10px] font-semibold uppercase tracking-[0.14em]` in `ProjectDetailDialog.tsx`.
50. `UI-050`: Strictly uphold Pure Live DOM Text Mandate: verify zero baked raster text in PNG/JPEG images.

### Cluster 06: Button System Standards & Hit Target Geometry
51. `UI-051`: Add canonical `glass` button variant in `button.tsx` (`backdrop-blur-md border border-white/10`).
52. `UI-052`: Add canonical `solid` button variant in `button.tsx` (`bg-primary text-primary-foreground`).
53. `UI-053`: Upgrade default button height to $44\text{px}$ (`h-11 px-5 text-sm`) satisfying mobile touch target standards.
54. `UI-054`: Standardize `icon` button sizing to $44\times 44\text{px}$ (`size-11`) in `button.tsx`.
55. `UI-055`: Add tactile active depression feedback (`active:scale-[0.98]`) across all button variants.
56. `UI-056`: Expand sub-44px touch targets on database tab buttons in `Database.tsx` lines 249–277.
57. `UI-057`: Upgrade dense admin `.btn` classes in `cli/cmdui/ui_assets_theme.go` to minimum 44px hit bounds on touch viewports.
58. `UI-058`: Add active scale feedback (`transform: scale(0.98)`) to embedded `cmdui` `.btn` styles.
59. `UI-059`: Introduce `.btn-glass` and `.btn-outline` token classes in `ui_assets_theme.go`.
60. `UI-060`: Ensure all clickable icon triggers throughout the interface provide a minimum $44\times 44\text{px}$ hit bounding box.

### Cluster 07: Header, Navigation & Safe-Region Physics
61. `UI-061`: Standardize `DocsLayout.tsx` header to fixed 72px glassmorphic geometry (`h-[72px] sticky top-0 z-50`).
62. `UI-062`: Elevate header blur to `backdrop-blur-xl` in `DocsLayout.tsx`.
63. `UI-063`: Implement 12px scroll trigger border and shadow elevation transition on docs header.
64. `UI-064`: Implement safe-region geometry buffer (`pad = 14px`) on navigation dropdown menus.
65. `UI-065`: Implement 220ms debounce timer on navigation menu close in `navigation-menu.tsx`.
66. `UI-066`: Ensure instant dismissal ($0\text{ms}$) of navigation menus on Escape key or window scroll.
67. `UI-067`: Upgrade embedded `cmdui` header to 72px glassmorphic geometry in `ui_assets_theme.go`.
68. `UI-068`: Support responsive mobile drawer safe-region handling with backdrop scrim in `DocsLayout.tsx`.
69. `UI-069`: Add active pill indicator physics on current navigation item in `DocsSidebar.tsx`.
70. `UI-070`: Enforce sticky navigation z-index layering hierarchy (`z-50` header, `z-40` sidebar, `z-30` sticky table headers).

### Cluster 08: Spacing Scales, Containers & Band Rhythm
71. `UI-071`: Align main docs container in `DocsLayout.tsx` to standard 1280px (`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8`).
72. `UI-072`: Eliminate non-modular `max-w-5xl` constraint on documentation pages.
73. `UI-073`: Standardize corner radius hierarchy: 8px (`rounded-lg`) inputs, 12px (`rounded-xl`) buttons, 20px (`rounded-2xl`) cards.
74. `UI-074`: Purge ad-hoc `rounded-[3px]` and `rounded-sm` from `DocsLayout.tsx` and `DocsSidebar.tsx`.
75. `UI-075`: Standardize regular card padding to 24px (`p-6` / `var(--sp-6)`) across all pages.
76. `UI-076`: Standardize hero/featured card padding to 32px (`p-8` / `var(--sp-8)`).
77. `UI-077`: Align embedded `cmdui` card padding to modular 24px (`var(--sp-6)`) in `ui_assets_theme.go`.
78. `UI-078`: Enforce consistent 8-point vertical section padding rhythm (`py-8 lg:py-12`).
79. `UI-079`: Calibrate gap spacing in grid cards to `gap-6` (24px) on desktop viewports.
80. `UI-080`: Ensure harmonic spacing tokens `var(--sp-1)` through `var(--sp-12)` are uniformly used across CSS files.

### Cluster 09: CSS3 Hardware-Accelerated Motion & Accessibility
81. `UI-081`: Eliminate layout-triggering `top` transition on `.skip-link` in `ui_assets_theme.go`, refactoring to `transform: translateY(...)`.
82. `UI-082`: Eliminate layout-triggering `padding-left` transition on `.docs-h3:hover` in `index.css`.
83. `UI-083`: Add mandatory `@media (prefers-reduced-motion: reduce)` block in `ui_assets_theme.go`.
84. `UI-084`: Clamp all keyframe animations (`pulseGlow`, `appEntrance`) to `0.01ms` or `none` when reduced motion is preferred.
85. `UI-085`: Add `will-change: transform` on interactive card hover states in `index.css`.
86. `UI-086`: Standardize transform duration to 150ms and color/opacity duration to 200ms with cubic-bezier easings.
87. `UI-087`: Add smooth spring physics on modal dialog entry and exit in `dialog.tsx`.
88. `UI-088`: Hardware-accelerate sidebar slide-over transitions using `transform: translateX(...)`.
89. `UI-089`: Restrict all animated transitions strictly to GPU-composited properties (`transform`, `opacity`).
90. `UI-090`: Provide visible keyboard focus rings (`focus-visible:ring-2 focus-visible:ring-primary`) on all interactive controls.

### Cluster 10: Frontend Code Hygiene & React Component Deconstruction
91. `UI-091`: Deconstruct monolithic `Database.tsx` (763 lines) into `DatabaseToolbar.tsx`, `DatabaseTableGrid.tsx`, and `CellInspectorModal.tsx`.
92. `UI-092`: Extract database state management from `Database.tsx` into custom hook `useDatabaseViewer.ts`.
93. `UI-093`: Deconstruct monolithic `RepoSettings.tsx` (588 lines) into modular cards (`RepoRemotesCard.tsx`, `RepoCronCard.tsx`).
94. `UI-094`: Deconstruct monolithic `DesignSystem.tsx` (477 lines) into dedicated cluster showcase components.
95. `UI-095`: Deconstruct `DocsSidebar.tsx` (200 lines) into `DocsSidebarNav.tsx` and `DocsSidebarFooter.tsx`.
96. `UI-096`: Extract static docs navigation configuration into `src/config/docsNavConfig.ts`.
97. `UI-097`: Extract `DocsHeader.tsx` (< 80 lines) from `DocsLayout.tsx`.
98. `UI-098`: Deconstruct `CommandCard.tsx` (141 lines) into `CommandFlagsTable.tsx` and `CommandActionGroup.tsx`.
99. `UI-099`: Separate custom terminal emulation hooks from presentation viewports in `InteractiveTUI.tsx` and `Diff.tsx`.
100. `UI-100`: Ensure all newly constructed and refactored React components adhere strictly to $\le 80$–100 lines.
