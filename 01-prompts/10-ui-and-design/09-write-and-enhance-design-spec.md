# Write and Enhance a Design Spec That a Blind AI Can Follow

> **Prompt Version:** 2.0.0
> **Trigger keywords:** `design-spec`, `enhance-design-spec`, `ui-spec`, `slide-spec`, `blind-spec`, `spec-creator`, `spec-enhancer`

**/goal** Autonomously ingest source designs (codebases, specifications, presentation decks, UI components, screenshots) and author or enhance public design-system specifications so comprehensive, mathematically exact, and rigorous that any "blind AI" or human engineer with zero external context can follow them blindly to design websites, presentations, blogs, animations, menus, buttons, builder modes, or images without guessing a single value or hallucinating.

**/learn** A specification is finished ONLY when every component, token, physics parameter, layout slot, and interaction timing in the source has its exact values written down in structured tables, every gap is audited, and the completion gate at the end passes 100%. A file that merely points at another file or leaves vague placeholders does not count as complete.

---

## 1. Required Execution Inputs

Before beginning analysis or writing, establish these parameters:

| Input Parameter | Definition & Constraints | Default Value |
|:---|:---|:---|
| `SOURCE` | Read-only folder paths or files containing reference implementations (code, specs, decks). | User-provided or repository inspection |
| `TARGET` | Spec directory where public design system files are created or enhanced. | `02-spec/07-design-system/` |
| `MODE` | `create` (authoring new specifications) or `enhance` (fixing and expanding existing specs). | `enhance` |
| `DOMAINS` | Any of: `website`, `blog`, `menu`, `button`, `color`, `type`, `spacing`, `motion`, `sections`, `slide-deck`, `slide-builder`, `site-builder`, `image`. | All domains present in source |
| `BANNED_NAMES` | Client, product, company, private repository, and individual person names that MUST NEVER appear in public specifications. | Mandatory anonymization |

---

## 2. Hard Architectural Rules (Zero-Tolerance Violations)

1. **Total Ban on Invented Values:** Every hex, HSL, RGB, OKLCH, pixel, millisecond, cubic-bezier easing, character stagger rate, keycode, and default limit MUST originate from `SOURCE` or grounded mathematics. If the source is silent, explicitly state: `Not specified in source. Do not invent.` Never write "about", "around", "roughly", or "~" in place of a concrete number.
2. **Strict Anonymization Mandate:** Never leak names from `BANNED_NAMES` (no client company names, internal repositories, or personal names) into public specifications, comments, sample text, or image paths. Always use generic enterprise designations (e.g., `{BRAND}`, `{COMPANY}`, `Executive Leadership`, `Enterprise Platform`).
3. **Pure DOM Text Rendering Mandate (Zero Baked-In Text):** In websites, presentations, and UI designs, all headlines, body copy, bullets, metrics, credentials, and captions MUST be rendered as live, selectable DOM HTML elements (`<h1>`, `<p>`, `<span>`) styled with CSS typography tokens. Raster images (`.png`, `.jpg`, `.webp`) are strictly reserved for photographic visual plates, author avatars, and partner logos.
4. **Strict Relative Git Paths Only:** All file paths, markdown links, subtask paths, and citations MUST be strictly relative paths starting from the repository root (e.g., `02-spec/07-design-system/11-button-system.md`). TOTAL BAN on absolute filesystem paths and file-scheme URIs.
5. **Lowercase File Naming Convention:** All generated specification files MUST use strictly lowercase kebab-case naming (e.g., `34-slide-layout-catalog.md`). No uppercase characters, spaces, or underscores in filenames.
6. **Strict Boolean Standards:** Positive booleans MUST ALWAYS be evaluated implicitly (`if isReady`). NEVER evaluate explicitly against true (`if isReady == true` is banned). NEVER combine positive and negative checks in the same condition.
7. **Scoped Domain Palettes:** Marketing websites, presentation decks, and branding graphics each maintain distinct, dedicated color palettes. Never contaminate a marketing page with slide presentation amber or vice-versa.
8. **Closed Sets Rule:** Every selectable options list (button variants, sizes, slide layout types, theme IDs, pill preset colors, transition kinds, hotkeys, builder modes) MUST be specified as a closed set followed by: `Anything else does not exist.`
9. **Values Beside the Rule:** Never simply write "use the primary color". Document the semantic token, HEX value, RGB coordinates, and HSL values in the same row.
10. **Strict File Size Cap:** Specification files MUST remain bounded (at most 300 lines per file). Complex topics must be broken into discrete, single-responsibility files.

---

## 3. Systematic 5-Phase Workflow

### Phase 1: Comprehensive Source Inventory (Read-Only)
1. Recursively scan `SOURCE` directories, skipping `node_modules`, build artifacts, `.git`, and lockfiles.
2. Read token files, stylesheets, component primitives, hooks, route templates, and existing specs.
3. Construct an exhaustive inventory table recording:
   - **Tokens:** Colors, 10-step gradients, shadows, radii, spacing, gutters, breakpoints, durations, easings.
   - **Typography:** Font families, weights, fluid clamp scales, line heights, tracking, tabular numeral rules.
   - **Components:** Button variants, sizes, sticky glass header, nav links, mega panel, 3D flip card, mobile drawer, footer.
   - **Sections:** Hero layouts, feature grids, drag rails, process stacks, pricing tables, FAQ accordions.
   - **Motion:** Keyframe definitions, reveal thresholds, stagger intervals, reduced-motion fallbacks.
   - **Slides:** 16:9 canvas (`1920×1080`), scale math, 10-step ramps ($S_0$–$S_9$), 10 master layouts, controller HUD, dots, hotkeys, step reveals.
   - **Builders:** Slide builder (`useDeckStore` and `useEditStore`) and the website builder in `26-visual-builder.md` only. `36-website-content-builder-mode.md` is a pointer.
   - **Content extracts:** A folder of page-copy changes is content only. Do not take colors, type, or motion from it.
   - **Images:** Canvas dimensions, safe zones, text caps, aspect ratios.

### Phase 2: Gap & Fidelity Audit
Compare every inventory item against existing specifications in `TARGET`:
- `exact`: Value is explicitly documented and matches source ground truth.
- `wrong`: Value conflicts with source implementation.
- `pointer-only`: Specification merely links to another file without providing the concrete numbers.
- `missing`: Component, token, or behavior is absent from specifications.

### Phase 3: Surgical Specification Authoring & Enhancement
1. Remediate all `wrong` values first, aligning them with ground truth.
2. Author dedicated specification files for all `missing` and `pointer-only` areas.
3. Follow the standardized Spec File Template (Section 5) for every file.

### Phase 4: Bi-Directional Cross-Referencing & Prompt Sync
1. Add every new specification file to the reading checklist in `02-spec/07-design-system/readme.md`.
2. Add every new file to the module inventory table in `readme.md`.
3. Update consuming prompts (`01-prompts/10-ui-and-design/07-follow-ui-ux-design-system.md` and `08-create-slide-deck.md`) so their reading sequences reference all new files.

### Phase 5: Verification & Completion Gate
Execute the automated lint checks and review against the Completion Gate (Section 7).

---

## 4. Domain Field Checklist

This section names the fields to fill. It does not supply the values. The source wins. If the source is silent, write `Not specified in source. Do not invent.` If the slide source and the website source disagree, write two palettes and name the domain of each. Do not merge them. Do not copy a number from this prompt into a spec.

### 4.1 Colors & 10-Step Gradient Precision System
- Every semantic token name, HEX code, RGB coordinate, and HSL value.
- Mathematical 10-step gradient tables ($S_0$ through $S_9$) for each of the 7 flagship themes:
  - Theme 1: Royal Violet (Sample White Presentation System)
  - Theme 2: Midnight Corporate Gold (Executive Keynote & Global PPT)
  - Theme 3: Enterprise Slate Blue (SaaS & Engineering)
  - Theme 4: Emerald Clinical Green (Healthcare & Evidence Metrics)
  - Theme 5: Sunset Crimson Orange (Sports, High-Energy Bike & Tabletop, Dynamic Growth)
  - Theme 6: Minimalist Monochrome Paper & Slate (Academic & Architectural Keynotes)
  - Theme 7: Cyan-Indigo Tech Gradient (Developer Infrastructure & AI Platforms)
- Per-slide gradient background editor (`slide.gradient`): type (`linear` | `radial`), angle (`0-360deg`), 2–4 color stops (`{ color: string, at: number }`), with toggle switch and live CSS preview.
- Character-by-character color stepping engine (leading glyph accent, intermediate glyph tint, terminal glyphs primary ink).
- 7 Semantic pill presets (`yellow`, `white`, `purple`, `green`, `blue`, `pink`, `red`) with semantic mappings.
- Relative luminance formula for custom pill ink calculation ($0.299R + 0.587G + 0.114B > 0.6 \implies$ `#0b0b12`, else `#ffffff`).

### 4.2 Type
- Family, weight, size, line height, and tracking, each only when the source states it.
- A fluid `clamp()` scale only when the source states the formula.

### 4.3 Button System
- 6 Standard variants: `primary` (shine-sweep on gradient accent), `solid` (shine-sweep on brand primary), `outline` (pointer-fill on hover), `glass` (gradient-ring, frosted acrylic), `ghost`, `link`.
- `ShineButton` (High-Conversion Header CTA): Near-black pill (`bg-night`, `border-white/10`) with continuous 45-degree linear sheen sweep (`@keyframes shine`) and orbiting 8-particle spark constellation (`SPARKS` array, `@keyframes spark-float`). Inverse variant: white pill on dark imagery.
- 4 Sizing scales: `sm` (36px), `md` (44px), `lg` (52px), `icon` (44×44px).
- Magnetic cursor pull physics (`strength: 0.22`).
- Hardware-accelerated CSS3 interactions: `.shine-sweep` keyframes, `.pointer-fill` radial fills, `.gradient-ring` border masks.
- Presentation controller buttons: 40×40px round action triggers, counter with tabular numbers, deep link share, fullscreen toggle, webcam PIP toggle.
- Capsule pill buttons: Tokenized `--capsule-gold`, `--capsule-ember`, and near-black gradient text (`0 0% 4%`) for WCAG 2.2 AA contrast.

### 4.4 Navigation
- Header height, scroll threshold, panel motion, and drawer behavior, each only when the source states it.

### 4.5 Slide Presentation Engine & 20 Master Layouts
- Pure DOM text mandate: Absolute ban on baked-in raster text.
- Virtual canvas: `1920 × 1080` (16:9), vector scale $\min(\text{w}/1920, \text{h}/1080)$, `transform-origin: center center`.
- Floating Controller HUD (`31-slide-controller-buttons.md`): 8-position mounting system (`ControllerPosition`), 2.5s idle fade with hover-reveal, Web Audio API sound synthesis (ticks and chords), background music toggle, first-run onboarding popup, top 4px progress bar, bottom dot row (active `28×8px`, inactive `8×8px`).
- The 3-Axis Layout Rule (`28-slide-layouts.md`): Axis A (vertical placement via `justify-content` on outer `<section>`), Axis B (internal sibling spacing via `mb-*`/`mt-*`/`gap-*` on children), Axis C (horizontal brand alignment via inline `style={{ paddingLeft: 'var(--brand-inset-x)', paddingRight: 'var(--brand-inset-x)' }}`).
- CSS Grid Track Packing: Container `:has(.slide-card.is-compact)` triggers `grid-auto-rows: min-content`, `align-content: start`, `row-gap: 1rem`; compact cards use `align-self: start`.
- 20 Master Layout Models:
  1. `title`: 78px headline, category pill, subtitle, presenter bio card (64×64 avatar), bottom organic SVG dual wave ribbon.
  2. `executive-persona`: Asymmetric portrait staging, halftone matrix, character-shaded hero name (104px), profile card, location tag.
  3. `key-player`: 3–4 member grid, 380×380 portraits, roles, bios, social badges.
  4. `before-after`: Rose-200 negative card vs Violet/Emerald positive card, pain points vs metrics, image wipe.
  5. `usp-strike`: 124px statement with 6px strikethrough rejecting industry practice, paired with 3-point proof cluster.
  6. `pricing`: 3-tier card grid, featured Hot tier with `scale: 1.03` and gradient border, price figures 48px, full-width CTA.
  7. `steps-chain`: 4-phase horizontal roadmap, numbered 48×48px step badges, progress horizon line, duration pills.
  8. `testimonials`: Dual quote cards in 26px italic Poppins, author avatars, bottom partner logo bar with grayscale hover.
  9. `talent-funnel`: 4 progressively narrowing capability bands (1640px down to 860px).
  10. `bullets`: Ground-truth 3 bullet cards with 48×48px icon containers paired with right-side photographic plate.
  11. `center`: Centered headline (88px) + subhead (28px) for punchy transitions and thesis statements.
  12. `left`: Asymmetric 2-column editorial narrative column paired with right-aligned code block or plate.
  13. `quote`: High-authority one-liner quote with 120px decorative quotation glyph and attribution pill.
  14. `process`: Connected phased milestone flow with animated progress line and milestone deliverables.
  15. `counter-stat`: Giant 110px tabular numeric counter with trending delta indicator and context deck.
  16. `reveal-grid`: 3D card depth stack (`perspective: 1200px`) and staggered disclosure cards.
  17. `embed`: Live sandboxed interactive prototype container with status bar and stage controls.
  18. `poll`: Real-time interactive audience voting bars with scannable QR room code and dynamic percentages.
  19. `tabletop-bike`: Physical hardware showcase with transparent product cutout, contact ground shadow, circular pulse hot-spots, and telemetry spec grid.
  20. `tech-stack`: 4-column categorized architecture capability matrix (Frontend, Backend, Data, AI/Pipelines).

### 4.6 Slide Builder Mode & Canvas Inspector
- Dual-store separation: `useDeckStore` (persisted to localStorage `deck-v1`) vs `useEditStore` (ephemeral editor state: selected element, undo/redo stacks).
- 7 Visual layers: Base -> Watermarks -> Media -> DOM Typography -> Ink Annotations -> Selection Overlays -> Inspector HUD.
- Floating Draggable Builder Panel (`BuilderPanel`): Width `264px`, drag physics with boundary clamping, minimizable.
- `convertSlideType` engine: Seamlessly converts between all 20 slide types while preserving primary and secondary text.
- Builder Sub-Panels: `SlideOptionsPanel`, `FontPanel`, `GradientPanel`, `InsertIconPanel`, `InsertImagePanel`, `StepCameraPanel` & `XYPad`, and `HistoryPanel`.
- Hotkeys: `B`/`E` toggle builder mode, `Tab` cycle elements, `Cmd/Ctrl+Z` undo, `Cmd/Ctrl+Shift+Z` redo, `Escape` deselect, `1`–`7` theme switch, `F` fullscreen.
- Bounding box overrides: `boxes: Record<string, EditBox>`.
- Multi-format exports: Headless Chromium print-ready PDF export (`slides.print.tsx`), 3-up executive handout (`slides.handout-3up.tsx`), and JSON deck manifest.

### 4.7 Website Content Builder Mode
- Follow `02-spec/07-design-system/26-visual-builder.md`:
- Zero-server, client-side visual review-only overlay, never auto-publishes to live, produces machine-applicable change lists for developers.
- Gate: Opens via URL query `?builder=1&email={OWNER_EMAIL}` with parameter preservation across client-side SPA navigation.
- Modes: Text, Images, Menu, Layout, Off.
- Inline editing: Double-click sets contentEditable (`plaintext-only`), commits on Ctrl+Enter / Cmd+Enter, restores on Escape.
- Allowed tags strictly: `STRONG`, `B`, `EM`, `I`, `BR`, `A`, `SPAN` (only with `class="gradient-text"`).
- Deterministic export: Downloadable ZIP archive containing `content-changes--all-pages--YYYY-MM-DD-HHmm.zip` with `SUMMARY.md`, page changes markdown, menu changes markdown, and assets.

### 4.8 Homepage & Blog Sections Catalog (`41-homepage-and-blog-sections.md`)
- Complete 15 Homepage Section Patterns:
  1. Social Proof Hero with Split Capability Console & dynamic toggle row
  2. Dual Marquee Logo Rail with reverse direction & pause on hover
  3. Feature Block Grid with cursor-tracking `.spotlight` radial highlight
  4. Tabbed Product Panel with sliding pill indicator (`layoutId="activeTab"`)
  5. Integration Constellation with animated SVG bezier connector paths
  6. Metric Counter Band with CountUp animation, tabular figures, and stat labels
  7. Value Card Deck with pointer-tracking 3D tilt (`TiltCard`, max 8deg)
  8. Industry Hover Grid with expandable accordion / sequential capability cards
  9. Testimonial Drag Rail with grab-cursor physics, progress bar, and card snap
  10. Sticky Process Stack with scroll-pinned sequential steps and drawing connector lines
  11. Pricing Plan Cards with monthly/annual toggle and featured tier elevation
  12. FAQ Split Accordion with zero-JS fluid CSS grid accordion (`grid-template-rows: 0fr -> 1fr`)
  13. Oversized CTA Footer with floating card render (`cta-float`) and shine sweep
  14. Motion Primitives Catalog (Fancy style: slide-swap, pointer-fill, spotlight, shine-sweep, note-connector)
  15. Scroll Stack Cards with sticky stacking deck, progressive scaling, and blur fade
- Blog Editorial System Specifications:
  - Blog Index (`/blog`): Hero search, category filter pills, featured article card with badge, 3-column article cards, pagination.
  - Blog Post Detail (`/blog/:slug`): Editorial hero, author credential pill, calculated reading time (`words / 200`), sticky Table of Contents sidebar with scroll-spy, prose typography, code blocks with syntax highlighting and copy button, callout blocks, author bio box, related articles.

---

## 5. Standardized Spec File Template

Every newly created or enhanced specification file MUST follow this structure:

```markdown
# NN — [Topic Title]

> **/goal** One sentence defining what an AI or human engineer can build using this specification.
> **/learn** One sentence summarizing key tokens, formulas, and references to sibling specifications.

**Version:** 4.3.0
**Status:** Active
**AI Confidence:** High only for numbers this file states. Unstated numbers are not invented.
**Ambiguity:** State it. Do not write None while another file disagrees.

---

## 1. System Overview & Scope
[When to use this specification, architectural pillars, and anti-hallucination boundaries]

---

## 2. Token Registry & Exact Mathematics
[Tables of tokens with Hex, HSL, RGB, OKLCH, dimensions, and mathematical formulas]

---

## 3. Component Taxonomy & Architectural Blueprint
[Component anatomy, geometric dimensions, states, motion curves, and TypeScript interfaces]

---

## 4. Closed Sets & Valid Options
[The complete enumerated list of allowed IDs, variants, or modes, followed by:]
*Anything else does not exist.*

---

## 5. Reference Implementation
[Tested, framework-agnostic code snippet or component blueprint]

---

## 6. Anti-Hallucination & Quality Verification Checklist
[Exhaustive checkboxes guaranteeing that every required number, property, and rule is satisfied]
```

---

## 6. Automated Verification Commands

Run these checks from the repository root:

```powershell
# 1. Verify no absolute paths or file URI protocols
Get-ChildItem -Path "02-spec/07-design-system", "01-prompts/10-ui-and-design" -Filter "*.md" | Select-String -Pattern "([A-Za-z]:[\\/]|file:[\\/]{3}|/home/|/Users/)"

# 2. Verify no vague numbers (about, roughly, approx)
Get-ChildItem -Path "02-spec/07-design-system", "01-prompts/10-ui-and-design" -Filter "*.md" | Select-String -Pattern "\b(about|around|roughly|approx\.?)\s+\d|~\s?\d"

# 3. Check line counts (must be <= 300 lines per file)
Get-ChildItem -Path "02-spec/07-design-system/*.md" | ForEach-Object { "{0,4} lines: {1}" -f (Get-Content $_.FullName).Count, $_.Name }
```

---

## 7. Completion Gate

An AI agent MUST NOT report completion until every gate is verified:

- [ ] Every inventory row from the source has been mapped to an exact specification with concrete numbers.
- [ ] NO `pointer-only` stubs remain; all referenced values are fully declared.
- [ ] Zero baked-in text: Pure DOM typography mandate is enforced across websites and slide decks.
- [ ] Navigation, buttons, and slides contain only numbers that appear in the source. A missing field says `Not specified in source. Do not invent.`
- [ ] Slide types match renderer `case` labels. A type with no component is not given sample pixels.
- [ ] The website builder matches `26-visual-builder.md`. It does not add `data-edit-id` or `?builder=true`.
- [ ] Page-copy extracts contributed no color, type, or motion value.
- [ ] No banned or private company names exist in public specification files.
- [ ] All file paths are strictly relative paths from the git root.
- [ ] All new files are indexed in `02-spec/07-design-system/readme.md`.
- [ ] Consuming prompts (`07-follow-ui-ux-design-system.md` and `08-create-slide-deck.md`) list all new specification files in their reading sequences.
