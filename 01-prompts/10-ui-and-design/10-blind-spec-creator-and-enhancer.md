# Blind AI Specification Creator & Enhancer Protocol

> **Prompt Version:** 3.0.0
> **Trigger keywords:** `blind-spec-creator`, `blind-ai-instruction`, `create-spec`, `enhance-spec`, `design-system-spec`, `ui-ux-spec`, `slide-system-spec`

**/goal** Autonomously ingest source codebases, presentation decks, UI components, and design assets to author or enhance public design-system specifications so comprehensive, mathematically exact, and rigorous that any "blind AI" or human engineer with zero prior context can follow them blindly to design websites, presentations, blogs, animations, menus, buttons, builder modes, or images without guessing a single value or hallucinating.

**/learn** A specification is complete ONLY when every token, component anatomy, motion curve, layout slot, and builder interaction is grounded in concrete numbers (hex, OKLCH, px, ms, cubic-bezier), structured in data tables, and validated against the 100% completion gate. Pointer-only stubs and vague approximations are strictly banned.

---

## 1. Execution Parameters & Anonymization Contract

Before analyzing or authoring, bind these execution parameters:

| Parameter | Type | Definition & Operational Boundary |
|:---|:---|:---|
| `SOURCE` | Path(s) | Read-only directory or workspace containing reference code, decks, or specifications (e.g. presentation repositories, design specifications). |
| `TARGET` | Relative Path | Target specification directory within the repository (default: `02-spec/07-design-system/`). |
| `MODE` | Enum | `create` (authoring new specifications) or `enhance` (deep audit and surgical remediation of existing files). |
| `DOMAINS` | Set | Any of: `navigation`, `mega-menu`, `buttons`, `colors-gradients`, `slide-engine`, `slide-layouts`, `slide-builder`, `web-builder`, `sections`, `blog`, `images`. |
| `BANNED_NAMES` | Token List | Strictly banned proprietary, client, or internal entity names. Every occurrence MUST be replaced with generic tokens (`{BRAND}`, `{COMPANY}`, `Executive Platform`). |

---

## 2. Hard Architectural Rules (Zero-Tolerance Violations)

1. **Total Ban on Invented Values & Approximations:** Every dimension, color coordinate, easing bezier, duration, and stagger MUST originate from `SOURCE` or verifiable mathematics. Never write "about", "roughly", "approx", or "~". If the source is silent, write: `Not specified in source. Do not invent.`
2. **Strict Anonymization Mandate:** Never leak proprietary names into public specs, markdown files, comments, or examples. Use generic enterprise designations.
3. **Pure DOM Text Mandate (Zero Baked-In Text):** In websites, presentations, and UI designs, all headlines, body copy, bullets, metrics, credentials, and captions MUST be rendered as live, selectable DOM HTML elements (`<h1>`, `<p>`, `<span>`) styled with CSS tokens. Raster images are strictly reserved for photographic plates, avatars, and partner logos.
4. **Strict Relative Git Paths Only:** All file paths, markdown links, and citations MUST be strictly relative paths starting from the repository root (e.g. `02-spec/07-design-system/11-button-system.md`). Absolute filesystem paths and `file:///` URIs are totally banned.
5. **Strict Lowercase File Naming:** All generated or modified files MUST use strictly lowercase kebab-case naming (e.g. `34-slide-layout-catalog.md`). No uppercase characters, spaces, or underscores.
6. **Boolean Principles:** Implicit positive booleans only (`if isReady`). Never evaluate explicitly against true (`if isReady == true` is banned). Never combine positive and negative checks in the same condition.
7. **Closed Sets Rule:** Every list of options (button variants, sizes, slide layout types, theme IDs, pill preset colors, hotkeys, builder modes) MUST be declared as a closed set followed by: `Anything else does not exist.`
8. **300-Line File Cap:** Every specification file must remain bounded (at most 300 lines per file). Complex architectures must be modularized into single-responsibility files.

---

## 3. Systematic 5-Phase Execution Workflow

### Phase 1: Exhaustive Source Inventory (Read-Only)
1. Recursively scan `SOURCE` directories, ignoring build artifacts, node modules, and VCS internals.
2. Ingest stylesheets, token registries, component primitives, route views, and presentation engines.
3. Populate a structured inventory covering:
   - **Tokens:** Colors (OKLCH, HEX, RGB, HSL), 10-step gradients ($S_0$–$S_9$), shadows, radii, durations, easings.
   - **Navigation & Menu:** Sticky header geometry, scroll threshold, safe region math (`pad = 14px`, 220ms timer), `SlideSwapLabel` keyframes, 3D flip card, mobile drawer scroll intercept.
   - **Buttons:** 6 variants, 4 sizing scales, magnetic spring physics (`strength: 0.22`), `.shine-sweep`, `.pointer-fill`, `ShineButton` 8-particle constellation, 40×40 HUD buttons, capsule pills.
   - **Slides Engine:** Canvas 1920×1080, scale math, 8-position HUD, sound synthesizer, 10 themes, 20 master layouts, dual-store Slide Builder (`useDeckStore` vs `useEditStore`), 7 visual canvas layers.
   - **Website Builder:** Review overlay gate (`?builder=1&email={EMAIL}`), `contentEditable` inline editor, gradient text preservation, diff tracking, deterministic ZIP export.
   - **Sections & Blog:** Band rhythm, 15 homepage section patterns, pricing tables, blog index, reading progress, sticky TOC.
   - **Images:** Safe zones, text caps, aspect ratios.

### Phase 2: Gap & Fidelity Audit
Compare every inventory item against existing specifications in `TARGET`:
- `exact`: Ground-truth values are fully documented and mathematically verified.
- `wrong`: Specification contradicts source implementation.
- `pointer-only`: Specification links to another file without providing concrete numbers.
- `missing`: Component, token, or behavior is absent from specifications.

### Phase 3: Surgical Specification Authoring & Enhancement
1. Remediate all `wrong` values first to establish absolute source parity.
2. Replace all `pointer-only` stubs with complete, self-contained data tables and blueprints.
3. Author dedicated specification files for all `missing` domains following the Standardized Spec Template.

### Phase 4: Bi-Directional Cross-Referencing & Prompt Synchronization
1. Register every new or updated file in the reading checklist and inventory of `02-spec/07-design-system/readme.md`.
2. Synchronize consuming prompts (`07-follow-ui-ux-design-system.md` and `08-create-slide-deck.md`) so their reading sequences reference all enhanced specifications via relative paths.

### Phase 5: Verification & Zero-Failure Completion Gate
Execute automated verification commands and certify compliance against the Completion Gate (Section 7).

---

## 4. Domain-by-Domain Strong Points & Checklists

### 4.1 Header & Precision Mega Menu System (`10-header-navigation.md`, `33-mega-menu-components.md`)
- [ ] Header height is strictly `72px` (`h-[72px]`), sticky `top-0 z-50 w-full`, surface `bg-background/90 backdrop-blur-xl`.
- [ ] Scroll trigger: Viewport scroll > `12px` activates `border-b border-border shadow-[var(--shadow-card)]` easing over `420ms`.
- [ ] `SlideSwapLabel`: Per-character vertical slide swap on hover; `stagger: 0.04s` for nav links (`0.018s` for buttons); top glyph slides to `translateY(-110%)`, duplicate enters from `translateY(100%)` to `translateY(0)` over `520ms`.
- [ ] Pointer safe region: `pad = 14px` buffer around header and dropdown bounds; pointer exit triggers `220ms` debounce close. Immediate `0ms` dismiss on `Escape` or scroll.
- [ ] MegaPanel entrance: Mounts `absolute left-0 right-0 top-full z-40 pt-3`, animates `opacity: 0, y: -8, scale: 0.985 -> 1, 0, 1` over `260ms cubic-bezier(0.16, 1, 0.3, 1)`.
- [ ] Link items: Left accent hairline `1px` grows `scaleY(0) -> scaleY(1)` with gradient accent over `420ms`. Trailing arrow transitions from `-translate-x-1 opacity-0` to `translate-x-0 opacity-100`.
- [ ] 3D Flip Promo Card: `min-h-[220px]`, `perspective: 1400px`, flips `rotateY(180deg)` over `820ms`.
- [ ] Sticky-safe mobile drawer: Intercepts `wheel` and `touchmove` on `window` **WITHOUT** setting `overflow: hidden` on `<body>` (preserves sticky header).
- [ ] In-Page Menu Editor Modal: Dialog editor for label, description, and target URL with protocol validation (`/`, `#`, `https://`, `http://`, `mailto:`, `tel:`).

### 4.2 Universal Button System & CSS3 Physics (`11-button-system.md`)
- [ ] 6 Variants: `primary`, `solid`, `outline`, `glass`, `ghost`, `link`. Closed set.
- [ ] 4 Sizing scales: `sm` (36px, `text-[13px]`), `md` (44px, `text-sm`), `lg` (52px, `text-[15px]`), `icon` (44×44px).
- [ ] Magnetic cursor pull physics: `<Magnetic strength={0.22} radius={90}>` using Framer Motion springs (`stiffness: 260, damping: 18, mass: 0.4`).
- [ ] `.shine-sweep`: 45-degree angled sheen on hover/focus using `@keyframes shine-sweep-action 3s linear infinite`.
- [ ] `.pointer-fill`: Radial cursor fill scaling smoothly from bottom (`scaleY(0) -> scaleY(1)` over `240ms`).
- [ ] `ShineButton`: Near-black pill (`bg-night`) with continuous sheen and split outer wrapper hosting an 8-particle constellation (`hidden md:block`, deterministic offsets, delay `0s` to `2.8s`, hover acceleration `8s -> 4s`).
- [ ] Presentation controller action buttons: 40×40px circular hit targets, counter with tabular numbers.
- [ ] Capsule gradient buttons: Tokenized `--capsule-gold`, `--capsule-ember`, and near-black gradient text (`0 0% 4%`) for WCAG 2.2 AA contrast.

### 4.3 Colors, Brand Gradients & 10-Step Precision System (`03-theme-variable-architecture.md`, `32-slide-color-options.md`, `27-slide-canvas-and-themes.md`)
- [ ] Complete 3-format color tables (OKLCH, HEX, RGB, HSL) for Blue, Violet, Cyan, Indigo, Teal, Neutral Slate, and Surface Dark.
- [ ] 9 Core brand gradients: `--gradient-hero`, `--gradient-accent`, `--gradient-accent-hover`, `--gradient-brand`, `--gradient-text`, `--gradient-halo`, `--gradient-card-dark`, `--gradient-border`, `--gradient-chrome`.
- [ ] 6 Elevation shadows: `--shadow-xs`, `--shadow-card`, `--shadow-lift`, `--shadow-glow`, `--shadow-dark-card`, `--shadow-3d`.
- [ ] 10-step gradient tables ($S_0$ through $S_9$) documented for each of the 7 flagship themes with per-slide gradient editor schema.
- [ ] 10 Production runtime themes (`bright-gold` default, `noir-gold`, `vscode-dark`, `dracula`, `monokai`, `github-light`, `paper-ink`, `macos-sonoma`, `windows-11`, `navy-blue`) with light-theme contrast inversion rules.
- [ ] Character-by-character color stepping: Leading glyph accent $S_4$, intermediate glyph tint $S_6$, remaining glyphs primary ink $S_0$.
- [ ] 7 Semantic pill presets (`yellow`, `white`, `purple`, `green`, `blue`, `pink`, `red`) with BT.709 luminance contrast formula ($0.299R + 0.587G + 0.114B > 0.6 \implies$ `#0b0b12`, else `#ffffff`).

### 4.4 Slide Presentation Engine & Layout Catalog (`24-slide-presentation-system.md`, `31-slide-controller-buttons.md`, `34-slide-layout-catalog.md`, `42-slide-step-and-sound-system.md`, `43-slide-webcam-overlay.md`, `44-slide-presenter-inspector-and-handouts.md`)
- [ ] Canvas is strictly `1920 × 1080` (16:9) centered with vector scale $\min(\text{viewportWidth}/1920, \text{viewportHeight}/1080)$.
- [ ] Floating controller HUD: 8-position mounting system (`ControllerPosition`), 2.5s idle fade with hover-reveal, Web Audio API sound synthesis (ticks and chords), background music toggle, top 4px progress bar, bottom dot pagination (active `28×8px`, inactive `8×8px`).
- [ ] Step System & Sound Synthesis (`42-slide-step-and-sound-system.md`): `StepTimelineSlide` (chain + active focus row + right description card + gold/ember capsules) vs `AdvanceStepSlide` (camera dolly per frame), plus Web Audio API singleton sound engine (audio asset table, 60ms debounce, procedural synth fallback, gain envelopes, ducking).
- [ ] Presenter Webcam PIP Overlay (`43-slide-webcam-overlay.md`): `PresenterWebcamProvider` state machine, 4 stepped presets (S/M/L/XL), squircle gold rim frame, free pointer drag, autoframe face tracker, and keyboard shortcuts (`C`, `M`, `F`, `+`, `-`, `O`, `H`, `1`).
- [ ] Presenter Inspector & Handouts (`44-slide-presenter-inspector-and-handouts.md`): Standalone speaker route (`/slides/inspector`) with 60% live slide + 40% next preview + speaker notes + elapsed timer, plus 3-up printable handouts (`/slides/handout-3up`) with ruled writing lines and 1-up print mode.
- [ ] 3-Axis Layout Rule: Axis A (vertical placement via `justify-content`), Axis B (internal sibling spacing via `gap-*`/`mb-*`), Axis C (horizontal brand alignment via `--brand-inset-x`).
- [ ] CSS Grid Track Packing: Container `:has(.slide-card.is-compact)` triggers `grid-auto-rows: min-content`, `align-content: start`, `row-gap: 1rem`; compact cards use `align-self: start`.
- [ ] 20 Master Slide Layout Models: Title Hero, Executive Persona, Key Player Bio, Before/After Split, USP Strikethrough, SaaS Pricing, Steps Chain Roadmap, Social Proof, Talent Funnel, 3-Point Master Cards, Center, Left, One-Liner Quote, Process/Timeline, Counter Stat, Reveal Grid/Depth Stack, Embed, Poll/Q&A, Tabletop Hardware & Bike Showcase, Tech Stack Matrix.

### 4.5 Dual Visual Builder Modes (`35-slide-builder-canvas-inspector.md`, `36-website-content-builder-mode.md`)
- [ ] Slide Builder Engine:
  - Decoupled dual-store architecture (`useDeckStore` persisted to `deck.draft.v1` vs `useEditStore` ephemeral).
  - 7 Visual canvas stacking layers (Base, Watermarks, Media, DOM Typography, Ink, Overlays, Inspector HUD).
  - Form field schemas (`SLIDE_TYPE_SCHEMAS`, `FieldKey` union, defaults) ensuring instant renderable preview.
  - Specialized canvas editors: `BoxDiagramCanvasEditor` (node/edge diagram editor), `HotspotCanvasEditor` (pin placement), `GuideMeasurementHUD` (magnetic snapping and pixel distance HUD), `AnimationPreviewPanel` & `ClickRevealToggle` (step preview).
  - `convertSlideType` engine across all 20 layouts preserving primary/secondary narrative copy.
  - Multi-format exports: Headless PDF (`slides.print.tsx`), 3-Up Executive Handout (`slides.handout-3up.tsx`), PPTX native presentation export (`exportPptx.ts`), and Audience Companion screen (`audience.$sessionId.tsx`).
- [ ] Website Content Review Builder:
  - Zero-server, client-side review-only overlay gated strictly by `?builder=1&email={OWNER_EMAIL}`.
  - DOM auto-tagging pipeline (`assignIds()`: Pass 1 menu links, Pass 2 images/logo, Pass 3 animated headings, Pass 4 generic text with block-child splitting and decorative preservation).
  - In-place `contentEditable` with gradient text flattening for caret visibility, floating keyboard helper chip, and selection gradient highlight toolbar.
  - Plain-text paste sanitizer allowlist (`STRONG`, `B`, `EM`, `I`, `BR`, `A`, `<span class="gradient-text">`).
  - Media & menu modals with 2 MB upload cap, mandatory alt-text, and link protocol validation.
  - Automated before/after section screenshot pairs (`shots.ts` via `html-to-image`) with red (`#ef4444`) and slate (`#64748b`) outlines.
  - Deterministic ZIP export (`content-changes--all-pages--*.zip`) containing `SUMMARY.md`, `manifest.json`, per-page Markdown diffs, screenshots, and assets.

### 4.6 Section Architecture, Blog System & Image Rules (`38-card-and-pricing-components.md`, `41-homepage-and-blog-sections.md`, `37-image-specifications.md`)
- [ ] Alternating band rhythm: Sections cycle tones (`light` -> `soft` -> `light` -> `soft`); adjacent sections NEVER share the same tone; dark tone used at most once per page.
- [ ] Word count boundaries: Eyebrow <= 3 words, H2 <= 8 words, lead <= 28 words, card body <= 32 words, at most ONE gradient accent phrase per heading.
- [ ] 15 Core Homepage Patterns: Split Hero, Dual Marquee, Spotlight Grid, Tabbed Product, Integration Constellation, Metric Counters, 3D Tilt Deck, Industry Hover Grid, Testimonial Drag Rail, Sticky Process Stack, Pricing Plan Cards, FAQ Fluid Accordion, Floating CTA Footer, Motion Primitives, Scroll Stack Cards.
- [ ] Blog Editorial System: Blog Index with search/filter pills/reading time; Post Detail with reading time, sticky 280px TOC sidebar with scroll-spy, prose typography, code block syntax highlighting with copy button, callout blocks, author bio.
- [ ] Image specifications: Exact dimensions, safe zones, text rules, and closed image palette.

---

## 5. Standardized Spec File Markdown Template

Every newly created or enhanced specification file MUST strictly adhere to this structure:

```markdown
# NN — [Topic Title]

> **/goal** One sentence defining what an AI or human engineer can build using this specification.
> **/learn** One sentence summarizing key tokens, formulas, and references to sibling specifications.

**Version:** 4.3.0
**Status:** Active
**AI Confidence:** High (only for numbers explicitly stated in this file)
**Ambiguity:** [Explicitly state any ambiguity or write None]

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

Execute these commands from the repository root:

```bash
# 1. Verify strict relative paths across target folders
python 03-ai-scripts/07-relative-path-fixer.py 02-spec/07-design-system
python 03-ai-scripts/07-relative-path-fixer.py 01-prompts/10-ui-and-design

# 2. Verify no vague numbers (about, roughly, approx) using GitMap search
gitmap search "roughly"
gitmap search "approx"

# 3. Check line counts (must be <= 300 lines per file)
python 03-ai-scripts/13-file-size-guard.py --path 02-spec/07-design-system
```

---

## 7. Zero-Failure Completion Gate

An AI agent MUST NOT declare completion until every gate is verified:

- [ ] Every component, token, physics parameter, and layout slot from the source has been mapped to an exact specification with concrete numbers.
- [ ] ZERO `pointer-only` stubs remain; all referenced values are fully declared in structured data tables.
- [ ] Pure DOM text mandate is enforced across websites and slide decks; zero baked-in text in images.
- [ ] Sizing scales, variants, themes, and layouts are defined as closed sets ending with `Anything else does not exist.`
- [ ] All file paths are strictly relative paths from the git repository root.
- [ ] ZERO proprietary, private, or company names exist in public specification files or examples.
- [ ] All generated or modified files use strictly lowercase kebab-case naming.
- [ ] Positive booleans are evaluated implicitly without explicit `== true` checks or mixed polarity conditions.
- [ ] All files remain bounded (<= 300 lines per file).
- [ ] `02-spec/07-design-system/readme.md` indexes all new and updated files.
- [ ] Consuming prompts (`07-follow-ui-ux-design-system.md` and `08-create-slide-deck.md`) list all new specification files in their reading sequences.
- [ ] `01-prompts/10-ui-and-design/readme.md` catalogs this prompt and all other prompts in the suite.
