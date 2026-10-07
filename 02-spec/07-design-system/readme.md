# AI-Adaptable Design System

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 07 Design System.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable AI Agent Reading & Learning Checklist

Before authoring components or generating code, AI agents MUST read and master the specifications in this directory in the following sequence:

- [ ] `/learn` **Phase 1: Visual Fundamentals**
  - Read [`02-design-principles.md`](./02-design-principles.md) — 60/30/10 visual balance, 4-plane depth hierarchy, single-accent Von Restorff.
  - Read [`03-theme-variable-architecture.md`](./03-theme-variable-architecture.md) — CSS custom properties single source of truth.
  - Read [`04-typography.md`](./04-typography.md) — Fluid typography scale (`clamp()`), line height ratios, tabular numerals.
- [ ] `/learn` **Phase 2: Core Components & Layout**
  - Read [`05-spacing-layout.md`](./05-spacing-layout.md) — Spacing scale, container max-widths, responsive breakpoints.
  - Read [`06-borders-shapes.md`](./06-borders-shapes.md) — Border radii tokens, hairlines over drop shadows.
  - Read [`08-motion-transitions.md`](./08-motion-transitions.md) & [`11-button-system.md`](./11-button-system.md) — Deceleration curves, button variants.
- [ ] `/learn` **Phase 3: Multi-Theme Catalog & Machine-Readable Tokens**
  - Read [`16-theme-catalogue-and-palettes.md`](./16-theme-catalogue-and-palettes.md) — 8 distinct theme families.
  - Inspect [`tokens/design-tokens.json`](./tokens/design-tokens.json) & [`17-theme-tokens.json`](./17-theme-tokens.json) — Programmatic token data.
  - Master [`tokens/theme-palette.less`](./tokens/theme-palette.less) — Parametric LESS mixins (**LESS is explicitly preferred over CSS**).
- [ ] `/learn` **Phase 4: Materiality, Interactions & AI Training**
  - Read [`18-dark-mode-and-materiality.md`](./18-dark-mode-and-materiality.md) — 4-plane depth steps, grain overlays, progressive blur.
  - Read [`19-modern-motion-and-sliding-interactions.md`](./19-modern-motion-and-sliding-interactions.md) — Multi-card sliding carousels, reset loops.
  - Read [`20-ai-training-and-checklist-guide.md`](./20-ai-training-and-checklist-guide.md) — Anti-AI-slop rubric, 5-question above-the-fold contract.
- [ ] `/learn` **Phase 5: Modern CSS3 & Building Blocks**
  - Read [`21-css3-animations-and-interactions.md`](./21-css3-animations-and-interactions.md) — Marquees, fluid accordions, neon borders.
  - Read [`22-native-css-select-and-border-shapes.md`](./22-native-css-select-and-border-shapes.md) — `base-select`, `::picker(select)`, `border-shape`.
  - Read [`23-building-block-components.md`](./23-building-block-components.md) — Curriculum cards, 4-cell subgrids, standalone SVGs.
  - Read [`38-card-and-pricing-components.md`](./38-card-and-pricing-components.md) — 3-tier pricing tables and floating editorial containers.
  - Read [`41-homepage-and-blog-sections.md`](./41-homepage-and-blog-sections.md) — The 15 flagship homepage sections and complete blog editorial templates.
- [ ] `/learn` **Phase 6: Slide Presentation Engine & Canvas Architecture**
  - Read [`24-slide-presentation-system.md`](./24-slide-presentation-system.md) — Fixed 16:9 canvas scaling (`1920×1080`), dual-screen presenter console (`/present`), BroadcastChannel sync, webcam PIP, step reveals, overview grid (`G`), top jumper (`J`).
  - Read [`27-slide-canvas-and-themes.md`](./27-slide-canvas-and-themes.md) — Virtual canvas scale math, 10 production themes (`bright-gold` default, `noir-gold`, `vscode-dark`, `dracula`, `monokai`, `github-light`, `paper-ink`, `macos-sonoma`, `windows-11`, `navy-blue`), and the Light Theme Contract.
  - Read [`40-theme-switch.md`](./40-theme-switch.md) — 8 slide themes, shared color variables, and the runtime color switch.
  - Read [`28-slide-layouts.md`](./28-slide-layouts.md) — Base slot models and closed layout union.
  - Read [`29-slide-navigation-and-builder.md`](./29-slide-navigation-and-builder.md) — HUD, transitions, keys, slide builder.
  - Read [`30-slide-palette-type-and-shell.md`](./30-slide-palette-type-and-shell.md) — Dark amber palette, type scale, shell layers.
  - Read [`31-slide-controller-buttons.md`](./31-slide-controller-buttons.md) — Floating HUD controller pill, action buttons, webcam PIP overlay, timer, dots.
  - Read [`42-slide-quiz-preview-chrome-and-default-shadows.md`](./42-slide-quiz-preview-chrome-and-default-shadows.md) — Default text/box shadow pair, presentation option cards, botanical-light quiz theme, keyboard shortcut matrix, center-stage layout.
  - Read [`32-slide-color-options.md`](./32-slide-color-options.md) — 10-step gradient precision system ($S_0$–$S_9$) across 7 flagship themes, per-slide gradient editor, pill presets, 9-cell align.
  - Read [`42-slide-step-and-sound-system.md`](./42-slide-step-and-sound-system.md) — `StepTimelineSlide` vs `AdvanceStepSlide`, active focus rows, and Web Audio API synthesis engine.
  - Read [`43-slide-webcam-overlay.md`](./43-slide-webcam-overlay.md) — Presenter webcam PIP overlay, squircle gold rim, S/M/L/XL presets, auto-frame, and hotkeys.
  - Read [`44-slide-presenter-inspector-and-handouts.md`](./44-slide-presenter-inspector-and-handouts.md) — Presenter inspector route (`/slides/inspector`), 3-up printable handouts, and print mode.
- [ ] `/learn` **Phase 6b: Master Slide Layout Catalog & Pure DOM Typography**
  - Read [`34-slide-layout-catalog.md`](./34-slide-layout-catalog.md) — Non-Image Text Mandate and 20 master enterprise slide layouts.
- [ ] `/learn` **Phase 6c: Live Slide Builder Mode & Canvas Inspector**
  - Read [`35-slide-builder-canvas-inspector.md`](./35-slide-builder-canvas-inspector.md) — Decoupled dual-store architecture, 7 visual layers, draggable `BuilderPanel`, layout switcher, and headless PDF.
- [ ] `/learn` **Phase 6d: Precision Navigation & Mega Menu System**
  - Read [`10-header-navigation.md`](./10-header-navigation.md) & [`33-mega-menu-components.md`](./33-mega-menu-components.md) — Sticky `72px` glass header, `SlideSwapLabel` CSS keyframes, growing left hairline, 3D flip promo card, and safe-region hover mechanics.
- [ ] `/learn` **Phase 6e: Website Content Builder Mode & In-Page Visual Editor**
  - Read [`36-website-content-builder-mode.md`](./36-website-content-builder-mode.md) & [`26-visual-builder.md`](./26-visual-builder.md) — In-page review-only editor: `contentEditable`, media replacement modal, menu editor, and ZIP export.
- [ ] `/learn` **Phase 6f0: Bright Gold Tech, Logo & Theme Switch**
  - Read [`05-bright-gold-tech/readme.md`](./05-bright-gold-tech/readme.md), [`39-logo-construction.md`](./39-logo-construction.md), and [`40-theme-switch.md`](./40-theme-switch.md).
- [ ] `/learn` **Phase 6f: Standalone Image & Banner Specifications**
  - Read [`37-image-specifications.md`](./37-image-specifications.md) — Canvas geometries, safe zones, text rules, and image palette.
- [ ] `/learn` **Phase 7: AI-Adaptable Modern SaaS Design System**
  - Read [`02-ai-system-design/readme.md`](./02-ai-system-design/readme.md) — Variable-driven theme swapping (Pink to Green), sticky nav, search module, "Join Us" CSS3 text-slide animation, "Climate AI" highlight, and "Team Greenhouse" section.
- [ ] `/learn` **Phase 8: Sweet Digs Design System & Interactive Theme Tester**
  - Read [`03-sweet-digs-design-system/readme.md`](./03-sweet-digs-design-system/readme.md) — Master index and token flow for the Sweet Digs system.
  - Explore interactive demo: [`theme-tester/index.html`](../../theme-tester/index.html) — Live interactive theme tester with theme switcher, split hero, search card, filter chips, expanding underline menu, and interactive buttons.
- [ ] `/learn` **Phase 9: White Blue Theme Design System (3-Format Color Standard)**
  - Read [`04-white-blue-theme/readme.md`](./04-white-blue-theme/readme.md) — Master index, 4-tier token flow, alternating band rhythm (`light` -> `soft` -> `dark` -> `void`), and blind-AI checklist.
  - Read [`04-white-blue-theme/01-colors-typography-and-tokens.md`](./04-white-blue-theme/01-colors-typography-and-tokens.md) — Complete 3-format color tables (`HEX`, `RGB`/`RGBA`, `HSL` + `OKLCH`), `Ubuntu` + `Poppins` + `JetBrains Mono` typography, gradients, and brand-tinted shadows.
  - Read [`04-white-blue-theme/02-header-mega-menu-and-footer.md`](./04-white-blue-theme/02-header-mega-menu-and-footer.md) — Sticky `72px` glass header, `SlideSwapLabel` nav links, safe-region hover Mega-Menu with 3D Flip Promo Card, and soft-band footer.
  - Read [`04-white-blue-theme/03-buttons-motion-and-interactions.md`](./04-white-blue-theme/03-buttons-motion-and-interactions.md) & [`04-white-blue-theme/04-cards-heroes-and-section-library.md`](./04-white-blue-theme/04-cards-heroes-and-section-library.md) — `WhiteBlueButton`, `.shine-sweep`, `.pointer-fill`, `card-premium`, Split Hero `CapabilityStack`, Sticky-Note Workflow Board, and Fluted Glass `ScrollStack`.

---

## Offered Design Systems & Multi-Theme Catalog

The slide switcher has 8 themes. The procedure, the count, and the CSS are [`40-theme-switch.md`](./40-theme-switch.md). The 9 rows below are reference families. They are not the switcher. Do not add them to `data-theme` unless file 40 is edited with all six colors.

| Short-Form ID | Theme Name | Base Ground | Accent | Mood & Best For |
|:---|:---|:---|:---|:---|
| `WHITE-BLUE` | White Blue Enterprise Editorial | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#2563EB` (`rgb(37, 99, 235)`, `hsl(221, 83%, 53%)`) / `#822EE8` | Light-first B2B SaaS, cloud ERP/CRM platforms, enterprise consulting |
| `LIGHT-TRUST` | Light High-Trust Editorial | `#ffffff` | `#6366f1` / `#f43f5e` | Clean SaaS, education, agency, certification portals |
| `DARK-NAVY` | Dark Obsidian Navy | `#0b1329` | `#8b5cf6` (Electric Violet) | AI engineering, developer tools, technical infrastructure |
| `CORP-GOLD` | Corporate Gold | `#0f172a` | `#facc15` (Gold) / `#10b981` | High-authority corporate events, keynote presentations |
| `CYBER-INDIGO`| Cyber Indigo Slate | `#030712` | `#6366f1` (Neon Indigo) | Dark-first developer dashboards, cyber security |
| `VSCODE-DARK` | VS Code Modern Dark | `#1e1e1e` | `#007acc` (Editor Blue) | Code-heavy developer tools, IDE extensions |
| `TOKYO-NIGHT` | Tokyo Night Storm | `#1a1b26` | `#bb9af7` (Magenta) | Modern developer portals, aesthetic CLI docs |
| `ONEDARK-PRO` | One Dark Pro | `#282c34` | `#61afef` (Chalky Blue) | Syntax themes, terminal output documentation |
| `WARM-PAPER`  | Warm Editorial Craft | `#fbf9f6` | `#ffad01` (Brand Amber) | Long-form reading, book publishers, research specs |

---

## Offered CSS3 Animations & Transitions

All animations are hardware-accelerated (`transform` and `opacity`) and include mandatory `prefers-reduced-motion` fallbacks:

| Animation Name | Keyframe / Method | Timing | Visual Effect & Purpose |
|:---|:---|:---|:---|
| **Hover Card Lift** | `translateY(-4px)` + shadow | 300ms `@ease-emphasized` | Tactile lift with non-white-blended darkish shade |
| **Interactive Row Highlight** | `padding-left` + left pill | 250ms `@ease-emphasized` | Vertical accent pill expands; row receives subtle darkish tint |
| **Infinite Marquee Ticker** | `translateX(-50%)` | 35s linear infinite | Continuous horizontal drift for tech badges (pauses on hover) |
| **Zero-JS Fluid Accordion** | `grid-template-rows: 0fr -> 1fr` | 300ms `@ease-emphasized` | Smooth height expansion without pixel calculation hacks |
| **Rotating Neon Border** | `conic-gradient` via `@property` | 6s linear infinite | Continuous electric border sweep around hero containers |
| **Ambient Status Glow** | `box-shadow` pulse | 2.5s infinite | Breathing green/violet halo around live system indicators |
| **Masked Text Reveal** | `clip-path: inset()` | 700ms `@ease-emphasized` | Editorial headline reveal sliding upward from mask |

---

## Styling Technology Preference: LESS Preferred Over CSS

> [!IMPORTANT]
> **LESS is explicitly preferred over CSS across this repository.**
> LESS provides modular parametric mixins, mathematical color tinting, nested scoping, and zero runtime overhead when compiled. Components must provide LESS mixins alongside raw CSS equivalents. Full implementations, anatomy breakdowns, and parametric LESS mixins are documented in:
- [`23-building-block-components.md`](./23-building-block-components.md) — Curriculum cards, 4-cell subgrids, and interactive list rows.
- [`38-card-and-pricing-components.md`](./38-card-and-pricing-components.md) — 3-Tier modular pricing tables and floating editorial cards.

---

> [!IMPORTANT]
> **CRITICAL AI INSTRUCTION:** This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** See [`98-confidence-report.md`](./98-confidence-report.md). Bright gold tech is the deck an agent can build today.
**Ambiguity:** The website builder is file 26 only. Slide types `usp-strike` and `bullets` have no component. A black slide catalog is not in this folder.

---

## Overview

This is the **canonical design system specification** for the project. It defines all visual behavior, interaction patterns, multi-theme architecture, color tokens, motion rules, and component construction guidance in a single, portable reference. Any AI agent or human contributor reading this specification should be able to:

1. **Reproduce** the current visual language on a new website
2. **Extend** the system with new pages and components that remain visually consistent
3. **Re-theme** the entire design by selecting from the standardized multi-theme catalog (Navy & Purple, VS Code themes, Heatmaps, or Warm Editorial)
4. **Migrate** the system to WordPress or any other CMS without rewriting component logic

The design system follows a **variable-driven architecture**: all colors, spacing, borders, and visual tokens are defined as CSS custom properties (HSL and OKLCH formats) in a single root file. Components never use hardcoded color values — they reference semantic tokens. Changing a token propagates to every component that uses it.

All animations and transitions prioritize **GPU-composited CSS3 transforms and opacity** — no heavy runtime libraries gating initial content paints.

---

## Design Philosophy

| Principle | Description |
|-----------|-------------|
| **Variable-First** | Every color, spacing, and visual property comes from a CSS custom property. No hardcoded values in components. |
| **Semantic Tokens** | Colors are named by purpose (`--primary`, `--accent`, `--muted`), not by value (`--purple`, `--pink`). |
| **HSL & OKLCH Color Models** | All colors use HSL space-separated format with OKLCH depth calculations for seamless lightness/saturation adjustments. |
| **4-Plane Depth Hierarchy** | Clear optical elevation: Base (Plane 0), Raised (Plane 1), Surface (Plane 2), and Elevated (Plane 3). |
| **60 / 30 / 10 Balance** | 60% dominant neutral, 30% structural surfaces/text, 10% purposeful accent. |
| **Single-Accent Von Restorff** | Exactly one prominent accent action per viewport to maintain unmistakable conversion focus. |
| **CSS3 Motion & Carousels** | High-performance transforms, controlled multi-card sliding loops, and instant reduced-motion fallbacks. |
| **Dark/Light Parity** | Every token has both light and dark values. Components never branch on theme — tokens handle it. |
| **Anti-AI-Slop Governance** | Concrete rejection rubrics that ban generic 3-card grids, unanchored heroes, and purple gradient soup. |
| **Portability** | Platform-agnostic. Works natively with React, Tailwind v4, static HTML, WordPress, or modern SSR frameworks. |

---

## Keywords

`design-system` · `multi-theme` · `navy-purple` · `vscode-themes` · `heatmaps` · `css-variables` · `4-plane-depth` · `sliding-carousels` · `anti-ai-slop` · `60-30-10-rule` · `von-restorff` · `typography` · `motion-system` · `ai-training`

---

## File Inventory

| # | File | Category | Description |
|---|------|----------|-------------|
| 00 | [readme.md](./readme.md) | Overview | This file — design system entry point and index |
| 01 | [02-design-principles.md](./02-design-principles.md) | Principles | Visual philosophy, consistency rules, interaction feel |
| 02 | [03-theme-variable-architecture.md](./03-theme-variable-architecture.md) | Theme | Complete CSS custom property registry — the single source of truth |
| 03 | [04-typography.md](./04-typography.md) | Typography | Font stacks, size hierarchy, weight rules, text spacing |
| 04 | [05-spacing-layout.md](./05-spacing-layout.md) | Layout | Spacing scale, container rules, grid/flex patterns, responsive breakpoints |
| 05 | [06-borders-shapes.md](./06-borders-shapes.md) | Borders | Border thickness, radius, color behavior, state changes |
| 06 | [08-motion-transitions.md](./08-motion-transitions.md) | Motion | CSS3 transition durations, easing, keyframe animations, state transforms |
| 07 | [09-code-blocks.md](./09-code-blocks.md) | Components | Code block rendering, language badges, line interaction, fullscreen |
| 08 | [10-header-navigation.md](./10-header-navigation.md) | Components | Header layout, menu structure, hover underlines, icon transitions |
| 09 | [11-button-system.md](./11-button-system.md) | Components | Button variants, slide text animation, highlight styles |
| 10 | [12-sidebar-system.md](./12-sidebar-system.md) | Components | Sidebar tree, active states, expand/collapse, search |
| 11 | [13-section-patterns.md](./13-section-patterns.md) | Patterns | Reusable section templates (hero, feature, team, CTA) |
| 12 | [14-page-creation-rules.md](./14-page-creation-rules.md) | Guide | Rules for building new pages from the design language |
| 13 | [15-wordpress-migration.md](./15-wordpress-migration.md) | Migration | CMS compatibility notes, block theme mapping, admin theming |
| 14 | [16-theme-catalogue-and-palettes.md](./16-theme-catalogue-and-palettes.md) | Multi-Theme | Navy & Purple, VS Code ecosystem, Heatmaps, and Warm Editorial palettes |
| 15 | [17-theme-tokens.json](./17-theme-tokens.json) | Token Data | Machine-readable theme JSON for programmatic consumption |
| 16 | [18-dark-mode-and-materiality.md](./18-dark-mode-and-materiality.md) | Materiality | 4-plane depth hierarchy, hairlines over shadows, progressive blur, grain, 60/30/10 |
| 17 | [19-modern-motion-and-sliding-interactions.md](./19-modern-motion-and-sliding-interactions.md) | Motion | Controlled sliding carousels, entrance grammar, section rhythm, reduced motion |
| 18 | [20-ai-training-and-checklist-guide.md](./20-ai-training-and-checklist-guide.md) | AI Training | Anti-slop rubric, 5-question above-the-fold contract, step-by-step design checklist |
| 19 | [21-css3-animations-and-interactions.md](./21-css3-animations-and-interactions.md) | Motion & Hover | CSS3 keyframes, cubic-bezier easing, darkish hover shades, line hover highlights |
| 20 | [22-native-css-select-and-border-shapes.md](./22-native-css-select-and-border-shapes.md) | Native Controls | Modern base-select, ::picker(select), and organic border-shape geometry |
| 21 | [23-building-block-components.md](./23-building-block-components.md) | Components | Curriculum card anatomy, 4-cell subgrids, standalone SVGs, LESS mixins |
| 22 | [24-slide-presentation-system.md](./24-slide-presentation-system.md) | Presentation | Webcam PIP, step reveals, dual-screen console |
| 23 | [25-page-assembly.md](./25-page-assembly.md) | Guide | Site and blog shells from the White Blue section library |
| 24 | [26-visual-builder.md](./26-visual-builder.md) | Builder | Review overlay for text, images, menu labels, in-group order |
| 25 | [27-slide-canvas-and-themes.md](./27-slide-canvas-and-themes.md) | Presentation | `1920×1080` canvas, amber tokens, JSON themes, noir-gold |
| 26 | [28-slide-layouts.md](./28-slide-layouts.md) | Presentation | Closed slide layout union |
| 27 | [29-slide-navigation-and-builder.md](./29-slide-navigation-and-builder.md) | Presentation | HUD, transition families, keys, slide builder |
| 30 | [30-slide-palette-type-and-shell.md](./30-slide-palette-type-and-shell.md) | Presentation | Dark amber palette, type scale, spotlight, shell layers |
| 31 | [31-slide-controller-buttons.md](./31-slide-controller-buttons.md) | Presentation | Floating HUD controller pill, 8-position mounting, audio synthesis, webcam PIP |
| 32 | [32-slide-color-options.md](./32-slide-color-options.md) | Presentation | 10-step gradient precision system ($S_0$–$S_9$), pill presets, relative luminance formula, 9-cell align |
| 33 | [33-mega-menu-components.md](./33-mega-menu-components.md) | Navigation | Precision mega panel dropdown, 3D flip card, growing left hairline, link staggers |
| 34 | [34-slide-layout-catalog.md](./34-slide-layout-catalog.md) | Presentation | Pure DOM text mandate and 20 master slide layout models |
| 35 | [35-slide-builder-canvas-inspector.md](./35-slide-builder-canvas-inspector.md) | Slide Builder | Dual-store architecture, 7 visual layers, selection overlays, hotkeys, acoustic audio cues, headless PDF |
| 36 | [36-website-content-builder-mode.md](./36-website-content-builder-mode.md) | Web Builder | In-page review-only builder overlay, `contentEditable` inline text editing, media modal, and ZIP diff exports |
| 37 | [37-image-specifications.md](./37-image-specifications.md) | Images | Exact canvas geometries, safe zones, text rules, and closed image palette |
| 38 | [38-card-and-pricing-components.md](./38-card-and-pricing-components.md) | Components | 3-Tier modular pricing tables, split pricing, What We Do, Ruled Paper workflow, Enterprise Hero |
| 39 | [39-logo-construction.md](./39-logo-construction.md) | Brand | Reusable SVG mark construction rules, archetypes, and theme variables |
| 40 | [40-theme-switch.md](./40-theme-switch.md) | Theme | 8 slide themes, shared color variables, and runtime color switcher |
| 41 | [41-homepage-and-blog-sections.md](./41-homepage-and-blog-sections.md) | Homepage & Blog | 15 modular homepage inspiration sections, blog index, reading progress, and editorial layout |
| 42 | [42-slide-step-and-sound-system.md](./42-slide-step-and-sound-system.md) | Slide Steps & Sound | `StepTimelineSlide` vs `AdvanceStepSlide`, focus row, Web Audio API synthesis engine |
| 43 | [43-slide-webcam-overlay.md](./43-slide-webcam-overlay.md) | Webcam PIP | Presenter webcam overlay, squircle gold rim, S/M/L/XL presets, auto-frame face tracker |
| 44 | [44-slide-presenter-inspector-and-handouts.md](./44-slide-presenter-inspector-and-handouts.md) | Presenter & Print | Presenter inspector route (`/slides/inspector`), 3-up printable handouts, print mode |
| 05c | [05-bright-gold-tech/readme.md](./05-bright-gold-tech/readme.md) | Theme | Bright gold tech deck and website bands |
| 98 | [98-confidence-report.md](./98-confidence-report.md) | Meta | What an agent can build, and what it must refuse |
| 02b | [02-ai-system-design/readme.md](./02-ai-system-design/readme.md) | Sub-System | AI-adaptable modern SaaS theme swapping and section blueprints |
| 03b | [03-sweet-digs-design-system/readme.md](./03-sweet-digs-design-system/readme.md) | Sub-System | Sweet Digs HSL multi-theme architecture and zoom-free interaction suite |
| 04b | [04-white-blue-theme/readme.md](./04-white-blue-theme/readme.md) | Sub-System | White Blue Theme 3-format color architecture (`HEX`, `RGB`/`RGBA`, `HSL` + `OKLCH`), 3D flip Mega-Menu, `WhiteBlueButton`, and enterprise section library |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | Testing | Testable criteria for design system compliance |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Meta | Consistency validation report |

---

## Variable Dependency Architecture

```
┌─────────────────────────────────────────┐
│         CSS Custom Properties           │
│  (index.css :root / .dark)              │
│  --primary, --accent, --background...   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│       Tailwind Config Mapping           │
│  (tailwind.config.ts)                   │
│  primary: "hsl(var(--primary))"         │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│       Component Token Layer             │
│  Semantic classes: .prose-spec,         │
│  .code-block-wrapper, .checklist-block  │
│  All use hsl(var(--token)) only         │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│       Component States                  │
│  :hover, :focus, :active, .dark         │
│  Transform, opacity, box-shadow shifts  │
│  All driven by tokens + CSS3 transitions│
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│       Page-Level Composition            │
│  Pages assemble components              │
│  No page-specific colors or overrides   │
│  Consistent spacing rhythm              │
└─────────────────────────────────────────┘
```

---

## How to Re-Theme

1. Open `index.css`
2. Change HSL values in `:root { }` and `.dark { }` blocks
3. Every component, page, and interaction automatically updates
4. No component files need editing

**Example — Switch from purple/pink to teal/amber:**

```css
:root {
  --primary: 175 85% 40%;        /* was: 252 85% 60% */
  --accent: 38 92% 50%;          /* was: 330 85% 60% */
  --heading-gradient-from: 175 85% 45%;
  --heading-gradient-to: 38 92% 55%;
}
```

Every heading gradient, link color, button, code block glow, and hover effect updates instantly.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CSS Variables Source | `src/index.css` |
| Tailwind Config | `tailwind.config.ts` |
| Spec Authoring Guide | `../01-spec-authoring-guide/readme.md` |
| Docs Viewer UI Spec | `../08-docs-viewer-ui/readme.md` |
| Visual Rendering Guide | `../08-docs-viewer-ui/02-features/07-visual-rendering-guide.md` |

---

## Verification

Full compliance criteria and automated verification suites are defined in [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) and tracked in [`99-consistency-report.md`](./99-consistency-report.md).

