# White Blue Theme: 3-Format Color Palettes, Fluid Typography & Design Tokens

> **/goal** Provide the single source of truth for all color ramps, semantic tokens, section band overrides, typography scales, gradients, shadows, border radii, neumorphic surfaces, and LESS mixins in the White Blue Theme (`04-white-blue-theme`).
> **/learn** Master the mandatory 3-format color specification (`HEX`, `RGB`/`RGBA`, `HSL` + `OKLCH`), the 3-font typographic hierarchy (`Ubuntu`, `Poppins`, `JetBrains Mono`), and the CSS custom property / Tailwind v4 `@theme inline` architecture.

---

## 1. Mandatory 3-Format Color Architecture

Every color token in the White Blue Theme is defined at runtime using perceptual **OKLCH** (`oklch(L C H)`) inside `:root` and mapped into Tailwind v4 via `@theme inline`. To guarantee universal compatibility across Figma, CSS/LESS preprocessors, canvas shaders, and blind AI code generators, **every single color in this specification is documented in three standard formats side-by-side**:

1. **Format 1 — HEX:** Standard 6-digit sRGB hexadecimal (`#2563EB`).
2. **Format 2 — RGB / RGBA:** Standard integer sRGB triplet or alpha quad (`rgb(37, 99, 235)` / `rgba(37, 99, 235, 0.34)`).
3. **Format 3 — HSL:** Standard hue-saturation-lightness (`hsl(221, 83%, 53%)`).
4. **Runtime CSS Token — OKLCH:** Perceptually uniform lightness-chroma-hue (`oklch(0.546 0.215 262.9)`), enabling smooth `color-mix(in oklab, ...)` alpha and tint derivation without muddy midtones.

---

## 2. Primitive Color Ramps (3-Format Reference Tables)

### 2.1 Blue Ramp (Primary Action & Corporate Trust)

The Blue ramp spans from icy porcelain tint (`--blue-50`) through vibrant Electric Cobalt (`--blue-500`) to Deep Royal Navy (`--blue-900`):

| CSS Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL | Role & Usage |
|:---|:---|:---|:---|:---|:---|
| `--blue-50` | `oklch(0.972 0.014 259)` | `#F0F6FF` | `rgb(240, 246, 255)` | `hsl(216, 100%, 97%)` | Hero top wash gradient start, brand pill background (`bg-blue-50`) |
| `--blue-100` | `oklch(0.935 0.032 262)` | `#DEEAFF` | `rgb(222, 234, 255)` | `hsl(218, 100%, 94%)` | Light badge hover fill, soft icon tile ring |
| `--blue-200` | `oklch(0.874 0.061 263)` | `#C1D6FF` | `rgb(193, 214, 255)` | `hsl(220, 100%, 88%)` | Chrome gradient stop, subtle focus border, panel-wash top highlight |
| `--blue-300` | `oklch(0.777 0.104 264)` | `#95B6FA` | `rgb(149, 182, 250)` | `hsl(220, 91%, 78%)` | Chrome gradient end stop, secondary data visualization series |
| `--blue-400` | `oklch(0.662 0.164 264)` | `#5E8EF6` | `rgb(94, 142, 246)` | `hsl(221, 89%, 67%)` | Dark-mode focus ring (`--ring`), gradient text end stop, retro grid line tint |
| `--blue-500` | `oklch(0.546 0.215 262.9)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | **Primary Action (`--primary`, `--brand-secondary`)**, active toggles, links, pin #1 |
| `--blue-600` | `oklch(0.489 0.208 263.5)` | `#1B51D3` | `rgb(27, 81, 211)` | `hsl(222, 77%, 47%)` | Brand gradient start (`--gradient-brand`), pressed primary button state |
| `--blue-700` | `oklch(0.409 0.176 264.5)` | `#173DA7` | `rgb(23, 61, 167)` | `hsl(224, 76%, 37%)` | Brand pill foreground (`text-blue-700`), high-contrast link hover |
| `--blue-800` | `oklch(0.369 0.147 264.8)` | `#17368B` | `rgb(23, 54, 139)` | `hsl(224, 72%, 32%)` | Deep interactive border, dark-on-light badge text |
| `--blue-900` | `oklch(0.317 0.135 264.5)` | `#0D2975` | `rgb(13, 41, 117)` | `hsl(224, 80%, 25%)` | **Core Brand Anchor (`--brand-primary`)**, solid button fill, shadow tint base |

### 2.2 Violet Ramp (Tertiary Accent & Editorial Gradient Partner)

The Violet ramp partners with Blue in all directional brand gradients (`--gradient-accent`, `--gradient-brand`, `--gradient-text`):

| CSS Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL | Role & Usage |
|:---|:---|:---|:---|:---|:---|
| `--violet-50` | `oklch(0.972 0.017 300)` | `#F7F3FF` | `rgb(247, 243, 255)` | `hsl(260, 100%, 98%)` | Light band semantic `--accent` surface fill |
| `--violet-100` | `oklch(0.946 0.032 300)` | `#F0E9FF` | `rgb(240, 233, 255)` | `hsl(259, 100%, 96%)` | Subtle secondary chip background |
| `--violet-200` | `oklch(0.885 0.068 299)` | `#DFD0FF` | `rgb(223, 208, 255)` | `hsl(259, 100%, 91%)` | Decorative accent border on light cards |
| `--violet-300` | `oklch(0.784 0.128 298)` | `#C3A6FF` | `rgb(195, 166, 255)` | `hsl(260, 100%, 83%)` | Soft glow highlight on dark surfaces |
| `--violet-400` | `oklch(0.659 0.187 297.5)` | `#A173F4` | `rgb(161, 115, 244)` | `hsl(261, 85%, 70%)` | Secondary chart accent, dark-mode badge text |
| `--violet-500` | `oklch(0.532 0.253 296.9)` | `#822EE8` | `rgb(130, 46, 232)` | `hsl(267, 80%, 55%)` | **Tertiary Brand Accent (`--brand-tertiary`)**, gradient end stop, workflow pin #2 |
| `--violet-600` | `oklch(0.481 0.243 296.5)` | `#721CD2` | `rgb(114, 28, 210)` | `hsl(268, 76%, 47%)` | `--gradient-brand` midpoint (55%), hero radial bottom-left wash |
| `--violet-700` | `oklch(0.417 0.213 295.9)` | `#5C14AF` | `rgb(92, 20, 175)` | `hsl(268, 79%, 38%)` | Light band semantic `--accent-foreground` text |
| `--violet-800` | `oklch(0.362 0.187 295.5)` | `#4B0D92` | `rgb(75, 13, 146)` | `hsl(268, 84%, 31%)` | Deep violet structural accent |
| `--violet-900` | `oklch(0.257 0.153 294.5)` | `#2D0063` | `rgb(45, 0, 99)` | `hsl(267, 100%, 19%)` | Maximum depth violet anchor |

### 2.3 Cyan, Indigo & Teal Decorative Accents

Used for button shine sweeps (`.shine-sweep`), gradient rings (`.gradient-ring`), workflow board push-pins, and scroll-stack card hues:

| CSS Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL | Role & Usage |
|:---|:---|:---|:---|:---|:---|
| `--cyan-50` | `oklch(0.984 0.02 197)` | `#EBFEFE` | `rgb(235, 254, 254)` | `hsl(180, 90%, 96%)` | Pale aqua highlight tint |
| `--cyan-200` | `oklch(0.907 0.075 195)` | `#A4F1F0` | `rgb(164, 241, 240)` | `hsl(179, 73%, 79%)` | Soft cyan badge border |
| `--cyan-400` | `oklch(0.797 0.136 205.6)` | `#03D5E7` | `rgb(3, 213, 231)` | `hsl(185, 97%, 46%)` | **Brand Highlight (`--brand-highlight`)**, `.shine-sweep` 45deg sheen, workflow pin #3 |
| `--cyan-600` | `oklch(0.6 0.116 214)` | `#0091AA` | `rgb(0, 145, 170)` | `hsl(189, 100%, 33%)` | Deep cyan technical label |
| `--cyan-800` | `oklch(0.44 0.077 218)` | `#075C6F` | `rgb(7, 92, 111)` | `hsl(191, 88%, 23%)` | Dark cyan structural accent |
| `--indigo-500` | `oklch(0.545 0.216 277)` | `#5856E9` | `rgb(88, 86, 233)` | `hsl(241, 77%, 63%)` | Workflow sticky-note pin #4, scroll-stack card hue |
| `--teal-500` | `oklch(0.63 0.118 195)` | `#009F9F` | `rgb(0, 159, 159)` | `hsl(180, 100%, 31%)` | Workflow sticky-note pin #5, scroll-stack card hue |

### 2.4 Neutral Cool Slate Ink Ramp

All neutrals carry a subtle cool sapphire-slate undertone (`hue 248–266`) so grey borders and body text harmonize with the Blue/Violet brand palette:

| CSS Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL | Role & Usage |
|:---|:---|:---|:---|:---|:---|
| `--neutral-50` | `oklch(0.984 0.003 247.9)` | `#F8FAFC` | `rgb(248, 250, 252)` | `hsl(210, 40%, 98%)` | Subtlest table header / code well wash |
| `--neutral-100` | `oklch(0.966 0.006 258)` | `#F1F4F8` | `rgb(241, 244, 248)` | `hsl(214, 33%, 96%)` | Secondary input well background |
| `--neutral-200` | `oklch(0.924 0.011 265)` | `#E2E6ED` | `rgb(226, 230, 237)` | `hsl(218, 23%, 91%)` | **Light Hairline Border (`--border-light`, `--border`, `--input`)** |
| `--neutral-300` | `oklch(0.853 0.023 264)` | `#C7CFDE` | `rgb(199, 207, 222)` | `hsl(219, 26%, 83%)` | Disabled control border, divider rule |
| `--neutral-400` | `oklch(0.688 0.036 264)` | `#8F9BB2` | `rgb(143, 155, 178)` | `hsl(219, 19%, 63%)` | **Dark-Band Muted Text (`--on-dark-muted`)**, placeholder text |
| `--neutral-500` | `oklch(0.535 0.041 265)` | `#626D86` | `rgb(98, 109, 134)` | `hsl(222, 16%, 45%)` | Chrome gradient metallic midpoint, secondary icon tone |
| `--neutral-600` | `oklch(0.441 0.043 265.9)` | `#48536B` | `rgb(72, 83, 107)` | `hsl(221, 20%, 35%)` | **Soft Ink (`--ink-soft`, `--muted-foreground`)** for lead copy, captions, eyebrows |
| `--neutral-700` | `oklch(0.32 0.041 266)` | `#293248` | `rgb(41, 50, 72)` | `hsl(223, 27%, 22%)` | Subheading and table body ink |
| `--neutral-800` | `oklch(0.216 0.033 266)` | `#121929` | `rgb(18, 25, 41)` | `hsl(222, 39%, 12%)` | High-density UI control ink |
| `--neutral-900` | `oklch(0.163 0.023 265)` | `#090E18` | `rgb(9, 14, 24)` | `hsl(220, 45%, 6%)` | **Primary Ink (`--ink`, `--foreground`)** for headings and primary copy |

---

## 3. Core Brand, Surface, Band & Status Tokens (3-Format Reference)

| CSS Token | Alias / OKLCH Value | Format 1: HEX | Format 2: RGB / RGBA | Format 3: HSL / HSLA | Architectural Purpose |
|:---|:---|:---|:---|:---|:---|
| `--brand-primary` | `var(--blue-900)` / `oklch(0.317 0.135 264.5)` | `#0D2975` | `rgb(13, 41, 117)` | `hsl(224, 80%, 25%)` | Deep royal anchor for solid buttons and brand-tinted elevation shadows |
| `--brand-secondary` | `var(--blue-500)` / `oklch(0.546 0.215 262.9)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | Primary interactive blue for links, underlines, spotlights, and active states |
| `--brand-tertiary` | `var(--violet-500)` / `oklch(0.532 0.253 296.9)` | `#822EE8` | `rgb(130, 46, 232)` | `hsl(267, 80%, 55%)` | Electric violet partner in accent gradients and glow filters |
| `--brand-highlight` | `var(--cyan-400)` / `oklch(0.797 0.136 205.6)` | `#03D5E7` | `rgb(3, 213, 231)` | `hsl(185, 97%, 46%)` | Luminous cyan specular highlight inside `.shine-sweep` and `.gradient-ring` |
| `--ink` | `var(--neutral-900)` / `oklch(0.163 0.023 265)` | `#090E18` | `rgb(9, 14, 24)` | `hsl(220, 45%, 6%)` | Primary editorial ink on light and soft bands |
| `--ink-soft` | `var(--neutral-600)` / `oklch(0.441 0.043 265.9)` | `#48536B` | `rgb(72, 83, 107)` | `hsl(221, 20%, 35%)` | Supporting body text, descriptions, and metadata labels |
| `--paper` | `oklch(1 0 0)` | `#FFFFFF` | `rgb(255, 255, 255)` | `hsl(0, 0%, 100%)` | Default page canvas (`--background`) and elevated card surface (`--card`) |
| `--paper-muted` | `oklch(0.972 0.006 264.5)` | `#F4F6FA` | `rgb(244, 246, 250)` | `hsl(220, 37%, 97%)` | Muted pill/chip well (`--muted`, `--secondary`) on light bands |
| `--surface-soft` | `oklch(0.962 0.011 258)` | `#EEF3FA` | `rgb(238, 243, 250)` | `hsl(215, 55%, 96%)` | Tinted cool-mist canvas for `.band-soft` sections and `.site-footer` |
| `--surface-dark` | `oklch(0.129 0.019 264)` | `#04070F` | `rgb(4, 7, 15)` | `hsl(224, 58%, 4%)` | Deep obsidian canvas for `.band-dark` and `.dark` scopes |
| `--surface-dark-elevated` | `oklch(0.181 0.027 265)` | `#0C121E` | `rgb(12, 18, 30)` | `hsl(220, 43%, 8%)` | Card and popover surface inside `.band-dark` |
| `--surface-dark-raised` | `oklch(0.221 0.031 265)` | `#141B29` | `rgb(20, 27, 41)` | `hsl(220, 34%, 12%)` | Raised interactive row or hover well inside `.band-dark` |
| `--neu-surface` | `oklch(0.951 0.009 265)` | `#ECEFF5` | `rgb(236, 239, 245)` | `hsl(220, 31%, 94%)` | Tactile base surface for `.neu-raised`, `.neu-pressed`, and `.neu-flat` cards |
| `--hairline` | `oklch(1 0 0 / 8%)` | `#FFFFFF14` | `rgba(255, 255, 255, 0.08)` | `hsla(0, 0%, 100%, 0.08)` | Translucent structural rule and grid backdrop line |
| `--border-light` | `var(--neutral-200)` / `oklch(0.924 0.011 265)` | `#E2E6ED` | `rgb(226, 230, 237)` | `hsl(218, 23%, 91%)` | Crisp 1px card and section divider border on light/soft bands |
| `--border-dark` | `oklch(1 0 0 / 10%)` | `#FFFFFF1A` | `rgba(255, 255, 255, 0.10)` | `hsla(0, 0%, 100%, 0.10)` | Translucent 1px border inside `.band-dark` |
| `--on-dark` | `oklch(0.877 0.021 265)` | `#D0D7E5` | `rgb(208, 215, 229)` | `hsl(220, 29%, 86%)` | Primary readable foreground text inside `.band-dark` |
| `--on-dark-muted` | `oklch(0.688 0.036 264)` | `#8F9BB2` | `rgb(143, 155, 178)` | `hsl(219, 19%, 63%)` | Muted supporting text inside `.band-dark` |
| `--success` | `oklch(0.696 0.149 162.5)` | `#10B981` | `rgb(16, 185, 129)` | `hsl(160, 84%, 39%)` | Positive SLA indicator, verified checkmarks, live telemetry status |
| `--warning` | `oklch(0.769 0.163 70.1)` | `#F59E13` | `rgb(245, 158, 19)` | `hsl(37, 92%, 52%)` | Advisory badge, threshold alert |
| `--destructive` | `oklch(0.577 0.245 27.325)` | `#E7000B` | `rgb(231, 0, 11)` | `hsl(357, 100%, 45%)` | Form validation error, destructive action |

---

## 4. Semantic Theme Aliases Across All Four Band Scopes

Components never branch on light vs. dark mode in JSX; they consume semantic tokens (`bg-background`, `bg-card`, `text-foreground`, `text-muted-foreground`, `border-border`) which automatically rebind inside `.band-soft`, `.band-dark`, and `.band-void`:

| Semantic Token | Light Band (`:root`) `HEX` / `RGB` / `HSL` | Soft Band (`.band-soft`) `HEX` / `RGB` / `HSL` | Dark Band (`.band-dark`, `.dark`) `HEX` / `RGB` / `HSL` | Void Band (`.band-void`) `HEX` / `RGB` / `HSL` |
|:---|:---|:---|:---|:---|
| `--background` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#EEF3FA` / `rgb(238, 243, 250)` / `hsl(215, 55%, 96%)` | `#04070F` / `rgb(4, 7, 15)` / `hsl(224, 58%, 4%)` | `#010204` / `rgb(1, 2, 4)` / `hsl(220, 60%, 1%)` |
| `--foreground` | `#090E18` / `rgb(9, 14, 24)` / `hsl(220, 45%, 6%)` | `#090E18` / `rgb(9, 14, 24)` / `hsl(220, 45%, 6%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` |
| `--card` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#0C121E` / `rgb(12, 18, 30)` / `hsl(220, 43%, 8%)` | `#020408` / `rgb(2, 4, 8)` / `hsl(220, 60%, 2%)` |
| `--card-foreground` | `#090E18` / `rgb(9, 14, 24)` / `hsl(220, 45%, 6%)` | `#090E18` / `rgb(9, 14, 24)` / `hsl(220, 45%, 6%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` |
| `--primary` | `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` | `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` | `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` | `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` |
| `--primary-foreground` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` | `#FFFFFF` / `rgb(255, 255, 255)` / `hsl(0, 0%, 100%)` |
| `--secondary` | `#F4F6FA` / `rgb(244, 246, 250)` / `hsl(220, 37%, 97%)` | `#F4F6FA` / `rgb(244, 246, 250)` / `hsl(220, 37%, 97%)` | `#0C121E` / `rgb(12, 18, 30)` / `hsl(220, 43%, 8%)` | `#0C121E` / `rgb(12, 18, 30)` / `hsl(220, 43%, 8%)` |
| `--secondary-foreground` | `#0D2975` / `rgb(13, 41, 117)` / `hsl(224, 80%, 25%)` | `#0D2975` / `rgb(13, 41, 117)` / `hsl(224, 80%, 25%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` |
| `--muted` | `#F4F6FA` / `rgb(244, 246, 250)` / `hsl(220, 37%, 97%)` | `#F5F8FC` / `rgb(245, 248, 252)` / `hsl(214, 54%, 97%)` | `#0C121E` / `rgb(12, 18, 30)` / `hsl(220, 43%, 8%)` | `#03050A` / `rgb(3, 5, 10)` / `hsl(223, 54%, 3%)` |
| `--muted-foreground` | `#48536B` / `rgb(72, 83, 107)` / `hsl(221, 20%, 35%)` | `#48536B` / `rgb(72, 83, 107)` / `hsl(221, 20%, 35%)` | `#8F9BB2` / `rgb(143, 155, 178)` / `hsl(219, 19%, 63%)` | `#8F9BB2` / `rgb(143, 155, 178)` / `hsl(219, 19%, 63%)` |
| `--accent` | `#F7F3FF` / `rgb(247, 243, 255)` / `hsl(260, 100%, 98%)` | `#F7F3FF` / `rgb(247, 243, 255)` / `hsl(260, 100%, 98%)` | `#0C121E` / `rgb(12, 18, 30)` / `hsl(220, 43%, 8%)` | `#0C121E` / `rgb(12, 18, 30)` / `hsl(220, 43%, 8%)` |
| `--accent-foreground` | `#5C14AF` / `rgb(92, 20, 175)` / `hsl(268, 79%, 38%)` | `#5C14AF` / `rgb(92, 20, 175)` / `hsl(268, 79%, 38%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` | `#D0D7E5` / `rgb(208, 215, 229)` / `hsl(220, 29%, 86%)` |
| `--border` / `--input` | `#E2E6ED` / `rgb(226, 230, 237)` / `hsl(218, 23%, 91%)` | `#E2E6ED` / `rgb(226, 230, 237)` / `hsl(218, 23%, 91%)` | `#FFFFFF1A` / `rgba(255, 255, 255, 0.10)` / `hsla(0, 0%, 100%, 0.10)` | `#FFFFFF1A` / `rgba(255, 255, 255, 0.10)` / `hsla(0, 0%, 100%, 0.10)` |
| `--ring` | `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` | `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` | `#5E8EF6` / `rgb(94, 142, 246)` / `hsl(221, 89%, 67%)` | `#5E8EF6` / `rgb(94, 142, 246)` / `hsl(221, 89%, 67%)` |

---

## 5. Typography & Three-Font Editorial System

### 5.1 Font Stack & Google Fonts Loading Contract

> [!IMPORTANT]
> **Explicit Ban on `Inter`:** The White Blue Theme strictly prohibits `Inter`, `Roboto`, or generic system defaults as primary fonts. Every page MUST load the canonical 3-font stack (`Ubuntu`, `Poppins`, `JetBrains Mono`).

| CSS Variable | Tailwind Utility | Font Family | Weights Loaded | Assigned UI Roles |
|:---|:---|:---|:---|:---|
| `--font-display` | `font-display` | `"Ubuntu", system-ui, sans-serif` | `400, 500, 700` | All headings (`h1`–`h6`), hero display copy, buttons (`WhiteBlueButton`), nav links, stat counters (`text-stat`, `text-numeral`), wordmarks |
| `--font-body` | `font-body` | `"Poppins", system-ui, sans-serif` | `400, 500, 600` | Body paragraphs, lead copy (`text-lead`), card descriptions, form inputs, table cells |
| `--font-mono` | `font-mono` | `"JetBrains Mono", ui-monospace, monospace` | `400, 500, 600` | Section eyebrows (`text-eyebrow`), pill badges (`Pill`), numbered row indices (`01`..`06`), telemetry counters (`3/5 enabled`), metadata kickers |

Include this exact `<link>` block inside `<head>`:

```html
<link rel="preconnect" href="https://fonts.googleapis.com" />
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
<link
  rel="stylesheet"
  href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Poppins:wght@400;500;600&family=Ubuntu:wght@400;500;700&display=swap"
/>
```

### 5.2 Fluid Typography `@utility` Scale (10 Canonical Classes)

| Utility Class | Font Token | Fluid Size (`clamp()`) | Weight | Letter Spacing | Line Height | Extra Properties |
|:---|:---|:---|:---|:---|:---|:---|
| `text-mega` | `var(--font-display)` | `clamp(44px, 9vw, 120px)` | `700` | `-0.032em` | `0.94` | Hero / monumental display statements |
| `text-h1` | `var(--font-display)` | `clamp(38px, 6vw, 80px)` | `700` | `-0.03em` | `1.0` | Page-level primary headlines |
| `text-h2` | `var(--font-display)` | `clamp(30px, 4.4vw, 60px)` | `700` | `-0.03em` | `1.06` | Section headlines (`SectionHeading`, `MaskedHeading`) |
| `text-h3` | `var(--font-display)` | `clamp(20px, 1.8vw, 26px)` | `700` | `-0.015em` | `1.2` | Flagship card titles, sub-band headers |
| `text-h4` | `var(--font-display)` | `20px` (`1.25rem`) | `500` | `-0.01em` | `1.3` | Compact card titles, typographic row headings |
| `text-lead` | `var(--font-body)` | `clamp(16px, 1.35vw, 20px)` | `400` | `normal` | `1.65` | Section lead paragraphs (`max-w-2xl text-muted-foreground`) |
| `text-eyebrow` | `var(--font-mono)` | `12px` (`0.75rem`) | `500` | `0.14em` | `1.0` | `text-transform: uppercase`, paired with `w-6 h-px` leading rule |
| `text-stat` | `var(--font-display)` | `clamp(44px, 7vw, 96px)` | `700` | `-0.038em` | `0.9` | `font-variant-numeric: tabular-nums` for animated KPI counters |
| `text-numeral` | `var(--font-display)` | `clamp(64px, 10vw, 160px)` | `700` | `-0.042em` | `0.8` | Oversized architectural index numerals |
| `text-wordmark` | `var(--font-display)` | `clamp(72px, 18vw, 280px)` | `700` | `-0.038em` | `0.8` | Monumental footer or backdrop watermarks |

---

## 6. Gradients, Brand-Tinted Shadows, Radii, Layout & Neumorphism

### 6.1 Multi-Stop Brand Gradients (with 3-Format Color Breakdown)

| Gradient Token | CSS Definition | Component Color Stops (`HEX` / `RGB` / `HSL`) | Primary Consumers |
|:---|:---|:---|:---|
| `--gradient-accent` | `linear-gradient(90deg, var(--blue-500) 0%, var(--violet-500) 100%)` | Start: `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)`<br>End: `#822EE8` / `rgb(130, 46, 232)` / `hsl(267, 80%, 55%)` | `WhiteBlueButton` (`primary`), `.row-premium::before`, progress bars, 3D flip promo card front |
| `--gradient-accent-hover` | `linear-gradient(90deg, color-mix(in oklab, var(--blue-500) 88%, white) 0%, color-mix(in oklab, var(--violet-500) 88%, white) 100%)` | Start: `#4277EE` / `rgb(66, 119, 238)` / `hsl(222, 83%, 60%)`<br>End: `#924AEB` / `rgb(146, 74, 235)` / `hsl(267, 80%, 61%)` | Hover state for `primary` and `solid` buttons |
| `--gradient-brand` | `linear-gradient(120deg, var(--blue-600) 0%, var(--violet-600) 55%, var(--blue-500) 100%)` | `0%`: `#1B51D3` / `rgb(27, 81, 211)` / `hsl(222, 77%, 47%)`<br>`55%`: `#721CD2` / `rgb(114, 28, 210)` / `hsl(268, 76%, 47%)`<br>`100%`: `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)` | `.card-premium::before` 2px left edge growth, default `.pointer-fill` |
| `--gradient-text` | `linear-gradient(92deg, var(--blue-500) 0%, var(--violet-500) 55%, var(--blue-400) 100%)` | `0%`: `#2563EB` / `rgb(37, 99, 235)` / `hsl(221, 83%, 53%)`<br>`55%`: `#822EE8` / `rgb(130, 46, 232)` / `hsl(267, 80%, 55%)`<br>`100%`: `#5E8EF6` / `rgb(94, 142, 246)` / `hsl(221, 89%, 67%)` | `.gradient-text` accent words inside `MaskedHeading` |
| `--gradient-halo` | `radial-gradient(55% 60% at 72% 12%, color-mix(in oklab, var(--blue-400) 26%, transparent) 0%, color-mix(in oklab, var(--violet-500) 10%, transparent) 45%, transparent 72%)` | Core: `rgba(94, 142, 246, 0.26)` (`#5E8EF642` / `hsla(221, 89%, 67%, 0.26)`)<br>Mid: `rgba(130, 46, 232, 0.10)` (`#822EE81A` / `hsla(267, 80%, 55%, 0.10)`) | `Halo` backdrop behind Hero and Flagship `SpotlightCard` |
| `--gradient-border` | `linear-gradient(120deg, color-mix(in oklab, var(--blue-500) 60%, transparent), color-mix(in oklab, var(--violet-500) 35%, transparent), color-mix(in oklab, var(--cyan-400) 40%, transparent))` | `rgba(37, 99, 235, 0.60)` -> `rgba(130, 46, 232, 0.35)` -> `rgba(3, 213, 231, 0.40)` | `.gradient-ring` 1px masked border on dark glass cards and pills |

### 6.2 Brand-Tinted Elevation Shadows

Shadows in the White Blue Theme are never pure grey/black on light bands; they are tinted with `--brand-primary` (`#0D2975` / `rgb(13, 41, 117)` / `hsl(224, 80%, 25%)`) so cards cast a crisp sapphire-ink ambient shadow:

| Shadow Token | CSS Value | 3-Format Shadow Color Equivalent | Usage |
|:---|:---|:---|:---|
| `--shadow-xs` | `0 1px 2px color-mix(in oklab, var(--ink) 5%, transparent)` | `#090E180D` / `rgba(9, 14, 24, 0.05)` / `hsla(220, 45%, 6%, 0.05)` | Toggle switches, small index badges |
| `--shadow-card` | `0 8px 24px -12px color-mix(in oklab, var(--brand-primary) 18%, transparent)` | `#0D29752E` / `rgba(13, 41, 117, 0.18)` / `hsla(224, 80%, 25%, 0.18)` | Resting elevation for `SurfaceCard`, `card-premium`, scrolled `SiteHeader` |
| `--shadow-lift` | `0 24px 60px -24px color-mix(in oklab, var(--brand-primary) 35%, transparent)` | `#0D297559` / `rgba(13, 41, 117, 0.35)` / `hsla(224, 80%, 25%, 0.35)` | Hover elevation for buttons, Mega-Menu panel, Hero `CapabilityStack`, `WorkflowNote` |
| `--shadow-glow` | `0 0 60px -12px color-mix(in oklab, var(--blue-500) 55%, transparent)` | `#2563EB8C` / `rgba(37, 99, 235, 0.55)` / `hsla(221, 83%, 53%, 0.55)` | High-emphasis focal halo |
| `--shadow-dark-card` | `0 24px 60px -30px oklch(0 0 0 / 85%)` | `#000000D9` / `rgba(0, 0, 0, 0.85)` / `hsla(0, 0%, 0%, 0.85)` | `GlassCard` hover elevation on dark bands |

### 6.3 Radii, Motion Durations, Layout & Neumorphism Tokens

| Category | Token | Value | Purpose |
|:---|:---|:---|:---|
| **Border Radii** | `--radius` | `14px` | Base radius (`--radius-sm: 10px`, `--radius-md: 12px`, `--radius-lg: 14px`, `--radius-xl: 18px`) |
| | `--radius-button` | `12px` | Rounded rectangle curvature on all `WhiteBlueButton` sizes |
| | `--radius-card` | `20px` | Standard card curvature (`SurfaceCard`, `card-premium`, `MegaPanel`, `ToggleRow`) |
| | `--radius-media` | `28px` | Outer stage containers (`CapabilityStack`, `stack-panel`, Pricing quote card) |
| | `--radius-neu` | `18px` | Neumorphic surface curvature (`NeuCard`) |
| **Motion Timing** | `--dur-instant` | `120ms` | Instant feedback & reduced-motion fallback duration |
| | `--dur-fast` | `240ms` | Button hover, chevron rotation, pointer-fill, underline scaleX |
| | `--dur-base` | `420ms` | Card border/background transitions, header scroll elevation, spotlight fade |
| | `--dur-slow` | `700ms` | Standard scroll reveal (`fadeUp`) |
| | `--dur-cine` | `1100ms` | Word/line `maskUp` headline reveal |
| | `--ease-out` | `cubic-bezier(0.16, 1, 0.3, 1)` | Primary exponential deceleration curve |
| | `--ease-in-out` | `cubic-bezier(0.65, 0, 0.35, 1)` | Symmetric oscillation curve |
| **Layout & Grid** | `--content-max` | `1280px` | Max container width (`Container`) |
| | `--gutter` | `20px` (`<400px`), `24px` (default), `40px` (`>=768px`), `64px` (`>=1280px`) | Responsive horizontal container padding |
| | `--section-tight` | `clamp(40px, 5vw, 64px)` | Compact section vertical padding |
| | `--section-base` | `clamp(72px, 8vw, 120px)` | Standard section vertical padding |
| | `--section-hero-top` | `clamp(48px, 6vw, 96px)` | Top hero clearance below sticky header |
| **Neumorphism** | `--neu-raised` | `-6px -6px 14px rgba(255,255,255,0.95), 7px 7px 18px rgba(13,41,117,0.14)` | Extruded soft-clay card shadow on `--neu-surface` (`#ECEFF5`) |
| | `--neu-pressed` | `inset -4px -4px 10px rgba(255,255,255,0.95), inset 5px 5px 12px rgba(13,41,117,0.14)` | Recessed well shadow on `--neu-surface` (`#ECEFF5`) |
| | `--neu-flat` | `-2px -2px 6px rgba(255,255,255,0.95), 3px 3px 8px rgba(13,41,117,0.14)` | Low-profile tactile tile shadow |

---

## 7. Complete Copy-Pasteable Root CSS & LESS Token Sheet

### 7.1 Canonical `styles.css` Root & `@theme inline` Block

```css
@import "tailwindcss" source(none);
@source "../src";
@import "tw-animate-css";

@custom-variant dark (&:is(.dark *));

:root {
  /* Typography */
  --font-display: "Ubuntu", system-ui, sans-serif;
  --font-body: "Poppins", system-ui, sans-serif;
  --font-mono: "JetBrains Mono", ui-monospace, monospace;

  /* Blue ramp (secondary / action) */
  --blue-50: oklch(0.972 0.014 259);      /* #F0F6FF | rgb(240, 246, 255) | hsl(216, 100%, 97%) */
  --blue-100: oklch(0.935 0.032 262);     /* #DEEAFF | rgb(222, 234, 255) | hsl(218, 100%, 94%) */
  --blue-200: oklch(0.874 0.061 263);     /* #C1D6FF | rgb(193, 214, 255) | hsl(220, 100%, 88%) */
  --blue-300: oklch(0.777 0.104 264);     /* #95B6FA | rgb(149, 182, 250) | hsl(220, 91%, 78%)  */
  --blue-400: oklch(0.662 0.164 264);     /* #5E8EF6 | rgb(94, 142, 246)  | hsl(221, 89%, 67%)  */
  --blue-500: oklch(0.546 0.215 262.9);   /* #2563EB | rgb(37, 99, 235)   | hsl(221, 83%, 53%)  */
  --blue-600: oklch(0.489 0.208 263.5);   /* #1B51D3 | rgb(27, 81, 211)   | hsl(222, 77%, 47%)  */
  --blue-700: oklch(0.409 0.176 264.5);   /* #173DA7 | rgb(23, 61, 167)   | hsl(224, 76%, 37%)  */
  --blue-800: oklch(0.369 0.147 264.8);   /* #17368B | rgb(23, 54, 139)   | hsl(224, 72%, 32%)  */
  --blue-900: oklch(0.317 0.135 264.5);   /* #0D2975 | rgb(13, 41, 117)   | hsl(224, 80%, 25%)  */

  /* Violet ramp (tertiary / accent) */
  --violet-50: oklch(0.972 0.017 300);    /* #F7F3FF | rgb(247, 243, 255) | hsl(260, 100%, 98%) */
  --violet-100: oklch(0.946 0.032 300);   /* #F0E9FF | rgb(240, 233, 255) | hsl(259, 100%, 96%) */
  --violet-200: oklch(0.885 0.068 299);   /* #DFD0FF | rgb(223, 208, 255) | hsl(259, 100%, 91%) */
  --violet-300: oklch(0.784 0.128 298);   /* #C3A6FF | rgb(195, 166, 255) | hsl(260, 100%, 83%) */
  --violet-400: oklch(0.659 0.187 297.5); /* #A173F4 | rgb(161, 115, 244) | hsl(261, 85%, 70%)  */
  --violet-500: oklch(0.532 0.253 296.9); /* #822EE8 | rgb(130, 46, 232)  | hsl(267, 80%, 55%)  */
  --violet-600: oklch(0.481 0.243 296.5); /* #721CD2 | rgb(114, 28, 210)  | hsl(268, 76%, 47%)  */
  --violet-700: oklch(0.417 0.213 295.9); /* #5C14AF | rgb(92, 20, 175)   | hsl(268, 79%, 38%)  */
  --violet-800: oklch(0.362 0.187 295.5); /* #4B0D92 | rgb(75, 13, 146)   | hsl(268, 84%, 31%)  */
  --violet-900: oklch(0.257 0.153 294.5); /* #2D0063 | rgb(45, 0, 99)     | hsl(267, 100%, 19%) */

  /* Cyan, Indigo & Teal accents */
  --cyan-50: oklch(0.984 0.02 197);       /* #EBFEFE | rgb(235, 254, 254) | hsl(180, 90%, 96%)  */
  --cyan-200: oklch(0.907 0.075 195);     /* #A4F1F0 | rgb(164, 241, 240) | hsl(179, 73%, 79%)  */
  --cyan-400: oklch(0.797 0.136 205.6);   /* #03D5E7 | rgb(3, 213, 231)   | hsl(185, 97%, 46%)  */
  --cyan-600: oklch(0.6 0.116 214);       /* #0091AA | rgb(0, 145, 170)   | hsl(189, 100%, 33%) */
  --cyan-800: oklch(0.44 0.077 218);      /* #075C6F | rgb(7, 92, 111)    | hsl(191, 88%, 23%)  */
  --indigo-500: oklch(0.545 0.216 277);   /* #5856E9 | rgb(88, 86, 233)   | hsl(241, 77%, 63%)  */
  --teal-500: oklch(0.63 0.118 195);      /* #009F9F | rgb(0, 159, 159)   | hsl(180, 100%, 31%) */

  /* Neutral cool slate ink ramp */
  --neutral-50: oklch(0.984 0.003 247.9); /* #F8FAFC | rgb(248, 250, 252) | hsl(210, 40%, 98%)  */
  --neutral-100: oklch(0.966 0.006 258);  /* #F1F4F8 | rgb(241, 244, 248) | hsl(214, 33%, 96%)  */
  --neutral-200: oklch(0.924 0.011 265);  /* #E2E6ED | rgb(226, 230, 237) | hsl(218, 23%, 91%)  */
  --neutral-300: oklch(0.853 0.023 264);  /* #C7CFDE | rgb(199, 207, 222) | hsl(219, 26%, 83%)  */
  --neutral-400: oklch(0.688 0.036 264);  /* #8F9BB2 | rgb(143, 155, 178) | hsl(219, 19%, 63%)  */
  --neutral-500: oklch(0.535 0.041 265);  /* #626D86 | rgb(98, 109, 134)  | hsl(222, 16%, 45%)  */
  --neutral-600: oklch(0.441 0.043 265.9);/* #48536B | rgb(72, 83, 107)   | hsl(221, 20%, 35%)  */
  --neutral-700: oklch(0.32 0.041 266);   /* #293248 | rgb(41, 50, 72)    | hsl(223, 27%, 22%)  */
  --neutral-800: oklch(0.216 0.033 266);  /* #121929 | rgb(18, 25, 41)    | hsl(222, 39%, 12%)  */
  --neutral-900: oklch(0.163 0.023 265);  /* #090E18 | rgb(9, 14, 24)     | hsl(220, 45%, 6%)   */

  /* Core brand & surface tokens */
  --brand-primary: var(--blue-900);
  --brand-secondary: var(--blue-500);
  --brand-tertiary: var(--violet-500);
  --brand-highlight: var(--cyan-400);
  --ink: var(--neutral-900);
  --ink-soft: var(--neutral-600);
  --paper: oklch(1 0 0);                        /* #FFFFFF | rgb(255, 255, 255) | hsl(0, 0%, 100%)   */
  --paper-muted: oklch(0.972 0.006 264.5);      /* #F4F6FA | rgb(244, 246, 250) | hsl(220, 37%, 97%) */
  --surface-soft: oklch(0.962 0.011 258);       /* #EEF3FA | rgb(238, 243, 250) | hsl(215, 55%, 96%) */
  --surface-dark: oklch(0.129 0.019 264);       /* #04070F | rgb(4, 7, 15)      | hsl(224, 58%, 4%)  */
  --surface-dark-elevated: oklch(0.181 0.027 265); /* #0C121E | rgb(12, 18, 30) | hsl(220, 43%, 8%)  */
  --surface-dark-raised: oklch(0.221 0.031 265);   /* #141B29 | rgb(20, 27, 41) | hsl(220, 34%, 12%) */
  --hairline: oklch(1 0 0 / 8%);                /* #FFFFFF14 | rgba(255, 255, 255, 0.08) */
  --border-light: var(--neutral-200);
  --border-dark: oklch(1 0 0 / 10%);            /* #FFFFFF1A | rgba(255, 255, 255, 0.10) */
  --on-dark: oklch(0.877 0.021 265);            /* #D0D7E5 | rgb(208, 215, 229) | hsl(220, 29%, 86%) */
  --on-dark-muted: oklch(0.688 0.036 264);      /* #8F9BB2 | rgb(143, 155, 178) | hsl(219, 19%, 63%) */
  --success: oklch(0.696 0.149 162.5);          /* #10B981 | rgb(16, 185, 129)  | hsl(160, 84%, 39%) */
  --warning: oklch(0.769 0.163 70.1);           /* #F59E13 | rgb(245, 158, 19)  | hsl(37, 92%, 52%)  */
  --destructive: oklch(0.577 0.245 27.325);     /* #E7000B | rgb(231, 0, 11)    | hsl(357, 100%, 45%)*/

  /* Semantic light-band defaults */
  --background: var(--paper);
  --foreground: var(--ink);
  --card: var(--paper);
  --card-foreground: var(--ink);
  --popover: var(--paper);
  --popover-foreground: var(--ink);
  --primary: var(--blue-500);
  --primary-foreground: oklch(1 0 0);
  --secondary: var(--paper-muted);
  --secondary-foreground: var(--brand-primary);
  --muted: var(--paper-muted);
  --muted-foreground: var(--ink-soft);
  --accent: var(--violet-50);
  --accent-foreground: var(--violet-700);
  --border: var(--border-light);
  --input: var(--border-light);
  --ring: var(--blue-500);
}
```

### 7.2 Parametric LESS Token & Typography Mixins (Preferred Preprocessor Pattern)

```less
// ============================================================================
// WHITE BLUE THEME — PARAMETRIC LESS TOKENS & MIXINS
// ============================================================================

// Core 3-Format Brand Variables
@wb-brand-primary:    #0D2975; // rgb(13, 41, 117)   | hsl(224, 80%, 25%)
@wb-brand-secondary:  #2563EB; // rgb(37, 99, 235)   | hsl(221, 83%, 53%)
@wb-brand-tertiary:   #822EE8; // rgb(130, 46, 232)  | hsl(267, 80%, 55%)
@wb-brand-highlight:  #03D5E7; // rgb(3, 213, 231)   | hsl(185, 97%, 46%)

@wb-ink:              #090E18; // rgb(9, 14, 24)     | hsl(220, 45%, 6%)
@wb-ink-soft:         #48536B; // rgb(72, 83, 107)   | hsl(221, 20%, 35%)
@wb-paper:            #FFFFFF; // rgb(255, 255, 255) | hsl(0, 0%, 100%)
@wb-paper-muted:      #F4F6FA; // rgb(244, 246, 250) | hsl(220, 37%, 97%)
@wb-surface-soft:     #EEF3FA; // rgb(238, 243, 250) | hsl(215, 55%, 96%)
@wb-surface-dark:     #04070F; // rgb(4, 7, 15)      | hsl(224, 58%, 4%)
@wb-border-light:     #E2E6ED; // rgb(226, 230, 237) | hsl(218, 23%, 91%)

// Typography Mixins
.wb-heading-h2() {
  font-family: "Ubuntu", system-ui, sans-serif;
  font-size: clamp(30px, 4.4vw, 60px);
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.06;
  color: @wb-ink;
}

.wb-eyebrow() {
  font-family: "JetBrains Mono", ui-monospace, monospace;
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  line-height: 1;
  color: @wb-ink-soft;
}

// Brand-Tinted Shadow Mixins
.wb-shadow-card() {
  box-shadow: 0 8px 24px -12px fade(@wb-brand-primary, 18%);
}

.wb-shadow-lift() {
  box-shadow: 0 24px 60px -24px fade(@wb-brand-primary, 35%);
}
```
