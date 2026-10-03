---
name: follow-ui-ux-design-system
description: Autonomously construct production-grade marketing websites, blogs, dashboards, navigation systems, avant-garde mega menus, buttons, and responsive sections by strictly following grounded design system specifications without inventing values or hallucinating aesthetics.
---

# Follow UI/UX Design System Specification & Component Assembly

> **[/goal](slashCommand;goal)** Autonomously construct production-grade marketing websites, blogs, dashboards, navigation systems, avant-garde mega menus, buttons, and responsive sections by strictly following grounded design specifications without inventing values or hallucinating aesthetics.
> **[/learn](slashCommand;learn)** Ingest the sequentially numbered design system specifications under `02-spec/07-design-system/` in exact order before authoring HTML, JSX, CSS, LESS, or Tailwind code.

**Version:** 2.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. When to Use

Activate this skill when:
- Building marketing pages, SaaS dashboards, blogs, or documentation layouts.
- Implementing navigation headers, sticky glass surfaces, or avant-garde mega panels.
- Constructing button component variants (`primary`, `solid`, `outline`, `glass`, `ghost`, `link`) and magnetic pill interactions.
- Designing responsive layouts, pricing tables, hero sections, and feature grids.
- Operating within website content builder modes (`?builder=1`).

---

## 2. Mandatory Reading Sequence (Strict Relative Paths)

AI agents MUST sequentially ingest these specification files:

1. `02-spec/07-design-system/02-design-principles.md` — 60/30/10 visual balance, 4-plane depth hierarchy, single-accent Von Restorff rule.
2. `02-spec/07-design-system/04-typography.md` — Fluid type scale (`clamp()`), Ubuntu headings, Poppins body, JetBrains Mono eyebrows.
3. `02-spec/07-design-system/05-spacing-layout.md` — Content max width (`1280px`), responsive gutters, vertical section rhythms.
4. `02-spec/07-design-system/10-header-navigation.md` — 72px sticky glass header, 12px scroll threshold, safe-region pointer physics (`pad = 14px`, 220ms debounce), `SlideSwapLabel` nav links, sticky-safe mobile drawer.
5. `02-spec/07-design-system/11-button-system.md` — 6 Button variants, 4 sizes (`sm: 36px`, `md: 44px`, `lg: 52px`, `icon: 44px`), magnetic cursor physics (`strength: 0.22`), `.shine-sweep`, capsule pills.
6. `02-spec/07-design-system/16-theme-catalogue-and-palettes.md` — Multi-theme catalog, semantic color token mappings, contrast verification.
7. `02-spec/07-design-system/21-css3-animations-and-interactions.md` — Hardware-accelerated CSS3 keyframes, marquees, accordions, neon borders.
8. `02-spec/07-design-system/25-page-assembly.md` — Site and blog shells, copy length constraints.
9. `02-spec/07-design-system/33-mega-menu-components.md` — Avant-garde mega panel dropdown: 260ms entrance, multi-column grid, growing left hairline, 3D flip card.
10. `02-spec/07-design-system/36-website-content-builder-mode.md` — Website visual builder overlay, in-place text editor, diff tracking, and zip export.
11. `02-spec/07-design-system/37-image-specifications.md` — Infographics, social banners, thumbnails, safe zones, aspect ratios.
12. `02-spec/07-design-system/38-card-and-pricing-components.md` — 3-tier pricing table matrix, featured badges, floating editorial sections.
13. `02-spec/07-design-system/41-homepage-and-blog-sections.md` — 15 Flagship homepage sections and complete blog editorial templates.

---

## 3. Core Architectural Rules

1. **Total Ban on Invented Values:** Every hex, pixel, millisecond, and cubic-bezier easing MUST originate from `02-spec/07-design-system/`.
2. **Pure DOM Text Rendering Mandate:** Never flatten text into raster images. All copy must be live selectable HTML elements (`<h1>`, `<h2>`, `<p>`, `<span>`).
3. **Strict Relative Git Paths:** No absolute paths or `file:///` URIs.
4. **Strict Lowercase File Naming:** All generated files must be lowercase kebab-case.
5. **Strict Boolean Standards:** Positive booleans evaluated implicitly (`if isReady`). Zero explicit `== true`.
