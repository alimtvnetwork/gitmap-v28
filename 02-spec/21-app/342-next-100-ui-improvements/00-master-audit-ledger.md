# Master Audit Ledger: Next 100 UI Improvements (Round 3)
**Task Slug:** `342-next-100-ui-improvements`  
**Standard:** WCAG 2.2 AA/AAA, 4-Plane Depth Hierarchy, 60/30/10 Rule, Single-Accent Von Restorff Principle, Pure DOM Live Typography, 44px Minimum Touch Bounds, 72px Glassmorphic Header Physics, GPU Composite Acceleration (`transform-gpu`), React Micro-Component Decomposition (<= 100 lines).

---

## 1. Cluster 01: Color Tokens, HSL Architecture & Zero Hardcoded Hex (Items 01–10)
1. `src/pages/Database.tsx:550,558,566`: Replace ad-hoc overlay `bg-black/40` and raw neon hexes (`text-sky-400`, `text-emerald-400`, `text-amber-300`) with semantic surface and code tokens.
2. `src/pages/Watch.tsx:28,33`: Replace raw `text-yellow-400` status indicator with semantic warning tokens (`text-warning` / `text-amber-600 dark:text-amber-400`).
3. `src/pages/Release.tsx:208,212`: Replace ad-hoc palette strings with semantic tokens.
4. `cli/cmdui/ui_assets_markup_misc.go:129`: Replace ad-hoc inline styles with theme CSS variables.
5. `src/components/docs/TerminalDemo.tsx:102`: Replace inline HSL literal `bg-[hsl(45,80%,50%)]/60` with standard `--warning` token.
6. `src/pages/Pipeline.tsx:82`: Standardize `yellow-500` status text to WCAG AA-compliant amber tokens (`text-amber-700 dark:text-amber-300`).
7. `src/pages/Settings.tsx:474-475`: Fix `text-zinc-400` on light background by upgrading to `text-muted-foreground` / `text-foreground/80`.
8. `src/pages/Doctor.tsx:111-117`: Standardize severity column badges to balanced neutral containers with semantic color dots.
9. `src/pages/DiffProfiles.tsx:82-88`: Balance pipeline step chips with muted neutral backgrounds, isolating the active stage.
10. `src/pages/Architecture.tsx:52-68`: Re-balance architecture pipeline step tags to 60/30 neutral tones.

---

## 2. Cluster 02: Contrast Ratios & WCAG 2.2 AA/AAA Compliance (Items 11–20)
11. `src/components/docs/TerminalDemo.tsx:28`: Remove double suppression `text-[hsl(var(--foreground))]/70` in favor of high-contrast foreground text.
12. `src/pages/Database.tsx:362-363`: Replace low-contrast `text-emerald-400` / `text-red-400` table cell badges with calibrated WCAG AA compliant colors (`text-emerald-700 dark:text-emerald-400` and `text-rose-700 dark:text-rose-400`).
13. `src/pages/Database.tsx:376`: Upgrade `text-purple-400` on `bg-purple-500/15` badge to `text-purple-700 dark:text-purple-300` (>= 4.5:1).
14. `src/pages/Database.tsx:389`: Upgrade `text-sky-400` on `bg-sky-500/15` badge to `text-sky-700 dark:text-sky-300` (>= 4.5:1).
15. `src/pages/Database.tsx:402`: Upgrade `text-amber-400` on `bg-amber-500/15` badge to `text-amber-800 dark:text-amber-300` (>= 4.5:1).
16. `src/pages/Database.tsx:352`: Eliminate degraded contrast `text-muted-foreground/60 italic` on `NULL` table values, raising contrast above 4.5:1.
17. `src/pages/AGYPrompts.tsx:359`: Eliminate double opacity suppression `<Search className="text-muted-foreground opacity-50" />` to single `text-muted-foreground`.
18. `src/components/spec/SpecSearchBar.tsx:29`: Eliminate double opacity suppression on `kbd` helper (`text-muted-foreground/50` -> `text-muted-foreground font-medium`).
19. `src/pages/InteractiveTUI.tsx:119`: Upgrade low-contrast badge labels to WCAG AA-compliant tokens.
20. `src/components/docs/troubleshooting/TroubleshootingIssueCard.tsx:66`: Elevate issue tags and metadata to minimum 4.5:1 contrast against surface background.

---

## 3. Cluster 03: 4-Plane Depth Hierarchy & Materiality (Items 21–30)
21. `src/pages/Database.tsx:311`: Upgrade database table container to Plane 1 recessed well base (`bg-plane-well` / `dark:bg-plane-well`).
22. `src/pages/Database.tsx:470`: Upgrade `CellInspectionModal` to Plane 3 elevation tokens (`dark:bg-plane-elevated`) with top hairline highlight `inset 0 1px 0 rgba(255,255,255,0.12)`.
23. `src/pages/Database.tsx:550,558,566`: Add semantic Plane 1 inset well shadows (`shadow-inner`) to hex dump, JSON, and text viewers.
24. `src/pages/Projects.tsx:122-142`: Elevate project filter metric summary cards with Plane 1 stat-well treatments and hairline borders.
25. `src/components/projects/ProjectDetailDialog.tsx:57`: Refine `DialogContent` container with Plane 3 modal elevation overrides and top hairline highlight.
26. `src/pages/Settings.tsx:397`: Elevate Antigravity Multi-Instance Panel to Plane 2 surface card with hairline lighting (`shadow-sm`, top refraction highlight).
27. `src/pages/Settings.tsx:581`: Refactor UI/UX & CSS3 Repository Index section card with Plane 2 card hierarchy tokens.
28. `src/pages/Dashboard.tsx:68`: Upgrade dashboard feature grid cards with Plane 2 card elevation and subtle hover glow.
29. `src/pages/Dashboard.tsx:90,112`: Add Plane 1 recessed well base to Flags and File Layout tables.
30. `src/components/docs/troubleshooting/TroubleshootingIssueCard.tsx:66`: Elevate issue cards with Plane 2 surface card styling and hairline highlights.

---

## 4. Cluster 04: 60/30/10 Visual Balance & Single-Accent Von Restorff Rule (Items 31–40)
31. `src/pages/Database.tsx:360-406`: Neutralize multi-chroma badge soup in table rows, re-weighting background surfaces to muted neutrals (60/30 balance).
32. `src/pages/Database.tsx:597-601`: Neutralize modal secondary buttons, reserving primary accent strictly for primary completion action.
33. `src/pages/Settings.tsx:328`: Isolate "Save Changes" as exclusive primary CTA, neutralizing secondary header icons.
34. `src/pages/AGYPrompts.tsx:339-342`: Calibrate category chip filters to neutral muted pills when active to avoid competing with search CTA.
35. `src/pages/Release.tsx:208,212`: De-saturate resolution step chips to subtle neutral borders.
36. `src/pages/DiffProfiles.tsx:82-88`: Balance pipeline step chips with muted neutral backgrounds, isolating the active stage.
37. `src/pages/Architecture.tsx:52-68`: Re-balance architecture pipeline step tags to 60/30 neutral tones.
38. `src/pages/Doctor.tsx:111-117`: Standardize severity column badges to balanced neutral containers with semantic color dots.
39. `src/pages/Projects.tsx:145-165`: Harmonize project card action buttons with single-accent hierarchy.
40. `src/pages/Changelog.tsx:85-95`: Ensure release tags maintain neutral ground, emphasizing release title hierarchy.

---

## 5. Cluster 05: Typography Stacks & Pure Live DOM Text Mandate (Items 41–50)
41. `src/pages/Database.tsx:208`: Upgrade heading to `font-heading font-bold tracking-[-0.02em]` for optical kerning.
42. `src/pages/Projects.tsx:101`: Replace raw `font-mono` on page title with `font-heading font-bold tracking-[-0.02em]`.
43. `src/pages/Dashboard.tsx:59,65,76,84,89,111`: Add optical heading kerning `tracking-[-0.02em]` across all major section headers.
44. `src/pages/Settings.tsx:315`: Calibrate main page title from loose `tracking-tight` to exact optical kerning `tracking-[-0.02em]`.
45. `src/pages/Troubleshooting.tsx:176`: Add explicit optical tracking `tracking-[-0.02em]` to `docs-h1` title.
46. `src/pages/AGYPrompts.tsx:285`: Bind `font-heading` token and `tracking-[-0.02em]` to page title.
47. `src/pages/InteractiveTUI.tsx:119`: Apply optical tracking and typography ramp to TUI view titles.
48. `src/pages/Architecture.tsx:49`: Standardize diagram labels to pure DOM live text with JetBrains Mono inline.
49. `src/pages/BatchActions.tsx:178`: Standardize batch action headers to `font-heading tracking-[-0.02em]`.
50. `src/pages/CloneNext.tsx:109`: Enforce pure live DOM text and optical kerning across version cards.

---

## 6. Cluster 06: Button System Standards & Hit Target Geometry (Items 51–60)
51. `src/pages/Database.tsx:249-257`: Expand database quick filter buttons to 44px hit targets and add `active:scale-[0.98]`.
52. `src/pages/Database.tsx:504-510`: Upgrade table reload button from `h-7 w-7` to standard 44px hit bounds (`size-11` or padded container) with tactile depression.
53. `src/components/docs/TerminalDemo.tsx:108-110`: Upgrade demo play button to 44px hit bounds and tactile `active:scale-[0.98]`.
54. `src/components/docs/TerminalDemo.tsx:111-113`: Upgrade reset button with 44px touch protection and active depression scaling.
55. `src/pages/Projects.tsx:145-165`: Ensure project quick actions have 44px touch protection on mobile viewport.
56. `src/pages/Settings.tsx:328`: Standardize "Save Changes" button hit target and active state.
57. `cli/cmdui/ui_assets_theme.go:120`: Add desktop tactile depression `active { transform: translateY(1px); filter: brightness(0.94); }` and mobile 44px bounds.
58. `src/pages/GenericCLI.tsx:32`: Expand CLI action buttons to minimum 44px bounds on mobile viewports.
59. `src/pages/BatchActions.tsx:120`: Standardize batch execution triggers to 44px touch bounding.
60. `src/components/spec/SpecSearchBar.tsx:24`: Ensure search clear button maintains 44px hit bounds.

---

## 7. Cluster 07: Header, Navigation & Safe-Region Physics (Items 61–70)
61. `cli/cmdui/ui_assets_theme.go:97`: Align embedded web UI `header` to 72px min-height with glassmorphic `backdrop-filter: blur(12px)`.
62. `cli/cmdui/ui_assets_theme.go:147`: Add glassmorphic `backdrop-filter: blur(12px)` and `env(safe-area-inset-bottom)` to sticky `.settings-savebar`.
63. `cli/cmdui/ui_assets_markup_misc.go:29-35`: Add safe-area inset padding (`env(safe-area-inset-top)`) to embedded CLI header.
64. `src/pages/Database.tsx:202-219`: Align `DatabaseHeader` with glassmorphic backdrop-blur when sticky.
65. `src/pages/Database.tsx:529-544`: Add backdrop-blur and safe pin offset to database tab strip.
66. `src/pages/Settings.tsx:337-375`: Add sticky glassmorphic styling to settings navigation bar.
67. `src/pages/ReleaseVersion.tsx:62-65`: Elevate release header with sticky navigation indicators.
68. `src/components/docs/troubleshooting/TroubleshootingIssueCard.tsx:19-23`: Ensure safe hit-padding and header alignment.
69. `src/components/ui/sidebar.tsx:174-184`: Apply 220ms safe-region debouncing to sidebar sheet close interactions.
70. `src/components/ui/dialog.tsx:39-44`: Add safe-area margin boundaries (`max-h-[calc(100vh-2rem)] overflow-y-auto`) to modal dialogs.

---

## 8. Cluster 08: Spacing Scales, Containers & Band Rhythm (Items 71–80)
71. `src/pages/Database.tsx:311`: Enforce 8pt cell padding rhythm in table grid headers and cells.
72. `src/pages/Projects.tsx:122`: Enforce uniform gap rhythm (`gap-4 sm:gap-6`) in metric summary cards.
73. `src/pages/CloneFixRepo.tsx:6`: Harmonize container with standard responsive gutter stepping (`px-4 sm:px-6 lg:px-8`).
74. `src/pages/CommitIn.tsx:12`: Realign container child band spacing to uniform 8pt grid scale (`space-y-8`).
75. `src/pages/GenericCLI.tsx:18`: Standardize sub-panel arbitrary padding `px-4 pb-5 pt-3` to 8pt rhythm (`p-4 sm:p-6`).
76. `src/pages/MakePublic.tsx:6`: Align container width and gutter stepping with documentation layouts.
77. `src/pages/MergeLeft.tsx:6`: Harmonize layout container with responsive 8pt rhythm.
78. `src/pages/Replace.tsx:6`: Align content wrapper to eliminate page-to-page band jitter.
79. `src/pages/BatchActions.tsx:40`: Enforce 8pt spacing rhythm between action cards.
80. `src/pages/DiffProfiles.tsx:60`: Harmonize profile panel padding to 8pt scale.

---

## 9. Cluster 09: CSS3 Hardware-Accelerated Motion & Accessibility (Items 81–90)
81. `src/components/ui/sidebar.tsx:174-184`: Add `transform-gpu` to off-canvas sidebar transitions.
82. `cli/cmdui/ui_assets_theme.go:91-92`: Add `will-change: transform` and active depression transform to `.nav-btn`.
83. `cli/cmdui/ui_assets_theme.go:108-109`: Isolate card hover transforms with GPU composite layers (`will-change: transform`).
84. `cli/cmdui/ui_assets_theme.go:156`: Incorporate `env(safe-area-inset-bottom)` and `env(safe-area-inset-right)` into toast notifications.
85. `src/components/docs/TerminalDemo.tsx:102`: Upgrade cursor pulse animation to pure opacity transitions with `will-change: opacity`.
86. `src/components/projects/ProjectDetailDialog.tsx:57`: Add `transform-gpu` to dialog modal entrance transition.
87. `src/pages/Projects.tsx:122-142`: Add hardware-accelerated transitions to metric card hover elevations.
88. `src/pages/Database.tsx:470`: Ensure modal animations leverage GPU composition (`transform-gpu`).
89. `src/pages/InteractiveTUI.tsx:80`: Apply hardware-accelerated transitions to interactive TUI view switches.
90. `src/pages/Dashboard.tsx:68`: Add `transform-gpu` and `transition-transform` to feature card hover physics.

---

## 10. Cluster 10: Frontend Code Hygiene & React Component Deconstruction (Items 91–100)
91. `src/pages/Database.tsx`: Deconstruct database visualizer, extracting `DatabaseHeader.tsx` (<= 80 lines).
92. `src/pages/Database.tsx`: Extract `DatabaseCellRenderer.tsx` (<= 90 lines).
93. `src/pages/Database.tsx`: Extract `DatabaseCellInspector.tsx` (<= 95 lines).
94. `src/pages/Database.tsx`: Extract `DatabaseQuickFilters.tsx` (<= 75 lines).
95. `src/pages/Projects.tsx`: Extract metric cards into `ProjectMetricCards.tsx` (<= 75 lines).
96. `src/pages/Settings.tsx`: Extract multi-instance settings card into `AntigravityInstancesCard.tsx` (<= 90 lines).
97. `src/components/projects/TypeBadge.tsx`: Verify clean modular props contract under 60 lines.
98. `src/components/docs/troubleshooting/TroubleshootingFilterButtons.tsx`: Verify single-responsibility structure under 80 lines.
99. `src/components/docs/CopyPaletteButton.tsx`: Ensure modular composition under 80 lines.
100. `src/components/terminal/TerminalDrawer.tsx`: Maintain lean deconstructed sub-components under 100 lines each.
