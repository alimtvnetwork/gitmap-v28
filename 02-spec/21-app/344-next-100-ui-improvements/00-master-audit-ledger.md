# Master Audit Ledger: Next 100 UI Improvements (Round 5)
**Task Slug:** `344-next-100-ui-improvements`  
**Standard:** WCAG 2.2 AA/AAA, 4-Plane Depth Hierarchy, 60/30/10 Rule, Single-Accent Von Restorff Principle, Pure DOM Live Typography, 44px Minimum Touch Bounds, 72px Glassmorphic Header Physics, GPU Composite Acceleration (`transform-gpu`), React Micro-Component Decomposition (<= 100 lines).

---

## 1. Cluster 01: Color Tokens, HSL Architecture & Zero Hardcoded Hex (Items 01–10)
1. `src/pages/Watch.tsx:28,33`: Audit and verify semantic warning token usage across watch indicators.
2. `src/pages/Release.tsx:208,212`: Replace ad-hoc palette strings with semantic tokens.
3. `cli/cmdui/ui_assets_markup_misc.go:129`: Replace ad-hoc inline styles with theme CSS variables.
4. `src/components/docs/TerminalDemo.tsx:102`: Replace inline HSL literal `bg-[hsl(45,80%,50%)]/60` with standard `--warning` token.
5. `src/pages/Pipeline.tsx:82`: Standardize `yellow-500` status text to WCAG AA-compliant amber tokens.
6. `src/pages/DiffProfiles.tsx:82-88`: Balance pipeline step chips with muted neutral backgrounds, isolating the active stage.
7. `src/pages/Architecture.tsx:52-68`: Re-balance architecture pipeline step tags to 60/30 neutral tones.
8. `src/pages/InteractiveTUI.tsx:119`: Upgrade low-contrast badge labels to WCAG AA-compliant tokens.
9. `src/components/docs/troubleshooting/TroubleshootingIssueCard.tsx:66`: Elevate issue tags and metadata to minimum 4.5:1 contrast against surface background.
10. `src/pages/History.tsx:55`: Ensure branch color labels utilize tokenized ramps.

---

## 2. Cluster 02: Contrast Ratios & WCAG 2.2 AA/AAA Compliance (Items 11–20)
11. `src/components/docs/TerminalDemo.tsx:28`: Remove double suppression `text-[hsl(var(--foreground))]/70` in favor of high-contrast foreground text.
12. `src/pages/ReleaseVersion.tsx:62-65`: Elevate release header with high-contrast text tokens.
13. `src/pages/GenericCLI.tsx:18`: Upgrade muted parameters to minimum 4.5:1 contrast ratio.
14. `src/pages/NotFound.tsx:15`: Verify 404 page contrast and semantic tokens.
15. `src/pages/Setup.tsx:45`: Verify setup steps contrast against dark canvas.
16. `src/pages/Commands.tsx:80`: Elevate command flag descriptions to WCAG AA contrast.
17. `src/pages/Prune.tsx:50`: Ensure prune candidate warning text satisfies 4.5:1 ratio.
18. `src/pages/Import.tsx:40`: Ensure repo import status labels maintain clear readability.
19. `src/pages/Export.tsx:45`: Calibrate export format options with high-contrast active state.
20. `src/pages/SSH.tsx:60`: Align SSH key fingerprint contrast with typography standards.

---

## 3. Cluster 03: 4-Plane Depth Hierarchy & Materiality (Items 21–30)
21. `src/pages/Release.tsx:150`: Elevate release candidate cards with Plane 2 surface treatment.
22. `src/pages/History.tsx:75`: Add hairline borders and subtle shadow to git commit timeline nodes.
23. `src/pages/DiffProfiles.tsx:60`: Add Plane 1 recessed well base to profile comparison tables.
24. `src/pages/Architecture.tsx:40`: Add Plane 1 inset container for architecture diagrams.
25. `src/pages/Bookmarks.tsx:50`: Elevate saved bookmarks with Plane 2 card tokens.
26. `src/pages/Alias.tsx:60`: Add subtle hairline border to shell alias cards.
27. `src/pages/Commands.tsx:70`: Elevate command cards with Plane 2 surface styling.
28. `src/pages/Setup.tsx:60`: Add Plane 1 well base to terminal setup snippets.
29. `src/pages/Doctor.tsx:60`: Add Plane 1 recessed well base to doctor output.
30. `src/pages/Watch.tsx:60`: Elevate watch monitor overview with Plane 2 depth hierarchy.

---

## 4. Cluster 04: 60/30/10 Visual Balance & Single-Accent Von Restorff Rule (Items 31–40)
31. `src/pages/Settings.tsx:328`: Isolate "Save Changes" as exclusive primary CTA, neutralizing secondary header icons.
32. `src/pages/Release.tsx:208,212`: De-saturate resolution step chips to subtle neutral borders.
33. `src/pages/Pipeline.tsx:110`: Neutralize secondary pipeline trigger buttons to outline/ghost variants.
34. `src/pages/Changelog.tsx:85-95`: Ensure release tags maintain neutral ground, emphasizing release title.
35. `src/pages/Install.tsx:65`: Isolate primary download button against secondary instructions.
36. `src/pages/FixRepo.tsx:70`: Neutralize secondary repair actions, keeping primary fix button prominent.
37. `src/pages/Profile.tsx:55`: Balance user profile configuration actions.
38. `src/pages/RepoSettings.tsx:80`: Enforce single accent button on repository settings form.
39. `src/pages/SSH.tsx:90`: Standardize SSH connection test button to single-accent hierarchy.
40. `src/pages/Export.tsx:60`: Isolate export file generation button as sole primary CTA.

---

## 5. Cluster 05: Typography Stacks & Pure Live DOM Text Mandate (Items 41–50)
41. `src/pages/Settings.tsx:315`: Calibrate main page title from loose `tracking-tight` to exact optical kerning `tracking-[-0.02em]`.
42. `src/pages/Troubleshooting.tsx:176`: Add explicit optical tracking `tracking-[-0.02em]` to `docs-h1` title.
43. `src/pages/AGYPrompts.tsx:285`: Bind `font-heading` token and `tracking-[-0.02em]` to page title.
44. `src/pages/InteractiveTUI.tsx:119`: Apply optical tracking and typography ramp to TUI view titles.
45. `src/pages/Architecture.tsx:49`: Standardize diagram labels to pure DOM live text with JetBrains Mono inline.
46. `src/pages/History.tsx:40`: Enforce optical kerning on commit history timeline.
47. `src/pages/Commands.tsx:45`: Enforce pure live DOM text and optical kerning on command dictionary.
48. `src/pages/Doctor.tsx:50`: Standardize diagnostic page title and subtitles to `font-heading`.
49. `src/pages/Setup.tsx:35`: Enforce `font-heading tracking-[-0.02em]` on setup guides.
50. `src/pages/Bookmarks.tsx:30`: Apply optical tracking to bookmark section headers.

---

## 6. Cluster 06: Button System Standards & Hit Target Geometry (Items 51–60)
51. `src/components/docs/TerminalDemo.tsx:108-110`: Upgrade demo play button to 44px hit bounds and tactile `active:scale-[0.98]`.
52. `src/components/docs/TerminalDemo.tsx:111-113`: Upgrade reset button with 44px touch protection and active depression scaling.
53. `cli/cmdui/ui_assets_theme.go:120`: Add desktop tactile depression `active { transform: translateY(1px); filter: brightness(0.94); }` and mobile 44px bounds.
54. `src/pages/GenericCLI.tsx:32`: Expand CLI action buttons to minimum 44px bounds on mobile viewports.
55. `src/pages/Release.tsx:160`: Ensure release action triggers maintain 44px touch protection.
56. `src/pages/Pipeline.tsx:120`: Expand pipeline stage restart buttons to 44px bounds.
57. `src/pages/History.tsx:90`: Expand commit cherry-pick/revert buttons to 44px touch targets.
58. `src/pages/SSH.tsx:75`: Upgrade SSH key copy triggers to 44px touch geometry.
59. `src/pages/Bookmarks.tsx:65`: Ensure bookmark remove buttons maintain safe 44px hit bounds.
60. `src/pages/Export.tsx:80`: Expand export copy/download buttons to 44px hit bounds.

---

## 7. Cluster 07: Header, Navigation & Safe-Region Physics (Items 61–70)
61. `cli/cmdui/ui_assets_theme.go:97`: Align embedded web UI `header` to 72px min-height with glassmorphic `backdrop-filter: blur(12px)`.
62. `cli/cmdui/ui_assets_theme.go:147`: Add glassmorphic `backdrop-filter: blur(12px)` and `env(safe-area-inset-bottom)` to sticky `.settings-savebar`.
63. `cli/cmdui/ui_assets_markup_misc.go:29-35`: Add safe-area inset padding (`env(safe-area-inset-top)`) to embedded CLI header.
64. `src/pages/Settings.tsx:337-375`: Add sticky glassmorphic styling to settings navigation bar.
65. `src/pages/ReleaseVersion.tsx:62-65`: Elevate release header with sticky navigation indicators.
66. `src/components/docs/troubleshooting/TroubleshootingIssueCard.tsx:19-23`: Ensure safe hit-padding and header alignment.
67. `src/components/ui/sidebar.tsx:174-184`: Apply 220ms safe-region debouncing to sidebar sheet close interactions.
68. `src/components/ui/dialog.tsx:39-44`: Add safe-area margin boundaries (`max-h-[calc(100vh-2rem)] overflow-y-auto`) to modal dialogs.
69. `src/pages/Pipeline.tsx:50`: Ensure pipeline progress header pins smoothly with backdrop-blur.
70. `src/pages/Watch.tsx:40`: Align watch refresh header with safe navigation padding.

---

## 8. Cluster 08: Spacing Scales, Containers & Band Rhythm (Items 71–80)
71. `src/pages/CloneFixRepo.tsx:6`: Harmonize container with standard responsive gutter stepping (`px-4 sm:px-6 lg:px-8`).
72. `src/pages/CommitIn.tsx:12`: Realign container child band spacing to uniform 8pt grid scale (`space-y-8`).
73. `src/pages/GenericCLI.tsx:18`: Standardize sub-panel arbitrary padding `px-4 pb-5 pt-3` to 8pt rhythm (`p-4 sm:p-6`).
74. `src/pages/MakePublic.tsx:6`: Align container width and gutter stepping with documentation layouts.
75. `src/pages/MergeLeft.tsx:6`: Harmonize layout container with responsive 8pt rhythm.
76. `src/pages/Replace.tsx:6`: Align content wrapper to eliminate page-to-page band jitter.
77. `src/pages/DiffProfiles.tsx:60`: Harmonize profile panel padding to 8pt scale.
78. `src/pages/History.tsx:60`: Align commit log timeline spacing to 8pt scale.
79. `src/pages/Commands.tsx:55`: Standardize command list gutters to 8pt scale.
80. `src/pages/Setup.tsx:50`: Harmonize setup step gaps to 8pt scale.

---

## 9. Cluster 09: CSS3 Hardware-Accelerated Motion & Accessibility (Items 81–90)
81. `src/components/ui/sidebar.tsx:174-184`: Add `transform-gpu` to off-canvas sidebar transitions.
82. `cli/cmdui/ui_assets_theme.go:91-92`: Add `will-change: transform` and active depression transform to `.nav-btn`.
83. `cli/cmdui/ui_assets_theme.go:108-109`: Isolate card hover transforms with GPU composite layers (`will-change: transform`).
84. `cli/cmdui/ui_assets_theme.go:156`: Incorporate `env(safe-area-inset-bottom)` and `env(safe-area-inset-right)` into toast notifications.
85. `src/components/docs/TerminalDemo.tsx:102`: Upgrade cursor pulse animation to pure opacity transitions with `will-change: opacity`.
86. `src/components/projects/ProjectDetailDialog.tsx:57`: Add `transform-gpu` to dialog modal entrance transition.
87. `src/pages/InteractiveTUI.tsx:80`: Apply hardware-accelerated transitions to interactive TUI view switches.
88. `src/pages/Pipeline.tsx:90`: Add hardware-accelerated spinner transitions to running tasks.
89. `src/pages/Stats.tsx:70`: Ensure chart hover tooltips utilize GPU composite layers.
90. `src/pages/Watch.tsx:80`: Apply GPU-accelerated pulses to live sync indicators.

---

## 10. Cluster 10: Frontend Code Hygiene & React Component Deconstruction (Items 91–100)
91. `src/pages/Settings.tsx`: Deconstruct settings sub-panels into modular cards (<= 90 lines each).
92. `src/pages/Pipeline.tsx`: Decompose pipeline run rows into `PipelineRunRow.tsx`.
93. `src/pages/History.tsx`: Extract commit timeline entry into `HistoryCommitRow.tsx`.
94. `src/pages/Release.tsx`: Extract release step card into `ReleaseStepCard.tsx`.
95. `src/pages/DiffProfiles.tsx`: Extract profile selector into `DiffProfileSelector.tsx`.
96. `src/pages/Commands.tsx`: Extract command dictionary item into `CommandItemCard.tsx`.
97. `src/pages/SSH.tsx`: Extract SSH node credentials into `SSHNodeCard.tsx`.
98. `src/pages/Doctor.tsx`: Maintain concise modular sections <= 100 lines.
99. `src/pages/Dashboard.tsx`: Maintain concise modular layout <= 100 lines.
100. `src/pages/Stats.tsx`: Maintain concise modular structure <= 120 lines.
