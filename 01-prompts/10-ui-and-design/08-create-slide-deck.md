# Create Presentation Slide Deck & Live Builder System

> **Prompt Version:** 2.1.0
> **Trigger keywords:** `create-slide-deck`, `presentation-spec`, `slide-builder`, `slide-layout`, `deck-system`, `presenter-view`

**/goal** Autonomously architect, construct, and style high-authority, declarative, 16:9 presentation slide decks, interactive presenter HUDs, 10 production themes, 10-step precision color ramps, and real-time Slide Builder Modes on the 1920×1080 virtual canvas with pure live DOM typography and zero hallucinations.

**/learn** Before writing slide schemas, components, or JSON decks, AI agents MUST read the sequentially numbered specification files listed below in exact order. All relative paths are from the git repository root.

---

## 1. Mandatory Reading Sequence (Strict Relative Paths)

AI agents MUST sequentially ingest these specification files:

1. `02-spec/07-design-system/24-slide-presentation-system.md` — Responsive coordinate transforms (`1920×1080` canvas), dual-screen presenter console (`/present`), BroadcastChannel sync, webcam PIP overlay, step reveals, overview grid (`G`), top jumper (`J`).
2. `02-spec/07-design-system/27-slide-canvas-and-themes.md` — Virtual canvas scale formula, 10 production themes (`bright-gold` default, `noir-gold`, `vscode-dark`, `dracula`, `monokai`, `github-light`, `paper-ink`, `macos-sonoma`, `windows-11`, `navy-blue`), and the Light Theme Contract.
3. `02-spec/07-design-system/28-slide-layouts.md` — Base slot model, closed layout union, and field constraints.
4. `02-spec/07-design-system/29-slide-navigation-and-builder.md` — HUD, transitions, hotkeys (`ArrowRight`, `Space`, `B`/`E`, `G`, `J`, `C`, `F`, `1`–`0`), and builder chrome.
5. `02-spec/07-design-system/30-slide-palette-type-and-shell.md` — Dark amber palette, typography scales, spotlight radial glow, shell layers.
6. `02-spec/07-design-system/31-slide-controller-buttons.md` — Floating top-right HUD controller pill, circular 40×40px action buttons, tabular numbers counter, deep link share, fullscreen toggle, theme popover with live WCAG contrast, webcam PIP overlay, timer, and bottom dot pagination.
7. `02-spec/07-design-system/32-slide-color-options.md` — 10-step gradient precision system ($S_0$ to $S_9$) across 7 flagship palettes, per-slide gradient editor, character-by-character shading engine, 7 semantic pill presets, relative luminance formula, and 9-cell text alignment.
8. `02-spec/07-design-system/34-slide-layout-catalog.md` — The Non-Image Text Mandate (pure DOM typography; zero baked-in text) and the 20 master enterprise slide layouts: Title Hero, Executive Persona, Key Player Bio, Before/After Split, USP Strikethrough, SaaS Pricing, Steps Chain Roadmap, Social Proof, Talent Funnel, 3-Point Master Cards, Center, Left, One-Liner Quote, Process/Timeline, Counter Stat, Reveal Grid/Depth Stack, Embed, Poll/Q&A, Tabletop Hardware & Bike Showcase, and Tech Stack Matrix.
9. `02-spec/07-design-system/35-slide-builder-canvas-inspector.md` — Decoupled dual-store architecture (`useDeckStore` persisted vs `useEditStore` ephemeral), 7 visual canvas stacking layers, draggable `BuilderPanel`, `convertSlideType` engine across all 20 layouts, builder hotkeys (`B`/`E`, `Tab`, `Cmd+Z`, `1`–`0`), bounding box coordinate overrides, and headless Chromium print-ready PDF/handout exports.
10. `02-spec/07-design-system/42-slide-step-and-sound-system.md` — `StepTimelineSlide` vs `AdvanceStepSlide`, vertical chain, active focus row, gold/ember capsules, and Web Audio API synthesis engine.
11. `02-spec/07-design-system/43-slide-webcam-overlay.md` — Presenter webcam PIP overlay, squircle gold rim, S/M/L/XL presets, autoframe face tracker, and keyboard shortcuts (`C`, `M`, `F`, `+`, `-`, `O`, `H`, `1`).
12. `02-spec/07-design-system/44-slide-presenter-inspector-and-handouts.md` — Dedicated Presenter Inspector route (`/slides/inspector`), 3-up printable handouts with ruled writing lines, and print mode.
13. `02-spec/07-design-system/37-image-specifications.md` — Non-image typography mandate, image plates, safe zones, avatar rules, and asset bounds.
14. `02-spec/07-design-system/38-card-and-pricing-components.md` — SaaS pricing card matrix, feature list hierarchy, and highlighted tier geometry.

---

## 2. The Non-Image Text Mandate (Pure DOM Typography)

> [!CRITICAL]
> **TOTAL BAN ON BAKED-IN TEXT:**
> All headlines, subheads, kickers, bullet points, metrics, author names, credentials, and captions MUST be rendered as live, selectable DOM HTML elements (`<h1>`, `<h2>`, `<p>`, `<span>`) styled with CSS typography tokens.
> Under no circumstances should text be flattened into raster images (`.png`, `.jpg`, `.webp`). Raster images are strictly reserved for photographic hero visual plates, author avatars, and partner logos.

---

## 3. 1920×1080 Virtual Canvas & Adaptive Uniform Scaling

All slide components author geometry on an authoritative reference canvas of `1920 × 1080` pixels:

1. **Runtime Scaling Formula:**
   $$\text{scale} = \min\left(\frac{\text{viewportWidth}}{1920}, \frac{\text{viewportHeight}}{1080}\right)$$
2. **Container CSS Vector Transform:**
   ```css
   .slide-stage {
     width: 1920px;
     height: 1080px;
     transform-origin: center center;
     transform: scale(var(--stage-scale));
     contain: layout paint;
     isolation: isolate;
     will-change: transform;
   }
   ```
3. **Stage Positioning:** Centered both vertically and horizontally in the browser viewport with letterboxing on non-16:9 displays.

---

## 4. Floating Controller HUD, 8-Position Mounting & Audio Chimes

The presenter HUD is mounted on the global viewport above all slides (`z-index: 60`), never inside the scaled stage:

1. **HUD Geometry & 8-Position System (`31-slide-controller-buttons.md`):**
   - Positions: Any of 8 edge anchors (`TopLeft`, `TopCenter`, `TopRight`, `BottomLeft`, `BottomCenter`, `BottomRight`, `LeftCenter`, `RightCenter`).
   - Default: `BottomCenter` (or `TopRight`), height `56px`, radius `9999px`, backdrop blur `12px`, border hairline.
   - Safe areas: Includes `env(safe-area-inset-*)`. Tooltip and expand directions derive from edge.
2. **Hover-Reveal & Auto-Hide:**
   - 2.5s idle fade to `opacity: 0.15` (or `0`). Any mouse movement restores `opacity: 1.0` over `150ms`.
3. **Action Clusters (Divided by `1×24px` borders):**
   - **Navigation Group:** Circular 40×40px `ChevronLeft` + Counter + Circular 40×40px `ChevronRight`. Counter format: `{current} / {total}` in Poppins 500, tabular figures.
   - **Share, Theme & Media Group:** 40×40px `Share2` (deep links to `#slide-{current}`, clipboard copy + toast) + 40×40px `Palette` (theme popover with live WCAG contrast readout) + 40×40px `Music` (ambient audio loop toggle).
   - **Presenter & Tools Group:** 40×40px `Video` (cycles webcam PIP presets: Hidden, 120px, 180px, 240px) + 40×40px `Pencil` (builder mode toggle) + 40×40px `Maximize2` (fullscreen toggle).
4. **Audio Chimes & Sound Synthesis:** Uses native Web Audio API oscillators for navigation ticks (800Hz decay) and jump chords (523Hz+659Hz dual chime) without asset loading.
5. **Top Progress Bar & Bottom Dots:** Fixed `top: 0` 4px gradient progress bar + fixed `bottom-8` dot pagination expanding active dot to `28×8px`.
6. **The 3-Axis Layout Rule (`28-slide-layouts.md`):**
   - Axis A: Vertical placement via `justify-content` on outer `<section>`.
   - Axis B: Internal sibling spacing via `mb-*`/`mt-*`/`gap-*` on children.
   - Axis C: Horizontal brand alignment via inline `style={{ paddingLeft: 'var(--brand-inset-x)', paddingRight: 'var(--brand-inset-x)' }}`.
7. **CSS Grid Track Packing:** Container `:has(.slide-card.is-compact)` triggers `grid-auto-rows: min-content`, `align-content: start`, `row-gap: 1rem`; compact cards use `align-self: start`.
8. **Dual-Screen Presenter Console (`/present`):** Dedicated `1280×800` window with speaker notes, elapsed timer, pacing clock, and next-slide preview, synchronized via `BroadcastChannel("deck-sync")`.

---

## 5. 10 Production Themes & 10-Step Precision Gradients

1. **10 Production Themes (`27-slide-canvas-and-themes.md`):**
   - `bright-gold` *(Default Dark)*: Vivid Gold `#F3A502`, Cream `#FFF1D6`, Obsidian `#0D0D0D`.
   - `noir-gold`: Classic Gold `#C9A84C`, Warm Cream `#F0D78C`, Obsidian `#0D0D0D`.
   - `vscode-dark`: Azure Blue `#007ACC`, Light Gray `#D4D4D4`, Dark Slate `#1E1E1E`.
   - `dracula`: Electric Purple `#BD93F9`, Pink Glow `#FF79C6`, Charcoal `#282A36`.
   - `monokai`: Neon Green `#A6E22E`, Orange `#FD971F`, Deep Ink `#272822`.
   - `github-light` *(Canonical Light)*: Blue `#0969DA`, Espresso Ink `#24292F`, White `#FFFFFF`.
   - `paper-ink` *(Print-Friendly Light)*: Deep Amber `#8A5A0E`, Espresso Ink `#1F1A12`, Warm Paper `#FAF6EC`.
   - `macos-sonoma`: Cupertino Blue `#007AFF`, Soft White `#F5F5F7`, Dark Glass `#1E1E24`.
   - `windows-11`: Fluent Cyan `#60CDFF`, Crisp White `#FFFFFF`, Mica Dark `#202020`.
   - `navy-blue`: High-Tech Cyan `#06B6D4`, Orange Ember `#F59E0B`, Deep Navy `#1A2840`.
2. **Light Theme Contract:** Inverts text into Espresso Ink (`#1F1A12` / `#24292F`), darkens amber accents for WCAG AA compliance, and preserves dark HUD chrome.
3. **Character-by-Character Shading Engine:** Hero names use per-glyph color stepping ($C_0$: accent $S_4$, $C_1$: warm intermediate $S_6$, $C_2..C_n$: primary ink $S_0$).
4. **7 Semantic Pill Presets:** `yellow`, `white`, `purple`, `green`, `blue`, `pink`, `red` with ITU-R BT.709 contrast math.

---

## 6. The 20 Master Slide Layouts

Use `02-spec/07-design-system/34-slide-layout-catalog.md` for coordinate geometries and slot specifications. When building presentation decks, select from the 20 master layout models:

1. **`title`:** 78px headline, category pill kicker, subtitle, presenter bio card (avatar 64×64), bottom organic SVG dual wave ribbon.
2. **`executive-persona`:** Asymmetric portrait staging, halftone dot matrix, character gradient hero name (104px), profile card, location tag.
3. **`key-player`:** 3–4 member grid, 380×380px portraits, roles, bios, social links.
4. **`before-after`:** High-contrast split cards: Left ("Before" negative, Rose-200 border, pain points) vs Right ("After" positive, Violet/Emerald border, proof metrics, image wipe).
5. **`usp-strike`:** 124px statement with 6px editorial strikethrough rejecting industry malpractice, paired with horizontal 3-point proof cluster.
6. **`pricing`:** 3-tier card grid, featured "Hot" tier with `scale: 1.03` and gradient border, price figures 48px, full-width CTA buttons.
7. **`steps-chain`:** 4-phase horizontal roadmap, numbered 48×48px step badges, connecting horizon progress line, duration pills, deliverable lists.
8. **`testimonials`:** Dual quote cards in 26px italic Poppins, author avatars, bottom partner logo bar with grayscale hover reveal.
9. **`talent-funnel`:** 4 progressively narrowing capability bands (1640px down to 860px).
10. **`bullets`:** Ground-truth 3-point bullet cards with 48×48px icon containers paired with right-side photographic plate.
11. **`center`:** Centered 88px headline + 28px subhead for keynote thesis statements, module transitions, and core metrics.
12. **`left`:** Asymmetric editorial narrative column (64px heading + 22px body) + right-side code block or media plate.
13. **`quote` / `priority`:** Bold 60px quotation headline, 120px quotation mark glyph, and attribution pill bar.
14. **`process` / `timeline`:** Connected milestone flow with animated track line and 3–5 milestone deliverable cards.
15. **`counter-stat`:** Giant 110px tabular numeric counter with trending delta badge and benchmark methodology deck.
16. **`reveal-grid` / `depth-stack`:** 3D card depth stack (`perspective: 1200px`) and staggered disclosure cards.
17. **`embed`:** Sandboxed interactive iframe container with source URL bar and stage controls.
18. **`poll` / `qa`:** Real-time interactive audience voting bars with scannable QR room code and dynamic percentages.
19. **`tabletop-bike`:** Physical hardware showcase with transparent product cutout, contact ground shadow, circular pulse hot-spots, and telemetry spec grid.
20. **`tech-stack`:** 4-column categorized architecture capability matrix (Frontend, Backend, Data, AI/Pipelines).

---

## 7. Slide Builder Mode & Live Canvas Inspector

1. **Dual-Store Separation:**
   - `useDeckStore` (persisted to localStorage `deck-v1`) holds the authoritative presentation JSON tree.
   - `useEditStore` (ephemeral) manages `selectedElementId`, bounding box coordinates, and undo/redo stacks.
2. **7 Canvas Stacking Layers:** Layer 0 (Base) -> Layer 1 (Watermarks) -> Layer 2 (Media) -> Layer 3 (DOM Typography) -> Layer 4 (Ink) -> Layer 5 (Selection Overlays) -> Layer 6 (Inspector HUD).
3. **`convertSlideType` Engine:** Seamlessly converts between all 20 slide layouts while preserving primary text, secondary text, and media.
4. **Hotkeys:** `B`/`E` toggle builder mode, `Tab` cycle elements, `Cmd/Ctrl+Z` undo, `Cmd/Ctrl+Shift+Z` redo, `Escape` deselect, `1`–`0` quick-switch theme palettes, `F` fullscreen, `G` grid, `J` jumper.

---

## 8. Blind-AI Execution Checklist

Before outputting code or completing your task, verify every item:

- [ ] Canvas is strictly `1920×1080` with uniform `min(w/1920, h/1080)` vector scaling.
- [ ] ZERO text is flattened into raster images; all text is live DOM elements.
- [ ] Controller HUD is anchored to fixed `top-8 right-8` on the viewport with 3 divided groups.
- [ ] Counter displays `{current} / {total}` with tabular figures.
- [ ] Dot pagination expands the active dot to `28×8px` and inactive to `8×8px`.
- [ ] Slide layout strictly matches one of the 20 master models in `34-slide-layout-catalog.md`.
- [ ] Theme binds to one of the 10 production themes in `27-slide-canvas-and-themes.md` (default `'bright-gold'`).
- [ ] Gradients adhere strictly to the 10-step tables ($S_0$ through $S_9$) in `32-slide-color-options.md`.
- [ ] Character-by-character color shading applies to hero names and title keywords.
- [ ] Builder mode maintains dual-store separation (`useDeckStore` vs `useEditStore`) and 7 visual layers.
- [ ] Presenter console syncs via BroadcastChannel with speaker notes and pacing timers.
- [ ] All relative paths referenced start from the git repository root.
- [ ] NO private or company names are present in public files or examples.
