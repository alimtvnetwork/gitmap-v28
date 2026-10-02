# Follow UI/UX Design System Specification & Component Assembly

> **Prompt Version:** 2.1.0
> **Trigger keywords:** `ui-ux`, `follow-design-system`, `build-website`, `avant-garde-menu`, `design-tokens`, `button-system`, `blog-system`

**/goal** Autonomously construct production-grade marketing websites, blogs, dashboards, navigation systems, avant-garde mega menus, buttons, and responsive sections by strictly following the grounded design specifications without inventing a single value or hallucinating aesthetics.

**/learn** Before authoring HTML, JSX, CSS, LESS, or Tailwind code, AI agents MUST read the sequentially numbered specification files listed below in exact order. All relative paths are from the git repository root.

---

## 1. Mandatory Reading Sequence (Strict Relative Paths)

AI agents MUST sequentially ingest these specification files:

1. `02-spec/07-design-system/02-design-principles.md` — 60/30/10 visual balance, 4-plane depth hierarchy, single-accent Von Restorff rule.
2. `02-spec/07-design-system/04-typography.md` — Fluid type scale (`clamp()`), Ubuntu headings, Poppins body, JetBrains Mono eyebrows.
3. `02-spec/07-design-system/05-spacing-layout.md` — Content max width (`1280px`), responsive gutters, vertical section rhythms.
4. `02-spec/07-design-system/10-header-navigation.md` — 72px sticky glass header, 12px scroll threshold, safe-region pointer physics (`pad = 14px`, 220ms debounce), `SlideSwapLabel` nav links, sticky-safe mobile drawer.
5. `02-spec/07-design-system/11-button-system.md` — 6 Button variants (`primary`, `solid`, `outline`, `glass`, `ghost`, `link`), 4 sizes (`sm: 36px`, `md: 44px`, `lg: 52px`, `icon: 44px`), magnetic cursor physics (`strength: 0.22`), `.shine-sweep`, `.pointer-fill`, capsule pills.
6. `02-spec/07-design-system/16-theme-catalogue-and-palettes.md` — Multi-theme catalog, semantic color token mappings, contrast verification.
7. `02-spec/07-design-system/21-css3-animations-and-interactions.md` — Hardware-accelerated CSS3 keyframes, marquees, accordions, neon borders.
8. `02-spec/07-design-system/25-page-assembly.md` — Site and blog shells, copy length constraints.
9. `02-spec/07-design-system/33-mega-menu-components.md` — Avant-garde mega panel dropdown: 260ms entrance, multi-column grid, growing left hairline over 420ms, 3D promotional flip card (`perspective: 1400px`, `rotateY(180deg)` over 820ms).
10. `02-spec/07-design-system/36-website-content-builder-mode.md` — Comprehensive website visual builder overlay: review-only gate (`?builder=1&email={OWNER_EMAIL}`), in-place `contentEditable` text editor, images/icon replacement modal, menu link editor, floating side panel diff tracking, and deterministic `content-changes--all-pages--*.zip` export.
11. `02-spec/07-design-system/37-image-specifications.md` — Infographics, social banners, YouTube thumbnails, safe zones, and aspect ratio standards.
12. `02-spec/07-design-system/38-card-and-pricing-components.md` — 3-tier pricing table matrix, featured badges, floating editorial sections.
13. `02-spec/07-design-system/41-homepage-and-blog-sections.md` — 15 Flagship homepage sections and complete blog editorial templates.

---

## 2. The Avant-Garde Navigation & Mega Menu Rules

When authoring header navigation, enforce these exact mechanics:

1. **Header Geometry & Surface:**
   - Height: Strictly `72px` (`h-[72px]`), sticky positioning `top-0 z-50 w-full`.
   - Surface: `bg-background/90 text-foreground backdrop-blur-xl`.
   - Scroll threshold: Viewport scroll > `12px` activates `border-b border-border shadow-[var(--shadow-card)]` easing over `420ms cubic-bezier(0.16, 1, 0.3, 1)`.
   - Alternative shrinking pill header: `1240px` transparent container shrinking to a `900px` white pill after scroll.
2. **Nav Link Micro-Interactions:**
   - `SlideSwapLabel`: Per-character vertical slide swap on hover with `stagger: 0.04s` for nav links (`0.018s` for buttons). Top glyph slides to `translateY(-110%)`; duplicate bottom glyph enters from `translateY(100%)` to `translateY(0)` over `520ms var(--ease-out)`. Under reduced motion, collapses to a plain span.
   - Underline indicator: `1px` rule, `origin-left scale-x-0 group-hover:scale-x-100` over `240ms`.
   - Chevron: `14px` (`size-3.5`), rotates `180deg` when dropdown panel is active over `240ms`.
3. **Pointer Safe Region & Debounced Close Geometry:**
   - Safe region buffer: `pad = 14px` around both header and dropdown panel rects.
   - Pointer moving within safe region cancels close; pointer moving outside schedules close after `220ms` debounce.
   - `Escape` key or page scroll dismisses dropdown immediately (`0ms`).
4. **MegaPanel Dropdown Surface:**
   - Mount point: `absolute left-0 right-0 top-full z-40 pt-3`.
   - Entrance: `initial: { opacity: 0, y: -8, scale: 0.985 }`, `animate: { opacity: 1, y: 0, scale: 1 }`, duration `260ms`, ease `cubic-bezier(0.16, 1, 0.3, 1)`. Exit: `opacity: 0, y: -6, scale: 0.99`. Reduced motion: opacity only, `120ms`.
   - Panel card: `rounded-[20px] border border-border bg-card shadow-[var(--shadow-lift)] p-8 gap-8`.
   - Group stagger: `opacity: 0, y: 8 -> 1, 0`, duration `280ms`, delay `0.05s + gi * 0.05s`.
   - Link stagger: `opacity: 0, x: -6 -> 1, 0`, duration `260ms`, delay `0.08s + gi * 0.05s + li * 0.03s`.
   - Link item: `rounded-[10px] px-3 py-2`. Left accent hairline `1px` grows `scaleY(0) -> scaleY(1)` with gradient accent over `420ms`. Trailing arrow `size-3.5` transitions from `-translate-x-1 opacity-0` to `translate-x-0 opacity-100`.
   - 3D Flip Promo Card: `min-h-[220px]`, `perspective: 1400px`, `transform-style: preserve-3d`, flips `rotateY(180deg)` over `820ms`. Front face has gradient accent; back face has card border, title, body, and CTA pill button.
5. **Sticky-Safe Mobile Drawer:**
   - Opens below `lg` breakpoint fixed `top-[72px] bottom-0`.
   - Intercepts window `wheel` and `touchmove` events outside the menu drawer **WITHOUT** setting `overflow: hidden` on `body` (which breaks `position: sticky`).

---

## 3. Universal Button System Standards

Every interactive button MUST implement one of the specified variants from `11-button-system.md`:

1. **`primary`:** `bg-[image:var(--gradient-accent)] text-white shadow-[var(--shadow-card)]`. Carries active `.shine-sweep` angled 45-degree sheen on hover/focus.
2. **`solid`:** `bg-brand-primary text-white shadow-[var(--shadow-card)]`. Active `.shine-sweep` on hover.
3. **`outline`:** `border border-border bg-transparent text-foreground`. Carries active `.pointer-fill` filling smoothly from cursor entry angle.
4. **`glass`:** `bg-[image:var(--gradient-card-dark)] text-on-dark backdrop-blur-md`. Carries active `.gradient-ring` border mask.
5. **`ghost`:** Zero-chrome utility, `hover:bg-muted`.
6. **`link`:** Text anchor, `hover:underline underline-offset-4 px-0`.
7. **`ShineButton` (Header High-Conversion CTA):** Near-black pill (`bg-night`) with continuous 45-degree linear sheen sweep (`@keyframes shine`) and an orbiting 8-particle floating spark constellation (`hidden md:block`, `@keyframes spark-float`). Inverse variant: white pill on dark imagery.
8. **Sizing:** Strictly `sm: 36px` (`h-9 px-4 text-[13px]`), `md: 44px` (`h-11 px-5 text-sm`), `lg: 52px` (`h-13 px-7 text-[15px]`), `icon: 44px` (`size-11 px-0`).
9. **Magnetic Physics:** Primary hero CTAs wrap in `<Magnetic strength={0.22}>`.
10. **Capsule Pill Buttons:** Follow tokenized `--capsule-gold`, `--capsule-ember`, and near-black gradient text (`0 0% 4%`) for WCAG 2.2 AA contrast.

---

## 4. Band Rhythm & Section Architecture

1. **Alternating Band Rhythm:**
   - Sections cycle tones: `light` (`--paper`) -> `soft` (`--surface-soft`) -> `light` -> `soft` -> `light`.
   - `dark` (`--surface-dark` / void) is used at most ONCE per page for an embedded product or testimonial panel.
   - **Total Ban:** Adjacent sections must NEVER share the same tone.
2. **Vertical Rhythm:**
   - Standard section padding: clamp(72px, 8vw, 120px) vertical padding.
   - Tight section padding: clamp(40px, 5vw, 64px).
3. **Content Limits:**
   - Section Eyebrow: At most 3 words.
   - Section H2: At most 8 words.
   - Lead paragraph: At most 28 words.
   - Card body copy: At most 32 words.
   - Accent phrases: At most ONE gradient accent phrase (`<GradientText>`) per heading.
4. **High-Conversion Showcase Sections:** Follow `38-card-and-pricing-components.md` and `41-homepage-and-blog-sections.md`:
   - **Interactive Split Pricing:** Interactive option buttons on left + warm cream price card (`#FDFBF7`) on right.
   - **Scroll-Driven Services Showcase (`WhatWeDo`):** Dark shell (`#0A0A0A`, radius 32px), sticky left column with `translateY(-active * panelHeight)` filmstrip copy slide, and paired staggered project screenshot blocks.
   - **Ruled-Paper Process Reassurance (`Workflow`):** `repeating-linear-gradient` paper overlay with 5 alternating tilted sticky notes (`-1.5deg` to `+1.8deg`), glowing pin circle, script numerals, and hover lift.
   - **Enterprise Ready Hero:** Split hero with gradient headline and 7-card depth-scaled toggle stack over blurred violet halo.
   - **15 Core Homepage Patterns:** Follow `41-homepage-and-blog-sections.md`.

---

## 5. Blog Editorial System Architecture

When authoring blog pages, follow `41-homepage-and-blog-sections.md`:
1. **Blog Index (`/blog`):** Hero search, category filter pills, featured article card with reading time, 3-column article cards, and pagination.
2. **Blog Post Detail (`/blog/:slug`):**
   - Header with category pill, title, author avatar + credential, publication date, calculated reading time (`Math.ceil(words / 200)` min read).
   - Sticky sidebar (`width: 280px`): Table of Contents with active scroll-spy indicator and social share bar.
   - Main prose column (`max-w-[760px]`): Poppins 18px body, line-height 1.75, code blocks with dark slate background, language badge, line numbers, and copy button.
   - Inline callout blocks (`info`, `warning`, `tip`) and pull quotes with left accent border.
   - Author bio box, related articles 3-card grid, and newsletter card.

---

## 6. Client-Side Website Content Builder Readiness

When building marketing or blog pages, follow `02-spec/07-design-system/36-website-content-builder-mode.md` and `26-visual-builder.md`:
1. Open the editor with `?builder=1&email={OWNER_EMAIL}`. Do not use `?builder=true`.
2. Id is the content-module source key, else the DJB2 hash in that file. Do not invent `data-edit-id`.
3. Accent stays a `SPAN` with `class="gradient-text"`. Do not invent `[[accent]]` markers.

---

## 7. Blind-AI Execution Checklist

Before outputting code or completing your task, verify every item:

- [ ] Header height is strictly `72px` with `top-0 z-50` sticky positioning and `12px` scroll threshold.
- [ ] Safe-region pointer buffer is `14px` with a `220ms` debounce timer.
- [ ] Primary nav links use `SlideSwapLabel` with `stagger: 0.04s` and `520ms` roll transition.
- [ ] MegaPanel entrance uses `260ms cubic-bezier(0.16, 1, 0.3, 1)` from `y: -8, scale: 0.985`.
- [ ] Growing left hairline in mega links uses `scaleY(0) -> scaleY(1)` over `420ms`.
- [ ] Promo flip card uses `perspective: 1400px` and flips `rotateY(180deg)` over `820ms`.
- [ ] Mobile drawer NEVER touches `document.body.style.overflow`.
- [ ] Buttons strictly use the 6 variants, 4 sizing tiers, and capsule pill buttons where required.
- [ ] Typography strictly follows `Ubuntu` (display/headings), `Poppins` (body/UI), `JetBrains Mono` (eyebrows/counters).
- [ ] Band tones alternate (`light -> soft -> light -> soft`); adjacent sections never share the same tone.
- [ ] ZERO hardcoded hex colors or inline style overrides are used in component classes; all colors bind to semantic CSS variables.
- [ ] All relative paths referenced start from the git repository root.
- [ ] NO private or company names are present in public files or examples.
