# 41 — Homepage & Blog Section Catalog Specification

> **/goal** Provide the definitive, generic, multi-domain section catalog for building enterprise websites, marketing homepages, blogs, and interactive web surfaces.
> **/learn** Master the geometry, layout, typography, CSS3 animations, and distance parameters across the 15 core homepage sections and complete blog editorial templates.

**Version:** 4.1.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. System Overview & The Alternating Band Rhythm

Web pages compose sequentially from isolated section blocks obeying the **Alternating Band Rhythm**:
- **Tone Cycling:** `light` (`--paper`, `#FFFFFF`) ➔ `soft` (`--surface-soft`, `#F8FAFC`) ➔ `light` ➔ `soft`.
- **Dark Tone Limit:** `dark` (`--surface-dark`, `#0F172A`) is permitted at most ONCE per page (for a high-impact product panel or CTA), never for the hero section.
- **Tone Rule:** Adjacent sections must NEVER share the same background tone.
- **Vertical Spacing:** Standard section padding is `clamp(72px, 8vw, 120px)` vertical. Tight section padding is `clamp(40px, 5vw, 64px)`. Content max-width is `1280px` with gutters: `20px` (<400px), `24px` (mobile), `40px` (tablet), `64px` (desktop).

---

## 2. The 15 Flagship Homepage Sections

### 2.1 Social Proof Hero with Split Capability Console
- **Layout:** 2-Column desktop grid (`lg:grid-cols-[1.1fr_0.9fr]`, gap `48px`).
- **Left Column:** Eyebrow (12px mono uppercase, tracking 0.14em) ➔ H1 Headline (`text-h1`, clamp 38–80px, bold) with one `<GradientText>` accent word ➔ Lead paragraph (`text-lead`, max 28 words) ➔ Button CTA pair (`primary` with `.shine-sweep` + `outline` with `.pointer-fill`) ➔ Social Proof Cluster (40px circular avatar stack with 12px overlap, 5× 16px gold stars, rating count).
- **Right Column (Capability Console):** Rounded `28px` card with ambient halo glow, housing an interactive feature switcher (`ToggleRow`) cycling 3 live capability previews with 420ms morph easing.

### 2.2 Dual Marquee Logo Rail
- **Layout:** Full-width container with feathered horizontal edge masks (`linear-gradient(to right, transparent, black 15%, black 85%, transparent)`).
- **Rails:** Two counter-rotating tracks: Top track scrolls left-to-right (`translateX(-50%)` to `0`), bottom track scrolls right-to-left (`0` to `translateX(-50%)`).
- **Animation:** Continuous `35s linear infinite`, pauses on hover (`animation-play-state: paused`).
- **Logos:** 12–16 vector partner logos, height `32px`, grayscale `opacity: 0.6`, transitions to full color on hover.

### 2.3 Feature Block Grid
- **Layout:** Responsive 6-card bento grid (`grid-cols-1 md:grid-cols-2 lg:grid-cols-3`, gap `32px`).
- **Card Anatomy:** Rounded `20px`, border `1px solid var(--border)`, background `var(--card)`, internal padding `32px`.
- **Interactions:** Cursor-tracking `.spotlight` radial highlight, `48×48px` icon container with gradient ring, card title in `text-h3` (20–26px Ubuntu bold), muted body copy, hover lift `-translate-y-1` with shadow elevation.

### 2.4 Tabbed Product Panel
- **Layout:** Centered tab switcher bar mounted above a full-bleed application preview card.
- **Tab Bar:** Pill container with active animated background indicator (`layoutId="activeTab"` over 260ms), 3–5 category tabs (`px-5 py-2.5`, text 14px font-medium).
- **Content Surface:** Dark glass or clean light surface, displaying high-fidelity UI dashboard mockup, interactive table, or syntax-highlighted code block with line numbers and copy button.

### 2.5 Integration Constellation
- **Layout:** Central core platform node surrounded by 8–12 orbiting third-party technology nodes.
- **Connectors:** Responsive SVG bezier paths connecting each satellite node to the center core.
- **Animation:** Flowing light pulses running along the paths using CSS `@keyframes` with `stroke-dasharray: 6 8` and `stroke-dashoffset` animation over `1.4s`.

### 2.6 Metric Counter Band
- **Layout:** 4-Column stats ribbon on `soft` band (`grid-cols-2 lg:grid-cols-4`, gap `32px`, py `64px`).
- **Stat Item:** Giant numeric counter (`text-mega`, 44–120px Ubuntu bold, tabular numbers) ➔ Animated on view entry via `CountUp` hook over `1.6s` ➔ Trending delta badge (`+142% YoY` in emerald pill) ➔ Stat label in 14px Poppins medium.

### 2.7 Value Card Deck with 3D Tilt
- **Layout:** 3-Column card deck (`grid-cols-1 lg:grid-cols-3`, gap `32px`).
- **Card Physics:** Wrapped in `TiltCard` with pointer 3D tilt (max rotation `8deg`, perspective `1200px`).
- **Card Decor:** Radial ambient halo (`--gradient-halo`), top numbered eyebrow badge (`01`, `02`, `03`), bold headline, bullet deliverables, and hover gradient border.

### 2.8 Industry Hover Grid
- **Layout:** 4–6 Horizontal expandable accordion rows or 2-column capability cards.
- **Hover Mechanics:** Hovering a row expands its vertical height from `72px` to `180px`, revealing sub-capabilities, verified case metrics, and an arrow link with `SlideSwapLabel`.

### 2.9 Testimonial Drag Rail
- **Layout:** Full-bleed horizontal slider container with draggable snap cards.
- **Physics:** Grab cursor (`cursor: grab`, `cursor: grabbing` on active), momentum touch drag, card width `440px`, gap `28px`.
- **Testimonial Card:** Quote in 18px Poppins italic, author avatar `56×56px`, author name in 16px bold, enterprise badge, and bottom progress track indicator.

### 2.10 Sticky Process Stack
- **Layout:** 2-Column pinned container (`lg:grid-cols-[1fr_1.2fr]`, min-h `200vh`).
- **Left Column:** Pinned sticky block (`position: sticky, top: 120px`) displaying section title, lead deck, and active stage counter.
- **Right Column:** 4 Sequential process cards scrolling past. As each card activates, an SVG connector line (`note-connector-draw`) draws itself to the next step.

### 2.11 Pricing Plan Cards
- **Layout:** 3-Column commercial tier grid (Starter, Growth, Enterprise).
- **Billing Switch:** Monthly / Annual billing pill toggle with "Save 20%" badge.
- **Featured Tier:** Center "Growth" card elevated with `scale: 1.03`, gradient border (`2.5px solid var(--primary)`), elevated shadow (`var(--shadow-lift)`), and prominent "Most Popular" ribbon.
- **Checklist:** 8–12 Verifiable feature rows with green checkmark SVGs.

### 2.12 FAQ Split Accordion
- **Layout:** 2-Column split section (`lg:grid-cols-[1fr_1.4fr]`, gap `64px`).
- **Left Column:** Section headline ("Frequently Asked Questions"), descriptive support copy, and "Still have questions? Contact our architects" card with direct CTA button.
- **Right Column:** 5–8 Accordion items using zero-JS fluid grid animation (`grid-template-rows: 0fr` to `1fr` over 240ms), plus/minus icon toggle, and accessible ARIA attributes.

### 2.13 Oversized CTA Footer Banner
- **Layout:** Boxed floating card (`max-w-[1200px] mx-auto`, rounded `28px`, py `80px`, px `48px`).
- **Surface:** Ambient gradient backdrop (`--gradient-hero` or `--gradient-brand`) with floating render animation (`cta-float`, 14px gentle translation over 6s infinite).
- **Content:** Bold 48px heading, explanatory deck, primary CTA button with `.shine-sweep`, and secondary text link.

### 2.14 Motion Primitives & Interaction Layer
- **Standardized CSS Interactions:**
  - `.shine-sweep`: Angled 45-degree light sweep running across CTAs on hover (`3s linear infinite`).
  - `.pointer-fill`: Smooth background expansion tracking cursor entry angle.
  - `.slide-swap`: Dual-layer vertical character swap on hover (`520ms ease-out`).
  - `.spotlight`: Cursor-following radial light highlight on card hover.
  - `.note-connector-draw`: Animated SVG dashed path drawing itself on viewport entry.

### 2.15 Scroll Stack Cards
- **Layout:** Stack of 3–5 full-width cards pinned sequentially to viewport top (`position: sticky, top: 100px`).
- **Stacking Physics:** Each incoming card overlaps the previous card. Previous cards progressively scale down (`scale: calc(1 - (index * 0.04))`), darken by 10%, and gain subtle backdrop blur (`blur(2px)`).

---

## 3. Blog Editorial System Specifications

### 3.1 Blog Index Layout (`/blog`)
1. **Editorial Hero:** Category title ("Engineering & Architecture Insights"), search input with instant filter, and category filter pills (`All`, `Architecture`, `AI & LLM`, `Performance`, `Design Systems`).
2. **Featured Post Card:** Full-width hero article card (height `420px`), split layout with 16:9 photographic plate on left, "FEATURED" pill, title in `text-h2`, excerpt, author avatar, reading time, and publish date.
3. **Article Grid:** 3-Column responsive card grid (`grid-cols-1 md:grid-cols-2 lg:grid-cols-3`, gap `32px`). Card: 16:9 cover image with zoom hover (`scale: 1.04`), category tag, headline (`text-h3`, hover color shift to primary), 2-line excerpt, author metadata row.
4. **Pagination:** Accessible numeric page selector (`‹ 1 2 3 ... 8 ›`).

### 3.2 Blog Post Detail Template (`/blog/:slug`)
1. **Header Zone:** Breadcrumb trail (`Home › Blog › Category`) ➔ Category pill ➔ Article H1 (`text-h1`, clamp 36–64px bold) ➔ Author meta block (author avatar 48×48px, author name with link, role credential, publication date, calculated reading time: `Math.ceil(words / 200)` min read) ➔ Hero photographic plate (aspect ratio 21:9 or 16:9, rounded `24px`).
2. **2-Column Body Layout:**
   - **Sticky Sidebar (Width: 280px, `hidden xl:block`):** `position: sticky, top: 100px`, Table of Contents (TOC) with active section scroll-spy indicator, estimated time remaining bar, and social share button cluster (Twitter/X, LinkedIn, Copy Link).
   - **Main Prose Column (Max Width: 760px):**
     - Typography: Poppins 18px body copy, line-height `1.75`, paragraph gap `28px`.
     - Headings: H2 (`text-h2`, mt 56px, mb 24px), H3 (`text-h3`, mt 40px, mb 16px).
     - Code Blocks: Dark slate background (`#0F172A`), rounded `16px`, JetBrains Mono 14px, top window control dots, language badge, line numbers, and clipboard copy button.
     - Callout Blocks: `info` (blue hairline + icon), `warning` (amber hairline), `tip` (emerald hairline).
     - Pull Quotes: Left accent border `3px solid var(--primary)`, 24px italic text, author attribution.
3. **Footer Zone:** Author bio box (avatar, bio, social links) ➔ "Share this article" bar ➔ Related Articles 3-card grid ➔ Newsletter subscription card.

---

## 4. Anti-Hallucination & Quality Verification Checklist

- [ ] All 15 homepage sections implement the alternating band rhythm (`light` -> `soft` -> `light`).
- [ ] No section uses hardcoded color values; all styles compose from semantic design tokens.
- [ ] Content limits are strictly enforced: Eyebrow <= 3 words, H2 <= 8 words, Lead <= 28 words, Card <= 32 words.
- [ ] Blog post detail template includes sticky Table of Contents sidebar with active scroll-spy.
- [ ] Code blocks in blog posts include language tag, line numbers, and copy button.
- [ ] All animations provide `@media (prefers-reduced-motion: reduce)` fallbacks.
- [ ] All interactive cards use standardized hover lifts (`-translate-y-1`) with non-blocking composited transforms.
