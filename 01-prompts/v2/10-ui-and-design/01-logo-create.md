# Logo Design & Branding Generation — Lovable Design Workflow

> **Prompt Version:** 3.4.0
> **Target Environment:** Lovable & Web Design AI Platforms
> **Synchronization:** Main Meta-Repo & Connected Workspaces

This prompt instructs Lovable (and similar web design AI environments) to generate logo, icon, and branding assets for a given product or brand.

---

## Strictly Avoid (Critical Negative Constraints)

> [!CAUTION]
> **TOTAL BAN ON HTML, WEBSITES, PAGES, AND APPLICATION CODE (AUTO-REJECT)**
> The AI MUST strictly avoid generating any web application or frontend/backend code. This prompt is ONLY for branding assets and icons.

1. **NO HTML Files:** NEVER create or output `.html` files (no `index.html`, `app.html`, `test.html`, or any other HTML files).
2. **NO Web Pages or Mockup Apps:** NEVER build landing pages, homepages, dashboard pages, or mock web applications.
3. **NO React / Vue / Svelte Components:** NEVER create `.tsx`, `.jsx`, `.vue`, or `.svelte` components.
4. **NO CSS Stylesheets:** NEVER create external CSS files (`styles.css`, `app.css`). Theme colors must live ONLY in `colors-themes/palette.md`.
5. **NO Server or Script Code:** NEVER create JavaScript/TypeScript scripts, API routes, or backend servers.
6. **NO Tokens.json:** NEVER generate `tokens.json` in `colors-themes/` (it has been completely removed as unnecessary; use `palette.md` only).
7. **NO Solid / Opaque Backgrounds on ANY PNG (TOTAL BAN):** NEVER generate ANY PNG icon with a solid background (no solid white, no solid black, no grey, no dark square backdrop). ALL PNG files without exception MUST have a 100% transparent alpha-channel background (`#00000000`). "Dark" and "Light" variants refer ONLY to the color/fill of the icon artwork itself—NEVER the background! The canvas/background is ALWAYS 100% transparent.
8. **NO Solid Background Shapes in SVGs:** NEVER add a background `<rect>` or container shape (e.g. `<rect fill="#000" .../>` or `<rect fill="#fff" .../>`). The SVG canvas must always be transparent.
9. **NO Messy Overlapping Lines:** NEVER create tangled, chaotic, or densely overlapping strokes. Overlapping lines make icons look muddy, unreadable at small resolutions, and amateurish. Maintain clean geometric separation, purposeful negative space, or continuous connected paths.
10. **NO Unsequenced or Uppercase Filenames:** ALL generated files and images inside project subfolders MUST use lowercase kebab-case preceded by a two-digit zero-padded sequence number (e.g. `01-logo.svg`, `02-logo-dark.svg`, `01-logo-052.png`, `01-prompt.md`, `02-plan.md`). Unsequenced names like `logo.svg` or uppercase names like `README.md` are strictly BANNED.
11. **NO Base64 Images in SVGs:** NEVER embed raster images or base64 data URLs inside SVG files.

---

## 1. Input Capture & Clarification

Before generating any assets, capture and validate the following inputs:

1. `product_name`: Name of the product, service, or brand.
2. `product_idea`: Brief explanation of the product, its purpose, target audience, and brand voice/tone (e.g., modern, playful, corporate, technical, minimalist).
3. `sample_colors`: Target color hints, hex codes, or preferred palette mood (e.g., "deep navy and electric cyan").
4. `inspiration_logos`: Optional reference brands or logo styles (e.g., "Apple", "Microsoft", "OpenAI", "Sephora", "Stripe"). If provided, prioritize the visual harmony, silhouette, and geometric simplicity of that reference.
5. `needs_dark_white_variants`: Boolean indicating if dark-theme and white monochrome variants are required (defaults to `true`).
6. `is_animated`: Boolean indicating whether animated assets (`gif-animation.gif`, animated SVG) are requested (defaults to `false` unless explicitly asked).

> [!IMPORTANT]
> If `product_name`, `product_idea`, or `sample_colors` are not provided, **STOP and ask the user** for clarification before proceeding with asset generation.

---

## 2. Brand Identity & Icon Design Principles

### A. Conveying the Brand Image

- The icon must be an immediate visual metaphor for the brand's core identity (e.g., speed, intelligence, security, luxury, simplicity, or interconnectedness).
- Design for instant memorability: a user should be able to sketch the basic concept from memory after seeing it once.
- Ensure micro-to-macro scalability: the icon must remain crystal-clear at 16x16 px (favicon) and 52x52 px (app icon) without losing definition, while looking balanced and elegant at 1024x1024 px.

### B. Anti-Overlapping Rule & Clean Geometry

- **Ban on Chaotic Overlaps:** Do not cross multiple lines or stack overlapping wireframe shapes. Chaos reduces contrast and creates rendering artifacts at small sizes.
- **Connect the Dots:** Connect visual elements purposefully through continuous strokes, modular alignment, or clean geometric joints rather than piling disconnected strokes on top of each other.
- **Purposeful Negative Space:** Use negative space intentionally to define secondary shapes or give the primary mark breathing room.

---

## 3. Inspirational Archetypes & Code Models

When designing, draw inspiration from four proven industry archetypes:

1. **The Silhouette & Negative Space Model (Apple Archetype):**
   - A single, bold silhouette with a distinctive subtracted shape or cut.
   - Clean outline, zero internal line clutter, immediately recognizable at 16px.
   - *Example Concept:* A solid geometric fruit, shield, or gem with an intentional negative-space cutout mask.

2. **The Modular Grid Model (Microsoft Archetype):**
   - Distinct, non-overlapping geometric shapes organized with harmonious negative-space gutters.
   - Communicates ecosystem, collaboration, and modular structure.
   - *Example Concept:* 4 rounded tiles or geometric shapes aligned in a 2x2 grid with consistent negative space between them.

3. **The Continuous Ribbon / Rotational Flow Model (OpenAI Archetype):**
   - Interlocking rotational symmetry formed by a single continuous path or connected ribbon segments.
   - Communicates intelligence, iteration, cycles, and connectivity without cluttered intersections.
   - *Example Concept:* Symmetrical curved segments radiating from a central core in rotational flow.

4. **The Fluid Organic Curve / Flame Model (Sephora Archetype):**
   - An elegant, single-stroke elongated "S" curve or flame motif.
   - Conveys luxury, beauty, cosmetics, fluidity, and sophistication.
   - Zero overlapping lines; relies on variable stroke weight, organic taper, and clean curvature.
   - *Example Concept:* A single tapered bezier curve that flows smoothly from thick to thin without any intersecting lines.

---

## 4. Directory & File Hierarchy (Strict Lowercase & Two-Digit Sequence)

All generated assets must follow strict lowercase naming and two-digit zero-padded sequence numbers. No unsequenced or uppercase files are permitted in filenames or folder paths.

```
/ (repo root)
├── 01-prompts/
└── 02-projects/
    ├── 01-{project-name}/
    │   ├── readme.md
    │   ├── prompts/
    │   │   ├── 01-prompt.md
    │   │   └── 02-plan.md
    │   ├── icons-svg/
    │   │   ├── 01-logo.svg
    │   │   ├── 02-logo-dark.svg
    │   │   └── 03-logo-white.svg
    │   ├── icons-image/
    │   │   ├── 01-logo-052.png
    │   │   ├── 02-logo-128.png
    │   │   ├── 03-logo-256.png
    │   │   ├── 04-logo-512.png
    │   │   ├── 05-logo-1024-light.png
    │   │   └── 06-logo-1024-dark.png
    │   └── colors-themes/
    │       └── 01-palette.md
    └── 02-{project-name}/
(repo root assets)
├── favicon.ico
├── favicon.png
└── gif-animation.gif (generated only when animation is requested)
```

### Hierarchy Rules

1. **Root Projects Folder:** All projects live under `02-projects/`.
2. **Project Folder Naming:** `{sequence}-{project-name}` using two-digit zero-padding and kebab-case (e.g., `01-acme-pay`, `02-cloud-sync`).
3. **Never Overwrite:** Never overwrite an existing project folder; always increment the sequence number (`01`, `02`, ...).
4. **Prompt & Plan Preservation (Mandatory):** When formulating the design strategy, every detail of the design plan (goals, composition, color palette, typography, execution steps, and quality checklist) MUST be saved directly to the file system at `prompts/02-plan.md`. The exact prompt given to the generation engine, user inputs, and model parameters MUST be saved in `prompts/01-prompt.md` so designs can be reproduced, audited, and iterated on.
5. **No Tokens.json:** Do not generate `tokens.json`. Color specifications live exclusively inside `colors-themes/01-palette.md`.
6. **Favicon Placement:** `favicon.ico` and `favicon.png` are placed at the repository root and regenerated per project.
7. **Relative Paths:** All links and image embeds in `readme.md` must use relative paths so GitHub renders them natively.

---

## 5. Asset Specifications

### A. Vector Icons (`icons-svg/`)

- Pure, clean XML/SVG vector code without editor junk or base64 embedded bitmaps.
- Responsive `viewBox` (e.g., `viewBox="0 0 100 100"`), omitting hardcoded `width` and `height` attributes on the root `<svg>`.
- Always 100% transparent canvas (never add a background `<rect>`).
- Standard variants:
  - `01-logo.svg`: Primary full-color logo on a transparent canvas.
  - `02-logo-dark.svg`: Optimized for dark themes using inverted, lighter, or vibrant strokes/fills—on a **100% transparent canvas** (NO dark background `<rect>`).
  - `03-logo-white.svg`: Pure white monochrome icon (`#ffffff`) for dark surfaces or overlays—on a **100% transparent canvas**.

### B. Transparent Raster Images (`icons-image/`)

- **100% Transparency Requirement (Zero Solid Backgrounds):** Every single PNG image in `icons-image/` MUST have a **100% transparent alpha-channel background**. NO solid white, solid black, solid grey, or colored background boxes are permitted.
- **Understanding "Dark" vs. "Light" Variants:**
  - **"Dark" does NOT mean a dark background!** It means the icon artwork itself uses inverted, lighter, or white strokes and fills so that when the transparent PNG is placed onto a dark website theme, the icon is clearly visible. The canvas/background is 100% transparent.
  - **"Light" means the icon artwork uses standard, darker, or saturated brand colors** intended to be placed on a light website theme. The canvas/background is 100% transparent.
  - Both `05-logo-1024-light.png` and `06-logo-1024-dark.png` are **100% transparent PNGs** with zero background fills.
- **Specification Listing (consecutive with zero blank lines between entries):**
  - `01-logo-052.png` — 52x52 px 100% transparent PNG icon
  - `02-logo-128.png` — 128x128 px 100% transparent PNG icon
  - `03-logo-256.png` — 256x256 px 100% transparent PNG icon
  - `04-logo-512.png` — 512x512 px 100% transparent PNG icon
  - `05-logo-1024-light.png` — 1024x1024 px 100% transparent PNG with dark/saturated artwork for light themes
  - `06-logo-1024-dark.png` — 1024x1024 px 100% transparent PNG with inverted/light artwork for dark themes

### C. Color Themes (`colors-themes/01-palette.md`)

- A concise Markdown document listing brand colors with HEX, RGB, and HSL values.
- Includes clear guidance on primary brand color, secondary accent, neutral dark, neutral light, and feedback/glow shades.

### D. Favicons (Repo Root)

- `favicon.ico`: Multi-resolution icon for browser tabs.
- `favicon.png`: High-resolution 32x32 px or 64x64 px 100% transparent PNG favicon.

### E. Generation Prompt & Plan Archive (`prompts/`)

- `01-prompt.md`: Contains the full generation prompt, model parameters (aspect ratio, style, negative prompts), and exact text strings used for the generation run.
- `02-plan.md`: Contains the full design plan, geometry breakdown, color palette, typography specification, and verification checklist.

---

## 6. Prompt Execution Modes

### Version 1: Static Logo & Branding Generation (Default)

Executes standard branding creation:
1. Validates user inputs (product name, idea, tone, colors, optional inspiration).
2. Applies clean geometry principles (anti-overlapping, connected dots, scalability).
3. Creates the project folder under `02-projects/{seq}-{project-name}/`.
4. Creates `icons-svg/` with `01-logo.svg`, `02-logo-dark.svg`, and `03-logo-white.svg` (all on transparent canvas).
5. Creates `icons-image/` with the 6 consecutive 100% transparent PNG sizes (`01-logo-052.png` through `06-logo-1024-dark.png`).
6. Creates `colors-themes/01-palette.md`.
7. Generates `favicon.ico` and `favicon.png` at the repository root.
8. Generates `readme.md` in the project directory displaying all assets in a GitHub-compatible table.

### Version 2: Animated Logo & Branding Generation

Executes static branding creation PLUS animated assets:
1. Completes all steps from Version 1.
2. Generates an animated SVG (`04-logo-animated.svg` in `icons-svg/`) on a transparent canvas using smooth CSS `@keyframes` or SMIL for subtle motion (e.g., stroke-dasharray draw effect, pulsing glow, or rotational geometry).
3. Generates `gif-animation.gif` placed at the repo root showcasing the animated icon loop (transparent or smooth background, 60–120 frames, seamless loop, 24–30 fps).
4. Embeds the animated preview into `readme.md`.

---

## 7. Concrete Asset Examples

### Example 1: Connected Geometric Logo (`icons-svg/01-logo.svg`)

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="logo-title logo-desc">
  <title id="logo-title">AcmePay Logo</title>
  <desc id="logo-desc">Connected geometric hexagon with directional arrow nodes in cyan and navy.</desc>
  <defs>
    <linearGradient id="primary-grad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#06b6d4" />
      <stop offset="100%" stop-color="#3b82f6" />
    </linearGradient>
  </defs>
  <polygon points="50,8 88,29 88,71 50,92 12,71 12,29" stroke="url(#primary-grad)" stroke-width="5" fill="none" stroke-linejoin="round" />
  <circle cx="50" cy="8" r="4" fill="#06b6d4" />
  <circle cx="88" cy="29" r="4" fill="#3b82f6" />
  <circle cx="88" cy="71" r="4" fill="#3b82f6" />
  <circle cx="50" cy="92" r="4" fill="#06b6d4" />
  <circle cx="12" cy="71" r="4" fill="#06b6d4" />
  <circle cx="12" cy="29" r="4" fill="#06b6d4" />
  <path d="M35 50 L50 35 L65 50 M50 35 L50 65" stroke="#06b6d4" stroke-width="5" stroke-linecap="round" stroke-linejoin="round" />
</svg>
```

### Example 2: Silhouette & Negative Space Archetype (Apple Inspiration)

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="sil-title sil-desc">
  <title id="sil-title">Silhouette with Negative Space</title>
  <desc id="sil-desc">Bold iconic silhouette with an intentional transparent negative-space cutout.</desc>
  <mask id="cutout-mask">
    <rect width="100" height="100" fill="#ffffff" />
    <circle cx="62" cy="44" r="14" fill="#000000" />
  </mask>
  <path d="M50 16 C34 16 22 28 22 46 C22 66 38 84 50 84 C62 84 78 66 78 46 C78 28 66 16 50 16 Z" fill="#06b6d4" mask="url(#cutout-mask)" />
  <path d="M50 8 C54 12 56 18 52 22 C48 18 46 12 50 8 Z" fill="#3b82f6" />
</svg>
```

### Example 3: Modular Grid Archetype (Microsoft Inspiration)

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="mod-title mod-desc">
  <title id="mod-title">Modular Grid Logo</title>
  <desc id="mod-desc">Four non-overlapping geometric quadrants separated by clean negative space.</desc>
  <rect x="12" y="12" width="34" height="34" rx="6" fill="#06b6d4" />
  <rect x="54" y="12" width="34" height="34" rx="6" fill="#3b82f6" />
  <rect x="12" y="54" width="34" height="34" rx="6" fill="#0284c7" />
  <rect x="54" y="54" width="34" height="34" rx="6" fill="#6366f1" />
</svg>
```

### Example 4: Continuous Rotational Flow (OpenAI Inspiration)

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round" role="img" aria-labelledby="rot-title rot-desc">
  <title id="rot-title">Continuous Rotational Flow</title>
  <desc id="rot-desc">Interlocking curved segments radiating in rotational symmetry without chaotic overlaps.</desc>
  <g transform="translate(50,50)">
    <path d="M0 -35 C15 -35, 30 -20, 30 0 L15 0 C15 -10, 8 -20, 0 -20 Z" fill="#06b6d4" stroke="none" />
    <path d="M0 -35 C15 -35, 30 -20, 30 0 L15 0 C15 -10, 8 -20, 0 -20 Z" fill="#3b82f6" stroke="none" transform="rotate(120)" />
    <path d="M0 -35 C15 -35, 30 -20, 30 0 L15 0 C15 -10, 8 -20, 0 -20 Z" fill="#6366f1" stroke="none" transform="rotate(240)" />
  </g>
</svg>
```

### Example 5: Fluid Organic Curve Archetype (Sephora Inspiration)

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="sep-title sep-desc">
  <title id="sep-title">Fluid Organic Curve Logo</title>
  <desc id="sep-desc">Elegant single-stroke elongated S-curve flame conveying beauty, luxury, and fluidity without overlapping lines.</desc>
  <path d="M52 10 C54 22 42 34 38 46 C34 58 40 70 50 82 C42 74 30 64 32 50 C34 36 48 24 52 10 Z" fill="#06b6d4" />
  <path d="M58 24 C62 36 54 48 50 58 C46 68 52 76 60 86 C52 78 42 70 44 58 C46 46 56 36 58 24 Z" fill="#3b82f6" />
</svg>
```

### Example 6: Animated SVG Logo (`icons-svg/04-logo-animated.svg`)

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="anim-title anim-desc">
  <title id="anim-title">AcmePay Animated Logo</title>
  <desc id="anim-desc">Animated geometric hexagon with glowing pulse and rotating core.</desc>
  <style>
    @keyframes pulse-glow {
      0%, 100% { stroke-opacity: 0.6; filter: drop-shadow(0 0 2px #06b6d4); }
      50% { stroke-opacity: 1; filter: drop-shadow(0 0 8px #06b6d4); }
    }
    @keyframes rotate-core {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
    .hex-pulse {
      animation: pulse-glow 3s ease-in-out infinite;
    }
    .core-spin {
      transform-origin: 50px 50px;
      animation: rotate-core 12s linear infinite;
    }
  </style>
  <polygon class="hex-pulse" points="50,8 88,29 88,71 50,92 12,71 12,29" stroke="#06b6d4" stroke-width="5" fill="none" stroke-linejoin="round" />
  <g class="core-spin">
    <circle cx="50" cy="50" r="16" stroke="#3b82f6" stroke-width="4" stroke-dasharray="6 4" fill="none" />
    <circle cx="50" cy="34" r="4" fill="#06b6d4" />
  </g>
</svg>
```

### Example 7: Color Palette (`colors-themes/01-palette.md`)

```markdown
# Brand Color Palette: AcmePay

## Primary Colors
| Role | Name | Swatch | HEX | RGB | HSL |
|---|---|---|---|---|---|
| Primary | Cyan Glow | ![#06b6d4](https://placehold.co/15x15/06b6d4/06b6d4.png) | `#06b6d4` | `rgb(6, 182, 212)` | `hsl(189, 94%, 43%)` |
| Accent | Electric Blue | ![#3b82f6](https://placehold.co/15x15/3b82f6/3b82f6.png) | `#3b82f6` | `rgb(59, 130, 246)` | `hsl(217, 91%, 60%)` |

## Neutral & Surface Colors
| Role | Name | Swatch | HEX | RGB | HSL |
|---|---|---|---|---|---|
| Background Dark | Midnight Navy | ![#0f172a](https://placehold.co/15x15/0f172a/0f172a.png) | `#0f172a` | `rgb(15, 23, 42)` | `hsl(222, 47%, 11%)` |
| Surface Light | Pure White | ![#ffffff](https://placehold.co/15x15/ffffff/ffffff.png) | `#ffffff` | `rgb(255, 255, 255)` | `hsl(0, 0%, 100%)` |
```

### Example 8: Project README Display (`02-projects/01-acme-pay/readme.md`)

```markdown
# AcmePay — Branding & Logo Assets

## Overview
AcmePay is a modern financial platform designed for effortless developer billing.

## Vector Logos
| Primary Logo (Transparent) | Dark Theme (Inverted, Transparent) | Monochrome White (Transparent) |
|:---:|:---:|:---:|
| ![Primary](icons-svg/01-logo.svg) | ![Dark](icons-svg/02-logo-dark.svg) | ![White](icons-svg/03-logo-white.svg) |

## Transparent Icon Sizes
| 52px | 128px | 256px | 512px |
|:---:|:---:|:---:|:---:|
| <img src="icons-image/01-logo-052.png" width="52" /> | <img src="icons-image/02-logo-128.png" width="128" /> | <img src="icons-image/03-logo-256.png" width="256" /> | <img src="icons-image/04-logo-512.png" width="256" /> |

## High-Resolution Transparent Icons (1024px)
| For Light Themes (Transparent PNG) | For Dark Themes (Inverted, Transparent PNG) |
|:---:|:---:|
| <img src="icons-image/05-logo-1024-light.png" width="300" /> | <img src="icons-image/06-logo-1024-dark.png" width="300" /> |

## Color Palette
See [Color Palette](colors-themes/01-palette.md) for hex codes, swatches, and usage rules.
```
