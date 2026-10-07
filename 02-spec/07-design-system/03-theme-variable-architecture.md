# Theme & Variable Architecture

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

This is the **single source of truth** for every color, spacing, and visual token in the design system. All components reference these variables — never hardcoded values. To re-theme the entire system, change values in this registry only.

Variables are defined as **CSS custom properties** in `src/index.css` within `:root` (light) and `.dark` (dark) blocks. Values use **HSL space-separated format** (e.g., `252 85% 60%`) to allow opacity modifiers.

---

## Variable Format

```css
--token-name: H S% L%;
```

Usage in components:

```css
color: hsl(var(--token-name));                /* full opacity */
background: hsl(var(--token-name) / 0.1);     /* 10% opacity */
```

Usage in Tailwind classes:

```html
<div class="bg-primary text-primary-foreground">
```

---

## Complete Token Registry

### 1. Primitive Color Ramps (3-Format Reference)

#### 1.1 Blue Ramp (Primary Action & Corporate Trust)
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
|:---|:---|:---|:---|:---|:---|
| `--blue-50` | `oklch(0.972 0.014 259)` | `#F0F6FF` | `rgb(240, 246, 255)` | `hsl(216, 100%, 97%)` | Hero wash start, brand pill bg |
| `--blue-100` | `oklch(0.935 0.032 262)` | `#DEEAFF` | `rgb(222, 234, 255)` | `hsl(218, 100%, 94%)` | Badge hover fill, soft tile ring |
| `--blue-200` | `oklch(0.874 0.061 263)` | `#C1D6FF` | `rgb(193, 214, 255)` | `hsl(220, 100%, 88%)` | Chrome gradient stop, subtle border |
| `--blue-300` | `oklch(0.777 0.104 264)` | `#95B6FA` | `rgb(149, 182, 250)` | `hsl(220, 91%, 78%)` | Chrome gradient end, data series |
| `--blue-400` | `oklch(0.662 0.164 264)` | `#5E8EF6` | `rgb(94, 142, 246)` | `hsl(221, 89%, 67%)` | Gradient text end stop, ring focus |
| `--blue-500` | `oklch(0.546 0.215 262.9)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | **Primary Action (`--brand-secondary`)** |
| `--blue-600` | `oklch(0.489 0.208 263.5)` | `#1B51D3` | `rgb(27, 81, 211)` | `hsl(222, 77%, 47%)` | Brand gradient start, active CTA |
| `--blue-700` | `oklch(0.409 0.176 264.5)` | `#173DA7` | `rgb(23, 61, 167)` | `hsl(224, 76%, 37%)` | Brand pill text, high-contrast link |
| `--blue-800` | `oklch(0.369 0.147 264.8)` | `#17368B` | `rgb(23, 54, 139)` | `hsl(224, 72%, 32%)` | Deep interactive border, badge text |
| `--blue-900` | `oklch(0.317 0.135 264.5)` | `#0D2975` | `rgb(13, 41, 117)` | `hsl(224, 80%, 25%)` | **Core Brand Anchor (`--brand-primary`)** |

#### 1.2 Violet Ramp (Tertiary Accent & Gradient Partner)
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
|:---|:---|:---|:---|:---|:---|
| `--violet-50` | `oklch(0.972 0.017 300)` | `#F7F3FF` | `rgb(247, 243, 255)` | `hsl(260, 100%, 98%)` | Accent surface fill |
| `--violet-100` | `oklch(0.946 0.032 300)` | `#F0E9FF` | `rgb(240, 233, 255)` | `hsl(259, 100%, 96%)` | Subtle chip background |
| `--violet-200` | `oklch(0.885 0.068 299)` | `#DFD0FF` | `rgb(223, 208, 255)` | `hsl(259, 100%, 91%)` | Accent border on light cards |
| `--violet-300` | `oklch(0.784 0.128 298)` | `#C3A6FF` | `rgb(195, 166, 255)` | `hsl(260, 100%, 83%)` | Soft glow highlight on dark |
| `--violet-400` | `oklch(0.659 0.187 297.5)` | `#A173F4` | `rgb(161, 115, 244)` | `hsl(261, 85%, 70%)` | Secondary chart accent |
| `--violet-500` | `oklch(0.532 0.253 296.9)` | `#822EE8` | `rgb(130, 46, 232)` | `hsl(267, 80%, 55%)` | **Tertiary Brand Accent (`--brand-tertiary`)** |
| `--violet-600` | `oklch(0.481 0.243 296.5)` | `#721CD2` | `rgb(114, 28, 210)` | `hsl(268, 76%, 47%)` | Brand gradient midpoint (55%) |
| `--violet-700` | `oklch(0.417 0.213 295.9)` | `#5C14AF` | `rgb(92, 20, 175)` | `hsl(268, 79%, 38%)` | Accent foreground text |
| `--violet-800` | `oklch(0.362 0.187 295.5)` | `#4B0D92` | `rgb(75, 13, 146)` | `hsl(268, 84%, 31%)` | Deep violet structural accent |
| `--violet-900` | `oklch(0.257 0.153 294.5)` | `#2D0063` | `rgb(45, 0, 99)` | `hsl(267, 100%, 19%)` | Maximum depth violet anchor |

#### 1.3 Cyan, Indigo & Neutral Cool Slate Ramps
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
|:---|:---|:---|:---|:---|:---|
| `--cyan-400` | `oklch(0.797 0.136 205.6)` | `#03D5E7` | `rgb(3, 213, 231)` | `hsl(185, 97%, 46%)` | **Brand Highlight (`--brand-highlight`)**, `.shine-sweep` |
| `--indigo-500` | `oklch(0.545 0.216 277)` | `#5856E9` | `rgb(88, 86, 233)` | `hsl(241, 77%, 63%)` | Workflow pin #4, scroll-stack hue |
| `--teal-500` | `oklch(0.63 0.118 195)` | `#009F9F` | `rgb(0, 159, 159)` | `hsl(180, 100%, 31%)` | Workflow pin #5, scroll-stack hue |
| `--neutral-200` | `oklch(0.924 0.011 265)` | `#E2E6ED` | `rgb(226, 230, 237)` | `hsl(218, 23%, 91%)` | **Light Hairline Border (`--border`, `--input`)** |
| `--neutral-400` | `oklch(0.688 0.036 264)` | `#8F9BB2` | `rgb(143, 155, 178)` | `hsl(219, 19%, 63%)` | Dark-band muted text (`--on-dark-muted`) |
| `--neutral-600` | `oklch(0.441 0.043 265.9)` | `#48536B` | `rgb(72, 83, 107)` | `hsl(221, 20%, 35%)` | **Soft Ink (`--ink-soft`, `--muted-foreground`)** |
| `--neutral-900` | `oklch(0.163 0.023 265)` | `#090E18` | `rgb(9, 14, 24)` | `hsl(220, 45%, 6%)` | **Primary Ink (`--ink`, `--foreground`)** |
| `--surface-dark` | `oklch(0.129 0.019 264)` | `#0A0E17` | `rgb(10, 14, 23)` | `hsl(222, 39%, 6%)` | Dark section canvas / header bar |

---

### 2. Core Gradients, Shadows & Geometry Tokens

#### 2.1 Brand Gradients
```css
:root {
  --gradient-hero: radial-gradient(70% 55% at 78% 8%, color-mix(in oklab, var(--blue-500) 34%, transparent) 0%, transparent 70%),
                   radial-gradient(50% 45% at 8% 96%, color-mix(in oklab, var(--violet-600) 20%, transparent) 0%, transparent 72%),
                   linear-gradient(180deg, var(--surface-dark) 0%, oklch(0.155 0.03 265) 62%, var(--surface-dark) 100%);
  --gradient-accent: linear-gradient(90deg, var(--blue-500) 0%, var(--violet-500) 100%);
  --gradient-accent-hover: linear-gradient(90deg, color-mix(in oklab, var(--blue-500) 88%, white) 0%, color-mix(in oklab, var(--violet-500) 88%, white) 100%);
  --gradient-brand: linear-gradient(120deg, var(--blue-600) 0%, var(--violet-600) 55%, var(--blue-500) 100%);
  --gradient-text: linear-gradient(92deg, var(--blue-500) 0%, var(--violet-500) 55%, var(--blue-400) 100%);
  --gradient-halo: radial-gradient(55% 60% at 72% 12%, color-mix(in oklab, var(--blue-400) 26%, transparent) 0%, color-mix(in oklab, var(--violet-500) 10%, transparent) 45%, transparent 72%);
  --gradient-card-dark: linear-gradient(160deg, oklch(1 0 0 / 7%) 0%, oklch(1 0 0 / 2%) 100%);
  --gradient-border: linear-gradient(120deg, color-mix(in oklab, var(--blue-500) 60%, transparent), color-mix(in oklab, var(--violet-500) 35%, transparent), color-mix(in oklab, var(--cyan-400) 40%, transparent));
  --gradient-chrome: linear-gradient(180deg, oklch(1 0 0) 0%, var(--blue-200) 28%, var(--neutral-500) 52%, oklch(1 0 0) 66%, var(--blue-300) 100%);
}
```

#### 2.2 Shadows & Radii
```css
:root {
  --shadow-xs: 0 1px 2px color-mix(in oklab, var(--ink) 5%, transparent);
  --shadow-card: 0 8px 24px -12px color-mix(in oklab, var(--brand-primary) 18%, transparent);
  --shadow-lift: 0 24px 60px -24px color-mix(in oklab, var(--brand-primary) 35%, transparent);
  --shadow-glow: 0 0 60px -12px color-mix(in oklab, var(--blue-500) 55%, transparent);
  --shadow-dark-card: 0 24px 60px -30px oklch(0 0 0 / 85%);
  --shadow-3d: 0 40px 90px -40px color-mix(in oklab, var(--brand-primary) 45%, transparent), 0 8px 20px -12px color-mix(in oklab, var(--brand-primary) 25%, transparent);

  --radius: 14px;
  --radius-button: 12px;
  --radius-card: 20px;
  --radius-media: 28px;
  --radius-neu: 18px;

  --dur-instant: 120ms;
  --dur-fast: 240ms;
  --dur-base: 420ms;
  --dur-slow: 700ms;
  --dur-cine: 1100ms;
}
```

### Reading / Prose Theme Tokens

| Token | Light Value | Dark Value | Purpose |
|-------|-------------|------------|---------|
| `--heading-gradient-from` | `252 85% 60%` | `252 85% 70%` | Heading gradient start |
| `--heading-gradient-to` | `330 85% 60%` | `330 85% 70%` | Heading gradient end |
| `--link-color` | `252 85% 55%` | `252 85% 72%` | Link text color |
| `--code-bg` | `250 25% 95%` | `230 20% 15%` | Inline code background |
| `--code-text` | `330 85% 45%` | `330 85% 70%` | Inline code text color |
| `--blockquote-border` | `252 85% 70%` | `252 60% 50%` | Blockquote left border |
| `--table-header-bg` | `252 85% 97%` | `230 20% 14%` | Table header background |
| `--table-row-hover` | `252 85% 97%` | `252 40% 15%` | Table row hover background |
| `--highlight-glow` | `252 85% 60%` | `252 85% 65%` | Inline code hover glow |

### Code Block Component Tokens

These are **NOT** CSS custom properties — they are fixed values used in the code block component. They are intentionally NOT themed because code blocks maintain a consistent dark appearance in both light and dark modes.

| Property | Value | Purpose |
|----------|-------|---------|
| Block background | `hsl(220, 14%, 11%)` | Always-dark code background |
| Header background | `hsl(220, 14%, 14%)` | Code header bar |
| Header border | `hsl(220, 13%, 20%)` | Header bottom border |
| Block border | `hsl(220, 13%, 22%)` | Outer border |
| Line number background | `hsl(220, 14%, 9%)` | Line number gutter |
| Line number border | `hsl(220, 13%, 18%)` | Gutter right border |
| Line number text | `hsl(220, 10%, 35%)` | Line number text |
| Tool button background | `hsl(220, 13%, 20%)` | Header action buttons |
| Tool button border | `hsl(220, 13%, 25%)` | Button borders |
| Tool button hover bg | `hsl(220, 13%, 28%)` | Button hover |
| Font controls bg | `hsl(220, 13%, 18%)` | Font size control group |
| `--lang-accent` | Per-language HSL | Language-specific glow color |
| `--code-font-size` | `18px` | Default code font size |
| `--code-line-height` | `1.6` | Code line height |

### Language Accent Colors

Each language gets a unique HSL accent used for the badge dot, hover glow, and fullscreen shadow:

| Language | HSL Value |
|----------|-----------|
| TypeScript | `210 80% 60%` |
| JavaScript | `50 90% 55%` |
| Go | `190 80% 50%` |
| Rust | `20 85% 55%` |
| CSS | `280 70% 60%` |
| JSON | `45 85% 55%` |
| Bash/Shell | `120 50% 50%` |
| SQL | `200 70% 55%` |
| Markdown | `252 60% 60%` |
| YAML | `340 60% 55%` |
| PHP | `240 55% 60%` |
| HTML/XML | `15 80% 55%` |
| Plain Text | `220 10% 65%` |
| Tree/Structure | `220 10% 65%` |

---

## Tailwind Config Mapping

Every CSS custom property is mapped to a Tailwind utility class in `tailwind.config.ts`:

```typescript
colors: {
  primary: {
    DEFAULT: "hsl(var(--primary))",
    foreground: "hsl(var(--primary-foreground))",
  },
  accent: {
    DEFAULT: "hsl(var(--accent))",
    foreground: "hsl(var(--accent-foreground))",
  },
  // ... all tokens follow same pattern
}
```

This enables usage like:

```html
<button class="bg-primary text-primary-foreground hover:bg-primary/90">
```

---

## How to Create a New Theme Variant

### Step 1: Copy the `:root` and `.dark` blocks

### Step 2: Adjust HSL values

| Adjustment | Technique |
|-----------|-----------|
| Change brand identity | Alter hue on `--primary` and `--accent` |
| Make warmer | Shift hues toward 30-60° range |
| Make cooler | Shift hues toward 200-240° range |
| Increase contrast | Increase saturation, widen lightness gaps |
| Soften design | Decrease saturation by 10-20% |
| Dark mode | Backgrounds: 5-12% lightness; Foregrounds: 85-95% lightness |

### Step 3: Validate

- Heading gradients should remain visually distinct
- Primary and accent should have minimum 60° hue difference
- Text on colored backgrounds must meet WCAG AA contrast (4.5:1)

---

## Multi-Theme Architecture & Available Ecosystems

The variable system natively supports dynamic re-theming across multiple complete design languages:
1. **Navy Blue & Purple:** Deep midnight navy base (`#0A0F1D`) with electric violet (`#A855F7`) and ice cyan (`#38BDF8`).
2. **VS Code Ecosystem:** Visual Studio Code Dark+, Tokyo Night Storm, One Dark Pro, GitHub Dark High Contrast, and Monokai Pro.
3. **Heatmap & Density Scales:** Continuous Plasma thermal, Diverging Red-Amber-Green, and GitHub-style activity matrices.
4. **Warm Editorial & Craft:** Warm-black dark base (`#0B0A09`), warm paper light base (`#FBF9F6`), and brand amber (`#FFAD01`).

See [16-theme-catalogue-and-palettes.md](./16-theme-catalogue-and-palettes.md) for full palette specs, and [17-theme-tokens.json](./17-theme-tokens.json) for raw machine-readable JSON tokens.

---

## Forbidden Patterns

| Pattern | Why Forbidden |
|---------|--------------|
| `color: #7c3aed` | Hardcoded hex — bypasses theme system |
| `bg-purple-600` | Tailwind literal color — not theme-aware |
| `rgb(124, 58, 237)` | RGB format — can't use opacity modifier |
| `hsl(252, 85%, 60%)` | Full hsl() — should be `hsl(var(--primary))` |
| `text-white` | Literal color — use `text-primary-foreground` |
| `bg-black` | Literal color — use `bg-background` in dark mode |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CSS Source File | `src/index.css` |
| Tailwind Config | `tailwind.config.ts` |
| Design Principles | [02-design-principles.md](./02-design-principles.md) |
| Theme Catalogue & Palettes | [16-theme-catalogue-and-palettes.md](./16-theme-catalogue-and-palettes.md) |
| Machine-Readable Theme Tokens | [17-theme-tokens.json](./17-theme-tokens.json) |
| Dark Mode & Materiality | [18-dark-mode-and-materiality.md](./18-dark-mode-and-materiality.md) |
| Code Block Tokens | [09-code-blocks.md](./09-code-blocks.md) |
