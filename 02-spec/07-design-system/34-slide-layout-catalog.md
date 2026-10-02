# 34 — Master Slide Layout Catalog & Pure DOM Typography Specification

> **/goal** Provide the definitive, comprehensive layout catalog for 16:9 presentation slides on the 1920×1080 virtual canvas with pure DOM typography enforcement and zero baked-in text.
> **/learn** Master the coordinate geometries, slot models, typography scales, and visual zones across the 20 master slide layouts across all presentations: Title Hero, Executive Persona, Key Player Bio, Before/After Split, USP Strikethrough, SaaS Pricing, Steps Chain Roadmap, Social Proof, Talent Funnel, 3-Point Master Cards, Center Punchy Headline, Left Editorial, One-Liner Quote, Process/Timeline, Counter Stat, Reveal Grid/Depth Stack, Embed, Poll/Q&A, Tabletop Hardware & Bike Showcase, and Tech Stack Matrix.

**Version:** 4.3.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. The Non-Image Text Mandate (Pure DOM Typography)

> [!CRITICAL]
> **TOTAL BAN ON BAKED-IN TEXT:**
> Headlines, subtitles, kickers, bullet points, numbered metrics, author bios, and captions MUST ALWAYS be rendered as live, selectable DOM HTML elements (`<h1>`, `<h2>`, `<p>`, `<span>`, `<div>`) styled with CSS typography tokens.
> Under no circumstances should text be flattened into raster images (`.png`, `.jpg`, `.webp`). Raster images are strictly reserved for photographic hero visual plates, author avatars, and partner logos.

---

## 2. 1920×1080 Coordinate Space & Scaling Architecture

All slide layouts operate on a virtual reference coordinate grid of `1920 × 1080` pixels:
- At runtime, `ScaledSlide` detects container dimensions using `ResizeObserver`.
- Computes uniform scale factor: $\text{scale} = \min(\text{width}/1920, \text{height}/1080)$.
- Applies vector scaling: `transform: scale(var(--stage-scale))` with `transform-origin: center center`.
- Container specifies: `contain: layout paint; isolation: isolate; will-change: transform;`.

---

## 3. The 20 Master Slide Layouts

### 3.1 Layout 1: Title & Hero Slide (`type: "title"`)
- **Top Kicker Pill:** Left `140px`, Top `120px`, `14px` mono uppercase, tracking `0.25em`.
- **Main Headline:** Left `140px`, Top `260px`, Max width `1400px`, `82px` Ubuntu black, leading `1.05`.
- **Subtitle:** Left `140px`, Top `460px`, Max width `1050px`, `26px` Poppins, leading `1.4`.
- **Presenter Bio Card:** Left `140px`, Top `640px`, Avatar `60×60px`, Name `22px`, Role `16px`.
- **Bottom Organic Wave:** Bottom `0px`, Width `1920px`, Height `140–180px` dual-gradient SVG ribbon.

### 3.2 Layout 2: Executive Persona & CEO Slide (`type: "executive-persona"` / `type: "persona"`)
- **Portrait Staging Zone:** Left `-150px`, Width `1100–1500px` with halftone dot matrix and radial halo.
- **Hero Name:** `68–104px` Ubuntu black, character-by-character color stepping ($S_4$ to $S_0$).
- **Credential Grid:** 2-column experience & impact metrics (`40px` stat, `14px` mono uppercase label).

### 3.3 Layout 3: Key Player Bio Slide (`type: "key-player"`)
- **Header:** Left `140px`, Top `120px`, Title `52–54px`, Subtitle `18–22px`.
- **3–4 Member Grid:** Gap `32–40px`, Top `280px`, Cards width `380–420px`.
- **Member Card:** Portrait `380×380px` to `420×480px`, Name `24px` Ubuntu bold, Role `15–16px`, Bio `14px`.

### 3.4 Layout 4: Before / After Split Showcase (`type: "before-after"`)
- **Left Card ("Before"):** Width `790px`, background `#FFF5F5`, border `2px solid #FECDD3`, pain points in rose badges.
- **Right Card ("After"):** Width `790px`, background `#FFFFFF`, border `2.5px solid #7C3AED` or `#10B981`, verified proof metrics with checkmarks. Optional interactive slider wipe handle.

### 3.5 Layout 5: USP Strikethrough Strike Slide (`type: "usp-strike"`)
- **Headline Typography:** `124px` Ubuntu bold, line-height `1.02`, letter-spacing `-0.03em`.
- **Editorial Strikethrough:** `textDecoration: "line-through"`, `textDecorationColor: "hsl(var(--pres-accent) / 0.7)"`, thickness `6px`.
- **3-Point Proof Cluster:** Bottom horizontal card strip, Height `180px`, gap `32px`.

### 3.6 Layout 6: SaaS Pricing & Metric Proof Slide (`type: "pricing"`)
- **Tiers:** 3 Columns, Width `500px` each, gap `32–40px`, Top `270px`, Height `700px`.
- **Featured Plan ("Hot"):** Border `2.5px solid #7C3AED`, shadow `0 24px 48px -12px rgba(124, 58, 237, 0.18)`, `scale: 1.03`, "MOST POPULAR" ribbon.
- **Price Figures:** `48–52px` Ubuntu bold. CTAs: `52px` height, rounded `12px`.

### 3.7 Layout 7: Steps Chain & Process Roadmap Slide (`type: "steps-chain"`)
- **Connecting Horizon Line:** Top `364px`, Width `1640px`, Height `3px`, gradient progress fill.
- **Step Badges:** `48×48px` circular badge centered on track with bold numeral.
- **4 Roadmap Cards:** Width `370px` each, gap `53px`, Top `410px`, Height `440px`, duration pill + deliverables.

### 3.8 Layout 8: Social Proof & Testimonials Slide (`type: "testimonials"`)
- **Dual Quote Cards:** Width `790px` each, gap `40px`, Top `280px`, Height `460px`.
- **Quote Typography:** `22–26px` Poppins italic, color `#1E293B`, line-height `1.5`.
- **Bottom Partner Logo Rail:** Top `820px`, Height `100px`, grayscale opacity `0.6` hover `1.0`.

### 3.9 Layout 9: Talent Funnel & Capability Stack Slide (`type: "talent-funnel"`)
- **Progressive Funnel Bands:** 4 Bands (1640px -> 1380px -> 1120px -> 860px), Height `110px` each.
- **Funnel Stages:** Pool Vetting -> Algorithmic Testing -> Architecture Simulation -> Deployment.

### 3.10 Layout 10: 3-Point Master Cards Slide (`type: "bullets"`)
- **Left Editorial Block:** Width `900px`, Kicker pill, Headline `56px` Ubuntu bold, Subtitle `24px`.
- **3 Bullet Cards:** Height `88px` each, circular icon container `48×48px`, text `20px` Poppins medium.
- **Right Visual Plate:** Width `700px`, photographic hero plate with feathered left mask and radial glow.

### 3.11 Layout 11: Center Punchy Headline Slide (`type: "center"`)
- **Center Focus:** Centered `88px` Ubuntu bold headline + `28px` Poppins subhead for keynote theses and transitions.

### 3.12 Layout 12: Left Editorial Narrative Column (`type: "left"`)
- **Left Narrative Block:** Left `140px`, Top `180px`, Width `800px`, Headline `64px`, Body `22px` Poppins.
- **Right Content Plate:** Left `1020px`, Width `760px`, Height `720px`, rounded `20px`, border `1px solid var(--border)`.

### 3.13 Layout 13: One-Liner Quote & Priority Attribution (`type: "quote"` / `type: "priority"`)
- **Quotation Mark Glyph:** `120px`, color `hsl(var(--pres-accent) / 0.25)`.
- **Quote Headline:** `60px` Ubuntu italic medium, max width `1400px`. Attribution pill with `52×52px` avatar.

### 3.14 Layout 14: Process & Phased Timeline Flow (`type: "process"` / `type: "timeline"`)
- **Connected Track:** Top `450px`, Height `4px`, gradient fill, with 3–5 Milestone pulse nodes (`56×56px`).
- **Milestone Cards:** Width `340px`, milestone title, date badge, and bullet deliverables.

### 3.15 Layout 15: Counter Stat & KPI Velocity Meter (`type: "counter-stat"`)
- **Metric Columns:** 3 Columns, Width `500px` each, Top `340px`, Height `480px`.
- **Giant Number:** `110px` Ubuntu bold, tabular figures, gradient accent fill (`gradient-text`), +142% delta badge.

### 3.16 Layout 16: Reveal Grid & 3D Depth Stack (`type: "reveal-grid"` / `type: "depth-stack"`)
- **Perspective Container:** `perspective: 1200px`.
- **3-Layer Depth Stack:** Layer 1 (scale 0.92, opacity 0.5) -> Layer 2 (scale 0.96, opacity 0.8) -> Layer 3 (scale 1.0, active).

### 3.17 Layout 17: Live Embed Stage & Interactive Prototype (`type: "embed"`)
- **Stage Frame:** Width `1640px`, Height `760px`, rounded `24px`, border `1px solid var(--border)`.
- **Caption Bar:** Bottom `40px` toolbar showing source URL, refresh trigger, and fullscreen expander.

### 3.18 Layout 18: Interactive Audience Poll & Real-Time Q&A (`type: "poll"` / `type: "qa"`)
- **Headline:** Top `140px`, Width `1100px`, `52px` Ubuntu bold.
- **QR Code Card:** Right `140px`, `280×280px` white card with scannable room code.
- **4 Vote Option Bars:** Width `1100px`, Height `72px`, spring-fill percentage bars with mono tallies.

### 3.19 Layout 19: Tabletop Hardware & Bike Physical Showcase (`type: "tabletop-bike"`)
- **Photographic Cutout Stage:** Left `140px`, Width `880px`, Height `720px`, transparent cutout with contact ground shadow (`radial-gradient`) and interactive pulsing beacon hot-spots (`28×28px`).
- **Spec & Telemetry Grid:** Left `1060px`, Width `720px`, 4 cards in 2×2 grid (`"750W Peak"`, `"60 Miles Range"`, `"48V 15Ah"`, `"IP67 Sealed"`).

### 3.20 Layout 20: Tech Stack & Capability Matrix (`type: "tech-stack"`)
- **4 Categorized Columns:** Frontend & UI, Backend & Distributed Core, Data & Storage, AI & CI/CD Pipelines.
- **Category Card:** Rounded `20px`, border `1px solid var(--border)`, containing 4–6 interactive tech pill badges.

---

## 4. Anti-Hallucination & Quality Verification Checklist

- [ ] All 20 layout models author dimensions strictly declared on the `1920×1080` canvas.
- [ ] ZERO text is flattened or baked into images; all headings and body copy are live DOM nodes.
- [ ] Headings strictly use `Ubuntu`, body text uses `Poppins`, and metadata uses `JetBrains Mono`.
- [ ] Step reveals preserve dimmed state (`opacity: 0.15; filter: blur(2px)`) until activated.
- [ ] Pricing Hot tiers use `scale: 1.03` with a gradient border and elevated shadow.
- [ ] Before/After split uses high-contrast rose muted vs violet/emerald vibrant styling.
- [ ] Tabletop hardware showcase includes transparent cutout bounds, contact shadow, and hot-spot badges.
