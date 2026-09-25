# SVG Icon & Vector Graphic Creation — Design Workflow

> **Prompt Version:** 3.4.0
> **Target Environment:** Lovable & Web Design AI Platforms
> **Synchronization:** Main Meta-Repo & Connected Workspaces

This prompt instructs Lovable (and similar web design AI environments) on the exact technical standards for generating scalable, responsive, and accessible SVG vector icons and graphics.

---

## Strictly Avoid (Critical Negative Constraints)

> [!CAUTION]
> **TOTAL BAN ON HTML WRAPPERS, EMBEDDED BITMAPS, AND SOLID BACKGROUNDS (AUTO-REJECT)**
> The AI MUST strictly avoid generating HTML, web pages, solid background shapes, or embedding raster bitmaps inside SVGs.

1. **NO HTML Files or Wrappers:** NEVER create `.html` files and NEVER wrap SVG code inside `<html>`, `<body>`, `<div>`, or any HTML container. Output ONLY the raw SVG code inside an `xml` or `svg` code block.
2. **NO Web Pages or UI Templates:** NEVER generate a website, landing page, mock application, or UI component when asked for an SVG logo/icon.
3. **NO React / Vue Components:** NEVER generate `.tsx`, `.jsx`, or component files unless explicitly requested.
4. **NO Base64 Raster Images:** NEVER embed base64-encoded bitmap images (`<image href="data:image/png;base64,...">`). All graphics MUST be pure mathematical vector paths, polygons, and curves.
5. **NO Hardcoded Fixed Dimensions on Root:** NEVER include hardcoded `width="..."` and `height="..."` attributes on the root `<svg>` tag that prevent responsive scaling; use `viewBox` instead.
6. **NO Solid Background Elements:** NEVER insert a background `<rect>` or container shape (e.g. `<rect width="100%" height="100%" fill="#000"/>` or `<rect fill="#fff"/>`). The canvas MUST remain 100% transparent. "Dark mode" SVGs invert the strokes/fills to lighter colors, but the background is ALWAYS transparent.
7. **NO Messy Overlapping Lines:** NEVER create tangled, chaotic, or densely overlapping strokes. Overlapping lines make icons look muddy, unreadable at small sizes, and amateurish. Maintain clean geometric separation, purposeful negative space, or continuous connected paths.
8. **NO Editor Bloat:** NEVER include third-party editor metadata (e.g., `xmlns:inkscape`, `sodipodi:docname`, `adobe:ns`).
9. **NO Unstyled Elements in Monochrome:** NEVER use hardcoded black/white hex fills on monochrome icons; use `currentColor` so the icon inherits text color dynamically.
10. **NO Unsequenced or Uppercase Filenames:** ALL generated files and images inside project subfolders MUST use lowercase kebab-case preceded by a two-digit zero-padded sequence number (e.g. `01-logo.svg`, `02-logo-dark.svg`, `03-logo-white.svg`, `01-prompt.md`, `02-plan.md`). Unsequenced names like `logo.svg` or uppercase names are strictly BANNED.

---

## 1. Clean SVG Architecture & Design Principles

- **Valid XML & SVG Syntax:** Ensure the output is fully valid XML and well-formed SVG markup.
- **Responsive `viewBox`:** Always define an appropriate `viewBox` (e.g., `viewBox="0 0 24 24"` for UI icons, `viewBox="0 0 100 100"` for logos/illustrations).
- **Omit Hardcoded Dimensions:** Remove hardcoded `width` and `height` attributes on the root `<svg>` element so the icon scales responsively within its CSS container.
- **Anti-Overlapping & Connected Geometry:** Rather than chaotic overlapping lines, connect visual elements purposefully ("connect the dots"). Use continuous smooth strokes, modular aligned shapes, or balanced negative space.
- **Inspirational Archetypes:**
  - *Silhouette & Cutout (Apple model):* Bold solid mark with distinctive negative space.
  - *Modular Grid (Microsoft model):* Clean, non-overlapping geometric tiles separated by consistent gutters.
  - *Rotational Flow (OpenAI model):* Interlocking, symmetrical curved ribbons forming a continuous loop.
  - *Fluid Organic Curve (Sephora model):* Single-stroke elongated "S" curve/flame with zero overlapping lines, conveying beauty, elegance, and fluidity.
- **Minimal Grouping:** Do not wrap elements in redundant `<g>` tags unless needed for shared transforms, styling, or animations.
- **No Editor Metadata:** Do not include Adobe Illustrator, Inkscape, or Figma metadata/namespaces (`xmlns:inkscape`, `sodipodi:docname`, etc.).
- **TOTAL BAN on Base64 Images:** NEVER embed base64-encoded raster images (`<image href="data:image/png;base64,...">`). All artwork must be pure vector paths, polygons, circles, and curves.
- **No HTML Wrapping:** Provide only the raw SVG code inside an `xml` or `svg` code block unless explicitly requested.

---

## 2. Accessibility & Semantics

- Always include a `<title>` and `<desc>` element for screen readers.
- Assign matching `id` attributes to `<title>` and `<desc>` and link them via `aria-labelledby` on the root `<svg>`.
- Add `role="img"` to the root `<svg>`.

---

## 3. Theming & Color Strategy

- **Monochrome Icons:** Use `fill="currentColor"` or `stroke="currentColor"` so the icon seamlessly inherits text color from parent CSS.
- **Multi-Color Logos:** Use semantic CSS variables (`var(--primary)`, `var(--accent)`) or clean standard hex codes.
- **Gradients:** Place `<linearGradient>` and `<radialGradient>` definitions inside `<defs>` with semantic, unique IDs.
- **Variants (All 100% Transparent Canvas):** When generating brand sets, provide:
  - `01-logo.svg`: Default full-color vector on a transparent canvas.
  - `02-logo-dark.svg`: Dark mode variant with inverted, lighter, or vibrant strokes/fills designed for dark themes—on a **100% transparent canvas** (NO dark background `<rect>`).
  - `03-logo-white.svg`: Pure white monochrome variant (`#ffffff`) strokes/fills on a **100% transparent canvas**.

---

## 4. Execution Modes & Concrete Examples

### Version 1: Static Vector SVG Icon (Default)

Generates a clean, production-grade vector icon with full responsive scaling, accessibility tags, and `currentColor` support.

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" role="img" aria-labelledby="shield-check-title shield-check-desc">
  <title id="shield-check-title">Security Verified</title>
  <desc id="shield-check-desc">A shield icon containing a checkmark indicating verified security status.</desc>
  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
  <path d="m9 12 2 2 4-4" />
</svg>
```

### Version 2: Fluid Organic Curve Archetype (Sephora Inspiration)

Generates an elegant single-stroke curved flame with variable weight and zero overlapping lines.

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="flame-title flame-desc">
  <title id="flame-title">Fluid Curve Logo</title>
  <desc id="flame-desc">Elegant single-stroke elongated S-curve flame conveying beauty and fluidity without overlapping lines.</desc>
  <path d="M52 10 C54 22 42 34 38 46 C34 58 40 70 50 82 C42 74 30 64 32 50 C34 36 48 24 52 10 Z" fill="#06b6d4" />
  <path d="M58 24 C62 36 54 48 50 58 C46 68 52 76 60 86 C52 78 42 70 44 58 C46 46 56 36 58 24 Z" fill="#3b82f6" />
</svg>
```

### Version 3: Animated Vector SVG Icon

Generates an interactive or looping animated SVG using self-contained CSS `@keyframes` embedded within a `<style>` block. Ideal for loading states, hero branding, or interactive micro-animations.

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none" role="img" aria-labelledby="sync-title sync-desc">
  <title id="sync-title">Syncing Data</title>
  <desc id="sync-desc">Two rotating curved arrows indicating real-time data synchronization.</desc>
  <style>
    @keyframes spin-clockwise {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
    @keyframes pulse-stroke {
      0%, 100% { stroke-opacity: 0.5; }
      50% { stroke-opacity: 1; }
    }
    .rotating-group {
      transform-origin: 50px 50px;
      animation: spin-clockwise 3s linear infinite;
    }
    .pulsing-arrow {
      animation: pulse-stroke 1.5s ease-in-out infinite;
    }
  </style>
  <g class="rotating-group">
    <path class="pulsing-arrow" d="M50 15 A35 35 0 0 1 85 50 L75 50 L90 65 L95 50 L85 50" stroke="#06b6d4" stroke-width="5" stroke-linecap="round" stroke-linejoin="round" fill="none" />
    <path class="pulsing-arrow" d="M50 85 A35 35 0 0 1 15 50 L25 50 L10 35 L5 50 L15 50" stroke="#3b82f6" stroke-width="5" stroke-linecap="round" stroke-linejoin="round" fill="none" />
  </g>
  <circle cx="50" cy="50" r="8" fill="#06b6d4" />
</svg>
```

---

## 5. Directory & File Hierarchy (Strict Lowercase & Two-Digit Sequence)

When generating SVG assets as part of a project, assets must follow strict lowercase naming under `02-projects/{sequence}-{project-name}/`. The AI must persist the design plan to `prompts/02-plan.md` and the generation prompt to `prompts/01-prompt.md`:

```
/ (repo root)
└── 02-projects/
    ├── 01-{project-name}/
    │   ├── readme.md (overview and vector previews)
    │   ├── prompts/
    │   │   ├── 01-prompt.md (the exact prompt, user inputs, and AI parameters used)
    │   │   └── 02-plan.md (design plan, geometry breakdown, and verification checklist)
    │   └── icons-svg/
    │       ├── 01-logo.svg (primary full-color vector on transparent canvas)
    │       ├── 02-logo-dark.svg (inverted vector for dark themes on transparent canvas)
    │       └── 03-logo-white.svg (pure white monochrome vector on transparent canvas)
    └── 02-{project-name}/
```

### Hierarchy Rules

1. **Root Projects Folder:** All projects live under `02-projects/`.
2. **Project Folder Naming:** `{sequence}-{project-name}` using two-digit zero-padding and kebab-case.
3. **Prompt & Plan Preservation (Mandatory):** When formulating the design strategy, every detail of the design plan (geometry, color palette, viewbox, accessibility tags, and verification checklist) MUST be saved directly to the file system at `prompts/02-plan.md`. The exact prompt given to the generation engine MUST be saved in `prompts/01-prompt.md`.
4. **Asset Organization:** All SVG vector files MUST be stored inside `icons-svg/` with zero-padded sequence numbers (`01-`, `02-`, etc.).
5. **Relative Paths:** All links and image embeds in `readme.md` must use relative paths.

---

## 6. Output Format

- Provide strictly the raw SVG code inside a fenced code block with language identifier `xml` or `svg`.
- Do not output HTML wrappers or surrounding page containers.

