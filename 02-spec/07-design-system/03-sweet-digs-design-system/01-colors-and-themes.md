# Sweet Digs Color & Theme Architecture

> **/goal** Define the centralized HSL color token architecture, theme variables, and multi-palette switching model for the Sweet Digs Design System.
> **/learn** Master the single source of truth for color tokens, positive state logic, and dark mode adaptations across all components.

---

## 1. HSL Token Architecture

All colors in the Sweet Digs Design System are defined in **HSL format** (hue saturation% lightness%) without the `hsl()` wrapper. This architecture enables flexible alpha channel compositing across all components:

```css
/* Alpha channel composition pattern */
color: hsl(var(--primary) / 0.8);
background-color: hsl(var(--secondary) / 0.5);
border-color: hsl(var(--primary) / 0.4);
```

Variables are mounted on the `:root` selector for the light theme and overridden under the `.dark` selector.

---

## 2. Light Theme Variables (Default Emerald Green)

The default brand color is **Emerald Green** (`142 71% 45%`, hex `#22c55e` equivalent), paired with soft mint tints and neutral surfaces:

```css
:root {
  /* ───── Brand Accent ───── */
  --primary:                142 71% 45%;      /* Emerald brand green */
  --primary-foreground:     0 0% 100%;        /* Crisp white text on primary */

  /* ───── Secondary / Tints ───── */
  --secondary:              142 60% 93%;      /* Subtle mint tint for chips, hovers */
  --secondary-foreground:   142 71% 30%;      /* Deep green text on secondary */

  /* ───── Surfaces & Elevation ───── */
  --background:             140 20% 97%;      /* Warm tinted page background */
  --foreground:             160 20% 10%;      /* High-contrast slate body text */
  --card:                   0 0% 100%;        /* Pure white elevated cards */
  --card-foreground:        160 20% 10%;      /* Text on card surfaces */
  --popover:                0 0% 100%;        /* Dropdowns and popovers */
  --popover-foreground:     160 20% 10%;      /* Text inside popovers */

  /* ───── Muted / Input Backgrounds ───── */
  --muted:                  150 14% 96%;      /* Soft background for inputs & chips */
  --muted-foreground:       160 9% 46%;       /* Subdued secondary descriptions */

  /* ───── Accent Elements ───── */
  --accent:                 142 60% 93%;      /* Accent highlight surface */
  --accent-foreground:      142 71% 30%;      /* Text on accent */

  /* ───── Destructive / Danger ───── */
  --destructive:            0 84% 60%;        /* System alerts and destructive actions */
  --destructive-foreground: 0 0% 100%;        /* Text on destructive */

  /* ───── Borders, Rings & Radius ───── */
  --border:                 150 13% 91%;      /* Subtle hairline borders */
  --input:                  150 13% 91%;      /* Form field borders */
  --ring:                   142 71% 45%;      /* Focus ring color (matches primary) */
  --radius:                 0.75rem;          /* 12px base radius */

  /* ───── Hero Ambient Gradients ───── */
  --gradient-start:         142 40% 92%;      /* Soft ambient hero top */
  --gradient-end:           160 30% 90%;      /* Ambient hero bottom */

  /* ───── Typography Stack ───── */
  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
```

---

## 3. Dark Theme Variables (`.dark`)

Dark mode preserves the vibrant `--primary` brand token while inverting surfaces into a deep obsidian slate:

```css
.dark {
  /* ───── Surfaces & Elevation ───── */
  --background:             160 20% 6%;       /* Deep obsidian slate */
  --foreground:             0 0% 95%;         /* Soft white readable text */
  --card:                   160 20% 10%;      /* Elevated dark card surface */
  --card-foreground:        0 0% 95%;         /* Text on dark card */
  --popover:                160 20% 10%;      /* Dark popover container */
  --popover-foreground:     0 0% 95%;         /* Popover text */

  /* ───── Brand Accent ───── */
  --primary:                142 71% 45%;      /* Brand emerald green remains vibrant */
  --primary-foreground:     0 0% 100%;        /* Text on primary */

  /* ───── Secondary / Tints ───── */
  --secondary:              142 30% 15%;      /* Deep muted green for hovers */
  --secondary-foreground:   142 71% 75%;      /* Bright green accent text */

  /* ───── Muted / Input Backgrounds ───── */
  --muted:                  160 15% 15%;      /* Dark muted container */
  --muted-foreground:       160 10% 60%;      /* Mid-tone slate text */

  /* ───── Accent Elements ───── */
  --accent:                 142 30% 15%;      /* Dark accent surface */
  --accent-foreground:      142 71% 75%;      /* Text on accent */

  /* ───── Destructive / Danger ───── */
  --destructive:            0 62% 30%;        /* Muted deep red */
  --destructive-foreground: 0 0% 95%;

  /* ───── Borders, Rings & Radius ───── */
  --border:                 160 15% 20%;      /* Dark hairline border */
  --input:                  160 15% 20%;      /* Form field border */
  --ring:                   142 71% 45%;      /* Focus ring */

  /* ───── Hero Ambient Gradients ───── */
  --gradient-start:         160 25% 12%;      /* Dark hero ambient start */
  --gradient-end:           160 20% 8%;       /* Dark hero ambient end */
}
```

---

## 4. Multi-Theme Switching Catalog

The Sweet Digs Design System supports dynamic runtime retheming by swapping token values. The 5 standard themes are:

| Theme Identifier | Name | `--primary` (HSL) | `--secondary` (HSL) | Surface Character | Best Used For |
|:---|:---|:---|:---|:---|:---|
| `emerald-green` | Emerald Green (Default) | `142 71% 45%` | `142 60% 93%` | Minty, natural, clean | Environmental, sustainability, eco-real estate |
| `pink-red` | Radiant Coral / Pink | `342 90% 55%` | `342 80% 94%` | Warm, energetic, lively | Modern consumer SaaS, creative agencies |
| `obsidian-navy` | Obsidian Navy | `222 80% 55%` | `222 40% 93%` | Deep oceanic, authoritative | Developer platforms, data science portals |
| `sunset-amber` | Sunset Amber | `38 92% 50%` | `38 80% 93%` | Golden warm editorial | Academic research, publishing, editorial |
| `cyber-indigo` | Cyber Indigo | `243 75% 59%` | `243 70% 94%` | High-tech vibrant electric | AI infrastructure, terminal tools, cloud dashboards |

### Theme Variable Swapping Definitions

```css
/* 1. Emerald Green (Default) */
[data-theme="emerald-green"] {
  --primary: 142 71% 45%;
  --primary-foreground: 0 0% 100%;
  --secondary: 142 60% 93%;
  --secondary-foreground: 142 71% 30%;
  --ring: 142 71% 45%;
  --gradient-start: 142 40% 92%;
  --gradient-end: 160 30% 90%;
}

/* 2. Pink / Red */
[data-theme="pink-red"] {
  --primary: 342 90% 55%;
  --primary-foreground: 0 0% 100%;
  --secondary: 342 80% 94%;
  --secondary-foreground: 342 80% 35%;
  --ring: 342 90% 55%;
  --gradient-start: 342 60% 95%;
  --gradient-end: 20 40% 92%;
}

/* 3. Obsidian Navy */
[data-theme="obsidian-navy"] {
  --primary: 222 80% 55%;
  --primary-foreground: 0 0% 100%;
  --secondary: 222 50% 93%;
  --secondary-foreground: 222 80% 30%;
  --ring: 222 80% 55%;
  --gradient-start: 222 45% 93%;
  --gradient-end: 240 30% 91%;
}

/* 4. Sunset Amber */
[data-theme="sunset-amber"] {
  --primary: 38 92% 50%;
  --primary-foreground: 0 0% 100%;
  --secondary: 38 75% 92%;
  --secondary-foreground: 38 92% 28%;
  --ring: 38 92% 50%;
  --gradient-start: 38 60% 93%;
  --gradient-end: 48 40% 91%;
}

/* 5. Cyber Indigo */
[data-theme="cyber-indigo"] {
  --primary: 243 75% 59%;
  --primary-foreground: 0 0% 100%;
  --secondary: 243 65% 94%;
  --secondary-foreground: 243 75% 35%;
  --ring: 243 75% 59%;
  --gradient-start: 243 50% 94%;
  --gradient-end: 270 35% 92%;
}
```

---

## 5. Token Usage Constraints & Rules

1. **Total Ban on Hardcoded Literals:** Components MUST NOT contain hardcoded color hex values (`#22c55e`, `#ff2d6f`). All color declarations MUST reference CSS variables via `hsl(var(--token))`.
2. **Predictable Semantic Tints:** When generating translucent overlays, use the CSS slash syntax: `hsl(var(--primary) / 0.1)`.
3. **Contrast Ratios:** Text on `--primary` MUST always use `--primary-foreground` to maintain WCAG AAA compliance across all theme families.
4. **Positive Boolean State Validation:** When evaluating theme status or dark mode activation in application logic, evaluate conditions implicitly (`if (isDarkModeActive) { ... }`). Explicit checks against true are strictly prohibited.
