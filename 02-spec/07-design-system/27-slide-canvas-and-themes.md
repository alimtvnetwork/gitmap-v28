# 27 — Slide Canvas, 10 Production Themes & Contrast Inversion

> **/goal** Specify the authoritative 16:9 virtual stage scaling mathematics, the complete 10-theme production catalog, the Light Theme Contrast Inversion contract, and the Theme Menu popover system.
> **/learn** Master the 1920×1080 stage scale formula, the 10 built-in theme definitions (`bright-gold` default, `noir-gold`, `vscode-dark`, `dracula`, `monokai`, `github-light`, `paper-ink`, `macos-sonoma`, `windows-11`, `navy-blue`), the Light Theme Capsule Contract, and Theme Manifest import/export.

**Version:** 4.2.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. 16:9 Virtual Canvas & Coordinate Scaling

Every slide in the deck is authored on an authoritative reference canvas of `1920 × 1080` pixels (16:9 aspect ratio):

```text
scale = min(viewportWidth / 1920, viewportHeight / 1080)
transform-origin: center center
transform: scale(scale)
```

- **Runtime Architecture:** `ScaledSlide` detects container dimensions using `ResizeObserver`, applies CSS vector scaling, and centers the stage with letterboxing on non-16:9 viewports.
- **Stage Container:** `contain: layout paint; isolation: isolate; will-change: transform; width: 1920px; height: 1080px;`.
- **Absolute Coordinate Mandate:** All typography, margins, paddings, and card dimensions MUST be specified in authoring pixels on the `1920×1080` grid. Never use viewport units (`vw`, `vh`) inside the stage.

---

## 2. Authoritative 10 Production Presentation Themes

The catalog consists of **10 production themes**:

| ID | Theme Name | Appearance | Primary Accent | Cream / Text | Background | Primary Mood & Application |
|:---|:---|:---:|:---|:---|:---|:---|
| **`bright-gold`** *(Default)* | Bright Gold | Dark | Vivid Gold `40 96% 48%` (`#F3A502`) | Cream `42 100% 94%` (`#FFF1D6`) | Obsidian `0 0% 5%` (`#0D0D0D`) | Default keynote authority, executive pitch decks |
| **`noir-gold`** | Noir & Gold | Dark | Muted Gold `43 56% 54%` (`#C9A84C`) | Warm Cream `40 75% 84%` (`#F0D78C`) | Obsidian `0 0% 5%` (`#0D0D0D`) | Classic high-end luxury, board briefings |
| **`vscode-dark`** | VS Code Dark+ | Dark | Azure Blue `207 100% 50%` (`#007ACC`) | Crisp Gray `0 0% 83%` (`#D4D4D4`) | Slate `#1E1E1E` | Code-heavy tech talks, developer tooling |
| **`dracula`** | Dracula Gothic | Dark | Electric Purple `265 89% 78%` (`#BD93F9`) | Pure Cream `60 30% 96%` (`#F8F8F2`) | Charcoal `#282A36` | Aesthetic dev conferences, AI demos |
| **`monokai`** | Monokai Vibrant| Dark | Neon Green `80 76% 53%` (`#A6E22E`) | Light Cream `60 30% 96%` (`#F8F8F2`) | Deep Ink `#272822` | High-energy technical architecture, CLI keynotes |
| **`github-light`** | GitHub Light | **Light** | Open Blue `212 92% 45%` (`#0969DA`) | Espresso Ink `210 12% 16%` (`#24292F`) | Pure White `#FFFFFF` | Daytime presentations, public documentation |
| **`paper-ink`** | Paper & Ink | **Light** | Deep Amber `38 80% 30%` (`#8A5A0E`) | Espresso Ink `36 25% 12%` (`#1F1A12`) | Warm Cream `#FAF6EC` | Academic presentations, print handouts, research |
| **`macos-sonoma`** | macOS Sonoma | Dark | System Blue `212 100% 50%` (`#007AFF`) | Soft White `240 7% 97%` (`#F5F5F7`) | Dark Glass `#1E1E24` | Product design keynotes, client showcases |
| **`windows-11`** | Windows Fluent | Dark | Cyan Accent `199 100% 69%` (`#60CDFF`) | Crisp White `0 0% 100%` (`#FFFFFF`) | Mica Dark `#202020` | Enterprise platform migrations, IT briefings |
| **`navy-blue`** | Deep Navy Tech | Dark | Electric Cyan `188 95% 43%` (`#06B6D4`) | Crisp Slate `210 40% 96%` (`#F1F5F9`) | Deep Navy `#1A2840` | Cloud telemetry, infrastructure, bike showcases |

`DEFAULT_THEME` constant across all engines: **`'bright-gold'`**.

---

## 3. Light Theme Contrast Inversion Contract & Capsule Rules

When switching to light appearance (`github-light` or `paper-ink`), the color contrast budget must invert completely:

### 3.1 Token Collisions on Light Themes
Tokens that change meaning between dark and light themes:

| Token | Dark Themes | Light Themes (`paper-ink`, `github-light`) |
|---|---|---|
| `--ink` | Dark surface plate (bg) | Dark body text (fg) |
| `--cream` | Warm light text | **Repurposed → Dark espresso ink** |
| `--white` | Pure white text | **Repurposed → Dark espresso ink** |
| `--gold` | Bright accent (L=48%) | Darkened accent (L=30%) for AA contrast on cream |
| `--ember` | Warm coral (L=57%) | Darkened rust (L=45%) |

### 3.2 The Light-Theme Capsule Contract (Total Ban on Inline Styles)
> **Capsules MUST be painted via the `.capsule-{tone}` className system.**
> Inline `style.background` / `style.color` on a chip or capsule is **strictly forbidden** because it bypasses per-theme CSS overrides and causes catastrophic contrast collapse (e.g. brown blob on cream or black pill with invisible text).

```css
/* Canonical className system with per-theme overrides in index.css */
.capsule-gold { background: hsl(var(--gold)); color: hsl(var(--ink)); }
.capsule-ember { background: hsl(var(--ember)); color: white; }
.capsule-cream { background: hsl(var(--cream)); color: hsl(var(--ink)); }
.capsule-meta { background: hsl(var(--meta-bg)); color: hsl(var(--meta-fg)); }

/* Light Theme Overrides */
[data-theme='paper-ink'] .capsule-gold { background: hsl(var(--gold)); color: white; }
[data-theme='paper-ink'] .capsule-ember { background: hsl(var(--ember)); color: white; }
[data-theme='paper-ink'] .capsule-cream { background: var(--capsule-cream-bg); color: var(--capsule-cream-fg); }
```

### 3.3 Audit Grep Command for AI
Run before shipping slide changes:
```bash
rg -n "style=\{\{[^}]*hsl\(var\(--(gold|ember|cream|ink|white)" src/
```
Zero matches permitted on pill/chip elements.

---

## 4. Theme Menu Popover Component & Manifest Import/Export

The Theme Menu is anchored to the Palette button in the Controller HUD:
- **Upward Opening:** Anchored at bottom-right or top-right, opening upward or downward away from viewport edge.
- **Chrome Isolation:** Uses dedicated `--chrome-bg: 240 10% 6%` and `--chrome-fg: 0 0% 100%` tokens so the popover remains dark glass on all themes, including `paper-ink` and `github-light`.
- **Live Announcer:** Updates `liveMessage` for screen readers on selection (`setLiveMessage("Theme set to " + label)`).

### 4.1 Theme Manifest Import/Export Format
Custom themes export and import as portable JSON manifests:

```json
{
  "$schema": "https://specs.local/schemas/theme-manifest.v1.json",
  "id": "custom-emerald",
  "label": "Custom Emerald",
  "appearance": "dark",
  "description": "High-contrast clinical emerald theme",
  "swatch": ["#063729", "#1CC491", "#F5FEFA", "#2EEBA3"],
  "vars": {
    "--primary": "155 75% 44%",
    "--background": "165 80% 6%",
    "--foreground": "130 80% 98%",
    "--border": "158 40% 18%"
  }
}
```

- `buildThemeManifest(id)`: Constructs manifest object from active theme.
- `downloadThemeManifest(manifest)`: Triggers client-side browser JSON download.
- `parseThemeManifest(jsonString)`: Validates schema and imports preset into localStorage registry.

---
## 5. Theme Swatch Arrays

Each theme publishes a 4-color swatch array used by the picker popover:
- **`bright-gold` Swatch:** `['#0D0D0D', '#F3A502', '#FFF1D6', '#E85D3A']`
- **`noir-gold` Swatch:** `['#0D0D0D', '#C9A84C', '#F0D78C', '#E85D3A']`
- **`github-light` Swatch:** `['#FFFFFF', '#0969DA', '#24292F', '#CF222E']`
- **`paper-ink` Swatch:** `['#FAF6EC', '#8A5A0E', '#1F1A12', '#C04A24']`
- **`navy-blue` Swatch:** `['#1A2840', '#06B6D4', '#F1F5F9', '#F59E0B']`
- **`vscode-dark` Swatch:** `['#1E1E1E', '#0A84FF', '#D4D4D4', '#CE9178']`
- **`dracula` Swatch:** `['#282A36', '#BD93F9', '#F8F8F2', '#FF79C6']`
- **`monokai` Swatch:** `['#272822', '#A6E22E', '#F8F8F2', '#FD971F']`

---

## 6. Slide Theme Switching & Shared Variables

The runtime theme switch mechanisms, the 8-theme slide switch implementation, and the shared variable contract (`--canvas`, `--ink`, `--ink-muted`, `--accent`, `--accent-ink`, `--card`) are specified in [`40-theme-switch.md`](./40-theme-switch.md).

- Read [`40-theme-switch.md`](./40-theme-switch.md) before changing or switching runtime themes.
- Presentation themes are strictly scoped to the slide stage. Marketing website colors (navy `#0D2975`, cobalt `#2563EB`, or violet `#822EE8`) must NEVER be imported to color slide canvases, and slide presentation amber `#F3A502` must NEVER be used as the primary action color on marketing homepages.

---

## 7. Anti-Hallucination & Quality Verification Checklist

- [ ] Canvas math strictly enforces `1920×1080` authoring with uniform `min()` scale vector scaling.
- [ ] Theme binds to the defined production theme IDs (default `'bright-gold'` / `bright-gold-tech`).
- [ ] Light themes (`github-light`, `paper-ink`, `paper`, `print`) enforce the Light Theme Contract with inverted ink.
- [ ] Zero inline styles used for capsule/chip background or text colors.
- [ ] Floating controller HUD and Theme Menu preserve dark chrome tokens across all slide themes.
- [ ] Theme manifests conform to the JSON schema with 4-swatch definitions and valid HSL vars.

---

## 8. Sibling References

- Standalone marketing images, social cards, and thumbnail specs: [`37-image-specifications.md`](./37-image-specifications.md)
- Slide layouts and constraints: [`28-slide-layouts.md`](./28-slide-layouts.md)
- Master slide layout catalog: [`34-slide-layout-catalog.md`](./34-slide-layout-catalog.md)
- Runtime 8-theme switcher: [`40-theme-switch.md`](./40-theme-switch.md)
