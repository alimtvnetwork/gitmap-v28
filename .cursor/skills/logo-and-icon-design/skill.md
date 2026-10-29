---
name: logo-and-icon-design
description: >-
  Autonomously design, generate, and organize logo, icon, and branding assets following strict lowercase conventions, transparent PNG specs, responsive SVG rules, and dual static/animated modes for Lovable and web design platforms.
---

# Logo & Icon Design Workflow — Asset & Branding Standards

> **Invocation:** `@[logo-and-icon-design](file;.cursor/skills/logo-and-icon-design)`

> **Skill Version:** 1.3.0
> **Target Environment:** Lovable, Web Design AI Platforms, & Local Design Workflows
> **Synchronization:** Main Meta-Repo & Connected Workspaces

This skill guides the creation, organization, and validation of production-ready logo, icon, and branding assets.

---

## Strictly Avoid (Critical Negative Constraints)

> [!CAUTION]
> **TOTAL BAN ON HTML, WEBSITES, PAGES, AND APPLICATION CODE (AUTO-REJECT)**
> The AI MUST strictly avoid generating any web application, frontend/backend code, or HTML files. This workflow is exclusively for branding assets and icons.

1. **NO HTML Files:** NEVER create or output `.html` files (no `index.html`, `app.html`, `test.html`, etc.).
2. **NO Web Pages or Mockup Apps:** NEVER build landing pages, homepages, dashboard pages, or mock web applications.
3. **NO React / Vue / Svelte Components:** NEVER create `.tsx`, `.jsx`, `.vue`, or `.svelte` components.
4. **NO CSS Stylesheets:** NEVER create external CSS files (`styles.css`, `app.css`). Theme colors must live ONLY in `colors-themes/palette.md`.
5. **NO Server or Script Code:** NEVER create JavaScript/TypeScript scripts, API routes, or backend servers.
6. **NO Tokens.json:** NEVER generate `tokens.json` in `colors-themes/` (use `palette.md` only).
7. **NO Solid / Opaque Backgrounds on ANY PNG (TOTAL BAN):** NEVER generate ANY PNG icon with a solid background (no solid white, no solid black, no grey, no dark square backdrop). ALL PNG files without exception MUST have a 100% transparent alpha-channel background (`alpha = 0`, `#00000000`). "Dark" and "Light" variants refer ONLY to the color/fill of the icon artwork itself—NEVER the background! The canvas/background is ALWAYS 100% transparent. Dark Mode strictly means inverted high-luminance strokes (luminance > 0.70, e.g. #FFFFFF, #FCD34D, #38BDF8) on a 100% transparent canvas—NEVER a dark box container (`<rect>`)! Light Mode uses rich saturated hues (luminance < 0.40).
8. **NO Solid Background Shapes in SVGs:** NEVER add a background `<rect>` or container shape (e.g. `<rect fill="#000" .../>` or `<rect fill="#fff" .../>`). The SVG canvas must always be 100% transparent (`alpha = 0`). Artwork must float freely within the viewBox with an 8–10% quiet zone (inner safe margin) to eliminate edge clipping at small favicon sizes (16px/32px).
9. **NO Messy Overlapping Lines:** NEVER create tangled, chaotic, or densely overlapping strokes. Overlapping lines make icons look muddy, unreadable at small resolutions, and amateurish. Maintain clean geometric separation or continuous connected paths.
10. **NO Uppercase Filenames:** NEVER use uppercase letters in folder names or file names (`README.md`, `Projects/`, `Logo.svg` are BANNED; use `readme.md`, `02-projects/`, `logo.svg`).
11. **NO Base64 Images in SVGs:** NEVER embed raster images or base64 data URLs inside SVG files.

---

## 1. Strict Boundary & Scope

### Input Capture & Validation

Always verify the following inputs before generating assets:
1. `product_name`: Brand/product name.
2. `product_idea`: Purpose, target audience, and brand voice/tone (e.g. minimalist, playful, high-tech).
3. `sample_colors`: Target color palette or mood.
4. `inspiration_logos`: Optional reference brands or styles (e.g. Apple, Microsoft, OpenAI, Sephora, Stripe). Prioritize referenced geometric style and simplicity.
5. `needs_dark_white_variants`: Boolean for dark and white variants (default: `true`).
6. `is_animated`: Boolean for animated assets (default: `false` unless requested).

If any core input is missing, prompt the user for clarification before generating files.

---

## 2. Brand Identity & Icon Design Principles

### A. Conveying the Brand Image

- The icon must serve as an instant visual metaphor for the brand's core purpose.
- Prioritize high memorability: a viewer should easily recognize and recall the mark after a single viewing.
- Guarantee micro-to-macro scalability: the icon must be legible and distinct at 16x16 px (favicon) and 52x52 px (app icon) while remaining balanced at 1024x1024 px.

### B. Anti-Overlapping Rule & Clean Geometry

- **Ban Chaotic Overlaps:** Avoid tangled strokes, cross-hatching, or visual clutter.
- **Connect the Dots:** Connect visual elements purposefully through continuous strokes, modular alignment, or clean geometric joints.
- **Inspirational Archetypes:**
  - *Silhouette & Cutout (Apple Archetype):* Bold solid mark with distinctive negative space.
  - *Modular Grid (Microsoft Archetype):* Clean, non-overlapping geometric tiles separated by consistent gutters.
  - *Rotational Flow (OpenAI Archetype):* Interlocking, symmetrical curved ribbons forming a continuous loop.
  - *Fluid Organic Curve (Sephora Archetype):* Single-stroke elongated "S" curve/flame with zero overlapping lines, conveying beauty, elegance, and fluidity.

### C. Dual-Mode Inverted Artwork Clarity & Luminance Guidelines

A critical anti-pattern in AI icon generation is confusing "Dark Mode" with a dark background box. Dark Mode assets MUST represent the artwork itself inverted for readability against dark host environments, floating freely on a 100% transparent canvas (`alpha = 0`).

- **Dark Mode Requirement (Inverted High-Luminance Strokes):**
  - Strokes and primary fills must utilize high-luminance tones (perceptual luminance $L > 0.70$), such as `#ffffff` ($L = 1.0$), `#fcd34d` ($L \approx 0.82$), `#38bdf8` ($L \approx 0.78$), or `#a7f3d0` ($L \approx 0.88$).
  - **TOTAL BAN ON DARK BOX CONTAINERS:** Dark mode icons must NEVER include a dark rectangular container (`<rect>`, `<circle>`, or `<path>` acting as a dark card backdrop). The background is ALWAYS 100% transparent.
- **Light Mode Requirement (Rich Saturated Hues):**
  - Light mode assets utilize deep, rich, high-contrast saturated hues (perceptual luminance $L < 0.40$), such as deep indigos (`#1e1b4b`), slate navy (`#0f172a`), deep cobalt (`#1d4ed8`), or dark emerald (`#064e3b`).
  - The canvas is ALWAYS 100% transparent.

| Mode | Canvas State | Stroke / Fill Luminance Range | Primary Palette Samples | Banned Elements |
|:---|:---|:---|:---|:---|
| **Dark Mode** | 100% Transparent (`alpha = 0`) | High Luminance ($L > 0.70$) | `#ffffff`, `#fcd34d`, `#38bdf8`, `#cbd5e1` | `<rect>` dark backdrops, dark cards |
| **Light Mode** | 100% Transparent (`alpha = 0`) | Low Luminance ($L < 0.40$) | `#0f172a`, `#1e1b4b`, `#1e3a8a`, `#14532d` | Artificial border boxes, matte backgrounds |

### D. SVG Precision, Quiet Zone & Coordinate Geometry

To ensure sharp, unclipped rendering across high-DPI displays, desktop favicons (16px/32px), and mobile app docks:

- **Quiet Zone (8–10% Inner Padding):**
  - All visual vector paths, strokes, and shapes must remain strictly within the inner 80–84% safe area of the `viewBox`.
  - For a standard `0 0 24 24` viewBox, maintain a minimum 2.00px quiet zone margin on all four boundaries (content bounds: `x: 2.00..22.00`, `y: 2.00..22.00`).
  - For a `0 0 48 48` viewBox, maintain a minimum 4.00px quiet zone margin (content bounds: `x: 4.00..44.00`, `y: 4.00..44.00`).
  - For a `0 0 100 100` viewBox, maintain an 8.00–10.00px quiet zone margin (content bounds: `x: 8.00..92.00`, `y: 8.00..92.00`).
  - This quiet zone prevents outer stroke clipping when rendered as favicons or small UI glyphs.
- **Coordinate Precision (2 Decimal Places):**
  - All path data coordinates (`d="..."`), transforms, and dimensions must be rounded to a maximum of 2 decimal places (e.g., `d="M 12.50 4.25 L 18.75 14.10 ..."`).
  - Eliminates floating-point bloat and guarantees deterministic SVG rendering.
- **Semantic Gradients & Defs:**
  - Complex fills must be grouped inside `<defs>` with unique, descriptive, lowercase IDs (e.g. `id="gradient-brand-accent"`).
  - Gradients must specify explicit coordinate space (`gradientUnits="userSpaceOnUse"` or `objectBoundingBox`) and linear stops with color and opacity.
- **Accessibility & Compliance:**
  - Root element must declare `role="img"`.
  - Provide semantic descriptions via `<title id="icon-title">` and `<desc id="icon-desc">` referenced via `aria-labelledby="icon-title icon-desc"`.
  - Monochrome glyphs must use `fill="currentColor"` or `stroke="currentColor"` to inherit parent CSS styling dynamically.

### E. Design System Token Alignment

Icons must map directly to established design tokens to ensure strict aesthetic coherence across web applications, presentations, and UI dashboards:

- **Presentation Engine Tokens:**
  - Primary Accent: `var(--pres-accent, #38bdf8)`
  - Primary Content / Typography: `var(--pres-text, #f8fafc)`
  - Subdued Accents: `var(--pres-muted, #94a3b8)`
- **Tailwind / Shadcn Tokens:**
  - Primary Stroke / Fill: `var(--primary, #0f172a)`
  - Foreground Contrast: `var(--foreground, #020817)`
  - Muted Foreground: `var(--muted-foreground, #64748b)`
- **Dynamic CSS Variable Binding:**
  When rendering inline SVGs in React or HTML components, color attributes should leverage CSS variable fallbacks (e.g. `stroke="var(--pres-accent, #38bdf8)"`) or leverage `currentColor` alongside CSS utility classes (`text-primary`, `text-accent`).

---

## 3. Directory Hierarchy (Strict Lowercase)

All files and folders must follow strict lowercase naming with zero-padded two-digit sequence numbers:

```
/ (repo root)
├── 01-prompts/
└── 02-projects/
    ├── 01-{project-name}/
    │   ├── readme.md
    │   ├── icons-svg/
    │   │   ├── logo.svg
    │   │   ├── logo-dark.svg
    │   │   └── logo-white.svg
    │   ├── icons-image/
    │   │   ├── logo-052.png
    │   │   ├── logo-128.png
    │   │   ├── logo-256.png
    │   │   ├── logo-512.png
    │   │   ├── logo-1024-light.png
    │   │   └── logo-1024-dark.png
    │   └── colors-themes/
    │       └── palette.md
    └── 02-{project-name}/
(repo root assets)
├── favicon.ico
├── favicon.png
└── gif-animation.gif (generated only when animation is requested)
```

### Hierarchy Rules

- Never overwrite an existing project folder; always increment the sequence number (`01`, `02`, ...).
- Do not create `tokens.json`. Color documentation lives exclusively in `colors-themes/palette.md`.
- `readme.md` must use relative links so GitHub renders all assets inline.

---

## 4. Asset Specifications

### A. Vector Graphics (`icons-svg/`)

- Pure, clean XML/SVG vector code without editor junk or base64 bitmaps.
- Responsive `viewBox` (e.g., `viewBox="0 0 100 100"` or `viewBox="0 0 48 48"`), omitting hardcoded `width` and `height` attributes on `<svg>`.
- Always 100% transparent canvas (`alpha = 0`), never add a background `<rect>` or container shape.
- Strict Quiet Zone: 8–10% inner padding around all boundaries to prevent clipping at small favicon sizes (16px/32px) and app dock sizes.
- Coordinate Precision: Maximum 2 decimal places on all coordinates, paths, and transformations.
- Semantic Defs: Gradients and filters declared inside `<defs>` with unique lowercase IDs.
- Accessibility: `<title>` and `<desc>` linked with `aria-labelledby` and `role="img"`.
- Theming & Tokens: `currentColor` for monochrome icons; direct mapping to design system tokens (`var(--pres-accent, #38bdf8)`, `var(--primary, #0f172a)`).
- Variants (all on transparent canvas):
  - `logo.svg`: Primary brand logo on a 100% transparent canvas.
  - `logo-dark.svg`: Optimized for dark themes using inverted high-luminance strokes/fills ($L > 0.70$) on a **100% transparent canvas** (NO dark background `<rect>`).
  - `logo-white.svg`: Pure white monochrome icon (`#ffffff`, $L = 1.0$) on a **100% transparent canvas**.

### B. Transparent Raster Icons (`icons-image/`)

- **100% Transparency Requirement (Zero Solid Backgrounds):** Every single PNG image in `icons-image/` MUST have a **100% transparent alpha-channel background** (`alpha = 0`). NO solid white, solid black, solid grey, or colored background boxes are permitted.
- **Understanding "Dark" vs. "Light" Variants:**
  - **"Dark" does NOT mean a dark background!** It means the icon artwork itself uses inverted high-luminance strokes and fills ($L > 0.70$, e.g. `#ffffff`, `#fcd34d`, `#38bdf8`) so that when the transparent PNG is placed onto a dark website theme, the icon is clearly visible. The canvas/background is 100% transparent.
  - **"Light" means the icon artwork uses rich, dark, or saturated brand colors** ($L < 0.40$, e.g. `#0f172a`, `#1e1b4b`, `#1d4ed8`) intended to be placed on a light website theme. The canvas/background is 100% transparent.
  - Both `logo-1024-light.png` and `logo-1024-dark.png` are **100% transparent PNGs** with zero background fills.
- Standard sizes:
  - `logo-052.png` — 52x52 px 100% transparent PNG icon
  - `logo-128.png` — 128x128 px 100% transparent PNG icon
  - `logo-256.png` — 256x256 px 100% transparent PNG icon
  - `logo-512.png` — 512x512 px 100% transparent PNG icon
  - `logo-1024-light.png` — 1024x1024 px 100% transparent PNG with dark/saturated artwork ($L < 0.40$) for light themes
  - `logo-1024-dark.png` — 1024x1024 px 100% transparent PNG with inverted/high-luminance artwork ($L > 0.70$) for dark themes

### C. Color Palette (`colors-themes/palette.md`)

- Document HEX, RGB, and HSL values with visual markdown swatch previews.

### D. Favicons (Repo Root)

- `favicon.ico` and `favicon.png` placed at the repo root and regenerated per project.

---

## 5. Execution Modes

### Mode 1: Static Logo & Branding Generation (Default)

Generates the complete static vector set, transparent PNG sizes (with inverted colors for dark theme), color palette, root favicons, and project `readme.md`.

### Mode 2: Animated Logo & Branding Generation

In addition to static assets, generates:
- Animated SVG (`icons-svg/logo-animated.svg`) on a transparent canvas using smooth CSS `@keyframes` (pulse, rotation, draw effect).
- Seamless looping GIF (`gif-animation.gif`) placed at repo root.
- Animated asset preview embedded in `readme.md`.

---

## 6. Quality Checklist

- [ ] All file and directory names are strictly lowercase.
- [ ] No HTML pages, React components, or website templates were generated.
- [ ] ALL PNG icons have 100% transparent backgrounds (`alpha = 0`) with ZERO solid background fills.
- [ ] Dark theme variants use inverted high-luminance strokes/fills ($L > 0.70$) on a 100% transparent background (NEVER a black/dark `<rect>` box).
- [ ] Light theme variants use rich saturated hues ($L < 0.40$) on a 100% transparent background.
- [ ] SVG files maintain an 8–10% quiet zone safe margin to prevent edge clipping at 16px/32px favicon sizes.
- [ ] SVG path coordinates, transforms, and dimensions use a maximum of 2 decimal places.
- [ ] SVG files have responsive `viewBox`, no hardcoded `width`/`height`, and no background `<rect>`.
- [ ] Design system tokens (`--pres-accent`, `--pres-text`, `--primary`, `--foreground`) or `currentColor` are used.
- [ ] Semantic `<defs><linearGradient>` with unique lowercase IDs and accessibility tags (`<title>`, `<desc>`, `role="img"`, `aria-labelledby`) are present.
- [ ] Overlapping clutter is eliminated; geometry is cleanly connected.
- [ ] Project `readme.md` uses relative paths to render all assets.
